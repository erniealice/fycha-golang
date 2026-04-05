package disbursement_report

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	disbreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/disbursement_report"
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

// PageData holds the data for the disbursement report page.
type PageData struct {
	types.PageData
	ContentTemplate   string
	Labels            fycha.DisbursementReportLabels
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

// NewView creates the disbursement report pivot-table view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.DisbursementReport
		pl := deps.Labels.Period

		// Parse query params
		primary := viewCtx.QueryParams["primary"]
		if primary == "" {
			primary = "monthly"
		}
		rows := viewCtx.QueryParams["rows"]
		if rows == "" {
			rows = "supplier"
		}
		period := viewCtx.QueryParams["period"]
		if period == "" {
			period = "thisMonth"
		}
		startDateStr := viewCtx.QueryParams["start"]
		endDateStr := viewCtx.QueryParams["end"]

		// Secondary filter IDs
		supplierID := viewCtx.QueryParams["supplier-id"]
		supplierCategoryID := viewCtx.QueryParams["supplier-category-id"]
		locationID := viewCtx.QueryParams["location-id"]
		expenditureCategoryID := viewCtx.QueryParams["expenditure-category-id"]
		disbursementType := viewCtx.QueryParams["disbursement-type"]
		disbursementMethodID := viewCtx.QueryParams["disbursement-method-id"]

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = deps.Routes.DisbursementReportURL
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
			return view.OK("disbursement-report-filter-sheet", &DisbursementReportFilterSheetData{
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
		req := &disbreportpb.DisbursementReportRequest{
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
		if supplierID != "" {
			req.SupplierId = &supplierID
		}
		if supplierCategoryID != "" {
			req.SupplierCategoryId = &supplierCategoryID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if expenditureCategoryID != "" {
			req.ExpenditureCategoryId = &expenditureCategoryID
		}
		if disbursementType != "" {
			req.DisbursementType = &disbursementType
		}
		if disbursementMethodID != "" {
			req.DisbursementMethodId = &disbursementMethodID
		}

		// Call data source
		resp, err := deps.DB.GetDisbursementReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get disbursement report: %v", err)
			resp = &disbreportpb.DisbursementReportResponse{
				ColumnKeys: []string{},
				Rows:       []*disbreportpb.DisbursementReportRow{},
				Summary:    &disbreportpb.DisbursementReportSummary{},
			}
		}

		// Build summary bar
		summary := buildSummary(resp.GetSummary(), l)

		// Build pivot table
		table := buildPivotTable(resp, l, deps.TableLabels, primary, rows)

		filter := fycha.FilterState{
			ActivePreset:  period,
			StartDate:     startDateStr,
			EndDate:       endDateStr,
			PeriodPresets: fycha.DefaultPeriodPresets(pl, period),
		}

		// Build export URL with current query params
		exportURL := buildExportURL(deps.Routes.DisbursementReportExportURL, primary, rows, period, startDateStr, endDateStr)

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
		if rows != "" && rows != "supplier" {
			activeCount++
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "report",
				ActiveSubNav: "disbursement-report",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-bar-chart",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate:   "disbursement-report-content",
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
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-disbursement-report"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("disbursement-report-content", pageData)
		}
		return view.OK("disbursement-report", pageData)
	})
}

// DisbursementReportFilterSheetData holds data for the disbursement report filter sheet template.
type DisbursementReportFilterSheetData struct {
	Filter           fycha.FilterState
	PeriodLabels     fycha.PeriodLabels
	Labels           fycha.DisbursementReportLabels
	ReportURL        string
	PrimaryDimension string
	RowDimension     string
	PrimaryOptions   []fycha.FilterOption
	RowOptions       []fycha.FilterOption
}

func buildSummary(s *disbreportpb.DisbursementReportSummary, l fycha.DisbursementReportLabels) []fycha.SummaryMetric {
	if s == nil {
		s = &disbreportpb.DisbursementReportSummary{}
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

func buildPivotTable(resp *disbreportpb.DisbursementReportResponse, l fycha.DisbursementReportLabels, tableLabels types.TableLabels, primary, rowDim string) *types.TableConfig {
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
		ID:              "disbursementReportTable",
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
		cellMap := make(map[string]*disbreportpb.DisbursementReportCell, len(row.GetCells()))
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
				val = c.GetTotalDisbursement()
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
		colTotalMap := make(map[string]*disbreportpb.DisbursementReportCell, len(summary.GetColumnTotals()))
		for _, ct := range summary.GetColumnTotals() {
			colTotalMap[ct.GetColumnKey()] = ct
		}

		totalsCells := []types.TableCell{
			{Type: "name", Value: l.Totals},
		}
		for _, ck := range columnKeys {
			var val int64
			if ct, ok := colTotalMap[ck]; ok {
				val = ct.GetTotalDisbursement()
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
