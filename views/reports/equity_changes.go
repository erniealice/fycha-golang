package reports

import (
	"context"
	"fmt"
	"time"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ---------------------------------------------------------------------------
// Data model
// ---------------------------------------------------------------------------

// ECColumn defines a column in the equity changes statement.
type ECColumn struct {
	AccountCode string // e.g. "3010"
	AccountName string // e.g. "Owner's Capital"
	IsTotal     bool   // true for the rightmost "Total" column
}

// ECCell is a single cell value in the equity matrix.
type ECCell struct {
	Value      string // formatted amount or empty string
	IsNegative bool   // true for negative/bracketed values
	IsBold     bool   // true for closing balance row cells
}

// ECRow is a row in the equity changes matrix.
type ECRow struct {
	Label     string   // e.g. "Opening Balance"
	SubLabel  string   // e.g. "Apr 1, 2025" (shown below Label)
	Cells     []ECCell // one cell per column (including Total)
	IsTotal   bool     // closing balance row (bold, double-underline)
	IsSpacer  bool     // blank spacer row for readability
}

// ---------------------------------------------------------------------------
// Deps + PageData
// ---------------------------------------------------------------------------

// EquityChangesDeps holds dependencies for the Equity Changes view.
type EquityChangesDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       fycha.ReportsLabels

	// GetEquityChanges fetches equity changes data for the given period.
	// Phase 8: set to nil — mock data is used automatically.
	GetEquityChanges func(ctx context.Context, startDate, endDate string) ([]ECColumn, []ECRow, error)
}

// EquityChangesPageData is the template data for the equity-changes page.
type EquityChangesPageData struct {
	types.PageData
	ContentTemplate string

	// Period filter state
	ActivePreset  string
	StartDate     string
	EndDate       string
	PeriodLabel   string
	PeriodPresets []fycha.FilterOption

	// KPI summary metrics
	OpeningEquity        string
	ClosingEquity        string
	EquityChangeTrend    string // "+20.3%"
	EquityChangeVariant  string // "success" or "danger"

	// Statement body
	Columns []ECColumn
	Rows    []ECRow
}

// ---------------------------------------------------------------------------
// View constructor
// ---------------------------------------------------------------------------

// NewEquityChangesView creates the Statement of Changes in Equity view.
func NewEquityChangesView(deps *EquityChangesDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		q := viewCtx.QueryParams

		preset := q["period"]
		if preset == "" {
			preset = "thisYear"
		}
		startDate := q["start"]
		endDate := q["end"]

		// Resolve date range from preset
		start, end := fycha.ParsePeriodPreset(preset)
		if preset == "custom" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				start = t
			}
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				end = t
			}
		}

		startDate = start.Format("2006-01-02")
		endDate = end.Format("2006-01-02")
		periodLabel := fmt.Sprintf("%s – %s",
			start.Format("January 2, 2006"),
			end.Format("January 2, 2006"),
		)

		pl := deps.Labels.Period
		periodPresets := fycha.DefaultPeriodPresets(pl, preset)

		// Fetch data
		var columns []ECColumn
		var rows []ECRow
		if deps.GetEquityChanges != nil {
			cols, rws, err := deps.GetEquityChanges(ctx, startDate, endDate)
			if err == nil {
				columns = cols
				rows = rws
			}
		}
		if columns == nil {
			columns, rows = mockECData(start, end)
		}

		// Extract opening and closing equity from rows
		openingEquity := ""
		closingEquity := ""
		for _, row := range rows {
			if len(row.Cells) > 0 {
				totalCell := row.Cells[len(row.Cells)-1]
				if row.Label == "Opening Balance" {
					openingEquity = totalCell.Value
				}
				if row.IsTotal {
					closingEquity = totalCell.Value
				}
			}
		}

		equityChangeVariant := "success"

		pageData := &EquityChangesPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Statement of Changes in Equity",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "reports",
				ActiveSubNav:   "equity-changes",
				HeaderTitle:    "Statement of Changes in Equity",
				HeaderSubtitle: "Equity movements for the period",
				HeaderIcon:     "icon-percent",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:     "equity-changes-content",
			ActivePreset:        preset,
			StartDate:           startDate,
			EndDate:             endDate,
			PeriodLabel:         periodLabel,
			PeriodPresets:       periodPresets,
			OpeningEquity:       openingEquity,
			ClosingEquity:       closingEquity,
			EquityChangeTrend:   "+20.3%",
			EquityChangeVariant: equityChangeVariant,
			Columns:             columns,
			Rows:                rows,
		}

		if viewCtx.IsHTMX {
			return view.OK("equity-changes-content", pageData)
		}
		return view.OK("equity-changes", pageData)
	})
}

// ---------------------------------------------------------------------------
// Mock data (Phase 8)
// ---------------------------------------------------------------------------

// mockECData returns realistic equity changes data.
// Based on plan doc ASCII mockup: partnership with two owners.
func mockECData(start, end time.Time) ([]ECColumn, []ECRow) {
	columns := []ECColumn{
		{AccountCode: "3010", AccountName: "Owner's Capital"},
		{AccountCode: "3020", AccountName: "Owner's Draw"},
		{AccountCode: "3030", AccountName: "Retained Earnings"},
		{AccountCode: "", AccountName: "Total", IsTotal: true},
	}

	openLabel := start.Format("Jan 2, 2006")
	closeLabel := end.Format("Jan 2, 2006")

	rows := []ECRow{
		{
			Label:    "Opening Balance",
			SubLabel: openLabel,
			Cells: []ECCell{
				{Value: "₱800,000.00"},
				{Value: "(₱40,000.00)", IsNegative: true},
				{Value: "(₱180,000.00)", IsNegative: true},
				{Value: "₱580,000.00", IsBold: false},
			},
		},
		{IsSpacer: true},
		{
			Label: "+ Net Income",
			Cells: []ECCell{
				{Value: ""},
				{Value: ""},
				{Value: "₱67,600.00"},
				{Value: "₱67,600.00"},
			},
		},
		{
			Label: "+ Contributions",
			Cells: []ECCell{
				{Value: "₱0.00"},
				{Value: ""},
				{Value: ""},
				{Value: "₱0.00"},
			},
		},
		{
			Label: "- Withdrawals",
			Cells: []ECCell{
				{Value: ""},
				{Value: "(₱10,000.00)", IsNegative: true},
				{Value: ""},
				{Value: "(₱10,000.00)", IsNegative: true},
			},
		},
		{
			Label: "+ Prior Period Adjustments",
			Cells: []ECCell{
				{Value: ""},
				{Value: ""},
				{Value: "₱60,000.00"},
				{Value: "₱60,000.00"},
			},
		},
		{IsSpacer: true},
		{
			Label:    "Closing Balance",
			SubLabel: closeLabel,
			IsTotal:  true,
			Cells: []ECCell{
				{Value: "₱800,000.00", IsBold: true},
				{Value: "(₱50,000.00)", IsNegative: true, IsBold: true},
				{Value: "(₱52,400.00)", IsNegative: true, IsBold: true},
				{Value: "₱697,600.00", IsBold: true},
			},
		},
	}

	return columns, rows
}
