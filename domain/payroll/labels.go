// ---------------------------------------------------------------------------
// Payroll labels
// ---------------------------------------------------------------------------

package payroll

// PayrollLabels holds all translatable strings for the Payroll module.
type PayrollLabels struct {
	Run        PayrollRunLabels        `json:"run"`
	Remittance PayrollRemittanceLabels `json:"remittance"`
	Employee   PayrollEmployeeLabels   `json:"employee"`
	Settings   PayrollSettingsLabels   `json:"settings"`
	Dashboard  PayrollDashboardLabels  `json:"dashboard"`
}

// PayrollDashboardLabels holds translatable strings for the Payroll live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type PayrollDashboardLabels struct {
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

// PayrollRunLabels holds labels for the Payroll Run sub-module.
type PayrollRunLabels struct {
	Page    PayrollRunPageLabels   `json:"page"`
	Tabs    PayrollRunTabLabels    `json:"tabs"`
	Buttons PayrollRunButtonLabels `json:"buttons"`
	Columns PayrollRunColumnLabels `json:"columns"`
	Empty   PayrollRunEmptyLabels  `json:"empty"`
	Actions PayrollRunActionLabels `json:"actions"`
}

type PayrollRunPageLabels struct {
	HeadingDraft       string `json:"headingDraft"`
	SubtitleDraft      string `json:"subtitleDraft"`
	HeadingCalculated  string `json:"headingCalculated"`
	SubtitleCalculated string `json:"subtitleCalculated"`
	HeadingApproved    string `json:"headingApproved"`
	SubtitleApproved   string `json:"subtitleApproved"`
	HeadingPosted      string `json:"headingPosted"`
	SubtitlePosted     string `json:"subtitlePosted"`
}

type PayrollRunTabLabels struct {
	Draft      string `json:"draft"`
	Calculated string `json:"calculated"`
	Approved   string `json:"approved"`
	Posted     string `json:"posted"`
}

type PayrollRunButtonLabels struct {
	NewRun string `json:"newRun"`
}

type PayrollRunColumnLabels struct {
	RunNumber       string `json:"runNumber"`
	PayPeriod       string `json:"payPeriod"`
	Employees       string `json:"employees"`
	TotalGross      string `json:"totalGross"`
	TotalDeductions string `json:"totalDeductions"`
	TotalNet        string `json:"totalNet"`
	Status          string `json:"status"`
	ApprovedBy      string `json:"approvedBy"`
	PostedAt        string `json:"postedAt"`
}

type PayrollRunEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PayrollRunActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// PayrollRemittanceLabels holds labels for the Payroll Remittance sub-module.
type PayrollRemittanceLabels struct {
	Page    PayrollRemittancePageLabels   `json:"page"`
	Tabs    PayrollRemittanceTabLabels    `json:"tabs"`
	Columns PayrollRemittanceColumnLabels `json:"columns"`
	Types   PayrollRemittanceTypeLabels   `json:"types"`
	Empty   PayrollRemittanceEmptyLabels  `json:"empty"`
}

type PayrollRemittancePageLabels struct {
	HeadingPending  string `json:"headingPending"`
	SubtitlePending string `json:"subtitlePending"`
	HeadingFiled    string `json:"headingFiled"`
	SubtitleFiled   string `json:"subtitleFiled"`
	HeadingPaid     string `json:"headingPaid"`
	SubtitlePaid    string `json:"subtitlePaid"`
}

type PayrollRemittanceTabLabels struct {
	Pending string `json:"pending"`
	Filed   string `json:"filed"`
	Paid    string `json:"paid"`
}

type PayrollRemittanceColumnLabels struct {
	RemittanceType  string `json:"remittanceType"`
	Amount          string `json:"amount"`
	DueDate         string `json:"dueDate"`
	Status          string `json:"status"`
	FiledAt         string `json:"filedAt"`
	ReferenceNumber string `json:"referenceNumber"`
}

type PayrollRemittanceTypeLabels struct {
	SSS            string `json:"sss"`
	PhilHealth     string `json:"philHealth"`
	PagIBIG        string `json:"pagIbig"`
	BIRWithholding string `json:"birWithholding"`
}

type PayrollRemittanceEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// PayrollEmployeeLabels holds labels for the Payroll Employee sub-module.
type PayrollEmployeeLabels struct {
	Page         PayrollEmployeePageLabels         `json:"page"`
	Columns      PayrollEmployeeColumnLabels       `json:"columns"`
	Status       PayrollEmployeeStatusLabels       `json:"status"`
	PayFrequency PayrollEmployeePayFrequencyLabels `json:"payFrequency"`
	Empty        PayrollEmployeeEmptyLabels        `json:"empty"`
}

type PayrollEmployeePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type PayrollEmployeeColumnLabels struct {
	Name         string `json:"name"`
	Position     string `json:"position"`
	Department   string `json:"department"`
	BasicSalary  string `json:"basicSalary"`
	PayFrequency string `json:"payFrequency"`
	Status       string `json:"status"`
}

type PayrollEmployeeStatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type PayrollEmployeePayFrequencyLabels struct {
	SemiMonthly string `json:"semiMonthly"`
	Monthly     string `json:"monthly"`
	Weekly      string `json:"weekly"`
}

type PayrollEmployeeEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// PayrollSettingsLabels holds labels for Payroll Settings pages.
type PayrollSettingsLabels struct {
	GovRates   PayrollGovRatesLabels   `json:"govRates"`
	PayPeriods PayrollPayPeriodsLabels `json:"payPeriods"`
}

type PayrollGovRatesLabels struct {
	Page   PayrollGovRatesPageLabels   `json:"page"`
	Agency PayrollGovRatesAgencyLabels `json:"agency"`
}

type PayrollGovRatesPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type PayrollGovRatesAgencyLabels struct {
	SSS            string `json:"sss"`
	PhilHealth     string `json:"philHealth"`
	PagIBIG        string `json:"pagIbig"`
	BIRWithholding string `json:"birWithholding"`
}

type PayrollPayPeriodsLabels struct {
	Page PayrollPayPeriodsPageLabels `json:"page"`
}

type PayrollPayPeriodsPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

// DefaultPayrollLabels returns PayrollLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultPayrollLabels() PayrollLabels {
	return PayrollLabels{
		Run: PayrollRunLabels{
			Page: PayrollRunPageLabels{
				HeadingDraft:       "Draft Payroll Runs",
				SubtitleDraft:      "Payroll runs in preparation — payslips not yet finalized",
				HeadingCalculated:  "Calculated Payroll Runs",
				SubtitleCalculated: "Amounts locked and pending approval",
				HeadingApproved:    "Approved Payroll Runs",
				SubtitleApproved:   "Approved and ready for disbursement",
				HeadingPosted:      "Posted Payroll Runs",
				SubtitlePosted:     "Disbursement completed and journal entry created",
			},
			Tabs: PayrollRunTabLabels{
				Draft:      "Draft",
				Calculated: "Calculated",
				Approved:   "Approved",
				Posted:     "Posted",
			},
			Buttons: PayrollRunButtonLabels{
				NewRun: "New Payroll Run",
			},
			Columns: PayrollRunColumnLabels{
				RunNumber:       "Run #",
				PayPeriod:       "Pay Period",
				Employees:       "Employees",
				TotalGross:      "Total Gross",
				TotalDeductions: "Deductions",
				TotalNet:        "Net Pay",
				Status:          "Status",
				ApprovedBy:      "Approved By",
				PostedAt:        "Posted At",
			},
			Empty: PayrollRunEmptyLabels{
				Title:   "No payroll runs found",
				Message: "Create a new payroll run to start processing employee salaries.",
			},
			Actions: PayrollRunActionLabels{
				View:         "View",
				NoPermission: "No permission",
			},
		},
		Remittance: PayrollRemittanceLabels{
			Page: PayrollRemittancePageLabels{
				HeadingPending:  "Pending Remittances",
				SubtitlePending: "Government contributions due for filing and payment",
				HeadingFiled:    "Filed Remittances",
				SubtitleFiled:   "Remittances filed with the government agency",
				HeadingPaid:     "Paid Remittances",
				SubtitlePaid:    "Remittances confirmed paid to the government agency",
			},
			Tabs: PayrollRemittanceTabLabels{
				Pending: "Pending",
				Filed:   "Filed",
				Paid:    "Paid",
			},
			Columns: PayrollRemittanceColumnLabels{
				RemittanceType:  "Agency",
				Amount:          "Amount",
				DueDate:         "Due Date",
				Status:          "Status",
				FiledAt:         "Filed At",
				ReferenceNumber: "Reference #",
			},
			Types: PayrollRemittanceTypeLabels{
				SSS:            "SSS",
				PhilHealth:     "PhilHealth",
				PagIBIG:        "Pag-IBIG",
				BIRWithholding: "BIR Withholding",
			},
			Empty: PayrollRemittanceEmptyLabels{
				Title:   "No remittances found",
				Message: "Government contribution remittances will appear here once payroll runs are processed.",
			},
		},
		Employee: PayrollEmployeeLabels{
			Page: PayrollEmployeePageLabels{
				Heading: "Payroll Employees",
				Caption: "Manage employees enrolled in payroll",
			},
			Columns: PayrollEmployeeColumnLabels{
				Name:         "Name",
				Position:     "Position",
				Department:   "Department",
				BasicSalary:  "Basic Salary",
				PayFrequency: "Pay Frequency",
				Status:       "Status",
			},
			Status: PayrollEmployeeStatusLabels{
				Active:   "Active",
				Inactive: "Inactive",
			},
			PayFrequency: PayrollEmployeePayFrequencyLabels{
				SemiMonthly: "Semi-Monthly",
				Monthly:     "Monthly",
				Weekly:      "Weekly",
			},
			Empty: PayrollEmployeeEmptyLabels{
				Title:   "No employees found",
				Message: "Add employees to payroll to begin processing salaries.",
			},
		},
		Settings: PayrollSettingsLabels{
			GovRates: PayrollGovRatesLabels{
				Page: PayrollGovRatesPageLabels{
					Heading: "Government Contribution Rates",
					Caption: "Philippine mandatory contribution rates — SSS, PhilHealth, Pag-IBIG, BIR",
				},
				Agency: PayrollGovRatesAgencyLabels{
					SSS:            "SSS (Social Security System)",
					PhilHealth:     "PhilHealth",
					PagIBIG:        "Pag-IBIG (HDMF)",
					BIRWithholding: "BIR Withholding Tax",
				},
			},
			PayPeriods: PayrollPayPeriodsLabels{
				Page: PayrollPayPeriodsPageLabels{
					Heading: "Pay Period Settings",
					Caption: "Configure payroll cut-off dates and pay schedules",
				},
			},
		},
		Dashboard: PayrollDashboardLabels{
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
		},
	}
}
