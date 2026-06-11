// Package list provides the Fund source-list view for fund owners.
// Route: GET /app/funding/sources
// Shows all Funds owned by the current user (owner_party_type=USER, owner_party_id=session.user_id).
package list

import (
	"context"
	"fmt"
	"log"

	fundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund"
	funding "github.com/erniealice/fycha-golang/domain/funding"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// Deps holds view dependencies.
type Deps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       funding.FundingFormLabels
	ListFunds    func(ctx context.Context, req *fundpb.ListFundsRequest) (*fundpb.ListFundsResponse, error)
}

// PageData holds data for the source list page.
type PageData struct {
	types.PageData
	Rows []FundRow
}

// FundRow is the view-model for a single fund row.
type FundRow struct {
	ID              string
	Name            string
	Kind            string
	Currency        string
	AuthorizedLimit string
	Status          string
}

// NewView creates the fund source list view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		rows := fetchFunds(ctx, deps)
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Source.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "funding",
				ActiveSubNav:   "sources-all",
				HeaderTitle:    deps.Labels.Source.Title,
				HeaderSubtitle: deps.Labels.Source.Subtitle,
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
			Kind:            fundKindLabel(deps.Labels.Source.Kind, f.GetKind()),
			Currency:        f.GetCurrency(),
			AuthorizedLimit: fmt.Sprintf("%.2f", float64(f.GetAuthorizedLimit())/100.0),
			Status:          fundStatusLabel(deps.Labels.Source.Status, f.GetStatus()),
		})
	}
	return rows
}

func fundKindLabel(l funding.FundingSourceKindLabels, k fundpb.FundKind) string {
	switch k {
	case fundpb.FundKind_FUND_KIND_CASH_ON_HAND:
		return l.CashOnHand
	case fundpb.FundKind_FUND_KIND_BANK_ACCOUNT:
		return l.BankAccount
	case fundpb.FundKind_FUND_KIND_PETTY_CASH:
		return l.PettyCash
	case fundpb.FundKind_FUND_KIND_CREDIT_CARD:
		return l.CreditCard
	case fundpb.FundKind_FUND_KIND_CREDIT_LINE:
		return l.CreditLine
	case fundpb.FundKind_FUND_KIND_PREPAID_CARD:
		return l.PrepaidCard
	case fundpb.FundKind_FUND_KIND_MOBILE_MONEY:
		return l.MobileMoney
	default:
		return l.Unknown
	}
}

func fundStatusLabel(l funding.FundingSourceStatusLabels, s fundpb.FundStatus) string {
	switch s {
	case fundpb.FundStatus_FUND_STATUS_DRAFT:
		return l.Draft
	case fundpb.FundStatus_FUND_STATUS_ACTIVE:
		return l.Active
	case fundpb.FundStatus_FUND_STATUS_SUSPENDED:
		return l.Suspended
	case fundpb.FundStatus_FUND_STATUS_ARCHIVED:
		return l.Archived
	default:
		return l.Unknown
	}
}
