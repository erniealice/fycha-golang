package forex_rate

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "forex_rate:read",
			Items: []compose.NavItem{
				{Key: "forex-rates", Route: "forex_rate.list", Params: map[string]string{"status": "active"}, Label: "Forex Rates", Icon: "icon-refresh-cw", Permission: "forex_rate:read", LabelKey: "forex_rates_label", IconKey: "forex_rates_icon"},
			},
		},
	}
}
