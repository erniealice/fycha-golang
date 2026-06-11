package treasury

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	loan "github.com/erniealice/fycha-golang/domain/treasury/loan"
	dashboardview "github.com/erniealice/fycha-golang/domain/treasury/loan/dashboard"
	loanlist "github.com/erniealice/fycha-golang/domain/treasury/loan/loanlist"
	loanpayments "github.com/erniealice/fycha-golang/domain/treasury/loan/loanpayments"
	loanpayment "github.com/erniealice/fycha-golang/domain/treasury/loan_payment"

	loanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan"
	loanpaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan_payment"
)

// LoanModuleDeps holds all dependencies for the loans module.
type LoanModuleDeps struct {
	// Routes
	Routes        loan.Routes
	PaymentRoutes loanpayment.Routes

	// Labels
	Labels        loan.Labels
	PaymentLabels loanpayment.Labels
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

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	GetLoanDashboardPageData func(ctx context.Context, req *dashboardview.Request) (*dashboardview.Response, error)
}

// LoanModule holds all constructed loans views.
type LoanModule struct {
	LoanList     view.View
	LoanPayments view.View

	// Phase 2 — Pyeza dashboard block + per-app live dashboards plan.
	Dashboard view.View
	routes    loan.Routes
}

// NewLoanModule creates a loans module with LoanList and LoanPayments views wired.
// Amortization remains a coming-soon placeholder.
func NewLoanModule(deps *LoanModuleDeps) *LoanModule {
	if deps == nil {
		deps = &LoanModuleDeps{}
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

	dashDeps := &dashboardview.Deps{
		Routes:               deps.Routes,
		PaymentRoutes:        deps.PaymentRoutes,
		Labels:               deps.Labels,
		CommonLabels:         deps.CommonLabels,
		GetDashboardPageData: deps.GetLoanDashboardPageData,
	}

	return &LoanModule{
		LoanList:     loanlist.NewView(listDeps),
		LoanPayments: loanpayments.NewView(paymentDeps),
		Dashboard:    dashboardview.NewView(dashDeps),
		routes:       deps.Routes,
	}
}

// RegisterRoutes registers all loans routes with the given route registrar.
func (m *LoanModule) RegisterRoutes(r view.RouteRegistrar) {
	if m.Dashboard != nil && m.routes.DashboardURL != "" {
		r.GET(m.routes.DashboardURL, m.Dashboard)
	}
	r.GET(loan.ListURL, m.LoanList)
	r.GET(loanpayment.ListURL, m.LoanPayments)
	r.GET(loan.AmortizationURL, comingSoonView("Amortization Schedules", "loans", "amortization"))
}
