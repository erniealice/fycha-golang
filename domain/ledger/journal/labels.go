package journal

// ---------------------------------------------------------------------------
// Journal labels (Journal Entries)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Journal Entries module.
type Labels struct {
	Page    PageLabels       `json:"page"`
	Tabs    TabLabels        `json:"tabs"`
	Buttons ButtonLabels     `json:"buttons"`
	Columns ColumnLabels     `json:"columns"`
	Empty   EmptyLabels      `json:"empty"`
	Actions ActionLabels     `json:"actions"`
	Lines   LineLabels       `json:"lines"`
	Form    FormLabels       `json:"form"`
	Detail  DetailLabels     `json:"detail"`
	Confirm ConfirmLabels    `json:"confirm"`
	Source  SourceTypeLabels `json:"source"`
}

// SourceTypeLabels maps JournalSourceType enum values to display strings.
// Keys mirror the proto JournalSourceType enum; used by the journal list and
// detail views to label the originating-transaction source.
type SourceTypeLabels struct {
	Manual                 string `json:"manual"`
	Revenue                string `json:"revenue"`
	Expenditure            string `json:"expenditure"`
	Collection             string `json:"collection"`
	Disbursement           string `json:"disbursement"`
	Depreciation           string `json:"depreciation"`
	AssetAcquisition       string `json:"assetAcquisition"`
	AssetDisposal          string `json:"assetDisposal"`
	Prepayment             string `json:"prepayment"`
	PrepaymentAmortization string `json:"prepaymentAmortization"`
	LoanReceipt            string `json:"loanReceipt"`
	LoanPayment            string `json:"loanPayment"`
	PettyCashReplenishment string `json:"pettyCashReplenishment"`
	BadDebtProvision       string `json:"badDebtProvision"`
	DeferredRevenue        string `json:"deferredRevenue"`
	EquityContribution     string `json:"equityContribution"`
	EquityWithdrawal       string `json:"equityWithdrawal"`
	EquityDistribution     string `json:"equityDistribution"`
	YearEndClose           string `json:"yearEndClose"`
	Recurring              string `json:"recurring"`
	Payroll                string `json:"payroll"`
}

// ConfirmLabels holds confirmation dialog strings for journal actions.
type ConfirmLabels struct {
	Post    string `json:"post"`
	Delete  string `json:"delete"`
	Reverse string `json:"reverse"`
}

// DetailLabels holds translatable strings for the journal detail page.
type DetailLabels struct {
	Stats       DetailStatLabels `json:"stats"`
	Info        DetailInfoLabels `json:"info"`
	SourceLabel string           `json:"sourceLabel"`
	ViewSource  string           `json:"viewSource"`
	// Balance status badges shown in totals row
	Balanced   string `json:"balanced"`
	Unbalanced string `json:"unbalanced"`
	Totals     string `json:"totals"`
	Difference string `json:"difference"`
	// Tab labels
	TabLines       string `json:"tabLines"`
	TabAttachments string `json:"tabAttachments"`
}

type DetailStatLabels struct {
	TotalDebit  string `json:"totalDebit"`
	TotalCredit string `json:"totalCredit"`
}

type DetailInfoLabels struct {
	Title         string `json:"title"`
	Date          string `json:"date"`
	Reference     string `json:"reference"`
	Description   string `json:"description"`
	Status        string `json:"status"`
	SourceType    string `json:"sourceType"`
	Notes         string `json:"notes"`
	PostedBy      string `json:"postedBy"`
	PostedAt      string `json:"postedAt"`
	ReversedBy    string `json:"reversedBy"`
	ReversedAt    string `json:"reversedAt"`
	ReversalEntry string `json:"reversalEntry"`
	Created       string `json:"created"`
	LastModified  string `json:"lastModified"`
}

type PageLabels struct {
	HeadingDraft     string `json:"headingDraft"`
	SubtitleDraft    string `json:"subtitleDraft"`
	HeadingPosted    string `json:"headingPosted"`
	SubtitlePosted   string `json:"subtitlePosted"`
	HeadingReversed  string `json:"headingReversed"`
	SubtitleReversed string `json:"subtitleReversed"`
}

type TabLabels struct {
	Draft    string `json:"draft"`
	Posted   string `json:"posted"`
	Reversed string `json:"reversed"`
}

type ButtonLabels struct {
	NewEntry string `json:"newEntry"`
}

type ColumnLabels struct {
	EntryNumber string `json:"entryNumber"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Source      string `json:"source"`
	Status      string `json:"status"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Post         string `json:"post"`
	Reverse      string `json:"reverse"`
	Delete       string `json:"delete"`
	NoPermission string `json:"noPermission"`
	// Error messages for action handlers
	SaveError    string `json:"saveError"`
	PostError    string `json:"postError"`
	ReverseError string `json:"reverseError"`
	DeleteError  string `json:"deleteError"`
}

// FormLabels holds translatable strings for the journal entry form.
type FormLabels struct {
	Date                   string `json:"date"`
	DatePlaceholder        string `json:"datePlaceholder"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Notes                  string `json:"notes"`
	NotesPlaceholder       string `json:"notesPlaceholder"`
	// Line table headers
	LineNumber         string `json:"lineNumber"`
	Account            string `json:"account"`
	AccountPlaceholder string `json:"accountPlaceholder"`
	Debit              string `json:"debit"`
	Credit             string `json:"credit"`
	Memo               string `json:"memo"`
	AddLine            string `json:"addLine"`
	RemoveLine         string `json:"removeLine"`
	// Balance alert
	BalancedMessage   string `json:"balancedMessage"`
	UnbalancedMessage string `json:"unbalancedMessage"`
	DifferenceLabel   string `json:"differenceLabel"`
	// Buttons
	SaveDraft string `json:"saveDraft"`
	PostEntry string `json:"postEntry"`
	// Section titles
	EntryDetails string `json:"entryDetails"`
	JournalLines string `json:"journalLines"`
	// Balance status hint (initial state before any values entered)
	BalanceHint string `json:"balanceHint"`
}

type LineLabels struct {
	AccountCode  string `json:"accountCode"`
	AccountName  string `json:"accountName"`
	Memo         string `json:"memo"`
	Debit        string `json:"debit"`
	Credit       string `json:"credit"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingDraft:     "Draft Journal Entries",
			SubtitleDraft:    "Review and post journal entries that are still in draft",
			HeadingPosted:    "Posted Journal Entries",
			SubtitlePosted:   "View journal entries that have been posted to the ledger",
			HeadingReversed:  "Reversed Journal Entries",
			SubtitleReversed: "View journal entries that have been reversed",
		},
		Tabs: TabLabels{
			Draft:    "Draft",
			Posted:   "Posted",
			Reversed: "Reversed",
		},
		Buttons: ButtonLabels{
			NewEntry: "New Entry",
		},
		Columns: ColumnLabels{
			EntryNumber: "Entry #",
			Date:        "Date",
			Description: "Description",
			Amount:      "Amount",
			Source:      "Source",
			Status:      "Status",
		},
		Empty: EmptyLabels{
			Title:   "No journal entries",
			Message: "No journal entries found for this status.",
		},
		Actions: ActionLabels{
			View:         "View",
			Edit:         "Edit",
			Post:         "Post",
			Reverse:      "Reverse",
			Delete:       "Delete",
			NoPermission: "No permission",
			SaveError:    "Failed to save journal entry",
			PostError:    "Failed to post journal entry",
			ReverseError: "Failed to reverse journal entry",
			DeleteError:  "Failed to delete journal entry",
		},
		Lines: LineLabels{
			AccountCode:  "Account Code",
			AccountName:  "Account Name",
			Memo:         "Memo",
			Debit:        "Debit",
			Credit:       "Credit",
			EmptyTitle:   "No journal lines",
			EmptyMessage: "This journal entry has no lines.",
		},
		Form: FormLabels{
			Date:                   "Date",
			DatePlaceholder:        "YYYY-MM-DD",
			Description:            "Description",
			DescriptionPlaceholder: "e.g. Office supplies purchase",
			Notes:                  "Notes",
			NotesPlaceholder:       "Optional notes or reference",
			LineNumber:             "#",
			Account:                "Account",
			AccountPlaceholder:     "Search by code or name…",
			Debit:                  "Debit",
			Credit:                 "Credit",
			Memo:                   "Memo",
			AddLine:                "+ Add Line",
			RemoveLine:             "Remove",
			BalancedMessage:        "Balanced — Total Debits equal Total Credits",
			UnbalancedMessage:      "Unbalanced — Debits and Credits do not match",
			DifferenceLabel:        "Difference",
			SaveDraft:              "Save as Draft",
			PostEntry:              "Post",
			EntryDetails:           "Entry Details",
			JournalLines:           "Journal Lines",
			BalanceHint:            "Enter debits and credits above",
		},
		Detail: DetailLabels{
			Stats: DetailStatLabels{
				TotalDebit:  "Total Debits",
				TotalCredit: "Total Credits",
			},
			Info: DetailInfoLabels{
				Title:         "Entry Details",
				Date:          "Date",
				Reference:     "Reference",
				Description:   "Description",
				Status:        "Status",
				SourceType:    "Source Type",
				Notes:         "Notes",
				PostedBy:      "Posted By",
				PostedAt:      "Posted At",
				ReversedBy:    "Reversed By",
				ReversedAt:    "Reversed At",
				ReversalEntry: "Reversal Entry",
				Created:       "Created",
				LastModified:  "Last Modified",
			},
			SourceLabel:    "Source",
			ViewSource:     "View Source →",
			Balanced:       "Balanced",
			Unbalanced:     "Unbalanced",
			Totals:         "TOTALS",
			Difference:     "DIFFERENCE",
			TabLines:       "Journal Lines",
			TabAttachments: "Attachments",
		},
		Confirm: ConfirmLabels{
			Post:    "Are you sure you want to post this journal entry? This action cannot be undone.",
			Delete:  "Are you sure you want to delete this journal entry? This action cannot be undone.",
			Reverse: "Are you sure you want to reverse this journal entry? A reversing entry will be created.",
		},
		Source: SourceTypeLabels{
			Manual:                 "Manual",
			Revenue:                "Revenue",
			Expenditure:            "Expenditure",
			Collection:             "Collection",
			Disbursement:           "Disbursement",
			Depreciation:           "Depreciation",
			AssetAcquisition:       "Asset Acquisition",
			AssetDisposal:          "Asset Disposal",
			Prepayment:             "Prepayment",
			PrepaymentAmortization: "Prepayment Amortization",
			LoanReceipt:            "Loan Receipt",
			LoanPayment:            "Loan Payment",
			PettyCashReplenishment: "Petty Cash Replenishment",
			BadDebtProvision:       "Bad Debt Provision",
			DeferredRevenue:        "Deferred Revenue",
			EquityContribution:     "Equity Contribution",
			EquityWithdrawal:       "Equity Withdrawal",
			EquityDistribution:     "Equity Distribution",
			YearEndClose:           "Year-End Close",
			Recurring:              "Recurring",
			Payroll:                "Payroll",
		},
	}
}
