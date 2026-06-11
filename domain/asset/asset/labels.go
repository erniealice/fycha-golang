package asset

// ---------------------------------------------------------------------------
// Asset labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the fixed asset module.
type Labels struct {
	Page      PageLabels      `json:"page"`
	Buttons   ButtonLabels    `json:"buttons"`
	Columns   ColumnLabels    `json:"columns"`
	Empty     EmptyLabels     `json:"empty"`
	Form      FormLabels      `json:"form"`
	Actions   ActionLabels    `json:"actions"`
	Detail    DetailLabels    `json:"detail"`
	Dashboard DashboardLabels `json:"dashboard"`
}

type PageLabels struct {
	Heading         string `json:"heading"`
	HeadingActive   string `json:"headingActive"`
	HeadingInactive string `json:"headingInactive"`
	Caption         string `json:"caption"`
	CaptionActive   string `json:"captionActive"`
	CaptionInactive string `json:"captionInactive"`
}

type ButtonLabels struct {
	AddAsset string `json:"addAsset"`
}

type ColumnLabels struct {
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

type EmptyLabels struct {
	ActiveTitle     string `json:"activeTitle"`
	ActiveMessage   string `json:"activeMessage"`
	InactiveTitle   string `json:"inactiveTitle"`
	InactiveMessage string `json:"inactiveMessage"`
}

type FormLabels struct {
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

type ActionLabels struct {
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

type DetailLabels struct {
	BasicInfo        DetailBasicInfoLabels `json:"basicInfo"`
	Tabs             DetailTabLabels       `json:"tabs"`
	EmptyStates      DetailEmptyLabels     `json:"emptyStates"`
	AttachmentUpload string                `json:"attachmentUpload"`
}

type DetailBasicInfoLabels struct {
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

type DetailTabLabels struct {
	Info                  string `json:"info"`
	LapsingActualSchedule string `json:"lapsingActualSchedule"`
	TransactionLedger     string `json:"transactionLedger"`
	Attachments           string `json:"attachments"`
	History               string `json:"history"`
}

type DetailEmptyLabels struct {
	DepreciationTitle string `json:"depreciationTitle"`
	DepreciationDesc  string `json:"depreciationDesc"`
	MaintenanceTitle  string `json:"maintenanceTitle"`
	MaintenanceDesc   string `json:"maintenanceDesc"`
	TransactionsTitle string `json:"transactionsTitle"`
	TransactionsDesc  string `json:"transactionsDesc"`
}

type DashboardLabels struct {
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

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading:         "Fixed Assets",
			HeadingActive:   "Active Assets",
			HeadingInactive: "Inactive Assets",
			Caption:         "Manage your fixed assets",
			CaptionActive:   "Active fixed assets in your register",
			CaptionInactive: "Inactive or disposed fixed assets",
		},
		Buttons: ButtonLabels{
			AddAsset: "Add Asset",
		},
		Columns: ColumnLabels{
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
		Empty: EmptyLabels{
			ActiveTitle:     "No active assets",
			ActiveMessage:   "Add your first fixed asset to start tracking depreciation and maintenance.",
			InactiveTitle:   "No inactive assets",
			InactiveMessage: "Deactivated or disposed assets will appear here.",
		},
		Form: FormLabels{
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
		Actions: ActionLabels{
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
		Detail: DetailLabels{
			BasicInfo: DetailBasicInfoLabels{
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
			Tabs: DetailTabLabels{
				Info:                  "Info",
				LapsingActualSchedule: "Lapsing Schedule",
				TransactionLedger:     "Transaction Ledger",
				Attachments:           "Attachments",
				History:               "History",
			},
			EmptyStates: DetailEmptyLabels{
				DepreciationTitle: "No depreciation records",
				DepreciationDesc:  "Depreciation schedule will appear here once configured.",
				MaintenanceTitle:  "No maintenance records",
				MaintenanceDesc:   "Maintenance history for this asset will appear here.",
				TransactionsTitle: "No transactions",
				TransactionsDesc:  "Transaction audit trail for this asset will appear here.",
			},
			AttachmentUpload: "Upload Attachment",
		},
		Dashboard: DashboardLabels{
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
