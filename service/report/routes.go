package report

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

	// Cash report routes
	CashBookURL = "/cash/reports/cash-book"

	// Reports — Financial Statements (business-stakeholder output)
	ReportsIncomeStatementURL = "/reports/income-statement"
	ReportsBalanceSheetURL    = "/reports/balance-sheet"
	ReportsCashFlowURL        = "/reports/cash-flow"
	ReportsEquityChangesURL   = "/reports/equity-changes"
)

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
		DashboardURL:                     ReportsDashboardURL,
		RevenueURL:                       ReportsRevenueURL,
		CostOfSalesURL:                   ReportsCostOfSalesURL,
		GrossProfitURL:                   ReportsGrossProfitURL,
		ExpensesURL:                      ReportsExpensesURL,
		NetProfitURL:                     ReportsNetProfitURL,
		IncomeStatementURL:               ReportsIncomeStatementURL,
		BalanceSheetURL:                  ReportsBalanceSheetURL,
		CashFlowURL:                      ReportsCashFlowURL,
		EquityChangesURL:                 ReportsEquityChangesURL,
		RevenueReportURL:                 ReportsRevenueReportURL,
		RevenueReportExportURL:           ReportsRevenueReportExportURL,
		ExpenditureReportURL:             ReportsExpenditureReportURL,
		ExpenditureReportExportURL:       ReportsExpenditureReportExportURL,
		DisbursementReportURL:            ReportsDisbursementReportURL,
		DisbursementReportExportURL:      ReportsDisbursementReportExportURL,
		ReceivablesAgingReportURL:        ReportsReceivablesAgingReportURL,
		ReceivablesAgingReportExportURL:  ReportsReceivablesAgingReportExportURL,
		PayablesAgingReportURL:           ReportsPayablesAgingReportURL,
		PayablesAgingReportExportURL:     ReportsPayablesAgingReportExportURL,
		CollectionSummaryReportURL:       ReportsCollectionSummaryReportURL,
		CollectionSummaryReportExportURL: ReportsCollectionSummaryReportExportURL,
	}
}

// RouteMap returns a map of dot-notation keys to route paths.
func (r ReportsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"reports.dashboard":                        r.DashboardURL,
		"reports.revenue":                          r.RevenueURL,
		"reports.cost_of_sales":                    r.CostOfSalesURL,
		"reports.gross_profit":                     r.GrossProfitURL,
		"reports.expenses":                         r.ExpensesURL,
		"reports.net_profit":                       r.NetProfitURL,
		"reports.income_statement":                 r.IncomeStatementURL,
		"reports.balance_sheet":                    r.BalanceSheetURL,
		"reports.cash_flow":                        r.CashFlowURL,
		"reports.equity_changes":                   r.EquityChangesURL,
		"reports.revenue_report":                   r.RevenueReportURL,
		"reports.revenue_report_export":            r.RevenueReportExportURL,
		"reports.expenditure_report":               r.ExpenditureReportURL,
		"reports.expenditure_report_export":        r.ExpenditureReportExportURL,
		"reports.disbursement_report":              r.DisbursementReportURL,
		"reports.disbursement_report_export":       r.DisbursementReportExportURL,
		"reports.receivables_aging_report":         r.ReceivablesAgingReportURL,
		"reports.receivables_aging_report_export":  r.ReceivablesAgingReportExportURL,
		"reports.payables_aging_report":            r.PayablesAgingReportURL,
		"reports.payables_aging_report_export":     r.PayablesAgingReportExportURL,
		"reports.collection_summary_report":        r.CollectionSummaryReportURL,
		"reports.collection_summary_report_export": r.CollectionSummaryReportExportURL,
	}
}
