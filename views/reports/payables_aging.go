package reports

import (
	"context"
	"database/sql"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NewPayablesAgingView creates the payables aging report with DB data.
func NewPayablesAgingView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "supplier",
		ActiveSubNav: "payables-aging",
		Title:        "Payables Aging",
		Subtitle:     "Aging analysis of outstanding payables by supplier",
		Icon:         "icon-file-text",
		TableID:      "payables-aging-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			return fetchPayablesAging(ctx, db)
		},
	})
}

func fetchPayablesAging(ctx context.Context, db *sql.DB) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "supplier", Label: "Supplier", Sortable: true},
		{Key: "current", Label: "Current", Sortable: true, Align: "right"},
		{Key: "days-30", Label: "1-30 Days", Sortable: true, Align: "right"},
		{Key: "days-60", Label: "31-60 Days", Sortable: true, Align: "right"},
		{Key: "days-90", Label: "61-90 Days", Sortable: true, Align: "right"},
		{Key: "over-90", Label: "Over 90 Days", Sortable: true, Align: "right"},
		{Key: "total", Label: "Total", Sortable: true, Align: "right"},
	}

	// vendor_name is stored directly on expenditure; fall back to supplier company_name via vendor_id
	query := `
		SELECT
			COALESCE(NULLIF(e.vendor_name, ''), s.company_name, 'Unknown') AS supplier_name,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(e.due_date, e.expenditure_date)::date <= 0 THEN e.total_amount ELSE 0 END), 0) AS current_amt,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(e.due_date, e.expenditure_date)::date BETWEEN 1 AND 30 THEN e.total_amount ELSE 0 END), 0) AS days_30,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(e.due_date, e.expenditure_date)::date BETWEEN 31 AND 60 THEN e.total_amount ELSE 0 END), 0) AS days_60,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(e.due_date, e.expenditure_date)::date BETWEEN 61 AND 90 THEN e.total_amount ELSE 0 END), 0) AS days_90,
			COALESCE(SUM(CASE WHEN CURRENT_DATE - COALESCE(e.due_date, e.expenditure_date)::date > 90 THEN e.total_amount ELSE 0 END), 0) AS over_90,
			COALESCE(SUM(e.total_amount), 0) AS total
		FROM expenditure e
		LEFT JOIN supplier s ON e.vendor_id = s.id
		WHERE e.expenditure_type = 'purchase'
		  AND e.status NOT IN ('paid', 'cancelled')
		GROUP BY supplier_name
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
			ID: fmt.Sprintf("pa-%d", idx),
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
