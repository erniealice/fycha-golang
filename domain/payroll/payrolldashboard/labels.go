// ---------------------------------------------------------------------------
// Payroll Dashboard labels
// ---------------------------------------------------------------------------

package payrolldashboard

// Labels holds translatable strings for the Payroll live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type Labels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	CurrentRunStatus   string `json:"currentRunStatus"`
	EmployeesInCurrent string `json:"employeesInCurrent"`
	TotalGrossMTD      string `json:"totalGrossMtd"`
	RemittancesDue     string `json:"remittancesDue"`
	// Widgets
	GrossPayByMonth     string `json:"grossPayByMonth"`
	RecentRuns          string `json:"recentRuns"`
	UpcomingRemittances string `json:"upcomingRemittances"`
	NoRecentRuns        string `json:"noRecentRuns"`
	NoUpcomingDeadlines string `json:"noUpcomingDeadlines"`
	// Quick actions
	QuickNewRun            string `json:"quickNewRun"`
	QuickProcessRun        string `json:"quickProcessRun"`
	QuickFileRemittance    string `json:"quickFileRemittance"`
	QuickPayPeriodSettings string `json:"quickPayPeriodSettings"`
	// Common
	ViewAll   string `json:"viewAll"`
	AxisGross string `json:"axisGross"`
	NoRunYet  string `json:"noRunYet"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Title:                  "Payroll Dashboard",
		Subtitle:               "Run status, monthly gross-pay trend, and upcoming government remittances",
		CurrentRunStatus:       "Current Run Status",
		EmployeesInCurrent:     "Employees in Current Run",
		TotalGrossMTD:          "Total Gross (MTD)",
		RemittancesDue:         "Remittances Due (30d)",
		GrossPayByMonth:        "Gross Pay by Month",
		RecentRuns:             "Recent Payroll Runs",
		UpcomingRemittances:    "Upcoming Remittance Deadlines",
		NoRecentRuns:           "No recent payroll runs",
		NoUpcomingDeadlines:    "No upcoming remittance deadlines",
		QuickNewRun:            "New Payroll Run",
		QuickProcessRun:        "Process Run",
		QuickFileRemittance:    "File Remittance",
		QuickPayPeriodSettings: "Pay Period Settings",
		ViewAll:                "View All",
		AxisGross:              "Gross",
		NoRunYet:               "No run yet",
	}
}
