// Package detail provides the Fund card-detail view for consuming workspaces.
// Route: GET /app/funding/cards/{allocation_id}
// Shows the parent Fund's info + this workspace's allocation + recent FundTransactions.
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
	ReadAllocation   func(ctx context.Context, req *fundallocationpb.ReadFundAllocationRequest) (*fundallocationpb.ReadFundAllocationResponse, error)
	ListTransactions func(ctx context.Context, req *fundtransactionpb.ListFundTransactionsRequest) (*fundtransactionpb.ListFundTransactionsResponse, error)
}

// TransactionRow is the view-model for a single transaction row.
type TransactionRow struct {
	ID     string
	Kind   string
	Amount string
	Status string
}

// PageData holds data for the card detail page.
type PageData struct {
	types.PageData
	AllocationID     string
	FundID           string
	FundName         string
	FundKind         string
	FundCurrency     string
	AllocatedLimit   string
	AllocationMode   string
	AllocationStatus string
	Transactions     []TransactionRow
}

// NewView creates the fund card detail view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		allocationID := viewCtx.Request.PathValue("allocation_id")

		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Funding Card",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "funding",
				ActiveSubNav:   "cards-active",
				HeaderTitle:    "Funding Card",
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			AllocationID: allocationID,
		}

		// Resolve allocation via Read (pass id in Data)
		if deps.ReadAllocation != nil && allocationID != "" {
			resp, err := deps.ReadAllocation(ctx, &fundallocationpb.ReadFundAllocationRequest{
				Data: &fundallocationpb.FundAllocation{Id: allocationID},
			})
			if err != nil {
				log.Printf("[funding] ReadAllocation error: %v", err)
			} else if resp != nil && len(resp.GetData()) > 0 {
				a := resp.GetData()[0]
				pd.FundID = a.GetFundId()
				pd.AllocatedLimit = fmt.Sprintf("%.2f", float64(a.GetAllocatedLimit())/100.0)
				pd.AllocationMode = a.GetMode().String()
				pd.AllocationStatus = a.GetStatus().String()
			}
		}

		// Resolve parent Fund
		if deps.ReadFund != nil && pd.FundID != "" {
			resp, err := deps.ReadFund(ctx, &fundpb.ReadFundRequest{
				Data: &fundpb.Fund{Id: pd.FundID},
			})
			if err != nil {
				log.Printf("[funding] ReadFund error: %v", err)
			} else if resp != nil && len(resp.GetData()) > 0 {
				f := resp.GetData()[0]
				pd.FundName = f.GetName()
				pd.FundKind = f.GetKind().String()
				pd.FundCurrency = f.GetCurrency()
				pd.PageData.Title = f.GetName()
				pd.PageData.HeaderTitle = f.GetName()
			}
		}

		// Resolve recent transactions for this allocation
		if deps.ListTransactions != nil && allocationID != "" {
			resp, err := deps.ListTransactions(ctx, &fundtransactionpb.ListFundTransactionsRequest{})
			if err != nil {
				log.Printf("[funding] ListTransactions (card) error: %v", err)
			} else if resp != nil {
				for _, t := range resp.GetData() {
					if t.GetAllocationId() != allocationID {
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

		tmpl := "funding-card-detail"
		if viewCtx.IsHTMX {
			tmpl = "funding-card-detail-content"
		}
		return view.OK(tmpl, pd)
	})
}
