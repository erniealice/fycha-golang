package equity

// ---------------------------------------------------------------------------
// Equity labels (Funding > Equity app)
// ---------------------------------------------------------------------------

// Labels is the top-level label container for the Equity app.
type Labels struct {
	Accounts     AccountLabels     `json:"accounts"`
	Transactions TransactionLabels `json:"transactions"`
	Sheet        SheetLabels       `json:"sheet"`
	Dashboard    DashboardLabels   `json:"dashboard"`
}

// DashboardLabels holds translatable strings for the Equity live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type DashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	TotalContributed string `json:"totalContributed"`
	ActiveOwners     string `json:"activeOwners"`
	DistributionsYTD string `json:"distributionsYtd"`
	NetMovementYTD   string `json:"netMovementYtd"`
	// Widgets
	EquityByOwner      string `json:"equityByOwner"`
	TopContributors    string `json:"topContributors"`
	RecentTransactions string `json:"recentTransactions"`
	NoRecentTxns       string `json:"noRecentTxns"`
	// Quick actions
	QuickRecordContribution string `json:"quickRecordContribution"`
	QuickRecordDistribution string `json:"quickRecordDistribution"`
	QuickOwnerStatement     string `json:"quickOwnerStatement"`
	QuickEquityReport       string `json:"quickEquityReport"`
	// Common
	ViewAll    string `json:"viewAll"`
	AxisAmount string `json:"axisAmount"`
}

// SheetLabels holds sheet-form title and button labels for equity pages.
type SheetLabels struct {
	AddCapitalAccount       string `json:"addCapitalAccount"`
	RecordTransaction       string `json:"recordTransaction"`
	RecordEquityTransaction string `json:"recordEquityTransaction"`
	PostTransaction         string `json:"postTransaction"`
}

// AccountLabels holds translatable strings for the capital accounts list.
type AccountLabels struct {
	Page    AccountPageLabels   `json:"page"`
	Buttons AccountButtonLabels `json:"buttons"`
	Columns AccountColumnLabels `json:"columns"`
	Empty   AccountEmptyLabels  `json:"empty"`
	Actions AccountActionLabels `json:"actions"`
}

type AccountPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type AccountButtonLabels struct {
	AddAccount string `json:"addAccount"`
}

type AccountColumnLabels struct {
	Name        string `json:"name"`
	OwnerName   string `json:"ownerName"`
	AccountType string `json:"accountType"`
	Balance     string `json:"balance"`
}

type AccountEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type AccountActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	NoPermission string `json:"noPermission"`
}

// TransactionLabels holds translatable strings for the equity transactions list.
type TransactionLabels struct {
	Page    TransactionPageLabels   `json:"page"`
	Buttons TransactionButtonLabels `json:"buttons"`
	Columns TransactionColumnLabels `json:"columns"`
	Empty   TransactionEmptyLabels  `json:"empty"`
	Actions TransactionActionLabels `json:"actions"`
	Form    TransactionFormLabels   `json:"form"`
}

type TransactionPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type TransactionButtonLabels struct {
	RecordTransaction string `json:"recordTransaction"`
}

type TransactionColumnLabels struct {
	Date            string `json:"date"`
	TransactionType string `json:"transactionType"`
	Amount          string `json:"amount"`
	Description     string `json:"description"`
}

type TransactionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type TransactionActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

type TransactionFormLabels struct {
	TransactionType         string `json:"transactionType"`
	TransactionContribution string `json:"transactionContribution"`
	TransactionWithdrawal   string `json:"transactionWithdrawal"`
	TransactionDistribution string `json:"transactionDistribution"`
	TransactionTransfer     string `json:"transactionTransfer"`
	EquityAccount           string `json:"equityAccount"`
	Amount                  string `json:"amount"`
	TransactionDate         string `json:"transactionDate"`
	Description             string `json:"description"`
	DescriptionPlaceholder  string `json:"descriptionPlaceholder"`
	JournalEntryHint        string `json:"journalEntryHint"`
	MemoPlaceholder         string `json:"memoPlaceholder"`
	JournalEntryNote        string `json:"journalEntryNote"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
func DefaultLabels() Labels {
	return Labels{
		Accounts: AccountLabels{
			Page: AccountPageLabels{
				Heading: "Capital Accounts",
				Caption: "Track owner equity and capital contributions",
			},
			Buttons: AccountButtonLabels{
				AddAccount: "Add Capital Account",
			},
			Columns: AccountColumnLabels{
				Name:        "Account Name",
				OwnerName:   "Owner",
				AccountType: "Type",
				Balance:     "Balance",
			},
			Empty: AccountEmptyLabels{
				Title:   "No capital accounts",
				Message: "Add your first capital account to start tracking owner equity.",
			},
			Actions: AccountActionLabels{
				View:         "View",
				Edit:         "Edit",
				NoPermission: "No permission",
			},
		},
		Transactions: TransactionLabels{
			Page: TransactionPageLabels{
				Heading: "Equity Transactions",
				Caption: "Capital contributions, withdrawals, and distributions",
			},
			Buttons: TransactionButtonLabels{
				RecordTransaction: "Record Transaction",
			},
			Columns: TransactionColumnLabels{
				Date:            "Date",
				TransactionType: "Type",
				Amount:          "Amount",
				Description:     "Description",
			},
			Empty: TransactionEmptyLabels{
				Title:   "No equity transactions",
				Message: "Record your first equity transaction to start tracking capital movements.",
			},
			Actions: TransactionActionLabels{
				View:         "View",
				NoPermission: "No permission",
			},
			Form: TransactionFormLabels{
				TransactionType:         "Transaction Type",
				EquityAccount:           "Capital Account",
				Amount:                  "Amount",
				TransactionDate:         "Transaction Date",
				Description:             "Description",
				DescriptionPlaceholder:  "Optional memo for this transaction",
				MemoPlaceholder:         "Optional memo for this transaction",
				JournalEntryHint:        "The corresponding journal entry will be auto-generated when you post this transaction.",
				JournalEntryNote:        "The corresponding journal entry will be auto-generated when you post this transaction. Debits and credits are determined by the transaction type selected above.",
				TransactionContribution: "Contribution — Owner adds capital",
				TransactionWithdrawal:   "Withdrawal — Owner draws cash",
				TransactionDistribution: "Distribution — Profit distributed",
				TransactionTransfer:     "Transfer — Between equity accounts",
			},
		},
		Sheet: SheetLabels{
			AddCapitalAccount:       "Add Capital Account",
			RecordTransaction:       "Record Transaction",
			RecordEquityTransaction: "Record Equity Transaction",
			PostTransaction:         "Post Transaction",
		},
		Dashboard: DashboardLabels{
			Title:                   "Equity Dashboard",
			Subtitle:                "Owner contributions, distributions, and movements across capital accounts",
			TotalContributed:        "Total Contributed",
			ActiveOwners:            "Active Owners",
			DistributionsYTD:        "Distributions YTD",
			NetMovementYTD:          "Net Movement YTD",
			EquityByOwner:           "Equity by Owner",
			TopContributors:         "Top Contributors",
			RecentTransactions:      "Recent Transactions",
			NoRecentTxns:            "No recent equity transactions",
			QuickRecordContribution: "Record Contribution",
			QuickRecordDistribution: "Record Distribution",
			QuickOwnerStatement:     "Owner Statement",
			QuickEquityReport:       "Equity Report",
			ViewAll:                 "View All",
			AxisAmount:              "Amount",
		},
	}
}
