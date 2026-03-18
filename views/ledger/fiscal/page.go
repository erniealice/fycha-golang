package fiscal

import (
	"context"
	"fmt"
	"log"

	fiscalpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the fiscal period list page.
type Deps struct {
	Routes       fycha.FiscalPeriodRoutes
	Labels       fycha.FiscalPeriodLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Fiscal period use cases
	GetFiscalPeriodListPageData func(ctx context.Context) ([]*fiscalpb.FiscalPeriod, error)
}

// PageData holds the data for the fiscal period list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// FiscalPeriodRow is the view-model for a single fiscal period row.
type FiscalPeriodRow struct {
	ID           string
	Name         string
	PeriodNumber int
	FiscalYear   int
	StartDate    string
	EndDate      string
	Status       string // "open", "closed", "locked"
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the fiscal period list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		periods := fetchPeriods(ctx, deps)

		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, periods, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-calendar",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "fiscal-periods-content",
			Table:           tableConfig,
		}

		return view.OK("fiscal-periods", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

// fetchPeriods calls the use case and converts the response to view-model rows.
// Falls back to mock data when no use case is wired (placeholder mode).
func fetchPeriods(ctx context.Context, deps *Deps) []FiscalPeriodRow {
	if deps.GetFiscalPeriodListPageData == nil {
		return mockPeriods()
	}

	periods, err := deps.GetFiscalPeriodListPageData(ctx)
	if err != nil {
		log.Printf("GetFiscalPeriodListPageData error: %v", err)
		return mockPeriods()
	}

	rows := make([]FiscalPeriodRow, 0, len(periods))
	for _, p := range periods {
		rows = append(rows, protoToRow(p))
	}
	return rows
}

// protoToRow converts a proto FiscalPeriod to a view-model FiscalPeriodRow.
func protoToRow(p *fiscalpb.FiscalPeriod) FiscalPeriodRow {
	return FiscalPeriodRow{
		ID:           p.GetId(),
		Name:         p.GetName(),
		PeriodNumber: int(p.GetPeriodNumber()),
		FiscalYear:   int(p.GetFiscalYear()),
		StartDate:    p.GetStartDateString(),
		EndDate:      p.GetEndDateString(),
		Status:       statusString(p.GetStatus()),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, periods []FiscalPeriodRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := periodColumns(l)
	rows := buildTableRows(periods, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                "fiscal-periods-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        false,
		ShowActions:       true,
		ShowExport:        false,
		ShowEntries:       true,
		DefaultSortColumn: "period",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddPeriod,
			ActionURL:       "#",
			Icon:            "icon-plus",
			Disabled:        !perms.Can("fiscal_period", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func periodColumns(l fycha.FiscalPeriodLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "period", Label: l.Columns.Period, Sortable: false},
		{Key: "year", Label: l.Columns.Year, Sortable: false, Width: "90px"},
		{Key: "start_date", Label: l.Columns.StartDate, Sortable: false, Width: "120px"},
		{Key: "end_date", Label: l.Columns.EndDate, Sortable: false, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: false, Width: "110px"},
	}
}

func buildTableRows(periods []FiscalPeriodRow, l fycha.FiscalPeriodLabels, routes fycha.FiscalPeriodRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, p := range periods {
		canClose := perms.Can("fiscal_period", "update") && p.Status == "open"

		var actions []types.TableAction
		if p.Status == "open" {
			closeURL := route.ResolveURL(routes.CloseURL, "id", p.ID)
			actions = append(actions, types.TableAction{
				Type:            "edit",
				Label:           l.Actions.Close,
				Action:          "post",
				URL:             closeURL,
				DrawerTitle:     l.Actions.Close,
				Disabled:        !canClose,
				DisabledTooltip: l.Actions.NoPermission,
			})
		}

		statusVariant := statusBadgeVariant(p.Status)
		statusLabel := statusDisplayLabel(l, p.Status)

		row := types.TableRow{
			ID: p.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: p.Name},
				{Type: "text", Value: fmt.Sprintf("%d", p.FiscalYear)},
				{Type: "text", Value: p.StartDate},
				{Type: "text", Value: p.EndDate},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			Actions: actions,
		}
		rows = append(rows, row)
	}
	return rows
}

// ---------------------------------------------------------------------------
// Mock data (placeholder until DB is wired)
// ---------------------------------------------------------------------------

func mockPeriods() []FiscalPeriodRow {
	return []FiscalPeriodRow{
		{ID: "fp-01", Name: "April 2025", PeriodNumber: 1, FiscalYear: 2026, StartDate: "2025-04-01", EndDate: "2025-04-30", Status: "closed"},
		{ID: "fp-02", Name: "May 2025", PeriodNumber: 2, FiscalYear: 2026, StartDate: "2025-05-01", EndDate: "2025-05-31", Status: "closed"},
		{ID: "fp-03", Name: "June 2025", PeriodNumber: 3, FiscalYear: 2026, StartDate: "2025-06-01", EndDate: "2025-06-30", Status: "closed"},
		{ID: "fp-04", Name: "July 2025", PeriodNumber: 4, FiscalYear: 2026, StartDate: "2025-07-01", EndDate: "2025-07-31", Status: "closed"},
		{ID: "fp-05", Name: "August 2025", PeriodNumber: 5, FiscalYear: 2026, StartDate: "2025-08-01", EndDate: "2025-08-31", Status: "closed"},
		{ID: "fp-06", Name: "September 2025", PeriodNumber: 6, FiscalYear: 2026, StartDate: "2025-09-01", EndDate: "2025-09-30", Status: "closed"},
		{ID: "fp-07", Name: "October 2025", PeriodNumber: 7, FiscalYear: 2026, StartDate: "2025-10-01", EndDate: "2025-10-31", Status: "closed"},
		{ID: "fp-08", Name: "November 2025", PeriodNumber: 8, FiscalYear: 2026, StartDate: "2025-11-01", EndDate: "2025-11-30", Status: "closed"},
		{ID: "fp-09", Name: "December 2025", PeriodNumber: 9, FiscalYear: 2026, StartDate: "2025-12-01", EndDate: "2025-12-31", Status: "open"},
		{ID: "fp-10", Name: "January 2026", PeriodNumber: 10, FiscalYear: 2026, StartDate: "2026-01-01", EndDate: "2026-01-31", Status: "open"},
		{ID: "fp-11", Name: "February 2026", PeriodNumber: 11, FiscalYear: 2026, StartDate: "2026-02-01", EndDate: "2026-02-28", Status: "open"},
		{ID: "fp-12", Name: "March 2026", PeriodNumber: 12, FiscalYear: 2026, StartDate: "2026-03-01", EndDate: "2026-03-31", Status: "open"},
	}
}

// ---------------------------------------------------------------------------
// Proto enum → display string converters
// ---------------------------------------------------------------------------

func statusString(s fiscalpb.FiscalPeriodStatus) string {
	switch s {
	case fiscalpb.FiscalPeriodStatus_FISCAL_PERIOD_STATUS_OPEN:
		return "open"
	case fiscalpb.FiscalPeriodStatus_FISCAL_PERIOD_STATUS_CLOSED:
		return "closed"
	case fiscalpb.FiscalPeriodStatus_FISCAL_PERIOD_STATUS_LOCKED:
		return "locked"
	default:
		return "open"
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func statusBadgeVariant(status string) string {
	switch status {
	case "open":
		return "success"
	case "closed":
		return "muted"
	case "locked":
		return "warning"
	default:
		return "default"
	}
}

func statusDisplayLabel(l fycha.FiscalPeriodLabels, status string) string {
	switch status {
	case "open":
		return l.Status.Open
	case "closed":
		return l.Status.Closed
	case "locked":
		return l.Status.Locked
	default:
		return status
	}
}
