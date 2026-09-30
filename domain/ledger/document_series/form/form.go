// Package form holds the drawer-form view model of the document series pages.
package form

import (
	ds "github.com/erniealice/fycha-golang/domain/ledger/document_series"
)

// Option is one <select> option.
type Option struct {
	Value    string
	Label    string
	Selected bool
	// Description feeds the pyeza select's per-option hint (optional).
	Description string
}

// SeriesData is the add/edit drawer. Code, document kind, fiscal reset and the
// next number are fixed once the series exists (edit shows them read-only).
type SeriesData struct {
	FormAction    string
	WorkspaceID   string // injected by the ViewAdapter for action_workspace_guard
	IsEdit        bool
	ID            string
	Code          string
	Name          string
	IssuerName    string
	IssuerTaxID   string
	DocumentKind  string
	KindLabel     string // edit: read-only display of the kind
	KindOptions   []Option
	Prefix        string
	BranchCode    string
	FiscalReset   string // display label (only NEVER is implemented)
	NextNumber    string
	NumberPadding string
	NextPreview   string // edit: formatted "Next document: …"
	Locked        bool   // edit: numbers were already issued
	Labels        ds.Labels
	CommonLabels  any
}
