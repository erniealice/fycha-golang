package charge_policy

import "strings"

// Charge policy route constants (ui-charge-policy.md §3.3, LOCKED 2026-09-28).
// Detail and version pages live under the Ledger › Settings namespace; the
// version page nests under its policy like price_schedule → plan.
const (
	ListURL   = "/ledger/settings/charge-policies/list/{status}"
	TableURL  = "/action/charge-policy/table/{status}"
	DetailURL = "/ledger/settings/charge-policies/detail/{id}"

	TabActionURL = "/action/charge-policy/{id}/tab/{tab}"

	AddURL        = "/action/charge-policy/add"
	EditURL       = "/action/charge-policy/edit/{id}"
	DeleteURL     = "/action/charge-policy/delete"
	BulkDeleteURL = "/action/charge-policy/bulk-delete"
	RetireURL     = "/action/charge-policy/retire/{id}"

	VersionDetailURL    = "/ledger/settings/charge-policies/detail/{id}/version/{vid}"
	VersionTabActionURL = "/action/charge-policy/{id}/version/{vid}/tab/{tab}"
	VersionAddURL       = "/action/charge-policy/{id}/version/add"
	VersionEditURL      = "/action/charge-policy/{id}/version/edit/{vid}"
	VersionDeleteURL    = "/action/charge-policy/{id}/version/delete/{vid}"
	VersionApproveURL   = "/action/charge-policy/{id}/version/approve/{vid}"

	ComponentAddURL    = "/action/charge-policy/{id}/version/{vid}/components/add"
	ComponentEditURL   = "/action/charge-policy/{id}/version/{vid}/components/edit/{cid}"
	ComponentDeleteURL = "/action/charge-policy/{id}/version/{vid}/components/delete"

	PostingAddURL    = "/action/charge-policy/{id}/version/{vid}/postings/add"
	PostingEditURL   = "/action/charge-policy/{id}/version/{vid}/postings/edit/{pid}"
	PostingDeleteURL = "/action/charge-policy/{id}/version/{vid}/postings/delete"

	AttachmentUploadURL = "/action/charge-policy/{id}/version/{vid}/attachments/upload"
	AttachmentDeleteURL = "/action/charge-policy/{id}/version/{vid}/attachments/delete"
)

// Routes holds every charge policy route. JSON tags are the lyngua route.json
// keys (general/route.json carries "charge_policy": {} — defaults live here).
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`

	ListURL             string `json:"list_url"`
	TableURL            string `json:"table_url"`
	DetailURL           string `json:"detail_url"`
	TabActionURL        string `json:"tab_action_url"`
	AddURL              string `json:"add_url"`
	EditURL             string `json:"edit_url"`
	DeleteURL           string `json:"delete_url"`
	BulkDeleteURL       string `json:"bulk_delete_url"`
	RetireURL           string `json:"retire_url"`
	VersionDetailURL    string `json:"version_detail_url"`
	VersionTabActionURL string `json:"version_tab_action_url"`
	VersionAddURL       string `json:"version_add_url"`
	VersionEditURL      string `json:"version_edit_url"`
	VersionDeleteURL    string `json:"version_delete_url"`
	VersionApproveURL   string `json:"version_approve_url"`
	ComponentAddURL     string `json:"component_add_url"`
	ComponentEditURL    string `json:"component_edit_url"`
	ComponentDeleteURL  string `json:"component_delete_url"`
	PostingAddURL       string `json:"posting_add_url"`
	PostingEditURL      string `json:"posting_edit_url"`
	PostingDeleteURL    string `json:"posting_delete_url"`
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

// DefaultRoutes returns the route table populated from the constants.
func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:           "ledger",
		ActiveSubNav:        "charge-policies",
		ListURL:             ListURL,
		TableURL:            TableURL,
		DetailURL:           DetailURL,
		TabActionURL:        TabActionURL,
		AddURL:              AddURL,
		EditURL:             EditURL,
		DeleteURL:           DeleteURL,
		BulkDeleteURL:       BulkDeleteURL,
		RetireURL:           RetireURL,
		VersionDetailURL:    VersionDetailURL,
		VersionTabActionURL: VersionTabActionURL,
		VersionAddURL:       VersionAddURL,
		VersionEditURL:      VersionEditURL,
		VersionDeleteURL:    VersionDeleteURL,
		VersionApproveURL:   VersionApproveURL,
		ComponentAddURL:     ComponentAddURL,
		ComponentEditURL:    ComponentEditURL,
		ComponentDeleteURL:  ComponentDeleteURL,
		PostingAddURL:       PostingAddURL,
		PostingEditURL:      PostingEditURL,
		PostingDeleteURL:    PostingDeleteURL,
		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
	}
}

// RouteMap returns dot-notation route keys (ui-charge-policy.md §3.3).
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"charge_policy.list":               r.ListURL,
		"charge_policy.table":              r.TableURL,
		"charge_policy.detail":             r.DetailURL,
		"charge_policy.tab_action":         r.TabActionURL,
		"charge_policy.add":                r.AddURL,
		"charge_policy.edit":               r.EditURL,
		"charge_policy.delete":             r.DeleteURL,
		"charge_policy.bulk_delete":        r.BulkDeleteURL,
		"charge_policy.retire":             r.RetireURL,
		"charge_policy.version.detail":     r.VersionDetailURL,
		"charge_policy.version.tab_action": r.VersionTabActionURL,
		"charge_policy.version.add":        r.VersionAddURL,
		"charge_policy.version.edit":       r.VersionEditURL,
		"charge_policy.version.delete":     r.VersionDeleteURL,
		"charge_policy.version.approve":    r.VersionApproveURL,
		"charge_policy.component.add":      r.ComponentAddURL,
		"charge_policy.component.edit":     r.ComponentEditURL,
		"charge_policy.component.delete":   r.ComponentDeleteURL,
		"charge_policy.posting.add":        r.PostingAddURL,
		"charge_policy.posting.edit":       r.PostingEditURL,
		"charge_policy.posting.delete":     r.PostingDeleteURL,
		"charge_policy.attachment.upload":  r.AttachmentUploadURL,
		"charge_policy.attachment.delete":  r.AttachmentDeleteURL,
	}
}

// Statuses are the canonical list-status path values.
var Statuses = []string{"active", "draft", "retired"}

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
