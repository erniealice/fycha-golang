package ledger

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	equity "github.com/erniealice/fycha-golang/domain/ledger/equity"
	capitalaccounts "github.com/erniealice/fycha-golang/domain/ledger/equity/capitalaccounts"
	dashboardview "github.com/erniealice/fycha-golang/domain/ledger/equity/dashboard"
	equitytransactions "github.com/erniealice/fycha-golang/domain/ledger/equity/equitytransactions"

	equityaccountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_account"
	equitytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_transaction"
)

// EquityModuleDeps holds all dependencies for the equity module.
type EquityModuleDeps struct {
	Routes       equity.Routes
	Labels       equity.Labels
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

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	GetEquityDashboardPageData func(ctx context.Context, req *dashboardview.Request) (*dashboardview.Response, error)
}

// EquityModule holds all constructed equity views.
type EquityModule struct {
	CapitalAccounts    view.View
	EquityTransactions view.View

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	Dashboard view.View
	routes    equity.Routes
}

// NewEquityModule creates an equity module with real view constructors.
func NewEquityModule(deps *EquityModuleDeps) *EquityModule {
	if deps == nil {
		deps = &EquityModuleDeps{}
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

	dashDeps := &dashboardview.Deps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		GetDashboardPageData: deps.GetEquityDashboardPageData,
	}

	return &EquityModule{
		CapitalAccounts:    capitalaccounts.NewView(accountDeps),
		EquityTransactions: equitytransactions.NewView(txnDeps),
		Dashboard:          dashboardview.NewView(dashDeps),
		routes:             deps.Routes,
	}
}

// RegisterRoutes registers all equity routes with the given route registrar.
func (m *EquityModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.Dashboard != nil && m.routes.DashboardURL != "" {
		r.GET(m.routes.DashboardURL, m.Dashboard)
	}
	r.GET(equity.AccountsURL, m.CapitalAccounts)
	r.GET(equity.TransactionsURL, m.EquityTransactions)
}
