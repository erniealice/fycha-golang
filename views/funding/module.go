// Package funding provides stub views for the funding domain:
// source list + detail (fund owner), card list + detail (workspace consumer),
// and four drawer-form stubs (allocation, draw, settlement, transfer).
package funding

import (
	"context"

	fundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund"
	fundallocationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_allocation"
	fundtransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_transaction"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	allocationform "github.com/erniealice/fycha-golang/views/funding/allocation/form"
	carddetail "github.com/erniealice/fycha-golang/views/funding/card/detail"
	cardlist "github.com/erniealice/fycha-golang/views/funding/card/list"
	drawform "github.com/erniealice/fycha-golang/views/funding/draw/form"
	settlementform "github.com/erniealice/fycha-golang/views/funding/settlement/form"
	sourcedetail "github.com/erniealice/fycha-golang/views/funding/source/detail"
	sourcelist "github.com/erniealice/fycha-golang/views/funding/source/list"
	transferform "github.com/erniealice/fycha-golang/views/funding/transfer/form"
)

// ModuleDeps holds all dependencies for the funding view module.
type ModuleDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Route URLs for the 8 funding views. Each defaults to its /app/* constant
	// when empty for backward compatibility during the P3–P11 parallel-mux window.
	SourceListURL      string // e.g. "/app/funding/sources"
	SourceDetailURL    string // e.g. "/app/funding/sources/{fund_id}"
	CardListURL        string // e.g. "/app/funding/cards"
	CardDetailURL      string // e.g. "/app/funding/cards/{allocation_id}"
	AllocationFormURL  string // e.g. "/app/funding/allocation/form"
	DrawFormURL        string // e.g. "/app/funding/draw/form"
	SettlementFormURL  string // e.g. "/app/funding/settlement/form"
	TransferFormURL    string // e.g. "/app/funding/transfer/form"

	// Fund use cases
	ReadFund  func(ctx context.Context, req *fundpb.ReadFundRequest) (*fundpb.ReadFundResponse, error)
	ListFunds func(ctx context.Context, req *fundpb.ListFundsRequest) (*fundpb.ListFundsResponse, error)

	// FundAllocation use cases
	ReadAllocation  func(ctx context.Context, req *fundallocationpb.ReadFundAllocationRequest) (*fundallocationpb.ReadFundAllocationResponse, error)
	ListAllocations func(ctx context.Context, req *fundallocationpb.ListFundAllocationsRequest) (*fundallocationpb.ListFundAllocationsResponse, error)

	// FundTransaction use cases
	ListTransactions func(ctx context.Context, req *fundtransactionpb.ListFundTransactionsRequest) (*fundtransactionpb.ListFundTransactionsResponse, error)
}

// Module holds all constructed funding views.
type Module struct {
	SourceList     view.View
	SourceDetail   view.View
	CardList       view.View
	CardDetail     view.View
	AllocationForm view.View
	DrawForm       view.View
	SettlementForm view.View
	TransferForm   view.View
	// Route URLs (set from ModuleDeps; used by RegisterRoutes).
	sourceListURL     string
	sourceDetailURL   string
	cardListURL       string
	cardDetailURL     string
	allocationFormURL string
	drawFormURL       string
	settlementFormURL string
	transferFormURL   string
}

// NewModule constructs a fully wired funding view module.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}
	return &Module{
		sourceListURL:     deps.SourceListURL,
		sourceDetailURL:   deps.SourceDetailURL,
		cardListURL:       deps.CardListURL,
		cardDetailURL:     deps.CardDetailURL,
		allocationFormURL: deps.AllocationFormURL,
		drawFormURL:       deps.DrawFormURL,
		settlementFormURL: deps.SettlementFormURL,
		transferFormURL:   deps.TransferFormURL,
		SourceList: sourcelist.NewView(&sourcelist.Deps{
			CommonLabels: deps.CommonLabels,
			TableLabels:  deps.TableLabels,
			ListFunds:    deps.ListFunds,
		}),
		SourceDetail: sourcedetail.NewView(&sourcedetail.Deps{
			CommonLabels:     deps.CommonLabels,
			ReadFund:         deps.ReadFund,
			ListAllocations:  deps.ListAllocations,
			ListTransactions: deps.ListTransactions,
		}),
		CardList: cardlist.NewView(&cardlist.Deps{
			CommonLabels:    deps.CommonLabels,
			TableLabels:     deps.TableLabels,
			ListAllocations: deps.ListAllocations,
		}),
		CardDetail: carddetail.NewView(&carddetail.Deps{
			CommonLabels:     deps.CommonLabels,
			ReadFund:         deps.ReadFund,
			ReadAllocation:   deps.ReadAllocation,
			ListTransactions: deps.ListTransactions,
		}),
		AllocationForm: allocationform.NewView(&allocationform.Deps{
			CommonLabels: deps.CommonLabels,
		}),
		DrawForm: drawform.NewView(&drawform.Deps{
			CommonLabels: deps.CommonLabels,
		}),
		SettlementForm: settlementform.NewView(&settlementform.Deps{
			CommonLabels: deps.CommonLabels,
		}),
		TransferForm: transferform.NewView(&transferform.Deps{
			CommonLabels: deps.CommonLabels,
		}),
	}
}

// resolveURL returns url if non-empty, otherwise fallback.
func resolveURL(url, fallback string) string {
	if url == "" {
		return fallback
	}
	return url
}

// RegisterRoutes registers all 8 funding routes with the given registrar.
// Route paths are read from the ModuleDeps URL fields (set during construction);
// each falls back to its /app/* default when empty.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Source-side views (fund owners — cross-workspace)
	r.GET(resolveURL(m.sourceListURL, "/app/funding/sources"), m.SourceList)
	r.GET(resolveURL(m.sourceDetailURL, "/app/funding/sources/{fund_id}"), m.SourceDetail)

	// Card-side views (fund consumers — within current workspace)
	r.GET(resolveURL(m.cardListURL, "/app/funding/cards"), m.CardList)
	r.GET(resolveURL(m.cardDetailURL, "/app/funding/cards/{allocation_id}"), m.CardDetail)

	// Drawer forms (GET — loaded into side sheet)
	r.GET(resolveURL(m.allocationFormURL, "/app/funding/allocation/form"), m.AllocationForm)
	r.GET(resolveURL(m.drawFormURL, "/app/funding/draw/form"), m.DrawForm)
	r.GET(resolveURL(m.settlementFormURL, "/app/funding/settlement/form"), m.SettlementForm)
	r.GET(resolveURL(m.transferFormURL, "/app/funding/transfer/form"), m.TransferForm)
}
