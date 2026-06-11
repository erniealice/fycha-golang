package tax

const (
	// Tax — Tax Rates (read-only; supersession via admin SQL recipe)
	TaxRateListURL   = "/tax/tax-rates/list/{status}"
	TaxRateDetailURL = "/tax/tax-rates/detail/{id}"
)

// TaxRateRoutes holds route paths for Tax Rate read-only views.
// Tax rates are read-only in the UI; supersession is done via admin SQL recipe.
type TaxRateRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	DetailURL string `json:"detail_url"`
}

// DefaultTaxRateRoutes returns a TaxRateRoutes populated from the package-level
// route constants.
func DefaultTaxRateRoutes() TaxRateRoutes {
	return TaxRateRoutes{
		ActiveNav: "tax_rate",
		ListURL:   TaxRateListURL,
		DetailURL: TaxRateDetailURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r TaxRateRoutes) RouteMap() map[string]string {
	return map[string]string{
		"tax_rate.list":   r.ListURL,
		"tax_rate.detail": r.DetailURL,
	}
}

// DetailFor returns the resolved detail URL for a given tax rate ID.
func (r TaxRateRoutes) DetailFor(id string) string {
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
