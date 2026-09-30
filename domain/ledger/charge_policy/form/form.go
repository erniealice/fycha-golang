// Package form holds the drawer-form view models of the charge policy pages.
package form

import (
	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
)

// PolicyData is the add/edit drawer of a policy header (code is immutable on edit).
type PolicyData struct {
	FormAction   string
	WorkspaceID  string // injected by the ViewAdapter for action_workspace_guard
	IsEdit       bool
	ID           string
	Code         string
	Name         string
	Description  string
	Labels       cp.Labels
	CommonLabels any
}

// VersionData is the classification drawer of a draft version.
type VersionData struct {
	FormAction          string
	WorkspaceID         string
	VersionID           string
	AssessmentScope     string
	AssessmentNote      string
	AccountingRoles     []cp.Option
	TaxPositions        []cp.Option
	BookPresentations   []cp.Option
	TaxTreatmentID      string
	TaxTreatmentOptions []cp.Option // empty when the workspace has no picker source
	Labels              cp.Labels
	CommonLabels        any
}

// ComponentData is the add/edit drawer of a component (document rule).
type ComponentData struct {
	FormAction        string
	WorkspaceID       string
	IsEdit            bool
	ID                string
	Roles             []cp.Option
	DocumentKinds     []cp.Option
	BookPresentations []cp.Option
	Labels            cp.Labels
	CommonLabels      any
}

// PostingData is the add/edit drawer of a posting (account mapping).
type PostingData struct {
	FormAction   string
	WorkspaceID  string
	IsEdit       bool
	ID           string
	Events       []cp.Option
	Roles        []cp.Option
	Accounts     []cp.Option
	Labels       cp.Labels
	CommonLabels any
}

// NewVersionData is the confirmation drawer of "New draft version".
type NewVersionData struct {
	FormAction   string
	WorkspaceID  string
	Message      string
	Labels       cp.Labels
	CommonLabels any
}
