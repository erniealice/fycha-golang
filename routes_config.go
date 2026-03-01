// routes_config.go defines configurable route structs for fycha views.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultXxxRoutes()
//     constructors return structs populated from the package-level route constants
//     defined in routes.go. These serve as sensible defaults for any consumer app.
//   - Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//     Apps can load a JSON config file that maps route keys to custom paths,
//     allowing industry templates (e.g. salon, retail) to rebrand URLs without
//     code changes. The json struct tags on each field support this workflow.
//   - Level 3: App-specific overrides via Go field assignment (optional).
//     After constructing defaults (and optionally applying JSON), consumer apps
//     can directly assign individual struct fields for one-off customizations.
//
// RouteMap() methods return a map[string]string of dot-notation keys to route
// paths, useful for template rendering and route resolution at runtime.
package fycha

// ReportsRoutes holds route paths for all reporting views.
type ReportsRoutes struct {
	DashboardURL   string `json:"dashboard_url"`
	GrossProfitURL string `json:"gross_profit_url"`
}

// DefaultReportsRoutes returns a ReportsRoutes populated from package-level consts.
func DefaultReportsRoutes() ReportsRoutes {
	return ReportsRoutes{
		DashboardURL:   ReportsDashboardURL,
		GrossProfitURL: ReportsGrossProfitURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ReportsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"reports.dashboard":    r.DashboardURL,
		"reports.gross_profit": r.GrossProfitURL,
	}
}
