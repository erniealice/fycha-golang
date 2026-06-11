// Package list provides the Fund card-list view for consuming workspaces.
// Route: GET /app/funding/cards
// Lists FundAllocation rows for the current workspace.
package list

import (
	"context"
	"fmt"
	"log"

	fundallocationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_allocation"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels    pyeza.CommonLabels
	TableLabels     types.TableLabels
	ListAllocations func(ctx context.Context, req *fundallocationpb.ListFundAllocationsRequest) (*fundallocationpb.ListFundAllocationsResponse, error)
}

// CardRow is the view-model for a single allocation (card) row.
type CardRow struct {
	AllocationID   string
	FundID         string
	FundName       string
	AllocatedLimit string
	Mode           string
	Status         string
}

// PageData holds data for the card list page.
type PageData struct {
	types.PageData
	Rows []CardRow
}

// NewView creates the fund card list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		rows := fetchCards(ctx, deps)
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Funding Cards",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "funding",
				ActiveSubNav:   "cards-active",
				HeaderTitle:    "Funding Cards",
				HeaderSubtitle: "Shared funds allocated to this workspace",
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			Rows: rows,
		}
		tmpl := "funding-card-list"
		if viewCtx.IsHTMX {
			tmpl = "funding-card-list-content"
		}
		return view.OK(tmpl, pd)
	})
}

func fetchCards(ctx context.Context, deps *Deps) []CardRow {
	if deps.ListAllocations == nil {
		return nil
	}
	resp, err := deps.ListAllocations(ctx, &fundallocationpb.ListFundAllocationsRequest{})
	if err != nil {
		log.Printf("[funding] ListAllocations (cards) error: %v", err)
		return nil
	}
	if resp == nil {
		return nil
	}
	rows := make([]CardRow, 0, len(resp.GetData()))
	for _, a := range resp.GetData() {
		rows = append(rows, CardRow{
			AllocationID:   a.GetId(),
			FundID:         a.GetFundId(),
			FundName:       a.GetFundId(), // resolved to fund name when ReadFund is wired
			AllocatedLimit: fmt.Sprintf("%.2f", float64(a.GetAllocatedLimit())/100.0),
			Mode:           a.GetMode().String(),
			Status:         a.GetStatus().String(),
		})
	}
	return rows
}
