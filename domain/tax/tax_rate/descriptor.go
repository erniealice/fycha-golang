package tax_rate

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "tax_rate:read",
			Items: []compose.NavItem{
				{Key: "tax-rates", Route: "tax_rate.list", Params: map[string]string{"status": "active"}, Label: "Tax Rates", Icon: "icon-percent", Permission: "tax_rate:read", LabelKey: "tax_rates_label", IconKey: "tax_rates_icon"},
			},
		},
	}
}
