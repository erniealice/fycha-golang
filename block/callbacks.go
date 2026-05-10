// Package block — lapsing-schedule, policy, and revaluation callback helpers.
//
// This file holds the 12 lapsing-schedule helpers called by wireAssetModule's
// callback closures (listCandidatesWorkspace, listPoliciesWithRollup,
// listCandidatesForPolicy, listDepreciationCandidatesForAsset,
// generateDepreciationRunForAsset, assetToRow, listCandidatesForCategory,
// generateDepreciationRunForCategory, listDepreciationRunsForWorkspace,
// readDepreciationRunWithEntries, depreciationScheduleToEntryRow,
// depreciationRunToRow) plus 2 Surface E revaluation helpers
// (revalueAssetCallback, previewRevaluationCallback).
//
// Rule of thumb: functions here are called from asset.go's wireAssetModule
// closures; they translate proto request/response shapes into view-layer types.
package block

import (
	"context"
	"fmt"
	"strings"

	consumer "github.com/erniealice/espyna-golang/consumer"

	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	deprunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation_run"
	depschpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"

	"github.com/erniealice/pyeza-golang/types"

	assetaction "github.com/erniealice/fycha-golang/views/asset/action"
	assetlist "github.com/erniealice/fycha-golang/views/asset/list"
	assetcataction "github.com/erniealice/fycha-golang/views/asset_category/action"
	assetcatpolicies "github.com/erniealice/fycha-golang/views/asset_category/policies"
	depreciationrunmod "github.com/erniealice/fycha-golang/views/depreciation_run"
	lapsinglist "github.com/erniealice/fycha-golang/views/lapsing_schedule/list"
)

// ---------------------------------------------------------------------------
// Lapsing schedule + policy helpers
// ---------------------------------------------------------------------------

// listCandidatesWorkspace calls ListDepreciationCandidates with scope_kind=WORKSPACE
// and maps the proto response to view-layer CandidateRow slices.
func listCandidatesWorkspace(
	ctx context.Context,
	useCases *consumer.UseCases,
	asOfDate, cursor string,
	limit int32,
) ([]lapsinglist.CandidateRow, string, error) {
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_WORKSPACE,
		AsOfDate:  asOfDate,
		Pagination: &commonpb.PaginationRequest{
			Limit: limit,
			Method: &commonpb.PaginationRequest_Cursor{
				Cursor: &commonpb.CursorPagination{Token: cursor},
			},
		},
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, "", err
	}
	rows := make([]lapsinglist.CandidateRow, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		row := lapsinglist.CandidateRow{
			AssetID:          c.GetAssetId(),
			AssetName:        c.GetAssetName(),
			Currency:         c.GetCurrency(),
			CurrentBookValue: c.GetCurrentBookValue(),
			CanRun:           len(c.GetBlockers()) == 0 && len(c.GetPeriods()) > 0,
		}
		if len(c.GetPeriods()) > 0 {
			row.PendingCount = len(c.GetPeriods())
			row.NextAmount = c.GetPeriods()[0].GetAmount()
			row.NextPendingPeriod = c.GetPeriods()[0].GetPeriodStartDate()
		}
		if len(c.GetBlockers()) > 0 {
			row.Status = "blocked"
			row.BlockerLabel = c.GetBlockers()[0].GetLabel()
		} else if row.PendingCount == 0 {
			row.Status = "up_to_date"
		} else if row.LastPostedPeriod == "" {
			row.Status = "not_started"
		} else {
			row.Status = "pending"
		}
		rows = append(rows, row)
	}
	nextCursor := ""
	if resp.GetPagination() != nil {
		nextCursor = resp.GetPagination().GetNextCursor()
	}
	return rows, nextCursor, nil
}

// listPoliciesWithRollup fetches all AssetCategory rows with per-category aggregate counts.
// Uses the new ListAssetCategoriesWithPolicyRollup use case (Wave 3 espyna enhancement).
// AssetsInPolicy and AssetsDeviating are real counts from the Postgres bulk query.
func listPoliciesWithRollup(
	ctx context.Context,
	useCases *consumer.UseCases,
) ([]assetcatpolicies.PolicyRow, error) {
	rollupRows, err := consumer.ListAssetCategoriesWithPolicyRollup(useCases, ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]assetcatpolicies.PolicyRow, 0, len(rollupRows))
	for _, r := range rollupRows {
		c := r.Category
		if c == nil {
			continue
		}
		method := c.GetDefaultDepreciationMethod()
		if c.DepreciationMethod != nil {
			method = c.GetDepreciationMethod()
		}
		usefulLife := c.GetDefaultUsefulLifeMonths()
		if c.UsefulLifeMonths != nil {
			usefulLife = c.GetUsefulLifeMonths()
		}
		salvage := c.GetDefaultSalvageValuePercent()
		if c.SalvagePct != nil {
			salvage = c.GetSalvagePct()
		}
		rows = append(rows, assetcatpolicies.PolicyRow{
			CategoryID:         c.GetId(),
			PolicyID:           c.GetId(),
			Name:               c.GetName(),
			DepreciationMethod: method,
			UsefulLifeMonths:   usefulLife,
			SalvagePct:         salvage,
			AssetsInPolicy:     r.AssetsInPolicy,
			AssetsDeviating:    r.AssetsDeviating,
		})
	}
	return rows, nil
}

// listCandidatesForPolicy calls ListDepreciationCandidates with scope_kind=POLICY
// and maps results to the preview drawer's PreviewCandidateRow.
func listCandidatesForPolicy(
	ctx context.Context,
	useCases *consumer.UseCases,
	categoryID, asOfDate string,
) ([]assetcataction.PreviewCandidateRow, error) {
	scopeID := categoryID
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_POLICY,
		ScopeId:   &scopeID,
		AsOfDate:  asOfDate,
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, err
	}
	rows := make([]assetcataction.PreviewCandidateRow, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		row := assetcataction.PreviewCandidateRow{
			AssetID:      c.GetAssetId(),
			AssetName:    c.GetAssetName(),
			Currency:     c.GetCurrency(),
			BookValue:    c.GetCurrentBookValue(),
			PendingCount: len(c.GetPeriods()),
		}
		if len(c.GetPeriods()) > 0 {
			row.NextAmount = c.GetPeriods()[0].GetAmount()
			row.NextAmountFmt = types.FormatMoney(row.NextAmount, row.Currency)
		}
		for _, b := range c.GetBlockers() {
			row.Blockers = append(row.Blockers, b.GetLabel())
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// ---------------------------------------------------------------------------
// Surface A helpers — per-asset depreciation-run wrappers
// ---------------------------------------------------------------------------

// listDepreciationCandidatesForAsset calls ListDepreciationCandidates with
// scope_kind=ASSET and maps the proto response to the view-layer DepreciationCandidate slice.
func listDepreciationCandidatesForAsset(
	ctx context.Context,
	useCases *consumer.UseCases,
	assetID, asOfDate string,
) ([]assetaction.DepreciationCandidate, error) {
	if asOfDate == "" {
		asOfDate = "today" // espyna engine accepts "today" as a sentinel
	}
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_ASSET,
		ScopeId:   &assetID,
		AsOfDate:  asOfDate,
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, err
	}
	var rows []assetaction.DepreciationCandidate
	for _, c := range resp.GetData() {
		assetCurrency := c.GetCurrency()
		for _, p := range c.GetPeriods() {
			rows = append(rows, assetaction.DepreciationCandidate{
				PeriodStart:        p.GetPeriodStartDate(),
				PeriodEnd:          p.GetPeriodEndDate(),
				ProjectedAmount:    p.GetAmount(),
				ProjectedAmountFmt: types.FormatMoney(p.GetAmount(), assetCurrency),
				ProjectedAccum:     p.GetRunningAccumulated(),
				ProjectedAccumFmt:  types.FormatMoney(p.GetRunningAccumulated(), assetCurrency),
			})
		}
		for _, b := range c.GetBlockers() {
			blockerKind := ""
			if b.GetKind() == deprunpb.DepreciationCandidateBlocker_DEPRECIATION_CANDIDATE_BLOCKER_KIND_UNITS_REQUIRED {
				blockerKind = "UNITS_REQUIRED"
			}
			rows = append(rows, assetaction.DepreciationCandidate{
				Blocked:       true,
				BlockerReason: b.GetLabel(),
				BlockerKind:   blockerKind,
			})
		}
	}
	return rows, nil
}

// generateDepreciationRunForAsset posts selected periods for a single asset
// and maps the proto response back to the view-layer DepreciationRunResult.
func generateDepreciationRunForAsset(
	ctx context.Context,
	useCases *consumer.UseCases,
	req assetaction.DepreciationRunRequest,
) (*assetaction.DepreciationRunResult, error) {
	protoReq := &deprunpb.GenerateDepreciationRunRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_ASSET,
		ScopeId:   &req.AssetID,
		AsOfDate:  req.AsOfDate,
		Selections: []*deprunpb.DepreciationRunSelection{
			{
				AssetId:          req.AssetID,
				PeriodStartDates: req.PeriodStartDates,
			},
		},
	}
	resp, err := consumer.GenerateDepreciationRun(useCases, ctx, protoReq)
	if err != nil {
		return nil, err
	}
	runID := ""
	if resp.GetRun() != nil {
		runID = resp.GetRun().GetId()
	}
	return &assetaction.DepreciationRunResult{
		RunID:        runID,
		CreatedCount: int(resp.GetCreatedCount()),
		SkippedCount: int(resp.GetSkippedCount()),
		ErroredCount: int(resp.GetErroredCount()),
		Success:      resp.GetSuccess(),
	}, nil
}

// assetToRow converts a proto Asset to the flat assetlist.AssetRow used by the list view.
func assetToRow(a *assetpb.Asset) assetlist.AssetRow {
	row := assetlist.AssetRow{
		ID:              a.GetId(),
		AssetNumber:     a.GetAssetNumber(),
		Name:            a.GetName(),
		AcquisitionCost: float64(a.GetAcquisitionCost()) / 100,
		BookValue:       float64(a.GetBookValue()) / 100,
		Active:          a.GetActive(),
	}
	if a.AssetCategory != nil {
		row.CategoryName = a.AssetCategory.GetName()
	}
	if a.Location != nil {
		row.LocationName = a.Location.GetName()
	}
	return row
}

// ---------------------------------------------------------------------------
// Surface C helpers — per-category / per-policy depreciation-run wrappers
// ---------------------------------------------------------------------------

// listCandidatesForCategory calls ListDepreciationCandidates with scope_kind=CATEGORY or POLICY
// and maps results to the Surface C CategoryDepreciationRunAssetRow slice.
// One row per asset (not per period — the drawer shows which assets to include, not individual periods).
func listCandidatesForCategory(
	ctx context.Context,
	useCases *consumer.UseCases,
	categoryID, scopeKind, asOfDate string,
) ([]assetcataction.CategoryDepreciationRunAssetRow, error) {
	if asOfDate == "" {
		asOfDate = "today"
	}
	proto := deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_CATEGORY
	if scopeKind == "policy" || scopeKind == "POLICY" {
		proto = deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_POLICY
	}
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: proto,
		ScopeId:   &categoryID,
		AsOfDate:  asOfDate,
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, err
	}
	rows := make([]assetcataction.CategoryDepreciationRunAssetRow, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		row := assetcataction.CategoryDepreciationRunAssetRow{
			AssetID:      c.GetAssetId(),
			AssetName:    c.GetAssetName(),
			Currency:     c.GetCurrency(),
			BookValue:    c.GetCurrentBookValue(),
			PendingCount: len(c.GetPeriods()),
		}
		if len(c.GetPeriods()) > 0 {
			row.NextAmount = c.GetPeriods()[0].GetAmount()
			row.NextAmountFmt = types.FormatMoney(row.NextAmount, row.Currency)
		}
		for _, b := range c.GetBlockers() {
			row.Blockers = append(row.Blockers, b.GetLabel())
		}
		row.CanRun = len(row.Blockers) == 0 && row.PendingCount > 0
		rows = append(rows, row)
	}
	return rows, nil
}

// generateDepreciationRunForCategory posts a depreciation run for all selected assets
// in a category/policy scope and maps the result back to the view-layer type.
func generateDepreciationRunForCategory(
	ctx context.Context,
	useCases *consumer.UseCases,
	req assetcataction.CategoryDepreciationRunRequest,
) (*assetcataction.CategoryDepreciationRunResult, error) {
	protoScope := deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_CATEGORY
	if req.ScopeKind == "POLICY" {
		protoScope = deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_POLICY
	}
	// Build per-asset selections: each selected asset contributes all its pending periods.
	// The use case engine resolves pending periods server-side from the as_of_date; we pass
	// an empty period list per asset so the engine computes them (same as the Revenue Run
	// "all-for-scope" pattern). If asset IDs are empty, scope covers ALL assets in the category.
	var selections []*deprunpb.DepreciationRunSelection
	for _, assetID := range req.AssetIDs {
		aid := assetID
		selections = append(selections, &deprunpb.DepreciationRunSelection{
			AssetId: aid,
			// PeriodStartDates empty → use case posts all pending periods for this asset.
		})
	}
	protoReq := &deprunpb.GenerateDepreciationRunRequest{
		ScopeKind:  protoScope,
		ScopeId:    &req.CategoryID,
		AsOfDate:   req.AsOfDate,
		Selections: selections,
	}
	resp, err := consumer.GenerateDepreciationRun(useCases, ctx, protoReq)
	if err != nil {
		return nil, err
	}
	runID := ""
	if resp.GetRun() != nil {
		runID = resp.GetRun().GetId()
	}
	return &assetcataction.CategoryDepreciationRunResult{
		RunID:        runID,
		CreatedCount: int(resp.GetCreatedCount()),
		SkippedCount: int(resp.GetSkippedCount()),
		ErroredCount: int(resp.GetErroredCount()),
		Success:      resp.GetSuccess(),
	}, nil
}

// ---------------------------------------------------------------------------
// Surface D helpers — depreciation-run history wrappers
// ---------------------------------------------------------------------------

// listDepreciationRunsForWorkspace fetches a page of DepreciationRun rows for Surface D.
func listDepreciationRunsForWorkspace(
	ctx context.Context,
	useCases *consumer.UseCases,
	scope depreciationrunmod.ListDepreciationRunsScope,
) ([]depreciationrunmod.DepreciationRunRow, string, error) {
	req := &deprunpb.ListDepreciationRunsRequest{}
	resp, err := consumer.ListDepreciationRuns(useCases, ctx, req)
	if err != nil {
		return nil, "", err
	}
	rows := make([]depreciationrunmod.DepreciationRunRow, 0, len(resp.GetData()))
	for _, r := range resp.GetData() {
		if scope.Status != "" {
			status := strings.ToLower(strings.TrimPrefix(r.GetStatus().String(), "DEPRECIATION_RUN_STATUS_"))
			if status != scope.Status {
				continue
			}
		}
		rows = append(rows, depreciationRunToRow(r))
	}
	return rows, "", nil
}

// readDepreciationRunWithEntries fetches a single DepreciationRun plus schedule entries.
// Phase 4 (codex H2): wires ListDepreciationRunEntries so the detail page
// selections/results tabs render real data instead of empty tables.
func readDepreciationRunWithEntries(
	ctx context.Context,
	useCases *consumer.UseCases,
	id string,
) (*depreciationrunmod.DepreciationRunWithEntries, error) {
	resp, err := consumer.ReadDepreciationRun(useCases, ctx, &deprunpb.ReadDepreciationRunRequest{
		Data: &deprunpb.DepreciationRun{Id: id},
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.GetData()) == 0 {
		return nil, fmt.Errorf("depreciation run %s not found", id)
	}
	run := depreciationRunToRow(resp.GetData()[0])

	// Fetch schedule entries (selections/results tabs).
	entriesResp, err := consumer.ListDepreciationRunEntries(useCases, ctx, id, nil)
	if err != nil {
		return nil, err
	}
	entries := make([]depreciationrunmod.DepreciationRunEntryRow, 0)
	if entriesResp != nil {
		for _, s := range entriesResp.GetData() {
			entries = append(entries, depreciationScheduleToEntryRow(s))
		}
	}

	return &depreciationrunmod.DepreciationRunWithEntries{
		Run:     run,
		Entries: entries,
	}, nil
}

// depreciationScheduleToEntryRow maps a proto DepreciationSchedule (scoped to a run)
// to the view-layer DepreciationRunEntryRow used by the selections and results tabs.
func depreciationScheduleToEntryRow(s *depschpb.DepreciationSchedule) depreciationrunmod.DepreciationRunEntryRow {
	if s == nil {
		return depreciationrunmod.DepreciationRunEntryRow{}
	}
	// outcome is stored as a string matching DepreciationRunOutcome enum values
	// (e.g. "DEPRECIATION_RUN_OUTCOME_CREATED"). Map to lowercase view status
	// ("created" | "skipped" | "errored") for CSS class hooks.
	outcome := strings.ToLower(strings.TrimPrefix(s.GetOutcome(), "DEPRECIATION_RUN_OUTCOME_"))
	return depreciationrunmod.DepreciationRunEntryRow{
		ID:                 s.GetId(),
		RunID:              s.GetDepreciationRunId(),
		AssetID:            s.GetAssetId(),
		PeriodStartDate:    s.GetPeriodStartDate(),
		DepreciationAmount: s.GetDepreciationAmount(),
		Outcome:            outcome,
		ErrorMessage:       s.GetErrorMessage(),
		IsPosted:           s.GetIsPosted(),
	}
}

// depreciationRunToRow maps a proto DepreciationRun to the view-layer DepreciationRunRow.
func depreciationRunToRow(r *deprunpb.DepreciationRun) depreciationrunmod.DepreciationRunRow {
	if r == nil {
		return depreciationrunmod.DepreciationRunRow{}
	}
	status := strings.ToLower(strings.TrimPrefix(r.GetStatus().String(), "DEPRECIATION_RUN_STATUS_"))
	scopeKind := strings.ToLower(strings.TrimPrefix(r.GetScopeKind().String(), "DEPRECIATION_RUN_SCOPE_KIND_"))
	return depreciationrunmod.DepreciationRunRow{
		ID:           r.GetId(),
		WorkspaceID:  r.GetWorkspaceId(),
		ScopeKind:    scopeKind,
		ScopeID:      r.GetScopeId(),
		AsOfDate:     r.GetAsOfDate(),
		InitiatorID:  r.GetInitiatorId(),
		Status:       status,
		CreatedCount: r.GetCreatedCount(),
		SkippedCount: r.GetSkippedCount(),
		ErroredCount: r.GetErroredCount(),
	}
}

// ---------------------------------------------------------------------------
// Surface E helpers — per-asset revaluation wrappers (Phase 3, codex C4)
// ---------------------------------------------------------------------------

// revalueAssetCallback bridges the view-layer RevaluationRequest/Result to the
// espyna RevalueAsset use case. The use case writes the AssetRevaluation +
// AssetTransaction + Asset update atomically inside one tx; this helper just
// translates types and renders the view-friendly toast strings.
func revalueAssetCallback(
	ctx context.Context,
	useCases *consumer.UseCases,
	req assetaction.RevaluationRequest,
) (*assetaction.RevaluationResult, error) {
	ucReq := consumer.RevalueAssetRequest{
		AssetID:         req.AssetID,
		NewFairValue:    req.NewFairValue,
		AppraiserName:   req.AppraiserName,
		ValuationMethod: req.ValuationMethod,
		Notes:           req.Notes,
	}
	ucResult, err := consumer.RevalueAsset(useCases, ctx, ucReq)
	if err != nil {
		return nil, err
	}
	if ucResult == nil || ucResult.Revaluation == nil || ucResult.Transaction == nil {
		return nil, fmt.Errorf("revalue_asset: use case returned an empty result")
	}

	rev := ucResult.Revaluation
	tx := ucResult.Transaction

	direction := "UP"
	if !rev.GetIsIncrease() {
		direction = "DOWN"
	}

	// Recognition picks the dominant side so the toast can summarise where
	// the entry landed. When both PnL and OCI receive a portion (mixed
	// case under IAS 16.39-40), we surface "OCI/P&L" so the user knows it
	// straddles both. The detail page's history tab carries the exact
	// per-side amounts.
	pnlMag := rev.GetRecognizedInPnl()
	if pnlMag < 0 {
		pnlMag = -pnlMag
	}
	ociMag := rev.GetRecognizedInOci()
	if ociMag < 0 {
		ociMag = -ociMag
	}
	recognition := "OCI"
	switch {
	case pnlMag > 0 && ociMag > 0:
		recognition = "OCI/P&L"
	case pnlMag > 0:
		recognition = "P&L"
	case ociMag > 0:
		recognition = "OCI"
	}

	absAmount := tx.GetAmount()
	if absAmount < 0 {
		absAmount = -absAmount
	}

	return &assetaction.RevaluationResult{
		TransactionID: tx.GetId(),
		Direction:     direction,
		Amount:        absAmount,
		AmountFmt:     types.FormatMoney(absAmount, getFunctionalCurrency(ctx, useCases)),
		Recognition:   recognition,
		Success:       true,
	}, nil
}

// previewRevaluationCallback bridges the view-layer preview signature to the
// espyna PreviewRevaluation use case. Read-only; no DB writes.
func previewRevaluationCallback(
	ctx context.Context,
	useCases *consumer.UseCases,
	assetID string,
	newFairValue int64,
) (*assetaction.RevaluationPreview, error) {
	ucResult, err := consumer.PreviewRevaluation(useCases, ctx, consumer.PreviewRevaluationRequest{
		AssetID:      assetID,
		NewFairValue: newFairValue,
	})
	if err != nil {
		return nil, err
	}
	if ucResult == nil {
		return nil, nil
	}

	direction := "UP"
	if !ucResult.IsIncrease {
		direction = "DOWN"
	}

	revAmount := ucResult.RevaluationAmount
	if revAmount < 0 {
		revAmount = -revAmount
	}
	pnlMag := ucResult.RecognizedInPnL
	if pnlMag < 0 {
		pnlMag = -pnlMag
	}
	ociMag := ucResult.RecognizedInOCI
	if ociMag < 0 {
		ociMag = -ociMag
	}

	currency := getFunctionalCurrency(ctx, useCases)
	return &assetaction.RevaluationPreview{
		RevaluationAmountFmt: types.FormatMoney(revAmount, currency),
		PnLAmountFmt:         types.FormatMoney(pnlMag, currency),
		OCIAmountFmt:         types.FormatMoney(ociMag, currency),
		Direction:            direction,
	}, nil
}
