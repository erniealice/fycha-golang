// Package list provides the Fund source-list view for fund owners.
// Route: GET /app/funding/sources
// Shows all Funds owned by the current user (owner_party_type=USER, owner_party_id=session.user_id).
package list

import (
	"context"
	"fmt"
	"log"

	fundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	ListFunds    func(ctx context.Context, req *fundpb.ListFundsRequest) (*fundpb.ListFundsResponse, error)
}

// PageData holds data for the source list page.
type PageData struct {
	types.PageData
	Rows []FundRow
}

// FundRow is the view-model for a single fund row.
type FundRow struct {
	ID             string
	Name           string
	Kind           string
	Currency       string
	AuthorizedLimit string
	Status         string
}

// NewView creates the fund source list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		rows := fetchFunds(ctx, deps)
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Fund Sources",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "funding",
				ActiveSubNav:   "sources-all",
				HeaderTitle:    "Fund Sources",
				HeaderSubtitle: "Funds you own and share with workspaces",
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			Rows: rows,
		}
		tmpl := "funding-source-list"
		if viewCtx.IsHTMX {
			tmpl = "funding-source-list-content"
		}
		return view.OK(tmpl, pd)
	})
}

func fetchFunds(ctx context.Context, deps *Deps) []FundRow {
	if deps.ListFunds == nil {
		return nil
	}
	resp, err := deps.ListFunds(ctx, &fundpb.ListFundsRequest{})
	if err != nil {
		log.Printf("[funding] ListFunds error: %v", err)
		return nil
	}
	if resp == nil {
		return nil
	}
	rows := make([]FundRow, 0, len(resp.GetData()))
	for _, f := range resp.GetData() {
		rows = append(rows, FundRow{
			ID:              f.GetId(),
			Name:            f.GetName(),
			Kind:            fundKindLabel(f.GetKind()),
			Currency:        f.GetCurrency(),
			AuthorizedLimit: fmt.Sprintf("%.2f", float64(f.GetAuthorizedLimit())/100.0),
			Status:          fundStatusLabel(f.GetStatus()),
		})
	}
	return rows
}

func fundKindLabel(k fundpb.FundKind) string {
	switch k {
	case fundpb.FundKind_FUND_KIND_CASH_ON_HAND:
		return "Cash on Hand"
	case fundpb.FundKind_FUND_KIND_BANK_ACCOUNT:
		return "Bank Account"
	case fundpb.FundKind_FUND_KIND_PETTY_CASH:
		return "Petty Cash"
	case fundpb.FundKind_FUND_KIND_CREDIT_CARD:
		return "Credit Card"
	case fundpb.FundKind_FUND_KIND_CREDIT_LINE:
		return "Credit Line"
	case fundpb.FundKind_FUND_KIND_PREPAID_CARD:
		return "Prepaid Card"
	case fundpb.FundKind_FUND_KIND_MOBILE_MONEY:
		return "Mobile Money"
	default:
		return "Unknown"
	}
}

func fundStatusLabel(s fundpb.FundStatus) string {
	switch s {
	case fundpb.FundStatus_FUND_STATUS_DRAFT:
		return "Draft"
	case fundpb.FundStatus_FUND_STATUS_ACTIVE:
		return "Active"
	case fundpb.FundStatus_FUND_STATUS_SUSPENDED:
		return "Suspended"
	case fundpb.FundStatus_FUND_STATUS_ARCHIVED:
		return "Archived"
	default:
		return "Unknown"
	}
}
