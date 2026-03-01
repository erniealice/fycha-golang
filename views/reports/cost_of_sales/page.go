package cost_of_sales

import (
	"context"
	"fmt"
	"log"
	"strconv"

	fycha "github.com/erniealice/fycha-golang"
	reportpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/reporting/gross_profit"
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

type SummaryMetric struct {
	Label     string
	Value     string
	Highlight bool
	Variant   string
}

type PageData struct {
	types.PageData
	ContentTemplate string
	Summary         []SummaryMetric
	Table           *types.TableConfig
}

func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		l := deps.Labels.CostOfSales

		// Reuse the gross profit report — it already has COGS data per group
		req := &reportpb.GrossProfitReportRequest{}
		groupBy := "product"
		req.GroupBy = &groupBy

		resp, err := deps.DB.GetGrossProfitReport(ctx, req)
		if err != nil {
			log.Printf("Failed to get cost of sales report: %v", err)
			resp = &reportpb.GrossProfitReportResponse{
				LineItems: []*reportpb.GrossProfitLineItem{},
				Summary:   &reportpb.GrossProfitSummary{},
			}
		}

		s := resp.GetSummary()
		if s == nil {
			s = &reportpb.GrossProfitSummary{}
		}

		cogsRatio := 0.0
		if s.GetNetRevenue() > 0 {
			cogsRatio = (s.GetTotalCogs() / s.GetNetRevenue()) * 100
		}

		summary := []SummaryMetric{
			{Label: l.SummaryTotalCOGS, Value: formatCurrency(s.GetTotalCogs()), Highlight: true},
			{Label: l.SummaryRevenue, Value: formatCurrency(s.GetNetRevenue())},
			{Label: l.SummaryCOGSRatio, Value: fmt.Sprintf("%.1f%%", cogsRatio)},
			{Label: l.SummaryUnits, Value: strconv.FormatInt(s.GetTotalUnitsSold(), 10)},
		}

		// Table: group key, COGS, revenue, ratio, units
		columns := []types.TableColumn{
			{Key: "group", Label: "Item", Sortable: true},
			{Key: "cogs", Label: "COGS", Sortable: true, Align: "right", MinWidth: "120px"},
			{Key: "revenue", Label: "Net Revenue", Sortable: true, Align: "right", MinWidth: "120px"},
			{Key: "ratio", Label: "COGS %", Sortable: true, Align: "right", MinWidth: "80px"},
			{Key: "units", Label: "Units", Sortable: true, Align: "right", MinWidth: "80px"},
		}

		rows := make([]types.TableRow, 0, len(resp.GetLineItems()))
		for _, item := range resp.GetLineItems() {
			ratio := 0.0
			if item.GetNetRevenue() > 0 {
				ratio = (item.GetCostOfGoodsSold() / item.GetNetRevenue()) * 100
			}
			rows = append(rows, types.TableRow{
				ID: item.GetGroupKey(),
				Cells: []types.TableCell{
					{Type: "name", Value: item.GetGroupKey()},
					{Type: "text", Value: formatCurrency(item.GetCostOfGoodsSold())},
					{Type: "text", Value: formatCurrency(item.GetNetRevenue())},
					{Type: "text", Value: fmt.Sprintf("%.1f%%", ratio)},
					{Type: "text", Value: strconv.FormatInt(item.GetUnitsSold(), 10)},
				},
				DataAttrs: map[string]string{
					"cogs":    fmt.Sprintf("%.2f", item.GetCostOfGoodsSold()),
					"revenue": fmt.Sprintf("%.2f", item.GetNetRevenue()),
					"ratio":   fmt.Sprintf("%.1f", ratio),
					"units":   strconv.FormatInt(item.GetUnitsSold(), 10),
				},
			})
		}

		types.ApplyColumnStyles(columns, rows)

		tableConfig := &types.TableConfig{
			ID:                   "cost-of-sales-table",
			Columns:              columns,
			Rows:                 rows,
			ShowSearch:           false,
			ShowFilters:          true,
			ShowSort:             true,
			ShowColumns:          true,
			ShowExport:           true,
			ShowDensity:          true,
			ShowEntries:          true,
			DefaultSortColumn:    "cogs",
			DefaultSortDirection: "desc",
			Labels:               deps.TableLabels,
			EmptyState: types.TableEmptyState{
				Title:   "No data",
				Message: "No cost of sales data found.",
			},
		}
		types.ApplyTableSettings(tableConfig)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          l.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "reports",
				ActiveSubNav:   "cost-of-sales",
				HeaderTitle:    l.Title,
				HeaderSubtitle: l.Subtitle,
				HeaderIcon:     "icon-package",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "cost-of-sales-content",
			Summary:         summary,
			Table:           tableConfig,
		}

		if viewCtx.IsHTMX {
			return view.OK("cost-of-sales-content", pageData)
		}
		return view.OK("cost-of-sales", pageData)
	})
}

func formatCurrency(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
	whole := int64(amount)
	frac := int64((amount - float64(whole)) * 100 + 0.5)
	if frac >= 100 {
		whole++
		frac -= 100
	}
	wholeStr := strconv.FormatInt(whole, 10)
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
	formatted := fmt.Sprintf("₱%s.%02d", wholeStr, frac)
	if negative {
		formatted = "-" + formatted
	}
	return formatted
}
