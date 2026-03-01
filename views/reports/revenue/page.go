package revenue

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
		l := deps.Labels.Revenue
		pl := deps.Labels.Period

		// Parse filter
		filter := parseFilter(viewCtx.QueryParams, pl)

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = fycha.ReportsRevenueURL
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

		records, err := deps.DB.ListRevenue(ctx, &start, &end)
		if err != nil {
			log.Printf("Failed to list revenue: %v", err)
			records = nil
		}

		// Build summary
		var totalAmount float64
		for _, r := range records {
			totalAmount += toFloat64(r["total_amount"])
		}
		avgAmount := 0.0
		if len(records) > 0 {
			avgAmount = totalAmount / float64(len(records))
		}
		summary := []fycha.SummaryMetric{
			{Label: l.SummaryTotal, Value: formatCurrency(totalAmount), Highlight: true},
			{Label: l.SummaryTransactions, Value: fmt.Sprintf("%d", len(records))},
			{Label: l.SummaryAverage, Value: formatCurrency(avgAmount)},
		}

		columns := []types.TableColumn{
			{Key: "reference", Label: l.Reference, Sortable: true},
			{Key: "customer", Label: l.Customer, Sortable: true},
			{Key: "amount", Label: l.Amount, Sortable: true, Width: "140px", Align: "right"},
			{Key: "status", Label: l.Status, Sortable: true, Width: "120px"},
		}

		rows := buildRows(records)
		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "revenue-table",
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           true,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "reference",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   "No revenue",
				Message: "No revenue records found for the selected period.",
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "reports",
				ActiveSubNav:   "revenue",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-trending-up",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "revenue-content",
			Summary:         summary,
			Table:           tableConfig,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: fycha.ActiveFilterCount(filter),
		}

		if viewCtx.IsHTMX {
			return view.OK("revenue-content", pageData)
		}
		return view.OK("revenue", pageData)
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
		customer := toString(r["customer_name"])
		currency := toString(r["currency"])
		status := toString(r["status"])
		amount := currency + " " + formatAmount(r["total_amount"])

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: ref},
				{Type: "text", Value: customer},
				{Type: "text", Value: amount},
				{Type: "badge", Value: status, Variant: statusVariant(status)},
			},
			DataAttrs: map[string]string{
				"reference": ref,
				"customer":  customer,
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
	case "completed", "paid":
		return "success"
	case "pending":
		return "warning"
	case "cancelled", "refunded":
		return "danger"
	default:
		return "default"
	}
}
