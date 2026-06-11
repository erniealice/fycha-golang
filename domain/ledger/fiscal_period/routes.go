// routes.go defines configurable route structs for fycha ledger fiscal_period views.
package fiscal_period

const (
	// Ledger — Fiscal Periods / Settings
	ListURL   = "/ledger/settings/fiscal-periods"
	DetailURL = "/ledger/settings/fiscal-periods/detail/{id}"
	AddURL    = "/action/ledger/fiscal-period/add"
	CloseURL  = "/action/ledger/fiscal-period/close/{id}"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for fiscal period management views.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	AddURL       string `json:"add_url"`
	CloseURL     string `json:"close_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "ledger",
		ActiveSubNav: "fiscal-periods",
		ListURL:      ListURL,
		DetailURL:    DetailURL,
		AddURL:       AddURL,
		CloseURL:     CloseURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.fiscal_period.list":   r.ListURL,
		"ledger.fiscal_period.detail": r.DetailURL,
		"ledger.fiscal_period.add":    r.AddURL,
		"ledger.fiscal_period.close":  r.CloseURL,
	}
}
