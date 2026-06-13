package loan

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.loan",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "loan"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "loan.json", Key: ""},
		LabelName: "LoanLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "loan:list",
			AppEntry: &compose.AppEntry{
				Key:        "loan",
				Route:      "loan.dashboard",
				Label:      "Loans",
				Icon:       "icon-credit-card",
				Permission: "loan:list",
			},
			Items: []compose.NavItem{
				{Key: "loans-active", Route: "loan.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-check-circle", Permission: "loan:list"},
				{Key: "loans-completed", Route: "loan.list", Params: map[string]string{"status": "completed"}, Label: "Complete", Icon: "icon-check-circle", Permission: "loan:list"},
				{Key: "payments-upcoming", Route: "loan_payment.list", Params: map[string]string{"status": "upcoming"}, Label: "Upcoming", Icon: "icon-clock", Permission: "loan:list"},
				{Key: "payments-history", Route: "loan_payment.list", Params: map[string]string{"status": "history"}, Label: "History", Icon: "icon-clock", Permission: "loan:list"},
				{Key: "amortization-schedules", Route: "loan.amortization", Label: "Amortization Schedules", Icon: "icon-calendar", Permission: "loan:list"},
			},
		},
	}
}
