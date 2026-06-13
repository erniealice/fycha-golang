package prepayment

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "expenditure.prepayment",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "prepayment"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "prepayment.json", Key: ""},
		LabelName: "PrepaymentLabels",
		Templates: TemplatesFS,
		// Prepayments section removed from service-admin sidebar 2026-05-17 (Plan B Decision 2).
		// Advance disbursements are now accessible via Cash → Advances.
		// Nav items retained here for composability in other consumer apps.
		Nav: compose.NavContrib{
			Permission: "expense:list",
			Items: []compose.NavItem{
				{Key: "prepayments-active", Route: "prepayment.list", Params: map[string]string{"status": "active"}, Label: "Active", Icon: "icon-check-circle", Permission: "expense:list"},
				{Key: "prepayments-amortization", Route: "prepayment.amortization", Label: "Amortization Schedule", Icon: "icon-calendar", Permission: "expense:list"},
			},
		},
	}
}
