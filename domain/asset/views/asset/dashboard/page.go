package dashboard

import (
	"context"
	"fmt"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	asset "github.com/erniealice/fycha-golang/domain/asset"
)

// Deps holds view dependencies.
type Deps struct {
	Routes       asset.AssetRoutes
	Labels       asset.AssetLabels
	CommonLabels pyeza.CommonLabels
}

// PageData is what the asset dashboard template receives.
type PageData struct {
	types.PageData
	ContentTemplate string
	Dashboard       types.DashboardData
}

// NewView creates the asset dashboard view.
//
// Phase 1c refactor (2026-05-02): wired onto the pyeza "dashboard" block.
// Stats remain placeholder/mock until a future phase wires real
// fixed-asset aggregate methods (count by status, sum book value, etc.).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard

		// Mock stats — preserved from pre-refactor implementation.
		totalAssets := 24
		fullyDepreciated := 3
		underMaintenance := 2
		// MoneyCell formats centavos -> "12,457,500.00" with currency "PHP".
		bookValueCell := types.MoneyCell(1_245_750_00, "PHP", true)
		bookValueDisplay := bookValueCell.Currency + " " + bookValueCell.Value

		// Synthetic 12-month asset value trend (placeholder until aggregate wires up).
		trend := &types.ChartData{
			Labels: []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"},
			Series: []types.ChartSeries{{
				Name:   l.TotalBookValue,
				Values: []float64{1_280_000, 1_268_000, 1_256_000, 1_244_000, 1_252_000, 1_245_750, 1_240_000, 1_230_000, 1_220_000, 1_215_000, 1_210_000, 1_205_000},
				Color:  "navy",
			}},
		}
		trend.AutoScale()

		// Activity items — preserved from pre-refactor implementation.
		recent := []types.ActivityItem{
			{
				IconName:    "icon-box",
				IconVariant: "client",
				Title:       l.ActivityAcquired,
				Description: "Office Laptop (Dell XPS 15) added to register",
				Time:        "2h ago",
				TestID:      "asset-activity-acquired",
			},
			{
				IconName:    "icon-tool",
				IconVariant: "award",
				Title:       l.ActivityMaintenance,
				Description: "Air Conditioning Unit — Annual servicing",
				Time:        "1d ago",
				TestID:      "asset-activity-maintenance",
			},
			{
				IconName:    "icon-trending-down",
				IconVariant: "integration",
				Title:       l.ActivityDepreciation,
				Description: "Monthly depreciation for 24 assets processed",
				Time:        "3d ago",
				TestID:      "asset-activity-depreciation",
			},
		}

		dash := types.DashboardData{
			QuickActions: []types.QuickAction{
				{Icon: "icon-plus", Label: l.QuickNewAsset, Href: deps.Routes.AddURL, Variant: "primary", TestID: "asset-action-new"},
				{Icon: "icon-list", Label: l.QuickViewAll, Href: deps.Routes.ListURL, TestID: "asset-action-list"},
				{Icon: "icon-trending-down", Label: l.QuickDepreciationSchedule, Href: deps.Routes.LapsingScheduleURL, TestID: "asset-action-lapsing"},
				{Icon: "icon-tool", Label: l.QuickMaintenanceLog, Href: deps.Routes.DepreciationPoliciesURL, TestID: "asset-action-policies"},
			},
			Stats: []types.StatCardData{
				{Icon: "icon-box", Value: fmt.Sprintf("%d", totalAssets), Label: l.TotalAssets, Color: "terracotta", TestID: "asset-stat-total"},
				{Icon: "icon-dollar-sign", Value: bookValueDisplay, Label: l.TotalBookValue, Color: "sage", TestID: "asset-stat-book-value"},
				{Icon: "icon-trending-down", Value: fmt.Sprintf("%d", fullyDepreciated), Label: l.FullyDepreciated, Color: "navy", TestID: "asset-stat-depreciated"},
				{Icon: "icon-tool", Value: fmt.Sprintf("%d", underMaintenance), Label: l.UnderMaintenance, Color: "amber", TestID: "asset-stat-maintenance"},
			},
			Widgets: []types.DashboardWidget{
				{
					ID: "trend", Title: l.AssetValueTrend, Type: "chart", ChartKind: "line",
					ChartData: trend, Span: 2,
				},
				{
					ID: "recent", Title: l.RecentActivity, Type: "list", Span: 1,
					HeaderActions: []types.QuickAction{
						{Label: l.ViewAll, Href: deps.Routes.ListURL},
					},
					ListItems: recent,
				},
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "asset",
				ActiveSubNav:   "assets-dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-box",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "asset-dashboard-content",
			Dashboard:       dash,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "asset-dashboard"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("asset-dashboard", pageData)
	})
}
