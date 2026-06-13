package forex_rate

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "finance.forex_rate",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "forex_rate"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "forex_rate.json", Key: ""},
		LabelName: "ForexRateLabels",
		Templates: TemplatesFS,
	}
}
