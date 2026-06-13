package asset_category

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	return compose.Unit{
		Key:       "asset.asset_category",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "asset_category_depreciation"},
		Templates: TemplatesFS,
	}
}
