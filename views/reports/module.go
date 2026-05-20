package reports

import (
	"context"
	"log"
	"net/http"
	"time"

	apagingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ap_aging"
	aragingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ar_aging"
	dspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/domain_specific"
	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	collectionsummaryreport "github.com/erniealice/fycha-golang/views/reports/collection_summary_report"
	costsales "github.com/erniealice/fycha-golang/views/reports/cost_of_sales"
	dashboardview "github.com/erniealice/fycha-golang/views/reports/dashboard"
	disbursementreport "github.com/erniealice/fycha-golang/views/reports/disbursement_report"
	expenditurereport "github.com/erniealice/fycha-golang/views/reports/expenditure_report"
	expensesview "github.com/erniealice/fycha-golang/views/reports/expenses"
	grossprofit "github.com/erniealice/fycha-golang/views/reports/gross_profit"
	netprofit "github.com/erniealice/fycha-golang/views/reports/net_profit"
	payablesagingreport "github.com/erniealice/fycha-golang/views/reports/payables_aging_report"
	receivablesagingreport "github.com/erniealice/fycha-golang/views/reports/receivables_aging_report"
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
//
// 20260520-21 Wave B P1.E.1-P1.E.5 — every legacy `fycha.DataSource` method
// for the 13 in-scope methods has been replaced with typed closures over
// the new `service.reporting.v1` proto packages. The `DB` field has been
// removed entirely; downstream view consumers receive the typed closures
// they actually need rather than asserting through a duck interface.
// Closures may be nil on mock builds — every view handles nil by
// returning empty responses.
type ModuleDeps struct {
	Routes       fycha.ReportsRoutes
	Labels       fycha.ReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Wave B P1.E.1 — service-driven AR aging closures.
	GetReceivablesAgingReport  func(context.Context, *aragingpb.GetReceivablesAgingRequest) (*aragingpb.GetReceivablesAgingResponse, error)
	GetCollectionSummaryReport func(context.Context, *aragingpb.GetCollectionSummaryRequest) (*aragingpb.GetCollectionSummaryResponse, error)
	// Wave B P1.E.2 — service-driven AP aging closures.
	GetPayablesAgingReport func(context.Context, *apagingpb.GetPayablesAgingRequest) (*apagingpb.GetPayablesAgingResponse, error)
	// Wave B P1.E.3 — service-driven gross/cashflow closures (cash book
	// gets wired separately in fycha block.go's Cash module).
	GetGrossProfitReport func(context.Context, *gcfpb.GetGrossProfitRequest) (*gcfpb.GetGrossProfitResponse, error)
	// Wave B P1.E.5 — service-driven domain-specific closures.
	GetRevenueReport      func(context.Context, *dspb.GetRevenueReportRequest) (*dspb.GetRevenueReportResponse, error)
	GetExpenditureReport  func(context.Context, *dspb.GetExpenditureReportRequest) (*dspb.GetExpenditureReportResponse, error)
	GetDisbursementReport func(context.Context, *dspb.GetDisbursementReportRequest) (*dspb.GetDisbursementReportResponse, error)
	ListRevenue           func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
	ListExpenses          func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
}

// Module holds all constructed report views.
type Module struct {
	routes                        fycha.ReportsRoutes
	Dashboard                     view.View
	Revenue                       view.View
	CostOfSales                   view.View
	GrossProfit                   view.View
	Expenses                      view.View
	NetProfit                     view.View
	RevenueReport                 view.View
	RevenueReportExport           http.HandlerFunc
	ExpenditureReport             view.View
	ExpenditureReportExport       http.HandlerFunc
	DisbursementReport            view.View
	DisbursementReportExport      http.HandlerFunc
	ReceivablesAgingReport        view.View
	ReceivablesAgingReportExport  http.HandlerFunc
	PayablesAgingReport           view.View
	PayablesAgingReportExport     http.HandlerFunc
	CollectionSummaryReport       view.View
	CollectionSummaryReportExport http.HandlerFunc
}

func NewModule(deps *ModuleDeps) *Module {
	grossProfitDeps := &grossprofit.Deps{
		GetGrossProfitReport: deps.GetGrossProfitReport,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
	}
	revenueReportDeps := &revenuereport.Deps{
		GetRevenueReport: deps.GetRevenueReport,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		Routes:           deps.Routes,
	}
	expenditureReportDeps := &expenditurereport.Deps{
		GetExpenditureReport: deps.GetExpenditureReport,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		Routes:               deps.Routes,
	}
	disbursementReportDeps := &disbursementreport.Deps{
		GetDisbursementReport: deps.GetDisbursementReport,
		Labels:                deps.Labels,
		CommonLabels:          deps.CommonLabels,
		TableLabels:           deps.TableLabels,
		Routes:                deps.Routes,
	}
	payablesAgingReportDeps := &payablesagingreport.Deps{
		GetPayablesAgingReport: deps.GetPayablesAgingReport,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		Routes:                 deps.Routes,
	}
	return &Module{
		routes: deps.Routes,
		Dashboard: dashboardview.NewView(&dashboardview.Deps{
			Routes:               deps.Routes,
			GetGrossProfitReport: deps.GetGrossProfitReport,
			ListExpenses:         deps.ListExpenses,
			Labels:               deps.Labels,
			CommonLabels:         deps.CommonLabels,
		}),
		Revenue: revenue.NewView(&revenue.Deps{
			ListRevenue:  deps.ListRevenue,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
		}),
		CostOfSales: costsales.NewView(&costsales.Deps{
			GetGrossProfitReport: deps.GetGrossProfitReport,
			Labels:               deps.Labels,
			CommonLabels:         deps.CommonLabels,
			TableLabels:          deps.TableLabels,
		}),
		GrossProfit: grossprofit.NewView(grossProfitDeps),
		Expenses: expensesview.NewView(&expensesview.Deps{
			ListExpenses: deps.ListExpenses,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
		}),
		NetProfit: netprofit.NewView(&netprofit.Deps{
			GetGrossProfitReport: deps.GetGrossProfitReport,
			ListExpenses:         deps.ListExpenses,
			Labels:               deps.Labels,
			CommonLabels:         deps.CommonLabels,
			TableLabels:          deps.TableLabels,
		}),
		RevenueReport:                 revenuereport.NewView(revenueReportDeps),
		RevenueReportExport:           revenuereport.NewExportHandler(revenueReportDeps),
		ExpenditureReport:             expenditurereport.NewView(expenditureReportDeps),
		ExpenditureReportExport:       expenditurereport.NewExportHandler(expenditureReportDeps),
		DisbursementReport:            disbursementreport.NewView(disbursementReportDeps),
		DisbursementReportExport:      disbursementreport.NewExportHandler(disbursementReportDeps),
		ReceivablesAgingReport: receivablesagingreport.NewView(&receivablesagingreport.Deps{
			Labels:                    deps.Labels,
			CommonLabels:              deps.CommonLabels,
			TableLabels:               deps.TableLabels,
			Routes:                    deps.Routes,
			GetReceivablesAgingReport: deps.GetReceivablesAgingReport,
		}),
		ReceivablesAgingReportExport: receivablesagingreport.NewExportHandler(&receivablesagingreport.Deps{
			Labels:                    deps.Labels,
			CommonLabels:              deps.CommonLabels,
			TableLabels:               deps.TableLabels,
			Routes:                    deps.Routes,
			GetReceivablesAgingReport: deps.GetReceivablesAgingReport,
		}),
		PayablesAgingReport:       payablesagingreport.NewView(payablesAgingReportDeps),
		PayablesAgingReportExport: payablesagingreport.NewExportHandler(payablesAgingReportDeps),
		CollectionSummaryReport: collectionsummaryreport.NewView(&collectionsummaryreport.Deps{
			Labels:                     deps.Labels,
			CommonLabels:               deps.CommonLabels,
			TableLabels:                deps.TableLabels,
			Routes:                     deps.Routes,
			GetCollectionSummaryReport: deps.GetCollectionSummaryReport,
		}),
		CollectionSummaryReportExport: collectionsummaryreport.NewExportHandler(&collectionsummaryreport.Deps{
			Labels:                     deps.Labels,
			CommonLabels:               deps.CommonLabels,
			TableLabels:                deps.TableLabels,
			Routes:                     deps.Routes,
			GetCollectionSummaryReport: deps.GetCollectionSummaryReport,
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
	r.GET(m.routes.ExpenditureReportURL, m.ExpenditureReport)
	handleFunc(r, "GET", m.routes.ExpenditureReportExportURL, m.ExpenditureReportExport)
	r.GET(m.routes.DisbursementReportURL, m.DisbursementReport)
	handleFunc(r, "GET", m.routes.DisbursementReportExportURL, m.DisbursementReportExport)
	r.GET(m.routes.ReceivablesAgingReportURL, m.ReceivablesAgingReport)
	handleFunc(r, "GET", m.routes.ReceivablesAgingReportExportURL, m.ReceivablesAgingReportExport)
	r.GET(m.routes.PayablesAgingReportURL, m.PayablesAgingReport)
	handleFunc(r, "GET", m.routes.PayablesAgingReportExportURL, m.PayablesAgingReportExport)
	r.GET(m.routes.CollectionSummaryReportURL, m.CollectionSummaryReport)
	handleFunc(r, "GET", m.routes.CollectionSummaryReportExportURL, m.CollectionSummaryReportExport)
}
