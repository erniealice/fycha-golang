package reports

import (
	"context"
	"fmt"
	"time"

	fycha "github.com/erniealice/fycha-golang"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ---------------------------------------------------------------------------
// Data model
// ---------------------------------------------------------------------------

// BSLine is one account line in the balance sheet.
type BSLine struct {
	Code        string // e.g. "1110"
	Name        string // e.g. "Cash on Hand"
	Amount      string // e.g. "₱45,200.00"
	IsNegative  bool   // true for contra-accounts (show in parentheses)
	IsSubtotal  bool   // classification subtotal (underlined)
	IsSeparator bool   // horizontal rule
}

// BSClassification is a classification sub-group (Current/Non-Current).
type BSClassification struct {
	Title    string
	Lines    []BSLine
	Subtotal string
}

// BSSection is a major element section (Assets, Liabilities, Equity).
type BSSection struct {
	Title           string
	Classifications []BSClassification // may be empty for Equity
	Lines           []BSLine           // direct lines for Equity section
	Total           string
	IsBold          bool
}

// ---------------------------------------------------------------------------
// Deps + PageData
// ---------------------------------------------------------------------------

// BalanceSheetDeps holds dependencies for the Balance Sheet view.
type BalanceSheetDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       fycha.ReportsLabels

	// GetBalanceSheet fetches balance sheet data as of the given date.
	// Phase 8: set to nil — mock data is used automatically.
	GetBalanceSheet func(ctx context.Context, asOfDate string) ([]BSSection, error)
}

// BalanceSheetPageData is the template data for the balance-sheet page.
type BalanceSheetPageData struct {
	types.PageData
	ContentTemplate string

	// Filter state
	AsOfDate string

	// KPI summary metrics
	TotalAssets      string
	TotalLiabilities string
	TotalEquity      string
	TotalLandE       string // Liabilities + Equity (must equal Assets)

	// Accounting equation verification
	IsBalanced      bool
	EquationMessage string

	// Statement body
	Sections []BSSection
}

// ---------------------------------------------------------------------------
// View constructor
// ---------------------------------------------------------------------------

// NewBalanceSheetView creates the Balance Sheet report view.
func NewBalanceSheetView(deps *BalanceSheetDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		q := viewCtx.QueryParams

		asOfDate := q["as_of"]
		if asOfDate == "" {
			now := time.Now()
			lastDay := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.UTC)
			asOfDate = lastDay.Format("2006-01-02")
		}

		// Fetch sections
		var sections []BSSection
		if deps.GetBalanceSheet != nil {
			ss, err := deps.GetBalanceSheet(ctx, asOfDate)
			if err == nil {
				sections = ss
			}
		}
		if sections == nil {
			sections = mockBSSections()
		}

		// Parse KPIs from sections
		totalAssets, totalLiab, totalEquity := calcBSKPIs(sections)
		totalLandE := totalLiab + totalEquity
		diff := totalAssets - totalLandE
		if diff < 0 {
			diff = -diff
		}
		isBalanced := diff < 0.01

		var equationMsg string
		if isBalanced {
			equationMsg = fmt.Sprintf("A = L + E verified: %s = %s + %s",
				formatCurrencyFS(totalAssets),
				formatCurrencyFS(totalLiab),
				formatCurrencyFS(totalEquity),
			)
		} else {
			equationMsg = fmt.Sprintf("Warning: Assets (%s) ≠ Liabilities + Equity (%s). Difference: %s",
				formatCurrencyFS(totalAssets),
				formatCurrencyFS(totalLandE),
				formatCurrencyFS(diff),
			)
		}

		pageData := &BalanceSheetPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.BalanceSheet.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "report",
				ActiveSubNav:   "balance-sheet",
				HeaderTitle:    deps.Labels.BalanceSheet.Title,
				HeaderSubtitle: deps.Labels.BalanceSheet.Subtitle,
				HeaderIcon:     "icon-layers",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:  "balance-sheet-content",
			AsOfDate:         asOfDate,
			TotalAssets:      formatCurrencyFS(totalAssets),
			TotalLiabilities: formatCurrencyFS(totalLiab),
			TotalEquity:      formatCurrencyFS(totalEquity),
			TotalLandE:       formatCurrencyFS(totalLandE),
			IsBalanced:       isBalanced,
			EquationMessage:  equationMsg,
			Sections:         sections,
		}

		if viewCtx.IsHTMX {
			return view.OK("balance-sheet-content", pageData)
		}
		return view.OK("balance-sheet", pageData)
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func calcBSKPIs(sections []BSSection) (totalAssets, totalLiab, totalEquity float64) {
	for _, s := range sections {
		switch s.Title {
		case "ASSETS":
			totalAssets = parseISAmount(s.Total)
		case "LIABILITIES":
			totalLiab = parseISAmount(s.Total)
		case "EQUITY":
			totalEquity = parseISAmount(s.Total)
		}
	}
	return
}

// ---------------------------------------------------------------------------
// Mock data (Phase 8)
// ---------------------------------------------------------------------------

// mockBSSections returns a realistic balance sheet for a Philippine salon/spa.
// Assets = ₱1,245,800 = Liabilities (₱548,200) + Equity (₱697,600)
func mockBSSections() []BSSection {
	return []BSSection{
		{
			Title: "ASSETS",
			Classifications: []BSClassification{
				{
					Title: "Current Assets",
					Lines: []BSLine{
						{Code: "1010", Name: "Cash in Bank", Amount: "₱182,500.00"},
						{Code: "1020", Name: "Petty Cash", Amount: "₱9,000.00"},
						{Code: "1030", Name: "Cash on Hand (GCash/Maya)", Amount: "₱15,000.00"},
						{Code: "1110", Name: "Accounts Receivable", Amount: "₱125,000.00"},
						{Code: "1150", Name: "Allowance for Doubtful Accounts", Amount: "(₱5,000.00)", IsNegative: true},
						{Code: "1210", Name: "Salon Supplies", Amount: "₱18,500.00"},
						{Code: "1220", Name: "Retail Products for Resale", Amount: "₱24,000.00"},
						{Code: "1310", Name: "Prepaid Rent", Amount: "₱15,000.00"},
						{Code: "1320", Name: "Prepaid Insurance", Amount: "₱24,000.00"},
					},
					Subtotal: "₱408,000.00",
				},
				{
					Title: "Non-Current Assets",
					Lines: []BSLine{
						{Code: "1510", Name: "Salon Chairs & Stations", Amount: "₱120,000.00"},
						{Code: "1515", Name: "Accum. Depreciation — Salon Chairs", Amount: "(₱30,000.00)", IsNegative: true},
						{Code: "1520", Name: "Spa Equipment", Amount: "₱250,000.00"},
						{Code: "1525", Name: "Accum. Depreciation — Spa Equipment", Amount: "(₱62,500.00)", IsNegative: true},
						{Code: "1530", Name: "Furniture & Fixtures", Amount: "₱85,000.00"},
						{Code: "1535", Name: "Accum. Depreciation — Furniture", Amount: "(₱21,250.00)", IsNegative: true},
						{Code: "1540", Name: "Leasehold Improvements", Amount: "₱510,000.00"},
						{Code: "1545", Name: "Accum. Amortization — Leasehold", Amount: "(₱13,450.00)", IsNegative: true},
					},
					Subtotal: "₱837,800.00",
				},
			},
			Total:  "₱1,245,800.00",
			IsBold: true,
		},
		{
			Title: "LIABILITIES",
			Classifications: []BSClassification{
				{
					Title: "Current Liabilities",
					Lines: []BSLine{
						{Code: "2010", Name: "Accounts Payable", Amount: "₱48,200.00"},
						{Code: "2110", Name: "Accrued Salaries & Wages", Amount: "₱85,000.00"},
						{Code: "2120", Name: "Accrued SSS/PhilHealth/HDMF", Amount: "₱12,400.00"},
						{Code: "2210", Name: "Output VAT / Percentage Tax", Amount: "₱28,000.00"},
						{Code: "2220", Name: "Withholding Tax Payable", Amount: "₱12,000.00"},
						{Code: "2230", Name: "Income Tax Payable", Amount: "₱5,600.00"},
					},
					Subtotal: "₱191,200.00",
				},
				{
					Title: "Non-Current Liabilities",
					Lines: []BSLine{
						{Code: "2510", Name: "Bank Loan — BDO", Amount: "₱200,000.00"},
						{Code: "2520", Name: "Equipment Finance Payable", Amount: "₱157,000.00"},
					},
					Subtotal: "₱357,000.00",
				},
			},
			Total:  "₱548,200.00",
			IsBold: true,
		},
		{
			Title: "EQUITY",
			Lines: []BSLine{
				{Code: "3010", Name: "Owner's Capital — Maria Santos", Amount: "₱500,000.00"},
				{Code: "3011", Name: "Owner's Capital — Juan dela Cruz", Amount: "₱300,000.00"},
				{Code: "3020", Name: "Owner's Draw — Maria Santos", Amount: "(₱50,000.00)", IsNegative: true},
				{Code: "3030", Name: "Retained Earnings (Prior Periods)", Amount: "(₱120,000.00)", IsNegative: true},
				{Code: "3031", Name: "Current Year Earnings", Amount: "₱67,600.00"},
			},
			Total:  "₱697,600.00",
			IsBold: true,
		},
	}
}
