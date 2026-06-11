// Package policies implements the Surface F actionable Depreciation Policies page
// (Phase 7.5 of the 20260507-asset-lapsing-revaluation plan).
//
// URL: /app/assets/settings/depreciation-policies
// Each row exposes [Preview] (read-only drawer) and [Run] (Surface C drawer).
// No bulk actions. Policy = AssetCategory in v1.
package policies

import (
	"context"
	"fmt"
	"log"
	"strconv"

	asset "github.com/erniealice/fycha-golang/domain/asset"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PolicyRow is the view-layer representation of one depreciation policy row.
// A policy corresponds to an AssetCategory in v1.
type PolicyRow struct {
	CategoryID         string
	PolicyID           string // = CategoryID in v1
	Name               string
	DepreciationMethod string
	UsefulLifeMonths   int32
	SalvagePct         float64
	AssetsInPolicy     int
	AssetsDeviating    int
}

// ViewDeps holds all dependencies for the depreciation policies page.
type ViewDeps struct {
	Routes       asset.AssetCategoryDepreciationRoutes
	Labels       asset.DepreciationPoliciesLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ListPolicies returns all asset categories enriched with assets-in-policy and
	// assets-deviating counts. Block.go wires this to espyna ListAssetCategories +
	// per-category aggregation queries. Nil = empty table (graceful degradation).
	ListPolicies func(ctx context.Context) ([]PolicyRow, error)
}

// PageData is the template context for the depreciation-policies page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	Labels          asset.DepreciationPoliciesLabels
}

// NewView creates the full-page depreciation policies view.
func NewView(deps *ViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		tableConfig, err := buildTableConfig(ctx, deps)
		if err != nil {
			return view.Error(err)
		}

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "assets",
				ActiveSubNav:   "depreciation-policies",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-settings",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "depreciation-policies-content",
			Table:           tableConfig,
			Labels:          l,
		}

		return view.OK("depreciation-policies", pageData)
	})
}

// buildTableConfig fetches policy rows and assembles a TableConfig.
func buildTableConfig(ctx context.Context, deps *ViewDeps) (*types.TableConfig, error) {
	l := deps.Labels
	perms := view.GetUserPermissions(ctx)

	var rows []PolicyRow
	if deps.ListPolicies != nil {
		var err error
		rows, err = deps.ListPolicies(ctx)
		if err != nil {
			log.Printf("depreciation-policies: ListPolicies error: %v", err)
			return nil, fmt.Errorf("failed to load depreciation policies: %w", err)
		}
	}
	if rows == nil {
		rows = []PolicyRow{}
	}

	columns := policiesColumns(l)
	tableRows := buildTableRows(rows, deps, perms)
	types.ApplyColumnStyles(columns, tableRows)

	tableConfig := &types.TableConfig{
		ID:          "depreciation-policies-table",
		Columns:     columns,
		Rows:        tableRows,
		ShowSearch:  false,
		ShowActions: true,
		ShowFilters: false,
		ShowSort:    false,
		ShowColumns: false,
		ShowExport:  false,
		ShowDensity: false,
		ShowEntries: false,
		Labels:      deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.EmptyTitle,
			Message: l.EmptyMessage,
		},
	}

	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

// policiesColumns defines the column layout for Surface F.
func policiesColumns(l asset.DepreciationPoliciesLabels) []types.TableColumn {
	lc := l.Columns
	return []types.TableColumn{
		{Key: "policy", Label: lc.Policy, WidthClass: "col-9xl"},
		{Key: "method", Label: lc.Method, WidthClass: "col-5xl", NoSort: true},
		{Key: "useful_life", Label: lc.UsefulLife, WidthClass: "col-3xl", Align: "right", NoSort: true},
		{Key: "salvage_pct", Label: lc.SalvagePct, WidthClass: "col-3xl", Align: "right", NoSort: true},
		{Key: "assets_in_policy", Label: lc.AssetsInPolicy, WidthClass: "col-4xl", Align: "right", NoSort: true},
		{Key: "assets_deviating", Label: lc.AssetsDeviating, WidthClass: "col-4xl", Align: "right", NoSort: true},
	}
}

// buildTableRows converts PolicyRows to pyeza TableRows.
func buildTableRows(
	rows []PolicyRow,
	deps *ViewDeps,
	perms *types.UserPermissions,
) []types.TableRow {
	l := deps.Labels
	tableRows := make([]types.TableRow, 0, len(rows))

	for _, r := range rows {
		previewURL := deps.Routes.PolicyPreviewFor(r.CategoryID)
		runURL := deps.Routes.PolicyRunFor(r.CategoryID)

		usefulLifeDisplay := fmt.Sprintf("%d", r.UsefulLifeMonths) + l.UsefulLifeMonthsSuffix
		salvagePctDisplay := fmt.Sprintf("%.0f%%", r.SalvagePct)

		deviatingVariant := "neutral"
		if r.AssetsDeviating > 0 {
			deviatingVariant = "warning"
		}

		actions := []types.TableAction{
			{
				Type:     "preview",
				Label:    l.ActionPreview,
				Action:   "preview",
				HxGet:    previewURL,
				HxTarget: "#sheetContent",
				HxSwap:   "innerHTML",
				OnClick:  "lf.ui.Sheet.open()",
				TestID:   "depreciation-policy-action-preview",
			},
			{
				Type:     "run",
				Label:    l.ActionRun,
				Action:   "run",
				HxGet:    runURL,
				HxTarget: "#sheetContent",
				HxSwap:   "innerHTML",
				OnClick:  "lf.ui.Sheet.open()",
				// 2026-05-14 permission-gates P3: re-key from non-catalog
				// `asset:depreciate` to catalog `depreciation_schedule:create`.
				Disabled:        !perms.Can("depreciation_schedule", "create"),
				DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "depreciation_schedule:create"),
				TestID:          "depreciation-policy-action-run",
			},
		}

		tableRows = append(tableRows, types.TableRow{
			ID: r.CategoryID,
			DataAttrs: map[string]string{
				"category-id": r.CategoryID,
				"policy-id":   r.PolicyID,
				"testid":      "depreciation-policy-row",
			},
			Cells: []types.TableCell{
				{Type: "text", Value: r.Name},
				{Type: "text", Value: r.DepreciationMethod},
				{Type: "text", Value: usefulLifeDisplay, Align: "right"},
				{Type: "text", Value: salvagePctDisplay, Align: "right"},
				{
					Type:  "text",
					Value: strconv.Itoa(r.AssetsInPolicy),
					Align: "right",
				},
				{
					Type:    "badge",
					Value:   strconv.Itoa(r.AssetsDeviating),
					Variant: deviatingVariant,
				},
			},
			Actions: actions,
		})
	}
	return tableRows
}
