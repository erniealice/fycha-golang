package payables_aging_report

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	payagingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/payables_aging"
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

// PageData holds the data for the payables aging report page.
type PageData struct {
	types.PageData
	ContentTemplate   string
	Labels            fycha.PayablesAgingReportLabels
	Summary           []fycha.SummaryMetric
	Table             *types.TableConfig
	AsOfDate          string
	RowDimension      string
	RowOptions        []fycha.FilterOption
	ReportURL         string
	FilterSheetURL    string
	ExportURL         string
	ActiveFilterCount int
}

// NewView creates the payables aging report view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.PayablesAging

		// Parse query params
		asOfDate := viewCtx.QueryParams["as-of-date"]
		if asOfDate == "" {
			asOfDate = time.Now().Format("2006-01-02")
		}
		rows := viewCtx.QueryParams["rows"]
		if rows == "" {
			rows = "supplier"
		}

		// Secondary filter IDs
		supplierID := viewCtx.QueryParams["supplier-id"]
		locationID := viewCtx.QueryParams["location-id"]
		expenditureCategoryID := viewCtx.QueryParams["expenditure-category-id"]

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = deps.Routes.PayablesAgingReportURL
		}

		// Build row dimension options
		rowOptions := []fycha.FilterOption{
			{Value: "supplier", Label: l.DimensionSupplier, Selected: rows == "supplier"},
			{Value: "supplierCategory", Label: l.DimensionSupplierCategory, Selected: rows == "supplierCategory"},
			{Value: "location", Label: l.DimensionLocation, Selected: rows == "location"},
			{Value: "locationArea", Label: l.DimensionLocationArea, Selected: rows == "locationArea"},
			{Value: "expenditureCategory", Label: l.DimensionExpenditureCategory, Selected: rows == "expenditureCategory"},
		}

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			return view.OK("payables-aging-report-filter-sheet", &FilterSheetData{
				Labels:       l,
				ReportURL:    reportURL,
				AsOfDate:     asOfDate,
				RowDimension: rows,
				RowOptions:   rowOptions,
			})
		}

		// Build proto request
		req := &payagingpb.PayablesAgingRequest{
			AsOfDate:     &asOfDate,
			RowDimension: rows,
		}

		// Apply optional secondary filters
		if supplierID != "" {
			req.SupplierId = &supplierID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if expenditureCategoryID != "" {
			req.ExpenditureCategoryId = &expenditureCategoryID
		}

		// Call data source
		resp, err := deps.DB.GetPayablesAgingReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get payables aging report: %v", err)
			resp = &payagingpb.PayablesAgingResponse{
				BucketLabels: []string{},
				Rows:         []*payagingpb.PayablesAgingRow{},
				Summary:      &payagingpb.PayablesAgingSummary{},
			}
		}

		// Build summary bar
		summary := buildSummary(resp.GetSummary(), l)

		// Build fixed-column table
		table := buildTable(resp, l, deps.TableLabels, rows)

		// Build export URL with current query params
		exportURL := buildExportURL(deps.Routes.PayablesAgingReportExportURL, asOfDate, rows)

		// Build filter sheet URL
		filterSheetURL := buildFilterSheetURL(reportURL, asOfDate, rows)

		// Count active filters
		activeCount := 0
		if rows != "" && rows != "supplier" {
			activeCount++
		}
		if asOfDate != "" && asOfDate != time.Now().Format("2006-01-02") {
			activeCount++
		}

		// Inject filter button + dimension chips into the table toolbar prefix
		table.ToolbarPrefixTemplate = "report-aging-toolbar-prefix"
		table.ToolbarPrefixData = fycha.AgingToolbarPrefixData{
			FilterSheetURL:    filterSheetURL,
			ActiveFilterCount: activeCount,
			AsOfDate:          asOfDate,
			GroupByValue:      rows,
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.PageTitle,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "supplier",
				ActiveSubNav: "payables-aging-report",
				HeaderTitle:  l.PageTitle,
				HeaderIcon:   "icon-bar-chart",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate:   "payables-aging-report-content",
			Labels:            l,
			Summary:           summary,
			Table:             table,
			AsOfDate:          asOfDate,
			RowDimension:      rows,
			RowOptions:        rowOptions,
			ReportURL:         reportURL,
			FilterSheetURL:    filterSheetURL,
			ExportURL:         exportURL,
			ActiveFilterCount: activeCount,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-payables-aging-report"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("payables-aging-report-content", pageData)
		}
		return view.OK("payables-aging-report", pageData)
	})
}

// FilterSheetData holds data for the payables aging report filter sheet template.
type FilterSheetData struct {
	Labels       fycha.PayablesAgingReportLabels
	ReportURL    string
	AsOfDate     string
	RowDimension string
	RowOptions   []fycha.FilterOption
}

func buildSummary(s *payagingpb.PayablesAgingSummary, l fycha.PayablesAgingReportLabels) []fycha.SummaryMetric {
	if s == nil {
		s = &payagingpb.PayablesAgingSummary{}
	}
	grandTotal := float64(s.GetGrandTotalOutstanding()) / 100.0
	invoiceCount := s.GetTotalInvoiceCount()

	// Overdue = everything past the current bucket
	var currentBucket int64
	if b := s.GetBuckets(); b != nil {
		currentBucket = b.GetCurrent()
	}
	overdueTotal := float64(s.GetGrandTotalOutstanding()-currentBucket) / 100.0

	return []fycha.SummaryMetric{
		{Label: l.SummaryGrandTotal, Value: formatCurrency(grandTotal), Highlight: true},
		{Label: l.SummaryInvoiceCount, Value: fmt.Sprintf("%d", invoiceCount)},
		{Label: l.SummaryOverdueAmount, Value: formatCurrency(overdueTotal), Variant: "danger"},
	}
}

func buildTable(resp *payagingpb.PayablesAgingResponse, l fycha.PayablesAgingReportLabels, tableLabels types.TableLabels, rowDim string) *types.TableConfig {
	// Fixed columns for aging buckets. The name column is listed first so that
	// ApplyColumnStyles maps columns[i] to cells[i] correctly (cells[0] is the
	// "name" type cell; columns[0] must correspond to it).
	columns := []types.TableColumn{
		{Key: "row_key", Label: rowDimensionLabel(l, rowDim), Sortable: true},
		{Key: "current", Label: l.BucketCurrent, Sortable: true, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_1_30", Label: l.Bucket1To30, Sortable: true, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_31_60", Label: l.Bucket31To60, Sortable: true, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_61_90", Label: l.Bucket61To90, Sortable: true, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_over_90", Label: l.BucketOver90, Sortable: true, Align: "right", MinWidth: "7.5rem"},
		{Key: "total", Label: l.TotalOutstanding, Sortable: true, Align: "right", MinWidth: "8.125rem"},
		{Key: "invoice_count", Label: l.InvoiceCount, Sortable: true, Align: "right", MinWidth: "6rem"},
	}

	table := &types.TableConfig{
		ID:      "payablesAgingReportTable",
		Columns: columns,
		ShowSearch:      false,
		ShowFilters:     false,
		ShowSort:        false,
		ShowColumns:     false,
		ShowExport:      true,
		ShowEntries:     true,
		ShowDensity:     true,
		Labels:          tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.EmptyTitle,
			Message: l.EmptyMessage,
		},
	}

	rows := make([]types.TableRow, 0, len(resp.GetRows()))
	for _, row := range resp.GetRows() {
		b := row.GetBuckets()
		if b == nil {
			b = &payagingpb.PayablesAgingBuckets{}
		}

		cells := []types.TableCell{
			{Type: "name", Value: row.GetRowKey()},
			{Type: "text", Value: formatCurrency(float64(b.GetCurrent()) / 100.0)},
			{Type: "text", Value: formatCurrency(float64(b.GetDays_1_30()) / 100.0)},
			{Type: "text", Value: formatCurrency(float64(b.GetDays_31_60()) / 100.0)},
			{Type: "text", Value: formatCurrency(float64(b.GetDays_61_90()) / 100.0)},
			{Type: "text", Value: formatCurrency(float64(b.GetDaysOver_90()) / 100.0)},
			{Type: "text", Value: formatCurrency(float64(row.GetTotalOutstanding()) / 100.0)},
			{Type: "text", Value: fmt.Sprintf("%d", row.GetInvoiceCount())},
		}

		dataAttrs := map[string]string{
			"current":       fmt.Sprintf("%.2f", float64(b.GetCurrent())/100.0),
			"days_1_30":     fmt.Sprintf("%.2f", float64(b.GetDays_1_30())/100.0),
			"days_31_60":    fmt.Sprintf("%.2f", float64(b.GetDays_31_60())/100.0),
			"days_61_90":    fmt.Sprintf("%.2f", float64(b.GetDays_61_90())/100.0),
			"days_over_90":  fmt.Sprintf("%.2f", float64(b.GetDaysOver_90())/100.0),
			"total":         fmt.Sprintf("%.2f", float64(row.GetTotalOutstanding())/100.0),
			"invoice_count": fmt.Sprintf("%d", row.GetInvoiceCount()),
		}

		rows = append(rows, types.TableRow{
			ID:        row.GetRowKey(),
			Cells:     cells,
			DataAttrs: dataAttrs,
		})
	}

	table.Rows = rows
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(table)

	// Build tfoot totals from summary
	summary := resp.GetSummary()
	if summary != nil && len(resp.GetRows()) > 0 {
		sb := summary.GetBuckets()
		if sb == nil {
			sb = &payagingpb.PayablesAgingBuckets{}
		}
		table.TotalsRow = []types.TableCell{
			{Value: "Total"},
			{Value: formatCurrency(float64(sb.GetCurrent()) / 100.0), Align: "right"},
			{Value: formatCurrency(float64(sb.GetDays_1_30()) / 100.0), Align: "right"},
			{Value: formatCurrency(float64(sb.GetDays_31_60()) / 100.0), Align: "right"},
			{Value: formatCurrency(float64(sb.GetDays_61_90()) / 100.0), Align: "right"},
			{Value: formatCurrency(float64(sb.GetDaysOver_90()) / 100.0), Align: "right"},
			{Value: formatCurrency(float64(summary.GetGrandTotalOutstanding()) / 100.0), Align: "right"},
			{Value: fmt.Sprintf("%d", summary.GetTotalInvoiceCount()), Align: "right"},
		}
	}

	return table
}

func rowDimensionLabel(l fycha.PayablesAgingReportLabels, dim string) string {
	return l.PrimaryGroupLabel(dim)
}

func buildExportURL(base, asOfDate, rows string) string {
	params := url.Values{}
	params.Set("as-of-date", asOfDate)
	params.Set("rows", rows)
	return base + "?" + params.Encode()
}

func buildFilterSheetURL(base, asOfDate, rows string) string {
	params := url.Values{}
	params.Set("sheet", "filters")
	params.Set("as-of-date", asOfDate)
	params.Set("rows", rows)
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
