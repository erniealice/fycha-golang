package equity

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "ledger.equity",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "equity"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "equity.json", Key: ""},
		LabelName: "EquityLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "equity:list",
			AppEntry: &compose.AppEntry{
				Key:        "equity",
				Route:      "equity.accounts",
				Label:      "Equity",
				Icon:       "icon-trending-up",
				Permission: "equity:list",
			},
			Items: []compose.NavItem{
				{Key: "capital-accounts", Route: "equity.accounts", Label: "Capital Accounts", Icon: "icon-users", Permission: "equity:list"},
				{Key: "equity-transactions", Route: "equity.transactions", Label: "Equity Transactions", Icon: "icon-repeat", Permission: "equity:list"},
			},
		},
	}
}
