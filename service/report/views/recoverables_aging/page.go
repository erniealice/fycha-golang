// Package recoverables_aging renders the recoverables aging report: unpaid
// issued recovery documents by days past due, per client and currency.
package recoverables_aging

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

	recoverydocumentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/recovery_document"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	report "github.com/erniealice/fycha-golang/service/report"
)

const dateLayout = "2006-01-02"

// Deps holds view dependencies. ListRecoverablesAging is nil when the espyna
// use case is not wired (the unit refuses boot; the view also fails closed).
type Deps struct {
	Labels                report.RecoveryReportsLabels
	CommonLabels          pyeza.CommonLabels
	TableLabels           types.TableLabels
	Routes                report.RecoveryReportsRoutes
	ListRecoverablesAging func(context.Context, *recoverydocumentpb.ListRecoverablesAgingRequest) (*recoverydocumentpb.ListRecoverablesAgingResponse, error)
}

// PageData is the report page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	AsOfDate        string
	ClientID        string
	ExportURL       string
}

func filters(deps *Deps, viewCtx *view.ViewContext) (asOf, clientID string, ok bool) {
	asOf = viewCtx.QueryParams["as-of-date"]
	if asOf == "" {
		asOf = time.Now().Format(dateLayout)
	}
	if _, err := time.Parse(dateLayout, asOf); err != nil {
		return asOf, "", false
	}
	return asOf, viewCtx.QueryParams["client-id"], true
}

func request(asOf, clientID string) *recoverydocumentpb.ListRecoverablesAgingRequest {
	req := &recoverydocumentpb.ListRecoverablesAgingRequest{AsOf: &asOf}
	if clientID != "" {
		req.ClientId = &clientID
	}
	return req
}

// NewView creates the recoverables aging report view.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("recovery_document", "list") {
			return view.Forbidden("recovery_document:list")
		}
		l := deps.Labels.RecoverablesAging
		if deps.ListRecoverablesAging == nil {
			return view.Error(fmt.Errorf("%s", deps.Labels.Shared.Unavailable))
		}
		asOf, clientID, ok := filters(deps, viewCtx)
		if !ok {
			return view.Error(fmt.Errorf("%s", deps.Labels.Shared.InvalidFilter))
		}
		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = deps.Routes.RecoverablesAgingURL
		}

		// Filter sheet request.
		if viewCtx.QueryParams["sheet"] == "filters" {
			return view.OK("recovery-report-filter-sheet", &report.RecoveryFilterSheetData{
				FormID: "recoverablesAgingFilterForm", ReportURL: reportURL,
				ApplyLabel: deps.Labels.Shared.Apply, ClearLabel: deps.Labels.Shared.Clear,
				Fields: []report.RecoveryFilterField{
					{Type: "date", Name: "as-of-date", ID: "recoverables-aging-as-of-date", Label: l.FilterAsOfDate, Value: asOf, TestID: "recoverables-aging-as-of-date"},
					{Type: "text", Name: "client-id", ID: "recoverables-aging-client-id", Label: l.FilterClient, Value: clientID, Hint: l.FilterClientHint, TestID: "recoverables-aging-client-id"},
				},
			})
		}

		resp, err := deps.ListRecoverablesAging(ctx, request(asOf, clientID))
		if err != nil {
			log.Printf("recoverables_aging: %v", err)
			return view.Error(fmt.Errorf("%s", deps.Labels.Shared.Unavailable))
		}

		table := buildTable(deps, resp)
		activeCount := 0
		chips := []report.RecoveryToolbarChip{{Label: l.ChipAsOf, Value: asOf, TestID: "rr-chip-as-of-date"}}
		if asOf != time.Now().Format(dateLayout) {
			activeCount++
		}
		if clientID != "" {
			activeCount++
			chips = append(chips, report.RecoveryToolbarChip{Label: l.ChipClient, Value: clientID, TestID: "rr-chip-client"})
		}
		q := url.Values{"sheet": {"filters"}, "as-of-date": {asOf}}
		if clientID != "" {
			q.Set("client-id", clientID)
		}
		table.ToolbarPrefixTemplate = "report-recovery-toolbar-prefix"
		table.ToolbarPrefixData = report.RecoveryToolbarPrefixData{
			FilterSheetURL: reportURL + "?" + q.Encode(), FiltersLabel: deps.Labels.Shared.Filters,
			ActiveFilterCount: activeCount, Chips: chips,
		}

		ex := url.Values{"as-of-date": {asOf}}
		if clientID != "" {
			ex.Set("client-id", clientID)
		}
		pd := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.PageTitle,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "report",
				ActiveSubNav:   "recoverables-aging",
				HeaderTitle:    l.PageTitle,
				HeaderSubtitle: l.PageDescription,
				HeaderIcon:     "icon-clock",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "recoverables-aging-report-content",
			Table:           table,
			AsOfDate:        asOf,
			ClientID:        clientID,
			ExportURL:       deps.Routes.RecoverablesAgingExportURL + "?" + ex.Encode(),
		}
		if viewCtx.IsHTMX {
			return view.OK("recoverables-aging-report-content", pd)
		}
		return view.OK("recoverables-aging-report", pd)
	})
}

func buildTable(deps *Deps, resp *recoverydocumentpb.ListRecoverablesAgingResponse) *types.TableConfig {
	l := deps.Labels.RecoverablesAging
	columns := []types.TableColumn{
		{Key: "client", Label: deps.Labels.Shared.Client},
		{Key: "currency", Label: deps.Labels.Shared.Currency, WidthClass: "col-xl"},
		{Key: "days_0_30", Label: l.Bucket0To30, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_31_60", Label: l.Bucket31To60, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_61_90", Label: l.Bucket61To90, Align: "right", MinWidth: "7.5rem"},
		{Key: "days_over_90", Label: l.BucketOver90, Align: "right", MinWidth: "7.5rem"},
		{Key: "total", Label: l.TotalOutstanding, Align: "right", MinWidth: "8.125rem"},
		{Key: "document_count", Label: l.DocumentCount, Align: "right", MinWidth: "6rem"},
	}
	table := &types.TableConfig{
		ID:          "recoverablesAgingTable",
		Columns:     columns,
		ShowSearch:  false,
		ShowFilters: false,
		ShowSort:    false,
		ShowColumns: false,
		ShowExport:  false,
		ShowEntries: true,
		ShowDensity: true,
		Labels:      deps.TableLabels,
		EmptyState:  types.TableEmptyState{Title: l.EmptyTitle, Message: l.EmptyMessage},
	}
	toRow := func(id, name string, r *recoverydocumentpb.RecoverablesAgingRow, testid string) types.TableRow {
		cur := r.GetCurrency()
		return types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "name", Value: name},
				{Type: "text", Value: cur},
				types.MoneyCell(float64(r.GetDays_0_30()), cur, true),
				types.MoneyCell(float64(r.GetDays_31_60()), cur, true),
				types.MoneyCell(float64(r.GetDays_61_90()), cur, true),
				types.MoneyCell(float64(r.GetDaysOver_90()), cur, true),
				types.MoneyCell(float64(r.GetTotal()), cur, true),
				{Type: "text", Value: fmt.Sprintf("%d", r.GetDocumentCount())},
			},
			DataAttrs: map[string]string{
				"testid": testid, "total": fmt.Sprintf("%d", r.GetTotal()), "currency": cur,
			},
		}
	}
	var rows []types.TableRow
	for _, r := range resp.GetRows() {
		id := r.GetClientId() + "|" + r.GetCurrency()
		rows = append(rows, toRow(id, r.GetClientId(), r, "recoverables-aging-row-"+id))
	}
	// One totals row per currency (a workspace may bill in several).
	if len(rows) > 0 {
		for _, t := range resp.GetTotals() {
			rows = append(rows, toRow("total|"+t.GetCurrency(), deps.Labels.Shared.Total, t, "recoverables-aging-total-"+t.GetCurrency()))
		}
	}
	table.Rows = rows
	types.ApplyColumnStyles(columns, rows)
	types.ApplyTableSettings(table)
	return table
}

func centavos(n int64) string { return fmt.Sprintf("%.2f", float64(n)/100.0) }

// NewExportHandler serves the CSV export with the same filters as the page.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !view.GetUserPermissions(ctx).Can("recovery_document", "list") {
			http.Error(w, deps.CommonLabels.Errors.PermissionDenied, http.StatusForbidden)
			return
		}
		if deps.ListRecoverablesAging == nil {
			http.Error(w, deps.Labels.Shared.Unavailable, http.StatusInternalServerError)
			return
		}
		q := r.URL.Query()
		asOf := q.Get("as-of-date")
		if asOf == "" {
			asOf = time.Now().Format(dateLayout)
		}
		if _, err := time.Parse(dateLayout, asOf); err != nil {
			http.Error(w, deps.Labels.Shared.InvalidFilter, http.StatusBadRequest)
			return
		}
		resp, err := deps.ListRecoverablesAging(ctx, request(asOf, q.Get("client-id")))
		if err != nil {
			log.Printf("recoverables_aging export: %v", err)
			http.Error(w, deps.Labels.Shared.Unavailable, http.StatusInternalServerError)
			return
		}
		l := deps.Labels.RecoverablesAging
		filename := l.ExportFilename
		if filename == "" {
			filename = "recoverables-aging.csv"
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
		cw := csv.NewWriter(w)
		defer cw.Flush()
		_ = cw.Write([]string{deps.Labels.Shared.Client, deps.Labels.Shared.Currency, l.Bucket0To30, l.Bucket31To60, l.Bucket61To90, l.BucketOver90, l.TotalOutstanding, l.DocumentCount})
		write := func(name string, x *recoverydocumentpb.RecoverablesAgingRow) {
			_ = cw.Write([]string{name, x.GetCurrency(), centavos(x.GetDays_0_30()), centavos(x.GetDays_31_60()), centavos(x.GetDays_61_90()),
				centavos(x.GetDaysOver_90()), centavos(x.GetTotal()), fmt.Sprintf("%d", x.GetDocumentCount())})
		}
		for _, x := range resp.GetRows() {
			write(x.GetClientId(), x)
		}
		if len(resp.GetRows()) > 0 {
			for _, t := range resp.GetTotals() {
				write(deps.Labels.Shared.Total, t)
			}
		}
	}
}
