// Package form provides the FundAllocation drawer-form stub.
// Renders a create/edit allocation form; submit handler is stubbed for FS-E.
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

// PageData holds data for the allocation form drawer.
type PageData struct {
	types.PageData
	FundID       string
	AllocationID string
}

// NewView creates the allocation drawer-form view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		fundID := viewCtx.Request.URL.Query().Get("fund_id")
		allocationID := viewCtx.Request.URL.Query().Get("allocation_id")
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Allocation",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "funding",
				CommonLabels: deps.CommonLabels,
			},
			FundID:       fundID,
			AllocationID: allocationID,
		}
		return view.OK("funding-allocation-form", pd)
	})
}
