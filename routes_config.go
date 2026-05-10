// routes_config.go defines configurable route structs for fycha views.
//
// Three-level routing system:
//   - Level 1: Generic defaults from Go consts (this file). DefaultXxxRoutes()
//     constructors return structs populated from the package-level route constants
//     defined in routes.go. These serve as sensible defaults for any consumer app.
//   - Level 2: Industry-specific overrides via JSON (loaded by consumer apps).
//     Apps can load a JSON config file that maps route keys to custom paths,
//     allowing industry templates (e.g. salon, retail) to rebrand URLs without
//     code changes. The json struct tags on each field support this workflow.
//   - Level 3: App-specific overrides via Go field assignment (optional).
//     After constructing defaults (and optionally applying JSON), consumer apps
//     can directly assign individual struct fields for one-off customizations.
//
// RouteMap() methods return a map[string]string of dot-notation keys to route
// paths, useful for template rendering and route resolution at runtime.
package fycha

// ReportsRoutes holds route paths for all reporting views.
type ReportsRoutes struct {
	DashboardURL   string `json:"dashboard_url"`
	RevenueURL     string `json:"revenue_url"`
	CostOfSalesURL string `json:"cost_of_sales_url"`
	GrossProfitURL string `json:"gross_profit_url"`
	ExpensesURL    string `json:"expenses_url"`
	NetProfitURL   string `json:"net_profit_url"`
	// Financial Statements (NEW — derived from ledger, exposed to business stakeholders)
	IncomeStatementURL string `json:"income_statement_url"`
	BalanceSheetURL    string `json:"balance_sheet_url"`
	CashFlowURL        string `json:"cash_flow_url"`
	EquityChangesURL   string `json:"equity_changes_url"`
	// Revenue Report pivot table
	RevenueReportURL       string `json:"revenue_report_url"`
	RevenueReportExportURL string `json:"revenue_report_export_url"`
	// Expenditure Report pivot table
	ExpenditureReportURL       string `json:"expenditure_report_url"`
	ExpenditureReportExportURL string `json:"expenditure_report_export_url"`
	// Disbursement Report pivot table
	DisbursementReportURL       string `json:"disbursement_report_url"`
	DisbursementReportExportURL string `json:"disbursement_report_export_url"`
	// Receivables Aging Report
	ReceivablesAgingReportURL       string `json:"receivables_aging_report_url"`
	ReceivablesAgingReportExportURL string `json:"receivables_aging_report_export_url"`
	// Payables Aging Report
	PayablesAgingReportURL       string `json:"payables_aging_report_url"`
	PayablesAgingReportExportURL string `json:"payables_aging_report_export_url"`
	// Collection Summary Report pivot table
	CollectionSummaryReportURL       string `json:"collection_summary_report_url"`
	CollectionSummaryReportExportURL string `json:"collection_summary_report_export_url"`
}

// DefaultReportsRoutes returns a ReportsRoutes populated from package-level consts.
func DefaultReportsRoutes() ReportsRoutes {
	return ReportsRoutes{
		DashboardURL:           ReportsDashboardURL,
		RevenueURL:             ReportsRevenueURL,
		CostOfSalesURL:         ReportsCostOfSalesURL,
		GrossProfitURL:         ReportsGrossProfitURL,
		ExpensesURL:            ReportsExpensesURL,
		NetProfitURL:           ReportsNetProfitURL,
		IncomeStatementURL:     ReportsIncomeStatementURL,
		BalanceSheetURL:        ReportsBalanceSheetURL,
		CashFlowURL:            ReportsCashFlowURL,
		EquityChangesURL:       ReportsEquityChangesURL,
		RevenueReportURL:       ReportsRevenueReportURL,
		RevenueReportExportURL: ReportsRevenueReportExportURL,
		ExpenditureReportURL:       ReportsExpenditureReportURL,
		ExpenditureReportExportURL: ReportsExpenditureReportExportURL,
		DisbursementReportURL:       ReportsDisbursementReportURL,
		DisbursementReportExportURL: ReportsDisbursementReportExportURL,
		ReceivablesAgingReportURL:       ReportsReceivablesAgingReportURL,
		ReceivablesAgingReportExportURL: ReportsReceivablesAgingReportExportURL,
		PayablesAgingReportURL:          ReportsPayablesAgingReportURL,
		PayablesAgingReportExportURL:    ReportsPayablesAgingReportExportURL,
		CollectionSummaryReportURL:       ReportsCollectionSummaryReportURL,
		CollectionSummaryReportExportURL: ReportsCollectionSummaryReportExportURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ReportsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"reports.dashboard":             r.DashboardURL,
		"reports.revenue":               r.RevenueURL,
		"reports.cost_of_sales":         r.CostOfSalesURL,
		"reports.gross_profit":          r.GrossProfitURL,
		"reports.expenses":              r.ExpensesURL,
		"reports.net_profit":            r.NetProfitURL,
		"reports.income_statement":      r.IncomeStatementURL,
		"reports.balance_sheet":         r.BalanceSheetURL,
		"reports.cash_flow":             r.CashFlowURL,
		"reports.equity_changes":        r.EquityChangesURL,
		"reports.revenue_report":        r.RevenueReportURL,
		"reports.revenue_report_export": r.RevenueReportExportURL,
		"reports.expenditure_report":        r.ExpenditureReportURL,
		"reports.expenditure_report_export": r.ExpenditureReportExportURL,
		"reports.disbursement_report":        r.DisbursementReportURL,
		"reports.disbursement_report_export": r.DisbursementReportExportURL,
		"reports.receivables_aging_report":        r.ReceivablesAgingReportURL,
		"reports.receivables_aging_report_export": r.ReceivablesAgingReportExportURL,
		"reports.payables_aging_report":           r.PayablesAgingReportURL,
		"reports.payables_aging_report_export":    r.PayablesAgingReportExportURL,
		"reports.collection_summary_report":        r.CollectionSummaryReportURL,
		"reports.collection_summary_report_export": r.CollectionSummaryReportExportURL,
	}
}

// ---------------------------------------------------------------------------
// AssetRoutes
// ---------------------------------------------------------------------------

// AssetRoutes holds route paths for fixed asset management views.
type AssetRoutes struct {
	DashboardURL     string `json:"dashboard_url"`
	ListURL          string `json:"list_url"`
	DetailURL        string `json:"detail_url"`
	TabActionURL     string `json:"tab_action_url"`
	TableURL         string `json:"table_url"`
	AddURL           string `json:"add_url"`
	EditURL          string `json:"edit_url"`
	DeleteURL        string `json:"delete_url"`
	BulkDeleteURL    string `json:"bulk_delete_url"`
	SetStatusURL     string `json:"set_status_url"`
	BulkSetStatusURL string `json:"bulk_set_status_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`

	// Report/settings routes (legacy mock paths)
	LapsingScheduleURL      string `json:"lapsing_schedule_url"`
	DepreciationPoliciesURL string `json:"depreciation_policies_url"`

	// Depreciation-run drawer routes (Surface A)
	DepreciationRunURL     string `json:"depreciation_run_url"`
	RevaluationURL         string `json:"revaluation_url"`
	RevaluationPreviewURL  string `json:"revaluation_preview_url"`
}

// DefaultAssetRoutes returns an AssetRoutes populated from package-level consts.
func DefaultAssetRoutes() AssetRoutes {
	return AssetRoutes{
		DashboardURL:     AssetDashboardURL,
		ListURL:          AssetListURL,
		DetailURL:        AssetDetailURL,
		TabActionURL:     AssetTabActionURL,
		TableURL:         AssetTableURL,
		AddURL:           AssetAddURL,
		EditURL:          AssetEditURL,
		DeleteURL:        AssetDeleteURL,
		BulkDeleteURL:    AssetBulkDeleteURL,
		SetStatusURL:     AssetSetStatusURL,
		BulkSetStatusURL: AssetBulkSetStatusURL,

		AttachmentUploadURL: AssetAttachmentUploadURL,
		AttachmentDeleteURL: AssetAttachmentDeleteURL,

		LapsingScheduleURL:      AssetLapsingScheduleURL,
		DepreciationPoliciesURL: AssetDepreciationPoliciesURL,

		DepreciationRunURL:    AssetDepreciationRunURL,
		RevaluationURL:        AssetRevaluationURL,
		RevaluationPreviewURL: AssetRevaluationPreviewURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r AssetRoutes) RouteMap() map[string]string {
	return map[string]string{
		"asset.dashboard":       r.DashboardURL,
		"asset.list":            r.ListURL,
		"asset.detail":          r.DetailURL,
		"asset.tab_action":      r.TabActionURL,
		"asset.table":           r.TableURL,
		"asset.add":             r.AddURL,
		"asset.edit":            r.EditURL,
		"asset.delete":          r.DeleteURL,
		"asset.bulk_delete":     r.BulkDeleteURL,
		"asset.set_status":      r.SetStatusURL,
		"asset.bulk_set_status": r.BulkSetStatusURL,

		"asset.attachment.upload": r.AttachmentUploadURL,
		"asset.attachment.delete": r.AttachmentDeleteURL,

		"asset.lapsing_schedule":      r.LapsingScheduleURL,
		"asset.depreciation_policies": r.DepreciationPoliciesURL,

		"asset.depreciation_run":     r.DepreciationRunURL,
		"asset.revaluation":          r.RevaluationURL,
		"asset.revaluation_preview":  r.RevaluationPreviewURL,
	}
}

// DepreciationRunFor returns the resolved Surface A depreciation-run drawer URL
// for the given asset ID.
func (r AssetRoutes) DepreciationRunFor(assetID string) string {
	return resolveParam(r.DepreciationRunURL, "asset_id", assetID)
}

// RevaluationFor returns the resolved Surface E revaluation drawer URL for the
// given asset ID.
func (r AssetRoutes) RevaluationFor(assetID string) string {
	return resolveParam(r.RevaluationURL, "asset_id", assetID)
}

// RevaluationPreviewFor returns the resolved HTMX revaluation-preview endpoint
// URL for the given asset ID.
func (r AssetRoutes) RevaluationPreviewFor(assetID string) string {
	return resolveParam(r.RevaluationPreviewURL, "asset_id", assetID)
}

// ---------------------------------------------------------------------------
// AccountRoutes
// ---------------------------------------------------------------------------

// AccountRoutes holds route paths for Chart of Accounts views.
type AccountRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	DashboardURL string `json:"dashboard_url"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	TreeURL      string `json:"tree_url"`
	TemplatesURL string `json:"templates_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	DeleteURL    string `json:"delete_url"`
}

func DefaultAccountRoutes() AccountRoutes {
	return AccountRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "chart-of-accounts",
		DashboardURL: LedgerDashboardURL,
		ListURL:      AccountListURL,
		DetailURL:    AccountDetailURL,
		TabActionURL: AccountTabActionURL,
		TreeURL:      AccountTreeURL,
		TemplatesURL: AccountTemplatesURL,
		AddURL:       AccountAddURL,
		EditURL:      AccountEditURL,
		DeleteURL:    AccountDeleteURL,
	}
}

func (r AccountRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.dashboard":         r.DashboardURL,
		"ledger.account.list":      r.ListURL,
		"ledger.account.detail":    r.DetailURL,
		"ledger.account.tree":      r.TreeURL,
		"ledger.account.templates": r.TemplatesURL,
		"ledger.account.add":       r.AddURL,
		"ledger.account.edit":      r.EditURL,
		"ledger.account.delete":    r.DeleteURL,
	}
}

// ---------------------------------------------------------------------------
// JournalRoutes
// ---------------------------------------------------------------------------

// JournalRoutes holds route paths for Journal Entry views.
type JournalRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	PostURL      string `json:"post_url"`
	ReverseURL   string `json:"reverse_url"`
	DeleteURL    string `json:"delete_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

func DefaultJournalRoutes() JournalRoutes {
	return JournalRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "journals-draft",
		ListURL:      JournalListURL,
		DetailURL:    JournalDetailURL,
		TabActionURL: JournalTabActionURL,
		AddURL:       JournalAddURL,
		EditURL:      JournalEditURL,
		PostURL:      JournalPostURL,
		ReverseURL:   JournalReverseURL,
		DeleteURL:    JournalDeleteURL,

		AttachmentUploadURL: JournalAttachmentUploadURL,
		AttachmentDeleteURL: JournalAttachmentDeleteURL,
	}
}

func (r JournalRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.journal.list":    r.ListURL,
		"ledger.journal.detail":  r.DetailURL,
		"ledger.journal.add":     r.AddURL,
		"ledger.journal.edit":    r.EditURL,
		"ledger.journal.post":    r.PostURL,
		"ledger.journal.reverse": r.ReverseURL,
		"ledger.journal.delete":  r.DeleteURL,

		"ledger.journal.attachment.upload": r.AttachmentUploadURL,
		"ledger.journal.attachment.delete": r.AttachmentDeleteURL,
	}
}

// ---------------------------------------------------------------------------
// LedgerStatementRoutes
// ---------------------------------------------------------------------------

// LedgerStatementRoutes holds route paths for accounting statement views
// (General Ledger, Trial Balance — internal accounting tools, not business reports).
type LedgerStatementRoutes struct {
	ActiveNav        string `json:"active_nav"`
	GeneralLedgerURL string `json:"general_ledger_url"`
	TrialBalanceURL  string `json:"trial_balance_url"`
}

func DefaultLedgerStatementRoutes() LedgerStatementRoutes {
	return LedgerStatementRoutes{
		ActiveNav:        "ledger",
		GeneralLedgerURL: LedgerGeneralLedgerURL,
		TrialBalanceURL:  LedgerTrialBalanceURL,
	}
}

func (r LedgerStatementRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.statement.general_ledger": r.GeneralLedgerURL,
		"ledger.statement.trial_balance":  r.TrialBalanceURL,
	}
}

// ---------------------------------------------------------------------------
// FiscalPeriodRoutes
// ---------------------------------------------------------------------------

// FiscalPeriodRoutes holds route paths for fiscal period management views.
type FiscalPeriodRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	AddURL       string `json:"add_url"`
	CloseURL     string `json:"close_url"`
}

func DefaultFiscalPeriodRoutes() FiscalPeriodRoutes {
	return FiscalPeriodRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "fiscal-periods",
		ListURL:      FiscalPeriodListURL,
		DetailURL:    FiscalPeriodDetailURL,
		AddURL:       FiscalPeriodAddURL,
		CloseURL:     FiscalPeriodCloseURL,
	}
}

func (r FiscalPeriodRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.fiscal_period.list":   r.ListURL,
		"ledger.fiscal_period.detail": r.DetailURL,
		"ledger.fiscal_period.add":    r.AddURL,
		"ledger.fiscal_period.close":  r.CloseURL,
	}
}

// ---------------------------------------------------------------------------
// LedgerSettingsRoutes
// ---------------------------------------------------------------------------

// LedgerSettingsRoutes holds route paths for ledger settings views
// (Bad Debt Policy, Recurring Templates).
type LedgerSettingsRoutes struct {
	ActiveNav             string `json:"active_nav"`
	BadDebtPolicyURL      string `json:"bad_debt_policy_url"`
	RecurringTemplatesURL string `json:"recurring_templates_url"`
}

func DefaultLedgerSettingsRoutes() LedgerSettingsRoutes {
	return LedgerSettingsRoutes{
		ActiveNav:             "ledger",
		BadDebtPolicyURL:      BadDebtPolicyURL,
		RecurringTemplatesURL: RecurringTemplatesURL,
	}
}

func (r LedgerSettingsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.settings.bad_debt_policy":     r.BadDebtPolicyURL,
		"ledger.settings.recurring_templates": r.RecurringTemplatesURL,
	}
}

// ---------------------------------------------------------------------------
// LoanRoutes
// ---------------------------------------------------------------------------

// LoanRoutes holds route paths for Loan views.
type LoanRoutes struct {
	ActiveNav       string `json:"active_nav"`
	DashboardURL    string `json:"dashboard_url"`
	ListURL         string `json:"list_url"`
	DetailURL       string `json:"detail_url"`
	AddURL          string `json:"add_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultLoanRoutes() LoanRoutes {
	return LoanRoutes{
		ActiveNav:       "loan",
		DashboardURL:    LoanDashboardURL,
		ListURL:         LoanListURL,
		DetailURL:       LoanDetailURL,
		AddURL:          LoanAddURL,
		AmortizationURL: LoanAmortizationURL,
	}
}

func (r LoanRoutes) RouteMap() map[string]string {
	return map[string]string{
		"loan.dashboard":    r.DashboardURL,
		"loan.list":         r.ListURL,
		"loan.detail":       r.DetailURL,
		"loan.add":          r.AddURL,
		"loan.amortization": r.AmortizationURL,
	}
}

// ---------------------------------------------------------------------------
// LoanPaymentRoutes
// ---------------------------------------------------------------------------

// LoanPaymentRoutes holds route paths for Loan Payment views.
type LoanPaymentRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	AddURL    string `json:"add_url"`
}

func DefaultLoanPaymentRoutes() LoanPaymentRoutes {
	return LoanPaymentRoutes{
		ActiveNav: "loan",
		ListURL:   LoanPaymentListURL,
		AddURL:    LoanPaymentAddURL,
	}
}

func (r LoanPaymentRoutes) RouteMap() map[string]string {
	return map[string]string{
		"loan_payment.list": r.ListURL,
		"loan_payment.add":  r.AddURL,
	}
}

// ---------------------------------------------------------------------------
// EquityRoutes
// ---------------------------------------------------------------------------

// EquityRoutes holds route paths for Equity views.
type EquityRoutes struct {
	ActiveNav         string `json:"active_nav"`
	DashboardURL      string `json:"dashboard_url"`
	AccountsURL       string `json:"accounts_url"`
	AccountDetailURL  string `json:"account_detail_url"`
	TransactionsURL   string `json:"transactions_url"`
	TransactionAddURL string `json:"transaction_add_url"`
}

func DefaultEquityRoutes() EquityRoutes {
	return EquityRoutes{
		ActiveNav:         "equity",
		DashboardURL:      EquityDashboardURL,
		AccountsURL:       EquityAccountsURL,
		AccountDetailURL:  EquityAccountDetailURL,
		TransactionsURL:   EquityTransactionsURL,
		TransactionAddURL: EquityTransactionAddURL,
	}
}

func (r EquityRoutes) RouteMap() map[string]string {
	return map[string]string{
		"equity.dashboard":       r.DashboardURL,
		"equity.accounts":        r.AccountsURL,
		"equity.account_detail":  r.AccountDetailURL,
		"equity.transactions":    r.TransactionsURL,
		"equity.transaction_add": r.TransactionAddURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollRunRoutes
// ---------------------------------------------------------------------------

// PayrollRunRoutes holds route paths for Payroll Run views.
type PayrollRunRoutes struct {
	ActiveNav    string `json:"active_nav"`
	DashboardURL string `json:"dashboard_url"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
}

func DefaultPayrollRunRoutes() PayrollRunRoutes {
	return PayrollRunRoutes{
		ActiveNav:    "payroll",
		DashboardURL: PayrollDashboardURL,
		ListURL:      PayrollRunListURL,
		DetailURL:    PayrollRunDetailURL,
	}
}

func (r PayrollRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.dashboard":  r.DashboardURL,
		"payroll.run.list":   r.ListURL,
		"payroll.run.detail": r.DetailURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollRemittanceRoutes
// ---------------------------------------------------------------------------

// PayrollRemittanceRoutes holds route paths for Payroll Remittance views.
type PayrollRemittanceRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
}

func DefaultPayrollRemittanceRoutes() PayrollRemittanceRoutes {
	return PayrollRemittanceRoutes{
		ActiveNav: "payroll",
		ListURL:   PayrollRemittanceListURL,
	}
}

func (r PayrollRemittanceRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.remittance.list": r.ListURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollEmployeeRoutes
// ---------------------------------------------------------------------------

// PayrollEmployeeRoutes holds route paths for Payroll Employee views.
type PayrollEmployeeRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
}

func DefaultPayrollEmployeeRoutes() PayrollEmployeeRoutes {
	return PayrollEmployeeRoutes{
		ActiveNav: "payroll",
		ListURL:   PayrollEmployeeListURL,
	}
}

func (r PayrollEmployeeRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.employee.list": r.ListURL,
	}
}

// ---------------------------------------------------------------------------
// PayrollSettingsRoutes
// ---------------------------------------------------------------------------

// PayrollSettingsRoutes holds route paths for Payroll Settings views.
type PayrollSettingsRoutes struct {
	ActiveNav     string `json:"active_nav"`
	GovRatesURL   string `json:"gov_rates_url"`
	PayPeriodsURL string `json:"pay_periods_url"`
}

func DefaultPayrollSettingsRoutes() PayrollSettingsRoutes {
	return PayrollSettingsRoutes{
		ActiveNav:     "payroll",
		GovRatesURL:   PayrollGovRatesURL,
		PayPeriodsURL: PayrollPayPeriodsURL,
	}
}

func (r PayrollSettingsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"payroll.settings.gov_rates":   r.GovRatesURL,
		"payroll.settings.pay_periods": r.PayPeriodsURL,
	}
}

// ---------------------------------------------------------------------------
// DepositRoutes
// ---------------------------------------------------------------------------

// DepositRoutes holds route paths for Cash app deposit views.
type DepositRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
}

func DefaultDepositRoutes() DepositRoutes {
	return DepositRoutes{
		ActiveNav: "cash",
		ListURL:   DepositListURL,
	}
}

func (r DepositRoutes) RouteMap() map[string]string {
	return map[string]string{
		"deposit.list": r.ListURL,
	}
}

// ---------------------------------------------------------------------------
// PettyCashRoutes
// ---------------------------------------------------------------------------

// PettyCashRoutes holds route paths for Cash app petty cash views.
type PettyCashRoutes struct {
	ActiveNav            string `json:"active_nav"`
	RegisterURL          string `json:"register_url"`
	ReplenishmentListURL string `json:"replenishment_list_url"`
	CustodianBalancesURL string `json:"custodian_balances_url"`
}

func DefaultPettyCashRoutes() PettyCashRoutes {
	return PettyCashRoutes{
		ActiveNav:            "cash",
		RegisterURL:          PettyCashRegisterURL,
		ReplenishmentListURL: PettyCashReplenishmentListURL,
		CustodianBalancesURL: PettyCashCustodianBalancesURL,
	}
}

func (r PettyCashRoutes) RouteMap() map[string]string {
	return map[string]string{
		"petty_cash.register":           r.RegisterURL,
		"petty_cash.replenishment_list": r.ReplenishmentListURL,
		"petty_cash.custodian_balances": r.CustodianBalancesURL,
	}
}

// ---------------------------------------------------------------------------
// PrepaymentRoutes
// ---------------------------------------------------------------------------

// PrepaymentRoutes holds route paths for Expenses app prepayment views.
type PrepaymentRoutes struct {
	ActiveNav       string `json:"active_nav"`
	ListURL         string `json:"list_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultPrepaymentRoutes() PrepaymentRoutes {
	return PrepaymentRoutes{
		ActiveNav:       "expense",
		ListURL:         PrepaymentListURL,
		AmortizationURL: PrepaymentAmortizationURL,
	}
}

func (r PrepaymentRoutes) RouteMap() map[string]string {
	return map[string]string{
		"prepayment.list":         r.ListURL,
		"prepayment.amortization": r.AmortizationURL,
	}
}

// ---------------------------------------------------------------------------
// LapsingScheduleRoutes
// ---------------------------------------------------------------------------

// LapsingScheduleRoutes holds route paths for the lapsing-schedule live page
// (Surface B) and its bulk-action endpoints.
type LapsingScheduleRoutes struct {
	// ListURL is the Surface B full-page lapsing-schedule list.
	ListURL string `json:"list_url"`

	// BulkRunSelectedURL is the endpoint for bulk-running selected rows.
	BulkRunSelectedURL string `json:"bulk_run_selected_url"`

	// BulkRunAllMatchingURL is the endpoint for bulk-running all matching rows.
	BulkRunAllMatchingURL string `json:"bulk_run_all_matching_url"`

	// AssetListBulkRunSelectedURL is the equivalent bulk-run endpoint on the
	// main assets list page.
	AssetListBulkRunSelectedURL string `json:"asset_list_bulk_run_selected_url"`
}

// DefaultLapsingScheduleRoutes returns a LapsingScheduleRoutes populated from
// the package-level route constants.
func DefaultLapsingScheduleRoutes() LapsingScheduleRoutes {
	return LapsingScheduleRoutes{
		ListURL:                     LapsingScheduleListURL,
		BulkRunSelectedURL:          LapsingScheduleBulkRunSelectedURL,
		BulkRunAllMatchingURL:       LapsingScheduleBulkRunAllMatchingURL,
		AssetListBulkRunSelectedURL: AssetListBulkRunSelectedURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r LapsingScheduleRoutes) RouteMap() map[string]string {
	return map[string]string{
		"lapsing_schedule.list":                         r.ListURL,
		"lapsing_schedule.bulk_run_selected":            r.BulkRunSelectedURL,
		"lapsing_schedule.bulk_run_all_matching":        r.BulkRunAllMatchingURL,
		"lapsing_schedule.asset_list_bulk_run_selected": r.AssetListBulkRunSelectedURL,
	}
}

// ---------------------------------------------------------------------------
// DepreciationRunRoutes
// ---------------------------------------------------------------------------

// DepreciationRunRoutes holds route paths for the depreciation-run history
// list and detail pages (Surface D).
type DepreciationRunRoutes struct {
	// ActiveNav is the sidebar key used to highlight the active nav item.
	ActiveNav string `json:"active_nav"`

	// ListURL is the Surface D list page; status is a path parameter
	// (complete | failed | pending, or empty for all).
	ListURL string `json:"list_url"`

	// ListTableURL is the HTMX inner-swap target for the list table.
	ListTableURL string `json:"list_table_url"`

	// DetailURL is the Surface D detail page; run_id is a path parameter.
	DetailURL string `json:"detail_url"`

	// DetailTabActionURL is the HTMX tab-swap target on the detail page.
	DetailTabActionURL string `json:"detail_tab_action_url"`
}

// DefaultDepreciationRunRoutes returns a DepreciationRunRoutes populated from
// the package-level route constants.
func DefaultDepreciationRunRoutes() DepreciationRunRoutes {
	return DepreciationRunRoutes{
		ActiveNav:          "depreciation-runs",
		ListURL:            DepreciationRunListURL,
		ListTableURL:       DepreciationRunListTableURL,
		DetailURL:          DepreciationRunDetailURL,
		DetailTabActionURL: DepreciationRunDetailTabActionURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r DepreciationRunRoutes) RouteMap() map[string]string {
	return map[string]string{
		"depreciation_run.list":             r.ListURL,
		"depreciation_run.list_table":       r.ListTableURL,
		"depreciation_run.detail":           r.DetailURL,
		"depreciation_run.detail_tab_action": r.DetailTabActionURL,
	}
}

// ListFor returns the resolved depreciation-run list URL for the given status
// (e.g. "complete", "failed", "pending", or empty string for all).
func (r DepreciationRunRoutes) ListFor(status string) string {
	return resolveParam(r.ListURL, "status", status)
}

// DetailFor returns the resolved depreciation-run detail URL for the given
// run ID.
func (r DepreciationRunRoutes) DetailFor(runID string) string {
	return resolveParam(r.DetailURL, "run_id", runID)
}

// ---------------------------------------------------------------------------
// AssetCategoryDepreciationRoutes
// ---------------------------------------------------------------------------

// AssetCategoryDepreciationRoutes holds route paths for the per-category and
// per-policy depreciation drawer endpoints (Surfaces C and F).
type AssetCategoryDepreciationRoutes struct {
	// CategoryRunURL is the Surface C per-category run drawer.
	CategoryRunURL string `json:"category_run_url"`

	// PolicyRunURL is the Surface C per-policy run drawer (policy breadcrumb).
	PolicyRunURL string `json:"policy_run_url"`

	// PolicyPreviewURL is the Surface F preview drawer (no writes).
	PolicyPreviewURL string `json:"policy_preview_url"`

	// DepreciationPoliciesURL is the Surface F actionable policies page.
	DepreciationPoliciesURL string `json:"depreciation_policies_url"`
}

// DefaultAssetCategoryDepreciationRoutes returns an
// AssetCategoryDepreciationRoutes populated from the package-level route
// constants.
func DefaultAssetCategoryDepreciationRoutes() AssetCategoryDepreciationRoutes {
	return AssetCategoryDepreciationRoutes{
		CategoryRunURL:          AssetCategoryDepreciationRunURL,
		PolicyRunURL:            AssetPolicyDepreciationRunURL,
		PolicyPreviewURL:        AssetPolicyDepreciationPreviewURL,
		DepreciationPoliciesURL: DepreciationPoliciesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r AssetCategoryDepreciationRoutes) RouteMap() map[string]string {
	return map[string]string{
		"asset_category_depreciation.category_run":          r.CategoryRunURL,
		"asset_category_depreciation.policy_run":            r.PolicyRunURL,
		"asset_category_depreciation.policy_preview":        r.PolicyPreviewURL,
		"asset_category_depreciation.depreciation_policies": r.DepreciationPoliciesURL,
	}
}

// CategoryRunFor returns the resolved Surface C per-category run drawer URL
// for the given category ID.
func (r AssetCategoryDepreciationRoutes) CategoryRunFor(categoryID string) string {
	return resolveParam(r.CategoryRunURL, "category_id", categoryID)
}

// PolicyRunFor returns the resolved Surface C per-policy run drawer URL for
// the given category ID (policy scope).
func (r AssetCategoryDepreciationRoutes) PolicyRunFor(categoryID string) string {
	return resolveParam(r.PolicyRunURL, "category_id", categoryID)
}

// PolicyPreviewFor returns the resolved Surface F preview drawer URL for the
// given category ID.
func (r AssetCategoryDepreciationRoutes) PolicyPreviewFor(categoryID string) string {
	return resolveParam(r.PolicyPreviewURL, "category_id", categoryID)
}

// ---------------------------------------------------------------------------
// TaxRateRoutes
// ---------------------------------------------------------------------------

// TaxRateRoutes holds route paths for Tax Rate read-only views.
// Tax rates are read-only in the UI; supersession is done via admin SQL recipe.
type TaxRateRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	DetailURL string `json:"detail_url"`
}

// DefaultTaxRateRoutes returns a TaxRateRoutes populated from the package-level
// route constants.
func DefaultTaxRateRoutes() TaxRateRoutes {
	return TaxRateRoutes{
		ActiveNav: "tax_rate",
		ListURL:   TaxRateListURL,
		DetailURL: TaxRateDetailURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r TaxRateRoutes) RouteMap() map[string]string {
	return map[string]string{
		"tax_rate.list":   r.ListURL,
		"tax_rate.detail": r.DetailURL,
	}
}

// DetailFor returns the resolved detail URL for a given tax rate ID.
func (r TaxRateRoutes) DetailFor(id string) string {
	return resolveParam(r.DetailURL, "id", id)
}

// ---------------------------------------------------------------------------
// ForexRateRoutes
// ---------------------------------------------------------------------------

// ForexRateRoutes holds route paths for Forex Rate read-only views.
// Forex rates are read-only in the UI; rows are appended only via RecordOperatorRate.
type ForexRateRoutes struct {
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	DetailURL string `json:"detail_url"`
}

// DefaultForexRateRoutes returns a ForexRateRoutes populated from the
// package-level route constants.
func DefaultForexRateRoutes() ForexRateRoutes {
	return ForexRateRoutes{
		ActiveNav: "forex_rate",
		ListURL:   ForexRateListURL,
		DetailURL: ForexRateDetailURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ForexRateRoutes) RouteMap() map[string]string {
	return map[string]string{
		"forex_rate.list":   r.ListURL,
		"forex_rate.detail": r.DetailURL,
	}
}

// DetailFor returns the resolved detail URL for a given forex rate ID.
func (r ForexRateRoutes) DetailFor(id string) string {
	return resolveParam(r.DetailURL, "id", id)
}

// ---------------------------------------------------------------------------
// WithholdingCertificateRoutes
// ---------------------------------------------------------------------------

// WithholdingCertificateRoutes holds route paths for Withholding Certificate
// CRUD views (Treasury domain — tax integration v1).
type WithholdingCertificateRoutes struct {
	ActiveNav     string `json:"active_nav"`
	ListURL       string `json:"list_url"`
	DetailURL     string `json:"detail_url"`
	TableURL      string `json:"table_url"`
	AddURL        string `json:"add_url"`
	EditURL       string `json:"edit_url"`
	DeleteURL     string `json:"delete_url"`
	BulkDeleteURL string `json:"bulk_delete_url"`
	SetStatusURL  string `json:"set_status_url"`
}

// DefaultWithholdingCertificateRoutes returns a WithholdingCertificateRoutes
// populated from the package-level route constants.
func DefaultWithholdingCertificateRoutes() WithholdingCertificateRoutes {
	return WithholdingCertificateRoutes{
		ActiveNav:     "withholding_certificate",
		ListURL:       WithholdingCertificateListURL,
		DetailURL:     WithholdingCertificateDetailURL,
		TableURL:      WithholdingCertificateTableURL,
		AddURL:        WithholdingCertificateAddURL,
		EditURL:       WithholdingCertificateEditURL,
		DeleteURL:     WithholdingCertificateDeleteURL,
		BulkDeleteURL: WithholdingCertificateBulkDeleteURL,
		SetStatusURL:  WithholdingCertificateSetStatusURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r WithholdingCertificateRoutes) RouteMap() map[string]string {
	return map[string]string{
		"withholding_certificate.list":        r.ListURL,
		"withholding_certificate.detail":      r.DetailURL,
		"withholding_certificate.table":       r.TableURL,
		"withholding_certificate.add":         r.AddURL,
		"withholding_certificate.edit":        r.EditURL,
		"withholding_certificate.delete":      r.DeleteURL,
		"withholding_certificate.bulk_delete": r.BulkDeleteURL,
		"withholding_certificate.set_status":  r.SetStatusURL,
	}
}

// DetailFor returns the resolved detail URL for a given withholding certificate ID.
func (r WithholdingCertificateRoutes) DetailFor(id string) string {
	return resolveParam(r.DetailURL, "id", id)
}

// EditFor returns the resolved edit drawer URL for a given withholding certificate ID.
func (r WithholdingCertificateRoutes) EditFor(id string) string {
	return resolveParam(r.EditURL, "id", id)
}

// ---------------------------------------------------------------------------
// resolveParam — internal URL template helper
// ---------------------------------------------------------------------------

// resolveParam replaces a single {placeholder} in a URL pattern with value.
// It is the internal single-parameter URL resolver; for multi-parameter URLs
// use route.ResolveURL from packages/pyeza-golang/route directly.
func resolveParam(pattern, placeholder, value string) string {
	token := "{" + placeholder + "}"
	n := len(token)
	for i := 0; i+n <= len(pattern); i++ {
		if pattern[i:i+n] == token {
			return pattern[:i] + value + pattern[i+n:]
		}
	}
	return pattern
}
