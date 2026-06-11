package financial

import (
	report "github.com/erniealice/fycha-golang/service/report"
	balancesheetview "github.com/erniealice/fycha-golang/service/report/views/balance_sheet"
	cashflowview "github.com/erniealice/fycha-golang/service/report/views/cash_flow"
	equitychangesview "github.com/erniealice/fycha-golang/service/report/views/equity_changes"
	incomestatementview "github.com/erniealice/fycha-golang/service/report/views/income_statement"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the financial statements module.
type ModuleDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       report.ReportsLabels
}

// Module holds all constructed financial statement views.
type Module struct {
	incomeStatement view.View
	balanceSheet    view.View
	cashFlow        view.View
	equityChanges   view.View
}

// NewModule creates a financial statements module with real report views.
func NewModule(deps *ModuleDeps) *Module {
	return &Module{
		incomeStatement: incomestatementview.NewIncomeStatementView(&incomestatementview.IncomeStatementDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
		balanceSheet: balancesheetview.NewBalanceSheetView(&balancesheetview.BalanceSheetDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
		cashFlow: cashflowview.NewCashFlowView(&cashflowview.CashFlowDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
		equityChanges: equitychangesview.NewEquityChangesView(&equitychangesview.EquityChangesDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
	}
}

// RegisterRoutes registers all financial statement routes with the given route registrar.
// These routes live under the Reports app (active nav: "reports").
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(report.ReportsIncomeStatementURL, m.incomeStatement)
	r.GET(report.ReportsBalanceSheetURL, m.balanceSheet)
	r.GET(report.ReportsCashFlowURL, m.cashFlow)
	r.GET(report.ReportsEquityChangesURL, m.equityChanges)
}
