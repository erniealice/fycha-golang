package asset

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "asset.asset",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "asset"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "asset.json", Key: ""},
		LabelName: "AssetLabels",
		Templates: TemplatesFS,
	}
}
