// Package block — typed wiring contract for fycha.Block.
//
// This file declares what fycha's Block() needs from outside.
// Service-admin's composition layer constructs a *UseCases value from
// espyna's consumer container; fycha's Block() consumes only this typed shape.
//
// Shape this struct by what FYCHA needs, NOT by mirroring espyna's
// *consumer.UseCases. Service-admin's adapter is the only place that
// knows both vocabularies.
package block

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	assetcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset_category"
	revaluationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset_revaluation"
	depschpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation"
	deprunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation_run"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	forexratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/finance/forex_rate"
	fundpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund"
	fundallocationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_allocation"
	fundtransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/funding/fund_transaction"
	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"
	journalentrypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
	taxratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/tax/tax_rate"
	withholdinpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	apagingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ap_aging"
	aragingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ar_aging"
	dspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/domain_specific"
	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"
	stmtspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/statements"

	equitydashboard "github.com/erniealice/fycha-golang/domain/ledger/equity/dashboard"
	ledgerdashboard "github.com/erniealice/fycha-golang/domain/ledger/ledger/dashboard"
	payrolldashboard "github.com/erniealice/fycha-golang/domain/payroll/payrolldashboard"
	loansdashboard "github.com/erniealice/fycha-golang/domain/treasury/loan/dashboard"
)

// UseCases declares everything fycha's Block() needs from outside.
// Construction is service-admin's job; fycha only declares the shape.
//
// Naming conventions mirror the entydad block convention:
//  1. Field names are singular matching the domain folder name.
//  2. Group struct types use the `<Entity>UseCases` suffix.
//  3. Closure signatures use proto request/response types directly.
//  4. Dashboard closures use the view-layer types.
type UseCases struct {
	// Workspace.ReadWorkspace is needed by getFunctionalCurrency.
	Workspace WorkspaceUseCases

	// Asset domain groups
	Asset       AssetUseCases
	DepRun      DepreciationRunUseCases
	Revaluation RevaluationUseCases

	// Ledger domain groups
	Ledger       LedgerUseCases
	FiscalPeriod FiscalPeriodUseCases

	// Funding domain group
	Funding FundingUseCases

	// Other accounting groups
	Tax      TaxUseCases
	Finance  FinanceUseCases
	Treasury TreasuryUseCases

	// Reports — service-driven report use case closures. Populated in
	// service-admin's buildFychaUseCases from
	// `consumer.UseCases.Service.Reporting.<Group>.<UseCase>.Execute`.
	// Each P1.E.* sub-candidate adds its own group struct here.
	Reports ReportsUseCases

	// Dashboard closures — typed; service-admin builds these from espyna's
	// internal use cases and maps their internal response types to fycha's
	// view-layer types. This removes the need for reflection in dashboard_wiring.go.
	GetLedgerDashboardPageData  func(ctx context.Context, req *ledgerdashboard.Request) (*ledgerdashboard.Response, error)
	GetEquityDashboardPageData  func(ctx context.Context, req *equitydashboard.Request) (*equitydashboard.Response, error)
	GetPayrollDashboardPageData func(ctx context.Context, req *payrolldashboard.Request) (*payrolldashboard.Response, error)
	GetLoanDashboardPageData    func(ctx context.Context, req *loansdashboard.Request) (*loansdashboard.Response, error)
}

// ReportsUseCases aggregates the service-driven report use case closures
// the fycha report views need. Wave B P1.E.* sub-candidates each fill in
// one of the sub-groups below.
type ReportsUseCases struct {
	// ARAging holds the AR-side reporting closures (Wave B P1.E.1).
	ARAging ARAgingUseCases
	// APAging holds the AP-side reporting closures (Wave B P1.E.2).
	APAging APAgingUseCases
	// GrossCashFlow holds the gross-profit + cash-book closures (Wave B P1.E.3).
	GrossCashFlow GrossCashFlowUseCases
	// DomainSpecific holds the revenue/expenditure/disbursement pivot
	// closures + the Go-only ListRevenue/ListExpenses feeders (Wave B P1.E.5).
	DomainSpecific DomainSpecificUseCases
	// Statements holds the counterparty statement + balance closures
	// (Wave B P1.E.4). The closures live on fycha's block UseCases for
	// symmetry with the other report groups; the actual consumers are
	// largely the entydad client/supplier detail/list views and the
	// service-admin adapter wires entydad block UseCases from these.
	Statements StatementsUseCases
}

// ARAgingUseCases — service-driven AR aging report closures consumed by
// `views/reports/receivables_aging_report` and
// `views/reports/collection_summary_report`. Migrated 2026-05-20 out of
// `fycha.DataSource` into the proto-shaped service layer.
type ARAgingUseCases struct {
	GetReceivablesAgingReport  func(context.Context, *aragingpb.GetReceivablesAgingRequest) (*aragingpb.GetReceivablesAgingResponse, error)
	GetCollectionSummaryReport func(context.Context, *aragingpb.GetCollectionSummaryRequest) (*aragingpb.GetCollectionSummaryResponse, error)
}

// APAgingUseCases — service-driven AP aging report closures consumed by
// `views/reports/payables_aging_report` (parameterized) and
// `views/reports/payables_aging` (simple). Migrated 2026-05-21 out of
// `fycha.DataSource` into the proto-shaped service layer (Wave B P1.E.2).
type APAgingUseCases struct {
	GetPayablesAgingReport       func(context.Context, *apagingpb.GetPayablesAgingRequest) (*apagingpb.GetPayablesAgingResponse, error)
	GetSimplePayablesAgingReport func(context.Context, *apagingpb.GetSimplePayablesAgingRequest) (*apagingpb.GetSimplePayablesAgingResponse, error)
}

// GrossCashFlowUseCases — service-driven gross-profit + cash-book closures
// consumed by `views/reports/gross_profit`, `views/reports/cost_of_sales`,
// `views/reports/net_profit`, `views/reports/dashboard`, and
// `views/reports/cash_book`. Migrated 2026-05-21 out of `fycha.DataSource`
// into the proto-shaped service layer (Wave B P1.E.3).
type GrossCashFlowUseCases struct {
	GetGrossProfitReport func(context.Context, *gcfpb.GetGrossProfitRequest) (*gcfpb.GetGrossProfitResponse, error)
	GetCashBookReport    func(context.Context, *gcfpb.GetCashBookRequest) (*gcfpb.GetCashBookResponse, error)
}

// DomainSpecificUseCases — service-driven domain-specific report closures
// consumed by `views/reports/revenue_report`, `views/reports/expenditure_report`,
// `views/reports/disbursement_report`, plus the Go-only feeders consumed by
// `views/reports/revenue`, `views/reports/expenses`, `views/reports/net_profit`,
// and `views/reports/dashboard`. Migrated 2026-05-21 out of `fycha.DataSource`
// into the proto-shaped service layer (Wave B P1.E.5).
//
// ListRevenue and ListExpenses retain the Go-only `[]map[string]any` shape
// per Q-SDM-MAP-SHAPES — see
// `packages/espyna-golang/internal/application/usecases/service/reporting/
// domain_specific/list_revenue.go` for the rationale.
type DomainSpecificUseCases struct {
	GetRevenueReport      func(context.Context, *dspb.GetRevenueReportRequest) (*dspb.GetRevenueReportResponse, error)
	GetExpenditureReport  func(context.Context, *dspb.GetExpenditureReportRequest) (*dspb.GetExpenditureReportResponse, error)
	GetDisbursementReport func(context.Context, *dspb.GetDisbursementReportRequest) (*dspb.GetDisbursementReportResponse, error)
	ListRevenue           func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
	ListExpenses          func(context.Context, *time.Time, *time.Time) ([]map[string]any, error)
}

// StatementsUseCases — service-driven counterparty statement + balance
// closures. Lives on the fycha block UseCases struct for symmetry with the
// other report groups, but the consumers are largely under entydad
// (client/supplier detail/list views). Migrated 2026-05-21 out of
// `centymo/entydad.LedgerReportingService` into the proto-shaped service
// layer (Wave B P1.E.4).
//
// **Map shim:** the legacy entydad views accept
// `func(ctx) (map[string]int64, error)` for balances; the typed closures
// here expose the new `[]BalanceRow` shape. Service-admin's adapter
// constructs map-returning wrappers for entydad and exposes the typed
// closures for any future fycha consumer.
type StatementsUseCases struct {
	GetClientStatement   func(context.Context, *stmtspb.GetClientStatementRequest) (*stmtspb.GetClientStatementResponse, error)
	GetSupplierStatement func(context.Context, *stmtspb.GetSupplierStatementRequest) (*stmtspb.GetSupplierStatementResponse, error)
	ListClientBalances   func(context.Context, *stmtspb.ListClientBalancesRequest) (*stmtspb.ListClientBalancesResponse, error)
	ListSupplierBalances func(context.Context, *stmtspb.ListSupplierBalancesRequest) (*stmtspb.ListSupplierBalancesResponse, error)
}

// WorkspaceUseCases — subset needed by fycha (functional currency lookup).
type WorkspaceUseCases struct {
	Read func(context.Context, *workspacepb.ReadWorkspaceRequest) (*workspacepb.ReadWorkspaceResponse, error)
}

// AssetUseCases — direct CRUD on the Asset entity.
type AssetUseCases struct {
	GetListPageData func(context.Context, *assetpb.GetAssetListPageDataRequest) (*assetpb.GetAssetListPageDataResponse, error)
	Create          func(context.Context, *assetpb.CreateAssetRequest) (*assetpb.CreateAssetResponse, error)
	Read            func(context.Context, *assetpb.ReadAssetRequest) (*assetpb.ReadAssetResponse, error)
	Update          func(context.Context, *assetpb.UpdateAssetRequest) (*assetpb.UpdateAssetResponse, error)
	SetActive       func(context.Context, *assetpb.SetAssetActiveRequest) (*assetpb.SetAssetActiveResponse, error)
	// Category wraps the asset_category sub-domain ops needed by fycha.
	Category AssetCategoryUseCases
}

type AssetCategoryUseCases struct {
	ListWithPolicyRollup func(context.Context, *assetcategorypb.ListAssetCategoriesWithPolicyRollupRequest) (*assetcategorypb.ListAssetCategoriesWithPolicyRollupResponse, error)
}

// DepreciationRunUseCases — depreciation-run sub-domain operations.
type DepreciationRunUseCases struct {
	ListCandidates func(context.Context, *deprunpb.ListDepreciationCandidatesRequest) (*deprunpb.ListDepreciationCandidatesResponse, error)
	Generate       func(context.Context, *deprunpb.GenerateDepreciationRunRequest) (*deprunpb.GenerateDepreciationRunResponse, error)
	List           func(context.Context, *deprunpb.ListDepreciationRunsRequest) (*deprunpb.ListDepreciationRunsResponse, error)
	Read           func(context.Context, *deprunpb.ReadDepreciationRunRequest) (*deprunpb.ReadDepreciationRunResponse, error)
	ListEntries    func(context.Context, *deprunpb.ListDepreciationRunEntriesRequest) (*depschpb.ListDepreciationSchedulesResponse, error)
}

// RevaluationUseCases — asset revaluation sub-domain operations.
type RevaluationUseCases struct {
	Revalue func(context.Context, *revaluationpb.RevalueAssetUseCaseRequest) (*revaluationpb.RevalueAssetUseCaseResponse, error)
	Preview func(context.Context, *revaluationpb.PreviewRevaluationUseCaseRequest) (*revaluationpb.PreviewRevaluationUseCaseResponse, error)
}

// LedgerUseCases — Chart of Accounts + Journal Entry operations.
type LedgerUseCases struct {
	Account      AccountUseCases
	JournalEntry JournalEntryUseCases
}

type AccountUseCases struct {
	GetListPageData func(context.Context, *accountpb.GetAccountListPageDataRequest) (*accountpb.GetAccountListPageDataResponse, error)
	Create          func(context.Context, *accountpb.CreateAccountRequest) (*accountpb.CreateAccountResponse, error)
	Read            func(context.Context, *accountpb.ReadAccountRequest) (*accountpb.ReadAccountResponse, error)
	Update          func(context.Context, *accountpb.UpdateAccountRequest) (*accountpb.UpdateAccountResponse, error)
	Delete          func(context.Context, *accountpb.DeleteAccountRequest) (*accountpb.DeleteAccountResponse, error)
}

type JournalEntryUseCases struct {
	GetListPageData func(context.Context, *journalentrypb.GetJournalEntryListPageDataRequest) (*journalentrypb.GetJournalEntryListPageDataResponse, error)
	Create          func(context.Context, *journalentrypb.CreateJournalEntryRequest) (*journalentrypb.CreateJournalEntryResponse, error)
	Read            func(context.Context, *journalentrypb.ReadJournalEntryRequest) (*journalentrypb.ReadJournalEntryResponse, error)
	Update          func(context.Context, *journalentrypb.UpdateJournalEntryRequest) (*journalentrypb.UpdateJournalEntryResponse, error)
	Delete          func(context.Context, *journalentrypb.DeleteJournalEntryRequest) (*journalentrypb.DeleteJournalEntryResponse, error)
	Post            func(context.Context, *journalentrypb.PostJournalEntryRequest) (*journalentrypb.PostJournalEntryResponse, error)
	Reverse         func(context.Context, *journalentrypb.ReverseJournalEntryRequest) (*journalentrypb.ReverseJournalEntryResponse, error)
}

type FiscalPeriodUseCases struct {
	GetListPageData func(context.Context, *fiscalperiodpb.GetFiscalPeriodListPageDataRequest) (*fiscalperiodpb.GetFiscalPeriodListPageDataResponse, error)
	Create          func(context.Context, *fiscalperiodpb.CreateFiscalPeriodRequest) (*fiscalperiodpb.CreateFiscalPeriodResponse, error)
	Close           func(context.Context, *fiscalperiodpb.CloseFiscalPeriodRequest) (*fiscalperiodpb.CloseFiscalPeriodResponse, error)
}

// TaxUseCases — tax rate read-only ops.
type TaxUseCases struct {
	ListTaxRates func(context.Context, *taxratepb.ListTaxRatesRequest) (*taxratepb.ListTaxRatesResponse, error)
}

// FinanceUseCases — forex rate read-only ops.
type FinanceUseCases struct {
	ListForexRates func(context.Context, *forexratepb.ListForexRatesRequest) (*forexratepb.ListForexRatesResponse, error)
}

// FundingUseCases — cross-workspace fund source/card/transaction read ops
// consumed by the funding view module (source list, source detail, card list,
// card detail). Write ops (create/update/delete) are not wired yet (FS-E stubs).
type FundingUseCases struct {
	ReadFund         func(context.Context, *fundpb.ReadFundRequest) (*fundpb.ReadFundResponse, error)
	ListFunds        func(context.Context, *fundpb.ListFundsRequest) (*fundpb.ListFundsResponse, error)
	ReadAllocation   func(context.Context, *fundallocationpb.ReadFundAllocationRequest) (*fundallocationpb.ReadFundAllocationResponse, error)
	ListAllocations  func(context.Context, *fundallocationpb.ListFundAllocationsRequest) (*fundallocationpb.ListFundAllocationsResponse, error)
	ListTransactions func(context.Context, *fundtransactionpb.ListFundTransactionsRequest) (*fundtransactionpb.ListFundTransactionsResponse, error)
}

// TreasuryUseCases — withholding certificate and other treasury ops.
type TreasuryUseCases struct {
	ListWithholdingCertificates func(context.Context, *withholdinpb.ListWithholdingCertificatesRequest) (*withholdinpb.ListWithholdingCertificatesResponse, error)
}

// RequireFor returns an error listing every needed-but-nil field for cfg's
// enabled modules. Called (via MustValidate) at Block() entry; a missing field
// for an enabled module → startup error rather than a silently-empty render.
//
// CRITICAL: this is the deterministic completeness check. Required-vs-optional
// lives ENTIRELY here and is correct by construction — a field is REQUIRED iff
// it is `check(...)`-asserted inside an enabled module's `if cfg.wantXxx()`
// block; a field never asserted is OPTIONAL (legitimately nil → graceful
// empty-state). The REQUIRED set is the render-path the module cannot function
// without: each report's report/list/get closure, the asset/ledger CRUD + list
// + page-data closures, and the read-only modules' single list closure.
//
// NOT checked (intentionally OPTIONAL, nil-able — each marked below):
//   - The four dashboard closures (Ledger/Equity/Payroll/Loan) — wired nil-safe
//     in dashboard_wiring.go; nil → the dashboard view renders empty-state.
//   - Loans / Equity / Payroll modules — dashboard-only today; their non-nil
//     surface is exactly the (optional) dashboard closures above.
//   - Cash / Expenses / Financial modules — registered with TODO-stub deps (no
//     espyna use case yet); nothing to assert until those are wired.
//   - Workspace.Read — functional-currency lookup; nil → default-currency path.
//   - FiscalPeriod.{Create,Close} — block.go nil-guards them independently of
//     the list closure (the list page renders without the mutators).
//   - Revaluation.{Revalue,Preview} — Surface-E preview/commit; nil-guarded in
//     asset.go, degrades to a disabled revaluation CTA.
func (u *UseCases) RequireFor(cfg *blockConfig) error {
	if u == nil {
		return fmt.Errorf("fycha.Block: WithUseCases(...) was not supplied")
	}

	var missing []string
	check := func(ok bool, name string) {
		if !ok {
			missing = append(missing, name)
		}
	}

	if cfg.wantReports() {
		// The report views' render path calls these report/list/get closures
		// directly (block.go reportmod.ModuleDeps). Each is REQUIRED whenever the
		// Reports module is enabled — without it the corresponding report renders
		// blank with no error.
		ar := &u.Reports.ARAging
		check(ar.GetReceivablesAgingReport != nil, "UseCases.Reports.ARAging.GetReceivablesAgingReport")
		check(ar.GetCollectionSummaryReport != nil, "UseCases.Reports.ARAging.GetCollectionSummaryReport")
		ap := &u.Reports.APAging
		check(ap.GetPayablesAgingReport != nil, "UseCases.Reports.APAging.GetPayablesAgingReport")
		gcf := &u.Reports.GrossCashFlow
		check(gcf.GetGrossProfitReport != nil, "UseCases.Reports.GrossCashFlow.GetGrossProfitReport")
		ds := &u.Reports.DomainSpecific
		check(ds.GetRevenueReport != nil, "UseCases.Reports.DomainSpecific.GetRevenueReport")
		check(ds.GetExpenditureReport != nil, "UseCases.Reports.DomainSpecific.GetExpenditureReport")
		check(ds.GetDisbursementReport != nil, "UseCases.Reports.DomainSpecific.GetDisbursementReport")
		check(ds.ListRevenue != nil, "UseCases.Reports.DomainSpecific.ListRevenue")
		check(ds.ListExpenses != nil, "UseCases.Reports.DomainSpecific.ListExpenses")
	}

	if cfg.wantAsset() {
		a := &u.Asset
		check(a.GetListPageData != nil, "UseCases.Asset.GetListPageData")
		check(a.Create != nil, "UseCases.Asset.Create")
		check(a.Read != nil, "UseCases.Asset.Read")
		check(a.Update != nil, "UseCases.Asset.Update")
		check(a.SetActive != nil, "UseCases.Asset.SetActive")
		check(a.Category.ListWithPolicyRollup != nil, "UseCases.Asset.Category.ListWithPolicyRollup")
		// Depreciation-run list/read/generate back the Surface-A run views.
		dr := &u.DepRun
		check(dr.ListCandidates != nil, "UseCases.DepRun.ListCandidates")
		check(dr.Generate != nil, "UseCases.DepRun.Generate")
		check(dr.List != nil, "UseCases.DepRun.List")
		check(dr.Read != nil, "UseCases.DepRun.Read")
		check(dr.ListEntries != nil, "UseCases.DepRun.ListEntries")
	}

	if cfg.wantLedger() {
		acc := &u.Ledger.Account
		check(acc.GetListPageData != nil, "UseCases.Ledger.Account.GetListPageData")
		check(acc.Create != nil, "UseCases.Ledger.Account.Create")
		check(acc.Read != nil, "UseCases.Ledger.Account.Read")
		check(acc.Update != nil, "UseCases.Ledger.Account.Update")
		check(acc.Delete != nil, "UseCases.Ledger.Account.Delete")
		je := &u.Ledger.JournalEntry
		check(je.GetListPageData != nil, "UseCases.Ledger.JournalEntry.GetListPageData")
		check(je.Create != nil, "UseCases.Ledger.JournalEntry.Create")
		check(je.Read != nil, "UseCases.Ledger.JournalEntry.Read")
		check(je.Update != nil, "UseCases.Ledger.JournalEntry.Update")
		check(je.Delete != nil, "UseCases.Ledger.JournalEntry.Delete")
		check(je.Post != nil, "UseCases.Ledger.JournalEntry.Post")
		check(je.Reverse != nil, "UseCases.Ledger.JournalEntry.Reverse")
		// The fiscal-period list page-data feeder backs the Ledger settings tab.
		check(u.FiscalPeriod.GetListPageData != nil, "UseCases.FiscalPeriod.GetListPageData")
	}

	if cfg.wantTaxRate() {
		check(u.Tax.ListTaxRates != nil, "UseCases.Tax.ListTaxRates")
	}

	if cfg.wantForexRate() {
		check(u.Finance.ListForexRates != nil, "UseCases.Finance.ListForexRates")
	}

	if cfg.wantWithholdingCertificate() {
		check(u.Treasury.ListWithholdingCertificates != nil, "UseCases.Treasury.ListWithholdingCertificates")
	}

	if len(missing) > 0 {
		return fmt.Errorf("fycha.Block: incomplete UseCases — missing %v", missing)
	}
	return nil
}

// MustValidate is the FAIL-CLOSED enforcement wrapper around RequireFor. It is
// the seam-level guard that makes a missing REQUIRED closure impossible to
// ignore — mirroring the AUTHZ_ENFORCE boot-guard in service-admin's
// container.go (a missing security precondition is a boot REFUSAL, never a
// silent degrade).
//
// Why a wrapper and not just `return RequireFor(...)`: a bare returned error is
// fail-OPEN by convention. A caller can drop it (`_ =`, an ignored value, a
// future app that doesn't check) and the block silently registers an empty
// feature — the exact nil-closure trap the architecture roast (burn #1) named.
// MustValidate removes that escape hatch:
//
//   - In dev/test (running under `go test`, OR FYCHA_BLOCK_STRICT truthy) a
//     missing REQUIRED closure PANICS with the full field list. A panic cannot
//     be silently dropped, prints a stack trace at the offending wiring site,
//     and fails the test/CI loudly. This is where a developer wiring a new
//     entity discovers a gap — at their desk, not in prod.
//   - In prod a missing REQUIRED closure logs a screaming FATAL line at the
//     seam (so even a caller that drops the returned error leaves an
//     unmissable log record) AND returns the error so Block() propagates it and
//     NewServiceAdmin halts boot with a clear "domain block failed" message.
//
// OPTIONAL ports (the four dashboard closures, the dashboard-only Loans/Equity/
// Payroll modules, the TODO-stub Cash/Expenses/Financial modules, Workspace.Read,
// FiscalPeriod mutators, Revaluation) are NEVER flagged — that required-vs-
// optional discrimination lives entirely in RequireFor, which only asserts a
// field when its enabling cfg.wantXxx() module is on. MustValidate adds posture,
// not policy: it changes HOW a gap fails, not WHICH fields gate.
func (u *UseCases) MustValidate(cfg *blockConfig) error {
	err := u.RequireFor(cfg)
	if err == nil {
		return nil
	}
	if blockStrictMode() {
		// Dev/test: loud, uncatchable-by-accident, stack-traced.
		panic("FATAL: " + err.Error() + " — REQUIRED block wiring is nil. " +
			"Fix the closure assignment in service-admin's buildFychaUseCases " +
			"(adapters.go) before this reaches prod.")
	}
	// Prod: scream at the seam, then return so boot halts. The log line is the
	// belt to the returned-error's suspenders (a dropped error still screams).
	log.Printf("FATAL: %v — refusing to register fycha modules with a nil "+
		"REQUIRED closure (fail-closed wiring).", err)
	return err
}

// blockStrictMode reports whether the fail-closed wiring guard should PANIC
// (dev/test) rather than return-and-log (prod) on a missing REQUIRED closure.
//
// True when running under `go test` (testing.Testing(), Go 1.21+ — zero env
// coupling, auto-on in every test + CI run) OR when FYCHA_BLOCK_STRICT is set to
// an explicit truthy value (the dev escape hatch for `go run` smoke tests).
// The env matching mirrors container.go's authzEnforceEnabled — anything else
// (unset, "", "0", "false") is prod posture.
func blockStrictMode() bool {
	if testing.Testing() {
		return true
	}
	switch os.Getenv("FYCHA_BLOCK_STRICT") {
	case "1", "true", "TRUE", "True", "yes", "on":
		return true
	default:
		return false
	}
}
