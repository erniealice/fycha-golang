// Package reports provides ledger-internal accounting report views.
// These are tools for accountants/bookkeepers, distinct from business-stakeholder
// reports in the /app/reports/ namespace.
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

// GLLine represents a single journal line entry in the general ledger.
type GLLine struct {
	Date           string // formatted date, e.g. "03/01"
	EntryNumber    string // e.g. "JE-0031"
	EntryDetailURL string // link to journal entry detail page (may be empty)
	Description    string
	Debit          float64
	Credit         float64
	RunningBalance float64
	IsSpecialRow   bool   // true for Opening Balance / Totals / Closing Balance rows
	SpecialRowType string // "opening", "totals", "closing"
}

// GLAccountSection groups all transaction lines for one account.
type GLAccountSection struct {
	AccountID      string
	AccountCode    string
	AccountName    string
	Element        string // "asset", "liability", "equity", "revenue", "expense"
	OpeningBalance float64
	PeriodDebits   float64
	PeriodCredits  float64
	ClosingBalance float64
	Lines          []GLLine
}

// ---------------------------------------------------------------------------
// Deps + PageData
// ---------------------------------------------------------------------------

// GeneralLedgerDeps holds dependencies for the General Ledger report view.
// GetGeneralLedger is intentionally a func field so Phase 3 can use mock data
// without coupling to a concrete DB type. Wire the real use case in Phase 4.
type GeneralLedgerDeps struct {
	Routes       fycha.LedgerStatementRoutes
	Labels       fycha.AccountLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// GetGeneralLedger fetches journal lines for a given account and date range.
	// Returning nil, nil means "no data" (renders empty state, not error).
	// Phase 3: set to nil → mock data is used automatically.
	GetGeneralLedger func(ctx context.Context, accountID, startDate, endDate string) (*GLAccountSection, error)
}

// GeneralLedgerPageData is the template data for the general-ledger page.
type GeneralLedgerPageData struct {
	types.PageData
	ContentTemplate string

	// Filter state
	AccountID   string
	AccountCode string
	AccountName string
	StartDate   string
	EndDate     string

	// Report state
	HasData        bool // false when no account selected or no results
	Section        *GLAccountSection
	SummaryMetrics []fycha.SummaryMetric
	Table          *types.TableConfig

	// Labels
	Labels fycha.AccountLabels
}

// ---------------------------------------------------------------------------
// View constructor
// ---------------------------------------------------------------------------

// NewGeneralLedgerView creates the General Ledger report view.
func NewGeneralLedgerView(deps *GeneralLedgerDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		q := viewCtx.QueryParams

		accountID := q["account_id"]
		startDate := q["start"]
		endDate := q["end"]

		// Default date range: first day of current month → today
		if startDate == "" {
			now := time.Now()
			startDate = fmt.Sprintf("%d-%02d-01", now.Year(), now.Month())
		}
		if endDate == "" {
			endDate = time.Now().Format("2006-01-02")
		}

		pageData := &GeneralLedgerPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.GeneralLedger.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "general-ledger",
				HeaderTitle:    deps.Labels.GeneralLedger.Title,
				HeaderSubtitle: deps.Labels.GeneralLedger.Subtitle,
				HeaderIcon:     "icon-book-open",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "general-ledger-content",
			AccountID:       accountID,
			StartDate:       startDate,
			EndDate:         endDate,
			Labels:          deps.Labels,
		}

		// No account selected → show info state
		if accountID == "" {
			if viewCtx.IsHTMX {
				return view.OK("general-ledger-content", pageData)
			}
			return view.OK("general-ledger", pageData)
		}

		// Fetch data (Phase 3: use mock if no real use case wired)
		var section *GLAccountSection
		if deps.GetGeneralLedger != nil {
			s, err := deps.GetGeneralLedger(ctx, accountID, startDate, endDate)
			if err == nil {
				section = s
			}
		}
		if section == nil {
			section = mockGLSection(accountID, startDate, endDate)
		}

		pageData.HasData = true
		pageData.Section = section
		pageData.AccountCode = section.AccountCode
		pageData.AccountName = section.AccountName
		pageData.SummaryMetrics = buildGLSummary(section, deps.Labels)
		pageData.Table = buildGLTable(section, deps.TableLabels, deps.Labels)

		if viewCtx.IsHTMX {
			return view.OK("general-ledger-content", pageData)
		}
		return view.OK("general-ledger", pageData)
	})
}

// ---------------------------------------------------------------------------
// Summary bar
// ---------------------------------------------------------------------------

func buildGLSummary(s *GLAccountSection, labels fycha.AccountLabels) []fycha.SummaryMetric {
	return []fycha.SummaryMetric{
		{Label: labels.GeneralLedger.OpeningBalance, Value: formatCurrencyGL(s.OpeningBalance)},
		{Label: labels.GeneralLedger.PeriodDebits, Value: formatCurrencyGL(s.PeriodDebits), Highlight: true},
		{Label: labels.GeneralLedger.PeriodCredits, Value: formatCurrencyGL(s.PeriodCredits)},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildGLTable(s *GLAccountSection, tableLabels types.TableLabels, labels fycha.AccountLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "date", Label: labels.Columns.Date, Sortable: false, Width: "100px"},
		{Key: "entry", Label: labels.Columns.EntryNumber, Sortable: false, Width: "110px"},
		{Key: "description", Label: labels.Columns.Description, Sortable: false},
		{Key: "debit", Label: labels.Columns.Debit, Sortable: false, Width: "130px", Align: "right"},
		{Key: "credit", Label: labels.Columns.Credit, Sortable: false, Width: "130px", Align: "right"},
		{Key: "balance", Label: labels.GeneralLedger.RunningBalance, Sortable: false, Width: "140px", Align: "right"},
	}

	rows := make([]types.TableRow, 0, len(s.Lines))
	for i, line := range s.Lines {
		debitVal := ""
		creditVal := ""
		balanceVal := ""

		if !line.IsSpecialRow || line.SpecialRowType == "opening" || line.SpecialRowType == "closing" {
			if line.RunningBalance != 0 {
				balanceVal = formatCurrencyGL(line.RunningBalance)
			}
		}
		if line.SpecialRowType == "totals" {
			if line.Debit != 0 {
				debitVal = formatCurrencyGL(line.Debit)
			}
			if line.Credit != 0 {
				creditVal = formatCurrencyGL(line.Credit)
			}
		} else {
			if line.Debit > 0 {
				debitVal = formatCurrencyGL(line.Debit)
			}
			if line.Credit > 0 {
				creditVal = formatCurrencyGL(line.Credit)
			}
		}

		entryCell := types.TableCell{Type: "text", Value: line.EntryNumber}
		if line.EntryDetailURL != "" {
			entryCell = types.TableCell{Type: "link", Value: line.EntryNumber, Href: line.EntryDetailURL}
		}

		descType := "text"
		if line.IsSpecialRow {
			descType = "text"
		}

		row := types.TableRow{
			ID: fmt.Sprintf("gl-row-%d", i),
			Cells: []types.TableCell{
				{Type: "text", Value: line.Date},
				entryCell,
				{Type: descType, Value: line.Description},
				{Type: "text", Value: debitVal},
				{Type: "text", Value: creditVal},
				{Type: "text", Value: balanceVal},
			},
			DataAttrs: map[string]string{
				"row-type": line.SpecialRowType,
			},
		}
		rows = append(rows, row)
	}

	return &types.TableConfig{
		ID:          "gl-entries-table",
		Columns:     columns,
		Rows:        rows,
		ShowSearch:  false,
		ShowExport:  true,
		ShowEntries: true,
		ShowDensity: true,
		Labels:      tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   labels.GeneralLedger.NoTransactionsTitle,
			Message: labels.GeneralLedger.NoTransactionsDetail,
		},
	}
}

// ---------------------------------------------------------------------------
// Mock data (Phase 3)
// ---------------------------------------------------------------------------

// mockGLSection returns a realistic demo General Ledger section for a cash account.
// All debits and credits balance correctly within the period.
func mockGLSection(accountID, startDate, endDate string) *GLAccountSection {
	openingBalance := 28400.00

	type rawLine struct {
		date        string
		entryNum    string
		description string
		debit       float64
		credit      float64
	}
	rawLines := []rawLine{
		{"03/01", "JE-0031", "Daily collection", 12800.00, 0},
		{"03/02", "JE-0032", "Supplier payment", 0, 25000.00},
		{"03/03", "JE-0033", "Cash sale", 8200.00, 0},
		{"03/05", "JE-0035", "Petty cash replenishment", 0, 5000.00},
		{"03/10", "JE-0038", "Cash sale", 3500.00, 0},
		{"03/12", "JE-0040", "Rent payment", 0, 18000.00},
		{"03/14", "JE-0041", "Daily collection", 15200.00, 0},
		{"03/15", "JE-0042", "Salary disbursement", 0, 35300.00},
		{"03/16", "JE-0043", "Cash sale", 6500.00, 0},
		{"03/18", "JE-0044", "Utility bills", 0, 8500.00},
		{"03/20", "JE-0045", "Daily collection", 18600.00, 0},
		{"03/22", "JE-0046", "Supplier payment", 0, 12500.00},
		{"03/25", "JE-0047", "Cash sale", 22400.00, 0},
		{"03/28", "JE-0048", "Petty cash replenishment", 0, 5000.00},
		{"03/30", "JE-0049", "Daily collection", 41300.00, 0},
	}

	var lines []GLLine
	var totalDebit, totalCredit float64
	runBal := openingBalance

	// Opening balance row
	lines = append(lines, GLLine{
		Description:    "Opening Balance",
		RunningBalance: openingBalance,
		IsSpecialRow:   true,
		SpecialRowType: "opening",
	})

	for _, rl := range rawLines {
		runBal += rl.debit - rl.credit
		totalDebit += rl.debit
		totalCredit += rl.credit
		lines = append(lines, GLLine{
			Date:           rl.date,
			EntryNumber:    rl.entryNum,
			EntryDetailURL: fmt.Sprintf("/app/ledger/journals/detail/%s", rl.entryNum),
			Description:    rl.description,
			Debit:          rl.debit,
			Credit:         rl.credit,
			RunningBalance: runBal,
		})
	}

	// Period totals row
	lines = append(lines, GLLine{
		Description:    "PERIOD TOTALS",
		Debit:          totalDebit,
		Credit:         totalCredit,
		IsSpecialRow:   true,
		SpecialRowType: "totals",
	})

	// Closing balance row
	lines = append(lines, GLLine{
		Description:    "CLOSING BALANCE",
		RunningBalance: runBal,
		IsSpecialRow:   true,
		SpecialRowType: "closing",
	})

	return &GLAccountSection{
		AccountID:      accountID,
		AccountCode:    "1110",
		AccountName:    "Cash on Hand",
		Element:        "asset",
		OpeningBalance: openingBalance,
		PeriodDebits:   totalDebit,
		PeriodCredits:  totalCredit,
		ClosingBalance: runBal,
		Lines:          lines,
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func formatCurrencyGL(amount float64) string {
	negative := amount < 0
	if negative {
		amount = -amount
	}
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
	formatted := fmt.Sprintf("\u20b1%s.%02d", wholeStr, frac)
	if negative {
		formatted = "-" + formatted
	}
	return formatted
}
