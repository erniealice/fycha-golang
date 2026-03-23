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

// NewSalesSummaryView creates the sales summary report with DB data.
func NewSalesSummaryView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "sales",
		ActiveSubNav: "sales-summary",
		Title:        "Sales Summary",
		Subtitle:     "Summary of sales performance and trends",
		Icon:         "icon-bar-chart",
		TableID:      "sales-summary-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			return fetchSalesSummary(ctx, db)
		},
	})
}

func fetchSalesSummary(ctx context.Context, db *sql.DB) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "period", Label: "Period", Sortable: true},
		{Key: "transactions", Label: "Transactions", Sortable: true, Align: "right"},
		{Key: "gross-sales", Label: "Gross Sales", Sortable: true, Align: "right"},
		{Key: "net-sales", Label: "Net Sales", Sortable: true, Align: "right"},
	}

	query := `
		SELECT
			TO_CHAR(date_created, 'Month YYYY') AS period,
			COUNT(*) AS transactions,
			COALESCE(SUM(total_amount), 0) AS gross_sales,
			COALESCE(SUM(CASE WHEN status != 'cancelled' THEN total_amount ELSE 0 END), 0) AS net_sales
		FROM revenue
		WHERE date_created >= NOW() - INTERVAL '12 months'
		GROUP BY DATE_TRUNC('month', date_created), TO_CHAR(date_created, 'Month YYYY')
		ORDER BY DATE_TRUNC('month', date_created) DESC
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
		var txCount int
		var grossSales, netSales int64
		if err := dbRows.Scan(&period, &txCount, &grossSales, &netSales); err != nil {
			continue
		}
		idx++
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("ss-%d", idx),
			Cells: []types.TableCell{
				{Value: strings.TrimSpace(period)},
				{Value: fmt.Sprintf("%d", txCount)},
				{Value: FormatCurrency(float64(grossSales) / 100)},
				{Value: FormatCurrency(float64(netSales) / 100)},
			},
		})
	}

	return columns, rows, nil
}
