package report

// AgingToolbarPrefixData holds data for the report-aging-toolbar-prefix template.
type AgingToolbarPrefixData struct {
	FilterSheetURL    string
	ActiveFilterCount int
	AsOfDate          string
	GroupByValue      string
}

// DimensionToolbarPrefixData holds data for the report-dimension-toolbar-prefix template.
type DimensionToolbarPrefixData struct {
	FilterSheetURL    string
	ActiveFilterCount int
	PrimaryLabel      string
	PrimaryValue      string
	RowsLabel         string
	RowsValue         string
}

// RecoveryToolbarChip is one active-filter chip of the recovery report toolbar.
type RecoveryToolbarChip struct {
	Label  string
	Value  string
	TestID string
}

// RecoveryToolbarPrefixData holds data for the report-recovery-toolbar-prefix
// template (labels come from Lyngua, never hardcoded in the template).
type RecoveryToolbarPrefixData struct {
	FilterSheetURL    string
	FiltersLabel      string
	ActiveFilterCount int
	Chips             []RecoveryToolbarChip
}

// RecoveryFilterField is one input of the recovery report filter sheet.
type RecoveryFilterField struct {
	Type   string // date | text
	Name   string
	ID     string
	Label  string
	Value  string
	Hint   string
	TestID string
}

// RecoveryFilterSheetData holds data for the recovery-report-filter-sheet template.
type RecoveryFilterSheetData struct {
	FormID     string
	ReportURL  string
	Fields     []RecoveryFilterField
	ApplyLabel string
	ClearLabel string
}
