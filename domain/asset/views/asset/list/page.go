package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	asset "github.com/erniealice/fycha-golang/domain/asset"
)

// AssetRow is a flat row returned by the list query. Exported so block.go
// can construct it from raw SQL without importing protobuf types.
type AssetRow struct {
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
	Routes       asset.AssetRoutes
	Labels       asset.AssetLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ListAssets returns asset rows filtered by status. Wired from block.go.
	ListAssets func(ctx context.Context, status string) ([]AssetRow, error)

	// GetAssetInUseIDs returns a map of asset IDs that have at least one
	// asset_transaction row. Used to set data-deletable=false and disable the
	// delete action on rows with posted transactions (H5 soft-delete gate).
	// Nil = skip the check (mock build or use cases not yet wired).
	GetAssetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
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
		// 2026-05-14 permission-gates P2a: reject direct-URL access when the
		// user lacks asset:list — sidebar visibility is not a security boundary.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "list") {
			return view.Forbidden("asset:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		tableConfig := buildTableConfig(ctx, deps, status, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "asset",
				ActiveSubNav:   "assets-fixed",
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-box",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "asset-list-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "asset"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("asset-list", pageData)
	})
}

// NewTableView creates a view that returns only the table-card HTML.
// Used as the refresh target after CRUD operations.
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: refresh-target partial inherits the
		// same gate as the full page.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "list") {
			return view.Forbidden("asset:list")
		}

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		tableConfig := buildTableConfig(ctx, deps, status, perms)
		return view.OK("table-card", tableConfig)
	})
}

func buildTableConfig(ctx context.Context, deps *ListViewDeps, status string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := assetColumns(l)

	var assets []AssetRow
	if deps.ListAssets != nil {
		var err error
		assets, err = deps.ListAssets(ctx, status)
		if err != nil {
			log.Printf("asset list query error: %v", err)
		}
	}

	// Batch-fetch which assets are in use (have any asset_transaction row).
	// Performed once per page render; nil map = no in-use assets (safe default).
	var inUseIDs map[string]bool
	if deps.GetAssetInUseIDs != nil && len(assets) > 0 {
		ids := make([]string, 0, len(assets))
		for _, a := range assets {
			ids = append(ids, a.ID)
		}
		var err error
		inUseIDs, err = deps.GetAssetInUseIDs(ctx, ids)
		if err != nil {
			// Non-fatal: log and fall back to permitting all deletes rather than
			// blocking all deletes. The server-side handler is the hard gate.
			log.Printf("asset in-use check error: %v", err)
		}
	}

	rows := buildTableRows(assets, l, deps.Routes, perms, status, inUseIDs)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	// 2026-05-14 permission-gates P2b: pass perms through so bulk actions
	// render disabled-with-tooltip when the user lacks asset:update / :delete.
	bulkCfg.Actions = buildBulkActions(l, deps.CommonLabels, status, deps.Routes, perms)

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

func assetColumns(l asset.AssetLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "asset_number", Label: l.Columns.AssetNumber, WidthClass: "col-2xl"},
		{Key: "name", Label: l.Columns.Name},
		{Key: "category", Label: l.Columns.Category},
		{Key: "location", Label: l.Columns.Location},
		{Key: "acquisition_cost", Label: l.Columns.AcquisitionCost, WidthClass: "col-5xl"},
		{Key: "book_value", Label: l.Columns.BookValue, WidthClass: "col-3xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-2xl"},
	}
}

func buildTableRows(assets []AssetRow, l asset.AssetLabels, routes asset.AssetRoutes, perms *types.UserPermissions, status string, inUseIDs map[string]bool) []types.TableRow {
	rows := []types.TableRow{}
	for _, asset := range assets {
		id := asset.ID
		name := asset.Name

		recordStatus := "active"
		if !asset.Active {
			recordStatus = "inactive"
		}

		canUpdate := perms.Can("asset", "update")
		canDelete := perms.Can("asset", "delete")

		// An asset is non-deletable when it has any posted asset_transaction row,
		// regardless of the operator's delete permission.
		isInUse := inUseIDs[id]
		deletable := canDelete && !isInUse

		// Determine tooltip for a disabled delete action:
		// - in-use takes priority over no-permission (more informative)
		var deleteTooltip string
		switch {
		case isInUse:
			deleteTooltip = l.Actions.CannotDeleteInUse
		case !canDelete:
			deleteTooltip = l.Actions.NoPermission
		}

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
			Type:            "delete",
			Label:           l.Actions.Delete,
			Action:          "delete",
			URL:             routes.DeleteURL,
			ItemName:        name,
			Disabled:        !deletable,
			DisabledTooltip: deleteTooltip,
		})

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: asset.AssetNumber},
				{Type: "text", Value: name},
				{Type: "text", Value: asset.CategoryName},
				{Type: "text", Value: asset.LocationName},
				// centMode=false is CORRECT here: assetToRow in block.go already
				// converts proto centavos → float64 pesos via float64(x)/100.
				types.MoneyCell(asset.AcquisitionCost, "PHP", false),
				types.MoneyCell(asset.BookValue, "PHP", false),
				{Type: "badge", Value: recordStatus, Variant: statusVariant(recordStatus)},
			},
			DataAttrs: map[string]string{
				"name":             name,
				"asset_number":     asset.AssetNumber,
				"category":         asset.CategoryName,
				"location":         asset.LocationName,
				"acquisition_cost": fmt.Sprintf("%.2f", asset.AcquisitionCost),
				"book_value":       fmt.Sprintf("%.2f", asset.BookValue),
				"status":           recordStatus,
				"deletable":        strconv.FormatBool(!isInUse),
			},
			Actions: actions,
		})
	}
	return rows
}

func statusTitle(l asset.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Page.HeadingActive
	case "inactive":
		return l.Page.HeadingInactive
	default:
		return l.Page.Heading
	}
}

func statusSubtitle(l asset.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Page.CaptionActive
	case "inactive":
		return l.Page.CaptionInactive
	default:
		return l.Page.Caption
	}
}

func statusEmptyTitle(l asset.AssetLabels, status string) string {
	switch status {
	case "active":
		return l.Empty.ActiveTitle
	case "inactive":
		return l.Empty.InactiveTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func statusEmptyMessage(l asset.AssetLabels, status string) string {
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

func buildBulkActions(l asset.AssetLabels, common pyeza.CommonLabels, status string, routes asset.AssetRoutes, perms *types.UserPermissions) []types.BulkAction {
	// 2026-05-14 permission-gates P2b: pyeza.BulkAction now exposes
	// Disabled + DisabledTooltip. Bulk activate/deactivate keys on
	// asset:update; bulk delete keys on asset:delete.
	canUpdate := perms.Can("asset", "update")
	canDelete := perms.Can("asset", "delete")
	updateTooltip := fmt.Sprintf(common.Errors.MissingPermission, "asset:update")
	deleteTooltip := fmt.Sprintf(common.Errors.MissingPermission, "asset:delete")

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
			Disabled:        !canUpdate,
			DisabledTooltip: updateTooltip,
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
			Disabled:        !canUpdate,
			DisabledTooltip: updateTooltip,
		})
	}

	actions = append(actions, types.BulkAction{
		Key:              "delete",
		Label:            common.Bulk.Delete,
		Icon:             "icon-trash-2",
		Variant:          "danger",
		Endpoint:         routes.BulkDeleteURL,
		ConfirmTitle:     common.Bulk.Delete,
		ConfirmMessage:   l.Actions.ConfirmBulkDelete,
		RequiresDataAttr: "deletable",
		Disabled:         !canDelete,
		DisabledTooltip:  deleteTooltip,
	})

	return actions
}
