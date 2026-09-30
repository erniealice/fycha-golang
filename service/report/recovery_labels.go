package report

// RecoveryReportsLabels holds the translatable strings of the recovery reports.
// Lyngua file: general/recovery_report.json, root key "recovery_report". JSON
// tags byte-match the key tree; every field has a general-tier key
// (TestDefaultRecoveryLabelsEqualGeneralLyngua).
type RecoveryReportsLabels struct {
	RecoverablesAging RecoverablesAgingLabels `json:"recoverables_aging"`
	Reconciliation    ReconciliationLabels    `json:"reconciliation"`
	Shared            RecoverySharedLabels    `json:"shared"`
}

// RecoverySharedLabels are the filter-bar strings both reports share.
type RecoverySharedLabels struct {
	Filters       string `json:"filters"`
	Apply         string `json:"apply"`
	Clear         string `json:"clear"`
	Total         string `json:"total"`
	Client        string `json:"client"`
	Currency      string `json:"currency"`
	Unavailable   string `json:"unavailable"`
	InvalidFilter string `json:"invalid_filter"`
}

// RecoverablesAgingLabels are the strings of the recoverables aging report.
type RecoverablesAgingLabels struct {
	PageTitle        string `json:"page_title"`
	PageDescription  string `json:"page_description"`
	Bucket0To30      string `json:"bucket_0_to_30"`
	Bucket31To60     string `json:"bucket_31_to_60"`
	Bucket61To90     string `json:"bucket_61_to_90"`
	BucketOver90     string `json:"bucket_over_90"`
	TotalOutstanding string `json:"total_outstanding"`
	DocumentCount    string `json:"document_count"`
	EmptyTitle       string `json:"empty_title"`
	EmptyMessage     string `json:"empty_message"`
	ExportFilename   string `json:"export_filename"`
	FilterAsOfDate   string `json:"filter_as_of_date"`
	FilterClient     string `json:"filter_client"`
	FilterClientHint string `json:"filter_client_hint"`
	ChipAsOf         string `json:"chip_as_of"`
	ChipClient       string `json:"chip_client"`
}

// ReconciliationLabels are the strings of the cost-source reconciliation report.
type ReconciliationLabels struct {
	PageTitle             string `json:"page_title"`
	PageDescription       string `json:"page_description"`
	Component             string `json:"component"`
	Amount                string `json:"amount"`
	SharesTotal           string `json:"shares_total"`
	RecoverableTotal      string `json:"recoverable_total"`
	ShareVariance         string `json:"share_variance"`
	ChargesOpen           string `json:"charges_open"`
	ChargesIssued         string `json:"charges_issued"`
	ChargesCancelled      string `json:"charges_cancelled"`
	ChargeVariance        string `json:"charge_variance"`
	Corrections           string `json:"corrections"`
	IssuedNet             string `json:"issued_net"`
	Status                string `json:"status"`
	StatusReconciled      string `json:"status_reconciled"`
	StatusMismatch        string `json:"status_mismatch"`
	StatusNoBatch         string `json:"status_no_batch"`
	PromptTitle           string `json:"prompt_title"`
	PromptMessage         string `json:"prompt_message"`
	EmptyTitle            string `json:"empty_title"`
	EmptyMessage          string `json:"empty_message"`
	ExportFilename        string `json:"export_filename"`
	FilterExpenditure     string `json:"filter_expenditure"`
	FilterExpenditureHint string `json:"filter_expenditure_hint"`
	ChipExpenditure       string `json:"chip_expenditure"`
}

// DefaultRecoveryReportsLabels returns the English defaults (identical to
// lyngua general/recovery_report.json).
func DefaultRecoveryReportsLabels() RecoveryReportsLabels {
	return RecoveryReportsLabels{
		Shared: RecoverySharedLabels{
			Filters: "Filters", Apply: "Apply", Clear: "Clear", Total: "Total", Client: "Client", Currency: "Currency",
			Unavailable:   "This report is not available.",
			InvalidFilter: "Check the filters and try again.",
		},
		RecoverablesAging: RecoverablesAgingLabels{
			PageTitle:        "Recoverables Aging",
			PageDescription:  "Issued recovery documents that are still unpaid, by how long they are past due",
			Bucket0To30:      "0-30 days",
			Bucket31To60:     "31-60 days",
			Bucket61To90:     "61-90 days",
			BucketOver90:     "Over 90 days",
			TotalOutstanding: "Outstanding",
			DocumentCount:    "Documents",
			EmptyTitle:       "Nothing outstanding",
			EmptyMessage:     "No unpaid recovery documents as of this date.",
			ExportFilename:   "recoverables-aging.csv",
			FilterAsOfDate:   "As of date",
			FilterClient:     "Client",
			FilterClientHint: "Leave blank for all clients",
			ChipAsOf:         "As of:",
			ChipClient:       "Client:",
		},
		Reconciliation: ReconciliationLabels{
			PageTitle:             "Cost Reconciliation",
			PageDescription:       "For each recoverable cost: the amount, how it was split, and what was charged",
			Component:             "Cost component",
			Amount:                "Amount",
			SharesTotal:           "Split total",
			RecoverableTotal:      "Recoverable",
			ShareVariance:         "Split difference",
			ChargesOpen:           "Charged (open)",
			ChargesIssued:         "Charged (issued)",
			ChargesCancelled:      "Charged (cancelled)",
			ChargeVariance:        "Charge difference",
			Corrections:           "Corrections",
			IssuedNet:             "Issued, net",
			Status:                "Status",
			StatusReconciled:      "Reconciled",
			StatusMismatch:        "Mismatch",
			StatusNoBatch:         "Not allocated",
			PromptTitle:           "Choose an expenditure",
			PromptMessage:         "Enter an expenditure to reconcile its recoverable costs.",
			EmptyTitle:            "No recoverable costs",
			EmptyMessage:          "This expenditure has no recoverable cost components.",
			ExportFilename:        "cost-reconciliation.csv",
			FilterExpenditure:     "Expenditure",
			FilterExpenditureHint: "The expenditure whose costs are reconciled",
			ChipExpenditure:       "Expenditure:",
		},
	}
}
