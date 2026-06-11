package action

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/view"

	assetform "github.com/erniealice/fycha-golang/domain/asset/asset/form"
)

// NewAddAction creates the asset add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "create") {
			// 2026-05-14 permission-gates P3: error-shape — emit a proper
			// HTMX error (HX-Error-Message header + 422) rather than the
			// generic view.Error(...) which surfaces as a 500 page.
			return view.HTMXError(deps.Labels.Actions.NoPermission)
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
			return view.HTMXError(deps.Labels.Actions.InvalidFormData)
		}

		name := viewCtx.Request.FormValue("name")
		if name == "" {
			return view.HTMXError("Name is required")
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
				return view.HTMXError("Failed to create asset")
			}
		} else {
			log.Printf("Mock create asset: %s (no CreateAsset wired)", name)
		}

		return view.HTMXSuccess("assets-table")
	})
}
