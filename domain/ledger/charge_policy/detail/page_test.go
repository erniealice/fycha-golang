package detail

import (
	"net/http"
	"testing"

	"github.com/erniealice/pyeza-golang/types"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	"github.com/erniealice/fycha-golang/domain/ledger/charge_policy/internal/cptest"
)

func newDeps(f *cptest.Fake) *Deps {
	return &Deps{Routes: cp.DefaultRoutes(), Labels: cp.DefaultLabels(), CommonLabels: cptest.Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases()}
}

func render(t *testing.T, f *cptest.Fake, id, tab string, codes ...string) *PageData {
	t.Helper()
	vc := cptest.Request(http.MethodGet, "/x?tab="+tab, "", "id", id)
	res := NewView(newDeps(f)).Handle(cptest.Ctx(codes...), vc)
	if res.Template != "charge-policy-detail" {
		t.Fatalf("template %q status %d err %v", res.Template, res.StatusCode, res.Error)
	}
	return res.Data.(*PageData)
}

func TestDetailTabsAndBuckets(t *testing.T) {
	f := cptest.New()
	f.AddPolicy("p1", "UTILITY", "Utility recovery", false, cptest.Approved("v1", 1), cptest.Draft("v2", 2, "u"))
	pd := render(t, f, "p1", "versions", "charge_policy:read")
	if len(pd.TabItems) != 4 || pd.TabItems[1].Count != 2 {
		t.Fatalf("tabs = %+v", pd.TabItems)
	}
	if pd.VersionsTable == nil || len(pd.VersionsTable.Rows) != 2 || pd.VersionsTable.Rows[0].ID != "v2" {
		t.Fatalf("versions table not newest-first: %+v", pd.VersionsTable)
	}
	// Unknown tab falls back to info.
	if got := render(t, f, "p1", "bogus", "charge_policy:read").ActiveTab; got != "info" {
		t.Errorf("fallback tab = %s", got)
	}
}

// One open draft at a time and permission gating on the header buttons.
func TestDetailButtonStates(t *testing.T) {
	f := cptest.New()
	f.AddPolicy("p1", "UTILITY", "Utility recovery", false, cptest.Approved("v1", 1), cptest.Draft("v2", 2, "u"))
	pd := render(t, f, "p1", "info", "charge_policy:read")
	if pd.CanEdit || pd.CanRetire || !pd.NewVersionOff {
		t.Errorf("read-only viewer: edit %v retire %v newVersionOff %v", pd.CanEdit, pd.CanRetire, pd.NewVersionOff)
	}
	if pd.RetireDisabled == "" || pd.NewVersionTip != cp.DefaultLabels().Errors.DraftExists {
		t.Errorf("tooltips: retire %q new %q", pd.RetireDisabled, pd.NewVersionTip)
	}
	full := render(t, f, "p1", "info", "charge_policy:read", "charge_policy:update", "charge_policy:retire")
	if !full.CanEdit || !full.CanRetire {
		t.Error("permitted viewer lost edit/retire")
	}
}

// AC-CP-05 (view half): a retired policy cannot be retired or versioned again.
func TestRetiredPolicyDetail(t *testing.T) {
	f := cptest.New()
	f.AddPolicy("p1", "OLD", "Old", true, cptest.Approved("v1", 1))
	pd := render(t, f, "p1", "info", "charge_policy:read", "charge_policy:update", "charge_policy:retire")
	if !pd.IsRetired || pd.CanRetire || !pd.NewVersionOff {
		t.Errorf("retired detail: retired %v canRetire %v newOff %v", pd.IsRetired, pd.CanRetire, pd.NewVersionOff)
	}
}

func TestDetailGatesAndNotFound(t *testing.T) {
	f := cptest.New()
	f.AddPolicy("p1", "X", "X", false, cptest.Draft("v1", 1, "u"))
	vc := cptest.Request(http.MethodGet, "/x", "", "id", "p1")
	if res := NewView(newDeps(f)).Handle(cptest.Ctx("charge_policy:list"), vc); res.StatusCode != http.StatusForbidden {
		t.Errorf("no read: %d", res.StatusCode)
	}
	vc = cptest.Request(http.MethodGet, "/x", "", "id", "missing")
	if res := NewView(newDeps(f)).Handle(cptest.Ctx("charge_policy:read"), vc); res.StatusCode != http.StatusNotFound {
		t.Errorf("unknown id: %d", res.StatusCode)
	}
}
