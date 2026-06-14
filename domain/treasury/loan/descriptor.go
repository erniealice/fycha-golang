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
				// NOTE: "payments-upcoming" / "payments-history" nav items removed —
				// they reference route key "loan_payment.list", but loan_payment is
				// scaffold-only (labels.go + routes.go, no descriptor / view module /
				// mounted unit in fycha block.AllUnits), so the key is never in the
				// route table and a dangling ref fail-closes the entire fycha engine.
				// Restore once loan_payment ships a view module and is mounted.
				{Key: "amortization-schedules", Route: "loan.amortization", Label: "Amortization Schedules", Icon: "icon-calendar", Permission: "loan:list"},
			},
		},
	}
}
