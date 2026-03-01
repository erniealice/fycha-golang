package fycha

// ReportsLabels holds all translatable strings for the reports module.
type ReportsLabels struct {
	GrossProfit GrossProfitLabels `json:"grossProfit"`
	Revenue     RevenueLabels     `json:"revenue"`
	CostOfSales CostOfSalesLabels `json:"costOfSales"`
	Expenses    ExpensesLabels    `json:"expenses"`
	NetProfit   NetProfitLabels   `json:"netProfit"`
	Dashboard   DashboardLabels   `json:"dashboard"`
	Period      PeriodLabels      `json:"period"`
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
	Title              string `json:"title"`
	Subtitle           string `json:"subtitle"`
	RevenueCard        string `json:"revenueCard"`
	ExpensesCard       string `json:"expensesCard"`
	NetProfitCard      string `json:"netProfitCard"`
	NetMarginCard      string `json:"netMarginCard"`
	RevenueDesc        string `json:"revenueDesc"`
	GrossProfitDesc    string `json:"grossProfitDesc"`
	CostOfSalesDesc    string `json:"costOfSalesDesc"`
	ExpensesDesc       string `json:"expensesDesc"`
	NetProfitDesc      string `json:"netProfitDesc"`
	ViewReport         string `json:"viewReport"`
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
}

// CostOfSalesLabels holds translatable strings for the cost of sales report.
type CostOfSalesLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Summary
	SummaryTotalCOGS string `json:"summaryTotalCogs"`
	SummaryRevenue   string `json:"summaryRevenue"`
	SummaryCOGSRatio string `json:"summaryCosRatio"`
	SummaryUnits     string `json:"summaryUnits"`
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
}

// NetProfitLabels holds translatable strings for the net profit report.
type NetProfitLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// P&L line items
	Revenue        string `json:"revenue"`
	CostOfSales    string `json:"costOfSales"`
	GrossProfit    string `json:"grossProfit"`
	GrossMargin    string `json:"grossMargin"`
	Expenses       string `json:"expenses"`
	NetProfit      string `json:"netProfit"`
	NetMargin      string `json:"netMargin"`
	// Summary
	SummaryRevenue    string `json:"summaryRevenue"`
	SummaryGross      string `json:"summaryGrossProfit"`
	SummaryExpenses   string `json:"summaryExpenses"`
	SummaryNetProfit  string `json:"summaryNetProfit"`
}
