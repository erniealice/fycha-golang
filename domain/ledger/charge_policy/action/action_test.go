package action

import (
	"net/http"
	"net/url"
	"testing"

	componentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_component"
	postingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_posting"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	"github.com/erniealice/pyeza-golang/view"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	"github.com/erniealice/fycha-golang/domain/ledger/charge_policy/internal/cptest"
)

func newDeps(f *cptest.Fake) *Deps {
	return &Deps{Routes: cp.DefaultRoutes(), Labels: cp.DefaultLabels(), CommonLabels: cptest.Common(), UseCases: f.UseCases()}
}

func seed() *cptest.Fake {
	f := cptest.New()
	f.AddPolicy("p1", "UTILITY", "Utility recovery", false, cptest.Approved("v-appr", 1), cptest.Draft("v-draft", 2, "u-actor"))
	f.Components["v-draft"] = []*componentpb.ChargePolicyComponent{{Id: "c1", ChargePolicyVersionId: "v-draft"}}
	f.Postings["v-draft"] = []*postingpb.ChargePolicyPosting{{Id: "x1", ChargePolicyVersionId: "v-draft"}}
	f.Components["v-appr"] = []*componentpb.ChargePolicyComponent{{Id: "c-appr", ChargePolicyVersionId: "v-appr"}}
	return f
}

type call struct {
	name   string
	build  func(*Deps) view.View
	method string
	path   string
	form   string
	pv     []string
}

func calls() []call {
	return []call{
		{"add", NewAddAction, http.MethodPost, "/action/charge-policy/add", "code=X&name=Y", nil},
		{"edit", NewEditAction, http.MethodPost, "/e", "name=Y", []string{"id", "p1"}},
		{"delete", NewDeleteAction, http.MethodPost, "/d?id=p1", "", nil},
		{"bulk-delete", NewBulkDeleteAction, http.MethodPost, "/d", "id=p1", nil},
		{"retire", NewRetireAction, http.MethodPost, "/r", "", []string{"id", "p1"}},
		{"version-add", NewVersionAddAction, http.MethodPost, "/va", "", []string{"id", "p1"}},
		{"version-edit", NewVersionEditAction, http.MethodPost, "/ve", "assessment_scope=x", []string{"id", "p1", "vid", "v-draft"}},
		{"version-delete", NewVersionDeleteAction, http.MethodPost, "/vd", "", []string{"id", "p1", "vid", "v-draft"}},
		{"version-approve", NewVersionApproveAction, http.MethodPost, "/vap", "", []string{"id", "p1", "vid", "v-draft"}},
		{"component-add", NewComponentAddAction, http.MethodPost, "/ca", "component_role=CHARGE_COMPONENT_ROLE_RECOVERY_COST&document_kind=CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT&book_presentation=BOOK_PRESENTATION_EXCLUDED_REIMBURSEMENT", []string{"id", "p1", "vid", "v-draft"}},
		{"component-edit", NewComponentEditAction, http.MethodPost, "/ce", "component_role=CHARGE_COMPONENT_ROLE_FEE&document_kind=CHARGE_DOCUMENT_KIND_INVOICE&book_presentation=BOOK_PRESENTATION_REVENUE", []string{"id", "p1", "vid", "v-draft", "cid", "c1"}},
		{"component-delete", NewComponentDeleteAction, http.MethodPost, "/cd?id=c1", "", []string{"id", "p1", "vid", "v-draft"}},
		{"posting-add", NewPostingAddAction, http.MethodPost, "/pa", "event=CHARGE_POSTING_EVENT_ISSUE&posting_role=CHARGE_POSTING_ROLE_RECEIVABLE&account_id=a1", []string{"id", "p1", "vid", "v-draft"}},
		{"posting-edit", NewPostingEditAction, http.MethodPost, "/pe", "event=CHARGE_POSTING_EVENT_ISSUE&posting_role=CHARGE_POSTING_ROLE_CLEARING&account_id=a1", []string{"id", "p1", "vid", "v-draft", "pid", "x1"}},
		{"posting-delete", NewPostingDeleteAction, http.MethodPost, "/pd?id=x1", "", []string{"id", "p1", "vid", "v-draft"}},
	}
}

// AC-PERM-05 (view half): every write handler refuses a forged POST that lacks
// its permission — 403 and no use-case mutation.
func TestForgedPostsWithoutPermissionAreRefused(t *testing.T) {
	for _, c := range calls() {
		t.Run(c.name, func(t *testing.T) {
			f := seed()
			deps := newDeps(f)
			// Holds read/list only: every write verb is missing.
			ctx := cptest.Ctx("charge_policy:read", "charge_policy:list")
			res := c.build(deps).Handle(ctx, cptest.Request(c.method, c.path, c.form, c.pv...))
			if res.StatusCode != http.StatusUnprocessableEntity {
				t.Fatalf("status = %d, want 422 (HTMXError)", res.StatusCode)
			}
			if res.Headers["HX-Error-Message"] != cptest.Common().Errors.PermissionDenied {
				t.Errorf("refusal message = %q, want the shared PermissionDenied label", res.Headers["HX-Error-Message"])
			}
			if n := f.MutationCalls(); n != 0 {
				t.Fatalf("use case mutated %d time(s) despite the refusal", n)
			}
		})
	}
}

func post(t *testing.T, v view.View, perms []string, target, form string, pv ...string) view.ViewResult {
	t.Helper()
	return v.Handle(cptest.Ctx(perms...), cptest.Request(http.MethodPost, target, form, pv...))
}

// AC-CP-01 (view half): create redirects to the new draft's version page.
func TestAddCreatesPolicyAndRedirectsToDraft(t *testing.T) {
	f := seed()
	res := post(t, NewAddAction(newDeps(f)), []string{"charge_policy:create"}, "/a", url.Values{"code": {"util_x"}, "name": {"Util"}}.Encode())
	if res.Redirect != "/ledger/settings/charge-policies/detail/p-new/version/v-new" {
		t.Fatalf("redirect = %q (status %d)", res.Redirect, res.StatusCode)
	}
	if f.Calls["CreatePolicy"] != 1 {
		t.Fatal("CreatePolicy not called once")
	}
	// Blank name/code is refused before the use case.
	f2 := seed()
	res = post(t, NewAddAction(newDeps(f2)), []string{"charge_policy:create"}, "/a", "code=&name=")
	if res.StatusCode != http.StatusUnprocessableEntity || f2.MutationCalls() != 0 {
		t.Errorf("blank create: status %d, mutations %d", res.StatusCode, f2.MutationCalls())
	}
}

// AC-CP-04 (view half): approve refusals map to their Lyngua messages.
func TestApproveRefusalsAreMapped(t *testing.T) {
	perms := []string{"charge_policy:approve"}
	l := cp.DefaultLabels()

	// Use-case refusal (checklist) -> message from Lyngua errors.checklist_failed.
	f := seed()
	f.ApproveErr = cptest.ErrChecklist
	res := post(t, NewVersionApproveAction(newDeps(f)), perms, "/a", "", "id", "p1", "vid", "v-draft")
	if res.StatusCode != http.StatusUnprocessableEntity || res.Headers["HX-Error-Message"] != l.Errors.ChecklistFailed {
		t.Errorf("checklist refusal = %d %q", res.StatusCode, res.Headers["HX-Error-Message"])
	}

	// Self-approval mirrored at the view: prepared_by == actor without approve_own.
	// (No identity in the test context => actor "" never equals a preparer, so the
	// use-case mapping is what proves the message.)
	f = seed()
	f.ApproveErr = cptest.ErrSelfApproval
	res = post(t, NewVersionApproveAction(newDeps(f)), perms, "/a", "", "id", "p1", "vid", "v-draft")
	if res.Headers["HX-Error-Message"] != l.Errors.SelfApproval {
		t.Errorf("self-approval message = %q", res.Headers["HX-Error-Message"])
	}

	// Already-approved version: refused as not_draft, use case never called.
	f = seed()
	res = post(t, NewVersionApproveAction(newDeps(f)), perms, "/a", "", "id", "p1", "vid", "v-appr")
	if res.Headers["HX-Error-Message"] != l.Errors.NotDraft || f.Calls["Approve"] != 0 {
		t.Errorf("approved-version approve: %q calls=%d", res.Headers["HX-Error-Message"], f.Calls["Approve"])
	}

	// Success redirects back to the approval tab.
	f = seed()
	res = post(t, NewVersionApproveAction(newDeps(f)), perms, "/a", "", "id", "p1", "vid", "v-draft")
	if res.Redirect != "/ledger/settings/charge-policies/detail/p1/version/v-draft?tab=approval" || f.Calls["Approve"] != 1 {
		t.Errorf("approve success: redirect %q calls %d", res.Redirect, f.Calls["Approve"])
	}
}

// AC-CP-02 (view half): approved versions refuse child edits at the view too.
func TestChildWritesRefusedOnApprovedVersion(t *testing.T) {
	perms := []string{"charge_policy:update"}
	l := cp.DefaultLabels()
	f := seed()
	deps := newDeps(f)
	form := "component_role=CHARGE_COMPONENT_ROLE_FEE&document_kind=CHARGE_DOCUMENT_KIND_INVOICE&book_presentation=BOOK_PRESENTATION_REVENUE"
	for name, res := range map[string]view.ViewResult{
		"component-add":    post(t, NewComponentAddAction(deps), perms, "/x", form, "id", "p1", "vid", "v-appr"),
		"component-edit":   post(t, NewComponentEditAction(deps), perms, "/x", form, "id", "p1", "vid", "v-appr", "cid", "c-appr"),
		"component-delete": post(t, NewComponentDeleteAction(deps), perms, "/x?id=c-appr", "", "id", "p1", "vid", "v-appr"),
		"version-edit":     post(t, NewVersionEditAction(deps), perms, "/x", "assessment_scope=x", "id", "p1", "vid", "v-appr"),
	} {
		if res.Headers["HX-Error-Message"] != l.Errors.NotDraft {
			t.Errorf("%s: message %q, want not_draft", name, res.Headers["HX-Error-Message"])
		}
	}
	if f.MutationCalls() != 0 {
		t.Errorf("approved version was mutated %d time(s)", f.MutationCalls())
	}
}

// IDOR: a child id that does not belong to the addressed version is not found.
func TestChildIDMustBelongToVersion(t *testing.T) {
	perms := []string{"charge_policy:update"}
	f := seed()
	deps := newDeps(f)
	res := post(t, NewComponentDeleteAction(deps), perms, "/x?id=c-appr", "", "id", "p1", "vid", "v-draft")
	if res.Headers["HX-Error-Message"] != deps.Labels.Errors.NotFound || f.Calls["DeleteComponent"] != 0 {
		t.Errorf("foreign component delete: %q calls %d", res.Headers["HX-Error-Message"], f.Calls["DeleteComponent"])
	}
	// A version addressed under the wrong policy is not found either.
	res = post(t, NewVersionDeleteAction(deps), []string{"charge_policy:delete"}, "/x", "", "id", "other", "vid", "v-draft")
	if res.Headers["HX-Error-Message"] != deps.Labels.Errors.NotFound || f.Calls["DeleteVersion"] != 0 {
		t.Errorf("cross-policy version delete: %q calls %d", res.Headers["HX-Error-Message"], f.Calls["DeleteVersion"])
	}
}

func TestVersionEditParsesEnumsAndRejectsUnknownValues(t *testing.T) {
	perms := []string{"charge_policy:update"}
	f := seed()
	deps := newDeps(f)
	form := url.Values{
		"accounting_role":  {"ACCOUNTING_ROLE_AGENT"},
		"tax_position":     {"TAX_POSITION_EXCLUDED_REIMBURSEMENT"},
		"assessment_scope": {" Utilities recharge "},
	}.Encode()
	res := post(t, NewVersionEditAction(deps), perms, "/x", form, "id", "p1", "vid", "v-draft")
	if res.Redirect == "" || f.LastVersion == nil {
		t.Fatalf("edit failed: %+v", res)
	}
	if f.LastVersion.GetAccountingRole() != enumspb.AccountingRole_ACCOUNTING_ROLE_AGENT ||
		f.LastVersion.GetAssessmentScope() != "Utilities recharge" {
		t.Errorf("parsed version = %+v", f.LastVersion)
	}
	f2 := seed()
	res = post(t, NewVersionEditAction(newDeps(f2)), perms, "/x", "accounting_role=NOPE", "id", "p1", "vid", "v-draft")
	if res.StatusCode != http.StatusUnprocessableEntity || f2.Calls["UpdateVersion"] != 0 {
		t.Errorf("unknown enum accepted: %d", res.StatusCode)
	}
}

func TestRetireRefusedWithoutRetirePermissionButAllowedWithIt(t *testing.T) {
	f := seed()
	deps := newDeps(f)
	if res := post(t, NewRetireAction(deps), []string{"charge_policy:update"}, "/x", "", "id", "p1"); res.Headers["HX-Error-Message"] != cptest.Common().Errors.PermissionDenied || f.Calls["RetirePolicy"] != 0 {
		t.Errorf("retire without verb: %d %q", res.StatusCode, res.Headers["HX-Error-Message"])
	}
	res := post(t, NewRetireAction(deps), []string{"charge_policy:retire"}, "/x", "", "id", "p1")
	if res.StatusCode != http.StatusOK || f.Calls["RetirePolicy"] != 1 {
		t.Errorf("retire: %d calls %d", res.StatusCode, f.Calls["RetirePolicy"])
	}
}

func TestUnwiredUseCasesFailClosed(t *testing.T) {
	deps := &Deps{Routes: cp.DefaultRoutes(), Labels: cp.DefaultLabels(), CommonLabels: cptest.Common(), UseCases: &cp.UseCases{}}
	res := post(t, NewVersionApproveAction(deps), []string{"charge_policy:approve"}, "/x", "", "id", "p1", "vid", "v-draft")
	if res.StatusCode != http.StatusServiceUnavailable {
		t.Errorf("unwired approve: %d", res.StatusCode)
	}
}

// A strict-gate denial after the view's permission check maps to the shared
// PermissionDenied message, not the generic error (R5 M1, C20) — on the direct
// refuse path and on the shared draft-version loader path.
func TestStrictGateDenialMapsToPermissionDenied(t *testing.T) {
	want := cptest.Common().Errors.PermissionDenied
	f := seed()
	f.ApproveErr = cptest.ErrPermissionDenied
	res := post(t, NewVersionApproveAction(newDeps(f)), []string{"charge_policy:approve"}, "/a", "", "id", "p1", "vid", "v-draft")
	if got := res.Headers["HX-Error-Message"]; got != want || want == "" {
		t.Errorf("approve denial -> %q, want %q", got, want)
	}
	if got := refuseKind(newDeps(seed()), cp.ErrKindPermissionDenied).Headers["HX-Error-Message"]; got != want {
		t.Errorf("refuseKind denial -> %q, want %q", got, want)
	}
}
