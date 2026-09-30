package document_series_test

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ledger "github.com/erniealice/fycha-golang/domain/ledger"
	documentseries "github.com/erniealice/fycha-golang/domain/ledger/document_series"
	"github.com/erniealice/fycha-golang/domain/ledger/document_series/internal/dstest"
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

func module(f *dstest.Fake) *ledger.DocumentSeriesModule {
	return ledger.NewDocumentSeriesModule(&ledger.DocumentSeriesModuleDeps{
		Routes: documentseries.DefaultRoutes(), Labels: documentseries.DefaultLabels(),
		CommonLabels: dstest.Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases(),
	})
}

// The live ServeMux panics on conflicting patterns, so registering every
// route on a real mux proves there is no self-conflict.
func TestRoutesRegisterWithoutServeMuxConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration conflict: %v", r)
		}
	}()
	module(dstest.New()).RegisterRoutes(muxRegistrar{http.NewServeMux()})
}

func TestRouteMapKeys(t *testing.T) {
	m := documentseries.DefaultRoutes().RouteMap()
	want := map[string]string{
		"document_series.list":   "/ledger/settings/document-series/list/{status}",
		"document_series.table":  "/action/document-series/table/{status}",
		"document_series.add":    "/action/document-series/add",
		"document_series.edit":   "/action/document-series/edit/{id}",
		"document_series.retire": "/action/document-series/retire/{id}",
	}
	if !reflect.DeepEqual(m, want) {
		t.Errorf("route map = %v", m)
	}
}

func renderer(t *testing.T) *pyeza.HTMLRenderer {
	t.Helper()
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, documentseries.TemplatesFS)
	r.SetRouteMap(documentseries.DefaultRoutes().RouteMap())
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

// Rendering every page/partial catches template execution errors and proves the
// stable testids the E2E click path depends on.
func TestTemplatesRenderWithViewData(t *testing.T) {
	f := dstest.New()
	f.Add("s1", "SOA", false, 4)
	f.Add("s2", "OLD", true, 9)
	m := module(f)
	ctx := dstest.Ctx("document_series:list", "document_series:create", "document_series:update")
	r := renderer(t)

	run := func(v view.View, method, target, want string, pv ...string) string {
		t.Helper()
		res := v.Handle(ctx, dstest.Request(method, target, "", pv...))
		if res.Template != want {
			t.Fatalf("%s: template %q (status %d, err %v), want %q", target, res.Template, res.StatusCode, res.Error, want)
		}
		if rv := reflect.ValueOf(res.Data); rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
			if fld := rv.Elem().FieldByName("CommonLabels"); fld.IsValid() && fld.Kind() == reflect.Interface && fld.CanSet() {
				fld.Set(reflect.ValueOf(dstest.Common()))
			}
		}
		if res.Template == "document-series-list" {
			return renderTo(t, r, res.Template+"-content", res.Data)
		}
		return renderTo(t, r, res.Template, res.Data)
	}

	body := run(m.List, http.MethodGet, "/l", "document-series-list", "status", "active")
	for _, id := range []string{"document-series-add", "document-series-row-s1", "document-series-status-tabs"} {
		if !strings.Contains(body, id) {
			t.Errorf("active list lacks %s", id)
		}
	}
	if strings.Contains(body, "document-series-row-s2") {
		t.Error("retired series leaked into the active tab")
	}
	body = run(m.List, http.MethodGet, "/l", "document-series-list", "status", "retired")
	if !strings.Contains(body, "document-series-row-s2") || strings.Contains(body, "document-series-add") {
		t.Error("retired tab must list the retired series and offer no Add")
	}
	run(m.Table, http.MethodGet, "/t", "table-card", "status", "active")

	body = run(m.Add, http.MethodGet, "/a", "document-series-drawer-form")
	for _, id := range []string{"document-series-form", "document-series-code", "document-series-name", "document-series-issuer-name",
		"document-series-issuer-tax-id", "document-series-document-kind", "document-series-prefix", "document-series-branch-code",
		"document-series-fiscal-reset", "document-series-number-padding", "document-series-next-number"} {
		if !strings.Contains(body, `id="`+id+`"`) {
			t.Errorf("add drawer lacks input id %s", id)
		}
	}
	body = run(m.Edit, http.MethodGet, "/e", "document-series-drawer-form", "id", "s1")
	for _, want := range []string{"document-series-next-preview", "SOA-000004", "document-series-locked-notice"} {
		if !strings.Contains(body, want) {
			t.Errorf("edit drawer lacks %s", want)
		}
	}
	// Code, kind, restart mode and next number are read-only on edit.
	if !strings.Contains(body, "code_display") || !strings.Contains(body, "next_number_display") {
		t.Error("edit drawer must show code and next number read-only")
	}
}
