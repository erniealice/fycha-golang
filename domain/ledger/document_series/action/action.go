// Package action holds the document series write handlers: add, edit, retire.
//
// Every handler re-checks its permission first (layer 2, fail closed) before
// any use-case call, so a forged POST is refused at the view; the use cases
// remain the authoritative gate.
package action

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	ds "github.com/erniealice/fycha-golang/domain/ledger/document_series"
	"github.com/erniealice/fycha-golang/domain/ledger/document_series/form"
)

// Deps holds the action handler dependencies.
type Deps struct {
	Routes       ds.Routes
	Labels       ds.Labels
	CommonLabels pyeza.CommonLabels
	UseCases     *ds.UseCases
}

const defaultPadding = "6"

func denied(deps *Deps) view.ViewResult {
	return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
}

func refuse(deps *Deps, err error) view.ViewResult {
	kind := ds.ErrorKind(err)
	if kind == ds.ErrKindPermissionDenied {
		return denied(deps)
	}
	if kind == ds.ErrKindUnknown {
		log.Printf("document_series action: %v", err)
	}
	return view.HTMXError(deps.Labels.ErrorMessage(kind))
}

func unavailable(deps *Deps) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusServiceUnavailable,
		Headers:    map[string]string{"HX-Error-Message": deps.Labels.Errors.Unavailable},
	}
}

func kindOptions(l ds.Labels, selected string) []form.Option {
	opts := []struct {
		v string
		k enumspb.ChargeDocumentKind
	}{
		{"INVOICE", enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_INVOICE},
		{"RECOVERY_DOCUMENT", enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT},
	}
	out := make([]form.Option, 0, len(opts))
	for _, o := range opts {
		out = append(out, form.Option{Value: o.v, Label: l.DocumentKindLabel(o.k), Selected: o.v == selected})
	}
	return out
}

func parseKind(v string) (enumspb.ChargeDocumentKind, bool) {
	switch strings.ToUpper(strings.TrimSpace(v)) {
	case "INVOICE":
		return enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_INVOICE, true
	case "RECOVERY_DOCUMENT":
		return enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT, true
	}
	return 0, false
}

func opt(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

// padding parses the number length (1–18); ok=false when out of range.
func padding(v string) (int32, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil || n < 1 || n > 18 {
		return 0, false
	}
	return int32(n), true
}

func newData(deps *Deps, action string) *form.SeriesData {
	return &form.SeriesData{
		FormAction: action, Labels: deps.Labels, NumberPadding: defaultPadding, NextNumber: "1",
		FiscalReset: deps.Labels.Enums.FiscalResetNone,
		KindOptions: kindOptions(deps.Labels, ""),
	}
}

// NewAddAction handles GET (drawer) / POST (create an ACTIVE series).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("document_series", "create") {
			return denied(deps)
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("document-series-drawer-form", newData(deps, deps.Routes.AddURL))
		}
		if deps.UseCases == nil || deps.UseCases.CreateDocumentSeries == nil {
			return unavailable(deps)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		r := viewCtx.Request
		code := strings.ToUpper(strings.TrimSpace(r.FormValue("code")))
		issuer := strings.TrimSpace(r.FormValue("issuer_name"))
		kind, kindOK := parseKind(r.FormValue("document_kind"))
		pad, padOK := padding(r.FormValue("number_padding"))
		if code == "" || issuer == "" || !kindOK || !padOK {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		var next int64 = 1
		if v := strings.TrimSpace(r.FormValue("next_number")); v != "" {
			n, err := strconv.ParseInt(v, 10, 64)
			if err != nil || n < 1 {
				return view.HTMXError(deps.Labels.Errors.FormInvalid)
			}
			next = n
		}
		_, err := deps.UseCases.Create(ctx, &pb.DocumentSeries{
			Code: code, Name: opt(r.FormValue("name")), IssuerName: issuer, IssuerTaxId: opt(r.FormValue("issuer_tax_id")),
			DocumentKind: kind, Prefix: opt(r.FormValue("prefix")), BranchCode: opt(r.FormValue("branch_code")),
			// Only NONE is implemented (a YEARLY reset is not yet available).
			FiscalReset:   pb.DocumentSeriesFiscalReset_DOCUMENT_SERIES_FISCAL_RESET_NONE,
			NextNumber:    next,
			NumberPadding: pad,
		})
		if err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess("document-series-table")
	})
}

// NewEditAction handles GET (drawer) / POST (mutable fields only).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("document_series", "update") {
			return denied(deps)
		}
		if deps.UseCases == nil || deps.UseCases.ReadDocumentSeries == nil || deps.UseCases.UpdateDocumentSeries == nil {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			s, err := deps.UseCases.ReadSeries(ctx, id)
			if err != nil || s == nil {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			if ds.IsRetired(s) {
				return view.HTMXError(deps.Labels.Errors.Retired)
			}
			return view.OK("document-series-drawer-form", &form.SeriesData{
				FormAction: route.ResolveURL(deps.Routes.EditURL, "id", id), IsEdit: true, ID: id,
				Code: s.GetCode(), Name: s.GetName(), IssuerName: s.GetIssuerName(), IssuerTaxID: s.GetIssuerTaxId(),
				KindLabel: deps.Labels.DocumentKindLabel(s.GetDocumentKind()),
				Prefix:    s.GetPrefix(), BranchCode: s.GetBranchCode(),
				FiscalReset:   deps.Labels.FiscalResetLabel(s.GetFiscalReset()),
				NextNumber:    strconv.FormatInt(s.GetNextNumber(), 10),
				NumberPadding: strconv.Itoa(int(s.GetNumberPadding())),
				NextPreview:   ds.Format(deps.Labels.Detail.NextPreview, ds.NextDocumentNumber(s)),
				Locked:        s.GetNextNumber() > 1,
				Labels:        deps.Labels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		r := viewCtx.Request
		issuer := strings.TrimSpace(r.FormValue("issuer_name"))
		cur, rerr := deps.UseCases.ReadSeries(ctx, id)
		if rerr != nil {
			return refuse(deps, rerr)
		}
		if cur == nil {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		prefix, pad, padOK := opt(r.FormValue("prefix")), int32(0), false
		if cur.GetNextNumber() > 1 {
			// Numbers were already issued: prefix and padding are locked (the
			// drawer renders them read-only); keep the stored values so issued
			// and new numbers never differ in format. The use case stays the
			// authority.
			prefix, pad, padOK = cur.Prefix, cur.GetNumberPadding(), true
		} else {
			pad, padOK = padding(r.FormValue("number_padding"))
		}
		if issuer == "" || !padOK {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		err := deps.UseCases.Update(ctx, &pb.DocumentSeries{
			Id: id, Name: opt(r.FormValue("name")), IssuerName: issuer, IssuerTaxId: opt(r.FormValue("issuer_tax_id")),
			Prefix: prefix, BranchCode: opt(r.FormValue("branch_code")), NumberPadding: pad,
		})
		if err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess("document-series-table")
	})
}

// NewRetireAction retires a series (POST; confirm dialog on the client).
func NewRetireAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("document_series", "update") {
			return denied(deps)
		}
		if deps.UseCases == nil || deps.UseCases.ReadDocumentSeries == nil || deps.UseCases.UpdateDocumentSeries == nil {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if err := deps.UseCases.Retire(ctx, id); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess("document-series-table")
	})
}
