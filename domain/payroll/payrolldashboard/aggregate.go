package payrolldashboard

// aggregate.go — dashboard-facing aggregate label/route type aliases.
//
// The payroll dashboard reads labels/routes that originate in sibling entity
// packages (run, remittance, payrollsettings) plus its own dashboard labels.
// These aliases let the dashboard view + module reference those shapes without
// importing the domain/payroll facade (which would create an import cycle:
// payroll -> payrolldashboard -> payroll).
//
// The domain/payroll facade re-exports these under the same names so external
// consumers keep writing payroll.PayrollRunRoutes, payroll.PayrollLabels, etc.

import (
	payrollsettings "github.com/erniealice/fycha-golang/domain/payroll/payrollsettings"
	remittance "github.com/erniealice/fycha-golang/domain/payroll/remittance"
	run "github.com/erniealice/fycha-golang/domain/payroll/run"
)

// PayrollRunRoutes is the payroll-run route shape (domain/payroll/run).
type PayrollRunRoutes = run.Routes

// PayrollRemittanceRoutes is the payroll-remittance route shape (domain/payroll/remittance).
type PayrollRemittanceRoutes = remittance.Routes

// PayrollSettingsRoutes is the payroll-settings route shape (domain/payroll/payrollsettings).
type PayrollSettingsRoutes = payrollsettings.Routes

// PayrollDashboardLabels is the dashboard's own label set (this package's Labels).
type PayrollDashboardLabels = Labels

// PayrollLabels is the dashboard-facing aggregate label set. The dashboard view
// reads its strings via the Dashboard field.
type PayrollLabels struct {
	Dashboard PayrollDashboardLabels `json:"dashboard"`
}
