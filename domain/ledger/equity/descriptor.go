package equity

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "ledger.equity",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "equity"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "equity.json", Key: ""},
		LabelName: "EquityLabels",
		Templates: TemplatesFS,
	}
}
