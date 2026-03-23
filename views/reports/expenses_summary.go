package reports

import (
	"context"
	"database/sql"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NewExpensesSummaryView creates the expenses summary report with DB data.
func NewExpensesSummaryView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "expenses",
		ActiveSubNav: "expenses-summary",
		Title:        "Expenses Summary",
		Subtitle:     "Summary of expenses by category and period",
		Icon:         "icon-bar-chart",
		TableID:      "expenses-summary-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			return fetchExpensesSummary(ctx, db)
		},
	})
}

func fetchExpensesSummary(ctx context.Context, db *sql.DB) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "category", Label: "Category", Sortable: true},
		{Key: "count", Label: "Count", Sortable: true, Align: "right"},
		{Key: "total", Label: "Total Amount", Sortable: true, Align: "right"},
		{Key: "percent", Label: "% of Total", Sortable: true, Align: "right"},
	}

	query := `
		WITH expense_totals AS (
			SELECT
				COALESCE(ec.name, 'Uncategorized') AS category,
				COUNT(*) AS expense_count,
				COALESCE(SUM(e.total_amount), 0) AS total_amount
			FROM expenditure e
			LEFT JOIN expenditure_category ec ON e.expenditure_category_id = ec.id
			WHERE e.expenditure_type = 'expense'
			  AND e.status != 'cancelled'
			GROUP BY ec.name
		),
		grand_total AS (
			SELECT COALESCE(SUM(total_amount), 0) AS grand FROM expense_totals
		)
		SELECT
			et.category,
			et.expense_count,
			et.total_amount,
			CASE WHEN gt.grand > 0 THEN (et.total_amount::numeric / gt.grand::numeric * 100) ELSE 0 END AS pct
		FROM expense_totals et, grand_total gt
		ORDER BY et.total_amount DESC
	`

	dbRows, err := db.QueryContext(ctx, query)
	if err != nil {
		return columns, nil, nil
	}
	defer dbRows.Close()

	var rows []types.TableRow
	idx := 0
	for dbRows.Next() {
		var category string
		var count int
		var total int64
		var pct float64
		if err := dbRows.Scan(&category, &count, &total, &pct); err != nil {
			continue
		}
		idx++
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("es-%d", idx),
			Cells: []types.TableCell{
				{Value: category},
				{Value: fmt.Sprintf("%d", count)},
				{Value: FormatCurrency(float64(total) / 100)},
				{Value: fmt.Sprintf("%.1f%%", pct)},
			},
		})
	}

	return columns, rows, nil
}
