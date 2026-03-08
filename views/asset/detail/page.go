package detail

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Mock data types (intermediate -- converted to TableConfig for rendering)
// ---------------------------------------------------------------------------

// DepreciationRow holds one row of the depreciation schedule.
type DepreciationRow struct {
	Period      string
	StartValue  string
	Expense     string
	EndValue    string
	Accumulated string
}

// MaintenanceRow holds one row of the maintenance history.
type MaintenanceRow struct {
	Date        string
	Type        string
	Description string
	Status      string
	StatusClass string
	Cost        string
}

// TransactionRow holds one row of the transaction history.
type TransactionRow struct {
	Date        string
	Type        string
	Description string
	Amount      string
	Reference   string
}

// MockAssetDetail holds mock data for the detail page.
type MockAssetDetail struct {
	ID                    string
	AssetNumber           string
	Name                  string
	Description           string
	CategoryName          string
	LocationName          string
	AcquisitionCost       string
	AcquisitionCostRaw    string
	SalvageValue          string
	SalvageValueRaw       string
	UsefulLifeMonths      string
	DepreciationMethod    string
	DepreciationMethodKey string
	BookValue             string
	Status                string
	DepreciationSchedule  []DepreciationRow
	MaintenanceRecords    []MaintenanceRow
	TransactionHistory    []TransactionRow
}

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies.
type Deps struct {
	Routes       fycha.AssetRoutes
	Labels       fycha.AssetLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the asset detail page.
type PageData struct {
	types.PageData
	ContentTemplate       string
	Labels                fycha.AssetLabels
	ActiveTab             string
	TabItems              []pyeza.TabItem
	ID                    string
	AssetName             string
	AssetNumber           string
	AssetDescription      string
	CategoryName          string
	LocationName          string
	AcquisitionCost       string
	AcquisitionCostRaw    string
	SalvageValue          string
	SalvageValueRaw       string
	UsefulLifeMonths      string
	DepreciationMethod    string
	DepreciationMethodKey string
	BookValue             string
	AssetStatus           string
	StatusVariant         string
	EditURL               string
	CanEdit               bool
	DepreciationTable     *types.TableConfig
	MaintenanceTable      *types.TableConfig
	TransactionTable      *types.TableConfig
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the asset detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		perms := view.GetUserPermissions(ctx)
		pageData := buildPageData(deps, id, activeTab, viewCtx, perms)
		return view.OK("asset-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial -- returns only the tab content).
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		tab := viewCtx.Request.PathValue("tab")
		if tab == "" {
			tab = "info"
		}

		perms := view.GetUserPermissions(ctx)
		pageData := buildPageData(deps, id, tab, viewCtx, perms)

		// Return only the tab partial template
		templateName := "asset-tab-" + tab
		return view.OK(templateName, pageData)
	})
}

// ---------------------------------------------------------------------------
// Page data builder
// ---------------------------------------------------------------------------

func buildPageData(deps *Deps, id, activeTab string, viewCtx *view.ViewContext, perms *types.UserPermissions) *PageData {
	asset := getMockAsset(id)

	statusVariant := "success"
	if asset.Status == "inactive" {
		statusVariant = "warning"
	}

	tabItems := buildTabItems(id, deps.Labels, deps.Routes)

	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          asset.Name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "assets",
			ActiveSubNav:   "assets-fixed",
			HeaderTitle:    asset.Name,
			HeaderSubtitle: fmt.Sprintf("%s | %s", asset.AssetNumber, asset.CategoryName),
			HeaderIcon:     "icon-box",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:       "asset-detail-content",
		Labels:                deps.Labels,
		ActiveTab:             activeTab,
		TabItems:              tabItems,
		ID:                    id,
		AssetName:             asset.Name,
		AssetNumber:           asset.AssetNumber,
		AssetDescription:      asset.Description,
		CategoryName:          asset.CategoryName,
		LocationName:          asset.LocationName,
		AcquisitionCost:       asset.AcquisitionCost,
		AcquisitionCostRaw:    asset.AcquisitionCostRaw,
		SalvageValue:          asset.SalvageValue,
		SalvageValueRaw:       asset.SalvageValueRaw,
		UsefulLifeMonths:      asset.UsefulLifeMonths,
		DepreciationMethod:    asset.DepreciationMethod,
		DepreciationMethodKey: asset.DepreciationMethodKey,
		BookValue:             asset.BookValue,
		AssetStatus:           asset.Status,
		StatusVariant:         statusVariant,
		EditURL:               route.ResolveURL(deps.Routes.EditURL, "id", id),
		CanEdit:               perms.Can("asset", "update"),
		DepreciationTable:     buildDepreciationTable(asset.DepreciationSchedule, deps.Labels, deps.TableLabels),
		MaintenanceTable:      buildMaintenanceTable(asset.MaintenanceRecords, deps.Labels, deps.TableLabels),
		TransactionTable:      buildTransactionTable(asset.TransactionHistory, deps.Labels, deps.TableLabels),
	}
}

func buildTabItems(id string, labels fycha.AssetLabels, routes fycha.AssetRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: labels.Detail.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "depreciation", Label: labels.Detail.Tabs.Depreciation, Href: base + "?tab=depreciation", HxGet: action + "depreciation", Icon: "icon-trending-down", Count: 0, Disabled: false},
		{Key: "maintenance", Label: labels.Detail.Tabs.Maintenance, Href: base + "?tab=maintenance", HxGet: action + "maintenance", Icon: "icon-tool", Count: 0, Disabled: false},
		{Key: "transactions", Label: labels.Detail.Tabs.Transactions, Href: base + "?tab=transactions", HxGet: action + "transactions", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Table builders -- convert mock data into pyeza TableConfig
// ---------------------------------------------------------------------------

func buildDepreciationTable(schedule []DepreciationRow, labels fycha.AssetLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "period", Label: labels.Columns.Period, Sortable: true},
		{Key: "start_value", Label: labels.Columns.StartValue, Align: "right"},
		{Key: "depreciation", Label: labels.Columns.Depreciation, Align: "right"},
		{Key: "end_value", Label: labels.Columns.EndValue, Align: "right"},
		{Key: "accumulated", Label: labels.Columns.Accumulated, Align: "right"},
	}

	rows := make([]types.TableRow, len(schedule))
	for i, s := range schedule {
		rows[i] = types.TableRow{
			ID: fmt.Sprintf("dep-%d", i+1),
			Cells: []types.TableCell{
				{Type: "text", Value: s.Period},
				{Type: "text", Value: s.StartValue},
				{Type: "text", Value: s.Expense},
				{Type: "text", Value: s.EndValue},
				{Type: "text", Value: s.Accumulated},
			},
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:      "depreciation-table",
		Minimal: true,
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   labels.Detail.EmptyStates.DepreciationTitle,
			Message: labels.Detail.EmptyStates.DepreciationDesc,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

func buildMaintenanceTable(records []MaintenanceRow, labels fycha.AssetLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: labels.Columns.Date, Sortable: true},
		{Key: "type", Label: labels.Columns.Type, Sortable: true},
		{Key: "description", Label: labels.Columns.Description},
		{Key: "status", Label: labels.Columns.Status, Width: "120px"},
		{Key: "cost", Label: labels.Columns.Cost, Align: "right"},
	}

	rows := make([]types.TableRow, len(records))
	for i, m := range records {
		rows[i] = types.TableRow{
			ID: fmt.Sprintf("mnt-%d", i+1),
			Cells: []types.TableCell{
				{Type: "text", Value: m.Date},
				{Type: "text", Value: m.Type},
				{Type: "text", Value: m.Description},
				{Type: "badge", Value: m.Status, Variant: m.StatusClass},
				{Type: "text", Value: m.Cost},
			},
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:      "maintenance-table",
		Minimal: true,
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   labels.Detail.EmptyStates.MaintenanceTitle,
			Message: labels.Detail.EmptyStates.MaintenanceDesc,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

func buildTransactionTable(history []TransactionRow, labels fycha.AssetLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: labels.Columns.Date, Sortable: true},
		{Key: "type", Label: labels.Columns.Type, Sortable: true},
		{Key: "description", Label: labels.Columns.Description},
		{Key: "amount", Label: labels.Columns.Amount, Align: "right"},
		{Key: "reference", Label: labels.Columns.Reference},
	}

	rows := make([]types.TableRow, len(history))
	for i, t := range history {
		rows[i] = types.TableRow{
			ID: fmt.Sprintf("txn-%d", i+1),
			Cells: []types.TableCell{
				{Type: "text", Value: t.Date},
				{Type: "text", Value: t.Type},
				{Type: "text", Value: t.Description},
				{Type: "text", Value: t.Amount},
				{Type: "text", Value: t.Reference},
			},
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:      "transactions-table",
		Minimal: true,
		Columns: columns,
		Rows:    rows,
		Labels:  tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   labels.Detail.EmptyStates.TransactionsTitle,
			Message: labels.Detail.EmptyStates.TransactionsDesc,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

// ---------------------------------------------------------------------------
// Mock data
// ---------------------------------------------------------------------------

func getMockAsset(id string) MockAssetDetail {
	assets := map[string]MockAssetDetail{
		"ast-001": {
			ID: "ast-001", AssetNumber: "FA-001",
			Name: "Office Laptop (Dell XPS 15)", Description: "15-inch laptop for administrative use",
			CategoryName: "IT Equipment", LocationName: "Main Office",
			AcquisitionCost: "\u20b185,000.00", AcquisitionCostRaw: "85000",
			SalvageValue: "\u20b15,000.00", SalvageValueRaw: "5000",
			UsefulLifeMonths: "60", DepreciationMethod: "Straight Line",
			DepreciationMethodKey: "straight_line",
			BookValue:             "\u20b142,500.00", Status: "active",
			DepreciationSchedule: []DepreciationRow{
				{Period: "Year 1 (2024)", StartValue: "\u20b185,000.00", Expense: "\u20b116,000.00", EndValue: "\u20b169,000.00", Accumulated: "\u20b116,000.00"},
				{Period: "Year 2 (2025)", StartValue: "\u20b169,000.00", Expense: "\u20b116,000.00", EndValue: "\u20b153,000.00", Accumulated: "\u20b132,000.00"},
				{Period: "Year 3 (2026)", StartValue: "\u20b153,000.00", Expense: "\u20b116,000.00", EndValue: "\u20b137,000.00", Accumulated: "\u20b148,000.00"},
				{Period: "Year 4 (2027)", StartValue: "\u20b137,000.00", Expense: "\u20b116,000.00", EndValue: "\u20b121,000.00", Accumulated: "\u20b164,000.00"},
				{Period: "Year 5 (2028)", StartValue: "\u20b121,000.00", Expense: "\u20b116,000.00", EndValue: "\u20b15,000.00", Accumulated: "\u20b180,000.00"},
			},
			MaintenanceRecords: []MaintenanceRow{
				{Date: "2026-02-15", Type: "Preventive", Description: "Battery replacement", Status: "Completed", StatusClass: "success", Cost: "\u20b13,500.00"},
				{Date: "2026-01-10", Type: "Repair", Description: "Screen hinge repair", Status: "Completed", StatusClass: "success", Cost: "\u20b18,200.00"},
				{Date: "2025-12-01", Type: "Preventive", Description: "Annual checkup and cleaning", Status: "Completed", StatusClass: "success", Cost: "\u20b11,500.00"},
			},
			TransactionHistory: []TransactionRow{
				{Date: "2026-02-15", Type: "Maintenance", Description: "Battery replacement", Amount: "\u20b13,500.00", Reference: "MNT-003"},
				{Date: "2026-01-10", Type: "Maintenance", Description: "Screen hinge repair", Amount: "\u20b18,200.00", Reference: "MNT-002"},
				{Date: "2025-12-01", Type: "Maintenance", Description: "Annual checkup and cleaning", Amount: "\u20b11,500.00", Reference: "MNT-001"},
				{Date: "2024-01-15", Type: "Acquisition", Description: "Initial purchase \u2014 PO #2024-0158", Amount: "\u20b185,000.00", Reference: "PO-2024-0158"},
			},
		},
		"ast-002": {
			ID: "ast-002", AssetNumber: "FA-002",
			Name: "Salon Chair (Hydraulic)", Description: "Hydraulic adjustable salon chair",
			CategoryName: "Furniture", LocationName: "Branch 1",
			AcquisitionCost: "\u20b125,000.00", AcquisitionCostRaw: "25000",
			SalvageValue: "\u20b12,000.00", SalvageValueRaw: "2000",
			UsefulLifeMonths: "120", DepreciationMethod: "Straight Line",
			DepreciationMethodKey: "straight_line",
			BookValue:             "\u20b118,750.00", Status: "active",
			DepreciationSchedule: []DepreciationRow{
				{Period: "Year 1 (2024)", StartValue: "\u20b125,000.00", Expense: "\u20b12,300.00", EndValue: "\u20b122,700.00", Accumulated: "\u20b12,300.00"},
				{Period: "Year 2 (2025)", StartValue: "\u20b122,700.00", Expense: "\u20b12,300.00", EndValue: "\u20b120,400.00", Accumulated: "\u20b14,600.00"},
				{Period: "Year 3 (2026)", StartValue: "\u20b120,400.00", Expense: "\u20b12,300.00", EndValue: "\u20b118,100.00", Accumulated: "\u20b16,900.00"},
			},
			MaintenanceRecords: []MaintenanceRow{
				{Date: "2026-01-20", Type: "Repair", Description: "Hydraulic pump replacement", Status: "Completed", StatusClass: "success", Cost: "\u20b14,800.00"},
			},
			TransactionHistory: []TransactionRow{
				{Date: "2026-01-20", Type: "Maintenance", Description: "Hydraulic pump replacement", Amount: "\u20b14,800.00", Reference: "MNT-004"},
				{Date: "2024-03-01", Type: "Acquisition", Description: "Initial purchase", Amount: "\u20b125,000.00", Reference: "PO-2024-0201"},
			},
		},
		"ast-003": {
			ID: "ast-003", AssetNumber: "FA-003",
			Name: "Hair Dryer (Professional)", Description: "Professional-grade hair dryer",
			CategoryName: "Equipment", LocationName: "Branch 1",
			AcquisitionCost: "\u20b112,000.00", AcquisitionCostRaw: "12000",
			SalvageValue: "\u20b11,000.00", SalvageValueRaw: "1000",
			UsefulLifeMonths: "36", DepreciationMethod: "Straight Line",
			DepreciationMethodKey: "straight_line",
			BookValue:             "\u20b16,000.00", Status: "active",
			DepreciationSchedule: []DepreciationRow{
				{Period: "Year 1 (2025)", StartValue: "\u20b112,000.00", Expense: "\u20b13,666.67", EndValue: "\u20b18,333.33", Accumulated: "\u20b13,666.67"},
				{Period: "Year 2 (2026)", StartValue: "\u20b18,333.33", Expense: "\u20b13,666.67", EndValue: "\u20b14,666.66", Accumulated: "\u20b17,333.34"},
				{Period: "Year 3 (2027)", StartValue: "\u20b14,666.66", Expense: "\u20b13,666.66", EndValue: "\u20b11,000.00", Accumulated: "\u20b111,000.00"},
			},
			MaintenanceRecords: nil,
			TransactionHistory: []TransactionRow{
				{Date: "2025-06-15", Type: "Acquisition", Description: "Initial purchase", Amount: "\u20b112,000.00", Reference: "PO-2025-0089"},
			},
		},
	}

	if asset, ok := assets[id]; ok {
		return asset
	}

	// Default mock — "Unknown Asset" is a dev fallback; live data comes from DB
	return MockAssetDetail{
		ID: id, AssetNumber: "FA-???", Name: "Unknown Asset",
		Description: "\u2014", CategoryName: "\u2014", LocationName: "\u2014",
		AcquisitionCost: "\u20b10.00", AcquisitionCostRaw: "0",
		SalvageValue: "\u20b10.00", SalvageValueRaw: "0",
		UsefulLifeMonths: "\u2014", DepreciationMethod: "\u2014",
		DepreciationMethodKey: "straight_line",
		BookValue:             "\u20b10.00", Status: "active",
	}
}
