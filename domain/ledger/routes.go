// routes.go defines configurable route structs for fycha ledger-domain views.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultXxxRoutes()
//     constructors return structs populated from the package-level route constants
//     defined in routes.go. These serve as sensible defaults for any consumer app.
//   - Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//     Apps can load a JSON config file that maps route keys to custom paths,
//     allowing industry templates (e.g. salon, retail) to rebrand URLs without
//     code changes. The json struct tags on each field support this workflow.
//   - Level 3: App-specific overrides via Go field assignment (optional).
//     After constructing defaults (and optionally applying JSON), consumer apps
//     can directly assign individual struct fields for one-off customizations.
//
// RouteMap() methods return a map[string]string of dot-notation keys to route
// paths, useful for template rendering and route resolution at runtime.
package ledger

const (
	LedgerBaseURL       = "/ledger/"
	LedgerDashboardURL  = "/ledger/dashboard"
	AccountListURL      = "/ledger/accounts/list"
	AccountDetailURL    = "/ledger/accounts/detail/{id}"
	AccountTabActionURL = "/action/ledger/account/{id}/tab/{tab}"
	AccountTreeURL      = "/ledger/accounts/tree"
	AccountTemplatesURL = "/ledger/settings/account-templates"
	AccountAddURL       = "/action/ledger/account/add"
	AccountEditURL      = "/action/ledger/account/edit/{id}"
	AccountDeleteURL    = "/action/ledger/account/delete"

	// Ledger — Journal Entries
	JournalListURL             = "/ledger/journals/list/{status}"
	JournalDetailURL           = "/ledger/journals/detail/{id}"
	JournalTabActionURL        = "/action/ledger/journal/{id}/tab/{tab}"
	JournalAttachmentUploadURL = "/action/ledger/journal/{id}/attachments/upload"
	JournalAttachmentDeleteURL = "/action/ledger/journal/{id}/attachments/delete"
	JournalAddURL              = "/action/ledger/journal/add"
	JournalEditURL             = "/action/ledger/journal/edit/{id}"
	JournalPostURL             = "/action/ledger/journal/post/{id}"
	JournalReverseURL          = "/action/ledger/journal/reverse/{id}"
	JournalDeleteURL           = "/action/ledger/journal/delete"

	// Ledger — Accounting Statements (internal tools)
	LedgerGeneralLedgerURL = "/ledger/reports/general-ledger"
	LedgerTrialBalanceURL  = "/ledger/reports/trial-balance"

	// Ledger — Fiscal Periods / Settings
	FiscalPeriodListURL   = "/ledger/settings/fiscal-periods"
	FiscalPeriodDetailURL = "/ledger/settings/fiscal-periods/detail/{id}"
	FiscalPeriodAddURL    = "/action/ledger/fiscal-period/add"
	FiscalPeriodCloseURL  = "/action/ledger/fiscal-period/close/{id}"

	// Ledger — Bad Debt Policy
	BadDebtPolicyURL = "/ledger/settings/bad-debt-policy"

	// Ledger — Recurring Templates
	RecurringTemplatesURL = "/ledger/settings/recurring"

	// Funding — Equity
	EquityDashboardURL      = "/funding/equity/dashboard"
	EquityAccountsURL       = "/funding/equity/accounts"
	EquityAccountDetailURL  = "/funding/equity/accounts/detail/{id}"
	EquityTransactionsURL   = "/funding/equity/transactions"
	EquityTransactionAddURL = "/action/funding/equity/transaction/add"
)

// ---------------------------------------------------------------------------
// AccountRoutes
// ---------------------------------------------------------------------------

// AccountRoutes holds route paths for Chart of Accounts views.
type AccountRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	DashboardURL string `json:"dashboard_url"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	TreeURL      string `json:"tree_url"`
	TemplatesURL string `json:"templates_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	DeleteURL    string `json:"delete_url"`
}

func DefaultAccountRoutes() AccountRoutes {
	return AccountRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "chart-of-accounts",
		DashboardURL: LedgerDashboardURL,
		ListURL:      AccountListURL,
		DetailURL:    AccountDetailURL,
		TabActionURL: AccountTabActionURL,
		TreeURL:      AccountTreeURL,
		TemplatesURL: AccountTemplatesURL,
		AddURL:       AccountAddURL,
		EditURL:      AccountEditURL,
		DeleteURL:    AccountDeleteURL,
	}
}

func (r AccountRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.dashboard":         r.DashboardURL,
		"ledger.account.list":      r.ListURL,
		"ledger.account.detail":    r.DetailURL,
		"ledger.account.tree":      r.TreeURL,
		"ledger.account.templates": r.TemplatesURL,
		"ledger.account.add":       r.AddURL,
		"ledger.account.edit":      r.EditURL,
		"ledger.account.delete":    r.DeleteURL,
	}
}

// ---------------------------------------------------------------------------
// JournalRoutes
// ---------------------------------------------------------------------------

// JournalRoutes holds route paths for Journal Entry views.
type JournalRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	PostURL      string `json:"post_url"`
	ReverseURL   string `json:"reverse_url"`
	DeleteURL    string `json:"delete_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

func DefaultJournalRoutes() JournalRoutes {
	return JournalRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "journals-draft",
		ListURL:      JournalListURL,
		DetailURL:    JournalDetailURL,
		TabActionURL: JournalTabActionURL,
		AddURL:       JournalAddURL,
		EditURL:      JournalEditURL,
		PostURL:      JournalPostURL,
		ReverseURL:   JournalReverseURL,
		DeleteURL:    JournalDeleteURL,

		AttachmentUploadURL: JournalAttachmentUploadURL,
		AttachmentDeleteURL: JournalAttachmentDeleteURL,
	}
}

func (r JournalRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.journal.list":    r.ListURL,
		"ledger.journal.detail":  r.DetailURL,
		"ledger.journal.add":     r.AddURL,
		"ledger.journal.edit":    r.EditURL,
		"ledger.journal.post":    r.PostURL,
		"ledger.journal.reverse": r.ReverseURL,
		"ledger.journal.delete":  r.DeleteURL,

		"ledger.journal.attachment.upload": r.AttachmentUploadURL,
		"ledger.journal.attachment.delete": r.AttachmentDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// LedgerStatementRoutes
// ---------------------------------------------------------------------------

// LedgerStatementRoutes holds route paths for accounting statement views
// (General Ledger, Trial Balance — internal accounting tools, not business reports).
type LedgerStatementRoutes struct {
	ActiveNav        string `json:"active_nav"`
	GeneralLedgerURL string `json:"general_ledger_url"`
	TrialBalanceURL  string `json:"trial_balance_url"`
}

func DefaultLedgerStatementRoutes() LedgerStatementRoutes {
	return LedgerStatementRoutes{
		ActiveNav:        "ledger",
		GeneralLedgerURL: LedgerGeneralLedgerURL,
		TrialBalanceURL:  LedgerTrialBalanceURL,
	}
}

func (r LedgerStatementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.statement.general_ledger": r.GeneralLedgerURL,
		"ledger.statement.trial_balance":  r.TrialBalanceURL,
	}
}

// ---------------------------------------------------------------------------
// FiscalPeriodRoutes
// ---------------------------------------------------------------------------

// FiscalPeriodRoutes holds route paths for fiscal period management views.
type FiscalPeriodRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	AddURL       string `json:"add_url"`
	CloseURL     string `json:"close_url"`
}

func DefaultFiscalPeriodRoutes() FiscalPeriodRoutes {
	return FiscalPeriodRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "fiscal-periods",
		ListURL:      FiscalPeriodListURL,
		DetailURL:    FiscalPeriodDetailURL,
		AddURL:       FiscalPeriodAddURL,
		CloseURL:     FiscalPeriodCloseURL,
	}
}

func (r FiscalPeriodRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.fiscal_period.list":   r.ListURL,
		"ledger.fiscal_period.detail": r.DetailURL,
		"ledger.fiscal_period.add":    r.AddURL,
		"ledger.fiscal_period.close":  r.CloseURL,
	}
}

// ---------------------------------------------------------------------------
// LedgerSettingsRoutes
// ---------------------------------------------------------------------------

// LedgerSettingsRoutes holds route paths for ledger settings views
// (Bad Debt Policy, Recurring Templates).
type LedgerSettingsRoutes struct {
	ActiveNav             string `json:"active_nav"`
	BadDebtPolicyURL      string `json:"bad_debt_policy_url"`
	RecurringTemplatesURL string `json:"recurring_templates_url"`
}

func DefaultLedgerSettingsRoutes() LedgerSettingsRoutes {
	return LedgerSettingsRoutes{
		ActiveNav:             "ledger",
		BadDebtPolicyURL:      BadDebtPolicyURL,
		RecurringTemplatesURL: RecurringTemplatesURL,
	}
}

func (r LedgerSettingsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.settings.bad_debt_policy":     r.BadDebtPolicyURL,
		"ledger.settings.recurring_templates": r.RecurringTemplatesURL,
	}
}

// ---------------------------------------------------------------------------
// EquityRoutes
// ---------------------------------------------------------------------------

// EquityRoutes holds route paths for Equity views.
type EquityRoutes struct {
	ActiveNav         string `json:"active_nav"`
	DashboardURL      string `json:"dashboard_url"`
	AccountsURL       string `json:"accounts_url"`
	AccountDetailURL  string `json:"account_detail_url"`
	TransactionsURL   string `json:"transactions_url"`
	TransactionAddURL string `json:"transaction_add_url"`
}

func DefaultEquityRoutes() EquityRoutes {
	return EquityRoutes{
		ActiveNav:         "equity",
		DashboardURL:      EquityDashboardURL,
		AccountsURL:       EquityAccountsURL,
		AccountDetailURL:  EquityAccountDetailURL,
		TransactionsURL:   EquityTransactionsURL,
		TransactionAddURL: EquityTransactionAddURL,
	}
}

func (r EquityRoutes) RouteMap() map[string]string {
	return map[string]string{
		"equity.dashboard":       r.DashboardURL,
		"equity.accounts":        r.AccountsURL,
		"equity.account_detail":  r.AccountDetailURL,
		"equity.transactions":    r.TransactionsURL,
		"equity.transaction_add": r.TransactionAddURL,
	}
}

// ---------------------------------------------------------------------------
// resolveParam — internal URL template helper
// ---------------------------------------------------------------------------

// resolveParam replaces a single {placeholder} in a URL pattern with value.
// It is the internal single-parameter URL resolver; for multi-parameter URLs
// use route.ResolveURL from packages/pyeza-golang/route directly.
func resolveParam(pattern, placeholder, value string) string {
	token := "{" + placeholder + "}"
	n := len(token)
	for i := 0; i+n <= len(pattern); i++ {
		if pattern[i:i+n] == token {
			return pattern[:i] + value + pattern[i+n:]
		}
	}
	return pattern
}
