// Package action holds the charge policy write handlers: policy header, draft
// versions, components, postings, retire, delete and approval.
//
// Every handler re-checks its permission first (layer 2, fail closed) and
// before any use-case call, so a forged POST is refused at the view; the use
// cases remain the authoritative gate (strict verbs, separation of duties).
package action

import (
	"context"
	"log"
	"net/http"
	"strings"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	componentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_component"
	postingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_posting"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	"github.com/erniealice/fycha-golang/domain/ledger/charge_policy/form"
)

// Deps holds the action handler dependencies.
type Deps struct {
	Routes       cp.Routes
	Labels       cp.Labels
	CommonLabels pyeza.CommonLabels
	UseCases     *cp.UseCases

	// AttachmentConfig builds the evidence attachment config (shared with the
	// version page so the table and the handlers agree).
	AttachmentConfig func() *attachment.Config
}

// ---------------------------------------------------------------------------
// Result helpers
// ---------------------------------------------------------------------------

// denied is the fail-closed refusal of a missing permission (HTMXError with the
// shared PermissionDenied message, infra-auth-permission-reflection.md layer 2).
func denied(deps *Deps, code string) view.ViewResult {
	_ = code
	return view.HTMXError(deps.CommonLabels.Errors.PermissionDenied)
}

func refuse(deps *Deps, err error) view.ViewResult {
	kind := cp.ErrorKind(err)
	if kind == cp.ErrKindUnknown {
		log.Printf("charge_policy action: %v", err)
	}
	return refuseKind(deps, kind)
}

// refuseKind maps an error kind to its refusal; a strict-gate denial is the
// shared PermissionDenied message, never the generic error.
func refuseKind(deps *Deps, kind string) view.ViewResult {
	if kind == cp.ErrKindPermissionDenied {
		return denied(deps, "")
	}
	return view.HTMXError(deps.Labels.ErrorMessage(kind))
}

func unavailable(deps *Deps) view.ViewResult {
	return view.ViewResult{
		StatusCode: http.StatusServiceUnavailable,
		Headers:    map[string]string{"HX-Error-Message": deps.Labels.Errors.Unavailable},
	}
}

func ready(deps *Deps, fns ...bool) bool {
	if deps.UseCases == nil {
		return false
	}
	for _, ok := range fns {
		if !ok {
			return false
		}
	}
	return true
}

func versionURL(deps *Deps, policyID, vid, tab string) string {
	u := route.ResolveURL(deps.Routes.VersionDetailURL, "id", policyID, "vid", vid)
	if tab != "" {
		u += "?tab=" + tab
	}
	return u
}

func policyURL(deps *Deps, policyID, tab string) string {
	u := route.ResolveURL(deps.Routes.DetailURL, "id", policyID)
	if tab != "" {
		u += "?tab=" + tab
	}
	return u
}

// ---------------------------------------------------------------------------
// Policy header
// ---------------------------------------------------------------------------

// NewAddAction handles GET (drawer) / POST (create policy + draft v1).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "create") {
			return denied(deps, "charge_policy:create")
		}
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("charge-policy-drawer-form", &form.PolicyData{
				FormAction: deps.Routes.AddURL, Labels: deps.Labels,
			})
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.CreateChargePolicy != nil) {
			return unavailable(deps)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		r := viewCtx.Request
		code := strings.ToUpper(strings.TrimSpace(r.FormValue("code")))
		name := strings.TrimSpace(r.FormValue("name"))
		if code == "" || name == "" {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		pol, draft, err := deps.UseCases.CreatePolicy(ctx, code, name, strings.TrimSpace(r.FormValue("description")))
		if err != nil {
			return refuse(deps, err)
		}
		if draft == nil || pol == nil {
			return view.HTMXError(deps.Labels.Errors.Generic)
		}
		return view.Redirect(versionURL(deps, pol.GetId(), draft.GetId(), ""))
	})
}

// NewEditAction handles GET (drawer) / POST (name + description only).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.ReadChargePolicy != nil && deps.UseCases.UpdateChargePolicy != nil) {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			d, err := deps.UseCases.ReadPolicy(ctx, id)
			if err != nil || d == nil || d.Policy == nil {
				return view.HTMXError(deps.Labels.Errors.NotFound)
			}
			return view.OK("charge-policy-drawer-form", &form.PolicyData{
				FormAction: route.ResolveURL(deps.Routes.EditURL, "id", id), IsEdit: true, ID: id,
				Code: d.Policy.GetCode(), Name: d.Policy.GetName(), Description: d.Policy.GetDescription(),
				Labels: deps.Labels,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		name := strings.TrimSpace(viewCtx.Request.FormValue("name"))
		if name == "" {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		if err := deps.UseCases.UpdatePolicy(ctx, id, name, strings.TrimSpace(viewCtx.Request.FormValue("description"))); err != nil {
			return refuse(deps, err)
		}
		if strings.Contains(viewCtx.Request.Header.Get("HX-Current-URL"), "/detail/") {
			return view.Redirect(policyURL(deps, id, ""))
		}
		return view.HTMXSuccess("charge-policies-table")
	})
}

// NewDeleteAction deletes one never-approved, unused policy (POST).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "delete") {
			return denied(deps, "charge_policy:delete")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.DeleteChargePolicy != nil) {
			return unavailable(deps)
		}
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if err := deleteOne(ctx, deps, id); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess("charge-policies-table")
	})
}

func deleteOne(ctx context.Context, deps *Deps, id string) error {
	if deps.UseCases.GetChargePolicyInUseIds != nil {
		if m, err := deps.UseCases.InUseIDs(ctx, []string{id}); err == nil && m[id] {
			return errInUse
		}
	}
	return deps.UseCases.DeletePolicy(ctx, id)
}

type localErr string

func (e localErr) Error() string { return string(e) }

// errInUse is the view-level in-use refusal (the use case re-checks).
const errInUse = localErr("charge_policy: in_use")

// NewBulkDeleteAction deletes the eligible drafts among the posted ids.
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "delete") {
			return denied(deps, "charge_policy:delete")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.DeleteChargePolicy != nil) {
			return unavailable(deps)
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		var deleted int
		var last error
		for _, id := range viewCtx.Request.Form["id"] {
			if id == "" {
				continue
			}
			if err := deleteOne(ctx, deps, id); err != nil {
				last = err
				continue
			}
			deleted++
		}
		if deleted == 0 && last != nil {
			return refuse(deps, last)
		}
		return view.HTMXSuccess("charge-policies-table")
	})
}

// NewRetireAction retires a policy (POST; confirm dialog on the client).
func NewRetireAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "retire") {
			return denied(deps, "charge_policy:retire")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.RetireChargePolicy != nil) {
			return unavailable(deps)
		}
		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if err := deps.UseCases.RetirePolicy(ctx, id); err != nil {
			return refuse(deps, err)
		}
		if strings.Contains(viewCtx.Request.Header.Get("HX-Current-URL"), "/detail/") {
			return view.Redirect(policyURL(deps, id, ""))
		}
		return view.HTMXSuccess("charge-policies-table")
	})
}

// ---------------------------------------------------------------------------
// Versions
// ---------------------------------------------------------------------------

// NewVersionAddAction handles GET (confirm drawer) / POST (clone → new draft).
func NewVersionAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		policyID := viewCtx.Request.PathValue("id")
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("charge-policy-new-version-drawer-form", &form.NewVersionData{
				FormAction: route.ResolveURL(deps.Routes.VersionAddURL, "id", policyID),
				Labels:     deps.Labels,
				Message:    deps.Labels.Detail.NewVersionMessage,
			})
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.CreateDraftChargePolicyVersion != nil) {
			return unavailable(deps)
		}
		v, err := deps.UseCases.CreateDraftVersion(ctx, policyID)
		if err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, v.GetId(), ""))
	})
}

// enumOf parses a posted proto enum name; ok=false when unknown.
func enumOf(m map[string]int32, v string) (int32, bool) {
	n, ok := m[v]
	return n, ok
}

// loadDraftVersion reads the version, checks it belongs to the policy and is a draft.
func loadDraftVersion(ctx context.Context, deps *Deps, policyID, vid string) (*cp.VersionDetail, string) {
	if deps.UseCases.ReadChargePolicyVersion == nil {
		return nil, cp.ErrKindUnknown
	}
	d, err := deps.UseCases.ReadVersion(ctx, vid)
	if err != nil {
		return nil, cp.ErrorKind(err)
	}
	if d == nil || d.Version == nil || d.Version.GetChargePolicyId() != policyID {
		return nil, cp.ErrKindNotFound
	}
	if !cp.IsDraft(d.Version) {
		return d, cp.ErrKindNotDraft
	}
	return d, ""
}

// NewVersionEditAction handles GET (classification drawer) / POST (update).
func NewVersionEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.UpdateChargePolicyVersion != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		d, kind := loadDraftVersion(ctx, deps, policyID, vid)
		if kind != "" {
			return refuseKind(deps, kind)
		}
		l := deps.Labels
		if viewCtx.Request.Method == http.MethodGet {
			v := d.Version
			return view.OK("charge-policy-version-drawer-form", &form.VersionData{
				FormAction:        route.ResolveURL(deps.Routes.VersionEditURL, "id", policyID, "vid", vid),
				VersionID:         vid,
				AssessmentScope:   v.GetAssessmentScope(),
				AssessmentNote:    v.GetAssessmentNote(),
				AccountingRoles:   l.AccountingRoleOptions(v.GetAccountingRole()),
				TaxPositions:      l.TaxPositionOptions(v.GetTaxPosition()),
				BookPresentations: l.BookPresentationOptions(v.GetBookPresentation()),
				Labels:            l,
			})
		}
		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.HTMXError(l.Errors.FormInvalid)
		}
		r := viewCtx.Request
		upd := &versionpb.ChargePolicyVersion{Id: vid, ChargePolicyId: policyID}
		if s := r.FormValue("accounting_role"); s != "" {
			n, ok := enumOf(enumspb.AccountingRole_value, s)
			if !ok {
				return view.HTMXError(l.Errors.FormInvalid)
			}
			e := enumspb.AccountingRole(n)
			upd.AccountingRole = &e
		}
		if s := r.FormValue("book_presentation"); s != "" {
			n, ok := enumOf(enumspb.BookPresentation_value, s)
			if !ok {
				return view.HTMXError(l.Errors.FormInvalid)
			}
			e := enumspb.BookPresentation(n)
			upd.BookPresentation = &e
		}
		if s := r.FormValue("tax_position"); s != "" {
			n, ok := enumOf(enumspb.TaxPosition_value, s)
			if !ok {
				return view.HTMXError(l.Errors.FormInvalid)
			}
			e := enumspb.TaxPosition(n)
			upd.TaxPosition = &e
		}
		if s := strings.TrimSpace(r.FormValue("tax_treatment_id")); s != "" {
			upd.TaxTreatmentId = &s
		}
		if s := strings.TrimSpace(r.FormValue("assessment_scope")); s != "" {
			upd.AssessmentScope = &s
		}
		if s := strings.TrimSpace(r.FormValue("assessment_note")); s != "" {
			upd.AssessmentNote = &s
		}
		if err := deps.UseCases.UpdateVersion(ctx, upd); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, vid, "classification"))
	})
}

// NewVersionDeleteAction deletes a draft version (POST).
func NewVersionDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "delete") {
			return denied(deps, "charge_policy:delete")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.DeleteChargePolicyVersion != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		if _, kind := loadDraftVersion(ctx, deps, policyID, vid); kind != "" {
			return refuseKind(deps, kind)
		}
		if err := deps.UseCases.DeleteVersion(ctx, vid); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(policyURL(deps, policyID, "versions"))
	})
}

// NewVersionApproveAction approves a draft version (POST; confirm dialog on the
// client). Refusals — permission, self-approval, checklist, unsupported
// combination — come from the use case and are mapped to Lyngua messages.
func NewVersionApproveAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "approve") {
			return denied(deps, "charge_policy:approve")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.ApproveChargePolicyVersion != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		d, kind := loadDraftVersion(ctx, deps, policyID, vid)
		if kind != "" {
			return refuseKind(deps, kind)
		}
		// Mirror of the use-case separation-of-duties rule so the operator gets
		// the reason before any write (the use case re-enforces it strictly).
		if d.Version.GetPreparedBy() != "" && d.Version.GetPreparedBy() == cp.ActorID(ctx) && !perms.Can("charge_policy", "approve_own") {
			return view.HTMXError(deps.Labels.Errors.SelfApproval)
		}
		if _, err := deps.UseCases.Approve(ctx, vid); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, vid, "approval"))
	})
}

// ---------------------------------------------------------------------------
// Components
// ---------------------------------------------------------------------------

func componentForm(deps *Deps, policyID, vid string, c *componentpb.ChargePolicyComponent, action string) *form.ComponentData {
	l := deps.Labels
	data := &form.ComponentData{
		FormAction: action, Labels: l,
		Roles:             l.ComponentRoleOptions(c.GetComponentRole()),
		DocumentKinds:     l.DocumentKindOptions(c.GetDocumentKind()),
		BookPresentations: l.BookPresentationOptions(c.GetBookPresentation()),
	}
	if c.GetId() != "" {
		data.IsEdit, data.ID = true, c.GetId()
	}
	return data
}

// parseComponent reads and validates the posted component fields.
func parseComponent(deps *Deps, r *http.Request, vid string) (*componentpb.ChargePolicyComponent, bool) {
	if err := r.ParseForm(); err != nil {
		return nil, false
	}
	role, ok1 := enumOf(enumspb.ChargeComponentRole_value, r.FormValue("component_role"))
	doc, ok2 := enumOf(enumspb.ChargeDocumentKind_value, r.FormValue("document_kind"))
	book, ok3 := enumOf(enumspb.BookPresentation_value, r.FormValue("book_presentation"))
	if !ok1 || !ok2 || !ok3 || role == 0 || doc == 0 || book == 0 {
		return nil, false
	}
	return &componentpb.ChargePolicyComponent{
		ChargePolicyVersionId: vid,
		ComponentRole:         enumspb.ChargeComponentRole(role),
		DocumentKind:          enumspb.ChargeDocumentKind(doc),
		BookPresentation:      enumspb.BookPresentation(book),
	}, true
}

// NewComponentAddAction handles GET (drawer) / POST (create) for a component.
func NewComponentAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.CreateChargePolicyComponent != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		if _, kind := loadDraftVersion(ctx, deps, policyID, vid); kind != "" {
			return refuseKind(deps, kind)
		}
		action := route.ResolveURL(deps.Routes.ComponentAddURL, "id", policyID, "vid", vid)
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("charge-policy-component-drawer-form", componentForm(deps, policyID, vid, &componentpb.ChargePolicyComponent{}, action))
		}
		c, ok := parseComponent(deps, viewCtx.Request, vid)
		if !ok {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		if err := deps.UseCases.CreateComponent(ctx, c); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, vid, "components"))
	})
}

func findComponent(d *cp.VersionDetail, id string) *componentpb.ChargePolicyComponent {
	for _, c := range d.Components {
		if c.GetId() == id {
			return c
		}
	}
	return nil
}

// NewComponentEditAction handles GET (drawer) / POST (update) for a component.
func NewComponentEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.UpdateChargePolicyComponent != nil) {
			return unavailable(deps)
		}
		policyID, vid, cid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid"), viewCtx.Request.PathValue("cid")
		d, kind := loadDraftVersion(ctx, deps, policyID, vid)
		if kind != "" {
			return refuseKind(deps, kind)
		}
		existing := findComponent(d, cid)
		if existing == nil {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		action := route.ResolveURL(deps.Routes.ComponentEditURL, "id", policyID, "vid", vid, "cid", cid)
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("charge-policy-component-drawer-form", componentForm(deps, policyID, vid, existing, action))
		}
		c, ok := parseComponent(deps, viewCtx.Request, vid)
		if !ok {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		c.Id = cid
		c.SequenceOrder = existing.GetSequenceOrder()
		if err := deps.UseCases.UpdateComponent(ctx, c); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, vid, "components"))
	})
}

// NewComponentDeleteAction deletes a component of a draft version (POST id).
func NewComponentDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.DeleteChargePolicyComponent != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		d, kind := loadDraftVersion(ctx, deps, policyID, vid)
		if kind != "" {
			return refuseKind(deps, kind)
		}
		if id == "" || findComponent(d, id) == nil {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if err := deps.UseCases.DeleteComponent(ctx, id); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess("charge-policy-components-table")
	})
}

// ---------------------------------------------------------------------------
// Postings
// ---------------------------------------------------------------------------

func accountOptions(ctx context.Context, deps *Deps, selected string) []cp.Option {
	if deps.UseCases == nil || deps.UseCases.GetAccountListPageData == nil {
		return nil
	}
	accs, err := deps.UseCases.ListAccounts(ctx)
	if err != nil {
		log.Printf("charge_policy action: ListAccounts: %v", err)
		return nil
	}
	opts := make([]cp.Option, 0, len(accs))
	for _, a := range accs {
		opts = append(opts, cp.Option{Value: a.GetId(), Label: accountLabel(a), Selected: a.GetId() == selected})
	}
	return opts
}

func accountLabel(a *accountpb.Account) string {
	if a.GetCode() != "" {
		return a.GetCode() + " · " + a.GetName()
	}
	return a.GetName()
}

func postingForm(ctx context.Context, deps *Deps, p *postingpb.ChargePolicyPosting, action string) *form.PostingData {
	l := deps.Labels
	data := &form.PostingData{
		FormAction: action, Labels: l,
		Events:   l.EventOptions(p.GetEvent()),
		Roles:    l.PostingRoleOptions(p.GetPostingRole()),
		Accounts: accountOptions(ctx, deps, p.GetAccountId()),
	}
	if p.GetId() != "" {
		data.IsEdit, data.ID = true, p.GetId()
	}
	return data
}

func parsePosting(r *http.Request, vid string) (*postingpb.ChargePolicyPosting, bool) {
	if err := r.ParseForm(); err != nil {
		return nil, false
	}
	ev, ok1 := enumOf(enumspb.ChargePostingEvent_value, r.FormValue("event"))
	ro, ok2 := enumOf(enumspb.ChargePostingRole_value, r.FormValue("posting_role"))
	acc := strings.TrimSpace(r.FormValue("account_id"))
	if !ok1 || !ok2 || ev == 0 || ro == 0 || acc == "" {
		return nil, false
	}
	return &postingpb.ChargePolicyPosting{
		ChargePolicyVersionId: vid,
		Event:                 enumspb.ChargePostingEvent(ev),
		PostingRole:           enumspb.ChargePostingRole(ro),
		AccountId:             acc,
	}, true
}

// NewPostingAddAction handles GET (drawer) / POST (create) for a posting.
func NewPostingAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.CreateChargePolicyPosting != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		if _, kind := loadDraftVersion(ctx, deps, policyID, vid); kind != "" {
			return refuseKind(deps, kind)
		}
		action := route.ResolveURL(deps.Routes.PostingAddURL, "id", policyID, "vid", vid)
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("charge-policy-posting-drawer-form", postingForm(ctx, deps, &postingpb.ChargePolicyPosting{}, action))
		}
		p, ok := parsePosting(viewCtx.Request, vid)
		if !ok {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		if err := deps.UseCases.CreatePosting(ctx, p); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, vid, "postings"))
	})
}

func findPosting(d *cp.VersionDetail, id string) *postingpb.ChargePolicyPosting {
	for _, p := range d.Postings {
		if p.GetId() == id {
			return p
		}
	}
	return nil
}

// NewPostingEditAction handles GET (drawer) / POST (update) for a posting.
func NewPostingEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.UpdateChargePolicyPosting != nil) {
			return unavailable(deps)
		}
		policyID, vid, pid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid"), viewCtx.Request.PathValue("pid")
		d, kind := loadDraftVersion(ctx, deps, policyID, vid)
		if kind != "" {
			return refuseKind(deps, kind)
		}
		existing := findPosting(d, pid)
		if existing == nil {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		action := route.ResolveURL(deps.Routes.PostingEditURL, "id", policyID, "vid", vid, "pid", pid)
		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("charge-policy-posting-drawer-form", postingForm(ctx, deps, existing, action))
		}
		p, ok := parsePosting(viewCtx.Request, vid)
		if !ok {
			return view.HTMXError(deps.Labels.Errors.FormInvalid)
		}
		p.Id = pid
		if err := deps.UseCases.UpdatePosting(ctx, p); err != nil {
			return refuse(deps, err)
		}
		return view.Redirect(versionURL(deps, policyID, vid, "postings"))
	})
}

// NewPostingDeleteAction deletes a posting of a draft version (POST id).
func NewPostingDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if !ready(deps, deps.UseCases != nil && deps.UseCases.DeleteChargePolicyPosting != nil) {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		d, kind := loadDraftVersion(ctx, deps, policyID, vid)
		if kind != "" {
			return refuseKind(deps, kind)
		}
		if id == "" || findPosting(d, id) == nil {
			return view.HTMXError(deps.Labels.Errors.NotFound)
		}
		if err := deps.UseCases.DeletePosting(ctx, id); err != nil {
			return refuse(deps, err)
		}
		return view.HTMXSuccess("charge-policy-postings-table")
	})
}

// ---------------------------------------------------------------------------
// Evidence attachments (draft versions only; update permission)
// ---------------------------------------------------------------------------

func guardedAttachment(deps *Deps, inner func(*attachment.Config) view.View) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "update") {
			return denied(deps, "charge_policy:update")
		}
		if deps.AttachmentConfig == nil || deps.UseCases == nil {
			return unavailable(deps)
		}
		policyID, vid := viewCtx.Request.PathValue("id"), viewCtx.Request.PathValue("vid")
		if _, kind := loadDraftVersion(ctx, deps, policyID, vid); kind != "" {
			return refuseKind(deps, kind)
		}
		return inner(deps.AttachmentConfig()).Handle(ctx, viewCtx)
	})
}

// NewAttachmentUploadAction wraps the hybra upload handler with the charge policy gates.
func NewAttachmentUploadAction(deps *Deps) view.View {
	return guardedAttachment(deps, attachment.NewUploadAction)
}

// NewAttachmentDeleteAction wraps the hybra delete handler with the charge policy gates.
func NewAttachmentDeleteAction(deps *Deps) view.View {
	return guardedAttachment(deps, attachment.NewDeleteAction)
}
