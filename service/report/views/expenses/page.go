package expenses

import (
	"context"
	"fmt"
	"log"
	"time"

	report "github.com/erniealice/fycha-golang/service/report"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// 20260521 Wave B P1.E.5 — DB replaced with the typed ListExpenses closure
// into the service-driven domain-specific use case.
type Deps struct {
	ListExpenses func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
	Labels       report.ReportsLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
}

type PageData struct {
	types.PageData
	ContentTemplate   string
	Summary           []report.SummaryMetric
	Table             *types.TableConfig
	Filter            report.FilterState
	PeriodLabels      report.PeriodLabels
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
			reportURL = report.ReportsExpensesURL
		}

		// Handle filter sheet request
		if viewCtx.QueryParams["sheet"] == "filters" {
			return view.OK("report-filter-sheet", &report.FilterSheetData{
				Filter:       filter,
				PeriodLabels: pl,
				ReportURL:    reportURL,
			})
		}

		// Resolve dates from period preset
		start, end := report.ParsePeriodPreset(filter.ActivePreset)
		if filter.ActivePreset == "custom" {
			if t, err := time.Parse("2006-01-02", filter.StartDate); err == nil {
				start = t
			}
			if t, err := time.Parse("2006-01-02", filter.EndDate); err == nil {
				end = t
			}
		}

		var records []map[string]any
		if deps.ListExpenses != nil {
			var err error
			records, err = deps.ListExpenses(ctx, &start, &end)
			if err != nil {
				log.Printf("Failed to list expenses: %v", err)
				records = nil
			}
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
		es0 := types.MoneyCell(totalAmount, "PHP", true)
		summary := []report.SummaryMetric{
			{Label: l.SummaryTotal, Value: es0.Currency + " " + es0.Value, Highlight: true},
			{Label: l.SummaryCount, Value: fmt.Sprintf("%d", len(records))},
			{Label: l.SummaryApproved, Value: fmt.Sprintf("%d", approvedCount), Variant: "success"},
			{Label: l.SummaryPending, Value: fmt.Sprintf("%d", pendingCount), Variant: "warning"},
		}

		columns := []types.TableColumn{
			{Key: "reference", Label: l.Reference},
			{Key: "vendor", Label: l.Vendor},
			{Key: "category", Label: l.Category},
			{Key: "date", Label: l.Date, WidthClass: "col-3xl"},
			{Key: "amount", Label: l.Amount, WidthClass: "col-3xl", Align: "right"},
			{Key: "status", Label: l.Status, WidthClass: "col-2xl"},
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
				Title:   l.EmptyTitle,
				Message: l.EmptyMessage,
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "report",
				ActiveSubNav:   "expenses",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-file-minus",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:   "expenses-report-content",
			Summary:           summary,
			Table:             tableConfig,
			Filter:            filter,
			PeriodLabels:      pl,
			ReportURL:         reportURL,
			ActiveFilterCount: report.ActiveFilterCount(filter),
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-expense"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("expenses-report-content", pageData)
		}
		return view.OK("expenses-report", pageData)
	})
}

func parseFilter(params map[string]string, pl report.PeriodLabels) report.FilterState {
	preset := params["period"]
	if preset == "" {
		preset = "thisMonth"
	}
	return report.FilterState{
		ActivePreset:  preset,
		StartDate:     params["start"],
		EndDate:       params["end"],
		PeriodPresets: report.DefaultPeriodPresets(pl, preset),
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
		amountCell := types.MoneyCell(toFloat64(r["total_amount"]), currency, true)
		amountStr := amountCell.Currency + " " + amountCell.Value

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: ref},
				{Type: "text", Value: vendor},
				{Type: "text", Value: category},
				{Type: "text", Value: date},
				amountCell,
				{Type: "badge", Value: status, Variant: statusVariant(status)},
			},
			DataAttrs: map[string]string{
				"reference": ref,
				"vendor":    vendor,
				"category":  category,
				"date":      date,
				"amount":    amountStr,
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
