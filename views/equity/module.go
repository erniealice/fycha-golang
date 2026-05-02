package equity

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	capitalaccounts "github.com/erniealice/fycha-golang/views/equity/capitalaccounts"
	dashboardview "github.com/erniealice/fycha-golang/views/equity/dashboard"
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

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	GetEquityDashboardPageData func(ctx context.Context, req *dashboardview.Request) (*dashboardview.Response, error)
}

// Module holds all constructed equity views.
type Module struct {
	CapitalAccounts    view.View
	EquityTransactions view.View

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	Dashboard view.View
	routes    fycha.EquityRoutes
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

	dashDeps := &dashboardview.Deps{
		Routes:               deps.Routes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		GetDashboardPageData: deps.GetEquityDashboardPageData,
	}

	return &Module{
		CapitalAccounts:    capitalaccounts.NewView(accountDeps),
		EquityTransactions: equitytransactions.NewView(txnDeps),
		Dashboard:          dashboardview.NewView(dashDeps),
		routes:             deps.Routes,
	}
}

// RegisterRoutes registers all equity routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	if m.Dashboard != nil && m.routes.DashboardURL != "" {
		r.GET(m.routes.DashboardURL, m.Dashboard)
	}
	r.GET(fycha.EquityAccountsURL, m.CapitalAccounts)
	r.GET(fycha.EquityTransactionsURL, m.EquityTransactions)
}
