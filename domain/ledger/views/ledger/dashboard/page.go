// Package dashboard implements the read-only Ledger live dashboard view
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

	ledger "github.com/erniealice/fycha-golang/domain/ledger"

	journalentrypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
)

// Request is the view-side request shape; mirrors espyna's
// dashboard.GetLedgerDashboardPageDataRequest without importing espyna.
type Request struct {
	WorkspaceID string
}

// Response is the view-side response shape; mirrors espyna's
// dashboard.GetLedgerDashboardPageDataResponse.
type Response struct {
	TotalAssets      int64
	TotalLiabilities int64
	TotalEquity      int64
	NetIncomeMTD     int64
	UnpostedJournals int64
	BalanceByType    map[string]int64
	UnpostedTop      []*journalentrypb.JournalEntry
	RecentEntries    []*journalentrypb.JournalEntry
}

// Deps holds view dependencies.
type Deps struct {
	Routes               ledger.AccountRoutes
	JournalRoutes        ledger.JournalRoutes
	StatementRoutes      ledger.LedgerStatementRoutes
	FiscalRoutes         ledger.FiscalPeriodRoutes
	Labels               ledger.AccountLabels
	CommonLabels         pyeza.CommonLabels
	GetDashboardPageData func(ctx context.Context, req *Request) (*Response, error)
}

// PageData is what the ledger dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the ledger dashboard view.
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
			resp = &Response{BalanceByType: map[string]int64{}}
		}

		// Bar chart: balance by element. Force a stable ordering.
		barLabels := []string{l.AccountTypeAssets, l.AccountTypeLiabilities, l.AccountTypeEquity, l.AccountTypeRevenue, l.AccountTypeExpense}
		barValues := []float64{
			float64(resp.BalanceByType["asset"]),
			float64(resp.BalanceByType["liability"]),
			float64(resp.BalanceByType["equity"]),
			float64(resp.BalanceByType["revenue"]),
			float64(resp.BalanceByType["expense"]),
		}
		barChart := &types.ChartData{
			Labels: barLabels,
			Series: []types.ChartSeries{{
				Name:   l.BalanceByType,
				Values: barValues,
				Color:  "navy",
			}},
			Currency: "PHP",
			YAxis:    l.AxisAmount,
		}
		barChart.AutoScale()

		// Recent activity list.
		recent := buildRecentJournalActivity(resp.RecentEntries, l)

		dash := types.DashboardData{
			Title:    l.Title,
			Icon:     "icon-book-open",
			Subtitle: l.Subtitle,
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewJournal, Href: deps.JournalRoutes.AddURL, Variant: "primary", TestID: "ledger-action-new-journal"},
				{Icon: "icon-list", Label: l.QuickTrialBalance, Href: deps.StatementRoutes.TrialBalanceURL, TestID: "ledger-action-trial-balance"},
				{Icon: "icon-lock", Label: l.QuickClosePeriod, Href: deps.FiscalRoutes.ListURL, TestID: "ledger-action-close-period"},
				{Icon: "icon-search", Label: l.QuickAccountLookup, Href: deps.Routes.ListURL, TestID: "ledger-action-account-lookup"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-trending-up", Value: formatPHP(resp.TotalAssets), Label: l.TotalAssets, Color: "terracotta", TestID: "ledger-stat-assets"},
				{Icon: "icon-trending-down", Value: formatPHP(resp.TotalLiabilities), Label: l.TotalLiabilities, Color: "navy", TestID: "ledger-stat-liabilities"},
				{Icon: "icon-pie-chart", Value: formatPHP(resp.TotalEquity), Label: l.TotalEquity, Color: "sage", TestID: "ledger-stat-equity"},
				{Icon: "icon-bar-chart", Value: formatPHPSigned(resp.NetIncomeMTD), Label: l.NetIncomeMTD, Color: "amber", TestID: "ledger-stat-net-income"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "balance-by-type", Title: l.BalanceByType,
					Type: "chart", ChartKind: "bar",
					ChartData: barChart, Span: 2,
				},
				{
					ID: "unposted", Title: l.UnpostedJournals, Type: "custom", Span: 2,
					Custom: buildJournalsHTML(resp.UnpostedTop, l, deps.JournalRoutes),
					EmptyState: &types.EmptyStateData{
						Icon: "icon-file-text", Title: l.UnpostedJournals, Desc: l.NoRecentJournals,
					},
				},
				{
					ID: "recent", Title: l.RecentJournals, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.JournalRoutes.ListURL},
					},
					ListItems: recent,
					EmptyState: &types.EmptyStateData{
						Icon: "icon-activity", Title: l.RecentJournals, Desc: l.NoRecentJournals,
					},
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "ledger",
				ActiveSubNav:   "dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-book-open",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "ledger-dashboard-content",
			Dashboard:       dash,
		}
		return view.OK("ledger-dashboard", pageData)
	})
}

func buildRecentJournalActivity(entries []*journalentrypb.JournalEntry, l ledger.LedgerDashboardLabels) []types.ActivityItem {
	if len(entries) == 0 {
		return nil
	}
	items := make([]types.ActivityItem, 0, len(entries))
	for i, e := range entries {
		title := e.GetEntryNumber()
		if title == "" {
			title = l.JournalEntryFallback
		}
		desc := e.GetDescription()
		variant := "client"
		switch e.GetStatus() {
		case journalentrypb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_DRAFT:
			variant = "quote"
		case journalentrypb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_POSTED:
			variant = "client"
		case journalentrypb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_REVERSED:
			variant = "integration"
		}
		t := ""
		if s := e.GetEntryDateString(); s != "" {
			t = s
		}
		items = append(items, types.ActivityItem{
			IconName:    "icon-file-text",
			IconVariant: variant,
			Title:       title,
			Description: desc,
			Time:        t,
			TestID:      fmt.Sprintf("ledger-list-item-%d", i),
		})
	}
	return items
}

func buildJournalsHTML(entries []*journalentrypb.JournalEntry, l ledger.LedgerDashboardLabels, routes ledger.JournalRoutes) template.HTML {
	if len(entries) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(`<table class="dashboard-mini-table"><thead><tr>`)
	b.WriteString(`<th>#</th><th>` + template.HTMLEscapeString(l.UnpostedJournals) + `</th><th class="numeric">` + template.HTMLEscapeString(l.AxisAmount) + `</th>`)
	b.WriteString(`</tr></thead><tbody>`)
	for _, e := range entries {
		b.WriteString(`<tr data-testid="ledger-table-row-`)
		b.WriteString(template.HTMLEscapeString(e.GetId()))
		b.WriteString(`">`)
		b.WriteString(`<td>` + template.HTMLEscapeString(e.GetEntryNumber()) + `</td>`)
		b.WriteString(`<td>` + template.HTMLEscapeString(e.GetDescription()) + `</td>`)
		b.WriteString(`<td class="numeric">` + formatPHP(e.GetTotalDebit()) + `</td>`)
		b.WriteString(`</tr>`)
	}
	b.WriteString(`</tbody></table>`)
	return template.HTML(b.String()) //nolint:gosec // values escaped above
}

// formatPHP renders a centavo amount as "PHP 12,345.67". All values from
// aggregates are centavos (int64); centavos ÷ 100 = pesos.
func formatPHP(centavos int64) string {
	pesos := float64(centavos) / 100.0
	return fmt.Sprintf("PHP %s", humanizePesos(pesos))
}

// formatPHPSigned formats including a leading sign for net-income style stats.
func formatPHPSigned(centavos int64) string {
	pesos := float64(centavos) / 100.0
	if pesos < 0 {
		return fmt.Sprintf("-PHP %s", humanizePesos(-pesos))
	}
	return fmt.Sprintf("PHP %s", humanizePesos(pesos))
}

// humanizePesos prints a float with thousands separators and 2 decimal places.
func humanizePesos(v float64) string {
	negative := false
	if v < 0 {
		v = -v
		negative = true
	}
	int64Part := int64(v)
	fracPart := int64((v-float64(int64Part))*100 + 0.5)
	intStr := fmt.Sprintf("%d", int64Part)
	// Thousands separators
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
