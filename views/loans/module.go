package loans

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	loanlist "github.com/erniealice/fycha-golang/views/loans/loanlist"
	loanpayments "github.com/erniealice/fycha-golang/views/loans/loanpayments"

	loanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan"
	loanpaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan_payment"
)

// ModuleDeps holds all dependencies for the loans module.
type ModuleDeps struct {
	// Routes
	Routes        fycha.LoanRoutes
	PaymentRoutes fycha.LoanPaymentRoutes

	// Labels
	Labels        fycha.LoanLabels
	PaymentLabels fycha.LoanPaymentLabels
	CommonLabels  pyeza.CommonLabels
	TableLabels   types.TableLabels

	// Loan use cases
	CreateLoan          func(ctx context.Context, req *loanpb.CreateLoanRequest) (*loanpb.CreateLoanResponse, error)
	ReadLoan            func(ctx context.Context, req *loanpb.ReadLoanRequest) (*loanpb.ReadLoanResponse, error)
	ListLoans           func(ctx context.Context, req *loanpb.ListLoansRequest) (*loanpb.ListLoansResponse, error)
	GetLoanListPageData func(ctx context.Context, req *loanpb.GetLoanListPageDataRequest) (*loanpb.GetLoanListPageDataResponse, error)

	// LoanPayment use cases
	CreateLoanPayment func(ctx context.Context, req *loanpaymentpb.CreateLoanPaymentRequest) (*loanpaymentpb.CreateLoanPaymentResponse, error)
	ListLoanPayments  func(ctx context.Context, req *loanpaymentpb.ListLoanPaymentsRequest) (*loanpaymentpb.ListLoanPaymentsResponse, error)
}

// Module holds all constructed loans views.
type Module struct {
	LoanList     view.View
	LoanPayments view.View
}

// NewModule creates a loans module with LoanList and LoanPayments views wired.
// Amortization remains a coming-soon placeholder.
func NewModule(deps *ModuleDeps) *Module {
	if deps == nil {
		deps = &ModuleDeps{}
	}

	listDeps := &loanlist.Deps{
		Routes:              deps.Routes,
		Labels:              deps.Labels,
		CommonLabels:        deps.CommonLabels,
		TableLabels:         deps.TableLabels,
		GetLoanListPageData: deps.GetLoanListPageData,
		ListLoans:           deps.ListLoans,
	}

	paymentDeps := &loanpayments.Deps{
		Routes:           deps.PaymentRoutes,
		Labels:           deps.PaymentLabels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		ListLoanPayments: deps.ListLoanPayments,
	}

	return &Module{
		LoanList:     loanlist.NewView(listDeps),
		LoanPayments: loanpayments.NewView(paymentDeps),
	}
}

// RegisterRoutes registers all loans routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(fycha.LoanListURL, m.LoanList)
	r.GET(fycha.LoanPaymentListURL, m.LoanPayments)
	r.GET(fycha.LoanAmortizationURL, comingSoonView("Amortization Schedules", "loans", "amortization"))
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
