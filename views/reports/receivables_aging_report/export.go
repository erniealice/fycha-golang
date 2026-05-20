package receivables_aging_report

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"time"

	aragingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ar_aging"
)

// NewExportHandler creates an http.HandlerFunc for CSV export of the receivables aging report.
// It applies the same filters as the page view.
func NewExportHandler(deps *Deps) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		q := r.URL.Query()

		// Parse the same query params as the page view
		asOfDate := q.Get("as-of-date")
		if asOfDate == "" {
			asOfDate = time.Now().Format("2006-01-02")
		}
		rows := q.Get("rows")
		if rows == "" {
			rows = "client"
		}

		// Secondary filter IDs
		clientID := q.Get("client-id")
		locationID := q.Get("location-id")
		revenueCategoryID := q.Get("revenue-category-id")

		// Build proto request
		req := &aragingpb.GetReceivablesAgingRequest{
			AsOfDate:     &asOfDate,
			RowDimension: rows,
		}

		// Apply optional secondary filters
		if clientID != "" {
			req.ClientId = &clientID
		}
		if locationID != "" {
			req.LocationId = &locationID
		}
		if revenueCategoryID != "" {
			req.RevenueCategoryId = &revenueCategoryID
		}

		// Call service-driven AR aging use case (Wave B P1.E.1).
		if deps.GetReceivablesAgingReport == nil {
			http.Error(w, "Failed to generate report", http.StatusInternalServerError)
			return
		}
		resp, err := deps.GetReceivablesAgingReport(ctx, req)
		if err != nil {
			log.Printf("receivables_aging_report export: failed to get report: %v", err)
			http.Error(w, "Failed to generate report", http.StatusInternalServerError)
			return
		}
		if resp == nil {
			resp = &aragingpb.GetReceivablesAgingResponse{
				BucketLabels: []string{},
				Rows:         []*aragingpb.ReceivablesAgingRow{},
				Summary:      &aragingpb.ReceivablesAgingSummary{},
			}
		}

		// Set CSV response headers
		l := deps.Labels.ReceivablesAging
		filename := l.ExportFilename
		if filename == "" {
			filename = fmt.Sprintf("receivables-aging-%s.csv", asOfDate)
		}
		w.Header().Set("Content-Type", "text/csv; charset=utf-8")
		w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))

		writer := csv.NewWriter(w)
		defer writer.Flush()

		// Resolve bucket labels: prefer response labels, fall back to static labels
		bucketLabels := resp.GetBucketLabels()
		var bCurrent, b1to30, b31to60, b61to90, bOver90 string
		if len(bucketLabels) >= 5 {
			bCurrent = bucketLabels[0]
			b1to30 = bucketLabels[1]
			b31to60 = bucketLabels[2]
			b61to90 = bucketLabels[3]
			bOver90 = bucketLabels[4]
		} else {
			bCurrent = "Current"
			b1to30 = "1-30 Days"
			b31to60 = "31-60 Days"
			b61to90 = "61-90 Days"
			bOver90 = "Over 90 Days"
		}

		// Write header row
		header := []string{
			rowDimensionLabel(l, rows),
			bCurrent,
			b1to30,
			b31to60,
			b61to90,
			bOver90,
			"Total Outstanding",
			"Invoice Count",
		}
		if err := writer.Write(header); err != nil {
			log.Printf("receivables_aging_report export: failed to write CSV header: %v", err)
			return
		}

		// Write data rows
		for _, row := range resp.GetRows() {
			b := row.GetBuckets()
			if b == nil {
				b = &aragingpb.AgingBuckets{}
			}

			record := []string{
				row.GetRowKey(),
				csvCurrency(b.GetCurrent()),
				csvCurrency(b.GetDays_1_30()),
				csvCurrency(b.GetDays_31_60()),
				csvCurrency(b.GetDays_61_90()),
				csvCurrency(b.GetDaysOver_90()),
				csvCurrency(row.GetTotalOutstanding()),
				fmt.Sprintf("%d", row.GetInvoiceCount()),
			}

			if err := writer.Write(record); err != nil {
				log.Printf("receivables_aging_report export: failed to write CSV row: %v", err)
				return
			}
		}

		// Write totals row
		summary := resp.GetSummary()
		if summary != nil && len(resp.GetRows()) > 0 {
			sb := summary.GetBuckets()
			if sb == nil {
				sb = &aragingpb.AgingBuckets{}
			}

			totalsRecord := []string{
				"TOTAL",
				csvCurrency(sb.GetCurrent()),
				csvCurrency(sb.GetDays_1_30()),
				csvCurrency(sb.GetDays_31_60()),
				csvCurrency(sb.GetDays_61_90()),
				csvCurrency(sb.GetDaysOver_90()),
				csvCurrency(summary.GetGrandTotalOutstanding()),
				fmt.Sprintf("%d", summary.GetTotalInvoiceCount()),
			}

			if err := writer.Write(totalsRecord); err != nil {
				log.Printf("receivables_aging_report export: failed to write CSV totals row: %v", err)
				return
			}
		}
	}
}

// csvCurrency formats a centavo integer as a plain decimal string (e.g. "15000.50").
// No commas, no currency symbol -- safe for CSV consumption.
func csvCurrency(centavos int64) string {
	return fmt.Sprintf("%.2f", float64(centavos)/100.0)
}
