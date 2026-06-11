// routes.go defines route constants and configurable route structs for the
// payroll dashboard entity.

package payrolldashboard

// Note: DashboardURL is defined in the run entity as it shares the same URL space.
// This entity does not own additional URL constants.

// Routes holds route paths for Payroll Dashboard views.
// Cross-entity route references (run, remittance, settings) are resolved
// by the facade layer at wiring time.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	DashboardURL string `json:"dashboard_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "payroll",
		DashboardURL: "/payroll/dashboard",
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.dashboard": r.DashboardURL,
	}
}
