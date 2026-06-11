package prepayment

const (
	// Expenses — Prepayments
	ListURL         = "/expenses/prepayments/{status}"
	AmortizationURL = "/expenses/prepayments/amortization"
)

// Routes holds route paths for Expenses app prepayment views.
type Routes struct {
	ActiveNav       string `json:"active_nav"`
	ListURL         string `json:"list_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:       "expense",
		ListURL:         ListURL,
		AmortizationURL: AmortizationURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"prepayment.list":         r.ListURL,
		"prepayment.amortization": r.AmortizationURL,
	}
}
