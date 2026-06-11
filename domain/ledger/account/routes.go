// routes.go defines configurable route structs for fycha ledger account views.
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
package account

const (
	LedgerBaseURL      = "/ledger/"
	LedgerDashboardURL = "/ledger/dashboard"
	ListURL            = "/ledger/accounts/list"
	DetailURL          = "/ledger/accounts/detail/{id}"
	TabActionURL       = "/action/ledger/account/{id}/tab/{tab}"
	TreeURL            = "/ledger/accounts/tree"
	TemplatesURL       = "/ledger/settings/account-templates"
	AddURL             = "/action/ledger/account/add"
	EditURL            = "/action/ledger/account/edit/{id}"
	DeleteURL          = "/action/ledger/account/delete"

	// Ledger — Accounting Statements (internal tools)
	LedgerGeneralLedgerURL = "/ledger/reports/general-ledger"
	LedgerTrialBalanceURL  = "/ledger/reports/trial-balance"

	// Ledger — Bad Debt Policy
	BadDebtPolicyURL = "/ledger/settings/bad-debt-policy"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Chart of Accounts views.
type Routes struct {
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

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "ledger",
		ActiveSubNav: "chart-of-accounts",
		DashboardURL: LedgerDashboardURL,
		ListURL:      ListURL,
		DetailURL:    DetailURL,
		TabActionURL: TabActionURL,
		TreeURL:      TreeURL,
		TemplatesURL: TemplatesURL,
		AddURL:       AddURL,
		EditURL:      EditURL,
		DeleteURL:    DeleteURL,
	}
}

func (r Routes) RouteMap() map[string]string {
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
// StatementRoutes
// ---------------------------------------------------------------------------

// StatementRoutes holds route paths for accounting statement views
// (General Ledger, Trial Balance — internal accounting tools, not business reports).
type StatementRoutes struct {
	ActiveNav        string `json:"active_nav"`
	GeneralLedgerURL string `json:"general_ledger_url"`
	TrialBalanceURL  string `json:"trial_balance_url"`
}

func DefaultStatementRoutes() StatementRoutes {
	return StatementRoutes{
		ActiveNav:        "ledger",
		GeneralLedgerURL: LedgerGeneralLedgerURL,
		TrialBalanceURL:  LedgerTrialBalanceURL,
	}
}

func (r StatementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.statement.general_ledger": r.GeneralLedgerURL,
		"ledger.statement.trial_balance":  r.TrialBalanceURL,
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
