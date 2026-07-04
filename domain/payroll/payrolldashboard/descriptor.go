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
				{Key: "payroll-draft", Route: "payroll.run.list", Params: map[string]string{"status": "draft"}, Label: "Draft", Icon: "icon-edit", Permission: "payroll:list", LabelKey: "draft_label", IconKey: "payroll_draft_icon"},
				{Key: "payroll-processed", Route: "payroll.run.list", Params: map[string]string{"status": "processed"}, Label: "Processed", Icon: "icon-check-circle", Permission: "payroll:list", LabelKey: "processed_label", IconKey: "payroll_processed_icon"},
				{Key: "remittances-pending", Route: "payroll.remittance.list", Params: map[string]string{"status": "pending"}, Label: "Pending", Icon: "icon-clock", Permission: "payroll:list", LabelKey: "pending_label", IconKey: "remittances_pending_icon"},
				{Key: "remittances-filed", Route: "payroll.remittance.list", Params: map[string]string{"status": "filed"}, Label: "Filed", Icon: "icon-check-circle", Permission: "payroll:list", LabelKey: "filed_label", IconKey: "remittances_filed_icon"},
				{Key: "payroll-employees", Route: "payroll.employee.list", Label: "Employees", Icon: "icon-users", Permission: "payroll:list", LabelKey: "employees_label", IconKey: "payroll_employees_icon"},
				{Key: "gov-contribution-rates", Route: "payroll.settings.gov_rates", Label: "Gov Contribution Rates", Icon: "icon-settings", Permission: "payroll:list", LabelKey: "gov_rates_label", IconKey: "gov_rates_icon"},
				{Key: "pay-periods", Route: "payroll.settings.pay_periods", Label: "Pay Periods", Icon: "icon-calendar", Permission: "payroll:list", LabelKey: "pay_periods_label", IconKey: "pay_periods_icon"},
			},
		},
	}
}
