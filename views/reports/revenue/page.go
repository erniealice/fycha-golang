package revenue

import (
	"context"
	"fmt"
	"log"

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
	ContentTemplate string
	Table           *types.TableConfig
}

func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.Revenue

		records, err := deps.DB.ListRevenue(ctx)
		if err != nil {
			log.Printf("Failed to list revenue: %v", err)
			records = nil
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
				Message: "No revenue records found.",
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion: viewCtx.CacheVersion,
				Title:        l.Title,
				CurrentPath:  viewCtx.CurrentPath,
				ActiveNav:    "reports",
				ActiveSubNav: "revenue",
				HeaderTitle:  l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:   "icon-trending-up",
				CommonLabels: deps.CommonLabels,
			},
			ContentTemplate: "revenue-content",
			Table:           tableConfig,
		}

		if viewCtx.IsHTMX {
			return view.OK("revenue-content", pageData)
		}
		return view.OK("revenue", pageData)
	})
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
