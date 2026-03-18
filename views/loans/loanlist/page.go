// Package loanlist provides the view for the Funding > Loans list page.
package loanlist

import (
	"context"
	"fmt"
	"log"
	"time"

	loanpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the loan list page.
type Deps struct {
	Routes       fycha.LoanRoutes
	Labels       fycha.LoanLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Loan use cases
	GetLoanListPageData func(ctx context.Context, req *loanpb.GetLoanListPageDataRequest) (*loanpb.GetLoanListPageDataResponse, error)
	ListLoans           func(ctx context.Context, req *loanpb.ListLoansRequest) (*loanpb.ListLoansResponse, error)
}

// PageData holds the data for the loan list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// LoanRow is the view-model for a single loan row.
type LoanRow struct {
	ID               string
	LoanNumber       string
	LenderName       string
	LoanType         string
	PrincipalAmount  string
	RemainingBalance string
	InterestRate     string
	Status           string
	StartDate        string
	MaturityDate     string
	Active           bool
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the loan list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		loans := fetchLoans(ctx, deps, status)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, loans, status, perms)

		heading, caption := headingForStatus(deps.Labels, status)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-banknotes",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "loan-list-content",
			Table:           tableConfig,
		}

		return view.OK("loan-list", pageData)
	})
}

// NewContentView creates the loan list HTMX partial view.
func NewContentView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		loans := fetchLoans(ctx, deps, status)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, loans, status, perms)

		heading, caption := headingForStatus(deps.Labels, status)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-banknotes",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "loan-list-content",
			Table:           tableConfig,
		}

		return view.OK("loan-list-content", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchLoans(ctx context.Context, deps *Deps, status string) []LoanRow {
	if deps.ListLoans == nil {
		return []LoanRow{}
	}

	resp, err := deps.ListLoans(ctx, &loanpb.ListLoansRequest{})
	if err != nil {
		log.Printf("ListLoans error: %v", err)
		return []LoanRow{}
	}
	if resp == nil {
		return []LoanRow{}
	}

	rows := make([]LoanRow, 0)
	for _, l := range resp.GetData() {
		row := protoToRow(l)
		// Filter by status tab
		if status == "active" && row.Status != "Active" {
			continue
		}
		if status == "completed" && row.Status != "Completed" {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func protoToRow(l *loanpb.Loan) LoanRow {
	startDate := ""
	if l.GetStartDate() > 0 {
		startDate = time.UnixMilli(l.GetStartDate()).Format("2006-01-02")
	}
	maturityDate := ""
	if l.GetMaturityDate() > 0 {
		maturityDate = time.UnixMilli(l.GetMaturityDate()).Format("2006-01-02")
	}

	return LoanRow{
		ID:               l.GetId(),
		LoanNumber:       l.GetLoanNumber(),
		LenderName:       l.GetLenderName(),
		LoanType:         loanTypeLabel(l.GetLoanType()),
		PrincipalAmount:  fmt.Sprintf("%.2f", l.GetPrincipalAmount()),
		RemainingBalance: fmt.Sprintf("%.2f", l.GetRemainingBalance()),
		InterestRate:     fmt.Sprintf("%.4f%%", l.GetInterestRate()),
		Status:           loanStatusLabel(l.GetStatus()),
		StartDate:        startDate,
		MaturityDate:     maturityDate,
		Active:           l.GetActive(),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, loans []LoanRow, status string, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "loan_number", Label: l.Columns.LoanNumber, Sortable: false, Width: "120px"},
		{Key: "lender", Label: l.Columns.LenderName, Sortable: false},
		{Key: "type", Label: l.Columns.LoanType, Sortable: false, Width: "120px"},
		{Key: "principal", Label: l.Columns.PrincipalAmount, Sortable: false, Width: "140px", Align: "right"},
		{Key: "balance", Label: l.Columns.RemainingBalance, Sortable: false, Width: "140px", Align: "right"},
		{Key: "rate", Label: l.Columns.InterestRate, Sortable: false, Width: "80px", Align: "right"},
		{Key: "maturity", Label: l.Columns.MaturityDate, Sortable: false, Width: "120px"},
	}

	rows := []types.TableRow{}
	for _, loan := range loans {
		actions := []types.TableAction{
			{Type: "view", Label: l.Actions.View, Action: "view", Href: route.ResolveURL(deps.Routes.DetailURL, "id", loan.ID)},
		}

		typeVariant := loanTypeBadgeVariant(loan.LoanType)

		row := types.TableRow{
			ID:   loan.ID,
			Href: route.ResolveURL(deps.Routes.DetailURL, "id", loan.ID),
			Cells: []types.TableCell{
				{Type: "text", Value: loan.LoanNumber},
				{Type: "text", Value: loan.LenderName},
				{Type: "badge", Value: loan.LoanType, Variant: typeVariant},
				{Type: "money", Value: loan.PrincipalAmount},
				{Type: "money", Value: loan.RemainingBalance},
				{Type: "text", Value: loan.InterestRate},
				{Type: "text", Value: loan.MaturityDate},
			},
			Actions: actions,
		}
		rows = append(rows, row)
	}

	types.ApplyColumnStyles(columns, rows)

	canCreate := perms.Can("loan", "create")
	tableConfig := &types.TableConfig{
		ID:                "loans-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowEntries:       true,
		DefaultSortColumn: "loan_number",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddLoan,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !canCreate,
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func headingForStatus(l fycha.LoanLabels, status string) (heading, caption string) {
	switch status {
	case "completed":
		return l.Page.HeadingCompleted, l.Page.CaptionCompleted
	default:
		return l.Page.HeadingActive, l.Page.CaptionActive
	}
}

func loanTypeLabel(t loanpb.LoanType) string {
	switch t {
	case loanpb.LoanType_LOAN_TYPE_PAYABLE:
		return "Payable"
	case loanpb.LoanType_LOAN_TYPE_RECEIVABLE:
		return "Receivable"
	default:
		return "Unspecified"
	}
}

func loanStatusLabel(s loanpb.LoanStatus) string {
	switch s {
	case loanpb.LoanStatus_LOAN_STATUS_DRAFT:
		return "Draft"
	case loanpb.LoanStatus_LOAN_STATUS_ACTIVE:
		return "Active"
	case loanpb.LoanStatus_LOAN_STATUS_COMPLETED:
		return "Completed"
	case loanpb.LoanStatus_LOAN_STATUS_DEFAULTED:
		return "Defaulted"
	default:
		return "Unknown"
	}
}

func loanTypeBadgeVariant(label string) string {
	switch label {
	case "Payable":
		return "amber"
	case "Receivable":
		return "sage"
	default:
		return "default"
	}
}
