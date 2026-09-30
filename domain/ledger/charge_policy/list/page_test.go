package list

import (
	"net/http"
	"strings"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	"github.com/erniealice/fycha-golang/domain/ledger/charge_policy/internal/cptest"
)

func newDeps(f *cptest.Fake) *Deps {
	return &Deps{
		Routes: cp.DefaultRoutes(), Labels: cp.DefaultLabels(),
		CommonLabels: cptest.Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases(),
	}
}

func seed() *cptest.Fake {
	f := cptest.New()
	f.AddPolicy("p-active", "UTILITY", "Utility recovery", false, cptest.Approved("v1", 1))
	f.AddPolicy("p-both", "WATER", "Water recovery", false, cptest.Approved("w1", 1), cptest.Draft("w2", 2, "u-prep"))
	f.AddPolicy("p-draft", "NEWPOL", "New policy", false, cptest.Draft("d1", 1, "u-prep"))
	f.AddPolicy("p-retired", "OLD", "Old policy", true, cptest.Approved("o1", 1))
	return f
}

func render(t *testing.T, deps *Deps, status string, codes ...string) *PageData {
	t.Helper()
	vc := cptest.Request(http.MethodGet, "/ledger/settings/charge-policies/list/"+status, "", "status", status)
	res := NewView(deps).Handle(cptest.Ctx(codes...), vc)
	if res.Template != "charge-policy-list" {
		t.Fatalf("template = %q (status %d, err %v)", res.Template, res.StatusCode, res.Error)
	}
	return res.Data.(*PageData)
}

func rowIDs(pd *PageData) []string {
	var ids []string
	for _, r := range pd.Table.Rows {
		ids = append(ids, r.ID)
	}
	return ids
}

func colKeys(pd *PageData) string {
	var keys []string
	for _, c := range pd.Table.Columns {
		keys = append(keys, c.Key)
	}
	return strings.Join(keys, ",")
}

// AC-CP-01 (view half): the Draft tab lists a policy that has no approved version.
func TestListTabsBucketPoliciesByComputedStatus(t *testing.T) {
	deps := newDeps(seed())
	codes := []string{"charge_policy:list", "charge_policy:create", "charge_policy:update", "charge_policy:retire", "charge_policy:delete"}

	active := render(t, deps, "active", codes...)
	if got := strings.Join(rowIDs(active), ","); got != "p-active,p-both" {
		t.Errorf("active rows = %s", got)
	}
	draft := render(t, deps, "draft", codes...)
	if got := strings.Join(rowIDs(draft), ","); got != "p-draft" {
		t.Errorf("draft rows = %s", got)
	}
	retired := render(t, deps, "retired", codes...)
	if got := strings.Join(rowIDs(retired), ","); got != "p-retired" {
		t.Errorf("retired rows = %s", got)
	}
	// Tab counts follow the same buckets.
	want := map[string]int{"active": 2, "draft": 1, "retired": 1}
	for _, tab := range active.StatusTabs {
		if tab.Count != want[tab.Key] {
			t.Errorf("tab %s count = %d, want %d", tab.Key, tab.Count, want[tab.Key])
		}
		if !strings.HasSuffix(tab.Href, "/list/"+tab.Key) {
			t.Errorf("tab %s href = %s", tab.Key, tab.Href)
		}
	}
}

func TestListColumnsPerStatus(t *testing.T) {
	deps := newDeps(seed())
	codes := []string{"charge_policy:list"}
	if got := colKeys(render(t, deps, "active", codes...)); got != "code,name,accounting_role,document_kind,current_version,in_use,modified" {
		t.Errorf("active columns = %s", got)
	}
	if got := colKeys(render(t, deps, "draft", codes...)); got != "code,name,draft_version,prepared_by,modified" {
		t.Errorf("draft columns = %s", got)
	}
	if got := colKeys(render(t, deps, "retired", codes...)); got != "code,name,last_version,retired_on,retired_by" {
		t.Errorf("retired columns = %s", got)
	}
}

// AC-PERM-02 (view): the list is refused without charge_policy:list.
func TestListForbiddenWithoutListPermission(t *testing.T) {
	deps := newDeps(seed())
	vc := cptest.Request(http.MethodGet, "/x", "", "status", "active")
	res := NewView(deps).Handle(cptest.Ctx("charge_policy:read"), vc)
	if res.StatusCode != http.StatusForbidden {
		t.Fatalf("status = %d, want 403", res.StatusCode)
	}
	if n := deps.UseCases.GetChargePolicyListPageData; n == nil {
		t.Fatal("ports lost")
	}
}

// AC-PERM-02: buttons are disabled (never hidden) without the verb, with the tooltip.
func TestListPermissionGatedButtons(t *testing.T) {
	deps := newDeps(seed())
	pd := render(t, deps, "active", "charge_policy:list")
	pa := pd.Table.PrimaryAction
	if pa == nil || !pa.Disabled || !strings.Contains(pa.DisabledTooltip, "charge_policy:create") {
		t.Fatalf("add button not disabled with tooltip: %+v", pa)
	}
	if pa.TestID != "charge-policy-add" {
		t.Errorf("add testid = %q", pa.TestID)
	}
	row := pd.Table.Rows[0]
	for _, a := range row.Actions {
		switch a.Action {
		case "view":
			if a.Disabled {
				t.Error("view must stay enabled")
			}
		default:
			if !a.Disabled {
				t.Errorf("action %s enabled without permission", a.Action)
			}
		}
	}

	all := render(t, deps, "active", "charge_policy:list", "charge_policy:create", "charge_policy:update", "charge_policy:retire")
	if all.Table.PrimaryAction.Disabled {
		t.Error("add disabled despite create permission")
	}
}

// AC-CP-05 (view half): a retired policy is not on Active and has only View.
func TestListRetiredRowsAreViewOnly(t *testing.T) {
	deps := newDeps(seed())
	pd := render(t, deps, "retired", "charge_policy:list", "charge_policy:update", "charge_policy:retire")
	if pd.Table.PrimaryAction != nil {
		t.Error("retired tab must not offer Add")
	}
	if got := len(pd.Table.Rows[0].Actions); got != 1 || pd.Table.Rows[0].Actions[0].Action != "view" {
		t.Errorf("retired actions = %+v", pd.Table.Rows[0].Actions)
	}
}

// One open draft at a time: New draft version is disabled when a draft exists.
func TestListNewDraftVersionDisabledWhenDraftOpen(t *testing.T) {
	deps := newDeps(seed())
	pd := render(t, deps, "active", "charge_policy:list", "charge_policy:update", "charge_policy:retire")
	var withDraft, without *types.TableAction
	for _, r := range pd.Table.Rows {
		for i := range r.Actions {
			if r.Actions[i].Action != "clone" {
				continue
			}
			switch r.ID {
			case "p-both":
				withDraft = &r.Actions[i]
			case "p-active":
				without = &r.Actions[i]
			}
		}
	}
	if withDraft == nil || !withDraft.Disabled || withDraft.DisabledTooltip != deps.Labels.Errors.DraftExists {
		t.Errorf("draft-open new-version action = %+v", withDraft)
	}
	if without == nil || without.Disabled {
		t.Errorf("no-draft new-version action = %+v", without)
	}
	for _, r := range pd.Table.Rows {
		for _, a := range r.Actions {
			if strings.Contains(strings.ToLower(a.Label), "approve") {
				t.Errorf("row %s carries an Approve action; approval is version-page only", r.ID)
			}
		}
	}
}

func TestListDeleteDisabledWhenInUse(t *testing.T) {
	f := seed()
	f.InUse["p-draft"] = true
	pd := render(t, newDeps(f), "draft", "charge_policy:list", "charge_policy:delete", "charge_policy:update")
	for _, a := range pd.Table.Rows[0].Actions {
		if a.Action == "delete" && (!a.Disabled || a.DisabledTooltip == "") {
			t.Errorf("delete enabled for an in-use policy: %+v", a)
		}
	}
	if pd.Table.Rows[0].DataAttrs["deletable"] != "false" {
		t.Error("in-use draft must not be bulk-deletable")
	}
}

func TestListFailsClosedWhenUnwired(t *testing.T) {
	deps := &Deps{Routes: cp.DefaultRoutes(), Labels: cp.DefaultLabels(), CommonLabels: pyeza.CommonLabels{}, UseCases: &cp.UseCases{}}
	vc := cptest.Request(http.MethodGet, "/x", "", "status", "active")
	res := NewView(deps).Handle(cptest.Ctx("charge_policy:list"), vc)
	if res.Error == nil || res.Template != "" {
		t.Fatalf("expected an error result, got template %q err %v", res.Template, res.Error)
	}
}
