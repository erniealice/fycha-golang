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
		Nav: compose.NavContrib{
			Permission: "asset:list",
			AppEntry: &compose.AppEntry{
				Key:        "asset",
				Route:      "asset.dashboard",
				Label:      "Assets",
				Icon:       "icon-archive",
				Permission: "asset:list",
			},
			Items: []compose.NavItem{
				{Key: "assets-fixed", Route: "asset.list", Params: map[string]string{"status": "active"}, Label: "Fixed Assets", Icon: "icon-archive", Permission: "asset:list"},
				{Key: "depreciation-policies", Route: "asset.depreciation_policies", Label: "Depreciation Policies", Icon: "icon-settings", Permission: "asset:list"},
				{Key: "lapsing-schedule", Route: "asset.lapsing_schedule", Label: "Lapsing Schedule", Icon: "icon-calendar", Permission: "asset:list"},
				{Key: "depreciation-runs", Route: "depreciation_run.list", Params: map[string]string{"status": "complete"}, Label: "Depreciation Runs", Icon: "icon-trending-down", Permission: "asset:list"},
			},
		},
	}
}
