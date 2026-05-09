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
package block

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"

	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"

	consumer "github.com/erniealice/espyna-golang/consumer"

	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	deprunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation_run"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"

	fycha "github.com/erniealice/fycha-golang"
	assetmod "github.com/erniealice/fycha-golang/views/asset"
	assetaction "github.com/erniealice/fycha-golang/views/asset/action"
	assetform "github.com/erniealice/fycha-golang/views/asset/form"
	assetlist "github.com/erniealice/fycha-golang/views/asset/list"
	cashmod "github.com/erniealice/fycha-golang/views/cash"
	equitymod "github.com/erniealice/fycha-golang/views/equity"
	expensesmod "github.com/erniealice/fycha-golang/views/expenses"
	financialmod "github.com/erniealice/fycha-golang/views/financial"
	ledgermod "github.com/erniealice/fycha-golang/views/ledger"
	loansmod "github.com/erniealice/fycha-golang/views/loans"
	payrollmod "github.com/erniealice/fycha-golang/views/payroll"
	reportmod "github.com/erniealice/fycha-golang/views/reports"
	cashbookview "github.com/erniealice/fycha-golang/views/reports/cash_book"
	assetcataction "github.com/erniealice/fycha-golang/views/asset_category/action"
	assetcatpolicies "github.com/erniealice/fycha-golang/views/asset_category/policies"
	depreciationrunmod "github.com/erniealice/fycha-golang/views/depreciation_run"
	lapsinglist "github.com/erniealice/fycha-golang/views/lapsing_schedule/list"
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
// BlockOption — per-module granular selection
// ---------------------------------------------------------------------------

// BlockOption enables specific fycha sub-modules within Block().
type BlockOption func(*blockConfig)

type blockConfig struct {
	enableAll bool
	reports   bool
	asset     bool
	ledger    bool
	loans     bool
	equity    bool
	payroll   bool
	cash      bool
	expenses  bool
	financial bool
	// assetDepreciationRunURL is the resolved run-detail URL template plumbed into
	// the Surface A drawer so toast links point to the correct run-detail page.
	// Set via WithAssetDepreciationRunURL (Wave 2 hard requirement).
	assetDepreciationRunURL string
}

func WithReports() BlockOption   { return func(c *blockConfig) { c.reports = true } }
func WithAsset() BlockOption     { return func(c *blockConfig) { c.asset = true } }
func WithLedger() BlockOption    { return func(c *blockConfig) { c.ledger = true } }
func WithLoans() BlockOption     { return func(c *blockConfig) { c.loans = true } }
func WithEquity() BlockOption    { return func(c *blockConfig) { c.equity = true } }
func WithPayroll() BlockOption   { return func(c *blockConfig) { c.payroll = true } }
func WithCash() BlockOption      { return func(c *blockConfig) { c.cash = true } }
func WithExpenses() BlockOption  { return func(c *blockConfig) { c.expenses = true } }
func WithFinancial() BlockOption { return func(c *blockConfig) { c.financial = true } }

// WithAssetDepreciationRunURL injects the run-detail URL into the block so
// the Surface A drawer can include a resolved link in its success toast payload.
// Wave 2 hard requirement — must be called before routes register.
//
// Example:
//
//	fychablock.Block(
//	    fychablock.WithAsset(),
//	    fychablock.WithAssetDepreciationRunURL(fycha.DepreciationRunDetailURL),
//	)
func WithAssetDepreciationRunURL(url string) BlockOption {
	return func(c *blockConfig) { c.assetDepreciationRunURL = url }
}

func (c *blockConfig) wantReports() bool   { return c.enableAll || c.reports }
func (c *blockConfig) wantAsset() bool     { return c.enableAll || c.asset }
func (c *blockConfig) wantLedger() bool    { return c.enableAll || c.ledger }
func (c *blockConfig) wantLoans() bool     { return c.enableAll || c.loans }
func (c *blockConfig) wantEquity() bool    { return c.enableAll || c.equity }
func (c *blockConfig) wantPayroll() bool   { return c.enableAll || c.payroll }
func (c *blockConfig) wantCash() bool      { return c.enableAll || c.cash }
func (c *blockConfig) wantExpenses() bool  { return c.enableAll || c.expenses }
func (c *blockConfig) wantFinancial() bool { return c.enableAll || c.financial }

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
		// Asset module (fycha)
		// =====================================================================

		if cfg.wantAsset() {
			assetDeps := &assetmod.ModuleDeps{
				Routes:       assetRoutes,
				CommonLabels: ctx.Common,
				Labels:       assetLabels,
				TableLabels:  ctx.Table,
				// Depreciation Run + Revaluation labels (Surface A / E)
				DepreciationRunLabels:   depreciationRunLabels,
				AssetRevaluationLabels:  assetRevaluationLabels,
				// Depreciation run routes + toast URL (Surface A)
				DepreciationRunRoutes:   depreciationRunRoutes,
				AssetDepreciationRunURL: cfg.assetDepreciationRunURL,
				// Attachments
				UploadFile:       uploadFile,
				ListAttachments:  listAttachments,
				CreateAttachment: createAttachment,
				DeleteAttachment: deleteAttachment,
			}

			// Typed asset stack (asset-stack buildout, 2026-05-03). Falls back to
			// nothing-wired if the asset use cases are unavailable (mock build, etc.) —
			// same graceful-degradation semantics as ledger/treasury.
			if useCases != nil && useCases.Asset != nil && useCases.Asset.Asset != nil {
				ua := useCases.Asset.Asset
				assetDeps.NewID = func() string {
					if newAttachmentID != nil {
						return newAttachmentID()
					}
					return "" // CreateAsset use case generates IDs internally via IDService
				}
				assetDeps.CreateAsset = func(fctx context.Context, a *assetform.Record) error {
					_, err := ua.CreateAsset.Execute(fctx, &assetpb.CreateAssetRequest{Data: recordToAsset(a)})
					return err
				}
				assetDeps.ReadAsset = func(fctx context.Context, id string) (*assetform.Record, error) {
					resp, err := ua.ReadAsset.Execute(fctx, &assetpb.ReadAssetRequest{Data: &assetpb.Asset{Id: id}})
					if err != nil {
						return nil, err
					}
					if resp == nil || len(resp.Data) == 0 {
						return nil, fmt.Errorf("asset %s not found", id)
					}
					return assetToRecord(resp.Data[0]), nil
				}
				assetDeps.UpdateAsset = func(fctx context.Context, a *assetform.Record) error {
					_, err := ua.UpdateAsset.Execute(fctx, &assetpb.UpdateAssetRequest{Data: recordToAsset(a)})
					return err
				}
				// DeleteAsset preserves the legacy soft-delete (active=false) semantic via
				// SetAssetActive — routes the change through audit/auth instead of bypass.
				assetDeps.DeleteAsset = func(fctx context.Context, id string) error {
					_, err := ua.SetAssetActive.Execute(fctx, &assetpb.SetAssetActiveRequest{AssetId: id, Active: false})
					return err
				}
				assetDeps.SetActive = func(fctx context.Context, id string, active bool) error {
					_, err := ua.SetAssetActive.Execute(fctx, &assetpb.SetAssetActiveRequest{AssetId: id, Active: active})
					return err
				}
				assetDeps.ListAssets = func(fctx context.Context, status string) ([]assetlist.AssetRow, error) {
					active := status == "active"
					resp, err := ua.GetAssetListPageData.Execute(fctx, &assetpb.GetAssetListPageDataRequest{
						Filters: &commonpb.FilterRequest{
							Filters: []*commonpb.TypedFilter{
								{
									Field: "active",
									FilterType: &commonpb.TypedFilter_BooleanFilter{
										BooleanFilter: &commonpb.BooleanFilter{Value: active},
									},
								},
							},
						},
					})
					if err != nil {
						return nil, err
					}
					rows := make([]assetlist.AssetRow, 0, len(resp.GetAssetList()))
					for _, a := range resp.GetAssetList() {
						rows = append(rows, assetToRow(a))
					}
					return rows, nil
				}
			}

			// Wire Surface A (depreciation-run) use cases.
			if useCases != nil && useCases.Asset != nil && useCases.Asset.DepreciationRun != nil {
				udr := useCases.Asset.DepreciationRun
				if udr.ListDepreciationCandidates != nil {
					ucsCopy := useCases // capture for closure
					assetDeps.ListDepreciationCandidates = func(fctx context.Context, assetID, asOfDate string) ([]assetaction.DepreciationCandidate, error) {
						return listDepreciationCandidatesForAsset(fctx, ucsCopy, assetID, asOfDate)
					}
				}
				if udr.GenerateDepreciationRun != nil {
					ucsCopy := useCases // capture for closure
					assetDeps.GenerateDepreciationRun = func(fctx context.Context, req assetaction.DepreciationRunRequest) (*assetaction.DepreciationRunResult, error) {
						return generateDepreciationRunForAsset(fctx, ucsCopy, req)
					}
				}
				// DepreciationFieldsLockedFn: fields are locked once any run has been posted.
				// A dedicated use case (GetAssetDepreciationLock) is pending Wave 3 espyna work.
				// For now we expose the hook so the edit form can call it — nil means unlocked.
			}

			// Wire Surface E (revaluation) use cases.
			// RevalueAsset and PreviewRevaluation are wired when the asset revaluation
			// use case is available (Wave 3 espyna wiring — nil-safe fallback for now).
			// assetDeps.RevalueAsset / assetDeps.PreviewRevaluation remain nil until wired.

			assetmod.NewModule(assetDeps).RegisterRoutes(ctx.Routes)

			// ---------------------------------------------------------------------------
			// Surface B — Lapsing Schedule live list page (replaces mock at /app/assets/reports/lapsing-schedule)
			// ---------------------------------------------------------------------------
			lapsingDeps := &lapsinglist.ViewDeps{
				Routes:                lapsingScheduleRoutes,
				AssetRoutes:           assetRoutes,
				DepreciationRunRoutes: depreciationRunRoutes,
				Labels:                depreciationRunLabels,
				CommonLabels:          ctx.Common,
				TableLabels:           fychaTableLabels,
			}
			// Wire ListCandidates when depreciation use cases are available.
			if useCases != nil && useCases.Asset != nil && useCases.Asset.DepreciationRun != nil &&
				useCases.Asset.DepreciationRun.ListDepreciationCandidates != nil {
				lapsingDeps.ListCandidates = func(fctx context.Context, asOfDate, cursor string, limit int32) ([]lapsinglist.CandidateRow, string, error) {
					return listCandidatesWorkspace(fctx, useCases, asOfDate, cursor, limit)
				}
			}
			// Register Surface B at new URL.
			ctx.Routes.GET(fycha.LapsingScheduleListURL, lapsinglist.NewView(lapsingDeps))
			// Redirect from legacy mock URL to new URL (preserves bookmarks).
			handleFunc(ctx.Routes, "GET", fycha.AssetLapsingScheduleURL, func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, fycha.LapsingScheduleListURL, http.StatusMovedPermanently)
			})

			// ---------------------------------------------------------------------------
			// Surface F — Depreciation Policies actionable page (replaces mock at /app/assets/settings/depreciation-policies)
			// ---------------------------------------------------------------------------
			policiesDeps := &assetcatpolicies.ViewDeps{
				Routes:       assetCategoryDepreciationRoutes,
				Labels:       depreciationPoliciesLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			// Wire ListPolicies when asset-category use cases are available.
			if useCases != nil && useCases.Asset != nil && useCases.Asset.AssetCategory != nil {
				policiesDeps.ListPolicies = func(fctx context.Context) ([]assetcatpolicies.PolicyRow, error) {
					return listPoliciesWithRollup(fctx, useCases)
				}
			}
			ctx.Routes.GET(fycha.DepreciationPoliciesURL, assetcatpolicies.NewView(policiesDeps))

			// ---------------------------------------------------------------------------
			// Surface F preview drawer (read-only /action/asset-policy/depreciation-preview/{category_id})
			// ---------------------------------------------------------------------------
			previewDeps := &assetcataction.DepreciationPreviewDeps{
				Routes:       assetCategoryDepreciationRoutes,
				Labels:       depreciationRunLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases != nil && useCases.Asset != nil && useCases.Asset.DepreciationRun != nil &&
				useCases.Asset.DepreciationRun.ListDepreciationCandidates != nil {
				previewDeps.ListPolicyCandidates = func(fctx context.Context, categoryID, asOfDate string) ([]assetcataction.PreviewCandidateRow, error) {
					return listCandidatesForPolicy(fctx, useCases, categoryID, asOfDate)
				}
			}
			ctx.Routes.GET(fycha.AssetPolicyDepreciationPreviewURL, assetcataction.NewDepreciationPreviewView(previewDeps))

			// ---------------------------------------------------------------------------
			// Surface C — per-category / per-policy depreciation-run drawer
			// Both URLs use the same handler; scope kind is inferred from the URL path.
			// ---------------------------------------------------------------------------
			categoryRunDeps := &assetcataction.CategoryDepreciationRunDeps{
				Routes:       assetCategoryDepreciationRoutes,
				RunRoutes:    depreciationRunRoutes,
				Labels:       depreciationRunLabels,
				CommonLabels: ctx.Common,
			}
			if useCases != nil && useCases.Asset != nil && useCases.Asset.DepreciationRun != nil {
				udr := useCases.Asset.DepreciationRun
				if udr.ListDepreciationCandidates != nil {
					ucsCopy := useCases // capture for closure
					categoryRunDeps.ListCategoryCandidates = func(fctx context.Context, categoryID, scopeKind, asOfDate string) ([]assetcataction.CategoryDepreciationRunAssetRow, error) {
						return listCandidatesForCategory(fctx, ucsCopy, categoryID, scopeKind, asOfDate)
					}
				}
				if udr.GenerateDepreciationRun != nil {
					ucsCopy := useCases // capture for closure
					categoryRunDeps.GenerateCategoryRun = func(fctx context.Context, req assetcataction.CategoryDepreciationRunRequest) (*assetcataction.CategoryDepreciationRunResult, error) {
						return generateDepreciationRunForCategory(fctx, ucsCopy, req)
					}
				}
			}
			categoryRunView := assetcataction.NewCategoryDepreciationRunAction(categoryRunDeps)
			ctx.Routes.GET(fycha.AssetCategoryDepreciationRunURL, categoryRunView)
			ctx.Routes.POST(fycha.AssetCategoryDepreciationRunURL, categoryRunView)
			ctx.Routes.GET(fycha.AssetPolicyDepreciationRunURL, categoryRunView)
			ctx.Routes.POST(fycha.AssetPolicyDepreciationRunURL, categoryRunView)

			// ---------------------------------------------------------------------------
			// Surface D — Depreciation Runs history list + detail module
			// ---------------------------------------------------------------------------
			drDeps := &depreciationrunmod.ModuleDeps{
				Routes:       depreciationRunRoutes,
				Labels:       depreciationRunLabels,
				CommonLabels: ctx.Common,
				TableLabels:  fychaTableLabels,
			}
			if useCases != nil && useCases.Asset != nil && useCases.Asset.DepreciationRun != nil {
				udr := useCases.Asset.DepreciationRun
				if udr.ListDepreciationRuns != nil {
					ucsCopy := useCases
					drDeps.ListDepreciationRuns = func(fctx context.Context, scope depreciationrunmod.ListDepreciationRunsScope) ([]depreciationrunmod.DepreciationRunRow, string, error) {
						return listDepreciationRunsForWorkspace(fctx, ucsCopy, scope)
					}
				}
				if udr.ReadDepreciationRun != nil {
					ucsCopy := useCases
					drDeps.ReadDepreciationRun = func(fctx context.Context, id string) (*depreciationrunmod.DepreciationRunWithEntries, error) {
						return readDepreciationRunWithEntries(fctx, ucsCopy, id)
					}
				}
			}
			depreciationrunmod.NewModule(drDeps).RegisterRoutes(ctx.Routes)
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

		log.Println("  fycha accounting domain initialized")
		return nil
	}
}

// ---------------------------------------------------------------------------
// Asset type-translation helpers (asset-stack buildout, 2026-05-03)
// ---------------------------------------------------------------------------

// recordToAsset converts a view-layer assetform.Record to the proto Asset type.
// Money fields are translated from float64 pesos → int64 centavos.
// Enum fields are translated from string → proto enum using the generated _value maps.
// Unknown enum strings map to 0 (*_UNSPECIFIED), preserving current behaviour.
func recordToAsset(r *assetform.Record) *assetpb.Asset {
	a := &assetpb.Asset{
		Id:                 r.ID,
		AssetNumber:        r.AssetNumber,
		Name:               r.Name,
		AssetType:          assetpb.AssetType(assetpb.AssetType_value[r.AssetType]),
		AssetCategoryId:    r.AssetCategoryID,
		AcquisitionCost:    int64(math.Round(r.AcquisitionCost * 100)),
		Currency:           r.Currency,
		SalvageValue:       int64(math.Round(r.SalvageValue * 100)),
		BookValue:          int64(math.Round(r.BookValue * 100)),
		UsefulLifeMonths:   int32(r.UsefulLifeMonths),
		DepreciationMethod: assetpb.DepreciationMethod(assetpb.DepreciationMethod_value[r.DepreciationMethod]),
		Status:             assetpb.AssetStatus(assetpb.AssetStatus_value[r.Status]),
		Active:             r.Active,
	}
	if r.Description != "" {
		a.Description = &r.Description
	}
	if r.LocationID != "" {
		a.LocationId = &r.LocationID
	}
	return a
}

// assetToRecord converts a proto Asset back to the view-layer assetform.Record.
// Money fields are translated from int64 centavos → float64 pesos.
// Enum strings are lowercased and stripped of their proto prefix so they round-trip
// to the form values the view layer expects (e.g. "DEPRECIATION_METHOD_STRAIGHT_LINE"
// → "straight_line").
func assetToRecord(a *assetpb.Asset) *assetform.Record {
	r := &assetform.Record{
		ID:                 a.GetId(),
		AssetNumber:        a.GetAssetNumber(),
		Name:               a.GetName(),
		AssetType:          strings.ToLower(strings.TrimPrefix(a.GetAssetType().String(), "ASSET_TYPE_")),
		AssetCategoryID:    a.GetAssetCategoryId(),
		LocationID:         a.GetLocationId(),
		AcquisitionCost:    float64(a.GetAcquisitionCost()) / 100,
		Currency:           a.GetCurrency(),
		SalvageValue:       float64(a.GetSalvageValue()) / 100,
		BookValue:          float64(a.GetBookValue()) / 100,
		UsefulLifeMonths:   int(a.GetUsefulLifeMonths()),
		DepreciationMethod: strings.ToLower(strings.TrimPrefix(a.GetDepreciationMethod().String(), "DEPRECIATION_METHOD_")),
		Status:             strings.ToLower(strings.TrimPrefix(a.GetStatus().String(), "ASSET_STATUS_")),
		Active:             a.GetActive(),
	}
	if a.Description != nil {
		r.Description = *a.Description
	}
	return r
}

// ---------------------------------------------------------------------------
// Lapsing schedule + policy helpers
// ---------------------------------------------------------------------------

// listCandidatesWorkspace calls ListDepreciationCandidates with scope_kind=WORKSPACE
// and maps the proto response to view-layer CandidateRow slices.
func listCandidatesWorkspace(
	ctx context.Context,
	useCases *consumer.UseCases,
	asOfDate, cursor string,
	limit int32,
) ([]lapsinglist.CandidateRow, string, error) {
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_WORKSPACE,
		AsOfDate:  asOfDate,
		Pagination: &commonpb.PaginationRequest{
			Limit: limit,
			Method: &commonpb.PaginationRequest_Cursor{
				Cursor: &commonpb.CursorPagination{Token: cursor},
			},
		},
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, "", err
	}
	rows := make([]lapsinglist.CandidateRow, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		row := lapsinglist.CandidateRow{
			AssetID:          c.GetAssetId(),
			AssetName:        c.GetAssetName(),
			Currency:         c.GetCurrency(),
			CurrentBookValue: c.GetCurrentBookValue(),
			CanRun:           len(c.GetBlockers()) == 0 && len(c.GetPeriods()) > 0,
		}
		if len(c.GetPeriods()) > 0 {
			row.PendingCount = len(c.GetPeriods())
			row.NextAmount = c.GetPeriods()[0].GetAmount()
			row.NextPendingPeriod = c.GetPeriods()[0].GetPeriodStartDate()
		}
		if len(c.GetBlockers()) > 0 {
			row.Status = "blocked"
			row.BlockerLabel = c.GetBlockers()[0].GetLabel()
		} else if row.PendingCount == 0 {
			row.Status = "up_to_date"
		} else if row.LastPostedPeriod == "" {
			row.Status = "not_started"
		} else {
			row.Status = "pending"
		}
		rows = append(rows, row)
	}
	nextCursor := ""
	if resp.GetPagination() != nil {
		nextCursor = resp.GetPagination().GetNextCursor()
	}
	return rows, nextCursor, nil
}

// listPoliciesWithRollup fetches all AssetCategory rows with per-category aggregate counts.
// Uses the new ListAssetCategoriesWithPolicyRollup use case (Wave 3 espyna enhancement).
// AssetsInPolicy and AssetsDeviating are real counts from the Postgres bulk query.
func listPoliciesWithRollup(
	ctx context.Context,
	useCases *consumer.UseCases,
) ([]assetcatpolicies.PolicyRow, error) {
	rollupRows, err := consumer.ListAssetCategoriesWithPolicyRollup(useCases, ctx)
	if err != nil {
		return nil, err
	}
	rows := make([]assetcatpolicies.PolicyRow, 0, len(rollupRows))
	for _, r := range rollupRows {
		c := r.Category
		if c == nil {
			continue
		}
		method := c.GetDefaultDepreciationMethod()
		if c.DepreciationMethod != nil {
			method = c.GetDepreciationMethod()
		}
		usefulLife := c.GetDefaultUsefulLifeMonths()
		if c.UsefulLifeMonths != nil {
			usefulLife = c.GetUsefulLifeMonths()
		}
		salvage := c.GetDefaultSalvageValuePercent()
		if c.SalvagePct != nil {
			salvage = c.GetSalvagePct()
		}
		rows = append(rows, assetcatpolicies.PolicyRow{
			CategoryID:         c.GetId(),
			PolicyID:           c.GetId(),
			Name:               c.GetName(),
			DepreciationMethod: method,
			UsefulLifeMonths:   usefulLife,
			SalvagePct:         salvage,
			AssetsInPolicy:     r.AssetsInPolicy,
			AssetsDeviating:    r.AssetsDeviating,
		})
	}
	return rows, nil
}

// listCandidatesForPolicy calls ListDepreciationCandidates with scope_kind=POLICY
// and maps results to the preview drawer's PreviewCandidateRow.
func listCandidatesForPolicy(
	ctx context.Context,
	useCases *consumer.UseCases,
	categoryID, asOfDate string,
) ([]assetcataction.PreviewCandidateRow, error) {
	scopeID := categoryID
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_POLICY,
		ScopeId:   &scopeID,
		AsOfDate:  asOfDate,
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, err
	}
	rows := make([]assetcataction.PreviewCandidateRow, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		row := assetcataction.PreviewCandidateRow{
			AssetID:      c.GetAssetId(),
			AssetName:    c.GetAssetName(),
			Currency:     c.GetCurrency(),
			BookValue:    c.GetCurrentBookValue(),
			PendingCount: len(c.GetPeriods()),
		}
		if len(c.GetPeriods()) > 0 {
			row.NextAmount = c.GetPeriods()[0].GetAmount()
		}
		for _, b := range c.GetBlockers() {
			row.Blockers = append(row.Blockers, b.GetLabel())
		}
		rows = append(rows, row)
	}
	return rows, nil
}

// ---------------------------------------------------------------------------
// Surface A helpers — per-asset depreciation-run wrappers
// ---------------------------------------------------------------------------

// listDepreciationCandidatesForAsset calls ListDepreciationCandidates with
// scope_kind=ASSET and maps the proto response to the view-layer DepreciationCandidate slice.
func listDepreciationCandidatesForAsset(
	ctx context.Context,
	useCases *consumer.UseCases,
	assetID, asOfDate string,
) ([]assetaction.DepreciationCandidate, error) {
	if asOfDate == "" {
		asOfDate = "today" // espyna engine accepts "today" as a sentinel
	}
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_ASSET,
		ScopeId:   &assetID,
		AsOfDate:  asOfDate,
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, err
	}
	var rows []assetaction.DepreciationCandidate
	for _, c := range resp.GetData() {
		for _, p := range c.GetPeriods() {
			rows = append(rows, assetaction.DepreciationCandidate{
				PeriodStart:        p.GetPeriodStartDate(),
				PeriodEnd:          p.GetPeriodEndDate(),
				ProjectedAmount:    p.GetAmount(),
				ProjectedAmountFmt: fmt.Sprintf("₱%.2f", float64(p.GetAmount())/100),
				ProjectedAccum:     p.GetRunningAccumulated(),
				ProjectedAccumFmt:  fmt.Sprintf("₱%.2f", float64(p.GetRunningAccumulated())/100),
			})
		}
		for _, b := range c.GetBlockers() {
			rows = append(rows, assetaction.DepreciationCandidate{
				Blocked:       true,
				BlockerReason: b.GetLabel(),
			})
		}
	}
	return rows, nil
}

// generateDepreciationRunForAsset posts selected periods for a single asset
// and maps the proto response back to the view-layer DepreciationRunResult.
func generateDepreciationRunForAsset(
	ctx context.Context,
	useCases *consumer.UseCases,
	req assetaction.DepreciationRunRequest,
) (*assetaction.DepreciationRunResult, error) {
	protoReq := &deprunpb.GenerateDepreciationRunRequest{
		ScopeKind: deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_ASSET,
		ScopeId:   &req.AssetID,
		AsOfDate:  req.AsOfDate,
		Selections: []*deprunpb.DepreciationRunSelection{
			{
				AssetId:          req.AssetID,
				PeriodStartDates: req.PeriodStartDates,
			},
		},
	}
	resp, err := consumer.GenerateDepreciationRun(useCases, ctx, protoReq)
	if err != nil {
		return nil, err
	}
	runID := ""
	if resp.GetRun() != nil {
		runID = resp.GetRun().GetId()
	}
	return &assetaction.DepreciationRunResult{
		RunID:        runID,
		CreatedCount: int(resp.GetCreatedCount()),
		SkippedCount: int(resp.GetSkippedCount()),
		ErroredCount: int(resp.GetErroredCount()),
		Success:      resp.GetSuccess(),
	}, nil
}

// assetToRow converts a proto Asset to the flat assetlist.AssetRow used by the list view.
func assetToRow(a *assetpb.Asset) assetlist.AssetRow {
	row := assetlist.AssetRow{
		ID:              a.GetId(),
		AssetNumber:     a.GetAssetNumber(),
		Name:            a.GetName(),
		AcquisitionCost: float64(a.GetAcquisitionCost()) / 100,
		BookValue:       float64(a.GetBookValue()) / 100,
		Active:          a.GetActive(),
	}
	if a.AssetCategory != nil {
		row.CategoryName = a.AssetCategory.GetName()
	}
	if a.Location != nil {
		row.LocationName = a.Location.GetName()
	}
	return row
}

// ---------------------------------------------------------------------------
// Surface C helpers — per-category / per-policy depreciation-run wrappers
// ---------------------------------------------------------------------------

// listCandidatesForCategory calls ListDepreciationCandidates with scope_kind=CATEGORY or POLICY
// and maps results to the Surface C CategoryDepreciationRunAssetRow slice.
// One row per asset (not per period — the drawer shows which assets to include, not individual periods).
func listCandidatesForCategory(
	ctx context.Context,
	useCases *consumer.UseCases,
	categoryID, scopeKind, asOfDate string,
) ([]assetcataction.CategoryDepreciationRunAssetRow, error) {
	if asOfDate == "" {
		asOfDate = "today"
	}
	proto := deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_CATEGORY
	if scopeKind == "policy" || scopeKind == "POLICY" {
		proto = deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_POLICY
	}
	req := &deprunpb.ListDepreciationCandidatesRequest{
		ScopeKind: proto,
		ScopeId:   &categoryID,
		AsOfDate:  asOfDate,
	}
	resp, err := consumer.ListDepreciationCandidates(useCases, ctx, req)
	if err != nil {
		return nil, err
	}
	rows := make([]assetcataction.CategoryDepreciationRunAssetRow, 0, len(resp.GetData()))
	for _, c := range resp.GetData() {
		row := assetcataction.CategoryDepreciationRunAssetRow{
			AssetID:      c.GetAssetId(),
			AssetName:    c.GetAssetName(),
			Currency:     c.GetCurrency(),
			BookValue:    c.GetCurrentBookValue(),
			PendingCount: len(c.GetPeriods()),
		}
		if len(c.GetPeriods()) > 0 {
			row.NextAmount = c.GetPeriods()[0].GetAmount()
		}
		for _, b := range c.GetBlockers() {
			row.Blockers = append(row.Blockers, b.GetLabel())
		}
		row.CanRun = len(row.Blockers) == 0 && row.PendingCount > 0
		rows = append(rows, row)
	}
	return rows, nil
}

// generateDepreciationRunForCategory posts a depreciation run for all selected assets
// in a category/policy scope and maps the result back to the view-layer type.
func generateDepreciationRunForCategory(
	ctx context.Context,
	useCases *consumer.UseCases,
	req assetcataction.CategoryDepreciationRunRequest,
) (*assetcataction.CategoryDepreciationRunResult, error) {
	protoScope := deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_CATEGORY
	if req.ScopeKind == "POLICY" {
		protoScope = deprunpb.DepreciationRunScopeKind_DEPRECIATION_RUN_SCOPE_KIND_POLICY
	}
	// Build per-asset selections: each selected asset contributes all its pending periods.
	// The use case engine resolves pending periods server-side from the as_of_date; we pass
	// an empty period list per asset so the engine computes them (same as the Revenue Run
	// "all-for-scope" pattern). If asset IDs are empty, scope covers ALL assets in the category.
	var selections []*deprunpb.DepreciationRunSelection
	for _, assetID := range req.AssetIDs {
		aid := assetID
		selections = append(selections, &deprunpb.DepreciationRunSelection{
			AssetId: aid,
			// PeriodStartDates empty → use case posts all pending periods for this asset.
		})
	}
	protoReq := &deprunpb.GenerateDepreciationRunRequest{
		ScopeKind:  protoScope,
		ScopeId:    &req.CategoryID,
		AsOfDate:   req.AsOfDate,
		Selections: selections,
	}
	resp, err := consumer.GenerateDepreciationRun(useCases, ctx, protoReq)
	if err != nil {
		return nil, err
	}
	runID := ""
	if resp.GetRun() != nil {
		runID = resp.GetRun().GetId()
	}
	return &assetcataction.CategoryDepreciationRunResult{
		RunID:        runID,
		CreatedCount: int(resp.GetCreatedCount()),
		SkippedCount: int(resp.GetSkippedCount()),
		ErroredCount: int(resp.GetErroredCount()),
		Success:      resp.GetSuccess(),
	}, nil
}

// ---------------------------------------------------------------------------
// Surface D helpers — depreciation-run history wrappers
// ---------------------------------------------------------------------------

// listDepreciationRunsForWorkspace fetches a page of DepreciationRun rows for Surface D.
func listDepreciationRunsForWorkspace(
	ctx context.Context,
	useCases *consumer.UseCases,
	scope depreciationrunmod.ListDepreciationRunsScope,
) ([]depreciationrunmod.DepreciationRunRow, string, error) {
	req := &deprunpb.ListDepreciationRunsRequest{}
	resp, err := consumer.ListDepreciationRuns(useCases, ctx, req)
	if err != nil {
		return nil, "", err
	}
	rows := make([]depreciationrunmod.DepreciationRunRow, 0, len(resp.GetData()))
	for _, r := range resp.GetData() {
		if scope.Status != "" {
			status := strings.ToLower(strings.TrimPrefix(r.GetStatus().String(), "DEPRECIATION_RUN_STATUS_"))
			if status != scope.Status {
				continue
			}
		}
		rows = append(rows, depreciationRunToRow(r))
	}
	return rows, "", nil
}

// readDepreciationRunWithEntries fetches a single DepreciationRun plus schedule entries.
func readDepreciationRunWithEntries(
	ctx context.Context,
	useCases *consumer.UseCases,
	id string,
) (*depreciationrunmod.DepreciationRunWithEntries, error) {
	resp, err := consumer.ReadDepreciationRun(useCases, ctx, &deprunpb.ReadDepreciationRunRequest{
		Data: &deprunpb.DepreciationRun{Id: id},
	})
	if err != nil {
		return nil, err
	}
	if resp == nil || len(resp.GetData()) == 0 {
		return nil, fmt.Errorf("depreciation run %s not found", id)
	}
	run := depreciationRunToRow(resp.GetData()[0])
	return &depreciationrunmod.DepreciationRunWithEntries{
		Run: run,
	}, nil
}

// depreciationRunToRow maps a proto DepreciationRun to the view-layer DepreciationRunRow.
func depreciationRunToRow(r *deprunpb.DepreciationRun) depreciationrunmod.DepreciationRunRow {
	if r == nil {
		return depreciationrunmod.DepreciationRunRow{}
	}
	status := strings.ToLower(strings.TrimPrefix(r.GetStatus().String(), "DEPRECIATION_RUN_STATUS_"))
	scopeKind := strings.ToLower(strings.TrimPrefix(r.GetScopeKind().String(), "DEPRECIATION_RUN_SCOPE_KIND_"))
	return depreciationrunmod.DepreciationRunRow{
		ID:           r.GetId(),
		WorkspaceID:  r.GetWorkspaceId(),
		ScopeKind:    scopeKind,
		ScopeID:      r.GetScopeId(),
		AsOfDate:     r.GetAsOfDate(),
		InitiatorID:  r.GetInitiatorId(),
		Status:       status,
		CreatedCount: r.GetCreatedCount(),
		SkippedCount: r.GetSkippedCount(),
		ErroredCount: r.GetErroredCount(),
	}
}
