// Package detail implements the depreciation-run detail page (Surface D).
// Mirror of packages/centymo-golang/views/revenue_run/detail/page.go.
package detail

import (
	"context"
	"fmt"
	"log"

	fycha "github.com/erniealice/fycha-golang"
	detailform "github.com/erniealice/fycha-golang/views/depreciation_run/detail/form"
	drshared "github.com/erniealice/fycha-golang/views/depreciation_run/shared"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// DetailViewDeps holds view dependencies for the detail page.
type DetailViewDeps struct {
	Routes       fycha.DepreciationRunRoutes
	Labels       fycha.DepreciationRunLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ReadDepreciationRun fetches a run + all schedule entries by run ID.
	ReadDepreciationRun func(ctx context.Context, id string) (*drshared.DepreciationRunWithEntries, error)

	// ListAssetTransactionsByRunID fetches asset_transaction rows scoped to a run.
	ListAssetTransactionsByRunID func(ctx context.Context, runID string) ([]drshared.AssetTransactionRow, error)
}

// NewView creates the full-page depreciation-run detail view.
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: view-package `depreciation_run`,
		// permission entity `depreciation_schedule`.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "read") {
			return view.Forbidden("depreciation_schedule:read")
		}
		_ = perms

		id := viewCtx.Request.PathValue("run_id")

		runWithEntries, err := deps.ReadDepreciationRun(ctx, id)
		if err != nil {
			log.Printf("Failed to read depreciation run %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load run: %w", err))
		}
		if runWithEntries == nil {
			log.Printf("Depreciation run %s not found", id)
			return view.Error(fmt.Errorf("run not found"))
		}

		l := deps.Labels
		run := runWithEntries.Run
		headerTitle := l.Detail.Title + " — " + run.ID

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "summary"
		}
		tabItems := buildTabItems(l, id, deps.Routes)

		pageData := &detailform.PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          headerTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    headerTitle,
				HeaderSubtitle: l.Detail.Title,
				HeaderIcon:     "icon-zap",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:       "depreciation-run-detail-content",
			Run:                   run,
			Entries:               runWithEntries.Entries,
			IsPossiblyInterrupted: run.IsStalePending,
			ActiveTab:             activeTab,
			TabItems:              tabItems,
			Labels:                l,
		}

		loadTabData(ctx, pageData, deps, runWithEntries, activeTab)

		return view.OK("depreciation-run-detail", pageData)
	})
}

// NewTabView creates a partial view that returns only the active tab content.
// Called via HTMX when the user clicks a tab button.
func NewTabView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates P2a: HTMX tab partial re-check.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "read") {
			return view.Forbidden("depreciation_schedule:read")
		}
		_ = perms

		id := viewCtx.Request.PathValue("run_id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "summary"
		}

		runWithEntries, err := deps.ReadDepreciationRun(ctx, id)
		if err != nil {
			log.Printf("Failed to read depreciation run %s: %v", id, err)
			return view.Error(fmt.Errorf("failed to load run: %w", err))
		}
		if runWithEntries == nil {
			log.Printf("Depreciation run %s not found", id)
			return view.Error(fmt.Errorf("run not found"))
		}

		l := deps.Labels
		pageData := &detailform.PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				CommonLabels: deps.CommonLabels,
			},
			Run:                   runWithEntries.Run,
			Entries:               runWithEntries.Entries,
			IsPossiblyInterrupted: runWithEntries.Run.IsStalePending,
			ActiveTab:             tab,
			TabItems:              buildTabItems(l, id, deps.Routes),
			Labels:                l,
		}

		loadTabData(ctx, pageData, deps, runWithEntries, tab)

		templateName := "depreciation-run-" + tab + "-tab"
		return view.OK(templateName, pageData)
	})
}

// loadTabData populates tab-specific fields on pageData.
func loadTabData(
	ctx context.Context,
	pageData *detailform.PageData,
	deps *DetailViewDeps,
	runWithEntries *drshared.DepreciationRunWithEntries,
	tab string,
) {
	l := deps.Labels
	entries := runWithEntries.Entries

	switch tab {
	case "summary":
		// All data is already in Run; nothing extra to load.

	case "selections":
		pageData.SelectionsTable = buildSelectionsTable(entries, l, deps.TableLabels)

	case "results":
		pageData.ResultsTable = buildResultsTable(entries, l, deps.TableLabels)

	case "transactions":
		if deps.ListAssetTransactionsByRunID != nil {
			txs, err := deps.ListAssetTransactionsByRunID(ctx, runWithEntries.Run.ID)
			if err != nil {
				log.Printf("Failed to load transactions for run %s: %v", runWithEntries.Run.ID, err)
				txs = []drshared.AssetTransactionRow{}
			}
			pageData.Transactions = txs
			pageData.TransactionsTable = buildTransactionsTable(txs, l, deps.TableLabels)
		} else {
			pageData.TransactionsTable = buildTransactionsTable(nil, l, deps.TableLabels)
		}

	case "history":
		// Deferred — rendered as an info alert in the template.
	}
}

// buildTabItems constructs the tab bar items for the detail page.
func buildTabItems(l fycha.DepreciationRunLabels, id string, routes fycha.DepreciationRunRoutes) []pyeza.TabItem {
	base := routes.DetailFor(id)
	action := route.ResolveURL(routes.DetailTabActionURL, "run_id", id, "tab", "")
	lt := l.Detail.Tabs
	return []pyeza.TabItem{
		{Key: "summary", Label: lt.Summary, Href: base + "?tab=summary", HxGet: action + "summary", Icon: "icon-info"},
		{Key: "selections", Label: lt.Selections, Href: base + "?tab=selections", HxGet: action + "selections", Icon: "icon-list"},
		{Key: "results", Label: lt.Results, Href: base + "?tab=results", HxGet: action + "results", Icon: "icon-check-circle"},
		{Key: "transactions", Label: lt.Transactions, Href: base + "?tab=transactions", HxGet: action + "transactions", Icon: "icon-activity"},
		{Key: "history", Label: lt.History, Href: base + "?tab=history", HxGet: action + "history", Icon: "icon-clock"},
	}
}

// buildSelectionsTable builds the TableConfig for the Selections tab.
// Shows all entries (all are selections — created, skipped, errored).
func buildSelectionsTable(entries []drshared.DepreciationRunEntryRow, l fycha.DepreciationRunLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "asset", Label: l.List.Columns.Scope, NoSort: true, WidthClass: "col-6xl"},
		{Key: "period_start_date", Label: l.Detail.Summary.AsOfDate, NoSort: true, WidthClass: "col-3xl"},
	}

	rows := make([]types.TableRow, 0, len(entries))
	for _, e := range entries {
		assetDisplay := e.AssetID
		if e.AssetName != "" {
			assetDisplay = e.AssetName
		}
		rows = append(rows, types.TableRow{
			ID: e.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: assetDisplay},
				{Type: "text", Value: e.PeriodStartDate},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "depreciation-run-selections-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.List.Empty.Pending.Title,
			Message: l.List.Empty.Pending.Message,
		},
	}
}

// buildResultsTable builds the TableConfig for the Results tab.
// Same data as selections but adds Outcome + Error columns.
func buildResultsTable(entries []drshared.DepreciationRunEntryRow, l fycha.DepreciationRunLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "asset", Label: l.List.Columns.Scope, NoSort: true, WidthClass: "col-6xl"},
		{Key: "period_start_date", Label: l.Detail.Summary.AsOfDate, NoSort: true, WidthClass: "col-3xl"},
		{Key: "outcome", Label: l.EntryOutcome.Created, NoSort: true, WidthClass: "col-3xl"},
		{Key: "error_message", Label: l.Errors.InvalidSelection, NoSort: true, WidthClass: "col-5xl"},
	}

	rows := make([]types.TableRow, 0, len(entries))
	for _, e := range entries {
		assetDisplay := e.AssetID
		if e.AssetName != "" {
			assetDisplay = e.AssetName
		}
		outcomeLabel, outcomeVariant := entryOutcomeCell(l, e.Outcome)
		rows = append(rows, types.TableRow{
			ID: e.ID,
			DataAttrs: map[string]string{
				"outcome": e.Outcome,
			},
			Cells: []types.TableCell{
				{Type: "text", Value: assetDisplay},
				{Type: "text", Value: e.PeriodStartDate},
				{Type: "badge", Value: outcomeLabel, Variant: outcomeVariant},
				{Type: "text", Value: e.ErrorMessage},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "depreciation-run-results-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.List.Empty.Complete.Title,
			Message: l.List.Empty.Complete.Message,
		},
	}
}

// buildTransactionsTable builds the TableConfig for the Transactions tab.
func buildTransactionsTable(txs []drshared.AssetTransactionRow, l fycha.DepreciationRunLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "asset", Label: l.List.Columns.Scope, NoSort: true, WidthClass: "col-6xl"},
		{Key: "transaction_date", Label: l.Detail.Summary.InitiatedAt, NoSort: true, WidthClass: "col-3xl"},
		{Key: "period_start_date", Label: l.Detail.Summary.AsOfDate, NoSort: true, WidthClass: "col-3xl"},
		{Key: "amount", Label: l.List.Columns.Created, NoSort: true, WidthClass: "col-3xl", Align: "right"},
	}

	rows := make([]types.TableRow, 0, len(txs))
	for _, tx := range txs {
		assetDisplay := tx.AssetID
		if tx.AssetName != "" {
			assetDisplay = tx.AssetName
		}
		rows = append(rows, types.TableRow{
			ID: tx.ID,
			DataAttrs: map[string]string{
				"transaction-id": tx.ID,
			},
			Cells: []types.TableCell{
				{Type: "text", Value: assetDisplay},
				{Type: "text", Value: tx.TransactionDate},
				{Type: "text", Value: tx.PeriodStartDate},
				types.MoneyCell(float64(tx.Amount), tx.Currency, true),
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	return &types.TableConfig{
		ID:      "depreciation-run-transactions-table",
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.List.Empty.Complete.Title,
			Message: l.List.Empty.Complete.Message,
		},
	}
}

func entryOutcomeCell(l fycha.DepreciationRunLabels, outcome string) (label, variant string) {
	switch outcome {
	case "created":
		return l.EntryOutcome.Created, "success"
	case "skipped":
		return l.EntryOutcome.Skipped, "info"
	case "errored":
		return l.EntryOutcome.Errored, "error"
	default:
		return outcome, "info"
	}
}
