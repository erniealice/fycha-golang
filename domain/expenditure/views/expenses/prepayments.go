package expenses

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	expenditure "github.com/erniealice/fycha-golang/domain/expenditure"
)

// MockPrepayment represents a prepayment record for mock data display.
type MockPrepayment struct {
	ID                 string
	Description        string
	VendorName         string
	TotalAmount        float64
	RemainingAmount    float64
	AmortizationMonths int
	StartDateString    string
	EndDateString      string
	Status             string
}

// PrepaymentDeps holds view dependencies for prepayment views.
type PrepaymentDeps struct {
	Routes       expenditure.PrepaymentRoutes
	Labels       expenditure.PrepaymentLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PrepaymentPageData holds the data for the prepayments list page.
type PrepaymentPageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewPrepaymentsView creates the prepayments list view (full page).
func NewPrepaymentsView(deps *PrepaymentDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildPrepaymentTableConfig(deps, perms)

		pageData := &PrepaymentPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "expense",
				ActiveSubNav:   "expenses-prepayments",
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "prepayments-content",
			Table:           tableConfig,
		}

		return view.OK("prepayments", pageData)
	})
}

// NewAmortizationScheduleView creates the amortization schedule view.
func NewAmortizationScheduleView(deps *PrepaymentDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		tableConfig := buildAmortizationTableConfig(deps)

		pageData := &PrepaymentPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.AmortizationHeading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "expense",
				ActiveSubNav:   "expenses-prepayments",
				HeaderTitle:    deps.Labels.Page.AmortizationHeading,
				HeaderSubtitle: deps.Labels.Page.AmortizationCaption,
				HeaderIcon:     "icon-calendar",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "prepayment-amortization-content",
			Table:           tableConfig,
		}

		return view.OK("prepayment-amortization", pageData)
	})
}

// mockPrepayments returns hardcoded prepayment data for initial UI development.
func mockPrepayments() []MockPrepayment {
	return []MockPrepayment{
		{ID: "pre-001", Description: "Annual Office Insurance", VendorName: "SafeGuard Insurance", TotalAmount: 120000, RemainingAmount: 80000, AmortizationMonths: 12, StartDateString: "2026-01-01", EndDateString: "2026-12-31", Status: "PREPAYMENT_ACTIVE"},
		{ID: "pre-002", Description: "Software License (Annual)", VendorName: "TechCorp Inc.", TotalAmount: 60000, RemainingAmount: 45000, AmortizationMonths: 12, StartDateString: "2026-01-01", EndDateString: "2026-12-31", Status: "PREPAYMENT_ACTIVE"},
		{ID: "pre-003", Description: "Advance Rent Deposit", VendorName: "Prime Properties LLC", TotalAmount: 200000, RemainingAmount: 0, AmortizationMonths: 6, StartDateString: "2025-07-01", EndDateString: "2025-12-31", Status: "PREPAYMENT_AMORTIZED"},
		{ID: "pre-004", Description: "3-Year Maintenance Contract", VendorName: "ServicePro Corp", TotalAmount: 360000, RemainingAmount: 270000, AmortizationMonths: 36, StartDateString: "2026-01-01", EndDateString: "2028-12-31", Status: "PREPAYMENT_ACTIVE"},
	}
}

func buildPrepaymentTableConfig(deps *PrepaymentDeps, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := prepaymentColumns(l)
	rows := buildPrepaymentRows(mockPrepayments(), l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := pyeza.MapBulkConfig(deps.CommonLabels)

	tableConfig := &types.TableConfig{
		ID:                   "prepayments-table",
		RefreshURL:           deps.Routes.ListURL,
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
		DefaultSortColumn:    "description",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddPrepayment,
			ActionURL:       deps.Routes.ListURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("prepayment", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func prepaymentColumns(l expenditure.PrepaymentLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "description", Label: l.Columns.Description},
		{Key: "vendor", Label: l.Columns.Vendor},
		{Key: "total_amount", Label: l.Columns.TotalAmount, WidthClass: "col-4xl"},
		{Key: "remaining_amount", Label: l.Columns.RemainingAmount, WidthClass: "col-3xl"},
		{Key: "amortization_months", Label: l.Columns.AmortizationMonths, WidthClass: "col-lg"},
		{Key: "start_date", Label: l.Columns.StartDate, WidthClass: "col-2xl"},
		{Key: "end_date", Label: l.Columns.EndDate, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, WidthClass: "col-3xl"},
	}
}

func buildPrepaymentRows(items []MockPrepayment, l expenditure.PrepaymentLabels, routes expenditure.PrepaymentRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		id := item.ID
		canDelete := perms.Can("prepayment", "delete")

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.AmortizationURL, "id", id)},
			{
				Type: "delete", Label: l.Actions.Delete, Action: "delete",
				URL: routes.ListURL, ItemName: item.Description,
				ConfirmTitle:   l.Actions.Delete,
				ConfirmMessage: l.Actions.ConfirmDelete,
				Disabled:       !canDelete, DisabledTooltip: l.Actions.NoPermission,
			},
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: item.Description},
				{Type: "text", Value: item.VendorName},
				types.MoneyCell(item.TotalAmount, "PHP", false),
				types.MoneyCell(item.RemainingAmount, "PHP", false),
				{Type: "text", Value: fmt.Sprintf("%d mo", item.AmortizationMonths)},
				{Type: "text", Value: item.StartDateString},
				{Type: "text", Value: item.EndDateString},
				{Type: "badge", Value: prepaymentStatusLabel(l, item.Status), Variant: prepaymentStatusVariant(item.Status)},
			},
			DataAttrs: map[string]string{
				"description": item.Description,
				"vendor":      item.VendorName,
				"status":      item.Status,
			},
			Actions: actions,
		})
	}
	return rows
}

// buildAmortizationTableConfig builds the amortization schedule table config with mock data.
func buildAmortizationTableConfig(deps *PrepaymentDeps) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "description", Label: l.Columns.Description},
		{Key: "vendor", Label: l.Columns.Vendor},
		{Key: "month", Label: l.Columns.Month, WidthClass: "col-2xl"},
		{Key: "opening", Label: l.Columns.Opening, NoSort: true, WidthClass: "col-3xl"},
		{Key: "expense", Label: l.Columns.Expense, NoSort: true, WidthClass: "col-3xl"},
		{Key: "closing", Label: l.Columns.Closing, NoSort: true, WidthClass: "col-3xl"},
	}

	// Mock amortization schedule rows
	type amortRow struct {
		Description string
		Vendor      string
		Month       string
		Opening     float64
		Expense     float64
		Closing     float64
	}
	mockRows := []amortRow{
		{Description: "Annual Office Insurance", Vendor: "SafeGuard Insurance", Month: "Jan 2026", Opening: 120000, Expense: 10000, Closing: 110000},
		{Description: "Annual Office Insurance", Vendor: "SafeGuard Insurance", Month: "Feb 2026", Opening: 110000, Expense: 10000, Closing: 100000},
		{Description: "Software License (Annual)", Vendor: "TechCorp Inc.", Month: "Jan 2026", Opening: 60000, Expense: 5000, Closing: 55000},
		{Description: "Software License (Annual)", Vendor: "TechCorp Inc.", Month: "Feb 2026", Opening: 55000, Expense: 5000, Closing: 50000},
	}

	rows := []types.TableRow{}
	for i, r := range mockRows {
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("amort-%d", i),
			Cells: []types.TableCell{
				{Type: "text", Value: r.Description},
				{Type: "text", Value: r.Vendor},
				{Type: "text", Value: r.Month},
				types.MoneyCell(r.Opening, "PHP", false),
				types.MoneyCell(r.Expense, "PHP", false),
				types.MoneyCell(r.Closing, "PHP", false),
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "amortization-table",
		RefreshURL:           deps.Routes.AmortizationURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "month",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func prepaymentStatusLabel(l expenditure.PrepaymentLabels, status string) string {
	switch status {
	case "PREPAYMENT_ACTIVE":
		return l.Status.Active
	case "PREPAYMENT_AMORTIZED":
		return l.Status.Amortized
	case "PREPAYMENT_CANCELLED":
		return l.Status.Cancelled
	default:
		return status
	}
}

func prepaymentStatusVariant(status string) string {
	switch status {
	case "PREPAYMENT_ACTIVE":
		return "success"
	case "PREPAYMENT_AMORTIZED":
		return "default"
	case "PREPAYMENT_CANCELLED":
		return "danger"
	default:
		return "default"
	}
}
