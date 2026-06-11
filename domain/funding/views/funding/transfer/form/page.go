// Package form provides the FundTransaction Transfer drawer-form stub.
// Renders a "transfer between funds" form; submit handler is stubbed for FS-E.
package form

import (
	"context"

	funding "github.com/erniealice/fycha-golang/domain/funding"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
	Labels       funding.FundingFormLabels
}

// PageData holds data for the transfer form drawer.
type PageData struct {
	types.PageData
	SourceFundID string
	Labels       funding.FundingFormLabels
}

// NewView creates the transfer drawer-form view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		sourceFundID := viewCtx.Request.URL.Query().Get("source_fund_id")
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        deps.Labels.Transfer.Submit,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "funding",
				CommonLabels: deps.CommonLabels,
			},
			SourceFundID: sourceFundID,
			Labels:       deps.Labels,
		}
		return view.OK("funding-transfer-form", pd)
	})
}
