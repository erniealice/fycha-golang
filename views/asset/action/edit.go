package action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	assetform "github.com/erniealice/fycha-golang/views/asset/form"
)

// NewEditAction creates the asset edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "update") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Load asset from DB for edit form pre-fill
			if deps.ReadAsset != nil {
				record, err := deps.ReadAsset(ctx, id)
				if err != nil {
					log.Printf("asset read error for edit: %v", err)
					return fycha.HTMXError("Failed to read asset")
				}
				return view.OK("asset-drawer-form", &assetform.Data{
					FormAction:         route.ResolveURL(deps.Routes.EditURL, "id", record.ID),
					IsEdit:             true,
					ID:                 record.ID,
					Name:               record.Name,
					AssetNumber:        record.AssetNumber,
					Description:        record.Description,
					CategoryID:         record.AssetCategoryID,
					LocationID:         record.LocationID,
					AcquisitionCost:    fmt.Sprintf("%.2f", record.AcquisitionCost),
					SalvageValue:       fmt.Sprintf("%.2f", record.SalvageValue),
					UsefulLifeMonths:   strconv.Itoa(record.UsefulLifeMonths),
					DepreciationMethod: record.DepreciationMethod,
					Active:             record.Active,
					Labels:             labelsFromDeps(deps),
					CommonLabels:       nil,
				})
			}

			// Fallback: mock data
			return view.OK("asset-drawer-form", &assetform.Data{
				FormAction:         route.ResolveURL(deps.Routes.EditURL, "id", id),
				IsEdit:             true,
				ID:                 id,
				Name:               "Mock Asset",
				AssetNumber:        "FA-001",
				Description:        "Mock asset for development",
				AcquisitionCost:    "85000.00",
				SalvageValue:       "5000.00",
				UsefulLifeMonths:   "60",
				DepreciationMethod: "straight_line",
				Active:             true,
				Labels:             labelsFromDeps(deps),
				CommonLabels:       nil,
			})
		}

		// POST — update asset
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

		depMethod := viewCtx.Request.FormValue("depreciation_method")
		if depMethod == "" {
			depMethod = "STRAIGHT_LINE"
		}

		record := &assetform.Record{
			ID:                 id,
			AssetNumber:        viewCtx.Request.FormValue("asset_number"),
			Name:               name,
			Description:        viewCtx.Request.FormValue("description"),
			AssetCategoryID:    viewCtx.Request.FormValue("asset_category_id"),
			LocationID:         viewCtx.Request.FormValue("location_id"),
			AcquisitionCost:    acqCost,
			SalvageValue:       salvage,
			BookValue:          acqCost - salvage,
			UsefulLifeMonths:   usefulLife,
			DepreciationMethod: depMethod,
			Currency:           "PHP",
		}

		if deps.UpdateAsset != nil {
			if err := deps.UpdateAsset(ctx, record); err != nil {
				log.Printf("asset update error: %v", err)
				return fycha.HTMXError("Failed to update asset")
			}
		} else {
			log.Printf("Mock update asset %s: %s (no UpdateAsset wired)", id, name)
		}

		return fycha.HTMXSuccess("assets-table")
	})
}
