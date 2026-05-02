package action

import (
	"context"
	"fmt"
	"log"

	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

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
