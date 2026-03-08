# fycha-golang

Recording and accounting domain package for Ryta OS. Provides reusable report views (P&L family), fixed asset management views, period-based filtering, storage handling, and HTMX helpers for Go HTMX admin applications built on pyeza-golang.

**Module path:** `github.com/erniealice/fycha-golang`

**Dependencies:**
- `github.com/erniealice/pyeza-golang` -- UI framework (views, types, templates)
- `github.com/erniealice/esqyma` -- protobuf schemas (ledger/reporting)

## Package Structure

```
packages/fycha-golang-ryta/
  go.mod
  go.sum
  package_dir.go          -- runtime.Caller(0) for resolving package directory
  datasource.go           -- DataSource interface (report data access)
  routes.go               -- Route path constants
  routes_config.go        -- Configurable route structs (ReportsRoutes, AssetRoutes)
  labels.go               -- All label structs + MapTableLabels/MapBulkConfig helpers
  report_filter.go        -- FilterState, period presets, date parsing
  htmx.go                 -- HTMXSuccess/HTMXError response helpers
  assets.go               -- CopyStyles/CopyStaticAssets for CSS/JS asset pipeline
  storage_handler.go      -- StorageHandler for serving files from object storage
  assets/
    css/
      fycha-report.css            -- Report page styles
      fycha-asset-dashboard.css   -- Asset dashboard styles
    js/
      fycha-report-filter.js      -- Client-side filter interaction
  views/
    reports/
      embed.go                    -- //go:embed templates/*.html
      templates/
        dashboard.html
        revenue.html
        expenses-report.html
        gross-profit.html
        cost-of-sales.html
        net-profit.html
        report-filter.html
        report-filter-btn.html
      dashboard/page.go           -- Reports dashboard view
      revenue/page.go             -- Revenue report view
      expenses/page.go            -- Expenses report view
      gross_profit/page.go        -- Gross profit report view
      cost_of_sales/page.go       -- Cost of sales report view
      net_profit/page.go          -- Net profit (P&L) report view
    asset/
      embed.go                    -- //go:embed templates/*.html
      templates/
        list.html
        detail.html
        dashboard.html
        asset-drawer-form.html
      list/page.go                -- Asset list + table-only refresh view
      detail/page.go              -- Asset detail + tab action views
      dashboard/page.go           -- Asset dashboard view
      action/action.go            -- CRUD action handlers (add, edit, delete, status)
```

## DataSource Interface

The `DataSource` interface abstracts report data access. Consumer apps provide an implementation (typically backed by espyna postgres adapters).

```go
type DataSource interface {
    // GetGrossProfitReport returns grouped P&L data with summary totals.
    // Uses protobuf request/response from esqyma ledger/reporting/gross_profit.
    GetGrossProfitReport(ctx context.Context, req *reportpb.GrossProfitReportRequest) (*reportpb.GrossProfitReportResponse, error)

    // ListRevenue returns revenue records for a date range as generic maps.
    // Map keys: id, reference_number, customer_name, currency, total_amount, status.
    ListRevenue(ctx context.Context, start, end *time.Time) ([]map[string]any, error)

    // ListExpenses returns expense records for a date range as generic maps.
    // Map keys: id, reference_number, vendor_name, category, expenditure_date, currency, total_amount, status.
    ListExpenses(ctx context.Context, start, end *time.Time) ([]map[string]any, error)
}
```

The `GetGrossProfitReport` method is used by dashboard, gross profit, cost of sales, and net profit views. The protobuf types come from `github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit`.

## Route Constants

All route paths are defined as package-level constants in `routes.go`:

### Report Routes

| Constant | Path |
|----------|------|
| `ReportsBaseURL` | `/app/reports/` |
| `ReportsDashboardURL` | `/app/reports/dashboard` |
| `ReportsRevenueURL` | `/app/reports/revenue` |
| `ReportsCostOfSalesURL` | `/app/reports/cost-of-sales` |
| `ReportsGrossProfitURL` | `/app/reports/gross-profit` |
| `ReportsExpensesURL` | `/app/reports/expenses` |
| `ReportsNetProfitURL` | `/app/reports/net-profit` |
| `CashBookURL` | `/app/cash/reports/cash-book` |

### Asset Routes

| Constant | Path |
|----------|------|
| `AssetDashboardURL` | `/app/assets/dashboard` |
| `AssetListURL` | `/app/assets/list/{status}` |
| `AssetDetailURL` | `/app/assets/detail/{id}` |
| `AssetTabActionURL` | `/action/assets/{id}/tab/{tab}` |
| `AssetTableURL` | `/action/assets/table/{status}` |
| `AssetAddURL` | `/action/assets/add` |
| `AssetEditURL` | `/action/assets/edit/{id}` |
| `AssetDeleteURL` | `/action/assets/delete` |
| `AssetBulkDeleteURL` | `/action/assets/bulk-delete` |
| `AssetSetStatusURL` | `/action/assets/set-status` |
| `AssetBulkSetStatusURL` | `/action/assets/bulk-set-status` |
| `AssetLapsingScheduleURL` | `/app/assets/reports/lapsing-schedule` |
| `AssetDepreciationPoliciesURL` | `/app/assets/settings/depreciation-policies` |

### Storage

| Constant | Path |
|----------|------|
| `StorageImagesPrefix` | `/storage/images` |

## Route Config Structs

Route paths are configurable via structs that support a three-level override system:

1. **Level 1 (Go defaults):** `DefaultReportsRoutes()` / `DefaultAssetRoutes()` populate from package constants
2. **Level 2 (JSON overrides):** Struct fields have `json` tags for loading from config files
3. **Level 3 (Go field assignment):** Direct field assignment for one-off customizations

Each config struct has a `RouteMap()` method that returns `map[string]string` of dot-notation keys to paths, useful for template rendering.

### ReportsRoutes

```go
type ReportsRoutes struct {
    DashboardURL    string `json:"dashboard_url"`
    RevenueURL      string `json:"revenue_url"`
    CostOfSalesURL  string `json:"cost_of_sales_url"`
    GrossProfitURL  string `json:"gross_profit_url"`
    ExpensesURL     string `json:"expenses_url"`
    NetProfitURL    string `json:"net_profit_url"`
}
```

`RouteMap()` keys: `reports.dashboard`, `reports.revenue`, `reports.cost_of_sales`, `reports.gross_profit`, `reports.expenses`, `reports.net_profit`.

### AssetRoutes

```go
type AssetRoutes struct {
    DashboardURL             string `json:"dashboard_url"`
    ListURL                  string `json:"list_url"`
    DetailURL                string `json:"detail_url"`
    TabActionURL             string `json:"tab_action_url"`
    TableURL                 string `json:"table_url"`
    AddURL                   string `json:"add_url"`
    EditURL                  string `json:"edit_url"`
    DeleteURL                string `json:"delete_url"`
    BulkDeleteURL            string `json:"bulk_delete_url"`
    SetStatusURL             string `json:"set_status_url"`
    BulkSetStatusURL         string `json:"bulk_set_status_url"`
    LapsingScheduleURL       string `json:"lapsing_schedule_url"`
    DepreciationPoliciesURL  string `json:"depreciation_policies_url"`
}
```

`RouteMap()` keys: `asset.dashboard`, `asset.list`, `asset.detail`, `asset.tab_action`, `asset.table`, `asset.add`, `asset.edit`, `asset.delete`, `asset.bulk_delete`, `asset.set_status`, `asset.bulk_set_status`, `asset.lapsing_schedule`, `asset.depreciation_policies`.

## Label Structs

All labels are fully translatable via lyngua JSON files. Consumer apps load labels at startup and pass them through the module deps.

### Report Labels Hierarchy

```
ReportsLabels
  +-- GrossProfit   GrossProfitLabels    -- title, column groups, group-by, filters, period presets, summary, empty state
  +-- Revenue       RevenueLabels        -- title, subtitle, column headers, summary, empty state
  +-- CostOfSales   CostOfSalesLabels    -- title, subtitle, column headers, summary, empty state
  +-- Expenses      ExpensesLabels       -- title, subtitle, column headers, summary, empty state
  +-- NetProfit     NetProfitLabels       -- title, subtitle, P&L line items, summary
  +-- Dashboard     DashboardLabels      -- title, subtitle, KPI card labels, navigation card descriptions
  +-- Period        PeriodLabels         -- shared period preset labels (thisMonth, lastMonth, etc.)
```

### PeriodLabels

Shared across all report views for the filter sheet:

```go
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
```

### GrossProfitLabels

The most complex label struct, covering column groups, group-by options, filters, period presets, and summary:

```go
type GrossProfitLabels struct {
    Title              string   // Page title
    // Column group headers
    RevenueGroup       string
    ProfitabilityGroup string
    VolumeGroup        string
    // Column labels
    GrossRevenue, Discount, NetRevenue string
    COGS, GrossProfit, Margin          string
    UnitsSold, Transactions            string
    // Group-by options
    GroupBy, GroupByProduct, GroupByLocation, GroupByCategory string
    GroupByMonthly, GroupByQuarterly                         string
    // Filter options
    FilterProduct, FilterLocation, FilterCategory, FilterAll string
    // Period presets (duplicated from PeriodLabels for backward compat)
    PeriodThisMonth, PeriodLastMonth, PeriodThisQuarter string
    PeriodLastQuarter, PeriodThisYear, PeriodLastYear   string
    PeriodCustom, DateStart, DateEnd, Apply             string
    // Summary bar
    SummaryNetRevenue, SummaryCogs, SummaryGrossProfit, SummaryMargin string
    // Table footer + empty state
    Totals, EmptyTitle, EmptyMessage string
}
```

### Asset Labels Hierarchy

```
AssetLabels
  +-- Page       AssetPageLabels          -- heading, caption (per status variant: active/inactive)
  +-- Buttons    AssetButtonLabels        -- addAsset
  +-- Columns    AssetColumnLabels        -- table columns for list + sub-tables (depreciation, maintenance, transactions, cost of sales)
  +-- Empty      AssetEmptyLabels         -- active/inactive empty states
  +-- Form       AssetFormLabels          -- drawer form field labels + placeholders (17 fields)
  +-- Actions    AssetActionLabels        -- CRUD labels, 6 confirm messages, 5 error messages
  +-- Detail     AssetDetailLabels        -- basic info labels (12 fields), tab labels (4 tabs), empty states per tab (3 pairs)
  +-- Dashboard  AssetDashboardLabels     -- stat card labels (4), activity feed labels (5)
```

`DefaultAssetLabels()` provides hardcoded English defaults. Consumer apps should override via lyngua JSON files.

### Helper Functions

```go
// MapTableLabels maps pyeza.CommonLabels into types.TableLabels for table rendering.
func MapTableLabels(common pyeza.CommonLabels) types.TableLabels

// MapBulkConfig returns a BulkActionsConfig with labels from CommonLabels.
func MapBulkConfig(common pyeza.CommonLabels) types.BulkActionsConfig
```

## Report Filter System

The filter system provides period-based date filtering with server-side preset resolution.

### Core Types

```go
// ReportFilter holds parsed filter parameters for data queries.
type ReportFilter struct {
    StartDate *time.Time
    EndDate   *time.Time
    Period    string   // "thisMonth", "lastMonth", "thisQuarter", "lastQuarter", "thisYear", "lastYear", "custom"
    GroupBy   string   // "product", "customer", "status", "monthly", "quarterly", "vendor", "category"
}

// FilterState holds the current filter state for template rendering.
type FilterState struct {
    ActivePreset   string
    StartDate      string
    EndDate        string
    GroupBy        string
    GroupByOptions []FilterOption
    PeriodPresets  []FilterOption
}

// FilterOption holds a dropdown or button option.
type FilterOption struct {
    Value    string
    Label    string
    Selected bool
}

// FilterSheetData is passed to the report-filter-sheet template.
type FilterSheetData struct {
    Filter       FilterState
    PeriodLabels PeriodLabels
    ReportURL    string
}
```

### Period Preset Resolution

`ParsePeriodPreset(preset string) (start, end time.Time)` computes a date range from a named preset in local timezone:

| Preset | Start | End |
|--------|-------|-----|
| `thisMonth` (default) | 1st of current month | now |
| `lastMonth` | 1st of previous month | last second of previous month |
| `thisQuarter` | 1st of current quarter | now |
| `lastQuarter` | 1st of previous quarter | last second of previous quarter |
| `thisYear` | January 1st | now |
| `lastYear` | January 1st last year | December 31st last year |

### Helper Functions

```go
// DefaultPeriodPresets returns the standard 7 period options with the active one marked.
func DefaultPeriodPresets(labels PeriodLabels, active string) []FilterOption

// ActiveFilterCount returns the number of non-default filters applied.
// Non-default period counts as 1, non-"product" group-by counts as 1.
func ActiveFilterCount(filter FilterState) int
```

### Supporting Display Types

```go
// SummaryMetric holds a single summary bar metric.
type SummaryMetric struct {
    Label     string
    Value     string
    Highlight bool
    Variant   string  // "success", "warning", "danger"
}

// PLLineItem holds a single line in a P&L statement (used by net profit view).
type PLLineItem struct {
    Label   string
    Value   string
    IsTotal bool
    Variant string
}
```

## Report View Sub-Packages

All report views follow the same pattern: a `Deps` struct for dependency injection, a `PageData` struct embedding `types.PageData`, and a `NewView(deps *Deps) view.View` constructor that returns a `view.ViewFunc`.

Each view handles two rendering modes:
- **Full page** (non-HTMX): renders the full page template (e.g., `"revenue"`)
- **Partial** (HTMX): renders only the content template (e.g., `"revenue-content"`)

Views that support filtering also handle `?sheet=filters` to return the filter sheet partial via the `"report-filter-sheet"` template.

### views/reports/dashboard

Reports landing page with KPI summary cards and navigation cards to individual reports.

```go
type Deps struct {
    Routes       fycha.ReportsRoutes
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
}
```

- Calls `DB.GetGrossProfitReport` and `DB.ListExpenses` for current month
- Computes net profit and net margin
- Renders 4 KPI summary cards: Revenue, Expenses, Net Profit, Net Margin
- Renders 5 navigation cards linking to individual report pages
- Net profit variant coloring: `"danger"` if negative, `"warning"` if margin < 10%, `"success"` otherwise

Templates: `reports-dashboard` (full page), `reports-dashboard-content` (HTMX partial).

### views/reports/revenue

Revenue transaction list with period filtering.

```go
type Deps struct {
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

- Calls `DB.ListRevenue` with resolved date range
- Summary bar: Total Revenue, Transaction Count, Average
- Table columns: Reference, Customer, Amount, Status
- Status badge variants: completed/paid=success, pending=warning, cancelled/refunded=danger
- Table ID: `revenue-table`

Templates: `revenue` / `revenue-content`.

### views/reports/gross_profit

Grouped gross profit analysis with column groups and group-by filtering.

```go
type Deps struct {
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

- Calls `DB.GetGrossProfitReport` with group-by, period, and optional product/location/category filters
- Supports group-by: `product`, `location`, `category`, `monthly`, `quarterly`
- Uses `ColumnGroups` for grouped table headers:
  - Revenue: Gross Revenue, Discount, Net Revenue
  - Profitability: COGS, Gross Profit, Margin
  - Volume: Units Sold, Transactions
- Summary bar: Net Revenue, COGS, Gross Profit, Margin
- Margin variant coloring: < 15% = danger, < 30% = warning, >= 30% = success
- Includes totals row with ID `"__totals__"`
- Table ID: `grossProfitTable`

Templates: `gross-profit` / `gross-profit-content`.

### views/reports/cost_of_sales

Cost of goods sold analysis per product.

```go
type Deps struct {
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

- Calls `DB.GetGrossProfitReport` grouped by product
- Summary bar: Total COGS, Revenue, COGS Ratio, Units
- Table columns: Item, COGS, Net Revenue, COGS%, Units
- Default sort: `cogs` descending
- Table ID: `cost-of-sales-table`

Templates: `cost-of-sales` / `cost-of-sales-content`.

### views/reports/expenses

Expense transaction list with period filtering.

```go
type Deps struct {
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

- Calls `DB.ListExpenses` with resolved date range
- Summary bar: Total Expenses, Count, Approved, Pending
- Table columns: Reference, Vendor, Category, Date, Amount, Status
- Status badge variants: paid=success, approved=info, pending=warning, cancelled=danger, draft=default
- Default sort: `date` descending
- Table ID: `expenses-report-table`

Templates: `expenses-report` / `expenses-report-content`.

### views/reports/net_profit

P&L income statement with line items.

```go
type Deps struct {
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

- Calls both `DB.GetGrossProfitReport` and `DB.ListExpenses`
- Computes: Net Profit = Gross Profit - Expenses, Net Margin, Gross Margin
- Summary bar: Revenue, Gross Profit, Expenses, Net Profit
- Renders 7 P&L line items:
  1. Revenue
  2. Cost of Sales
  3. **Gross Profit** (IsTotal)
  4. Gross Margin
  5. Expenses
  6. **Net Profit** (IsTotal)
  7. Net Margin (with variant coloring)

Templates: `net-profit` / `net-profit-content`.

## Asset View Sub-Packages

Asset views use mock data for initial UI development. The current implementation uses hardcoded `MockAsset` structs; live data will come from a database layer once the asset data model is wired.

### views/asset/list

Asset list with status-based filtering (active/inactive tabs).

```go
type Deps struct {
    Routes       fycha.AssetRoutes
    Labels       fycha.AssetLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

**Constructors:**
- `NewView(deps *Deps) view.View` -- full page list
- `NewTableView(deps *Deps) view.View` -- table-card only (HTMX refresh target after CRUD)

Features:
- Status from URL path parameter `{status}` (defaults to `"active"`)
- RBAC-aware: checks `asset/create`, `asset/update`, `asset/delete` permissions
- Row actions: View (link), Edit (drawer), Activate/Deactivate (confirm dialog), Delete (confirm dialog)
- Bulk actions: Activate/Deactivate (context-dependent on current status tab), Delete
- Primary action: "Add Asset" button (disabled if no create permission)
- Table columns: Asset Number, Name, Category, Location, Acquisition Cost, Book Value, Status
- Default sort: `asset_number` ascending
- Table ID: `assets-table`

**Mock data:** 7 assets spanning IT Equipment, Furniture, Equipment, Building Equipment categories across Main Office, Branch 1, Branch 2 locations. Book values range from 0 (fully depreciated) to 42,500.

Templates: `asset-list` (full page), uses pyeza `table-card` for table-only refresh.

### views/asset/detail

Asset detail page with tabbed content (info, depreciation, maintenance, transactions).

```go
type Deps struct {
    Routes       fycha.AssetRoutes
    Labels       fycha.AssetLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}
```

**Constructors:**
- `NewView(deps *Deps) view.View` -- full detail page
- `NewTabAction(deps *Deps) view.View` -- tab content partial (HTMX tab switching, returns `asset-tab-{tab}` template)

**Tabs:**

| Tab | Table ID | Columns |
|-----|----------|---------|
| Info | -- | Basic info key-value pairs (not a table) |
| Depreciation | `depreciation-table` | Period, Start Value, Depreciation, End Value, Accumulated |
| Maintenance | `maintenance-table` | Date, Type, Description, Status (badge), Cost |
| Transactions | `transactions-table` | Date, Type, Description, Amount, Reference |

Each sub-table is rendered as a minimal `types.TableConfig` (no search/filter toolbar). Empty states use per-tab labels.

**Mock data:** 3 detailed assets (ast-001 through ast-003) with:
- ast-001: Dell XPS 15 laptop, 5-year depreciation schedule, 3 maintenance records, 4 transactions
- ast-002: Hydraulic salon chair, 3-year partial schedule, 1 maintenance record, 2 transactions
- ast-003: Professional hair dryer, 3-year schedule, no maintenance, 1 transaction

Unknown IDs return a fallback "Unknown Asset" record with em-dash placeholders.

Templates: `asset-detail` (full page), `asset-tab-{info|depreciation|maintenance|transactions}` (tab partials).

### views/asset/dashboard

Asset register dashboard with KPI stat cards and activity feed.

```go
type Deps struct {
    Routes       fycha.AssetRoutes
    Labels       fycha.AssetLabels
    CommonLabels pyeza.CommonLabels
}
```

- Mock statistics: Total Assets (24), Total Book Value (P1,245,750.00), Fully Deprecated (3), Under Maintenance (2)
- Recent activity feed (3 items) with SVG icon HTML, title, description, and time ago
- Activity types: acquired (box icon), maintenance (tool icon), depreciation (trending-down icon)

Templates: `asset-dashboard` (full page).

### views/asset/action

CRUD action handlers for asset management.

```go
type Deps struct {
    Routes fycha.AssetRoutes
    Labels fycha.AssetLabels
}
```

**Action constructors (all return `view.View`):**

| Constructor | HTTP Methods | Permission | Description |
|-------------|-------------|-----------|-------------|
| `NewAddAction(deps)` | GET, POST | `asset/create` | GET returns drawer form; POST creates asset (mock) |
| `NewEditAction(deps)` | GET, POST | `asset/update` | GET returns pre-filled drawer form; POST updates asset (mock) |
| `NewDeleteAction(deps)` | POST | `asset/delete` | Deletes single asset by ID (query param or form value) |
| `NewBulkDeleteAction(deps)` | POST | `asset/delete` | Deletes multiple assets from multipart form `id` field |
| `NewSetStatusAction(deps)` | POST | `asset/update` | Sets status to `active` or `inactive` |
| `NewBulkSetStatusAction(deps)` | POST | `asset/update` | Bulk status change with `target_status` form field |

All actions are RBAC-protected via `view.GetUserPermissions(ctx)` and use:
- `fycha.HTMXSuccess("assets-table")` on success -- triggers sheet close + table refresh
- `fycha.HTMXError(message)` on validation failure -- returns 422 with error header

Templates: `asset-drawer-form` (for GET on add/edit actions).

## Template System

Templates are embedded via `//go:embed` in two `embed.go` files:

- `views/reports/embed.go` -- `reports.TemplatesFS` embeds `views/reports/templates/*.html`
- `views/asset/embed.go` -- `asset.TemplatesFS` embeds `views/asset/templates/*.html`

Consumer apps register these embedded filesystems with pyeza's template engine during container initialization.

### Report Templates

| Template File | Template Names Used in Go |
|---------------|---------------------------|
| `dashboard.html` | `reports-dashboard`, `reports-dashboard-content` |
| `revenue.html` | `revenue`, `revenue-content` |
| `expenses-report.html` | `expenses-report`, `expenses-report-content` |
| `gross-profit.html` | `gross-profit`, `gross-profit-content` |
| `cost-of-sales.html` | `cost-of-sales`, `cost-of-sales-content` |
| `net-profit.html` | `net-profit`, `net-profit-content` |
| `report-filter.html` | `report-filter-sheet` |
| `report-filter-btn.html` | (included by report page templates) |

### Asset Templates

| Template File | Template Names Used in Go |
|---------------|---------------------------|
| `list.html` | `asset-list`, `asset-list-content` |
| `detail.html` | `asset-detail`, `asset-detail-content`, `asset-tab-info`, `asset-tab-depreciation`, `asset-tab-maintenance`, `asset-tab-transactions` |
| `dashboard.html` | `asset-dashboard`, `asset-dashboard-content` |
| `asset-drawer-form.html` | `asset-drawer-form` |

## HTMX Helpers

```go
// HTMXSuccess returns a header-only 200 response that triggers sheet close and table refresh.
// Sets HX-Trigger: {"formSuccess":true,"refreshTable":"<tableID>"}
func HTMXSuccess(tableID string) view.ViewResult

// HTMXError returns a 422 response with an error message header.
// Sets HX-Error-Message: <message>
func HTMXError(message string) view.ViewResult
```

## Storage Handler

The `StorageHandler` serves files from any storage backend via HTTP.

```go
// StorageReader is provider-agnostic (wraps espyna, GCS, S3, mock, etc.)
type StorageReader interface {
    ReadObject(ctx context.Context, containerName, objectKey string) (*StorageReadResult, error)
}

type StorageReadResult struct {
    Content     []byte
    ContentType string
}

// Create and register:
handler := fycha.NewStorageHandler(storageReader, "my-bucket", "/storage/images")
handler.RegisterRoutes(routeRegistrar)
```

Features:
- Path traversal protection (rejects `..`)
- MIME type detection from metadata or file extension
- 24-hour cache headers (`Cache-Control: public, max-age=86400`)
- Supports common image formats (JPEG, PNG, WebP, GIF, SVG, AVIF) and PDF
- `ErrObjectNotFound` sentinel error for 404 responses

## Asset Pipeline

`CopyStyles(targetDir)` and `CopyStaticAssets(targetDir)` copy CSS and JS files from the package's `assets/` directory to the consumer app's static file directory at startup. Files are namespaced under a `fycha/` subdirectory.

```go
// In container.go initialization:
fycha.CopyStyles(cssDir)       // copies assets/css/*.css  -> {cssDir}/fycha/
fycha.CopyStaticAssets(jsDir)  // copies assets/js/*.js    -> {jsDir}/fycha/
```

Uses `runtime.Caller(0)` via `packageDir()` to discover the package source directory at runtime (same approach as centymo and entydad packages).

## Consumer App Wiring Guide

### Step 1: Create a DataSource implementation

The consumer app creates an adapter that implements `fycha.DataSource`. In service-admin, this is done via espyna's ledger reporting service:

```go
ledgerReportingService := consumer.NewLedgerReportingService(sqlDB, consumer.LedgerReportingTableConfig{
    Revenue:              "revenue",
    RevenueLineItem:      "revenue_line_item",
    InventoryTransaction: "inventory_transaction",
    InventoryItem:        "inventory_item",
    Product:              "product",
})
```

### Step 2: Initialize routes and labels

```go
import fycha "github.com/erniealice/fycha-golang"

// Routes (use defaults or override)
reportsRoutes := fycha.DefaultReportsRoutes()
assetRoutes := fycha.DefaultAssetRoutes()

// Labels (load from lyngua JSON or use defaults for assets)
var reportsLabels fycha.ReportsLabels
// ... load via lyngua.LoadPath("fycha.reports", &reportsLabels) ...
assetLabels := fycha.DefaultAssetLabels()
```

### Step 3: Register template filesystems

```go
import (
    fycha_reports "github.com/erniealice/fycha-golang/views/reports"
    fycha_asset "github.com/erniealice/fycha-golang/views/asset"
)

// Pass to pyeza template engine alongside other embedded FSes:
templateFSes := []embed.FS{
    // ... other package template FSes ...
    fycha_reports.TemplatesFS,
    fycha_asset.TemplatesFS,
}
```

### Step 4: Copy static assets

```go
fycha.CopyStyles(cssDir)
fycha.CopyStaticAssets(jsDir)
```

### Step 5: Inject route map for template access

```go
for k, v := range reportsRoutes.RouteMap() {
    routeMap[k] = v
}
for k, v := range assetRoutes.RouteMap() {
    routeMap[k] = v
}
```

### Step 6: Create app-level modules

The recommended pattern is to create thin app-level module wrappers in `apps/{app}/internal/presentation/`:

**Report module** (`report/module.go`):

```go
type ModuleDeps struct {
    Routes       fycha.ReportsRoutes
    DB           fycha.DataSource
    Labels       fycha.ReportsLabels
    CommonLabels pyeza.CommonLabels
    TableLabels  types.TableLabels
}

func NewModule(deps *ModuleDeps) *Module {
    return &Module{
        routes:      deps.Routes,
        Dashboard:   dashboardview.NewView(&dashboardview.Deps{...}),
        Revenue:     revenue.NewView(&revenue.Deps{...}),
        CostOfSales: costsales.NewView(&costsales.Deps{...}),
        GrossProfit: grossprofit.NewView(&grossprofit.Deps{...}),
        Expenses:    expensesview.NewView(&expensesview.Deps{...}),
        NetProfit:   netprofit.NewView(&netprofit.Deps{...}),
    }
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
    r.GET(m.routes.DashboardURL, m.Dashboard)
    r.GET(m.routes.RevenueURL, m.Revenue)
    r.GET(m.routes.CostOfSalesURL, m.CostOfSales)
    r.GET(m.routes.GrossProfitURL, m.GrossProfit)
    r.GET(m.routes.ExpensesURL, m.Expenses)
    r.GET(m.routes.NetProfitURL, m.NetProfit)
}
```

**Asset module** (`asset/module.go`):

```go
type ModuleDeps struct {
    Routes       fycha.AssetRoutes
    CommonLabels pyeza.CommonLabels
    Labels       fycha.AssetLabels
    TableLabels  types.TableLabels
}

func NewModule(deps *ModuleDeps) *Module {
    // Constructs: Dashboard, List, Table (refresh), Detail, TabAction,
    //             Add, Edit, Delete, BulkDelete, SetStatus, BulkSetStatus
}

func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
    r.GET(m.routes.DashboardURL, m.Dashboard)
    r.GET(m.routes.ListURL, m.List)
    r.GET(m.routes.TableURL, m.Table)
    r.GET(m.routes.DetailURL, m.Detail)
    r.GET(m.routes.TabActionURL, m.TabAction)
    r.GET(m.routes.AddURL, m.Add)
    r.POST(m.routes.AddURL, m.Add)
    r.GET(m.routes.EditURL, m.Edit)
    r.POST(m.routes.EditURL, m.Edit)
    r.POST(m.routes.DeleteURL, m.Delete)
    r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
    r.POST(m.routes.SetStatusURL, m.SetStatus)
    r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
}
```

### Step 7: Wire in views.go (RegisterAllRoutes)

```go
reportmod.NewModule(&reportmod.ModuleDeps{
    Routes:       deps.ReportsRoutes,
    DB:           deps.FychaDB,
    Labels:       deps.ReportsLabels,
    CommonLabels: deps.CommonLabels,
    TableLabels:  deps.TableLabels,
}).RegisterRoutes(routes)

assetmod.NewModule(&assetmod.ModuleDeps{
    Routes:       deps.AssetRoutes,
    CommonLabels: deps.CommonLabels,
    Labels:       deps.AssetLabels,
    TableLabels:  deps.TableLabels,
}).RegisterRoutes(routes)

// Add redirects for base URLs
routes.Redirect(fycha.ReportsBaseURL, deps.ReportsRoutes.DashboardURL)
routes.Redirect("/app/assets/", deps.AssetRoutes.DashboardURL)
```

## Consumer App Usage

### retail-admin

Uses the **reports module only** (no asset module). Wired in:
- `apps/retail-admin/internal/composition/container.go` -- initializes routes, labels, template FSes, asset pipeline
- `apps/retail-admin/internal/composition/views.go` -- creates report module, registers routes
- `apps/retail-admin/internal/presentation/report/module.go` -- thin module wrapper

### service-admin

Uses **both reports and asset modules**. Wired in:
- `apps/service-admin/internal/composition/container.go` -- initializes routes, labels (including `DefaultAssetLabels()`), template FSes, asset pipeline
- `apps/service-admin/internal/composition/views.go` -- creates both report and asset modules, plus additional report views:
  - Lapsing schedule (`AssetLapsingScheduleURL`)
  - Depreciation policies (`AssetDepreciationPoliciesURL`)
  - Cash book (`CashBookURL`)
  - Payables/receivables aging
  - Sales/purchases/expenses summaries
- `apps/service-admin/internal/presentation/report/module.go` -- report module wrapper
- `apps/service-admin/internal/presentation/asset/module.go` -- asset module wrapper

## Entities Summary

| Entity | Sub-package | View Count | Current Data Source |
|--------|-------------|-----------|---------------------|
| Reports Dashboard | `views/reports/dashboard` | 1 | Live (DataSource) |
| Revenue Report | `views/reports/revenue` | 1 | Live (DataSource) |
| Gross Profit Report | `views/reports/gross_profit` | 1 | Live (DataSource) |
| Cost of Sales Report | `views/reports/cost_of_sales` | 1 | Live (DataSource) |
| Expenses Report | `views/reports/expenses` | 1 | Live (DataSource) |
| Net Profit Report | `views/reports/net_profit` | 1 | Live (DataSource) |
| Asset Dashboard | `views/asset/dashboard` | 1 | Mock data |
| Asset List | `views/asset/list` | 2 (page + table refresh) | Mock data |
| Asset Detail | `views/asset/detail` | 2 (page + tab action) | Mock data |
| Asset Actions | `views/asset/action` | 6 (add, edit, delete, bulk-delete, set-status, bulk-set-status) | Mock handlers |
