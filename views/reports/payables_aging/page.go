package payables_aging

import (
	"context"
	"database/sql"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	reports "github.com/erniealice/fycha-golang/views/reports"
)

// NewPayablesAgingView creates the payables aging report with DB data.
func NewPayablesAgingView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return reports.NewReportView(reports.ReportConfig{
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
		BuildTotals: payablesAgingTotals,
	})
}

// payablesAgingTotals computes column totals for the payables aging tfoot.
// Columns: Supplier | Current | 1-30 | 31-60 | 61-90 | Over 90 | Total
func payablesAgingTotals(rows []types.TableRow) []types.TableCell {
	if len(rows) == 0 {
		return nil
	}
	// Parse and sum the 6 numeric columns (indices 1-6); index 0 is supplier name.
	var current, d30, d60, d90, over90, total float64
	for _, row := range rows {
		if len(row.Cells) < 7 {
			continue
		}
		current += reports.ParseCurrency(row.Cells[1].Value)
		d30 += reports.ParseCurrency(row.Cells[2].Value)
		d60 += reports.ParseCurrency(row.Cells[3].Value)
		d90 += reports.ParseCurrency(row.Cells[4].Value)
		over90 += reports.ParseCurrency(row.Cells[5].Value)
		total += reports.ParseCurrency(row.Cells[6].Value)
	}
	return []types.TableCell{
		{Value: "Total"},
		{Value: reports.FormatCurrency(current), Align: "right"},
		{Value: reports.FormatCurrency(d30), Align: "right"},
		{Value: reports.FormatCurrency(d60), Align: "right"},
		{Value: reports.FormatCurrency(d90), Align: "right"},
		{Value: reports.FormatCurrency(over90), Align: "right"},
		{Value: reports.FormatCurrency(total), Align: "right"},
	}
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

	// Compute outstanding = total_amount - paid disbursements per expenditure.
	// Use COALESCE(supplier_id, vendor_id) to support both old and new P2P schema.
	// Include both 'purchase' and 'expense' types that create payables.
	query := `
		WITH outstanding AS (
			SELECT
				e.id,
				COALESCE(NULLIF(TRIM(s.company_name), ''), NULLIF(TRIM(e.name), ''), 'Unknown') AS supplier_name,
				e.total_amount - COALESCE(paid.total_paid, 0) AS outstanding_amount,
				CURRENT_DATE - COALESCE(e.due_date, e.expenditure_date)::date AS days_overdue
			FROM expenditure e
			LEFT JOIN supplier s ON s.id = e.supplier_id
			LEFT JOIN (
				SELECT d.expenditure_id, SUM(d.amount) AS total_paid
				FROM treasury_disbursement d
				WHERE d.active = true AND d.status IN ('paid', 'completed')
				GROUP BY d.expenditure_id
			) paid ON paid.expenditure_id = e.id
			WHERE e.active = true
			  AND e.expenditure_type IN ('purchase', 'expense')
			  AND e.status NOT IN ('paid', 'cancelled')
			  AND e.total_amount - COALESCE(paid.total_paid, 0) > 0
		)
		SELECT
			supplier_name,
			COALESCE(SUM(CASE WHEN days_overdue <= 0 THEN outstanding_amount ELSE 0 END), 0) AS current_amt,
			COALESCE(SUM(CASE WHEN days_overdue BETWEEN 1 AND 30 THEN outstanding_amount ELSE 0 END), 0) AS days_30,
			COALESCE(SUM(CASE WHEN days_overdue BETWEEN 31 AND 60 THEN outstanding_amount ELSE 0 END), 0) AS days_60,
			COALESCE(SUM(CASE WHEN days_overdue BETWEEN 61 AND 90 THEN outstanding_amount ELSE 0 END), 0) AS days_90,
			COALESCE(SUM(CASE WHEN days_overdue > 90 THEN outstanding_amount ELSE 0 END), 0) AS over_90,
			COALESCE(SUM(outstanding_amount), 0) AS total
		FROM outstanding
		GROUP BY supplier_name
		ORDER BY total DESC
	`

	dbRows, err := db.QueryContext(ctx, query)
	if err != nil {
		return columns, nil, fmt.Errorf("payables aging query: %w", err)
	}
	defer dbRows.Close()

	var rows []types.TableRow
	idx := 0
	for dbRows.Next() {
		var name string
		var current, d30, d60, d90, over90, total float64
		if err := dbRows.Scan(&name, &current, &d30, &d60, &d90, &over90, &total); err != nil {
			return columns, nil, fmt.Errorf("payables aging scan: %w", err)
		}
		idx++
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("pa-%d", idx),
			Cells: []types.TableCell{
				{Value: name},
				{Value: reports.FormatCurrency(current / 100)},
				{Value: reports.FormatCurrency(d30 / 100)},
				{Value: reports.FormatCurrency(d60 / 100)},
				{Value: reports.FormatCurrency(d90 / 100)},
				{Value: reports.FormatCurrency(over90 / 100)},
				{Value: reports.FormatCurrency(total / 100)},
			},
		})
	}

	return columns, rows, nil
}
