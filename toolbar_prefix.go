package fycha

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
