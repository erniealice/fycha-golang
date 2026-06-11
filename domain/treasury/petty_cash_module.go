package treasury

import (
	"context"

	"github.com/erniealice/pyeza-golang/view"

	pettycashfundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/petty_cash_fund"
	pettycash "github.com/erniealice/fycha-golang/domain/treasury/petty_cash"
)

// PettyCashModuleDeps holds all dependencies for the cash expansion module.
// Use case fields are nil until Phase 4-8 treasury use cases are implemented in espyna.
//
// NOTE: Deposit-related fields were removed 2026-05-17 as part of the Plan B Advance Cash Events
// rollout — the Deposits sidebar section was retired and the mock deposits view orphan-deleted.
// Operators access UNSCHEDULED advances via Cash → Advances → filter. The proto-level
// SecurityDeposit entity + use cases remain in espyna for future integrations.
type PettyCashModuleDeps struct {
	// PettyCashFund use cases
	CreatePettyCashFund func(ctx context.Context, req *pettycashfundpb.CreatePettyCashFundRequest) (*pettycashfundpb.CreatePettyCashFundResponse, error)
	ListPettyCashFunds  func(ctx context.Context, req *pettycashfundpb.ListPettyCashFundsRequest) (*pettycashfundpb.ListPettyCashFundsResponse, error)
}

// PettyCashModule holds all constructed cash expansion views.
type PettyCashModule struct {
	deps *PettyCashModuleDeps
}

// NewPettyCashModule creates a cash expansion module.
func NewPettyCashModule(deps *PettyCashModuleDeps) *PettyCashModule {
	if deps == nil {
		deps = &PettyCashModuleDeps{}
	}
	return &PettyCashModule{deps: deps}
}

// RegisterRoutes registers all cash expansion routes with the given route registrar.
// These routes extend the existing Cash app (active nav: "cash").
// Routes render "Coming Soon" placeholders until view constructors are wired.
func (m *PettyCashModule) RegisterRoutes(r view.RouteRegistrar) {
	// Petty Cash
	r.GET(pettycash.RegisterURL, comingSoonView("Petty Cash Register", "cash", "petty-cash-register"))
	r.GET(pettycash.ReplenishmentListURL, comingSoonView("Replenishments", "cash", "replenishments"))
	r.GET(pettycash.CustodianBalancesURL, comingSoonView("Custodian Balances", "cash", "custodian-balances"))
}
