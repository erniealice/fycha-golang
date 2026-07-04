package depreciationrun

import "github.com/erniealice/espyna-golang/consumer/compose"

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
		Nav: compose.NavContrib{
			Permission: "asset:list",
			Items: []compose.NavItem{
				{Key: "depreciation-runs", Route: "depreciation_run.list", Params: map[string]string{"status": "complete"}, Label: "Depreciation Runs", Icon: "icon-trending-down", Permission: "asset:list", LabelKey: "depreciation_runs_label", IconKey: "depreciation_runs_icon"},
			},
		},
	}
}
