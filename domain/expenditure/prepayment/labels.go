package prepayment

// Labels (Expenses — Prepayments)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Prepayments module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Status  StatusLabels `json:"status"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	Heading             string `json:"heading"`
	Caption             string `json:"caption"`
	AmortizationHeading string `json:"amortizationHeading"`
	AmortizationCaption string `json:"amortizationCaption"`
}

type ButtonLabels struct {
	AddPrepayment string `json:"addPrepayment"`
}

type ColumnLabels struct {
	Description        string `json:"description"`
	Vendor             string `json:"vendor"`
	TotalAmount        string `json:"totalAmount"`
	RemainingAmount    string `json:"remainingAmount"`
	AmortizationMonths string `json:"amortizationMonths"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	Status             string `json:"status"`
	// Amortization schedule sub-table
	Month   string `json:"month"`
	Opening string `json:"opening"`
	Expense string `json:"expense"`
	Closing string `json:"closing"`
}

type StatusLabels struct {
	Active    string `json:"active"`
	Amortized string `json:"amortized"`
	Cancelled string `json:"cancelled"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FormLabels struct {
	Description             string `json:"description"`
	DescriptionPlaceholder  string `json:"descriptionPlaceholder"`
	Vendor                  string `json:"vendor"`
	VendorPlaceholder       string `json:"vendorPlaceholder"`
	TotalAmount             string `json:"totalAmount"`
	TotalAmountPlaceholder  string `json:"totalAmountPlaceholder"`
	AmortizationMonths      string `json:"amortizationMonths"`
	AmortizationPlaceholder string `json:"amortizationPlaceholder"`
	StartDate               string `json:"startDate"`
	EndDate                 string `json:"endDate"`
	PrepaidAccount          string `json:"prepaidAccount"`
	ExpenseAccount          string `json:"expenseAccount"`
}

type ActionLabels struct {
	View          string `json:"view"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:             "Prepayments",
			Caption:             "Track prepaid expenses and their amortization schedules",
			AmortizationHeading: "Amortization Schedule",
			AmortizationCaption: "Monthly expense recognition for active prepayments",
		},
		Buttons: ButtonLabels{
			AddPrepayment: "Add Prepayment",
		},
		Columns: ColumnLabels{
			Description:        "Description",
			Vendor:             "Vendor",
			TotalAmount:        "Total Amount",
			RemainingAmount:    "Remaining",
			AmortizationMonths: "Months",
			StartDate:          "Start Date",
			EndDate:            "End Date",
			Status:             "Status",
			Month:              "Month",
			Opening:            "Opening Balance",
			Expense:            "Monthly Expense",
			Closing:            "Closing Balance",
		},
		Status: StatusLabels{
			Active:    "Active",
			Amortized: "Fully Amortized",
			Cancelled: "Cancelled",
		},
		Empty: EmptyLabels{
			Title:   "No prepayments found",
			Message: "Record prepaid expenses such as insurance, rent, and subscriptions paid in advance.",
		},
		Form: FormLabels{
			Description:             "Description",
			DescriptionPlaceholder:  "e.g. Annual insurance premium",
			Vendor:                  "Vendor",
			VendorPlaceholder:       "e.g. ABC Insurance Co.",
			TotalAmount:             "Total Amount",
			TotalAmountPlaceholder:  "0.00",
			AmortizationMonths:      "Amortization Period (Months)",
			AmortizationPlaceholder: "e.g. 12",
			StartDate:               "Start Date",
			EndDate:                 "End Date",
			PrepaidAccount:          "Prepaid Account (Asset)",
			ExpenseAccount:          "Expense Account",
		},
		Actions: ActionLabels{
			View:          "View",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this prepayment? This action cannot be undone.",
		},
	}
}
