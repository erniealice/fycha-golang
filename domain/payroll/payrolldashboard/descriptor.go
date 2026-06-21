package payrolldashboard

import "github.com/erniealice/espyna-golang/consumer/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "payroll.payrolldashboard",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "payroll"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "payroll.json", Key: "dashboard"},
		LabelName: "PayrollDashboardLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "payroll:list",
			AppEntry: &compose.AppEntry{
				Key:        "payroll",
				Route:      "payroll.dashboard",
				Label:      "Payroll",
				Icon:       "icon-users",
				Permission: "payroll:list",
			},
			Items: []compose.NavItem{
				{Key: "payroll-draft", Route: "payroll.run.list", Params: map[string]string{"status": "draft"}, Label: "Draft", Icon: "icon-edit", Permission: "payroll:list"},
				{Key: "payroll-processed", Route: "payroll.run.list", Params: map[string]string{"status": "processed"}, Label: "Processed", Icon: "icon-check-circle", Permission: "payroll:list"},
				{Key: "remittances-pending", Route: "payroll.remittance.list", Params: map[string]string{"status": "pending"}, Label: "Pending", Icon: "icon-clock", Permission: "payroll:list"},
				{Key: "remittances-filed", Route: "payroll.remittance.list", Params: map[string]string{"status": "filed"}, Label: "Filed", Icon: "icon-check-circle", Permission: "payroll:list"},
				{Key: "payroll-employees", Route: "payroll.employee.list", Label: "Employees", Icon: "icon-users", Permission: "payroll:list"},
				{Key: "gov-contribution-rates", Route: "payroll.settings.gov_rates", Label: "Gov Contribution Rates", Icon: "icon-settings", Permission: "payroll:list"},
				{Key: "pay-periods", Route: "payroll.settings.pay_periods", Label: "Pay Periods", Icon: "icon-calendar", Permission: "payroll:list"},
			},
		},
	}
}
