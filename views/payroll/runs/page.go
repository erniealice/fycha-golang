package runs

import (
	"context"
	"fmt"
	"log"

	payrollrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_run"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
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
	Routes       fycha.PayrollRunRoutes
	Labels       fycha.PayrollLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// PayrollRun use case
	GetPayrollRunListPageData func(ctx context.Context, req *payrollrunpb.GetPayrollRunListPageDataRequest) (*payrollrunpb.GetPayrollRunListPageDataResponse, error)
}

// PageData holds the data for the payroll runs list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveStatus    string
	StatusTabs      []pyeza.TabItem
	Table           *types.TableConfig
}

// PayrollRunRow is the view-model for a single payroll run row (mapped from proto).
type PayrollRunRow struct {
	ID              string
	RunNumber       string
	PayPeriodStart  string
	PayPeriodEnd    string
	EmployeeCount   int32
	TotalGross      float64
	TotalDeductions float64
	TotalNet        float64
	Status          string
	ApprovedBy      string
	PostedAt        string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the payroll runs list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "draft"
		}

		runs := fetchPayrollRuns(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		statusTabs := buildStatusTabs(deps)
		tableConfig := buildTableConfig(deps, status, runs, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          statusTitle(deps.Labels, status),
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "payroll-runs",
				HeaderTitle:    statusTitle(deps.Labels, status),
				HeaderSubtitle: statusSubtitle(deps.Labels, status),
				HeaderIcon:     "icon-users",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "payroll-runs-content",
			ActiveStatus:    status,
			StatusTabs:      statusTabs,
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "payroll-run"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("payroll-runs", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchPayrollRuns(ctx context.Context, deps *Deps) []PayrollRunRow {
	if deps.GetPayrollRunListPageData == nil {
		return mockPayrollRuns()
	}

	resp, err := deps.GetPayrollRunListPageData(ctx, &payrollrunpb.GetPayrollRunListPageDataRequest{})
	if err != nil {
		log.Printf("GetPayrollRunListPageData error: %v", err)
		return mockPayrollRuns()
	}
	if resp == nil || !resp.GetSuccess() {
		return mockPayrollRuns()
	}

	rows := make([]PayrollRunRow, 0, len(resp.GetPayrollRunList()))
	for _, r := range resp.GetPayrollRunList() {
		rows = append(rows, protoToRow(r))
	}
	return rows
}

func protoToRow(r *payrollrunpb.PayrollRun) PayrollRunRow {
	periodStart := r.GetPayPeriodStart()
	periodEnd := r.GetPayPeriodEnd()
	approvedBy := r.GetApprovedBy()
	postedAt := r.GetPostedAtString()
	return PayrollRunRow{
		ID:              r.GetId(),
		RunNumber:       r.GetRunNumber(),
		PayPeriodStart:  periodStart,
		PayPeriodEnd:    periodEnd,
		EmployeeCount:   r.GetEmployeeCount(),
		TotalGross:      float64(r.GetTotalGross()) / 100.0,
		TotalDeductions: float64(r.GetTotalDeductions()) / 100.0,
		TotalNet:        float64(r.GetTotalNet()) / 100.0,
		Status:          statusString(r.GetStatus()),
		ApprovedBy:      approvedBy,
		PostedAt:        postedAt,
	}
}

// ---------------------------------------------------------------------------
// Mock data (UI development)
// ---------------------------------------------------------------------------

func mockPayrollRuns() []PayrollRunRow {
	return []PayrollRunRow{
		{ID: "pr-001", RunNumber: "PR-2026-001", PayPeriodStart: "Mar 1, 2026", PayPeriodEnd: "Mar 15, 2026", EmployeeCount: 12, TotalGross: 312000, TotalDeductions: 48600, TotalNet: 263400, Status: "draft"},
		{ID: "pr-002", RunNumber: "PR-2026-002", PayPeriodStart: "Feb 16, 2026", PayPeriodEnd: "Feb 28, 2026", EmployeeCount: 12, TotalGross: 308500, TotalDeductions: 47200, TotalNet: 261300, Status: "calculated"},
		{ID: "pr-003", RunNumber: "PR-2026-003", PayPeriodStart: "Feb 1, 2026", PayPeriodEnd: "Feb 15, 2026", EmployeeCount: 11, TotalGross: 290000, TotalDeductions: 44800, TotalNet: 245200, Status: "approved", ApprovedBy: "Maria Santos"},
		{ID: "pr-004", RunNumber: "PR-2026-004", PayPeriodStart: "Jan 16, 2026", PayPeriodEnd: "Jan 31, 2026", EmployeeCount: 11, TotalGross: 285000, TotalDeductions: 43900, TotalNet: 241100, Status: "posted", PostedAt: "Feb 1, 2026"},
		{ID: "pr-005", RunNumber: "PR-2026-005", PayPeriodStart: "Jan 1, 2026", PayPeriodEnd: "Jan 15, 2026", EmployeeCount: 10, TotalGross: 260000, TotalDeductions: 40100, TotalNet: 219900, Status: "posted", PostedAt: "Jan 16, 2026"},
	}
}

// ---------------------------------------------------------------------------
// Tab builder
// ---------------------------------------------------------------------------

func buildStatusTabs(deps *Deps) []pyeza.TabItem {
	l := deps.Labels.Run.Tabs
	base := deps.Routes.ListURL
	return []pyeza.TabItem{
		{Key: "draft", Label: l.Draft, Href: route.ResolveURL(base, "status", "draft"), Icon: "", Count: 0, Disabled: false},
		{Key: "calculated", Label: l.Calculated, Href: route.ResolveURL(base, "status", "calculated"), Icon: "", Count: 0, Disabled: false},
		{Key: "approved", Label: l.Approved, Href: route.ResolveURL(base, "status", "approved"), Icon: "", Count: 0, Disabled: false},
		{Key: "posted", Label: l.Posted, Href: route.ResolveURL(base, "status", "posted"), Icon: "", Count: 0, Disabled: false},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, status string, runs []PayrollRunRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := payrollRunColumns(l)
	rows := buildTableRows(runs, status, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                   "payroll-runs-table",
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowExport:           true,
		ShowEntries:          true,
		DefaultSortColumn:    "run_number",
		DefaultSortDirection: "desc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Run.Empty.Title,
			Message: l.Run.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Run.Buttons.NewRun,
			ActionURL:       "#",
			Icon:            "icon-plus",
			Disabled:        !perms.Can("payroll_run", "create"),
			DisabledTooltip: l.Run.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func payrollRunColumns(l fycha.PayrollLabels) []types.TableColumn {
	c := l.Run.Columns
	return []types.TableColumn{
		{Key: "run_number", Label: c.RunNumber, Sortable: true, Width: "130px"},
		{Key: "pay_period", Label: c.PayPeriod, Sortable: false},
		{Key: "employees", Label: c.Employees, Sortable: true, Width: "110px", Align: "right"},
		{Key: "total_gross", Label: c.TotalGross, Sortable: true, Width: "150px", Align: "right"},
		{Key: "total_deductions", Label: c.TotalDeductions, Sortable: true, Width: "150px", Align: "right"},
		{Key: "total_net", Label: c.TotalNet, Sortable: true, Width: "150px", Align: "right"},
		{Key: "status", Label: c.Status, Sortable: true, Width: "120px"},
	}
}

func buildTableRows(runs []PayrollRunRow, status string, l fycha.PayrollLabels, routes fycha.PayrollRunRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, r := range runs {
		if status != "all" && r.Status != status {
			continue
		}

		canView := perms.Can("payroll_run", "read")

		actions := []types.TableAction{
			{
				Type:     "view",
				Label:    l.Run.Actions.View,
				Action:   "view",
				Href:     route.ResolveURL(routes.DetailURL, "id", r.ID),
				Disabled: !canView, DisabledTooltip: l.Run.Actions.NoPermission,
			},
		}

		payPeriod := r.PayPeriodStart
		if r.PayPeriodEnd != "" {
			payPeriod += " – " + r.PayPeriodEnd
		}

		rows = append(rows, types.TableRow{
			ID:   r.ID,
			Href: route.ResolveURL(routes.DetailURL, "id", r.ID),
			Cells: []types.TableCell{
				{Type: "text", Value: r.RunNumber},
				{Type: "text", Value: payPeriod},
				{Type: "text", Value: fmt.Sprintf("%d", r.EmployeeCount)},
				{Type: "text", Value: formatCurrency(r.TotalGross)},
				{Type: "text", Value: formatCurrency(r.TotalDeductions)},
				{Type: "text", Value: formatCurrency(r.TotalNet)},
				{Type: "badge", Value: statusLabel(l, r.Status), Variant: runStatusVariant(r.Status)},
			},
			DataAttrs: map[string]string{
				"run_number": r.RunNumber,
				"status":     r.Status,
			},
			Actions: actions,
		})
	}
	return rows
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func statusString(s payrollrunpb.PayrollRunStatus) string {
	switch s {
	case payrollrunpb.PayrollRunStatus_PAYROLL_RUN_STATUS_DRAFT:
		return "draft"
	case payrollrunpb.PayrollRunStatus_PAYROLL_RUN_STATUS_CALCULATED:
		return "calculated"
	case payrollrunpb.PayrollRunStatus_PAYROLL_RUN_STATUS_APPROVED:
		return "approved"
	case payrollrunpb.PayrollRunStatus_PAYROLL_RUN_STATUS_POSTED:
		return "posted"
	default:
		return "draft"
	}
}

func statusTitle(l fycha.PayrollLabels, status string) string {
	switch status {
	case "draft":
		return l.Run.Page.HeadingDraft
	case "calculated":
		return l.Run.Page.HeadingCalculated
	case "approved":
		return l.Run.Page.HeadingApproved
	case "posted":
		return l.Run.Page.HeadingPosted
	default:
		return l.Run.Page.HeadingDraft
	}
}

func statusSubtitle(l fycha.PayrollLabels, status string) string {
	switch status {
	case "draft":
		return l.Run.Page.SubtitleDraft
	case "calculated":
		return l.Run.Page.SubtitleCalculated
	case "approved":
		return l.Run.Page.SubtitleApproved
	case "posted":
		return l.Run.Page.SubtitlePosted
	default:
		return l.Run.Page.SubtitleDraft
	}
}

func statusLabel(l fycha.PayrollLabels, status string) string {
	switch status {
	case "draft":
		return l.Run.Tabs.Draft
	case "calculated":
		return l.Run.Tabs.Calculated
	case "approved":
		return l.Run.Tabs.Approved
	case "posted":
		return l.Run.Tabs.Posted
	default:
		return status
	}
}

func runStatusVariant(status string) string {
	switch status {
	case "draft":
		return "default"
	case "calculated":
		return "amber"
	case "approved":
		return "navy"
	case "posted":
		return "success"
	default:
		return "default"
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
