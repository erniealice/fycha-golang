// Package shared holds view-typed data shapes used by both the list and detail
// sub-packages of the depreciation-run view module.
// Having a dedicated leaf package (no intra-module imports) breaks the import
// cycle that would otherwise form if list/page.go or detail/page.go imported
// from the parent views/depreciation_run package.
package shared

// DepreciationRunRow is the view-layer representation of a single depreciation run.
type DepreciationRunRow struct {
	ID            string
	WorkspaceID   string
	ScopeKind     string // "asset" | "category" | "policy" | "workspace"
	ScopeID       string // asset_id | category_id | policy_id — empty for WORKSPACE
	ScopeLabel    string // human-readable scope display name (resolved by block.go)
	AsOfDate      string // YYYY-MM-DD
	InitiatorID   string // FK to entity.User
	InitiatorName string // resolved display name (resolved by block.go)
	InitiatedAt   string // RFC3339 or ""
	CompletedAt   string // RFC3339 or ""
	Status        string // "pending" | "complete" | "failed"
	CreatedCount  int32
	SkippedCount  int32
	ErroredCount  int32
	ErrorSummary  string
	Notes         string
	// IsStalePending is true when status=pending AND now()-initiated_at > stale threshold.
	// Computed by block.go shim using DEPRECIATION_RUN_PENDING_STALE_MINUTES env (default 5).
	IsStalePending bool
}

// DepreciationRunWithEntries bundles a run and its schedule entry list for the detail page.
type DepreciationRunWithEntries struct {
	Run     DepreciationRunRow
	Entries []DepreciationRunEntryRow
}

// DepreciationRunEntryRow is the view-layer representation of a single
// DepreciationSchedule row scoped to a run (used by selections, results, and
// transactions tabs).
type DepreciationRunEntryRow struct {
	ID                  string
	RunID               string
	AssetID             string
	AssetName           string
	PeriodStartDate     string // YYYY-MM-DD
	DepreciationAmount  int64  // centavos
	Currency            string
	Outcome             string // "created" | "skipped" | "errored"
	ErrorMessage        string
	AssetTransactionID  string // populated when outcome=created
	IsPosted            bool
}

// AssetTransactionRow is the view-layer representation of an asset_transaction
// row scoped to a depreciation run (used by the Transactions tab).
type AssetTransactionRow struct {
	ID              string
	AssetID         string
	AssetName       string
	TransactionType string // e.g. "ASSET_TRANSACTION_TYPE_DEPRECIATION"
	TransactionDate string // YYYY-MM-DD
	Amount          int64  // centavos
	Currency        string
	PeriodStartDate string // YYYY-MM-DD — set when type=DEPRECIATION
	RunID           string
}

// ListDepreciationRunsScope carries filter parameters for the list page.
type ListDepreciationRunsScope struct {
	WorkspaceID string
	Status      string // "" = all
	Limit       int32
}
