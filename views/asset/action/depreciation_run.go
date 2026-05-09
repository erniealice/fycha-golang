package action

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Surface A — per-asset depreciation-run drawer
// URL: /action/asset/depreciation-run/{asset_id}
// ---------------------------------------------------------------------------

// DepreciationRunDeps holds dependencies for the Surface A drawer.
type DepreciationRunDeps struct {
	Routes                fycha.AssetRoutes
	DepreciationRunRoutes fycha.DepreciationRunRoutes
	Labels                fycha.DepreciationRunLabels
	// ListDepreciationCandidates calls the espyna dry-run engine.
	// Signature intentionally wide ([]byte) so we don't import espyna protos here.
	ListDepreciationCandidates func(ctx context.Context, assetID, asOfDate string) ([]DepreciationCandidate, error)
	// GenerateDepreciationRun posts selected periods via espyna.
	GenerateDepreciationRun func(ctx context.Context, req DepreciationRunRequest) (*DepreciationRunResult, error)
	// AssetDepreciationRunURL is the resolved run-detail URL template (e.g. for toast link).
	// Injected by block.go via WithAssetDepreciationRunURL.
	AssetDepreciationRunURL string
}

// DepreciationCandidate is one pending period for a single asset.
type DepreciationCandidate struct {
	PeriodStart        string // YYYY-MM-DD
	PeriodEnd          string // YYYY-MM-DD
	PeriodLabel        string // e.g. "Jan 2025"
	ProjectedAmount    int64  // centavos
	ProjectedAmountFmt string // pre-formatted for display
	ProjectedAccum     int64  // centavos — running accumulated after this period
	ProjectedAccumFmt  string
	Blocked            bool
	BlockerReason      string
}

// DepreciationRunRequest is the POST payload to GenerateDepreciationRun.
type DepreciationRunRequest struct {
	AssetID          string
	AsOfDate         string
	PeriodStartDates []string // selected period_start_dates
}

// DepreciationRunResult is the response from GenerateDepreciationRun.
type DepreciationRunResult struct {
	RunID        string
	CreatedCount int
	SkippedCount int
	ErroredCount int
	Success      bool
	ErrorMessage string
}

// depreciationRunFormData is the template data for the Surface A drawer.
type depreciationRunFormData struct {
	FormAction   string
	AssetID      string
	AsOfDate     string
	MaxAsOfDate  string
	FragmentURL  string // HTMX inner-swap target URL
	Candidates   []DepreciationCandidate
	EligibleCount int
	Labels       fycha.DepreciationRunLabels
}

// NewDepreciationRunAction creates the Surface A per-asset depreciation-run drawer.
func NewDepreciationRunAction(deps *DepreciationRunDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		assetID := viewCtx.Request.PathValue("asset_id")
		if assetID == "" {
			assetID = viewCtx.Request.PathValue("id")
		}

		if viewCtx.Request.Method == http.MethodGet {
			return handleDepreciationRunGET(ctx, viewCtx, deps, assetID)
		}
		return handleDepreciationRunPOST(ctx, viewCtx, deps, assetID)
	})
}

func handleDepreciationRunGET(ctx context.Context, viewCtx *view.ViewContext, deps *DepreciationRunDeps, assetID string) view.ViewResult {
	asOfDate := viewCtx.Request.URL.Query().Get("as_of_date")
	if asOfDate == "" {
		asOfDate = time.Now().Format("2006-01-02")
	}

	formURL := deps.Routes.DepreciationRunFor(assetID)

	data := &depreciationRunFormData{
		FormAction:  formURL,
		AssetID:     assetID,
		AsOfDate:    asOfDate,
		MaxAsOfDate: time.Now().Format("2006-01-02"),
		FragmentURL: formURL + "?fragment=periods&as_of_date=" + asOfDate,
		Labels:      deps.Labels,
	}

	if deps.ListDepreciationCandidates != nil {
		candidates, err := deps.ListDepreciationCandidates(ctx, assetID, asOfDate)
		if err != nil {
			log.Printf("ListDepreciationCandidates error: %v", err)
		} else {
			data.Candidates = candidates
			for _, c := range candidates {
				if !c.Blocked {
					data.EligibleCount++
				}
			}
		}
	}

	return view.OK("asset-depreciation-run-drawer-form", data)
}

func handleDepreciationRunPOST(ctx context.Context, viewCtx *view.ViewContext, deps *DepreciationRunDeps, assetID string) view.ViewResult {
	if err := viewCtx.Request.ParseForm(); err != nil {
		return fycha.HTMXError(deps.Labels.Errors.InvalidSelection)
	}

	asOfDate := viewCtx.Request.FormValue("as_of_date")
	if asOfDate == "" {
		asOfDate = time.Now().Format("2006-01-02")
	}
	selections := viewCtx.Request.Form["selection"]

	var result *DepreciationRunResult
	if deps.GenerateDepreciationRun != nil {
		var err error
		result, err = deps.GenerateDepreciationRun(ctx, DepreciationRunRequest{
			AssetID:          assetID,
			AsOfDate:         asOfDate,
			PeriodStartDates: selections,
		})
		if err != nil {
			log.Printf("GenerateDepreciationRun error: %v", err)
			return fycha.HTMXError(deps.Labels.Errors.UseCaseUnavailable)
		}
	} else {
		// Graceful mock result when use case not yet wired
		result = &DepreciationRunResult{
			RunID:        "run-mock",
			CreatedCount: len(selections),
			Success:      true,
		}
	}

	// Build toast payload per pyeza:toast contract
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

	runDetailURL := deps.AssetDepreciationRunURL
	if deps.DepreciationRunRoutes.DetailURL != "" && result.RunID != "" {
		runDetailURL = deps.DepreciationRunRoutes.DetailFor(result.RunID)
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
		"refreshTable": "asset-detail-lapsing-actual-schedule",
	})

	return view.ViewResult{
		StatusCode: 200,
		Headers: map[string]string{
			"HX-Trigger": string(triggerPayload),
		},
	}
}
