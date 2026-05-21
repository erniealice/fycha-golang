// Package detail provides the Fund source-detail view.
// Route: GET /app/funding/sources/{fund_id}
// Shows Info / Allocations / Activity tabs for a single Fund.
package detail

import (
	"context"
	"fmt"
	"log"

	fundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund"
	fundallocationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_allocation"
	fundtransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_transaction"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels     pyeza.CommonLabels
	ReadFund         func(ctx context.Context, req *fundpb.ReadFundRequest) (*fundpb.ReadFundResponse, error)
	ListAllocations  func(ctx context.Context, req *fundallocationpb.ListFundAllocationsRequest) (*fundallocationpb.ListFundAllocationsResponse, error)
	ListTransactions func(ctx context.Context, req *fundtransactionpb.ListFundTransactionsRequest) (*fundtransactionpb.ListFundTransactionsResponse, error)
}

// AllocationRow is the view-model for a single allocation row.
type AllocationRow struct {
	ID             string
	WorkspaceID    string
	AllocatedLimit string
	Mode           string
	Status         string
}

// TransactionRow is the view-model for a single transaction row.
type TransactionRow struct {
	ID     string
	Kind   string
	Amount string
	Status string
}

// PageData holds data for the source detail page.
type PageData struct {
	types.PageData
	FundID          string
	FundName        string
	FundKind        string
	FundCurrency    string
	AuthorizedLimit string
	FundStatus      string
	ActiveTab       string
	Allocations     []AllocationRow
	Transactions    []TransactionRow
}

// NewView creates the fund source detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		fundID := viewCtx.Request.PathValue("fund_id")
		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        "Fund Source",
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "funding",
				ActiveSubNav: "sources-all",
				HeaderTitle:  "Fund Source",
				HeaderIcon:   "icon-credit-card",
				CommonLabels: deps.CommonLabels,
			},
			FundID:    fundID,
			ActiveTab: activeTab,
		}

		// Resolve fund via Read (pass id in Data)
		if deps.ReadFund != nil && fundID != "" {
			resp, err := deps.ReadFund(ctx, &fundpb.ReadFundRequest{
				Data: &fundpb.Fund{Id: fundID},
			})
			if err != nil {
				log.Printf("[funding] ReadFund error: %v", err)
			} else if resp != nil && len(resp.GetData()) > 0 {
				f := resp.GetData()[0]
				pd.FundName = f.GetName()
				pd.FundKind = f.GetKind().String()
				pd.FundCurrency = f.GetCurrency()
				pd.AuthorizedLimit = fmt.Sprintf("%.2f", float64(f.GetAuthorizedLimit())/100.0)
				pd.FundStatus = f.GetStatus().String()
				pd.PageData.Title = f.GetName()
				pd.PageData.HeaderTitle = f.GetName()
			}
		}

		if activeTab == "allocations" && deps.ListAllocations != nil {
			resp, err := deps.ListAllocations(ctx, &fundallocationpb.ListFundAllocationsRequest{})
			if err != nil {
				log.Printf("[funding] ListAllocations error: %v", err)
			} else if resp != nil {
				for _, a := range resp.GetData() {
					if a.GetFundId() != fundID {
						continue
					}
					pd.Allocations = append(pd.Allocations, AllocationRow{
						ID:             a.GetId(),
						WorkspaceID:    a.GetWorkspaceId(),
						AllocatedLimit: fmt.Sprintf("%.2f", float64(a.GetAllocatedLimit())/100.0),
						Mode:           a.GetMode().String(),
						Status:         a.GetStatus().String(),
					})
				}
			}
		}

		if activeTab == "activity" && deps.ListTransactions != nil {
			resp, err := deps.ListTransactions(ctx, &fundtransactionpb.ListFundTransactionsRequest{})
			if err != nil {
				log.Printf("[funding] ListTransactions error: %v", err)
			} else if resp != nil {
				for _, t := range resp.GetData() {
					if t.GetFundId() != fundID {
						continue
					}
					pd.Transactions = append(pd.Transactions, TransactionRow{
						ID:     t.GetId(),
						Kind:   t.GetKind().String(),
						Amount: fmt.Sprintf("%.2f", float64(t.GetAmount())/100.0),
						Status: t.GetStatus().String(),
					})
				}
			}
		}

		tmpl := "funding-source-detail"
		if viewCtx.IsHTMX {
			tmpl = "funding-source-detail-content"
		}
		return view.OK(tmpl, pd)
	})
}
