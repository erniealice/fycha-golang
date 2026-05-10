// Package forex_rate provides read-only views for the Forex Rate entity.
// Forex rates are read-only in the UI; rows are appended only via RecordOperatorRate.
package forex_rate

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	forexratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/finance/forex_rate"
	fycha "github.com/erniealice/fycha-golang"
	listview "github.com/erniealice/fycha-golang/views/forex_rate/list"
)

// ModuleDeps holds all dependencies for the forex_rate module.
type ModuleDeps struct {
	Routes       fycha.ForexRateRoutes
	Labels       fycha.ForexRateLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Forex rate use cases (read-only)
	ListForexRates func(ctx context.Context, req *forexratepb.ListForexRatesRequest) (*forexratepb.ListForexRatesResponse, error)
}

// Module holds all constructed forex_rate views.
type Module struct {
	List   view.View
	routes fycha.ForexRateRoutes
}

// NewModule creates a forex_rate module with the List view wired.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}

	listDeps := &listview.Deps{
		Routes:         deps.Routes,
		Labels:         deps.Labels,
		CommonLabels:   deps.CommonLabels,
		TableLabels:    deps.TableLabels,
		ListForexRates: deps.ListForexRates,
	}

	return &Module{
		List:   listview.NewView(listDeps),
		routes: deps.Routes,
	}
}

// RegisterRoutes registers all forex_rate routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
