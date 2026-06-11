package payroll

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	payrollemployee "github.com/erniealice/fycha-golang/domain/payroll/payrollemployee"
	employeelist "github.com/erniealice/fycha-golang/domain/payroll/payrollemployee/list"
)

// PayrollEmployeeModuleDeps holds all dependencies for the payroll employee module.
type PayrollEmployeeModuleDeps struct {
	Routes       payrollemployee.Routes
	Labels       payrollemployee.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PayrollEmployeeModule holds all constructed payroll employee views.
type PayrollEmployeeModule struct {
	routes payrollemployee.Routes
	List   view.View
}

// NewPayrollEmployeeModule creates a payroll employee module.
func NewPayrollEmployeeModule(deps *PayrollEmployeeModuleDeps) *PayrollEmployeeModule {
	if deps == nil {
		deps = &PayrollEmployeeModuleDeps{}
	}

	listDeps := &employeelist.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}

	return &PayrollEmployeeModule{
		routes: deps.Routes,
		List:   employeelist.NewView(listDeps),
	}
}

// RegisterRoutes registers all payroll employee routes with the given route registrar.
func (m *PayrollEmployeeModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.List != nil && m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
}
