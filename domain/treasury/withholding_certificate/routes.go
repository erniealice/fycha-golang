// routes.go defines configurable route structs for the withholding_certificate entity.
package withholding_certificate

const (
	// Treasury — Withholding Certificates (full CRUD)
	ListURL       = "/treasury/withholding-certificates/list/{status}"
	DetailURL     = "/treasury/withholding-certificates/detail/{id}"
	TableURL      = "/action/withholding-certificate/table/{status}"
	AddURL        = "/action/withholding-certificate/add"
	EditURL       = "/action/withholding-certificate/edit/{id}"
	DeleteURL     = "/action/withholding-certificate/delete"
	BulkDeleteURL = "/action/withholding-certificate/bulk-delete"
	SetStatusURL  = "/action/withholding-certificate/set-status"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Withholding Certificate CRUD views
// (Treasury domain — tax integration v1).
type Routes struct {
	ActiveNav     string `json:"active_nav"`
	ListURL       string `json:"list_url"`
	DetailURL     string `json:"detail_url"`
	TableURL      string `json:"table_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	SetStatusURL  string `json:"set_status_url"`
}

// DefaultRoutes returns a Routes populated from the package-level route constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:     "withholding_certificate",
		ListURL:       ListURL,
		DetailURL:     DetailURL,
		TableURL:      TableURL,
		AddURL:        AddURL,
		EditURL:       EditURL,
		DeleteURL:     DeleteURL,
		BulkDeleteURL: BulkDeleteURL,
		SetStatusURL:  SetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"withholding_certificate.list":        r.ListURL,
		"withholding_certificate.detail":      r.DetailURL,
		"withholding_certificate.table":       r.TableURL,
		"withholding_certificate.add":         r.AddURL,
		"withholding_certificate.edit":        r.EditURL,
		"withholding_certificate.delete":      r.DeleteURL,
		"withholding_certificate.bulk_delete": r.BulkDeleteURL,
		"withholding_certificate.set_status":  r.SetStatusURL,
	}
}

// DetailFor returns the resolved detail URL for a given withholding certificate ID.
func (r Routes) DetailFor(id string) string {
	return resolveParam(r.DetailURL, "id", id)
}

// EditFor returns the resolved edit drawer URL for a given withholding certificate ID.
func (r Routes) EditFor(id string) string {
	return resolveParam(r.EditURL, "id", id)
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
