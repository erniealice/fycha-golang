// Package cost_source_reconciliation renders the cost-source reconciliation
// report for one expenditure: per recoverable cost component, the amount versus
// the sum of its allocation shares versus the charges issued from them.
// Mismatches are highlighted.
package cost_source_reconciliation

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	costsourcecomponentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/expenditure/cost_source_component"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	report "github.com/erniealice/fycha-golang/service/report"
)

// Deps holds view dependencies. ReconcileCostSource is nil when the espyna use
// case is not wired (the unit refuses boot; the view also fails closed).
type Deps struct {
	Labels              report.RecoveryReportsLabels
	CommonLabels        pyeza.CommonLabels
	TableLabels         types.TableLabels
	Routes              report.RecoveryReportsRoutes
	ReconcileCostSource func(context.Context, *costsourcecomponentpb.ReconcileCostSourceRequest) (*costsourcecomponentpb.ReconcileCostSourceResponse, error)
}

// PageData is the report page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	ExpenditureID   string
	ExportURL       string
}

// Status is the reconciliation state of one component.
type Status string

const (
	StatusReconciled Status = "reconciled"
	StatusMismatch   Status = "mismatch"
	StatusNoBatch    Status = "no_batch"
)

// StatusOf classifies a component: no published allocation, reconciled, or a
// mismatch (either variance non-zero).
func StatusOf(r *costsourcecomponentpb.CostSourceComponentReconciliation) Status {
	switch {
	case !r.GetHasPublishedBatch():
		return StatusNoBatch
	case r.GetReconciled():
		return StatusReconciled
	default:
		return StatusMismatch
	}
}

func statusLabel(l report.ReconciliationLabels, s Status) string {
	switch s {
	case StatusReconciled:
		return l.StatusReconciled
	case StatusMismatch:
		return l.StatusMismatch
	default:
		return l.StatusNoBatch
	}
}

// NewView creates the cost-source reconciliation report view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("expenditure", "read") {
			return view.Forbidden("expenditure:read")
		}
		l := deps.Labels.Reconciliation
		if deps.ReconcileCostSource == nil {
			return view.Error(fmt.Errorf("%s", deps.Labels.Shared.Unavailable))
		}
		expID := strings.TrimSpace(viewCtx.QueryParams["expenditure-id"])
		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = deps.Routes.CostSourceReconciliationURL
		}

		if viewCtx.QueryParams["sheet"] == "filters" {
			return view.OK("recovery-report-filter-sheet", &report.RecoveryFilterSheetData{
				FormID: "costSourceReconciliationFilterForm", ReportURL: reportURL,
				ApplyLabel: deps.Labels.Shared.Apply, ClearLabel: deps.Labels.Shared.Clear,
				Fields: []report.RecoveryFilterField{
					{Type: "text", Name: "expenditure-id", ID: "cost-source-reconciliation-expenditure-id", Label: l.FilterExpenditure,
						Value: expID, Hint: l.FilterExpenditureHint, TestID: "cost-source-reconciliation-expenditure-id"},
				},
			})
		}

		var resp *costsourcecomponentpb.ReconcileCostSourceResponse
		if expID != "" {
			var err error
			resp, err = deps.ReconcileCostSource(ctx, &costsourcecomponentpb.ReconcileCostSourceRequest{ExpenditureId: &expID})
			if err != nil {
				log.Printf("cost_source_reconciliation: %v", err)
				return view.Error(fmt.Errorf("%s", deps.Labels.Shared.Unavailable))
			}
		}

		table := buildTable(deps, resp, expID != "")
		chips := []report.RecoveryToolbarChip{}
		active := 0
		if expID != "" {
			active = 1
			chips = append(chips, report.RecoveryToolbarChip{Label: l.ChipExpenditure, Value: expID, TestID: "rr-chip-expenditure"})
		}
		q := url.Values{"sheet": {"filters"}}
		ex := url.Values{}
		if expID != "" {
			q.Set("expenditure-id", expID)
			ex.Set("expenditure-id", expID)
		}
		table.ToolbarPrefixTemplate = "report-recovery-toolbar-prefix"
		table.ToolbarPrefixData = report.RecoveryToolbarPrefixData{
			FilterSheetURL: reportURL + "?" + q.Encode(), FiltersLabel: deps.Labels.Shared.Filters,
			ActiveFilterCount: active, Chips: chips,
		}
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.PageTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "report",
				ActiveSubNav:   "cost-source-reconciliation",
				HeaderTitle:    l.PageTitle,
				HeaderSubtitle: l.PageDescription,
				HeaderIcon:     "icon-check-circle",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "cost-source-reconciliation-report-content",
			Table:           table,
			ExpenditureID:   expID,
			ExportURL:       deps.Routes.CostSourceReconciliationExportURL + "?" + ex.Encode(),
		}
		if viewCtx.IsHTMX {
			return view.OK("cost-source-reconciliation-report-content", pd)
		}
		return view.OK("cost-source-reconciliation-report", pd)
	})
}

// varianceCell renders a difference: a plain money cell at zero, a danger badge
// (the highlighted mismatch) otherwise.
func varianceCell(n int64, currency string) types.TableCell {
	m := types.MoneyCell(float64(n), currency, true)
	if n == 0 {
		return m
	}
	return types.TableCell{Type: "badge", Value: m.Currency + " " + m.Value, Variant: "danger"}
}

func buildTable(deps *Deps, resp *costsourcecomponentpb.ReconcileCostSourceResponse, chosen bool) *types.TableConfig {
	l := deps.Labels.Reconciliation
	columns := []types.TableColumn{
		{Key: "component", Label: l.Component},
		{Key: "amount", Label: l.Amount, Align: "right", MinWidth: "7.5rem"},
		{Key: "shares_total", Label: l.SharesTotal, Align: "right", MinWidth: "7.5rem"},
		{Key: "recoverable_total", Label: l.RecoverableTotal, Align: "right", MinWidth: "7.5rem"},
		{Key: "share_variance", Label: l.ShareVariance, Align: "right", MinWidth: "7.5rem"},
		{Key: "charges_open", Label: l.ChargesOpen, Align: "right", MinWidth: "7.5rem"},
		{Key: "charges_issued", Label: l.ChargesIssued, Align: "right", MinWidth: "7.5rem"},
		{Key: "charges_cancelled", Label: l.ChargesCancelled, Align: "right", MinWidth: "7.5rem"},
		{Key: "charge_variance", Label: l.ChargeVariance, Align: "right", MinWidth: "7.5rem"},
		{Key: "corrections", Label: l.Corrections, Align: "right", MinWidth: "7.5rem"},
		{Key: "issued_net", Label: l.IssuedNet, Align: "right", MinWidth: "7.5rem"},
		{Key: "status", Label: l.Status, WidthClass: "col-2xl"},
	}
	empty := types.TableEmptyState{Title: l.EmptyTitle, Message: l.EmptyMessage}
	if !chosen {
		empty = types.TableEmptyState{Title: l.PromptTitle, Message: l.PromptMessage}
	}
	table := &types.TableConfig{
		ID:          "costSourceReconciliationTable",
		Columns:     columns,
		ShowSearch:  false,
		ShowFilters: false,
		ShowSort:    false,
		ShowColumns: false,
		ShowExport:  false,
		ShowEntries: true,
		ShowDensity: true,
		Labels:      deps.TableLabels,
		EmptyState:  empty,
	}
	var rows []types.TableRow
	for _, r := range resp.GetData() {
		cur := r.GetCurrency()
		st := StatusOf(r)
		variant := map[Status]string{StatusReconciled: "success", StatusMismatch: "danger", StatusNoBatch: "warning"}[st]
		rows = append(rows, types.TableRow{
			ID: r.GetComponentId(),
			Cells: []types.TableCell{
				{Type: "name", Value: r.GetComponentId()},
				types.MoneyCell(float64(r.GetAmount()), cur, true),
				types.MoneyCell(float64(r.GetSharesTotal()), cur, true),
				types.MoneyCell(float64(r.GetRecoverableTotal()), cur, true),
				varianceCell(r.GetShareVariance(), cur),
				types.MoneyCell(float64(r.GetChargesOpen()), cur, true),
				types.MoneyCell(float64(r.GetChargesIssued()), cur, true),
				types.MoneyCell(float64(r.GetChargesCancelled()), cur, true),
				varianceCell(r.GetChargeVariance(), cur),
				types.MoneyCell(float64(r.GetCorrections()), cur, true),
				types.MoneyCell(float64(r.GetIssuedNet()), cur, true),
				{Type: "badge", Value: statusLabel(l, st), Variant: variant},
			},
			DataAttrs: map[string]string{
				"testid": "cost-source-reconciliation-row-" + r.GetComponentId(),
				"status": string(st),
			},
		})
	}
	table.Rows = rows
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(table)
	return table
}

func centavos(n int64) string { return fmt.Sprintf("%.2f", float64(n)/100.0) }

// NewExportHandler serves the CSV export with the same filter as the page.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !view.GetUserPermissions(ctx).Can("expenditure", "read") {
			http.Error(w, deps.CommonLabels.Errors.PermissionDenied, http.StatusForbidden)
			return
		}
		expID := strings.TrimSpace(r.URL.Query().Get("expenditure-id"))
		if deps.ReconcileCostSource == nil || expID == "" {
			http.Error(w, deps.Labels.Shared.InvalidFilter, http.StatusBadRequest)
			return
		}
		resp, err := deps.ReconcileCostSource(ctx, &costsourcecomponentpb.ReconcileCostSourceRequest{ExpenditureId: &expID})
		if err != nil {
			log.Printf("cost_source_reconciliation export: %v", err)
			http.Error(w, deps.Labels.Shared.Unavailable, http.StatusInternalServerError)
			return
		}
		l := deps.Labels.Reconciliation
		filename := l.ExportFilename
		if filename == "" {
			filename = "cost-reconciliation.csv"
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		cw := csv.NewWriter(w)
		defer cw.Flush()
		_ = cw.Write([]string{l.Component, deps.Labels.Shared.Currency, l.Amount, l.SharesTotal, l.RecoverableTotal, l.ShareVariance,
			l.ChargesOpen, l.ChargesIssued, l.ChargesCancelled, l.ChargeVariance, l.Corrections, l.IssuedNet, l.Status})
		for _, x := range resp.GetData() {
			_ = cw.Write([]string{x.GetComponentId(), x.GetCurrency(), centavos(x.GetAmount()), centavos(x.GetSharesTotal()),
				centavos(x.GetRecoverableTotal()), centavos(x.GetShareVariance()), centavos(x.GetChargesOpen()), centavos(x.GetChargesIssued()),
				centavos(x.GetChargesCancelled()), centavos(x.GetChargeVariance()), centavos(x.GetCorrections()), centavos(x.GetIssuedNet()),
				statusLabel(l, StatusOf(x))})
		}
	}
}
