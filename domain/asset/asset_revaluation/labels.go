package asset_revaluation

// ---------------------------------------------------------------------------
// Asset Revaluation labels (Surface E)
// Lyngua root key: "assetRevaluation"
// ---------------------------------------------------------------------------

// ErrorLabels holds error message strings for the asset revaluation drawer.
type ErrorLabels struct {
	UseCaseUnavailable string `json:"useCaseUnavailable"`
	FormParseFailed    string `json:"formParseFailed"`
	RevaluateFailed    string `json:"revaluateFailed"`
	InvalidAmount      string `json:"invalidAmount"`
	// 2026-05-14 permission-gates: AWS-style permission tooltip surface for the
	// revaluation drawer. Emitted when the caller lacks asset_revaluation:*.
	PermissionDenied string `json:"permissionDenied"`
}

// Labels holds all translatable strings for the Asset Revaluation drawer.
type Labels struct {
	Title string `json:"title"`
	// Form field labels
	NewFairValue    string `json:"newFairValue"`
	AppraiserName   string `json:"appraiserName"`
	ValuationMethod string `json:"valuationMethod"`
	Notes           string `json:"notes"`
	// Preview section labels
	PreviewTitle      string `json:"previewTitle"`
	RevaluationAmount string `json:"revaluationAmount"`
	PnlSplit          string `json:"pnlSplit"`
	OciSplit          string `json:"ociSplit"`
	// Submit and cancel
	SubmitLabel string `json:"submitLabel"`
	CancelLabel string `json:"cancelLabel"`
	// Toast message template (supports {{.Direction}}/{{.Amount}}/{{.Recognition}} placeholders)
	ToastSuccessTemplate string `json:"toastSuccessTemplate"`
	// Errors
	Errors ErrorLabels `json:"errors"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Title:                "Revalue Asset",
		NewFairValue:         "New fair value",
		AppraiserName:        "Appraiser name",
		ValuationMethod:      "Valuation method",
		Notes:                "Notes",
		PreviewTitle:         "Revaluation preview",
		RevaluationAmount:    "Revaluation amount",
		PnlSplit:             "Recognized in P&L",
		OciSplit:             "Recognized in OCI",
		SubmitLabel:          "Revalue",
		CancelLabel:          "Cancel",
		ToastSuccessTemplate: "Asset revalued: {{.Direction}}{{.Amount}} recognized in {{.Recognition}}",
		Errors: ErrorLabels{
			UseCaseUnavailable: "Service unavailable. Please try again.",
			FormParseFailed:    "Form data could not be read.",
			RevaluateFailed:    "Failed to record the revaluation.",
			InvalidAmount:      "Invalid amount format. Use a number with up to 2 decimal places.",
			PermissionDenied:   "You do not have permission to revalue this asset.",
		},
	}
}
