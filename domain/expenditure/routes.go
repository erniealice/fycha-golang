package expenditure

const (
	// Expenses — Prepayments
	PrepaymentListURL         = "/expenses/prepayments/{status}"
	PrepaymentAmortizationURL = "/expenses/prepayments/amortization"
)

// PrepaymentRoutes holds route paths for Expenses app prepayment views.
type PrepaymentRoutes struct {
	ActiveNav       string `json:"active_nav"`
	ListURL         string `json:"list_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultPrepaymentRoutes() PrepaymentRoutes {
	return PrepaymentRoutes{
		ActiveNav:       "expense",
		ListURL:         PrepaymentListURL,
		AmortizationURL: PrepaymentAmortizationURL,
	}
}

func (r PrepaymentRoutes) RouteMap() map[string]string {
	return map[string]string{
		"prepayment.list":         r.ListURL,
		"prepayment.amortization": r.AmortizationURL,
	}
}
