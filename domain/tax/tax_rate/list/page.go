// Package list provides the read-only Tax Rate list view.
package list

import (
	"context"
	"fmt"
	"log"

	taxratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_rate"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	"github.com/erniealice/fycha-golang/domain/tax/tax_rate"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the tax rate list page.
type Deps struct {
	Routes       tax_rate.Routes
	Labels       tax_rate.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Tax rate use cases
	ListTaxRates func(ctx context.Context, req *taxratepb.ListTaxRatesRequest) (*taxratepb.ListTaxRatesResponse, error)
}

// PageData holds the data for the tax rate list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveStatus    string
	StatusTabs      []pyeza.TabItem
	Table           *types.TableConfig
	Labels          tax_rate.Labels
}

// TaxRateRow is the view-model for a single tax rate row.
type TaxRateRow struct {
	ID            string
	Jurisdiction  string
	AuthorityCode string
	Kind          string
	TreatmentCode string
	Direction     string
	RateBps       string
	EffectiveFrom string
	Status        string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the tax rate list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("tax_rate", "list") {
			return view.Forbidden("tax_rate:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		rows := fetchTaxRates(ctx, deps, status)
		statusTabs := buildStatusTabs(deps)
		tableConfig := buildTableConfig(deps, rows, perms)

		heading, caption := headingForStatus(deps.Labels, status)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-percent",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "tax-rate-list-content",
			ActiveStatus:    status,
			StatusTabs:      statusTabs,
			Table:           tableConfig,
			Labels:          deps.Labels,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "tax_rate"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("tax-rate-list", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchTaxRates(ctx context.Context, deps *Deps, status string) []TaxRateRow {
	if deps.ListTaxRates == nil {
		return []TaxRateRow{}
	}

	resp, err := deps.ListTaxRates(ctx, &taxratepb.ListTaxRatesRequest{})
	if err != nil {
		log.Printf("ListTaxRates error: %v", err)
		return []TaxRateRow{}
	}
	if resp == nil {
		return []TaxRateRow{}
	}

	rows := make([]TaxRateRow, 0)
	for _, tr := range resp.GetData() {
		row := protoToRow(tr)
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

func protoToRow(tr *taxratepb.TaxRate) TaxRateRow {
	return TaxRateRow{
		ID:            tr.GetId(),
		Jurisdiction:  tr.GetJurisdiction(),
		AuthorityCode: tr.GetAuthorityCode(),
		Kind:          tr.GetKind(),
		TreatmentCode: tr.GetTreatmentCode(),
		Direction:     directionLabel(tr.GetDirection()),
		RateBps:       fmt.Sprintf("%d", tr.GetRateBasisPoints()),
		EffectiveFrom: tr.GetEffectiveFrom(),
		Status:        statusString(tr.GetStatus()),
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

func buildTableConfig(deps *Deps, rows []TaxRateRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "jurisdiction", Label: l.Columns.Jurisdiction, WidthClass: "col-3xl"},
		{Key: "authority", Label: l.Columns.AuthorityCode, WidthClass: "col-2xl"},
		{Key: "kind", Label: l.Columns.Kind, WidthClass: "col-2xl"},
		{Key: "treatment", Label: l.Columns.TreatmentCode, WidthClass: "col-2xl"},
		{Key: "direction", Label: l.Columns.Direction, WidthClass: "col-2xl"},
		{Key: "rate_bps", Label: l.Columns.RateBps, WidthClass: "col-xl", Align: "right"},
		{Key: "effective_from", Label: l.Columns.EffectiveFrom, WidthClass: "col-3xl"},
	}

	tableRows := []types.TableRow{}
	canView := perms.Can("tax_rate", "read")
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

		dirVariant := directionVariant(r.Direction)

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: deps.Routes.DetailFor(r.ID),
			Cells: []types.TableCell{
				{Type: "text", Value: r.Jurisdiction},
				{Type: "text", Value: r.AuthorityCode},
				{Type: "text", Value: r.Kind},
				{Type: "text", Value: r.TreatmentCode},
				{Type: "badge", Value: r.Direction, Variant: dirVariant},
				{Type: "text", Value: r.RateBps},
				{Type: "text", Value: r.EffectiveFrom},
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tableConfig := &types.TableConfig{
		ID:                "tax-rates-table",
		Columns:           columns,
		Rows:              tableRows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowEntries:       true,
		DefaultSortColumn: "jurisdiction",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		// No primary action — tax rates are read-only
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func headingForStatus(l tax_rate.Labels, status string) (heading, caption string) {
	switch status {
	case "superseded":
		return l.Page.HeadingSuperseded, l.Page.CaptionSuperseded
	default:
		return l.Page.HeadingActive, l.Page.CaptionActive
	}
}

func statusString(s taxratepb.TaxRateStatus) string {
	switch s {
	case taxratepb.TaxRateStatus_TAX_RATE_STATUS_ACTIVE:
		return "active"
	case taxratepb.TaxRateStatus_TAX_RATE_STATUS_SUPERSEDED:
		return "superseded"
	case taxratepb.TaxRateStatus_TAX_RATE_STATUS_VOIDED:
		return "voided"
	default:
		return "active"
	}
}

func directionLabel(d taxratepb.TaxRateDirection) string {
	switch d {
	case taxratepb.TaxRateDirection_TAX_RATE_DIRECTION_SURCHARGE:
		return "Surcharge"
	case taxratepb.TaxRateDirection_TAX_RATE_DIRECTION_WITHHOLDING:
		return "Withholding"
	default:
		return "Unknown"
	}
}

func directionVariant(label string) string {
	switch label {
	case "Surcharge":
		return "navy"
	case "Withholding":
		return "amber"
	default:
		return "default"
	}
}
