package cash

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// MockDeposit represents a security deposit record for mock data display.
type MockDeposit struct {
	ID               string
	CounterpartyName string
	Direction        string // "DEPOSIT_PAID" or "DEPOSIT_RECEIVED"
	Amount           float64
	DepositDateString string
	Status           string
	Notes            string
}

// DepositDeps holds view dependencies for deposit views.
type DepositDeps struct {
	Routes       fycha.DepositRoutes
	Labels       fycha.DepositLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// DepositPageData holds the data for the deposits list page.
type DepositPageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// NewDepositsView creates the security deposits list view (full page).
func NewDepositsView(deps *DepositDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "all"
		}

		perms := view.GetUserPermissions(ctx)
		tableConfig := buildDepositTableConfig(deps, status, perms)

		pageData := &DepositPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "cash",
				ActiveSubNav:   "cash-deposits",
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-shield",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "deposits-content",
			Table:           tableConfig,
		}

		return view.OK("deposits", pageData)
	})
}

// mockDeposits returns hardcoded deposit data for initial UI development.
func mockDeposits() []MockDeposit {
	return []MockDeposit{
		{ID: "dep-001", CounterpartyName: "Prime Properties LLC", Direction: "DEPOSIT_PAID", Amount: 150000, DepositDateString: "2026-01-15", Status: "DEPOSIT_HELD", Notes: "Office lease deposit"},
		{ID: "dep-002", CounterpartyName: "TechSupply Inc.", Direction: "DEPOSIT_PAID", Amount: 25000, DepositDateString: "2025-11-01", Status: "DEPOSIT_RETURNED", Notes: "Equipment rental deposit"},
		{ID: "dep-003", CounterpartyName: "XYZ Corp", Direction: "DEPOSIT_RECEIVED", Amount: 50000, DepositDateString: "2026-02-01", Status: "DEPOSIT_HELD", Notes: "Service contract deposit"},
		{ID: "dep-004", CounterpartyName: "ABC Retailer", Direction: "DEPOSIT_RECEIVED", Amount: 20000, DepositDateString: "2025-09-15", Status: "DEPOSIT_FORFEITED", Notes: "Customer cancelled contract"},
		{ID: "dep-005", CounterpartyName: "City Mall", Direction: "DEPOSIT_PAID", Amount: 200000, DepositDateString: "2026-03-01", Status: "DEPOSIT_HELD", Notes: "Retail space deposit"},
	}
}

func buildDepositTableConfig(deps *DepositDeps, statusFilter string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := depositColumns(l)

	filtered := filterDeposits(mockDeposits(), statusFilter)
	rows := buildDepositRows(filtered, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	bulkCfg := fycha.MapBulkConfig(deps.CommonLabels)

	tableConfig := &types.TableConfig{
		ID:                   "deposits-table",
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
		DefaultSortColumn:    "deposit_date",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.RecordDeposit,
			ActionURL:       deps.Routes.ListURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("security_deposit", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
		BulkActions: &bulkCfg,
	}
	types.ApplyTableSettings(tableConfig)
	return tableConfig
}

func depositColumns(l fycha.DepositLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "counterparty", Label: l.Columns.Counterparty, Sortable: true},
		{Key: "direction", Label: l.Columns.Direction, Sortable: true, Width: "160px"},
		{Key: "amount", Label: l.Columns.Amount, Sortable: true, Width: "140px"},
		{Key: "deposit_date", Label: l.Columns.DepositDate, Sortable: true, Width: "130px"},
		{Key: "status", Label: l.Columns.Status, Sortable: true, Width: "130px"},
		{Key: "notes", Label: l.Columns.Notes, Sortable: false},
	}
}

func filterDeposits(items []MockDeposit, statusFilter string) []MockDeposit {
	if statusFilter == "all" || statusFilter == "" {
		return items
	}
	var filtered []MockDeposit
	for _, d := range items {
		if statusFilter == "paid" && d.Direction == "DEPOSIT_PAID" {
			filtered = append(filtered, d)
		} else if statusFilter == "received" && d.Direction == "DEPOSIT_RECEIVED" {
			filtered = append(filtered, d)
		}
	}
	return filtered
}

func buildDepositRows(items []MockDeposit, l fycha.DepositLabels, routes fycha.DepositRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, item := range items {
		id := item.ID
		canDelete := perms.Can("security_deposit", "delete")

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: routes.ListURL},
			{
				Type: "delete", Label: l.Actions.Delete, Action: "delete",
				URL: routes.ListURL, ItemName: item.CounterpartyName,
				ConfirmTitle:   l.Actions.Delete,
				ConfirmMessage: l.Actions.ConfirmDelete,
				Disabled:       !canDelete, DisabledTooltip: l.Actions.NoPermission,
			},
		}

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: item.CounterpartyName},
				{Type: "badge", Value: depositDirectionLabel(l, item.Direction), Variant: depositDirectionVariant(item.Direction)},
				{Type: "text", Value: formatDepositCurrency(item.Amount)},
				{Type: "text", Value: item.DepositDateString},
				{Type: "badge", Value: depositStatusLabel(l, item.Status), Variant: depositStatusVariant(item.Status)},
				{Type: "text", Value: item.Notes},
			},
			DataAttrs: map[string]string{
				"counterparty": item.CounterpartyName,
				"direction":    item.Direction,
				"status":       item.Status,
			},
			Actions: actions,
		})
	}
	return rows
}

func depositDirectionLabel(l fycha.DepositLabels, direction string) string {
	switch direction {
	case "DEPOSIT_PAID":
		return l.Form.DirectionPaid
	case "DEPOSIT_RECEIVED":
		return l.Form.DirectionReceived
	default:
		return direction
	}
}

func depositDirectionVariant(direction string) string {
	switch direction {
	case "DEPOSIT_PAID":
		return "warning"
	case "DEPOSIT_RECEIVED":
		return "info"
	default:
		return "default"
	}
}

func depositStatusLabel(l fycha.DepositLabels, status string) string {
	switch status {
	case "DEPOSIT_HELD":
		return l.Status.Held
	case "DEPOSIT_RETURNED":
		return l.Status.Returned
	case "DEPOSIT_FORFEITED":
		return l.Status.Forfeited
	default:
		return status
	}
}

func depositStatusVariant(status string) string {
	switch status {
	case "DEPOSIT_HELD":
		return "info"
	case "DEPOSIT_RETURNED":
		return "success"
	case "DEPOSIT_FORFEITED":
		return "danger"
	default:
		return "default"
	}
}

func formatDepositCurrency(amount float64) string {
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
