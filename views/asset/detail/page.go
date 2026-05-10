package detail

import (
	"context"
	"fmt"
	"log"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Mock data types (intermediate -- converted to TableConfig for rendering)
// ---------------------------------------------------------------------------

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
	// MeasurementModel controls which page-pillar actions are shown.
	// Values: "COST" (default) or "REVALUATION".
	MeasurementModel   string
	TransactionHistory []TransactionRow
}

// ---------------------------------------------------------------------------
// Asset action resolution (Hard rule #4 — CTA branching is a struct + helper)
// ---------------------------------------------------------------------------

// AssetActions holds the resolved page-pillar action buttons for a given asset state.
type AssetActions struct {
	// Primary is the main CTA button shown in the page pillar.
	// Empty URL means no primary button.
	Primary AssetAction
	// Secondary holds optional secondary actions (e.g. Edit when Revalue is primary).
	Secondary []AssetAction
}

// AssetAction is a single CTA button in the page pillar.
type AssetAction struct {
	URL      string
	Label    string
	TestID   string
	Variant  string // "primary" or "secondary"
}

// resolveAssetActions derives the page-pillar actions for a given asset.
// Locked decision (flow.md §Action pillar):
//
//	IN_SERVICE + COST         → primary [Edit], no secondary
//	IN_SERVICE + REVALUATION  → primary [Revalue], secondary [Edit]
//	FULLY_DEPRECIATED         → primary [Edit], no secondary
//	DISPOSED                  → empty
//	[Dispose] is HIDDEN this plan.
func resolveAssetActions(asset MockAssetDetail, routes fycha.AssetRoutes, labels fycha.AssetLabels) AssetActions {
	editURL := route.ResolveURL(routes.EditURL, "id", asset.ID)
	revalURL := routes.RevaluationFor(asset.ID)

	switch asset.Status {
	case "active", "in_service", "IN_SERVICE":
		if asset.MeasurementModel == "REVALUATION" {
			return AssetActions{
				Primary: AssetAction{
					URL:     revalURL,
					Label:   labels.Actions.Revalue,
					TestID:  "asset-detail-action-revalue",
					Variant: "primary",
				},
				Secondary: []AssetAction{
					{
						URL:     editURL,
						Label:   labels.Actions.Edit,
						TestID:  "asset-detail-action-edit",
						Variant: "secondary",
					},
				},
			}
		}
		return AssetActions{
			Primary: AssetAction{
				URL:     editURL,
				Label:   labels.Actions.Edit,
				TestID:  "asset-detail-action-edit",
				Variant: "primary",
			},
		}
	case "fully_depreciated", "FULLY_DEPRECIATED":
		return AssetActions{
			Primary: AssetAction{
				URL:     editURL,
				Label:   labels.Actions.Edit,
				TestID:  "asset-detail-action-edit",
				Variant: "primary",
			},
		}
	default: // disposed, inactive, unknown
		return AssetActions{}
	}
}

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// DetailViewDeps holds view dependencies.
type DetailViewDeps struct {
	attachment.AttachmentOps
	auditlog.AuditOps

	Routes       fycha.AssetRoutes
	Labels       fycha.AssetLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// DepreciationRunLabels for the lapsing-actual-schedule tab
	DepreciationRunLabels fycha.DepreciationRunLabels
	// AssetRevaluationLabels for the revaluation drawer
	AssetRevaluationLabels fycha.AssetRevaluationLabels
}

// PageData holds the data for the asset detail page.
type PageData struct {
	types.PageData
	ContentTemplate       string
	Labels                fycha.AssetLabels
	DepreciationRunLabels fycha.DepreciationRunLabels
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
	MeasurementModel      string
	StatusVariant         string
	// Page-pillar resolved actions
	AssetActions     AssetActions
	EditURL          string
	RevaluationURL   string
	DeprecRunURL     string
	CanEdit          bool
	TransactionTable *types.TableConfig
	AttachmentTable  *types.TableConfig
	// Audit history tab
	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the asset detail view (full page).
func NewView(deps *DetailViewDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")

		activeTab := viewCtx.Request.URL.Query().Get("tab")
		if activeTab == "" {
			activeTab = "info"
		}

		perms := view.GetUserPermissions(ctx)
		pageData := buildPageData(deps, id, activeTab, viewCtx, perms)

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "asset-detail"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("asset-detail", pageData)
	})
}

// NewTabAction creates the tab action view (partial -- returns only the tab content).
func NewTabAction(deps *DetailViewDeps) view.View {
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
		switch tab {
		case "attachments":
			templateName = "attachment-tab"
		case "history":
			templateName = "asset-tab-history"
		case "lapsing-actual-schedule":
			templateName = "asset-tab-lapsing-actual-schedule"
		case "transaction-ledger":
			templateName = "asset-tab-transaction-ledger"
		}
		return view.OK(templateName, pageData)
	})
}

// ---------------------------------------------------------------------------
// Page data builder
// ---------------------------------------------------------------------------

func buildPageData(deps *DetailViewDeps, id, activeTab string, viewCtx *view.ViewContext, perms *types.UserPermissions) *PageData {
	asset := getMockAsset(id)

	statusVariant := "success"
	if asset.Status == "inactive" {
		statusVariant = "warning"
	}

	tabItems := buildTabItems(id, deps.Labels, deps.Routes)
	assetActions := resolveAssetActions(asset, deps.Routes, deps.Labels)

	pageData := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          asset.Name,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      "asset",
			ActiveSubNav:   "assets-fixed",
			HeaderTitle:    asset.Name,
			HeaderSubtitle: fmt.Sprintf("%s | %s", asset.AssetNumber, asset.CategoryName),
			HeaderIcon:     "icon-box",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate:       "asset-detail-content",
		Labels:                deps.Labels,
		DepreciationRunLabels: deps.DepreciationRunLabels,
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
		MeasurementModel:      asset.MeasurementModel,
		StatusVariant:         statusVariant,
		AssetActions:          assetActions,
		EditURL:               route.ResolveURL(deps.Routes.EditURL, "id", id),
		RevaluationURL:        deps.Routes.RevaluationFor(id),
		DeprecRunURL:          deps.Routes.DepreciationRunFor(id),
		CanEdit:               perms.Can("asset", "update"),
		TransactionTable:      buildTransactionTable(asset.TransactionHistory, deps.Labels, deps.TableLabels),
	}

	if activeTab == "attachments" {
		if deps.ListAttachments != nil {
			cfg := attachmentConfig(deps)
			resp, err := deps.ListAttachments(viewCtx.Request.Context(), cfg.EntityType, id)
			if err != nil {
				log.Printf("Failed to list attachments: %v", err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pageData.AttachmentTable = attachment.BuildTable(items, cfg, id)
		}
	}

	if activeTab == "history" {
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			auditResp, err := deps.ListAuditHistory(viewCtx.Request.Context(), &auditlog.ListAuditRequest{
				EntityType:  "asset",
				EntityID:    id,
				Limit:       20,
				CursorToken: cursor,
			})
			if err != nil {
				log.Printf("Failed to load audit history: %v", err)
			}
			if auditResp != nil {
				pageData.AuditEntries = auditResp.Entries
				pageData.AuditHasNext = auditResp.HasNext
				pageData.AuditNextCursor = auditResp.NextCursor
			}
		}
		pageData.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "history"
	}

	return pageData
}

// buildTabItems constructs the 5-tab nav for the asset detail page.
// Locked tab order (flow.md §Tab order, 2026-05-10):
//
//	info | lapsing-actual-schedule | transaction-ledger | attachments | history
//
// Tabs removed: depreciation, maintenance.
// Renamed: transactions → transaction-ledger, audit-history → history.
func buildTabItems(id string, labels fycha.AssetLabels, routes fycha.AssetRoutes) []pyeza.TabItem {
	base := route.ResolveURL(routes.DetailURL, "id", id)
	action := route.ResolveURL(routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: labels.Detail.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info", Count: 0, Disabled: false},
		{Key: "lapsing-actual-schedule", Label: labels.Detail.Tabs.LapsingActualSchedule, Href: base + "?tab=lapsing-actual-schedule", HxGet: action + "lapsing-actual-schedule", Icon: "icon-trending-down", Count: 0, Disabled: false},
		{Key: "transaction-ledger", Label: labels.Detail.Tabs.TransactionLedger, Href: base + "?tab=transaction-ledger", HxGet: action + "transaction-ledger", Icon: "icon-list", Count: 0, Disabled: false},
		{Key: "attachments", Label: labels.Detail.Tabs.Attachments, Href: base + "?tab=attachments", HxGet: action + "attachments", Icon: "icon-paperclip", Count: 0, Disabled: false},
		{Key: "history", Label: labels.Detail.Tabs.History, Href: base + "?tab=history", HxGet: action + "history", Icon: "icon-clock", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Table builders -- convert mock data into pyeza TableConfig
// ---------------------------------------------------------------------------

func buildTransactionTable(history []TransactionRow, labels fycha.AssetLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: labels.Columns.Date},
		{Key: "type", Label: labels.Columns.Type},
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
			AcquisitionCost: "₱85,000.00", AcquisitionCostRaw: "85000",
			SalvageValue: "₱5,000.00", SalvageValueRaw: "5000",
			UsefulLifeMonths: "60", DepreciationMethod: "Straight Line",
			DepreciationMethodKey: "straight_line",
			BookValue:             "₱42,500.00", Status: "active",
			MeasurementModel:      "COST",
			TransactionHistory: []TransactionRow{
				{Date: "2026-02-15", Type: "Maintenance", Description: "Battery replacement", Amount: "₱3,500.00", Reference: "MNT-003"},
				{Date: "2026-01-10", Type: "Maintenance", Description: "Screen hinge repair", Amount: "₱8,200.00", Reference: "MNT-002"},
				{Date: "2025-12-01", Type: "Maintenance", Description: "Annual checkup and cleaning", Amount: "₱1,500.00", Reference: "MNT-001"},
				{Date: "2024-01-15", Type: "Acquisition", Description: "Initial purchase — PO #2024-0158", Amount: "₱85,000.00", Reference: "PO-2024-0158"},
			},
		},
		"ast-002": {
			ID: "ast-002", AssetNumber: "FA-002",
			Name: "Salon Chair (Hydraulic)", Description: "Hydraulic adjustable salon chair",
			CategoryName: "Furniture", LocationName: "Branch 1",
			AcquisitionCost: "₱25,000.00", AcquisitionCostRaw: "25000",
			SalvageValue: "₱2,000.00", SalvageValueRaw: "2000",
			UsefulLifeMonths: "120", DepreciationMethod: "Straight Line",
			DepreciationMethodKey: "straight_line",
			BookValue:             "₱18,750.00", Status: "active",
			MeasurementModel:      "COST",
			TransactionHistory: []TransactionRow{
				{Date: "2026-01-20", Type: "Maintenance", Description: "Hydraulic pump replacement", Amount: "₱4,800.00", Reference: "MNT-004"},
				{Date: "2024-03-01", Type: "Acquisition", Description: "Initial purchase", Amount: "₱25,000.00", Reference: "PO-2024-0201"},
			},
		},
		"ast-003": {
			ID: "ast-003", AssetNumber: "FA-003",
			Name: "Hair Dryer (Professional)", Description: "Professional-grade hair dryer",
			CategoryName: "Equipment", LocationName: "Branch 1",
			AcquisitionCost: "₱12,000.00", AcquisitionCostRaw: "12000",
			SalvageValue: "₱1,000.00", SalvageValueRaw: "1000",
			UsefulLifeMonths: "36", DepreciationMethod: "Straight Line",
			DepreciationMethodKey: "straight_line",
			BookValue:             "₱6,000.00", Status: "active",
			MeasurementModel:      "COST",
			TransactionHistory: []TransactionRow{
				{Date: "2025-06-15", Type: "Acquisition", Description: "Initial purchase", Amount: "₱12,000.00", Reference: "PO-2025-0089"},
			},
		},
	}

	if asset, ok := assets[id]; ok {
		return asset
	}

	// Default mock — dev fallback; live data comes from DB
	return MockAssetDetail{
		ID: id, AssetNumber: "FA-???", Name: "—",
		Description: "—", CategoryName: "—", LocationName: "—",
		AcquisitionCost: "₱0.00", AcquisitionCostRaw: "0",
		SalvageValue: "₱0.00", SalvageValueRaw: "0",
		UsefulLifeMonths: "—", DepreciationMethod: "—",
		DepreciationMethodKey: "straight_line",
		BookValue:             "₱0.00", Status: "active",
		MeasurementModel:      "COST",
	}
}
