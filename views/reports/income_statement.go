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

// ISStatementLine is one account line in the income statement.
type ISStatementLine struct {
	Code          string // e.g. "4010"
	Name          string // e.g. "Service Revenue"
	CurrentPeriod string // e.g. "₱380,000.00"
	PriorPeriod   string // e.g. "₱345,000.00"
	Change        string // e.g. "+10.1%"
	IsTotal       bool   // true for bold total lines (Gross Profit, Net Income)
	IsSeparator   bool   // horizontal rule between sections
	IsNegative    bool   // true when amount should be styled as negative
}

// ISStatementGroup is a sub-group within a section (e.g. Selling / G&A).
type ISStatementGroup struct {
	Title    string
	Lines    []ISStatementLine
	Subtotal string
}

// ISStatementSection is a major section of the income statement.
type ISStatementSection struct {
	Title    string
	Lines    []ISStatementLine
	Subtotal string
	Bold     bool // major total line (Gross Profit, Net Income)
	Groups   []ISStatementGroup
}

// ---------------------------------------------------------------------------
// Deps + PageData
// ---------------------------------------------------------------------------

// IncomeStatementDeps holds dependencies for the Income Statement view.
type IncomeStatementDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       fycha.ReportsLabels

	// GetIncomeStatement fetches income statement data for the period.
	// Phase 8: set to nil — mock data is used automatically.
	GetIncomeStatement func(ctx context.Context, startDate, endDate string) ([]ISStatementSection, error)
}

// IncomeStatementPageData is the template data for the income-statement page.
type IncomeStatementPageData struct {
	types.PageData
	ContentTemplate string

	// Period filter state
	ActivePreset  string
	StartDate     string
	EndDate       string
	PeriodLabel   string
	PeriodPresets []fycha.FilterOption

	// KPI summary metrics
	TotalRevenue     string
	TotalExpenses    string
	NetIncome        string
	NetIncomeVariant string // "success" or "danger"
	NetIncomeTrend   string // "+12%"

	// Statement body
	Sections []ISStatementSection
}

// ---------------------------------------------------------------------------
// View constructor
// ---------------------------------------------------------------------------

// NewIncomeStatementView creates the Income Statement report view.
func NewIncomeStatementView(deps *IncomeStatementDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		q := viewCtx.QueryParams

		preset := q["period"]
		if preset == "" {
			preset = "thisMonth"
		}
		startDate := q["start"]
		endDate := q["end"]

		// Resolve date range from preset
		start, end := fycha.ParsePeriodPreset(preset)
		if preset == "custom" {
			if t, err := time.Parse("2006-01-02", startDate); err == nil {
				start = t
			}
			if t, err := time.Parse("2006-01-02", endDate); err == nil {
				end = t
			}
		}

		startDate = start.Format("2006-01-02")
		endDate = end.Format("2006-01-02")
		periodLabel := fmt.Sprintf("%s – %s",
			start.Format("January 2, 2006"),
			end.Format("January 2, 2006"),
		)

		pl := deps.Labels.Period
		periodPresets := fycha.DefaultPeriodPresets(pl, preset)

		// Fetch sections
		var sections []ISStatementSection
		if deps.GetIncomeStatement != nil {
			ss, err := deps.GetIncomeStatement(ctx, startDate, endDate)
			if err == nil {
				sections = ss
			}
		}
		if sections == nil {
			sections = mockISSections()
		}

		// Calculate KPIs from sections
		totalRevenue, totalExpenses, netIncome := calcISKPIs(sections)

		netIncomeVariant := "success"
		if netIncome < 0 {
			netIncomeVariant = "danger"
		}

		pageData := &IncomeStatementPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.IncomeStatement.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "report",
				ActiveSubNav:   "income-statement",
				HeaderTitle:    deps.Labels.IncomeStatement.Title,
				HeaderSubtitle: deps.Labels.IncomeStatement.Subtitle,
				HeaderIcon:     "icon-trending-up",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:  "income-statement-content",
			ActivePreset:     preset,
			StartDate:        startDate,
			EndDate:          endDate,
			PeriodLabel:      periodLabel,
			PeriodPresets:    periodPresets,
			TotalRevenue:     formatCurrencyFS(totalRevenue),
			TotalExpenses:    formatCurrencyFS(totalExpenses),
			NetIncome:        formatCurrencyFS(netIncome),
			NetIncomeVariant: netIncomeVariant,
			NetIncomeTrend:   "+12%",
			Sections:         sections,
		}

		if viewCtx.IsHTMX {
			return view.OK("income-statement-content", pageData)
		}
		return view.OK("income-statement", pageData)
	})
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func calcISKPIs(sections []ISStatementSection) (totalRevenue, totalExpenses, netIncome float64) {
	// Sections are: Revenue, Cost of Sales, (Gross Profit calc), Operating Expenses,
	// (Operating Income calc), Other Expenses, (Net Income calc).
	// We identify by title prefix.
	for _, s := range sections {
		switch s.Title {
		case "REVENUE":
			totalRevenue = parseISAmount(s.Subtotal)
		case "OPERATING EXPENSES", "OTHER EXPENSES", "COST OF SALES":
			totalExpenses += parseISAmount(s.Subtotal)
		}
	}
	netIncome = totalRevenue - totalExpenses
	return
}

func parseISAmount(s string) float64 {
	// Extract digits, decimal point, and leading minus sign from formatted currency strings.
	var result float64
	clean := ""
	for _, ch := range s {
		if ch == '-' || (ch >= '0' && ch <= '9') || ch == '.' {
			clean += string(ch)
		}
	}
	fmt.Sscanf(clean, "%f", &result)
	return result
}

// ---------------------------------------------------------------------------
// Mock data (Phase 8)
// ---------------------------------------------------------------------------

// mockISSections returns a realistic income statement for a Philippine salon/spa.
// Numbers are based on the plan doc examples.
func mockISSections() []ISStatementSection {
	return []ISStatementSection{
		{
			Title: "REVENUE",
			Lines: []ISStatementLine{
				{Code: "4010", Name: "Hair Services Revenue", CurrentPeriod: "₱380,000.00", PriorPeriod: "₱345,000.00", Change: "+10.1%"},
				{Code: "4020", Name: "Nail Services Revenue", CurrentPeriod: "₱52,000.00", PriorPeriod: "₱48,000.00", Change: "+8.3%"},
				{Code: "4030", Name: "Spa & Body Services Revenue", CurrentPeriod: "₱18,000.00", PriorPeriod: "₱10,800.00", Change: "+66.7%"},
			},
			Subtotal: "₱450,000.00",
			Bold:     true,
		},
		{
			Title: "COST OF SALES",
			Lines: []ISStatementLine{
				{Code: "5010", Name: "Salon Supplies Used", CurrentPeriod: "₱28,000.00", PriorPeriod: "₱25,500.00", Change: "+9.8%"},
				{Code: "5020", Name: "Cost of Products Sold", CurrentPeriod: "₱12,000.00", PriorPeriod: "₱11,200.00", Change: "+7.1%"},
			},
			Subtotal: "₱40,000.00",
			Bold:     false,
		},
		{
			Title:    "GROSS PROFIT",
			Lines:    nil,
			Subtotal: "₱410,000.00",
			Bold:     true,
		},
		{
			Title: "OPERATING EXPENSES",
			Groups: []ISStatementGroup{
				{
					Title: "Selling Expenses",
					Lines: []ISStatementLine{
						{Code: "6410", Name: "Marketing & Advertising", CurrentPeriod: "₱18,000.00", PriorPeriod: "₱15,000.00", Change: "+20.0%"},
						{Code: "5030", Name: "Service Commission Expense", CurrentPeriod: "₱38,000.00", PriorPeriod: "₱34,500.00", Change: "+10.1%"},
					},
					Subtotal: "₱56,000.00",
				},
				{
					Title: "General & Administrative",
					Lines: []ISStatementLine{
						{Code: "6110", Name: "Salaries & Wages", CurrentPeriod: "₱180,000.00", PriorPeriod: "₱175,000.00", Change: "+2.9%"},
						{Code: "6120", Name: "SSS / PhilHealth / Pag-IBIG", CurrentPeriod: "₱12,400.00", PriorPeriod: "₱12,400.00", Change: "0.0%"},
						{Code: "6210", Name: "Rent Expense", CurrentPeriod: "₱45,000.00", PriorPeriod: "₱45,000.00", Change: "0.0%"},
						{Code: "6220", Name: "Utilities (Electric/Water)", CurrentPeriod: "₱12,500.00", PriorPeriod: "₱11,800.00", Change: "+5.9%"},
						{Code: "6510", Name: "Depreciation Expense", CurrentPeriod: "₱16,000.00", PriorPeriod: "₱16,000.00", Change: "0.0%"},
						{Code: "6620", Name: "Insurance", CurrentPeriod: "₱8,000.00", PriorPeriod: "₱8,000.00", Change: "0.0%"},
						{Code: "6630", Name: "Office Supplies & Miscellaneous", CurrentPeriod: "₱5,100.00", PriorPeriod: "₱3,200.00", Change: "+59.4%"},
					},
					Subtotal: "₱279,000.00",
				},
			},
			Subtotal: "₱335,000.00",
			Bold:     false,
		},
		{
			Title:    "OPERATING INCOME",
			Lines:    nil,
			Subtotal: "₱75,000.00",
			Bold:     true,
		},
		{
			Title: "OTHER EXPENSES",
			Lines: []ISStatementLine{
				{Code: "7010", Name: "Bank Charges & Fees", CurrentPeriod: "₱5,500.00", PriorPeriod: "₱4,800.00", Change: "+14.6%"},
				{Code: "7020", Name: "Interest Expense", CurrentPeriod: "₱2,000.00", PriorPeriod: "₱2,000.00", Change: "0.0%"},
			},
			Subtotal: "₱7,500.00",
			Bold:     false,
		},
		{
			Title:    "NET INCOME",
			Lines:    nil,
			Subtotal: "₱67,500.00",
			Bold:     true,
		},
	}
}

// ---------------------------------------------------------------------------
// Currency formatter (package-level, shared by all financial statement views)
// ---------------------------------------------------------------------------

func formatCurrencyFS(amount float64) string {
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
