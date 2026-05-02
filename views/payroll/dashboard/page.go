// Package dashboard implements the read-only Payroll live dashboard view
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

	fycha "github.com/erniealice/fycha-golang"

	payrollremittancepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_remittance"
	payrollrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_run"
)

// Request mirrors espyna's request.
type Request struct {
	WorkspaceID string
}

// Response mirrors espyna's response.
type Response struct {
	CurrentRunStatus    string
	EmployeesInCurrent  int32
	TotalGrossMTD       int64
	RemittancesDue30Cnt int64
	LatestRun           *payrollrunpb.PayrollRun
	RecentRuns          []*payrollrunpb.PayrollRun
	UpcomingDeadlines   []*payrollremittancepb.PayrollRemittance
	GrossTrendLabels    []string
	GrossTrendValues    []float64
}

// Deps holds view dependencies.
type Deps struct {
	Routes               fycha.PayrollRunRoutes
	RemittanceRoutes     fycha.PayrollRemittanceRoutes
	SettingsRoutes       fycha.PayrollSettingsRoutes
	Labels               fycha.PayrollLabels
	CommonLabels         pyeza.CommonLabels
	GetDashboardPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is the dashboard template payload.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the payroll dashboard view.
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

		// 12-month gross-pay trend.
		labels := resp.GrossTrendLabels
		values := resp.GrossTrendValues
		if len(labels) == 0 {
			labels = []string{"Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec", "Jan", "Feb", "Mar", "Apr", "May"}
			values = make([]float64, 12)
		}
		trend := &types.ChartData{
			Labels: labels,
			Series: []types.ChartSeries{{
				Name:   l.GrossPayByMonth,
				Values: values,
				Color:  "sage",
			}},
			Currency: "PHP",
			YAxis:    l.AxisGross,
		}
		trend.AutoScale()

		// Stats — current run status (text label) + counters.
		runStatus := resp.CurrentRunStatus
		if runStatus == "" {
			runStatus = l.NoRunYet
		}

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-users",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewRun, Href: deps.Routes.ListURL, Variant: "primary", TestID: "payroll-action-new-run"},
				{Icon: "icon-zap", Label: l.QuickProcessRun, Href: deps.Routes.ListURL, TestID: "payroll-action-process"},
				{Icon: "icon-file", Label: l.QuickFileRemittance, Href: deps.RemittanceRoutes.ListURL, TestID: "payroll-action-remittance"},
				{Icon: "icon-settings", Label: l.QuickPayPeriodSettings, Href: deps.SettingsRoutes.PayPeriodsURL, TestID: "payroll-action-pay-periods"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-activity", Value: humanizeStatus(runStatus), Label: l.CurrentRunStatus, Color: "navy", TestID: "payroll-stat-status"},
				{Icon: "icon-users", Value: fmt.Sprintf("%d", resp.EmployeesInCurrent), Label: l.EmployeesInCurrent, Color: "sage", TestID: "payroll-stat-employees"},
				{Icon: "icon-dollar-sign", Value: formatPHP(resp.TotalGrossMTD), Label: l.TotalGrossMTD, Color: "terracotta", TestID: "payroll-stat-gross-mtd"},
				{Icon: "icon-clock", Value: fmt.Sprintf("%d", resp.RemittancesDue30Cnt), Label: l.RemittancesDue, Color: "amber", TestID: "payroll-stat-remittances"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "gross-trend", Title: l.GrossPayByMonth,
					Type: "chart", ChartKind: "bar",
					ChartData: trend, Span: 2,
				},
				{
					ID: "recent-runs", Title: l.RecentRuns, Type: "custom", Span: 2,
					Custom: buildRecentRunsHTML(resp.RecentRuns, l, deps.Routes),
					EmptyState: &types.EmptyStateData{
						Icon: "icon-users", Title: l.RecentRuns, Desc: l.NoRecentRuns,
					},
				},
				{
					ID: "remittances", Title: l.UpcomingRemittances, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.RemittanceRoutes.ListURL},
					},
					ListItems: buildRemittanceList(resp.UpcomingDeadlines, l),
					EmptyState: &types.EmptyStateData{
						Icon: "icon-clock", Title: l.UpcomingRemittances, Desc: l.NoUpcomingDeadlines,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "payroll",
				ActiveSubNav:   "dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "payroll-dashboard-content",
			Dashboard:       dash,
		}
		return view.OK("payroll-dashboard", pageData)
	})
}

func buildRecentRunsHTML(runs []*payrollrunpb.PayrollRun, l fycha.PayrollDashboardLabels, _ fycha.PayrollRunRoutes) template.HTML {
	if len(runs) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<table class="dashboard-mini-table"><thead><tr>`)
	b.WriteString(`<th>#</th><th>` + template.HTMLEscapeString(l.GrossPayByMonth) + `</th><th class="numeric">` + template.HTMLEscapeString(l.AxisGross) + `</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	for _, r := range runs {
		b.WriteString(`<tr data-testid="payroll-table-row-`)
		b.WriteString(template.HTMLEscapeString(r.GetId()))
		b.WriteString(`">`)
		b.WriteString(`<td>` + template.HTMLEscapeString(r.GetRunNumber()) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(r.GetPayPeriodEnd()) + `</td>`)
		b.WriteString(`<td class="numeric">` + formatPHP(r.GetTotalGross()) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return template.HTML(b.String()) //nolint:gosec
}

func buildRemittanceList(items []*payrollremittancepb.PayrollRemittance, l fycha.PayrollDashboardLabels) []types.ActivityItem {
	if len(items) == 0 {
		return nil
	}
	out := make([]types.ActivityItem, 0, len(items))
	for i, r := range items {
		out = append(out, types.ActivityItem{
			IconName:    "icon-file",
			IconVariant: "quote",
			Title:       r.GetRemittanceType().String(),
			Description: formatPHP(r.GetAmount()),
			Time:        r.GetDueDate(),
			TestID:      fmt.Sprintf("payroll-list-item-%d", i),
		})
	}
	return out
}

func humanizeStatus(s string) string {
	if s == "" {
		return ""
	}
	// E.g. PAYROLL_RUN_STATUS_DRAFT → "Draft"
	parts := strings.Split(s, "_")
	if len(parts) == 0 {
		return s
	}
	last := parts[len(parts)-1]
	if last == "" || last == "UNSPECIFIED" {
		return s
	}
	return strings.Title(strings.ToLower(last)) //nolint:staticcheck // case-stable; not user input
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
