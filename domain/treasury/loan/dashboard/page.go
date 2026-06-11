// Package dashboard implements the read-only Loan live dashboard view
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
package dashboard

import (
	"context"
	"fmt"
	"html/template"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	loan "github.com/erniealice/fycha-golang/domain/treasury/loan"
	loanpayment "github.com/erniealice/fycha-golang/domain/treasury/loan_payment"

	loanpaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan_payment"
)

// LoanRow mirrors espyna.dashboard.LoanSlice (top-loans widget).
type LoanRow struct {
	ID               string
	LoanNumber       string
	LenderName       string
	RemainingBalance int64
	PrincipalAmount  int64
	Status           string
}

// Request mirrors espyna's request.
type Request struct {
	WorkspaceID string
}

// Response mirrors espyna's response.
type Response struct {
	TotalOutstanding int64
	InterestYTD      int64
	PaymentsDue30    int64
	DefaultedCount   int64
	TrendLabels      []string
	TrendValues      []float64
	TopLoans         []LoanRow
	RecentPayments   []*loanpaymentpb.LoanPayment
}

// Deps holds view dependencies.
type Deps struct {
	Routes               loan.Routes
	PaymentRoutes        loanpayment.Routes
	Labels               loan.Labels
	CommonLabels         pyeza.CommonLabels
	GetDashboardPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is the dashboard template payload.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the loan dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard

		var resp *Response
		if deps.GetDashboardPageData != nil {
			r, err := deps.GetDashboardPageData(ctx, &Request{WorkspaceID: ""})
			if err == nil && r != nil {
				resp = r
			}
		}
		if resp == nil {
			resp = &Response{}
		}

		// 6-month outstanding-principal trend.
		labels := resp.TrendLabels
		values := resp.TrendValues
		if len(labels) == 0 {
			labels = []string{"Dec", "Jan", "Feb", "Mar", "Apr", "May"}
			values = []float64{0, 0, 0, 0, 0, 0}
		}
		trend := &types.ChartData{
			Labels: labels,
			Series: []types.ChartSeries{{
				Name:   l.OutstandingTrend,
				Values: values,
				Color:  "terracotta",
			}},
			Currency: "PHP",
			YAxis:    l.AxisAmount,
		}
		trend.AutoScale()

		// Recent payments list.
		recent := buildRecentPaymentsList(resp.RecentPayments, l)

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-credit-card",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewLoan, Href: deps.Routes.AddURL, Variant: "primary", TestID: "loan-action-new"},
				{Icon: "icon-dollar-sign", Label: l.QuickRecordPay, Href: deps.PaymentRoutes.AddURL, TestID: "loan-action-record-pay"},
				{Icon: "icon-list", Label: l.QuickAmortization, Href: deps.Routes.AmortizationURL, TestID: "loan-action-amortization"},
				{Icon: "icon-calendar", Label: l.QuickLoanCalendar, Href: deps.Routes.ListURL, TestID: "loan-action-calendar"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-trending-up", Value: formatPHP(resp.TotalOutstanding), Label: l.TotalOutstanding, Color: "terracotta", TestID: "loan-stat-outstanding"},
				{Icon: "icon-percent", Value: formatPHP(resp.InterestYTD), Label: l.InterestYTD, Color: "sage", TestID: "loan-stat-interest"},
				{Icon: "icon-clock", Value: formatPHP(resp.PaymentsDue30), Label: l.PaymentsDue30, Color: "amber", TestID: "loan-stat-due"},
				{Icon: "icon-alert-triangle", Value: fmt.Sprintf("%d", resp.DefaultedCount), Label: l.DefaultedCount, Color: "navy", TestID: "loan-stat-defaulted"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "outstanding-trend", Title: l.OutstandingTrend,
					Type: "chart", ChartKind: "line",
					ChartData: trend, Span: 2,
				},
				{
					ID: "top-loans", Title: l.TopLoans, Type: "custom", Span: 2,
					Custom: buildTopLoansHTML(resp.TopLoans, l, deps.Routes),
					EmptyState: &types.EmptyStateData{
						Icon: "icon-credit-card", Title: l.TopLoans, Desc: l.NoRecentPayments,
					},
				},
				{
					ID: "recent-payments", Title: l.RecentPayments, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.PaymentRoutes.ListURL},
					},
					ListItems: recent,
					EmptyState: &types.EmptyStateData{
						Icon: "icon-activity", Title: l.RecentPayments, Desc: l.NoRecentPayments,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "loan",
				ActiveSubNav:   "dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "loan-dashboard-content",
			Dashboard:       dash,
		}
		return view.OK("loan-dashboard", pageData)
	})
}

func buildRecentPaymentsList(payments []*loanpaymentpb.LoanPayment, l loan.DashboardLabels) []types.ActivityItem {
	if len(payments) == 0 {
		return nil
	}
	items := make([]types.ActivityItem, 0, len(payments))
	for i, p := range payments {
		items = append(items, types.ActivityItem{
			IconName:    "icon-dollar-sign",
			IconVariant: "client",
			Title:       p.GetPaymentNumber(),
			Description: fmt.Sprintf("%s · %s", l.OutstandingTrend, formatPHP(p.GetTotalAmount())),
			Time:        p.GetPaymentDate(),
			TestID:      fmt.Sprintf("loan-list-item-%d", i),
		})
	}
	return items
}

func buildTopLoansHTML(rows []LoanRow, l loan.DashboardLabels, routes loan.Routes) template.HTML {
	if len(rows) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<table class="dashboard-mini-table"><thead><tr>`)
	b.WriteString(`<th>#</th><th>` + template.HTMLEscapeString(l.TopLoans) + `</th><th class="numeric">` + template.HTMLEscapeString(l.AxisAmount) + `</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	for _, r := range rows {
		b.WriteString(`<tr data-testid="loan-table-row-`)
		b.WriteString(template.HTMLEscapeString(r.ID))
		b.WriteString(`">`)
		b.WriteString(`<td>` + template.HTMLEscapeString(r.LoanNumber) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(r.LenderName) + `</td>`)
		b.WriteString(`<td class="numeric">` + formatPHP(r.RemainingBalance) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return template.HTML(b.String()) //nolint:gosec
}

func formatPHP(centavos int64) string {
	pesos := float64(centavos) / 100.0
	return fmt.Sprintf("PHP %s", humanizePesos(pesos))
}

func humanizePesos(v float64) string {
	negative := false
	if v < 0 {
		v = -v
		negative = true
	}
	int64Part := int64(v)
	fracPart := int64((v-float64(int64Part))*100 + 0.5)
	intStr := fmt.Sprintf("%d", int64Part)
	n := len(intStr)
	if n > 3 {
		var parts []string
		for n > 3 {
			parts = append([]string{intStr[n-3:]}, parts...)
			n -= 3
		}
		parts = append([]string{intStr[:n]}, parts...)
		intStr = strings.Join(parts, ",")
	}
	out := fmt.Sprintf("%s.%02d", intStr, fracPart)
	if negative {
		out = "-" + out
	}
	return out
}
