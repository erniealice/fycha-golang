package charge_policy_test

import (
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"
	"testing/fstest"

	componentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_component"
	postingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_posting"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ledger "github.com/erniealice/fycha-golang/domain/ledger"
	chargepolicy "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	"github.com/erniealice/fycha-golang/domain/ledger/charge_policy/internal/cptest"
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

func module(f *cptest.Fake) *ledger.ChargePolicyModule {
	return ledger.NewChargePolicyModule(&ledger.ChargePolicyModuleDeps{
		Routes: chargepolicy.DefaultRoutes(), Labels: chargepolicy.DefaultLabels(),
		CommonLabels: cptest.Common(), TableLabels: types.TableLabels{}, UseCases: f.UseCases(),
	})
}

// The live ServeMux panics on conflicting patterns, so registering every
// charge policy route on a real mux proves there is no self-conflict.
func TestRoutesRegisterWithoutServeMuxConflict(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("route registration conflict: %v", r)
		}
	}()
	module(cptest.New()).RegisterRoutes(muxRegistrar{http.NewServeMux()})
}

func TestRouteMapKeysMatchSpec(t *testing.T) {
	m := chargepolicy.DefaultRoutes().RouteMap()
	want := map[string]string{
		"charge_policy.list":            "/ledger/settings/charge-policies/list/{status}",
		"charge_policy.table":           "/action/charge-policy/table/{status}",
		"charge_policy.detail":          "/ledger/settings/charge-policies/detail/{id}",
		"charge_policy.tab_action":      "/action/charge-policy/{id}/tab/{tab}",
		"charge_policy.version.detail":  "/ledger/settings/charge-policies/detail/{id}/version/{vid}",
		"charge_policy.version.approve": "/action/charge-policy/{id}/version/approve/{vid}",
	}
	for k, v := range want {
		if m[k] != v {
			t.Errorf("%s = %q, want %q", k, m[k], v)
		}
	}
	// C14: sub-resources are dotted, never flat `<sub>_<verb>` keys.
	for k := range m {
		for _, flat := range []string{"version_", "component_", "posting_", "attachment_"} {
			if strings.HasPrefix(k, "charge_policy."+flat) {
				t.Errorf("flat sub-resource route key %q (use charge_policy.%s.<verb>)", k, strings.TrimSuffix(flat, "_"))
			}
		}
	}
	for k, v := range m {
		if v == "" || !strings.HasPrefix(k, "charge_policy.") {
			t.Errorf("bad route entry %q=%q", k, v)
		}
	}
}

// renderer builds a real pyeza renderer over the shared components plus the
// charge policy templates; rendering every page/partial catches template
// execution errors (missing fields, bad component arguments).
func renderer(t *testing.T) *pyeza.HTMLRenderer {
	t.Helper()
	// The app shell is app-owned; a stub keeps the full-page templates parseable.
	shell := fstest.MapFS{"app-shell.html": {Data: []byte(`{{define "app-shell"}}[shell]{{end}}`)}}
	r := pyeza.NewHTMLRendererFromFS(pyeza.SharedFS, shell, chargepolicy.TemplatesFS)
	r.SetRouteMap(chargepolicy.DefaultRoutes().RouteMap())
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

func TestTemplatesRenderWithViewData(t *testing.T) {
	f := cptest.New()
	f.AddPolicy("p1", "UTILITY", "Utility recovery", false, cptest.Approved("v-appr", 1), cptest.Draft("v-draft", 2, "u-prep"))
	f.Components["v-draft"] = []*componentpb.ChargePolicyComponent{{Id: "c1", ChargePolicyVersionId: "v-draft",
		ComponentRole: enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_RECOVERY_COST}}
	f.Postings["v-draft"] = []*postingpb.ChargePolicyPosting{{Id: "x1", ChargePolicyVersionId: "v-draft",
		Event: enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ISSUE, PostingRole: enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_RECEIVABLE, AccountId: "a1"}}
	f.Validation = &versionpb.ChargePolicyApprovalChecklist{Supported: false, Reasons: []string{"own_supply"},
		Items: []*versionpb.ChargePolicyChecklistItem{{Code: "assessment_scope", Passed: false}}}
	m := module(f)
	ctx := cptest.Ctx("charge_policy:list", "charge_policy:read", "charge_policy:create", "charge_policy:update", "charge_policy:retire", "charge_policy:approve", "charge_policy:delete")
	r := renderer(t)

	run := func(v view.View, method, target string, want string, pv ...string) string {
		t.Helper()
		res := v.Handle(ctx, cptest.Request(method, target, "", pv...))
		if res.Template != want {
			t.Fatalf("%s: template %q (status %d, err %v), want %q", target, res.Template, res.StatusCode, res.Error, want)
		}
		// Drawer data carries CommonLabels injected by the ViewAdapter in production.
		if rv := reflect.ValueOf(res.Data); rv.Kind() == reflect.Ptr && rv.Elem().Kind() == reflect.Struct {
			if f := rv.Elem().FieldByName("CommonLabels"); f.IsValid() && f.Kind() == reflect.Interface && f.CanSet() {
				f.Set(reflect.ValueOf(cptest.Common()))
			}
		}
		out := renderTo(t, r, res.Template, res.Data)
		// Full pages wrap a "-content" partial: render it too (the shell is stubbed).
		if strings.HasPrefix(res.Template, "charge-policy-") && strings.Count(res.Template, "-") <= 3 &&
			(res.Template == "charge-policy-list" || res.Template == "charge-policy-detail" || res.Template == "charge-policy-version-detail") {
			out = renderTo(t, r, res.Template+"-content", res.Data)
		}
		return out
	}

	body := run(m.List, http.MethodGet, "/l", "charge-policy-list", "status", "active")
	if !strings.Contains(body, `data-testid="charge-policy-add"`) || !strings.Contains(body, "charge-policy-row-p1") {
		t.Error("list page lacks its testids")
	}
	run(m.Table, http.MethodGet, "/t", "table-card", "status", "draft")

	body = run(m.Detail, http.MethodGet, "/d?tab=info", "charge-policy-detail", "id", "p1")
	for _, id := range []string{"charge-policy-edit", "charge-policy-retire", "charge-policy-status"} {
		if !strings.Contains(body, `data-testid="`+id+`"`) {
			t.Errorf("detail info lacks testid %s", id)
		}
	}
	for _, tab := range []string{"versions", "usage"} {
		run(m.TabAction, http.MethodGet, "/x", "charge-policy-tab-"+tab, "id", "p1", "tab", tab)
	}
	run(m.TabAction, http.MethodGet, "/x", "audit-history-tab", "id", "p1", "tab", "audit-history")

	body = run(m.VersionDetail, http.MethodGet, "/v?tab=approval", "charge-policy-version-detail", "id", "p1", "vid", "v-draft")
	for _, id := range []string{"charge-policy-approve", "charge-policy-checklist", "charge-policy-unsupported", "charge-policy-version-status"} {
		if !strings.Contains(body, `data-testid="`+id+`"`) {
			t.Errorf("version approval lacks testid %s", id)
		}
	}
	// The unsupported combination disables Approve.
	if !regexp.MustCompile(`(?s)<button[^>]*data-testid="charge-policy-approve"[^>]*disabled`).MatchString(body) {
		t.Error("Approve is not disabled for an unsupported combination")
	}
	for _, tab := range []string{"classification", "components", "postings", "approval"} {
		run(m.VersionTabAction, http.MethodGet, "/x", "charge-policy-version-tab-"+tab, "id", "p1", "vid", "v-draft", "tab", tab)
	}
	if body = run(m.VersionDetail, http.MethodGet, "/v", "charge-policy-version-detail", "id", "p1", "vid", "v-appr"); !strings.Contains(body, "charge-policy-locked-notice") {
		t.Error("approved version page lacks the locked notice")
	}

	// Drawers.
	run(m.Add, http.MethodGet, "/a", "charge-policy-drawer-form")
	run(m.Edit, http.MethodGet, "/e", "charge-policy-drawer-form", "id", "p1")
	run(m.VersionAdd, http.MethodGet, "/va", "charge-policy-new-version-drawer-form", "id", "p1")
	run(m.VersionEdit, http.MethodGet, "/ve", "charge-policy-version-drawer-form", "id", "p1", "vid", "v-draft")
	run(m.ComponentAdd, http.MethodGet, "/ca", "charge-policy-component-drawer-form", "id", "p1", "vid", "v-draft")
	run(m.ComponentEdit, http.MethodGet, "/ce", "charge-policy-component-drawer-form", "id", "p1", "vid", "v-draft", "cid", "c1")
	run(m.PostingAdd, http.MethodGet, "/pa", "charge-policy-posting-drawer-form", "id", "p1", "vid", "v-draft")
	run(m.PostingEdit, http.MethodGet, "/pe", "charge-policy-posting-drawer-form", "id", "p1", "vid", "v-draft", "pid", "x1")
}
