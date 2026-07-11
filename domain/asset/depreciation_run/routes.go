package depreciationrun

// Depreciation run history page routes (Surface D)

const (
	ListURL            = "/asset/depreciation-runs/list/{status}"
	ListTableURL       = "/action/depreciation-run/table/{status}"
	DetailURL          = "/asset/depreciation-runs/detail/{run_id}"
	DetailTabActionURL = "/action/depreciation-run/detail/{run_id}/tab/{tab}"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for the depreciation-run history
// list and detail pages (Surface D).
type Routes struct {
	// ActiveNav is the sidebar key used to highlight the active nav item.
	ActiveNav string `json:"active_nav"`

	// ListURL is the Surface D list page; status is a path parameter
	// (complete | failed | pending, or empty for all).
	ListURL string `json:"list_url"`

	// ListTableURL is the HTMX inner-swap target for the list table.
	ListTableURL string `json:"list_table_url"`

	// DetailURL is the Surface D detail page; run_id is a path parameter.
	DetailURL string `json:"detail_url"`

	// DetailTabActionURL is the HTMX tab-swap target on the detail page.
	DetailTabActionURL string `json:"detail_tab_action_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:          "asset",
		ListURL:            ListURL,
		ListTableURL:       ListTableURL,
		DetailURL:          DetailURL,
		DetailTabActionURL: DetailTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"depreciation_run.list":              r.ListURL,
		"depreciation_run.list_table":        r.ListTableURL,
		"depreciation_run.detail":            r.DetailURL,
		"depreciation_run.detail_tab_action": r.DetailTabActionURL,
	}
}

// ListFor returns the resolved depreciation-run list URL for the given status
// (e.g. "complete", "failed", "pending", or empty string for all).
func (r Routes) ListFor(status string) string {
	return resolveParam(r.ListURL, "status", status)
}

// DetailFor returns the resolved depreciation-run detail URL for the given
// run ID.
func (r Routes) DetailFor(runID string) string {
	return resolveParam(r.DetailURL, "run_id", runID)
}

// resolveParam replaces a single {placeholder} in a URL pattern with value.
func resolveParam(pattern, placeholder, value string) string {
	token := "{" + placeholder + "}"
	n := len(token)
	for i := 0; i+n <= len(pattern); i++ {
		if pattern[i:i+n] == token {
			return pattern[:i] + value + pattern[i+n:]
		}
	}
	return pattern
}
