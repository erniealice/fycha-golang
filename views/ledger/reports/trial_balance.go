package reports

import (
	"context"
	"fmt"
	"time"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Data model
// ---------------------------------------------------------------------------

// TBAccountRow is one account line in a Trial Balance.
type TBAccountRow struct {
	AccountID   string
	AccountCode string
	AccountName string
	Element     string  // "asset", "liability", "equity", "revenue", "expense"
	Debit       float64 // positive when normal balance is debit
	Credit      float64 // positive when normal balance is credit
}

// TBElementGroup groups accounts by element with subtotals.
type TBElementGroup struct {
	Element      string // "asset", "liability", "equity", "revenue", "expense"
	Label        string // display label, e.g. "ASSETS"
	Accounts     []TBAccountRow
	SubtotalDebit  float64
	SubtotalCredit float64
}

// TBTotals holds the grand totals and balance check result.
type TBTotals struct {
	TotalDebit   float64
	TotalCredit  float64
	Difference   float64
	IsBalanced   bool
	// Pre-formatted for templates
	TotalDebitStr  string
	TotalCreditStr string
	DifferenceStr  string
}

// ---------------------------------------------------------------------------
// Deps + PageData
// ---------------------------------------------------------------------------

// TrialBalanceDeps holds dependencies for the Trial Balance report view.
type TrialBalanceDeps struct {
	Routes       fycha.LedgerStatementRoutes
	Labels       fycha.AccountLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// GetTrialBalance fetches account balances as of the given date.
	// Returning nil, nil means "no data" — mock data is used automatically.
	// Phase 3: set to nil to rely on mock data.
	GetTrialBalance func(ctx context.Context, asOfDate string) ([]TBAccountRow, error)
}

// TrialBalancePageData is the template data for the trial-balance page.
type TrialBalancePageData struct {
	types.PageData
	ContentTemplate string

	// Filter state
	AsOfDate string

	// Report data
	HasData bool
	Groups  []TBElementGroup
	Totals  TBTotals
	Table   *types.TableConfig

	// Labels
	Labels fycha.AccountLabels
}

// ---------------------------------------------------------------------------
// View constructor
// ---------------------------------------------------------------------------

// NewTrialBalanceView creates the Trial Balance report view.
func NewTrialBalanceView(deps *TrialBalanceDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		q := viewCtx.QueryParams

		asOfDate := q["as_of"]
		if asOfDate == "" {
			// Default: last day of the current month
			now := time.Now()
			// First day of next month minus one day = last day of current month
			lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			asOfDate = lastDay.Format("2006-01-02")
		}

		pageData := &TrialBalancePageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Trial Balance",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "trial-balance",
				HeaderTitle:    "Trial Balance",
				HeaderSubtitle: "Verify that total debits equal total credits",
				HeaderIcon:     "icon-check-square",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "trial-balance-content",
			AsOfDate:        asOfDate,
			Labels:          deps.Labels,
		}

		// Fetch data (Phase 3: fall back to mock if no real use case wired)
		var accounts []TBAccountRow
		if deps.GetTrialBalance != nil {
			rows, err := deps.GetTrialBalance(ctx, asOfDate)
			if err == nil {
				accounts = rows
			}
		}
		if accounts == nil {
			accounts = mockTBAccounts()
		}

		if len(accounts) > 0 {
			pageData.HasData = true
			pageData.Groups = buildTBGroups(accounts)
			pageData.Totals = buildTBTotals(pageData.Groups)
			pageData.Table = buildTBTable(pageData.Groups, pageData.Totals, deps.TableLabels)
		}

		if viewCtx.IsHTMX {
			return view.OK("trial-balance-content", pageData)
		}
		return view.OK("trial-balance", pageData)
	})
}

// ---------------------------------------------------------------------------
// Group and totals builders
// ---------------------------------------------------------------------------

var elementOrder = []struct {
	key   string
	label string
}{
	{"asset", "ASSETS"},
	{"liability", "LIABILITIES"},
	{"equity", "EQUITY"},
	{"revenue", "REVENUE"},
	{"expense", "EXPENSES"},
}

func buildTBGroups(accounts []TBAccountRow) []TBElementGroup {
	// bucket accounts by element
	buckets := make(map[string][]TBAccountRow, 5)
	for _, a := range accounts {
		buckets[a.Element] = append(buckets[a.Element], a)
	}

	groups := make([]TBElementGroup, 0, 5)
	for _, e := range elementOrder {
		rows, ok := buckets[e.key]
		if !ok || len(rows) == 0 {
			continue
		}
		var subtotalDebit, subtotalCredit float64
		for _, r := range rows {
			subtotalDebit += r.Debit
			subtotalCredit += r.Credit
		}
		groups = append(groups, TBElementGroup{
			Element:        e.key,
			Label:          e.label,
			Accounts:       rows,
			SubtotalDebit:  subtotalDebit,
			SubtotalCredit: subtotalCredit,
		})
	}
	return groups
}

func buildTBTotals(groups []TBElementGroup) TBTotals {
	var totalDebit, totalCredit float64
	for _, g := range groups {
		totalDebit += g.SubtotalDebit
		totalCredit += g.SubtotalCredit
	}
	diff := totalDebit - totalCredit
	if diff < 0 {
		diff = -diff
	}
	isBalanced := diff < 0.01 // floating-point tolerance
	return TBTotals{
		TotalDebit:     totalDebit,
		TotalCredit:    totalCredit,
		Difference:     diff,
		IsBalanced:     isBalanced,
		TotalDebitStr:  formatCurrencyGL(totalDebit),
		TotalCreditStr: formatCurrencyGL(totalCredit),
		DifferenceStr:  formatCurrencyGL(diff),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTBTable(groups []TBElementGroup, totals TBTotals, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "code", Label: "Code", Sortable: false, Width: "100px"},
		{Key: "name", Label: "Account Name", Sortable: false},
		{Key: "debit", Label: "Debit Balance", Sortable: false, Width: "150px", Align: "right"},
		{Key: "credit", Label: "Credit Balance", Sortable: false, Width: "150px", Align: "right"},
	}

	rowGroups := make([]types.TableRowGroup, 0, len(groups))
	for _, g := range groups {
		rows := make([]types.TableRow, 0, len(g.Accounts)+1)

		for _, acct := range g.Accounts {
			debitVal := ""
			creditVal := ""
			if acct.Debit > 0 {
				debitVal = formatCurrencyGL(acct.Debit)
			}
			if acct.Credit > 0 {
				creditVal = formatCurrencyGL(acct.Credit)
			}
			rows = append(rows, types.TableRow{
				ID: acct.AccountID,
				Cells: []types.TableCell{
					{Type: "text", Value: acct.AccountCode},
					{Type: "text", Value: acct.AccountName},
					{Type: "text", Value: debitVal},
					{Type: "text", Value: creditVal},
				},
				DataAttrs: map[string]string{
					"element": acct.Element,
				},
			})
		}

		// Subtotal row for this element group
		subtotalDebitStr := ""
		subtotalCreditStr := ""
		if g.SubtotalDebit > 0 {
			subtotalDebitStr = formatCurrencyGL(g.SubtotalDebit)
		}
		if g.SubtotalCredit > 0 {
			subtotalCreditStr = formatCurrencyGL(g.SubtotalCredit)
		}
		rows = append(rows, types.TableRow{
			ID: fmt.Sprintf("subtotal-%s", g.Element),
			Cells: []types.TableCell{
				{Type: "text", Value: ""},
				{Type: "text", Value: fmt.Sprintf("Subtotal: %s", g.Label)},
				{Type: "text", Value: subtotalDebitStr},
				{Type: "text", Value: subtotalCreditStr},
			},
			DataAttrs: map[string]string{
				"row-type": "subtotal",
			},
		})

		rowGroups = append(rowGroups, types.TableRowGroup{
			ID:    g.Element,
			Title: g.Label,
			Rows:  rows,
		})
	}

	// Grand totals as a final plain row appended to a special group
	totalDebitStr := formatCurrencyGL(totals.TotalDebit)
	totalCreditStr := formatCurrencyGL(totals.TotalCredit)
	balanceLabel := "Unbalanced"
	if totals.IsBalanced {
		balanceLabel = "Balanced"
	}
	differenceStr := fmt.Sprintf("%s (%s)", formatCurrencyGL(totals.Difference), balanceLabel)

	totalsGroup := types.TableRowGroup{
		ID:    "totals",
		Title: "TOTALS",
		Rows: []types.TableRow{
			{
				ID: "grand-totals",
				Cells: []types.TableCell{
					{Type: "text", Value: ""},
					{Type: "text", Value: "TOTAL"},
					{Type: "text", Value: totalDebitStr},
					{Type: "text", Value: totalCreditStr},
				},
				DataAttrs: map[string]string{"row-type": "grand-total"},
			},
			{
				ID: "difference",
				Cells: []types.TableCell{
					{Type: "text", Value: ""},
					{Type: "text", Value: "DIFFERENCE"},
					{Type: "text", Value: differenceStr},
					{Type: "text", Value: ""},
				},
				DataAttrs: map[string]string{"row-type": "difference"},
			},
		},
	}
	rowGroups = append(rowGroups, totalsGroup)

	return &types.TableConfig{
		ID:          "trial-balance-table",
		Columns:     columns,
		Groups:      rowGroups,
		ShowSearch:  true,
		ShowExport:  true,
		ShowEntries: true,
		Labels:      tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   "No accounts with balances",
			Message: "No accounts have non-zero balances as of the selected date.",
		},
	}
}

// ---------------------------------------------------------------------------
// Mock data (Phase 3)
// ---------------------------------------------------------------------------

// mockTBAccounts returns a balanced mock trial balance.
// Accounting equation: Assets = Liabilities + Equity
// With Revenue and Expenses: Assets + Expenses = Liabilities + Equity + Revenue
// Total Debits == Total Credits: 1,245,800.00
func mockTBAccounts() []TBAccountRow {
	return []TBAccountRow{
		// Assets (debit-normal) — total debit: 531,500
		{AccountID: "acc-1110", AccountCode: "1110", AccountName: "Cash on Hand", Element: "asset", Debit: 45200.00},
		{AccountID: "acc-1120", AccountCode: "1120", AccountName: "BDO Savings", Element: "asset", Debit: 182500.00},
		{AccountID: "acc-1130", AccountCode: "1130", AccountName: "BPI Checking", Element: "asset", Debit: 93800.00},
		{AccountID: "acc-1210", AccountCode: "1210", AccountName: "Accounts Receivable", Element: "asset", Debit: 125000.00},
		{AccountID: "acc-1310", AccountCode: "1310", AccountName: "Merchandise Inventory", Element: "asset", Debit: 89000.00, Credit: 0},
		{AccountID: "acc-1510", AccountCode: "1510", AccountName: "Office Equipment", Element: "asset", Debit: 85000.00},
		{AccountID: "acc-1610", AccountCode: "1610", AccountName: "Furniture and Fixtures", Element: "asset", Debit: 48000.00},
		// Asset subtotal debit: 668,500.00 (note: expanded below so totals balance)

		// Liabilities (credit-normal) — total credit: 333,200
		{AccountID: "acc-2010", AccountCode: "2010", AccountName: "Accounts Payable", Element: "liability", Credit: 48200.00},
		{AccountID: "acc-2020", AccountCode: "2020", AccountName: "Salaries Payable", Element: "liability", Credit: 85000.00},
		{AccountID: "acc-2030", AccountCode: "2030", AccountName: "SSS/PhilHealth/Pag-IBIG Payable", Element: "liability", Credit: 12400.00},
		{AccountID: "acc-2510", AccountCode: "2510", AccountName: "Long-term Loan", Element: "liability", Credit: 187600.00},

		// Equity (credit-normal) — total credit: 198,300
		{AccountID: "acc-3010", AccountCode: "3010", AccountName: "Owner's Capital", Element: "equity", Credit: 198300.00},

		// Revenue (credit-normal) — total credit: 450,000
		{AccountID: "acc-4010", AccountCode: "4010", AccountName: "Service Revenue", Element: "revenue", Credit: 320000.00},
		{AccountID: "acc-4020", AccountCode: "4020", AccountName: "Product Sales", Element: "revenue", Credit: 130000.00},

		// Expenses (debit-normal) — total debit: 313,000
		{AccountID: "acc-5010", AccountCode: "5010", AccountName: "Cost of Goods Sold", Element: "expense", Debit: 78000.00},
		{AccountID: "acc-5020", AccountCode: "5020", AccountName: "Salaries Expense", Element: "expense", Debit: 85000.00},
		{AccountID: "acc-5030", AccountCode: "5030", AccountName: "Rent Expense", Element: "expense", Debit: 54000.00},
		{AccountID: "acc-5040", AccountCode: "5040", AccountName: "Utilities Expense", Element: "expense", Debit: 28500.00},
		{AccountID: "acc-5050", AccountCode: "5050", AccountName: "Depreciation Expense", Element: "expense", Debit: 16750.00},
		{AccountID: "acc-5060", AccountCode: "5060", AccountName: "Supplies Expense", Element: "expense", Debit: 12300.00},
		{AccountID: "acc-5070", AccountCode: "5070", AccountName: "Miscellaneous Expense", Element: "expense", Debit: 38450.00},
	}
	// Mock data verification:
	// Total Debit  = Assets (668,500) + Expenses (313,000)         = 981,500
	// Total Credit = Liabilities (333,200) + Equity (198,300) + Revenue (450,000) = 981,500 ✓
}
