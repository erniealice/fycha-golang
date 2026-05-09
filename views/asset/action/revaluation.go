package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Surface E — per-asset revaluation drawer
// URL: /action/asset/revaluation/{asset_id}
// Preview: /action/asset/revaluation-preview/{asset_id}
// Visible only when asset.measurement_model = REVALUATION.
// ---------------------------------------------------------------------------

// RevaluationDeps holds dependencies for the Surface E drawer.
type RevaluationDeps struct {
	Routes fycha.AssetRoutes
	Labels fycha.AssetRevaluationLabels
	// RevalueAsset posts the revaluation via espyna.
	RevalueAsset func(ctx context.Context, req RevaluationRequest) (*RevaluationResult, error)
	// PreviewRevaluation computes the PnL/OCI split without writing.
	// If nil, the preview is skipped (split shown as unknown).
	PreviewRevaluation func(ctx context.Context, assetID string, newFairValue int64) (*RevaluationPreview, error)
}

// RevaluationRequest is the POST payload to RevalueAsset.
type RevaluationRequest struct {
	AssetID         string
	NewFairValue    int64  // centavos
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
	FormAction          string
	PreviewURL          string
	AssetID             string
	Preview             *RevaluationPreview
	Labels              fycha.AssetRevaluationLabels
}

// NewRevaluationAction creates the Surface E per-asset revaluation drawer.
func NewRevaluationAction(deps *RevaluationDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
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
		assetID := viewCtx.Request.PathValue("asset_id")
		if assetID == "" {
			assetID = viewCtx.Request.PathValue("id")
		}

		if err := viewCtx.Request.ParseForm(); err != nil {
			return view.OK("asset-revaluation-preview-partial", &revaluationFormData{Labels: deps.Labels})
		}

		newFairValueCents := parseCents(viewCtx.Request.FormValue("new_fair_value"))
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
	if err := viewCtx.Request.ParseForm(); err != nil {
		return fycha.HTMXError(deps.Labels.SubmitLabel)
	}

	req := RevaluationRequest{
		AssetID:         assetID,
		NewFairValue:    parseCents(viewCtx.Request.FormValue("new_fair_value")),
		AppraiserName:   viewCtx.Request.FormValue("appraiser_name"),
		ValuationMethod: viewCtx.Request.FormValue("valuation_method"),
		Notes:           viewCtx.Request.FormValue("notes"),
	}

	var result *RevaluationResult
	if deps.RevalueAsset != nil {
		var err error
		result, err = deps.RevalueAsset(ctx, req)
		if err != nil {
			log.Printf("RevalueAsset error: %v", err)
			return fycha.HTMXError(deps.Labels.SubmitLabel)
		}
	} else {
		result = &RevaluationResult{
			TransactionID: "tx-mock",
			Direction:     "UP",
			Amount:        req.NewFairValue,
			AmountFmt:     fmt.Sprintf("₱%.2f", float64(req.NewFairValue)/100),
			Recognition:   "OCI",
			Success:       true,
		}
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

// parseCents parses a decimal peso string into centavos (int64).
// Input format: "1234.56" → 123456. Returns 0 on parse failure.
func parseCents(s string) int64 {
	if s == "" {
		return 0
	}
	var pesos float64
	if _, err := fmt.Sscanf(s, "%f", &pesos); err != nil {
		return 0
	}
	return int64(pesos * 100)
}
