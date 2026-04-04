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
	})
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

	// revenue stores customer name in the `name` column
	query := `
		SELECT
			COALESCE(NULLIF(TRIM(r.name), ''), 'Unknown') AS customer_name,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(TO_TIMESTAMP(r.due_date / 1000.0)::date, r.date_created::date) <= 0 THEN r.total_amount ELSE 0 END), 0) AS current_amt,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(TO_TIMESTAMP(r.due_date / 1000.0)::date, r.date_created::date) BETWEEN 1 AND 30 THEN r.total_amount ELSE 0 END), 0) AS days_30,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(TO_TIMESTAMP(r.due_date / 1000.0)::date, r.date_created::date) BETWEEN 31 AND 60 THEN r.total_amount ELSE 0 END), 0) AS days_60,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(TO_TIMESTAMP(r.due_date / 1000.0)::date, r.date_created::date) BETWEEN 61 AND 90 THEN r.total_amount ELSE 0 END), 0) AS days_90,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(TO_TIMESTAMP(r.due_date / 1000.0)::date, r.date_created::date) > 90 THEN r.total_amount ELSE 0 END), 0) AS over_90,
			COALESCE(SUM(r.total_amount), 0) AS total
		FROM revenue r
		WHERE r.status NOT IN ('paid', 'cancelled')
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
		var current, d30, d60, d90, over90, total float64
		if err := dbRows.Scan(&name, &current, &d30, &d60, &d90, &over90, &total); err != nil {
			continue
		}
		idx++
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("ra-%d", idx),
			Cells: []types.TableCell{
				{Value: name},
				{Value: FormatCurrency(current / 100)},
				{Value: FormatCurrency(d30 / 100)},
				{Value: FormatCurrency(d60 / 100)},
				{Value: FormatCurrency(d90 / 100)},
				{Value: FormatCurrency(over90 / 100)},
				{Value: FormatCurrency(total / 100)},
			},
		})
	}

	return columns, rows, nil
}
