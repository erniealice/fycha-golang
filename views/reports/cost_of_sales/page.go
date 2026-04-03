package cost_of_sales

import (
	"context"
	"fmt"
	"log"
	"strconv"

	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	fycha "github.com/erniealice/fycha-golang"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

type Deps struct {
	DB           fycha.DataSource
	Labels       fycha.ReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

type PageData struct {
	types.PageData
	ContentTemplate   string
	Summary           []fycha.SummaryMetric
	Table             *types.TableConfig
	Filter            fycha.FilterState
	PeriodLabels      fycha.PeriodLabels
	ReportURL         string
	ActiveFilterCount int
}

func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.CostOfSales
		pl := deps.Labels.Period

		// Parse filter
		period := viewCtx.QueryParams["period"]
		if period == "" {
			period = "thisMonth"
		}
		startDateStr := viewCtx.QueryParams["start"]
		endDateStr := viewCtx.QueryParams["end"]

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = fycha.ReportsCostOfSalesURL
		}

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			sheetFilter := fycha.FilterState{
				ActivePreset:  period,
				StartDate:     startDateStr,
				EndDate:       endDateStr,
				PeriodPresets: fycha.DefaultPeriodPresets(pl, period),
			}
			return view.OK("report-filter-sheet", &fycha.FilterSheetData{
				Filter:       sheetFilter,
				PeriodLabels: pl,
				ReportURL:    reportURL,
			})
		}

		// Build proto request with date filtering
		req := &reportpb.GrossProfitReportRequest{}
		groupBy := "product"
		req.GroupBy = &groupBy

		// Resolve dates
		if period == "custom" && startDateStr != "" {
			req.StartDate = &startDateStr
		}
		if period == "custom" && endDateStr != "" {
			req.EndDate = &endDateStr
		}
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

		resp, err := deps.DB.GetGrossProfitReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get cost of sales report: %v", err)
			resp = &reportpb.GrossProfitReportResponse{
				LineItems: []*reportpb.GrossProfitLineItem{},
				Summary:   &reportpb.GrossProfitSummary{},
			}
		}

		s := resp.GetSummary()
		if s == nil {
			s = &reportpb.GrossProfitSummary{}
		}

		cogsRatio := 0.0
		if s.GetNetRevenue() > 0 {
			cogsRatio = (float64(s.GetTotalCogs()) / float64(s.GetNetRevenue())) * 100
		}

		summary := []fycha.SummaryMetric{
			{Label: l.SummaryTotalCOGS, Value: formatCurrency(float64(s.GetTotalCogs()) / 100.0), Highlight: true},
			{Label: l.SummaryRevenue, Value: formatCurrency(float64(s.GetNetRevenue()) / 100.0)},
			{Label: l.SummaryCOGSRatio, Value: fmt.Sprintf("%.1f%%", cogsRatio)},
			{Label: l.SummaryUnits, Value: strconv.FormatInt(s.GetTotalUnitsSold(), 10)},
		}

		// Table
		columns := []types.TableColumn{
			{Key: "group", Label: l.Item, Sortable: true},
			{Key: "cogs", Label: l.COGS, Sortable: true, Align: "right", MinWidth: "120px"},
			{Key: "revenue", Label: l.NetRevenue, Sortable: true, Align: "right", MinWidth: "120px"},
			{Key: "ratio", Label: l.COGSPct, Sortable: true, Align: "right", MinWidth: "80px"},
			{Key: "units", Label: l.Units, Sortable: true, Align: "right", MinWidth: "80px"},
		}

		rows := make([]types.TableRow, 0, len(resp.GetLineItems()))
		for _, item := range resp.GetLineItems() {
			ratio := 0.0
			if item.GetNetRevenue() > 0 {
				ratio = (float64(item.GetCostOfGoodsSold()) / float64(item.GetNetRevenue())) * 100
			}
			rows = append(rows, types.TableRow{
				ID: item.GetGroupKey(),
				Cells: []types.TableCell{
					{Type: "name", Value: item.GetGroupKey()},
					{Type: "text", Value: formatCurrency(float64(item.GetCostOfGoodsSold()) / 100.0)},
					{Type: "text", Value: formatCurrency(float64(item.GetNetRevenue()) / 100.0)},
					{Type: "text", Value: fmt.Sprintf("%.1f%%", ratio)},
					{Type: "text", Value: strconv.FormatInt(item.GetUnitsSold(), 10)},
				},
				DataAttrs: map[string]string{
					"cogs":    fmt.Sprintf("%.2f", float64(item.GetCostOfGoodsSold())/100.0),
					"revenue": fmt.Sprintf("%.2f", float64(item.GetNetRevenue())/100.0),
					"ratio":   fmt.Sprintf("%.1f", ratio),
					"units":   strconv.FormatInt(item.GetUnitsSold(), 10),
				},
			})
		}

		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "cost-of-sales-table",
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           false,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "cogs",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   l.EmptyTitle,
				Message: l.EmptyMessage,
			},
		}
		types.ApplyTableSettings(tableConfig)

		filter := fycha.FilterState{
			ActivePreset:  period,
			StartDate:     startDateStr,
			EndDate:       endDateStr,
			PeriodPresets: fycha.DefaultPeriodPresets(pl, period),
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "reports",
				ActiveSubNav:   "cost-of-sales",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:   "cost-of-sales-content",
			Summary:           summary,
			Table:             tableConfig,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: fycha.ActiveFilterCount(filter),
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-cost-of-sale"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("cost-of-sales-content", pageData)
		}
		return view.OK("cost-of-sales", pageData)
	})
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
		formatted = "-" + formatted
	}
	return formatted
}
