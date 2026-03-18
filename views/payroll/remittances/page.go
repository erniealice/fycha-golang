package remittances

import (
	"context"
	"fmt"
	"log"

	payrollremittancepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_remittance"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies.
type Deps struct {
	Routes       fycha.PayrollRemittanceRoutes
	Labels       fycha.PayrollLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// PayrollRemittance use case
	GetPayrollRemittanceListPageData func(ctx context.Context, req *payrollremittancepb.GetPayrollRemittanceListPageDataRequest) (*payrollremittancepb.GetPayrollRemittanceListPageDataResponse, error)
}

// PageData holds the data for the payroll remittances list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveStatus    string
	StatusTabs      []pyeza.TabItem
	Table           *types.TableConfig
}

// RemittanceRow is the view-model for a single remittance row.
type RemittanceRow struct {
	ID              string
	RemittanceType  string
	Amount          float64
	DueDate         string
	Status          string
	FiledAt         string
	PaidAt          string
	ReferenceNumber string
	PayrollRunID    string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the payroll remittances list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "pending"
		}

		remittances := fetchRemittances(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		statusTabs := buildStatusTabs(deps)
		tableConfig := buildTableConfig(deps, status, remittances, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          remittanceStatusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "payroll-remittances",
				HeaderTitle:    remittanceStatusTitle(deps.Labels, status),
				HeaderSubtitle: remittanceStatusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-landmark",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "remittances-content",
			ActiveStatus:    status,
			StatusTabs:      statusTabs,
			Table:           tableConfig,
		}

		return view.OK("remittances", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchRemittances(ctx context.Context, deps *Deps) []RemittanceRow {
	if deps.GetPayrollRemittanceListPageData == nil {
		return mockRemittances()
	}

	resp, err := deps.GetPayrollRemittanceListPageData(ctx, &payrollremittancepb.GetPayrollRemittanceListPageDataRequest{})
	if err != nil {
		log.Printf("GetPayrollRemittanceListPageData error: %v", err)
		return mockRemittances()
	}
	if resp == nil || !resp.GetSuccess() {
		return mockRemittances()
	}

	rows := make([]RemittanceRow, 0, len(resp.GetPayrollRemittanceList()))
	for _, r := range resp.GetPayrollRemittanceList() {
		rows = append(rows, protoToRow(r))
	}
	return rows
}

func protoToRow(r *payrollremittancepb.PayrollRemittance) RemittanceRow {
	return RemittanceRow{
		ID:              r.GetId(),
		RemittanceType:  remittanceTypeString(r.GetRemittanceType()),
		Amount:          r.GetAmount(),
		DueDate:         r.GetDueDateString(),
		Status:          remittanceStatusString(r.GetStatus()),
		FiledAt:         r.GetFiledAtString(),
		PaidAt:          r.GetPaidAtString(),
		ReferenceNumber: r.GetReferenceNumber(),
		PayrollRunID:    r.GetPayrollRunId(),
	}
}

// ---------------------------------------------------------------------------
// Mock data (UI development)
// ---------------------------------------------------------------------------

func mockRemittances() []RemittanceRow {
	return []RemittanceRow{
		{ID: "rem-001", RemittanceType: "sss", Amount: 18240, DueDate: "Apr 15, 2026", Status: "pending", PayrollRunID: "pr-001"},
		{ID: "rem-002", RemittanceType: "philhealth", Amount: 9120, DueDate: "Apr 15, 2026", Status: "pending", PayrollRunID: "pr-001"},
		{ID: "rem-003", RemittanceType: "pagibig", Amount: 3040, DueDate: "Apr 15, 2026", Status: "pending", PayrollRunID: "pr-001"},
		{ID: "rem-004", RemittanceType: "bir_withholding", Amount: 14850, DueDate: "Apr 10, 2026", Status: "pending", PayrollRunID: "pr-001"},
		{ID: "rem-005", RemittanceType: "sss", Amount: 17900, DueDate: "Mar 15, 2026", Status: "filed", FiledAt: "Mar 14, 2026", ReferenceNumber: "SSS-2026-03-001", PayrollRunID: "pr-004"},
		{ID: "rem-006", RemittanceType: "philhealth", Amount: 8950, DueDate: "Mar 15, 2026", Status: "filed", FiledAt: "Mar 14, 2026", ReferenceNumber: "PH-2026-03-001", PayrollRunID: "pr-004"},
		{ID: "rem-007", RemittanceType: "pagibig", Amount: 2980, DueDate: "Mar 15, 2026", Status: "paid", FiledAt: "Mar 14, 2026", PaidAt: "Mar 14, 2026", ReferenceNumber: "HDMF-2026-03-001", PayrollRunID: "pr-004"},
	}
}

// ---------------------------------------------------------------------------
// Tab builder
// ---------------------------------------------------------------------------

func buildStatusTabs(deps *Deps) []pyeza.TabItem {
	l := deps.Labels.Remittance.Tabs
	base := deps.Routes.ListURL
	return []pyeza.TabItem{
		{Key: "pending", Label: l.Pending, Href: route.ResolveURL(base, "status", "pending"), Icon: "", Count: 0, Disabled: false},
		{Key: "filed", Label: l.Filed, Href: route.ResolveURL(base, "status", "filed"), Icon: "", Count: 0, Disabled: false},
		{Key: "paid", Label: l.Paid, Href: route.ResolveURL(base, "status", "paid"), Icon: "", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, status string, remittances []RemittanceRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := remittanceColumns(l)
	rows := buildTableRows(remittances, status, l, perms)
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "remittances-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowExport:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "due_date",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Remittance.Empty.Title,
			Message: l.Remittance.Empty.Message,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func remittanceColumns(l fycha.PayrollLabels) []types.TableColumn {
	c := l.Remittance.Columns
	return []types.TableColumn{
		{Key: "remittance_type", Label: c.RemittanceType, Sortable: true, Width: "150px"},
		{Key: "amount", Label: c.Amount, Sortable: true, Width: "140px", Align: "right"},
		{Key: "due_date", Label: c.DueDate, Sortable: true, Width: "130px"},
		{Key: "status", Label: c.Status, Sortable: true, Width: "110px"},
		{Key: "filed_at", Label: c.FiledAt, Sortable: true, Width: "130px"},
		{Key: "reference_number", Label: c.ReferenceNumber, Sortable: false},
	}
}

func buildTableRows(remittances []RemittanceRow, status string, l fycha.PayrollLabels, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range remittances {
		if status != "all" && r.Status != status {
			continue
		}

		rows = append(rows, types.TableRow{
			ID: r.ID,
			Cells: []types.TableCell{
				{Type: "badge", Value: remittanceTypeLabel(l, r.RemittanceType), Variant: remittanceTypeVariant(r.RemittanceType)},
				{Type: "text", Value: formatCurrency(r.Amount)},
				{Type: "text", Value: r.DueDate},
				{Type: "badge", Value: remittanceStatusLabel(l, r.Status), Variant: remittanceStatusVariant(r.Status)},
				{Type: "text", Value: r.FiledAt},
				{Type: "text", Value: r.ReferenceNumber},
			},
			DataAttrs: map[string]string{
				"type":   r.RemittanceType,
				"status": r.Status,
			},
		})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func remittanceTypeString(t payrollremittancepb.RemittanceType) string {
	switch t {
	case payrollremittancepb.RemittanceType_REMITTANCE_TYPE_SSS:
		return "sss"
	case payrollremittancepb.RemittanceType_REMITTANCE_TYPE_PHILHEALTH:
		return "philhealth"
	case payrollremittancepb.RemittanceType_REMITTANCE_TYPE_PAGIBIG:
		return "pagibig"
	case payrollremittancepb.RemittanceType_REMITTANCE_TYPE_BIR_WITHHOLDING:
		return "bir_withholding"
	default:
		return "unknown"
	}
}

func remittanceStatusString(s payrollremittancepb.RemittanceStatus) string {
	switch s {
	case payrollremittancepb.RemittanceStatus_REMITTANCE_STATUS_PENDING:
		return "pending"
	case payrollremittancepb.RemittanceStatus_REMITTANCE_STATUS_FILED:
		return "filed"
	case payrollremittancepb.RemittanceStatus_REMITTANCE_STATUS_PAID:
		return "paid"
	default:
		return "pending"
	}
}

func remittanceTypeLabel(l fycha.PayrollLabels, t string) string {
	switch t {
	case "sss":
		return l.Remittance.Types.SSS
	case "philhealth":
		return l.Remittance.Types.PhilHealth
	case "pagibig":
		return l.Remittance.Types.PagIBIG
	case "bir_withholding":
		return l.Remittance.Types.BIRWithholding
	default:
		return t
	}
}

func remittanceTypeVariant(t string) string {
	switch t {
	case "sss":
		return "sage"
	case "philhealth":
		return "terracotta"
	case "pagibig":
		return "amber"
	case "bir_withholding":
		return "navy"
	default:
		return "default"
	}
}

func remittanceStatusLabel(l fycha.PayrollLabels, status string) string {
	switch status {
	case "pending":
		return l.Remittance.Tabs.Pending
	case "filed":
		return l.Remittance.Tabs.Filed
	case "paid":
		return l.Remittance.Tabs.Paid
	default:
		return status
	}
}

func remittanceStatusVariant(status string) string {
	switch status {
	case "pending":
		return "amber"
	case "filed":
		return "navy"
	case "paid":
		return "success"
	default:
		return "default"
	}
}

func remittanceStatusTitle(l fycha.PayrollLabels, status string) string {
	switch status {
	case "pending":
		return l.Remittance.Page.HeadingPending
	case "filed":
		return l.Remittance.Page.HeadingFiled
	case "paid":
		return l.Remittance.Page.HeadingPaid
	default:
		return l.Remittance.Page.HeadingPending
	}
}

func remittanceStatusSubtitle(l fycha.PayrollLabels, status string) string {
	switch status {
	case "pending":
		return l.Remittance.Page.SubtitlePending
	case "filed":
		return l.Remittance.Page.SubtitleFiled
	case "paid":
		return l.Remittance.Page.SubtitlePaid
	default:
		return l.Remittance.Page.SubtitlePending
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
