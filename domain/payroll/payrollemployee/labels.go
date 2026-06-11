// ---------------------------------------------------------------------------
// Payroll Employee labels
// ---------------------------------------------------------------------------

package payrollemployee

// Labels holds all translatable strings for the Payroll Employee sub-module.
type Labels struct {
	Page         PageLabels         `json:"page"`
	Columns      ColumnLabels       `json:"columns"`
	Status       StatusLabels       `json:"status"`
	PayFrequency PayFrequencyLabels `json:"payFrequency"`
	Empty        EmptyLabels        `json:"empty"`
}

type PageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type ColumnLabels struct {
	Name         string `json:"name"`
	Position     string `json:"position"`
	Department   string `json:"department"`
	BasicSalary  string `json:"basicSalary"`
	PayFrequency string `json:"payFrequency"`
	Status       string `json:"status"`
}

type StatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type PayFrequencyLabels struct {
	SemiMonthly string `json:"semiMonthly"`
	Monthly     string `json:"monthly"`
	Weekly      string `json:"weekly"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading: "Payroll Employees",
			Caption: "Manage employees enrolled in payroll",
		},
		Columns: ColumnLabels{
			Name:         "Name",
			Position:     "Position",
			Department:   "Department",
			BasicSalary:  "Basic Salary",
			PayFrequency: "Pay Frequency",
			Status:       "Status",
		},
		Status: StatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		PayFrequency: PayFrequencyLabels{
			SemiMonthly: "Semi-Monthly",
			Monthly:     "Monthly",
			Weekly:      "Weekly",
		},
		Empty: EmptyLabels{
			Title:   "No employees found",
			Message: "Add employees to payroll to begin processing salaries.",
		},
	}
}
