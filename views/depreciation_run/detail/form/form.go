// Package form holds the data types shared between the depreciation-run detail
// page view and its templates.
package form

import (
	fycha "github.com/erniealice/fycha-golang"
	drshared "github.com/erniealice/fycha-golang/views/depreciation_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// PageData is the full data context passed to the depreciation-run-detail template.
type PageData struct {
	types.PageData
	ContentTemplate string

	// Run holds the view-typed run row.
	Run drshared.DepreciationRunRow

	// Entries holds all schedule entry rows for the run.
	Entries []drshared.DepreciationRunEntryRow

	// IsPossiblyInterrupted is true when Status=pending AND initiated_at is stale.
	IsPossiblyInterrupted bool

	// ActiveTab is the currently active tab key.
	ActiveTab string

	// TabItems is the slice of tab buttons rendered by {{template "tabs" ...}}.
	TabItems []pyeza.TabItem

	// Labels is the depreciation-run label bundle.
	Labels fycha.DepreciationRunLabels

	// SelectionsTable is the TableConfig for the Selections tab.
	SelectionsTable *types.TableConfig

	// ResultsTable is the TableConfig for the Results tab.
	ResultsTable *types.TableConfig

	// TransactionsTable is the TableConfig for the Transactions tab.
	TransactionsTable *types.TableConfig

	// Transactions holds asset_transaction rows scoped to the run.
	Transactions []drshared.AssetTransactionRow

	// HistoryTable is the TableConfig for the History (audit log) tab.
	// Populated in future iterations; rendered as coming-soon alert for now.
	HistoryTable *types.TableConfig
}
