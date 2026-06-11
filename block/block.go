// Package block exposes fycha.Block() — the Lego composition entry point
// for the fycha accounting domain (reports, asset, ledger, loans, equity,
// payroll, financial statements, cash, and expenses/prepayments). Consumer
// apps import this package and optionally alias it:
//
//	import fychablock "github.com/erniealice/fycha-golang/block"
//	// ...
//	fychablock.Block()               // all modules
//	fychablock.Block(
//	    fychablock.WithReports(),
//	    fychablock.WithLedger(),
//	)                                 // selective modules
//
// This package lives in a sub-package (not the fycha root) to avoid a Go
// import cycle: fycha/views/* imports fycha (root) for route/label types,
// so Block() cannot live in the root package while also importing fycha/views/*.
//
// Companion files in this package:
//   - options.go       — BlockOption / blockConfig / WithX / wantX accessors (incl. WithUseCases)
//   - usecases.go      — typed UseCases struct (wiring contract for fycha.Block)
//   - asset.go         — Asset domain wiring (wireAssetModule + proto<->record translators)
//   - callbacks.go     — lapsing-schedule + revaluation callback helpers called by asset.go
//   - helpers.go       — workspace/currency helpers (getDefaultWorkspaceID, getFunctionalCurrency)
//   - dashboard_wiring.go — closure-based dashboard wiring helpers (wireLedgerDashboard, etc.)
package block

import (
	"context"
	"fmt"
	"log"
	"net/http"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"

	topref "github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"

	asset "github.com/erniealice/fycha-golang/domain/asset"
	expenditure "github.com/erniealice/fycha-golang/domain/expenditure"
	finance "github.com/erniealice/fycha-golang/domain/finance"
	ledger "github.com/erniealice/fycha-golang/domain/ledger"
	payroll "github.com/erniealice/fycha-golang/domain/payroll"
	tax "github.com/erniealice/fycha-golang/domain/tax"
	treasury "github.com/erniealice/fycha-golang/domain/treasury"
	report "github.com/erniealice/fycha-golang/service/report"
	reportmod "github.com/erniealice/fycha-golang/service/report/views"
	cashbookview "github.com/erniealice/fycha-golang/service/report/views/cash_book"
	financialmod "github.com/erniealice/fycha-golang/service/report/views/financial"
)

// ---------------------------------------------------------------------------
// routeRegistrarFull — optional extension for raw http.HandlerFunc routes
// ---------------------------------------------------------------------------

// routeRegistrarFull extends pyeza.RouteRegistrar with HandleFunc support.
// Consumer apps whose RouteRegistrar implements this interface can register raw
// http.HandlerFunc routes (e.g. cash book export). Apps that do not implement
// HandleFunc will skip those routes with a log warning.
type routeRegistrarFull interface {
	pyeza.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// handleFunc is a nil-safe helper that registers an http.HandlerFunc route if the
// RouteRegistrar supports it, otherwise logs a warning and skips.
func handleFunc(r pyeza.RouteRegistrar, method, path string, handler http.HandlerFunc) {
	if handler == nil {
		return
	}
	if full, ok := r.(routeRegistrarFull); ok {
		full.HandleFunc(method, path, handler)
		return
	}
	log.Printf("fycha.Block: RouteRegistrar does not support HandleFunc — skipping %s %s", method, path)
}

// ---------------------------------------------------------------------------
// Block — the main Lego entry point
// ---------------------------------------------------------------------------

// Block registers fycha accounting domain modules (reports, asset, ledger,
// loans, equity, payroll, financial statements, cash deposits + petty cash,
// and expenses/prepayments). Call with no options to register ALL modules.
// Call with specific WithX() options for a subset.
func Block(opts ...BlockOption) pyeza.AppOption {
	cfg := &blockConfig{}
	for _, opt := range opts {
		opt(cfg)
	}
	// "Enable all modules" is derived — true when no module-toggling option was
	// passed. Non-module options (WithUseCases, WithAssetDepreciationRunURL)
	// must NOT flip this off, otherwise `Block(WithUseCases(...), WithAssetDepreciationRunURL(...))`
	// silently registers zero modules. Matches the centymo block pattern.
	moduleSelected := cfg.reports || cfg.asset || cfg.ledger || cfg.loans ||
		cfg.equity || cfg.payroll || cfg.cash || cfg.expenses || cfg.financial ||
		cfg.taxRate || cfg.forexRate || cfg.withholdingCertificate
	cfg.enableAll = !moduleSelected

	return func(ctx *pyeza.AppContext) error {
		// --- Type-assert translations ---
		translations, ok := ctx.Translations.(*lynguaV1.TranslationProvider)
		if !ok || translations == nil {
			return fmt.Errorf("fycha.Block: ctx.Translations must be *lynguaV1.TranslationProvider")
		}

		// --- Typed use cases supplied via WithUseCases() ---
		if cfg.useCases == nil {
			return fmt.Errorf("fycha.Block: WithUseCases(...) is required")
		}
		useCases := cfg.useCases

		// 20260521 Wave B P1.E.1-P1.E.5 — fycha report views consume
		// service-driven typed closures via `useCases.Reports.<Group>`
		// instead of `ctx.LedgerReportingSvc`. The legacy
		// `fycha.DataSource` duck interface no longer ships; the
		// LedgerReportingSvc context field is retained on pyeza
		// AppContext but unread by fycha.

		// --- Type-assert attachment operations ---
		uploadFile, _ := ctx.UploadFile.(func(context.Context, string, string, []byte, string) error)
		listAttachments, _ := ctx.ListAttachments.(func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error))
		createAttachment, _ := ctx.CreateAttachment.(func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		deleteAttachment, _ := ctx.DeleteAttachment.(func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		newAttachmentID, _ := ctx.NewAttachmentID.(func() string)

		// --- Type-assert reference checker (optional — nil-safe) ---
		// Used to wire in-use checks for deletable entities (e.g. asset H5 gate).
		var refChecker topref.Checker
		if ctx.RefChecker != nil {
			refChecker, _ = ctx.RefChecker.(topref.Checker)
		}

		// --- Fycha-specific table labels ---
		fychaTableLabels := pyeza.MapTableLabels(ctx.Common)

		// --- Load routes (defaults + optional lyngua overrides) ---
		reportsRoutes := report.DefaultReportsRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "reports", &reportsRoutes)

		assetRoutes := asset.DefaultAssetRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "asset", &assetRoutes)

		accountRoutes := ledger.DefaultAccountRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_account", &accountRoutes)

		journalRoutes := ledger.DefaultJournalRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_journal", &journalRoutes)

		statementRoutes := ledger.DefaultLedgerStatementRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_statement", &statementRoutes)

		fiscalPeriodRoutes := ledger.DefaultFiscalPeriodRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "fiscal_period", &fiscalPeriodRoutes)

		ledgerSettingsRoutes := ledger.DefaultLedgerSettingsRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_settings", &ledgerSettingsRoutes)

		loanRoutes := treasury.DefaultLoanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "loan", &loanRoutes)

		loanPaymentRoutes := treasury.DefaultLoanPaymentRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "loan_payment", &loanPaymentRoutes)

		equityRoutes := ledger.DefaultEquityRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "equity", &equityRoutes)

		lapsingScheduleRoutes := asset.DefaultLapsingScheduleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "lapsing_schedule", &lapsingScheduleRoutes)

		depreciationRunRoutes := asset.DefaultDepreciationRunRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "depreciation_run", &depreciationRunRoutes)

		assetCategoryDepreciationRoutes := asset.DefaultAssetCategoryDepreciationRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "asset_category_depreciation", &assetCategoryDepreciationRoutes)

		taxRateRoutes := tax.DefaultTaxRateRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "tax_rate", &taxRateRoutes)

		forexRateRoutes := finance.DefaultForexRateRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "forex_rate", &forexRateRoutes)

		withholdingCertificateRoutes := treasury.DefaultWithholdingCertificateRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "withholding_certificate", &withholdingCertificateRoutes)

		// --- Load labels ---
		var reportsLabels report.ReportsLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "report.json", "", &reportsLabels); err != nil {
			log.Printf("fycha.Block: warning loading reports labels: %v", err)
		}

		assetLabels := asset.DefaultAssetLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "asset.json", "", &assetLabels)

		accountLabels := ledger.DefaultAccountLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "account.json", "", &accountLabels)

		journalLabels := ledger.DefaultJournalLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "journal.json", "", &journalLabels)

		fiscalPeriodLabels := ledger.DefaultFiscalPeriodLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "fiscal_period.json", "", &fiscalPeriodLabels)

		recurringTemplateLabels := ledger.DefaultRecurringTemplateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "recurring_template.json", "", &recurringTemplateLabels)

		loanLabels := treasury.DefaultLoanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "loan.json", "", &loanLabels)

		depreciationRunLabels := asset.DefaultDepreciationRunLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "depreciation_run.json", "", &depreciationRunLabels)

		assetRevaluationLabels := asset.DefaultAssetRevaluationLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "asset_revaluation.json", "", &assetRevaluationLabels)

		depreciationPoliciesLabels := asset.DefaultDepreciationPoliciesLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "depreciation_policies.json", "", &depreciationPoliciesLabels)

		loanPaymentLabels := treasury.DefaultLoanPaymentLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "loan_payment.json", "", &loanPaymentLabels)

		taxRateLabels := tax.DefaultTaxRateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "tax_rate.json", "", &taxRateLabels)

		forexRateLabels := finance.DefaultForexRateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "forex_rate.json", "", &forexRateLabels)

		withholdingCertificateLabels := treasury.DefaultWithholdingCertificateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "withholding_certificate.json", "", &withholdingCertificateLabels)

		// =====================================================================
		// Reports module (fycha)
		// =====================================================================

		if cfg.wantReports() {
			reportmod.NewModule(&reportmod.ModuleDeps{
				Routes:       reportsRoutes,
				Labels:       reportsLabels,
				CommonLabels: ctx.Common,
				TableLabels:  ctx.Table,
				// 20260520-21 Wave B P1.E.1-P1.E.5 — every report view now
				// consumes the service-driven use case via a typed closure
				// instead of `fycha.DataSource`. The 13 in-scope methods
				// are routed through their respective `Reports.<Group>`
				// sub-aggregate on the typed block UseCases.
				GetReceivablesAgingReport:  useCases.Reports.ARAging.GetReceivablesAgingReport,
				GetCollectionSummaryReport: useCases.Reports.ARAging.GetCollectionSummaryReport,
				GetPayablesAgingReport:     useCases.Reports.APAging.GetPayablesAgingReport,
				GetGrossProfitReport:       useCases.Reports.GrossCashFlow.GetGrossProfitReport,
				GetRevenueReport:           useCases.Reports.DomainSpecific.GetRevenueReport,
				GetExpenditureReport:       useCases.Reports.DomainSpecific.GetExpenditureReport,
				GetDisbursementReport:      useCases.Reports.DomainSpecific.GetDisbursementReport,
				ListRevenue:                useCases.Reports.DomainSpecific.ListRevenue,
				ListExpenses:               useCases.Reports.DomainSpecific.ListExpenses,
			}).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Asset module (fycha) — see asset.go for full wiring
		// =====================================================================

		if cfg.wantAsset() {
			wireAssetModule(ctx, cfg, useCases, assetWiring{
				assetRoutes:                     assetRoutes,
				lapsingScheduleRoutes:           lapsingScheduleRoutes,
				depreciationRunRoutes:           depreciationRunRoutes,
				assetCategoryDepreciationRoutes: assetCategoryDepreciationRoutes,
				assetLabels:                     assetLabels,
				depreciationRunLabels:           depreciationRunLabels,
				depreciationPoliciesLabels:      depreciationPoliciesLabels,
				assetRevaluationLabels:          assetRevaluationLabels,
				fychaTableLabels:                fychaTableLabels,
				common:                          ctx.Common,
				refChecker:                      refChecker,
				newAttachmentID:                 newAttachmentID,
				uploadFile:                      uploadFile,
				listAttachments:                 listAttachments,
				createAttachment:                createAttachment,
				deleteAttachment:                deleteAttachment,
			})
		}

		// =====================================================================
		// Ledger module (fycha — Chart of Accounts + Journals + FiscalPeriod + Settings)
		// =====================================================================

		if cfg.wantLedger() {
			ledgerDeps := &ledger.LedgerModuleDeps{
				Routes:                  accountRoutes,
				StatementRoutes:         statementRoutes,
				JournalRoutes:           journalRoutes,
				FiscalPeriodRoutes:      fiscalPeriodRoutes,
				LedgerSettingsRoutes:    ledgerSettingsRoutes,
				CommonLabels:            ctx.Common,
				Labels:                  accountLabels,
				JournalLabels:           journalLabels,
				FiscalPeriodLabels:      fiscalPeriodLabels,
				RecurringTemplateLabels: recurringTemplateLabels,
				TableLabels:             fychaTableLabels,
				// Attachments
				NewAttachmentID:  newAttachmentID,
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
			}
			if useCases.Ledger.Account.GetListPageData != nil {
				ledgerDeps.GetAccountListPageData = useCases.Ledger.Account.GetListPageData
				ledgerDeps.CreateAccount = useCases.Ledger.Account.Create
				ledgerDeps.ReadAccount = useCases.Ledger.Account.Read
				ledgerDeps.UpdateAccount = useCases.Ledger.Account.Update
				ledgerDeps.DeleteAccount = useCases.Ledger.Account.Delete
			}
			if useCases.Ledger.JournalEntry.GetListPageData != nil {
				ledgerDeps.GetJournalEntryListPageData = useCases.Ledger.JournalEntry.GetListPageData
				ledgerDeps.CreateJournalEntry = useCases.Ledger.JournalEntry.Create
				ledgerDeps.ReadJournalEntry = useCases.Ledger.JournalEntry.Read
				ledgerDeps.UpdateJournalEntry = useCases.Ledger.JournalEntry.Update
				ledgerDeps.DeleteJournalEntry = useCases.Ledger.JournalEntry.Delete
				ledgerDeps.PostJournalEntry = useCases.Ledger.JournalEntry.Post
				ledgerDeps.ReverseJournalEntry = useCases.Ledger.JournalEntry.Reverse
			}
			if useCases.FiscalPeriod.GetListPageData != nil {
				ledgerDeps.GetFiscalPeriodListPageData = func(fctx context.Context) ([]*fiscalperiodpb.FiscalPeriod, error) {
					resp, err := useCases.FiscalPeriod.GetListPageData(fctx, &fiscalperiodpb.GetFiscalPeriodListPageDataRequest{})
					if err != nil {
						return nil, err
					}
					if resp == nil {
						return nil, nil
					}
					return resp.GetFiscalPeriodList(), nil
				}
				if useCases.FiscalPeriod.Create != nil {
					ledgerDeps.CreateFiscalPeriod = useCases.FiscalPeriod.Create
				}
				if useCases.FiscalPeriod.Close != nil {
					ledgerDeps.CloseFiscalPeriod = useCases.FiscalPeriod.Close
				}
			}
			// Wire ledger dashboard use case.
			wireLedgerDashboard(ledgerDeps, useCases)
			ledger.NewLedgerModule(ledgerDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Loans module (fycha)
		// =====================================================================

		if cfg.wantLoans() {
			loansDeps := &treasury.LoanModuleDeps{
				Routes:        loanRoutes,
				PaymentRoutes: loanPaymentRoutes,
				Labels:        loanLabels,
				PaymentLabels: loanPaymentLabels,
				CommonLabels:  ctx.Common,
				TableLabels:   fychaTableLabels,
			}
			// Wire loan dashboard use case.
			wireLoansDashboard(loansDeps, useCases)
			treasury.NewLoanModule(loansDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Equity module (fycha)
		// =====================================================================

		if cfg.wantEquity() {
			equityDeps := &ledger.EquityModuleDeps{
				Routes:       equityRoutes,
				Labels:       ledger.DefaultEquityLabels(),
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			// Wire equity dashboard use case.
			wireEquityDashboard(equityDeps, useCases)
			ledger.NewEquityModule(equityDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Payroll module (fycha)
		// =====================================================================

		if cfg.wantPayroll() {
			payrollDeps := &payroll.PayrollDashboardModuleDeps{}
			// Wire payroll dashboard use case.
			wirePayrollDashboard(payrollDeps, useCases)
			payroll.NewPayrollDashboardModule(payrollDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Financial statements module (fycha)
		// =====================================================================

		if cfg.wantFinancial() {
			financialmod.NewModule(&financialmod.ModuleDeps{
				CommonLabels: ctx.Common,
				TableLabels:  ctx.Table,
				Labels:       reportsLabels,
			}).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Cash expansion module — Deposits + Petty Cash (fycha)
		// =====================================================================

		if cfg.wantCash() {
			treasury.NewPettyCashModule(&treasury.PettyCashModuleDeps{
				// TODO: wire when useCases.Treasury.SecurityDeposit / PettyCashFund are available
			}).RegisterRoutes(ctx.Routes)

			// Cash → Reports → Cash Book — Wave B P1.E.3 service-driven closure.
			ctx.Routes.GET(report.CashBookURL, cashbookview.NewCashBookView(
				useCases.Reports.GrossCashFlow.GetCashBookReport,
				ctx.Common,
				ctx.Table,
			))
		}

		// =====================================================================
		// Expenses expansion module — Prepayments (fycha)
		// =====================================================================

		if cfg.wantExpenses() {
			expenditure.NewPrepaymentModule(&expenditure.PrepaymentModuleDeps{
				// TODO: wire when useCases.Expenditure.Prepayment is available
			}).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Tax Rate module (fycha — read-only)
		// =====================================================================

		if cfg.wantTaxRate() {
			taxRateDeps := &tax.TaxRateModuleDeps{
				Routes:       taxRateRoutes,
				Labels:       taxRateLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases.Tax.ListTaxRates != nil {
				taxRateDeps.ListTaxRates = useCases.Tax.ListTaxRates
			}
			tax.NewTaxRateModule(taxRateDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Forex Rate module (fycha — read-only)
		// =====================================================================

		if cfg.wantForexRate() {
			forexRateDeps := &finance.ForexRateModuleDeps{
				Routes:       forexRateRoutes,
				Labels:       forexRateLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases.Finance.ListForexRates != nil {
				forexRateDeps.ListForexRates = useCases.Finance.ListForexRates
			}
			finance.NewForexRateModule(forexRateDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Withholding Certificate module (fycha — full CRUD)
		// =====================================================================

		if cfg.wantWithholdingCertificate() {
			withholdingCertDeps := &treasury.WithholdingCertificateModuleDeps{
				Routes:       withholdingCertificateRoutes,
				Labels:       withholdingCertificateLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases.Treasury.ListWithholdingCertificates != nil {
				withholdingCertDeps.ListWithholdingCertificates = useCases.Treasury.ListWithholdingCertificates
			}
			treasury.NewWithholdingCertificateModule(withholdingCertDeps).RegisterRoutes(ctx.Routes)
		}

		log.Println("  fycha accounting domain initialized")
		return nil
	}
}
