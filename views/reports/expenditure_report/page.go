package expenditure_report

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/url"
	"strconv"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	expreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/expenditure_report"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	DB           fycha.DataSource
	Labels       fycha.ReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Routes       fycha.ReportsRoutes
}

// PageData holds the data for the expenditure report page.
type PageData struct {
	types.PageData
	ContentTemplate   string
	Labels            fycha.ExpenditureReportLabels
	Summary           []fycha.SummaryMetric
	Table             *types.TableConfig
	Filter            fycha.FilterState
	PeriodLabels      fycha.PeriodLabels
	ReportURL         string
	FilterSheetURL    string
	ExportURL         string
	ActiveFilterCount int
	PrimaryDimension  string
	RowDimension      string
	PrimaryOptions    []fycha.FilterOption
	RowOptions        []fycha.FilterOption
}

// NewView creates the expenditure report pivot-table view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.ExpenditureReport
		pl := deps.Labels.Period

		// Parse query params
		primary := viewCtx.QueryParams["primary"]
		if primary == "" {
			primary = "monthly"
		}
		rows := viewCtx.QueryParams["rows"]
		if rows == "" {
			rows = "category"
		}
		period := viewCtx.QueryParams["period"]
		if period == "" {
			period = "thisMonth"
		}
		startDateStr := viewCtx.QueryParams["start"]
		endDateStr := viewCtx.QueryParams["end"]

		// Secondary filter IDs
		productID := viewCtx.QueryParams["product-id"]
		locationID := viewCtx.QueryParams["location-id"]
		locationAreaID := viewCtx.QueryParams["location-area-id"]
		expenditureCategoryID := viewCtx.QueryParams["expenditure-category-id"]
		supplierID := viewCtx.QueryParams["supplier-id"]
		expenditureType := viewCtx.QueryParams["expenditure-type"]

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = deps.Routes.ExpenditureReportURL
		}

		// Build dimension options
		primaryOptions := l.DimensionOptions(primary)
		rowOptions := l.DimensionOptions(rows)

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			sheetFilter := fycha.FilterState{
				ActivePreset:  period,
				StartDate:     startDateStr,
				EndDate:       endDateStr,
				PeriodPresets: fycha.DefaultPeriodPresets(pl, period),
			}
			return view.OK("expenditure-report-filter-sheet", &ExpenditureReportFilterSheetData{
				Filter:           sheetFilter,
				PeriodLabels:     pl,
				Labels:           l,
				ReportURL:        reportURL,
				PrimaryDimension: primary,
				RowDimension:     rows,
				PrimaryOptions:   primaryOptions,
				RowOptions:       rowOptions,
			})
		}

		// Build proto request
		req := &expreportpb.ExpenditureReportRequest{
			PrimaryDimension: primary,
			RowDimension:     rows,
		}

		// Resolve dates
		if period == "custom" && startDateStr != "" {
			if _, err := time.Parse("2006-01-02", startDateStr); err == nil {
				req.StartDate = &startDateStr
			}
		}
		if period == "custom" && endDateStr != "" {
			if _, err := time.Parse("2006-01-02", endDateStr); err == nil {
				req.EndDate = &endDateStr
			}
		}

		// Apply period preset if not custom
		if req.StartDate == nil {
			start, _ := fycha.ParsePeriodPreset(period)
			s := start.Format("2006-01-02")
			req.StartDate = &s
		}
		if req.EndDate == nil {
			_, end := fycha.ParsePeriodPreset(period)
			e := end.Format("2006-01-02")
			req.EndDate = &e
		}

		// Apply optional secondary filters
		if productID != "" {
			req.ProductId = &productID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if locationAreaID != "" {
			req.LocationAreaId = &locationAreaID
		}
		if expenditureCategoryID != "" {
			req.ExpenditureCategoryId = &expenditureCategoryID
		}
		if supplierID != "" {
			req.SupplierId = &supplierID
		}
		if expenditureType != "" {
			req.ExpenditureType = &expenditureType
		}

		// Call data source
		resp, err := deps.DB.GetExpenditureReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get expenditure report: %v", err)
			resp = &expreportpb.ExpenditureReportResponse{
				ColumnKeys: []string{},
				Rows:       []*expreportpb.ExpenditureReportRow{},
				Summary:    &expreportpb.ExpenditureReportSummary{},
			}
		}

		// Build summary bar
		summary := buildSummary(resp.GetSummary(), l)

		// Build pivot table
		table := buildPivotTable(resp, l, deps.TableLabels, primary, rows)

		// Build filter sheet URL
		filterSheetURL := buildFilterSheetURL(reportURL, primary, rows, period, startDateStr, endDateStr)

		// Count active filters
		activeCount := 0
		if period != "" && period != "thisMonth" {
			activeCount++
		}
		if primary != "" && primary != "monthly" {
			activeCount++
		}
		if rows != "" && rows != "category" {
			activeCount++
		}

		// Inject filter button + dimension chips into the table toolbar prefix
		table.ToolbarPrefix = buildToolbarPrefix(filterSheetURL, activeCount, "Columns:", primary, "Rows:", rows)

		filter := fycha.FilterState{
			ActivePreset:  period,
			StartDate:     startDateStr,
			EndDate:       endDateStr,
			PeriodPresets: fycha.DefaultPeriodPresets(pl, period),
		}

		// Build export URL with current query params
		exportURL := buildExportURL(deps.Routes.ExpenditureReportExportURL, primary, rows, period, startDateStr, endDateStr)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "expense",
				ActiveSubNav: "expenditure-report",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-bar-chart",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate:   "expenditure-report",
			Labels:            l,
			Summary:           summary,
			Table:             table,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			FilterSheetURL:    filterSheetURL,
			ExportURL:         exportURL,
			ActiveFilterCount: activeCount,
			PrimaryDimension:  primary,
			RowDimension:      rows,
			PrimaryOptions:    primaryOptions,
			RowOptions:        rowOptions,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-expenditure-report"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("expenditure-report-content", pageData)
		}
		return view.OK("expenditure-report", pageData)
	})
}

// ExpenditureReportFilterSheetData holds data for the expenditure report filter sheet template.
type ExpenditureReportFilterSheetData struct {
	Filter           fycha.FilterState
	PeriodLabels     fycha.PeriodLabels
	Labels           fycha.ExpenditureReportLabels
	ReportURL        string
	PrimaryDimension string
	RowDimension     string
	PrimaryOptions   []fycha.FilterOption
	RowOptions       []fycha.FilterOption
}

func buildSummary(s *expreportpb.ExpenditureReportSummary, l fycha.ExpenditureReportLabels) []fycha.SummaryMetric {
	if s == nil {
		s = &expreportpb.ExpenditureReportSummary{}
	}
	grandTotal := float64(s.GetGrandTotal()) / 100.0
	txnCount := s.GetTotalTransactions()
	avgTxn := 0.0
	if txnCount > 0 {
		avgTxn = grandTotal / float64(txnCount)
	}
	return []fycha.SummaryMetric{
		{Label: l.SummaryGrandTotal, Value: formatCurrency(grandTotal), Highlight: true},
		{Label: l.SummaryTransactions, Value: fmt.Sprintf("%d", txnCount)},
		{Label: l.SummaryAverage, Value: formatCurrency(avgTxn)},
	}
}

func buildPivotTable(resp *expreportpb.ExpenditureReportResponse, l fycha.ExpenditureReportLabels, tableLabels types.TableLabels, primary, rowDim string) *types.TableConfig {
	columnKeys := resp.GetColumnKeys()

	// Build dynamic columns
	dynamicColumns := make([]types.TableColumn, 0, len(columnKeys))
	for _, ck := range columnKeys {
		dynamicColumns = append(dynamicColumns, types.TableColumn{
			Key:      ck,
			Label:    ck,
			Sortable: true,
			Align:    "right",
			MinWidth: "7.5rem",
		})
	}

	table := &types.TableConfig{
		ID:              "expenditureReportTable",
		NameColumnLabel: l.PrimaryGroupLabel(rowDim),
		ShowSearch:      false,
		ShowFilters:     false,
		ShowSort:        false,
		ShowColumns:     false,
		ShowExport:      true,
		ShowEntries:     true,
		ShowDensity:     true,
		Labels:          tableLabels,
		ColumnGroups: []types.ColumnGroup{
			{
				Label:   l.PrimaryGroupLabel(primary),
				Columns: dynamicColumns,
			},
			{
				Label: "",
				Columns: []types.TableColumn{
					{Key: "total", Label: l.Total, Sortable: true, Align: "right", MinWidth: "8.125rem"},
				},
			},
		},
		EmptyState: types.TableEmptyState{
			Title:   l.EmptyTitle,
			Message: l.EmptyMessage,
		},
	}

	// Flatten columns for ApplyColumnStyles
	var allColumns []types.TableColumn
	for _, group := range table.ColumnGroups {
		allColumns = append(allColumns, group.Columns...)
	}

	rows := make([]types.TableRow, 0, len(resp.GetRows()))
	for _, row := range resp.GetRows() {
		cellMap := make(map[string]*expreportpb.ExpenditureReportCell, len(row.GetCells()))
		for _, c := range row.GetCells() {
			cellMap[c.GetColumnKey()] = c
		}

		cells := []types.TableCell{
			{Type: "name", Value: row.GetRowKey()},
		}
		dataAttrs := map[string]string{}

		for _, ck := range columnKeys {
			var val int64
			if c, ok := cellMap[ck]; ok {
				val = c.GetTotalExpenditure()
			}
			cells = append(cells, types.TableCell{
				Type:  "text",
				Value: formatCurrency(float64(val) / 100.0),
			})
			dataAttrs[ck] = fmt.Sprintf("%.2f", float64(val)/100.0)
		}

		// Total cell
		cells = append(cells, types.TableCell{
			Type:  "text",
			Value: formatCurrency(float64(row.GetRowTotal()) / 100.0),
		})
		dataAttrs["total"] = fmt.Sprintf("%.2f", float64(row.GetRowTotal())/100.0)

		rows = append(rows, types.TableRow{
			ID:        row.GetRowKey(),
			Cells:     cells,
			DataAttrs: dataAttrs,
		})
	}

	// Add totals row from summary.column_totals
	summary := resp.GetSummary()
	if summary != nil && len(resp.GetRows()) > 0 {
		colTotalMap := make(map[string]*expreportpb.ExpenditureReportCell, len(summary.GetColumnTotals()))
		for _, ct := range summary.GetColumnTotals() {
			colTotalMap[ct.GetColumnKey()] = ct
		}

		totalsCells := []types.TableCell{
			{Type: "name", Value: l.Totals},
		}
		for _, ck := range columnKeys {
			var val int64
			if ct, ok := colTotalMap[ck]; ok {
				val = ct.GetTotalExpenditure()
			}
			totalsCells = append(totalsCells, types.TableCell{
				Type:  "text",
				Value: formatCurrency(float64(val) / 100.0),
			})
		}
		totalsCells = append(totalsCells, types.TableCell{
			Type:  "text",
			Value: formatCurrency(float64(summary.GetGrandTotal()) / 100.0),
		})

		rows = append(rows, types.TableRow{
			ID:    "__totals__",
			Cells: totalsCells,
		})
	}

	table.Rows = rows
	types.ApplyColumnStyles(allColumns, rows)
	types.ApplyTableSettings(table)

	return table
}

func buildExportURL(base, primary, rows, period, start, end string) string {
	params := url.Values{}
	params.Set("primary", primary)
	params.Set("rows", rows)
	params.Set("period", period)
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	return base + "?" + params.Encode()
}

func buildFilterSheetURL(base, primary, rows, period, start, end string) string {
	params := url.Values{}
	params.Set("sheet", "filters")
	params.Set("primary", primary)
	params.Set("rows", rows)
	params.Set("period", period)
	if start != "" {
		params.Set("start", start)
	}
	if end != "" {
		params.Set("end", end)
	}
	return base + "?" + params.Encode()
}

// buildToolbarPrefix builds the filter button + chips HTML for the table toolbar prefix slot.
func buildToolbarPrefix(filterSheetURL string, activeCount int, primaryLabel, primary, rowsLabel, rows string) template.HTML {
	badgeHTML := ""
	if activeCount > 0 {
		badgeHTML = fmt.Sprintf(`<span class="filter-count-badge">%d</span>`, activeCount)
	}
	return template.HTML(fmt.Sprintf(
		`<div class="report-header-actions"><button type="button" class="fycha-filter-btn" data-testid="report-filters-open-btn" aria-controls="sheetContent" aria-haspopup="dialog" hx-get="%s" hx-target="#sheetContent" hx-swap="innerHTML" hx-push-url="false" onclick="Sheet.open('Filters')"><svg class="icon" aria-hidden="true"><use href="#icon-filter"></use></svg><span>Filters</span>%s</button></div><div class="rr-active-filters"><span class="rr-chip" data-testid="rr-chip-primary"><span class="rr-chip-label">%s</span> <span class="rr-chip-value">%s</span></span><span class="rr-chip-sep">&times;</span><span class="rr-chip" data-testid="rr-chip-rows"><span class="rr-chip-label">%s</span> <span class="rr-chip-value">%s</span></span></div>`,
		template.HTMLEscapeString(filterSheetURL),
		badgeHTML,
		template.HTMLEscapeString(primaryLabel),
		template.HTMLEscapeString(primary),
		template.HTMLEscapeString(rowsLabel),
		template.HTMLEscapeString(rows),
	))
}

func formatCurrency(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
	whole := int64(amount)
	frac := int64((amount-float64(whole))*100 + 0.5)
	if frac >= 100 {
		whole++
		frac -= 100
	}
	wholeStr := strconv.FormatInt(whole, 10)
	n := len(wholeStr)
	if n > 3 {
		var result []byte
		for i, ch := range wholeStr {
			if i > 0 && (n-i)%3 == 0 {
				result = append(result, ',')
			}
			result = append(result, byte(ch))
		}
		wholeStr = string(result)
	}
	formatted := fmt.Sprintf("\u20b1%s.%02d", wholeStr, frac)
	if negative {
		formatted = "(" + formatted + ")"
	}
	return formatted
}
