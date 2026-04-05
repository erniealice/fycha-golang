package expenditure_report

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	expreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/expenditure_report"
)

// NewExportHandler returns an http.HandlerFunc that exports the expenditure report as CSV.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := deps.Labels.ExpenditureReport

		// Parse query params (same as page.go)
		q := r.URL.Query()
		primary := q.Get("primary")
		if primary == "" {
			primary = "monthly"
		}
		rows := q.Get("rows")
		if rows == "" {
			rows = "category"
		}
		period := q.Get("period")
		if period == "" {
			period = "thisMonth"
		}
		startDateStr := q.Get("start")
		endDateStr := q.Get("end")

		// Secondary filters
		productID := q.Get("product-id")
		locationID := q.Get("location-id")
		locationAreaID := q.Get("location-area-id")
		expenditureCategoryID := q.Get("expenditure-category-id")
		supplierID := q.Get("supplier-id")
		expenditureType := q.Get("expenditure-type")

		// Build proto request
		req := &expreportpb.ExpenditureReportRequest{
			PrimaryDimension: primary,
			RowDimension:     rows,
		}

		if period == "custom" && startDateStr != "" {
			if _, err := time.Parse("2006-01-02", startDateStr); err == nil {
				req.StartDate = &startDateStr
			}
		}
		if period == "custom" && endDateStr != "" {
			if _, err := time.Parse("2006-01-02", endDateStr); err == nil {
				req.EndDate = &endDateStr
			}
		}

		if req.StartDate == nil {
			start, _ := fycha.ParsePeriodPreset(period)
			s := start.Format("2006-01-02")
			req.StartDate = &s
		}
		if req.EndDate == nil {
			_, end := fycha.ParsePeriodPreset(period)
			e := end.Format("2006-01-02")
			req.EndDate = &e
		}

		if productID != "" {
			req.ProductId = &productID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if locationAreaID != "" {
			req.LocationAreaId = &locationAreaID
		}
		if expenditureCategoryID != "" {
			req.ExpenditureCategoryId = &expenditureCategoryID
		}
		if supplierID != "" {
			req.SupplierId = &supplierID
		}
		if expenditureType != "" {
			req.ExpenditureType = &expenditureType
		}

		resp, err := deps.DB.GetExpenditureReport(ctx, req)
		if err != nil {
			http.Error(w, "Failed to generate expenditure report", http.StatusInternalServerError)
			return
		}

		// CSV response
		filename := fmt.Sprintf("expenditure-report-%s.csv", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Header row
		columnKeys := resp.GetColumnKeys()
		header := make([]string, 0, len(columnKeys)+2)
		header = append(header, l.PrimaryGroupLabel(rows))
		header = append(header, columnKeys...)
		header = append(header, "Total")
		_ = writer.Write(header)

		// Data rows
		for _, row := range resp.GetRows() {
			cellMap := make(map[string]*expreportpb.ExpenditureReportCell, len(row.GetCells()))
			for _, c := range row.GetCells() {
				cellMap[c.GetColumnKey()] = c
			}

			record := make([]string, 0, len(columnKeys)+2)
			record = append(record, row.GetRowKey())
			for _, ck := range columnKeys {
				var val int64
				if c, ok := cellMap[ck]; ok {
					val = c.GetTotalExpenditure()
				}
				record = append(record, csvCurrency(val))
			}
			record = append(record, csvCurrency(row.GetRowTotal()))
			_ = writer.Write(record)
		}

		// Totals row
		summary := resp.GetSummary()
		if summary != nil && len(resp.GetRows()) > 0 {
			colTotalMap := make(map[string]*expreportpb.ExpenditureReportCell, len(summary.GetColumnTotals()))
			for _, ct := range summary.GetColumnTotals() {
				colTotalMap[ct.GetColumnKey()] = ct
			}
			totalsRecord := make([]string, 0, len(columnKeys)+2)
			totalsRecord = append(totalsRecord, "TOTAL")
			for _, ck := range columnKeys {
				var val int64
				if ct, ok := colTotalMap[ck]; ok {
					val = ct.GetTotalExpenditure()
				}
				totalsRecord = append(totalsRecord, csvCurrency(val))
			}
			totalsRecord = append(totalsRecord, csvCurrency(summary.GetGrandTotal()))
			_ = writer.Write(totalsRecord)
		}
	}
}

// csvCurrency formats centavos as a plain decimal string (no symbol, no commas).
func csvCurrency(centavos int64) string {
	return fmt.Sprintf("%.2f", float64(centavos)/100.0)
}
