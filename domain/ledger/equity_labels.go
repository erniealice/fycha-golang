package ledger

// ---------------------------------------------------------------------------
// Equity labels (Funding > Equity app)
// ---------------------------------------------------------------------------

// EquityLabels is the top-level label container for the Equity app.
type EquityLabels struct {
	Accounts     EquityAccountLabels     `json:"accounts"`
	Transactions EquityTransactionLabels `json:"transactions"`
	Sheet        EquitySheetLabels       `json:"sheet"`
	Dashboard    EquityDashboardLabels   `json:"dashboard"`
}

// EquityDashboardLabels holds translatable strings for the Equity live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type EquityDashboardLabels struct {
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

// EquitySheetLabels holds sheet-form title and button labels for equity pages.
type EquitySheetLabels struct {
	AddCapitalAccount       string `json:"addCapitalAccount"`
	RecordTransaction       string `json:"recordTransaction"`
	RecordEquityTransaction string `json:"recordEquityTransaction"`
	PostTransaction         string `json:"postTransaction"`
}

// EquityAccountLabels holds translatable strings for the capital accounts list.
type EquityAccountLabels struct {
	Page    EquityAccountPageLabels   `json:"page"`
	Buttons EquityAccountButtonLabels `json:"buttons"`
	Columns EquityAccountColumnLabels `json:"columns"`
	Empty   EquityAccountEmptyLabels  `json:"empty"`
	Actions EquityAccountActionLabels `json:"actions"`
}

type EquityAccountPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type EquityAccountButtonLabels struct {
	AddAccount string `json:"addAccount"`
}

type EquityAccountColumnLabels struct {
	Name        string `json:"name"`
	OwnerName   string `json:"ownerName"`
	AccountType string `json:"accountType"`
	Balance     string `json:"balance"`
}

type EquityAccountEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type EquityAccountActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	NoPermission string `json:"noPermission"`
}

// EquityTransactionLabels holds translatable strings for the equity transactions list.
type EquityTransactionLabels struct {
	Page    EquityTransactionPageLabels   `json:"page"`
	Buttons EquityTransactionButtonLabels `json:"buttons"`
	Columns EquityTransactionColumnLabels `json:"columns"`
	Empty   EquityTransactionEmptyLabels  `json:"empty"`
	Actions EquityTransactionActionLabels `json:"actions"`
	Form    EquityTransactionFormLabels   `json:"form"`
}

type EquityTransactionPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type EquityTransactionButtonLabels struct {
	RecordTransaction string `json:"recordTransaction"`
}

type EquityTransactionColumnLabels struct {
	Date            string `json:"date"`
	TransactionType string `json:"transactionType"`
	Amount          string `json:"amount"`
	Description     string `json:"description"`
}

type EquityTransactionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type EquityTransactionActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

type EquityTransactionFormLabels struct {
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

// DefaultEquityLabels returns EquityLabels with hardcoded English defaults.
func DefaultEquityLabels() EquityLabels {
	return EquityLabels{
		Accounts: EquityAccountLabels{
			Page: EquityAccountPageLabels{
				Heading: "Capital Accounts",
				Caption: "Track owner equity and capital contributions",
			},
			Buttons: EquityAccountButtonLabels{
				AddAccount: "Add Capital Account",
			},
			Columns: EquityAccountColumnLabels{
				Name:        "Account Name",
				OwnerName:   "Owner",
				AccountType: "Type",
				Balance:     "Balance",
			},
			Empty: EquityAccountEmptyLabels{
				Title:   "No capital accounts",
				Message: "Add your first capital account to start tracking owner equity.",
			},
			Actions: EquityAccountActionLabels{
				View:         "View",
				Edit:         "Edit",
				NoPermission: "No permission",
			},
		},
		Transactions: EquityTransactionLabels{
			Page: EquityTransactionPageLabels{
				Heading: "Equity Transactions",
				Caption: "Capital contributions, withdrawals, and distributions",
			},
			Buttons: EquityTransactionButtonLabels{
				RecordTransaction: "Record Transaction",
			},
			Columns: EquityTransactionColumnLabels{
				Date:            "Date",
				TransactionType: "Type",
				Amount:          "Amount",
				Description:     "Description",
			},
			Empty: EquityTransactionEmptyLabels{
				Title:   "No equity transactions",
				Message: "Record your first equity transaction to start tracking capital movements.",
			},
			Actions: EquityTransactionActionLabels{
				View:         "View",
				NoPermission: "No permission",
			},
			Form: EquityTransactionFormLabels{
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
		Sheet: EquitySheetLabels{
			AddCapitalAccount:       "Add Capital Account",
			RecordTransaction:       "Record Transaction",
			RecordEquityTransaction: "Record Equity Transaction",
			PostTransaction:         "Post Transaction",
		},
		Dashboard: EquityDashboardLabels{
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
