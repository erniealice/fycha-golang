package cash_book

import (
	"context"
	"fmt"

	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	fycha "github.com/erniealice/fycha-golang"
	reports "github.com/erniealice/fycha-golang/views/reports"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// NewCashBookView creates the cash book report with typed service data.
func NewCashBookView(db fycha.DataSource, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return reports.NewReportView(reports.ReportConfig{
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

func fetchCashBook(ctx context.Context, db fycha.DataSource) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "date", Label: "Date"},
		{Key: "description", Label: "Description"},
		{Key: "reference", Label: "Reference"},
		{Key: "type", Label: "Type", WidthClass: "col-2xl"},
		{Key: "amount", Label: "Amount", Align: "right"},
	}

	if db == nil {
		return columns, nil, nil
	}

	resp, err := db.GetCashBookReport(ctx, &reportpb.CashBookReportRequest{})
	if err != nil || resp == nil {
		return columns, nil, nil
	}

	var rows []types.TableRow
	for idx, row := range resp.Data {
		variant := "info"
		if row.TxType == "Receipt" {
			variant = "success"
		} else if row.TxType == "Expense" {
			variant = "warning"
		}

		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("cb-%d", idx+1),
			Cells: []types.TableCell{
				{Value: row.TxDate},
				{Value: row.Description},
				{Value: row.Reference},
				{Type: "badge", Value: row.TxType, Variant: variant},
				types.MoneyCell(float64(row.Amount)/100, "PHP", true),
			},
		})
	}

	return columns, rows, nil
}
