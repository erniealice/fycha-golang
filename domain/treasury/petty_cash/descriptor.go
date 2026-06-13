package petty_cash

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.petty_cash",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "petty_cash"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "petty_cash.json", Key: ""},
		LabelName: "PettyCashLabels",
		Templates: TemplatesFS,
	}
}
