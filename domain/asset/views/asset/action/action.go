package action

import (
	"context"

	asset "github.com/erniealice/fycha-golang/domain/asset"
	assetform "github.com/erniealice/fycha-golang/domain/asset/views/asset/form"
)

// Deps holds dependencies for asset action handlers.
type Deps struct {
	Routes asset.AssetRoutes
	Labels asset.AssetLabels

	// CRUD operations (wired from block.go)
	CreateAsset func(ctx context.Context, asset *assetform.Record) error
	ReadAsset   func(ctx context.Context, id string) (*assetform.Record, error)
	UpdateAsset func(ctx context.Context, asset *assetform.Record) error
	DeleteAsset func(ctx context.Context, id string) error
	SetActive   func(ctx context.Context, id string, active bool) error
	NewID       func() string

	// DepreciationFieldsLockedFn returns true when a posted depreciation_schedule
	// row exists for the asset, triggering locked-field rendering in the edit drawer.
	// Nil = always unlocked (mock build or use cases not yet wired).
	DepreciationFieldsLockedFn func(ctx context.Context, assetID string) (bool, error)

	// GetAssetInUseIDs returns a map of asset IDs that have any asset_transaction
	// row. Used by the delete handler to enforce the H5 server-side soft-delete
	// gate before calling DeleteAsset.
	// Nil = skip the check (mock build or use cases not yet wired).
	GetAssetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)
}

// labelsFromDeps builds form.Labels from the deps label tree. The inline
// construction at each call site was identical; this helper eliminates the
// repetition without introducing a mapper that transforms no values.
func labelsFromDeps(deps *Deps) assetform.Labels {
	return assetform.Labels{
		Name:                             deps.Labels.Form.Name,
		NamePlaceholder:                  deps.Labels.Form.NamePlaceholder,
		AssetNumber:                      deps.Labels.Form.AssetNumber,
		AssetNumberPlaceholder:           deps.Labels.Form.AssetNumberPlaceholder,
		Description:                      deps.Labels.Form.Description,
		DescriptionPlaceholder:           deps.Labels.Form.DescriptionPlaceholder,
		Category:                         deps.Labels.Form.Category,
		CategoryPlaceholder:              deps.Labels.Form.CategoryPlaceholder,
		Location:                         deps.Labels.Form.Location,
		LocationPlaceholder:              deps.Labels.Form.LocationPlaceholder,
		AcquisitionCost:                  deps.Labels.Form.AcquisitionCost,
		AcquisitionCostPlaceholder:       deps.Labels.Form.AcquisitionCostPlaceholder,
		SalvageValue:                     deps.Labels.Form.SalvageValue,
		SalvageValuePlaceholder:          deps.Labels.Form.SalvageValuePlaceholder,
		UsefulLifeMonths:                 deps.Labels.Form.UsefulLifeMonths,
		UsefulLifePlaceholder:            deps.Labels.Form.UsefulLifePlaceholder,
		DepreciationMethod:               deps.Labels.Form.DepreciationMethod,
		DepMethodStraightLine:            deps.Labels.Form.DepMethodStraightLine,
		DepMethodDecliningBalance:        deps.Labels.Form.DepMethodDecliningBalance,
		DepMethodSumOfYears:              deps.Labels.Form.DepMethodSumOfYears,
		DepMethodUnitsOfProduction:       deps.Labels.Form.DepMethodUnitsOfProduction,
		Active:                           deps.Labels.Form.Active,
		AssetNumberInfo:                  deps.Labels.Form.AssetNumberInfo,
		AcquisitionCostInfo:              deps.Labels.Form.AcquisitionCostInfo,
		SalvageValueInfo:                 deps.Labels.Form.SalvageValueInfo,
		UsefulLifeMonthsInfo:             deps.Labels.Form.UsefulLifeMonthsInfo,
		DepreciationMethodInfo:           deps.Labels.Form.DepreciationMethodInfo,
		UnitsOfProductionDisabledTooltip: deps.Labels.Form.UnitsOfProductionDisabledTooltip,
	}
}
