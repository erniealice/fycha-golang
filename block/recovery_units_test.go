package block

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"github.com/erniealice/espyna-golang/consumer/compose"
	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	documentseriespb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	"github.com/erniealice/pyeza-golang/view"
)

type nopRegistrar struct{ mux *http.ServeMux }

func (n nopRegistrar) reg(method, path string) {
	n.mux.HandleFunc(method+" "+path, func(http.ResponseWriter, *http.Request) {})
}
func (n nopRegistrar) GET(path string, _ view.View, _ ...string)  { n.reg("GET", path) }
func (n nopRegistrar) POST(path string, _ view.View, _ ...string) { n.reg("POST", path) }
func (n nopRegistrar) HandleFunc(method, path string, _ http.HandlerFunc, _ ...string) {
	n.reg(method, path)
}

func unitKeys(us []compose.Unit) map[string]bool {
	m := map[string]bool{}
	for _, u := range us {
		m[u.Key] = true
	}
	return m
}

// The document series and recovery report units are opt-in: without their
// options the unit set is unchanged (service-admin / school-admin stay
// byte-identical), with them exactly one unit each is appended.
func TestDocumentSeriesAndRecoveryReportsAreOptIn(t *testing.T) {
	base := AllUnits(nil, nil)
	keys := unitKeys(base)
	if keys["ledger.document_series"] || keys["report.recovery"] {
		t.Fatal("opt-in units present without their options")
	}
	if got := len(AllUnits(nil, nil, WithDocumentSeries(false), WithRecoveryReports(false))); got != len(base) {
		t.Fatalf("disabled options changed the unit set: %d vs %d", got, len(base))
	}
	on := AllUnits(nil, nil, WithDocumentSeries(true), WithRecoveryReports(true))
	if len(on) != len(base)+2 {
		t.Fatalf("enabled options added %d units, want 2", len(on)-len(base))
	}
	onKeys := unitKeys(on)
	if !onKeys["ledger.document_series"] || !onKeys["report.recovery"] {
		t.Fatalf("unit keys = %v", onKeys)
	}
	for i := range base {
		if base[i].Key != on[i].Key {
			t.Fatalf("default unit order changed at %d", i)
		}
	}
	for _, u := range on[len(base):] {
		switch u.Key {
		case "ledger.document_series":
			if u.Nav.Permission != "document_series:list" || len(u.Nav.Items) != 1 || u.Nav.Items[0].Route != "document_series.list" {
				t.Errorf("document series nav = %+v", u.Nav)
			}
		case "report.recovery":
			if len(u.Nav.Items) != 2 || u.Nav.Items[0].Route != "reports.recoverables_aging" || u.Nav.Items[1].Route != "reports.cost_source_reconciliation" {
				t.Errorf("recovery report nav = %+v", u.Nav)
			}
		}
	}
}

// R4 M1: compose-v2 Mount does not run RequireFor, so each new Unit must refuse
// boot itself when a required closure is unbound (or the UseCases are absent).
func TestDocumentSeriesUnitMountFailsClosed(t *testing.T) {
	for name, uc := range map[string]*UseCases{"empty": {}, "nil": nil} {
		units := AllUnits(uc, nil, WithDocumentSeries(true))
		var u compose.Unit
		for _, x := range units {
			if x.Key == "ledger.document_series" {
				u = x
			}
		}
		if err := u.Mount(&compose.MountContext{}); err == nil {
			t.Errorf("%s: Mount accepted unbound document series closures", name)
		}
	}
}

func TestRecoveryReportsUnitMountFailsClosed(t *testing.T) {
	for name, uc := range map[string]*UseCases{"empty": {}, "nil": nil} {
		units := AllUnits(uc, nil, WithRecoveryReports(true))
		var u compose.Unit
		for _, x := range units {
			if x.Key == "report.recovery" {
				u = x
			}
		}
		if err := u.Mount(&compose.MountContext{}); err == nil {
			t.Errorf("%s: Mount accepted unbound recovery report closures", name)
		}
	}
	half := &UseCases{}
	half.RecoveryReports.ListRecoverablesAging = func(context.Context, *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error) {
		return nil, nil
	}
	if err := RecoveryReportsUnit(half, nil).Mount(&compose.MountContext{}); err == nil {
		t.Error("Mount accepted a half-bound recovery report set")
	}
}

func TestRequireHelpersListExactlyTheMissingClosures(t *testing.T) {
	if err := requireDocumentSeries(&UseCases{}); err == nil || strings.Count(err.Error(), "UseCases.DocumentSeries.") != 4 {
		t.Errorf("empty document series: %v, want 4 missing", err)
	}
	if err := requireRecoveryReports(&UseCases{}); err == nil || strings.Count(err.Error(), "UseCases.RecoveryReports.") != 2 {
		t.Errorf("empty recovery reports: %v, want 2 missing", err)
	}
	if err := requireRecoveryReports(&UseCases{RecoveryReports: RecoveryReportUseCases{}}); err == nil || !strings.Contains(err.Error(), "ListRecoverablesAging, ") {
		t.Errorf("missing names must be sorted: %v", err)
	}
}

// With every closure bound, Mount registers routes on a real ServeMux (a
// conflict would panic) — the composition-level proof for the new units.
func TestBoundUnitsMountAndRegisterRoutes(t *testing.T) {
	uc := &UseCases{}
	uc.DocumentSeries = DocumentSeriesUseCases{
		CreateDocumentSeries: func(context.Context, *documentseriespb.CreateDocumentSeriesRequest) (*documentseriespb.CreateDocumentSeriesResponse, error) {
			return nil, nil
		},
		ReadDocumentSeries: func(context.Context, *documentseriespb.ReadDocumentSeriesRequest) (*documentseriespb.ReadDocumentSeriesResponse, error) {
			return nil, nil
		},
		UpdateDocumentSeries: func(context.Context, *documentseriespb.UpdateDocumentSeriesRequest) (*documentseriespb.UpdateDocumentSeriesResponse, error) {
			return nil, nil
		},
		GetDocumentSeriesListPageData: func(context.Context, *documentseriespb.GetDocumentSeriesListPageDataRequest) (*documentseriespb.GetDocumentSeriesListPageDataResponse, error) {
			return nil, nil
		},
	}
	uc.RecoveryReports = RecoveryReportUseCases{
		ListRecoverablesAging: func(context.Context, *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error) {
			return nil, nil
		},
		ReconcileCostSource: func(context.Context, *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error) {
			return nil, nil
		},
	}
	reg := nopRegistrar{http.NewServeMux()}
	for _, u := range []compose.Unit{DocumentSeriesUnit(uc, nil), RecoveryReportsUnit(uc, nil)} {
		if err := u.Mount(&compose.MountContext{Routes: reg}); err != nil {
			t.Errorf("%s: Mount = %v", u.Key, err)
		}
	}
}

func TestRequireUnitSortedAndNilWhenComplete(t *testing.T) {
	if err := requireUnit("x", map[string]bool{"b": true, "a": true}); err != nil {
		t.Errorf("complete set = %v, want nil", err)
	}
	err := requireUnit("x", map[string]bool{"z": false, "a": false, "m": true})
	if err == nil || !strings.HasSuffix(err.Error(), "a, z") {
		t.Errorf("missing names not sorted: %v", err)
	}
}
