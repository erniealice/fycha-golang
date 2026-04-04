package reports

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NewPurchasesSummaryView creates the purchases summary report with DB data.
func NewPurchasesSummaryView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "purchase",
		ActiveSubNav: "purchases-summary",
		Title:        "Purchases Summary",
		Subtitle:     "Summary of purchase orders and spending",
		Icon:         "icon-bar-chart",
		TableID:      "purchases-summary-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			return fetchPurchasesSummary(ctx, db)
		},
	})
}

func fetchPurchasesSummary(ctx context.Context, db *sql.DB) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "period", Label: "Period", Sortable: true},
		{Key: "orders", Label: "Orders", Sortable: true, Align: "right"},
		{Key: "total", Label: "Total Amount", Sortable: true, Align: "right"},
		{Key: "paid", Label: "Paid", Sortable: true, Align: "right"},
		{Key: "outstanding", Label: "Outstanding", Sortable: true, Align: "right"},
	}

	query := `
		SELECT
			TO_CHAR(expenditure_date, 'Month YYYY') AS period,
			COUNT(*) AS orders,
			COALESCE(SUM(total_amount), 0) AS total_amount,
			COALESCE(SUM(CASE WHEN status = 'paid' THEN total_amount ELSE 0 END), 0) AS paid,
			COALESCE(SUM(CASE WHEN status NOT IN ('paid', 'cancelled') THEN total_amount ELSE 0 END), 0) AS outstanding
		FROM expenditure
		WHERE expenditure_type = 'purchase'
		  AND expenditure_date >= NOW() - INTERVAL '12 months'
		GROUP BY DATE_TRUNC('month', expenditure_date), TO_CHAR(expenditure_date, 'Month YYYY')
		ORDER BY DATE_TRUNC('month', expenditure_date) DESC
	`

	dbRows, err := db.QueryContext(ctx, query)
	if err != nil {
		return columns, nil, nil
	}
	defer dbRows.Close()

	var rows []types.TableRow
	idx := 0
	for dbRows.Next() {
		var period string
		var orderCount int
		var total, paid, outstanding int64
		if err := dbRows.Scan(&period, &orderCount, &total, &paid, &outstanding); err != nil {
			continue
		}
		idx++
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("ps-%d", idx),
			Cells: []types.TableCell{
				{Value: strings.TrimSpace(period)},
				{Value: fmt.Sprintf("%d", orderCount)},
				{Value: FormatCurrency(float64(total) / 100)},
				{Value: FormatCurrency(float64(paid) / 100)},
				{Value: FormatCurrency(float64(outstanding) / 100)},
			},
		})
	}

	return columns, rows, nil
}
