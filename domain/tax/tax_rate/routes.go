package tax_rate

const (
	// Tax — Tax Rates (read-only; supersession via admin SQL recipe)
	ListURL   = "/tax/tax-rates/list/{status}"
	DetailURL = "/tax/tax-rates/detail/{id}"
)

// Routes holds route paths for Tax Rate read-only views.
// Tax rates are read-only in the UI; supersession is done via admin SQL recipe.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
}

// DefaultRoutes returns a Routes populated from the package-level
// route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "ledger",
		ActiveSubNav: "tax-rates",
		ListURL:      ListURL,
		DetailURL:    DetailURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"tax_rate.list":   r.ListURL,
		"tax_rate.detail": r.DetailURL,
	}
}

// DetailFor returns the resolved detail URL for a given tax rate ID.
func (r Routes) DetailFor(id string) string {
	return resolveParam(r.DetailURL, "id", id)
}

// ---------------------------------------------------------------------------
// resolveParam — internal URL template helper
// ---------------------------------------------------------------------------

// resolveParam replaces a single {placeholder} in a URL pattern with value.
// It is the internal single-parameter URL resolver; for multi-parameter URLs
// use route.ResolveURL from packages/pyeza-golang/route directly.
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
