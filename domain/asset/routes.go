package asset

// Asset routes — URL consts copied verbatim from packages/fycha-golang/routes.go.

const (
	// Asset routes
	AssetDashboardURL        = "/asset/dashboard"
	AssetListURL             = "/asset/list/{status}"
	AssetDetailURL           = "/asset/detail/{id}"
	AssetTabActionURL        = "/action/asset/{id}/tab/{tab}"
	AssetAttachmentUploadURL = "/action/asset/{id}/attachments/upload"
	AssetAttachmentDeleteURL = "/action/asset/{id}/attachments/delete"
	AssetTableURL            = "/action/asset/table/{status}"
	AssetAddURL              = "/action/asset/add"
	AssetEditURL             = "/action/asset/edit/{id}"
	AssetDeleteURL           = "/action/asset/delete"
	AssetBulkDeleteURL       = "/action/asset/bulk-delete"
	AssetSetStatusURL        = "/action/asset/set-status"
	AssetBulkSetStatusURL    = "/action/asset/bulk-set-status"

	// Asset report/settings routes (legacy mock paths — kept for backwards compat)
	AssetLapsingScheduleURL      = "/asset/reports/lapsing-schedule"
	AssetDepreciationPoliciesURL = "/asset/settings/depreciation-policies"

	// Asset depreciation-run drawer routes (Surface A + E)
	AssetDepreciationRunURL    = "/action/asset/depreciation-run/{asset_id}"
	AssetRevaluationURL        = "/action/asset/revaluation/{asset_id}"
	AssetRevaluationPreviewURL = "/action/asset/revaluation-preview/{asset_id}"

	// Asset category / policy depreciation drawer routes (Surface C + F)
	AssetCategoryDepreciationRunURL   = "/action/asset-category/depreciation-run/{category_id}"
	AssetPolicyDepreciationRunURL     = "/action/asset-policy/depreciation-run/{category_id}"
	AssetPolicyDepreciationPreviewURL = "/action/asset-policy/depreciation-preview/{category_id}"

	// Lapsing schedule page routes (Surface B)
	LapsingScheduleListURL               = "/asset/lapsing-schedule/list"
	LapsingScheduleBulkRunSelectedURL    = "/action/lapsing-schedule/bulk-run-selected"
	LapsingScheduleBulkRunAllMatchingURL = "/action/lapsing-schedule/bulk-run-all-matching"
	AssetListBulkRunSelectedURL          = "/action/asset/bulk-run-selected"

	// Depreciation run history page routes (Surface D)
	DepreciationRunListURL            = "/asset/depreciation-runs/list/{status}"
	DepreciationRunListTableURL       = "/action/depreciation-run/table/{status}"
	DepreciationRunDetailURL          = "/asset/depreciation-runs/detail/{run_id}"
	DepreciationRunDetailTabActionURL = "/action/depreciation-run/detail/{run_id}/tab/{tab}"

	// Depreciation policies page route (Surface F — replaces mock)
	DepreciationPoliciesURL = "/asset/settings/depreciation-policies"
)

// ---------------------------------------------------------------------------
// AssetRoutes
// ---------------------------------------------------------------------------

// AssetRoutes holds route paths for fixed asset management views.
type AssetRoutes struct {
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

// DefaultAssetRoutes returns an AssetRoutes populated from package-level consts.
func DefaultAssetRoutes() AssetRoutes {
	return AssetRoutes{
		DashboardURL:     AssetDashboardURL,
		ListURL:          AssetListURL,
		DetailURL:        AssetDetailURL,
		TabActionURL:     AssetTabActionURL,
		TableURL:         AssetTableURL,
		AddURL:           AssetAddURL,
		EditURL:          AssetEditURL,
		DeleteURL:        AssetDeleteURL,
		BulkDeleteURL:    AssetBulkDeleteURL,
		SetStatusURL:     AssetSetStatusURL,
		BulkSetStatusURL: AssetBulkSetStatusURL,

		AttachmentUploadURL: AssetAttachmentUploadURL,
		AttachmentDeleteURL: AssetAttachmentDeleteURL,

		LapsingScheduleURL:      AssetLapsingScheduleURL,
		DepreciationPoliciesURL: AssetDepreciationPoliciesURL,

		DepreciationRunURL:    AssetDepreciationRunURL,
		RevaluationURL:        AssetRevaluationURL,
		RevaluationPreviewURL: AssetRevaluationPreviewURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r AssetRoutes) RouteMap() map[string]string {
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
func (r AssetRoutes) DepreciationRunFor(assetID string) string {
	return resolveParam(r.DepreciationRunURL, "asset_id", assetID)
}

// RevaluationFor returns the resolved Surface E revaluation drawer URL for the
// given asset ID.
func (r AssetRoutes) RevaluationFor(assetID string) string {
	return resolveParam(r.RevaluationURL, "asset_id", assetID)
}

// RevaluationPreviewFor returns the resolved HTMX revaluation-preview endpoint
// URL for the given asset ID.
func (r AssetRoutes) RevaluationPreviewFor(assetID string) string {
	return resolveParam(r.RevaluationPreviewURL, "asset_id", assetID)
}

// ---------------------------------------------------------------------------
// LapsingScheduleRoutes
// ---------------------------------------------------------------------------

// LapsingScheduleRoutes holds route paths for the lapsing-schedule live page
// (Surface B) and its bulk-action endpoints.
type LapsingScheduleRoutes struct {
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

// DefaultLapsingScheduleRoutes returns a LapsingScheduleRoutes populated from
// the package-level route constants.
func DefaultLapsingScheduleRoutes() LapsingScheduleRoutes {
	return LapsingScheduleRoutes{
		ListURL:                     LapsingScheduleListURL,
		BulkRunSelectedURL:          LapsingScheduleBulkRunSelectedURL,
		BulkRunAllMatchingURL:       LapsingScheduleBulkRunAllMatchingURL,
		AssetListBulkRunSelectedURL: AssetListBulkRunSelectedURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r LapsingScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"lapsing_schedule.list":                         r.ListURL,
		"lapsing_schedule.bulk_run_selected":            r.BulkRunSelectedURL,
		"lapsing_schedule.bulk_run_all_matching":        r.BulkRunAllMatchingURL,
		"lapsing_schedule.asset_list_bulk_run_selected": r.AssetListBulkRunSelectedURL,
	}
}

// ---------------------------------------------------------------------------
// DepreciationRunRoutes
// ---------------------------------------------------------------------------

// DepreciationRunRoutes holds route paths for the depreciation-run history
// list and detail pages (Surface D).
type DepreciationRunRoutes struct {
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

// DefaultDepreciationRunRoutes returns a DepreciationRunRoutes populated from
// the package-level route constants.
func DefaultDepreciationRunRoutes() DepreciationRunRoutes {
	return DepreciationRunRoutes{
		ActiveNav:          "depreciation-runs",
		ListURL:            DepreciationRunListURL,
		ListTableURL:       DepreciationRunListTableURL,
		DetailURL:          DepreciationRunDetailURL,
		DetailTabActionURL: DepreciationRunDetailTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r DepreciationRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"depreciation_run.list":              r.ListURL,
		"depreciation_run.list_table":        r.ListTableURL,
		"depreciation_run.detail":            r.DetailURL,
		"depreciation_run.detail_tab_action": r.DetailTabActionURL,
	}
}

// ListFor returns the resolved depreciation-run list URL for the given status
// (e.g. "complete", "failed", "pending", or empty string for all).
func (r DepreciationRunRoutes) ListFor(status string) string {
	return resolveParam(r.ListURL, "status", status)
}

// DetailFor returns the resolved depreciation-run detail URL for the given
// run ID.
func (r DepreciationRunRoutes) DetailFor(runID string) string {
	return resolveParam(r.DetailURL, "run_id", runID)
}

// ---------------------------------------------------------------------------
// AssetCategoryDepreciationRoutes
// ---------------------------------------------------------------------------

// AssetCategoryDepreciationRoutes holds route paths for the per-category and
// per-policy depreciation drawer endpoints (Surfaces C and F).
type AssetCategoryDepreciationRoutes struct {
	// CategoryRunURL is the Surface C per-category run drawer.
	CategoryRunURL string `json:"category_run_url"`

	// PolicyRunURL is the Surface C per-policy run drawer (policy breadcrumb).
	PolicyRunURL string `json:"policy_run_url"`

	// PolicyPreviewURL is the Surface F preview drawer (no writes).
	PolicyPreviewURL string `json:"policy_preview_url"`

	// DepreciationPoliciesURL is the Surface F actionable policies page.
	DepreciationPoliciesURL string `json:"depreciation_policies_url"`
}

// DefaultAssetCategoryDepreciationRoutes returns an
// AssetCategoryDepreciationRoutes populated from the package-level route
// constants.
func DefaultAssetCategoryDepreciationRoutes() AssetCategoryDepreciationRoutes {
	return AssetCategoryDepreciationRoutes{
		CategoryRunURL:          AssetCategoryDepreciationRunURL,
		PolicyRunURL:            AssetPolicyDepreciationRunURL,
		PolicyPreviewURL:        AssetPolicyDepreciationPreviewURL,
		DepreciationPoliciesURL: DepreciationPoliciesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r AssetCategoryDepreciationRoutes) RouteMap() map[string]string {
	return map[string]string{
		"asset_category_depreciation.category_run":          r.CategoryRunURL,
		"asset_category_depreciation.policy_run":            r.PolicyRunURL,
		"asset_category_depreciation.policy_preview":        r.PolicyPreviewURL,
		"asset_category_depreciation.depreciation_policies": r.DepreciationPoliciesURL,
	}
}

// CategoryRunFor returns the resolved Surface C per-category run drawer URL
// for the given category ID.
func (r AssetCategoryDepreciationRoutes) CategoryRunFor(categoryID string) string {
	return resolveParam(r.CategoryRunURL, "category_id", categoryID)
}

// PolicyRunFor returns the resolved Surface C per-policy run drawer URL for
// the given category ID (policy scope).
func (r AssetCategoryDepreciationRoutes) PolicyRunFor(categoryID string) string {
	return resolveParam(r.PolicyRunURL, "category_id", categoryID)
}

// PolicyPreviewFor returns the resolved Surface F preview drawer URL for the
// given category ID.
func (r AssetCategoryDepreciationRoutes) PolicyPreviewFor(categoryID string) string {
	return resolveParam(r.PolicyPreviewURL, "category_id", categoryID)
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
