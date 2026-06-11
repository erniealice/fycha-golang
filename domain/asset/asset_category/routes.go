package asset_category

// Asset category / policy depreciation drawer routes (Surfaces C and F)

const (
	// Asset category / policy depreciation drawer routes (Surface C + F)
	DepreciationRunURL      = "/action/asset-category/depreciation-run/{category_id}"
	PolicyRunURL            = "/action/asset-policy/depreciation-run/{category_id}"
	PolicyPreviewURL        = "/action/asset-policy/depreciation-preview/{category_id}"
	DepreciationPoliciesURL = "/asset/settings/depreciation-policies"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for the per-category and
// per-policy depreciation drawer endpoints (Surfaces C and F).
type Routes struct {
	// DepreciationRunURL is the Surface C per-category run drawer.
	DepreciationRunURL string `json:"category_run_url"`

	// PolicyRunURL is the Surface C per-policy run drawer (policy breadcrumb).
	PolicyRunURL string `json:"policy_run_url"`

	// PolicyPreviewURL is the Surface F preview drawer (no writes).
	PolicyPreviewURL string `json:"policy_preview_url"`

	// DepreciationPoliciesURL is the Surface F actionable policies page.
	DepreciationPoliciesURL string `json:"depreciation_policies_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		DepreciationRunURL:      DepreciationRunURL,
		PolicyRunURL:            PolicyRunURL,
		PolicyPreviewURL:        PolicyPreviewURL,
		DepreciationPoliciesURL: DepreciationPoliciesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"asset_category_depreciation.category_run":          r.DepreciationRunURL,
		"asset_category_depreciation.policy_run":            r.PolicyRunURL,
		"asset_category_depreciation.policy_preview":        r.PolicyPreviewURL,
		"asset_category_depreciation.depreciation_policies": r.DepreciationPoliciesURL,
	}
}

// CategoryRunFor returns the resolved Surface C per-category run drawer URL
// for the given category ID.
func (r Routes) CategoryRunFor(categoryID string) string {
	return resolveParam(r.DepreciationRunURL, "category_id", categoryID)
}

// PolicyRunFor returns the resolved Surface C per-policy run drawer URL for
// the given category ID (policy scope).
func (r Routes) PolicyRunFor(categoryID string) string {
	return resolveParam(r.PolicyRunURL, "category_id", categoryID)
}

// PolicyPreviewFor returns the resolved Surface F preview drawer URL for the
// given category ID.
func (r Routes) PolicyPreviewFor(categoryID string) string {
	return resolveParam(r.PolicyPreviewURL, "category_id", categoryID)
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
