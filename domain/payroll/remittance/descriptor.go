package remittance

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "payroll.remittance",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payroll_remittance"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payroll.json", Key: ""},
		LabelName: "RemittanceLabels",
		Templates: TemplatesFS,
	}
}
