package reports

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// MockLapsingAsset represents an asset in the lapsing schedule.
type MockLapsingAsset struct {
	Name             string
	Cost             float64
	UsefulLifeMonths int
	MonthlyDepr      float64
	Accumulated      float64
	BookValue        float64
}

func mockLapsingAssets() []MockLapsingAsset {
	return []MockLapsingAsset{
		{Name: "Salon Chair (x4)", Cost: 120000, UsefulLifeMonths: 60, MonthlyDepr: 2000, Accumulated: 24000, BookValue: 96000},
		{Name: "Hair Dryer Station", Cost: 45000, UsefulLifeMonths: 36, MonthlyDepr: 1250, Accumulated: 15000, BookValue: 30000},
		{Name: "Reception Desk", Cost: 35000, UsefulLifeMonths: 60, MonthlyDepr: 583.33, Accumulated: 7000, BookValue: 28000},
		{Name: "POS Terminal (x2)", Cost: 70000, UsefulLifeMonths: 36, MonthlyDepr: 1944.44, Accumulated: 23333, BookValue: 46667},
		{Name: "Air Conditioning Unit", Cost: 65000, UsefulLifeMonths: 120, MonthlyDepr: 541.67, Accumulated: 6500, BookValue: 58500},
		{Name: "UV Sterilizer Cabinet", Cost: 8500, UsefulLifeMonths: 36, MonthlyDepr: 236.11, Accumulated: 2833, BookValue: 5667},
		{Name: "Massage Table (x3)", Cost: 54000, UsefulLifeMonths: 60, MonthlyDepr: 900, Accumulated: 10800, BookValue: 43200},
	}
}

// NewLapsingScheduleView creates the lapsing schedule report (mock data).
func NewLapsingScheduleView(commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "assets",
		ActiveSubNav: "lapsing-schedule",
		Title:        "Lapsing Schedule",
		Subtitle:     "Asset depreciation lapsing schedule and projections",
		Icon:         "icon-calendar",
		TableID:      "lapsing-schedule-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			columns := []types.TableColumn{
				{Key: "asset", Label: "Asset", Sortable: true},
				{Key: "cost", Label: "Cost", Sortable: true, Align: "right"},
				{Key: "useful-life", Label: "Useful Life", Sortable: true, Align: "right"},
				{Key: "monthly-depr", Label: "Monthly Depr.", Sortable: true, Align: "right"},
				{Key: "accumulated", Label: "Accumulated", Sortable: true, Align: "right"},
				{Key: "book-value", Label: "Book Value", Sortable: true, Align: "right"},
			}
			assets := mockLapsingAssets()
			rows := make([]types.TableRow, len(assets))
			for i, a := range assets {
				rows[i] = types.TableRow{
					ID: fmt.Sprintf("lap-%d", i+1),
					Cells: []types.TableCell{
						{Value: a.Name},
						{Value: FormatCurrency(a.Cost)},
						{Value: fmt.Sprintf("%d months", a.UsefulLifeMonths)},
						{Value: FormatCurrency(a.MonthlyDepr)},
						{Value: FormatCurrency(a.Accumulated)},
						{Value: FormatCurrency(a.BookValue)},
					},
				}
			}
			return columns, rows, nil
		},
	})
}
