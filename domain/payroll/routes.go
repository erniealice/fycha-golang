// routes.go defines route constants and configurable route structs for the
// payroll domain.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultXxxRoutes()
//     constructors return structs populated from the package-level route constants
//     defined below. These serve as sensible defaults for any consumer app.
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

package payroll

const (
	// Payroll
	PayrollDashboardURL      = "/payroll/dashboard"
	PayrollRunListURL        = "/payroll/runs/{status}"
	PayrollRunDetailURL      = "/payroll/runs/detail/{id}"
	PayrollRemittanceListURL = "/payroll/remittances/{status}"
	PayrollEmployeeListURL   = "/payroll/employees"
	PayrollGovRatesURL       = "/payroll/settings/gov-rates"
	PayrollPayPeriodsURL     = "/payroll/settings/pay-periods"
)

// ---------------------------------------------------------------------------
// PayrollRunRoutes
// ---------------------------------------------------------------------------

// PayrollRunRoutes holds route paths for Payroll Run views.
type PayrollRunRoutes struct {
	ActiveNav    string `json:"active_nav"`
	DashboardURL string `json:"dashboard_url"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
}

func DefaultPayrollRunRoutes() PayrollRunRoutes {
	return PayrollRunRoutes{
		ActiveNav:    "payroll",
		DashboardURL: PayrollDashboardURL,
		ListURL:      PayrollRunListURL,
		DetailURL:    PayrollRunDetailURL,
	}
}

func (r PayrollRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.dashboard":  r.DashboardURL,
		"payroll.run.list":   r.ListURL,
		"payroll.run.detail": r.DetailURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollRemittanceRoutes
// ---------------------------------------------------------------------------

// PayrollRemittanceRoutes holds route paths for Payroll Remittance views.
type PayrollRemittanceRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
}

func DefaultPayrollRemittanceRoutes() PayrollRemittanceRoutes {
	return PayrollRemittanceRoutes{
		ActiveNav: "payroll",
		ListURL:   PayrollRemittanceListURL,
	}
}

func (r PayrollRemittanceRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.remittance.list": r.ListURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollEmployeeRoutes
// ---------------------------------------------------------------------------

// PayrollEmployeeRoutes holds route paths for Payroll Employee views.
type PayrollEmployeeRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
}

func DefaultPayrollEmployeeRoutes() PayrollEmployeeRoutes {
	return PayrollEmployeeRoutes{
		ActiveNav: "payroll",
		ListURL:   PayrollEmployeeListURL,
	}
}

func (r PayrollEmployeeRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.employee.list": r.ListURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollSettingsRoutes
// ---------------------------------------------------------------------------

// PayrollSettingsRoutes holds route paths for Payroll Settings views.
type PayrollSettingsRoutes struct {
	ActiveNav     string `json:"active_nav"`
	GovRatesURL   string `json:"gov_rates_url"`
	PayPeriodsURL string `json:"pay_periods_url"`
}

func DefaultPayrollSettingsRoutes() PayrollSettingsRoutes {
	return PayrollSettingsRoutes{
		ActiveNav:     "payroll",
		GovRatesURL:   PayrollGovRatesURL,
		PayPeriodsURL: PayrollPayPeriodsURL,
	}
}

func (r PayrollSettingsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.settings.gov_rates":   r.GovRatesURL,
		"payroll.settings.pay_periods": r.PayPeriodsURL,
	}
}
