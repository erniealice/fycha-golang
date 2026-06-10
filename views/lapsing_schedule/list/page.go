// Package list implements the Surface B Lapsing Schedule workspace list page
// (Phase 7 of the 20260507-asset-lapsing-revaluation plan).
//
// One ListDepreciationCandidates(scope_kind=WORKSPACE) call per page request.
// No per-asset fan-out. Request-scoped memoisation keyed by
// (workspace_id, as_of_date) is handled at the block.go shim layer via the
// ListCandidates callback — views never call espyna directly.
package list

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// CandidateRow is the view-layer shape for one row on the lapsing-schedule list.
// Populated by the ListCandidates callback injected via block.go.
type CandidateRow struct {
	AssetID           string
	AssetCode         string
	AssetName         string
	CategoryID        string
	CategoryName      string
	PolicyID          string // = CategoryID in v1
	PolicyName        string // = CategoryName in v1
	Currency          string
	CurrentBookValue  int64 // centavos
	LastPostedPeriod  string
	NextPendingPeriod string
	PendingCount      int
	NextAmount        int64 // centavos
	// Status: "up_to_date" | "pending" | "not_started" | "blocked"
	Status       string
	BlockerLabel string // populated when Status == "blocked"
	CanRun       bool
}

// ViewDeps holds all dependencies for the lapsing-schedule list page.
type ViewDeps struct {
	Routes                fycha.LapsingScheduleRoutes
	AssetRoutes           fycha.AssetRoutes
	DepreciationRunRoutes fycha.DepreciationRunRoutes
	Labels                fycha.DepreciationRunLabels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels

	// ListCandidates fetches all workspace-scoped depreciation candidates for the
	// given as_of_date + pagination cursor. Block.go wires this to espyna.
	ListCandidates func(ctx context.Context, asOfDate, cursor string, limit int32) ([]CandidateRow, string, error)
}

// PageData is the template context for the full lapsing-schedule-list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	AsOfDate        string
	AsOfDateMax     string
	RefreshURL      string
	Labels          fycha.DepreciationRunLapsingScheduleLabels
	CommonLabels    pyeza.CommonLabels
}

// NewView creates the full-page lapsing-schedule list view.
func NewView(deps *ViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: view-package is `lapsing_schedule`,
		// permission entity is `depreciation_schedule`.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "list") {
			return view.Forbidden("depreciation_schedule:list")
		}

		asOfDate, asOfDateMax := resolveAsOfDate(viewCtx.Request.URL.Query().Get("as_of_date"))
		cursor := viewCtx.Request.URL.Query().Get("cursor")

		tableConfig, err := buildTableConfig(ctx, deps, asOfDate, cursor)
		if err != nil {
			return view.Error(err)
		}

		l := deps.Labels.LapsingSchedule
		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "assets",
				ActiveSubNav:   "lapsing-schedule",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-calendar",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "lapsing-schedule-list-content",
			Table:           tableConfig,
			AsOfDate:        asOfDate,
			AsOfDateMax:     asOfDateMax,
			RefreshURL:      deps.Routes.ListURL,
			Labels:          l,
			CommonLabels:    deps.CommonLabels,
		}

		return view.OK("lapsing-schedule-list", pageData)
	})
}

// NewTableView returns only the table-card HTML for HTMX inner-swap on AsOfDate change.
func NewTableView(deps *ViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: inherit the same gate as full page.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "list") {
			return view.Forbidden("depreciation_schedule:list")
		}
		_ = perms

		asOfDate, _ := resolveAsOfDate(viewCtx.Request.URL.Query().Get("as_of_date"))
		cursor := viewCtx.Request.URL.Query().Get("cursor")

		tableConfig, err := buildTableConfig(ctx, deps, asOfDate, cursor)
		if err != nil {
			return view.Error(err)
		}

		return view.OK("table-card", tableConfig)
	})
}

// buildTableConfig fetches candidates and assembles a TableConfig.
func buildTableConfig(
	ctx context.Context,
	deps *ViewDeps,
	asOfDate, cursor string,
) (*types.TableConfig, error) {
	l := deps.Labels.LapsingSchedule
	perms := view.GetUserPermissions(ctx)

	var rows []CandidateRow
	var nextCursor string
	if deps.ListCandidates != nil {
		var err error
		rows, nextCursor, err = deps.ListCandidates(ctx, asOfDate, cursor, 50)
		if err != nil {
			log.Printf("lapsing-schedule list: ListCandidates error: %v", err)
			return nil, fmt.Errorf("failed to load lapsing schedule: %w", err)
		}
	}
	if rows == nil {
		rows = []CandidateRow{}
	}

	columns := lapsingScheduleColumns(l)
	tableRows := buildTableRows(rows, deps, perms)
	types.ApplyColumnStyles(columns, tableRows)

	// Bulk actions
	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)
	bulkCfg.Actions = []types.BulkAction{
		{
			Key:             "run-selected",
			Label:           l.BulkRunForSelected,
			Icon:            "icon-zap",
			Variant:         "primary",
			Endpoint:        deps.Routes.BulkRunSelectedURL,
			ExtraParamsJSON: `{"selection_mode":"selected","as_of_date":"` + asOfDate + `"}`,
		},
		{
			Key:             "run-all-matching",
			Label:           l.BulkRunForAllMatching,
			Icon:            "icon-zap",
			Variant:         "warning",
			Endpoint:        deps.Routes.BulkRunAllMatchingURL,
			ExtraParamsJSON: `{"selection_mode":"all_matching","as_of_date":"` + asOfDate + `"}`,
		},
	}
	// 2026-05-14 permission-gates P3: re-key from non-catalog `asset:depreciate`
	// to catalog `depreciation_schedule:create`. P2b: render disabled-with-tooltip
	// instead of removing the actions outright (per plan §"Pyeza primitive
	// contract" — surfaces stay visible so users know what perm to request).
	canBulkRun := perms.Can("depreciation_schedule", "create")
	if !canBulkRun {
		bulkTooltip := fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "depreciation_schedule:create")
		for i := range bulkCfg.Actions {
			bulkCfg.Actions[i].Disabled = true
			bulkCfg.Actions[i].DisabledTooltip = bulkTooltip
		}
	}

	sp := &types.ServerPagination{
		Enabled:       true,
		Mode:          "cursor",
		PaginationURL: deps.Routes.ListURL,
	}
	if nextCursor != "" {
		sp.NextCursor = nextCursor
	}
	sp.BuildDisplay()

	tableConfig := &types.TableConfig{
		ID:          "lapsing-schedule-list-table",
		RefreshURL:  deps.Routes.ListURL,
		Columns:     columns,
		Rows:        tableRows,
		ShowSearch:  false,
		ShowActions: true,
		ShowFilters: false,
		ShowSort:    false,
		ShowColumns: false,
		ShowExport:  false,
		ShowDensity: true,
		ShowEntries: false,
		Labels:      deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.EmptyTitle,
			Message: l.EmptyMessage,
		},
		BulkActions:      &bulkCfg,
		ServerPagination: sp,
	}

	types.ApplyTableSettings(tableConfig)
	return tableConfig, nil
}

// lapsingScheduleColumns defines the column layout for Surface B.
func lapsingScheduleColumns(l fycha.DepreciationRunLapsingScheduleLabels) []types.TableColumn {
	lc := l.Columns
	return []types.TableColumn{
		{Key: "asset", Label: lc.Asset, WidthClass: "col-9xl"},
		{Key: "category", Label: lc.Category, WidthClass: "col-6xl", NoSort: true},
		{Key: "policy", Label: lc.Policy, WidthClass: "col-6xl", NoSort: true},
		{Key: "currency", Label: lc.Currency, WidthClass: "col-2xl", NoSort: true},
		{Key: "book_value", Label: lc.CurrentBookValue, WidthClass: "col-5xl", Align: "right", NoSort: true},
		{Key: "last_posted", Label: lc.LastPostedPeriod, WidthClass: "col-4xl", NoSort: true},
		{Key: "next_pending", Label: lc.NextPendingPeriod, WidthClass: "col-4xl", NoSort: true},
		{Key: "pending", Label: lc.Pending, WidthClass: "col-md", Align: "right", NoSort: true},
		{Key: "next_amount", Label: lc.NextAmount, WidthClass: "col-4xl", Align: "right", NoSort: true},
		{Key: "status", Label: lc.Status, WidthClass: "col-4xl", NoSort: true},
	}
}

// buildTableRows converts CandidateRows to pyeza TableRows.
func buildTableRows(
	rows []CandidateRow,
	deps *ViewDeps,
	perms *types.UserPermissions,
) []types.TableRow {
	l := deps.Labels
	ls := l.LapsingSchedule
	tableRows := make([]types.TableRow, 0, len(rows))

	for _, r := range rows {
		assetDetailURL := route.ResolveURL(deps.AssetRoutes.DetailURL, "id", r.AssetID) + "?tab=lapsing-actual-schedule"
		drawerURL := deps.AssetRoutes.DepreciationRunFor(r.AssetID)

		// Asset cell — linked to detail page.
		assetLabel := r.AssetCode
		if r.AssetName != "" {
			assetLabel = r.AssetCode + " — " + r.AssetName
		}
		assetCell := types.TableCell{Type: "link", Value: assetLabel, Href: assetDetailURL}

		// Status badge
		statusLabel, statusVariant := resolveStatusBadge(l, r)
		statusCell := types.TableCell{
			Type:    "badge",
			Value:   statusLabel,
			Variant: statusVariant,
		}

		// Actions: per-row [Run] opens Surface A drawer via HTMX.
		// 2026-05-14 permission-gates P3: re-key to catalog verb.
		actions := []types.TableAction{}
		canRun := r.CanRun && perms.Can("depreciation_schedule", "create")
		if drawerURL != "" {
			actions = append(actions, types.TableAction{
				Type:            "run",
				Label:           ls.Columns.Actions,
				Action:          "run",
				HxGet:           drawerURL,
				HxTarget:        "#sheetContent",
				HxSwap:          "innerHTML",
				OnClick:         "lf.ui.Sheet.open()",
				Disabled:        !canRun,
				DisabledTooltip: l.Errors.PermissionDenied,
				TestID:          "lapsing-schedule-list-row-action-run",
			})
		}

		canRunStr := "false"
		if r.CanRun {
			canRunStr = "true"
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.AssetID,
			Href: assetDetailURL,
			DataAttrs: map[string]string{
				"asset-id":   r.AssetID,
				"status":     r.Status,
				"can-run":    canRunStr,
				"in-service": strconv.FormatBool(r.CanRun),
				"testid":     "lapsing-schedule-list-row",
			},
			Cells: []types.TableCell{
				assetCell,
				{Type: "text", Value: r.CategoryName},
				{Type: "text", Value: r.PolicyName},
				{Type: "badge", Value: r.Currency, Variant: "info"},
				types.MoneyCell(float64(r.CurrentBookValue), r.Currency, true),
				{Type: "text", Value: r.LastPostedPeriod},
				{Type: "text", Value: r.NextPendingPeriod},
				{Type: "text", Value: strconv.Itoa(r.PendingCount), Align: "right"},
				types.MoneyCell(float64(r.NextAmount), r.Currency, true),
				statusCell,
			},
			Actions: actions,
		})
	}
	return tableRows
}

// resolveStatusBadge returns the badge label and variant for a candidate row.
func resolveStatusBadge(l fycha.DepreciationRunLabels, r CandidateRow) (label, variant string) {
	ls := l.LapsingSchedule
	switch r.Status {
	case "up_to_date":
		return ls.StatusUpToDate, "success"
	case "not_started":
		return ls.StatusNotStarted, "neutral"
	case "blocked":
		msg := ls.StatusBlockedTemplate
		if r.BlockerLabel != "" {
			msg = ls.BlockedPrefix + r.BlockerLabel
		}
		return msg, "error"
	default:
		// "pending" — N periods pending
		tmpl := ls.StatusNPeriodsPendingTemplate
		_ = tmpl // lyngua template — block.go can format with .Count; view returns raw template key
		return fmt.Sprintf("%d %s", r.PendingCount, ls.Columns.Pending), "warning"
	}
}

// resolveAsOfDate returns the as_of_date and max date (today).
func resolveAsOfDate(input string) (asOfDate, maxDate string) {
	today := time.Now().UTC().Format("2006-01-02")
	if input == "" {
		return today, today
	}
	if _, err := time.Parse("2006-01-02", input); err != nil {
		return today, today
	}
	return input, today
}
