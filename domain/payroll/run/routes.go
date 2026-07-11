// routes.go defines route constants and configurable route structs for the
// payroll run entity.

package run

const (
	DashboardURL = "/payroll/dashboard"
	ListURL      = "/payroll/runs/{status}"
	DetailURL    = "/payroll/runs/detail/{id}"
)

// Routes holds route paths for Payroll Run views.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	DashboardURL string `json:"dashboard_url"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav: "payroll",
		// bare prefix — list page composes <base>-<status> to match sidebar item keys
		ActiveSubNav: "payroll",
		DashboardURL: DashboardURL,
		ListURL:      ListURL,
		DetailURL:    DetailURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.dashboard":  r.DashboardURL,
		"payroll.run.list":   r.ListURL,
		"payroll.run.detail": r.DetailURL,
	}
}
