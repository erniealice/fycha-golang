package form

import (
	fiscalperiodpkg "github.com/erniealice/fycha-golang/domain/ledger/fiscal_period"
)

// Data is the template data for the fiscal period add drawer form.
type Data struct {
	FormAction   string
	WorkspaceID  string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Labels       fiscalperiodpkg.FormLabels
	CommonLabels any
}
