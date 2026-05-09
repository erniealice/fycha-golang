package fycha

const (
	ReportsBaseURL                = "/app/reports/"
	ReportsDashboardURL           = "/app/reports/dashboard"
	ReportsRevenueURL             = "/app/reports/revenue"
	ReportsCostOfSalesURL         = "/app/reports/cost-of-sales"
	ReportsGrossProfitURL         = "/app/reports/gross-profit"
	ReportsExpensesURL            = "/app/reports/expenses"
	ReportsNetProfitURL           = "/app/reports/net-profit"
	ReportsRevenueReportURL       = "/app/reports/revenue-report"
	ReportsRevenueReportExportURL = "/app/reports/revenue-report/export"
	ReportsExpenditureReportURL       = "/app/reports/expenditure-report"
	ReportsExpenditureReportExportURL = "/app/reports/expenditure-report/export"
	ReportsDisbursementReportURL       = "/app/reports/disbursement-report"
	ReportsDisbursementReportExportURL = "/app/reports/disbursement-report/export"
	ReportsReceivablesAgingReportURL       = "/app/reports/receivables-aging"
	ReportsReceivablesAgingReportExportURL = "/app/reports/receivables-aging/export"
	ReportsPayablesAgingReportURL          = "/app/suppliers/reports/payables-aging"
	ReportsPayablesAgingReportExportURL    = "/app/suppliers/reports/payables-aging/export"
	ReportsCollectionSummaryReportURL       = "/app/reports/collection-summary"
	ReportsCollectionSummaryReportExportURL = "/app/reports/collection-summary/export"

	// StorageImagesPrefix is the default route prefix for image serving.
	StorageImagesPrefix = "/storage/images"

	// Cash report routes
	CashBookURL = "/app/cash/reports/cash-book"

	// Asset routes
	AssetDashboardURL        = "/app/assets/dashboard"
	AssetListURL             = "/app/assets/list/{status}"
	AssetDetailURL           = "/app/assets/detail/{id}"
	AssetTabActionURL        = "/action/asset/{id}/tab/{tab}"
	AssetAttachmentUploadURL = "/action/asset/{id}/attachments/upload"
	AssetAttachmentDeleteURL = "/action/asset/{id}/attachments/delete"
	AssetTableURL            = "/action/asset/table/{status}"
	AssetAddURL              = "/action/asset/add"
	AssetEditURL             = "/action/asset/edit/{id}"
	AssetDeleteURL           = "/action/asset/delete"
	AssetBulkDeleteURL       = "/action/asset/bulk-delete"
	AssetSetStatusURL        = "/action/asset/set-status"
	AssetBulkSetStatusURL    = "/action/asset/bulk-set-status"

	// Asset report/settings routes (legacy mock paths — kept for backwards compat)
	AssetLapsingScheduleURL      = "/app/assets/reports/lapsing-schedule"
	AssetDepreciationPoliciesURL = "/app/assets/settings/depreciation-policies"

	// Asset depreciation-run drawer routes (Surface A + E)
	AssetDepreciationRunURL      = "/action/asset/depreciation-run/{asset_id}"
	AssetRevaluationURL          = "/action/asset/revaluation/{asset_id}"
	AssetRevaluationPreviewURL   = "/action/asset/revaluation-preview/{asset_id}"

	// Asset category / policy depreciation drawer routes (Surface C + F)
	AssetCategoryDepreciationRunURL    = "/action/asset-category/depreciation-run/{category_id}"
	AssetPolicyDepreciationRunURL      = "/action/asset-policy/depreciation-run/{category_id}"
	AssetPolicyDepreciationPreviewURL  = "/action/asset-policy/depreciation-preview/{category_id}"

	// Lapsing schedule page routes (Surface B)
	LapsingScheduleListURL              = "/app/assets/lapsing-schedule/list"
	LapsingScheduleBulkRunSelectedURL   = "/action/lapsing-schedule/bulk-run-selected"
	LapsingScheduleBulkRunAllMatchingURL = "/action/lapsing-schedule/bulk-run-all-matching"
	AssetListBulkRunSelectedURL         = "/action/asset/bulk-run-selected"

	// Depreciation run history page routes (Surface D)
	DepreciationRunListURL          = "/app/assets/depreciation-runs/list/{status}"
	DepreciationRunListTableURL     = "/action/depreciation-run/table/{status}"
	DepreciationRunDetailURL        = "/app/assets/depreciation-runs/detail/{run_id}"
	DepreciationRunDetailTabActionURL = "/action/depreciation-run/detail/{run_id}/tab/{tab}"

	// Depreciation policies page route (Surface F — replaces mock)
	DepreciationPoliciesURL = "/app/assets/settings/depreciation-policies"

	// Ledger — Chart of Accounts
	LedgerBaseURL       = "/app/ledger/"
	LedgerDashboardURL  = "/app/ledger/dashboard"
	AccountListURL      = "/app/ledger/accounts/list"
	AccountDetailURL    = "/app/ledger/accounts/detail/{id}"
	AccountTabActionURL = "/action/ledger/account/{id}/tab/{tab}"
	AccountTreeURL      = "/app/ledger/accounts/tree"
	AccountTemplatesURL = "/app/ledger/settings/account-templates"
	AccountAddURL       = "/action/ledger/account/add"
	AccountEditURL      = "/action/ledger/account/edit/{id}"
	AccountDeleteURL    = "/action/ledger/account/delete"

	// Ledger — Journal Entries
	JournalListURL                    = "/app/ledger/journals/list/{status}"
	JournalDetailURL                  = "/app/ledger/journals/detail/{id}"
	JournalTabActionURL               = "/action/ledger/journal/{id}/tab/{tab}"
	JournalAttachmentUploadURL        = "/action/ledger/journal/{id}/attachments/upload"
	JournalAttachmentDeleteURL        = "/action/ledger/journal/{id}/attachments/delete"
	JournalAddURL                     = "/action/ledger/journal/add"
	JournalEditURL                    = "/action/ledger/journal/edit/{id}"
	JournalPostURL                    = "/action/ledger/journal/post/{id}"
	JournalReverseURL                 = "/action/ledger/journal/reverse/{id}"
	JournalDeleteURL                  = "/action/ledger/journal/delete"

	// Ledger — Accounting Statements (internal tools)
	LedgerGeneralLedgerURL = "/app/ledger/reports/general-ledger"
	LedgerTrialBalanceURL  = "/app/ledger/reports/trial-balance"

	// Ledger — Fiscal Periods / Settings
	FiscalPeriodListURL   = "/app/ledger/settings/fiscal-periods"
	FiscalPeriodDetailURL = "/app/ledger/settings/fiscal-periods/detail/{id}"
	FiscalPeriodAddURL    = "/action/ledger/fiscal-period/add"
	FiscalPeriodCloseURL  = "/action/ledger/fiscal-period/close/{id}"

	// Ledger — Bad Debt Policy
	BadDebtPolicyURL = "/app/ledger/settings/bad-debt-policy"

	// Ledger — Recurring Templates
	RecurringTemplatesURL = "/app/ledger/settings/recurring"

	// Reports — Financial Statements (business-stakeholder output)
	ReportsIncomeStatementURL = "/app/reports/income-statement"
	ReportsBalanceSheetURL    = "/app/reports/balance-sheet"
	ReportsCashFlowURL        = "/app/reports/cash-flow"
	ReportsEquityChangesURL   = "/app/reports/equity-changes"

	// Funding — Loans
	LoanDashboardURL    = "/app/funding/loans/dashboard"
	LoanListURL         = "/app/funding/loans/list/{status}"
	LoanDetailURL       = "/app/funding/loans/detail/{id}"
	LoanAddURL          = "/action/funding/loan/add"
	LoanAmortizationURL = "/app/funding/loans/amortization"
	LoanPaymentAddURL   = "/action/funding/loan/payment/add"
	LoanPaymentListURL  = "/app/funding/loans/payments/{status}"

	// Funding — Equity
	EquityDashboardURL      = "/app/funding/equity/dashboard"
	EquityAccountsURL       = "/app/funding/equity/accounts"
	EquityAccountDetailURL  = "/app/funding/equity/accounts/detail/{id}"
	EquityTransactionsURL   = "/app/funding/equity/transactions"
	EquityTransactionAddURL = "/action/funding/equity/transaction/add"

	// Payroll
	PayrollDashboardURL      = "/app/payroll/dashboard"
	PayrollRunListURL        = "/app/payroll/runs/{status}"
	PayrollRunDetailURL      = "/app/payroll/runs/detail/{id}"
	PayrollRemittanceListURL = "/app/payroll/remittances/{status}"
	PayrollEmployeeListURL   = "/app/payroll/employees"
	PayrollGovRatesURL       = "/app/payroll/settings/gov-rates"
	PayrollPayPeriodsURL     = "/app/payroll/settings/pay-periods"

	// Cash — Deposits and Petty Cash
	DepositListURL                = "/app/cash/deposits/{status}"
	PettyCashRegisterURL          = "/app/cash/petty-cash/register"
	PettyCashReplenishmentListURL = "/app/cash/petty-cash/replenishments/{status}"
	PettyCashCustodianBalancesURL = "/app/cash/petty-cash/custodian-balances"

	// Expenses — Prepayments
	PrepaymentListURL         = "/app/expenses/prepayments/{status}"
	PrepaymentAmortizationURL = "/app/expenses/prepayments/amortization"
)
