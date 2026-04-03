package expenses

import (
	"context"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	prepaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/prepayment"
	fycha "github.com/erniealice/fycha-golang"
)

// ModuleDeps holds all dependencies for the expenses expansion module.
// Use case fields are nil until Phase 4-8 prepayment use cases are implemented in espyna.
type ModuleDeps struct {
	// Prepayment use cases
	CreatePrepayment          func(ctx context.Context, req *prepaymentpb.CreatePrepaymentRequest) (*prepaymentpb.CreatePrepaymentResponse, error)
	ReadPrepayment            func(ctx context.Context, req *prepaymentpb.ReadPrepaymentRequest) (*prepaymentpb.ReadPrepaymentResponse, error)
	ListPrepayments           func(ctx context.Context, req *prepaymentpb.ListPrepaymentsRequest) (*prepaymentpb.ListPrepaymentsResponse, error)
	GetPrepaymentListPageData func(ctx context.Context, req *prepaymentpb.GetPrepaymentListPageDataRequest) (*prepaymentpb.GetPrepaymentListPageDataResponse, error)
}

// Module holds all constructed expenses expansion views.
type Module struct {
	deps *ModuleDeps
}

// NewModule creates an expenses expansion module.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}
	return &Module{deps: deps}
}

// RegisterRoutes registers all expenses expansion routes with the given route registrar.
// These routes extend the existing Expenses app (active nav: "expenses").
// Routes render "Coming Soon" placeholders until view constructors are wired.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Prepayments
	r.GET(fycha.PrepaymentListURL, comingSoonView("Prepayments", "expenses", "prepayments"))
	r.GET(fycha.PrepaymentAmortizationURL, comingSoonView("Amortization Schedule", "expenses", "prepayment-amortization"))
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
