package fycha

import "time"

// ReportFilter holds date range and grouping parameters for report queries.
// Period presets are resolved server-side via ParsePeriodPreset.
type ReportFilter struct {
	StartDate *time.Time
	EndDate   *time.Time
	Period    string // "thisMonth", "lastMonth", "thisQuarter", etc.
	GroupBy   string // "product", "customer", "status", "monthly", "quarterly", "vendor", "category"
}

// SummaryMetric holds a single summary bar metric.
type SummaryMetric struct {
	Label     string
	Value     string
	Highlight bool
	Variant   string // "success", "warning", "danger" — for badge coloring
}

// FilterOption holds a dropdown or button option for filters.
type FilterOption struct {
	Value    string
	Label    string
	Selected bool
}

// FilterState holds the current filter state for template rendering.
type FilterState struct {
	ActivePreset   string
	StartDate      string
	EndDate        string
	GroupBy        string
	GroupByOptions []FilterOption
	PeriodPresets  []FilterOption
}

// PLLineItem holds a single line in a P&L statement.
type PLLineItem struct {
	Label   string
	Value   string
	IsTotal bool
	Variant string
}

// ParsePeriodPreset computes a date range from a named preset.
// Returns start and end times in local timezone.
func ParsePeriodPreset(preset string) (start, end time.Time) {
	now := time.Now()
	year, month, _ := now.Date()
	loc := now.Location()

	switch preset {
	case "lastMonth":
		start = time.Date(year, month-1, 1, 0, 0, 0, 0, loc)
		end = start.AddDate(0, 1, 0).Add(-time.Second)
	case "thisQuarter":
		qMonth := time.Month(((int(month)-1)/3)*3 + 1)
		start = time.Date(year, qMonth, 1, 0, 0, 0, 0, loc)
		end = now
	case "lastQuarter":
		qMonth := time.Month(((int(month)-1)/3)*3 + 1)
		start = time.Date(year, qMonth-3, 1, 0, 0, 0, 0, loc)
		end = time.Date(year, qMonth, 1, 0, 0, 0, 0, loc).Add(-time.Second)
	case "thisYear":
		start = time.Date(year, 1, 1, 0, 0, 0, 0, loc)
		end = now
	case "lastYear":
		start = time.Date(year-1, 1, 1, 0, 0, 0, 0, loc)
		end = time.Date(year, 1, 1, 0, 0, 0, 0, loc).Add(-time.Second)
	default: // "thisMonth" or unknown
		start = time.Date(year, month, 1, 0, 0, 0, 0, loc)
		end = now
	}
	return start, end
}

// FilterSheetData holds the data passed to the report-filter-sheet template.
type FilterSheetData struct {
	Filter       FilterState
	PeriodLabels PeriodLabels
	ReportURL    string
}

// ActiveFilterCount computes how many filters differ from defaults.
// A non-default period preset counts as 1. Custom dates count as 1.
// A non-empty group-by that differs from "product" counts as 1.
func ActiveFilterCount(filter FilterState) int {
	count := 0
	if filter.ActivePreset != "" && filter.ActivePreset != "thisMonth" {
		count++
	}
	if filter.GroupBy != "" && filter.GroupBy != "product" {
		count++
	}
	return count
}

// DefaultPeriodPresets returns the standard period preset options.
func DefaultPeriodPresets(labels PeriodLabels, active string) []FilterOption {
	presets := []struct {
		value string
		label string
	}{
		{"thisMonth", labels.ThisMonth},
		{"lastMonth", labels.LastMonth},
		{"thisQuarter", labels.ThisQuarter},
		{"lastQuarter", labels.LastQuarter},
		{"thisYear", labels.ThisYear},
		{"lastYear", labels.LastYear},
		{"custom", labels.Custom},
	}
	opts := make([]FilterOption, len(presets))
	for i, p := range presets {
		opts[i] = FilterOption{
			Value:    p.value,
			Label:    p.label,
			Selected: p.value == active,
		}
	}
	return opts
}
