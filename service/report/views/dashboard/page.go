package dashboard

import (
	"context"
	"fmt"
	"log"
	"time"

	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"
	report "github.com/erniealice/fycha-golang/service/report"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// 20260521 Wave B P1.E.3 + P1.E.5 — DB replaced with typed closures into
// the service-driven gross/cashflow + domain-specific use cases.
type Deps struct {
	Routes               report.ReportsRoutes
	GetGrossProfitReport func(context.Context, *gcfpb.GetGrossProfitRequest) (*gcfpb.GetGrossProfitResponse, error)
	ListExpenses         func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
	Labels               report.ReportsLabels
	CommonLabels         pyeza.CommonLabels
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the reports dashboard view.
//
// Phase 1c refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// Real data flow via GetGrossProfitReport and ListExpenses preserved — no
// new aggregate methods.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard

		// Get this month's data for KPIs.
		start, end := report.ParsePeriodPreset("thisMonth")

		// Get gross profit data (contains revenue + COGS).
		req := &gcfpb.GetGrossProfitRequest{}
		startStr := start.Format("2006-01-02")
		endStr := end.Format("2006-01-02")
		req.StartDate = &startStr
		req.EndDate = &endStr

		var resp *gcfpb.GetGrossProfitResponse
		if deps.GetGrossProfitReport != nil {
			var err error
			resp, err = deps.GetGrossProfitReport(ctx, req)
			if err != nil {
				log.Printf("Failed to get dashboard report: %v", err)
				resp = nil
			}
		}
		if resp == nil {
			resp = &gcfpb.GetGrossProfitResponse{
				Summary: &gcfpb.GrossProfitSummary{},
			}
		}
		s := resp.GetSummary()
		if s == nil {
			s = &gcfpb.GrossProfitSummary{}
		}

		// Get expenses total.
		var expenseRecords []map[string]any
		if deps.ListExpenses != nil {
			var err error
			expenseRecords, err = deps.ListExpenses(ctx, &start, &end)
			if err != nil {
				log.Printf("Failed to list expenses for dashboard: %v", err)
			}
		}
		var totalExpenses float64
		for _, r := range expenseRecords {
			totalExpenses += toFloat64(r["total_amount"])
		}

		// Compute KPIs.
		netRevenue := float64(s.GetNetRevenue()) / 100.0
		totalCOGS := float64(s.GetTotalCogs()) / 100.0
		totalGrossProfit := float64(s.GetTotalGrossProfit()) / 100.0
		netProfit := totalGrossProfit - totalExpenses
		netMargin := 0.0
		if netRevenue > 0 {
			netMargin = (netProfit / netRevenue) * 100
		}

		// Stat-card colors track the existing summary-bar variant logic:
		// terracotta for revenue, sage for cost, navy for profit, amber for margin.
		netColor := "navy"
		if netProfit < 0 {
			netColor = "terracotta"
		}
		marginColor := "amber"
		if netProfit < 0 {
			marginColor = "terracotta"
		} else if netMargin >= 10 {
			marginColor = "sage"
		}

		revenueCell := types.MoneyCell(netRevenue, "PHP", false)
		cogsCell := types.MoneyCell(totalCOGS, "PHP", false)
		profitCell := types.MoneyCell(netProfit, "PHP", false)

		stats := []types.StatCardData{
			{Icon: "icon-trending-up", Value: revenueCell.Currency + " " + revenueCell.Value, Label: l.RevenueCard, Color: "terracotta", TestID: "report-stat-revenue"},
			{Icon: "icon-package", Value: cogsCell.Currency + " " + cogsCell.Value, Label: deps.Labels.CostOfSales.Title, Color: "sage", TestID: "report-stat-cogs"},
			{Icon: "icon-dollar-sign", Value: profitCell.Currency + " " + profitCell.Value, Label: l.NetProfitCard, Color: netColor, TestID: "report-stat-profit"},
			{Icon: "icon-percent", Value: fmt.Sprintf("%.1f%%", netMargin), Label: l.NetMarginCard, Color: marginColor, TestID: "report-stat-margin"},
		}

		// One widget per report type — Type=list with a single ActivityItem per
		// report. The block dispatches into activity-list.html which renders each
		// item as a clickable row.
		r := deps.Routes
		reportItems := []types.ActivityItem{
			{IconName: "icon-trending-up", IconVariant: "client", Title: deps.Labels.Revenue.Title, Description: l.RevenueDesc, Href: r.RevenueURL, TestID: "report-link-revenue"},
			{IconName: "icon-bar-chart", IconVariant: "quote", Title: deps.Labels.GrossProfit.Title, Description: l.GrossProfitDesc, Href: r.GrossProfitURL, TestID: "report-link-gross-profit"},
			{IconName: "icon-package", IconVariant: "award", Title: deps.Labels.CostOfSales.Title, Description: l.CostOfSalesDesc, Href: r.CostOfSalesURL, TestID: "report-link-cost-of-sales"},
			{IconName: "icon-file-minus", IconVariant: "integration", Title: deps.Labels.Expenses.Title, Description: l.ExpensesDesc, Href: r.ExpensesURL, TestID: "report-link-expenses"},
			{IconName: "icon-dollar-sign", IconVariant: "client", Title: deps.Labels.NetProfit.Title, Description: l.NetProfitDesc, Href: r.NetProfitURL, TestID: "report-link-net-profit"},
		}

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-trending-up", Label: deps.Labels.Revenue.Title, Href: r.RevenueURL, Variant: "primary", TestID: "report-action-revenue"},
				{Icon: "icon-clock", Label: deps.Labels.ReceivablesAging.PageTitle, Href: r.ReceivablesAgingReportURL, TestID: "report-action-receivables"},
				{Icon: "icon-credit-card", Label: deps.Labels.CollectionSummary.PageTitle, Href: r.CollectionSummaryReportURL, TestID: "report-action-collections"},
				{Icon: "icon-file-minus", Label: deps.Labels.Expenses.Title, Href: r.ExpenditureReportURL, TestID: "report-action-expenditure"},
			},
			Stats: stats,
			Widgets: []types.DashboardWidget{
				{
					ID: "reports", Title: l.Subtitle, Type: "list", Span: 3,
					ListItems: reportItems,
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "report",
				ActiveSubNav:   "dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-pie-chart",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "reports-dashboard-content",
			Dashboard:       dash,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-dashboard"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
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
