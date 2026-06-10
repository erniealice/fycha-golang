package fycha

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

// ---------------------------------------------------------------------------
// Asset labels
// ---------------------------------------------------------------------------

// AssetLabels holds all translatable strings for the fixed asset module.
type AssetLabels struct {
	Page      AssetPageLabels      `json:"page"`
	Buttons   AssetButtonLabels    `json:"buttons"`
	Columns   AssetColumnLabels    `json:"columns"`
	Empty     AssetEmptyLabels     `json:"empty"`
	Form      AssetFormLabels      `json:"form"`
	Actions   AssetActionLabels    `json:"actions"`
	Detail    AssetDetailLabels    `json:"detail"`
	Dashboard AssetDashboardLabels `json:"dashboard"`
}

type AssetPageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type AssetButtonLabels struct {
	AddAsset string `json:"addAsset"`
}

type AssetColumnLabels struct {
	AssetNumber     string `json:"assetNumber"`
	Name            string `json:"name"`
	Category        string `json:"category"`
	Location        string `json:"location"`
	AcquisitionCost string `json:"acquisitionCost"`
	BookValue       string `json:"bookValue"`
	Status          string `json:"status"`
	// Sub-table columns (depreciation)
	Period       string `json:"period"`
	StartValue   string `json:"startValue"`
	Depreciation string `json:"depreciation"`
	EndValue     string `json:"endValue"`
	Accumulated  string `json:"accumulated"`
	// Sub-table columns (maintenance)
	Date        string `json:"date"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Cost        string `json:"cost"`
	// Sub-table columns (transactions)
	Amount    string `json:"amount"`
	Reference string `json:"reference"`
	// Cost of sales columns
	Item       string `json:"item"`
	COGS       string `json:"cogs"`
	NetRevenue string `json:"netRevenue"`
	COGSPct    string `json:"cogsPct"`
	Units      string `json:"units"`
	// Summary row
	Totals string `json:"totals"`
}

type AssetEmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type AssetFormLabels struct {
	Name                       string `json:"name"`
	NamePlaceholder            string `json:"namePlaceholder"`
	AssetNumber                string `json:"assetNumber"`
	AssetNumberPlaceholder     string `json:"assetNumberPlaceholder"`
	Description                string `json:"description"`
	DescriptionPlaceholder     string `json:"descriptionPlaceholder"`
	Category                   string `json:"category"`
	CategoryPlaceholder        string `json:"categoryPlaceholder"`
	Location                   string `json:"location"`
	LocationPlaceholder        string `json:"locationPlaceholder"`
	AcquisitionCost            string `json:"acquisitionCost"`
	AcquisitionCostPlaceholder string `json:"acquisitionCostPlaceholder"`
	SalvageValue               string `json:"salvageValue"`
	SalvageValuePlaceholder    string `json:"salvageValuePlaceholder"`
	UsefulLifeMonths           string `json:"usefulLifeMonths"`
	UsefulLifePlaceholder      string `json:"usefulLifePlaceholder"`
	DepreciationMethod         string `json:"depreciationMethod"`
	Active                     string `json:"active"`
	// Depreciation method option labels
	DepMethodStraightLine      string `json:"depMethodStraightLine"`
	DepMethodDecliningBalance  string `json:"depMethodDecliningBalance"`
	DepMethodSumOfYears        string `json:"depMethodSumOfYears"`
	DepMethodUnitsOfProduction string `json:"depMethodUnitsOfProduction"`
	// Info popover text
	AssetNumberInfo        string `json:"assetNumberInfo"`
	AcquisitionCostInfo    string `json:"acquisitionCostInfo"`
	SalvageValueInfo       string `json:"salvageValueInfo"`
	UsefulLifeMonthsInfo   string `json:"usefulLifeMonthsInfo"`
	DepreciationMethodInfo string `json:"depreciationMethodInfo"`
	// UnitsOfProductionDisabledTooltip is the title shown on the disabled UoP
	// depreciation-method option in the asset edit drawer.
	UnitsOfProductionDisabledTooltip string `json:"unitsOfProductionDisabledTooltip"`
}

type AssetActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Revalue    string `json:"revalue"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
	// Confirm messages
	ConfirmActivate       string `json:"confirmActivate"`
	ConfirmDeactivate     string `json:"confirmDeactivate"`
	ConfirmDelete         string `json:"confirmDelete"`
	ConfirmBulkActivate   string `json:"confirmBulkActivate"`
	ConfirmBulkDeactivate string `json:"confirmBulkDeactivate"`
	ConfirmBulkDelete     string `json:"confirmBulkDelete"`
	// Error messages
	InvalidFormData     string `json:"invalidFormData"`
	IDRequired          string `json:"idRequired"`
	NoIDsProvided       string `json:"noIDsProvided"`
	InvalidStatus       string `json:"invalidStatus"`
	InvalidTargetStatus string `json:"invalidTargetStatus"`
	NoPermission        string `json:"noPermission"`
	// CannotDeleteInUse is shown (as a tooltip and server-side error) when a
	// delete is attempted on an asset that has one or more asset_transaction rows.
	CannotDeleteInUse string `json:"cannotDeleteInUse"`
}

type AssetDetailLabels struct {
	BasicInfo        AssetDetailBasicInfoLabels `json:"basicInfo"`
	Tabs             AssetDetailTabLabels       `json:"tabs"`
	EmptyStates      AssetDetailEmptyLabels     `json:"emptyStates"`
	AttachmentUpload string                     `json:"attachmentUpload"`
}

type AssetDetailBasicInfoLabels struct {
	Title              string `json:"title"`
	Name               string `json:"name"`
	AssetNumber        string `json:"assetNumber"`
	Description        string `json:"description"`
	Category           string `json:"category"`
	Location           string `json:"location"`
	AcquisitionCost    string `json:"acquisitionCost"`
	SalvageValue       string `json:"salvageValue"`
	UsefulLifeMonths   string `json:"usefulLifeMonths"`
	DepreciationMethod string `json:"depreciationMethod"`
	BookValue          string `json:"bookValue"`
	Status             string `json:"status"`
}

type AssetDetailTabLabels struct {
	Info                  string `json:"info"`
	LapsingActualSchedule string `json:"lapsingActualSchedule"`
	TransactionLedger     string `json:"transactionLedger"`
	Attachments           string `json:"attachments"`
	History               string `json:"history"`
}

type AssetDetailEmptyLabels struct {
	DepreciationTitle string `json:"depreciationTitle"`
	DepreciationDesc  string `json:"depreciationDesc"`
	MaintenanceTitle  string `json:"maintenanceTitle"`
	MaintenanceDesc   string `json:"maintenanceDesc"`
	TransactionsTitle string `json:"transactionsTitle"`
	TransactionsDesc  string `json:"transactionsDesc"`
}

type AssetDashboardLabels struct {
	Title            string `json:"title"`
	Subtitle         string `json:"subtitle"`
	TotalAssets      string `json:"totalAssets"`
	TotalBookValue   string `json:"totalBookValue"`
	FullyDepreciated string `json:"fullyDepreciated"`
	UnderMaintenance string `json:"underMaintenance"`
	// Activity feed
	ActivityAcquired     string `json:"activityAcquired"`
	ActivityMaintenance  string `json:"activityMaintenance"`
	ActivityDepreciation string `json:"activityDepreciation"`
	RecentActivity       string `json:"recentActivity"`
	NoRecentActivity     string `json:"noRecentActivity"`
	UnknownAsset         string `json:"unknownAsset"`
	// Pyeza dashboard block — quick actions + chart widget title
	AssetValueTrend           string `json:"assetValueTrend"`
	ViewAll                   string `json:"viewAll"`
	QuickNewAsset             string `json:"quickNewAsset"`
	QuickViewAll              string `json:"quickViewAll"`
	QuickDepreciationSchedule string `json:"quickDepreciationSchedule"`
	QuickMaintenanceLog       string `json:"quickMaintenanceLog"`
}

// DefaultAssetLabels returns AssetLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultAssetLabels() AssetLabels {
	return AssetLabels{
		Page: AssetPageLabels{
			Heading:         "Fixed Assets",
			HeadingActive:   "Active Assets",
			HeadingInactive: "Inactive Assets",
			Caption:         "Manage your fixed assets",
			CaptionActive:   "Active fixed assets in your register",
			CaptionInactive: "Inactive or disposed fixed assets",
		},
		Buttons: AssetButtonLabels{
			AddAsset: "Add Asset",
		},
		Columns: AssetColumnLabels{
			AssetNumber:     "Asset Number",
			Name:            "Name",
			Category:        "Category",
			Location:        "Location",
			AcquisitionCost: "Acquisition Cost",
			BookValue:       "Book Value",
			Status:          "Status",
			Period:          "Period",
			StartValue:      "Start Value",
			Depreciation:    "Depreciation",
			EndValue:        "End Value",
			Accumulated:     "Accumulated",
			Date:            "Date",
			Type:            "Type",
			Description:     "Description",
			Cost:            "Cost",
			Amount:          "Amount",
			Reference:       "Reference",
			Item:            "Item",
			COGS:            "COGS",
			NetRevenue:      "Net Revenue",
			COGSPct:         "COGS %",
			Units:           "Units",
			Totals:          "TOTALS",
		},
		Empty: AssetEmptyLabels{
			ActiveTitle:     "No active assets",
			ActiveMessage:   "Add your first fixed asset to start tracking depreciation and maintenance.",
			InactiveTitle:   "No inactive assets",
			InactiveMessage: "Deactivated or disposed assets will appear here.",
		},
		Form: AssetFormLabels{
			Name:                             "Name",
			NamePlaceholder:                  "e.g. Office Laptop",
			AssetNumber:                      "Asset Number",
			AssetNumberPlaceholder:           "e.g. FA-001",
			Description:                      "Description",
			DescriptionPlaceholder:           "Brief description of the asset",
			Category:                         "Category",
			CategoryPlaceholder:              "Select a category",
			Location:                         "Location",
			LocationPlaceholder:              "Select a location",
			AcquisitionCost:                  "Acquisition Cost",
			AcquisitionCostPlaceholder:       "0.00",
			SalvageValue:                     "Salvage Value",
			SalvageValuePlaceholder:          "0.00",
			UsefulLifeMonths:                 "Useful Life (Months)",
			UsefulLifePlaceholder:            "e.g. 60",
			DepreciationMethod:               "Depreciation Method",
			Active:                           "Active",
			DepMethodStraightLine:            "Straight Line",
			DepMethodDecliningBalance:        "Declining Balance",
			DepMethodSumOfYears:              "Sum of Years' Digits",
			DepMethodUnitsOfProduction:       "Units of Production",
			AssetNumberInfo:                  "Unique identifier for this asset in your register (e.g. FA-001).",
			AcquisitionCostInfo:              "Total cost to acquire or construct the asset, including installation.",
			SalvageValueInfo:                 "Estimated residual value at the end of the asset's useful life.",
			UsefulLifeMonthsInfo:             "Expected productive life of the asset in months, used to calculate depreciation.",
			DepreciationMethodInfo:           "The accounting method used to allocate the asset's cost over its useful life.",
			UnitsOfProductionDisabledTooltip: "Units of Production depreciation is not yet supported.",
		},
		Actions: AssetActionLabels{
			View:                  "View",
			Edit:                  "Edit",
			Revalue:               "Revalue",
			Delete:                "Delete",
			Activate:              "Activate",
			Deactivate:            "Deactivate",
			ConfirmActivate:       "Are you sure you want to activate %s?",
			ConfirmDeactivate:     "Are you sure you want to deactivate %s?",
			ConfirmDelete:         "Are you sure you want to delete %s? This action cannot be undone.",
			ConfirmBulkActivate:   "Are you sure you want to activate {{count}} asset(s)?",
			ConfirmBulkDeactivate: "Are you sure you want to deactivate {{count}} asset(s)?",
			ConfirmBulkDelete:     "Are you sure you want to delete {{count}} asset(s)? This action cannot be undone.",
			InvalidFormData:       "Invalid form data",
			IDRequired:            "Asset ID is required",
			NoIDsProvided:         "No asset IDs provided",
			InvalidStatus:         "Invalid status",
			InvalidTargetStatus:   "Invalid target status",
			NoPermission:          "No permission",
			CannotDeleteInUse:     "Cannot delete: asset has posted transactions.",
		},
		Detail: AssetDetailLabels{
			BasicInfo: AssetDetailBasicInfoLabels{
				Title:              "Asset Information",
				Name:               "Name",
				AssetNumber:        "Asset Number",
				Description:        "Description",
				Category:           "Category",
				Location:           "Location",
				AcquisitionCost:    "Acquisition Cost",
				SalvageValue:       "Salvage Value",
				UsefulLifeMonths:   "Useful Life (Months)",
				DepreciationMethod: "Depreciation Method",
				BookValue:          "Book Value",
				Status:             "Status",
			},
			Tabs: AssetDetailTabLabels{
				Info:                  "Info",
				LapsingActualSchedule: "Lapsing Schedule",
				TransactionLedger:     "Transaction Ledger",
				Attachments:           "Attachments",
				History:               "History",
			},
			EmptyStates: AssetDetailEmptyLabels{
				DepreciationTitle: "No depreciation records",
				DepreciationDesc:  "Depreciation schedule will appear here once configured.",
				MaintenanceTitle:  "No maintenance records",
				MaintenanceDesc:   "Maintenance history for this asset will appear here.",
				TransactionsTitle: "No transactions",
				TransactionsDesc:  "Transaction audit trail for this asset will appear here.",
			},
			AttachmentUpload: "Upload Attachment",
		},
		Dashboard: AssetDashboardLabels{
			Title:                     "Assets Dashboard",
			Subtitle:                  "Overview of your fixed asset register",
			TotalAssets:               "Total Assets",
			TotalBookValue:            "Total Book Value",
			FullyDepreciated:          "Fully Depreciated",
			UnderMaintenance:          "Under Maintenance",
			ActivityAcquired:          "New asset acquired",
			ActivityMaintenance:       "Maintenance completed",
			ActivityDepreciation:      "Depreciation recorded",
			RecentActivity:            "Recent Activity",
			NoRecentActivity:          "No recent asset activity",
			UnknownAsset:              "Unknown Asset",
			AssetValueTrend:           "Asset Value Trend",
			ViewAll:                   "View All",
			QuickNewAsset:             "New Asset",
			QuickViewAll:              "View All Assets",
			QuickDepreciationSchedule: "Depreciation Schedule",
			QuickMaintenanceLog:       "Maintenance Log",
		},
	}
}

// ---------------------------------------------------------------------------
// Account labels (Chart of Accounts)
// ---------------------------------------------------------------------------

// AccountLabels holds all translatable strings for the Chart of Accounts module.
type AccountLabels struct {
	Page          AccountPageLabels          `json:"page"`
	Buttons       AccountButtonLabels        `json:"buttons"`
	Columns       AccountColumnLabels        `json:"columns"`
	Tabs          AccountTabLabels           `json:"tabs"`
	Empty         AccountEmptyLabels         `json:"empty"`
	Form          AccountFormLabels          `json:"form"`
	Actions       AccountActionLabels        `json:"actions"`
	Detail        AccountDetailLabels        `json:"detail"`
	Templates     AccountTemplatesLabels     `json:"templates"`
	GeneralLedger AccountGeneralLedgerLabels `json:"generalLedger"`
	BadDebt       BadDebtLabels              `json:"badDebt"`
	Dashboard     LedgerDashboardLabels      `json:"dashboard"`
}

// LedgerDashboardLabels holds translatable strings for the Ledger live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type LedgerDashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	TotalAssets      string `json:"totalAssets"`
	TotalLiabilities string `json:"totalLiabilities"`
	TotalEquity      string `json:"totalEquity"`
	NetIncomeMTD     string `json:"netIncomeMtd"`
	// Widgets
	BalanceByType    string `json:"balanceByType"`
	UnpostedJournals string `json:"unpostedJournals"`
	RecentJournals   string `json:"recentJournals"`
	NoRecentJournals string `json:"noRecentJournals"`
	// Quick actions
	QuickNewJournal    string `json:"quickNewJournal"`
	QuickTrialBalance  string `json:"quickTrialBalance"`
	QuickClosePeriod   string `json:"quickClosePeriod"`
	QuickAccountLookup string `json:"quickAccountLookup"`
	// Account type labels — used as bar-chart axis labels
	AccountTypeAssets      string `json:"accountTypeAssets"`
	AccountTypeLiabilities string `json:"accountTypeLiabilities"`
	AccountTypeEquity      string `json:"accountTypeEquity"`
	AccountTypeRevenue     string `json:"accountTypeRevenue"`
	AccountTypeExpense     string `json:"accountTypeExpense"`
	// Fallback label for journal entries with no entry number
	JournalEntryFallback string `json:"journalEntryFallback"`
	// Common
	ViewAll    string `json:"viewAll"`
	AxisAmount string `json:"axisAmount"`
}

// BadDebtLabels holds translatable strings for the Bad Debt Policy settings page.
type BadDebtLabels struct {
	Title string `json:"title"`
}

// AccountTemplatesLabels holds translatable strings for the Account Templates settings page.
type AccountTemplatesLabels struct {
	PageTitle           string `json:"pageTitle"`
	PageSubtitle        string `json:"pageSubtitle"`
	CurrentAccountCount string `json:"currentAccountCount"`
	ApplyWarning        string `json:"applyWarning"`
	Empty               string `json:"empty"`
	EmptyDesc           string `json:"emptyDesc"`
	AccountsSuffix      string `json:"accountsSuffix"`
	ComingSoon          string `json:"comingSoon"`
	BadgeApplied        string `json:"badgeApplied"`
	BadgeAssets         string `json:"badgeAssets"`
	BadgeLiabilities    string `json:"badgeLiabilities"`
	BadgeEquity         string `json:"badgeEquity"`
	BadgeRevenue        string `json:"badgeRevenue"`
	BadgeExpenses       string `json:"badgeExpenses"`
	Preview             string `json:"preview"`
	AlreadyApplied      string `json:"alreadyApplied"`
	ApplyTemplate       string `json:"applyTemplate"`
	PreviewTitle        string `json:"previewTitle"`
	PreviewDesc         string `json:"previewDesc"`
	ColCode             string `json:"colCode"`
	ColAccountName      string `json:"colAccountName"`
	ColElement          string `json:"colElement"`
	ColClass            string `json:"colClass"`
	ColIsGroup          string `json:"colIsGroup"`
	Yes                 string `json:"yes"`
	SkipNote            string `json:"skipNote"`
}

type AccountGeneralLedgerLabels struct {
	Title                 string `json:"title"`
	Subtitle              string `json:"subtitle"`
	Account               string `json:"account"`
	AccountPlaceholder    string `json:"accountPlaceholder"`
	StartDate             string `json:"startDate"`
	EndDate               string `json:"endDate"`
	Apply                 string `json:"apply"`
	Clear                 string `json:"clear"`
	Print                 string `json:"print"`
	SelectAccountMessage  string `json:"selectAccountMessage"`
	NoTransactionsMessage string `json:"noTransactionsMessage"`
	DateRangeSeparator    string `json:"dateRangeSeparator"`
	OpeningBalance        string `json:"openingBalance"`
	PeriodDebits          string `json:"periodDebits"`
	PeriodCredits         string `json:"periodCredits"`
	RunningBalance        string `json:"runningBalance"`
	NoTransactionsTitle   string `json:"noTransactionsTitle"`
	NoTransactionsDetail  string `json:"noTransactionsDetail"`
}

type AccountPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type AccountButtonLabels struct {
	AddAccount string `json:"addAccount"`
}

type AccountColumnLabels struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	Element        string `json:"element"`
	Classification string `json:"classification"`
	Group          string `json:"group"`
	Balance        string `json:"balance"`
	// Entry sub-table columns
	Date        string `json:"date"`
	EntryNumber string `json:"entryNumber"`
	Description string `json:"description"`
	Debit       string `json:"debit"`
	Credit      string `json:"credit"`
	// Status
	Status string `json:"status"`
}

type AccountTabLabels struct {
	All       string `json:"all"`
	Asset     string `json:"asset"`
	Liability string `json:"liability"`
	Equity    string `json:"equity"`
	Revenue   string `json:"revenue"`
	Expense   string `json:"expense"`
}

type AccountEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type AccountFormLabels struct {
	Code                     string `json:"code"`
	CodePlaceholder          string `json:"codePlaceholder"`
	Name                     string `json:"name"`
	NamePlaceholder          string `json:"namePlaceholder"`
	Element                  string `json:"element"`
	Classification           string `json:"classification"`
	ParentAccount            string `json:"parentAccount"`
	ParentAccountPlaceholder string `json:"parentAccountPlaceholder"`
	Group                    string `json:"group"`
	IsGroup                  string `json:"isGroup"`
	Active                   string `json:"active"`
	Description              string `json:"description"`
	DescriptionPlaceholder   string `json:"descriptionPlaceholder"`
	CashFlowSection          string `json:"cashFlowSection"`
	CashFlowClassification   string `json:"cashFlowClassification"`
	// Element option labels
	ElementAsset     string `json:"elementAsset"`
	ElementLiability string `json:"elementLiability"`
	ElementEquity    string `json:"elementEquity"`
	ElementRevenue   string `json:"elementRevenue"`
	ElementExpense   string `json:"elementExpense"`
	// Class option labels
	ClassCurrentAsset        string `json:"classCurrentAsset"`
	ClassNonCurrentAsset     string `json:"classNonCurrentAsset"`
	ClassCurrentLiability    string `json:"classCurrentLiability"`
	ClassNonCurrentLiability string `json:"classNonCurrentLiability"`
	ClassEquity              string `json:"classEquity"`
	ClassOperatingRevenue    string `json:"classOperatingRevenue"`
	ClassOtherIncome         string `json:"classOtherIncome"`
	ClassCostOfSales         string `json:"classCostOfSales"`
	ClassOperatingExpense    string `json:"classOperatingExpense"`
	ClassFinanceCost         string `json:"classFinanceCost"`
	ClassIncomeTax           string `json:"classIncomeTax"`
	ClassOtherExpense        string `json:"classOtherExpense"`
	// Cash flow option labels
	CashFlowNone      string `json:"cashFlowNone"`
	CashFlowOperating string `json:"cashFlowOperating"`
	CashFlowInvesting string `json:"cashFlowInvesting"`
	CashFlowFinancing string `json:"cashFlowFinancing"`
	// Element group header labels (used in BuildAccountTree and preview)
	GroupAssets      string `json:"groupAssets"`
	GroupLiabilities string `json:"groupLiabilities"`
	GroupEquity      string `json:"groupEquity"`
	GroupRevenue     string `json:"groupRevenue"`
	GroupExpenses    string `json:"groupExpenses"`
	// Normal balance value labels
	NormalBalanceDebit  string `json:"normalBalanceDebit"`
	NormalBalanceCredit string `json:"normalBalanceCredit"`
	// Info popover text
	ElementInfo        string `json:"elementInfo"`
	ClassificationInfo string `json:"classificationInfo"`
	CashFlowClassInfo  string `json:"cashFlowClassInfo"`
}

type AccountActionLabels struct {
	View   string `json:"view"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
	// Confirm messages
	ConfirmDelete     string `json:"confirmDelete"`
	ConfirmBulkDelete string `json:"confirmBulkDelete"`
	// Error messages
	NoPermission string `json:"noPermission"`
}

type AccountDetailLabels struct {
	Tabs        AccountDetailTabLabels   `json:"tabs"`
	EmptyStates AccountDetailEmptyLabels `json:"emptyStates"`
	Stats       AccountDetailStatLabels  `json:"stats"`
	Info        AccountDetailInfoLabels  `json:"info"`
}

type AccountDetailTabLabels struct {
	JournalEntries string `json:"journalEntries"`
	Details        string `json:"details"`
}

type AccountDetailEmptyLabels struct {
	EntriesTitle   string `json:"entriesTitle"`
	EntriesMessage string `json:"entriesMessage"`
}

type AccountDetailStatLabels struct {
	CurrentBalance string `json:"currentBalance"`
	PeriodDebits   string `json:"periodDebits"`
	PeriodCredits  string `json:"periodCredits"`
}

type AccountDetailInfoLabels struct {
	Title          string `json:"title"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Element        string `json:"element"`
	Classification string `json:"classification"`
	Group          string `json:"group"`
	ParentAccount  string `json:"parentAccount"`
	NormalBalance  string `json:"normalBalance"`
	CashFlowTag    string `json:"cashFlowTag"`
	Description    string `json:"description"`
	Status         string `json:"status"`
	Created        string `json:"created"`
	LastModified   string `json:"lastModified"`
}

// DefaultAccountLabels returns AccountLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultAccountLabels() AccountLabels {
	return AccountLabels{
		Page: AccountPageLabels{
			Heading: "Chart of Accounts",
			Caption: "Manage your organization's account structure",
		},
		Buttons: AccountButtonLabels{
			AddAccount: "Add Account",
		},
		Columns: AccountColumnLabels{
			Code:           "Code",
			Name:           "Account Name",
			Element:        "Element",
			Classification: "Class",
			Group:          "Group",
			Balance:        "Balance",
			Date:           "Date",
			EntryNumber:    "Entry #",
			Description:    "Description",
			Debit:          "Debit",
			Credit:         "Credit",
			Status:         "Status",
		},
		Tabs: AccountTabLabels{
			All:       "All",
			Asset:     "Assets",
			Liability: "Liabilities",
			Equity:    "Equity",
			Revenue:   "Revenue",
			Expense:   "Expenses",
		},
		Empty: AccountEmptyLabels{
			Title:   "No accounts found",
			Message: "Add your first account or apply a template to get started.",
		},
		Form: AccountFormLabels{
			Code:                     "Account Code",
			CodePlaceholder:          "e.g. 1140",
			Name:                     "Account Name",
			NamePlaceholder:          "e.g. Petty Cash",
			Element:                  "Element",
			Classification:           "Class",
			ParentAccount:            "Parent Account",
			ParentAccountPlaceholder: "Search groups\u2026",
			Group:                    "Group",
			IsGroup:                  "Is Group",
			Active:                   "Active",
			Description:              "Description",
			DescriptionPlaceholder:   "Brief description of this account",
			CashFlowSection:          "Cash Flow Tag (optional)",
			CashFlowClassification:   "Cash Flow Classification",
			ElementAsset:             "Asset",
			ElementLiability:         "Liability",
			ElementEquity:            "Equity",
			ElementRevenue:           "Revenue",
			ElementExpense:           "Expense",
			ClassCurrentAsset:        "Current Asset",
			ClassNonCurrentAsset:     "Non-Current Asset",
			ClassCurrentLiability:    "Current Liability",
			ClassNonCurrentLiability: "Non-Current Liability",
			ClassEquity:              "Equity",
			ClassOperatingRevenue:    "Operating Revenue",
			ClassOtherIncome:         "Other Income",
			ClassCostOfSales:         "Cost of Sales",
			ClassOperatingExpense:    "Operating Expense",
			ClassFinanceCost:         "Finance Cost",
			ClassIncomeTax:           "Income Tax",
			ClassOtherExpense:        "Other Expense",
			CashFlowNone:             "None",
			CashFlowOperating:        "Operating Activities",
			CashFlowInvesting:        "Investing Activities",
			CashFlowFinancing:        "Financing Activities",
			GroupAssets:              "ASSETS",
			GroupLiabilities:         "LIABILITIES",
			GroupEquity:              "EQUITY",
			GroupRevenue:             "REVENUE",
			GroupExpenses:            "EXPENSES",
			NormalBalanceDebit:       "Debit",
			NormalBalanceCredit:      "Credit",
			ElementInfo:              "The broad accounting category this account belongs to (Asset, Liability, Equity, Revenue, or Expense).",
			ClassificationInfo:       "Sub-category within the element that determines where this account appears on financial statements.",
			CashFlowClassInfo:        "Tags this account's movements for cash flow statement classification. Leave empty if not applicable.",
		},
		Actions: AccountActionLabels{
			View:              "View",
			Edit:              "Edit",
			Delete:            "Delete",
			ConfirmDelete:     "Are you sure you want to delete %s? This action cannot be undone.",
			ConfirmBulkDelete: "Are you sure you want to delete {{count}} account(s)? This action cannot be undone.",
			NoPermission:      "No permission",
		},
		Detail: AccountDetailLabels{
			Tabs: AccountDetailTabLabels{
				JournalEntries: "Journal Entries",
				Details:        "Details",
			},
			EmptyStates: AccountDetailEmptyLabels{
				EntriesTitle:   "No journal entries",
				EntriesMessage: "Journal entries that touch this account will appear here.",
			},
			Stats: AccountDetailStatLabels{
				CurrentBalance: "Current Balance",
				PeriodDebits:   "Period Debits",
				PeriodCredits:  "Period Credits",
			},
			Info: AccountDetailInfoLabels{
				Title:          "Account Information",
				Code:           "Account Code",
				Name:           "Account Name",
				Element:        "Element",
				Classification: "Classification",
				Group:          "Group",
				ParentAccount:  "Parent Account",
				NormalBalance:  "Normal Balance",
				CashFlowTag:    "Cash Flow Tag",
				Description:    "Description",
				Status:         "Status",
				Created:        "Created",
				LastModified:   "Last Modified",
			},
		},
		Templates: AccountTemplatesLabels{
			PageTitle:           "Account Templates",
			PageSubtitle:        "Pre-built Chart of Accounts for your business type",
			CurrentAccountCount: "Your Chart of Accounts currently has {{.CurrentAccountCount}} accounts.",
			ApplyWarning:        "Applying a template will add new accounts. Existing accounts with matching codes will be skipped.",
			Empty:               "Your Chart of Accounts is empty.",
			EmptyDesc:           "Apply a template below to get started with a standard set of accounts for your business type.",
			AccountsSuffix:      "accounts \u00b7 PFRS-compliant",
			ComingSoon:          "Coming soon",
			BadgeApplied:        "Applied",
			BadgeAssets:         "Assets",
			BadgeLiabilities:    "Liabilities",
			BadgeEquity:         "Equity",
			BadgeRevenue:        "Revenue",
			BadgeExpenses:       "Expenses",
			Preview:             "Preview",
			AlreadyApplied:      "Already applied",
			ApplyTemplate:       "Apply Template",
			PreviewTitle:        "Preview",
			PreviewDesc:         "This template will create {{.AccountCount}} accounts organized as follows:",
			ColCode:             "Code",
			ColAccountName:      "Account Name",
			ColElement:          "Element",
			ColClass:            "Class",
			ColIsGroup:          "Is Group",
			Yes:                 "Yes",
			SkipNote:            "Accounts with matching codes in your existing Chart of Accounts will be skipped.",
		},
		GeneralLedger: AccountGeneralLedgerLabels{
			Title:                 "General Ledger",
			Subtitle:              "Detailed transaction history by account",
			Account:               "Account",
			AccountPlaceholder:    "Select an account",
			StartDate:             "Start Date",
			EndDate:               "End Date",
			Apply:                 "Apply",
			Clear:                 "Clear",
			Print:                 "Print",
			SelectAccountMessage:  "Select an account above to view its detailed transaction history.",
			NoTransactionsMessage: "No transactions found for the selected account and date range.",
			DateRangeSeparator:    "to",
			OpeningBalance:        "Opening Balance",
			PeriodDebits:          "Period Debits",
			PeriodCredits:         "Period Credits",
			RunningBalance:        "Running Balance",
			NoTransactionsTitle:   "No transactions",
			NoTransactionsDetail:  "No journal entries found for this account in the selected date range.",
		},
		Dashboard: LedgerDashboardLabels{
			Title:                  "Ledger Dashboard",
			Subtitle:               "Live position of your books — assets, liabilities, equity, and net income",
			TotalAssets:            "Total Assets",
			TotalLiabilities:       "Total Liabilities",
			TotalEquity:            "Total Equity",
			NetIncomeMTD:           "Net Income (MTD)",
			BalanceByType:          "Account Balance by Type",
			UnpostedJournals:       "Unposted Journals",
			RecentJournals:         "Recent Journals",
			NoRecentJournals:       "No recent journal entries",
			QuickNewJournal:        "New Journal Entry",
			QuickTrialBalance:      "Trial Balance",
			QuickClosePeriod:       "Close Period",
			QuickAccountLookup:     "Account Lookup",
			AccountTypeAssets:      "Assets",
			AccountTypeLiabilities: "Liabilities",
			AccountTypeEquity:      "Equity",
			AccountTypeRevenue:     "Revenue",
			AccountTypeExpense:     "Expense",
			JournalEntryFallback:   "Journal",
			ViewAll:                "View All",
			AxisAmount:             "Amount",
		},
	}
}

// ---------------------------------------------------------------------------
// Journal labels (Journal Entries)
// ---------------------------------------------------------------------------

// JournalLabels holds all translatable strings for the Journal Entries module.
type JournalLabels struct {
	Page    JournalPageLabels       `json:"page"`
	Tabs    JournalTabLabels        `json:"tabs"`
	Buttons JournalButtonLabels     `json:"buttons"`
	Columns JournalColumnLabels     `json:"columns"`
	Empty   JournalEmptyLabels      `json:"empty"`
	Actions JournalActionLabels     `json:"actions"`
	Lines   JournalLineLabels       `json:"lines"`
	Form    JournalFormLabels       `json:"form"`
	Detail  JournalDetailLabels     `json:"detail"`
	Confirm JournalConfirmLabels    `json:"confirm"`
	Source  JournalSourceTypeLabels `json:"source"`
}

// JournalSourceTypeLabels maps JournalSourceType enum values to display strings.
// Keys mirror the proto JournalSourceType enum; used by the journal list and
// detail views to label the originating-transaction source.
type JournalSourceTypeLabels struct {
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

// JournalConfirmLabels holds confirmation dialog strings for journal actions.
type JournalConfirmLabels struct {
	Post    string `json:"post"`
	Delete  string `json:"delete"`
	Reverse string `json:"reverse"`
}

// JournalDetailLabels holds translatable strings for the journal detail page.
type JournalDetailLabels struct {
	Stats       JournalDetailStatLabels `json:"stats"`
	Info        JournalDetailInfoLabels `json:"info"`
	SourceLabel string                  `json:"sourceLabel"`
	ViewSource  string                  `json:"viewSource"`
	// Balance status badges shown in totals row
	Balanced   string `json:"balanced"`
	Unbalanced string `json:"unbalanced"`
	Totals     string `json:"totals"`
	Difference string `json:"difference"`
	// Tab labels
	TabLines       string `json:"tabLines"`
	TabAttachments string `json:"tabAttachments"`
}

type JournalDetailStatLabels struct {
	TotalDebit  string `json:"totalDebit"`
	TotalCredit string `json:"totalCredit"`
}

type JournalDetailInfoLabels struct {
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

type JournalPageLabels struct {
	HeadingDraft     string `json:"headingDraft"`
	SubtitleDraft    string `json:"subtitleDraft"`
	HeadingPosted    string `json:"headingPosted"`
	SubtitlePosted   string `json:"subtitlePosted"`
	HeadingReversed  string `json:"headingReversed"`
	SubtitleReversed string `json:"subtitleReversed"`
}

type JournalTabLabels struct {
	Draft    string `json:"draft"`
	Posted   string `json:"posted"`
	Reversed string `json:"reversed"`
}

type JournalButtonLabels struct {
	NewEntry string `json:"newEntry"`
}

type JournalColumnLabels struct {
	EntryNumber string `json:"entryNumber"`
	Date        string `json:"date"`
	Description string `json:"description"`
	Amount      string `json:"amount"`
	Source      string `json:"source"`
	Status      string `json:"status"`
}

type JournalEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type JournalActionLabels struct {
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

// JournalFormLabels holds translatable strings for the journal entry form.
type JournalFormLabels struct {
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

type JournalLineLabels struct {
	AccountCode  string `json:"accountCode"`
	AccountName  string `json:"accountName"`
	Memo         string `json:"memo"`
	Debit        string `json:"debit"`
	Credit       string `json:"credit"`
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// DefaultJournalLabels returns JournalLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultJournalLabels() JournalLabels {
	return JournalLabels{
		Page: JournalPageLabels{
			HeadingDraft:     "Draft Journal Entries",
			SubtitleDraft:    "Review and post journal entries that are still in draft",
			HeadingPosted:    "Posted Journal Entries",
			SubtitlePosted:   "View journal entries that have been posted to the ledger",
			HeadingReversed:  "Reversed Journal Entries",
			SubtitleReversed: "View journal entries that have been reversed",
		},
		Tabs: JournalTabLabels{
			Draft:    "Draft",
			Posted:   "Posted",
			Reversed: "Reversed",
		},
		Buttons: JournalButtonLabels{
			NewEntry: "New Entry",
		},
		Columns: JournalColumnLabels{
			EntryNumber: "Entry #",
			Date:        "Date",
			Description: "Description",
			Amount:      "Amount",
			Source:      "Source",
			Status:      "Status",
		},
		Empty: JournalEmptyLabels{
			Title:   "No journal entries",
			Message: "No journal entries found for this status.",
		},
		Actions: JournalActionLabels{
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
		Lines: JournalLineLabels{
			AccountCode:  "Account Code",
			AccountName:  "Account Name",
			Memo:         "Memo",
			Debit:        "Debit",
			Credit:       "Credit",
			EmptyTitle:   "No journal lines",
			EmptyMessage: "This journal entry has no lines.",
		},
		Form: JournalFormLabels{
			Date:                   "Date",
			DatePlaceholder:        "YYYY-MM-DD",
			Description:            "Description",
			DescriptionPlaceholder: "e.g. Office supplies purchase",
			Notes:                  "Notes",
			NotesPlaceholder:       "Optional notes or reference",
			LineNumber:             "#",
			Account:                "Account",
			AccountPlaceholder:     "Search by code or name\u2026",
			Debit:                  "Debit",
			Credit:                 "Credit",
			Memo:                   "Memo",
			AddLine:                "+ Add Line",
			RemoveLine:             "Remove",
			BalancedMessage:        "Balanced \u2014 Total Debits equal Total Credits",
			UnbalancedMessage:      "Unbalanced \u2014 Debits and Credits do not match",
			DifferenceLabel:        "Difference",
			SaveDraft:              "Save as Draft",
			PostEntry:              "Post",
			EntryDetails:           "Entry Details",
			JournalLines:           "Journal Lines",
			BalanceHint:            "Enter debits and credits above",
		},
		Detail: JournalDetailLabels{
			Stats: JournalDetailStatLabels{
				TotalDebit:  "Total Debits",
				TotalCredit: "Total Credits",
			},
			Info: JournalDetailInfoLabels{
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
			ViewSource:     "View Source \u2192",
			Balanced:       "Balanced",
			Unbalanced:     "Unbalanced",
			Totals:         "TOTALS",
			Difference:     "DIFFERENCE",
			TabLines:       "Journal Lines",
			TabAttachments: "Attachments",
		},
		Confirm: JournalConfirmLabels{
			Post:    "Are you sure you want to post this journal entry? This action cannot be undone.",
			Delete:  "Are you sure you want to delete this journal entry? This action cannot be undone.",
			Reverse: "Are you sure you want to reverse this journal entry? A reversing entry will be created.",
		},
		Source: JournalSourceTypeLabels{
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

// ---------------------------------------------------------------------------
// FiscalPeriod labels
// ---------------------------------------------------------------------------

// FiscalPeriodLabels holds all translatable strings for the fiscal period module.
type FiscalPeriodLabels struct {
	Page    FiscalPeriodPageLabels   `json:"page"`
	Buttons FiscalPeriodButtonLabels `json:"buttons"`
	Columns FiscalPeriodColumnLabels `json:"columns"`
	Status  FiscalPeriodStatusLabels `json:"status"`
	Empty   FiscalPeriodEmptyLabels  `json:"empty"`
	Actions FiscalPeriodActionLabels `json:"actions"`
	Form    FiscalPeriodFormLabels   `json:"form"`
}

// FiscalPeriodFormLabels holds field-level labels for the fiscal period add/edit form.
type FiscalPeriodFormLabels struct {
	Name         string `json:"name"`
	PeriodNumber string `json:"period_number"`
	FiscalYear   string `json:"fiscal_year"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Status       string `json:"status"`
	// Info popover text
	PeriodNumberInfo string `json:"periodNumberInfo"`
	FiscalYearInfo   string `json:"fiscalYearInfo"`
	StartDateInfo    string `json:"startDateInfo"`
	EndDateInfo      string `json:"endDateInfo"`
}

type FiscalPeriodPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type FiscalPeriodButtonLabels struct {
	AddPeriod    string `json:"addPeriod"`
	CloseYearEnd string `json:"closeYearEnd"`
}

type FiscalPeriodColumnLabels struct {
	Period    string `json:"period"`
	Year      string `json:"year"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Status    string `json:"status"`
}

type FiscalPeriodStatusLabels struct {
	Open   string `json:"open"`
	Closed string `json:"closed"`
	Locked string `json:"locked"`
}

type FiscalPeriodEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FiscalPeriodActionLabels struct {
	Close        string `json:"close"`
	NoPermission string `json:"noPermission"`
	// Confirm messages
	ConfirmClose string `json:"confirmClose"`
}

// DefaultFiscalPeriodLabels returns FiscalPeriodLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultFiscalPeriodLabels() FiscalPeriodLabels {
	return FiscalPeriodLabels{
		Page: FiscalPeriodPageLabels{
			Heading: "Fiscal Periods",
			Caption: "Manage accounting periods and year-end close",
		},
		Buttons: FiscalPeriodButtonLabels{
			AddPeriod:    "Add Period",
			CloseYearEnd: "Close Year-End",
		},
		Columns: FiscalPeriodColumnLabels{
			Period:    "Period",
			Year:      "Year",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Status:    "Status",
		},
		Status: FiscalPeriodStatusLabels{
			Open:   "Open",
			Closed: "Closed",
			Locked: "Locked",
		},
		Empty: FiscalPeriodEmptyLabels{
			Title:   "No fiscal periods found",
			Message: "Add your first fiscal period to start tracking accounting periods.",
		},
		Actions: FiscalPeriodActionLabels{
			Close:        "Close",
			NoPermission: "No permission",
			ConfirmClose: "Are you sure you want to close %s? This will prevent new journal entries from being posted to this period.",
		},
		Form: FiscalPeriodFormLabels{
			Name:             "Name",
			PeriodNumber:     "Period Number",
			FiscalYear:       "Fiscal Year",
			StartDate:        "Start Date",
			EndDate:          "End Date",
			Status:           "Status",
			PeriodNumberInfo: "Sequential number of this period within the fiscal year (e.g. 1 for January).",
			FiscalYearInfo:   "The fiscal year this period belongs to (e.g. 2026).",
			StartDateInfo:    "First day of this accounting period. Journal entries dated on or after this date fall within the period.",
			EndDateInfo:      "Last day of this accounting period. Journal entries dated on or before this date fall within the period.",
		},
	}
}

// ---------------------------------------------------------------------------
// RecurringTemplate labels
// ---------------------------------------------------------------------------

// RecurringTemplateLabels holds all translatable strings for the recurring journal template module.
type RecurringTemplateLabels struct {
	Page      RecurringTemplatePageLabels      `json:"page"`
	Buttons   RecurringTemplateButtonLabels    `json:"buttons"`
	Columns   RecurringTemplateColumnLabels    `json:"columns"`
	Frequency RecurringTemplateFrequencyLabels `json:"frequency"`
	Status    RecurringTemplateStatusLabels    `json:"status"`
	Empty     RecurringTemplateEmptyLabels     `json:"empty"`
	Actions   RecurringTemplateActionLabels    `json:"actions"`
}

type RecurringTemplatePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type RecurringTemplateButtonLabels struct {
	AddTemplate string `json:"addTemplate"`
}

type RecurringTemplateColumnLabels struct {
	Name      string `json:"name"`
	Frequency string `json:"frequency"`
	NextRun   string `json:"nextRun"`
	Status    string `json:"status"`
}

type RecurringTemplateFrequencyLabels struct {
	Daily     string `json:"daily"`
	Weekly    string `json:"weekly"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Yearly    string `json:"yearly"`
}

type RecurringTemplateStatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type RecurringTemplateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RecurringTemplateActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Pause        string `json:"pause"`
	Resume       string `json:"resume"`
	NoPermission string `json:"noPermission"`
}

// DefaultRecurringTemplateLabels returns RecurringTemplateLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultRecurringTemplateLabels() RecurringTemplateLabels {
	return RecurringTemplateLabels{
		Page: RecurringTemplatePageLabels{
			Heading: "Recurring Journal Entries",
			Caption: "Automated entries for depreciation, amortization, accruals",
		},
		Buttons: RecurringTemplateButtonLabels{
			AddTemplate: "Add Recurring Entry",
		},
		Columns: RecurringTemplateColumnLabels{
			Name:      "Name",
			Frequency: "Frequency",
			NextRun:   "Next Run",
			Status:    "Status",
		},
		Frequency: RecurringTemplateFrequencyLabels{
			Daily:     "Daily",
			Weekly:    "Weekly",
			Monthly:   "Monthly",
			Quarterly: "Quarterly",
			Yearly:    "Yearly",
		},
		Status: RecurringTemplateStatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		Empty: RecurringTemplateEmptyLabels{
			Title:   "No recurring templates",
			Message: "Add your first recurring journal entry template to automate periodic entries.",
		},
		Actions: RecurringTemplateActionLabels{
			Edit:         "Edit",
			Delete:       "Delete",
			Pause:        "Pause",
			Resume:       "Resume",
			NoPermission: "No permission",
		},
	}
}

// ---------------------------------------------------------------------------
// Payroll labels
// ---------------------------------------------------------------------------

// PayrollLabels holds all translatable strings for the Payroll module.
type PayrollLabels struct {
	Run        PayrollRunLabels        `json:"run"`
	Remittance PayrollRemittanceLabels `json:"remittance"`
	Employee   PayrollEmployeeLabels   `json:"employee"`
	Settings   PayrollSettingsLabels   `json:"settings"`
	Dashboard  PayrollDashboardLabels  `json:"dashboard"`
}

// PayrollDashboardLabels holds translatable strings for the Payroll live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type PayrollDashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	CurrentRunStatus   string `json:"currentRunStatus"`
	EmployeesInCurrent string `json:"employeesInCurrent"`
	TotalGrossMTD      string `json:"totalGrossMtd"`
	RemittancesDue     string `json:"remittancesDue"`
	// Widgets
	GrossPayByMonth     string `json:"grossPayByMonth"`
	RecentRuns          string `json:"recentRuns"`
	UpcomingRemittances string `json:"upcomingRemittances"`
	NoRecentRuns        string `json:"noRecentRuns"`
	NoUpcomingDeadlines string `json:"noUpcomingDeadlines"`
	// Quick actions
	QuickNewRun            string `json:"quickNewRun"`
	QuickProcessRun        string `json:"quickProcessRun"`
	QuickFileRemittance    string `json:"quickFileRemittance"`
	QuickPayPeriodSettings string `json:"quickPayPeriodSettings"`
	// Common
	ViewAll   string `json:"viewAll"`
	AxisGross string `json:"axisGross"`
	NoRunYet  string `json:"noRunYet"`
}

// PayrollRunLabels holds labels for the Payroll Run sub-module.
type PayrollRunLabels struct {
	Page    PayrollRunPageLabels   `json:"page"`
	Tabs    PayrollRunTabLabels    `json:"tabs"`
	Buttons PayrollRunButtonLabels `json:"buttons"`
	Columns PayrollRunColumnLabels `json:"columns"`
	Empty   PayrollRunEmptyLabels  `json:"empty"`
	Actions PayrollRunActionLabels `json:"actions"`
}

type PayrollRunPageLabels struct {
	HeadingDraft       string `json:"headingDraft"`
	SubtitleDraft      string `json:"subtitleDraft"`
	HeadingCalculated  string `json:"headingCalculated"`
	SubtitleCalculated string `json:"subtitleCalculated"`
	HeadingApproved    string `json:"headingApproved"`
	SubtitleApproved   string `json:"subtitleApproved"`
	HeadingPosted      string `json:"headingPosted"`
	SubtitlePosted     string `json:"subtitlePosted"`
}

type PayrollRunTabLabels struct {
	Draft      string `json:"draft"`
	Calculated string `json:"calculated"`
	Approved   string `json:"approved"`
	Posted     string `json:"posted"`
}

type PayrollRunButtonLabels struct {
	NewRun string `json:"newRun"`
}

type PayrollRunColumnLabels struct {
	RunNumber       string `json:"runNumber"`
	PayPeriod       string `json:"payPeriod"`
	Employees       string `json:"employees"`
	TotalGross      string `json:"totalGross"`
	TotalDeductions string `json:"totalDeductions"`
	TotalNet        string `json:"totalNet"`
	Status          string `json:"status"`
	ApprovedBy      string `json:"approvedBy"`
	PostedAt        string `json:"postedAt"`
}

type PayrollRunEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PayrollRunActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// PayrollRemittanceLabels holds labels for the Payroll Remittance sub-module.
type PayrollRemittanceLabels struct {
	Page    PayrollRemittancePageLabels   `json:"page"`
	Tabs    PayrollRemittanceTabLabels    `json:"tabs"`
	Columns PayrollRemittanceColumnLabels `json:"columns"`
	Types   PayrollRemittanceTypeLabels   `json:"types"`
	Empty   PayrollRemittanceEmptyLabels  `json:"empty"`
}

type PayrollRemittancePageLabels struct {
	HeadingPending  string `json:"headingPending"`
	SubtitlePending string `json:"subtitlePending"`
	HeadingFiled    string `json:"headingFiled"`
	SubtitleFiled   string `json:"subtitleFiled"`
	HeadingPaid     string `json:"headingPaid"`
	SubtitlePaid    string `json:"subtitlePaid"`
}

type PayrollRemittanceTabLabels struct {
	Pending string `json:"pending"`
	Filed   string `json:"filed"`
	Paid    string `json:"paid"`
}

type PayrollRemittanceColumnLabels struct {
	RemittanceType  string `json:"remittanceType"`
	Amount          string `json:"amount"`
	DueDate         string `json:"dueDate"`
	Status          string `json:"status"`
	FiledAt         string `json:"filedAt"`
	ReferenceNumber string `json:"referenceNumber"`
}

type PayrollRemittanceTypeLabels struct {
	SSS            string `json:"sss"`
	PhilHealth     string `json:"philHealth"`
	PagIBIG        string `json:"pagIbig"`
	BIRWithholding string `json:"birWithholding"`
}

type PayrollRemittanceEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// PayrollEmployeeLabels holds labels for the Payroll Employee sub-module.
type PayrollEmployeeLabels struct {
	Page         PayrollEmployeePageLabels         `json:"page"`
	Columns      PayrollEmployeeColumnLabels       `json:"columns"`
	Status       PayrollEmployeeStatusLabels       `json:"status"`
	PayFrequency PayrollEmployeePayFrequencyLabels `json:"payFrequency"`
	Empty        PayrollEmployeeEmptyLabels        `json:"empty"`
}

type PayrollEmployeePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type PayrollEmployeeColumnLabels struct {
	Name         string `json:"name"`
	Position     string `json:"position"`
	Department   string `json:"department"`
	BasicSalary  string `json:"basicSalary"`
	PayFrequency string `json:"payFrequency"`
	Status       string `json:"status"`
}

type PayrollEmployeeStatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type PayrollEmployeePayFrequencyLabels struct {
	SemiMonthly string `json:"semiMonthly"`
	Monthly     string `json:"monthly"`
	Weekly      string `json:"weekly"`
}

type PayrollEmployeeEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// PayrollSettingsLabels holds labels for Payroll Settings pages.
type PayrollSettingsLabels struct {
	GovRates   PayrollGovRatesLabels   `json:"govRates"`
	PayPeriods PayrollPayPeriodsLabels `json:"payPeriods"`
}

type PayrollGovRatesLabels struct {
	Page   PayrollGovRatesPageLabels   `json:"page"`
	Agency PayrollGovRatesAgencyLabels `json:"agency"`
}

type PayrollGovRatesPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type PayrollGovRatesAgencyLabels struct {
	SSS            string `json:"sss"`
	PhilHealth     string `json:"philHealth"`
	PagIBIG        string `json:"pagIbig"`
	BIRWithholding string `json:"birWithholding"`
}

type PayrollPayPeriodsLabels struct {
	Page PayrollPayPeriodsPageLabels `json:"page"`
}

type PayrollPayPeriodsPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

// DefaultPayrollLabels returns PayrollLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultPayrollLabels() PayrollLabels {
	return PayrollLabels{
		Run: PayrollRunLabels{
			Page: PayrollRunPageLabels{
				HeadingDraft:       "Draft Payroll Runs",
				SubtitleDraft:      "Payroll runs in preparation \u2014 payslips not yet finalized",
				HeadingCalculated:  "Calculated Payroll Runs",
				SubtitleCalculated: "Amounts locked and pending approval",
				HeadingApproved:    "Approved Payroll Runs",
				SubtitleApproved:   "Approved and ready for disbursement",
				HeadingPosted:      "Posted Payroll Runs",
				SubtitlePosted:     "Disbursement completed and journal entry created",
			},
			Tabs: PayrollRunTabLabels{
				Draft:      "Draft",
				Calculated: "Calculated",
				Approved:   "Approved",
				Posted:     "Posted",
			},
			Buttons: PayrollRunButtonLabels{
				NewRun: "New Payroll Run",
			},
			Columns: PayrollRunColumnLabels{
				RunNumber:       "Run #",
				PayPeriod:       "Pay Period",
				Employees:       "Employees",
				TotalGross:      "Total Gross",
				TotalDeductions: "Deductions",
				TotalNet:        "Net Pay",
				Status:          "Status",
				ApprovedBy:      "Approved By",
				PostedAt:        "Posted At",
			},
			Empty: PayrollRunEmptyLabels{
				Title:   "No payroll runs found",
				Message: "Create a new payroll run to start processing employee salaries.",
			},
			Actions: PayrollRunActionLabels{
				View:         "View",
				NoPermission: "No permission",
			},
		},
		Remittance: PayrollRemittanceLabels{
			Page: PayrollRemittancePageLabels{
				HeadingPending:  "Pending Remittances",
				SubtitlePending: "Government contributions due for filing and payment",
				HeadingFiled:    "Filed Remittances",
				SubtitleFiled:   "Remittances filed with the government agency",
				HeadingPaid:     "Paid Remittances",
				SubtitlePaid:    "Remittances confirmed paid to the government agency",
			},
			Tabs: PayrollRemittanceTabLabels{
				Pending: "Pending",
				Filed:   "Filed",
				Paid:    "Paid",
			},
			Columns: PayrollRemittanceColumnLabels{
				RemittanceType:  "Agency",
				Amount:          "Amount",
				DueDate:         "Due Date",
				Status:          "Status",
				FiledAt:         "Filed At",
				ReferenceNumber: "Reference #",
			},
			Types: PayrollRemittanceTypeLabels{
				SSS:            "SSS",
				PhilHealth:     "PhilHealth",
				PagIBIG:        "Pag-IBIG",
				BIRWithholding: "BIR Withholding",
			},
			Empty: PayrollRemittanceEmptyLabels{
				Title:   "No remittances found",
				Message: "Government contribution remittances will appear here once payroll runs are processed.",
			},
		},
		Employee: PayrollEmployeeLabels{
			Page: PayrollEmployeePageLabels{
				Heading: "Payroll Employees",
				Caption: "Manage employees enrolled in payroll",
			},
			Columns: PayrollEmployeeColumnLabels{
				Name:         "Name",
				Position:     "Position",
				Department:   "Department",
				BasicSalary:  "Basic Salary",
				PayFrequency: "Pay Frequency",
				Status:       "Status",
			},
			Status: PayrollEmployeeStatusLabels{
				Active:   "Active",
				Inactive: "Inactive",
			},
			PayFrequency: PayrollEmployeePayFrequencyLabels{
				SemiMonthly: "Semi-Monthly",
				Monthly:     "Monthly",
				Weekly:      "Weekly",
			},
			Empty: PayrollEmployeeEmptyLabels{
				Title:   "No employees found",
				Message: "Add employees to payroll to begin processing salaries.",
			},
		},
		Settings: PayrollSettingsLabels{
			GovRates: PayrollGovRatesLabels{
				Page: PayrollGovRatesPageLabels{
					Heading: "Government Contribution Rates",
					Caption: "Philippine mandatory contribution rates \u2014 SSS, PhilHealth, Pag-IBIG, BIR",
				},
				Agency: PayrollGovRatesAgencyLabels{
					SSS:            "SSS (Social Security System)",
					PhilHealth:     "PhilHealth",
					PagIBIG:        "Pag-IBIG (HDMF)",
					BIRWithholding: "BIR Withholding Tax",
				},
			},
			PayPeriods: PayrollPayPeriodsLabels{
				Page: PayrollPayPeriodsPageLabels{
					Heading: "Pay Period Settings",
					Caption: "Configure payroll cut-off dates and pay schedules",
				},
			},
		},
		Dashboard: PayrollDashboardLabels{
			Title:                  "Payroll Dashboard",
			Subtitle:               "Run status, monthly gross-pay trend, and upcoming government remittances",
			CurrentRunStatus:       "Current Run Status",
			EmployeesInCurrent:     "Employees in Current Run",
			TotalGrossMTD:          "Total Gross (MTD)",
			RemittancesDue:         "Remittances Due (30d)",
			GrossPayByMonth:        "Gross Pay by Month",
			RecentRuns:             "Recent Payroll Runs",
			UpcomingRemittances:    "Upcoming Remittance Deadlines",
			NoRecentRuns:           "No recent payroll runs",
			NoUpcomingDeadlines:    "No upcoming remittance deadlines",
			QuickNewRun:            "New Payroll Run",
			QuickProcessRun:        "Process Run",
			QuickFileRemittance:    "File Remittance",
			QuickPayPeriodSettings: "Pay Period Settings",
			ViewAll:                "View All",
			AxisGross:              "Gross",
			NoRunYet:               "No run yet",
		},
	}
}

// ---------------------------------------------------------------------------
// PrepaymentLabels (Expenses — Prepayments)
// ---------------------------------------------------------------------------

// PrepaymentLabels holds all translatable strings for the Prepayments module.
type PrepaymentLabels struct {
	Page    PrepaymentPageLabels   `json:"page"`
	Buttons PrepaymentButtonLabels `json:"buttons"`
	Columns PrepaymentColumnLabels `json:"columns"`
	Status  PrepaymentStatusLabels `json:"status"`
	Empty   PrepaymentEmptyLabels  `json:"empty"`
	Form    PrepaymentFormLabels   `json:"form"`
	Actions PrepaymentActionLabels `json:"actions"`
}

type PrepaymentPageLabels struct {
	Heading             string `json:"heading"`
	Caption             string `json:"caption"`
	AmortizationHeading string `json:"amortizationHeading"`
	AmortizationCaption string `json:"amortizationCaption"`
}

type PrepaymentButtonLabels struct {
	AddPrepayment string `json:"addPrepayment"`
}

type PrepaymentColumnLabels struct {
	Description        string `json:"description"`
	Vendor             string `json:"vendor"`
	TotalAmount        string `json:"totalAmount"`
	RemainingAmount    string `json:"remainingAmount"`
	AmortizationMonths string `json:"amortizationMonths"`
	StartDate          string `json:"startDate"`
	EndDate            string `json:"endDate"`
	Status             string `json:"status"`
	// Amortization schedule sub-table
	Month   string `json:"month"`
	Opening string `json:"opening"`
	Expense string `json:"expense"`
	Closing string `json:"closing"`
}

type PrepaymentStatusLabels struct {
	Active    string `json:"active"`
	Amortized string `json:"amortized"`
	Cancelled string `json:"cancelled"`
}

type PrepaymentEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type PrepaymentFormLabels struct {
	Description             string `json:"description"`
	DescriptionPlaceholder  string `json:"descriptionPlaceholder"`
	Vendor                  string `json:"vendor"`
	VendorPlaceholder       string `json:"vendorPlaceholder"`
	TotalAmount             string `json:"totalAmount"`
	TotalAmountPlaceholder  string `json:"totalAmountPlaceholder"`
	AmortizationMonths      string `json:"amortizationMonths"`
	AmortizationPlaceholder string `json:"amortizationPlaceholder"`
	StartDate               string `json:"startDate"`
	EndDate                 string `json:"endDate"`
	PrepaidAccount          string `json:"prepaidAccount"`
	ExpenseAccount          string `json:"expenseAccount"`
}

type PrepaymentActionLabels struct {
	View          string `json:"view"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultPrepaymentLabels returns PrepaymentLabels with hardcoded English defaults.
func DefaultPrepaymentLabels() PrepaymentLabels {
	return PrepaymentLabels{
		Page: PrepaymentPageLabels{
			Heading:             "Prepayments",
			Caption:             "Track prepaid expenses and their amortization schedules",
			AmortizationHeading: "Amortization Schedule",
			AmortizationCaption: "Monthly expense recognition for active prepayments",
		},
		Buttons: PrepaymentButtonLabels{
			AddPrepayment: "Add Prepayment",
		},
		Columns: PrepaymentColumnLabels{
			Description:        "Description",
			Vendor:             "Vendor",
			TotalAmount:        "Total Amount",
			RemainingAmount:    "Remaining",
			AmortizationMonths: "Months",
			StartDate:          "Start Date",
			EndDate:            "End Date",
			Status:             "Status",
			Month:              "Month",
			Opening:            "Opening Balance",
			Expense:            "Monthly Expense",
			Closing:            "Closing Balance",
		},
		Status: PrepaymentStatusLabels{
			Active:    "Active",
			Amortized: "Fully Amortized",
			Cancelled: "Cancelled",
		},
		Empty: PrepaymentEmptyLabels{
			Title:   "No prepayments found",
			Message: "Record prepaid expenses such as insurance, rent, and subscriptions paid in advance.",
		},
		Form: PrepaymentFormLabels{
			Description:             "Description",
			DescriptionPlaceholder:  "e.g. Annual insurance premium",
			Vendor:                  "Vendor",
			VendorPlaceholder:       "e.g. ABC Insurance Co.",
			TotalAmount:             "Total Amount",
			TotalAmountPlaceholder:  "0.00",
			AmortizationMonths:      "Amortization Period (Months)",
			AmortizationPlaceholder: "e.g. 12",
			StartDate:               "Start Date",
			EndDate:                 "End Date",
			PrepaidAccount:          "Prepaid Account (Asset)",
			ExpenseAccount:          "Expense Account",
		},
		Actions: PrepaymentActionLabels{
			View:          "View",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this prepayment? This action cannot be undone.",
		},
	}
}

// ---------------------------------------------------------------------------
// DepositLabels (Cash — Security Deposits)
// ---------------------------------------------------------------------------

// DepositLabels holds all translatable strings for the Security Deposits module.
type DepositLabels struct {
	Page    DepositPageLabels   `json:"page"`
	Buttons DepositButtonLabels `json:"buttons"`
	Columns DepositColumnLabels `json:"columns"`
	Tabs    DepositTabLabels    `json:"tabs"`
	Status  DepositStatusLabels `json:"status"`
	Empty   DepositEmptyLabels  `json:"empty"`
	Form    DepositFormLabels   `json:"form"`
	Actions DepositActionLabels `json:"actions"`
}

type DepositPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type DepositButtonLabels struct {
	RecordDeposit string `json:"recordDeposit"`
}

type DepositColumnLabels struct {
	Counterparty string `json:"counterparty"`
	Direction    string `json:"direction"`
	Amount       string `json:"amount"`
	DepositDate  string `json:"depositDate"`
	Status       string `json:"status"`
	Account      string `json:"account"`
	Notes        string `json:"notes"`
}

type DepositTabLabels struct {
	Paid     string `json:"paid"`
	Received string `json:"received"`
	All      string `json:"all"`
}

type DepositStatusLabels struct {
	Held      string `json:"held"`
	Returned  string `json:"returned"`
	Forfeited string `json:"forfeited"`
}

type DepositEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type DepositFormLabels struct {
	Counterparty            string `json:"counterparty"`
	CounterpartyPlaceholder string `json:"counterpartyPlaceholder"`
	Direction               string `json:"direction"`
	DirectionPaid           string `json:"directionPaid"`
	DirectionReceived       string `json:"directionReceived"`
	Amount                  string `json:"amount"`
	AmountPlaceholder       string `json:"amountPlaceholder"`
	DepositDate             string `json:"depositDate"`
	Account                 string `json:"account"`
	Notes                   string `json:"notes"`
	NotesPlaceholder        string `json:"notesPlaceholder"`
}

type DepositActionLabels struct {
	View          string `json:"view"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultDepositLabels returns DepositLabels with hardcoded English defaults.
func DefaultDepositLabels() DepositLabels {
	return DepositLabels{
		Page: DepositPageLabels{
			Heading: "Security Deposits",
			Caption: "Track deposits paid to landlords, suppliers, and those received from customers",
		},
		Buttons: DepositButtonLabels{
			RecordDeposit: "Record Deposit",
		},
		Columns: DepositColumnLabels{
			Counterparty: "Counterparty",
			Direction:    "Direction",
			Amount:       "Amount",
			DepositDate:  "Deposit Date",
			Status:       "Status",
			Account:      "Account",
			Notes:        "Notes",
		},
		Tabs: DepositTabLabels{
			Paid:     "Paid (Asset)",
			Received: "Received (Liability)",
			All:      "All",
		},
		Status: DepositStatusLabels{
			Held:      "Held",
			Returned:  "Returned",
			Forfeited: "Forfeited",
		},
		Empty: DepositEmptyLabels{
			Title:   "No security deposits found",
			Message: "Record deposits paid to landlords or suppliers, and deposits received from customers.",
		},
		Form: DepositFormLabels{
			Counterparty:            "Counterparty",
			CounterpartyPlaceholder: "e.g. ABC Landlord / Customer name",
			Direction:               "Direction",
			DirectionPaid:           "Paid (We paid the deposit)",
			DirectionReceived:       "Received (We received the deposit)",
			Amount:                  "Amount",
			AmountPlaceholder:       "0.00",
			DepositDate:             "Deposit Date",
			Account:                 "GL Account",
			Notes:                   "Notes",
			NotesPlaceholder:        "Optional notes or reference number",
		},
		Actions: DepositActionLabels{
			View:          "View",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this security deposit? This action cannot be undone.",
		},
	}
}

// ---------------------------------------------------------------------------
// PettyCashLabels (Cash — Petty Cash)
// ---------------------------------------------------------------------------

// PettyCashLabels holds all translatable strings for the Petty Cash module.
type PettyCashLabels struct {
	Page    PettyCashPageLabels   `json:"page"`
	Buttons PettyCashButtonLabels `json:"buttons"`
	Columns PettyCashColumnLabels `json:"columns"`
	Status  PettyCashStatusLabels `json:"status"`
	Empty   PettyCashEmptyLabels  `json:"empty"`
	Form    PettyCashFormLabels   `json:"form"`
	Actions PettyCashActionLabels `json:"actions"`
}

type PettyCashPageLabels struct {
	RegisterHeading          string `json:"registerHeading"`
	RegisterCaption          string `json:"registerCaption"`
	ReplenishmentsHeading    string `json:"replenishmentsHeading"`
	ReplenishmentsCaption    string `json:"replenishmentsCaption"`
	CustodianBalancesHeading string `json:"custodianBalancesHeading"`
	CustodianBalancesCaption string `json:"custodianBalancesCaption"`
}

type PettyCashButtonLabels struct {
	AddFund   string `json:"addFund"`
	Replenish string `json:"replenish"`
}

type PettyCashColumnLabels struct {
	// Register columns
	Name             string `json:"name"`
	AuthorizedAmount string `json:"authorizedAmount"`
	CurrentBalance   string `json:"currentBalance"`
	Custodian        string `json:"custodian"`
	Location         string `json:"location"`
	Status           string `json:"status"`
	// Replenishment columns
	Fund   string `json:"fund"`
	Amount string `json:"amount"`
	Date   string `json:"date"`
	Notes  string `json:"notes"`
	// Custodian balance columns
	TotalFunds   string `json:"totalFunds"`
	TotalBalance string `json:"totalBalance"`
}

type PettyCashStatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type PettyCashEmptyLabels struct {
	RegisterTitle         string `json:"registerTitle"`
	RegisterMessage       string `json:"registerMessage"`
	ReplenishmentsTitle   string `json:"replenishmentsTitle"`
	ReplenishmentsMessage string `json:"replenishmentsMessage"`
	CustodianTitle        string `json:"custodianTitle"`
	CustodianMessage      string `json:"custodianMessage"`
}

type PettyCashFormLabels struct {
	Name                  string `json:"name"`
	NamePlaceholder       string `json:"namePlaceholder"`
	AuthorizedAmount      string `json:"authorizedAmount"`
	AuthorizedPlaceholder string `json:"authorizedPlaceholder"`
	CustodianID           string `json:"custodianId"`
	LocationID            string `json:"locationId"`
}

type PettyCashActionLabels struct {
	View          string `json:"view"`
	Replenish     string `json:"replenish"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultPettyCashLabels returns PettyCashLabels with hardcoded English defaults.
func DefaultPettyCashLabels() PettyCashLabels {
	return PettyCashLabels{
		Page: PettyCashPageLabels{
			RegisterHeading:          "Petty Cash Register",
			RegisterCaption:          "Manage petty cash funds across locations and custodians",
			ReplenishmentsHeading:    "Petty Cash Replenishments",
			ReplenishmentsCaption:    "Track fund replenishments and reimbursements",
			CustodianBalancesHeading: "Custodian Balances",
			CustodianBalancesCaption: "Current balance summary by custodian",
		},
		Buttons: PettyCashButtonLabels{
			AddFund:   "Add Fund",
			Replenish: "Replenish",
		},
		Columns: PettyCashColumnLabels{
			Name:             "Fund Name",
			AuthorizedAmount: "Authorized Amount",
			CurrentBalance:   "Current Balance",
			Custodian:        "Custodian",
			Location:         "Location",
			Status:           "Status",
			Fund:             "Fund",
			Amount:           "Amount",
			Date:             "Date",
			Notes:            "Notes",
			TotalFunds:       "Total Funds",
			TotalBalance:     "Total Balance",
		},
		Status: PettyCashStatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		Empty: PettyCashEmptyLabels{
			RegisterTitle:         "No petty cash funds",
			RegisterMessage:       "Set up petty cash funds for each location or department.",
			ReplenishmentsTitle:   "No replenishments",
			ReplenishmentsMessage: "Replenishment records will appear here when funds are restocked.",
			CustodianTitle:        "No custodian data",
			CustodianMessage:      "Assign custodians to funds to see balance summaries here.",
		},
		Form: PettyCashFormLabels{
			Name:                  "Fund Name",
			NamePlaceholder:       "e.g. Main Office Petty Cash",
			AuthorizedAmount:      "Authorized Amount",
			AuthorizedPlaceholder: "0.00",
			CustodianID:           "Custodian",
			LocationID:            "Location",
		},
		Actions: PettyCashActionLabels{
			View:          "View",
			Replenish:     "Replenish",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this petty cash fund? This action cannot be undone.",
		},
	}
}

// ---------------------------------------------------------------------------
// DeferredRevenueLabels (Revenue — Deferred Revenue)
// ---------------------------------------------------------------------------

// DeferredRevenueLabels holds all translatable strings for the Deferred Revenue module.
type DeferredRevenueLabels struct {
	Page    DeferredRevenuePageLabels   `json:"page"`
	Buttons DeferredRevenueButtonLabels `json:"buttons"`
	Columns DeferredRevenueColumnLabels `json:"columns"`
	Status  DeferredRevenueStatusLabels `json:"status"`
	Empty   DeferredRevenueEmptyLabels  `json:"empty"`
	Form    DeferredRevenueFormLabels   `json:"form"`
	Actions DeferredRevenueActionLabels `json:"actions"`
}

type DeferredRevenuePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type DeferredRevenueButtonLabels struct {
	AddDeferredRevenue string `json:"addDeferredRevenue"`
}

type DeferredRevenueColumnLabels struct {
	Description       string `json:"description"`
	Customer          string `json:"customer"`
	TotalAmount       string `json:"totalAmount"`
	RecognizedAmount  string `json:"recognizedAmount"`
	RemainingAmount   string `json:"remainingAmount"`
	RecognitionMonths string `json:"recognitionMonths"`
	StartDate         string `json:"startDate"`
	EndDate           string `json:"endDate"`
	Status            string `json:"status"`
}

type DeferredRevenueStatusLabels struct {
	Pending    string `json:"pending"`
	Active     string `json:"active"`
	Recognized string `json:"recognized"`
	Cancelled  string `json:"cancelled"`
}

type DeferredRevenueEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type DeferredRevenueFormLabels struct {
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
	Customer               string `json:"customer"`
	CustomerPlaceholder    string `json:"customerPlaceholder"`
	TotalAmount            string `json:"totalAmount"`
	TotalAmountPlaceholder string `json:"totalAmountPlaceholder"`
	RecognitionMonths      string `json:"recognitionMonths"`
	RecognitionPlaceholder string `json:"recognitionPlaceholder"`
	StartDate              string `json:"startDate"`
	EndDate                string `json:"endDate"`
	LiabilityAccount       string `json:"liabilityAccount"`
	RevenueAccount         string `json:"revenueAccount"`
}

type DeferredRevenueActionLabels struct {
	View          string `json:"view"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultDeferredRevenueLabels returns DeferredRevenueLabels with hardcoded English defaults.
func DefaultDeferredRevenueLabels() DeferredRevenueLabels {
	return DeferredRevenueLabels{
		Page: DeferredRevenuePageLabels{
			Heading: "Deferred Revenue",
			Caption: "Track revenue received in advance and its recognition schedule",
		},
		Buttons: DeferredRevenueButtonLabels{
			AddDeferredRevenue: "Record Deferred Revenue",
		},
		Columns: DeferredRevenueColumnLabels{
			Description:       "Description",
			Customer:          "Customer",
			TotalAmount:       "Total Amount",
			RecognizedAmount:  "Recognized",
			RemainingAmount:   "Remaining",
			RecognitionMonths: "Months",
			StartDate:         "Start Date",
			EndDate:           "End Date",
			Status:            "Status",
		},
		Status: DeferredRevenueStatusLabels{
			Pending:    "Pending",
			Active:     "Active",
			Recognized: "Fully Recognized",
			Cancelled:  "Cancelled",
		},
		Empty: DeferredRevenueEmptyLabels{
			Title:   "No deferred revenue found",
			Message: "Record advance payments from customers that will be earned over future periods.",
		},
		Form: DeferredRevenueFormLabels{
			Description:            "Description",
			DescriptionPlaceholder: "e.g. 12-month service contract",
			Customer:               "Customer",
			CustomerPlaceholder:    "e.g. XYZ Corp",
			TotalAmount:            "Total Amount",
			TotalAmountPlaceholder: "0.00",
			RecognitionMonths:      "Recognition Period (Months)",
			RecognitionPlaceholder: "e.g. 12",
			StartDate:              "Start Date",
			EndDate:                "End Date",
			LiabilityAccount:       "Deferred Revenue Account (Liability)",
			RevenueAccount:         "Revenue Account",
		},
		Actions: DeferredRevenueActionLabels{
			View:          "View",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this deferred revenue record? This action cannot be undone.",
		},
	}
}

// ---------------------------------------------------------------------------
// Equity labels (Funding > Equity app)
// ---------------------------------------------------------------------------

// EquityLabels is the top-level label container for the Equity app.
type EquityLabels struct {
	Accounts     EquityAccountLabels     `json:"accounts"`
	Transactions EquityTransactionLabels `json:"transactions"`
	Sheet        EquitySheetLabels       `json:"sheet"`
	Dashboard    EquityDashboardLabels   `json:"dashboard"`
}

// EquityDashboardLabels holds translatable strings for the Equity live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type EquityDashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	TotalContributed string `json:"totalContributed"`
	ActiveOwners     string `json:"activeOwners"`
	DistributionsYTD string `json:"distributionsYtd"`
	NetMovementYTD   string `json:"netMovementYtd"`
	// Widgets
	EquityByOwner      string `json:"equityByOwner"`
	TopContributors    string `json:"topContributors"`
	RecentTransactions string `json:"recentTransactions"`
	NoRecentTxns       string `json:"noRecentTxns"`
	// Quick actions
	QuickRecordContribution string `json:"quickRecordContribution"`
	QuickRecordDistribution string `json:"quickRecordDistribution"`
	QuickOwnerStatement     string `json:"quickOwnerStatement"`
	QuickEquityReport       string `json:"quickEquityReport"`
	// Common
	ViewAll    string `json:"viewAll"`
	AxisAmount string `json:"axisAmount"`
}

// EquitySheetLabels holds sheet-form title and button labels for equity pages.
type EquitySheetLabels struct {
	AddCapitalAccount       string `json:"addCapitalAccount"`
	RecordTransaction       string `json:"recordTransaction"`
	RecordEquityTransaction string `json:"recordEquityTransaction"`
	PostTransaction         string `json:"postTransaction"`
}

// EquityAccountLabels holds translatable strings for the capital accounts list.
type EquityAccountLabels struct {
	Page    EquityAccountPageLabels   `json:"page"`
	Buttons EquityAccountButtonLabels `json:"buttons"`
	Columns EquityAccountColumnLabels `json:"columns"`
	Empty   EquityAccountEmptyLabels  `json:"empty"`
	Actions EquityAccountActionLabels `json:"actions"`
}

type EquityAccountPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type EquityAccountButtonLabels struct {
	AddAccount string `json:"addAccount"`
}

type EquityAccountColumnLabels struct {
	Name        string `json:"name"`
	OwnerName   string `json:"ownerName"`
	AccountType string `json:"accountType"`
	Balance     string `json:"balance"`
}

type EquityAccountEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type EquityAccountActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	NoPermission string `json:"noPermission"`
}

// EquityTransactionLabels holds translatable strings for the equity transactions list.
type EquityTransactionLabels struct {
	Page    EquityTransactionPageLabels   `json:"page"`
	Buttons EquityTransactionButtonLabels `json:"buttons"`
	Columns EquityTransactionColumnLabels `json:"columns"`
	Empty   EquityTransactionEmptyLabels  `json:"empty"`
	Actions EquityTransactionActionLabels `json:"actions"`
	Form    EquityTransactionFormLabels   `json:"form"`
}

type EquityTransactionPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type EquityTransactionButtonLabels struct {
	RecordTransaction string `json:"recordTransaction"`
}

type EquityTransactionColumnLabels struct {
	Date            string `json:"date"`
	TransactionType string `json:"transactionType"`
	Amount          string `json:"amount"`
	Description     string `json:"description"`
}

type EquityTransactionEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type EquityTransactionActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

type EquityTransactionFormLabels struct {
	TransactionType         string `json:"transactionType"`
	TransactionContribution string `json:"transactionContribution"`
	TransactionWithdrawal   string `json:"transactionWithdrawal"`
	TransactionDistribution string `json:"transactionDistribution"`
	TransactionTransfer     string `json:"transactionTransfer"`
	EquityAccount           string `json:"equityAccount"`
	Amount                  string `json:"amount"`
	TransactionDate         string `json:"transactionDate"`
	Description             string `json:"description"`
	DescriptionPlaceholder  string `json:"descriptionPlaceholder"`
	JournalEntryHint        string `json:"journalEntryHint"`
	MemoPlaceholder         string `json:"memoPlaceholder"`
	JournalEntryNote        string `json:"journalEntryNote"`
}

// DefaultEquityLabels returns EquityLabels with hardcoded English defaults.
func DefaultEquityLabels() EquityLabels {
	return EquityLabels{
		Accounts: EquityAccountLabels{
			Page: EquityAccountPageLabels{
				Heading: "Capital Accounts",
				Caption: "Track owner equity and capital contributions",
			},
			Buttons: EquityAccountButtonLabels{
				AddAccount: "Add Capital Account",
			},
			Columns: EquityAccountColumnLabels{
				Name:        "Account Name",
				OwnerName:   "Owner",
				AccountType: "Type",
				Balance:     "Balance",
			},
			Empty: EquityAccountEmptyLabels{
				Title:   "No capital accounts",
				Message: "Add your first capital account to start tracking owner equity.",
			},
			Actions: EquityAccountActionLabels{
				View:         "View",
				Edit:         "Edit",
				NoPermission: "No permission",
			},
		},
		Transactions: EquityTransactionLabels{
			Page: EquityTransactionPageLabels{
				Heading: "Equity Transactions",
				Caption: "Capital contributions, withdrawals, and distributions",
			},
			Buttons: EquityTransactionButtonLabels{
				RecordTransaction: "Record Transaction",
			},
			Columns: EquityTransactionColumnLabels{
				Date:            "Date",
				TransactionType: "Type",
				Amount:          "Amount",
				Description:     "Description",
			},
			Empty: EquityTransactionEmptyLabels{
				Title:   "No equity transactions",
				Message: "Record your first equity transaction to start tracking capital movements.",
			},
			Actions: EquityTransactionActionLabels{
				View:         "View",
				NoPermission: "No permission",
			},
			Form: EquityTransactionFormLabels{
				TransactionType:         "Transaction Type",
				EquityAccount:           "Capital Account",
				Amount:                  "Amount",
				TransactionDate:         "Transaction Date",
				Description:             "Description",
				DescriptionPlaceholder:  "Optional memo for this transaction",
				MemoPlaceholder:         "Optional memo for this transaction",
				JournalEntryHint:        "The corresponding journal entry will be auto-generated when you post this transaction.",
				JournalEntryNote:        "The corresponding journal entry will be auto-generated when you post this transaction. Debits and credits are determined by the transaction type selected above.",
				TransactionContribution: "Contribution \u2014 Owner adds capital",
				TransactionWithdrawal:   "Withdrawal \u2014 Owner draws cash",
				TransactionDistribution: "Distribution \u2014 Profit distributed",
				TransactionTransfer:     "Transfer \u2014 Between equity accounts",
			},
		},
		Sheet: EquitySheetLabels{
			AddCapitalAccount:       "Add Capital Account",
			RecordTransaction:       "Record Transaction",
			RecordEquityTransaction: "Record Equity Transaction",
			PostTransaction:         "Post Transaction",
		},
		Dashboard: EquityDashboardLabels{
			Title:                   "Equity Dashboard",
			Subtitle:                "Owner contributions, distributions, and movements across capital accounts",
			TotalContributed:        "Total Contributed",
			ActiveOwners:            "Active Owners",
			DistributionsYTD:        "Distributions YTD",
			NetMovementYTD:          "Net Movement YTD",
			EquityByOwner:           "Equity by Owner",
			TopContributors:         "Top Contributors",
			RecentTransactions:      "Recent Transactions",
			NoRecentTxns:            "No recent equity transactions",
			QuickRecordContribution: "Record Contribution",
			QuickRecordDistribution: "Record Distribution",
			QuickOwnerStatement:     "Owner Statement",
			QuickEquityReport:       "Equity Report",
			ViewAll:                 "View All",
			AxisAmount:              "Amount",
		},
	}
}

// ---------------------------------------------------------------------------
// Loan labels (Funding > Loans app)
// ---------------------------------------------------------------------------

// LoanLabels is the top-level label container for the Loans app.
type LoanLabels struct {
	Page      LoanPageLabels      `json:"page"`
	Tabs      LoanTabLabels       `json:"tabs"`
	Buttons   LoanButtonLabels    `json:"buttons"`
	Columns   LoanColumnLabels    `json:"columns"`
	Empty     LoanEmptyLabels     `json:"empty"`
	Actions   LoanActionLabels    `json:"actions"`
	Form      LoanFormLabels      `json:"form"`
	Status    LoanStatusLabels    `json:"status"`
	Type      LoanTypeLabels      `json:"type"`
	Sheet     LoanSheetLabels     `json:"sheet"`
	Dashboard LoanDashboardLabels `json:"dashboard"`
}

// LoanDashboardLabels holds translatable strings for the Loan live dashboard
// (Phase 2 — Pyeza dashboard block + per-app live dashboards plan).
type LoanDashboardLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Stats
	TotalOutstanding string `json:"totalOutstanding"`
	InterestYTD      string `json:"interestYtd"`
	PaymentsDue30    string `json:"paymentsDue30"`
	DefaultedCount   string `json:"defaultedCount"`
	// Widgets
	OutstandingTrend string `json:"outstandingTrend"`
	TopLoans         string `json:"topLoans"`
	RecentPayments   string `json:"recentPayments"`
	NoRecentPayments string `json:"noRecentPayments"`
	// Quick actions
	QuickNewLoan      string `json:"quickNewLoan"`
	QuickRecordPay    string `json:"quickRecordPay"`
	QuickAmortization string `json:"quickAmortization"`
	QuickLoanCalendar string `json:"quickLoanCalendar"`
	// Common
	ViewAll    string `json:"viewAll"`
	AxisAmount string `json:"axisAmount"`
}

// LoanSheetLabels holds sheet-form title and button labels for loan list page.
type LoanSheetLabels struct {
	AddLoan  string `json:"addLoan"`
	SaveLoan string `json:"saveLoan"`
}

type LoanPageLabels struct {
	HeadingActive    string `json:"headingActive"`
	CaptionActive    string `json:"captionActive"`
	HeadingCompleted string `json:"headingCompleted"`
	CaptionCompleted string `json:"captionCompleted"`
}

type LoanTabLabels struct {
	Active    string `json:"active"`
	Completed string `json:"completed"`
}

type LoanButtonLabels struct {
	AddLoan string `json:"addLoan"`
}

type LoanColumnLabels struct {
	LoanNumber       string `json:"loanNumber"`
	LenderName       string `json:"lenderName"`
	LoanType         string `json:"loanType"`
	PrincipalAmount  string `json:"principalAmount"`
	RemainingBalance string `json:"remainingBalance"`
	InterestRate     string `json:"interestRate"`
	TermMonths       string `json:"termMonths"`
	StartDate        string `json:"startDate"`
	MaturityDate     string `json:"maturityDate"`
	Status           string `json:"status"`
}

type LoanEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type LoanActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
	SaveError    string `json:"saveError"`
}

type LoanFormLabels struct {
	LoanNumber             string `json:"loanNumber"`
	LoanNumberPlaceholder  string `json:"loanNumberPlaceholder"`
	LenderName             string `json:"lenderName"`
	LenderNamePlaceholder  string `json:"lenderNamePlaceholder"`
	LoanType               string `json:"loanType"`
	PrincipalAmount        string `json:"principalAmount"`
	InterestRate           string `json:"interestRate"`
	TermMonths             string `json:"termMonths"`
	StartDate              string `json:"startDate"`
	MaturityDate           string `json:"maturityDate"`
	Description            string `json:"description"`
	DescriptionPlaceholder string `json:"descriptionPlaceholder"`
}

type LoanStatusLabels struct {
	Draft     string `json:"draft"`
	Active    string `json:"active"`
	Completed string `json:"completed"`
	Defaulted string `json:"defaulted"`
}

type LoanTypeLabels struct {
	Payable    string `json:"payable"`
	Receivable string `json:"receivable"`
}

// DefaultLoanLabels returns LoanLabels with hardcoded English defaults.
func DefaultLoanLabels() LoanLabels {
	return LoanLabels{
		Page: LoanPageLabels{
			HeadingActive:    "Active Loans",
			CaptionActive:    "Loans currently being serviced",
			HeadingCompleted: "Completed Loans",
			CaptionCompleted: "Fully paid or closed loans",
		},
		Tabs: LoanTabLabels{
			Active:    "Active",
			Completed: "Completed",
		},
		Buttons: LoanButtonLabels{
			AddLoan: "Add Loan",
		},
		Columns: LoanColumnLabels{
			LoanNumber:       "Loan #",
			LenderName:       "Lender / Borrower",
			LoanType:         "Type",
			PrincipalAmount:  "Principal",
			RemainingBalance: "Balance",
			InterestRate:     "Rate",
			TermMonths:       "Term (mo.)",
			StartDate:        "Start Date",
			MaturityDate:     "Maturity",
			Status:           "Status",
		},
		Empty: LoanEmptyLabels{
			Title:   "No loans found",
			Message: "Add your first loan to start tracking borrowings and repayments.",
		},
		Actions: LoanActionLabels{
			View:         "View",
			NoPermission: "No permission",
			SaveError:    "Failed to save loan",
		},
		Form: LoanFormLabels{
			LoanNumber:             "Loan Number",
			LoanNumberPlaceholder:  "e.g. LN-001",
			LenderName:             "Lender / Borrower",
			LenderNamePlaceholder:  "Name of the lender or borrower",
			LoanType:               "Loan Type",
			PrincipalAmount:        "Principal Amount",
			InterestRate:           "Annual Interest Rate (%)",
			TermMonths:             "Term (Months)",
			StartDate:              "Start Date",
			MaturityDate:           "Maturity Date",
			Description:            "Description",
			DescriptionPlaceholder: "Brief description or purpose of the loan",
		},
		Status: LoanStatusLabels{
			Draft:     "Draft",
			Active:    "Active",
			Completed: "Completed",
			Defaulted: "Defaulted",
		},
		Type: LoanTypeLabels{
			Payable:    "Payable",
			Receivable: "Receivable",
		},
		Sheet: LoanSheetLabels{
			AddLoan:  "Add Loan",
			SaveLoan: "Save Loan",
		},
		Dashboard: LoanDashboardLabels{
			Title:             "Loans Dashboard",
			Subtitle:          "Outstanding balance, interest accrued, upcoming payments, and recent activity",
			TotalOutstanding:  "Total Outstanding",
			InterestYTD:       "Interest YTD",
			PaymentsDue30:     "Payments Due (30d)",
			DefaultedCount:    "Defaulted Loans",
			OutstandingTrend:  "Outstanding Principal Trend",
			TopLoans:          "Top Loans by Outstanding",
			RecentPayments:    "Recent Payments",
			NoRecentPayments:  "No recent payments",
			QuickNewLoan:      "New Loan",
			QuickRecordPay:    "Record Payment",
			QuickAmortization: "Amortization Schedule",
			QuickLoanCalendar: "Loan Calendar",
			ViewAll:           "View All",
			AxisAmount:        "Outstanding",
		},
	}
}

// ---------------------------------------------------------------------------
// LoanPayment labels (Funding > Loans > Payments)
// ---------------------------------------------------------------------------

// LoanPaymentLabels holds all translatable strings for the loan payments view.
type LoanPaymentLabels struct {
	Page    LoanPaymentPageLabels   `json:"page"`
	Buttons LoanPaymentButtonLabels `json:"buttons"`
	Columns LoanPaymentColumnLabels `json:"columns"`
	Empty   LoanPaymentEmptyLabels  `json:"empty"`
	Actions LoanPaymentActionLabels `json:"actions"`
	Form    LoanPaymentFormLabels   `json:"form"`
	Sheet   LoanPaymentSheetLabels  `json:"sheet"`
}

// LoanPaymentSheetLabels holds sheet-form title and button labels for loan payments page.
type LoanPaymentSheetLabels struct {
	RecordPayment string `json:"recordPayment"`
	PostPayment   string `json:"postPayment"`
}

type LoanPaymentPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type LoanPaymentButtonLabels struct {
	RecordPayment string `json:"recordPayment"`
}

type LoanPaymentColumnLabels struct {
	PaymentNumber    string `json:"paymentNumber"`
	PaymentDate      string `json:"paymentDate"`
	PrincipalAmount  string `json:"principalAmount"`
	InterestAmount   string `json:"interestAmount"`
	FeeAmount        string `json:"feeAmount"`
	TotalAmount      string `json:"totalAmount"`
	RemainingBalance string `json:"remainingBalance"`
}

type LoanPaymentEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type LoanPaymentActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

type LoanPaymentFormLabels struct {
	PaymentNumber            string `json:"paymentNumber"`
	PaymentNumberPlaceholder string `json:"paymentNumberPlaceholder"`
	PaymentDate              string `json:"paymentDate"`
	PrincipalAmount          string `json:"principalAmount"`
	InterestAmount           string `json:"interestAmount"`
	FeeAmount                string `json:"feeAmount"`
	TotalAmount              string `json:"totalAmount"`
	RemainingBalance         string `json:"remainingBalance"`
	Notes                    string `json:"notes"`
	NotesPlaceholder         string `json:"notesPlaceholder"`
	PaymentBreakdown         string `json:"paymentBreakdown"`
}

// DefaultLoanPaymentLabels returns LoanPaymentLabels with hardcoded English defaults.
func DefaultLoanPaymentLabels() LoanPaymentLabels {
	return LoanPaymentLabels{
		Page: LoanPaymentPageLabels{
			Heading: "Loan Payments",
			Caption: "Payment history for this loan",
		},
		Buttons: LoanPaymentButtonLabels{
			RecordPayment: "Record Payment",
		},
		Columns: LoanPaymentColumnLabels{
			PaymentNumber:    "Payment #",
			PaymentDate:      "Date",
			PrincipalAmount:  "Principal",
			InterestAmount:   "Interest",
			FeeAmount:        "Fees",
			TotalAmount:      "Total",
			RemainingBalance: "Balance",
		},
		Empty: LoanPaymentEmptyLabels{
			Title:   "No payments recorded",
			Message: "Record the first payment against this loan.",
		},
		Actions: LoanPaymentActionLabels{
			View:         "View",
			NoPermission: "No permission",
		},
		Form: LoanPaymentFormLabels{
			PaymentNumber:            "Payment Number",
			PaymentNumberPlaceholder: "e.g. PAY-001",
			PaymentDate:              "Payment Date",
			PrincipalAmount:          "Principal Amount",
			InterestAmount:           "Interest Amount",
			FeeAmount:                "Fees (PFRS 9)",
			TotalAmount:              "Total Payment",
			RemainingBalance:         "Remaining Balance After Payment",
			Notes:                    "Notes",
			NotesPlaceholder:         "Optional payment notes or reference",
			PaymentBreakdown:         "Payment Breakdown",
		},
		Sheet: LoanPaymentSheetLabels{
			RecordPayment: "Record Payment",
			PostPayment:   "Post Payment",
		},
	}
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

// ---------------------------------------------------------------------------
// Depreciation Run labels (Surfaces A, B, C, D, F) + Revaluation labels (Surface E)
// Lyngua root key: "depreciationRun" / "assetRevaluation" / "depreciationPolicies"
// Naming: depreciationRun / DepreciationRun / depreciation_run / depreciation-run everywhere
// except user-visible VALUES supplied by lyngua (e.g. "Lapsing Schedule").
// ---------------------------------------------------------------------------

// DepreciationRunLabels is the top-level struct for all Depreciation Run copy.
type DepreciationRunLabels struct {
	AppLabel string `json:"appLabel"`
	// Surface A — per-asset drawer
	RunForm DepreciationRunFormLabels `json:"runForm"`
	// Surface B — Lapsing Schedule list page (workspace overview)
	LapsingSchedule DepreciationRunLapsingScheduleLabels `json:"lapsingSchedule"`
	// Surface C — per-category / per-policy drawer
	CategoryRunForm DepreciationRunCategoryFormLabels `json:"categoryRunForm"`
	// Surface D — run history list + detail
	List   DepreciationRunListLabels   `json:"list"`
	Detail DepreciationRunDetailLabels `json:"detail"`
	// Status badges shared across Surfaces B and D
	StatusBadges DepreciationRunStatusBadgeLabels `json:"statusBadges"`
	// Scope kind display labels
	ScopeKind DepreciationRunScopeKindLabels `json:"scopeKind"`
	// Entry outcome labels
	EntryOutcome DepreciationRunEntryOutcomeLabels `json:"entryOutcome"`
	// Cross-cutting toast labels
	Toast DepreciationRunToastLabels `json:"toast"`
	// Errors
	Errors DepreciationRunErrorLabels `json:"errors"`
	// Cross-cutting asset-edit field labels
	AssetEditDepreciationFieldsLockedWarning  string `json:"assetEditDepreciationFieldsLockedWarning"`
	AssetEditUnitsOfProductionDisabledTooltip string `json:"assetEditUnitsOfProductionDisabledTooltip"`
}

// DepreciationRunFormLabels holds copy for the Surface A per-asset drawer.
type DepreciationRunFormLabels struct {
	Title            string `json:"title"`
	SubtitleTemplate string `json:"subtitleTemplate"`
	// AsOfDate input
	AsOfDateLabel string `json:"asOfDateLabel"`
	AsOfDateHint  string `json:"asOfDateHint"`
	// Pending-periods table column headers
	ColPeriod         string `json:"colPeriod"`
	ColAmount         string `json:"colAmount"`
	ColAccumulated    string `json:"colAccumulated"`
	ColBookValueAfter string `json:"colBookValueAfter"`
	// Generate button (label and count-template variant)
	GenerateLabel             string `json:"generateLabel"`
	GenerateWithCountTemplate string `json:"generateWithCountTemplate"`
	CancelLabel               string `json:"cancelLabel"`
	// Blocker chip labels
	BlockerNotInService     string `json:"blockerNotInService"`
	BlockerFullyDepreciated string `json:"blockerFullyDepreciated"`
	BlockerMissingMethod    string `json:"blockerMissingMethod"`
	BlockerUnitsRequired    string `json:"blockerUnitsRequired"`
	// UoP-specific blocker messaging (rendered as a translated message + link in the drawer)
	BlockerUnitsRequiredMessage string `json:"blockerUnitsRequiredMessage"`
	BlockerUnitsRequiredLink    string `json:"blockerUnitsRequiredLink"`
	// Empty state (no pending periods)
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// DepreciationRunLapsingScheduleLabels holds copy for the Surface B workspace
// lapsing-schedule list page.
type DepreciationRunLapsingScheduleLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// AsOfDate input
	AsOfDateLabel string `json:"asOfDateLabel"`
	// Table column labels
	Columns DepreciationRunLapsingScheduleColumnLabels `json:"columns"`
	// Status badge variants
	StatusUpToDate                string `json:"statusUpToDate"`
	StatusNPeriodsPendingTemplate string `json:"statusNPeriodsPendingTemplate"`
	StatusNotStarted              string `json:"statusNotStarted"`
	StatusBlockedTemplate         string `json:"statusBlockedTemplate"`
	// BlockedPrefix is the human-readable prefix for blocked-status badges
	// when a specific BlockerLabel is provided (e.g. "Blocked: Units required").
	// The trailing space is intentional — it is concatenated with the reason string.
	BlockedPrefix string `json:"blockedPrefix"`
	// Bulk action labels
	BulkRunForSelected    string `json:"bulkRunForSelected"`
	BulkRunForAllMatching string `json:"bulkRunForAllMatching"`
	// Filter chip labels
	FilterCategory string `json:"filterCategory"`
	FilterPolicy   string `json:"filterPolicy"`
	FilterStatus   string `json:"filterStatus"`
	FilterCurrency string `json:"filterCurrency"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// DepreciationRunLapsingScheduleColumnLabels holds column headers for Surface B.
type DepreciationRunLapsingScheduleColumnLabels struct {
	Asset             string `json:"asset"`
	Category          string `json:"category"`
	Policy            string `json:"policy"`
	Currency          string `json:"currency"`
	CurrentBookValue  string `json:"currentBookValue"`
	LastPostedPeriod  string `json:"lastPostedPeriod"`
	NextPendingPeriod string `json:"nextPendingPeriod"`
	Pending           string `json:"pending"`
	NextAmount        string `json:"nextAmount"`
	Status            string `json:"status"`
	Actions           string `json:"actions"`
}

// DepreciationRunCategoryFormLabels holds copy for the Surface C per-category /
// per-policy drawer. The same drawer serves both entry points; only the breadcrumb differs.
type DepreciationRunCategoryFormLabels struct {
	// Category breadcrumb variant
	TitleCategory string `json:"titleCategory"`
	// Policy breadcrumb variant
	TitlePolicy      string `json:"titlePolicy"`
	SubtitleTemplate string `json:"subtitleTemplate"`
	// Per-asset row column headers
	ColAsset    string `json:"colAsset"`
	ColMethod   string `json:"colMethod"`
	ColPending  string `json:"colPending"`
	ColAmount   string `json:"colAmount"`
	ColBlockers string `json:"colBlockers"`
	// Submit and cancel
	SubmitLabel string `json:"submitLabel"`
	CancelLabel string `json:"cancelLabel"`
}

// DepreciationRunListLabels holds copy for the Surface D run-history list page.
type DepreciationRunListLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Table column labels
	Columns DepreciationRunListColumnLabels `json:"columns"`
	// Status filter chip labels
	FilterPending  string `json:"filterPending"`
	FilterComplete string `json:"filterComplete"`
	FilterFailed   string `json:"filterFailed"`
	// Stale-pending warning
	StalePendingWarning string `json:"stalePendingWarning"`
	// Per-status empty states
	Empty DepreciationRunListEmptyLabels `json:"empty"`
}

// DepreciationRunListColumnLabels holds column headers for Surface D list.
type DepreciationRunListColumnLabels struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AsOfDate    string `json:"asOfDate"`
	Initiator   string `json:"initiator"`
	InitiatedAt string `json:"initiatedAt"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Skipped     string `json:"skipped"`
	Errored     string `json:"errored"`
	Actions     string `json:"actions"`
}

// DepreciationRunListEmptyLabels holds per-status empty-state copy for Surface D list.
type DepreciationRunListEmptyLabels struct {
	Pending  DepreciationRunListEmptyStateLabels `json:"pending"`
	Complete DepreciationRunListEmptyStateLabels `json:"complete"`
	Failed   DepreciationRunListEmptyStateLabels `json:"failed"`
}

// DepreciationRunListEmptyStateLabels holds title + message for one empty-state variant.
type DepreciationRunListEmptyStateLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DepreciationRunDetailLabels holds copy for the Surface D run-detail page.
type DepreciationRunDetailLabels struct {
	Title   string                         `json:"title"`
	Tabs    DepreciationRunDetailTabLabels `json:"tabs"`
	Summary DepreciationRunSummaryLabels   `json:"summary"`
}

// DepreciationRunDetailTabLabels holds tab labels for the Surface D detail page.
type DepreciationRunDetailTabLabels struct {
	Summary      string `json:"summary"`
	Selections   string `json:"selections"`
	Results      string `json:"results"`
	Transactions string `json:"transactions"`
	History      string `json:"history"`
}

// DepreciationRunSummaryLabels holds stat-card labels for the Surface D summary tab.
type DepreciationRunSummaryLabels struct {
	Scope                   string `json:"scope"`
	AsOfDate                string `json:"asOfDate"`
	Initiator               string `json:"initiator"`
	InitiatedAt             string `json:"initiatedAt"`
	CompletedAt             string `json:"completedAt"`
	Status                  string `json:"status"`
	Created                 string `json:"created"`
	Skipped                 string `json:"skipped"`
	Errored                 string `json:"errored"`
	Totals                  string `json:"totals"`
	PossiblyInterruptedNote string `json:"possiblyInterruptedNote"`
}

// DepreciationRunStatusBadgeLabels holds display labels for each run status value.
type DepreciationRunStatusBadgeLabels struct {
	Pending             string `json:"pending"`
	Complete            string `json:"complete"`
	Failed              string `json:"failed"`
	PossiblyInterrupted string `json:"possiblyInterrupted"`
}

// DepreciationRunScopeKindLabels holds display labels for each scope kind enum value.
type DepreciationRunScopeKindLabels struct {
	Asset     string `json:"asset"`
	Category  string `json:"category"`
	Policy    string `json:"policy"`
	Workspace string `json:"workspace"`
}

// DepreciationRunEntryOutcomeLabels holds display labels for per-entry outcome values.
type DepreciationRunEntryOutcomeLabels struct {
	Created string `json:"created"`
	Skipped string `json:"skipped"`
	Errored string `json:"errored"`
}

// DepreciationRunToastLabels holds toast message templates for all Depreciation Run surfaces.
type DepreciationRunToastLabels struct {
	// SuccessTemplate supports {{.Created}}/{{.Skipped}}/{{.Errored}} placeholders.
	SuccessTemplate string `json:"successTemplate"`
	// SkippedTemplate is shown when created_count=0 and skipped_count>0.
	SkippedTemplate string `json:"skippedTemplate"`
	// ErroredTemplate is shown when errored_count>0.
	ErroredTemplate string `json:"erroredTemplate"`
	// ViewRunLink is the link label used on single-run toasts.
	ViewRunLink string `json:"viewRunLink"`
}

// DepreciationRunErrorLabels holds error message strings for the depreciation-run module.
type DepreciationRunErrorLabels struct {
	RunForSelectedCapExceededError string `json:"runForSelectedCapExceededError"`
	PermissionDenied               string `json:"permissionDenied"`
	UseCaseUnavailable             string `json:"useCaseUnavailable"`
	FormParseFailed                string `json:"formParseFailed"`
	GenerateFailed                 string `json:"generateFailed"`
	InvalidSelection               string `json:"invalidSelection"`
	IdempotencyConflict            string `json:"idempotencyConflict"`
	WorkspaceMismatch              string `json:"workspaceMismatch"`
}

// DefaultDepreciationRunLabels returns DepreciationRunLabels with sensible English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultDepreciationRunLabels() DepreciationRunLabels {
	return DepreciationRunLabels{
		AppLabel: "Lapsing Schedule",
		RunForm: DepreciationRunFormLabels{
			Title:                       "Run Depreciation",
			SubtitleTemplate:            "Post depreciation for {{.AssetName}} through {{.AsOfDate}}",
			AsOfDateLabel:               "As of date",
			AsOfDateHint:                "Periods up to and including this date will be posted.",
			ColPeriod:                   "Period",
			ColAmount:                   "Amount",
			ColAccumulated:              "Accumulated",
			ColBookValueAfter:           "Book value after",
			GenerateLabel:               "Generate",
			GenerateWithCountTemplate:   "Generate ({{.Count}})",
			CancelLabel:                 "Cancel",
			BlockerNotInService:         "Not in service",
			BlockerFullyDepreciated:     "Fully depreciated",
			BlockerMissingMethod:        "Missing depreciation method",
			BlockerUnitsRequired:        "Units of Production not yet supported",
			BlockerUnitsRequiredMessage: "Units of Production depreciation requires per-period units input. See the future-release plan.",
			BlockerUnitsRequiredLink:    "Open future-release plan",
			EmptyTitle:                  "No pending periods",
			EmptyMessage:                "All periods up to the selected date have been posted.",
		},
		LapsingSchedule: DepreciationRunLapsingScheduleLabels{
			Title:         "Lapsing Schedule",
			Subtitle:      "In-service assets with pending depreciation periods",
			AsOfDateLabel: "As of date",
			Columns: DepreciationRunLapsingScheduleColumnLabels{
				Asset:             "Asset",
				Category:          "Category",
				Policy:            "Policy",
				Currency:          "Currency",
				CurrentBookValue:  "Book Value",
				LastPostedPeriod:  "Last posted",
				NextPendingPeriod: "Next pending",
				Pending:           "Pending",
				NextAmount:        "Next amount",
				Status:            "Status",
				Actions:           "Actions",
			},
			StatusUpToDate:                "Up to date",
			StatusNPeriodsPendingTemplate: "{{.Count}} periods pending",
			StatusNotStarted:              "Not started",
			StatusBlockedTemplate:         "Blocked: {{.Reason}}",
			BlockedPrefix:                 "Blocked: ",
			BulkRunForSelected:            "Run for selected",
			BulkRunForAllMatching:         "Run for all matching",
			FilterCategory:                "Category",
			FilterPolicy:                  "Policy",
			FilterStatus:                  "Status",
			FilterCurrency:                "Currency",
			EmptyTitle:                    "No assets in service",
			EmptyMessage:                  "Add an asset to get started.",
		},
		CategoryRunForm: DepreciationRunCategoryFormLabels{
			TitleCategory:    "Run depreciation for category",
			TitlePolicy:      "Run depreciation for policy",
			SubtitleTemplate: "{{.Count}} assets eligible",
			ColAsset:         "Asset",
			ColMethod:        "Method",
			ColPending:       "Pending",
			ColAmount:        "Amount",
			ColBlockers:      "Blockers",
			SubmitLabel:      "Run depreciation",
			CancelLabel:      "Cancel",
		},
		List: DepreciationRunListLabels{
			Title:    "Depreciation Runs",
			Subtitle: "History of depreciation run batches",
			Columns: DepreciationRunListColumnLabels{
				ID:          "Run ID",
				Scope:       "Scope",
				AsOfDate:    "As of date",
				Initiator:   "Initiator",
				InitiatedAt: "Initiated",
				Status:      "Status",
				Created:     "Created",
				Skipped:     "Skipped",
				Errored:     "Errored",
				Actions:     "Actions",
			},
			FilterPending:       "Pending",
			FilterComplete:      "Complete",
			FilterFailed:        "Failed",
			StalePendingWarning: "This run has been pending for an unusually long time and may have been interrupted.",
			Empty: DepreciationRunListEmptyLabels{
				Pending: DepreciationRunListEmptyStateLabels{
					Title:   "No pending runs",
					Message: "There are no depreciation runs currently in progress.",
				},
				Complete: DepreciationRunListEmptyStateLabels{
					Title:   "No completed runs",
					Message: "No depreciation runs have completed yet.",
				},
				Failed: DepreciationRunListEmptyStateLabels{
					Title:   "No failed runs",
					Message: "No depreciation runs have failed.",
				},
			},
		},
		Detail: DepreciationRunDetailLabels{
			Title: "Depreciation Run",
			Tabs: DepreciationRunDetailTabLabels{
				Summary:      "Summary",
				Selections:   "Selections",
				Results:      "Results",
				Transactions: "Transactions",
				History:      "History",
			},
			Summary: DepreciationRunSummaryLabels{
				Scope:                   "Scope",
				AsOfDate:                "As of date",
				Initiator:               "Initiator",
				InitiatedAt:             "Initiated",
				CompletedAt:             "Completed",
				Status:                  "Status",
				Created:                 "Created",
				Skipped:                 "Skipped",
				Errored:                 "Errored",
				Totals:                  "Totals",
				PossiblyInterruptedNote: "This run may have been interrupted before completing. Some periods may be missing.",
			},
		},
		StatusBadges: DepreciationRunStatusBadgeLabels{
			Pending:             "Pending",
			Complete:            "Complete",
			Failed:              "Failed",
			PossiblyInterrupted: "Possibly interrupted",
		},
		ScopeKind: DepreciationRunScopeKindLabels{
			Asset:     "Asset",
			Category:  "Category",
			Policy:    "Policy",
			Workspace: "Workspace",
		},
		EntryOutcome: DepreciationRunEntryOutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		Toast: DepreciationRunToastLabels{
			SuccessTemplate: "{{.Created}} periods posted, {{.Skipped}} skipped, {{.Errored}} errored",
			SkippedTemplate: "{{.Skipped}} periods already posted (skipped)",
			ErroredTemplate: "{{.Errored}} periods failed to post",
			ViewRunLink:     "View run",
		},
		Errors: DepreciationRunErrorLabels{
			RunForSelectedCapExceededError: "Batch cap exceeded — maximum 500 assets per run. Narrow the filter to run the rest.",
			PermissionDenied:               "You do not have permission to run depreciation.",
			UseCaseUnavailable:             "Service unavailable. Please try again.",
			FormParseFailed:                "Form data could not be read.",
			GenerateFailed:                 "Failed to record the depreciation run.",
			InvalidSelection:               "One or more selected assets are invalid.",
			IdempotencyConflict:            "Depreciation for one or more periods has already been posted.",
			WorkspaceMismatch:              "Selected assets belong to a different workspace.",
		},
		AssetEditDepreciationFieldsLockedWarning:  "Posted depreciation exists for this asset. Changing depreciation configuration requires a Useful Life Change action (not yet available — see Run history for posted periods).",
		AssetEditUnitsOfProductionDisabledTooltip: "Units of Production depreciation is not yet supported.",
	}
}

// ---------------------------------------------------------------------------
// Asset Revaluation labels (Surface E)
// Lyngua root key: "assetRevaluation"
// ---------------------------------------------------------------------------

// AssetRevaluationErrorLabels holds error message strings for the asset revaluation drawer.
type AssetRevaluationErrorLabels struct {
	UseCaseUnavailable string `json:"useCaseUnavailable"`
	FormParseFailed    string `json:"formParseFailed"`
	RevaluateFailed    string `json:"revaluateFailed"`
	InvalidAmount      string `json:"invalidAmount"`
	// 2026-05-14 permission-gates: AWS-style permission tooltip surface for the
	// revaluation drawer. Emitted when the caller lacks asset_revaluation:*.
	PermissionDenied string `json:"permissionDenied"`
}

// AssetRevaluationLabels holds all translatable strings for the Asset Revaluation drawer.
type AssetRevaluationLabels struct {
	Title string `json:"title"`
	// Form field labels
	NewFairValue    string `json:"newFairValue"`
	AppraiserName   string `json:"appraiserName"`
	ValuationMethod string `json:"valuationMethod"`
	Notes           string `json:"notes"`
	// Preview section labels
	PreviewTitle      string `json:"previewTitle"`
	RevaluationAmount string `json:"revaluationAmount"`
	PnlSplit          string `json:"pnlSplit"`
	OciSplit          string `json:"ociSplit"`
	// Submit and cancel
	SubmitLabel string `json:"submitLabel"`
	CancelLabel string `json:"cancelLabel"`
	// Toast message template (supports {{.Direction}}/{{.Amount}}/{{.Recognition}} placeholders)
	ToastSuccessTemplate string `json:"toastSuccessTemplate"`
	// Errors
	Errors AssetRevaluationErrorLabels `json:"errors"`
}

// DefaultAssetRevaluationLabels returns AssetRevaluationLabels with sensible English defaults.
func DefaultAssetRevaluationLabels() AssetRevaluationLabels {
	return AssetRevaluationLabels{
		Title:                "Revalue Asset",
		NewFairValue:         "New fair value",
		AppraiserName:        "Appraiser name",
		ValuationMethod:      "Valuation method",
		Notes:                "Notes",
		PreviewTitle:         "Revaluation preview",
		RevaluationAmount:    "Revaluation amount",
		PnlSplit:             "Recognized in P&L",
		OciSplit:             "Recognized in OCI",
		SubmitLabel:          "Revalue",
		CancelLabel:          "Cancel",
		ToastSuccessTemplate: "Asset revalued: {{.Direction}}{{.Amount}} recognized in {{.Recognition}}",
		Errors: AssetRevaluationErrorLabels{
			UseCaseUnavailable: "Service unavailable. Please try again.",
			FormParseFailed:    "Form data could not be read.",
			RevaluateFailed:    "Failed to record the revaluation.",
			InvalidAmount:      "Invalid amount format. Use a number with up to 2 decimal places.",
			PermissionDenied:   "You do not have permission to revalue this asset.",
		},
	}
}

// ---------------------------------------------------------------------------
// Depreciation Policies labels (Surface F)
// Lyngua root key: "depreciationPolicies"
// ---------------------------------------------------------------------------

// DepreciationPoliciesLabels holds all translatable strings for the actionable
// Depreciation Policies page (Surface F).
type DepreciationPoliciesLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Table column headers
	Columns DepreciationPoliciesColumnLabels `json:"columns"`
	// Row action labels
	ActionPreview string `json:"actionPreview"`
	ActionRun     string `json:"actionRun"`
	// Preview drawer subtitle (read-only candidate preview, no DB writes)
	PreviewSubtitle string `json:"previewSubtitle"`
	// Run drawer subtitle (opens Surface C drawer with policy breadcrumb)
	RunSubtitle string `json:"runSubtitle"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
	// UsefulLifeMonthsSuffix is appended after the numeric useful_life_months value
	// in the table (e.g. "60" + " mo" → "60 mo"). Leading space is intentional.
	UsefulLifeMonthsSuffix string `json:"usefulLifeMonthsSuffix"`
}

// DepreciationPoliciesColumnLabels holds column headers for the Surface F policies table.
type DepreciationPoliciesColumnLabels struct {
	Policy          string `json:"policy"`
	Method          string `json:"method"`
	UsefulLife      string `json:"usefulLife"`
	SalvagePct      string `json:"salvagePct"`
	AssetsInPolicy  string `json:"assetsInPolicy"`
	AssetsDeviating string `json:"assetsDeviating"`
	Actions         string `json:"actions"`
}

// DefaultDepreciationPoliciesLabels returns DepreciationPoliciesLabels with sensible English defaults.
func DefaultDepreciationPoliciesLabels() DepreciationPoliciesLabels {
	return DepreciationPoliciesLabels{
		Title:    "Depreciation Policies",
		Subtitle: "Manage depreciation policies across asset categories",
		Columns: DepreciationPoliciesColumnLabels{
			Policy:          "Policy",
			Method:          "Method",
			UsefulLife:      "Useful Life",
			SalvagePct:      "Salvage %",
			AssetsInPolicy:  "Assets in policy",
			AssetsDeviating: "Assets deviating",
			Actions:         "Actions",
		},
		ActionPreview:          "Preview",
		ActionRun:              "Run",
		PreviewSubtitle:        "Preview depreciation amounts for this policy (no changes will be posted)",
		RunSubtitle:            "Post depreciation for all in-service assets under this policy",
		EmptyTitle:             "No depreciation policies",
		EmptyMessage:           "Add an asset category to define a depreciation policy.",
		UsefulLifeMonthsSuffix: " mo",
	}
}

// ---------------------------------------------------------------------------
// TaxRateLabels
// Lyngua root key: "taxRate"
// ---------------------------------------------------------------------------

// TaxRateLabels holds all translatable strings for the Tax Rate read-only views.
// Tax rates are read-only in the UI; supersession is via admin SQL recipe.
type TaxRateLabels struct {
	Page    TaxRatePageLabels   `json:"page"`
	Columns TaxRateColumnLabels `json:"columns"`
	Actions TaxRateActionLabels `json:"actions"`
	Empty   TaxRateEmptyLabels  `json:"empty"`
	Detail  TaxRateDetailLabels `json:"detail"`
}

// TaxRatePageLabels holds page heading strings.
type TaxRatePageLabels struct {
	HeadingActive     string `json:"headingActive"`
	CaptionActive     string `json:"captionActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	CaptionSuperseded string `json:"captionSuperseded"`
}

// TaxRateColumnLabels holds table column headers.
type TaxRateColumnLabels struct {
	Jurisdiction  string `json:"jurisdiction"`
	AuthorityCode string `json:"authorityCode"`
	Kind          string `json:"kind"`
	TreatmentCode string `json:"treatmentCode"`
	Direction     string `json:"direction"`
	RateBps       string `json:"rateBps"`
	EffectiveFrom string `json:"effectiveFrom"`
	Status        string `json:"status"`
}

// TaxRateActionLabels holds action button labels.
type TaxRateActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// TaxRateEmptyLabels holds empty-state strings.
type TaxRateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// TaxRateDetailLabels holds detail-page field labels.
type TaxRateDetailLabels struct {
	Title         string `json:"title"`
	Jurisdiction  string `json:"jurisdiction"`
	AuthorityCode string `json:"authorityCode"`
	Kind          string `json:"kind"`
	TreatmentCode string `json:"treatmentCode"`
	Direction     string `json:"direction"`
	RateBps       string `json:"rateBps"`
	EffectiveFrom string `json:"effectiveFrom"`
	EffectiveTo   string `json:"effectiveTo"`
	Status        string `json:"status"`
	SupersedesID  string `json:"supersedesId"`
	WorkspaceID   string `json:"workspaceId"`
}

// DefaultTaxRateLabels returns TaxRateLabels with sensible English defaults.
func DefaultTaxRateLabels() TaxRateLabels {
	return TaxRateLabels{
		Page: TaxRatePageLabels{
			HeadingActive:     "Tax Rates",
			CaptionActive:     "Active tax rates applied during computation",
			HeadingSuperseded: "Superseded Tax Rates",
			CaptionSuperseded: "Historical tax rates (superseded by newer versions)",
		},
		Columns: TaxRateColumnLabels{
			Jurisdiction:  "Jurisdiction",
			AuthorityCode: "Authority",
			Kind:          "Kind",
			TreatmentCode: "Treatment",
			Direction:     "Direction",
			RateBps:       "Rate (bps)",
			EffectiveFrom: "Effective From",
			Status:        "Status",
		},
		Actions: TaxRateActionLabels{
			View:         "View",
			NoPermission: "You do not have permission to view tax rates",
		},
		Empty: TaxRateEmptyLabels{
			Title:   "No tax rates found",
			Message: "Tax rates are added via the seed CSVs and superseded via admin SQL recipe.",
		},
		Detail: TaxRateDetailLabels{
			Title:         "Tax Rate Detail",
			Jurisdiction:  "Jurisdiction",
			AuthorityCode: "Authority Code",
			Kind:          "Kind",
			TreatmentCode: "Treatment",
			Direction:     "Direction",
			RateBps:       "Rate (basis points)",
			EffectiveFrom: "Effective From",
			EffectiveTo:   "Effective To",
			Status:        "Status",
			SupersedesID:  "Supersedes",
			WorkspaceID:   "Workspace",
		},
	}
}

// ---------------------------------------------------------------------------
// ForexRateLabels
// Lyngua root key: "forexRate"
// ---------------------------------------------------------------------------

// ForexRateLabels holds all translatable strings for the Forex Rate read-only views.
// Forex rates are read-only in the UI; rows are appended only via RecordOperatorRate.
type ForexRateLabels struct {
	Page    ForexRatePageLabels   `json:"page"`
	Columns ForexRateColumnLabels `json:"columns"`
	Actions ForexRateActionLabels `json:"actions"`
	Empty   ForexRateEmptyLabels  `json:"empty"`
	Detail  ForexRateDetailLabels `json:"detail"`
}

// ForexRatePageLabels holds page heading strings.
type ForexRatePageLabels struct {
	HeadingActive     string `json:"headingActive"`
	CaptionActive     string `json:"captionActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	CaptionSuperseded string `json:"captionSuperseded"`
}

// ForexRateColumnLabels holds table column headers.
type ForexRateColumnLabels struct {
	FromCurrency   string `json:"fromCurrency"`
	ToCurrency     string `json:"toCurrency"`
	RateMicroUnits string `json:"rateMicroUnits"`
	Source         string `json:"source"`
	EffectiveFrom  string `json:"effectiveFrom"`
	Status         string `json:"status"`
}

// ForexRateActionLabels holds action button labels.
type ForexRateActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// ForexRateEmptyLabels holds empty-state strings.
type ForexRateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// ForexRateDetailLabels holds detail-page field labels.
type ForexRateDetailLabels struct {
	Title           string `json:"title"`
	FromCurrency    string `json:"fromCurrency"`
	ToCurrency      string `json:"toCurrency"`
	RateMicroUnits  string `json:"rateMicroUnits"`
	Source          string `json:"source"`
	EffectiveFrom   string `json:"effectiveFrom"`
	EffectiveTo     string `json:"effectiveTo"`
	Status          string `json:"status"`
	SupersedesID    string `json:"supersedesId"`
	WorkspaceID     string `json:"workspaceId"`
	CreatedByUserID string `json:"createdByUserId"`
	Notes           string `json:"notes"`
}

// DefaultForexRateLabels returns ForexRateLabels with sensible English defaults.
func DefaultForexRateLabels() ForexRateLabels {
	return ForexRateLabels{
		Page: ForexRatePageLabels{
			HeadingActive:     "Forex Rates",
			CaptionActive:     "Active foreign exchange rates used during billing",
			HeadingSuperseded: "Superseded Forex Rates",
			CaptionSuperseded: "Historical foreign exchange rates (superseded by newer versions)",
		},
		Columns: ForexRateColumnLabels{
			FromCurrency:   "From",
			ToCurrency:     "To",
			RateMicroUnits: "Rate (micro-units)",
			Source:         "Source",
			EffectiveFrom:  "Effective From",
			Status:         "Status",
		},
		Actions: ForexRateActionLabels{
			View:         "View",
			NoPermission: "You do not have permission to view forex rates",
		},
		Empty: ForexRateEmptyLabels{
			Title:   "No forex rates found",
			Message: "Forex rates are recorded automatically when revenue is recognized with a foreign currency.",
		},
		Detail: ForexRateDetailLabels{
			Title:           "Forex Rate Detail",
			FromCurrency:    "From Currency",
			ToCurrency:      "To Currency",
			RateMicroUnits:  "Rate (micro-units)",
			Source:          "Source",
			EffectiveFrom:   "Effective From",
			EffectiveTo:     "Effective To",
			Status:          "Status",
			SupersedesID:    "Supersedes",
			WorkspaceID:     "Workspace",
			CreatedByUserID: "Created By",
			Notes:           "Notes",
		},
	}
}

// ---------------------------------------------------------------------------
// WithholdingCertificateLabels
// Lyngua root key: "withholdingCertificate"
// ---------------------------------------------------------------------------

// WithholdingCertificateLabels holds all translatable strings for the
// Withholding Certificate CRUD views.
type WithholdingCertificateLabels struct {
	Page    WithholdingCertificatePageLabels   `json:"page"`
	Columns WithholdingCertificateColumnLabels `json:"columns"`
	Buttons WithholdingCertificateButtonLabels `json:"buttons"`
	Actions WithholdingCertificateActionLabels `json:"actions"`
	Empty   WithholdingCertificateEmptyLabels  `json:"empty"`
	Fields  WithholdingCertificateFieldLabels  `json:"fields"`
}

// WithholdingCertificatePageLabels holds page heading strings.
type WithholdingCertificatePageLabels struct {
	HeadingActive string `json:"headingActive"`
	CaptionActive string `json:"captionActive"`
	HeadingVoided string `json:"headingVoided"`
	CaptionVoided string `json:"captionVoided"`
}

// WithholdingCertificateColumnLabels holds table column headers.
type WithholdingCertificateColumnLabels struct {
	CertificateNumber  string `json:"certificateNumber"`
	RevenueID          string `json:"revenueId"`
	PeriodYear         string `json:"periodYear"`
	PeriodQuarter      string `json:"periodQuarter"`
	WhtAmountCertified string `json:"whtAmountCertified"`
	Status             string `json:"status"`
	DateIssued         string `json:"dateIssued"`
}

// WithholdingCertificateButtonLabels holds button text.
type WithholdingCertificateButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
	Void   string `json:"void"`
}

// WithholdingCertificateActionLabels holds action dropdown labels.
type WithholdingCertificateActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	NoPermission string `json:"noPermission"`
}

// WithholdingCertificateEmptyLabels holds empty-state strings.
type WithholdingCertificateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// WithholdingCertificateFieldLabels holds drawer form field labels.
type WithholdingCertificateFieldLabels struct {
	CertificateNumber  string `json:"certificateNumber"`
	RevenueID          string `json:"revenueId"`
	TaxAuthorityID     string `json:"taxAuthorityId"`
	PayorTin           string `json:"payorTin"`
	PayorName          string `json:"payorName"`
	PayeeTin           string `json:"payeeTin"`
	PayeeName          string `json:"payeeName"`
	PeriodYear         string `json:"periodYear"`
	PeriodQuarter      string `json:"periodQuarter"`
	WhtAmountCertified string `json:"whtAmountCertified"`
	Status             string `json:"status"`
	DateIssued         string `json:"dateIssued"`
	Notes              string `json:"notes"`
}

// DefaultWithholdingCertificateLabels returns WithholdingCertificateLabels with
// sensible English defaults.
func DefaultWithholdingCertificateLabels() WithholdingCertificateLabels {
	return WithholdingCertificateLabels{
		Page: WithholdingCertificatePageLabels{
			HeadingActive: "Withholding Certificates",
			CaptionActive: "BIR Form 2307 withholding tax certificates",
			HeadingVoided: "Voided Certificates",
			CaptionVoided: "Voided withholding tax certificates",
		},
		Columns: WithholdingCertificateColumnLabels{
			CertificateNumber:  "Certificate No.",
			RevenueID:          "Invoice",
			PeriodYear:         "Period Year",
			PeriodQuarter:      "Quarter",
			WhtAmountCertified: "WHT Certified",
			Status:             "Status",
			DateIssued:         "Date Issued",
		},
		Buttons: WithholdingCertificateButtonLabels{
			Add:    "Add Certificate",
			Edit:   "Edit",
			Delete: "Delete",
			Void:   "Void",
		},
		Actions: WithholdingCertificateActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			NoPermission: "You do not have permission to manage withholding certificates",
		},
		Empty: WithholdingCertificateEmptyLabels{
			Title:   "No withholding certificates",
			Message: "Add a withholding certificate received from a customer to record creditable WHT.",
		},
		Fields: WithholdingCertificateFieldLabels{
			CertificateNumber:  "Certificate Number",
			RevenueID:          "Invoice",
			TaxAuthorityID:     "Tax Authority",
			PayorTin:           "Payor TIN",
			PayorName:          "Payor Name",
			PayeeTin:           "Payee TIN",
			PayeeName:          "Payee Name",
			PeriodYear:         "Period Year",
			PeriodQuarter:      "Quarter",
			WhtAmountCertified: "WHT Amount Certified",
			Status:             "Status",
			DateIssued:         "Date Issued",
			Notes:              "Notes",
		},
	}
}

// ---------------------------------------------------------------------------
// FundingFormLabels (Funding > drawer forms)
// Lyngua root key: "funding"
// ---------------------------------------------------------------------------

// FundingFormLabels holds all translatable strings for the four funding
// drawer forms: allocation, draw (charge), settlement, and transfer.
// The lyngua root key is "funding"; JSON lives in
// packages/lyngua/translations/en/{general,professional}/funding.json.
type FundingFormLabels struct {
	Allocation FundingAllocationFormLabels `json:"allocation"`
	Draw       FundingDrawFormLabels       `json:"draw"`
	Settlement FundingSettlementFormLabels `json:"settlement"`
	Transfer   FundingTransferFormLabels   `json:"transfer"`
	Source     FundingSourceListLabels     `json:"source"`
}

// FundingAllocationFormLabels holds field/button labels for the allocation drawer.
type FundingAllocationFormLabels struct {
	AllocatedLimit string `json:"allocatedLimit"`
	Mode           string `json:"mode"`
	ModeHardLimit  string `json:"modeHardLimit"`
	ModeSoftLimit  string `json:"modeSoftLimit"`
}

// FundingDrawFormLabels holds field/button labels for the draw (charge) drawer.
type FundingDrawFormLabels struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Submit      string `json:"submit"`
}

// FundingSettlementFormLabels holds field/button labels for the settlement drawer.
type FundingSettlementFormLabels struct {
	Amount string `json:"amount"`
	Submit string `json:"submit"`
}

// FundingTransferFormLabels holds field/button labels for the transfer drawer.
type FundingTransferFormLabels struct {
	DestinationFundID string `json:"destinationFundId"`
	Amount            string `json:"amount"`
	Submit            string `json:"submit"`
}

// FundingSourceListLabels holds page-level strings for the fund source list view.
type FundingSourceListLabels struct {
	Title    string                    `json:"title"`
	Subtitle string                    `json:"subtitle"`
	Kind     FundingSourceKindLabels   `json:"kind"`
	Status   FundingSourceStatusLabels `json:"status"`
}

// FundingSourceKindLabels maps FundKind enum values to display strings for the
// fund source list. Keys mirror the proto FundKind enum.
type FundingSourceKindLabels struct {
	CashOnHand  string `json:"cashOnHand"`
	BankAccount string `json:"bankAccount"`
	PettyCash   string `json:"pettyCash"`
	CreditCard  string `json:"creditCard"`
	CreditLine  string `json:"creditLine"`
	PrepaidCard string `json:"prepaidCard"`
	MobileMoney string `json:"mobileMoney"`
	Unknown     string `json:"unknown"`
}

// FundingSourceStatusLabels maps FundStatus enum values to display strings for
// the fund source list. Keys mirror the proto FundStatus enum.
type FundingSourceStatusLabels struct {
	Draft     string `json:"draft"`
	Active    string `json:"active"`
	Suspended string `json:"suspended"`
	Archived  string `json:"archived"`
	Unknown   string `json:"unknown"`
}

// DefaultFundingFormLabels returns FundingFormLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files
// (packages/lyngua/translations/en/{general,professional}/funding.json).
func DefaultFundingFormLabels() FundingFormLabels {
	return FundingFormLabels{
		Allocation: FundingAllocationFormLabels{
			AllocatedLimit: "Allocated Limit",
			Mode:           "Mode",
			ModeHardLimit:  "Hard Limit",
			ModeSoftLimit:  "Soft Limit",
		},
		Draw: FundingDrawFormLabels{
			Amount:      "Amount",
			Description: "Description",
			Submit:      "Charge",
		},
		Settlement: FundingSettlementFormLabels{
			Amount: "Settlement Amount",
			Submit: "Settle",
		},
		Transfer: FundingTransferFormLabels{
			DestinationFundID: "Destination Fund ID",
			Amount:            "Amount",
			Submit:            "Transfer",
		},
		Source: FundingSourceListLabels{
			Title:    "Fund Sources",
			Subtitle: "Funds you own and share with workspaces",
			Kind: FundingSourceKindLabels{
				CashOnHand:  "Cash on Hand",
				BankAccount: "Bank Account",
				PettyCash:   "Petty Cash",
				CreditCard:  "Credit Card",
				CreditLine:  "Credit Line",
				PrepaidCard: "Prepaid Card",
				MobileMoney: "Mobile Money",
				Unknown:     "Unknown",
			},
			Status: FundingSourceStatusLabels{
				Draft:     "Draft",
				Active:    "Active",
				Suspended: "Suspended",
				Archived:  "Archived",
				Unknown:   "Unknown",
			},
		},
	}
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
