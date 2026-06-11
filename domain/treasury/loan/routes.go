// routes.go defines configurable route structs for the loan entity.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultRoutes()
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
package loan

const (
	// Funding — Loans
	DashboardURL    = "/funding/loans/dashboard"
	ListURL         = "/funding/loans/list/{status}"
	DetailURL       = "/funding/loans/detail/{id}"
	AddURL          = "/action/funding/loan/add"
	AmortizationURL = "/funding/loans/amortization"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Loan views.
type Routes struct {
	ActiveNav       string `json:"active_nav"`
	DashboardURL    string `json:"dashboard_url"`
	ListURL         string `json:"list_url"`
	DetailURL       string `json:"detail_url"`
	AddURL          string `json:"add_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:       "loan",
		DashboardURL:    DashboardURL,
		ListURL:         ListURL,
		DetailURL:       DetailURL,
		AddURL:          AddURL,
		AmortizationURL: AmortizationURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"loan.dashboard":    r.DashboardURL,
		"loan.list":         r.ListURL,
		"loan.detail":       r.DetailURL,
		"loan.add":          r.AddURL,
		"loan.amortization": r.AmortizationURL,
	}
}
