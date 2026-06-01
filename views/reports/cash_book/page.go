package cash_book

import (
	"context"
	"fmt"

	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"
	reports "github.com/erniealice/fycha-golang/views/reports"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// CashBookFetcher is the typed closure consumed by the cash book view —
// service-driven gross/cashflow use case (Wave B P1.E.3).
type CashBookFetcher func(context.Context, *gcfpb.GetCashBookRequest) (*gcfpb.GetCashBookResponse, error)

// NewCashBookView creates the cash book report with typed service data.
//
// 20260521 Wave B P1.E.3 — `db fycha.DataSource` replaced with the typed
// `GetCashBookReport` closure into the service-driven gross/cashflow use
// case. Nil-safe: when the closure is nil the view renders an empty table.
func NewCashBookView(getCashBook CashBookFetcher, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
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
			return fetchCashBook(ctx, getCashBook, commonLabels)
		},
	})
}

func fetchCashBook(ctx context.Context, getCashBook CashBookFetcher, commonLabels pyeza.CommonLabels) ([]types.TableColumn, []types.TableRow, error) {
	cols := commonLabels.Columns
	columns := []types.TableColumn{
		{Key: "date", Label: cols.Date},
		{Key: "description", Label: cols.Description},
		{Key: "reference", Label: cols.Reference},
		{Key: "type", Label: cols.Type, WidthClass: "col-2xl"},
		{Key: "amount", Label: cols.Amount, Align: "right"},
	}

	if getCashBook == nil {
		return columns, nil, nil
	}

	resp, err := getCashBook(ctx, &gcfpb.GetCashBookRequest{})
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
				types.MoneyCell(float64(row.Amount), "PHP", true),
			},
		})
	}

	return columns, rows, nil
}
