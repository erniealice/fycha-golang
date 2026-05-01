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

		cogsCell := types.MoneyCell(float64(s.GetTotalCogs()), "PHP", true)
		revenueCell := types.MoneyCell(float64(s.GetNetRevenue()), "PHP", true)
		summary := []fycha.SummaryMetric{
			{Label: l.SummaryTotalCOGS, Value: cogsCell.Currency + " " + cogsCell.Value, Highlight: true},
			{Label: l.SummaryRevenue, Value: revenueCell.Currency + " " + revenueCell.Value},
			{Label: l.SummaryCOGSRatio, Value: fmt.Sprintf("%.1f%%", cogsRatio)},
			{Label: l.SummaryUnits, Value: strconv.FormatInt(s.GetTotalUnitsSold(), 10)},
		}

		// Table
		columns := []types.TableColumn{
			{Key: "group", Label: l.Item},
			{Key: "cogs", Label: l.COGS, Align: "right", MinWidth: "7.5rem"},
			{Key: "revenue", Label: l.NetRevenue, Align: "right", MinWidth: "7.5rem"},
			{Key: "ratio", Label: l.COGSPct, Align: "right", MinWidth: "5rem"},
			{Key: "units", Label: l.Units, Align: "right", MinWidth: "5rem"},
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
					types.MoneyCell(float64(item.GetCostOfGoodsSold()), "PHP", true),
					types.MoneyCell(float64(item.GetNetRevenue()), "PHP", true),
					{Type: "text", Value: fmt.Sprintf("%.1f%%", ratio)},
					{Type: "text", Value: strconv.FormatInt(item.GetUnitsSold(), 10)},
				},
				DataAttrs: map[string]string{
					"cogs":    fmt.Sprintf("%d", item.GetCostOfGoodsSold()),
					"revenue": fmt.Sprintf("%d", item.GetNetRevenue()),
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
				ActiveNav:      "report",
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
