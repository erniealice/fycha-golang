package document_series

import "strings"

// Document series settings routes (build-spec §6.5: Ledger › Settings, beside
// charge policies). Retire is an edit-permission status change; there is no
// detail page (the list row and drawer carry everything).
const (
	ListURL   = "/ledger/settings/document-series/list/{status}"
	TableURL  = "/action/document-series/table/{status}"
	AddURL    = "/action/document-series/add"
	EditURL   = "/action/document-series/edit/{id}"
	RetireURL = "/action/document-series/retire/{id}"
)

// Routes holds every document series route. JSON tags are the lyngua route.json
// keys (general/route.json carries "document_series": {} — defaults live here).
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL   string `json:"list_url"`
	TableURL  string `json:"table_url"`
	AddURL    string `json:"add_url"`
	EditURL   string `json:"edit_url"`
	RetireURL string `json:"retire_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "ledger",
		ActiveSubNav: "document-series",
		ListURL:      ListURL,
		TableURL:     TableURL,
		AddURL:       AddURL,
		EditURL:      EditURL,
		RetireURL:    RetireURL,
	}
}

// RouteMap returns dot-notation route keys.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"document_series.list":   r.ListURL,
		"document_series.table":  r.TableURL,
		"document_series.add":    r.AddURL,
		"document_series.edit":   r.EditURL,
		"document_series.retire": r.RetireURL,
	}
}

// Statuses are the canonical list-status path values.
var Statuses = []string{"active", "retired"}

// NormalizeStatus maps an unknown list status to "active".
func NormalizeStatus(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	for _, v := range Statuses {
		if s == v {
			return s
		}
	}
	return "active"
}
