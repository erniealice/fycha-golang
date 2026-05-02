package action

import (
	"context"

	fycha "github.com/erniealice/fycha-golang"
	assetform "github.com/erniealice/fycha-golang/views/asset/form"
)

// Deps holds dependencies for asset action handlers.
type Deps struct {
	Routes fycha.AssetRoutes
	Labels fycha.AssetLabels

	// CRUD operations (wired from block.go)
	CreateAsset func(ctx context.Context, asset *assetform.Record) error
	ReadAsset   func(ctx context.Context, id string) (*assetform.Record, error)
	UpdateAsset func(ctx context.Context, asset *assetform.Record) error
	DeleteAsset func(ctx context.Context, id string) error
	SetActive   func(ctx context.Context, id string, active bool) error
	NewID       func() string
}

// labelsFromDeps builds form.Labels from the deps label tree. The inline
// construction at each call site was identical; this helper eliminates the
// repetition without introducing a mapper that transforms no values.
func labelsFromDeps(deps *Deps) assetform.Labels {
	return assetform.Labels{
		Name:                       deps.Labels.Form.Name,
		NamePlaceholder:            deps.Labels.Form.NamePlaceholder,
		AssetNumber:                deps.Labels.Form.AssetNumber,
		AssetNumberPlaceholder:     deps.Labels.Form.AssetNumberPlaceholder,
		Description:                deps.Labels.Form.Description,
		DescriptionPlaceholder:     deps.Labels.Form.DescriptionPlaceholder,
		Category:                   deps.Labels.Form.Category,
		CategoryPlaceholder:        deps.Labels.Form.CategoryPlaceholder,
		Location:                   deps.Labels.Form.Location,
		LocationPlaceholder:        deps.Labels.Form.LocationPlaceholder,
		AcquisitionCost:            deps.Labels.Form.AcquisitionCost,
		AcquisitionCostPlaceholder: deps.Labels.Form.AcquisitionCostPlaceholder,
		SalvageValue:               deps.Labels.Form.SalvageValue,
		SalvageValuePlaceholder:    deps.Labels.Form.SalvageValuePlaceholder,
		UsefulLifeMonths:           deps.Labels.Form.UsefulLifeMonths,
		UsefulLifePlaceholder:      deps.Labels.Form.UsefulLifePlaceholder,
		DepreciationMethod:         deps.Labels.Form.DepreciationMethod,
		DepMethodStraightLine:      deps.Labels.Form.DepMethodStraightLine,
		DepMethodDecliningBalance:  deps.Labels.Form.DepMethodDecliningBalance,
		DepMethodSumOfYears:        deps.Labels.Form.DepMethodSumOfYears,
		DepMethodUnitsOfProduction: deps.Labels.Form.DepMethodUnitsOfProduction,
		Active:                     deps.Labels.Form.Active,
		AssetNumberInfo:            deps.Labels.Form.AssetNumberInfo,
		AcquisitionCostInfo:        deps.Labels.Form.AcquisitionCostInfo,
		SalvageValueInfo:           deps.Labels.Form.SalvageValueInfo,
		UsefulLifeMonthsInfo:       deps.Labels.Form.UsefulLifeMonthsInfo,
		DepreciationMethodInfo:     deps.Labels.Form.DepreciationMethodInfo,
	}
}
