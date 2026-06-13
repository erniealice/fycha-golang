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
	}
}
