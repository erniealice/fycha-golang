package lapsing_schedule

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "asset.lapsing_schedule",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "lapsing_schedule"},
		Templates: TemplatesFS,
	}
}
