// ---------------------------------------------------------------------------
// Payroll Run labels
// ---------------------------------------------------------------------------

package run

// Labels holds all translatable strings for the Payroll Run sub-module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Tabs    TabLabels    `json:"tabs"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	HeadingDraft       string `json:"headingDraft"`
	SubtitleDraft      string `json:"subtitleDraft"`
	HeadingCalculated  string `json:"headingCalculated"`
	SubtitleCalculated string `json:"subtitleCalculated"`
	HeadingApproved    string `json:"headingApproved"`
	SubtitleApproved   string `json:"subtitleApproved"`
	HeadingPosted      string `json:"headingPosted"`
	SubtitlePosted     string `json:"subtitlePosted"`
}

type TabLabels struct {
	Draft      string `json:"draft"`
	Calculated string `json:"calculated"`
	Approved   string `json:"approved"`
	Posted     string `json:"posted"`
}

type ButtonLabels struct {
	NewRun string `json:"newRun"`
}

type ColumnLabels struct {
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

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingDraft:       "Draft Payroll Runs",
			SubtitleDraft:      "Payroll runs in preparation — payslips not yet finalized",
			HeadingCalculated:  "Calculated Payroll Runs",
			SubtitleCalculated: "Amounts locked and pending approval",
			HeadingApproved:    "Approved Payroll Runs",
			SubtitleApproved:   "Approved and ready for disbursement",
			HeadingPosted:      "Posted Payroll Runs",
			SubtitlePosted:     "Disbursement completed and journal entry created",
		},
		Tabs: TabLabels{
			Draft:      "Draft",
			Calculated: "Calculated",
			Approved:   "Approved",
			Posted:     "Posted",
		},
		Buttons: ButtonLabels{
			NewRun: "New Payroll Run",
		},
		Columns: ColumnLabels{
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
		Empty: EmptyLabels{
			Title:   "No payroll runs found",
			Message: "Create a new payroll run to start processing employee salaries.",
		},
		Actions: ActionLabels{
			View:         "View",
			NoPermission: "No permission",
		},
	}
}
