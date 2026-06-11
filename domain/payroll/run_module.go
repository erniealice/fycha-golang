package payroll

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	payrollrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_run"
	run "github.com/erniealice/fycha-golang/domain/payroll/run"
	runlist "github.com/erniealice/fycha-golang/domain/payroll/run/list"
)

// RunModuleDeps holds all dependencies for the payroll run module.
type RunModuleDeps struct {
	Routes       run.Routes
	Labels       run.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// PayrollRun use cases
	CreatePayrollRun          func(ctx context.Context, req *payrollrunpb.CreatePayrollRunRequest) (*payrollrunpb.CreatePayrollRunResponse, error)
	ReadPayrollRun            func(ctx context.Context, req *payrollrunpb.ReadPayrollRunRequest) (*payrollrunpb.ReadPayrollRunResponse, error)
	ListPayrollRuns           func(ctx context.Context, req *payrollrunpb.ListPayrollRunsRequest) (*payrollrunpb.ListPayrollRunsResponse, error)
	GetPayrollRunListPageData func(ctx context.Context, req *payrollrunpb.GetPayrollRunListPageDataRequest) (*payrollrunpb.GetPayrollRunListPageDataResponse, error)
}

// RunModule holds all constructed payroll run views.
type RunModule struct {
	routes run.Routes
	List   view.View
}

// NewRunModule creates a payroll run module.
func NewRunModule(deps *RunModuleDeps) *RunModule {
	if deps == nil {
		deps = &RunModuleDeps{}
	}

	listDeps := &runlist.Deps{
		Routes:                    deps.Routes,
		Labels:                    deps.Labels,
		CommonLabels:              deps.CommonLabels,
		TableLabels:               deps.TableLabels,
		GetPayrollRunListPageData: deps.GetPayrollRunListPageData,
	}

	return &RunModule{
		routes: deps.Routes,
		List:   runlist.NewView(listDeps),
	}
}

// RegisterRoutes registers all payroll run routes with the given route registrar.
func (m *RunModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.List != nil && m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
}
