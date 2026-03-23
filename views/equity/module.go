package equity

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	capitalaccounts "github.com/erniealice/fycha-golang/views/equity/capitalaccounts"
	equitytransactions "github.com/erniealice/fycha-golang/views/equity/equitytransactions"

	equityaccountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_account"
	equitytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_transaction"
)

// ModuleDeps holds all dependencies for the equity module.
type ModuleDeps struct {
	Routes       fycha.EquityRoutes
	Labels       fycha.EquityLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// EquityAccount use cases
	CreateEquityAccount          func(ctx context.Context, req *equityaccountpb.CreateEquityAccountRequest) (*equityaccountpb.CreateEquityAccountResponse, error)
	ReadEquityAccount            func(ctx context.Context, req *equityaccountpb.ReadEquityAccountRequest) (*equityaccountpb.ReadEquityAccountResponse, error)
	ListEquityAccounts           func(ctx context.Context, req *equityaccountpb.ListEquityAccountsRequest) (*equityaccountpb.ListEquityAccountsResponse, error)
	GetEquityAccountListPageData func(ctx context.Context, req *equityaccountpb.GetEquityAccountListPageDataRequest) (*equityaccountpb.GetEquityAccountListPageDataResponse, error)

	// EquityTransaction use cases
	CreateEquityTransaction          func(ctx context.Context, req *equitytransactionpb.CreateEquityTransactionRequest) (*equitytransactionpb.CreateEquityTransactionResponse, error)
	ListEquityTransactions           func(ctx context.Context, req *equitytransactionpb.ListEquityTransactionsRequest) (*equitytransactionpb.ListEquityTransactionsResponse, error)
	GetEquityTransactionListPageData func(ctx context.Context, req *equitytransactionpb.GetEquityTransactionListPageDataRequest) (*equitytransactionpb.GetEquityTransactionListPageDataResponse, error)
}

// Module holds all constructed equity views.
type Module struct {
	CapitalAccounts    view.View
	EquityTransactions view.View
}

// NewModule creates an equity module with real view constructors.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}

	accountDeps := &capitalaccounts.Deps{
		Routes:                       deps.Routes,
		Labels:                       deps.Labels,
		CommonLabels:                 deps.CommonLabels,
		TableLabels:                  deps.TableLabels,
		GetEquityAccountListPageData: deps.GetEquityAccountListPageData,
		ListEquityAccounts:           deps.ListEquityAccounts,
	}

	txnDeps := &equitytransactions.Deps{
		Routes:                 deps.Routes,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		ListEquityTransactions: deps.ListEquityTransactions,
	}

	return &Module{
		CapitalAccounts:    capitalaccounts.NewView(accountDeps),
		EquityTransactions: equitytransactions.NewView(txnDeps),
	}
}

// RegisterRoutes registers all equity routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(fycha.EquityAccountsURL, m.CapitalAccounts)
	r.GET(fycha.EquityTransactionsURL, m.EquityTransactions)
}
