package payrollemployee

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "payroll.payrollemployee",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payroll_employee"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payroll.json", Key: ""},
		LabelName: "PayrollEmployeeLabels",
		Templates: TemplatesFS,
	}
}
