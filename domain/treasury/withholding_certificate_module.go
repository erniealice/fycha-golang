package treasury

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	withholdingcertificatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	withholdingcertificate "github.com/erniealice/fycha-golang/domain/treasury/withholding_certificate"
	"github.com/erniealice/fycha-golang/domain/treasury/withholding_certificate/action"
	listview "github.com/erniealice/fycha-golang/domain/treasury/withholding_certificate/list"
)

// WithholdingCertificateModuleDeps holds all dependencies for the withholding_certificate module.
type WithholdingCertificateModuleDeps struct {
	Routes       withholdingcertificate.Routes
	Labels       withholdingcertificate.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Withholding certificate use cases
	ListWithholdingCertificates  func(ctx context.Context, req *withholdingcertificatepb.ListWithholdingCertificatesRequest) (*withholdingcertificatepb.ListWithholdingCertificatesResponse, error)
	ReadWithholdingCertificate   func(ctx context.Context, req *withholdingcertificatepb.ReadWithholdingCertificateRequest) (*withholdingcertificatepb.ReadWithholdingCertificateResponse, error)
	CreateWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.CreateWithholdingCertificateRequest) (*withholdingcertificatepb.CreateWithholdingCertificateResponse, error)
	UpdateWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.UpdateWithholdingCertificateRequest) (*withholdingcertificatepb.UpdateWithholdingCertificateResponse, error)
	DeleteWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.DeleteWithholdingCertificateRequest) (*withholdingcertificatepb.DeleteWithholdingCertificateResponse, error)
}

// WithholdingCertificateModule holds all constructed withholding_certificate views.
type WithholdingCertificateModule struct {
	List   view.View
	Add    view.View
	Edit   view.View
	Delete view.View
	routes withholdingcertificate.Routes
}

// NewWithholdingCertificateModule creates a withholding_certificate module with the List view wired.
func NewWithholdingCertificateModule(deps *WithholdingCertificateModuleDeps) *WithholdingCertificateModule {
	if deps == nil {
		deps = &WithholdingCertificateModuleDeps{}
	}

	listDeps := &listview.Deps{
		Routes:                      deps.Routes,
		Labels:                      deps.Labels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		ListWithholdingCertificates: deps.ListWithholdingCertificates,
	}

	actionDeps := &action.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		CreateWithholdingCertificate: deps.CreateWithholdingCertificate,
		ReadWithholdingCertificate:   deps.ReadWithholdingCertificate,
		UpdateWithholdingCertificate: deps.UpdateWithholdingCertificate,
		DeleteWithholdingCertificate: deps.DeleteWithholdingCertificate,
	}

	return &WithholdingCertificateModule{
		List:   listview.NewView(listDeps),
		Add:    action.NewCreateAction(actionDeps),
		Edit:   action.NewEditAction(actionDeps),
		Delete: action.NewDeleteAction(actionDeps),
		routes: deps.Routes,
	}
}

// RegisterRoutes registers all withholding_certificate routes with the given route registrar.
func (m *WithholdingCertificateModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.List != nil && m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
	if m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, comingSoonView("Withholding Certificate Detail", "withholding_certificate", "wht-cert-detail"))
	}
	// Phase 2 H2 — full CRUD handlers
	if m.Add != nil && m.routes.AddURL != "" {
		r.GET(m.routes.AddURL, m.Add)
		r.POST(m.routes.AddURL, m.Add)
	}
	if m.Edit != nil && m.routes.EditURL != "" {
		r.GET(m.routes.EditURL, m.Edit)
		r.POST(m.routes.EditURL, m.Edit)
	}
	if m.Delete != nil && m.routes.DeleteURL != "" {
		r.POST(m.routes.DeleteURL, m.Delete)
	}
}
