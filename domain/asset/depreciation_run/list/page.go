// Package list implements the depreciation-run history list page (Surface D).
// Mirror of packages/centymo-golang/views/revenue_run/list/page.go.
package list

import (
	"context"
	"fmt"
	"log"
	"strconv"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	depreciationrun "github.com/erniealice/fycha-golang/domain/asset/depreciation_run"
	drshared "github.com/erniealice/fycha-golang/domain/asset/depreciation_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// DepreciationRunPendingStaleMinutes is the threshold (minutes) after which a
// pending run is considered possibly-interrupted and shown with a warning chip.
const DepreciationRunPendingStaleMinutes = 5

// ListViewDeps holds view dependencies for the list page.
type ListViewDeps struct {
	Routes               depreciationrun.Routes
	Labels               depreciationrun.Labels
	CommonLabels         pyeza.CommonLabels
	TableLabels          types.TableLabels
	ListDepreciationRuns func(ctx context.Context, scope drshared.ListDepreciationRunsScope) ([]drshared.DepreciationRunRow, string, error)
}

// PageData is the full data context passed to the depreciation-run-list template.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewView creates the full-page depreciation-run list view.
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: view-package is `depreciation_run`,
		// permission entity is `depreciation_schedule`.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "list") {
			return view.Forbidden("depreciation_schedule:list")
		}
		_ = perms

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := depreciationRunColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(
			viewCtx.Request,
			types.SortableKeys(columns),
			types.FilterableKeys(columns),
			"initiated_at",
			"desc",
		)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		l := deps.Labels
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusPageTitle(l, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   status,
				HeaderTitle:    statusPageTitle(l, status),
				HeaderSubtitle: l.List.Subtitle,
				HeaderIcon:     "icon-zap",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "depreciation-run-list-content",
			Table:           tableConfig,
		}

		return view.OK("depreciation-run-list", pageData)
	})
}

// NewTableView returns only the table-card HTML (used as HTMX refresh target).
func NewTableView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: inherit gate from full page.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "list") {
			return view.Forbidden("depreciation_schedule:list")
		}
		_ = perms

		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		columns := depreciationRunColumns(deps.Labels)
		p, err := espynahttp.ParseTableParamsWithFilters(
			viewCtx.Request,
			types.SortableKeys(columns),
			types.FilterableKeys(columns),
			"initiated_at",
			"desc",
		)
		if err != nil {
			return view.Error(err)
		}

		tableConfig, err := buildTableConfig(ctx, deps, columns, status, p)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches depreciation-run data and builds the table configuration.
func buildTableConfig(
	ctx context.Context,
	deps *ListViewDeps,
	columns []types.TableColumn,
	status string,
	p tableparams.TableQueryParams,
) (*types.TableConfig, error) {
	if deps.ListDepreciationRuns == nil {
		log.Printf("depreciation-run list: ListDepreciationRuns callback is nil — returning empty table")
	}

	var rows []drshared.DepreciationRunRow
	var nextCursor string
	if deps.ListDepreciationRuns != nil {
		var err error
		rows, nextCursor, err = deps.ListDepreciationRuns(ctx, drshared.ListDepreciationRunsScope{
			Status: status,
			Limit:  int32(p.PageSize),
		})
		if err != nil {
			log.Printf("Failed to list depreciation runs: %v", err)
			return nil, fmt.Errorf("failed to load depreciation runs: %w", err)
		}
	}
	if rows == nil {
		rows = []drshared.DepreciationRunRow{}
	}

	l := deps.Labels
	// 2026-05-14 permission-gates P2b: depreciation_run list previously did not
	// load perms at all. Gate view-row actions on depreciation_schedule:read
	// (view-package name vs permission-entity name diverges — see plan §C1).
	perms := view.GetUserPermissions(ctx)
	tableRows := buildTableRows(rows, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, tableRows)

	refreshURL := route.ResolveURL(deps.Routes.ListTableURL, "status", status)

	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "cursor",
		SortColumn:    p.SortColumn,
		SortDirection: p.SortDir,
		FiltersJSON:   p.FiltersRaw,
		PaginationURL: refreshURL,
	}
	if nextCursor != "" {
		sp.NextCursor = nextCursor
	}
	sp.BuildDisplay()

	tableConfig := &types.TableConfig{
		ID:                   "depreciation-run-list-table",
		RefreshURL:           refreshURL,
		Columns:              columns,
		Rows:                 tableRows,
		ShowSearch:           false,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "initiated_at",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   statusEmptyTitle(l, status),
			Message: statusEmptyMessage(l, status),
		},
		ServerPagination: sp,
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig, nil
}

func depreciationRunColumns(l depreciationrun.Labels) []types.TableColumn {
	lc := l.List.Columns
	return []types.TableColumn{
		{Key: "id", Label: lc.ID, WidthClass: "col-5xl"},
		{Key: "scope", Label: lc.Scope, NoSort: true, NoFilter: true},
		{Key: "as_of_date", Label: lc.AsOfDate, WidthClass: "col-3xl"},
		{Key: "initiator", Label: lc.Initiator, WidthClass: "col-6xl", NoSort: true},
		{Key: "initiated_at", Label: lc.InitiatedAt, WidthClass: "col-4xl"},
		{Key: "status", Label: lc.Status, WidthClass: "col-3xl", NoSort: true},
		{Key: "created", Label: lc.Created, WidthClass: "col-md", Align: "right"},
		{Key: "skipped", Label: lc.Skipped, WidthClass: "col-md", Align: "right"},
		{Key: "errored", Label: lc.Errored, WidthClass: "col-md", Align: "right"},
	}
}

func buildTableRows(rows []drshared.DepreciationRunRow, l depreciationrun.Labels, routes depreciationrun.Routes, perms *types.UserPermissions) []types.TableRow {
	// 2026-05-14 permission-gates P2b: row view action gated on
	// depreciation_schedule:read (catalog entity name).
	canRead := perms.Can("depreciation_schedule", "read")

	tableRows := make([]types.TableRow, 0, len(rows))
	for _, r := range rows {
		detailURL := routes.DetailFor(r.ID)

		var statusCell types.TableCell
		if r.IsStalePending {
			statusCell = types.TableCell{
				Type:    "badge",
				Value:   l.StatusBadges.PossiblyInterrupted,
				Variant: "warning",
			}
		} else {
			label, variant := statusBadge(l, r.Status)
			statusCell = types.TableCell{
				Type:    "badge",
				Value:   label,
				Variant: variant,
			}
		}

		scopeDisplay := scopeKindLabel(l, r.ScopeKind)
		if r.ScopeLabel != "" {
			scopeDisplay = scopeKindLabel(l, r.ScopeKind) + ": " + r.ScopeLabel
		}

		actions := []types.TableAction{
			{
				Type:            "view",
				Label:           l.List.Columns.Actions,
				Action:          "view",
				Href:            detailURL,
				Disabled:        !canRead,
				DisabledTooltip: l.Errors.PermissionDenied,
			},
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: detailURL,
			DataAttrs: map[string]string{
				"run-id": r.ID,
			},
			Cells: []types.TableCell{
				{Type: "text", Value: r.ID},
				{Type: "text", Value: scopeDisplay},
				{Type: "text", Value: r.AsOfDate},
				{Type: "text", Value: r.InitiatorName},
				types.DateTimeCell(r.InitiatedAt, types.DateTimeFull),
				statusCell,
				{Type: "text", Value: strconv.Itoa(int(r.CreatedCount)), Align: "right"},
				{Type: "text", Value: strconv.Itoa(int(r.SkippedCount)), Align: "right"},
				{Type: "text", Value: strconv.Itoa(int(r.ErroredCount)), Align: "right"},
			},
			Actions: actions,
		})
	}
	return tableRows
}

func statusBadge(l depreciationrun.Labels, status string) (label, variant string) {
	switch status {
	case "pending":
		return l.StatusBadges.Pending, "warning"
	case "complete":
		return l.StatusBadges.Complete, "success"
	case "failed":
		return l.StatusBadges.Failed, "error"
	default:
		return status, "info"
	}
}

func scopeKindLabel(l depreciationrun.Labels, kind string) string {
	switch kind {
	case "asset":
		return l.ScopeKind.Asset
	case "category":
		return l.ScopeKind.Category
	case "policy":
		return l.ScopeKind.Policy
	case "workspace":
		return l.ScopeKind.Workspace
	default:
		return kind
	}
}

func statusPageTitle(l depreciationrun.Labels, status string) string {
	switch status {
	case "pending":
		return l.List.Title + " — " + l.List.FilterPending
	case "complete":
		return l.List.Title + " — " + l.List.FilterComplete
	case "failed":
		return l.List.Title + " — " + l.List.FilterFailed
	default:
		return l.List.Title
	}
}

func statusEmptyTitle(l depreciationrun.Labels, status string) string {
	switch status {
	case "pending":
		return l.List.Empty.Pending.Title
	case "complete":
		return l.List.Empty.Complete.Title
	case "failed":
		return l.List.Empty.Failed.Title
	default:
		return l.List.Empty.Pending.Title
	}
}

func statusEmptyMessage(l depreciationrun.Labels, status string) string {
	switch status {
	case "pending":
		return l.List.Empty.Pending.Message
	case "complete":
		return l.List.Empty.Complete.Message
	case "failed":
		return l.List.Empty.Failed.Message
	default:
		return l.List.Empty.Pending.Message
	}
}
