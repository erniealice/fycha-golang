package reports

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// MockDepreciationPolicy represents a depreciation policy.
type MockDepreciationPolicy struct {
	Name        string
	Method      string
	UsefulLife  int // months
	SalvageRate float64
	Active      bool
}

func mockDepreciationPolicies() []MockDepreciationPolicy {
	return []MockDepreciationPolicy{
		{Name: "Furniture & Fixtures", Method: "Straight Line", UsefulLife: 60, SalvageRate: 5, Active: true},
		{Name: "Equipment", Method: "Straight Line", UsefulLife: 36, SalvageRate: 10, Active: true},
		{Name: "Vehicles", Method: "Declining Balance", UsefulLife: 60, SalvageRate: 15, Active: true},
		{Name: "IT Equipment", Method: "Straight Line", UsefulLife: 24, SalvageRate: 5, Active: true},
		{Name: "Building Improvements", Method: "Straight Line", UsefulLife: 120, SalvageRate: 0, Active: true},
		{Name: "Leasehold Improvements", Method: "Straight Line", UsefulLife: 60, SalvageRate: 0, Active: false},
	}
}

// NewDepreciationPoliciesView creates the depreciation policies report (mock data).
func NewDepreciationPoliciesView(commonLabels pyeza.CommonLabels, tableLabels types.TableLabels) view.View {
	return NewReportView(ReportConfig{
		ActiveNav:    "assets",
		ActiveSubNav: "depreciation-policies",
		Title:        "Depreciation Policies",
		Subtitle:     "Configure asset depreciation methods and useful life settings",
		Icon:         "icon-settings",
		TableID:      "depreciation-policies-table",
		CommonLabels: commonLabels,
		TableLabels:  tableLabels,
		BuildData: func(ctx context.Context) ([]types.TableColumn, []types.TableRow, error) {
			columns := []types.TableColumn{
				{Key: "name", Label: "Policy Name", Sortable: true},
				{Key: "method", Label: "Method", Sortable: true},
				{Key: "useful-life", Label: "Default Useful Life", Sortable: true, Align: "right"},
				{Key: "salvage", Label: "Salvage Rate", Sortable: true, Align: "right"},
				{Key: "status", Label: "Status", Sortable: true, Width: "120px"},
			}
			policies := mockDepreciationPolicies()
			rows := make([]types.TableRow, len(policies))
			for i, p := range policies {
				status := "active"
				variant := "success"
				if !p.Active {
					status = "inactive"
					variant = "warning"
				}
				rows[i] = types.TableRow{
					ID: fmt.Sprintf("pol-%d", i+1),
					Cells: []types.TableCell{
						{Value: p.Name},
						{Value: p.Method},
						{Value: fmt.Sprintf("%d months", p.UsefulLife)},
						{Value: fmt.Sprintf("%.0f%%", p.SalvageRate)},
						{Type: "badge", Value: status, Variant: variant},
					},
				}
			}
			return columns, rows, nil
		},
	})
}
