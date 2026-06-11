package tax

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	taxratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_rate"
	taxrate "github.com/erniealice/fycha-golang/domain/tax/tax_rate"
	listview "github.com/erniealice/fycha-golang/domain/tax/tax_rate/list"
)

// TaxRateModuleDeps holds all dependencies for the tax_rate module.
type TaxRateModuleDeps struct {
	Routes       taxrate.Routes
	Labels       taxrate.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Tax rate use cases (read-only)
	ListTaxRates func(ctx context.Context, req *taxratepb.ListTaxRatesRequest) (*taxratepb.ListTaxRatesResponse, error)
}

// TaxRateModule holds all constructed tax_rate views.
type TaxRateModule struct {
	List   view.View
	routes taxrate.Routes
}

// NewTaxRateModule creates a tax_rate module with the List view wired.
func NewTaxRateModule(deps *TaxRateModuleDeps) *TaxRateModule {
	if deps == nil {
		deps = &TaxRateModuleDeps{}
	}

	listDeps := &listview.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
		ListTaxRates: deps.ListTaxRates,
	}

	return &TaxRateModule{
		List:   listview.NewView(listDeps),
		routes: deps.Routes,
	}
}

// RegisterRoutes registers all tax_rate routes with the given route registrar.
func (m *TaxRateModule) RegisterRoutes(r view.RouteRegistrar) {
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
