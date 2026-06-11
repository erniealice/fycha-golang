package lapsing_schedule

// Lapsing schedule page routes (Surface B)

const (
	ListURL               = "/asset/lapsing-schedule/list"
	BulkRunSelectedURL    = "/action/lapsing-schedule/bulk-run-selected"
	BulkRunAllMatchingURL = "/action/lapsing-schedule/bulk-run-all-matching"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for the lapsing-schedule live page
// (Surface B) and its bulk-action endpoints.
type Routes struct {
	// ListURL is the Surface B full-page lapsing-schedule list.
	ListURL string `json:"list_url"`

	// BulkRunSelectedURL is the endpoint for bulk-running selected rows.
	BulkRunSelectedURL string `json:"bulk_run_selected_url"`

	// BulkRunAllMatchingURL is the endpoint for bulk-running all matching rows.
	BulkRunAllMatchingURL string `json:"bulk_run_all_matching_url"`

	// AssetListBulkRunSelectedURL is the equivalent bulk-run endpoint on the
	// main assets list page.
	AssetListBulkRunSelectedURL string `json:"asset_list_bulk_run_selected_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ListURL:               ListURL,
		BulkRunSelectedURL:    BulkRunSelectedURL,
		BulkRunAllMatchingURL: BulkRunAllMatchingURL,
		// AssetListBulkRunSelectedURL comes from the asset entity package const.
		AssetListBulkRunSelectedURL: "/action/asset/bulk-run-selected",
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"lapsing_schedule.list":                         r.ListURL,
		"lapsing_schedule.bulk_run_selected":            r.BulkRunSelectedURL,
		"lapsing_schedule.bulk_run_all_matching":        r.BulkRunAllMatchingURL,
		"lapsing_schedule.asset_list_bulk_run_selected": r.AssetListBulkRunSelectedURL,
	}
}
