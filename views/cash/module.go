package cash

import (
	"context"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	securitydepositpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/security_deposit"
	pettycashfundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/petty_cash_fund"
)

// ModuleDeps holds all dependencies for the cash expansion module.
// Use case fields are nil until Phase 4-8 treasury use cases are implemented in espyna.
type ModuleDeps struct {
	// SecurityDeposit use cases
	CreateSecurityDeposit func(ctx context.Context, req *securitydepositpb.CreateSecurityDepositRequest) (*securitydepositpb.CreateSecurityDepositResponse, error)
	ListSecurityDeposits  func(ctx context.Context, req *securitydepositpb.ListSecurityDepositsRequest) (*securitydepositpb.ListSecurityDepositsResponse, error)

	// PettyCashFund use cases
	CreatePettyCashFund func(ctx context.Context, req *pettycashfundpb.CreatePettyCashFundRequest) (*pettycashfundpb.CreatePettyCashFundResponse, error)
	ListPettyCashFunds  func(ctx context.Context, req *pettycashfundpb.ListPettyCashFundsRequest) (*pettycashfundpb.ListPettyCashFundsResponse, error)
}

// Module holds all constructed cash expansion views.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates a cash expansion module.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}
	return &Module{deps: deps}
}

// RegisterRoutes registers all cash expansion routes with the given route registrar.
// These routes extend the existing Cash app (active nav: "cash").
// Routes render "Coming Soon" placeholders until view constructors are wired.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Deposits
	r.GET(fycha.DepositListURL, comingSoonView("Deposits", "cash", "deposits"))

	// Petty Cash
	r.GET(fycha.PettyCashRegisterURL, comingSoonView("Petty Cash Register", "cash", "petty-cash-register"))
	r.GET(fycha.PettyCashReplenishmentListURL, comingSoonView("Replenishments", "cash", "replenishments"))
	r.GET(fycha.PettyCashCustodianBalancesURL, comingSoonView("Custodian Balances", "cash", "custodian-balances"))
}

// comingSoonView returns a placeholder view that renders a "Coming Soon" page.
func comingSoonView(title, activeNav, activeSubNav string) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		templateName := "coming-soon"
		if viewCtx.IsHTMX {
			templateName = "coming-soon-content"
		}
		return view.OK(templateName, &types.PageData{
			CacheVersion: viewCtx.CacheVersion,
			Title:        title,
			CurrentPath:  viewCtx.CurrentPath,
			ActiveNav:    activeNav,
			ActiveSubNav: activeSubNav,
			HeaderTitle:  title,
		})
	})
}
