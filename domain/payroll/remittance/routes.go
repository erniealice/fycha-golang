// routes.go defines route constants and configurable route structs for the
// payroll remittance entity.

package remittance

const (
	ListURL = "/payroll/remittances/{status}"
)

// Routes holds route paths for Payroll Remittance views.
type Routes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav: "payroll",
		ListURL:   ListURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.remittance.list": r.ListURL,
	}
}
