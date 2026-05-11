// Package depreciationrun wires the Depreciation Run view module:
// Surface D — run history list + detail pages.
package depreciationrun

import (
	"context"

	fycha "github.com/erniealice/fycha-golang"
	depreciationrundetail "github.com/erniealice/fycha-golang/views/depreciation_run/detail"
	depreciationrunlist "github.com/erniealice/fycha-golang/views/depreciation_run/list"
	drshared "github.com/erniealice/fycha-golang/views/depreciation_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ---------------------------------------------------------------------------
// Re-export shared view-typed data shapes so block.go callers can reference
// them via the top-level depreciationrun package.
// ---------------------------------------------------------------------------

// DepreciationRunRow is the view-layer representation of a single depreciation run.
type DepreciationRunRow = drshared.DepreciationRunRow

// DepreciationRunWithEntries bundles a run and its entry list for the detail page.
type DepreciationRunWithEntries = drshared.DepreciationRunWithEntries

// DepreciationRunEntryRow is the view-layer representation of a single schedule entry.
type DepreciationRunEntryRow = drshared.DepreciationRunEntryRow

// AssetTransactionRow is the view-layer representation of an asset_transaction row.
type AssetTransactionRow = drshared.AssetTransactionRow

// ListDepreciationRunsScope carries filter parameters for the list page.
type ListDepreciationRunsScope = drshared.ListDepreciationRunsScope

// ---------------------------------------------------------------------------
// ModuleDeps — typed callbacks; no espyna/proto types cross this boundary.
// ---------------------------------------------------------------------------

// ModuleDeps holds all dependencies for the depreciation-run view module.
type ModuleDeps struct {
	Routes       fycha.DepreciationRunRoutes
	Labels       fycha.DepreciationRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Surface D callbacks — run history list + detail.

	// ListDepreciationRuns returns a page of run rows matching the given scope.
	ListDepreciationRuns func(ctx context.Context, scope ListDepreciationRunsScope) ([]DepreciationRunRow, string, error)

	// ReadDepreciationRun fetches a single run plus all its schedule entries by run ID.
	ReadDepreciationRun func(ctx context.Context, id string) (*DepreciationRunWithEntries, error)

	// ListAssetTransactionsByRunID fetches asset_transaction rows scoped to a run.
	// Used to populate the Transactions tab on the detail page.
	ListAssetTransactionsByRunID func(ctx context.Context, runID string) ([]AssetTransactionRow, error)
}

// ---------------------------------------------------------------------------
// Module — holds constructed view instances.
// ---------------------------------------------------------------------------

// Module holds all constructed depreciation-run views.
type Module struct {
	routes fycha.DepreciationRunRoutes
	// Surface D.
	List    view.View
	Table   view.View
	Detail  view.View
	TabView view.View
}

// NewModule constructs the depreciation-run module from the given deps.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &depreciationrunlist.ListViewDeps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		TableLabels:          deps.TableLabels,
		ListDepreciationRuns: deps.ListDepreciationRuns,
	}
	detailDeps := &depreciationrundetail.DetailViewDeps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		ReadDepreciationRun:          deps.ReadDepreciationRun,
		ListAssetTransactionsByRunID: deps.ListAssetTransactionsByRunID,
	}

	return &Module{
		routes:  deps.Routes,
		List:    depreciationrunlist.NewView(listDeps),
		Table:   depreciationrunlist.NewTableView(listDeps),
		Detail:  depreciationrundetail.NewView(detailDeps),
		TabView: depreciationrundetail.NewTabView(detailDeps),
	}
}

// RegisterRoutes registers all depreciation-run routes on the given registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Surface D — run history list + detail.
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.ListTableURL, m.Table)
	r.POST(m.routes.ListTableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.DetailTabActionURL, m.TabView)
}
