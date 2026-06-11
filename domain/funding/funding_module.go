package funding

import (
	"context"

	fundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund"
	fundallocationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_allocation"
	fundtransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_transaction"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	allocationform "github.com/erniealice/fycha-golang/domain/funding/funding/allocation/form"
	carddetail "github.com/erniealice/fycha-golang/domain/funding/funding/card/detail"
	cardlist "github.com/erniealice/fycha-golang/domain/funding/funding/card/list"
	drawform "github.com/erniealice/fycha-golang/domain/funding/funding/draw/form"
	settlementform "github.com/erniealice/fycha-golang/domain/funding/funding/settlement/form"
	sourcedetail "github.com/erniealice/fycha-golang/domain/funding/funding/source/detail"
	sourcelist "github.com/erniealice/fycha-golang/domain/funding/funding/source/list"
	transferform "github.com/erniealice/fycha-golang/domain/funding/funding/transfer/form"
)

// FundingModuleDeps holds all dependencies for the funding view module.
type FundingModuleDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// Labels holds all translatable strings for the funding views.
	// Defaults to DefaultFundingFormLabels() when zero-value.
	Labels FundingFormLabels

	// Route URLs for the 8 funding views. Each defaults to its /app/* constant
	// when empty for backward compatibility during the P3–P11 parallel-mux window.
	SourceListURL     string // e.g. "/app/funding/sources"
	SourceDetailURL   string // e.g. "/app/funding/sources/{fund_id}"
	CardListURL       string // e.g. "/app/funding/cards"
	CardDetailURL     string // e.g. "/app/funding/cards/{allocation_id}"
	AllocationFormURL string // e.g. "/app/funding/allocation/form"
	DrawFormURL       string // e.g. "/app/funding/draw/form"
	SettlementFormURL string // e.g. "/app/funding/settlement/form"
	TransferFormURL   string // e.g. "/app/funding/transfer/form"

	// Fund use cases
	ReadFund  func(ctx context.Context, req *fundpb.ReadFundRequest) (*fundpb.ReadFundResponse, error)
	ListFunds func(ctx context.Context, req *fundpb.ListFundsRequest) (*fundpb.ListFundsResponse, error)

	// FundAllocation use cases
	ReadAllocation  func(ctx context.Context, req *fundallocationpb.ReadFundAllocationRequest) (*fundallocationpb.ReadFundAllocationResponse, error)
	ListAllocations func(ctx context.Context, req *fundallocationpb.ListFundAllocationsRequest) (*fundallocationpb.ListFundAllocationsResponse, error)

	// FundTransaction use cases
	ListTransactions func(ctx context.Context, req *fundtransactionpb.ListFundTransactionsRequest) (*fundtransactionpb.ListFundTransactionsResponse, error)
}

// FundingModule holds all constructed funding views.
type FundingModule struct {
	SourceList     view.View
	SourceDetail   view.View
	CardList       view.View
	CardDetail     view.View
	AllocationForm view.View
	DrawForm       view.View
	SettlementForm view.View
	TransferForm   view.View
	// Route URLs (set from FundingModuleDeps; used by RegisterRoutes).
	sourceListURL     string
	sourceDetailURL   string
	cardListURL       string
	cardDetailURL     string
	allocationFormURL string
	drawFormURL       string
	settlementFormURL string
	transferFormURL   string
}

// NewFundingModule constructs a fully wired funding view module.
func NewFundingModule(deps *FundingModuleDeps) *FundingModule {
	if deps == nil {
		deps = &FundingModuleDeps{}
	}
	// Apply defaults so callers that don't set Labels still get valid text.
	lbls := deps.Labels
	if lbls.Source.Title == "" {
		lbls = DefaultFundingFormLabels()
	}
	return &FundingModule{
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
			Labels:       lbls,
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
			Labels:       lbls,
		}),
		DrawForm: drawform.NewView(&drawform.Deps{
			CommonLabels: deps.CommonLabels,
			Labels:       lbls,
		}),
		SettlementForm: settlementform.NewView(&settlementform.Deps{
			CommonLabels: deps.CommonLabels,
			Labels:       lbls,
		}),
		TransferForm: transferform.NewView(&transferform.Deps{
			CommonLabels: deps.CommonLabels,
			Labels:       lbls,
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
// Route paths are read from the FundingModuleDeps URL fields (set during construction);
// each falls back to its /app/* default when empty.
func (m *FundingModule) RegisterRoutes(r view.RouteRegistrar) {
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
