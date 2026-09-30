package reports

import (
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	report "github.com/erniealice/fycha-golang/service/report"
)

type muxRegistrar struct{ mux *http.ServeMux }

func (m muxRegistrar) handle(method, path string, v view.View) {
	m.mux.HandleFunc(method+" "+path, func(w http.ResponseWriter, r *http.Request) {
		res := v.Handle(r.Context(), &view.ViewContext{Request: r})
		w.WriteHeader(res.StatusCode)
	})
}
func (m muxRegistrar) GET(path string, v view.View, _ ...string)  { m.handle(http.MethodGet, path, v) }
func (m muxRegistrar) POST(path string, v view.View, _ ...string) { m.handle(http.MethodPost, path, v) }
func (m muxRegistrar) HandleFunc(method, path string, h http.HandlerFunc, _ ...string) {
	m.mux.HandleFunc(method+" "+path, h)
}

func aging(rows ...*recoverydocumentpb.RecoverablesAgingRow) func(context.Context, *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error) {
	return func(_ context.Context, req *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error) {
		return &recoverydocumentpb.ListRecoverablesAgingResponse{AsOf: req.GetAsOf(), Rows: rows,
			Totals: []*recoverydocumentpb.RecoverablesAgingRow{{Currency: "PHP", Days_0_30: 1000, Total: 1000, DocumentCount: 1}}, Success: true}, nil
	}
}

func recon(rs ...*costsourcecomponentpb.CostSourceComponentReconciliation) func(context.Context, *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error) {
	return func(_ context.Context, _ *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error) {
		return &costsourcecomponentpb.ReconcileCostSourceResponse{Data: rs, Success: true}, nil
	}
}

func module(agingFn func(context.Context, *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error),
	reconFn func(context.Context, *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error)) *RecoveryModule {
	return NewRecoveryModule(&RecoveryModuleDeps{
		Routes: report.DefaultRecoveryReportsRoutes(), Labels: report.DefaultRecoveryReportsLabels(),
		TableLabels: types.TableLabels{}, ListRecoverablesAging: agingFn, ReconcileCostSource: reconFn,
	})
}

func ctxWith(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

func vc(target string, q map[string]string) *view.ViewContext {
	req := httptest.NewRequest(http.MethodGet, target, nil)
	return &view.ViewContext{Request: req, CurrentPath: req.URL.Path, QueryParams: q}
}

func renderer(t *testing.T) *pyeza.HTMLRenderer {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, TemplatesFS)
	r.SetRouteMap(report.DefaultRecoveryReportsRoutes().RouteMap())
	if err := r.Init(); err != nil {
		t.Fatalf("renderer init: %v", err)
	}
	return r
}

func renderTo(t *testing.T, r *pyeza.HTMLRenderer, name string, data any) string {
	t.Helper()
	w := httptest.NewRecorder()
	if err := r.Render(w, name, data); err != nil {
		t.Fatalf("render %s: %v", name, err)
	}
	return w.Body.String()
}

func TestRecoveryRoutesRegisterWithoutServeMuxConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration conflict: %v", r)
		}
	}()
	module(aging(), recon()).RegisterRoutes(muxRegistrar{http.NewServeMux()})
}

func TestRecoverablesAgingRendersTitleFiltersAndTotals(t *testing.T) {
	m := module(aging(&recoverydocumentpb.RecoverablesAgingRow{ClientId: strp("c1"), Currency: "PHP", Days_0_30: 1000, Total: 1000, DocumentCount: 1}), recon())
	res := m.RecoverablesAging.Handle(ctxWith("recovery_document:list"), vc("/reports/recoverables-aging", map[string]string{"as-of-date": "2026-09-30", "client-id": "c1"}))
	if res.Template != "recoverables-aging-report" {
		t.Fatalf("template = %q (err %v)", res.Template, res.Error)
	}
	r := renderer(t)
	body := renderTo(t, r, "recoverables-aging-report-content", res.Data)
	for _, want := range []string{"recoverables-aging-report", "recoverables-aging-row-c1|PHP", "recoverables-aging-total-PHP", "rr-chip-as-of-date", "rr-chip-client", "report-filters-open-btn", "2026-09-30"} {
		if !strings.Contains(body, want) {
			t.Errorf("aging page lacks %q", want)
		}
	}
	if pageTitle(res.Data) != "Recoverables Aging" {
		t.Errorf("page title = %q (must be a real title)", pageTitle(res.Data))
	}
}

func TestRecoverablesAgingEmptyStateAndFilterSheet(t *testing.T) {
	m := module(aging(), recon())
	res := m.RecoverablesAging.Handle(ctxWith("recovery_document:list"), vc("/reports/recoverables-aging", nil))
	body := renderTo(t, renderer(t), "recoverables-aging-report-content", res.Data)
	if !strings.Contains(body, "Nothing outstanding") {
		t.Error("empty aging report must render its empty state")
	}
	res = m.RecoverablesAging.Handle(ctxWith("recovery_document:list"), vc("/reports/recoverables-aging", map[string]string{"sheet": "filters"}))
	if res.Template != "recovery-report-filter-sheet" {
		t.Fatalf("filter sheet template = %q", res.Template)
	}
	body = renderTo(t, renderer(t), res.Template, res.Data)
	for _, want := range []string{"recoverables-aging-as-of-date", "recoverables-aging-client-id", "report-filters-apply", "report-filters-clear"} {
		if !strings.Contains(body, want) {
			t.Errorf("filter sheet lacks %q", want)
		}
	}
}

func TestRecoveryReportsFailClosed(t *testing.T) {
	m := module(aging(), recon())
	if res := m.RecoverablesAging.Handle(ctxWith(), vc("/x", nil)); res.StatusCode != http.StatusForbidden {
		t.Errorf("aging without permission = %d, want 403", res.StatusCode)
	}
	if res := m.CostSourceReconciliation.Handle(ctxWith("recovery_document:list"), vc("/x", nil)); res.StatusCode != http.StatusForbidden {
		t.Errorf("reconciliation without expenditure:read = %d, want 403", res.StatusCode)
	}
	nilMod := module(nil, nil)
	if res := nilMod.RecoverablesAging.Handle(ctxWith("recovery_document:list"), vc("/x", nil)); res.Error == nil {
		t.Error("unbound aging closure must fail closed")
	}
	if res := nilMod.CostSourceReconciliation.Handle(ctxWith("expenditure:read"), vc("/x", nil)); res.Error == nil {
		t.Error("unbound reconciliation closure must fail closed")
	}
	if res := m.RecoverablesAging.Handle(ctxWith("recovery_document:list"), vc("/x", map[string]string{"as-of-date": "not-a-date"})); res.Error == nil {
		t.Error("an invalid as-of date must be refused")
	}
}

func TestCostSourceReconciliationHighlightsMismatch(t *testing.T) {
	ok := &costsourcecomponentpb.CostSourceComponentReconciliation{ComponentId: "c-ok", Currency: "PHP", Amount: 1000, HasPublishedBatch: true, SharesTotal: 1000, RecoverableTotal: 600, Reconciled: true}
	bad := &costsourcecomponentpb.CostSourceComponentReconciliation{ComponentId: "c-bad", Currency: "PHP", Amount: 1000, HasPublishedBatch: true, SharesTotal: 1000, RecoverableTotal: 600, ChargeVariance: 100}
	none := &costsourcecomponentpb.CostSourceComponentReconciliation{ComponentId: "c-none", Currency: "PHP", Amount: 500}
	m := module(aging(), recon(ok, bad, none))
	res := m.CostSourceReconciliation.Handle(ctxWith("expenditure:read"), vc("/reports/cost-source-reconciliation", map[string]string{"expenditure-id": "e1"}))
	if res.Template != "cost-source-reconciliation-report" {
		t.Fatalf("template = %q (err %v)", res.Template, res.Error)
	}
	body := renderTo(t, renderer(t), "cost-source-reconciliation-report-content", res.Data)
	for _, want := range []string{"cost-source-reconciliation-row-c-ok", "cost-source-reconciliation-row-c-bad", "cost-source-reconciliation-row-c-none", "rr-chip-expenditure", `data-status="mismatch"`, `data-status="reconciled"`, `data-status="no_batch"`, "status-badge danger"} {
		if !strings.Contains(body, want) {
			t.Errorf("reconciliation page lacks %q", want)
		}
	}
	if pageTitle(res.Data) != "Cost Reconciliation" {
		t.Errorf("page title = %q", pageTitle(res.Data))
	}
}

func TestCostSourceReconciliationPromptsForAnExpenditure(t *testing.T) {
	called := false
	m := module(aging(), func(context.Context, *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error) {
		called = true
		return &costsourcecomponentpb.ReconcileCostSourceResponse{}, nil
	})
	res := m.CostSourceReconciliation.Handle(ctxWith("expenditure:read"), vc("/x", nil))
	body := renderTo(t, renderer(t), "cost-source-reconciliation-report-content", res.Data)
	if called || !strings.Contains(body, "Choose an expenditure") {
		t.Errorf("without an expenditure the report must prompt and not call the use case (called=%v)", called)
	}
}

func TestRecoveryExportsAreGatedAndFilteredCSV(t *testing.T) {
	m := module(aging(&recoverydocumentpb.RecoverablesAgingRow{ClientId: strp("c1"), Currency: "PHP", Total: 12345, DocumentCount: 2}),
		recon(&costsourcecomponentpb.CostSourceComponentReconciliation{ComponentId: "c1", Currency: "PHP", Amount: 5000, HasPublishedBatch: true, Reconciled: true}))
	do := func(h http.HandlerFunc, target string, ctx context.Context) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		h(w, httptest.NewRequest(http.MethodGet, target, nil).WithContext(ctx))
		return w
	}
	if w := do(m.RecoverablesAgingExport, "/e", ctxWith()); w.Code != http.StatusForbidden {
		t.Errorf("aging export without permission = %d", w.Code)
	}
	w := do(m.RecoverablesAgingExport, "/e?as-of-date=2026-09-30", ctxWith("recovery_document:list"))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "123.45") || !strings.Contains(w.Header().Get("Content-Type"), "text/csv") {
		t.Errorf("aging export = %d %q", w.Code, w.Body.String())
	}
	if w := do(m.CostSourceReconciliationExport, "/e", ctxWith("expenditure:read")); w.Code != http.StatusBadRequest {
		t.Errorf("reconciliation export without an expenditure = %d, want 400", w.Code)
	}
	w = do(m.CostSourceReconciliationExport, "/e?expenditure-id=e1", ctxWith("expenditure:read"))
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "50.00") {
		t.Errorf("reconciliation export = %d %q", w.Code, w.Body.String())
	}
}

func strp(s string) *string { return &s }

func pageTitle(data any) string { return reflectTitle(data) }

func reflectTitle(data any) string {
	v := reflect.ValueOf(data)
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	f := v.FieldByName("Title")
	if f.IsValid() && f.Kind() == reflect.String {
		return f.String()
	}
	return ""
}

// R5 m6: export error bodies come from Lyngua (Shared / common errors), never
// hard-coded English.
func TestRecoveryExportErrorBodiesAreLocalized(t *testing.T) {
	var cl pyeza.CommonLabels
	cl.Errors.PermissionDenied = "PD-LABEL"
	labels := report.DefaultRecoveryReportsLabels()
	labels.Shared.Unavailable, labels.Shared.InvalidFilter = "UNAVAILABLE-LABEL", "INVALID-LABEL"
	m := NewRecoveryModule(&RecoveryModuleDeps{Routes: report.DefaultRecoveryReportsRoutes(), Labels: labels, CommonLabels: cl, TableLabels: types.TableLabels{}})
	do := func(h http.HandlerFunc, target string, ctx context.Context) (int, string) {
		w := httptest.NewRecorder()
		h(w, httptest.NewRequest(http.MethodGet, target, nil).WithContext(ctx))
		return w.Code, strings.TrimSpace(w.Body.String())
	}
	cases := []struct {
		name string
		h    http.HandlerFunc
		url  string
		ctx  context.Context
		code int
		body string
	}{
		{"aging forbidden", m.RecoverablesAgingExport, "/e", ctxWith(), http.StatusForbidden, "PD-LABEL"},
		{"aging unwired", m.RecoverablesAgingExport, "/e", ctxWith("recovery_document:list"), http.StatusInternalServerError, "UNAVAILABLE-LABEL"},
		{"recon forbidden", m.CostSourceReconciliationExport, "/e", ctxWith(), http.StatusForbidden, "PD-LABEL"},
		{"recon bad filter", m.CostSourceReconciliationExport, "/e", ctxWith("expenditure:read"), http.StatusBadRequest, "INVALID-LABEL"},
	}
	for _, c := range cases {
		if code, body := do(c.h, c.url, c.ctx); code != c.code || body != c.body {
			t.Errorf("%s = %d %q, want %d %q", c.name, code, body, c.code, c.body)
		}
	}
}
