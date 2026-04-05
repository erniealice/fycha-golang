package disbursement_report

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	disbreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/reporting/disbursement_report"
)

// NewExportHandler returns an http.HandlerFunc that exports the disbursement report as CSV.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		l := deps.Labels.DisbursementReport

		// Parse query params (same as page.go)
		q := r.URL.Query()
		primary := q.Get("primary")
		if primary == "" {
			primary = "monthly"
		}
		rows := q.Get("rows")
		if rows == "" {
			rows = "supplier"
		}
		period := q.Get("period")
		if period == "" {
			period = "thisMonth"
		}
		startDateStr := q.Get("start")
		endDateStr := q.Get("end")

		// Secondary filters
		supplierID := q.Get("supplier-id")
		supplierCategoryID := q.Get("supplier-category-id")
		locationID := q.Get("location-id")
		expenditureCategoryID := q.Get("expenditure-category-id")
		disbursementType := q.Get("disbursement-type")
		disbursementMethodID := q.Get("disbursement-method-id")

		// Build proto request
		req := &disbreportpb.DisbursementReportRequest{
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

		if supplierID != "" {
			req.SupplierId = &supplierID
		}
		if supplierCategoryID != "" {
			req.SupplierCategoryId = &supplierCategoryID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if expenditureCategoryID != "" {
			req.ExpenditureCategoryId = &expenditureCategoryID
		}
		if disbursementType != "" {
			req.DisbursementType = &disbursementType
		}
		if disbursementMethodID != "" {
			req.DisbursementMethodId = &disbursementMethodID
		}

		resp, err := deps.DB.GetDisbursementReport(ctx, req)
		if err != nil {
			http.Error(w, "Failed to generate disbursement report", http.StatusInternalServerError)
			return
		}

		// CSV response
		filename := fmt.Sprintf("disbursement-report-%s.csv", time.Now().Format("2006-01-02"))
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
			cellMap := make(map[string]*disbreportpb.DisbursementReportCell, len(row.GetCells()))
			for _, c := range row.GetCells() {
				cellMap[c.GetColumnKey()] = c
			}

			record := make([]string, 0, len(columnKeys)+2)
			record = append(record, row.GetRowKey())
			for _, ck := range columnKeys {
				var val int64
				if c, ok := cellMap[ck]; ok {
					val = c.GetTotalDisbursement()
				}
				record = append(record, csvCurrency(val))
			}
			record = append(record, csvCurrency(row.GetRowTotal()))
			_ = writer.Write(record)
		}

		// Totals row
		summary := resp.GetSummary()
		if summary != nil && len(resp.GetRows()) > 0 {
			colTotalMap := make(map[string]*disbreportpb.DisbursementReportCell, len(summary.GetColumnTotals()))
			for _, ct := range summary.GetColumnTotals() {
				colTotalMap[ct.GetColumnKey()] = ct
			}
			totalsRecord := make([]string, 0, len(columnKeys)+2)
			totalsRecord = append(totalsRecord, "TOTAL")
			for _, ck := range columnKeys {
				var val int64
				if ct, ok := colTotalMap[ck]; ok {
					val = ct.GetTotalDisbursement()
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
