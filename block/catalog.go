package block

// catalog.go — composition-v2 unit binders for the fycha accounting domain.
//
// Each XxxUnit function wraps a domain descriptor (Describe()) with a Mount
// closure that wires the typed UseCases + Infra deps into the view module.
// AllUnits() returns the complete ordered unit list for the fycha block.
//
// The design mirrors the fayna catalog pattern:
//   - Each binder calls <entity>.Describe() for the Unit skeleton.
//   - Mount builds the ModuleDeps from the post-overlay Routes/Labels pointers.
//   - wireAssetModule is reused directly (it was already split out of Block).
//   - The reports service and the multi-entity ledger module are handled by
//     their own binders that produce a single Unit per logical surface; the
//     actual module wiring is copied faithfully from block.go.
//
// Entities with embed.go + routes.go + labels.go each get a Unit binder here.
// Entities without embed.go (ledger/account, journal, fiscal_period,
// recurring_template) are wired inline in LedgerUnit via LedgerModuleDeps —
// they have no descriptor, but their routes/labels ride on the ledger Unit.

import (
	"context"

	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"

	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/espyna-golang/consumer/compose"

	assetpkg "github.com/erniealice/fycha-golang/domain/asset"
	assetentity "github.com/erniealice/fycha-golang/domain/asset/asset"
	assetcategory "github.com/erniealice/fycha-golang/domain/asset/asset_category"
	depreciationrun "github.com/erniealice/fycha-golang/domain/asset/depreciation_run"
	lapsingschedule "github.com/erniealice/fycha-golang/domain/asset/lapsing_schedule"
	expenditure "github.com/erniealice/fycha-golang/domain/expenditure"
	prepayment "github.com/erniealice/fycha-golang/domain/expenditure/prepayment"
	finance "github.com/erniealice/fycha-golang/domain/finance"
	forexrate "github.com/erniealice/fycha-golang/domain/finance/forex_rate"
	fundingdom "github.com/erniealice/fycha-golang/domain/funding"
	fundingpkg "github.com/erniealice/fycha-golang/domain/funding/funding"
	fundinglabels "github.com/erniealice/fycha-golang/domain/funding/funding/labels"
	ledger "github.com/erniealice/fycha-golang/domain/ledger"
	equity "github.com/erniealice/fycha-golang/domain/ledger/equity"
	ledgerview "github.com/erniealice/fycha-golang/domain/ledger/ledger"
	payroll "github.com/erniealice/fycha-golang/domain/payroll"
	payrolldashboard "github.com/erniealice/fycha-golang/domain/payroll/payrolldashboard"
	payrollemployee "github.com/erniealice/fycha-golang/domain/payroll/payrollemployee"
	payrollsettings "github.com/erniealice/fycha-golang/domain/payroll/payrollsettings"
	remittance "github.com/erniealice/fycha-golang/domain/payroll/remittance"
	run "github.com/erniealice/fycha-golang/domain/payroll/run"
	tax "github.com/erniealice/fycha-golang/domain/tax"
	taxrate "github.com/erniealice/fycha-golang/domain/tax/tax_rate"
	treasury "github.com/erniealice/fycha-golang/domain/treasury"
	loan "github.com/erniealice/fycha-golang/domain/treasury/loan"
	pettycash "github.com/erniealice/fycha-golang/domain/treasury/petty_cash"
	withholdingcert "github.com/erniealice/fycha-golang/domain/treasury/withholding_certificate"
	report "github.com/erniealice/fycha-golang/service/report"
	reportmod "github.com/erniealice/fycha-golang/service/report/views"
)

// ---------------------------------------------------------------------------
// Asset domain units
// ---------------------------------------------------------------------------

// AssetUnit wires the asset entity module. The asset module is complex
// (it spans Surfaces A-F via wireAssetModule) so the Mount closure
// delegates to the existing wireAssetModule helper.
func AssetUnit(uc *UseCases, infra *Infra) compose.Unit {
	u := assetentity.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*assetentity.Routes)

		var depRunRoutes depreciationrun.Routes
		if drUnit, ok := compose.RoutesOf[*depreciationrun.Routes](mc, "asset.depreciation_run"); ok {
			depRunRoutes = *drUnit
		}
		var lapsingRoutes lapsingschedule.Routes
		if lsUnit, ok := compose.RoutesOf[*lapsingschedule.Routes](mc, "asset.lapsing_schedule"); ok {
			lapsingRoutes = *lsUnit
		}
		var assetCategoryRoutes assetcategory.Routes
		if acUnit, ok := compose.RoutesOf[*assetcategory.Routes](mc, "asset.asset_category"); ok {
			assetCategoryRoutes = *acUnit
		}

		var depRunLabels depreciationrun.Labels
		if drLabels, ok := compose.LabelsOf[*depreciationrun.Labels](mc, "asset.depreciation_run"); ok {
			depRunLabels = *drLabels
		}

		var assetDepRunURL string
		if infra != nil {
			assetDepRunURL = infra.AssetDepreciationRunURL
		}

		w := assetWiring{
			assetRoutes:                     *r,
			lapsingScheduleRoutes:           lapsingRoutes,
			depreciationRunRoutes:           depRunRoutes,
			assetCategoryDepreciationRoutes: assetCategoryRoutes,
			assetLabels:                     *u.Labels.(*assetentity.Labels),
			depreciationRunLabels:           depRunLabels,
			depreciationPoliciesLabels:      assetpkg.DefaultDepreciationPoliciesLabels(),
			assetRevaluationLabels:          assetpkg.DefaultAssetRevaluationLabels(),
			fychaTableLabels:                mc.Table,
			common:                          mc.Common,
		}
		if infra != nil {
			w.newAttachmentID = infra.NewAttachmentID
			w.uploadFile = infra.UploadFile
			w.listAttachments = infra.ListAttachments
			w.createAttachment = infra.CreateAttachment
			w.deleteAttachment = infra.DeleteAttachment
			if infra.RefChecker != nil {
				w.refChecker = infra.RefChecker
			}
		}

		// Wire the asset module directly without going through blockConfig/Block();
		// use a minimal blockConfig with asset enabled so wireAssetModule gates work.
		// wireAssetModule only reads ctx.Routes, ctx.Table, and ctx.Common, so a
		// minimal consumerapp.AppContext with those fields populated is sufficient.
		cfg := &blockConfig{asset: true, assetDepreciationRunURL: assetDepRunURL}
		minCtx := &consumerapp.AppContext{Routes: mc.Routes, Table: mc.Table, Common: mc.Common}
		wireAssetModule(minCtx, cfg, uc, w)
		return nil
	}
	return u
}

// DepreciationRunUnit wires the depreciation-run history list + detail module.
// This unit contributes routes/labels/templates; the actual Surface D registration
// happens inside AssetUnit via wireAssetModule. The Unit is still listed so its
// routes/labels are available to sibling units via compose.RoutesOf/LabelsOf.
func DepreciationRunUnit(_ *UseCases, _ *Infra) compose.Unit {
	return depreciationrun.Describe()
}

// LapsingScheduleUnit contributes the lapsing-schedule routes/templates to the
// composition. The route handler registration happens inside AssetUnit via
// wireAssetModule (Surface B).
func LapsingScheduleUnit(_ *UseCases, _ *Infra) compose.Unit {
	return lapsingschedule.Describe()
}

// AssetCategoryUnit contributes asset-category depreciation routes/templates.
// Handler registration (Surfaces C and F) happens inside AssetUnit.
func AssetCategoryUnit(_ *UseCases, _ *Infra) compose.Unit {
	return assetcategory.Describe()
}

// ---------------------------------------------------------------------------
// Ledger domain unit (multi-entity)
// ---------------------------------------------------------------------------

// LedgerUnit wires the Chart of Accounts + Journal Entry + FiscalPeriod +
// LedgerSettings module. The ledger domain spans multiple sub-entities that
// share a single LedgerModuleDeps — they have no individual descriptors.
func LedgerUnit(uc *UseCases, infra *Infra) compose.Unit {
	// The ledger module has no single entity descriptor: its routes and labels
	// are declared across domain/ledger/account, journal, fiscal_period, and
	// recurring_template, none of which have embed.go. We use a synthetic unit
	// key that matches the pattern and attach the register-only Mount closure.
	u := compose.Unit{
		Key:       "ledger.ledger",
		Templates: ledgerview.TemplatesFS,
	}
	u.Mount = func(mc *compose.MountContext) error {
		accountRoutes := ledger.DefaultAccountRoutes()
		journalRoutes := ledger.DefaultJournalRoutes()
		statementRoutes := ledger.DefaultLedgerStatementRoutes()
		fiscalPeriodRoutes := ledger.DefaultFiscalPeriodRoutes()
		ledgerSettingsRoutes := ledger.DefaultLedgerSettingsRoutes()

		accountLabels := ledger.DefaultAccountLabels()
		journalLabels := ledger.DefaultJournalLabels()
		fiscalPeriodLabels := ledger.DefaultFiscalPeriodLabels()
		recurringTemplateLabels := ledger.DefaultRecurringTemplateLabels()

		fychaTableLabels := mc.Table

		deps := &ledger.LedgerModuleDeps{
			Routes:                  accountRoutes,
			StatementRoutes:         statementRoutes,
			JournalRoutes:           journalRoutes,
			FiscalPeriodRoutes:      fiscalPeriodRoutes,
			LedgerSettingsRoutes:    ledgerSettingsRoutes,
			CommonLabels:            mc.Common,
			Labels:                  accountLabels,
			JournalLabels:           journalLabels,
			FiscalPeriodLabels:      fiscalPeriodLabels,
			RecurringTemplateLabels: recurringTemplateLabels,
			TableLabels:             fychaTableLabels,
		}
		if infra != nil {
			deps.NewAttachmentID = infra.NewAttachmentID
			deps.UploadFile = infra.UploadFile
			deps.ListAttachments = infra.ListAttachments
			deps.CreateAttachment = infra.CreateAttachment
			deps.DeleteAttachment = infra.DeleteAttachment
		}

		if uc != nil {
			if uc.Ledger.Account.GetListPageData != nil {
				deps.GetAccountListPageData = uc.Ledger.Account.GetListPageData
				deps.CreateAccount = uc.Ledger.Account.Create
				deps.ReadAccount = uc.Ledger.Account.Read
				deps.UpdateAccount = uc.Ledger.Account.Update
				deps.DeleteAccount = uc.Ledger.Account.Delete
			}
			if uc.Ledger.JournalEntry.GetListPageData != nil {
				deps.GetJournalEntryListPageData = uc.Ledger.JournalEntry.GetListPageData
				deps.CreateJournalEntry = uc.Ledger.JournalEntry.Create
				deps.ReadJournalEntry = uc.Ledger.JournalEntry.Read
				deps.UpdateJournalEntry = uc.Ledger.JournalEntry.Update
				deps.DeleteJournalEntry = uc.Ledger.JournalEntry.Delete
				deps.PostJournalEntry = uc.Ledger.JournalEntry.Post
				deps.ReverseJournalEntry = uc.Ledger.JournalEntry.Reverse
			}
			if uc.FiscalPeriod.GetListPageData != nil {
				deps.GetFiscalPeriodListPageData = func(fctx context.Context) ([]*fiscalperiodpb.FiscalPeriod, error) {
					resp, err := uc.FiscalPeriod.GetListPageData(fctx, &fiscalperiodpb.GetFiscalPeriodListPageDataRequest{})
					if err != nil {
						return nil, err
					}
					if resp == nil {
						return nil, nil
					}
					return resp.GetFiscalPeriodList(), nil
				}
				if uc.FiscalPeriod.Create != nil {
					deps.CreateFiscalPeriod = uc.FiscalPeriod.Create
				}
				if uc.FiscalPeriod.Close != nil {
					deps.CloseFiscalPeriod = uc.FiscalPeriod.Close
				}
			}
			wireLedgerDashboard(deps, uc)
		}

		ledger.NewLedgerModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Equity unit
// ---------------------------------------------------------------------------

// EquityUnit wires the equity account + transaction module.
func EquityUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := equity.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*equity.Routes)
		l := u.Labels.(*equity.Labels)

		deps := &ledger.EquityModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		if uc != nil {
			wireEquityDashboard(deps, uc)
		}
		ledger.NewEquityModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Treasury units
// ---------------------------------------------------------------------------

// LoanUnit wires the loan + loan-payment module.
func LoanUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := loan.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*loan.Routes)
		l := u.Labels.(*loan.Labels)

		loanPaymentRoutes := treasury.DefaultLoanPaymentRoutes()
		loanPaymentLabels := treasury.DefaultLoanPaymentLabels()

		deps := &treasury.LoanModuleDeps{
			Routes:        *r,
			PaymentRoutes: loanPaymentRoutes,
			Labels:        *l,
			PaymentLabels: loanPaymentLabels,
			CommonLabels:  mc.Common,
			TableLabels:   mc.Table,
		}
		if uc != nil {
			wireLoansDashboard(deps, uc)
		}
		treasury.NewLoanModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// WithholdingCertificateUnit wires the withholding-certificate CRUD module.
func WithholdingCertificateUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := withholdingcert.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*withholdingcert.Routes)
		l := u.Labels.(*withholdingcert.Labels)

		deps := &treasury.WithholdingCertificateModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		if uc != nil && uc.Treasury.ListWithholdingCertificates != nil {
			deps.ListWithholdingCertificates = uc.Treasury.ListWithholdingCertificates
		}
		treasury.NewWithholdingCertificateModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// PettyCashUnit wires the petty-cash module (stub deps — use cases TBD).
func PettyCashUnit(_ *UseCases, _ *Infra) compose.Unit {
	u := pettycash.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		treasury.NewPettyCashModule(&treasury.PettyCashModuleDeps{}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Tax unit
// ---------------------------------------------------------------------------

// TaxRateUnit wires the tax-rate read-only list module.
func TaxRateUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := taxrate.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*taxrate.Routes)
		l := u.Labels.(*taxrate.Labels)

		deps := &tax.TaxRateModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		if uc != nil && uc.Tax.ListTaxRates != nil {
			deps.ListTaxRates = uc.Tax.ListTaxRates
		}
		tax.NewTaxRateModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Finance unit
// ---------------------------------------------------------------------------

// ForexRateUnit wires the forex-rate read-only list module.
func ForexRateUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := forexrate.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		r := u.Routes.(*forexrate.Routes)
		l := u.Labels.(*forexrate.Labels)

		deps := &finance.ForexRateModuleDeps{
			Routes:       *r,
			Labels:       *l,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
		}
		if uc != nil && uc.Finance.ListForexRates != nil {
			deps.ListForexRates = uc.Finance.ListForexRates
		}
		finance.NewForexRateModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Expenditure unit
// ---------------------------------------------------------------------------

// PrepaymentUnit wires the prepayment module (stub deps — use cases TBD).
func PrepaymentUnit(_ *UseCases, _ *Infra) compose.Unit {
	u := prepayment.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		expenditure.NewPrepaymentModule(&expenditure.PrepaymentModuleDeps{}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Payroll units
// ---------------------------------------------------------------------------

// PayrollUnit wires the payroll dashboard module. The payroll dashboard
// aggregates routes from sibling payroll entities (run, remittance,
// payrollsettings) so those are fetched via compose.RoutesOf.
func PayrollUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := payrolldashboard.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		l := u.Labels.(*payrolldashboard.Labels)

		runRoutes := run.DefaultRoutes()
		if rr, ok := compose.RoutesOf[*run.Routes](mc, "payroll.run"); ok {
			runRoutes = *rr
		}
		remittanceRoutes := remittance.DefaultRoutes()
		if rr, ok := compose.RoutesOf[*remittance.Routes](mc, "payroll.remittance"); ok {
			remittanceRoutes = *rr
		}
		settingsRoutes := payrollsettings.DefaultRoutes()
		if sr, ok := compose.RoutesOf[*payrollsettings.Routes](mc, "payroll.payrollsettings"); ok {
			settingsRoutes = *sr
		}

		deps := &payroll.PayrollDashboardModuleDeps{
			Routes:           payroll.PayrollRunRoutes(runRoutes),
			RemittanceRoutes: payroll.PayrollRemittanceRoutes(remittanceRoutes),
			SettingsRoutes:   payroll.PayrollSettingsRoutes(settingsRoutes),
			Labels:           payroll.PayrollLabels{Dashboard: payroll.PayrollDashboardLabels(*l)},
			CommonLabels:     mc.Common,
		}
		if uc != nil {
			wirePayrollDashboard(deps, uc)
		}
		payroll.NewPayrollDashboardModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// PayrollRunUnit contributes payroll-run routes/labels/templates to the
// composition. Handler registration happens inside PayrollUnit.
func PayrollRunUnit(_ *UseCases, _ *Infra) compose.Unit {
	return run.Describe()
}

// PayrollRemittanceUnit contributes payroll-remittance routes/labels/templates.
// Handler registration happens inside PayrollUnit.
func PayrollRemittanceUnit(_ *UseCases, _ *Infra) compose.Unit {
	return remittance.Describe()
}

// PayrollEmployeeUnit contributes payroll-employee routes/labels/templates.
func PayrollEmployeeUnit(_ *UseCases, _ *Infra) compose.Unit {
	u := payrollemployee.Describe()
	// TODO: wire PayrollEmployeeModule when PayrollEmployeeModuleDeps is available.
	return u
}

// PayrollSettingsUnit contributes payroll-settings routes/labels/templates.
func PayrollSettingsUnit(_ *UseCases, _ *Infra) compose.Unit {
	u := payrollsettings.Describe()
	// TODO: wire PayrollSettingsModule when PayrollSettingsModuleDeps is available.
	return u
}

// ---------------------------------------------------------------------------
// Funding unit
// ---------------------------------------------------------------------------

// FundingUnit wires the funding view module (8 views: source list/detail,
// card list/detail, 4 drawer form stubs). The funding descriptor contributes
// the AppEntry + sidebar items; the Mount closure builds the FundingModuleDeps
// from the block's UseCases and registers the HTTP routes.
//
// Route URLs use the defaults from FundingModule.RegisterRoutes (no /app/
// prefix — workspace_path middleware dispatches /w/{slug}/funding/* to the
// bare /funding/* handlers). Labels are overlaid via LabelJSON in phase 1.
func FundingUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := fundingpkg.Describe()
	u.Mount = func(mc *compose.MountContext) error {
		// Type-assert the post-overlay labels pointer back to concrete type.
		var lbls fundinglabels.FundingFormLabels
		if u.Labels != nil {
			if p, ok := u.Labels.(*fundinglabels.FundingFormLabels); ok {
				lbls = *p
			}
		}
		if lbls.Source.Title == "" {
			lbls = fundinglabels.DefaultFundingFormLabels()
		}

		deps := &fundingdom.FundingModuleDeps{
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,
			Labels:       fundingdom.FundingFormLabels(lbls),
			// Route URLs match DefaultFundRoutes() (workspace_path strips /app/).
			SourceListURL:     "/funding/sources",
			SourceDetailURL:   "/funding/sources/{fund_id}",
			CardListURL:       "/funding/cards",
			CardDetailURL:     "/funding/cards/{allocation_id}",
			AllocationFormURL: "/funding/allocation/form",
			DrawFormURL:       "/funding/draw/form",
			SettlementFormURL: "/funding/settlement/form",
			TransferFormURL:   "/funding/transfer/form",
		}

		if uc != nil {
			f := &uc.Funding
			deps.ReadFund = f.ReadFund
			deps.ListFunds = f.ListFunds
			deps.ReadAllocation = f.ReadAllocation
			deps.ListAllocations = f.ListAllocations
			deps.ListTransactions = f.ListTransactions
		}

		fundingdom.NewFundingModule(deps).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// Reports unit
// ---------------------------------------------------------------------------

// ReportsUnit wires the fycha service/report module. The report service is a
// cross-cutting view surface — it has no entity descriptor of its own (routes
// and labels are in service/report, not in a domain entity package). A
// synthetic unit is used.
func ReportsUnit(uc *UseCases, _ *Infra) compose.Unit {
	u := compose.Unit{
		Key:       "service.reports",
		Templates: reportmod.TemplatesFS,
	}
	u.Mount = func(mc *compose.MountContext) error {
		if uc == nil {
			return nil
		}

		reportsRoutes := report.DefaultReportsRoutes()
		var reportsLabels report.ReportsLabels

		reportmod.NewModule(&reportmod.ModuleDeps{
			Routes:       reportsRoutes,
			Labels:       reportsLabels,
			CommonLabels: mc.Common,
			TableLabels:  mc.Table,

			GetReceivablesAgingReport:  uc.Reports.ARAging.GetReceivablesAgingReport,
			GetCollectionSummaryReport: uc.Reports.ARAging.GetCollectionSummaryReport,
			GetPayablesAgingReport:     uc.Reports.APAging.GetPayablesAgingReport,
			GetGrossProfitReport:       uc.Reports.GrossCashFlow.GetGrossProfitReport,
			GetRevenueReport:           uc.Reports.DomainSpecific.GetRevenueReport,
			GetExpenditureReport:       uc.Reports.DomainSpecific.GetExpenditureReport,
			GetDisbursementReport:      uc.Reports.DomainSpecific.GetDisbursementReport,
			ListRevenue:                uc.Reports.DomainSpecific.ListRevenue,
			ListExpenses:               uc.Reports.DomainSpecific.ListExpenses,
		}).RegisterRoutes(mc.Routes)
		return nil
	}
	return u
}

// ---------------------------------------------------------------------------
// AllUnits — ordered unit list for the fycha block
// ---------------------------------------------------------------------------

// AllUnits returns the complete curated unit list for the fycha accounting
// domain, in the same logical grouping order as Block(). Service-admin's
// composition root calls this to obtain the unit slice for the engine.
func AllUnits(uc *UseCases, infra *Infra) []compose.Unit {
	return []compose.Unit{
		// Reports service surface
		ReportsUnit(uc, infra),

		// Asset domain — order matters: DepreciationRun + LapsingSchedule +
		// AssetCategory must precede AssetUnit so their routes/labels are
		// available via compose.RoutesOf/LabelsOf inside AssetUnit.Mount.
		DepreciationRunUnit(uc, infra),
		LapsingScheduleUnit(uc, infra),
		AssetCategoryUnit(uc, infra),
		AssetUnit(uc, infra),

		// Ledger domain (multi-entity module)
		LedgerUnit(uc, infra),
		EquityUnit(uc, infra),

		// Treasury domain
		LoanUnit(uc, infra),
		WithholdingCertificateUnit(uc, infra),
		PettyCashUnit(uc, infra),

		// Tax + Finance read-only modules
		TaxRateUnit(uc, infra),
		ForexRateUnit(uc, infra),

		// Expenditure (stub)
		PrepaymentUnit(uc, infra),

		// Funding domain (cross-workspace fund sources + cards + drawer forms)
		FundingUnit(uc, infra),

		// Payroll domain — siblings precede dashboard so routes are available
		PayrollRunUnit(uc, infra),
		PayrollRemittanceUnit(uc, infra),
		PayrollSettingsUnit(uc, infra),
		PayrollEmployeeUnit(uc, infra),
		PayrollUnit(uc, infra),
	}
}
