package payrolldashboard

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "payroll.payrolldashboard",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payroll"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payroll.json", Key: "dashboard"},
		LabelName: "PayrollDashboardLabels",
		Templates: TemplatesFS,
	}
}
