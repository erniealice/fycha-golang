package reports

import (
	"context"
	"database/sql"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NewReceivablesAgingView creates the receivables aging report with DB data.
func NewReceivablesAgingView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "client",
		ActiveSubNav: "receivables-aging",
		Title:        "Receivables Aging",
		Subtitle:     "Aging analysis of outstanding receivables by customer",
		Icon:         "icon-file-text",
		TableID:      "receivables-aging-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			return fetchReceivablesAging(ctx, db)
		},
		BuildTotals: receivablesAgingTotals,
	})
}

// receivablesAgingTotals computes column totals for the receivables aging tfoot.
// Columns: Customer | Current | 1-30 | 31-60 | 61-90 | Over 90 | Total
func receivablesAgingTotals(rows []types.TableRow) []types.TableCell {
	if len(rows) == 0 {
		return nil
	}
	var current, d30, d60, d90, over90, total float64
	for _, row := range rows {
		if len(row.Cells) < 7 {
			continue
		}
		current += parseCurrency(row.Cells[1].Value)
		d30 += parseCurrency(row.Cells[2].Value)
		d60 += parseCurrency(row.Cells[3].Value)
		d90 += parseCurrency(row.Cells[4].Value)
		over90 += parseCurrency(row.Cells[5].Value)
		total += parseCurrency(row.Cells[6].Value)
	}
	return []types.TableCell{
		{Value: "Total"},
		{Value: FormatCurrency(current), Align: "right"},
		{Value: FormatCurrency(d30), Align: "right"},
		{Value: FormatCurrency(d60), Align: "right"},
		{Value: FormatCurrency(d90), Align: "right"},
		{Value: FormatCurrency(over90), Align: "right"},
		{Value: FormatCurrency(total), Align: "right"},
	}
}

func fetchReceivablesAging(ctx context.Context, db *sql.DB) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "customer", Label: "Customer", Sortable: true},
		{Key: "current", Label: "Current", Sortable: true, Align: "right"},
		{Key: "days-30", Label: "1-30 Days", Sortable: true, Align: "right"},
		{Key: "days-60", Label: "31-60 Days", Sortable: true, Align: "right"},
		{Key: "days-90", Label: "61-90 Days", Sortable: true, Align: "right"},
		{Key: "over-90", Label: "Over 90 Days", Sortable: true, Align: "right"},
		{Key: "total", Label: "Total", Sortable: true, Align: "right"},
	}

	// Compute outstanding = total_amount - collected payments per revenue.
	// treasury_collection.amount is in centavos; revenue.total_amount is also in centavos.
	// due_date is stored as BIGINT epoch millis; fall back to date_created when absent.
	query := `
		WITH outstanding AS (
			SELECT
				r.id,
				COALESCE(NULLIF(TRIM(r.name), ''), 'Unknown') AS customer_name,
				r.total_amount - COALESCE(collected.total_collected, 0) AS outstanding_amount,
				CURRENT_DATE - COALESCE(TO_TIMESTAMP(r.due_date / 1000.0)::date, r.date_created::date) AS days_overdue
			FROM revenue r
			LEFT JOIN (
				SELECT c.revenue_id, SUM(c.amount) AS total_collected
				FROM treasury_collection c
				WHERE c.active = true AND c.status IN ('paid', 'completed')
				GROUP BY c.revenue_id
			) collected ON collected.revenue_id = r.id
			WHERE r.active = true
			  AND r.status NOT IN ('paid', 'cancelled')
			  AND r.total_amount - COALESCE(collected.total_collected, 0) > 0
		)
		SELECT
			customer_name,
			COALESCE(SUM(CASE WHEN days_overdue <= 0 THEN outstanding_amount ELSE 0 END), 0) AS current_amt,
			COALESCE(SUM(CASE WHEN days_overdue BETWEEN 1 AND 30 THEN outstanding_amount ELSE 0 END), 0) AS days_30,
			COALESCE(SUM(CASE WHEN days_overdue BETWEEN 31 AND 60 THEN outstanding_amount ELSE 0 END), 0) AS days_60,
			COALESCE(SUM(CASE WHEN days_overdue BETWEEN 61 AND 90 THEN outstanding_amount ELSE 0 END), 0) AS days_90,
			COALESCE(SUM(CASE WHEN days_overdue > 90 THEN outstanding_amount ELSE 0 END), 0) AS over_90,
			COALESCE(SUM(outstanding_amount), 0) AS total
		FROM outstanding
		GROUP BY customer_name
		ORDER BY total DESC
	`

	dbRows, err := db.QueryContext(ctx, query)
	if err != nil {
		return columns, nil, nil
	}
	defer dbRows.Close()

	var rows []types.TableRow
	idx := 0
	for dbRows.Next() {
		var name string
		var current, d30, d60, d90, over90, total int64
		if err := dbRows.Scan(&name, &current, &d30, &d60, &d90, &over90, &total); err != nil {
			continue
		}
		idx++
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("ra-%d", idx),
			Cells: []types.TableCell{
				{Value: name},
				{Value: FormatCurrency(float64(current) / 100)},
				{Value: FormatCurrency(float64(d30) / 100)},
				{Value: FormatCurrency(float64(d60) / 100)},
				{Value: FormatCurrency(float64(d90) / 100)},
				{Value: FormatCurrency(float64(over90) / 100)},
				{Value: FormatCurrency(float64(total) / 100)},
			},
		})
	}

	return columns, rows, nil
}
