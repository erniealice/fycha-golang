package revenue

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

// 20260521 Wave B P1.E.5 — DB replaced with the typed ListRevenue closure
// into the service-driven domain-specific use case.
type Deps struct {
	ListRevenue  func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
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
		l := deps.Labels.Revenue
		pl := deps.Labels.Period

		// Parse filter
		filter := parseFilter(viewCtx.QueryParams, pl)

		reportURL := viewCtx.CurrentPath
		if reportURL == "" {
			reportURL = report.ReportsRevenueURL
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
		if deps.ListRevenue != nil {
			var err error
			records, err = deps.ListRevenue(ctx, &start, &end)
			if err != nil {
				log.Printf("Failed to list revenue: %v", err)
				records = nil
			}
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
		rs0 := types.MoneyCell(totalAmount, "PHP", true)
		rs1 := types.MoneyCell(avgAmount, "PHP", true)
		summary := []report.SummaryMetric{
			{Label: l.SummaryTotal, Value: rs0.Currency + " " + rs0.Value, Highlight: true},
			{Label: l.SummaryTransactions, Value: fmt.Sprintf("%d", len(records))},
			{Label: l.SummaryAverage, Value: rs1.Currency + " " + rs1.Value},
		}

		columns := []types.TableColumn{
			{Key: "reference", Label: l.Reference},
			{Key: "customer", Label: l.Customer},
			{Key: "amount", Label: l.Amount, WidthClass: "col-3xl", Align: "right"},
			{Key: "status", Label: l.Status, WidthClass: "col-2xl"},
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
				ActiveSubNav:   "revenue",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-trending-up",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:   "revenue-content",
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
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "report-revenue"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		if viewCtx.IsHTMX {
			return view.OK("revenue-content", pageData)
		}
		return view.OK("revenue", pageData)
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
		customer := toString(r["customer_name"])
		currency := toString(r["currency"])
		status := toString(r["status"])
		amountCell := types.MoneyCell(toFloat64(r["total_amount"]), currency, true)
		amountStr := amountCell.Currency + " " + amountCell.Value

		rows = append(rows, types.TableRow{
			ID: id,
			Cells: []types.TableCell{
				{Type: "text", Value: ref},
				{Type: "text", Value: customer},
				amountCell,
				{Type: "badge", Value: status, Variant: statusVariant(status)},
			},
			DataAttrs: map[string]string{
				"reference": ref,
				"customer":  customer,
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
