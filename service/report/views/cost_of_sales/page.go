package cost_of_sales

import (
	"context"
	"fmt"
	"log"
	"strconv"

	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"
	report "github.com/erniealice/fycha-golang/service/report"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// 20260521 Wave B P1.E.3 — DB replaced with the typed GetGrossProfitReport
// closure into the service-driven gross/cashflow use case.
type Deps struct {
	GetGrossProfitReport func(context.Context, *gcfpb.GetGrossProfitRequest) (*gcfpb.GetGrossProfitResponse, error)
	Labels               report.ReportsLabels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
}

type PageData struct {
	types.PageData
	ContentTemplate   string
	Summary           []report.SummaryMetric
	Table             *types.TableConfig
	Filter            report.FilterState
	PeriodLabels      report.PeriodLabels
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
			reportURL = report.ReportsCostOfSalesURL
		}

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			sheetFilter := report.FilterState{
				ActivePreset:  period,
				StartDate:     startDateStr,
				EndDate:       endDateStr,
				PeriodPresets: report.DefaultPeriodPresets(pl, period),
			}
			return view.OK("report-filter-sheet", &report.FilterSheetData{
				Filter:       sheetFilter,
				PeriodLabels: pl,
				ReportURL:    reportURL,
			})
		}

		// Build proto request with date filtering
		req := &gcfpb.GetGrossProfitRequest{}
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
			start, _ := report.ParsePeriodPreset(period)
			s := start.Format("2006-01-02")
			req.StartDate = &s
		}
		if req.EndDate == nil {
			_, end := report.ParsePeriodPreset(period)
			e := end.Format("2006-01-02")
			req.EndDate = &e
		}

		var resp *gcfpb.GetGrossProfitResponse
		if deps.GetGrossProfitReport != nil {
			var err error
			resp, err = deps.GetGrossProfitReport(ctx, req)
			if err != nil {
				log.Printf("Failed to get cost of sales report: %v", err)
				resp = nil
			}
		}
		if resp == nil {
			resp = &gcfpb.GetGrossProfitResponse{
				LineItems: []*gcfpb.GrossProfitLineItem{},
				Summary:   &gcfpb.GrossProfitSummary{},
			}
		}

		s := resp.GetSummary()
		if s == nil {
			s = &gcfpb.GrossProfitSummary{}
		}

		cogsRatio := 0.0
		if s.GetNetRevenue() > 0 {
			cogsRatio = (float64(s.GetTotalCogs()) / float64(s.GetNetRevenue())) * 100
		}

		cogsCell := types.MoneyCell(float64(s.GetTotalCogs()), "PHP", true)
		revenueCell := types.MoneyCell(float64(s.GetNetRevenue()), "PHP", true)
		summary := []report.SummaryMetric{
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

		filter := report.FilterState{
			ActivePreset:  period,
			StartDate:     startDateStr,
			EndDate:       endDateStr,
			PeriodPresets: report.DefaultPeriodPresets(pl, period),
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
			ActiveFilterCount: report.ActiveFilterCount(filter),
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
