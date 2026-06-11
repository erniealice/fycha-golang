package payroll

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	payrollremittancepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_remittance"
	payrollrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_run"
	payroll "github.com/erniealice/fycha-golang/domain/payroll"
	dashboardview "github.com/erniealice/fycha-golang/domain/payroll/views/payroll/dashboard"
)

// ModuleDeps holds all dependencies for the payroll module.
// Use case fields are nil until Phase 4-8 payroll use cases are implemented in espyna.
type ModuleDeps struct {
	// PayrollRun use cases
	CreatePayrollRun          func(ctx context.Context, req *payrollrunpb.CreatePayrollRunRequest) (*payrollrunpb.CreatePayrollRunResponse, error)
	ReadPayrollRun            func(ctx context.Context, req *payrollrunpb.ReadPayrollRunRequest) (*payrollrunpb.ReadPayrollRunResponse, error)
	ListPayrollRuns           func(ctx context.Context, req *payrollrunpb.ListPayrollRunsRequest) (*payrollrunpb.ListPayrollRunsResponse, error)
	GetPayrollRunListPageData func(ctx context.Context, req *payrollrunpb.GetPayrollRunListPageDataRequest) (*payrollrunpb.GetPayrollRunListPageDataResponse, error)

	// PayrollRemittance use cases
	CreatePayrollRemittance func(ctx context.Context, req *payrollremittancepb.CreatePayrollRemittanceRequest) (*payrollremittancepb.CreatePayrollRemittanceResponse, error)
	ListPayrollRemittances  func(ctx context.Context, req *payrollremittancepb.ListPayrollRemittancesRequest) (*payrollremittancepb.ListPayrollRemittancesResponse, error)

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	Routes           payroll.PayrollRunRoutes
	RemittanceRoutes payroll.PayrollRemittanceRoutes
	SettingsRoutes   payroll.PayrollSettingsRoutes
	Labels           payroll.PayrollLabels
	CommonLabels     pyeza.CommonLabels

	GetPayrollDashboardPageData func(ctx context.Context, req *dashboardview.Request) (*dashboardview.Response, error)
}

// Module holds all constructed payroll views.
type Module struct {
	deps *ModuleDeps

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	Dashboard view.View
	routes    payroll.PayrollRunRoutes
}

// NewModule creates a payroll module.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}

	dashDeps := &dashboardview.Deps{
		Routes:               deps.Routes,
		RemittanceRoutes:     deps.RemittanceRoutes,
		SettingsRoutes:       deps.SettingsRoutes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		GetDashboardPageData: deps.GetPayrollDashboardPageData,
	}

	return &Module{
		deps:      deps,
		Dashboard: dashboardview.NewView(dashDeps),
		routes:    deps.Routes,
	}
}

// RegisterRoutes registers all payroll routes with the given route registrar.
// Routes render "Coming Soon" placeholders until view constructors are wired.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	if m.Dashboard != nil && m.routes.DashboardURL != "" {
		r.GET(m.routes.DashboardURL, m.Dashboard)
	}
	r.GET(payroll.PayrollRunListURL, comingSoonView("Payroll Runs", "payroll", "payroll-runs"))
	r.GET(payroll.PayrollRemittanceListURL, comingSoonView("Remittances", "payroll", "remittances"))
	r.GET(payroll.PayrollEmployeeListURL, comingSoonView("Employees", "payroll", "employees"))
	r.GET(payroll.PayrollGovRatesURL, comingSoonView("Gov Contribution Rates", "payroll", "gov-rates"))
	r.GET(payroll.PayrollPayPeriodsURL, comingSoonView("Pay Periods", "payroll", "pay-periods"))
}

// comingSoonView returns a placeholder view that renders a "Coming Soon" page.
func comingSoonView(title, activeNav, activeSubNav string) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		templateName := "coming-soon"
		if viewCtx.IsHTMX {
			templateName = "coming-soon-content"
		}
		return view.OK(templateName, &types.PageData{
			CacheVersion: viewCtx.CacheVersion,
			Title:        title,
			CurrentPath:  viewCtx.CurrentPath,
			ActiveNav:    activeNav,
			ActiveSubNav: activeSubNav,
			HeaderTitle:  title,
		})
	})
}
