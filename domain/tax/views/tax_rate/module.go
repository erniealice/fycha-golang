// Package tax_rate provides read-only views for the Tax Rate entity.
// Tax rates are read-only in the UI; supersession is done via admin SQL recipe.
package tax_rate

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	taxratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_rate"
	tax "github.com/erniealice/fycha-golang/domain/tax"
	listview "github.com/erniealice/fycha-golang/domain/tax/views/tax_rate/list"
)

// ModuleDeps holds all dependencies for the tax_rate module.
type ModuleDeps struct {
	Routes       tax.TaxRateRoutes
	Labels       tax.TaxRateLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Tax rate use cases (read-only)
	ListTaxRates func(ctx context.Context, req *taxratepb.ListTaxRatesRequest) (*taxratepb.ListTaxRatesResponse, error)
}

// Module holds all constructed tax_rate views.
type Module struct {
	List   view.View
	routes tax.TaxRateRoutes
}

// NewModule creates a tax_rate module with the List view wired.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}

	listDeps := &listview.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
		ListTaxRates: deps.ListTaxRates,
	}

	return &Module{
		List:   listview.NewView(listDeps),
		routes: deps.Routes,
	}
}

// RegisterRoutes registers all tax_rate routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	if m.List != nil && m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
	// Detail page is a coming-soon placeholder for now
	if m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, comingSoonView("Tax Rate Detail", "tax_rate", "tax-rate-detail"))
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
