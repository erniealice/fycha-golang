// routes.go defines route constants and configurable route structs for the
// payroll settings entity.

package payrollsettings

const (
	GovRatesURL   = "/payroll/settings/gov-rates"
	PayPeriodsURL = "/payroll/settings/pay-periods"
)

// Routes holds route paths for Payroll Settings views.
type Routes struct {
	ActiveNav     string `json:"active_nav"`
	GovRatesURL   string `json:"gov_rates_url"`
	PayPeriodsURL string `json:"pay_periods_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:     "payroll",
		GovRatesURL:   GovRatesURL,
		PayPeriodsURL: PayPeriodsURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.settings.gov_rates":   r.GovRatesURL,
		"payroll.settings.pay_periods": r.PayPeriodsURL,
	}
}
