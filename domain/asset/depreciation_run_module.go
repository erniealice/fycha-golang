package asset

import (
	"context"

	depreciationrun "github.com/erniealice/fycha-golang/domain/asset/depreciation_run"
	depreciationrundetail "github.com/erniealice/fycha-golang/domain/asset/depreciation_run/detail"
	depreciationrunlist "github.com/erniealice/fycha-golang/domain/asset/depreciation_run/list"
	drshared "github.com/erniealice/fycha-golang/domain/asset/depreciation_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ---------------------------------------------------------------------------
// Re-export shared view-typed data shapes so block.go callers can reference
// them via the top-level asset domain package.
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
// DepreciationRunModuleDeps — typed callbacks; no espyna/proto types cross this boundary.
// ---------------------------------------------------------------------------

// DepreciationRunModuleDeps holds all dependencies for the depreciation-run view module.
type DepreciationRunModuleDeps struct {
	Routes       depreciationrun.Routes
	Labels       depreciationrun.Labels
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
// DepreciationRunModule — holds constructed view instances.
// ---------------------------------------------------------------------------

// DepreciationRunModule holds all constructed depreciation-run views.
type DepreciationRunModule struct {
	routes depreciationrun.Routes
	// Surface D.
	List    view.View
	Table   view.View
	Detail  view.View
	TabView view.View
}

// NewDepreciationRunModule constructs the depreciation-run module from the given deps.
func NewDepreciationRunModule(deps *DepreciationRunModuleDeps) *DepreciationRunModule {
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

	return &DepreciationRunModule{
		routes:  deps.Routes,
		List:    depreciationrunlist.NewView(listDeps),
		Table:   depreciationrunlist.NewTableView(listDeps),
		Detail:  depreciationrundetail.NewView(detailDeps),
		TabView: depreciationrundetail.NewTabView(detailDeps),
	}
}

// RegisterRoutes registers all depreciation-run routes on the given registrar.
func (m *DepreciationRunModule) RegisterRoutes(r view.RouteRegistrar) {
	// Surface D — run history list + detail.
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.ListTableURL, m.Table)
	r.POST(m.routes.ListTableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.DetailTabActionURL, m.TabView)
}
