package reports

import (
	"context"
	"database/sql"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NewCashBookView creates the cash book report with DB data.
func NewCashBookView(db *sql.DB, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "cash",
		ActiveSubNav: "cash-book",
		Title:        "Cash Book",
		Subtitle:     "Record of all cash receipts and disbursements",
		Icon:         "icon-book",
		TableID:      "cash-book-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			return fetchCashBook(ctx, db)
		},
	})
}

func fetchCashBook(ctx context.Context, db *sql.DB) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "date", Label: "Date", Sortable: true},
		{Key: "description", Label: "Description", Sortable: true},
		{Key: "reference", Label: "Reference", Sortable: true},
		{Key: "type", Label: "Type", Sortable: true, Width: "120px"},
		{Key: "amount", Label: "Amount", Sortable: true, Align: "right"},
	}

	// Combine revenue (receipts) and expenditure (payments) into a single ledger
	query := `
		SELECT tx_date, description, reference, tx_type, amount
		FROM (
			SELECT
				TO_CHAR(date_created, 'YYYY-MM-DD') AS tx_date,
				COALESCE(NULLIF(TRIM(customer_first_name || ' ' || customer_last_name), ''), 'Collection') AS description,
				COALESCE(NULLIF(reference_number, ''), '-') AS reference,
				'Receipt' AS tx_type,
				total_amount AS amount
			FROM revenue
			WHERE status NOT IN ('cancelled', 'draft')

			UNION ALL

			SELECT
				TO_CHAR(expenditure_date, 'YYYY-MM-DD') AS tx_date,
				COALESCE(NULLIF(name, ''), 'Payment') AS description,
				COALESCE(NULLIF(reference_number, ''), '-') AS reference,
				CASE WHEN expenditure_type = 'purchase' THEN 'Purchase' ELSE 'Expense' END AS tx_type,
				total_amount AS amount
			FROM expenditure
			WHERE status NOT IN ('cancelled', 'draft')
		) combined
		ORDER BY tx_date DESC, reference
		LIMIT 200
	`

	dbRows, err := db.QueryContext(ctx, query)
	if err != nil {
		return columns, nil, nil
	}
	defer dbRows.Close()

	var rows []types.TableRow
	idx := 0
	for dbRows.Next() {
		var date, desc, ref, txType string
		var amount int64
		if err := dbRows.Scan(&date, &desc, &ref, &txType, &amount); err != nil {
			continue
		}
		idx++

		variant := "info"
		if txType == "Receipt" {
			variant = "success"
		} else if txType == "Expense" {
			variant = "warning"
		}

		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("cb-%d", idx),
			Cells: []types.TableCell{
				{Value: date},
				{Value: desc},
				{Value: ref},
				{Type: "badge", Value: txType, Variant: variant},
				{Value: FormatCurrency(float64(amount) / 100)},
			},
		})
	}

	return columns, rows, nil
}
