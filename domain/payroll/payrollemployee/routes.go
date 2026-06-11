// routes.go defines route constants and configurable route structs for the
// payroll employee entity.

package payrollemployee

const (
	ListURL = "/payroll/employees"
)

// Routes holds route paths for Payroll Employee views.
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
		"payroll.employee.list": r.ListURL,
	}
}
