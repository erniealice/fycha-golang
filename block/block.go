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
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"

	fycha "github.com/erniealice/fycha-golang"
	assetmod "github.com/erniealice/fycha-golang/views/asset"
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
	depreciationview "github.com/erniealice/fycha-golang/views/reports/depreciation_policies"
	lapsingview "github.com/erniealice/fycha-golang/views/reports/lapsing_schedule"
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

			assetmod.NewModule(assetDeps).RegisterRoutes(ctx.Routes)

			// Assets → Reports → Lapsing Schedule
			ctx.Routes.GET(assetRoutes.LapsingScheduleURL, lapsingview.NewLapsingScheduleView(ctx.Common, ctx.Table))
			// Assets → Settings → Depreciation Policies
			ctx.Routes.GET(assetRoutes.DepreciationPoliciesURL, depreciationview.NewDepreciationPoliciesView(ctx.Common, ctx.Table))
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
