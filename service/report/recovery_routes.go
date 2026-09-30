package report

// Recovery report routes (build-spec §6.5): the recoverables aging report and
// the cost-source reconciliation report. They are mounted by fycha's opt-in
// RecoveryReportsUnit (block.WithRecoveryReports), not by the default reports
// module, so apps that do not opt in keep byte-identical compositions.
const (
	ReportsRecoverablesAgingURL              = "/reports/recoverables-aging"
	ReportsRecoverablesAgingExportURL        = "/reports/recoverables-aging/export"
	ReportsCostSourceReconciliationURL       = "/reports/cost-source-reconciliation"
	ReportsCostSourceReconciliationExportURL = "/reports/cost-source-reconciliation/export"
)

// RecoveryReportsRoutes holds the recovery report routes. JSON tags are the
// lyngua route.json keys (general/route.json carries "recovery_reports": {} —
// defaults live here).
type RecoveryReportsRoutes struct {
	RecoverablesAgingURL              string `json:"recoverables_aging_url"`
	RecoverablesAgingExportURL        string `json:"recoverables_aging_export_url"`
	CostSourceReconciliationURL       string `json:"cost_source_reconciliation_url"`
	CostSourceReconciliationExportURL string `json:"cost_source_reconciliation_export_url"`
}

// DefaultRecoveryReportsRoutes returns the routes populated from the constants.
func DefaultRecoveryReportsRoutes() RecoveryReportsRoutes {
	return RecoveryReportsRoutes{
		RecoverablesAgingURL:              ReportsRecoverablesAgingURL,
		RecoverablesAgingExportURL:        ReportsRecoverablesAgingExportURL,
		CostSourceReconciliationURL:       ReportsCostSourceReconciliationURL,
		CostSourceReconciliationExportURL: ReportsCostSourceReconciliationExportURL,
	}
}

// RouteMap returns dot-notation route keys.
func (r RecoveryReportsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"reports.recoverables_aging":                r.RecoverablesAgingURL,
		"reports.recoverables_aging_export":         r.RecoverablesAgingExportURL,
		"reports.cost_source_reconciliation":        r.CostSourceReconciliationURL,
		"reports.cost_source_reconciliation_export": r.CostSourceReconciliationExportURL,
	}
}
