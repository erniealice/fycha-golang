package asset

// asset.go — asset-domain facade.
//
// Re-exports every entity-local type/const with the entity-prefixed name so
// that consumers (block/, service-admin) keep writing asset.AssetLabels,
// asset.DefaultAssetRoutes(), asset.LapsingScheduleListURL, etc.
//
// Route struct/label types are type aliases pointing at the entity packages;
// route URL consts are re-exported by value. This eliminates the duplicate
// struct/const definitions that previously lived in this package and breaks
// import cycles: entity packages import nothing from domain/asset; only the
// domain/asset facade + the hoisted *_module.go files import entity packages.

import (
	entityasset "github.com/erniealice/fycha-golang/domain/asset/asset"
	assetcategory "github.com/erniealice/fycha-golang/domain/asset/asset_category"
	assetrevaluation "github.com/erniealice/fycha-golang/domain/asset/asset_revaluation"
	depreciationpolicies "github.com/erniealice/fycha-golang/domain/asset/depreciation_policies"
	depreciationrun "github.com/erniealice/fycha-golang/domain/asset/depreciation_run"
	lapsingschedule "github.com/erniealice/fycha-golang/domain/asset/lapsing_schedule"
)

// ---------------------------------------------------------------------------
// Asset (domain/asset/asset)
// ---------------------------------------------------------------------------

type AssetRoutes = entityasset.Routes
type AssetLabels = entityasset.Labels

func DefaultAssetRoutes() AssetRoutes { return entityasset.DefaultRoutes() }
func DefaultAssetLabels() AssetLabels { return entityasset.DefaultLabels() }

// Asset route URL consts (re-exported by value).
const (
	AssetLapsingScheduleURL = entityasset.LapsingScheduleURL
	DepreciationPoliciesURL = entityasset.DepreciationPoliciesURL
)

// ---------------------------------------------------------------------------
// LapsingSchedule (domain/asset/lapsing_schedule)
// ---------------------------------------------------------------------------

type LapsingScheduleRoutes = lapsingschedule.Routes

func DefaultLapsingScheduleRoutes() LapsingScheduleRoutes { return lapsingschedule.DefaultRoutes() }

const LapsingScheduleListURL = lapsingschedule.ListURL

// ---------------------------------------------------------------------------
// DepreciationRun (domain/asset/depreciation_run)
// ---------------------------------------------------------------------------

type DepreciationRunRoutes = depreciationrun.Routes
type DepreciationRunLabels = depreciationrun.Labels

func DefaultDepreciationRunRoutes() DepreciationRunRoutes { return depreciationrun.DefaultRoutes() }
func DefaultDepreciationRunLabels() DepreciationRunLabels { return depreciationrun.DefaultLabels() }

// ---------------------------------------------------------------------------
// AssetCategoryDepreciation (domain/asset/asset_category)
// ---------------------------------------------------------------------------

type AssetCategoryDepreciationRoutes = assetcategory.Routes

func DefaultAssetCategoryDepreciationRoutes() AssetCategoryDepreciationRoutes {
	return assetcategory.DefaultRoutes()
}

const (
	AssetCategoryDepreciationRunURL   = assetcategory.DepreciationRunURL
	AssetPolicyDepreciationRunURL     = assetcategory.PolicyRunURL
	AssetPolicyDepreciationPreviewURL = assetcategory.PolicyPreviewURL
)

// ---------------------------------------------------------------------------
// AssetRevaluation (domain/asset/asset_revaluation)
// ---------------------------------------------------------------------------

type AssetRevaluationLabels = assetrevaluation.Labels

func DefaultAssetRevaluationLabels() AssetRevaluationLabels { return assetrevaluation.DefaultLabels() }

// ---------------------------------------------------------------------------
// DepreciationPolicies (domain/asset/depreciation_policies)
// ---------------------------------------------------------------------------

type DepreciationPoliciesLabels = depreciationpolicies.Labels

func DefaultDepreciationPoliciesLabels() DepreciationPoliciesLabels {
	return depreciationpolicies.DefaultLabels()
}
