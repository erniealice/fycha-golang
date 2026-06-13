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
	}
}
