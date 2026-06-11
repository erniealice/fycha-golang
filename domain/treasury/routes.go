// routes.go defines configurable route structs for fycha treasury domain views.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultXxxRoutes()
//     constructors return structs populated from the package-level route constants
//     defined in routes.go. These serve as sensible defaults for any consumer app.
//   - Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//     Apps can load a JSON config file that maps route keys to custom paths,
//     allowing industry templates (e.g. salon, retail) to rebrand URLs without
//     code changes. The json struct tags on each field support this workflow.
//   - Level 3: App-specific overrides via Go field assignment (optional).
//     After constructing defaults (and optionally applying JSON), consumer apps
//     can directly assign individual struct fields for one-off customizations.
//
// RouteMap() methods return a map[string]string of dot-notation keys to route
// paths, useful for template rendering and route resolution at runtime.
package treasury

const (
	// Funding — Loans
	LoanDashboardURL    = "/funding/loans/dashboard"
	LoanListURL         = "/funding/loans/list/{status}"
	LoanDetailURL       = "/funding/loans/detail/{id}"
	LoanAddURL          = "/action/funding/loan/add"
	LoanAmortizationURL = "/funding/loans/amortization"
	LoanPaymentAddURL   = "/action/funding/loan/payment/add"
	LoanPaymentListURL  = "/funding/loans/payments/{status}"

	// Cash — Petty Cash
	PettyCashRegisterURL          = "/cash/petty-cash/register"
	PettyCashReplenishmentListURL = "/cash/petty-cash/replenishments/{status}"
	PettyCashCustodianBalancesURL = "/cash/petty-cash/custodian-balances"

	// Treasury — Withholding Certificates (full CRUD)
	WithholdingCertificateListURL       = "/treasury/withholding-certificates/list/{status}"
	WithholdingCertificateDetailURL     = "/treasury/withholding-certificates/detail/{id}"
	WithholdingCertificateTableURL      = "/action/withholding-certificate/table/{status}"
	WithholdingCertificateAddURL        = "/action/withholding-certificate/add"
	WithholdingCertificateEditURL       = "/action/withholding-certificate/edit/{id}"
	WithholdingCertificateDeleteURL     = "/action/withholding-certificate/delete"
	WithholdingCertificateBulkDeleteURL = "/action/withholding-certificate/bulk-delete"
	WithholdingCertificateSetStatusURL  = "/action/withholding-certificate/set-status"
)

// ---------------------------------------------------------------------------
// LoanRoutes
// ---------------------------------------------------------------------------

// LoanRoutes holds route paths for Loan views.
type LoanRoutes struct {
	ActiveNav       string `json:"active_nav"`
	DashboardURL    string `json:"dashboard_url"`
	ListURL         string `json:"list_url"`
	DetailURL       string `json:"detail_url"`
	AddURL          string `json:"add_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultLoanRoutes() LoanRoutes {
	return LoanRoutes{
		ActiveNav:       "loan",
		DashboardURL:    LoanDashboardURL,
		ListURL:         LoanListURL,
		DetailURL:       LoanDetailURL,
		AddURL:          LoanAddURL,
		AmortizationURL: LoanAmortizationURL,
	}
}

func (r LoanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"loan.dashboard":    r.DashboardURL,
		"loan.list":         r.ListURL,
		"loan.detail":       r.DetailURL,
		"loan.add":          r.AddURL,
		"loan.amortization": r.AmortizationURL,
	}
}

// ---------------------------------------------------------------------------
// LoanPaymentRoutes
// ---------------------------------------------------------------------------

// LoanPaymentRoutes holds route paths for Loan Payment views.
type LoanPaymentRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	AddURL    string `json:"add_url"`
}

func DefaultLoanPaymentRoutes() LoanPaymentRoutes {
	return LoanPaymentRoutes{
		ActiveNav: "loan",
		ListURL:   LoanPaymentListURL,
		AddURL:    LoanPaymentAddURL,
	}
}

func (r LoanPaymentRoutes) RouteMap() map[string]string {
	return map[string]string{
		"loan_payment.list": r.ListURL,
		"loan_payment.add":  r.AddURL,
	}
}

// ---------------------------------------------------------------------------
// PettyCashRoutes
// ---------------------------------------------------------------------------

// PettyCashRoutes holds route paths for Cash app petty cash views.
type PettyCashRoutes struct {
	ActiveNav            string `json:"active_nav"`
	RegisterURL          string `json:"register_url"`
	ReplenishmentListURL string `json:"replenishment_list_url"`
	CustodianBalancesURL string `json:"custodian_balances_url"`
}

func DefaultPettyCashRoutes() PettyCashRoutes {
	return PettyCashRoutes{
		ActiveNav:            "cash",
		RegisterURL:          PettyCashRegisterURL,
		ReplenishmentListURL: PettyCashReplenishmentListURL,
		CustodianBalancesURL: PettyCashCustodianBalancesURL,
	}
}

func (r PettyCashRoutes) RouteMap() map[string]string {
	return map[string]string{
		"petty_cash.register":           r.RegisterURL,
		"petty_cash.replenishment_list": r.ReplenishmentListURL,
		"petty_cash.custodian_balances": r.CustodianBalancesURL,
	}
}

// ---------------------------------------------------------------------------
// WithholdingCertificateRoutes
// ---------------------------------------------------------------------------

// WithholdingCertificateRoutes holds route paths for Withholding Certificate
// CRUD views (Treasury domain — tax integration v1).
type WithholdingCertificateRoutes struct {
	ActiveNav     string `json:"active_nav"`
	ListURL       string `json:"list_url"`
	DetailURL     string `json:"detail_url"`
	TableURL      string `json:"table_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	SetStatusURL  string `json:"set_status_url"`
}

// DefaultWithholdingCertificateRoutes returns a WithholdingCertificateRoutes
// populated from the package-level route constants.
func DefaultWithholdingCertificateRoutes() WithholdingCertificateRoutes {
	return WithholdingCertificateRoutes{
		ActiveNav:     "withholding_certificate",
		ListURL:       WithholdingCertificateListURL,
		DetailURL:     WithholdingCertificateDetailURL,
		TableURL:      WithholdingCertificateTableURL,
		AddURL:        WithholdingCertificateAddURL,
		EditURL:       WithholdingCertificateEditURL,
		DeleteURL:     WithholdingCertificateDeleteURL,
		BulkDeleteURL: WithholdingCertificateBulkDeleteURL,
		SetStatusURL:  WithholdingCertificateSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r WithholdingCertificateRoutes) RouteMap() map[string]string {
	return map[string]string{
		"withholding_certificate.list":        r.ListURL,
		"withholding_certificate.detail":      r.DetailURL,
		"withholding_certificate.table":       r.TableURL,
		"withholding_certificate.add":         r.AddURL,
		"withholding_certificate.edit":        r.EditURL,
		"withholding_certificate.delete":      r.DeleteURL,
		"withholding_certificate.bulk_delete": r.BulkDeleteURL,
		"withholding_certificate.set_status":  r.SetStatusURL,
	}
}

// DetailFor returns the resolved detail URL for a given withholding certificate ID.
func (r WithholdingCertificateRoutes) DetailFor(id string) string {
	return resolveParam(r.DetailURL, "id", id)
}

// EditFor returns the resolved edit drawer URL for a given withholding certificate ID.
func (r WithholdingCertificateRoutes) EditFor(id string) string {
	return resolveParam(r.EditURL, "id", id)
}

// ---------------------------------------------------------------------------
// resolveParam — internal URL template helper
// ---------------------------------------------------------------------------

// resolveParam replaces a single {placeholder} in a URL pattern with value.
// It is the internal single-parameter URL resolver; for multi-parameter URLs
// use route.ResolveURL from packages/pyeza-golang/route directly.
func resolveParam(pattern, placeholder, value string) string {
	token := "{" + placeholder + "}"
	n := len(token)
	for i := 0; i+n <= len(pattern); i++ {
		if pattern[i:i+n] == token {
			return pattern[:i] + value + pattern[i+n:]
		}
	}
	return pattern
}
