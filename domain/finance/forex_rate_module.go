package finance

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	forexratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/finance/forex_rate"
	forexrate "github.com/erniealice/fycha-golang/domain/finance/forex_rate"
	listview "github.com/erniealice/fycha-golang/domain/finance/forex_rate/list"
)

// ForexRateModuleDeps holds all dependencies for the forex_rate module.
type ForexRateModuleDeps struct {
	Routes       forexrate.Routes
	Labels       forexrate.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Forex rate use cases (read-only)
	ListForexRates func(ctx context.Context, req *forexratepb.ListForexRatesRequest) (*forexratepb.ListForexRatesResponse, error)
}

// ForexRateModule holds all constructed forex_rate views.
type ForexRateModule struct {
	List   view.View
	routes forexrate.Routes
}

// NewForexRateModule creates a forex_rate module with the List view wired.
func NewForexRateModule(deps *ForexRateModuleDeps) *ForexRateModule {
	if deps == nil {
		deps = &ForexRateModuleDeps{}
	}

	listDeps := &listview.Deps{
		Routes:         deps.Routes,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
		TableLabels:    deps.TableLabels,
		ListForexRates: deps.ListForexRates,
	}

	return &ForexRateModule{
		List:   listview.NewView(listDeps),
		routes: deps.Routes,
	}
}

// RegisterRoutes registers all forex_rate routes with the given route registrar.
func (m *ForexRateModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.List != nil && m.routes.ListURL != "" {
		r.GET(m.routes.ListURL, m.List)
	}
	// Detail page is a coming-soon placeholder for now
	if m.routes.DetailURL != "" {
		r.GET(m.routes.DetailURL, comingSoonView("Forex Rate Detail", "forex_rate", "forex-rate-detail"))
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
