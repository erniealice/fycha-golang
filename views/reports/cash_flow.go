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

// CFLine is one line item in the cash flow statement.
type CFLine struct {
	Code        string // contra-account code (may be empty for calculated lines)
	Name        string // description
	Amount      string // formatted amount (positive = inflow, negative = outflow)
	IsNegative  bool   // true when the amount represents an outflow
	IsSubtotal  bool   // section net total (bold)
	IsLabel     bool   // section sub-header (no amount, just a label row)
	IndentLevel int    // 0 = section header, 1 = item, 2 = sub-item
}

// CFActivity is one of the three cash flow activity sections.
type CFActivity struct {
	Title    string   // "OPERATING ACTIVITIES"
	Lines    []CFLine // individual line items
	NetTotal string   // "Net Cash from Operating Activities"
	NetLabel string   // Label for the subtotal row
	IsPositive bool    // true if net total >= 0 (for styling)
}

// CFVerification holds the beginning/ending cash reconciliation.
type CFVerification struct {
	BeginningBalance string
	NetChange        string
	EndingBalance    string
	IsVerified       bool
	// Per-account reconciliation
	CashAccounts []CFLine
}

// ---------------------------------------------------------------------------
// Deps + PageData
// ---------------------------------------------------------------------------

// CashFlowDeps holds dependencies for the Cash Flow Statement view.
type CashFlowDeps struct {
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	Labels       fycha.ReportsLabels

	// GetCashFlow fetches cash flow data for the given period.
	// Phase 8: set to nil — mock data is used automatically.
	GetCashFlow func(ctx context.Context, startDate, endDate string) ([]CFActivity, *CFVerification, error)
}

// CashFlowPageData is the template data for the cash-flow page.
type CashFlowPageData struct {
	types.PageData
	ContentTemplate string

	// Period filter state
	ActivePreset  string
	StartDate     string
	EndDate       string
	PeriodLabel   string
	PeriodPresets []fycha.FilterOption

	// KPI summary metrics
	OperatingCF    string
	NetChange      string
	EndingCash     string
	OperatingTrend string // "+15%"

	// Statement body
	Activities   []CFActivity
	Verification *CFVerification
}

// ---------------------------------------------------------------------------
// View constructor
// ---------------------------------------------------------------------------

// NewCashFlowView creates the Cash Flow Statement view.
func NewCashFlowView(deps *CashFlowDeps) view.View {
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

		// Fetch data
		var activities []CFActivity
		var verification *CFVerification
		if deps.GetCashFlow != nil {
			acts, vfy, err := deps.GetCashFlow(ctx, startDate, endDate)
			if err == nil {
				activities = acts
				verification = vfy
			}
		}
		if activities == nil {
			activities, verification = mockCFData()
		}

		// Extract KPIs
		operatingCF := ""
		netChange := ""
		if verification != nil {
			netChange = verification.NetChange
		}
		endingCash := ""
		if verification != nil {
			endingCash = verification.EndingBalance
		}
		for _, act := range activities {
			if act.Title == "OPERATING ACTIVITIES" {
				operatingCF = act.NetTotal
				break
			}
		}

		pageData := &CashFlowPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.CashFlow.Title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      "reports",
				ActiveSubNav:   "cash-flow",
				HeaderTitle:    deps.Labels.CashFlow.Title,
				HeaderSubtitle: deps.Labels.CashFlow.Subtitle,
				HeaderIcon:     "icon-activity",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:  "cash-flow-content",
			ActivePreset:     preset,
			StartDate:        startDate,
			EndDate:          endDate,
			PeriodLabel:      periodLabel,
			PeriodPresets:    periodPresets,
			OperatingCF:      operatingCF,
			NetChange:        netChange,
			EndingCash:       endingCash,
			OperatingTrend:   "+15%",
			Activities:       activities,
			Verification:     verification,
		}

		if viewCtx.IsHTMX {
			return view.OK("cash-flow-content", pageData)
		}
		return view.OK("cash-flow", pageData)
	})
}

// ---------------------------------------------------------------------------
// Mock data (Phase 8)
// ---------------------------------------------------------------------------

// mockCFData returns realistic cash flow statement data (direct method).
// Based on plan doc Section 5 example output.
func mockCFData() ([]CFActivity, *CFVerification) {
	activities := []CFActivity{
		{
			Title: "OPERATING ACTIVITIES",
			Lines: []CFLine{
				{Name: "Collections from customers (AR)", Amount: "₱500,000.00", IndentLevel: 1},
				{Name: "Subscription/package payments received", Amount: "₱48,000.00", IndentLevel: 1},
				{Name: "Payments to suppliers (AP)", Amount: "(₱120,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Salaries & wages paid", Amount: "(₱180,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Government contributions (SSS/PhilHealth/HDMF)", Amount: "(₱25,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Rent paid", Amount: "(₱45,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Utilities paid", Amount: "(₱15,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Marketing & advertising", Amount: "(₱18,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Other operating expenses", Amount: "(₱20,500.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Income tax payments", Amount: "(₱12,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Bank charges", Amount: "(₱3,000.00)", IsNegative: true, IndentLevel: 1},
			},
			NetTotal:   "₱109,500.00",
			NetLabel:   "Net Cash from Operating Activities",
			IsPositive: true,
		},
		{
			Title: "INVESTING ACTIVITIES",
			Lines: []CFLine{
				{Name: "Purchase of salon equipment", Amount: "(₱25,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Purchase of office furniture", Amount: "(₱12,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Proceeds from sale of old equipment", Amount: "₱10,000.00", IndentLevel: 1},
			},
			NetTotal:   "(₱27,000.00)",
			NetLabel:   "Net Cash used in Investing Activities",
			IsPositive: false,
		},
		{
			Title: "FINANCING ACTIVITIES",
			Lines: []CFLine{
				{Name: "Bank loan received", Amount: "₱0.00", IndentLevel: 1},
				{Name: "Loan principal repayment", Amount: "(₱10,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Interest paid on loan", Amount: "(₱2,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Owner withdrawals", Amount: "(₱20,000.00)", IsNegative: true, IndentLevel: 1},
				{Name: "Owner contributions", Amount: "₱0.00", IndentLevel: 1},
			},
			NetTotal:   "(₱32,000.00)",
			NetLabel:   "Net Cash used in Financing Activities",
			IsPositive: false,
		},
	}

	verification := &CFVerification{
		BeginningBalance: "₱269,000.00",
		NetChange:        "₱50,500.00",
		EndingBalance:    "₱319,500.00",
		IsVerified:       true,
		CashAccounts: []CFLine{
			{Code: "1010", Name: "Cash in Bank", Amount: "₱295,500.00"},
			{Code: "1020", Name: "Petty Cash", Amount: "₱9,000.00"},
			{Code: "1030", Name: "Cash on Hand (GCash/Maya)", Amount: "₱15,000.00"},
		},
	}

	return activities, verification
}
