package form

import (
	fycha "github.com/erniealice/fycha-golang"
)

// Data is the template data for the fiscal period add drawer form.
type Data struct {
	FormAction   string
	WorkspaceID   string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	Labels       fycha.FiscalPeriodFormLabels
	CommonLabels any
}
