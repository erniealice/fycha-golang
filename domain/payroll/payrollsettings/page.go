package payrollsettings

import (
	"context"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
)

// ---------------------------------------------------------------------------
// Gov Rates View
// ---------------------------------------------------------------------------

// GovRatesDeps holds view dependencies for the government contribution rates page.
type GovRatesDeps struct {
	Routes       Routes
	Labels       Labels
	CommonLabels pyeza.CommonLabels
}

// GovRatesPageData holds the data for the gov rates settings page.
type GovRatesPageData struct {
	types.PageData
	ContentTemplate string
	Rates           []GovContributionRate
}

// GovContributionRate represents a single government agency contribution rate row.
type GovContributionRate struct {
	Agency        string // SSS, PhilHealth, Pag-IBIG, BIR
	RateType      string // "percentage" or "fixed"
	EmployeeRate  string
	EmployerRate  string
	EffectiveDate string
	Notes         string
}

// NewGovRatesView creates the government contribution rates settings view.
func NewGovRatesView(deps *GovRatesDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &GovRatesPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.GovRates.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "gov-contribution-rates",
				HeaderTitle:    deps.Labels.GovRates.Page.Heading,
				HeaderSubtitle: deps.Labels.GovRates.Page.Caption,
				HeaderIcon:     "icon-shield",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "gov-rates-content",
			Rates:           govContributionRates(deps.Labels),
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "payroll-setting"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("gov-rates", pageData)
	})
}

// govContributionRates returns the Philippine mandatory contribution rate table.
// These are the 2024–2026 contribution rates as published by each agency.
func govContributionRates(l Labels) []GovContributionRate {
	return []GovContributionRate{
		{
			Agency:        l.GovRates.Agency.SSS,
			RateType:      "percentage",
			EmployeeRate:  "4.5%",
			EmployerRate:  "9.5%",
			EffectiveDate: "Jan 1, 2025",
			Notes:         "Max MSC ₱30,000 (effective Jan 2025)",
		},
		{
			Agency:        l.GovRates.Agency.PhilHealth,
			RateType:      "percentage",
			EmployeeRate:  "2.5%",
			EmployerRate:  "2.5%",
			EffectiveDate: "Jan 1, 2024",
			Notes:         "Based on basic monthly salary; split equally",
		},
		{
			Agency:        l.GovRates.Agency.PagIBIG,
			RateType:      "fixed + percentage",
			EmployeeRate:  "₱100 – ₱200",
			EmployerRate:  "₱100",
			EffectiveDate: "Jan 1, 2023",
			Notes:         "₱100 for salary ≤₱1,500; ₱200 for >₱1,500",
		},
		{
			Agency:        l.GovRates.Agency.BIRWithholding,
			RateType:      "graduated",
			EmployeeRate:  "0% – 35%",
			EmployerRate:  "N/A",
			EffectiveDate: "Jan 1, 2023",
			Notes:         "TRAIN Law graduated table; employer withholds on behalf",
		},
	}
}

// ---------------------------------------------------------------------------
// Pay Periods View
// ---------------------------------------------------------------------------

// PayPeriodsDeps holds view dependencies for the pay periods settings page.
type PayPeriodsDeps struct {
	Routes       Routes
	Labels       Labels
	CommonLabels pyeza.CommonLabels
}

// PayPeriodsPageData holds the data for the pay periods settings page.
type PayPeriodsPageData struct {
	types.PageData
	ContentTemplate string
	Schedules       []PayPeriodSchedule
}

// PayPeriodSchedule represents a pay period schedule configuration.
type PayPeriodSchedule struct {
	Name       string
	Frequency  string // "semi-monthly", "monthly", "weekly"
	CutOffDay1 string
	CutOffDay2 string
	PayDay1    string
	PayDay2    string
	Active     bool
}

// NewPayPeriodsView creates the pay periods settings view.
func NewPayPeriodsView(deps *PayPeriodsDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		pageData := &PayPeriodsPageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.PayPeriods.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "pay-periods",
				HeaderTitle:    deps.Labels.PayPeriods.Page.Heading,
				HeaderSubtitle: deps.Labels.PayPeriods.Page.Caption,
				HeaderIcon:     "icon-calendar",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "pay-periods-content",
			Schedules:       mockPayPeriods(),
		}

		return view.OK("pay-periods", pageData)
	})
}

func mockPayPeriods() []PayPeriodSchedule {
	return []PayPeriodSchedule{
		{
			Name:       "Semi-Monthly (1st–15th)",
			Frequency:  "semi-monthly",
			CutOffDay1: "15th of the month",
			CutOffDay2: "Last day of the month",
			PayDay1:    "20th of the month",
			PayDay2:    "5th of the following month",
			Active:     true,
		},
	}
}
