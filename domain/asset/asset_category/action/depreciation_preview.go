// Package action implements drawer actions for the asset-category domain.
//
// depreciation_preview.go — Surface F read-only preview drawer.
// URL: /action/asset-policy/depreciation-preview/{category_id}
// GET only — no writes. Calls ListDepreciationCandidates(scope_kind=POLICY)
// and renders per-asset projected amounts + blockers.
package action

import (
	"context"
	"fmt"
	"log"
	"time"

	assetcategory "github.com/erniealice/fycha-golang/domain/asset/asset_category"
	depreciationrun "github.com/erniealice/fycha-golang/domain/asset/depreciation_run"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// PreviewCandidateRow holds per-asset preview data for the read-only drawer.
type PreviewCandidateRow struct {
	AssetID       string
	AssetName     string
	Currency      string
	BookValue     int64 // centavos — current
	PendingCount  int
	NextAmount    int64  // centavos — next period projected
	NextAmountFmt string // pre-formatted via types.FormatMoney (e.g. "PHP 50,000.00")
	Blockers      []string
}

// DepreciationPreviewDeps holds dependencies for the read-only preview drawer.
type DepreciationPreviewDeps struct {
	Routes       assetcategory.Routes
	Labels       depreciationrun.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// ListPolicyCandidates returns candidates for the given category (policy scope),
	// dry-run only — no writes. Block.go wires this to espyna.
	ListPolicyCandidates func(ctx context.Context, categoryID, asOfDate string) ([]PreviewCandidateRow, error)
}

// DepreciationPreviewPageData is the template context for the preview drawer.
type DepreciationPreviewPageData struct {
	CategoryID   string
	AsOfDate     string
	Rows         []PreviewCandidateRow
	Labels       depreciationrun.Labels
	CommonLabels pyeza.CommonLabels
}

// NewDepreciationPreviewView creates the read-only policy preview drawer.
// GET only — no POST route is registered for this handler.
func NewDepreciationPreviewView(deps *DepreciationPreviewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		// 2026-05-14 permission-gates: preview is read-only, gate on
		// depreciation_schedule:read.
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("depreciation_schedule", "read") {
			return view.HTMXError(deps.Labels.Errors.PermissionDenied)
		}

		categoryID := viewCtx.Request.PathValue("category_id")
		if categoryID == "" {
			return view.Error(fmt.Errorf("category_id is required"))
		}

		asOfDate := time.Now().UTC().Format("2006-01-02")

		var rows []PreviewCandidateRow
		if deps.ListPolicyCandidates != nil {
			var err error
			rows, err = deps.ListPolicyCandidates(ctx, categoryID, asOfDate)
			if err != nil {
				log.Printf("depreciation-preview: ListPolicyCandidates error for category %s: %v", categoryID, err)
				return view.Error(fmt.Errorf("failed to load preview data: %w", err))
			}
		}
		if rows == nil {
			rows = []PreviewCandidateRow{}
		}

		pageData := &DepreciationPreviewPageData{
			CategoryID:   categoryID,
			AsOfDate:     asOfDate,
			Rows:         rows,
			Labels:       deps.Labels,
			CommonLabels: deps.CommonLabels,
		}

		return view.OK("depreciation-policy-preview", pageData)
	})
}
