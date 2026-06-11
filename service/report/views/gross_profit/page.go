package gross_profit

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	report "github.com/erniealice/fycha-golang/service/report"

	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
//
// 20260521 Wave B P1.E.3 — `DB report.DataSource` replaced with the typed
// `GetGrossProfitReport` closure into the service-driven gross/cashflow
// use case (proto package `service.reporting.v1`). Nil-safe: when the
// closure is nil the view renders an empty report instead of crashing.
type Deps struct {
	GetGrossProfitReport func(context.Context, *gcfpb.GetGrossProfitRequest) (*gcfpb.GetGrossProfitResponse, error)
	Labels               report.ReportsLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
}

// PageData holds the data for the gross profit report page.
type PageData struct {
	types.PageData
	ContentTemplate   string
	Labels            report.GrossProfitLabels
	Summary           []report.SummaryMetric
	Table             *types.TableConfig
	Filter            report.FilterState
	PeriodLabels      report.PeriodLabels
	ReportURL         string
	ActiveFilterCount int
	// Legacy fields used by gross profit specific filters
	ProductID  string
	LocationID string
	CategoryID string
}

// NewView creates the gross profit report view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.GrossProfit
		pl := deps.Labels.Period

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

		// Period preset
		period := viewCtx.QueryParams["period"]
		if period == "" {
			period = "thisMonth"
		}

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = report.ReportsGrossProfitURL
		}

		// Build group-by options for filter sheet
		groupByOptions := []report.FilterOption{
			{Value: "product", Label: l.GroupByProduct, Selected: groupBy == "product"},
			{Value: "location", Label: l.GroupByLocation, Selected: groupBy == "location"},
			{Value: "category", Label: l.GroupByCategory, Selected: groupBy == "category"},
			{Value: "monthly", Label: l.GroupByMonthly, Selected: groupBy == "monthly"},
			{Value: "quarterly", Label: l.GroupByQuarterly, Selected: groupBy == "quarterly"},
		}

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			sheetFilter := report.FilterState{
				ActivePreset:   period,
				StartDate:      startDateStr,
				EndDate:        endDateStr,
				GroupBy:        groupBy,
				GroupByOptions: groupByOptions,
				PeriodPresets:  report.DefaultPeriodPresets(pl, period),
			}
			return view.OK("report-filter-sheet", &report.FilterSheetData{
				Filter:       sheetFilter,
				PeriodLabels: pl,
				ReportURL:    reportURL,
			})
		}

		// Build proto request
		req := &gcfpb.GetGrossProfitRequest{}
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

		// Resolve dates
		if period == "custom" && startDateStr != "" {
			// Validate as date string; pass through if valid
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
			start, _ := report.ParsePeriodPreset(period)
			s := start.Format("2006-01-02")
			req.StartDate = &s
		}
		if req.EndDate == nil {
			_, end := report.ParsePeriodPreset(period)
			e := end.Format("2006-01-02")
			req.EndDate = &e
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

		// Call service-driven gross/cashflow use case (Wave B P1.E.3).
		var resp *gcfpb.GetGrossProfitResponse
		if deps.GetGrossProfitReport != nil {
			var err error
			resp, err = deps.GetGrossProfitReport(ctx, req)
			if err != nil {
				log.Printf("Failed to get gross profit report: %v", err)
				resp = nil
			}
		}
		if resp == nil {
			resp = &gcfpb.GetGrossProfitResponse{
				LineItems: []*gcfpb.GrossProfitLineItem{},
				Summary:   &gcfpb.GrossProfitSummary{},
			}
		}

		// Build summary bar
		summary := buildSummary(resp.GetSummary(), l)

		// Build table
		table := buildTable(resp.GetLineItems(), resp.GetSummary(), l, deps.TableLabels, groupBy)

		filter := report.FilterState{
			ActivePreset:   period,
			StartDate:      startDateStr,
			EndDate:        endDateStr,
			GroupBy:        groupBy,
			GroupByOptions: groupByOptions,
			PeriodPresets:  report.DefaultPeriodPresets(pl, period),
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "report",
				ActiveSubNav: "gross-profit",
				HeaderTitle:  l.Title,
				HeaderIcon:   "icon-bar-chart",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate:   "gross-profit-content",
			Labels:            l,
			Summary:           summary,
			Table:             table,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: report.ActiveFilterCount(filter),
			ProductID:         productID,
			LocationID:        locationID,
			CategoryID:        categoryID,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-gross-profit"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("gross-profit-content", pageData)
		}
		return view.OK("gross-profit", pageData)
	})
}

func buildSummary(s *gcfpb.GrossProfitSummary, l report.GrossProfitLabels) []report.SummaryMetric {
	if s == nil {
		s = &gcfpb.GrossProfitSummary{}
	}
	marginVariant := "success"
	if s.GetOverallMargin() < 15 {
		marginVariant = "danger"
	} else if s.GetOverallMargin() < 30 {
		marginVariant = "warning"
	}
	netRevenueCell := types.MoneyCell(float64(s.GetNetRevenue()), "PHP", true)
	cogsCell := types.MoneyCell(float64(s.GetTotalCogs()), "PHP", true)
	grossProfitCell := types.MoneyCell(float64(s.GetTotalGrossProfit()), "PHP", true)
	return []report.SummaryMetric{
		{Label: l.SummaryNetRevenue, Value: netRevenueCell.Currency + " " + netRevenueCell.Value},
		{Label: l.SummaryCogs, Value: cogsCell.Currency + " " + cogsCell.Value},
		{Label: l.SummaryGrossProfit, Value: grossProfitCell.Currency + " " + grossProfitCell.Value, Highlight: true},
		{Label: l.SummaryMargin, Value: fmt.Sprintf("%.1f%%", s.GetOverallMargin()), Variant: marginVariant},
	}
}

func buildTable(items []*gcfpb.GrossProfitLineItem, summary *gcfpb.GrossProfitSummary, l report.GrossProfitLabels, tableLabels types.TableLabels, groupBy string) *types.TableConfig {
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
					{Key: "totalRevenue", Label: l.GrossRevenue, Align: "right", MinWidth: "7.5rem"},
					{Key: "totalDiscount", Label: l.Discount, Align: "right", MinWidth: "6.25rem"},
					{Key: "netRevenue", Label: l.NetRevenue, Align: "right", MinWidth: "7.5rem"},
				},
			},
			{
				Label: l.ProfitabilityGroup,
				Columns: []types.TableColumn{
					{Key: "cogs", Label: l.COGS, Align: "right", MinWidth: "7.5rem"},
					{Key: "grossProfit", Label: l.GrossProfit, Align: "right", MinWidth: "7.5rem"},
					{Key: "margin", Label: l.Margin, Align: "right", MinWidth: "5rem"},
				},
			},
			{
				Label: l.VolumeGroup,
				Columns: []types.TableColumn{
					{Key: "unitsSold", Label: l.UnitsSold, Align: "right", MinWidth: "5rem"},
					{Key: "txnCount", Label: l.Transactions, Align: "right", MinWidth: "5rem"},
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

	rows := make([]types.TableRow, 0, len(items))
	for _, item := range items {
		marginVariant := "success"
		if item.GetGrossProfitMargin() < 15 {
			marginVariant = "danger"
		} else if item.GetGrossProfitMargin() < 30 {
			marginVariant = "warning"
		}

		row := types.TableRow{
			ID: item.GetGroupKey(),
			DataAttrs: map[string]string{
				"totalRevenue":  fmt.Sprintf("%d", item.GetTotalRevenue()),
				"totalDiscount": fmt.Sprintf("%d", item.GetTotalDiscount()),
				"netRevenue":    fmt.Sprintf("%d", item.GetNetRevenue()),
				"cogs":          fmt.Sprintf("%d", item.GetCostOfGoodsSold()),
				"grossProfit":   fmt.Sprintf("%d", item.GetGrossProfit()),
				"margin":        fmt.Sprintf("%.1f", item.GetGrossProfitMargin()),
				"unitsSold":     strconv.FormatInt(item.GetUnitsSold(), 10),
				"txnCount":      strconv.FormatInt(item.GetTransactionCount(), 10),
			},
			Cells: []types.TableCell{
				{Type: "name", Value: item.GetGroupKey()},
				types.MoneyCell(float64(item.GetTotalRevenue()), "PHP", true),
				types.MoneyCell(float64(item.GetTotalDiscount()), "PHP", true),
				types.MoneyCell(float64(item.GetNetRevenue()), "PHP", true),
				types.MoneyCell(float64(item.GetCostOfGoodsSold()), "PHP", true),
				types.MoneyCell(float64(item.GetGrossProfit()), "PHP", true),
				{Type: "badge", Value: fmt.Sprintf("%.1f%%", item.GetGrossProfitMargin()), Variant: marginVariant},
				{Type: "text", Value: strconv.FormatInt(item.GetUnitsSold(), 10)},
				{Type: "text", Value: strconv.FormatInt(item.GetTransactionCount(), 10)},
			},
		}
		rows = append(rows, row)
	}

	// Add totals row
	if summary != nil && len(items) > 0 {
		marginVariant := "success"
		if summary.GetOverallMargin() < 15 {
			marginVariant = "danger"
		} else if summary.GetOverallMargin() < 30 {
			marginVariant = "warning"
		}
		table.TotalsRow = []types.TableCell{
			{Type: "name", Value: l.Totals},
			types.MoneyCell(float64(summary.GetTotalRevenue()), "PHP", true),
			types.MoneyCell(float64(summary.GetTotalDiscount()), "PHP", true),
			types.MoneyCell(float64(summary.GetNetRevenue()), "PHP", true),
			types.MoneyCell(float64(summary.GetTotalCogs()), "PHP", true),
			types.MoneyCell(float64(summary.GetTotalGrossProfit()), "PHP", true),
			{Type: "badge", Value: fmt.Sprintf("%.1f%%", summary.GetOverallMargin()), Variant: marginVariant},
			{Type: "text", Value: strconv.FormatInt(summary.GetTotalUnitsSold(), 10)},
			{Type: "text", Value: strconv.FormatInt(summary.GetTotalTransactions(), 10)},
		}
	}

	table.Rows = rows
	types.ApplyColumnStyles(allColumns, rows)
	types.ApplyTableSettings(table)

	return table
}
