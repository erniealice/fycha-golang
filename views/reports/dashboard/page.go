package dashboard

import (
	"context"
	"fmt"
	"log"
	"strconv"

	fycha "github.com/erniealice/fycha-golang"
	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

type Deps struct {
	Routes       fycha.ReportsRoutes
	DB           fycha.DataSource
	Labels       fycha.ReportsLabels
	CommonLabels pyeza.CommonLabels
}

// ReportCard holds navigation card data for the dashboard.
type ReportCard struct {
	Title       string
	Description string
	Icon        string
	URL         string
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Summary         []fycha.SummaryMetric
	ReportCards     []ReportCard
	Labels          fycha.DashboardLabels
}

func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard

		// Get this month's data for KPIs
		start, end := fycha.ParsePeriodPreset("thisMonth")

		// Get gross profit data (contains revenue + COGS)
		req := &reportpb.GrossProfitReportRequest{}
		startTS := start.Unix()
		endTS := end.Unix()
		req.StartDate = &startTS
		req.EndDate = &endTS

		resp, err := deps.DB.GetGrossProfitReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get dashboard report: %v", err)
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
			log.Printf("Failed to list expenses for dashboard: %v", err)
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

		// KPI summary
		netVariant := "success"
		if netProfit < 0 {
			netVariant = "danger"
		} else if netMargin < 10 {
			netVariant = "warning"
		}

		summary := []fycha.SummaryMetric{
			{Label: l.RevenueCard, Value: formatCurrency(s.GetNetRevenue())},
			{Label: l.ExpensesCard, Value: formatCurrency(totalExpenses)},
			{Label: l.NetProfitCard, Value: formatCurrency(netProfit), Highlight: true, Variant: netVariant},
			{Label: l.NetMarginCard, Value: fmt.Sprintf("%.1f%%", netMargin), Variant: netVariant},
		}

		// Navigation cards
		r := deps.Routes
		reportCards := []ReportCard{
			{Title: deps.Labels.Revenue.Title, Description: l.RevenueDesc, Icon: "icon-trending-up", URL: r.RevenueURL},
			{Title: deps.Labels.GrossProfit.Title, Description: l.GrossProfitDesc, Icon: "icon-bar-chart", URL: r.GrossProfitURL},
			{Title: deps.Labels.CostOfSales.Title, Description: l.CostOfSalesDesc, Icon: "icon-package", URL: r.CostOfSalesURL},
			{Title: deps.Labels.Expenses.Title, Description: l.ExpensesDesc, Icon: "icon-file-minus", URL: r.ExpensesURL},
			{Title: deps.Labels.NetProfit.Title, Description: l.NetProfitDesc, Icon: "icon-dollar-sign", URL: r.NetProfitURL},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "reports",
				ActiveSubNav: "dashboard",
				HeaderTitle:  l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:   "icon-pie-chart",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "reports-dashboard-content",
			Summary:         summary,
			ReportCards:     reportCards,
			Labels:          l,
		}

		if viewCtx.IsHTMX {
			return view.OK("reports-dashboard-content", pageData)
		}
		return view.OK("reports-dashboard", pageData)
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
