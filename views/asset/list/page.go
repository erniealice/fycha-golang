package list

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// MockAsset represents a fixed asset for mock data display.
type MockAsset struct {
	ID              string
	AssetNumber     string
	Name            string
	CategoryName    string
	LocationName    string
	AcquisitionCost float64
	BookValue       float64
	Active          bool
}

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes       fycha.AssetRoutes
	Labels       fycha.AssetLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the asset list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the asset list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, status, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "assets",
				ActiveSubNav:   "assets-fixed",
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-box",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "asset-list-content",
			Table:           tableConfig,
		}

		return view.OK("asset-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, status, perms)
		return view.OK("table-card", tableConfig)
	})
}

// mockAssets returns hardcoded asset data for initial UI development.
func mockAssets() []MockAsset {
	return []MockAsset{
		{ID: "ast-001", AssetNumber: "FA-001", Name: "Office Laptop (Dell XPS 15)", CategoryName: "IT Equipment", LocationName: "Main Office", AcquisitionCost: 85000, BookValue: 42500, Active: true},
		{ID: "ast-002", AssetNumber: "FA-002", Name: "Salon Chair (Hydraulic)", CategoryName: "Furniture", LocationName: "Branch 1", AcquisitionCost: 25000, BookValue: 18750, Active: true},
		{ID: "ast-003", AssetNumber: "FA-003", Name: "Hair Dryer (Professional)", CategoryName: "Equipment", LocationName: "Branch 1", AcquisitionCost: 12000, BookValue: 6000, Active: true},
		{ID: "ast-004", AssetNumber: "FA-004", Name: "POS Terminal", CategoryName: "IT Equipment", LocationName: "Main Office", AcquisitionCost: 35000, BookValue: 17500, Active: true},
		{ID: "ast-005", AssetNumber: "FA-005", Name: "Air Conditioning Unit", CategoryName: "Building Equipment", LocationName: "Main Office", AcquisitionCost: 65000, BookValue: 0, Active: false},
		{ID: "ast-006", AssetNumber: "FA-006", Name: "Massage Table", CategoryName: "Furniture", LocationName: "Branch 2", AcquisitionCost: 18000, BookValue: 13500, Active: true},
		{ID: "ast-007", AssetNumber: "FA-007", Name: "UV Sterilizer Cabinet", CategoryName: "Equipment", LocationName: "Branch 1", AcquisitionCost: 8500, BookValue: 2125, Active: false},
	}
}

func buildTableConfig(deps *ListViewDeps, status string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := assetColumns(l)
	rows := buildTableRows(mockAssets(), status, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := fycha.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status, deps.Routes)

	refreshURL := route.ResolveURL(deps.Routes.TableURL, "status", status)

	tableConfig := &types.TableConfig{
		ID:                   "assets-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "asset_number",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddAsset,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("asset", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func assetColumns(l fycha.AssetLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "asset_number", Label: l.Columns.AssetNumber, Sortable: true, Width: "120px"},
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "category", Label: l.Columns.Category, Sortable: true},
		{Key: "location", Label: l.Columns.Location, Sortable: true},
		{Key: "acquisition_cost", Label: l.Columns.AcquisitionCost, Sortable: true, Width: "160px"},
		{Key: "book_value", Label: l.Columns.BookValue, Sortable: true, Width: "140px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(assets []MockAsset, status string, l fycha.AssetLabels, routes fycha.AssetRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, asset := range assets {
		recordStatus := "active"
		if !asset.Active {
			recordStatus = "inactive"
		}
		if recordStatus != status {
			continue
		}

		id := asset.ID
		name := asset.Name

		canUpdate := perms.Can("asset", "update")
		canDelete := perms.Can("asset", "delete")

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", id)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit", URL: route.ResolveURL(routes.EditURL, "id", id), DrawerTitle: l.Actions.Edit, Disabled: !canUpdate, DisabledTooltip: l.Actions.NoPermission},
		}
		if asset.Active {
			actions = append(actions, types.TableAction{
				Type: "deactivate", Label: l.Actions.Deactivate, Action: "deactivate",
				URL: routes.SetStatusURL + "?status=inactive", ItemName: name,
				ConfirmTitle:   l.Actions.Deactivate,
				ConfirmMessage: fmt.Sprintf(l.Actions.ConfirmDeactivate, name),
				Disabled:       !canUpdate, DisabledTooltip: l.Actions.NoPermission,
			})
		} else {
			actions = append(actions, types.TableAction{
				Type: "activate", Label: l.Actions.Activate, Action: "activate",
				URL: routes.SetStatusURL + "?status=active", ItemName: name,
				ConfirmTitle:   l.Actions.Activate,
				ConfirmMessage: fmt.Sprintf(l.Actions.ConfirmActivate, name),
				Disabled:       !canUpdate, DisabledTooltip: l.Actions.NoPermission,
			})
		}
		actions = append(actions, types.TableAction{
			Type:     "delete",
			Label:    l.Actions.Delete,
			Action:   "delete",
			URL:      routes.DeleteURL,
			ItemName: name,
			Disabled: !canDelete, DisabledTooltip: l.Actions.NoPermission,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: asset.AssetNumber},
				{Type: "text", Value: name},
				{Type: "text", Value: asset.CategoryName},
				{Type: "text", Value: asset.LocationName},
				{Type: "text", Value: formatCurrency(asset.AcquisitionCost)},
				{Type: "text", Value: formatCurrency(asset.BookValue)},
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":             name,
				"asset_number":     asset.AssetNumber,
				"category":         asset.CategoryName,
				"location":         asset.LocationName,
				"acquisition_cost": formatCurrency(asset.AcquisitionCost),
				"book_value":       formatCurrency(asset.BookValue),
				"status":           recordStatus,
			},
			Actions: actions,
		})
	}
	return rows
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

func statusTitle(l fycha.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l fycha.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l fycha.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l fycha.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveMessage
	case "inactive":
		return l.Empty.InactiveMessage
	default:
		return l.Empty.ActiveMessage
	}
}

func statusVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	default:
		return "default"
	}
}

func buildBulkActions(l fycha.AssetLabels, common pyeza.CommonLabels, status string, routes fycha.AssetRoutes) []types.BulkAction {
	actions := []types.BulkAction{}

	switch status {
	case "active":
		actions = append(actions, types.BulkAction{
			Key:             "deactivate",
			Label:           l.Actions.Deactivate,
			Icon:            "icon-archive",
			Variant:         "warning",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Deactivate,
			ConfirmMessage:  l.Actions.ConfirmBulkDeactivate,
			ExtraParamsJSON: `{"target_status":"inactive"}`,
		})
	case "inactive":
		actions = append(actions, types.BulkAction{
			Key:             "activate",
			Label:           l.Actions.Activate,
			Icon:            "icon-box",
			Variant:         "primary",
			Endpoint:        routes.BulkSetStatusURL,
			ConfirmTitle:    l.Actions.Activate,
			ConfirmMessage:  l.Actions.ConfirmBulkActivate,
			ExtraParamsJSON: `{"target_status":"active"}`,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:            "delete",
		Label:          common.Bulk.Delete,
		Icon:           "icon-trash-2",
		Variant:        "danger",
		Endpoint:       routes.BulkDeleteURL,
		ConfirmTitle:   common.Bulk.Delete,
		ConfirmMessage: l.Actions.ConfirmBulkDelete,
	})

	return actions
}
