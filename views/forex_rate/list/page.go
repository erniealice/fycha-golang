// Package list provides the read-only Forex Rate list view.
package list

import (
	"context"
	"fmt"
	"log"

	forexratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/finance/forex_rate"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the forex rate list page.
type Deps struct {
	Routes       fycha.ForexRateRoutes
	Labels       fycha.ForexRateLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Forex rate use cases (read-only)
	ListForexRates func(ctx context.Context, req *forexratepb.ListForexRatesRequest) (*forexratepb.ListForexRatesResponse, error)
}

// PageData holds the data for the forex rate list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveStatus    string
	StatusTabs      []pyeza.TabItem
	Table           *types.TableConfig
	Labels          fycha.ForexRateLabels
}

// ForexRateRow is the view-model for a single forex rate row.
type ForexRateRow struct {
	ID             string
	FromCurrency   string
	ToCurrency     string
	RateMicroUnits string
	Source         string
	EffectiveFrom  string
	Status         string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the forex rate list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("forex_rate", "list") {
			return view.Forbidden("forex_rate:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		rows := fetchForexRates(ctx, deps, status)
		statusTabs := buildStatusTabs(deps)
		tableConfig := buildTableConfig(deps, rows, perms)

		heading, caption := headingForStatus(deps.Labels, status)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "forex-rates",
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-refresh-cw",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "forex-rate-list-content",
			ActiveStatus:    status,
			StatusTabs:      statusTabs,
			Table:           tableConfig,
			Labels:          deps.Labels,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "forex_rate"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("forex-rate-list", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchForexRates(ctx context.Context, deps *Deps, status string) []ForexRateRow {
	if deps.ListForexRates == nil {
		return []ForexRateRow{}
	}

	resp, err := deps.ListForexRates(ctx, &forexratepb.ListForexRatesRequest{})
	if err != nil {
		log.Printf("ListForexRates error: %v", err)
		return []ForexRateRow{}
	}
	if resp == nil {
		return []ForexRateRow{}
	}

	rows := make([]ForexRateRow, 0)
	for _, fr := range resp.GetData() {
		row := protoToRow(fr)
		if status == "active" && row.Status != "active" {
			continue
		}
		if status == "superseded" && row.Status != "superseded" {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func protoToRow(fr *forexratepb.ForexRate) ForexRateRow {
	return ForexRateRow{
		ID:             fr.GetId(),
		FromCurrency:   fr.GetFromCurrency(),
		ToCurrency:     fr.GetToCurrency(),
		RateMicroUnits: fmt.Sprintf("%d", fr.GetRateMicroUnits()),
		Source:         sourceLabel(fr.GetSource()),
		EffectiveFrom:  fr.GetEffectiveFrom(),
		Status:         statusString(fr.GetStatus()),
	}
}

// ---------------------------------------------------------------------------
// Tab builder
// ---------------------------------------------------------------------------

func buildStatusTabs(deps *Deps) []pyeza.TabItem {
	base := deps.Routes.ListURL
	return []pyeza.TabItem{
		{Key: "active", Label: "Active", Href: route.ResolveURL(base, "status", "active"), Icon: ""},
		{Key: "superseded", Label: "Superseded", Href: route.ResolveURL(base, "status", "superseded"), Icon: ""},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, rows []ForexRateRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "from_currency", Label: l.Columns.FromCurrency, WidthClass: "col-xl"},
		{Key: "to_currency", Label: l.Columns.ToCurrency, WidthClass: "col-xl"},
		{Key: "rate", Label: l.Columns.RateMicroUnits, WidthClass: "col-3xl", Align: "right"},
		{Key: "source", Label: l.Columns.Source, WidthClass: "col-2xl"},
		{Key: "effective_from", Label: l.Columns.EffectiveFrom, WidthClass: "col-3xl"},
	}

	tableRows := []types.TableRow{}
	canView := perms.Can("forex_rate", "read")
	for _, r := range rows {
		actions := []types.TableAction{
			{
				Type:            "view",
				Label:           l.Actions.View,
				Action:          "view",
				Href:            deps.Routes.DetailFor(r.ID),
				Disabled:        !canView,
				DisabledTooltip: l.Actions.NoPermission,
			},
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: deps.Routes.DetailFor(r.ID),
			Cells: []types.TableCell{
				{Type: "badge", Value: r.FromCurrency, Variant: "default"},
				{Type: "badge", Value: r.ToCurrency, Variant: "navy"},
				{Type: "text", Value: r.RateMicroUnits},
				{Type: "text", Value: r.Source},
				{Type: "text", Value: r.EffectiveFrom},
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tableConfig := &types.TableConfig{
		ID:                "forex-rates-table",
		Columns:           columns,
		Rows:              tableRows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowEntries:       true,
		DefaultSortColumn: "from_currency",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		// No primary action — forex rates are read-only (appended via RecordOperatorRate)
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func headingForStatus(l fycha.ForexRateLabels, status string) (heading, caption string) {
	switch status {
	case "superseded":
		return l.Page.HeadingSuperseded, l.Page.CaptionSuperseded
	default:
		return l.Page.HeadingActive, l.Page.CaptionActive
	}
}

func statusString(s forexratepb.ForexRateStatus) string {
	switch s {
	case forexratepb.ForexRateStatus_FOREX_RATE_STATUS_ACTIVE:
		return "active"
	case forexratepb.ForexRateStatus_FOREX_RATE_STATUS_SUPERSEDED:
		return "superseded"
	default:
		return "active"
	}
}

func sourceLabel(s forexratepb.ForexRateSource) string {
	switch s {
	case forexratepb.ForexRateSource_FOREX_RATE_SOURCE_OPERATOR:
		return "Operator"
	case forexratepb.ForexRateSource_FOREX_RATE_SOURCE_BSP_REF:
		return "BSP"
	default:
		return "Operator"
	}
}
