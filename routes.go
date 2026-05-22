package fycha

const (
	ReportsBaseURL                          = "/reports/"
	ReportsDashboardURL                     = "/reports/dashboard"
	ReportsRevenueURL                       = "/reports/revenue"
	ReportsCostOfSalesURL                   = "/reports/cost-of-sales"
	ReportsGrossProfitURL                   = "/reports/gross-profit"
	ReportsExpensesURL                      = "/reports/expenses"
	ReportsNetProfitURL                     = "/reports/net-profit"
	ReportsRevenueReportURL                 = "/reports/revenue-report"
	ReportsRevenueReportExportURL           = "/reports/revenue-report/export"
	ReportsExpenditureReportURL             = "/reports/expenditure-report"
	ReportsExpenditureReportExportURL       = "/reports/expenditure-report/export"
	ReportsDisbursementReportURL            = "/reports/disbursement-report"
	ReportsDisbursementReportExportURL      = "/reports/disbursement-report/export"
	ReportsReceivablesAgingReportURL        = "/reports/receivables-aging"
	ReportsReceivablesAgingReportExportURL  = "/reports/receivables-aging/export"
	ReportsPayablesAgingReportURL           = "/suppliers/reports/payables-aging"
	ReportsPayablesAgingReportExportURL     = "/suppliers/reports/payables-aging/export"
	ReportsCollectionSummaryReportURL       = "/reports/collection-summary"
	ReportsCollectionSummaryReportExportURL = "/reports/collection-summary/export"

	// StorageImagesPrefix is the default route prefix for image serving.
	StorageImagesPrefix = "/storage/images"

	// Cash report routes
	CashBookURL = "/cash/reports/cash-book"

	// Asset routes
	AssetDashboardURL        = "/asset/dashboard"
	AssetListURL             = "/asset/list/{status}"
	AssetDetailURL           = "/asset/detail/{id}"
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
	AssetLapsingScheduleURL      = "/asset/reports/lapsing-schedule"
	AssetDepreciationPoliciesURL = "/asset/settings/depreciation-policies"

	// Asset depreciation-run drawer routes (Surface A + E)
	AssetDepreciationRunURL    = "/action/asset/depreciation-run/{asset_id}"
	AssetRevaluationURL        = "/action/asset/revaluation/{asset_id}"
	AssetRevaluationPreviewURL = "/action/asset/revaluation-preview/{asset_id}"

	// Asset category / policy depreciation drawer routes (Surface C + F)
	AssetCategoryDepreciationRunURL   = "/action/asset-category/depreciation-run/{category_id}"
	AssetPolicyDepreciationRunURL     = "/action/asset-policy/depreciation-run/{category_id}"
	AssetPolicyDepreciationPreviewURL = "/action/asset-policy/depreciation-preview/{category_id}"

	// Lapsing schedule page routes (Surface B)
	LapsingScheduleListURL               = "/asset/lapsing-schedule/list"
	LapsingScheduleBulkRunSelectedURL    = "/action/lapsing-schedule/bulk-run-selected"
	LapsingScheduleBulkRunAllMatchingURL = "/action/lapsing-schedule/bulk-run-all-matching"
	AssetListBulkRunSelectedURL          = "/action/asset/bulk-run-selected"

	// Depreciation run history page routes (Surface D)
	DepreciationRunListURL            = "/asset/depreciation-runs/list/{status}"
	DepreciationRunListTableURL       = "/action/depreciation-run/table/{status}"
	DepreciationRunDetailURL          = "/asset/depreciation-runs/detail/{run_id}"
	DepreciationRunDetailTabActionURL = "/action/depreciation-run/detail/{run_id}/tab/{tab}"

	// Depreciation policies page route (Surface F — replaces mock)
	DepreciationPoliciesURL = "/asset/settings/depreciation-policies"

	// Ledger — Chart of Accounts
	LedgerBaseURL       = "/ledger/"
	LedgerDashboardURL  = "/ledger/dashboard"
	AccountListURL      = "/ledger/accounts/list"
	AccountDetailURL    = "/ledger/accounts/detail/{id}"
	AccountTabActionURL = "/action/ledger/account/{id}/tab/{tab}"
	AccountTreeURL      = "/ledger/accounts/tree"
	AccountTemplatesURL = "/ledger/settings/account-templates"
	AccountAddURL       = "/action/ledger/account/add"
	AccountEditURL      = "/action/ledger/account/edit/{id}"
	AccountDeleteURL    = "/action/ledger/account/delete"

	// Ledger — Journal Entries
	JournalListURL             = "/ledger/journals/list/{status}"
	JournalDetailURL           = "/ledger/journals/detail/{id}"
	JournalTabActionURL        = "/action/ledger/journal/{id}/tab/{tab}"
	JournalAttachmentUploadURL = "/action/ledger/journal/{id}/attachments/upload"
	JournalAttachmentDeleteURL = "/action/ledger/journal/{id}/attachments/delete"
	JournalAddURL              = "/action/ledger/journal/add"
	JournalEditURL             = "/action/ledger/journal/edit/{id}"
	JournalPostURL             = "/action/ledger/journal/post/{id}"
	JournalReverseURL          = "/action/ledger/journal/reverse/{id}"
	JournalDeleteURL           = "/action/ledger/journal/delete"

	// Ledger — Accounting Statements (internal tools)
	LedgerGeneralLedgerURL = "/ledger/reports/general-ledger"
	LedgerTrialBalanceURL  = "/ledger/reports/trial-balance"

	// Ledger — Fiscal Periods / Settings
	FiscalPeriodListURL   = "/ledger/settings/fiscal-periods"
	FiscalPeriodDetailURL = "/ledger/settings/fiscal-periods/detail/{id}"
	FiscalPeriodAddURL    = "/action/ledger/fiscal-period/add"
	FiscalPeriodCloseURL  = "/action/ledger/fiscal-period/close/{id}"

	// Ledger — Bad Debt Policy
	BadDebtPolicyURL = "/ledger/settings/bad-debt-policy"

	// Ledger — Recurring Templates
	RecurringTemplatesURL = "/ledger/settings/recurring"

	// Reports — Financial Statements (business-stakeholder output)
	ReportsIncomeStatementURL = "/reports/income-statement"
	ReportsBalanceSheetURL    = "/reports/balance-sheet"
	ReportsCashFlowURL        = "/reports/cash-flow"
	ReportsEquityChangesURL   = "/reports/equity-changes"

	// Funding — Loans
	LoanDashboardURL    = "/funding/loans/dashboard"
	LoanListURL         = "/funding/loans/list/{status}"
	LoanDetailURL       = "/funding/loans/detail/{id}"
	LoanAddURL          = "/action/funding/loan/add"
	LoanAmortizationURL = "/funding/loans/amortization"
	LoanPaymentAddURL   = "/action/funding/loan/payment/add"
	LoanPaymentListURL  = "/funding/loans/payments/{status}"

	// Funding — Equity
	EquityDashboardURL      = "/funding/equity/dashboard"
	EquityAccountsURL       = "/funding/equity/accounts"
	EquityAccountDetailURL  = "/funding/equity/accounts/detail/{id}"
	EquityTransactionsURL   = "/funding/equity/transactions"
	EquityTransactionAddURL = "/action/funding/equity/transaction/add"

	// Payroll
	PayrollDashboardURL      = "/payroll/dashboard"
	PayrollRunListURL        = "/payroll/runs/{status}"
	PayrollRunDetailURL      = "/payroll/runs/detail/{id}"
	PayrollRemittanceListURL = "/payroll/remittances/{status}"
	PayrollEmployeeListURL   = "/payroll/employees"
	PayrollGovRatesURL       = "/payroll/settings/gov-rates"
	PayrollPayPeriodsURL     = "/payroll/settings/pay-periods"

	// Cash — Deposits and Petty Cash
	DepositListURL                = "/cash/deposits/{status}"
	PettyCashRegisterURL          = "/cash/petty-cash/register"
	PettyCashReplenishmentListURL = "/cash/petty-cash/replenishments/{status}"
	PettyCashCustodianBalancesURL = "/cash/petty-cash/custodian-balances"

	// Expenses — Prepayments
	PrepaymentListURL         = "/expenses/prepayments/{status}"
	PrepaymentAmortizationURL = "/expenses/prepayments/amortization"

	// Tax — Tax Rates (read-only; supersession via admin SQL recipe)
	TaxRateListURL   = "/tax/tax-rates/list/{status}"
	TaxRateDetailURL = "/tax/tax-rates/detail/{id}"

	// Finance — Forex Rates (read-only in UI; appended only via RecordOperatorRate)
	ForexRateListURL   = "/finance/forex-rates/list/{status}"
	ForexRateDetailURL = "/finance/forex-rates/detail/{id}"

	// Treasury — Withholding Certificates (full CRUD)
	WithholdingCertificateListURL       = "/treasury/withholding-certificates/list/{status}"
	WithholdingCertificateDetailURL     = "/treasury/withholding-certificates/detail/{id}"
	WithholdingCertificateTableURL      = "/action/withholding-certificate/table/{status}"
	WithholdingCertificateAddURL        = "/action/withholding-certificate/add"
	WithholdingCertificateEditURL       = "/action/withholding-certificate/edit/{id}"
	WithholdingCertificateDeleteURL     = "/action/withholding-certificate/delete"
	WithholdingCertificateBulkDeleteURL = "/action/withholding-certificate/bulk-delete"
	WithholdingCertificateSetStatusURL  = "/action/withholding-certificate/set-status"
)
