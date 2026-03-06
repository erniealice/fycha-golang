package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

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

		// POST -- create asset (mock — just return success)
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.InvalidFormData)
		}

		log.Printf("Mock create asset: %s", viewCtx.Request.FormValue("name"))
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
			// Mock data for edit form
			return view.OK("asset-drawer-form", &FormData{
				FormAction:         deps.Routes.EditURL,
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

		// POST -- update asset (mock — just return success)
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.InvalidFormData)
		}

		log.Printf("Mock update asset %s: %s", id, viewCtx.Request.FormValue("name"))
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

		log.Printf("Mock delete asset: %s", id)
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

		log.Printf("Mock bulk delete assets: %v", ids)
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

		log.Printf("Mock set asset status %s: %s", id, targetStatus)
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

		log.Printf("Mock bulk set asset status %v: %s", ids, targetStatus)
		return fycha.HTMXSuccess("assets-table")
	})
}
