package depreciationrun

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "asset.depreciation_run",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "depreciation_run"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "depreciation_run.json", Key: ""},
		LabelName: "DepreciationRunLabels",
		Templates: TemplatesFS,
	}
}
