package expenditure

import (
	"context"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	prepaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/prepayment"
	prepayment "github.com/erniealice/fycha-golang/domain/expenditure/prepayment"
)

// PrepaymentModuleDeps holds all dependencies for the expenses expansion module.
// Use case fields are nil until Phase 4-8 prepayment use cases are implemented in espyna.
type PrepaymentModuleDeps struct {
	// Prepayment use cases
	CreatePrepayment          func(ctx context.Context, req *prepaymentpb.CreatePrepaymentRequest) (*prepaymentpb.CreatePrepaymentResponse, error)
	ReadPrepayment            func(ctx context.Context, req *prepaymentpb.ReadPrepaymentRequest) (*prepaymentpb.ReadPrepaymentResponse, error)
	ListPrepayments           func(ctx context.Context, req *prepaymentpb.ListPrepaymentsRequest) (*prepaymentpb.ListPrepaymentsResponse, error)
	GetPrepaymentListPageData func(ctx context.Context, req *prepaymentpb.GetPrepaymentListPageDataRequest) (*prepaymentpb.GetPrepaymentListPageDataResponse, error)
}

// PrepaymentModule holds all constructed expenses expansion views.
type PrepaymentModule struct {
	deps *PrepaymentModuleDeps
}

// NewPrepaymentModule creates an expenses expansion module.
func NewPrepaymentModule(deps *PrepaymentModuleDeps) *PrepaymentModule {
	if deps == nil {
		deps = &PrepaymentModuleDeps{}
	}
	return &PrepaymentModule{deps: deps}
}

// RegisterRoutes registers all expenses expansion routes with the given route registrar.
// These routes extend the existing Expenses app (active nav: "expenses").
// Routes render "Coming Soon" placeholders until view constructors are wired.
func (m *PrepaymentModule) RegisterRoutes(r view.RouteRegistrar) {
	// Prepayments
	r.GET(prepayment.ListURL, comingSoonView("Prepayments", "expenses", "prepayments"))
	r.GET(prepayment.AmortizationURL, comingSoonView("Amortization Schedule", "expenses", "prepayment-amortization"))
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
