// routes.go defines configurable route structs for the petty_cash entity.
package petty_cash

const (
	// Cash — Petty Cash
	RegisterURL          = "/cash/petty-cash/register"
	ReplenishmentListURL = "/cash/petty-cash/replenishments/{status}"
	CustodianBalancesURL = "/cash/petty-cash/custodian-balances"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Cash app petty cash views.
type Routes struct {
	ActiveNav            string `json:"active_nav"`
	RegisterURL          string `json:"register_url"`
	ReplenishmentListURL string `json:"replenishment_list_url"`
	CustodianBalancesURL string `json:"custodian_balances_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:            "cash",
		RegisterURL:          RegisterURL,
		ReplenishmentListURL: ReplenishmentListURL,
		CustodianBalancesURL: CustodianBalancesURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"petty_cash.register":           r.RegisterURL,
		"petty_cash.replenishment_list": r.ReplenishmentListURL,
		"petty_cash.custodian_balances": r.CustodianBalancesURL,
	}
}
