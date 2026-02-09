package gross_profit

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
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
}

// SummaryMetric holds a single summary bar metric.
type SummaryMetric struct {
	Label     string
	Value     string
	Highlight bool
	Variant   string // "success", "warning", "danger" — for margin badge
}

// FilterOption holds a dropdown option for filters.
type FilterOption struct {
	Value    string
	Label    string
	Selected bool
}

// PageData holds the data for the gross profit report page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          fycha.GrossProfitLabels
	Summary         []SummaryMetric
	Table           *types.TableConfig
	// Filter state
	GroupBy       string
	StartDate     string
	EndDate       string
	ProductID     string
	LocationID    string
	CategoryID    string
	ActivePreset  string
	GroupByOptions []FilterOption
	PeriodPresets  []FilterOption
}

// NewView creates the gross profit report view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.GrossProfit

		// Parse filter query params
		groupBy := viewCtx.QueryParams["group-by"]
		if groupBy == "" {
			groupBy = "product"
		}
		startDateStr := viewCtx.QueryParams["start"]
		endDateStr := viewCtx.QueryParams["end"]
		productID := viewCtx.QueryParams["product-id"]
		locationID := viewCtx.QueryParams["location-id"]
		categoryID := viewCtx.QueryParams["category-id"]

		// Build proto request
		req := &reportpb.GrossProfitReportRequest{}
		req.GroupBy = &groupBy

		// Handle period granularity for monthly/quarterly group-by
		if groupBy == "monthly" {
			gb := "period"
			req.GroupBy = &gb
			gran := "monthly"
			req.PeriodGranularity = &gran
		} else if groupBy == "quarterly" {
			gb := "period"
			req.GroupBy = &gb
			gran := "quarterly"
			req.PeriodGranularity = &gran
		}

		// Parse date filters
		activePreset := "thisMonth"
		if startDateStr != "" {
			if ts, err := strconv.ParseInt(startDateStr, 10, 64); err == nil {
				req.StartDate = &ts
				activePreset = "custom"
			}
		}
		if endDateStr != "" {
			if ts, err := strconv.ParseInt(endDateStr, 10, 64); err == nil {
				req.EndDate = &ts
				activePreset = "custom"
			}
		}

		// Apply default period (this month) if no dates provided
		if req.StartDate == nil {
			now := time.Now()
			startOfMonth := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local)
			ts := startOfMonth.Unix()
			req.StartDate = &ts
		}
		if req.EndDate == nil {
			now := time.Now()
			ts := now.Unix()
			req.EndDate = &ts
		}

		// Apply optional filters
		if productID != "" {
			req.ProductId = &productID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if categoryID != "" {
			req.RevenueCategoryId = &categoryID
		}

		// Call data source
		resp, err := deps.DB.GetGrossProfitReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get gross profit report: %v", err)
			resp = &reportpb.GrossProfitReportResponse{
				LineItems: []*reportpb.GrossProfitLineItem{},
				Summary:   &reportpb.GrossProfitSummary{},
			}
		}

		// Build summary bar
		summary := buildSummary(resp.GetSummary(), l)

		// Build table
		table := buildTable(resp.GetLineItems(), resp.GetSummary(), l, deps.TableLabels, groupBy)

		// Build group-by options
		groupByOptions := []FilterOption{
			{Value: "product", Label: l.GroupByProduct, Selected: groupBy == "product"},
			{Value: "location", Label: l.GroupByLocation, Selected: groupBy == "location"},
			{Value: "category", Label: l.GroupByCategory, Selected: groupBy == "category"},
			{Value: "monthly", Label: l.GroupByMonthly, Selected: groupBy == "monthly"},
			{Value: "quarterly", Label: l.GroupByQuarterly, Selected: groupBy == "quarterly"},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "reports",
				ActiveSubNav: "gross-profit",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-bar-chart",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "gross-profit-content",
			Labels:          l,
			Summary:         summary,
			Table:           table,
			GroupBy:         groupBy,
			StartDate:       startDateStr,
			EndDate:         endDateStr,
			ProductID:       productID,
			LocationID:      locationID,
			CategoryID:      categoryID,
			ActivePreset:    activePreset,
			GroupByOptions:  groupByOptions,
		}

		if viewCtx.IsHTMX {
			return view.OK("gross-profit-content", pageData)
		}
		return view.OK("gross-profit", pageData)
	})
}

func buildSummary(s *reportpb.GrossProfitSummary, l fycha.GrossProfitLabels) []SummaryMetric {
	if s == nil {
		s = &reportpb.GrossProfitSummary{}
	}
	marginVariant := "success"
	if s.GetOverallMargin() < 15 {
		marginVariant = "danger"
	} else if s.GetOverallMargin() < 30 {
		marginVariant = "warning"
	}
	return []SummaryMetric{
		{Label: l.SummaryNetRevenue, Value: formatCurrency(s.GetNetRevenue())},
		{Label: l.SummaryCogs, Value: formatCurrency(s.GetTotalCogs())},
		{Label: l.SummaryGrossProfit, Value: formatCurrency(s.GetTotalGrossProfit()), Highlight: true},
		{Label: l.SummaryMargin, Value: fmt.Sprintf("%.1f%%", s.GetOverallMargin()), Variant: marginVariant},
	}
}

func buildTable(items []*reportpb.GrossProfitLineItem, summary *reportpb.GrossProfitSummary, l fycha.GrossProfitLabels, tableLabels types.TableLabels, groupBy string) *types.TableConfig {
	table := &types.TableConfig{
		ID:          "grossProfitTable",
		ShowSearch:  false,
		ShowFilters: true,
		ShowSort:    true,
		ShowColumns: true,
		ShowExport:  true,
		ShowEntries: true,
		ShowDensity: true,
		Labels:      tableLabels,
		ColumnGroups: []types.ColumnGroup{
			{
				Label: l.RevenueGroup,
				Columns: []types.TableColumn{
					{Key: "totalRevenue", Label: l.GrossRevenue, Sortable: true, Align: "right", MinWidth: "120px"},
					{Key: "totalDiscount", Label: l.Discount, Sortable: true, Align: "right", MinWidth: "100px"},
					{Key: "netRevenue", Label: l.NetRevenue, Sortable: true, Align: "right", MinWidth: "120px"},
				},
			},
			{
				Label: l.ProfitabilityGroup,
				Columns: []types.TableColumn{
					{Key: "cogs", Label: l.COGS, Sortable: true, Align: "right", MinWidth: "120px"},
					{Key: "grossProfit", Label: l.GrossProfit, Sortable: true, Align: "right", MinWidth: "120px"},
					{Key: "margin", Label: l.Margin, Sortable: true, Align: "right", MinWidth: "80px"},
				},
			},
			{
				Label: l.VolumeGroup,
				Columns: []types.TableColumn{
					{Key: "unitsSold", Label: l.UnitsSold, Sortable: true, Align: "right", MinWidth: "80px"},
					{Key: "txnCount", Label: l.Transactions, Sortable: true, Align: "right", MinWidth: "80px"},
				},
			},
		},
		EmptyState: types.TableEmptyState{
			Title:   "No data",
			Message: "No gross profit data found for the selected period.",
		},
	}

	// Flatten columns for ApplyColumnStyles
	var allColumns []types.TableColumn
	for _, group := range table.ColumnGroups {
		allColumns = append(allColumns, group.Columns...)
	}

	rows := make([]types.TableRow, 0, len(items))
	for _, item := range items {
		// Margin badge variant
		marginVariant := "success"
		if item.GetGrossProfitMargin() < 15 {
			marginVariant = "danger"
		} else if item.GetGrossProfitMargin() < 30 {
			marginVariant = "warning"
		}

		row := types.TableRow{
			ID: item.GetGroupKey(),
			DataAttrs: map[string]string{
				"totalRevenue":  fmt.Sprintf("%.2f", item.GetTotalRevenue()),
				"totalDiscount": fmt.Sprintf("%.2f", item.GetTotalDiscount()),
				"netRevenue":    fmt.Sprintf("%.2f", item.GetNetRevenue()),
				"cogs":          fmt.Sprintf("%.2f", item.GetCostOfGoodsSold()),
				"grossProfit":   fmt.Sprintf("%.2f", item.GetGrossProfit()),
				"margin":        fmt.Sprintf("%.1f", item.GetGrossProfitMargin()),
				"unitsSold":     strconv.FormatInt(item.GetUnitsSold(), 10),
				"txnCount":      strconv.FormatInt(item.GetTransactionCount(), 10),
			},
			Cells: []types.TableCell{
				// Revenue group
				{Type: "name", Value: item.GetGroupKey()},
				{Type: "text", Value: formatCurrency(item.GetTotalRevenue())},
				{Type: "text", Value: formatCurrency(item.GetTotalDiscount())},
				{Type: "text", Value: formatCurrency(item.GetNetRevenue())},
				// Profitability group
				{Type: "text", Value: formatCurrency(item.GetCostOfGoodsSold())},
				{Type: "text", Value: formatCurrency(item.GetGrossProfit())},
				{Type: "badge", Value: fmt.Sprintf("%.1f%%", item.GetGrossProfitMargin()), Variant: marginVariant},
				// Volume group
				{Type: "text", Value: strconv.FormatInt(item.GetUnitsSold(), 10)},
				{Type: "text", Value: strconv.FormatInt(item.GetTransactionCount(), 10)},
			},
		}
		rows = append(rows, row)
	}

	// Add totals row if we have summary data
	if summary != nil && len(items) > 0 {
		marginVariant := "success"
		if summary.GetOverallMargin() < 15 {
			marginVariant = "danger"
		} else if summary.GetOverallMargin() < 30 {
			marginVariant = "warning"
		}
		totalsRow := types.TableRow{
			ID: "__totals__",
			Cells: []types.TableCell{
				{Type: "name", Value: "TOTALS"},
				{Type: "text", Value: formatCurrency(summary.GetTotalRevenue())},
				{Type: "text", Value: formatCurrency(summary.GetTotalDiscount())},
				{Type: "text", Value: formatCurrency(summary.GetNetRevenue())},
				{Type: "text", Value: formatCurrency(summary.GetTotalCogs())},
				{Type: "text", Value: formatCurrency(summary.GetTotalGrossProfit())},
				{Type: "badge", Value: fmt.Sprintf("%.1f%%", summary.GetOverallMargin()), Variant: marginVariant},
				{Type: "text", Value: strconv.FormatInt(summary.GetTotalUnitsSold(), 10)},
				{Type: "text", Value: strconv.FormatInt(summary.GetTotalTransactions(), 10)},
			},
		}
		rows = append(rows, totalsRow)
	}

	table.Rows = rows
	types.ApplyColumnStyles(allColumns, rows)
	types.ApplyTableSettings(table)

	return table
}

// formatCurrency formats a float64 as PHP currency with comma separators.
func formatCurrency(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}

	// Format with 2 decimal places
	whole := int64(amount)
	frac := int64((amount - float64(whole)) * 100 + 0.5)
	if frac >= 100 {
		whole++
		frac -= 100
	}

	// Add comma separators to whole part
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

	formatted := fmt.Sprintf("₱%s.%02d", wholeStr, frac)
	if negative {
		formatted = "-" + formatted
	}
	return formatted
}
