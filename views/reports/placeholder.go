package reports

import (
	"context"
	"fmt"
	"strings"
	"strconv"
	"unicode"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ReportConfig holds the configuration for a table-based report view.
type ReportConfig struct {
	ActiveNav    string
	ActiveSubNav string
	Title        string
	Subtitle     string
	Icon         string
	TableID      string
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// BuildData fetches columns and rows dynamically (DB query or mock data).
	// Called on every request so data is always fresh.
	BuildData func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error)
	// BuildTotals computes the totals row from the fetched rows (optional).
	// When set, a sticky <tfoot> with bold accounting totals is rendered.
	BuildTotals func(rows []types.TableRow) []types.TableCell
}

// ReportPageData holds the data for a report list page.
type ReportPageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewReportView creates a read-only table report view.
func NewReportView(cfg ReportConfig) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		columns, rows, err := cfg.BuildData(ctx)
		if err != nil {
			return view.Error(fmt.Errorf("report data: %w", err))
		}

		types.ApplyColumnStyles(columns, rows)

		defaultSort := ""
		if len(columns) > 0 {
			defaultSort = columns[0].Key
		}

		tableConfig := &types.TableConfig{
			ID:                   cfg.TableID,
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           false,
			ShowActions:          false,
			ShowFilters:          false,
			ShowSort:             true,
			ShowColumns:          false,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          false,
			DefaultSortColumn:    defaultSort,
			DefaultSortDirection: "asc",
			Labels:               cfg.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   "No data available",
				Message: "Report data will appear here once transactions are recorded.",
			},
		}
		if cfg.BuildTotals != nil {
			tableConfig.TotalsRow = cfg.BuildTotals(rows)
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &ReportPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          cfg.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      cfg.ActiveNav,
				ActiveSubNav:   cfg.ActiveSubNav,
				HeaderTitle:    cfg.Title,
				HeaderSubtitle: cfg.Subtitle,
				HeaderIcon:     cfg.Icon,
				CommonLabels:   cfg.CommonLabels,
			},
			ContentTemplate: "report-list-content",
			Table:           tableConfig,
		}

		return view.OK("report-list", pageData)
	})
}

// ParseCurrency parses a FormatCurrency string (e.g. "₱1,234.56") back to float64.
func ParseCurrency(s string) float64 {
	// Strip currency symbol and any non-numeric chars except '.' and '-'
	clean := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) || r == '.' || r == '-' {
			return r
		}
		return -1
	}, s)
	v, _ := strconv.ParseFloat(clean, 64)
	return v
}

