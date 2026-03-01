package net_profit

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
		startTS := start.Unix()
		endTS := end.Unix()
		req.StartDate = &startTS
		req.EndDate = &endTS

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

		netProfit := s.GetTotalGrossProfit() - totalExpenses
		netMargin := 0.0
		if s.GetNetRevenue() > 0 {
			netMargin = (netProfit / s.GetNetRevenue()) * 100
		}
		grossMargin := 0.0
		if s.GetNetRevenue() > 0 {
			grossMargin = (s.GetTotalGrossProfit() / s.GetNetRevenue()) * 100
		}

		// Summary bar
		netVariant := "success"
		if netProfit < 0 {
			netVariant = "danger"
		} else if netMargin < 10 {
			netVariant = "warning"
		}

		summary := []fycha.SummaryMetric{
			{Label: l.SummaryRevenue, Value: formatCurrency(s.GetNetRevenue())},
			{Label: l.SummaryGross, Value: formatCurrency(s.GetTotalGrossProfit())},
			{Label: l.SummaryExpenses, Value: formatCurrency(totalExpenses)},
			{Label: l.SummaryNetProfit, Value: formatCurrency(netProfit), Highlight: true, Variant: netVariant},
		}

		// P&L statement line items
		lineItems := []fycha.PLLineItem{
			{Label: l.Revenue, Value: formatCurrency(s.GetNetRevenue())},
			{Label: l.CostOfSales, Value: formatCurrency(s.GetTotalCogs())},
			{Label: l.GrossProfit, Value: formatCurrency(s.GetTotalGrossProfit()), IsTotal: true},
			{Label: l.GrossMargin, Value: fmt.Sprintf("%.1f%%", grossMargin)},
			{Label: l.Expenses, Value: formatCurrency(totalExpenses)},
			{Label: l.NetProfit, Value: formatCurrency(netProfit), IsTotal: true},
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
				ActiveNav:      "reports",
				ActiveSubNav:   "net-profit",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-dollar-sign",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "net-profit-content",
			Summary:         summary,
			LineItems:       lineItems,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: fycha.ActiveFilterCount(filter),
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

func formatCurrency(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
	whole := int64(amount)
	frac := int64((amount - float64(whole)) * 100 + 0.5)
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
