// Package block — asset domain wiring.
//
// Holds wireAssetModule (the lifted body of the `if cfg.wantAsset()` branch
// of Block()) plus the proto<->record translators that only the asset wiring
// calls. Co-locating translator + caller keeps the row-shape converters on
// the next page when a reader is following the wiring code.
package block

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strings"

	topref "github.com/erniealice/espyna-golang/reference"

	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"

	pyeza "github.com/erniealice/pyeza-golang"
	pyezatypes "github.com/erniealice/pyeza-golang/types"

	fycha "github.com/erniealice/fycha-golang"
	assetmod "github.com/erniealice/fycha-golang/views/asset"
	assetaction "github.com/erniealice/fycha-golang/views/asset/action"
	assetform "github.com/erniealice/fycha-golang/views/asset/form"
	assetlist "github.com/erniealice/fycha-golang/views/asset/list"
	assetcataction "github.com/erniealice/fycha-golang/views/asset_category/action"
	assetcatpolicies "github.com/erniealice/fycha-golang/views/asset_category/policies"
	depreciationrunmod "github.com/erniealice/fycha-golang/views/depreciation_run"
	lapsinglist "github.com/erniealice/fycha-golang/views/lapsing_schedule/list"
)

// assetWiring holds everything wireAssetModule needs from the surrounding
// Block() scope. Kept private; never re-exported.
//
// More than 6 fields → use a struct (per the convention in block.go's header
// doc). Asset has 14 dependencies, so a struct is unavoidable here.
type assetWiring struct {
	assetRoutes                     fycha.AssetRoutes
	lapsingScheduleRoutes           fycha.LapsingScheduleRoutes
	depreciationRunRoutes           fycha.DepreciationRunRoutes
	assetCategoryDepreciationRoutes fycha.AssetCategoryDepreciationRoutes
	assetLabels                     fycha.AssetLabels
	depreciationRunLabels           fycha.DepreciationRunLabels
	depreciationPoliciesLabels      fycha.DepreciationPoliciesLabels
	assetRevaluationLabels          fycha.AssetRevaluationLabels
	fychaTableLabels                pyezatypes.TableLabels
	common                          pyeza.CommonLabels
	refChecker                      topref.Checker
	// Attachments
	newAttachmentID  func() string
	uploadFile       func(ctx context.Context, name string, mime string, body []byte, refType string) error
	listAttachments  func(ctx context.Context, refType, refID string) (*attachmentpb.ListAttachmentsResponse, error)
	createAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	deleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
}

// wireAssetModule lifts the body of `if cfg.wantAsset()` from Block().
// Behaviour-preserving: same construction order, same registration order,
// same callbacks. block.go calls this exactly once when cfg.wantAsset().
func wireAssetModule(
	ctx *pyeza.AppContext, // route registrar + translation provider live on this
	cfg *blockConfig,
	useCases *UseCases,
	w assetWiring,
) {
	assetDeps := &assetmod.ModuleDeps{
		Routes:       w.assetRoutes,
		CommonLabels: w.common,
		Labels:       w.assetLabels,
		TableLabels:  ctx.Table,
		// Depreciation Run + Revaluation labels (Surface A / E)
		DepreciationRunLabels:  w.depreciationRunLabels,
		AssetRevaluationLabels: w.assetRevaluationLabels,
		// Depreciation run routes + toast URL (Surface A)
		DepreciationRunRoutes:   w.depreciationRunRoutes,
		AssetDepreciationRunURL: cfg.assetDepreciationRunURL,
		// Attachments
		UploadFile:       w.uploadFile,
		ListAttachments:  w.listAttachments,
		CreateAttachment: w.createAttachment,
		DeleteAttachment: w.deleteAttachment,
	}

	// H5 soft-delete gate: wire asset transaction reference checker.
	// The reference checker is provided by the consumer app via ctx.RefChecker
	// (a postgres implementation in prod; NoOp in tests).
	if w.refChecker != nil {
		assetDeps.GetAssetInUseIDs = w.refChecker.GetAssetInUseIDs
	}

	// Typed asset stack (asset-stack buildout, 2026-05-03). Falls back to
	// nothing-wired if the asset use cases are unavailable (mock build, etc.) —
	// same graceful-degradation semantics as ledger/treasury.
	if useCases != nil && useCases.Asset.Create != nil {
		assetDeps.NewID = func() string {
			if w.newAttachmentID != nil {
				return w.newAttachmentID()
			}
			return "" // CreateAsset use case generates IDs internally via IDService
		}
		assetDeps.CreateAsset = func(fctx context.Context, a *assetform.Record) error {
			_, err := useCases.Asset.Create(fctx, &assetpb.CreateAssetRequest{Data: recordToAsset(a)})
			return err
		}
		assetDeps.ReadAsset = func(fctx context.Context, id string) (*assetform.Record, error) {
			resp, err := useCases.Asset.Read(fctx, &assetpb.ReadAssetRequest{Data: &assetpb.Asset{Id: id}})
			if err != nil {
				return nil, err
			}
			if resp == nil || len(resp.Data) == 0 {
				return nil, fmt.Errorf("asset %s not found", id)
			}
			return assetToRecord(resp.Data[0]), nil
		}
		assetDeps.UpdateAsset = func(fctx context.Context, a *assetform.Record) error {
			_, err := useCases.Asset.Update(fctx, &assetpb.UpdateAssetRequest{Data: recordToAsset(a)})
			return err
		}
		// DeleteAsset preserves the legacy soft-delete (active=false) semantic via
		// SetAssetActive — routes the change through audit/auth instead of bypass.
		assetDeps.DeleteAsset = func(fctx context.Context, id string) error {
			_, err := useCases.Asset.SetActive(fctx, &assetpb.SetAssetActiveRequest{AssetId: id, Active: false})
			return err
		}
		assetDeps.SetActive = func(fctx context.Context, id string, active bool) error {
			_, err := useCases.Asset.SetActive(fctx, &assetpb.SetAssetActiveRequest{AssetId: id, Active: active})
			return err
		}
		assetDeps.ListAssets = func(fctx context.Context, status string) ([]assetlist.AssetRow, error) {
			active := status == "active"
			resp, err := useCases.Asset.GetListPageData(fctx, &assetpb.GetAssetListPageDataRequest{
				Filters: &commonpb.FilterRequest{
					Filters: []*commonpb.TypedFilter{
						{
							Field: "active",
							FilterType: &commonpb.TypedFilter_BooleanFilter{
								BooleanFilter: &commonpb.BooleanFilter{Value: active},
							},
						},
					},
				},
			})
			if err != nil {
				return nil, err
			}
			rows := make([]assetlist.AssetRow, 0, len(resp.GetAssetList()))
			for _, a := range resp.GetAssetList() {
				rows = append(rows, assetToRow(a))
			}
			return rows, nil
		}
	}

	// Wire Surface A (depreciation-run) use cases.
	if useCases != nil && useCases.DepRun.ListCandidates != nil {
		assetDeps.ListDepreciationCandidates = func(fctx context.Context, assetID, asOfDate string) ([]assetaction.DepreciationCandidate, error) {
			return listDepreciationCandidatesForAsset(fctx, useCases, assetID, asOfDate)
		}
	}
	if useCases != nil && useCases.DepRun.Generate != nil {
		assetDeps.GenerateDepreciationRun = func(fctx context.Context, req assetaction.DepreciationRunRequest) (*assetaction.DepreciationRunResult, error) {
			return generateDepreciationRunForAsset(fctx, useCases, req)
		}
		// DepreciationFieldsLockedFn: fields are locked once any run has been posted.
		// A dedicated use case (GetAssetDepreciationLock) is pending Wave 3 espyna work.
		// For now we expose the hook so the edit form can call it — nil means unlocked.
	}

	// Wire Surface E (revaluation) use cases — Phase 3 (codex C4).
	// RevalueAsset performs the IAS 16 revaluation atomically (history read +
	// inserts + asset update inside one tx). PreviewRevaluation computes the
	// predicted PnL/OCI split read-only. Both nil-safe when the use case
	// aggregate is missing (mock builds, etc.).
	if useCases != nil && useCases.Revaluation.Revalue != nil {
		assetDeps.RevalueAsset = func(fctx context.Context, req assetaction.RevaluationRequest) (*assetaction.RevaluationResult, error) {
			return revalueAssetCallback(fctx, useCases, req)
		}
	}
	if useCases != nil && useCases.Revaluation.Preview != nil {
		assetDeps.PreviewRevaluation = func(fctx context.Context, assetID string, newFairValue int64) (*assetaction.RevaluationPreview, error) {
			return previewRevaluationCallback(fctx, useCases, assetID, newFairValue)
		}
	}

	assetmod.NewModule(assetDeps).RegisterRoutes(ctx.Routes)

	// ---------------------------------------------------------------------------
	// Surface B — Lapsing Schedule live list page (replaces mock at /app/assets/reports/lapsing-schedule)
	// ---------------------------------------------------------------------------
	lapsingDeps := &lapsinglist.ViewDeps{
		Routes:                w.lapsingScheduleRoutes,
		AssetRoutes:           w.assetRoutes,
		DepreciationRunRoutes: w.depreciationRunRoutes,
		Labels:                w.depreciationRunLabels,
		CommonLabels:          w.common,
		TableLabels:           w.fychaTableLabels,
	}
	// Wire ListCandidates when depreciation use cases are available.
	if useCases != nil && useCases.DepRun.ListCandidates != nil {
		lapsingDeps.ListCandidates = func(fctx context.Context, asOfDate, cursor string, limit int32) ([]lapsinglist.CandidateRow, string, error) {
			return listCandidatesWorkspace(fctx, useCases, asOfDate, cursor, limit)
		}
	}
	// Register Surface B at new URL.
	ctx.Routes.GET(fycha.LapsingScheduleListURL, lapsinglist.NewView(lapsingDeps))
	// Redirect from legacy mock URL to new URL (preserves bookmarks).
	handleFunc(ctx.Routes, "GET", fycha.AssetLapsingScheduleURL, func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, fycha.LapsingScheduleListURL, http.StatusMovedPermanently)
	})

	// ---------------------------------------------------------------------------
	// Surface F — Depreciation Policies actionable page (replaces mock at /app/assets/settings/depreciation-policies)
	// ---------------------------------------------------------------------------
	policiesDeps := &assetcatpolicies.ViewDeps{
		Routes:       w.assetCategoryDepreciationRoutes,
		Labels:       w.depreciationPoliciesLabels,
		CommonLabels: w.common,
		TableLabels:  w.fychaTableLabels,
	}
	// Wire ListPolicies when asset-category use cases are available.
	if useCases != nil && useCases.Asset.Category.ListWithPolicyRollup != nil {
		policiesDeps.ListPolicies = func(fctx context.Context) ([]assetcatpolicies.PolicyRow, error) {
			return listPoliciesWithRollup(fctx, useCases)
		}
	}
	ctx.Routes.GET(fycha.DepreciationPoliciesURL, assetcatpolicies.NewView(policiesDeps))

	// ---------------------------------------------------------------------------
	// Surface F preview drawer (read-only /action/asset-policy/depreciation-preview/{category_id})
	// ---------------------------------------------------------------------------
	previewDeps := &assetcataction.DepreciationPreviewDeps{
		Routes:       w.assetCategoryDepreciationRoutes,
		Labels:       w.depreciationRunLabels,
		CommonLabels: w.common,
		TableLabels:  w.fychaTableLabels,
	}
	if useCases != nil && useCases.DepRun.ListCandidates != nil {
		previewDeps.ListPolicyCandidates = func(fctx context.Context, categoryID, asOfDate string) ([]assetcataction.PreviewCandidateRow, error) {
			return listCandidatesForPolicy(fctx, useCases, categoryID, asOfDate)
		}
	}
	ctx.Routes.GET(fycha.AssetPolicyDepreciationPreviewURL, assetcataction.NewDepreciationPreviewView(previewDeps))

	// ---------------------------------------------------------------------------
	// Surface C — per-category / per-policy depreciation-run drawer
	// Both URLs use the same handler; scope kind is inferred from the URL path.
	// ---------------------------------------------------------------------------
	categoryRunDeps := &assetcataction.CategoryDepreciationRunDeps{
		Routes:       w.assetCategoryDepreciationRoutes,
		RunRoutes:    w.depreciationRunRoutes,
		Labels:       w.depreciationRunLabels,
		CommonLabels: w.common,
	}
	if useCases != nil && useCases.DepRun.ListCandidates != nil {
		categoryRunDeps.ListCategoryCandidates = func(fctx context.Context, categoryID, scopeKind, asOfDate string) ([]assetcataction.CategoryDepreciationRunAssetRow, error) {
			return listCandidatesForCategory(fctx, useCases, categoryID, scopeKind, asOfDate)
		}
	}
	if useCases != nil && useCases.DepRun.Generate != nil {
		categoryRunDeps.GenerateCategoryRun = func(fctx context.Context, req assetcataction.CategoryDepreciationRunRequest) (*assetcataction.CategoryDepreciationRunResult, error) {
			return generateDepreciationRunForCategory(fctx, useCases, req)
		}
	}
	categoryRunView := assetcataction.NewCategoryDepreciationRunAction(categoryRunDeps)
	ctx.Routes.GET(fycha.AssetCategoryDepreciationRunURL, categoryRunView)
	ctx.Routes.POST(fycha.AssetCategoryDepreciationRunURL, categoryRunView)
	ctx.Routes.GET(fycha.AssetPolicyDepreciationRunURL, categoryRunView)
	ctx.Routes.POST(fycha.AssetPolicyDepreciationRunURL, categoryRunView)

	// ---------------------------------------------------------------------------
	// Surface D — Depreciation Runs history list + detail module
	// ---------------------------------------------------------------------------
	drDeps := &depreciationrunmod.ModuleDeps{
		Routes:       w.depreciationRunRoutes,
		Labels:       w.depreciationRunLabels,
		CommonLabels: w.common,
		TableLabels:  w.fychaTableLabels,
	}
	if useCases != nil && useCases.DepRun.List != nil {
		drDeps.ListDepreciationRuns = func(fctx context.Context, scope depreciationrunmod.ListDepreciationRunsScope) ([]depreciationrunmod.DepreciationRunRow, string, error) {
			return listDepreciationRunsForWorkspace(fctx, useCases, scope)
		}
	}
	if useCases != nil && useCases.DepRun.Read != nil {
		drDeps.ReadDepreciationRun = func(fctx context.Context, id string) (*depreciationrunmod.DepreciationRunWithEntries, error) {
			return readDepreciationRunWithEntries(fctx, useCases, id)
		}
		// TODO: add ListAssetTransactionsByRunID use case (Followup — Phase 4 codex H2).
		// When ListAssetTransactionsByRunID is available in espyna-golang, wire it here.
		// Until then, the transactions tab renders an empty table (nil guard in detail/page.go:loadTabData).
	}
	depreciationrunmod.NewModule(drDeps).RegisterRoutes(ctx.Routes)
}

// ---------------------------------------------------------------------------
// Asset type-translation helpers (asset-stack buildout, 2026-05-03)
// ---------------------------------------------------------------------------

// recordToAsset converts a view-layer assetform.Record to the proto Asset type.
// Money fields are translated from float64 pesos → int64 centavos.
// Enum fields are translated from string → proto enum using the generated _value maps.
// Unknown enum strings map to 0 (*_UNSPECIFIED), preserving current behaviour.
func recordToAsset(r *assetform.Record) *assetpb.Asset {
	a := &assetpb.Asset{
		Id:                 r.ID,
		AssetNumber:        r.AssetNumber,
		Name:               r.Name,
		AssetType:          assetpb.AssetType(assetpb.AssetType_value[r.AssetType]),
		AssetCategoryId:    r.AssetCategoryID,
		AcquisitionCost:    int64(math.Round(r.AcquisitionCost * 100)),
		Currency:           r.Currency,
		SalvageValue:       int64(math.Round(r.SalvageValue * 100)),
		BookValue:          int64(math.Round(r.BookValue * 100)),
		UsefulLifeMonths:   int32(r.UsefulLifeMonths),
		DepreciationMethod: assetpb.DepreciationMethod(assetpb.DepreciationMethod_value[r.DepreciationMethod]),
		Status:             assetpb.AssetStatus(assetpb.AssetStatus_value[r.Status]),
		Active:             r.Active,
	}
	if r.Description != "" {
		a.Description = &r.Description
	}
	if r.LocationID != "" {
		a.LocationId = &r.LocationID
	}
	return a
}

// assetToRecord converts a proto Asset back to the view-layer assetform.Record.
// Money fields are translated from int64 centavos → float64 pesos.
// Enum strings are lowercased and stripped of their proto prefix so they round-trip
// to the form values the view layer expects (e.g. "DEPRECIATION_METHOD_STRAIGHT_LINE"
// → "straight_line").
func assetToRecord(a *assetpb.Asset) *assetform.Record {
	r := &assetform.Record{
		ID:                 a.GetId(),
		AssetNumber:        a.GetAssetNumber(),
		Name:               a.GetName(),
		AssetType:          strings.ToLower(strings.TrimPrefix(a.GetAssetType().String(), "ASSET_TYPE_")),
		AssetCategoryID:    a.GetAssetCategoryId(),
		LocationID:         a.GetLocationId(),
		AcquisitionCost:    float64(a.GetAcquisitionCost()) / 100,
		Currency:           a.GetCurrency(),
		SalvageValue:       float64(a.GetSalvageValue()) / 100,
		BookValue:          float64(a.GetBookValue()) / 100,
		UsefulLifeMonths:   int(a.GetUsefulLifeMonths()),
		DepreciationMethod: strings.ToLower(strings.TrimPrefix(a.GetDepreciationMethod().String(), "DEPRECIATION_METHOD_")),
		Status:             strings.ToLower(strings.TrimPrefix(a.GetStatus().String(), "ASSET_STATUS_")),
		Active:             a.GetActive(),
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}
