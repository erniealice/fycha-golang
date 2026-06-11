package payroll

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"

	payrollremittancepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_remittance"
	payrollrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_run"

	payrolldashboard "github.com/erniealice/fycha-golang/domain/payroll/payrolldashboard"
)

// PayrollDashboardModuleDeps holds all dependencies for the payroll dashboard module.
// Use case fields are nil until Phase 4-8 payroll use cases are implemented in espyna.
type PayrollDashboardModuleDeps struct {
	// PayrollRun use cases
	CreatePayrollRun          func(ctx context.Context, req *payrollrunpb.CreatePayrollRunRequest) (*payrollrunpb.CreatePayrollRunResponse, error)
	ReadPayrollRun            func(ctx context.Context, req *payrollrunpb.ReadPayrollRunRequest) (*payrollrunpb.ReadPayrollRunResponse, error)
	ListPayrollRuns           func(ctx context.Context, req *payrollrunpb.ListPayrollRunsRequest) (*payrollrunpb.ListPayrollRunsResponse, error)
	GetPayrollRunListPageData func(ctx context.Context, req *payrollrunpb.GetPayrollRunListPageDataRequest) (*payrollrunpb.GetPayrollRunListPageDataResponse, error)

	// PayrollRemittance use cases
	CreatePayrollRemittance func(ctx context.Context, req *payrollremittancepb.CreatePayrollRemittanceRequest) (*payrollremittancepb.CreatePayrollRemittanceResponse, error)
	ListPayrollRemittances  func(ctx context.Context, req *payrollremittancepb.ListPayrollRemittancesRequest) (*payrollremittancepb.ListPayrollRemittancesResponse, error)

	// Cross-entity route types (re-exported by the payroll facade).
	Routes           PayrollRunRoutes
	RemittanceRoutes PayrollRemittanceRoutes
	SettingsRoutes   PayrollSettingsRoutes
	Labels           PayrollLabels
	CommonLabels     pyeza.CommonLabels

	GetPayrollDashboardPageData func(ctx context.Context, req *payrolldashboard.Request) (*payrolldashboard.Response, error)
}

// PayrollDashboardModule holds all constructed payroll dashboard views.
type PayrollDashboardModule struct {
	deps      *PayrollDashboardModuleDeps
	Dashboard view.View
	routes    PayrollRunRoutes
}

// NewPayrollDashboardModule creates a payroll dashboard module.
func NewPayrollDashboardModule(deps *PayrollDashboardModuleDeps) *PayrollDashboardModule {
	if deps == nil {
		deps = &PayrollDashboardModuleDeps{}
	}

	dashDeps := &payrolldashboard.Deps{
		Routes:               deps.Routes,
		RemittanceRoutes:     deps.RemittanceRoutes,
		SettingsRoutes:       deps.SettingsRoutes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		GetDashboardPageData: deps.GetPayrollDashboardPageData,
	}

	return &PayrollDashboardModule{
		deps:      deps,
		Dashboard: payrolldashboard.NewView(dashDeps),
		routes:    deps.Routes,
	}
}

// RegisterRoutes registers all payroll dashboard routes with the given route registrar.
func (m *PayrollDashboardModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.Dashboard != nil && m.routes.DashboardURL != "" {
		r.GET(m.routes.DashboardURL, m.Dashboard)
	}
}
