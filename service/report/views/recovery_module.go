package reports

import (
	"context"
	"net/http"

	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	report "github.com/erniealice/fycha-golang/service/report"
	costsourcereconciliation "github.com/erniealice/fycha-golang/service/report/views/cost_source_reconciliation"
	recoverablesaging "github.com/erniealice/fycha-golang/service/report/views/recoverables_aging"
)

// RecoveryModuleDeps holds the dependencies of the recovery reports module
// (recoverables aging + cost-source reconciliation). It is separate from
// ModuleDeps so apps that do not opt in keep the default reports surface.
type RecoveryModuleDeps struct {
	Routes       report.RecoveryReportsRoutes
	Labels       report.RecoveryReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	ListRecoverablesAging func(context.Context, *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error)
	ReconcileCostSource   func(context.Context, *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error)
}

// RecoveryModule holds the constructed recovery report views.
type RecoveryModule struct {
	routes report.RecoveryReportsRoutes

	RecoverablesAging              view.View
	RecoverablesAgingExport        http.HandlerFunc
	CostSourceReconciliation       view.View
	CostSourceReconciliationExport http.HandlerFunc
}

// NewRecoveryModule builds the recovery report views.
func NewRecoveryModule(deps *RecoveryModuleDeps) *RecoveryModule {
	agingDeps := &recoverablesaging.Deps{
		Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, Routes: deps.Routes,
		ListRecoverablesAging: deps.ListRecoverablesAging,
	}
	reconDeps := &costsourcereconciliation.Deps{
		Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, Routes: deps.Routes,
		ReconcileCostSource: deps.ReconcileCostSource,
	}
	return &RecoveryModule{
		routes:                         deps.Routes,
		RecoverablesAging:              recoverablesaging.NewView(agingDeps),
		RecoverablesAgingExport:        recoverablesaging.NewExportHandler(agingDeps),
		CostSourceReconciliation:       costsourcereconciliation.NewView(reconDeps),
		CostSourceReconciliationExport: costsourcereconciliation.NewExportHandler(reconDeps),
	}
}

// RegisterRoutes registers the recovery report routes.
func (m *RecoveryModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.RecoverablesAgingURL, m.RecoverablesAging)
	handleFunc(r, "GET", m.routes.RecoverablesAgingExportURL, m.RecoverablesAgingExport)
	r.GET(m.routes.CostSourceReconciliationURL, m.CostSourceReconciliation)
	handleFunc(r, "GET", m.routes.CostSourceReconciliationExportURL, m.CostSourceReconciliationExport)
}
