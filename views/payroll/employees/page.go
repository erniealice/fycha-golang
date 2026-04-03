package employees

import (
	"context"
	"fmt"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies.
type Deps struct {
	Routes       fycha.PayrollEmployeeRoutes
	Labels       fycha.PayrollLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

// PageData holds the data for the employees list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// EmployeeRow is the view-model for a single employee row.
type EmployeeRow struct {
	ID           string
	Name         string
	Position     string
	Department   string
	BasicSalary  float64
	PayFrequency string
	Status       string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the payroll employees list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Employee.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "payroll-employees",
				HeaderTitle:    deps.Labels.Employee.Page.Heading,
				HeaderSubtitle: deps.Labels.Employee.Page.Caption,
				HeaderIcon:     "icon-user-check",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "employees-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "payroll-employee"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("employees", pageData)
	})
}

// ---------------------------------------------------------------------------
// Mock data (UI development)
// ---------------------------------------------------------------------------

func mockEmployees() []EmployeeRow {
	return []EmployeeRow{
		{ID: "emp-001", Name: "Maria Santos", Position: "Senior Stylist", Department: "Operations", BasicSalary: 28000, PayFrequency: "semi-monthly", Status: "active"},
		{ID: "emp-002", Name: "Juan dela Cruz", Position: "Salon Manager", Department: "Management", BasicSalary: 45000, PayFrequency: "semi-monthly", Status: "active"},
		{ID: "emp-003", Name: "Ana Reyes", Position: "Junior Stylist", Department: "Operations", BasicSalary: 20000, PayFrequency: "semi-monthly", Status: "active"},
		{ID: "emp-004", Name: "Pedro Garcia", Position: "Receptionist", Department: "Front Desk", BasicSalary: 18000, PayFrequency: "semi-monthly", Status: "active"},
		{ID: "emp-005", Name: "Luisa Torres", Position: "Nail Technician", Department: "Operations", BasicSalary: 22000, PayFrequency: "semi-monthly", Status: "active"},
		{ID: "emp-006", Name: "Carlos Mendoza", Position: "Massage Therapist", Department: "Operations", BasicSalary: 24000, PayFrequency: "semi-monthly", Status: "active"},
		{ID: "emp-007", Name: "Rosa Bautista", Position: "Aesthetician", Department: "Operations", BasicSalary: 25000, PayFrequency: "semi-monthly", Status: "inactive"},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := employeeColumns(l)
	rows := buildTableRows(mockEmployees(), l, perms)
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "payroll-employees-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "name",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Employee.Empty.Title,
			Message: l.Employee.Empty.Message,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func employeeColumns(l fycha.PayrollLabels) []types.TableColumn {
	c := l.Employee.Columns
	return []types.TableColumn{
		{Key: "name", Label: c.Name, Sortable: true},
		{Key: "position", Label: c.Position, Sortable: true},
		{Key: "department", Label: c.Department, Sortable: true, Width: "150px"},
		{Key: "basic_salary", Label: c.BasicSalary, Sortable: true, Width: "150px", Align: "right"},
		{Key: "pay_frequency", Label: c.PayFrequency, Sortable: true, Width: "140px"},
		{Key: "status", Label: c.Status, Sortable: true, Width: "110px"},
	}
}

func buildTableRows(employees []EmployeeRow, l fycha.PayrollLabels, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, emp := range employees {
		statusVariant := "success"
		if emp.Status == "inactive" {
			statusVariant = "warning"
		}
		statusLabel := l.Employee.Status.Active
		if emp.Status == "inactive" {
			statusLabel = l.Employee.Status.Inactive
		}

		rows = append(rows, types.TableRow{
			ID: emp.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: emp.Name},
				{Type: "text", Value: emp.Position},
				{Type: "text", Value: emp.Department},
				{Type: "text", Value: formatCurrency(emp.BasicSalary)},
				{Type: "text", Value: payFrequencyLabel(l, emp.PayFrequency)},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			DataAttrs: map[string]string{
				"name":       emp.Name,
				"department": emp.Department,
				"status":     emp.Status,
			},
		})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func payFrequencyLabel(l fycha.PayrollLabels, freq string) string {
	switch freq {
	case "semi-monthly":
		return l.Employee.PayFrequency.SemiMonthly
	case "monthly":
		return l.Employee.PayFrequency.Monthly
	case "weekly":
		return l.Employee.PayFrequency.Weekly
	default:
		return freq
	}
}

func formatCurrency(amount float64) string {
	whole := int64(amount)
	frac := int64((amount-float64(whole))*100 + 0.5)
	if frac >= 100 {
		whole++
		frac -= 100
	}
	wholeStr := fmt.Sprintf("%d", whole)
	n := len(wholeStr)
	if n > 3 {
		var result []byte
		for i, ch := range wholeStr {
			if i > 0 && (n-i)%3 == 0 {
				result = append(result, ',')
			}
			result = append(result, byte(ch))
		}
		wholeStr = string(result)
	}
	return fmt.Sprintf("\u20b1%s.%02d", wholeStr, frac)
}
