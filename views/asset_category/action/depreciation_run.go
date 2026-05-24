// Package action implements drawer actions for the asset-category domain.
//
// depreciation_run.go — Surface C per-category / per-policy depreciation-run drawer.
//
// URLs served by this handler:
//   - GET/POST /action/asset-category/depreciation-run/{category_id} — category entry (category breadcrumb)
//   - GET/POST /action/asset-policy/depreciation-run/{category_id}   — policy entry (policy breadcrumb)
//
// The same drawer logic serves both entry points; entry point is distinguished
// by the "scope" query parameter ("category" | "policy") or inferred from the
// URL path prefix. Only the breadcrumb title differs.
package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/view"
)

// CategoryDepreciationRunAssetRow is one asset row inside the Surface C drawer.
type CategoryDepreciationRunAssetRow struct {
	AssetID       string
	AssetName     string
	Currency      string
	BookValue     int64 // centavos
	PendingCount  int
	NextAmount    int64  // centavos — next period projected
	NextAmountFmt string // pre-formatted via types.FormatMoney (e.g. "PHP 50,000.00")
	Blockers      []string
	CanRun        bool // true when no blockers and pendingCount > 0
}

// CategoryDepreciationRunRequest is the POST payload for Surface C.
type CategoryDepreciationRunRequest struct {
	CategoryID string
	ScopeKind  string // "CATEGORY" or "POLICY"
	AsOfDate   string
	AssetIDs   []string // selected asset IDs
}

// CategoryDepreciationRunResult is the response from GenerateDepreciationRun (CATEGORY/POLICY scope).
type CategoryDepreciationRunResult struct {
	RunID        string
	CreatedCount int
	SkippedCount int
	ErroredCount int
	Success      bool
	ErrorMessage string
}

// CategoryDepreciationRunDeps holds dependencies for the Surface C drawer.
type CategoryDepreciationRunDeps struct {
	Routes       fycha.AssetCategoryDepreciationRoutes
	RunRoutes    fycha.DepreciationRunRoutes
	Labels       fycha.DepreciationRunLabels
	CommonLabels pyeza.CommonLabels

	// ListCategoryCandidates returns per-asset candidate rows for the given category/policy.
	// block.go wires this to espyna ListDepreciationCandidates(scope=CATEGORY|POLICY).
	// Nil = empty table (graceful degradation).
	ListCategoryCandidates func(ctx context.Context, categoryID, scopeKind, asOfDate string) ([]CategoryDepreciationRunAssetRow, error)

	// GenerateCategoryRun posts selected periods for all requested assets in the category.
	// block.go wires this to espyna GenerateDepreciationRun(scope=CATEGORY|POLICY).
	// Nil = mock result (graceful degradation during dev).
	GenerateCategoryRun func(ctx context.Context, req CategoryDepreciationRunRequest) (*CategoryDepreciationRunResult, error)
}

// categoryRunFormData is the template context for the Surface C drawer.
type categoryRunFormData struct {
	FormAction    string
	WorkspaceID    string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	CategoryID    string
	ScopeKind     string // "category" | "policy" — controls breadcrumb
	AsOfDate      string
	MaxAsOfDate   string
	Rows          []CategoryDepreciationRunAssetRow
	EligibleCount int
	Labels        fycha.DepreciationRunLabels
	CommonLabels  pyeza.CommonLabels
}

// NewCategoryDepreciationRunAction creates the Surface C per-category/per-policy depreciation-run drawer.
// Handles both the category entry point and the policy entry point — scope is set by the URL
// pattern (/action/asset-category/... vs /action/asset-policy/...) and passed as a query param
// "scope" that block.go injects when registering the route.
func NewCategoryDepreciationRunAction(deps *CategoryDepreciationRunDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates: mutating verb keyed on the
		// depreciation_schedule entity (the verb that creates schedule rows).
		// asset_category:read parent context is enforced by the page-level
		// gate; here the action handler enforces the mutating verb.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "create") {
			return fycha.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		categoryID := viewCtx.Request.PathValue("category_id")
		if categoryID == "" {
			return view.Error(fmt.Errorf("category_id path parameter is required"))
		}

		// Determine scope kind from URL path prefix.
		scopeKind := "category"
		if strings.Contains(viewCtx.Request.URL.Path, "/asset-policy/") {
			scopeKind = "policy"
		}

		if viewCtx.Request.Method == http.MethodGet {
			return handleCategoryRunGET(ctx, viewCtx, deps, categoryID, scopeKind)
		}
		return handleCategoryRunPOST(ctx, viewCtx, deps, categoryID, scopeKind)
	})
}

func handleCategoryRunGET(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *CategoryDepreciationRunDeps,
	categoryID, scopeKind string,
) view.ViewResult {
	asOfDate := viewCtx.Request.URL.Query().Get("as_of_date")
	if asOfDate == "" {
		asOfDate = time.Now().Format("2006-01-02")
	}

	// Determine form action URL based on scope kind.
	var formAction string
	if scopeKind == "policy" {
		formAction = deps.Routes.PolicyRunFor(categoryID)
	} else {
		formAction = deps.Routes.CategoryRunFor(categoryID)
	}

	data := &categoryRunFormData{
		FormAction:   formAction,
		CategoryID:   categoryID,
		ScopeKind:    scopeKind,
		AsOfDate:     asOfDate,
		MaxAsOfDate:  time.Now().Format("2006-01-02"),
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
	}

	if deps.ListCategoryCandidates != nil {
		rows, err := deps.ListCategoryCandidates(ctx, categoryID, scopeKind, asOfDate)
		if err != nil {
			log.Printf("surface-c: ListCategoryCandidates error for category %s: %v", categoryID, err)
		} else {
			data.Rows = rows
			for _, r := range rows {
				if r.CanRun {
					data.EligibleCount++
				}
			}
		}
	}

	return view.OK("asset-category-depreciation-run-form", data)
}

func handleCategoryRunPOST(
	ctx context.Context,
	viewCtx *view.ViewContext,
	deps *CategoryDepreciationRunDeps,
	categoryID, scopeKind string,
) view.ViewResult {
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
		return fycha.HTMXError(deps.Labels.Errors.InvalidSelection)
	}

	asOfDate := viewCtx.Request.FormValue("as_of_date")
	if asOfDate == "" {
		asOfDate = time.Now().Format("2006-01-02")
	}

	selectedAssets := viewCtx.Request.Form["asset_id"]

	// Map "category"/"policy" to proto enum string expected by espyna.
	protoScopeKind := "CATEGORY"
	if scopeKind == "policy" {
		protoScopeKind = "POLICY"
	}

	if deps.GenerateCategoryRun == nil {
		log.Printf("surface-c: GenerateCategoryRun callback not wired (service unavailable)")
		return fycha.HTMXError(deps.Labels.Errors.UseCaseUnavailable)
	}
	result, err := deps.GenerateCategoryRun(ctx, CategoryDepreciationRunRequest{
		CategoryID: categoryID,
		ScopeKind:  protoScopeKind,
		AsOfDate:   asOfDate,
		AssetIDs:   selectedAssets,
	})
	if err != nil {
		log.Printf("surface-c: GenerateCategoryRun error for category %s: %v", categoryID, err)
		return fycha.HTMXError(deps.Labels.Errors.UseCaseUnavailable)
	}

	// Build toast payload per pyeza:toast contract.
	toastState := "success"
	toastMsg := strings.NewReplacer(
		"{{.Created}}", fmt.Sprintf("%d", result.CreatedCount),
		"{{.Skipped}}", fmt.Sprintf("%d", result.SkippedCount),
		"{{.Errored}}", fmt.Sprintf("%d", result.ErroredCount),
	).Replace(deps.Labels.Toast.SuccessTemplate)
	if result.CreatedCount == 0 && result.SkippedCount > 0 {
		toastState = "warning"
		toastMsg = strings.ReplaceAll(deps.Labels.Toast.SkippedTemplate,
			"{{.Skipped}}", fmt.Sprintf("%d", result.SkippedCount))
	}
	if result.ErroredCount > 0 {
		toastState = "errored"
		toastMsg = strings.ReplaceAll(deps.Labels.Toast.ErroredTemplate,
			"{{.Errored}}", fmt.Sprintf("%d", result.ErroredCount))
	}

	// Toast link points to the Surface D run-detail page.
	runDetailURL := ""
	if deps.RunRoutes.DetailURL != "" && result.RunID != "" {
		runDetailURL = deps.RunRoutes.DetailFor(result.RunID)
	}

	triggerPayload, _ := json.Marshal(map[string]any{
		"pyeza:toast": map[string]any{
			"state":   toastState,
			"message": toastMsg,
			"link": map[string]string{
				"url":   runDetailURL,
				"label": deps.Labels.Toast.ViewRunLink,
			},
		},
	})

	return view.ViewResult{
		StatusCode: 200,
		Headers: map[string]string{
			"HX-Trigger": string(triggerPayload),
		},
	}
}
