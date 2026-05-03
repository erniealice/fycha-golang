package payables_aging

import (
	"context"
	"fmt"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
	reports "github.com/erniealice/fycha-golang/views/reports"
)

// NewPayablesAgingView creates the payables aging report with typed service data.
func NewPayablesAgingView(db fycha.DataSource, commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
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
		types.MoneyCell(current, "PHP", false),
		types.MoneyCell(d30, "PHP", false),
		types.MoneyCell(d60, "PHP", false),
		types.MoneyCell(d90, "PHP", false),
		types.MoneyCell(over90, "PHP", false),
		types.MoneyCell(total, "PHP", false),
	}
}

func fetchPayablesAging(ctx context.Context, db fycha.DataSource) ([]types.TableColumn, []types.TableRow, error) {
	columns := []types.TableColumn{
		{Key: "supplier", Label: "Supplier"},
		{Key: "current", Label: "Current", Align: "right"},
		{Key: "days-30", Label: "1-30 Days", Align: "right"},
		{Key: "days-60", Label: "31-60 Days", Align: "right"},
		{Key: "days-90", Label: "61-90 Days", Align: "right"},
		{Key: "over-90", Label: "Over 90 Days", Align: "right"},
		{Key: "total", Label: "Total", Align: "right"},
	}

	if db == nil {
		return columns, nil, nil
	}

	resp, err := db.GetSimplePayablesAgingReport(ctx, &reportpb.PayablesAgingReportRequest{})
	if err != nil {
		return columns, nil, fmt.Errorf("payables aging query: %w", err)
	}
	if resp == nil {
		return columns, nil, nil
	}

	var rows []types.TableRow
	for idx, row := range resp.Data {
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("pa-%d", idx+1),
			Cells: []types.TableCell{
				{Value: row.SupplierName},
				types.MoneyCell(float64(row.Current)/100, "PHP", true),
				types.MoneyCell(float64(row.Days_30)/100, "PHP", true),
				types.MoneyCell(float64(row.Days_60)/100, "PHP", true),
				types.MoneyCell(float64(row.Days_90)/100, "PHP", true),
				types.MoneyCell(float64(row.Over_90)/100, "PHP", true),
				types.MoneyCell(float64(row.Total)/100, "PHP", true),
			},
		})
	}

	return columns, rows, nil
}
