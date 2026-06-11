// routes.go defines configurable route structs for the loan_payment entity.
package loan_payment

const (
	// Funding — Loan Payments
	AddURL  = "/action/funding/loan/payment/add"
	ListURL = "/funding/loans/payments/{status}"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Loan Payment views.
type Routes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	AddURL    string `json:"add_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav: "loan",
		ListURL:   ListURL,
		AddURL:    AddURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"loan_payment.list": r.ListURL,
		"loan_payment.add":  r.AddURL,
	}
}
