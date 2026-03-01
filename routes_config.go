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
	DashboardURL    string `json:"dashboard_url"`
	RevenueURL      string `json:"revenue_url"`
	CostOfSalesURL  string `json:"cost_of_sales_url"`
	GrossProfitURL  string `json:"gross_profit_url"`
	ExpensesURL     string `json:"expenses_url"`
	NetProfitURL    string `json:"net_profit_url"`
}

// DefaultReportsRoutes returns a ReportsRoutes populated from package-level consts.
func DefaultReportsRoutes() ReportsRoutes {
	return ReportsRoutes{
		DashboardURL:    ReportsDashboardURL,
		RevenueURL:      ReportsRevenueURL,
		CostOfSalesURL:  ReportsCostOfSalesURL,
		GrossProfitURL:  ReportsGrossProfitURL,
		ExpensesURL:     ReportsExpensesURL,
		NetProfitURL:    ReportsNetProfitURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ReportsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"reports.dashboard":      r.DashboardURL,
		"reports.revenue":        r.RevenueURL,
		"reports.cost_of_sales":  r.CostOfSalesURL,
		"reports.gross_profit":   r.GrossProfitURL,
		"reports.expenses":       r.ExpensesURL,
		"reports.net_profit":     r.NetProfitURL,
	}
}
