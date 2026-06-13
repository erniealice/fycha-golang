package run

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "payroll.run",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payroll_run"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payroll.json", Key: ""},
		LabelName: "RunLabels",
		Templates: TemplatesFS,
	}
}
