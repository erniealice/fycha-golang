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
)

// AssetRecord is a flat struct for passing asset data between action handlers
// and the DB layer. It avoids a dependency on protobuf types.
type AssetRecord struct {
	ID                 string
	AssetNumber        string
	Name               string
	Description        string
	AssetType          string
	AssetCategoryID    string
	LocationID         string
	AcquisitionCost    float64
	SalvageValue       float64
	BookValue          float64
	UsefulLifeMonths   int
	DepreciationMethod string
	Currency           string
	Status             string
	Active             bool
}

// FormLabels holds i18n labels for the drawer form template.
type FormLabels struct {
	Name                       string
	NamePlaceholder            string
	AssetNumber                string
	AssetNumberPlaceholder     string
	Description                string
	DescriptionPlaceholder     string
	Category                   string
	CategoryPlaceholder        string
	Location                   string
	LocationPlaceholder        string
	AcquisitionCost            string
	AcquisitionCostPlaceholder string
	SalvageValue               string
	SalvageValuePlaceholder    string
	UsefulLifeMonths           string
	UsefulLifePlaceholder      string
	DepreciationMethod         string
	DepMethodStraightLine      string
	DepMethodDecliningBalance  string
	DepMethodSumOfYears        string
	DepMethodUnitsOfProduction string
	Active                     string
}

// FormData is the template data for the asset drawer form.
type FormData struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	Name               string
	AssetNumber        string
	Description        string
	CategoryID         string
	LocationID         string
	AcquisitionCost    string
	SalvageValue       string
	UsefulLifeMonths   string
	DepreciationMethod string
	Active             bool
	Labels             FormLabels
	CommonLabels       any
}

// Deps holds dependencies for asset action handlers.
type Deps struct {
	Routes fycha.AssetRoutes
	Labels fycha.AssetLabels

	// CRUD operations (wired from block.go)
	CreateAsset func(ctx context.Context, asset *AssetRecord) error
	ReadAsset   func(ctx context.Context, id string) (*AssetRecord, error)
	UpdateAsset func(ctx context.Context, asset *AssetRecord) error
	DeleteAsset func(ctx context.Context, id string) error
	SetActive   func(ctx context.Context, id string, active bool) error
	NewID       func() string
}

func formLabelsFromStruct(l fycha.AssetFormLabels) FormLabels {
	return FormLabels{
		Name:                       l.Name,
		NamePlaceholder:            l.NamePlaceholder,
		AssetNumber:                l.AssetNumber,
		AssetNumberPlaceholder:     l.AssetNumberPlaceholder,
		Description:                l.Description,
		DescriptionPlaceholder:     l.DescriptionPlaceholder,
		Category:                   l.Category,
		CategoryPlaceholder:        l.CategoryPlaceholder,
		Location:                   l.Location,
		LocationPlaceholder:        l.LocationPlaceholder,
		AcquisitionCost:            l.AcquisitionCost,
		AcquisitionCostPlaceholder: l.AcquisitionCostPlaceholder,
		SalvageValue:               l.SalvageValue,
		SalvageValuePlaceholder:    l.SalvageValuePlaceholder,
		UsefulLifeMonths:           l.UsefulLifeMonths,
		UsefulLifePlaceholder:      l.UsefulLifePlaceholder,
		DepreciationMethod:         l.DepreciationMethod,
		DepMethodStraightLine:      l.DepMethodStraightLine,
		DepMethodDecliningBalance:  l.DepMethodDecliningBalance,
		DepMethodSumOfYears:        l.DepMethodSumOfYears,
		DepMethodUnitsOfProduction: l.DepMethodUnitsOfProduction,
		Active:                     l.Active,
	}
}

// NewAddAction creates the asset add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "create") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("asset-drawer-form", &FormData{
				FormAction:         deps.Routes.AddURL,
				Active:             true,
				DepreciationMethod: "straight_line",
				Labels:             formLabelsFromStruct(deps.Labels.Form),
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

		record := &AssetRecord{
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
				return view.OK("asset-drawer-form", &FormData{
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
					Labels:             formLabelsFromStruct(deps.Labels.Form),
					CommonLabels:       nil,
				})
			}

			// Fallback: mock data
			return view.OK("asset-drawer-form", &FormData{
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
				Labels:             formLabelsFromStruct(deps.Labels.Form),
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

		record := &AssetRecord{
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

// NewDeleteAction creates the asset delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "delete") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
		}
		if id == "" {
			return fycha.HTMXError(deps.Labels.Actions.IDRequired)
		}

		if deps.DeleteAsset != nil {
			if err := deps.DeleteAsset(ctx, id); err != nil {
				log.Printf("asset delete error: %v", err)
				return fycha.HTMXError("Failed to delete asset")
			}
		} else {
			log.Printf("Mock delete asset: %s", id)
		}

		return fycha.HTMXSuccess("assets-table")
	})
}

// NewBulkDeleteAction creates the asset bulk delete action (POST only).
func NewBulkDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "delete") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		if len(ids) == 0 {
			return fycha.HTMXError(deps.Labels.Actions.NoIDsProvided)
		}

		if deps.DeleteAsset != nil {
			for _, id := range ids {
				if err := deps.DeleteAsset(ctx, id); err != nil {
					log.Printf("asset bulk delete error for %s: %v", id, err)
				}
			}
		} else {
			log.Printf("Mock bulk delete assets: %v", ids)
		}

		return fycha.HTMXSuccess("assets-table")
	})
}

// NewSetStatusAction creates the asset activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "update") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.URL.Query().Get("id")
		targetStatus := viewCtx.Request.URL.Query().Get("status")

		if id == "" {
			_ = viewCtx.Request.ParseForm()
			id = viewCtx.Request.FormValue("id")
			targetStatus = viewCtx.Request.FormValue("status")
		}
		if id == "" {
			return fycha.HTMXError(deps.Labels.Actions.IDRequired)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return fycha.HTMXError(deps.Labels.Actions.InvalidStatus)
		}

		active := targetStatus == "active"

		if deps.SetActive != nil {
			if err := deps.SetActive(ctx, id, active); err != nil {
				log.Printf("asset set-status error: %v", err)
				return fycha.HTMXError("Failed to update asset")
			}
		} else {
			log.Printf("Mock set asset status %s: %s", id, targetStatus)
		}

		return fycha.HTMXSuccess("assets-table")
	})
}

// NewBulkSetStatusAction creates the asset bulk activate/deactivate action (POST only).
func NewBulkSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "update") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		_ = viewCtx.Request.ParseMultipartForm(32 << 20)

		ids := viewCtx.Request.Form["id"]
		targetStatus := viewCtx.Request.FormValue("target_status")

		if len(ids) == 0 {
			return fycha.HTMXError(deps.Labels.Actions.NoIDsProvided)
		}
		if targetStatus != "active" && targetStatus != "inactive" {
			return fycha.HTMXError(deps.Labels.Actions.InvalidTargetStatus)
		}

		active := targetStatus == "active"

		if deps.SetActive != nil {
			for _, id := range ids {
				if err := deps.SetActive(ctx, id, active); err != nil {
					log.Printf("asset bulk set-status error for %s: %v", id, err)
				}
			}
		} else {
			log.Printf("Mock bulk set asset status %v: %s", ids, targetStatus)
		}

		return fycha.HTMXSuccess("assets-table")
	})
}
