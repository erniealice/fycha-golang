package form

import (
	ledger "github.com/erniealice/fycha-golang/domain/ledger"
)

// Data is the template data for the journal entry drawer form.
type Data struct {
	FormAction   string
	WorkspaceID   string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit       bool
	ID           string
	Date         string
	Description  string
	Notes        string
	Lines        []Line
	Labels       ledger.JournalFormLabels
	CommonLabels any
}

// Line represents one editable journal line in the form.
// The account selector submits account_id[N] (hidden) alongside debit[N], credit[N], memo[N].
// Consumer apps wire account lookup to populate AccountCode and AccountName for display.
type Line struct {
	Index       int    // 1-based line number for display
	AccountID   string // selected account ID (stored as hidden input)
	AccountCode string // display code (e.g. "1110")
	AccountName string // display name (e.g. "Cash on Hand")
	Debit       string
	Credit      string
	Memo        string
}

// ParsedLine holds one parsed journal line from the form submission.
// Consumer apps create JournalLine protos from these after creating the JournalEntry.
type ParsedLine struct {
	AccountID string
	Debit     float64
	Credit    float64
	Memo      string
	Order     int32
}
