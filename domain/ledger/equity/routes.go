// routes.go defines configurable route structs for fycha ledger equity views.
package equity

const (
	// Funding — Equity
	DashboardURL      = "/funding/equity/dashboard"
	AccountsURL       = "/funding/equity/accounts"
	AccountDetailURL  = "/funding/equity/accounts/detail/{id}"
	TransactionsURL   = "/funding/equity/transactions"
	TransactionAddURL = "/action/funding/equity/transaction/add"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Equity views.
type Routes struct {
	ActiveNav         string `json:"active_nav"`
	DashboardURL      string `json:"dashboard_url"`
	AccountsURL       string `json:"accounts_url"`
	AccountDetailURL  string `json:"account_detail_url"`
	TransactionsURL   string `json:"transactions_url"`
	TransactionAddURL string `json:"transaction_add_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:         "equity",
		DashboardURL:      DashboardURL,
		AccountsURL:       AccountsURL,
		AccountDetailURL:  AccountDetailURL,
		TransactionsURL:   TransactionsURL,
		TransactionAddURL: TransactionAddURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"equity.dashboard":       r.DashboardURL,
		"equity.accounts":        r.AccountsURL,
		"equity.account_detail":  r.AccountDetailURL,
		"equity.transactions":    r.TransactionsURL,
		"equity.transaction_add": r.TransactionAddURL,
	}
}
