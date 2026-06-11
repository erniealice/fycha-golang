package treasury

// ---------------------------------------------------------------------------
// Loan labels (Funding > Loans app)
// ---------------------------------------------------------------------------

// LoanLabels is the top-level label container for the Loans app.
type LoanLabels struct {
	Page      LoanPageLabels      `json:"page"`
	Tabs      LoanTabLabels       `json:"tabs"`
	Buttons   LoanButtonLabels    `json:"buttons"`
	Columns   LoanColumnLabels    `json:"columns"`
	Empty     LoanEmptyLabels     `json:"empty"`
	Actions   LoanActionLabels    `json:"actions"`
	Form      LoanFormLabels      `json:"form"`
	Status    LoanStatusLabels    `json:"status"`
	Type      LoanTypeLabels      `json:"type"`
	Sheet     LoanSheetLabels     `json:"sheet"`
	Dashboard LoanDashboardLabels `json:"dashboard"`
}

// LoanDashboardLabels holds translatable strings for the Loan live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type LoanDashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	TotalOutstanding string `json:"totalOutstanding"`
	InterestYTD      string `json:"interestYtd"`
	PaymentsDue30    string `json:"paymentsDue30"`
	DefaultedCount   string `json:"defaultedCount"`
	// Widgets
	OutstandingTrend string `json:"outstandingTrend"`
	TopLoans         string `json:"topLoans"`
	RecentPayments   string `json:"recentPayments"`
	NoRecentPayments string `json:"noRecentPayments"`
	// Quick actions
	QuickNewLoan      string `json:"quickNewLoan"`
	QuickRecordPay    string `json:"quickRecordPay"`
	QuickAmortization string `json:"quickAmortization"`
	QuickLoanCalendar string `json:"quickLoanCalendar"`
	// Common
	ViewAll    string `json:"viewAll"`
	AxisAmount string `json:"axisAmount"`
}

// LoanSheetLabels holds sheet-form title and button labels for loan list page.
type LoanSheetLabels struct {
	AddLoan  string `json:"addLoan"`
	SaveLoan string `json:"saveLoan"`
}

type LoanPageLabels struct {
	HeadingActive    string `json:"headingActive"`
	CaptionActive    string `json:"captionActive"`
	HeadingCompleted string `json:"headingCompleted"`
	CaptionCompleted string `json:"captionCompleted"`
}

type LoanTabLabels struct {
	Active    string `json:"active"`
	Completed string `json:"completed"`
}

type LoanButtonLabels struct {
	AddLoan string `json:"addLoan"`
}

type LoanColumnLabels struct {
	LoanNumber       string `json:"loanNumber"`
	LenderName       string `json:"lenderName"`
	LoanType         string `json:"loanType"`
	PrincipalAmount  string `json:"principalAmount"`
	RemainingBalance string `json:"remainingBalance"`
	InterestRate     string `json:"interestRate"`
	TermMonths       string `json:"termMonths"`
	StartDate        string `json:"startDate"`
	MaturityDate     string `json:"maturityDate"`
	Status           string `json:"status"`
}

type LoanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type LoanActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
	SaveError    string `json:"saveError"`
}

type LoanFormLabels struct {
	LoanNumber             string `json:"loanNumber"`
	LoanNumberPlaceholder  string `json:"loanNumberPlaceholder"`
	LenderName             string `json:"lenderName"`
	LenderNamePlaceholder  string `json:"lenderNamePlaceholder"`
	LoanType               string `json:"loanType"`
	PrincipalAmount        string `json:"principalAmount"`
	InterestRate           string `json:"interestRate"`
	TermMonths             string `json:"termMonths"`
	StartDate              string `json:"startDate"`
	MaturityDate           string `json:"maturityDate"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
}

type LoanStatusLabels struct {
	Draft     string `json:"draft"`
	Active    string `json:"active"`
	Completed string `json:"completed"`
	Defaulted string `json:"defaulted"`
}

type LoanTypeLabels struct {
	Payable    string `json:"payable"`
	Receivable string `json:"receivable"`
}

// DefaultLoanLabels returns LoanLabels with hardcoded English defaults.
func DefaultLoanLabels() LoanLabels {
	return LoanLabels{
		Page: LoanPageLabels{
			HeadingActive:    "Active Loans",
			CaptionActive:    "Loans currently being serviced",
			HeadingCompleted: "Completed Loans",
			CaptionCompleted: "Fully paid or closed loans",
		},
		Tabs: LoanTabLabels{
			Active:    "Active",
			Completed: "Completed",
		},
		Buttons: LoanButtonLabels{
			AddLoan: "Add Loan",
		},
		Columns: LoanColumnLabels{
			LoanNumber:       "Loan #",
			LenderName:       "Lender / Borrower",
			LoanType:         "Type",
			PrincipalAmount:  "Principal",
			RemainingBalance: "Balance",
			InterestRate:     "Rate",
			TermMonths:       "Term (mo.)",
			StartDate:        "Start Date",
			MaturityDate:     "Maturity",
			Status:           "Status",
		},
		Empty: LoanEmptyLabels{
			Title:   "No loans found",
			Message: "Add your first loan to start tracking borrowings and repayments.",
		},
		Actions: LoanActionLabels{
			View:         "View",
			NoPermission: "No permission",
			SaveError:    "Failed to save loan",
		},
		Form: LoanFormLabels{
			LoanNumber:             "Loan Number",
			LoanNumberPlaceholder:  "e.g. LN-001",
			LenderName:             "Lender / Borrower",
			LenderNamePlaceholder:  "Name of the lender or borrower",
			LoanType:               "Loan Type",
			PrincipalAmount:        "Principal Amount",
			InterestRate:           "Annual Interest Rate (%)",
			TermMonths:             "Term (Months)",
			StartDate:              "Start Date",
			MaturityDate:           "Maturity Date",
			Description:            "Description",
			DescriptionPlaceholder: "Brief description or purpose of the loan",
		},
		Status: LoanStatusLabels{
			Draft:     "Draft",
			Active:    "Active",
			Completed: "Completed",
			Defaulted: "Defaulted",
		},
		Type: LoanTypeLabels{
			Payable:    "Payable",
			Receivable: "Receivable",
		},
		Sheet: LoanSheetLabels{
			AddLoan:  "Add Loan",
			SaveLoan: "Save Loan",
		},
		Dashboard: LoanDashboardLabels{
			Title:             "Loans Dashboard",
			Subtitle:          "Outstanding balance, interest accrued, upcoming payments, and recent activity",
			TotalOutstanding:  "Total Outstanding",
			InterestYTD:       "Interest YTD",
			PaymentsDue30:     "Payments Due (30d)",
			DefaultedCount:    "Defaulted Loans",
			OutstandingTrend:  "Outstanding Principal Trend",
			TopLoans:          "Top Loans by Outstanding",
			RecentPayments:    "Recent Payments",
			NoRecentPayments:  "No recent payments",
			QuickNewLoan:      "New Loan",
			QuickRecordPay:    "Record Payment",
			QuickAmortization: "Amortization Schedule",
			QuickLoanCalendar: "Loan Calendar",
			ViewAll:           "View All",
			AxisAmount:        "Outstanding",
		},
	}
}
