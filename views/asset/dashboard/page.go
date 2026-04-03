package dashboard

import (
	"context"
	"fmt"
	"html/template"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// DashboardStats holds count/value data for stat cards.
type DashboardStats struct {
	TotalAssets      int
	TotalBookValue   string
	FullyDepreciated int
	UnderMaintenance int
}

// ActivityItem represents a single entry in the recent activity feed.
type ActivityItem struct {
	IconHTML    template.HTML
	Title       string
	Description string
	TimeAgo     string
}

// Deps holds view dependencies.
type Deps struct {
	Routes       fycha.AssetRoutes
	Labels       fycha.AssetLabels
	CommonLabels pyeza.CommonLabels
}

// PageData holds the data for the asset dashboard page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Stats           DashboardStats
	RecentActivity  []ActivityItem
	Labels          fycha.AssetDashboardLabels
}

// NewView creates the asset dashboard view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Dashboard

		// Mock statistics
		stats := DashboardStats{
			TotalAssets:      24,
			TotalBookValue:   formatCurrency(1_245_750.00),
			FullyDepreciated: 3,
			UnderMaintenance: 2,
		}

		// Mock recent activity
		recentActivity := []ActivityItem{
			{
				IconHTML:    template.HTML(`<svg class="icon"><use href="#icon-box"></use></svg>`),
				Title:       l.ActivityAcquired,
				Description: "Office Laptop (Dell XPS 15) added to register",
				TimeAgo:     "2 hours ago",
			},
			{
				IconHTML:    template.HTML(`<svg class="icon"><use href="#icon-tool"></use></svg>`),
				Title:       l.ActivityMaintenance,
				Description: "Air Conditioning Unit - Annual servicing",
				TimeAgo:     "1 day ago",
			},
			{
				IconHTML:    template.HTML(`<svg class="icon"><use href="#icon-trending-down"></use></svg>`),
				Title:       l.ActivityDepreciation,
				Description: "Monthly depreciation for 24 assets processed",
				TimeAgo:     "3 days ago",
			},
		}

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "assets",
				ActiveSubNav:   "assets-dashboard",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-box",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "asset-dashboard-content",
			Stats:           stats,
			RecentActivity:  recentActivity,
			Labels:          l,
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

func formatCurrency(amount float64) string {
	whole := int64(amount)
	frac := int64((amount-float64(whole))*100 + 0.5)
	if frac >= 100 {
		whole++
		frac -= 100
	}
	wholeStr := fmt.Sprintf("%d", whole)
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
	return fmt.Sprintf("\u20b1%s.%02d", wholeStr, frac)
}
