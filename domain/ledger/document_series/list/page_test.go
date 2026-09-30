package list

import (
	"net/http"
	"strings"
	"testing"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"

	ds "github.com/erniealice/fycha-golang/domain/ledger/document_series"
	"github.com/erniealice/fycha-golang/domain/ledger/document_series/internal/dstest"
)

func newDeps(f *dstest.Fake) *Deps {
	return &Deps{Routes: ds.DefaultRoutes(), Labels: ds.DefaultLabels(), CommonLabels: dstest.Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases()}
}

func render(t *testing.T, deps *Deps, status string, codes ...string) *PageData {
	t.Helper()
	res := NewView(deps).Handle(dstest.Ctx(codes...), dstest.Request(http.MethodGet, "/l/"+status, "", "status", status))
	if res.Template != "document-series-list" {
		t.Fatalf("template = %q (status %d, err %v)", res.Template, res.StatusCode, res.Error)
	}
	return res.Data.(*PageData)
}

func ids(pd *PageData) string {
	var out []string
	for _, r := range pd.Table.Rows {
		out = append(out, r.ID)
	}
	return strings.Join(out, ",")
}

func seed() *dstest.Fake {
	f := dstest.New()
	f.Add("a", "AAA", false, 1)
	f.Add("b", "BBB", false, 12)
	f.Add("r", "RRR", true, 3)
	return f
}

func TestTabsBucketSeriesByStatus(t *testing.T) {
	deps := newDeps(seed())
	perms := []string{"document_series:list", "document_series:create", "document_series:update"}
	active := render(t, deps, "active", perms...)
	if got := ids(active); got != "a,b" {
		t.Errorf("active rows = %s", got)
	}
	if got := ids(render(t, deps, "retired", perms...)); got != "r" {
		t.Errorf("retired rows = %s", got)
	}
	want := map[string]int{"active": 2, "retired": 1}
	for _, tab := range active.StatusTabs {
		if tab.Count != want[tab.Key] {
			t.Errorf("tab %s count = %d, want %d", tab.Key, tab.Count, want[tab.Key])
		}
	}
	if got := ids(render(t, deps, "bogus", perms...)); got != "a,b" {
		t.Errorf("unknown status must fall back to active, got %s", got)
	}
}

func TestPermissionsGateThePageAndDisableActions(t *testing.T) {
	deps := newDeps(seed())
	res := NewView(deps).Handle(dstest.Ctx(), dstest.Request(http.MethodGet, "/l", "", "status", "active"))
	if res.StatusCode != http.StatusForbidden {
		t.Errorf("list without permission = %d, want 403", res.StatusCode)
	}
	pd := render(t, deps, "active", "document_series:list")
	if pd.Table.PrimaryAction == nil || !pd.Table.PrimaryAction.Disabled || pd.Table.PrimaryAction.TestID != "document-series-add" {
		t.Errorf("Add must be present but disabled: %+v", pd.Table.PrimaryAction)
	}
	for _, r := range pd.Table.Rows {
		for _, a := range r.Actions {
			if !a.Disabled || a.DisabledTooltip == "" {
				t.Errorf("row %s action %s must be disabled with a tooltip", r.ID, a.Type)
			}
		}
	}
	retired := render(t, deps, "retired", "document_series:list", "document_series:update")
	for _, r := range retired.Table.Rows {
		if len(r.Actions) != 0 {
			t.Errorf("retired row %s must offer no actions", r.ID)
		}
	}
}

func TestListFailsClosedWithoutUseCase(t *testing.T) {
	deps := newDeps(seed())
	deps.UseCases = &ds.UseCases{}
	res := NewView(deps).Handle(dstest.Ctx("document_series:list"), dstest.Request(http.MethodGet, "/l", "", "status", "active"))
	if res.Error == nil {
		t.Error("an unbound list closure must fail closed")
	}
}

var _ = pyeza.TabItem{}
