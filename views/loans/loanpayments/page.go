// Package loanpayments provides the view for the Funding > Loans > Payments list page.
package loanpayments

import (
	"context"
	"fmt"
	"log"

	loanpaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan_payment"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the loan payments list.
type Deps struct {
	Routes       fycha.LoanPaymentRoutes
	Labels       fycha.LoanPaymentLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Loan payment use cases
	ListLoanPayments func(ctx context.Context, req *loanpaymentpb.ListLoanPaymentsRequest) (*loanpaymentpb.ListLoanPaymentsResponse, error)
}

// PageData holds the data for the loan payments list page.
type PageData struct {
	types.PageData
	ContentTemplate  string
	LoanID           string
	Table            *types.TableConfig
	RecordPaymentURL string
	Labels           fycha.LoanPaymentLabels
}

// LoanPaymentRow is the view-model for a single loan payment row.
type LoanPaymentRow struct {
	ID               string
	LoanID           string
	PaymentNumber    string
	PaymentDate      string
	PrincipalAmount  string
	InterestAmount   string
	FeeAmount        string
	TotalAmount      string
	RemainingBalance string
	Notes            string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the loan payments list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		loanID := viewCtx.Request.URL.Query().Get("loan_id")
		payments := fetchPayments(ctx, deps, loanID)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, payments, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:  "loan-payments-content",
			LoanID:           loanID,
			Table:            tableConfig,
			RecordPaymentURL: deps.Routes.AddURL,
			Labels:           deps.Labels,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "loan-payment"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("loan-payments", pageData)
	})
}

// NewContentView creates the loan payments HTMX partial view.
func NewContentView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		loanID := viewCtx.Request.URL.Query().Get("loan_id")
		payments := fetchPayments(ctx, deps, loanID)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, payments, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-credit-card",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:  "loan-payments-content",
			LoanID:           loanID,
			Table:            tableConfig,
			RecordPaymentURL: deps.Routes.AddURL,
			Labels:           deps.Labels,
		}

		return view.OK("loan-payments-content", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchPayments(ctx context.Context, deps *Deps, loanID string) []LoanPaymentRow {
	if deps.ListLoanPayments == nil {
		return []LoanPaymentRow{}
	}

	req := &loanpaymentpb.ListLoanPaymentsRequest{}
	resp, err := deps.ListLoanPayments(ctx, req)
	if err != nil {
		log.Printf("ListLoanPayments error: %v", err)
		return []LoanPaymentRow{}
	}
	if resp == nil {
		return []LoanPaymentRow{}
	}

	rows := make([]LoanPaymentRow, 0)
	for _, p := range resp.GetData() {
		row := protoToRow(p)
		// Filter by loan ID if provided
		if loanID != "" && row.LoanID != loanID {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func protoToRow(p *loanpaymentpb.LoanPayment) LoanPaymentRow {
	payDate := p.GetPaymentDate()

	return LoanPaymentRow{
		ID:               p.GetId(),
		LoanID:           p.GetLoanId(),
		PaymentNumber:    p.GetPaymentNumber(),
		PaymentDate:      payDate,
		PrincipalAmount:  fmt.Sprintf("%.2f", float64(p.GetPrincipalAmount())/100.0),
		InterestAmount:   fmt.Sprintf("%.2f", float64(p.GetInterestAmount())/100.0),
		FeeAmount:        fmt.Sprintf("%.2f", float64(p.GetFeeAmount())/100.0),
		TotalAmount:      fmt.Sprintf("%.2f", float64(p.GetTotalAmount())/100.0),
		RemainingBalance: fmt.Sprintf("%.2f", float64(p.GetRemainingBalance())/100.0),
		Notes:            p.GetNotes(),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, payments []LoanPaymentRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "payment_number", Label: l.Columns.PaymentNumber, Sortable: false, Width: "120px"},
		{Key: "date", Label: l.Columns.PaymentDate, Sortable: false, Width: "120px"},
		{Key: "principal", Label: l.Columns.PrincipalAmount, Sortable: false, Width: "130px", Align: "right"},
		{Key: "interest", Label: l.Columns.InterestAmount, Sortable: false, Width: "130px", Align: "right"},
		{Key: "total", Label: l.Columns.TotalAmount, Sortable: false, Width: "130px", Align: "right"},
		{Key: "balance", Label: l.Columns.RemainingBalance, Sortable: false, Width: "130px", Align: "right"},
	}

	rows := []types.TableRow{}
	for _, p := range payments {
		row := types.TableRow{
			ID: p.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: p.PaymentNumber},
				{Type: "text", Value: p.PaymentDate},
				{Type: "money", Value: p.PrincipalAmount},
				{Type: "money", Value: p.InterestAmount},
				{Type: "money", Value: p.TotalAmount},
				{Type: "money", Value: p.RemainingBalance},
			},
		}
		rows = append(rows, row)
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                "loan-payments-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        false,
		ShowActions:       false,
		ShowEntries:       true,
		DefaultSortColumn: "payment_number",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.RecordPayment,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("loan_payment", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}
