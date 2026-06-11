package treasury

// ---------------------------------------------------------------------------
// LoanPayment labels (Funding > Loans > Payments)
// ---------------------------------------------------------------------------

// LoanPaymentLabels holds all translatable strings for the loan payments view.
type LoanPaymentLabels struct {
	Page    LoanPaymentPageLabels   `json:"page"`
	Buttons LoanPaymentButtonLabels `json:"buttons"`
	Columns LoanPaymentColumnLabels `json:"columns"`
	Empty   LoanPaymentEmptyLabels  `json:"empty"`
	Actions LoanPaymentActionLabels `json:"actions"`
	Form    LoanPaymentFormLabels   `json:"form"`
	Sheet   LoanPaymentSheetLabels  `json:"sheet"`
}

// LoanPaymentSheetLabels holds sheet-form title and button labels for loan payments page.
type LoanPaymentSheetLabels struct {
	RecordPayment string `json:"recordPayment"`
	PostPayment   string `json:"postPayment"`
}

type LoanPaymentPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type LoanPaymentButtonLabels struct {
	RecordPayment string `json:"recordPayment"`
}

type LoanPaymentColumnLabels struct {
	PaymentNumber    string `json:"paymentNumber"`
	PaymentDate      string `json:"paymentDate"`
	PrincipalAmount  string `json:"principalAmount"`
	InterestAmount   string `json:"interestAmount"`
	FeeAmount        string `json:"feeAmount"`
	TotalAmount      string `json:"totalAmount"`
	RemainingBalance string `json:"remainingBalance"`
}

type LoanPaymentEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type LoanPaymentActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

type LoanPaymentFormLabels struct {
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

// DefaultLoanPaymentLabels returns LoanPaymentLabels with hardcoded English defaults.
func DefaultLoanPaymentLabels() LoanPaymentLabels {
	return LoanPaymentLabels{
		Page: LoanPaymentPageLabels{
			Heading: "Loan Payments",
			Caption: "Payment history for this loan",
		},
		Buttons: LoanPaymentButtonLabels{
			RecordPayment: "Record Payment",
		},
		Columns: LoanPaymentColumnLabels{
			PaymentNumber:    "Payment #",
			PaymentDate:      "Date",
			PrincipalAmount:  "Principal",
			InterestAmount:   "Interest",
			FeeAmount:        "Fees",
			TotalAmount:      "Total",
			RemainingBalance: "Balance",
		},
		Empty: LoanPaymentEmptyLabels{
			Title:   "No payments recorded",
			Message: "Record the first payment against this loan.",
		},
		Actions: LoanPaymentActionLabels{
			View:         "View",
			NoPermission: "No permission",
		},
		Form: LoanPaymentFormLabels{
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
		Sheet: LoanPaymentSheetLabels{
			RecordPayment: "Record Payment",
			PostPayment:   "Post Payment",
		},
	}
}
