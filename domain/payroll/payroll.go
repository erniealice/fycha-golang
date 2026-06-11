package payroll

// payroll.go — payroll-domain facade.
//
// Re-exports the payroll entity packages' Routes/Labels with the entity prefix
// so consumers keep writing payroll.PayrollRunRoutes, payroll.PayrollLabels,
// payroll.DefaultPayrollRunRoutes(), etc.
//
// The dashboard-facing aggregate types (PayrollRunRoutes, PayrollRemittanceRoutes,
// PayrollSettingsRoutes, PayrollDashboardLabels, PayrollLabels) originate in the
// payrolldashboard entity package — defining them there (rather than here) keeps
// the entity DAG acyclic: payroll -> payrolldashboard, never the reverse.

import (
	payrolldashboard "github.com/erniealice/fycha-golang/domain/payroll/payrolldashboard"
	payrollemployee "github.com/erniealice/fycha-golang/domain/payroll/payrollemployee"
	payrollsettings "github.com/erniealice/fycha-golang/domain/payroll/payrollsettings"
	remittance "github.com/erniealice/fycha-golang/domain/payroll/remittance"
	run "github.com/erniealice/fycha-golang/domain/payroll/run"
)

// ---------------------------------------------------------------------------
// Dashboard-facing aggregate types (defined in payrolldashboard).
// ---------------------------------------------------------------------------

type PayrollRunRoutes = payrolldashboard.PayrollRunRoutes
type PayrollRemittanceRoutes = payrolldashboard.PayrollRemittanceRoutes
type PayrollSettingsRoutes = payrolldashboard.PayrollSettingsRoutes
type PayrollDashboardLabels = payrolldashboard.PayrollDashboardLabels
type PayrollLabels = payrolldashboard.PayrollLabels

// ---------------------------------------------------------------------------
// Run (domain/payroll/run)
// ---------------------------------------------------------------------------

type PayrollRunLabels = run.Labels

func DefaultPayrollRunRoutes() PayrollRunRoutes { return run.DefaultRoutes() }
func DefaultPayrollRunLabels() PayrollRunLabels { return run.DefaultLabels() }

// ---------------------------------------------------------------------------
// Remittance (domain/payroll/remittance)
// ---------------------------------------------------------------------------

type PayrollRemittanceLabels = remittance.Labels

func DefaultPayrollRemittanceRoutes() PayrollRemittanceRoutes { return remittance.DefaultRoutes() }
func DefaultPayrollRemittanceLabels() PayrollRemittanceLabels { return remittance.DefaultLabels() }

// ---------------------------------------------------------------------------
// Settings (domain/payroll/payrollsettings)
// ---------------------------------------------------------------------------

type PayrollSettingsLabels = payrollsettings.Labels

func DefaultPayrollSettingsRoutes() PayrollSettingsRoutes { return payrollsettings.DefaultRoutes() }
func DefaultPayrollSettingsLabels() PayrollSettingsLabels { return payrollsettings.DefaultLabels() }

// ---------------------------------------------------------------------------
// Dashboard (domain/payroll/payrolldashboard)
// ---------------------------------------------------------------------------

func DefaultPayrollDashboardLabels() PayrollDashboardLabels { return payrolldashboard.DefaultLabels() }

// ---------------------------------------------------------------------------
// Employee (domain/payroll/payrollemployee)
// ---------------------------------------------------------------------------

type PayrollEmployeeRoutes = payrollemployee.Routes
type PayrollEmployeeLabels = payrollemployee.Labels

func DefaultPayrollEmployeeRoutes() PayrollEmployeeRoutes { return payrollemployee.DefaultRoutes() }
func DefaultPayrollEmployeeLabels() PayrollEmployeeLabels { return payrollemployee.DefaultLabels() }
