package fycha

const (
	ReportsBaseURL         = "/app/reports/"
	ReportsDashboardURL    = "/app/reports/dashboard"
	ReportsRevenueURL      = "/app/reports/revenue"
	ReportsCostOfSalesURL  = "/app/reports/cost-of-sales"
	ReportsGrossProfitURL  = "/app/reports/gross-profit"
	ReportsExpensesURL     = "/app/reports/expenses"
	ReportsNetProfitURL    = "/app/reports/net-profit"

	// StorageImagesPrefix is the default route prefix for image serving.
	StorageImagesPrefix = "/storage/images"

	// Asset routes
	AssetDashboardURL     = "/app/assets/dashboard"
	AssetListURL          = "/app/assets/list/{status}"
	AssetDetailURL        = "/app/assets/detail/{id}"
	AssetTabActionURL     = "/action/assets/{id}/tab/{tab}"
	AssetTableURL         = "/action/assets/table/{status}"
	AssetAddURL           = "/action/assets/add"
	AssetEditURL          = "/action/assets/edit/{id}"
	AssetDeleteURL        = "/action/assets/delete"
	AssetBulkDeleteURL    = "/action/assets/bulk-delete"
	AssetSetStatusURL     = "/action/assets/set-status"
	AssetBulkSetStatusURL = "/action/assets/bulk-set-status"
)
