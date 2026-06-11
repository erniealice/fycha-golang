// Package action provides HTTP action handlers for the WithholdingCertificate
// CRUD views (add/edit/delete form submissions and GET form renders).
package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	withholdingcertificatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	treasury "github.com/erniealice/fycha-golang/domain/treasury"
	"github.com/erniealice/fycha-golang/domain/treasury/views/withholding_certificate/form"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds the use-case callbacks required by the create/edit/delete handlers.
type Deps struct {
	Routes                       treasury.WithholdingCertificateRoutes
	Labels                       treasury.WithholdingCertificateLabels
	CommonLabels                 any
	CreateWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.CreateWithholdingCertificateRequest) (*withholdingcertificatepb.CreateWithholdingCertificateResponse, error)
	ReadWithholdingCertificate   func(ctx context.Context, req *withholdingcertificatepb.ReadWithholdingCertificateRequest) (*withholdingcertificatepb.ReadWithholdingCertificateResponse, error)
	UpdateWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.UpdateWithholdingCertificateRequest) (*withholdingcertificatepb.UpdateWithholdingCertificateResponse, error)
	DeleteWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.DeleteWithholdingCertificateRequest) (*withholdingcertificatepb.DeleteWithholdingCertificateResponse, error)
}

// NewCreateAction handles GET (form) and POST (submit) for adding a new
// WithholdingCertificate.
//
// GET query params:
//   - revenue_id    — pre-fill the Revenue field
//   - expected_amount — pre-fill from revenue.wht_amount_expected (centavos)
func NewCreateAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("withholding_certificate", "create") {
			return view.HTMXError("You do not have permission to create withholding certificates")
		}

		labels := form.BuildLabels(deps.Labels)

		if viewCtx.Request.Method == http.MethodGet {
			q := viewCtx.Request.URL.Query()
			revenueID := q.Get("revenue_id")
			expectedCents := q.Get("expected_amount")
			expectedDisplay := ""
			if expectedCents != "" {
				if cents, err := strconv.ParseInt(expectedCents, 10, 64); err == nil {
					expectedDisplay = fmt.Sprintf("%.2f", float64(cents)/100.0)
				}
			}
			return view.OK("withholding-certificate-drawer-form", &form.Data{
				FormAction:     deps.Routes.AddURL,
				IsEdit:         false,
				RevenueID:      revenueID,
				ExpectedAmount: expectedDisplay,
				StatusOptions:  form.BuildStatusOptions("WITHHOLDING_CERTIFICATE_STATUS_PENDING_RECEIPT"),
				Labels:         labels,
				CommonLabels:   deps.CommonLabels,
			})
		}

		// POST — create
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError("Invalid form data")
		}
		r := viewCtx.Request

		certifiedStr := r.FormValue("wht_amount_certified")
		var certifiedCents int64
		if certifiedStr != "" {
			f, err := strconv.ParseFloat(certifiedStr, 64)
			if err == nil {
				certifiedCents = int64(f * 100)
			}
		}

		statusVal := r.FormValue("status")
		var status withholdingcertificatepb.WithholdingCertificateStatus
		if sv, ok := withholdingcertificatepb.WithholdingCertificateStatus_value[statusVal]; ok {
			status = withholdingcertificatepb.WithholdingCertificateStatus(sv)
		}

		revenueID := r.FormValue("revenue_id")
		record := &withholdingcertificatepb.WithholdingCertificate{
			RevenueId:         revenueID,
			CertificateNumber: r.FormValue("certificate_number"),
			BuyerTinSnapshot:  optionalString(r.FormValue("payor_tin")),
			IssuedDate:        optionalString(r.FormValue("date_issued")),
			ActualAmount:      certifiedCents,
			Status:            status,
		}
		if notes := r.FormValue("notes"); notes != "" {
			_ = notes // notes field stored in source_citation for now
		}

		if deps.CreateWithholdingCertificate == nil {
			return view.HTMXError("Create use case not available")
		}
		if _, err := deps.CreateWithholdingCertificate(ctx, &withholdingcertificatepb.CreateWithholdingCertificateRequest{
			Data: record,
		}); err != nil {
			log.Printf("CreateWithholdingCertificate error: %v", err)
			return view.HTMXError(err.Error())
		}

		return view.HTMXSuccess("withholding-certs-table")
	})
}

// optionalString converts a string to *string.
// Returns nil when s is empty so proto optional fields remain unset.
func optionalString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
