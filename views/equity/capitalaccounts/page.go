// Package capitalaccounts provides the view for the Funding > Equity > Capital Accounts list page.
package capitalaccounts

import (
	"context"
	"fmt"
	"log"

	equityaccountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_account"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the capital accounts list.
type Deps struct {
	Routes       fycha.EquityRoutes
	Labels       fycha.EquityLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Equity account use cases
	GetEquityAccountListPageData func(ctx context.Context, req *equityaccountpb.GetEquityAccountListPageDataRequest) (*equityaccountpb.GetEquityAccountListPageDataResponse, error)
	ListEquityAccounts           func(ctx context.Context, req *equityaccountpb.ListEquityAccountsRequest) (*equityaccountpb.ListEquityAccountsResponse, error)
}

// PageData holds the data for the capital accounts list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	Labels          fycha.EquityLabels
}

// EquityAccountRow is the view-model for a single equity account row.
type EquityAccountRow struct {
	ID          string
	Name        string
	OwnerName   string
	AccountType string
	Balance     string
	Active      bool
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the capital accounts list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		accounts := fetchAccounts(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, accounts, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Accounts.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    deps.Labels.Accounts.Page.Heading,
				HeaderSubtitle: deps.Labels.Accounts.Page.Caption,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "capital-accounts-content",
			Labels:          deps.Labels,
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "equity"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("capital-accounts", pageData)
	})
}

// NewContentView creates the capital accounts list HTMX partial view.
func NewContentView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		accounts := fetchAccounts(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, accounts, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Accounts.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    deps.Labels.Accounts.Page.Heading,
				HeaderSubtitle: deps.Labels.Accounts.Page.Caption,
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "capital-accounts-content",
			Labels:          deps.Labels,
			Table:           tableConfig,
		}

		return view.OK("capital-accounts-content", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchAccounts(ctx context.Context, deps *Deps) []EquityAccountRow {
	if deps.ListEquityAccounts == nil {
		return []EquityAccountRow{}
	}

	resp, err := deps.ListEquityAccounts(ctx, &equityaccountpb.ListEquityAccountsRequest{})
	if err != nil {
		log.Printf("ListEquityAccounts error: %v", err)
		return []EquityAccountRow{}
	}
	if resp == nil {
		return []EquityAccountRow{}
	}

	rows := make([]EquityAccountRow, 0, len(resp.GetData()))
	for _, a := range resp.GetData() {
		rows = append(rows, protoToRow(a))
	}
	return rows
}

func protoToRow(a *equityaccountpb.EquityAccount) EquityAccountRow {
	return EquityAccountRow{
		ID:          a.GetId(),
		Name:        a.GetName(),
		OwnerName:   a.GetOwnerName(),
		AccountType: accountTypeLabel(a.GetAccountType()),
		Balance:     fmt.Sprintf("%.2f", float64(a.GetBalance())/100.0),
		Active:      a.GetActive(),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, accounts []EquityAccountRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels.Accounts
	columns := []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: false},
		{Key: "owner", Label: l.Columns.OwnerName, Sortable: false},
		{Key: "type", Label: l.Columns.AccountType, Sortable: false, Width: "160px"},
		{Key: "balance", Label: l.Columns.Balance, Sortable: false, Width: "140px", Align: "right"},
	}

	rows := []types.TableRow{}
	for _, acct := range accounts {
		canUpdate := perms.Can("equity_account", "update")

		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(deps.Routes.AccountDetailURL, "id", acct.ID)},
			{Type: "edit", Label: l.Actions.Edit, Action: "edit",
				URL: route.ResolveURL(deps.Routes.AccountDetailURL, "id", acct.ID), DrawerTitle: l.Actions.Edit,
				Disabled: !canUpdate, DisabledTooltip: l.Actions.NoPermission},
		}

		typeVariant := accountTypeBadgeVariant(acct.AccountType)

		row := types.TableRow{
			ID:   acct.ID,
			Href: route.ResolveURL(deps.Routes.AccountDetailURL, "id", acct.ID),
			Cells: []types.TableCell{
				{Type: "text", Value: acct.Name},
				{Type: "text", Value: acct.OwnerName},
				{Type: "badge", Value: acct.AccountType, Variant: typeVariant},
				{Type: "money", Value: acct.Balance},
			},
			Actions: actions,
		}
		rows = append(rows, row)
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                "capital-accounts-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowEntries:       true,
		DefaultSortColumn: "name",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddAccount,
			ActionURL:       deps.Routes.AccountsURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("equity_account", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Proto enum helpers
// ---------------------------------------------------------------------------

func accountTypeLabel(t equityaccountpb.EquityAccountType) string {
	switch t {
	case equityaccountpb.EquityAccountType_EQUITY_ACCOUNT_TYPE_OWNERS_CAPITAL:
		return "Owner's Capital"
	case equityaccountpb.EquityAccountType_EQUITY_ACCOUNT_TYPE_OWNERS_DRAW:
		return "Owner's Draw"
	case equityaccountpb.EquityAccountType_EQUITY_ACCOUNT_TYPE_RETAINED_EARNINGS:
		return "Retained Earnings"
	case equityaccountpb.EquityAccountType_EQUITY_ACCOUNT_TYPE_ADDITIONAL_PAID_IN_CAPITAL:
		return "Additional Paid-In Capital"
	default:
		return "Unspecified"
	}
}

func accountTypeBadgeVariant(label string) string {
	switch label {
	case "Owner's Capital", "Additional Paid-In Capital":
		return "navy"
	case "Owner's Draw":
		return "amber"
	case "Retained Earnings":
		return "sage"
	default:
		return "default"
	}
}
