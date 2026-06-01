// Package form provides the FundTransaction Draw drawer-form stub.
// Renders a "charge to card" form; submit handler is stubbed for FS-E.
package form

import (
	"context"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
	Labels       fycha.FundingFormLabels
}

// PageData holds data for the draw form drawer.
type PageData struct {
	types.PageData
	AllocationID string
	Labels       fycha.FundingFormLabels
}

// NewView creates the draw drawer-form view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		allocationID := viewCtx.Request.URL.Query().Get("allocation_id")
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Draw.Submit,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "funding",
				CommonLabels: deps.CommonLabels,
			},
			AllocationID: allocationID,
			Labels:       deps.Labels,
		}
		return view.OK("funding-draw-form", pd)
	})
}
