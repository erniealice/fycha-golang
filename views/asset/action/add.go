package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	assetform "github.com/erniealice/fycha-golang/views/asset/form"
)

// NewAddAction creates the asset add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "create") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("asset-drawer-form", &assetform.Data{
				FormAction:         deps.Routes.AddURL,
				Active:             true,
				DepreciationMethod: "straight_line",
				Labels:             labelsFromDeps(deps),
				CommonLabels:       nil, // injected by ViewAdapter
			})
		}

		// POST — create asset
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.InvalidFormData)
		}

		name := viewCtx.Request.FormValue("name")
		if name == "" {
			return fycha.HTMXError("Name is required")
		}

		acqCost, _ := strconv.ParseFloat(viewCtx.Request.FormValue("acquisition_cost"), 64)
		salvage, _ := strconv.ParseFloat(viewCtx.Request.FormValue("salvage_value"), 64)
		usefulLife, _ := strconv.Atoi(viewCtx.Request.FormValue("useful_life_months"))

		id := ""
		if deps.NewID != nil {
			id = deps.NewID()
		}

		assetNumber := viewCtx.Request.FormValue("asset_number")
		if assetNumber == "" {
			assetNumber = id
		}

		depMethod := viewCtx.Request.FormValue("depreciation_method")
		if depMethod == "" {
			depMethod = "STRAIGHT_LINE"
		}

		record := &assetform.Record{
			ID:                 id,
			AssetNumber:        assetNumber,
			Name:               name,
			Description:        viewCtx.Request.FormValue("description"),
			AssetType:          "PPE",
			AssetCategoryID:    viewCtx.Request.FormValue("asset_category_id"),
			LocationID:         viewCtx.Request.FormValue("location_id"),
			AcquisitionCost:    acqCost,
			SalvageValue:       salvage,
			BookValue:          acqCost - salvage,
			UsefulLifeMonths:   usefulLife,
			DepreciationMethod: depMethod,
			Currency:           "PHP",
			Status:             "IN_SERVICE",
			Active:             true,
		}

		if deps.CreateAsset != nil {
			if err := deps.CreateAsset(ctx, record); err != nil {
				log.Printf("asset create error: %v", err)
				return fycha.HTMXError("Failed to create asset")
			}
		} else {
			log.Printf("Mock create asset: %s (no CreateAsset wired)", name)
		}

		return fycha.HTMXSuccess("assets-table")
	})
}
