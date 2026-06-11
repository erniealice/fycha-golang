package action

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	asset "github.com/erniealice/fycha-golang/domain/asset/asset"
	assetrevaluation "github.com/erniealice/fycha-golang/domain/asset/asset_revaluation"
)

// ---------------------------------------------------------------------------
// Surface E — per-asset revaluation drawer
// URL: /action/asset/revaluation/{asset_id}
// Preview: /action/asset/revaluation-preview/{asset_id}
// Visible only when asset.measurement_model = REVALUATION.
// ---------------------------------------------------------------------------

// RevaluationDeps holds dependencies for the Surface E drawer.
type RevaluationDeps struct {
	Routes asset.Routes
	Labels assetrevaluation.Labels
	// RevalueAsset posts the revaluation via espyna.
	RevalueAsset func(ctx context.Context, req RevaluationRequest) (*RevaluationResult, error)
	// PreviewRevaluation computes the PnL/OCI split without writing.
	// If nil, the preview is skipped (split shown as unknown).
	PreviewRevaluation func(ctx context.Context, assetID string, newFairValue int64) (*RevaluationPreview, error)
}

// RevaluationRequest is the POST payload to RevalueAsset.
type RevaluationRequest struct {
	AssetID         string
	NewFairValue    int64 // centavos
	AppraiserName   string
	ValuationMethod string
	Notes           string
}

// RevaluationResult is the response from RevalueAsset.
type RevaluationResult struct {
	TransactionID string
	Direction     string // "UP" or "DOWN"
	Amount        int64  // centavos (absolute)
	AmountFmt     string
	Recognition   string // "OCI" or "P&L"
	Success       bool
	ErrorMessage  string
}

// RevaluationPreview holds the live preview of a pending revaluation.
type RevaluationPreview struct {
	RevaluationAmountFmt string
	PnLAmountFmt         string
	OCIAmountFmt         string
	Direction            string // "UP" or "DOWN"
}

// revaluationFormData is the template data for the Surface E drawer.
type revaluationFormData struct {
	FormAction  string
	WorkspaceID string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	PreviewURL  string
	AssetID     string
	Preview     *RevaluationPreview
	Labels      assetrevaluation.Labels
}

// NewRevaluationAction creates the Surface E per-asset revaluation drawer.
func NewRevaluationAction(deps *RevaluationDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates: gate GET drawer + POST submit on
		// asset_revaluation:create (catalog verb).
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset_revaluation", "create") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		assetID := viewCtx.Request.PathValue("asset_id")
		if assetID == "" {
			assetID = viewCtx.Request.PathValue("id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			return handleRevaluationGET(ctx, viewCtx, deps, assetID)
		}
		return handleRevaluationPOST(ctx, viewCtx, deps, assetID)
	})
}

// NewRevaluationPreviewAction creates the HTMX preview partial endpoint.
func NewRevaluationPreviewAction(deps *RevaluationDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates: preview is read-only, gate on
		// asset_revaluation:read.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset_revaluation", "read") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		assetID := viewCtx.Request.PathValue("asset_id")
		if assetID == "" {
			assetID = viewCtx.Request.PathValue("id")
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.OK("asset-revaluation-preview-partial", &revaluationFormData{Labels: deps.Labels})
		}

		newFairValueCents, _ := types.ParseCentavos(viewCtx.Request.FormValue("new_fair_value"))
		var preview *RevaluationPreview
		if deps.PreviewRevaluation != nil && newFairValueCents > 0 {
			var err error
			preview, err = deps.PreviewRevaluation(ctx, assetID, newFairValueCents)
			if err != nil {
				log.Printf("PreviewRevaluation error: %v", err)
			}
		}

		return view.OK("asset-revaluation-preview-partial", &revaluationFormData{
			AssetID: assetID,
			Preview: preview,
			Labels:  deps.Labels,
		})
	})
}

func handleRevaluationGET(_ context.Context, viewCtx *view.ViewContext, deps *RevaluationDeps, assetID string) view.ViewResult {
	return view.OK("asset-revaluation-drawer-form", &revaluationFormData{
		FormAction: deps.Routes.RevaluationFor(assetID),
		PreviewURL: deps.Routes.RevaluationPreviewFor(assetID),
		AssetID:    assetID,
		Labels:     deps.Labels,
	})
}

func handleRevaluationPOST(ctx context.Context, viewCtx *view.ViewContext, deps *RevaluationDeps, assetID string) view.ViewResult { //nolint:unparam
	// Soft-block path (deferred — no field-level validation needed today):
	// For validation errors that need the form re-rendered with a field-level
	// chip rather than a toast, return:
	//     view.ViewResult{
	//         StatusCode: 422,
	//         Headers: map[string]string{
	//             "HX-Reswap":   "outerHTML",
	//             "HX-Retarget": "#sheet form",
	//         },
	//         Body: rerenderedFormHTML,
	//     }
	// lf.Sheet.handleResponse (sheet.js:208-225) honors these headers on non-2xx
	// and swaps the body in. Canonical example: subscription/recognize/action.go:335-345.

	if err := viewCtx.Request.ParseForm(); err != nil {
		return view.HTMXError(deps.Labels.Errors.FormParseFailed)
	}

	newFairValueCents, err := types.ParseCentavos(viewCtx.Request.FormValue("new_fair_value"))
	if err != nil {
		log.Printf("revaluation: invalid amount input: %v", err)
		return view.HTMXError(deps.Labels.Errors.InvalidAmount)
	}

	req := RevaluationRequest{
		AssetID:         assetID,
		NewFairValue:    newFairValueCents,
		AppraiserName:   viewCtx.Request.FormValue("appraiser_name"),
		ValuationMethod: viewCtx.Request.FormValue("valuation_method"),
		Notes:           viewCtx.Request.FormValue("notes"),
	}

	// Phase 3 (codex C4 part A): RevalueAsset is wired by fycha block.go
	// against the espyna consumer. The previous tx-mock fallback was
	// removed — when the wiring is missing in production it indicates a
	// real bootstrap defect and we surface a service-unavailable error.
	if deps.RevalueAsset == nil {
		log.Printf("revalue_asset: RevalueAsset callback not wired (codex C4 — service unavailable)")
		return view.HTMXError(deps.Labels.Errors.UseCaseUnavailable)
	}
	result, err := deps.RevalueAsset(ctx, req)
	if err != nil {
		log.Printf("RevalueAsset error: %v", err)
		return view.HTMXError(deps.Labels.Errors.RevaluateFailed)
	}

	detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", assetID)
	if result.TransactionID != "" {
		detailURL += "#tx-" + result.TransactionID
	}

	toastMsg := strings.NewReplacer(
		"{{.Direction}}", result.Direction,
		"{{.Amount}}", result.AmountFmt,
		"{{.Recognition}}", result.Recognition,
	).Replace(deps.Labels.ToastSuccessTemplate)

	triggerPayload, _ := json.Marshal(map[string]any{
		"pyeza:toast": map[string]any{
			"state":   "success",
			"message": toastMsg,
			"link": map[string]string{
				"url":   detailURL,
				"label": deps.Labels.SubmitLabel,
			},
		},
		"refreshTable": "asset-detail-history",
	})

	return view.ViewResult{
		StatusCode: 200,
		Headers: map[string]string{
			"HX-Trigger": string(triggerPayload),
		},
	}
}
