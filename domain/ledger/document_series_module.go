package ledger

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	documentseries "github.com/erniealice/fycha-golang/domain/ledger/document_series"
	dsaction "github.com/erniealice/fycha-golang/domain/ledger/document_series/action"
	dslist "github.com/erniealice/fycha-golang/domain/ledger/document_series/list"
)

// DocumentSeriesModuleDeps holds all dependencies for the document series
// settings module. UseCases is nil when the espyna document series use cases
// are not wired; the module then fails closed (routes answer 503, never fake data).
type DocumentSeriesModuleDeps struct {
	Routes       documentseries.Routes
	Labels       documentseries.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *documentseries.UseCases
}

// DocumentSeriesModule holds all constructed document series views.
type DocumentSeriesModule struct {
	routes documentseries.Routes

	List, Table       view.View
	Add, Edit, Retire view.View
}

// NewDocumentSeriesModule wires every document series view.
func NewDocumentSeriesModule(deps *DocumentSeriesModuleDeps) *DocumentSeriesModule {
	if deps == nil {
		deps = &DocumentSeriesModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &documentseries.UseCases{} // fail closed: every feature nil
	}
	listDeps := &dslist.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs}
	actionDeps := &dsaction.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, UseCases: ucs}
	return &DocumentSeriesModule{
		routes: deps.Routes,
		List:   dslist.NewView(listDeps),
		Table:  dslist.NewTableView(listDeps),
		Add:    dsaction.NewAddAction(actionDeps),
		Edit:   dsaction.NewEditAction(actionDeps),
		Retire: dsaction.NewRetireAction(actionDeps),
	}
}

// RegisterRoutes registers every document series route.
func (m *DocumentSeriesModule) RegisterRoutes(r view.RouteRegistrar) {
	rt := m.routes
	r.GET(rt.ListURL, m.List)
	r.GET(rt.TableURL, m.Table)
	r.GET(rt.AddURL, m.Add)
	r.POST(rt.AddURL, m.Add)
	r.GET(rt.EditURL, m.Edit)
	r.POST(rt.EditURL, m.Edit)
	r.POST(rt.RetireURL, m.Retire)
}
