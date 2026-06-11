package block

// dashboard_wiring.go wires dashboard closures from the fycha UseCases struct
// into the module ModuleDeps callbacks.
//
// Previously this file used reflection to call espyna's internal use cases
// (which had unreachable request/response types). The closure approach replaces
// that: service-admin's adapter builds the correctly-typed closures, and
// this file simply assigns them to the module deps.
//
// All assignments are nil-safe: if a closure is not supplied, the dashboard
// view renders empty state (the existing fallback behaviour).

import (
	equitymod "github.com/erniealice/fycha-golang/domain/ledger/views/equity"
	ledgermod "github.com/erniealice/fycha-golang/domain/ledger/views/ledger"
	loansmod "github.com/erniealice/fycha-golang/domain/treasury/views/loans"
	payrollmod "github.com/erniealice/fycha-golang/domain/payroll/views/payroll"
)

// wireLedgerDashboard sets ledgerDeps.GetLedgerDashboardPageData from the
// typed closure on the UseCases struct.
func wireLedgerDashboard(deps *ledgermod.ModuleDeps, uc *UseCases) {
	if uc == nil || uc.GetLedgerDashboardPageData == nil {
		return
	}
	deps.GetLedgerDashboardPageData = uc.GetLedgerDashboardPageData
}

// wireEquityDashboard sets equityDeps.GetEquityDashboardPageData from the
// typed closure on the UseCases struct.
func wireEquityDashboard(deps *equitymod.ModuleDeps, uc *UseCases) {
	if uc == nil || uc.GetEquityDashboardPageData == nil {
		return
	}
	deps.GetEquityDashboardPageData = uc.GetEquityDashboardPageData
}

// wirePayrollDashboard sets payrollDeps.GetPayrollDashboardPageData from the
// typed closure on the UseCases struct.
func wirePayrollDashboard(deps *payrollmod.ModuleDeps, uc *UseCases) {
	if uc == nil || uc.GetPayrollDashboardPageData == nil {
		return
	}
	deps.GetPayrollDashboardPageData = uc.GetPayrollDashboardPageData
}

// wireLoansDashboard sets loansDeps.GetLoanDashboardPageData from the
// typed closure on the UseCases struct.
func wireLoansDashboard(deps *loansmod.ModuleDeps, uc *UseCases) {
	if uc == nil || uc.GetLoanDashboardPageData == nil {
		return
	}
	deps.GetLoanDashboardPageData = uc.GetLoanDashboardPageData
}
