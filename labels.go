package fycha

import (
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
)

// MapTableLabels maps common labels into the flat types.TableLabels structure.
func MapTableLabels(common pyeza.CommonLabels) types.TableLabels {
	return types.TableLabels{
		Search:             common.Table.Search,
		SearchPlaceholder:  common.Table.SearchPlaceholder,
		Filters:            common.Table.Filters,
		FilterConditions:   common.Table.FilterConditions,
		ClearAll:           common.Table.ClearAll,
		AddCondition:       common.Table.AddCondition,
		Clear:              common.Table.Clear,
		ApplyFilters:       common.Table.ApplyFilters,
		Sort:               common.Table.Sort,
		Columns:            common.Table.Columns,
		Export:             common.Table.Export,
		DensityLabel:       common.Table.Density.Title,
		DensityDefault:     common.Table.Density.Default,
		DensityComfortable: common.Table.Density.Comfortable,
		DensityCompact:     common.Table.Density.Compact,
		EntriesPerPage:     common.Table.EntriesLabel,
		Show:               common.Table.Show,
		Entries:            common.Table.Entries,
		Showing:            common.Table.Showing,
		To:                 common.Table.To,
		Of:                 common.Table.Of,
		EntriesLabel:       common.Table.EntriesLabel,
		SelectAll:          common.Table.SelectAll,
		Actions:            common.Table.Actions,
		Prev:               common.Pagination.Prev,
		Next:               common.Pagination.Next,
	}
}

// MapBulkConfig returns a BulkActionsConfig with labels from common bulk labels.
func MapBulkConfig(common pyeza.CommonLabels) types.BulkActionsConfig {
	return types.BulkActionsConfig{
		Enabled:        true,
		SelectAllLabel: common.Bulk.SelectAll,
		SelectedLabel:  common.Bulk.Selected,
		CancelLabel:    common.Bulk.ClearSelection,
	}
}

// ReportsLabels holds all translatable strings for the reports module.
type ReportsLabels struct {
	GrossProfit     GrossProfitLabels     `json:"grossProfit"`
	Revenue         RevenueLabels         `json:"revenue"`
	RevenueReport   RevenueReportLabels   `json:"revenueReport"`
	CostOfSales     CostOfSalesLabels     `json:"costOfSales"`
	Expenses        ExpensesLabels        `json:"expenses"`
	NetProfit       NetProfitLabels       `json:"netProfit"`
	Dashboard       DashboardLabels       `json:"dashboard"`
	Period          PeriodLabels          `json:"period"`
	IncomeStatement IncomeStatementLabels `json:"incomeStatement"`
	BalanceSheet    BalanceSheetLabels    `json:"balanceSheet"`
	CashFlow        CashFlowLabels        `json:"cashFlow"`
	EquityChanges   EquityChangesLabels   `json:"equityChanges"`
}

// IncomeStatementLabels holds translatable strings for the Income Statement page.
type IncomeStatementLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

// BalanceSheetLabels holds translatable strings for the Balance Sheet page.
type BalanceSheetLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

// CashFlowLabels holds translatable strings for the Cash Flow Statement page.
type CashFlowLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
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
}

type AssetActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
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
	Info         string `json:"info"`
	Depreciation string `json:"depreciation"`
	Maintenance  string `json:"maintenance"`
	Transactions string `json:"transactions"`
	Attachments  string `json:"attachments"`
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
			Name:                       "Name",
			NamePlaceholder:            "e.g. Office Laptop",
			AssetNumber:                "Asset Number",
			AssetNumberPlaceholder:     "e.g. FA-001",
			Description:                "Description",
			DescriptionPlaceholder:     "Brief description of the asset",
			Category:                   "Category",
			CategoryPlaceholder:        "Select a category",
			Location:                   "Location",
			LocationPlaceholder:        "Select a location",
			AcquisitionCost:            "Acquisition Cost",
			AcquisitionCostPlaceholder: "0.00",
			SalvageValue:               "Salvage Value",
			SalvageValuePlaceholder:    "0.00",
			UsefulLifeMonths:           "Useful Life (Months)",
			UsefulLifePlaceholder:      "e.g. 60",
			DepreciationMethod:         "Depreciation Method",
			Active:                     "Active",
			DepMethodStraightLine:      "Straight Line",
			DepMethodDecliningBalance:  "Declining Balance",
			DepMethodSumOfYears:        "Sum of Years' Digits",
			DepMethodUnitsOfProduction: "Units of Production",
		},
		Actions: AssetActionLabels{
			View:                  "View",
			Edit:                  "Edit",
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
				Info:         "Info",
				Depreciation: "Depreciation",
				Maintenance:  "Maintenance",
				Transactions: "Transactions",
				Attachments:  "Attachments",
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
			Title:                "Assets Dashboard",
			Subtitle:             "Overview of your fixed asset register",
			TotalAssets:          "Total Assets",
			TotalBookValue:       "Total Book Value",
			FullyDepreciated:     "Fully Depreciated",
			UnderMaintenance:     "Under Maintenance",
			ActivityAcquired:     "New asset acquired",
			ActivityMaintenance:  "Maintenance completed",
			ActivityDepreciation: "Depreciation recorded",
			RecentActivity:       "Recent Activity",
			NoRecentActivity:     "No recent asset activity",
			UnknownAsset:         "Unknown Asset",
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
	}
}

// ---------------------------------------------------------------------------
// Journal labels (Journal Entries)
// ---------------------------------------------------------------------------

// JournalLabels holds all translatable strings for the Journal Entries module.
type JournalLabels struct {
	Page    JournalPageLabels    `json:"page"`
	Tabs    JournalTabLabels     `json:"tabs"`
	Buttons JournalButtonLabels  `json:"buttons"`
	Columns JournalColumnLabels  `json:"columns"`
	Empty   JournalEmptyLabels   `json:"empty"`
	Actions JournalActionLabels  `json:"actions"`
	Lines   JournalLineLabels    `json:"lines"`
	Form    JournalFormLabels    `json:"form"`
	Detail  JournalDetailLabels  `json:"detail"`
	Confirm JournalConfirmLabels `json:"confirm"`
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
			SourceLabel: "Source",
			ViewSource:  "View Source \u2192",
			Balanced:    "Balanced",
			Unbalanced:  "Unbalanced",
			Totals:      "TOTALS",
			Difference:  "DIFFERENCE",
		},
		Confirm: JournalConfirmLabels{
			Post:    "Are you sure you want to post this journal entry? This action cannot be undone.",
			Delete:  "Are you sure you want to delete this journal entry? This action cannot be undone.",
			Reverse: "Are you sure you want to reverse this journal entry? A reversing entry will be created.",
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
			Name:         "Name",
			PeriodNumber: "Period Number",
			FiscalYear:   "Fiscal Year",
			StartDate:    "Start Date",
			EndDate:      "End Date",
			Status:       "Status",
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
	}
}

// ---------------------------------------------------------------------------
// Loan labels (Funding > Loans app)
// ---------------------------------------------------------------------------

// LoanLabels is the top-level label container for the Loans app.
type LoanLabels struct {
	Page    LoanPageLabels   `json:"page"`
	Tabs    LoanTabLabels    `json:"tabs"`
	Buttons LoanButtonLabels `json:"buttons"`
	Columns LoanColumnLabels `json:"columns"`
	Empty   LoanEmptyLabels  `json:"empty"`
	Actions LoanActionLabels `json:"actions"`
	Form    LoanFormLabels   `json:"form"`
	Status  LoanStatusLabels `json:"status"`
	Type    LoanTypeLabels   `json:"type"`
	Sheet   LoanSheetLabels  `json:"sheet"`
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
	Title                 string `json:"title"`
	Subtitle              string `json:"subtitle"`
	ColumnDimension       string `json:"columnDimension"`
	RowDimension          string `json:"rowDimension"`
	DimensionMonthly      string `json:"dimensionMonthly"`
	DimensionQuarterly    string `json:"dimensionQuarterly"`
	DimensionYearly       string `json:"dimensionYearly"`
	DimensionProduct      string `json:"dimensionProduct"`
	DimensionProductLine  string `json:"dimensionProductLine"`
	DimensionLocation     string `json:"dimensionLocation"`
	DimensionLocationArea string `json:"dimensionLocationArea"`
	SummaryGrandTotal     string `json:"summaryGrandTotal"`
	SummaryTransactions   string `json:"summaryTransactions"`
	SummaryAverage        string `json:"summaryAverage"`
	Total                 string `json:"total"`
	Totals                string `json:"totals"`
	ExportCsv             string `json:"exportCsv"`
	Apply                 string `json:"apply"`
	Clear                 string `json:"clear"`
	EmptyTitle            string `json:"emptyTitle"`
	EmptyMessage          string `json:"emptyMessage"`
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
	}
	opts := make([]FilterOption, len(dims))
	for i, d := range dims {
		opts[i] = FilterOption{Value: d.value, Label: d.label, Selected: d.value == active}
	}
	return opts
}
