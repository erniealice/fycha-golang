package expenses

import (
	"context"
	"fmt"
	"log"
	"time"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

type Deps struct {
	DB           fycha.DataSource
	Labels       fycha.ReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

type PageData struct {
	types.PageData
	ContentTemplate   string
	Summary           []fycha.SummaryMetric
	Table             *types.TableConfig
	Filter            fycha.FilterState
	PeriodLabels      fycha.PeriodLabels
	ReportURL         string
	ActiveFilterCount int
}

func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Expenses
		pl := deps.Labels.Period

		// Parse filter
		filter := parseFilter(viewCtx.QueryParams, pl)

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = fycha.ReportsExpensesURL
		}

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			return view.OK("report-filter-sheet", &fycha.FilterSheetData{
				Filter:       filter,
				PeriodLabels: pl,
				ReportURL:    reportURL,
			})
		}

		// Resolve dates from period preset
		start, end := fycha.ParsePeriodPreset(filter.ActivePreset)
		if filter.ActivePreset == "custom" {
			if t, err := time.Parse("2006-01-02", filter.StartDate); err == nil {
				start = t
			}
			if t, err := time.Parse("2006-01-02", filter.EndDate); err == nil {
				end = t
			}
		}

		records, err := deps.DB.ListExpenses(ctx, &start, &end)
		if err != nil {
			log.Printf("Failed to list expenses: %v", err)
			records = nil
		}

		// Build summary
		var totalAmount float64
		var approvedCount, pendingCount int
		for _, r := range records {
			totalAmount += toFloat64(r["total_amount"])
			switch toString(r["status"]) {
			case "approved", "paid":
				approvedCount++
			case "pending":
				pendingCount++
			}
		}
		summary := []fycha.SummaryMetric{
			{Label: l.SummaryTotal, Value: formatCurrency(totalAmount), Highlight: true},
			{Label: l.SummaryCount, Value: fmt.Sprintf("%d", len(records))},
			{Label: l.SummaryApproved, Value: fmt.Sprintf("%d", approvedCount), Variant: "success"},
			{Label: l.SummaryPending, Value: fmt.Sprintf("%d", pendingCount), Variant: "warning"},
		}

		columns := []types.TableColumn{
			{Key: "reference", Label: l.Reference, Sortable: true},
			{Key: "vendor", Label: l.Vendor, Sortable: true},
			{Key: "category", Label: l.Category, Sortable: true},
			{Key: "date", Label: l.Date, Sortable: true, Width: "140px"},
			{Key: "amount", Label: l.Amount, Sortable: true, Width: "140px", Align: "right"},
			{Key: "status", Label: l.Status, Sortable: true, Width: "120px"},
		}

		rows := buildRows(records)
		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "expenses-report-table",
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           true,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "date",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   "No expenses",
				Message: "No expense records found for the selected period.",
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "reports",
				ActiveSubNav:   "expenses",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-file-minus",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "expenses-report-content",
			Summary:         summary,
			Table:           tableConfig,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: fycha.ActiveFilterCount(filter),
		}

		if viewCtx.IsHTMX {
			return view.OK("expenses-report-content", pageData)
		}
		return view.OK("expenses-report", pageData)
	})
}

func parseFilter(params map[string]string, pl fycha.PeriodLabels) fycha.FilterState {
	preset := params["period"]
	if preset == "" {
		preset = "thisMonth"
	}
	return fycha.FilterState{
		ActivePreset:  preset,
		StartDate:     params["start"],
		EndDate:       params["end"],
		PeriodPresets: fycha.DefaultPeriodPresets(pl, preset),
	}
}

func buildRows(records []map[string]any) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range records {
		id := toString(r["id"])
		ref := toString(r["reference_number"])
		vendor := toString(r["vendor_name"])
		category := toString(r["category"])
		date := toString(r["expenditure_date"])
		currency := toString(r["currency"])
		status := toString(r["status"])
		amount := currency + " " + formatAmount(r["total_amount"])

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: ref},
				{Type: "text", Value: vendor},
				{Type: "text", Value: category},
				{Type: "text", Value: date},
				{Type: "text", Value: amount},
				{Type: "badge", Value: status, Variant: statusVariant(status)},
			},
			DataAttrs: map[string]string{
				"reference": ref,
				"vendor":    vendor,
				"category":  category,
				"date":      date,
				"amount":    amount,
				"status":    status,
			},
		})
	}
	return rows
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	switch s := v.(type) {
	case string:
		return s
	case []byte:
		return string(s)
	default:
		return fmt.Sprintf("%v", v)
	}
}

func toFloat64(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int64:
		return float64(n)
	case int:
		return float64(n)
	default:
		return 0
	}
}

func formatAmount(v any) string {
	switch n := v.(type) {
	case float64:
		return fmt.Sprintf("%.2f", n)
	case float32:
		return fmt.Sprintf("%.2f", n)
	case int64:
		return fmt.Sprintf("%d.00", n)
	case int:
		return fmt.Sprintf("%d.00", n)
	case string:
		return n
	default:
		return fmt.Sprintf("%v", v)
	}
}

func formatCurrency(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
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
	formatted := fmt.Sprintf("\u20b1%s.%02d", wholeStr, frac)
	if negative {
		formatted = "-" + formatted
	}
	return formatted
}

func statusVariant(status string) string {
	switch status {
	case "paid":
		return "success"
	case "approved":
		return "info"
	case "pending":
		return "warning"
	case "cancelled":
		return "danger"
	case "draft":
		return "default"
	default:
		return "default"
	}
}
