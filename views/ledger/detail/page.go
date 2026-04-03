package detail

import (
	"context"
	"fmt"
	"log"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	"github.com/erniealice/hybra-golang/views/auditlog"
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

// Deps holds view dependencies.
type Deps struct {
	auditlog.AuditOps

	Routes       fycha.AccountRoutes
	Labels       fycha.AccountLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Account use cases
	ReadAccount func(ctx context.Context, req *accountpb.ReadAccountRequest) (*accountpb.ReadAccountResponse, error)
}

// PageData holds the data for the account detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          fycha.AccountLabels
	ActiveTab       string
	TabItems        []pyeza.TabItem
	// Account identity
	ID             string
	AccountCode    string
	AccountName    string
	Element        string
	ElementVariant string
	Classification string
	Group          string
	ParentAccount  string
	NormalBalance  string
	CashFlowTag    string
	Description    string
	AccountStatus  string
	StatusVariant  string
	Created        string
	LastModified   string
	// Stats
	CurrentBalance string
	BalanceColor   string
	PeriodDebits   string
	PeriodCredits  string
	// Tables
	EntriesTable *types.TableConfig
	// Actions
	EditURL string
	CanEdit bool
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the account detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "entries"
		}

		perms := view.GetUserPermissions(ctx)
		pageData := buildPageData(ctx, deps, id, activeTab, viewCtx, perms)

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "ledger-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("account-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial — returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "entries"
		}

		perms := view.GetUserPermissions(ctx)
		pageData := buildPageData(ctx, deps, id, tab, viewCtx, perms)

		templateName := "account-tab-" + tab
		if tab == "audit-history" {
			templateName = "audit-history-tab"
		}
		return view.OK(templateName, pageData)
	})
}

// ---------------------------------------------------------------------------
// Page data builder
// ---------------------------------------------------------------------------

func buildPageData(ctx context.Context, deps *Deps, id, activeTab string, viewCtx *view.ViewContext, perms *types.UserPermissions) *PageData {
	acct := fetchAccount(ctx, deps, id, deps.Labels)

	tabItems := buildTabItems(id, deps.Labels, deps.Routes)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          fmt.Sprintf("%s \u2013 %s", acct.AccountCode, acct.AccountName),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    fmt.Sprintf("%s \u2013 %s", acct.AccountCode, acct.AccountName),
			HeaderSubtitle: acct.Classification,
			HeaderIcon:     "icon-book-open",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "account-detail-content",
		Labels:          deps.Labels,
		ActiveTab:       activeTab,
		TabItems:        tabItems,
		ID:              id,
		AccountCode:     acct.AccountCode,
		AccountName:     acct.AccountName,
		Element:         acct.Element,
		ElementVariant:  acct.ElementVariant,
		Classification:  acct.Classification,
		Group:           acct.Group,
		ParentAccount:   acct.ParentAccount,
		NormalBalance:   acct.NormalBalance,
		CashFlowTag:     acct.CashFlowTag,
		Description:     acct.Description,
		AccountStatus:   acct.AccountStatus,
		StatusVariant:   acct.StatusVariant,
		Created:         acct.Created,
		LastModified:    acct.LastModified,
		CurrentBalance:  acct.CurrentBalance,
		BalanceColor:    acct.BalanceColor,
		PeriodDebits:    acct.PeriodDebits,
		PeriodCredits:   acct.PeriodCredits,
		EditURL:         route.ResolveURL(deps.Routes.EditURL, "id", id),
		CanEdit:         perms.Can("account", "update"),
	}

	if activeTab == "entries" {
		pageData.EntriesTable = buildEntriesTable(nil, deps.Labels, deps.TableLabels, deps.Routes)
	}

	if activeTab == "audit-history" {
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			auditResp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
				EntityType:  "account",
				EntityID:    id,
				Limit:       20,
				CursorToken: cursor,
			})
			if err != nil {
				log.Printf("Failed to load audit history: %v", err)
			}
			if auditResp != nil {
				pageData.AuditEntries = auditResp.Entries
				pageData.AuditHasNext = auditResp.HasNext
				pageData.AuditNextCursor = auditResp.NextCursor
			}
		}
		pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
	}

	return pageData
}

func buildTabItems(id string, labels fycha.AccountLabels, routes fycha.AccountRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "entries", Label: labels.Detail.Tabs.JournalEntries, Href: base + "?tab=entries", HxGet: action + "entries", Icon: "icon-file-text", Count: 0, Disabled: false},
		{Key: "details", Label: labels.Detail.Tabs.Details, Href: base + "?tab=details", HxGet: action + "details", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "audit-history", Label: "History", Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

// accountViewModel is the normalised view-model for a single account.
type accountViewModel struct {
	AccountCode    string
	AccountName    string
	Element        string
	ElementVariant string
	Classification string
	Group          string
	ParentAccount  string
	NormalBalance  string
	CashFlowTag    string
	Description    string
	AccountStatus  string
	StatusVariant  string
	Created        string
	LastModified   string
	CurrentBalance string
	BalanceColor   string
	PeriodDebits   string
	PeriodCredits  string
}

// fetchAccount loads a single account by ID via ReadAccount use case.
// Returns a placeholder view-model on error so the page renders with an empty state.
func fetchAccount(ctx context.Context, deps *Deps, id string, l fycha.AccountLabels) accountViewModel {
	placeholder := accountViewModel{
		AccountCode:    "\u2014",
		AccountName:    "Account not found",
		Element:        "asset",
		ElementVariant: "default",
		Classification: "\u2014",
		Group:          "\u2014",
		ParentAccount:  "\u2014",
		NormalBalance:  "\u2014",
		CashFlowTag:    "\u2014",
		Description:    "\u2014",
		AccountStatus:  "unknown",
		StatusVariant:  "default",
		Created:        "\u2014",
		LastModified:   "\u2014",
		CurrentBalance: "\u20b10.00",
		BalanceColor:   "default",
		PeriodDebits:   "\u20b10.00",
		PeriodCredits:  "\u20b10.00",
	}

	if deps.ReadAccount == nil {
		return placeholder
	}

	resp, err := deps.ReadAccount(ctx, &accountpb.ReadAccountRequest{
		Data: &accountpb.Account{Id: id},
	})
	if err != nil {
		log.Printf("ReadAccount error for %s: %v", id, err)
		return placeholder
	}
	if resp == nil || !resp.GetSuccess() || len(resp.GetData()) == 0 {
		return placeholder
	}

	return protoToViewModel(resp.GetData()[0], l)
}

// protoToViewModel converts a proto Account to accountViewModel.
func protoToViewModel(a *accountpb.Account, l fycha.AccountLabels) accountViewModel {
	element := elementString(a.GetElement())
	elementVariant := elementBadgeVariant(element)

	status := statusString(a.GetStatus())
	statusVariant := statusBadgeVariant(status)

	balanceColor := "sage"
	if a.GetNormalBalance() == accountpb.NormalBalance_NORMAL_BALANCE_CREDIT {
		balanceColor = "terracotta"
	}

	return accountViewModel{
		AccountCode:    a.GetCode(),
		AccountName:    a.GetName(),
		Element:        element,
		ElementVariant: elementVariant,
		Classification: classificationLabel(a.GetClassification(), l.Form),
		Group:          a.GetGroupId(),
		ParentAccount:  a.GetParentId(),
		NormalBalance:  normalBalanceLabel(a.GetNormalBalance(), l.Form),
		CashFlowTag:    cashFlowLabel(a.GetCashFlowActivity(), l.Form),
		Description:    a.GetDescription(),
		AccountStatus:  status,
		StatusVariant:  statusVariant,
		Created:        a.GetDateCreatedString(),
		LastModified:   a.GetDateModifiedString(),
		CurrentBalance: "\u20b10.00", // running balance not in proto yet
		BalanceColor:   balanceColor,
		PeriodDebits:   "\u20b10.00",
		PeriodCredits:  "\u20b10.00",
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

// EntryRow holds one row of the journal entries sub-table.
type EntryRow struct {
	Date        string
	EntryNumber string
	EntryID     string
	Description string
	Debit       string
	Credit      string
}

func buildEntriesTable(entries []EntryRow, labels fycha.AccountLabels, tableLabels types.TableLabels, routes fycha.AccountRoutes) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: labels.Columns.Date, Sortable: true, Width: "100px"},
		{Key: "entry_number", Label: labels.Columns.EntryNumber, Sortable: true, Width: "110px"},
		{Key: "description", Label: labels.Columns.Description, Sortable: false},
		{Key: "debit", Label: labels.Columns.Debit, Sortable: false, Width: "130px", Align: "right"},
		{Key: "credit", Label: labels.Columns.Credit, Sortable: false, Width: "130px", Align: "right"},
	}

	rows := make([]types.TableRow, len(entries))
	for i, e := range entries {
		entryHref := route.ResolveURL(routes.DetailURL, "id", e.EntryID)
		rows[i] = types.TableRow{
			ID: fmt.Sprintf("entry-%d", i+1),
			Cells: []types.TableCell{
				{Type: "text", Value: e.Date},
				{Type: "link", Value: e.EntryNumber, Href: entryHref},
				{Type: "text", Value: e.Description},
				{Type: "text", Value: e.Debit},
				{Type: "text", Value: e.Credit},
			},
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:         "account-entries-table",
		Minimal:    true,
		Columns:    columns,
		Rows:       rows,
		ShowSearch: true,
		Labels:     tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   labels.Detail.EmptyStates.EntriesTitle,
			Message: labels.Detail.EmptyStates.EntriesMessage,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

// ---------------------------------------------------------------------------
// Proto enum → display string converters (shared with list package via duplication)
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

func normalBalanceLabel(n accountpb.NormalBalance, l fycha.AccountFormLabels) string {
	switch n {
	case accountpb.NormalBalance_NORMAL_BALANCE_DEBIT:
		return l.NormalBalanceDebit
	case accountpb.NormalBalance_NORMAL_BALANCE_CREDIT:
		return l.NormalBalanceCredit
	default:
		return "\u2014"
	}
}

func cashFlowLabel(c accountpb.CashFlowActivity, l fycha.AccountFormLabels) string {
	switch c {
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_OPERATING:
		return l.CashFlowOperating
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_INVESTING:
		return l.CashFlowInvesting
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_FINANCING:
		return l.CashFlowFinancing
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_NONE:
		return l.CashFlowNone
	default:
		return "\u2014"
	}
}

func statusString(s accountpb.AccountStatus) string {
	switch s {
	case accountpb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return "active"
	case accountpb.AccountStatus_ACCOUNT_STATUS_INACTIVE:
		return "inactive"
	case accountpb.AccountStatus_ACCOUNT_STATUS_LOCKED:
		return "locked"
	default:
		return "unknown"
	}
}

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

func statusBadgeVariant(status string) string {
	switch status {
	case "active":
		return "success"
	case "inactive":
		return "warning"
	case "locked":
		return "danger"
	default:
		return "default"
	}
}
