package expenditure_report

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	dspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/domain_specific"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
//
// 20260521 Wave B P1.E.5 — DB replaced with the typed GetExpenditureReport
// closure into the service-driven domain-specific use case.
type Deps struct {
	GetExpenditureReport func(context.Context, *dspb.GetExpenditureReportRequest) (*dspb.GetExpenditureReportResponse, error)
	Labels               fycha.ReportsLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
	Routes               fycha.ReportsRoutes
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
		req := &dspb.GetExpenditureReportRequest{
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
		var resp *dspb.GetExpenditureReportResponse
		var err error
		if deps.GetExpenditureReport != nil {
			resp, err = deps.GetExpenditureReport(ctx, req)
		}
		if err != nil {
			log.Printf("Failed to get expenditure report: %v", err)
			resp = nil
		}
		if resp == nil {
			resp = &dspb.GetExpenditureReportResponse{
				ColumnKeys: []string{},
				Rows:       []*dspb.ExpenditureReportRow{},
				Summary:    &dspb.ExpenditureReportSummary{},
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
		table.ToolbarPrefixTemplate = "report-dimension-toolbar-prefix"
		table.ToolbarPrefixData = fycha.DimensionToolbarPrefixData{
			FilterSheetURL:    filterSheetURL,
			ActiveFilterCount: activeCount,
			PrimaryLabel:      "Columns:",
			PrimaryValue:      primary,
			RowsLabel:         "Rows:",
			RowsValue:         rows,
		}

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
			ContentTemplate:   "expenditure-report-content",
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
	Nonce            string // injected by ViewAdapter.injectPageData via reflection
}

func buildSummary(s *dspb.ExpenditureReportSummary, l fycha.ExpenditureReportLabels) []fycha.SummaryMetric {
	if s == nil {
		s = &dspb.ExpenditureReportSummary{}
	}
	grandTotalCents := s.GetGrandTotal()
	txnCount := s.GetTotalTransactions()
	avgCents := 0.0
	if txnCount > 0 {
		avgCents = float64(grandTotalCents) / float64(txnCount)
	}
	grandCell := types.MoneyCell(float64(grandTotalCents), "PHP", true)
	avgCell := types.MoneyCell(avgCents, "PHP", true)
	return []fycha.SummaryMetric{
		{Label: l.SummaryGrandTotal, Value: grandCell.Currency + " " + grandCell.Value, Highlight: true},
		{Label: l.SummaryTransactions, Value: fmt.Sprintf("%d", txnCount)},
		{Label: l.SummaryAverage, Value: avgCell.Currency + " " + avgCell.Value},
	}
}

func buildPivotTable(resp *dspb.GetExpenditureReportResponse, l fycha.ExpenditureReportLabels, tableLabels types.TableLabels, primary, rowDim string) *types.TableConfig {
	columnKeys := resp.GetColumnKeys()

	// Build dynamic columns
	dynamicColumns := make([]types.TableColumn, 0, len(columnKeys))
	for _, ck := range columnKeys {
		dynamicColumns = append(dynamicColumns, types.TableColumn{
			Key:      ck,
			Label:    ck,
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
					{Key: "total", Label: l.Total, Align: "right", MinWidth: "8.125rem"},
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

	currency := "PHP"
	tableRows := make([]types.TableRow, 0, len(resp.GetRows()))
	for i, row := range resp.GetRows() {
		cellMap := make(map[string]*dspb.ExpenditureReportCell, len(row.GetCells()))
		for _, c := range row.GetCells() {
			cellMap[c.GetColumnKey()] = c
		}

		rowCurrency := ""
		if i == 0 {
			rowCurrency = currency
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
			cells = append(cells, types.MoneyCell(float64(val), rowCurrency, true))
			dataAttrs[ck] = fmt.Sprintf("%d", val)
		}

		// Total cell
		cells = append(cells, types.MoneyCell(float64(row.GetRowTotal()), rowCurrency, true))
		dataAttrs["total"] = fmt.Sprintf("%d", row.GetRowTotal())

		tableRows = append(tableRows, types.TableRow{
			ID:        row.GetRowKey(),
			Cells:     cells,
			DataAttrs: dataAttrs,
		})
	}

	// Add totals row from summary.column_totals
	summary := resp.GetSummary()
	if summary != nil && len(resp.GetRows()) > 0 {
		colTotalMap := make(map[string]*dspb.ExpenditureReportCell, len(summary.GetColumnTotals()))
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
			totalsCells = append(totalsCells, types.MoneyCell(float64(val), currency, true))
		}
		totalsCells = append(totalsCells, types.MoneyCell(float64(summary.GetGrandTotal()), currency, true))

		table.TotalsRow = totalsCells
	}

	table.Rows = tableRows
	types.ApplyColumnStyles(allColumns, tableRows)
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
