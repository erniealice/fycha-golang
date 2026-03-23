package financial

import (
	fycha "github.com/erniealice/fycha-golang"
	fsviews "github.com/erniealice/fycha-golang/views/reports"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ModuleDeps holds all dependencies for the financial statements module.
type ModuleDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       fycha.ReportsLabels
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
		incomeStatement: fsviews.NewIncomeStatementView(&fsviews.IncomeStatementDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
		balanceSheet: fsviews.NewBalanceSheetView(&fsviews.BalanceSheetDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
		cashFlow: fsviews.NewCashFlowView(&fsviews.CashFlowDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
		equityChanges: fsviews.NewEquityChangesView(&fsviews.EquityChangesDeps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Labels:       deps.Labels,
		}),
	}
}

// RegisterRoutes registers all financial statement routes with the given route registrar.
// These routes live under the Reports app (active nav: "reports").
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(fycha.ReportsIncomeStatementURL, m.incomeStatement)
	r.GET(fycha.ReportsBalanceSheetURL, m.balanceSheet)
	r.GET(fycha.ReportsCashFlowURL, m.cashFlow)
	r.GET(fycha.ReportsEquityChangesURL, m.equityChanges)
}
