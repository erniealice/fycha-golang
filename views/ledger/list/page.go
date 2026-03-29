package list

import (
	"context"
	"log"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// ListViewDeps holds view dependencies.
type ListViewDeps struct {
	Routes       fycha.AccountRoutes
	Labels       fycha.AccountLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Account use cases
	GetAccountListPageData func(ctx context.Context, req *accountpb.GetAccountListPageDataRequest) (*accountpb.GetAccountListPageDataResponse, error)
}

// PageData holds the data for the account list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveElement   string          // "all", "asset", "liability", "equity", "revenue", "expense"
	ElementTabs     []pyeza.TabItem // category tabs for element filter
	Table           *types.TableConfig
}

// AccountRow is the view-model for a single account row (mapped from proto).
type AccountRow struct {
	ID             string
	Code           string
	Name           string
	Element        string // "asset", "liability", "equity", "revenue", "expense"
	Classification string
	IsGroup        bool
	Level          int
	ParentCode     string
	Balance        string
	Active         bool
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the account list view (full page).
func NewView(deps *ListViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		element := viewCtx.Request.URL.Query().Get("tab")
		if element == "" {
			element = "all"
		}

		accounts := fetchAccounts(ctx, deps)

		perms := view.GetUserPermissions(ctx)
		elementTabs := buildElementTabs(deps)
		tableConfig := buildTableConfig(deps, element, accounts, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-book-open",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "account-list-content",
			ActiveElement:   element,
			ElementTabs:     elementTabs,
			Table:           tableConfig,
		}

		return view.OK("account-list", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

// fetchAccounts calls the use case and converts the response to view-model rows.
// Falls back to empty slice on error (shows empty state rather than crashing).
// Element filtering is applied at the view layer after fetching.
func fetchAccounts(ctx context.Context, deps *ListViewDeps) []AccountRow {
	if deps.GetAccountListPageData == nil {
		return []AccountRow{}
	}

	resp, err := deps.GetAccountListPageData(ctx, &accountpb.GetAccountListPageDataRequest{})
	if err != nil {
		log.Printf("GetAccountListPageData error: %v", err)
		return []AccountRow{}
	}
	if resp == nil || !resp.GetSuccess() {
		return []AccountRow{}
	}

	rows := make([]AccountRow, 0, len(resp.GetAccountList()))
	for _, a := range resp.GetAccountList() {
		rows = append(rows, protoToRow(a, deps.Labels))
	}
	return rows
}

// protoToRow converts a proto Account to a view-model AccountRow.
func protoToRow(a *accountpb.Account, l fycha.AccountLabels) AccountRow {
	return AccountRow{
		ID:             a.GetId(),
		Code:           a.GetCode(),
		Name:           a.GetName(),
		Element:        elementString(a.GetElement()),
		Classification: classificationLabel(a.GetClassification(), l.Form),
		IsGroup:        false, // leaf accounts returned by use case
		Level:          3,
		ParentCode:     a.GetParentId(),
		Balance:        "", // balance not yet in proto; placeholder
		Active:         a.GetActive(),
	}
}

// ---------------------------------------------------------------------------
// Tab builder
// ---------------------------------------------------------------------------

func buildElementTabs(deps *ListViewDeps) []pyeza.TabItem {
	l := deps.Labels.Tabs
	base := deps.Routes.ListURL
	return []pyeza.TabItem{
		{Key: "all", Label: l.All, Href: base, Icon: "", Count: 0, Disabled: false},
		{Key: "asset", Label: l.Asset, Href: base + "?tab=asset", Icon: "", Count: 0, Disabled: false},
		{Key: "liability", Label: l.Liability, Href: base + "?tab=liability", Icon: "", Count: 0, Disabled: false},
		{Key: "equity", Label: l.Equity, Href: base + "?tab=equity", Icon: "", Count: 0, Disabled: false},
		{Key: "revenue", Label: l.Revenue, Href: base + "?tab=revenue", Icon: "", Count: 0, Disabled: false},
		{Key: "expense", Label: l.Expense, Href: base + "?tab=expense", Icon: "", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *ListViewDeps, element string, accounts []AccountRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := accountColumns(l)
	rows := buildTableRows(accounts, element, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                "accounts-tree-table",
		RefreshURL:        deps.Routes.ListURL,
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowExport:        true,
		ShowEntries:       true,
		DefaultSortColumn: "code",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddAccount,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("account", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func accountColumns(l fycha.AccountLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "code", Label: l.Columns.Code, Sortable: false, Width: "100px"},
		{Key: "name", Label: l.Columns.Name, Sortable: false},
		{Key: "element", Label: l.Columns.Element, Sortable: false, Width: "110px"},
		{Key: "class", Label: l.Columns.Classification, Sortable: false, Width: "120px"},
		{Key: "balance", Label: l.Columns.Balance, Sortable: false, Width: "140px", Align: "right"},
	}
}

func buildTableRows(accounts []AccountRow, element string, l fycha.AccountLabels, routes fycha.AccountRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, acct := range accounts {
		if element != "all" && acct.Element != element {
			continue
		}

		canUpdate := perms.Can("account", "update")
		canDelete := perms.Can("account", "delete")

		var actions []types.TableAction
		if !acct.IsGroup {
			actions = append(actions,
				types.TableAction{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(routes.DetailURL, "id", acct.ID)},
				types.TableAction{
					Type: "edit", Label: l.Actions.Edit, Action: "edit",
					URL: route.ResolveURL(routes.EditURL, "id", acct.ID), DrawerTitle: l.Actions.Edit,
					Disabled: !canUpdate, DisabledTooltip: l.Actions.NoPermission,
				},
				types.TableAction{
					Type: "delete", Label: l.Actions.Delete, Action: "delete",
					URL:      routes.DeleteURL,
					ItemName: acct.Name,
					Disabled: !canDelete, DisabledTooltip: l.Actions.NoPermission,
				},
			)
		}

		elementVariant := elementBadgeVariant(acct.Element)

		var balanceCell types.TableCell
		if acct.IsGroup {
			balanceCell = types.TableCell{Type: "text", Value: ""}
		} else {
			balanceCell = types.TableCell{Type: "text", Value: acct.Balance}
		}

		row := types.TableRow{
			ID: acct.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: acct.Code},
				{Type: "text", Value: acct.Name},
				{Type: "badge", Value: elementLabel(l, acct.Element), Variant: elementVariant},
				{Type: "text", Value: acct.Classification},
				balanceCell,
			},
			DataAttrs: map[string]string{
				"level":    levelStr(acct.Level),
				"is-group": boolStr(acct.IsGroup),
				"parent":   acct.ParentCode,
				"element":  acct.Element,
				"code":     acct.Code,
				"name":     acct.Name,
			},
			Actions: actions,
		}
		if !acct.IsGroup {
			row.Href = route.ResolveURL(routes.DetailURL, "id", acct.ID)
		}
		rows = append(rows, row)
	}
	return rows
}

// ---------------------------------------------------------------------------
// Proto enum → display string converters
// ---------------------------------------------------------------------------

func elementString(e accountpb.AccountElement) string {
	switch e {
	case accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET:
		return "asset"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY:
		return "liability"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY:
		return "equity"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE:
		return "revenue"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE:
		return "expense"
	default:
		return "asset"
	}
}

func classificationLabel(c accountpb.AccountClassification, l fycha.AccountFormLabels) string {
	switch c {
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET:
		return l.ClassCurrentAsset
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET:
		return l.ClassNonCurrentAsset
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY:
		return l.ClassCurrentLiability
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_LIABILITY:
		return l.ClassNonCurrentLiability
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_EQUITY:
		return l.ClassEquity
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE:
		return l.ClassOperatingRevenue
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OTHER_INCOME:
		return l.ClassOtherIncome
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_COST_OF_SALES:
		return l.ClassCostOfSales
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE:
		return l.ClassOperatingExpense
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_FINANCE_COST:
		return l.ClassFinanceCost
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_INCOME_TAX:
		return l.ClassIncomeTax
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OTHER_EXPENSE:
		return l.ClassOtherExpense
	default:
		return ""
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func elementBadgeVariant(element string) string {
	switch element {
	case "asset":
		return "sage"
	case "liability":
		return "amber"
	case "equity":
		return "navy"
	case "revenue":
		return "terracotta"
	case "expense":
		return "default"
	default:
		return "default"
	}
}

func elementLabel(l fycha.AccountLabels, element string) string {
	switch element {
	case "asset":
		return l.Tabs.Asset
	case "liability":
		return l.Tabs.Liability
	case "equity":
		return l.Tabs.Equity
	case "revenue":
		return l.Tabs.Revenue
	case "expense":
		return l.Tabs.Expense
	default:
		return element
	}
}

func levelStr(level int) string {
	switch level {
	case 0:
		return "0"
	case 1:
		return "1"
	case 2:
		return "2"
	default:
		return "3"
	}
}

func boolStr(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

// ---------------------------------------------------------------------------
// BuildAccountTree
// ---------------------------------------------------------------------------

// BuildAccountTree takes a flat list of AccountRow values (as returned from
// the use case / protoToRow) and returns a tree-ordered flat list with
// Level metadata computed from the parent_id chain.
//
// The flat list from the use case has no guaranteed order and no group rows.
// BuildAccountTree:
//  1. Groups accounts by element to infer root groups (Level 0).
//  2. Sorts each element group by Code.
//  3. Computes Level by counting hops from the root along ParentCode chains.
//  4. Synthesises IsGroup=true rows for element headers (Level 0) and
//     classification headers (Level 1) — these are display-only and carry
//     no ID.
//  5. Returns a pre-order flat list: root first, then children immediately
//     after their parent, mirroring how the tree renders.
//
// When the DB is fully wired this function replaces the hardcoded mockAccounts()
// flow. Until then the view falls back to mockAccounts() via fetchAccounts().
func BuildAccountTree(accounts []AccountRow, l fycha.AccountFormLabels) []AccountRow {
	if len(accounts) == 0 {
		return accounts
	}

	// Build a map of code → account for quick parent lookup.
	byCode := make(map[string]*AccountRow, len(accounts))
	for i := range accounts {
		if accounts[i].Code != "" {
			byCode[accounts[i].Code] = &accounts[i]
		}
	}

	// Compute Level for each account by walking the parent chain.
	levelCache := make(map[string]int, len(accounts))
	var computeLevel func(code string, visited map[string]bool) int
	computeLevel = func(code string, visited map[string]bool) int {
		if l, ok := levelCache[code]; ok {
			return l
		}
		if visited[code] {
			// Cycle guard — shouldn't happen in well-formed data
			return 2
		}
		visited[code] = true
		acct, ok := byCode[code]
		if !ok || acct.ParentCode == "" {
			levelCache[code] = 2 // leaf at base depth if no parent found
			return 2
		}
		parentLevel := computeLevel(acct.ParentCode, visited)
		levelCache[code] = parentLevel + 1
		return levelCache[code]
	}

	for i := range accounts {
		if accounts[i].Code != "" {
			accounts[i].Level = computeLevel(accounts[i].Code, make(map[string]bool))
		}
	}

	// Group by element for ordering: Assets → Liabilities → Equity → Revenue → Expenses
	elementOrder := []string{"asset", "liability", "equity", "revenue", "expense"}

	type elementGroup struct {
		element  string
		label    string
		codeBase string
		accounts []AccountRow
	}
	elementGroups := map[string]*elementGroup{
		"asset":     {element: "asset", label: l.GroupAssets, codeBase: "1000"},
		"liability": {element: "liability", label: l.GroupLiabilities, codeBase: "2000"},
		"equity":    {element: "equity", label: l.GroupEquity, codeBase: "3000"},
		"revenue":   {element: "revenue", label: l.GroupRevenue, codeBase: "4000"},
		"expense":   {element: "expense", label: l.GroupExpenses, codeBase: "5000"},
	}

	for _, a := range accounts {
		eg, ok := elementGroups[a.Element]
		if !ok {
			continue
		}
		eg.accounts = append(eg.accounts, a)
	}

	var result []AccountRow

	for _, elemKey := range elementOrder {
		eg, ok := elementGroups[elemKey]
		if !ok || len(eg.accounts) == 0 {
			continue
		}

		// Emit element group header (Level 0, IsGroup=true, no ID)
		result = append(result, AccountRow{
			ID:      "",
			Code:    eg.codeBase,
			Name:    eg.label,
			Element: eg.element,
			IsGroup: true,
			Level:   0,
			Active:  true,
		})

		// Group leaf accounts by Classification to emit classification headers (Level 1)
		// Maintain insertion order using a slice of seen classifications.
		seenClasses := make([]string, 0)
		classAccounts := make(map[string][]AccountRow)

		for _, a := range eg.accounts {
			cls := a.Classification
			if _, exists := classAccounts[cls]; !exists {
				seenClasses = append(seenClasses, cls)
			}
			classAccounts[cls] = append(classAccounts[cls], a)
		}

		for _, cls := range seenClasses {
			if cls != "" {
				// Emit classification header (Level 1, IsGroup=true)
				result = append(result, AccountRow{
					ID:             "",
					Code:           "",
					Name:           cls,
					Element:        eg.element,
					Classification: cls,
					IsGroup:        true,
					Level:          1,
					Active:         true,
				})
			}

			// Emit leaf accounts sorted by code (Level 2)
			leafAccounts := classAccounts[cls]
			// Simple insertion sort by Code (small N in practice)
			for i := 1; i < len(leafAccounts); i++ {
				for j := i; j > 0 && leafAccounts[j].Code < leafAccounts[j-1].Code; j-- {
					leafAccounts[j], leafAccounts[j-1] = leafAccounts[j-1], leafAccounts[j]
				}
			}
			for _, a := range leafAccounts {
				a.Level = 2
				result = append(result, a)
			}
		}
	}

	return result
}
