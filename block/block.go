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
//   - options.go       — BlockOption / blockConfig / WithX / wantX accessors
//   - asset.go         — Asset domain wiring (wireAssetModule + proto<->record translators)
//   - callbacks.go     — lapsing-schedule + revaluation callback helpers called by asset.go
//   - helpers.go       — workspace/currency helpers (getDefaultWorkspaceID, getFunctionalCurrency)
//   - dashboard_wiring.go — reflective dashboard wiring helpers (wireLedgerDashboard, etc.)
package block

import (
	"context"
	"fmt"
	"log"
	"net/http"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"

	consumer "github.com/erniealice/espyna-golang/consumer"
	topref "github.com/erniealice/espyna-golang/reference"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"

	fycha "github.com/erniealice/fycha-golang"
	cashmod "github.com/erniealice/fycha-golang/views/cash"
	equitymod "github.com/erniealice/fycha-golang/views/equity"
	expensesmod "github.com/erniealice/fycha-golang/views/expenses"
	financialmod "github.com/erniealice/fycha-golang/views/financial"
	ledgermod "github.com/erniealice/fycha-golang/views/ledger"
	loansmod "github.com/erniealice/fycha-golang/views/loans"
	payrollmod "github.com/erniealice/fycha-golang/views/payroll"
	reportmod "github.com/erniealice/fycha-golang/views/reports"
	cashbookview "github.com/erniealice/fycha-golang/views/reports/cash_book"
	forexratemod "github.com/erniealice/fycha-golang/views/forex_rate"
	taxratemod "github.com/erniealice/fycha-golang/views/tax_rate"
	withholdingcertmod "github.com/erniealice/fycha-golang/views/withholding_certificate"
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
	cfg := &blockConfig{enableAll: len(opts) == 0}
	for _, opt := range opts {
		opt(cfg)
	}

	return func(ctx *pyeza.AppContext) error {
		// --- Type-assert translations ---
		translations, ok := ctx.Translations.(*lynguaV1.TranslationProvider)
		if !ok || translations == nil {
			return fmt.Errorf("fycha.Block: ctx.Translations must be *lynguaV1.TranslationProvider")
		}

		// --- Type-assert use cases ---
		useCases, ok := ctx.UseCases.(*consumer.UseCases)
		if !ok || useCases == nil {
			return fmt.Errorf("fycha.Block: ctx.UseCases must be *consumer.UseCases")
		}

		// --- Type-assert LedgerReportingSvc (optional — nil-safe) ---
		var ledgerReportingSvc fycha.DataSource
		if ctx.LedgerReportingSvc != nil {
			ledgerReportingSvc, _ = ctx.LedgerReportingSvc.(fycha.DataSource)
		}

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
		fychaTableLabels := fycha.MapTableLabels(ctx.Common)

		// --- Load routes (defaults + optional lyngua overrides) ---
		reportsRoutes := fycha.DefaultReportsRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "reports", &reportsRoutes)

		assetRoutes := fycha.DefaultAssetRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "asset", &assetRoutes)

		accountRoutes := fycha.DefaultAccountRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_account", &accountRoutes)

		journalRoutes := fycha.DefaultJournalRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_journal", &journalRoutes)

		statementRoutes := fycha.DefaultLedgerStatementRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_statement", &statementRoutes)

		fiscalPeriodRoutes := fycha.DefaultFiscalPeriodRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "fiscal_period", &fiscalPeriodRoutes)

		ledgerSettingsRoutes := fycha.DefaultLedgerSettingsRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "ledger_settings", &ledgerSettingsRoutes)

		loanRoutes := fycha.DefaultLoanRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "loan", &loanRoutes)

		loanPaymentRoutes := fycha.DefaultLoanPaymentRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "loan_payment", &loanPaymentRoutes)

		equityRoutes := fycha.DefaultEquityRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "equity", &equityRoutes)

		lapsingScheduleRoutes := fycha.DefaultLapsingScheduleRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "lapsing_schedule", &lapsingScheduleRoutes)

		depreciationRunRoutes := fycha.DefaultDepreciationRunRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "depreciation_run", &depreciationRunRoutes)

		assetCategoryDepreciationRoutes := fycha.DefaultAssetCategoryDepreciationRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "asset_category_depreciation", &assetCategoryDepreciationRoutes)

		taxRateRoutes := fycha.DefaultTaxRateRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "tax_rate", &taxRateRoutes)

		forexRateRoutes := fycha.DefaultForexRateRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "forex_rate", &forexRateRoutes)

		withholdingCertificateRoutes := fycha.DefaultWithholdingCertificateRoutes()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "route.json", "withholding_certificate", &withholdingCertificateRoutes)

		// --- Load labels ---
		var reportsLabels fycha.ReportsLabels
		if err := translations.LoadPath("en", ctx.BusinessType, "report.json", "", &reportsLabels); err != nil {
			log.Printf("fycha.Block: warning loading reports labels: %v", err)
		}

		assetLabels := fycha.DefaultAssetLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "asset.json", "", &assetLabels)

		accountLabels := fycha.DefaultAccountLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "account.json", "", &accountLabels)

		journalLabels := fycha.DefaultJournalLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "journal.json", "", &journalLabels)

		fiscalPeriodLabels := fycha.DefaultFiscalPeriodLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "fiscal_period.json", "", &fiscalPeriodLabels)

		recurringTemplateLabels := fycha.DefaultRecurringTemplateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "recurring_template.json", "", &recurringTemplateLabels)

		loanLabels := fycha.DefaultLoanLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "loan.json", "", &loanLabels)

		depreciationRunLabels := fycha.DefaultDepreciationRunLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "depreciation_run.json", "", &depreciationRunLabels)

		assetRevaluationLabels := fycha.DefaultAssetRevaluationLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "asset_revaluation.json", "", &assetRevaluationLabels)

		depreciationPoliciesLabels := fycha.DefaultDepreciationPoliciesLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "depreciation_policies.json", "", &depreciationPoliciesLabels)

		loanPaymentLabels := fycha.DefaultLoanPaymentLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "loan_payment.json", "", &loanPaymentLabels)

		taxRateLabels := fycha.DefaultTaxRateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "tax_rate.json", "", &taxRateLabels)

		forexRateLabels := fycha.DefaultForexRateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "forex_rate.json", "", &forexRateLabels)

		withholdingCertificateLabels := fycha.DefaultWithholdingCertificateLabels()
		_ = translations.LoadPathIfExists("en", ctx.BusinessType, "withholding_certificate.json", "", &withholdingCertificateLabels)

		// =====================================================================
		// Reports module (fycha)
		// =====================================================================

		if cfg.wantReports() {
			reportmod.NewModule(&reportmod.ModuleDeps{
				Routes:       reportsRoutes,
				DB:           ledgerReportingSvc,
				Labels:       reportsLabels,
				CommonLabels: ctx.Common,
				TableLabels:  ctx.Table,
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
			ledgerDeps := &ledgermod.ModuleDeps{
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
			if useCases != nil && useCases.Ledger != nil && useCases.Ledger.Account != nil {
				ledgerDeps.GetAccountListPageData = useCases.Ledger.Account.GetAccountListPageData.Execute
				ledgerDeps.CreateAccount = useCases.Ledger.Account.CreateAccount.Execute
				ledgerDeps.ReadAccount = useCases.Ledger.Account.ReadAccount.Execute
				ledgerDeps.UpdateAccount = useCases.Ledger.Account.UpdateAccount.Execute
				ledgerDeps.DeleteAccount = useCases.Ledger.Account.DeleteAccount.Execute
			}
			if useCases != nil && useCases.Ledger != nil && useCases.Ledger.JournalEntry != nil {
				uje := useCases.Ledger.JournalEntry
				ledgerDeps.GetJournalEntryListPageData = uje.GetJournalEntryListPageData.Execute
				ledgerDeps.CreateJournalEntry = uje.CreateJournalEntry.Execute
				ledgerDeps.ReadJournalEntry = uje.ReadJournalEntry.Execute
				ledgerDeps.UpdateJournalEntry = uje.UpdateJournalEntry.Execute
				ledgerDeps.DeleteJournalEntry = uje.DeleteJournalEntry.Execute
				ledgerDeps.PostJournalEntry = uje.PostJournalEntry.Execute
				ledgerDeps.ReverseJournalEntry = uje.ReverseJournalEntry.Execute
			}
			if useCases != nil && useCases.Ledger != nil && useCases.Ledger.FiscalPeriod != nil {
				ufp := useCases.Ledger.FiscalPeriod
				ledgerDeps.GetFiscalPeriodListPageData = func(fctx context.Context) ([]*fiscalperiodpb.FiscalPeriod, error) {
					resp, err := ufp.GetFiscalPeriodListPageData.Execute(fctx, &fiscalperiodpb.GetFiscalPeriodListPageDataRequest{})
					if err != nil {
						return nil, err
					}
					if resp == nil {
						return nil, nil
					}
					return resp.GetFiscalPeriodList(), nil
				}
				if ufp.CreateFiscalPeriod != nil {
					ledgerDeps.CreateFiscalPeriod = ufp.CreateFiscalPeriod.Execute
				}
				if ufp.CloseFiscalPeriod != nil {
					ledgerDeps.CloseFiscalPeriod = ufp.CloseFiscalPeriod.Execute
				}
			}
			// Wire ledger dashboard use case.
			wireLedgerDashboard(ledgerDeps, useCases)
			ledgermod.NewModule(ledgerDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Loans module (fycha)
		// =====================================================================

		if cfg.wantLoans() {
			loansDeps := &loansmod.ModuleDeps{
				Routes:        loanRoutes,
				PaymentRoutes: loanPaymentRoutes,
				Labels:        loanLabels,
				PaymentLabels: loanPaymentLabels,
				CommonLabels:  ctx.Common,
				TableLabels:   fychaTableLabels,
			}
			// Wire loan dashboard use case.
			wireLoansDashboard(loansDeps, useCases)
			loansmod.NewModule(loansDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Equity module (fycha)
		// =====================================================================

		if cfg.wantEquity() {
			equityDeps := &equitymod.ModuleDeps{
				Routes:       equityRoutes,
				Labels:       fycha.DefaultEquityLabels(),
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			// Wire equity dashboard use case.
			wireEquityDashboard(equityDeps, useCases)
			equitymod.NewModule(equityDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Payroll module (fycha)
		// =====================================================================

		if cfg.wantPayroll() {
			payrollDeps := &payrollmod.ModuleDeps{}
			// Wire payroll dashboard use case.
			wirePayrollDashboard(payrollDeps, useCases)
			payrollmod.NewModule(payrollDeps).RegisterRoutes(ctx.Routes)
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
			cashmod.NewModule(&cashmod.ModuleDeps{
				// TODO: wire when useCases.Treasury.SecurityDeposit / PettyCashFund are available
			}).RegisterRoutes(ctx.Routes)

			// Cash → Reports → Cash Book
			ctx.Routes.GET(fycha.CashBookURL, cashbookview.NewCashBookView(ledgerReportingSvc, ctx.Common, ctx.Table))
		}

		// =====================================================================
		// Expenses expansion module — Prepayments (fycha)
		// =====================================================================

		if cfg.wantExpenses() {
			expensesmod.NewModule(&expensesmod.ModuleDeps{
				// TODO: wire when useCases.Expenditure.Prepayment is available
			}).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Tax Rate module (fycha — read-only)
		// =====================================================================

		if cfg.wantTaxRate() {
			taxRateDeps := &taxratemod.ModuleDeps{
				Routes:       taxRateRoutes,
				Labels:       taxRateLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases.Tax != nil && useCases.Tax.TaxRate != nil {
				taxRateDeps.ListTaxRates = useCases.Tax.TaxRate.ListTaxRates.Execute
			}
			taxratemod.NewModule(taxRateDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Forex Rate module (fycha — read-only)
		// =====================================================================

		if cfg.wantForexRate() {
			forexRateDeps := &forexratemod.ModuleDeps{
				Routes:       forexRateRoutes,
				Labels:       forexRateLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases.Finance != nil && useCases.Finance.ForexRate != nil {
				forexRateDeps.ListForexRates = useCases.Finance.ForexRate.ListForexRates.Execute
			}
			forexratemod.NewModule(forexRateDeps).RegisterRoutes(ctx.Routes)
		}

		// =====================================================================
		// Withholding Certificate module (fycha — full CRUD)
		// =====================================================================

		if cfg.wantWithholdingCertificate() {
			withholdingCertDeps := &withholdingcertmod.ModuleDeps{
				Routes:       withholdingCertificateRoutes,
				Labels:       withholdingCertificateLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases.Treasury != nil && useCases.Treasury.WithholdingCertificate != nil {
				withholdingCertDeps.ListWithholdingCertificates = useCases.Treasury.WithholdingCertificate.ListWithholdingCertificates.Execute
			}
			withholdingcertmod.NewModule(withholdingCertDeps).RegisterRoutes(ctx.Routes)
		}

		log.Println("  fycha accounting domain initialized")
		return nil
	}
}
