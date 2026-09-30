package version

import (
	"context"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	"net/http"
	"strings"
	"testing"

	"github.com/erniealice/espyna-golang/shared/identity"
	componentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_component"
	"github.com/erniealice/pyeza-golang/types"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	"github.com/erniealice/fycha-golang/domain/ledger/charge_policy/internal/cptest"
)

func newDeps(f *cptest.Fake) *Deps {
	return &Deps{Routes: cp.DefaultRoutes(), Labels: cp.DefaultLabels(), CommonLabels: cptest.Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases()}
}

func seed() *cptest.Fake {
	f := cptest.New()
	f.AddPolicy("p1", "UTILITY", "Utility recovery", false, cptest.Approved("v-appr", 1), cptest.Draft("v-draft", 2, "u-prep"))
	f.Components["v-draft"] = []*componentpb.ChargePolicyComponent{{Id: "c1", ChargePolicyVersionId: "v-draft"}}
	f.Validation = &versionpb.ChargePolicyApprovalChecklist{Supported: true, Passed: true, Items: []*versionpb.ChargePolicyChecklistItem{{Code: "assessment_scope", Passed: true}}}
	return f
}

func render(t *testing.T, f *cptest.Fake, vid, tab string, codes ...string) *PageData {
	t.Helper()
	vc := cptest.Request(http.MethodGet, "/x?tab="+tab, "", "id", "p1", "vid", vid)
	res := NewView(newDeps(f)).Handle(cptest.Ctx(codes...), vc)
	if res.Template != "charge-policy-version-detail" {
		t.Fatalf("template %q status %d err %v", res.Template, res.StatusCode, res.Error)
	}
	return res.Data.(*PageData)
}

// AC-CP-04 (view half): the Approve button's disabled states and reasons.
func TestApproveButtonStates(t *testing.T) {
	l := cp.DefaultLabels()
	base := []string{"charge_policy:read", "charge_policy:update"}
	withApprove := append(append([]string{}, base...), "charge_policy:approve")

	cases := []struct {
		name    string
		mutate  func(*cptest.Fake)
		vid     string
		perms   []string
		off     bool
		tipPart string
	}{
		{"missing approve permission", nil, "v-draft", base, true, "charge_policy:approve"},
		{"ready", nil, "v-draft", withApprove, false, ""},
		{"checklist failed", func(f *cptest.Fake) {
			f.Validation = &versionpb.ChargePolicyApprovalChecklist{Supported: true, Passed: false}
		}, "v-draft", withApprove, true, l.Errors.ChecklistFailed},
		{"unsupported combination", func(f *cptest.Fake) {
			f.Validation = &versionpb.ChargePolicyApprovalChecklist{Supported: false, Reasons: []string{"own_supply"}}
		}, "v-draft", withApprove, true, l.Errors.UnsupportedCombination},
		{"already approved", nil, "v-appr", withApprove, true, l.Errors.NotDraft},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			f := seed()
			if c.mutate != nil {
				c.mutate(f)
			}
			pd := render(t, f, c.vid, "approval", c.perms...)
			if pd.ApproveOff != c.off {
				t.Fatalf("ApproveOff = %v, want %v (tip %q)", pd.ApproveOff, c.off, pd.ApproveTip)
			}
			if c.off && !strings.Contains(pd.ApproveTip, c.tipPart) {
				t.Errorf("tip %q lacks %q", pd.ApproveTip, c.tipPart)
			}
		})
	}
}

// Self-approval is disabled unless charge_policy:approve_own is held.
func TestApproveDisabledForOwnDraftWithoutApproveOwn(t *testing.T) {
	f := seed()
	deps := newDeps(f)
	ctx := identity.WithRequestIdentity(context.Background(), &identity.RequestIdentity{UserID: "u-prep"})
	own := cptest.Draft("v-own", 3, "u-prep")
	ok := &versionpb.ChargePolicyApprovalChecklist{Supported: true, Passed: true}

	off, tip := approveState(ctx, deps, types.NewUserPermissions([]string{"charge_policy:approve"}), own, ok, true, false)
	if !off || tip != cp.DefaultLabels().Errors.SelfApproval {
		t.Errorf("own draft without approve_own: off=%v tip=%q", off, tip)
	}
	off, _ = approveState(ctx, deps, types.NewUserPermissions([]string{"charge_policy:approve", "charge_policy:approve_own"}), own, ok, true, false)
	if off {
		t.Error("approve_own holder must be able to approve their own draft")
	}
	// Someone else's draft needs only approve.
	other := cptest.Draft("v-other", 4, "u-else")
	off, _ = approveState(ctx, deps, types.NewUserPermissions([]string{"charge_policy:approve"}), other, ok, true, false)
	if off {
		t.Error("another preparer's draft blocked by the self-approval rule")
	}
}

// The locked notice shows for approved versions and the checklist card for drafts.
func TestLockedNoticeAndChecklist(t *testing.T) {
	f := seed()
	perms := []string{"charge_policy:read"}
	appr := render(t, f, "v-appr", "classification", perms...)
	if !appr.ShowLocked || appr.IsDraft {
		t.Error("approved version must be locked")
	}
	if appr.CanEdit || appr.EditTip != cp.DefaultLabels().Errors.NotDraft {
		t.Errorf("approved version editable: %v %q", appr.CanEdit, appr.EditTip)
	}
	draft := render(t, f, "v-draft", "classification", perms...)
	if !draft.IsDraft || !draft.ChecklistReady || len(draft.Checklist) != 1 {
		t.Errorf("draft checklist missing: %+v", draft.Checklist)
	}
}

// A missing required posting pair shows as a red "Missing" row on a draft.
func TestPostingsTableShowsMissingRows(t *testing.T) {
	f := seed()
	pd := render(t, f, "v-draft", "postings", "charge_policy:read", "charge_policy:update")
	missing := 0
	for _, r := range pd.PostingsTable.Rows {
		if strings.HasPrefix(r.ID, "missing-") {
			missing++
			if last := r.Cells[len(r.Cells)-1]; last.Variant != "danger" || last.Value != cp.DefaultLabels().Detail.MissingMapping {
				t.Errorf("missing cell = %+v", last)
			}
		}
	}
	if missing != 4 {
		t.Errorf("missing rows = %d, want 4 (the S1 required pairs)", missing)
	}
	if !strings.Contains(pd.PostingsTable.RefreshURL, "table=1") {
		t.Errorf("postings refresh url = %q", pd.PostingsTable.RefreshURL)
	}
	// Approved versions never show missing rows.
	ap := render(t, f, "v-appr", "postings", "charge_policy:read")
	for _, r := range ap.PostingsTable.Rows {
		if strings.HasPrefix(r.ID, "missing-") {
			t.Error("approved version shows a missing row")
		}
	}
}

func TestComponentActionsDisabledWithoutUpdate(t *testing.T) {
	f := seed()
	pd := render(t, f, "v-draft", "components", "charge_policy:read")
	if !pd.ComponentsTable.PrimaryAction.Disabled {
		t.Error("add component enabled without update")
	}
	for _, a := range pd.ComponentsTable.Rows[0].Actions {
		if !a.Disabled {
			t.Errorf("component action %s enabled without update", a.Action)
		}
	}
	on := render(t, f, "v-draft", "components", "charge_policy:read", "charge_policy:update")
	if on.ComponentsTable.PrimaryAction.Disabled {
		t.Error("add component disabled with update")
	}
}

func TestVersionPageGates(t *testing.T) {
	f := seed()
	vc := cptest.Request(http.MethodGet, "/x", "", "id", "p1", "vid", "v-draft")
	if res := NewView(newDeps(f)).Handle(cptest.Ctx("charge_policy:list"), vc); res.StatusCode != http.StatusForbidden {
		t.Errorf("no read permission: %d", res.StatusCode)
	}
	// A version addressed under another policy is not found.
	vc = cptest.Request(http.MethodGet, "/x", "", "id", "other", "vid", "v-draft")
	if res := NewView(newDeps(f)).Handle(cptest.Ctx("charge_policy:read"), vc); res.StatusCode != http.StatusNotFound {
		t.Errorf("cross-policy version: %d", res.StatusCode)
	}
}

func TestChecklistTextMapsCodes(t *testing.T) {
	l := cp.DefaultLabels()
	if got := checklistText(l, "assessment_scope"); got != l.Detail.ChecklistCodeLabel["assessment_scope"] {
		t.Errorf("assessment_scope -> %q", got)
	}
	if got := checklistText(l, "posting_account:ISSUE:RECEIVABLE"); !strings.Contains(got, l.Enums.EventIssue) || !strings.Contains(got, l.Enums.PostingRoleReceivable) {
		t.Errorf("posting_account text = %q", got)
	}
}
