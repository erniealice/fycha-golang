package petty_cash

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.petty_cash",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "petty_cash"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "petty_cash.json", Key: ""},
		LabelName: "PettyCashLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "petty_cash:list",
			Items: []compose.NavItem{
				{Key: "petty-cash-register", Route: "petty_cash.register", Label: "Register", Icon: "icon-clipboard", Permission: "petty_cash:list"},
				{Key: "petty-cash-replenishments", Route: "petty_cash.replenishment_list", Params: map[string]string{"status": "all"}, Label: "Replenishments", Icon: "icon-refresh-cw", Permission: "petty_cash:list"},
				{Key: "petty-cash-custodians", Route: "petty_cash.custodian_balances", Label: "Custodian Balances", Icon: "icon-users", Permission: "petty_cash:list"},
			},
		},
	}
}
