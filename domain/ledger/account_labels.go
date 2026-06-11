package ledger

// ---------------------------------------------------------------------------
// Account labels (Chart of Accounts)
// ---------------------------------------------------------------------------

// AccountLabels holds all translatable strings for the Chart of Accounts module.
type AccountLabels struct {
	Page          AccountPageLabels          `json:"page"`
	Buttons       AccountButtonLabels        `json:"buttons"`
	Columns       AccountColumnLabels        `json:"columns"`
	Tabs          AccountTabLabels           `json:"tabs"`
	Empty         AccountEmptyLabels         `json:"empty"`
	Form          AccountFormLabels          `json:"form"`
	Actions       AccountActionLabels        `json:"actions"`
	Detail        AccountDetailLabels        `json:"detail"`
	Templates     AccountTemplatesLabels     `json:"templates"`
	GeneralLedger AccountGeneralLedgerLabels `json:"generalLedger"`
	BadDebt       BadDebtLabels              `json:"badDebt"`
	Dashboard     LedgerDashboardLabels      `json:"dashboard"`
}

// LedgerDashboardLabels holds translatable strings for the Ledger live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type LedgerDashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	TotalAssets      string `json:"totalAssets"`
	TotalLiabilities string `json:"totalLiabilities"`
	TotalEquity      string `json:"totalEquity"`
	NetIncomeMTD     string `json:"netIncomeMtd"`
	// Widgets
	BalanceByType    string `json:"balanceByType"`
	UnpostedJournals string `json:"unpostedJournals"`
	RecentJournals   string `json:"recentJournals"`
	NoRecentJournals string `json:"noRecentJournals"`
	// Quick actions
	QuickNewJournal    string `json:"quickNewJournal"`
	QuickTrialBalance  string `json:"quickTrialBalance"`
	QuickClosePeriod   string `json:"quickClosePeriod"`
	QuickAccountLookup string `json:"quickAccountLookup"`
	// Account type labels — used as bar-chart axis labels
	AccountTypeAssets      string `json:"accountTypeAssets"`
	AccountTypeLiabilities string `json:"accountTypeLiabilities"`
	AccountTypeEquity      string `json:"accountTypeEquity"`
	AccountTypeRevenue     string `json:"accountTypeRevenue"`
	AccountTypeExpense     string `json:"accountTypeExpense"`
	// Fallback label for journal entries with no entry number
	JournalEntryFallback string `json:"journalEntryFallback"`
	// Common
	ViewAll    string `json:"viewAll"`
	AxisAmount string `json:"axisAmount"`
}

// BadDebtLabels holds translatable strings for the Bad Debt Policy settings page.
type BadDebtLabels struct {
	Title string `json:"title"`
}

// AccountTemplatesLabels holds translatable strings for the Account Templates settings page.
type AccountTemplatesLabels struct {
	PageTitle           string `json:"pageTitle"`
	PageSubtitle        string `json:"pageSubtitle"`
	CurrentAccountCount string `json:"currentAccountCount"`
	ApplyWarning        string `json:"applyWarning"`
	Empty               string `json:"empty"`
	EmptyDesc           string `json:"emptyDesc"`
	AccountsSuffix      string `json:"accountsSuffix"`
	ComingSoon          string `json:"comingSoon"`
	BadgeApplied        string `json:"badgeApplied"`
	BadgeAssets         string `json:"badgeAssets"`
	BadgeLiabilities    string `json:"badgeLiabilities"`
	BadgeEquity         string `json:"badgeEquity"`
	BadgeRevenue        string `json:"badgeRevenue"`
	BadgeExpenses       string `json:"badgeExpenses"`
	Preview             string `json:"preview"`
	AlreadyApplied      string `json:"alreadyApplied"`
	ApplyTemplate       string `json:"applyTemplate"`
	PreviewTitle        string `json:"previewTitle"`
	PreviewDesc         string `json:"previewDesc"`
	ColCode             string `json:"colCode"`
	ColAccountName      string `json:"colAccountName"`
	ColElement          string `json:"colElement"`
	ColClass            string `json:"colClass"`
	ColIsGroup          string `json:"colIsGroup"`
	Yes                 string `json:"yes"`
	SkipNote            string `json:"skipNote"`
}

type AccountGeneralLedgerLabels struct {
	Title                 string `json:"title"`
	Subtitle              string `json:"subtitle"`
	Account               string `json:"account"`
	AccountPlaceholder    string `json:"accountPlaceholder"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	Apply                 string `json:"apply"`
	Clear                 string `json:"clear"`
	Print                 string `json:"print"`
	SelectAccountMessage  string `json:"selectAccountMessage"`
	NoTransactionsMessage string `json:"noTransactionsMessage"`
	DateRangeSeparator    string `json:"dateRangeSeparator"`
	OpeningBalance        string `json:"openingBalance"`
	PeriodDebits          string `json:"periodDebits"`
	PeriodCredits         string `json:"periodCredits"`
	RunningBalance        string `json:"runningBalance"`
	NoTransactionsTitle   string `json:"noTransactionsTitle"`
	NoTransactionsDetail  string `json:"noTransactionsDetail"`
}

type AccountPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type AccountButtonLabels struct {
	AddAccount string `json:"addAccount"`
}

type AccountColumnLabels struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Element        string `json:"element"`
	Classification string `json:"classification"`
	Group          string `json:"group"`
	Balance        string `json:"balance"`
	// Entry sub-table columns
	Date        string `json:"date"`
	EntryNumber string `json:"entryNumber"`
	Description string `json:"description"`
	Debit       string `json:"debit"`
	Credit      string `json:"credit"`
	// Status
	Status string `json:"status"`
}

type AccountTabLabels struct {
	All       string `json:"all"`
	Asset     string `json:"asset"`
	Liability string `json:"liability"`
	Equity    string `json:"equity"`
	Revenue   string `json:"revenue"`
	Expense   string `json:"expense"`
}

type AccountEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type AccountFormLabels struct {
	Code                     string `json:"code"`
	CodePlaceholder          string `json:"codePlaceholder"`
	Name                     string `json:"name"`
	NamePlaceholder          string `json:"namePlaceholder"`
	Element                  string `json:"element"`
	Classification           string `json:"classification"`
	ParentAccount            string `json:"parentAccount"`
	ParentAccountPlaceholder string `json:"parentAccountPlaceholder"`
	Group                    string `json:"group"`
	IsGroup                  string `json:"isGroup"`
	Active                   string `json:"active"`
	Description              string `json:"description"`
	DescriptionPlaceholder   string `json:"descriptionPlaceholder"`
	CashFlowSection          string `json:"cashFlowSection"`
	CashFlowClassification   string `json:"cashFlowClassification"`
	// Element option labels
	ElementAsset     string `json:"elementAsset"`
	ElementLiability string `json:"elementLiability"`
	ElementEquity    string `json:"elementEquity"`
	ElementRevenue   string `json:"elementRevenue"`
	ElementExpense   string `json:"elementExpense"`
	// Class option labels
	ClassCurrentAsset        string `json:"classCurrentAsset"`
	ClassNonCurrentAsset     string `json:"classNonCurrentAsset"`
	ClassCurrentLiability    string `json:"classCurrentLiability"`
	ClassNonCurrentLiability string `json:"classNonCurrentLiability"`
	ClassEquity              string `json:"classEquity"`
	ClassOperatingRevenue    string `json:"classOperatingRevenue"`
	ClassOtherIncome         string `json:"classOtherIncome"`
	ClassCostOfSales         string `json:"classCostOfSales"`
	ClassOperatingExpense    string `json:"classOperatingExpense"`
	ClassFinanceCost         string `json:"classFinanceCost"`
	ClassIncomeTax           string `json:"classIncomeTax"`
	ClassOtherExpense        string `json:"classOtherExpense"`
	// Cash flow option labels
	CashFlowNone      string `json:"cashFlowNone"`
	CashFlowOperating string `json:"cashFlowOperating"`
	CashFlowInvesting string `json:"cashFlowInvesting"`
	CashFlowFinancing string `json:"cashFlowFinancing"`
	// Element group header labels (used in BuildAccountTree and preview)
	GroupAssets      string `json:"groupAssets"`
	GroupLiabilities string `json:"groupLiabilities"`
	GroupEquity      string `json:"groupEquity"`
	GroupRevenue     string `json:"groupRevenue"`
	GroupExpenses    string `json:"groupExpenses"`
	// Normal balance value labels
	NormalBalanceDebit  string `json:"normalBalanceDebit"`
	NormalBalanceCredit string `json:"normalBalanceCredit"`
	// Info popover text
	ElementInfo        string `json:"elementInfo"`
	ClassificationInfo string `json:"classificationInfo"`
	CashFlowClassInfo  string `json:"cashFlowClassInfo"`
}

type AccountActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
	// Confirm messages
	ConfirmDelete     string `json:"confirmDelete"`
	ConfirmBulkDelete string `json:"confirmBulkDelete"`
	// Error messages
	NoPermission string `json:"noPermission"`
}

type AccountDetailLabels struct {
	Tabs        AccountDetailTabLabels   `json:"tabs"`
	EmptyStates AccountDetailEmptyLabels `json:"emptyStates"`
	Stats       AccountDetailStatLabels  `json:"stats"`
	Info        AccountDetailInfoLabels  `json:"info"`
}

type AccountDetailTabLabels struct {
	JournalEntries string `json:"journalEntries"`
	Details        string `json:"details"`
}

type AccountDetailEmptyLabels struct {
	EntriesTitle   string `json:"entriesTitle"`
	EntriesMessage string `json:"entriesMessage"`
}

type AccountDetailStatLabels struct {
	CurrentBalance string `json:"currentBalance"`
	PeriodDebits   string `json:"periodDebits"`
	PeriodCredits  string `json:"periodCredits"`
}

type AccountDetailInfoLabels struct {
	Title          string `json:"title"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Element        string `json:"element"`
	Classification string `json:"classification"`
	Group          string `json:"group"`
	ParentAccount  string `json:"parentAccount"`
	NormalBalance  string `json:"normalBalance"`
	CashFlowTag    string `json:"cashFlowTag"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	Created        string `json:"created"`
	LastModified   string `json:"lastModified"`
}

// DefaultAccountLabels returns AccountLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultAccountLabels() AccountLabels {
	return AccountLabels{
		Page: AccountPageLabels{
			Heading: "Chart of Accounts",
			Caption: "Manage your organization's account structure",
		},
		Buttons: AccountButtonLabels{
			AddAccount: "Add Account",
		},
		Columns: AccountColumnLabels{
			Code:           "Code",
			Name:           "Account Name",
			Element:        "Element",
			Classification: "Class",
			Group:          "Group",
			Balance:        "Balance",
			Date:           "Date",
			EntryNumber:    "Entry #",
			Description:    "Description",
			Debit:          "Debit",
			Credit:         "Credit",
			Status:         "Status",
		},
		Tabs: AccountTabLabels{
			All:       "All",
			Asset:     "Assets",
			Liability: "Liabilities",
			Equity:    "Equity",
			Revenue:   "Revenue",
			Expense:   "Expenses",
		},
		Empty: AccountEmptyLabels{
			Title:   "No accounts found",
			Message: "Add your first account or apply a template to get started.",
		},
		Form: AccountFormLabels{
			Code:                     "Account Code",
			CodePlaceholder:          "e.g. 1140",
			Name:                     "Account Name",
			NamePlaceholder:          "e.g. Petty Cash",
			Element:                  "Element",
			Classification:           "Class",
			ParentAccount:            "Parent Account",
			ParentAccountPlaceholder: "Search groups…",
			Group:                    "Group",
			IsGroup:                  "Is Group",
			Active:                   "Active",
			Description:              "Description",
			DescriptionPlaceholder:   "Brief description of this account",
			CashFlowSection:          "Cash Flow Tag (optional)",
			CashFlowClassification:   "Cash Flow Classification",
			ElementAsset:             "Asset",
			ElementLiability:         "Liability",
			ElementEquity:            "Equity",
			ElementRevenue:           "Revenue",
			ElementExpense:           "Expense",
			ClassCurrentAsset:        "Current Asset",
			ClassNonCurrentAsset:     "Non-Current Asset",
			ClassCurrentLiability:    "Current Liability",
			ClassNonCurrentLiability: "Non-Current Liability",
			ClassEquity:              "Equity",
			ClassOperatingRevenue:    "Operating Revenue",
			ClassOtherIncome:         "Other Income",
			ClassCostOfSales:         "Cost of Sales",
			ClassOperatingExpense:    "Operating Expense",
			ClassFinanceCost:         "Finance Cost",
			ClassIncomeTax:           "Income Tax",
			ClassOtherExpense:        "Other Expense",
			CashFlowNone:             "None",
			CashFlowOperating:        "Operating Activities",
			CashFlowInvesting:        "Investing Activities",
			CashFlowFinancing:        "Financing Activities",
			GroupAssets:              "ASSETS",
			GroupLiabilities:         "LIABILITIES",
			GroupEquity:              "EQUITY",
			GroupRevenue:             "REVENUE",
			GroupExpenses:            "EXPENSES",
			NormalBalanceDebit:       "Debit",
			NormalBalanceCredit:      "Credit",
			ElementInfo:              "The broad accounting category this account belongs to (Asset, Liability, Equity, Revenue, or Expense).",
			ClassificationInfo:       "Sub-category within the element that determines where this account appears on financial statements.",
			CashFlowClassInfo:        "Tags this account's movements for cash flow statement classification. Leave empty if not applicable.",
		},
		Actions: AccountActionLabels{
			View:              "View",
			Edit:              "Edit",
			Delete:            "Delete",
			ConfirmDelete:     "Are you sure you want to delete %s? This action cannot be undone.",
			ConfirmBulkDelete: "Are you sure you want to delete {{count}} account(s)? This action cannot be undone.",
			NoPermission:      "No permission",
		},
		Detail: AccountDetailLabels{
			Tabs: AccountDetailTabLabels{
				JournalEntries: "Journal Entries",
				Details:        "Details",
			},
			EmptyStates: AccountDetailEmptyLabels{
				EntriesTitle:   "No journal entries",
				EntriesMessage: "Journal entries that touch this account will appear here.",
			},
			Stats: AccountDetailStatLabels{
				CurrentBalance: "Current Balance",
				PeriodDebits:   "Period Debits",
				PeriodCredits:  "Period Credits",
			},
			Info: AccountDetailInfoLabels{
				Title:          "Account Information",
				Code:           "Account Code",
				Name:           "Account Name",
				Element:        "Element",
				Classification: "Classification",
				Group:          "Group",
				ParentAccount:  "Parent Account",
				NormalBalance:  "Normal Balance",
				CashFlowTag:    "Cash Flow Tag",
				Description:    "Description",
				Status:         "Status",
				Created:        "Created",
				LastModified:   "Last Modified",
			},
		},
		Templates: AccountTemplatesLabels{
			PageTitle:           "Account Templates",
			PageSubtitle:        "Pre-built Chart of Accounts for your business type",
			CurrentAccountCount: "Your Chart of Accounts currently has {{.CurrentAccountCount}} accounts.",
			ApplyWarning:        "Applying a template will add new accounts. Existing accounts with matching codes will be skipped.",
			Empty:               "Your Chart of Accounts is empty.",
			EmptyDesc:           "Apply a template below to get started with a standard set of accounts for your business type.",
			AccountsSuffix:      "accounts · PFRS-compliant",
			ComingSoon:          "Coming soon",
			BadgeApplied:        "Applied",
			BadgeAssets:         "Assets",
			BadgeLiabilities:    "Liabilities",
			BadgeEquity:         "Equity",
			BadgeRevenue:        "Revenue",
			BadgeExpenses:       "Expenses",
			Preview:             "Preview",
			AlreadyApplied:      "Already applied",
			ApplyTemplate:       "Apply Template",
			PreviewTitle:        "Preview",
			PreviewDesc:         "This template will create {{.AccountCount}} accounts organized as follows:",
			ColCode:             "Code",
			ColAccountName:      "Account Name",
			ColElement:          "Element",
			ColClass:            "Class",
			ColIsGroup:          "Is Group",
			Yes:                 "Yes",
			SkipNote:            "Accounts with matching codes in your existing Chart of Accounts will be skipped.",
		},
		GeneralLedger: AccountGeneralLedgerLabels{
			Title:                 "General Ledger",
			Subtitle:              "Detailed transaction history by account",
			Account:               "Account",
			AccountPlaceholder:    "Select an account",
			StartDate:             "Start Date",
			EndDate:               "End Date",
			Apply:                 "Apply",
			Clear:                 "Clear",
			Print:                 "Print",
			SelectAccountMessage:  "Select an account above to view its detailed transaction history.",
			NoTransactionsMessage: "No transactions found for the selected account and date range.",
			DateRangeSeparator:    "to",
			OpeningBalance:        "Opening Balance",
			PeriodDebits:          "Period Debits",
			PeriodCredits:         "Period Credits",
			RunningBalance:        "Running Balance",
			NoTransactionsTitle:   "No transactions",
			NoTransactionsDetail:  "No journal entries found for this account in the selected date range.",
		},
		Dashboard: LedgerDashboardLabels{
			Title:                  "Ledger Dashboard",
			Subtitle:               "Live position of your books — assets, liabilities, equity, and net income",
			TotalAssets:            "Total Assets",
			TotalLiabilities:       "Total Liabilities",
			TotalEquity:            "Total Equity",
			NetIncomeMTD:           "Net Income (MTD)",
			BalanceByType:          "Account Balance by Type",
			UnpostedJournals:       "Unposted Journals",
			RecentJournals:         "Recent Journals",
			NoRecentJournals:       "No recent journal entries",
			QuickNewJournal:        "New Journal Entry",
			QuickTrialBalance:      "Trial Balance",
			QuickClosePeriod:       "Close Period",
			QuickAccountLookup:     "Account Lookup",
			AccountTypeAssets:      "Assets",
			AccountTypeLiabilities: "Liabilities",
			AccountTypeEquity:      "Equity",
			AccountTypeRevenue:     "Revenue",
			AccountTypeExpense:     "Expense",
			JournalEntryFallback:   "Journal",
			ViewAll:                "View All",
			AxisAmount:             "Amount",
		},
	}
}
