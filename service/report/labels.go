package report

// ReportsLabels holds all translatable strings for the reports module.
type ReportsLabels struct {
	GrossProfit        GrossProfitLabels             `json:"grossProfit"`
	Revenue            RevenueLabels                 `json:"revenue"`
	RevenueReport      RevenueReportLabels           `json:"revenueReport"`
	ExpenditureReport  ExpenditureReportLabels       `json:"expenditureReport"`
	DisbursementReport DisbursementReportLabels      `json:"disbursementReport"`
	ReceivablesAging   ReceivablesAgingReportLabels  `json:"receivablesAging"`
	PayablesAging      PayablesAgingReportLabels     `json:"payablesAging"`
	CollectionSummary  CollectionSummaryReportLabels `json:"collectionSummary"`
	CostOfSales        CostOfSalesLabels             `json:"costOfSales"`
	Expenses           ExpensesLabels                `json:"expenses"`
	NetProfit          NetProfitLabels               `json:"netProfit"`
	Dashboard          DashboardLabels               `json:"dashboard"`
	Period             PeriodLabels                  `json:"period"`
	IncomeStatement    IncomeStatementLabels         `json:"incomeStatement"`
	BalanceSheet       BalanceSheetLabels            `json:"balanceSheet"`
	CashFlow           CashFlowLabels                `json:"cashFlow"`
	EquityChanges      EquityChangesLabels           `json:"equityChanges"`
	TrialBalance       TrialBalanceLabels            `json:"trialBalance"`
}

// TrialBalanceLabels holds translatable strings for the Trial Balance report page.
type TrialBalanceLabels struct {
	Title        string `json:"title"`
	Subtitle     string `json:"subtitle"`
	Totals       string `json:"totals"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	// Filter bar
	FilterAsOfDate string `json:"filterAsOfDate"`
	Generate       string `json:"generate"`
	Print          string `json:"print"`
	// Balance indicator
	BalancedTitle   string `json:"balancedTitle"`
	BalancedMessage string `json:"balancedMessage"`
	UnbalancedTitle string `json:"unbalancedTitle"`
	DebitsLabel     string `json:"debitsLabel"`
	CreditsLabel    string `json:"creditsLabel"`
	DifferenceLabel string `json:"differenceLabel"`
	AsOfPrefix      string `json:"asOfPrefix"`
	GeneratePrompt  string `json:"generatePrompt"`
	// Domain-specific table column headers (Code/Account-Name fall back to the
	// shared common Columns block; these two are statement-specific).
	ColAccountName   string `json:"colAccountName"`
	ColDebitBalance  string `json:"colDebitBalance"`
	ColCreditBalance string `json:"colCreditBalance"`
}

// IncomeStatementLabels holds translatable strings for the Income Statement page.
type IncomeStatementLabels struct {
	Title    string                       `json:"title"`
	Subtitle string                       `json:"subtitle"`
	Sections IncomeStatementSectionLabels `json:"sections"`
	// Filter bar
	Showing    string `json:"showing"`
	Export     string `json:"export"`
	ExportSoon string `json:"exportSoon"`
	// KPI summary
	KpiTotalRevenue  string `json:"kpiTotalRevenue"`
	KpiTotalExpenses string `json:"kpiTotalExpenses"`
	KpiNetIncome     string `json:"kpiNetIncome"`
	// Statement body
	StatementTitle string `json:"statementTitle"`
	ColCode        string `json:"colCode"`
	ColAccount     string `json:"colAccount"`
	ColThisPeriod  string `json:"colThisPeriod"`
	ColPriorPeriod string `json:"colPriorPeriod"`
	ColChange      string `json:"colChange"`
	SubtotalPrefix string `json:"subtotalPrefix"`
	TotalPrefix    string `json:"totalPrefix"`
}

// IncomeStatementSectionLabels holds the section titles for the income statement.
type IncomeStatementSectionLabels struct {
	Revenue           string `json:"revenue"`
	CostOfSales       string `json:"costOfSales"`
	GrossProfit       string `json:"grossProfit"`
	OperatingExpenses string `json:"operatingExpenses"`
	SellingExpenses   string `json:"sellingExpenses"`
	GeneralAdmin      string `json:"generalAdmin"`
	OperatingIncome   string `json:"operatingIncome"`
	OtherExpenses     string `json:"otherExpenses"`
	NetIncome         string `json:"netIncome"`
}

// BalanceSheetLabels holds translatable strings for the Balance Sheet page.
type BalanceSheetLabels struct {
	Title    string                    `json:"title"`
	Subtitle string                    `json:"subtitle"`
	Sections BalanceSheetSectionLabels `json:"sections"`
	// Filter bar
	AsOf       string `json:"asOf"`
	Generate   string `json:"generate"`
	Export     string `json:"export"`
	ExportSoon string `json:"exportSoon"`
	// KPI summary
	KpiTotalAssets      string `json:"kpiTotalAssets"`
	KpiTotalLiabilities string `json:"kpiTotalLiabilities"`
	KpiTotalEquity      string `json:"kpiTotalEquity"`
	// Statement body
	StatementTitle     string `json:"statementTitle"`
	AsOfPrefix         string `json:"asOfPrefix"`
	ColCode            string `json:"colCode"`
	ColAccount         string `json:"colAccount"`
	ColAmount          string `json:"colAmount"`
	TotalPrefix        string `json:"totalPrefix"`
	TotalUpperPrefix   string `json:"totalUpperPrefix"`
	TotalLiabAndEquity string `json:"totalLiabAndEquity"`
	// Accounting-equation verification banner. Both carry three %s verbs
	// (assets, liabilities/diff, equity/landE) — every business-type override
	// MUST preserve all three %s placeholders or fmt.Sprintf renders EXTRA.
	EquationVerified string `json:"equationVerified"`
	EquationWarning  string `json:"equationWarning"`
}

// BalanceSheetSectionLabels holds the section and classification titles for the balance sheet.
type BalanceSheetSectionLabels struct {
	Assets                string `json:"assets"`
	CurrentAssets         string `json:"currentAssets"`
	NonCurrentAssets      string `json:"nonCurrentAssets"`
	Liabilities           string `json:"liabilities"`
	CurrentLiabilities    string `json:"currentLiabilities"`
	NonCurrentLiabilities string `json:"nonCurrentLiabilities"`
	Equity                string `json:"equity"`
}

// CashFlowLabels holds translatable strings for the Cash Flow Statement page.
type CashFlowLabels struct {
	Title    string                `json:"title"`
	Subtitle string                `json:"subtitle"`
	Sections CashFlowSectionLabels `json:"sections"`
	// Filter bar
	Showing    string `json:"showing"`
	Export     string `json:"export"`
	ExportSoon string `json:"exportSoon"`
	// KPI summary
	KpiOperating string `json:"kpiOperating"`
	KpiNetChange string `json:"kpiNetChange"`
	KpiEnding    string `json:"kpiEnding"`
	// Statement body
	StatementTitle        string `json:"statementTitle"`
	ColDescription        string `json:"colDescription"`
	ColAmount             string `json:"colAmount"`
	NetChangeRow          string `json:"netChangeRow"`
	CashBeginning         string `json:"cashBeginning"`
	CashEnding            string `json:"cashEnding"`
	Reconciliation        string `json:"reconciliation"`
	TotalCashBalanceSheet string `json:"totalCashBalanceSheet"`
	Verified              string `json:"verified"`
	Mismatch              string `json:"mismatch"`
}

// CashFlowSectionLabels holds the activity section titles for the cash flow statement.
type CashFlowSectionLabels struct {
	OperatingActivities string `json:"operatingActivities"`
	InvestingActivities string `json:"investingActivities"`
	FinancingActivities string `json:"financingActivities"`
	NetCashOperating    string `json:"netCashOperating"`
	NetCashInvesting    string `json:"netCashInvesting"`
	NetCashFinancing    string `json:"netCashFinancing"`
}

// EquityChangesLabels holds translatable strings for the Statement of Changes in Equity page.
type EquityChangesLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

// PeriodLabels holds shared period preset labels used across all reports.
type PeriodLabels struct {
	ThisMonth   string `json:"thisMonth"`
	LastMonth   string `json:"lastMonth"`
	ThisQuarter string `json:"thisQuarter"`
	LastQuarter string `json:"lastQuarter"`
	ThisYear    string `json:"thisYear"`
	LastYear    string `json:"lastYear"`
	Custom      string `json:"custom"`
	DateStart   string `json:"dateStart"`
	DateEnd     string `json:"dateEnd"`
	GroupBy     string `json:"groupBy"`
}

// DashboardLabels holds translatable strings for the reports dashboard.
type DashboardLabels struct {
	Title           string `json:"title"`
	Subtitle        string `json:"subtitle"`
	RevenueCard     string `json:"revenueCard"`
	ExpensesCard    string `json:"expensesCard"`
	NetProfitCard   string `json:"netProfitCard"`
	NetMarginCard   string `json:"netMarginCard"`
	RevenueDesc     string `json:"revenueDesc"`
	GrossProfitDesc string `json:"grossProfitDesc"`
	CostOfSalesDesc string `json:"costOfSalesDesc"`
	ExpensesDesc    string `json:"expensesDesc"`
	NetProfitDesc   string `json:"netProfitDesc"`
	ViewReport      string `json:"viewReport"`
}

// GrossProfitLabels holds translatable strings for the gross profit report.
type GrossProfitLabels struct {
	Title              string `json:"title"`
	RevenueGroup       string `json:"revenueGroup"`
	ProfitabilityGroup string `json:"profitabilityGroup"`
	VolumeGroup        string `json:"volumeGroup"`
	GrossRevenue       string `json:"grossRevenue"`
	Discount           string `json:"discount"`
	NetRevenue         string `json:"netRevenue"`
	COGS               string `json:"cogs"`
	GrossProfit        string `json:"profit"`
	Margin             string `json:"margin"`
	UnitsSold          string `json:"unitsSold"`
	Transactions       string `json:"transactions"`
	// Group by
	GroupBy          string `json:"groupBy"`
	GroupByProduct   string `json:"groupByProduct"`
	GroupByLocation  string `json:"groupByLocation"`
	GroupByCategory  string `json:"groupByCategory"`
	GroupByMonthly   string `json:"groupByMonthly"`
	GroupByQuarterly string `json:"groupByQuarterly"`
	// Filters
	FilterProduct  string `json:"filterProduct"`
	FilterLocation string `json:"filterLocation"`
	FilterCategory string `json:"filterCategory"`
	FilterAll      string `json:"filterAll"`
	// Period presets
	PeriodThisMonth   string `json:"periodThisMonth"`
	PeriodLastMonth   string `json:"periodLastMonth"`
	PeriodThisQuarter string `json:"periodThisQuarter"`
	PeriodLastQuarter string `json:"periodLastQuarter"`
	PeriodThisYear    string `json:"periodThisYear"`
	PeriodLastYear    string `json:"periodLastYear"`
	PeriodCustom      string `json:"periodCustom"`
	DateStart         string `json:"dateStart"`
	DateEnd           string `json:"dateEnd"`
	Apply             string `json:"apply"`
	// Summary
	SummaryNetRevenue  string `json:"summaryNetRevenue"`
	SummaryCogs        string `json:"summaryCogs"`
	SummaryGrossProfit string `json:"summaryGrossProfit"`
	SummaryMargin      string `json:"summaryMargin"`
	// Summary row + empty state
	Totals       string `json:"totals"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// RevenueLabels holds translatable strings for the revenue report.
type RevenueLabels struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Reference string `json:"reference"`
	Customer  string `json:"customer"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
	// Summary
	SummaryTotal        string `json:"summaryTotal"`
	SummaryTransactions string `json:"summaryTransactions"`
	SummaryAverage      string `json:"summaryAverage"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// CostOfSalesLabels holds translatable strings for the cost of sales report.
type CostOfSalesLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Column headers
	Item       string `json:"item"`
	COGS       string `json:"cogs"`
	NetRevenue string `json:"netRevenue"`
	COGSPct    string `json:"cogsPct"`
	Units      string `json:"units"`
	// Summary
	SummaryTotalCOGS string `json:"summaryTotalCogs"`
	SummaryRevenue   string `json:"summaryRevenue"`
	SummaryCOGSRatio string `json:"summaryCosRatio"`
	SummaryUnits     string `json:"summaryUnits"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// ExpensesLabels holds translatable strings for the expenses report.
type ExpensesLabels struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Reference string `json:"reference"`
	Vendor    string `json:"vendor"`
	Category  string `json:"category"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	Status    string `json:"status"`
	// Summary
	SummaryTotal    string `json:"summaryTotal"`
	SummaryCount    string `json:"summaryCount"`
	SummaryApproved string `json:"summaryApproved"`
	SummaryPending  string `json:"summaryPending"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// NetProfitLabels holds translatable strings for the net profit report.
type NetProfitLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// P&L line items
	Revenue     string `json:"revenue"`
	CostOfSales string `json:"costOfSales"`
	GrossProfit string `json:"grossProfit"`
	GrossMargin string `json:"grossMargin"`
	Expenses    string `json:"expenses"`
	NetProfit   string `json:"netProfit"`
	NetMargin   string `json:"netMargin"`
	// Summary
	SummaryRevenue   string `json:"summaryRevenue"`
	SummaryGross     string `json:"summaryGrossProfit"`
	SummaryExpenses  string `json:"summaryExpenses"`
	SummaryNetProfit string `json:"summaryNetProfit"`
}

// RevenueReportLabels holds translatable strings for the revenue pivot-table report.
type RevenueReportLabels struct {
	Title                   string `json:"title"`
	Subtitle                string `json:"subtitle"`
	ColumnDimension         string `json:"columnDimension"`
	RowDimension            string `json:"rowDimension"`
	DimensionMonthly        string `json:"dimensionMonthly"`
	DimensionQuarterly      string `json:"dimensionQuarterly"`
	DimensionYearly         string `json:"dimensionYearly"`
	DimensionProduct        string `json:"dimensionProduct"`
	DimensionProductLine    string `json:"dimensionProductLine"`
	DimensionLocation       string `json:"dimensionLocation"`
	DimensionLocationArea   string `json:"dimensionLocationArea"`
	DimensionClient         string `json:"dimensionClient"`
	DimensionClientCategory string `json:"dimensionClientCategory"`
	SummaryGrandTotal       string `json:"summaryGrandTotal"`
	SummaryTransactions     string `json:"summaryTransactions"`
	SummaryAverage          string `json:"summaryAverage"`
	Total                   string `json:"total"`
	Totals                  string `json:"totals"`
	ExportCsv               string `json:"exportCsv"`
	Apply                   string `json:"apply"`
	Clear                   string `json:"clear"`
	EmptyTitle              string `json:"emptyTitle"`
	EmptyMessage            string `json:"emptyMessage"`
}

// PrimaryGroupLabel returns the display label for the given dimension string.
func (l RevenueReportLabels) PrimaryGroupLabel(dim string) string {
	switch dim {
	case "monthly":
		return l.DimensionMonthly
	case "quarterly":
		return l.DimensionQuarterly
	case "yearly":
		return l.DimensionYearly
	case "product":
		return l.DimensionProduct
	case "productLine":
		return l.DimensionProductLine
	case "location":
		return l.DimensionLocation
	case "locationArea":
		return l.DimensionLocationArea
	case "client", "client_category":
		return l.DimensionClient
	case "clientCategory":
		return l.DimensionClientCategory
	default:
		return dim
	}
}

// DimensionOptions returns all seven dimension choices as FilterOption slices.
func (l RevenueReportLabels) DimensionOptions(active string) []FilterOption {
	dims := []struct {
		value string
		label string
	}{
		{"monthly", l.DimensionMonthly},
		{"quarterly", l.DimensionQuarterly},
		{"yearly", l.DimensionYearly},
		{"product", l.DimensionProduct},
		{"productLine", l.DimensionProductLine},
		{"location", l.DimensionLocation},
		{"locationArea", l.DimensionLocationArea},
		{"client", l.DimensionClient},
		{"clientCategory", l.DimensionClientCategory},
	}
	opts := make([]FilterOption, len(dims))
	for i, d := range dims {
		opts[i] = FilterOption{Value: d.value, Label: d.label, Selected: d.value == active}
	}
	return opts
}

// ExpenditureReportLabels holds labels for the expenditure report pivot table view.
type ExpenditureReportLabels struct {
	Title                    string `json:"title"`
	Subtitle                 string `json:"subtitle"`
	ColumnDimension          string `json:"columnDimension"`
	RowDimension             string `json:"rowDimension"`
	DimensionMonthly         string `json:"dimensionMonthly"`
	DimensionQuarterly       string `json:"dimensionQuarterly"`
	DimensionYearly          string `json:"dimensionYearly"`
	DimensionProduct         string `json:"dimensionProduct"`
	DimensionProductLine     string `json:"dimensionProductLine"`
	DimensionLocation        string `json:"dimensionLocation"`
	DimensionLocationArea    string `json:"dimensionLocationArea"`
	DimensionCategory        string `json:"dimensionCategory"`
	DimensionSupplier        string `json:"dimensionSupplier"`
	DimensionExpenditureType string `json:"dimensionExpenditureType"`
	SummaryGrandTotal        string `json:"summaryGrandTotal"`
	SummaryTransactions      string `json:"summaryTransactions"`
	SummaryAverage           string `json:"summaryAverage"`
	Total                    string `json:"total"`
	Totals                   string `json:"totals"`
	ExportCsv                string `json:"exportCsv"`
	Apply                    string `json:"apply"`
	Clear                    string `json:"clear"`
	EmptyTitle               string `json:"emptyTitle"`
	EmptyMessage             string `json:"emptyMessage"`
}

// PrimaryGroupLabel returns the display label for the given dimension string.
func (l ExpenditureReportLabels) PrimaryGroupLabel(dim string) string {
	switch dim {
	case "monthly":
		return l.DimensionMonthly
	case "quarterly":
		return l.DimensionQuarterly
	case "yearly":
		return l.DimensionYearly
	case "product":
		return l.DimensionProduct
	case "productLine":
		return l.DimensionProductLine
	case "location":
		return l.DimensionLocation
	case "locationArea":
		return l.DimensionLocationArea
	case "category":
		return l.DimensionCategory
	case "supplier":
		return l.DimensionSupplier
	case "expenditureType":
		return l.DimensionExpenditureType
	default:
		return dim
	}
}

// DimensionOptions returns all dimension choices as FilterOption slices.
func (l ExpenditureReportLabels) DimensionOptions(active string) []FilterOption {
	dims := []struct {
		value string
		label string
	}{
		{"monthly", l.DimensionMonthly},
		{"quarterly", l.DimensionQuarterly},
		{"yearly", l.DimensionYearly},
		{"product", l.DimensionProduct},
		{"productLine", l.DimensionProductLine},
		{"location", l.DimensionLocation},
		{"locationArea", l.DimensionLocationArea},
		{"category", l.DimensionCategory},
		{"supplier", l.DimensionSupplier},
		{"expenditureType", l.DimensionExpenditureType},
	}
	opts := make([]FilterOption, len(dims))
	for i, d := range dims {
		opts[i] = FilterOption{Value: d.value, Label: d.label, Selected: d.value == active}
	}
	return opts
}

// DisbursementReportLabels holds labels for the disbursement report pivot table view.
type DisbursementReportLabels struct {
	Title                        string `json:"title"`
	Subtitle                     string `json:"subtitle"`
	ColumnDimension              string `json:"columnDimension"`
	RowDimension                 string `json:"rowDimension"`
	DimensionMonthly             string `json:"dimensionMonthly"`
	DimensionQuarterly           string `json:"dimensionQuarterly"`
	DimensionYearly              string `json:"dimensionYearly"`
	DimensionSupplier            string `json:"dimensionSupplier"`
	DimensionSupplierCategory    string `json:"dimensionSupplierCategory"`
	DimensionLocation            string `json:"dimensionLocation"`
	DimensionLocationArea        string `json:"dimensionLocationArea"`
	DimensionExpenditureCategory string `json:"dimensionExpenditureCategory"`
	DimensionDisbursementType    string `json:"dimensionDisbursementType"`
	DimensionDisbursementMethod  string `json:"dimensionDisbursementMethod"`
	SummaryGrandTotal            string `json:"summaryGrandTotal"`
	SummaryTransactions          string `json:"summaryTransactions"`
	SummaryAverage               string `json:"summaryAverage"`
	Total                        string `json:"total"`
	Totals                       string `json:"totals"`
	ExportCsv                    string `json:"exportCsv"`
	Apply                        string `json:"apply"`
	Clear                        string `json:"clear"`
	EmptyTitle                   string `json:"emptyTitle"`
	EmptyMessage                 string `json:"emptyMessage"`
}

// PrimaryGroupLabel returns the display label for the given dimension string.
func (l DisbursementReportLabels) PrimaryGroupLabel(dim string) string {
	switch dim {
	case "monthly":
		return l.DimensionMonthly
	case "quarterly":
		return l.DimensionQuarterly
	case "yearly":
		return l.DimensionYearly
	case "supplier":
		return l.DimensionSupplier
	case "supplierCategory":
		return l.DimensionSupplierCategory
	case "location":
		return l.DimensionLocation
	case "locationArea":
		return l.DimensionLocationArea
	case "expenditureCategory":
		return l.DimensionExpenditureCategory
	case "disbursementType":
		return l.DimensionDisbursementType
	case "disbursementMethod":
		return l.DimensionDisbursementMethod
	default:
		return dim
	}
}

// DimensionOptions returns all dimension choices as FilterOption slices.
func (l DisbursementReportLabels) DimensionOptions(active string) []FilterOption {
	dims := []struct {
		value string
		label string
	}{
		{"monthly", l.DimensionMonthly},
		{"quarterly", l.DimensionQuarterly},
		{"yearly", l.DimensionYearly},
		{"supplier", l.DimensionSupplier},
		{"supplierCategory", l.DimensionSupplierCategory},
		{"location", l.DimensionLocation},
		{"locationArea", l.DimensionLocationArea},
		{"expenditureCategory", l.DimensionExpenditureCategory},
		{"disbursementType", l.DimensionDisbursementType},
		{"disbursementMethod", l.DimensionDisbursementMethod},
	}
	opts := make([]FilterOption, len(dims))
	for i, d := range dims {
		opts[i] = FilterOption{Value: d.value, Label: d.label, Selected: d.value == active}
	}
	return opts
}

// ReceivablesAgingReportLabels holds translatable strings for the receivables aging report.
type ReceivablesAgingReportLabels struct {
	PageTitle               string `json:"page_title"`
	PageDescription         string `json:"page_description"`
	BucketCurrent           string `json:"bucket_current"`
	Bucket1To30             string `json:"bucket_1_to_30"`
	Bucket31To60            string `json:"bucket_31_to_60"`
	Bucket61To90            string `json:"bucket_61_to_90"`
	BucketOver90            string `json:"bucket_over_90"`
	TotalOutstanding        string `json:"total_outstanding"`
	InvoiceCount            string `json:"invoice_count"`
	SummaryGrandTotal       string `json:"summary_grand_total"`
	SummaryInvoiceCount     string `json:"summary_invoice_count"`
	SummaryOverdueAmount    string `json:"summary_overdue_amount"`
	EmptyTitle              string `json:"empty_title"`
	EmptyMessage            string `json:"empty_message"`
	ExportFilename          string `json:"export_filename"`
	FilterAsOfDate          string `json:"filter_as_of_date"`
	FilterRowDimension      string `json:"filter_row_dimension"`
	DimensionClient         string `json:"dimension_client"`
	DimensionClientCategory string `json:"dimension_client_category"`
	DimensionLocation       string `json:"dimension_location"`
	DimensionLocationArea   string `json:"dimension_location_area"`
}

// PrimaryGroupLabel returns the display label for the given dimension string.
func (l ReceivablesAgingReportLabels) PrimaryGroupLabel(dim string) string {
	switch dim {
	case "client":
		return l.DimensionClient
	case "clientCategory":
		return l.DimensionClientCategory
	case "location":
		return l.DimensionLocation
	case "locationArea":
		return l.DimensionLocationArea
	default:
		return dim
	}
}

// PayablesAgingReportLabels holds translatable strings for the payables aging report.
type PayablesAgingReportLabels struct {
	PageTitle                    string `json:"page_title"`
	PageDescription              string `json:"page_description"`
	BucketCurrent                string `json:"bucket_current"`
	Bucket1To30                  string `json:"bucket_1_to_30"`
	Bucket31To60                 string `json:"bucket_31_to_60"`
	Bucket61To90                 string `json:"bucket_61_to_90"`
	BucketOver90                 string `json:"bucket_over_90"`
	TotalOutstanding             string `json:"total_outstanding"`
	InvoiceCount                 string `json:"invoice_count"`
	SummaryGrandTotal            string `json:"summary_grand_total"`
	SummaryInvoiceCount          string `json:"summary_invoice_count"`
	SummaryOverdueAmount         string `json:"summary_overdue_amount"`
	EmptyTitle                   string `json:"empty_title"`
	EmptyMessage                 string `json:"empty_message"`
	ExportFilename               string `json:"export_filename"`
	FilterAsOfDate               string `json:"filter_as_of_date"`
	FilterRowDimension           string `json:"filter_row_dimension"`
	DimensionSupplier            string `json:"dimension_supplier"`
	DimensionSupplierCategory    string `json:"dimension_supplier_category"`
	DimensionLocation            string `json:"dimension_location"`
	DimensionLocationArea        string `json:"dimension_location_area"`
	DimensionExpenditureCategory string `json:"dimension_expenditure_category"`
}

// PrimaryGroupLabel returns the display label for the given dimension string.
func (l PayablesAgingReportLabels) PrimaryGroupLabel(dim string) string {
	switch dim {
	case "supplier":
		return l.DimensionSupplier
	case "supplierCategory":
		return l.DimensionSupplierCategory
	case "location":
		return l.DimensionLocation
	case "locationArea":
		return l.DimensionLocationArea
	case "expenditureCategory":
		return l.DimensionExpenditureCategory
	default:
		return dim
	}
}

// CollectionSummaryReportLabels holds translatable strings for the collection summary pivot-table report.
type CollectionSummaryReportLabels struct {
	PageTitle                 string `json:"page_title"`
	PageDescription           string `json:"page_description"`
	DimensionMonthly          string `json:"dimensionMonthly"`
	DimensionQuarterly        string `json:"dimensionQuarterly"`
	DimensionYearly           string `json:"dimensionYearly"`
	DimensionLocation         string `json:"dimensionLocation"`
	DimensionLocationArea     string `json:"dimensionLocationArea"`
	DimensionClient           string `json:"dimensionClient"`
	DimensionClientCategory   string `json:"dimensionClientCategory"`
	DimensionCollectionMethod string `json:"dimensionCollectionMethod"`
	DimensionCollectionType   string `json:"dimensionCollectionType"`
	SummaryGrandTotal         string `json:"summaryGrandTotal"`
	SummaryTransactions       string `json:"summaryTransactions"`
	SummaryAverage            string `json:"summaryAverage"`
	Total                     string `json:"total"`
	Totals                    string `json:"totals"`
	ExportCsv                 string `json:"exportCsv"`
	Apply                     string `json:"apply"`
	Clear                     string `json:"clear"`
	EmptyTitle                string `json:"emptyTitle"`
	EmptyMessage              string `json:"emptyMessage"`
}

// DimensionOptions returns all dimension choices as FilterOption slices.
func (l CollectionSummaryReportLabels) DimensionOptions(active string) []FilterOption {
	dims := []struct {
		value string
		label string
	}{
		{"monthly", l.DimensionMonthly},
		{"quarterly", l.DimensionQuarterly},
		{"yearly", l.DimensionYearly},
		{"location", l.DimensionLocation},
		{"locationArea", l.DimensionLocationArea},
		{"client", l.DimensionClient},
		{"clientCategory", l.DimensionClientCategory},
		{"collectionMethod", l.DimensionCollectionMethod},
		{"collectionType", l.DimensionCollectionType},
	}
	opts := make([]FilterOption, len(dims))
	for i, d := range dims {
		opts[i] = FilterOption{Value: d.value, Label: d.label, Selected: d.value == active}
	}
	return opts
}

// PrimaryGroupLabel returns the display label for the given dimension string.
func (l CollectionSummaryReportLabels) PrimaryGroupLabel(dim string) string {
	switch dim {
	case "monthly":
		return l.DimensionMonthly
	case "quarterly":
		return l.DimensionQuarterly
	case "yearly":
		return l.DimensionYearly
	case "location":
		return l.DimensionLocation
	case "locationArea":
		return l.DimensionLocationArea
	case "client":
		return l.DimensionClient
	case "clientCategory":
		return l.DimensionClientCategory
	case "collectionMethod":
		return l.DimensionCollectionMethod
	case "collectionType":
		return l.DimensionCollectionType
	default:
		return dim
	}
}
