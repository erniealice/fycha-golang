package revenue_report

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"time"

	fycha "github.com/erniealice/fycha-golang"

	revreportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/revenue_report"
)

// NewExportHandler creates an http.HandlerFunc for CSV export of the revenue report.
// It applies the same filters as the page view.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		// Parse the same query params as the page view
		primary := q.Get("primary")
		if primary == "" {
			primary = "monthly"
		}
		rows := q.Get("rows")
		if rows == "" {
			rows = "product"
		}
		period := q.Get("period")
		if period == "" {
			period = "thisMonth"
		}
		startDateStr := q.Get("start")
		endDateStr := q.Get("end")

		// Secondary filter IDs
		productID := q.Get("product-id")
		collectionID := q.Get("collection-id")
		locationID := q.Get("location-id")
		locationAreaID := q.Get("location-area-id")
		revenueCategoryID := q.Get("revenue-category-id")

		// Build proto request
		req := &revreportpb.RevenueReportRequest{
			PrimaryDimension: primary,
			RowDimension:     rows,
		}

		// Resolve dates from custom range or period preset
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

		// Apply optional secondary filters
		if productID != "" {
			req.ProductId = &productID
		}
		if collectionID != "" {
			req.CollectionId = &collectionID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if locationAreaID != "" {
			req.LocationAreaId = &locationAreaID
		}
		if revenueCategoryID != "" {
			req.RevenueCategoryId = &revenueCategoryID
		}

		// Call data source
		resp, err := deps.DB.GetRevenueReport(ctx, req)
		if err != nil {
			log.Printf("revenue_report export: failed to get revenue report: %v", err)
			http.Error(w, "Failed to generate report", http.StatusInternalServerError)
			return
		}
		if resp == nil {
			resp = &revreportpb.RevenueReportResponse{
				ColumnKeys: []string{},
				Rows:       []*revreportpb.RevenueReportRow{},
				Summary:    &revreportpb.RevenueReportSummary{},
			}
		}

		// Set CSV response headers
		filename := fmt.Sprintf("revenue-report-%s.csv", time.Now().Format("2006-01-02"))
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		writer := csv.NewWriter(w)
		defer writer.Flush()

		columnKeys := resp.GetColumnKeys()
		l := deps.Labels.RevenueReport

		// Write header row: row dimension label, then each column_key, then "Total"
		header := make([]string, 0, len(columnKeys)+2)
		header = append(header, l.PrimaryGroupLabel(rows))
		header = append(header, columnKeys...)
		header = append(header, l.Total)
		if err := writer.Write(header); err != nil {
			log.Printf("revenue_report export: failed to write CSV header: %v", err)
			return
		}

		// Write data rows
		for _, row := range resp.GetRows() {
			// Build cell map for quick lookup by column key
			cellMap := make(map[string]*revreportpb.RevenueReportCell, len(row.GetCells()))
			for _, c := range row.GetCells() {
				cellMap[c.GetColumnKey()] = c
			}

			record := make([]string, 0, len(columnKeys)+2)
			record = append(record, row.GetRowKey())
			for _, ck := range columnKeys {
				var val int64
				if c, ok := cellMap[ck]; ok {
					val = c.GetTotalRevenue()
				}
				record = append(record, csvCurrency(val))
			}
			record = append(record, csvCurrency(row.GetRowTotal()))

			if err := writer.Write(record); err != nil {
				log.Printf("revenue_report export: failed to write CSV row: %v", err)
				return
			}
		}

		// Write totals row: "TOTAL", then column totals, then grand_total
		summary := resp.GetSummary()
		if summary != nil && len(resp.GetRows()) > 0 {
			colTotalMap := make(map[string]*revreportpb.RevenueReportCell, len(summary.GetColumnTotals()))
			for _, ct := range summary.GetColumnTotals() {
				colTotalMap[ct.GetColumnKey()] = ct
			}

			totalsRecord := make([]string, 0, len(columnKeys)+2)
			totalsRecord = append(totalsRecord, "TOTAL")
			for _, ck := range columnKeys {
				var val int64
				if ct, ok := colTotalMap[ck]; ok {
					val = ct.GetTotalRevenue()
				}
				totalsRecord = append(totalsRecord, csvCurrency(val))
			}
			totalsRecord = append(totalsRecord, csvCurrency(summary.GetGrandTotal()))

			if err := writer.Write(totalsRecord); err != nil {
				log.Printf("revenue_report export: failed to write CSV totals row: %v", err)
				return
			}
		}
	}
}

// csvCurrency formats a centavo integer as a plain decimal string (e.g. "15000.50").
// No commas, no currency symbol — safe for CSV consumption.
func csvCurrency(centavos int64) string {
	return fmt.Sprintf("%.2f", float64(centavos)/100.0)
}
