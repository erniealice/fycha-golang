// Package form provides the FundTransaction Settlement drawer-form stub.
// Renders a "settle outstanding" form; submit handler is stubbed for FS-E.
package form

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
}

// PageData holds data for the settlement form drawer.
type PageData struct {
	types.PageData
	AllocationID string
}

// NewView creates the settlement drawer-form view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		allocationID := viewCtx.Request.URL.Query().Get("allocation_id")
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Settle Outstanding",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "funding",
				CommonLabels: deps.CommonLabels,
			},
			AllocationID: allocationID,
		}
		return view.OK("funding-settlement-form", pd)
	})
}
