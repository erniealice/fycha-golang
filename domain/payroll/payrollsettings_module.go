package payroll

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	payrollsettings "github.com/erniealice/fycha-golang/domain/payroll/payrollsettings"
)

// PayrollSettingsModuleDeps holds all dependencies for the payroll settings module.
type PayrollSettingsModuleDeps struct {
	Routes       payrollsettings.Routes
	Labels       payrollsettings.Labels
	CommonLabels pyeza.CommonLabels
}

// PayrollSettingsModule holds all constructed payroll settings views.
type PayrollSettingsModule struct {
	routes     payrollsettings.Routes
	GovRates   view.View
	PayPeriods view.View
}

// NewPayrollSettingsModule creates a payroll settings module.
func NewPayrollSettingsModule(deps *PayrollSettingsModuleDeps) *PayrollSettingsModule {
	if deps == nil {
		deps = &PayrollSettingsModuleDeps{}
	}

	govRatesDeps := &payrollsettings.GovRatesDeps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
	}

	payPeriodsDeps := &payrollsettings.PayPeriodsDeps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
	}

	return &PayrollSettingsModule{
		routes:     deps.Routes,
		GovRates:   payrollsettings.NewGovRatesView(govRatesDeps),
		PayPeriods: payrollsettings.NewPayPeriodsView(payPeriodsDeps),
	}
}

// RegisterRoutes registers all payroll settings routes with the given route registrar.
func (m *PayrollSettingsModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.GovRates != nil && m.routes.GovRatesURL != "" {
		r.GET(m.routes.GovRatesURL, m.GovRates)
	}
	if m.PayPeriods != nil && m.routes.PayPeriodsURL != "" {
		r.GET(m.routes.PayPeriodsURL, m.PayPeriods)
	}
}
