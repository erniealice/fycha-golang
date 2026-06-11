package asset

// ---------------------------------------------------------------------------
// Depreciation Policies labels (Surface F)
// Lyngua root key: "depreciationPolicies"
// ---------------------------------------------------------------------------

// DepreciationPoliciesLabels holds all translatable strings for the actionable
// Depreciation Policies page (Surface F).
type DepreciationPoliciesLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Table column headers
	Columns DepreciationPoliciesColumnLabels `json:"columns"`
	// Row action labels
	ActionPreview string `json:"actionPreview"`
	ActionRun     string `json:"actionRun"`
	// Preview drawer subtitle (read-only candidate preview, no DB writes)
	PreviewSubtitle string `json:"previewSubtitle"`
	// Run drawer subtitle (opens Surface C drawer with policy breadcrumb)
	RunSubtitle string `json:"runSubtitle"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	// UsefulLifeMonthsSuffix is appended after the numeric useful_life_months value
	// in the table (e.g. "60" + " mo" → "60 mo"). Leading space is intentional.
	UsefulLifeMonthsSuffix string `json:"usefulLifeMonthsSuffix"`
}

// DepreciationPoliciesColumnLabels holds column headers for the Surface F policies table.
type DepreciationPoliciesColumnLabels struct {
	Policy          string `json:"policy"`
	Method          string `json:"method"`
	UsefulLife      string `json:"usefulLife"`
	SalvagePct      string `json:"salvagePct"`
	AssetsInPolicy  string `json:"assetsInPolicy"`
	AssetsDeviating string `json:"assetsDeviating"`
	Actions         string `json:"actions"`
}

// DefaultDepreciationPoliciesLabels returns DepreciationPoliciesLabels with sensible English defaults.
func DefaultDepreciationPoliciesLabels() DepreciationPoliciesLabels {
	return DepreciationPoliciesLabels{
		Title:    "Depreciation Policies",
		Subtitle: "Manage depreciation policies across asset categories",
		Columns: DepreciationPoliciesColumnLabels{
			Policy:          "Policy",
			Method:          "Method",
			UsefulLife:      "Useful Life",
			SalvagePct:      "Salvage %",
			AssetsInPolicy:  "Assets in policy",
			AssetsDeviating: "Assets deviating",
			Actions:         "Actions",
		},
		ActionPreview:          "Preview",
		ActionRun:              "Run",
		PreviewSubtitle:        "Preview depreciation amounts for this policy (no changes will be posted)",
		RunSubtitle:            "Post depreciation for all in-service assets under this policy",
		EmptyTitle:             "No depreciation policies",
		EmptyMessage:           "Add an asset category to define a depreciation policy.",
		UsefulLifeMonthsSuffix: " mo",
	}
}
