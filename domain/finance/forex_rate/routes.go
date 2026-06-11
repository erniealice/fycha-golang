// routes.go defines configurable route structs for the forex_rate entity views.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultRoutes()
//     constructors return structs populated from the package-level route constants
//     defined in routes.go. These serve as sensible defaults for any consumer app.
//   - Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//     Apps can load a JSON config file that maps route keys to custom paths,
//     allowing industry templates (e.g. salon, retail) to rebrand URLs without
//     code changes. The json struct tags on each field support this workflow.
//   - Level 3: App-specific overrides via Go field assignment (optional).
//     After constructing defaults (and optionally applying JSON), consumer apps
//     can directly assign individual struct fields for one-off customizations.
//
// RouteMap() methods return a map[string]string of dot-notation keys to route
// paths, useful for template rendering and route resolution at runtime.
package forex_rate

const (
	// Finance — Forex Rates (read-only in UI; appended only via RecordOperatorRate)
	ListURL   = "/finance/forex-rates/list/{status}"
	DetailURL = "/finance/forex-rates/detail/{id}"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Forex Rate read-only views.
// Forex rates are read-only in the UI; rows are appended only via RecordOperatorRate.
type Routes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	DetailURL string `json:"detail_url"`
}

// DefaultRoutes returns a Routes populated from the
// package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav: "forex_rate",
		ListURL:   ListURL,
		DetailURL: DetailURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"forex_rate.list":   r.ListURL,
		"forex_rate.detail": r.DetailURL,
	}
}

// DetailFor returns the resolved detail URL for a given forex rate ID.
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
