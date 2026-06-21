package payrollsettings

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "payroll.payrollsettings",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payroll_settings"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payroll.json", Key: ""},
		LabelName: "PayrollSettingsLabels",
		Templates: TemplatesFS,
	}
}
