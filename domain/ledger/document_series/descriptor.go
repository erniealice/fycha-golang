package document_series

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the document series settings
// pages. The unit is mounted by fycha's DocumentSeriesUnit only when the app
// opts in (block.WithDocumentSeries); it is not part of the default unit set.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "ledger.document_series",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "document_series"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "document_series.json", Key: "document_series"},
		LabelName: "DocumentSeriesLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "document_series:list",
			Items: []compose.NavItem{
				{Key: "document-series", Route: "document_series.list", Params: map[string]string{"status": "active"},
					Label: "Document Series", Icon: "icon-hash", Permission: "document_series:list",
					LabelKey: "document_series_label", IconKey: "document_series_icon"},
			},
		},
	}
}
