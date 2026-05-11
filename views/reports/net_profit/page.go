package net_profit

import (
	"context"
	"fmt"
	"log"
	"time"

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
	LineItems         []fycha.PLLineItem
	Filter            fycha.FilterState
	PeriodLabels      fycha.PeriodLabels
	ReportURL         string
	ActiveFilterCount int
}

func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.NetProfit
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
			reportURL = fycha.ReportsNetProfitURL
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

		// Resolve dates
		start, end := fycha.ParsePeriodPreset(period)
		if period == "custom" {
			if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
				start = t
			}
			if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
				end = t
			}
		}

		// Get gross profit data (contains revenue + COGS)
		req := &reportpb.GrossProfitReportRequest{}
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")
		req.StartDate = &startStr
		req.EndDate = &endStr

		resp, err := deps.DB.GetGrossProfitReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get profit report: %v", err)
			resp = &reportpb.GrossProfitReportResponse{
				Summary: &reportpb.GrossProfitSummary{},
			}
		}
		s := resp.GetSummary()
		if s == nil {
			s = &reportpb.GrossProfitSummary{}
		}

		// Get expenses total
		expenseRecords, err := deps.DB.ListExpenses(ctx, &start, &end)
		if err != nil {
			log.Printf("Failed to list expenses: %v", err)
		}
		var totalExpenses float64
		for _, r := range expenseRecords {
			totalExpenses += toFloat64(r["total_amount"])
		}

		grossProfitF := float64(s.GetTotalGrossProfit()) / 100.0
		netRevenueF := float64(s.GetNetRevenue()) / 100.0
		totalCogsF := float64(s.GetTotalCogs()) / 100.0
		netProfit := grossProfitF - totalExpenses
		netMargin := 0.0
		if netRevenueF > 0 {
			netMargin = (netProfit / netRevenueF) * 100
		}
		grossMargin := 0.0
		if netRevenueF > 0 {
			grossMargin = (grossProfitF / netRevenueF) * 100
		}

		// Summary bar
		netVariant := "success"
		if netProfit < 0 {
			netVariant = "danger"
		} else if netMargin < 10 {
			netVariant = "warning"
		}

		c0 := types.MoneyCell(netRevenueF, "PHP", true)
		c1 := types.MoneyCell(grossProfitF, "PHP", true)
		c2 := types.MoneyCell(totalExpenses, "PHP", true)
		c3 := types.MoneyCell(netProfit, "PHP", true)
		summary := []fycha.SummaryMetric{
			{Label: l.SummaryRevenue, Value: c0.Currency + " " + c0.Value},
			{Label: l.SummaryGross, Value: c1.Currency + " " + c1.Value},
			{Label: l.SummaryExpenses, Value: c2.Currency + " " + c2.Value},
			{Label: l.SummaryNetProfit, Value: c3.Currency + " " + c3.Value, Highlight: true, Variant: netVariant},
		}

		// P&L statement line items
		li0 := types.MoneyCell(netRevenueF, "PHP", true)
		li1 := types.MoneyCell(totalCogsF, "PHP", true)
		li2 := types.MoneyCell(grossProfitF, "PHP", true)
		li3 := types.MoneyCell(totalExpenses, "PHP", true)
		li4 := types.MoneyCell(netProfit, "PHP", true)
		lineItems := []fycha.PLLineItem{
			{Label: l.Revenue, Value: li0.Currency + " " + li0.Value},
			{Label: l.CostOfSales, Value: li1.Currency + " " + li1.Value},
			{Label: l.GrossProfit, Value: li2.Currency + " " + li2.Value, IsTotal: true},
			{Label: l.GrossMargin, Value: fmt.Sprintf("%.1f%%", grossMargin)},
			{Label: l.Expenses, Value: li3.Currency + " " + li3.Value},
			{Label: l.NetProfit, Value: li4.Currency + " " + li4.Value, IsTotal: true},
			{Label: l.NetMargin, Value: fmt.Sprintf("%.1f%%", netMargin), Variant: netVariant},
		}

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
				ActiveSubNav:   "net-profit",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-dollar-sign",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:   "net-profit-content",
			Summary:           summary,
			LineItems:         lineItems,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: fycha.ActiveFilterCount(filter),
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-net-profit"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("net-profit-content", pageData)
		}
		return view.OK("net-profit", pageData)
	})
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int64:
		return float64(n)
	case int:
		return float64(n)
	default:
		return 0
	}
}
