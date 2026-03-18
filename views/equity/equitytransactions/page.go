// Package equitytransactions provides the view for the Funding > Equity > Transactions list page.
package equitytransactions

import (
	"context"
	"fmt"
	"log"
	"time"

	equitytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_transaction"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the equity transactions list.
type Deps struct {
	Routes       fycha.EquityRoutes
	Labels       fycha.EquityLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Equity transaction use cases
	ListEquityTransactions func(ctx context.Context, req *equitytransactionpb.ListEquityTransactionsRequest) (*equitytransactionpb.ListEquityTransactionsResponse, error)
}

// PageData holds the data for the equity transactions list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	// Guided transaction form data
	AddFormURL string
}

// TransactionRow is the view-model for a single equity transaction row.
type TransactionRow struct {
	ID              string
	EquityAccountID string
	TransactionType string
	Amount          string
	Description     string
	TransactionDate string
	JournalEntryID  string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the equity transactions list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		txns := fetchTransactions(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, txns, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Transactions.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    deps.Labels.Transactions.Page.Heading,
				HeaderSubtitle: deps.Labels.Transactions.Page.Caption,
				HeaderIcon:     "icon-arrows-right-left",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "equity-transactions-content",
			Table:           tableConfig,
			AddFormURL:      deps.Routes.TransactionAddURL,
		}

		return view.OK("equity-transactions", pageData)
	})
}

// NewContentView creates the equity transactions list HTMX partial view.
func NewContentView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		txns := fetchTransactions(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, txns, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Transactions.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				HeaderTitle:    deps.Labels.Transactions.Page.Heading,
				HeaderSubtitle: deps.Labels.Transactions.Page.Caption,
				HeaderIcon:     "icon-arrows-right-left",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "equity-transactions-content",
			Table:           tableConfig,
			AddFormURL:      deps.Routes.TransactionAddURL,
		}

		return view.OK("equity-transactions-content", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchTransactions(ctx context.Context, deps *Deps) []TransactionRow {
	if deps.ListEquityTransactions == nil {
		return []TransactionRow{}
	}

	resp, err := deps.ListEquityTransactions(ctx, &equitytransactionpb.ListEquityTransactionsRequest{})
	if err != nil {
		log.Printf("ListEquityTransactions error: %v", err)
		return []TransactionRow{}
	}
	if resp == nil {
		return []TransactionRow{}
	}

	rows := make([]TransactionRow, 0, len(resp.GetData()))
	for _, t := range resp.GetData() {
		rows = append(rows, protoToRow(t))
	}
	return rows
}

func protoToRow(t *equitytransactionpb.EquityTransaction) TransactionRow {
	dateStr := t.GetTransactionDateString()
	if dateStr == "" && t.GetTransactionDate() > 0 {
		dateStr = time.UnixMilli(t.GetTransactionDate()).Format("2006-01-02")
	}

	jeID := ""
	if t.JournalEntryId != nil {
		jeID = *t.JournalEntryId
	}

	return TransactionRow{
		ID:              t.GetId(),
		EquityAccountID: t.GetEquityAccountId(),
		TransactionType: transactionTypeLabel(t.GetTransactionType()),
		Amount:          fmt.Sprintf("%.2f", t.GetAmount()),
		Description:     t.GetDescription(),
		TransactionDate: dateStr,
		JournalEntryID:  jeID,
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, txns []TransactionRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels.Transactions
	columns := []types.TableColumn{
		{Key: "date", Label: l.Columns.Date, Sortable: false, Width: "120px"},
		{Key: "type", Label: l.Columns.TransactionType, Sortable: false, Width: "160px"},
		{Key: "amount", Label: l.Columns.Amount, Sortable: false, Width: "140px", Align: "right"},
		{Key: "description", Label: l.Columns.Description, Sortable: false},
	}

	rows := []types.TableRow{}
	for _, txn := range txns {
		typeVariant := transactionTypeBadgeVariant(txn.TransactionType)

		row := types.TableRow{
			ID: txn.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: txn.TransactionDate},
				{Type: "badge", Value: txn.TransactionType, Variant: typeVariant},
				{Type: "money", Value: txn.Amount},
				{Type: "text", Value: txn.Description},
			},
		}
		rows = append(rows, row)
	}

	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                "equity-transactions-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowEntries:       true,
		DefaultSortColumn: "date",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.RecordTransaction,
			ActionURL:       deps.Routes.TransactionAddURL,
			Icon:            "icon-plus",
			Disabled:        !perms.Can("equity_transaction", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Proto enum helpers
// ---------------------------------------------------------------------------

func transactionTypeLabel(t equitytransactionpb.EquityTransactionType) string {
	switch t {
	case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_CONTRIBUTION:
		return "Contribution"
	case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_WITHDRAWAL:
		return "Withdrawal"
	case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_DISTRIBUTION:
		return "Distribution"
	case equitytransactionpb.EquityTransactionType_EQUITY_TRANSACTION_TYPE_TRANSFER:
		return "Transfer"
	default:
		return "Unspecified"
	}
}

func transactionTypeBadgeVariant(label string) string {
	switch label {
	case "Contribution":
		return "sage"
	case "Withdrawal", "Distribution":
		return "amber"
	case "Transfer":
		return "navy"
	default:
		return "default"
	}
}
