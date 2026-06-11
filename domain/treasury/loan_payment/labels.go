package loan_payment

// ---------------------------------------------------------------------------
// LoanPayment labels (Funding > Loans > Payments)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the loan payments view.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Empty   EmptyLabels  `json:"empty"`
	Actions ActionLabels `json:"actions"`
	Form    FormLabels   `json:"form"`
	Sheet   SheetLabels  `json:"sheet"`
}

// SheetLabels holds sheet-form title and button labels for loan payments page.
type SheetLabels struct {
	RecordPayment string `json:"recordPayment"`
	PostPayment   string `json:"postPayment"`
}

type PageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type ButtonLabels struct {
	RecordPayment string `json:"recordPayment"`
}

type ColumnLabels struct {
	PaymentNumber    string `json:"paymentNumber"`
	PaymentDate      string `json:"paymentDate"`
	PrincipalAmount  string `json:"principalAmount"`
	InterestAmount   string `json:"interestAmount"`
	FeeAmount        string `json:"feeAmount"`
	TotalAmount      string `json:"totalAmount"`
	RemainingBalance string `json:"remainingBalance"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

type FormLabels struct {
	PaymentNumber            string `json:"paymentNumber"`
	PaymentNumberPlaceholder string `json:"paymentNumberPlaceholder"`
	PaymentDate              string `json:"paymentDate"`
	PrincipalAmount          string `json:"principalAmount"`
	InterestAmount           string `json:"interestAmount"`
	FeeAmount                string `json:"feeAmount"`
	TotalAmount              string `json:"totalAmount"`
	RemainingBalance         string `json:"remainingBalance"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	PaymentBreakdown         string `json:"paymentBreakdown"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading: "Loan Payments",
			Caption: "Payment history for this loan",
		},
		Buttons: ButtonLabels{
			RecordPayment: "Record Payment",
		},
		Columns: ColumnLabels{
			PaymentNumber:    "Payment #",
			PaymentDate:      "Date",
			PrincipalAmount:  "Principal",
			InterestAmount:   "Interest",
			FeeAmount:        "Fees",
			TotalAmount:      "Total",
			RemainingBalance: "Balance",
		},
		Empty: EmptyLabels{
			Title:   "No payments recorded",
			Message: "Record the first payment against this loan.",
		},
		Actions: ActionLabels{
			View:         "View",
			NoPermission: "No permission",
		},
		Form: FormLabels{
			PaymentNumber:            "Payment Number",
			PaymentNumberPlaceholder: "e.g. PAY-001",
			PaymentDate:              "Payment Date",
			PrincipalAmount:          "Principal Amount",
			InterestAmount:           "Interest Amount",
			FeeAmount:                "Fees (PFRS 9)",
			TotalAmount:              "Total Payment",
			RemainingBalance:         "Remaining Balance After Payment",
			Notes:                    "Notes",
			NotesPlaceholder:         "Optional payment notes or reference",
			PaymentBreakdown:         "Payment Breakdown",
		},
		Sheet: SheetLabels{
			RecordPayment: "Record Payment",
			PostPayment:   "Post Payment",
		},
	}
}
