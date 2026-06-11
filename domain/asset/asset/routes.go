package asset

// Asset routes — URL consts copied verbatim from domain/asset/routes.go.

const (
	// Asset routes
	DashboardURL        = "/asset/dashboard"
	ListURL             = "/asset/list/{status}"
	DetailURL           = "/asset/detail/{id}"
	TabActionURL        = "/action/asset/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/asset/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/asset/{id}/attachments/delete"
	TableURL            = "/action/asset/table/{status}"
	AddURL              = "/action/asset/add"
	EditURL             = "/action/asset/edit/{id}"
	DeleteURL           = "/action/asset/delete"
	BulkDeleteURL       = "/action/asset/bulk-delete"
	SetStatusURL        = "/action/asset/set-status"
	BulkSetStatusURL    = "/action/asset/bulk-set-status"

	// Asset report/settings routes (legacy mock paths — kept for backwards compat)
	LapsingScheduleURL      = "/asset/reports/lapsing-schedule"
	DepreciationPoliciesURL = "/asset/settings/depreciation-policies"

	// Asset depreciation-run drawer routes (Surface A + E)
	DepreciationRunURL    = "/action/asset/depreciation-run/{asset_id}"
	RevaluationURL        = "/action/asset/revaluation/{asset_id}"
	RevaluationPreviewURL = "/action/asset/revaluation-preview/{asset_id}"

	// Lapsing schedule bulk-run endpoint on the main assets list page
	ListBulkRunSelectedURL = "/action/asset/bulk-run-selected"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for fixed asset management views.
type Routes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Report/settings routes (legacy mock paths)
	LapsingScheduleURL      string `json:"lapsing_schedule_url"`
	DepreciationPoliciesURL string `json:"depreciation_policies_url"`

	// Depreciation-run drawer routes (Surface A)
	DepreciationRunURL    string `json:"depreciation_run_url"`
	RevaluationURL        string `json:"revaluation_url"`
	RevaluationPreviewURL string `json:"revaluation_preview_url"`
}

// DefaultRoutes returns a Routes populated from package-level consts.
func DefaultRoutes() Routes {
	return Routes{
		DashboardURL:     DashboardURL,
		ListURL:          ListURL,
		DetailURL:        DetailURL,
		TabActionURL:     TabActionURL,
		TableURL:         TableURL,
		AddURL:           AddURL,
		EditURL:          EditURL,
		DeleteURL:        DeleteURL,
		BulkDeleteURL:    BulkDeleteURL,
		SetStatusURL:     SetStatusURL,
		BulkSetStatusURL: BulkSetStatusURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,

		LapsingScheduleURL:      LapsingScheduleURL,
		DepreciationPoliciesURL: DepreciationPoliciesURL,

		DepreciationRunURL:    DepreciationRunURL,
		RevaluationURL:        RevaluationURL,
		RevaluationPreviewURL: RevaluationPreviewURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"asset.dashboard":       r.DashboardURL,
		"asset.list":            r.ListURL,
		"asset.detail":          r.DetailURL,
		"asset.tab_action":      r.TabActionURL,
		"asset.table":           r.TableURL,
		"asset.add":             r.AddURL,
		"asset.edit":            r.EditURL,
		"asset.delete":          r.DeleteURL,
		"asset.bulk_delete":     r.BulkDeleteURL,
		"asset.set_status":      r.SetStatusURL,
		"asset.bulk_set_status": r.BulkSetStatusURL,

		"asset.attachment.upload": r.AttachmentUploadURL,
		"asset.attachment.delete": r.AttachmentDeleteURL,

		"asset.lapsing_schedule":      r.LapsingScheduleURL,
		"asset.depreciation_policies": r.DepreciationPoliciesURL,

		"asset.depreciation_run":    r.DepreciationRunURL,
		"asset.revaluation":         r.RevaluationURL,
		"asset.revaluation_preview": r.RevaluationPreviewURL,
	}
}

// DepreciationRunFor returns the resolved Surface A depreciation-run drawer URL
// for the given asset ID.
func (r Routes) DepreciationRunFor(assetID string) string {
	return resolveParam(r.DepreciationRunURL, "asset_id", assetID)
}

// RevaluationFor returns the resolved Surface E revaluation drawer URL for the
// given asset ID.
func (r Routes) RevaluationFor(assetID string) string {
	return resolveParam(r.RevaluationURL, "asset_id", assetID)
}

// RevaluationPreviewFor returns the resolved HTMX revaluation-preview endpoint
// URL for the given asset ID.
func (r Routes) RevaluationPreviewFor(assetID string) string {
	return resolveParam(r.RevaluationPreviewURL, "asset_id", assetID)
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
