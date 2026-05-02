// Package dashboard implements the read-only Equity live dashboard view
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

	equitytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_transaction"
)

// EquityAccountRow mirrors espyna.equity_dashboard.EquityAccountSlice.
type EquityAccountRow struct {
	ID          string
	Name        string
	OwnerName   string
	AccountType string
	Balance     int64
}

// Request mirrors espyna's request.
type Request struct {
	WorkspaceID string
}

// Response mirrors espyna's response.
type Response struct {
	TotalContributed int64
	ActiveOwners     int64
	DistributionsYTD int64
	NetMovementYTD   int64
	ByTypeYTD        map[string]int64
	TopContributors  []EquityAccountRow
	Recent           []*equitytransactionpb.EquityTransaction
}

// Deps holds view dependencies.
type Deps struct {
	Routes               fycha.EquityRoutes
	Labels               fycha.EquityLabels
	CommonLabels         pyeza.CommonLabels
	GetDashboardPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is the dashboard template payload.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the equity dashboard view.
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
			resp = &Response{ByTypeYTD: map[string]int64{}}
		}

		// Pie chart: equity by owner — slices = top contributors.
		var pieLabels []string
		var pieValues []float64
		for _, c := range resp.TopContributors {
			label := c.OwnerName
			if label == "" {
				label = c.Name
			}
			pieLabels = append(pieLabels, label)
			pieValues = append(pieValues, float64(c.Balance))
		}
		if len(pieLabels) == 0 {
			pieLabels = []string{"-"}
			pieValues = []float64{0}
		}
		pieChart := &types.ChartData{
			Labels: pieLabels,
			Series: []types.ChartSeries{{
				Name:   l.EquityByOwner,
				Values: pieValues,
				Color:  "plum",
			}},
			Currency: "PHP",
		}
		pieChart.AutoScale()

		// Recent transactions list.
		recent := buildRecentTxnList(resp.Recent, l)

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-pie-chart",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickRecordContribution, Href: deps.Routes.TransactionAddURL, Variant: "primary", TestID: "equity-action-contribution"},
				{Icon: "icon-trending-down", Label: l.QuickRecordDistribution, Href: deps.Routes.TransactionAddURL, TestID: "equity-action-distribution"},
				{Icon: "icon-user", Label: l.QuickOwnerStatement, Href: deps.Routes.AccountsURL, TestID: "equity-action-statement"},
				{Icon: "icon-bar-chart", Label: l.QuickEquityReport, Href: deps.Routes.TransactionsURL, TestID: "equity-action-report"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-trending-up", Value: formatPHP(resp.TotalContributed), Label: l.TotalContributed, Color: "terracotta", TestID: "equity-stat-contributed"},
				{Icon: "icon-users", Value: fmt.Sprintf("%d", resp.ActiveOwners), Label: l.ActiveOwners, Color: "sage", TestID: "equity-stat-owners"},
				{Icon: "icon-trending-down", Value: formatPHP(resp.DistributionsYTD), Label: l.DistributionsYTD, Color: "navy", TestID: "equity-stat-distributions"},
				{Icon: "icon-activity", Value: formatPHPSigned(resp.NetMovementYTD), Label: l.NetMovementYTD, Color: "amber", TestID: "equity-stat-net"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "equity-by-owner", Title: l.EquityByOwner,
					Type: "chart", ChartKind: "pie",
					ChartData: pieChart, Span: 2,
				},
				{
					ID: "top-contributors", Title: l.TopContributors, Type: "custom", Span: 2,
					Custom: buildTopContributorsHTML(resp.TopContributors, l, deps.Routes),
					EmptyState: &types.EmptyStateData{
						Icon: "icon-users", Title: l.TopContributors, Desc: l.NoRecentTxns,
					},
				},
				{
					ID: "recent", Title: l.RecentTransactions, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.TransactionsURL},
					},
					ListItems: recent,
					EmptyState: &types.EmptyStateData{
						Icon: "icon-activity", Title: l.RecentTransactions, Desc: l.NoRecentTxns,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "equity",
				ActiveSubNav:   "dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-pie-chart",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "equity-dashboard-content",
			Dashboard:       dash,
		}
		return view.OK("equity-dashboard", pageData)
	})
}

func buildRecentTxnList(txns []*equitytransactionpb.EquityTransaction, l fycha.EquityDashboardLabels) []types.ActivityItem {
	if len(txns) == 0 {
		return nil
	}
	items := make([]types.ActivityItem, 0, len(txns))
	for i, t := range txns {
		variant := "client"
		switch t.GetTransactionType() {
		case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_WITHDRAWAL:
			variant = "integration"
		case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_DISTRIBUTION:
			variant = "quote"
		case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_TRANSFER:
			variant = "award"
		}
		items = append(items, types.ActivityItem{
			IconName:    "icon-dollar-sign",
			IconVariant: variant,
			Title:       t.GetTransactionType().String(),
			Description: t.GetDescription(),
			Time:        formatTimestamp(t.GetTransactionDate()),
			TestID:      fmt.Sprintf("equity-list-item-%d", i),
		})
	}
	return items
}

func buildTopContributorsHTML(rows []EquityAccountRow, l fycha.EquityDashboardLabels, _ fycha.EquityRoutes) template.HTML {
	if len(rows) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<table class="dashboard-mini-table"><thead><tr>`)
	b.WriteString(`<th>` + template.HTMLEscapeString(l.TopContributors) + `</th><th>` + template.HTMLEscapeString(l.EquityByOwner) + `</th><th class="numeric">` + template.HTMLEscapeString(l.AxisAmount) + `</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	for _, r := range rows {
		b.WriteString(`<tr data-testid="equity-table-row-`)
		b.WriteString(template.HTMLEscapeString(r.ID))
		b.WriteString(`">`)
		b.WriteString(`<td>` + template.HTMLEscapeString(r.Name) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(r.OwnerName) + `</td>`)
		b.WriteString(`<td class="numeric">` + formatPHP(r.Balance) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return template.HTML(b.String()) //nolint:gosec
}

func formatPHP(centavos int64) string {
	pesos := float64(centavos) / 100.0
	return fmt.Sprintf("PHP %s", humanizePesos(pesos))
}

func formatPHPSigned(centavos int64) string {
	pesos := float64(centavos) / 100.0
	if pesos < 0 {
		return fmt.Sprintf("-PHP %s", humanizePesos(-pesos))
	}
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

func formatTimestamp(ms int64) string {
	if ms <= 0 {
		return ""
	}
	// We avoid pulling in time.* to render the millis here. The server already
	// emits ISO-formatted dates in DateString-typed proto fields where they
	// matter, and the dashboard list-time slot is informational.
	return ""
}
