package cash

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// MockPettyCashFund represents a petty cash fund record for mock data display.
type MockPettyCashFund struct {
	ID               string
	Name             string
	AuthorizedAmount float64
	CurrentBalance   float64
	CustodianName    string
	LocationName     string
	Active           bool
}

// PettyCashDeps holds view dependencies for petty cash views.
type PettyCashDeps struct {
	Routes       fycha.PettyCashRoutes
	Labels       fycha.PettyCashLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PettyCashPageData holds the data for petty cash pages.
type PettyCashPageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewPettyCashRegisterView creates the petty cash register view (full page).
func NewPettyCashRegisterView(deps *PettyCashDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildPettyCashRegisterTableConfig(deps, perms)

		pageData := &PettyCashPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.RegisterHeading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				ActiveSubNav:   "cash-petty-cash",
				HeaderTitle:    deps.Labels.Page.RegisterHeading,
				HeaderSubtitle: deps.Labels.Page.RegisterCaption,
				HeaderIcon:     "icon-dollar-sign",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "petty-cash-register-content",
			Table:           tableConfig,
		}

		return view.OK("petty-cash-register", pageData)
	})
}

// NewPettyCashReplenishmentsView creates the petty cash replenishments view.
func NewPettyCashReplenishmentsView(deps *PettyCashDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		tableConfig := buildReplenishmentTableConfig(deps)

		pageData := &PettyCashPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.ReplenishmentsHeading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				ActiveSubNav:   "cash-petty-cash",
				HeaderTitle:    deps.Labels.Page.ReplenishmentsHeading,
				HeaderSubtitle: deps.Labels.Page.ReplenishmentsCaption,
				HeaderIcon:     "icon-refresh-cw",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "petty-cash-replenishments-content",
			Table:           tableConfig,
		}

		return view.OK("petty-cash-replenishments", pageData)
	})
}

// NewCustodianBalancesView creates the custodian balances view.
func NewCustodianBalancesView(deps *PettyCashDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		tableConfig := buildCustodianBalancesTableConfig(deps)

		pageData := &PettyCashPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.CustodianBalancesHeading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				ActiveSubNav:   "cash-petty-cash",
				HeaderTitle:    deps.Labels.Page.CustodianBalancesHeading,
				HeaderSubtitle: deps.Labels.Page.CustodianBalancesCaption,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "custodian-balances-content",
			Table:           tableConfig,
		}

		return view.OK("custodian-balances", pageData)
	})
}

// mockPettyCashFunds returns hardcoded petty cash fund data for initial UI development.
func mockPettyCashFunds() []MockPettyCashFund {
	return []MockPettyCashFund{
		{ID: "pcf-001", Name: "Main Office Petty Cash", AuthorizedAmount: 10000, CurrentBalance: 4500, CustodianName: "Maria Santos", LocationName: "Main Office", Active: true},
		{ID: "pcf-002", Name: "Branch 1 Petty Cash", AuthorizedAmount: 5000, CurrentBalance: 2200, CustodianName: "Juan Cruz", LocationName: "Branch 1", Active: true},
		{ID: "pcf-003", Name: "Branch 2 Petty Cash", AuthorizedAmount: 5000, CurrentBalance: 5000, CustodianName: "Ana Reyes", LocationName: "Branch 2", Active: true},
		{ID: "pcf-004", Name: "Warehouse Petty Cash", AuthorizedAmount: 3000, CurrentBalance: 0, CustodianName: "Pedro Lim", LocationName: "Warehouse", Active: false},
	}
}

func buildPettyCashRegisterTableConfig(deps *PettyCashDeps, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := pettyCashRegisterColumns(l)
	rows := buildPettyCashRegisterRows(mockPettyCashFunds(), l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := fycha.MapBulkConfig(deps.CommonLabels)

	tableConfig := &types.TableConfig{
		ID:                   "petty-cash-register-table",
		RefreshURL:           deps.Routes.RegisterURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.RegisterTitle,
			Message: l.Empty.RegisterMessage,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddFund,
			ActionURL:       deps.Routes.RegisterURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("petty_cash_fund", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func pettyCashRegisterColumns(l fycha.PettyCashLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: true},
		{Key: "authorized_amount", Label: l.Columns.AuthorizedAmount, Sortable: true, Width: "160px"},
		{Key: "current_balance", Label: l.Columns.CurrentBalance, Sortable: true, Width: "150px"},
		{Key: "custodian", Label: l.Columns.Custodian, Sortable: true},
		{Key: "location", Label: l.Columns.Location, Sortable: true},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "110px"},
	}
}

func buildPettyCashRegisterRows(items []MockPettyCashFund, l fycha.PettyCashLabels, routes fycha.PettyCashRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		id := item.ID
		canDelete := perms.Can("petty_cash_fund", "delete")

		statusLabel := l.Status.Active
		statusVariant := "success"
		if !item.Active {
			statusLabel = l.Status.Inactive
			statusVariant = "warning"
		}

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: routes.RegisterURL},
			{Type: "edit", Label: l.Actions.Replenish, Action: "replenish", URL: routes.ReplenishmentListURL, DrawerTitle: l.Actions.Replenish, Disabled: !perms.Can("petty_cash_fund", "update"), DisabledTooltip: l.Actions.NoPermission},
			{
				Type: "delete", Label: l.Actions.Delete, Action: "delete",
				URL: routes.RegisterURL, ItemName: item.Name,
				ConfirmTitle:   l.Actions.Delete,
				ConfirmMessage: l.Actions.ConfirmDelete,
				Disabled:       !canDelete, DisabledTooltip: l.Actions.NoPermission,
			},
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: item.Name},
				{Type: "text", Value: formatPCCurrency(item.AuthorizedAmount)},
				{Type: "text", Value: formatPCCurrency(item.CurrentBalance)},
				{Type: "text", Value: item.CustodianName},
				{Type: "text", Value: item.LocationName},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			DataAttrs: map[string]string{
				"name":      item.Name,
				"custodian": item.CustodianName,
				"location":  item.LocationName,
			},
			Actions: actions,
		})
	}
	return rows
}

// buildReplenishmentTableConfig builds the replenishment history table config.
func buildReplenishmentTableConfig(deps *PettyCashDeps) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "fund", Label: l.Columns.Fund, Sortable: true},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, Width: "140px"},
		{Key: "date", Label: l.Columns.Date, Sortable: true, Width: "130px"},
		{Key: "notes", Label: l.Columns.Notes, Sortable: false},
	}

	type replenRow struct {
		FundName string
		Amount   float64
		Date     string
		Notes    string
	}
	mockRows := []replenRow{
		{FundName: "Main Office Petty Cash", Amount: 5500, Date: "2026-03-01", Notes: "Monthly replenishment"},
		{FundName: "Branch 1 Petty Cash", Amount: 2800, Date: "2026-03-01", Notes: "Monthly replenishment"},
		{FundName: "Main Office Petty Cash", Amount: 4800, Date: "2026-02-01", Notes: "Monthly replenishment"},
	}

	rows := []types.TableRow{}
	for i, r := range mockRows {
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("replen-%d", i),
			Cells: []types.TableCell{
				{Type: "text", Value: r.FundName},
				{Type: "text", Value: formatPCCurrency(r.Amount)},
				{Type: "text", Value: r.Date},
				{Type: "text", Value: r.Notes},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "replenishments-table",
		RefreshURL:           deps.Routes.ReplenishmentListURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "date",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.ReplenishmentsTitle,
			Message: l.Empty.ReplenishmentsMessage,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

// buildCustodianBalancesTableConfig builds the custodian balance summary table config.
func buildCustodianBalancesTableConfig(deps *PettyCashDeps) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "custodian", Label: l.Columns.Custodian, Sortable: true},
		{Key: "location", Label: l.Columns.Location, Sortable: true},
		{Key: "total_funds", Label: l.Columns.TotalFunds, Sortable: true, Width: "120px"},
		{Key: "total_balance", Label: l.Columns.TotalBalance, Sortable: true, Width: "150px"},
	}

	type custodianRow struct {
		Custodian    string
		Location     string
		TotalFunds   int
		TotalBalance float64
	}
	mockRows := []custodianRow{
		{Custodian: "Maria Santos", Location: "Main Office", TotalFunds: 1, TotalBalance: 4500},
		{Custodian: "Juan Cruz", Location: "Branch 1", TotalFunds: 1, TotalBalance: 2200},
		{Custodian: "Ana Reyes", Location: "Branch 2", TotalFunds: 1, TotalBalance: 5000},
	}

	rows := []types.TableRow{}
	for i, r := range mockRows {
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("cust-%d", i),
			Cells: []types.TableCell{
				{Type: "text", Value: r.Custodian},
				{Type: "text", Value: r.Location},
				{Type: "text", Value: fmt.Sprintf("%d", r.TotalFunds)},
				{Type: "text", Value: formatPCCurrency(r.TotalBalance)},
			},
		})
	}
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "custodian-balances-table",
		RefreshURL:           deps.Routes.CustodianBalancesURL,
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           true,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "custodian",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.CustodianTitle,
			Message: l.Empty.CustodianMessage,
		},
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func formatPCCurrency(amount float64) string {
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
