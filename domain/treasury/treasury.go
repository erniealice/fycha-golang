package treasury

// treasury.go — treasury-domain facade.
//
// Re-exports every treasury entity package's Routes/Labels with the entity
// prefix so consumers (block/) keep writing treasury.DefaultLoanRoutes(),
// treasury.DefaultWithholdingCertificateLabels(), etc.

import (
	"context"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	loan "github.com/erniealice/fycha-golang/domain/treasury/loan"
	loanpayment "github.com/erniealice/fycha-golang/domain/treasury/loan_payment"
	pettycash "github.com/erniealice/fycha-golang/domain/treasury/petty_cash"
	withholdingcertificate "github.com/erniealice/fycha-golang/domain/treasury/withholding_certificate"
)

// ---------------------------------------------------------------------------
// Loan (domain/treasury/loan)
// ---------------------------------------------------------------------------

type LoanRoutes = loan.Routes
type LoanLabels = loan.Labels

func DefaultLoanRoutes() LoanRoutes { return loan.DefaultRoutes() }
func DefaultLoanLabels() LoanLabels { return loan.DefaultLabels() }

// ---------------------------------------------------------------------------
// LoanPayment (domain/treasury/loan_payment)
// ---------------------------------------------------------------------------

type LoanPaymentRoutes = loanpayment.Routes
type LoanPaymentLabels = loanpayment.Labels

func DefaultLoanPaymentRoutes() LoanPaymentRoutes { return loanpayment.DefaultRoutes() }
func DefaultLoanPaymentLabels() LoanPaymentLabels { return loanpayment.DefaultLabels() }

// ---------------------------------------------------------------------------
// PettyCash (domain/treasury/petty_cash)
// ---------------------------------------------------------------------------

type PettyCashRoutes = pettycash.Routes
type PettyCashLabels = pettycash.Labels

func DefaultPettyCashRoutes() PettyCashRoutes { return pettycash.DefaultRoutes() }
func DefaultPettyCashLabels() PettyCashLabels { return pettycash.DefaultLabels() }

// ---------------------------------------------------------------------------
// WithholdingCertificate (domain/treasury/withholding_certificate)
// ---------------------------------------------------------------------------

type WithholdingCertificateRoutes = withholdingcertificate.Routes
type WithholdingCertificateLabels = withholdingcertificate.Labels

func DefaultWithholdingCertificateRoutes() WithholdingCertificateRoutes {
	return withholdingcertificate.DefaultRoutes()
}
func DefaultWithholdingCertificateLabels() WithholdingCertificateLabels {
	return withholdingcertificate.DefaultLabels()
}

// ---------------------------------------------------------------------------
// comingSoonView — shared placeholder view used by hoisted treasury modules.
// ---------------------------------------------------------------------------

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
