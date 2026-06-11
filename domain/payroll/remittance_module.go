package payroll

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	payrollremittancepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_remittance"
	remittance "github.com/erniealice/fycha-golang/domain/payroll/remittance"
	remittancelist "github.com/erniealice/fycha-golang/domain/payroll/remittance/list"
)

// RemittanceModuleDeps holds all dependencies for the payroll remittance module.
type RemittanceModuleDeps struct {
	Routes       remittance.Routes
	Labels       remittance.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// PayrollRemittance use cases
	CreatePayrollRemittance          func(ctx context.Context, req *payrollremittancepb.CreatePayrollRemittanceRequest) (*payrollremittancepb.CreatePayrollRemittanceResponse, error)
	ListPayrollRemittances           func(ctx context.Context, req *payrollremittancepb.ListPayrollRemittancesRequest) (*payrollremittancepb.ListPayrollRemittancesResponse, error)
	GetPayrollRemittanceListPageData func(ctx context.Context, req *payrollremittancepb.GetPayrollRemittanceListPageDataRequest) (*payrollremittancepb.GetPayrollRemittanceListPageDataResponse, error)
}

// RemittanceModule holds all constructed payroll remittance views.
type RemittanceModule struct {
	routes remittance.Routes
	List   view.View
}

// NewRemittanceModule creates a payroll remittance module.
func NewRemittanceModule(deps *RemittanceModuleDeps) *RemittanceModule {
	if deps == nil {
		deps = &RemittanceModuleDeps{}
	}

	listDeps := &remittancelist.Deps{
		Routes:                           deps.Routes,
		Labels:                           deps.Labels,
		CommonLabels:                     deps.CommonLabels,
		TableLabels:                      deps.TableLabels,
		GetPayrollRemittanceListPageData: deps.GetPayrollRemittanceListPageData,
	}

	return &RemittanceModule{
		routes: deps.Routes,
		List:   remittancelist.NewView(listDeps),
	}
}

// RegisterRoutes registers all payroll remittance routes with the given route registrar.
func (m *RemittanceModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.List != nil && m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
}
