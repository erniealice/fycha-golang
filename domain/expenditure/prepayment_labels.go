package expenditure

// PrepaymentLabels (Expenses — Prepayments)
// ---------------------------------------------------------------------------

// PrepaymentLabels holds all translatable strings for the Prepayments module.
type PrepaymentLabels struct {
	Page    PrepaymentPageLabels   `json:"page"`
	Buttons PrepaymentButtonLabels `json:"buttons"`
	Columns PrepaymentColumnLabels `json:"columns"`
	Status  PrepaymentStatusLabels `json:"status"`
	Empty   PrepaymentEmptyLabels  `json:"empty"`
	Form    PrepaymentFormLabels   `json:"form"`
	Actions PrepaymentActionLabels `json:"actions"`
}

type PrepaymentPageLabels struct {
	Heading             string `json:"heading"`
	Caption             string `json:"caption"`
	AmortizationHeading string `json:"amortizationHeading"`
	AmortizationCaption string `json:"amortizationCaption"`
}

type PrepaymentButtonLabels struct {
	AddPrepayment string `json:"addPrepayment"`
}

type PrepaymentColumnLabels struct {
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

type PrepaymentStatusLabels struct {
	Active    string `json:"active"`
	Amortized string `json:"amortized"`
	Cancelled string `json:"cancelled"`
}

type PrepaymentEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PrepaymentFormLabels struct {
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

type PrepaymentActionLabels struct {
	View          string `json:"view"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultPrepaymentLabels returns PrepaymentLabels with hardcoded English defaults.
func DefaultPrepaymentLabels() PrepaymentLabels {
	return PrepaymentLabels{
		Page: PrepaymentPageLabels{
			Heading:             "Prepayments",
			Caption:             "Track prepaid expenses and their amortization schedules",
			AmortizationHeading: "Amortization Schedule",
			AmortizationCaption: "Monthly expense recognition for active prepayments",
		},
		Buttons: PrepaymentButtonLabels{
			AddPrepayment: "Add Prepayment",
		},
		Columns: PrepaymentColumnLabels{
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
		Status: PrepaymentStatusLabels{
			Active:    "Active",
			Amortized: "Fully Amortized",
			Cancelled: "Cancelled",
		},
		Empty: PrepaymentEmptyLabels{
			Title:   "No prepayments found",
			Message: "Record prepaid expenses such as insurance, rent, and subscriptions paid in advance.",
		},
		Form: PrepaymentFormLabels{
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
		Actions: PrepaymentActionLabels{
			View:          "View",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this prepayment? This action cannot be undone.",
		},
	}
}
