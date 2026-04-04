package reports

import (
	"log"
	"net/http"

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
	revenuereport "github.com/erniealice/fycha-golang/views/reports/revenue_report"
)

// routeRegistrarFull extends view.RouteRegistrar with HandleFunc support.
// Consumer apps whose RouteRegistrar implements this interface can register raw
// http.HandlerFunc routes. Apps that do not implement HandleFunc will skip those
// routes with a log warning.
type routeRegistrarFull interface {
	view.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// handleFunc is a nil-safe helper that registers an http.HandlerFunc route if the
// RouteRegistrar supports it, otherwise logs a warning and skips.
func handleFunc(r view.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if handler == nil {
		return
	}
	if full, ok := r.(routeRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("fycha/reports: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}

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
	routes              fycha.ReportsRoutes
	Dashboard           view.View
	Revenue             view.View
	CostOfSales         view.View
	GrossProfit         view.View
	Expenses            view.View
	NetProfit           view.View
	RevenueReport       view.View
	RevenueReportExport http.HandlerFunc
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
		RevenueReport: revenuereport.NewView(&revenuereport.Deps{
			DB:           deps.DB,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Routes:       deps.Routes,
		}),
		RevenueReportExport: revenuereport.NewExportHandler(&revenuereport.Deps{
			DB:           deps.DB,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			Routes:       deps.Routes,
		}),
	}
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.RevenueURL, m.Revenue)
	r.GET(m.routes.CostOfSalesURL, m.CostOfSales)
	r.GET(m.routes.GrossProfitURL, m.GrossProfit)
	r.GET(m.routes.ExpensesURL, m.Expenses)
	r.GET(m.routes.NetProfitURL, m.NetProfit)
	r.GET(m.routes.RevenueReportURL, m.RevenueReport)
	handleFunc(r, "GET", m.routes.RevenueReportExportURL, m.RevenueReportExport)
}
