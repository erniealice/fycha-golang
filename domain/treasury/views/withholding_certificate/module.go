// Package withholding_certificate provides views for the Withholding Certificate entity.
// BIR Form 2307 — creditable withholding tax certificates.
package withholding_certificate

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	withholdingcertificatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	treasury "github.com/erniealice/fycha-golang/domain/treasury"
	"github.com/erniealice/fycha-golang/domain/treasury/views/withholding_certificate/action"
	listview "github.com/erniealice/fycha-golang/domain/treasury/views/withholding_certificate/list"
)

// ModuleDeps holds all dependencies for the withholding_certificate module.
type ModuleDeps struct {
	Routes       treasury.WithholdingCertificateRoutes
	Labels       treasury.WithholdingCertificateLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Withholding certificate use cases
	ListWithholdingCertificates  func(ctx context.Context, req *withholdingcertificatepb.ListWithholdingCertificatesRequest) (*withholdingcertificatepb.ListWithholdingCertificatesResponse, error)
	ReadWithholdingCertificate   func(ctx context.Context, req *withholdingcertificatepb.ReadWithholdingCertificateRequest) (*withholdingcertificatepb.ReadWithholdingCertificateResponse, error)
	CreateWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.CreateWithholdingCertificateRequest) (*withholdingcertificatepb.CreateWithholdingCertificateResponse, error)
	UpdateWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.UpdateWithholdingCertificateRequest) (*withholdingcertificatepb.UpdateWithholdingCertificateResponse, error)
	DeleteWithholdingCertificate func(ctx context.Context, req *withholdingcertificatepb.DeleteWithholdingCertificateRequest) (*withholdingcertificatepb.DeleteWithholdingCertificateResponse, error)
}

// Module holds all constructed withholding_certificate views.
type Module struct {
	List   view.View
	Add    view.View
	Edit   view.View
	Delete view.View
	routes treasury.WithholdingCertificateRoutes
}

// NewModule creates a withholding_certificate module with the List view wired.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
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

	return &Module{
		List:   listview.NewView(listDeps),
		Add:    action.NewCreateAction(actionDeps),
		Edit:   action.NewEditAction(actionDeps),
		Delete: action.NewDeleteAction(actionDeps),
		routes: deps.Routes,
	}
}

// RegisterRoutes registers all withholding_certificate routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
