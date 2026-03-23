package reports

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	costsales "github.com/erniealice/fycha-golang/views/reports/cost_of_sales"
	dashboardview "github.com/erniealice/fycha-golang/views/reports/dashboard"
	expensesview "github.com/erniealice/fycha-golang/views/reports/expenses"
	grossprofit "github.com/erniealice/fycha-golang/views/reports/gross_profit"
	netprofit "github.com/erniealice/fycha-golang/views/reports/net_profit"
	revenue "github.com/erniealice/fycha-golang/views/reports/revenue"
)

// ModuleDeps holds all dependencies for the report module.
type ModuleDeps struct {
	Routes       fycha.ReportsRoutes
	DB           fycha.DataSource
	Labels       fycha.ReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// Module holds all constructed report views.
type Module struct {
	routes      fycha.ReportsRoutes
	Dashboard   view.View
	Revenue     view.View
	CostOfSales view.View
	GrossProfit view.View
	Expenses    view.View
	NetProfit   view.View
}

func NewModule(deps *ModuleDeps) *Module {
	viewDeps := &grossprofit.Deps{
		DB:           deps.DB,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}
	return &Module{
		routes:      deps.Routes,
		Dashboard:   dashboardview.NewView(&dashboardview.Deps{Routes: deps.Routes, DB: deps.DB, Labels: deps.Labels, CommonLabels: deps.CommonLabels}),
		Revenue:     revenue.NewView(&revenue.Deps{DB: deps.DB, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels}),
		CostOfSales: costsales.NewView(&costsales.Deps{DB: deps.DB, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels}),
		GrossProfit: grossprofit.NewView(viewDeps),
		Expenses:    expensesview.NewView(&expensesview.Deps{DB: deps.DB, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels}),
		NetProfit:   netprofit.NewView(&netprofit.Deps{DB: deps.DB, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels}),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.RevenueURL, m.Revenue)
	r.GET(m.routes.CostOfSalesURL, m.CostOfSales)
	r.GET(m.routes.GrossProfitURL, m.GrossProfit)
	r.GET(m.routes.ExpensesURL, m.Expenses)
	r.GET(m.routes.NetProfitURL, m.NetProfit)
}
