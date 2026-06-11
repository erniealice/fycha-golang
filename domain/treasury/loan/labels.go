package loan

// ---------------------------------------------------------------------------
// Loan labels (Funding > Loans app)
// ---------------------------------------------------------------------------

// Labels is the top-level label container for the Loans app.
type Labels struct {
	Page      PageLabels      `json:"page"`
	Tabs      TabLabels       `json:"tabs"`
	Buttons   ButtonLabels    `json:"buttons"`
	Columns   ColumnLabels    `json:"columns"`
	Empty     EmptyLabels     `json:"empty"`
	Actions   ActionLabels    `json:"actions"`
	Form      FormLabels      `json:"form"`
	Status    StatusLabels    `json:"status"`
	Type      TypeLabels      `json:"type"`
	Sheet     SheetLabels     `json:"sheet"`
	Dashboard DashboardLabels `json:"dashboard"`
}

// DashboardLabels holds translatable strings for the Loan live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type DashboardLabels struct {
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

// SheetLabels holds sheet-form title and button labels for loan list page.
type SheetLabels struct {
	AddLoan  string `json:"addLoan"`
	SaveLoan string `json:"saveLoan"`
}

type PageLabels struct {
	HeadingActive    string `json:"headingActive"`
	CaptionActive    string `json:"captionActive"`
	HeadingCompleted string `json:"headingCompleted"`
	CaptionCompleted string `json:"captionCompleted"`
}

type TabLabels struct {
	Active    string `json:"active"`
	Completed string `json:"completed"`
}

type ButtonLabels struct {
	AddLoan string `json:"addLoan"`
}

type ColumnLabels struct {
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

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
	SaveError    string `json:"saveError"`
}

type FormLabels struct {
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

type StatusLabels struct {
	Draft     string `json:"draft"`
	Active    string `json:"active"`
	Completed string `json:"completed"`
	Defaulted string `json:"defaulted"`
}

type TypeLabels struct {
	Payable    string `json:"payable"`
	Receivable string `json:"receivable"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingActive:    "Active Loans",
			CaptionActive:    "Loans currently being serviced",
			HeadingCompleted: "Completed Loans",
			CaptionCompleted: "Fully paid or closed loans",
		},
		Tabs: TabLabels{
			Active:    "Active",
			Completed: "Completed",
		},
		Buttons: ButtonLabels{
			AddLoan: "Add Loan",
		},
		Columns: ColumnLabels{
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
		Empty: EmptyLabels{
			Title:   "No loans found",
			Message: "Add your first loan to start tracking borrowings and repayments.",
		},
		Actions: ActionLabels{
			View:         "View",
			NoPermission: "No permission",
			SaveError:    "Failed to save loan",
		},
		Form: FormLabels{
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
		Status: StatusLabels{
			Draft:     "Draft",
			Active:    "Active",
			Completed: "Completed",
			Defaulted: "Defaulted",
		},
		Type: TypeLabels{
			Payable:    "Payable",
			Receivable: "Receivable",
		},
		Sheet: SheetLabels{
			AddLoan:  "Add Loan",
			SaveLoan: "Save Loan",
		},
		Dashboard: DashboardLabels{
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
