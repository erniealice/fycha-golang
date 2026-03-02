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
		Export:              common.Table.Export,
		DensityDefault:     common.Table.Density.Default,
		DensityComfortable: common.Table.Density.Comfortable,
		DensityCompact:     common.Table.Density.Compact,
		Show:               common.Table.Show,
		Entries:             common.Table.Entries,
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

// ---------------------------------------------------------------------------
// Asset labels
// ---------------------------------------------------------------------------

// AssetLabels holds all translatable strings for the fixed asset module.
type AssetLabels struct {
	Page    AssetPageLabels    `json:"page"`
	Buttons AssetButtonLabels  `json:"buttons"`
	Columns AssetColumnLabels  `json:"columns"`
	Empty   AssetEmptyLabels   `json:"empty"`
	Form    AssetFormLabels    `json:"form"`
	Actions AssetActionLabels  `json:"actions"`
	Detail  AssetDetailLabels  `json:"detail"`
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
}

type AssetActionLabels struct {
	View       string `json:"view"`
	Edit       string `json:"edit"`
	Delete     string `json:"delete"`
	Activate   string `json:"activate"`
	Deactivate string `json:"deactivate"`
}

type AssetDetailLabels struct {
	BasicInfo   AssetDetailBasicInfoLabels `json:"basicInfo"`
	Tabs        AssetDetailTabLabels       `json:"tabs"`
	EmptyStates AssetDetailEmptyLabels     `json:"emptyStates"`
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
	Info          string `json:"info"`
	Depreciation  string `json:"depreciation"`
	Maintenance   string `json:"maintenance"`
	Transactions  string `json:"transactions"`
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
		},
		Actions: AssetActionLabels{
			View:       "View",
			Edit:       "Edit",
			Delete:     "Delete",
			Activate:   "Activate",
			Deactivate: "Deactivate",
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
			},
			EmptyStates: AssetDetailEmptyLabels{
				DepreciationTitle: "No depreciation records",
				DepreciationDesc:  "Depreciation schedule will appear here once configured.",
				MaintenanceTitle:  "No maintenance records",
				MaintenanceDesc:   "Maintenance history for this asset will appear here.",
				TransactionsTitle: "No transactions",
				TransactionsDesc:  "Transaction audit trail for this asset will appear here.",
			},
		},
		Dashboard: AssetDashboardLabels{
			Title:            "Assets Dashboard",
			Subtitle:         "Overview of your fixed asset register",
			TotalAssets:      "Total Assets",
			TotalBookValue:   "Total Book Value",
			FullyDepreciated: "Fully Depreciated",
			UnderMaintenance: "Under Maintenance",
		},
	}
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
