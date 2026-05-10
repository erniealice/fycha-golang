package form

// Record is a flat struct for passing asset data between action handlers
// and the DB layer. It avoids a dependency on protobuf types.
type Record struct {
	ID                 string
	AssetNumber        string
	Name               string
	Description        string
	AssetType          string
	AssetCategoryID    string
	LocationID         string
	AcquisitionCost    float64
	SalvageValue       float64
	BookValue          float64
	UsefulLifeMonths   int
	DepreciationMethod string
	Currency           string
	Status             string
	Active             bool
}

// Labels holds i18n labels for the drawer form template.
type Labels struct {
	Name                       string
	NamePlaceholder            string
	AssetNumber                string
	AssetNumberPlaceholder     string
	Description                string
	DescriptionPlaceholder     string
	Category                   string
	CategoryPlaceholder        string
	Location                   string
	LocationPlaceholder        string
	AcquisitionCost            string
	AcquisitionCostPlaceholder string
	SalvageValue               string
	SalvageValuePlaceholder    string
	UsefulLifeMonths           string
	UsefulLifePlaceholder      string
	DepreciationMethod         string
	DepMethodStraightLine      string
	DepMethodDecliningBalance  string
	DepMethodSumOfYears        string
	DepMethodUnitsOfProduction string
	Active                     string
	// Info popover text
	AssetNumberInfo        string
	AcquisitionCostInfo    string
	SalvageValueInfo       string
	UsefulLifeMonthsInfo   string
	DepreciationMethodInfo string
	// UoP-specific tooltip rendered on the disabled Units of Production option.
	UnitsOfProductionDisabledTooltip string
}

// Data is the template data for the asset drawer form.
type Data struct {
	FormAction         string
	IsEdit             bool
	ID                 string
	Name               string
	AssetNumber        string
	Description        string
	CategoryID         string
	LocationID         string
	AcquisitionCost    string
	SalvageValue       string
	UsefulLifeMonths   string
	DepreciationMethod string
	Active             bool
	// DepreciationFieldsLocked is true when posted depreciation_schedule rows
	// exist for this asset, causing useful_life_months / depreciation_method /
	// salvage_value / depreciation_start_date to render as read-only with a warning.
	DepreciationFieldsLocked bool
	Labels                   Labels
	CommonLabels             any
}
