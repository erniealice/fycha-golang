package action

import (
	"context"
	"log"

	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// NewSetStatusAction creates the asset activate/deactivate action (POST only).
func NewSetStatusAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("asset", "update") {
			// 2026-05-14 permission-gates P3: error-shape fix.
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
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
			// 2026-05-14 permission-gates P3: error-shape fix.
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
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
