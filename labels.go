package fycha

// ReportsLabels holds all translatable strings for the reports module.
type ReportsLabels struct {
	GrossProfit GrossProfitLabels `json:"grossProfit"`
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
