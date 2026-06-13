package tax_rate

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "tax.tax_rate",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "tax_rate"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "tax_rate.json", Key: ""},
		LabelName: "TaxRateLabels",
		Templates: TemplatesFS,
	}
}
