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
}

// DefaultReportsRoutes returns a ReportsRoutes populated from package-level consts.
func DefaultReportsRoutes() ReportsRoutes {
	return ReportsRoutes{
		DashboardURL:       ReportsDashboardURL,
		RevenueURL:         ReportsRevenueURL,
		CostOfSalesURL:     ReportsCostOfSalesURL,
		GrossProfitURL:     ReportsGrossProfitURL,
		ExpensesURL:        ReportsExpensesURL,
		NetProfitURL:       ReportsNetProfitURL,
		IncomeStatementURL: ReportsIncomeStatementURL,
		BalanceSheetURL:    ReportsBalanceSheetURL,
		CashFlowURL:        ReportsCashFlowURL,
		EquityChangesURL:   ReportsEquityChangesURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ReportsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"reports.dashboard":        r.DashboardURL,
		"reports.revenue":          r.RevenueURL,
		"reports.cost_of_sales":    r.CostOfSalesURL,
		"reports.gross_profit":     r.GrossProfitURL,
		"reports.expenses":         r.ExpensesURL,
		"reports.net_profit":       r.NetProfitURL,
		"reports.income_statement": r.IncomeStatementURL,
		"reports.balance_sheet":    r.BalanceSheetURL,
		"reports.cash_flow":        r.CashFlowURL,
		"reports.equity_changes":   r.EquityChangesURL,
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

	// Report/settings routes
	LapsingScheduleURL      string `json:"lapsing_schedule_url"`
	DepreciationPoliciesURL string `json:"depreciation_policies_url"`
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
	}
}

// ---------------------------------------------------------------------------
// AccountRoutes
// ---------------------------------------------------------------------------

// AccountRoutes holds route paths for Chart of Accounts views.
type AccountRoutes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
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
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	PostURL      string `json:"post_url"`
	ReverseURL   string `json:"reverse_url"`
	DeleteURL    string `json:"delete_url"`
}

func DefaultJournalRoutes() JournalRoutes {
	return JournalRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "journals-draft",
		ListURL:      JournalListURL,
		DetailURL:    JournalDetailURL,
		AddURL:       JournalAddURL,
		EditURL:      JournalEditURL,
		PostURL:      JournalPostURL,
		ReverseURL:   JournalReverseURL,
		DeleteURL:    JournalDeleteURL,
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
	CloseURL     string `json:"close_url"`
}

func DefaultFiscalPeriodRoutes() FiscalPeriodRoutes {
	return FiscalPeriodRoutes{
		ActiveNav:    "ledger",
		ActiveSubNav: "fiscal-periods",
		ListURL:      FiscalPeriodListURL,
		DetailURL:    FiscalPeriodDetailURL,
		CloseURL:     FiscalPeriodCloseURL,
	}
}

func (r FiscalPeriodRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.fiscal_period.list":   r.ListURL,
		"ledger.fiscal_period.detail": r.DetailURL,
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
	ListURL         string `json:"list_url"`
	DetailURL       string `json:"detail_url"`
	AddURL          string `json:"add_url"`
	AmortizationURL string `json:"amortization_url"`
}

func DefaultLoanRoutes() LoanRoutes {
	return LoanRoutes{
		ActiveNav:       "loans",
		ListURL:         LoanListURL,
		DetailURL:       LoanDetailURL,
		AddURL:          LoanAddURL,
		AmortizationURL: LoanAmortizationURL,
	}
}

func (r LoanRoutes) RouteMap() map[string]string {
	return map[string]string{
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
		ActiveNav: "loans",
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
	AccountsURL       string `json:"accounts_url"`
	AccountDetailURL  string `json:"account_detail_url"`
	TransactionsURL   string `json:"transactions_url"`
	TransactionAddURL string `json:"transaction_add_url"`
}

func DefaultEquityRoutes() EquityRoutes {
	return EquityRoutes{
		ActiveNav:         "equity",
		AccountsURL:       EquityAccountsURL,
		AccountDetailURL:  EquityAccountDetailURL,
		TransactionsURL:   EquityTransactionsURL,
		TransactionAddURL: EquityTransactionAddURL,
	}
}

func (r EquityRoutes) RouteMap() map[string]string {
	return map[string]string{
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
	ActiveNav string `json:"active_nav"`
	ListURL   string `json:"list_url"`
	DetailURL string `json:"detail_url"`
}

func DefaultPayrollRunRoutes() PayrollRunRoutes {
	return PayrollRunRoutes{
		ActiveNav: "payroll",
		ListURL:   PayrollRunListURL,
		DetailURL: PayrollRunDetailURL,
	}
}

func (r PayrollRunRoutes) RouteMap() map[string]string {
	return map[string]string{
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
		ActiveNav:       "expenses",
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
