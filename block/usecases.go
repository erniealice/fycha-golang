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
	"time"

	assetpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset"
	assetcategorypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset_category"
	revaluationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/asset_revaluation"
	depschpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation"
	deprunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/asset/depreciation_run"
	workspacepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/workspace"
	forexratepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/finance/forex_rate"
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

	equitydashboard "github.com/erniealice/fycha-golang/views/equity/dashboard"
	ledgerdashboard "github.com/erniealice/fycha-golang/views/ledger/dashboard"
	loansdashboard "github.com/erniealice/fycha-golang/views/loans/dashboard"
	payrolldashboard "github.com/erniealice/fycha-golang/views/payroll/dashboard"
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

// TreasuryUseCases — withholding certificate and other treasury ops.
type TreasuryUseCases struct {
	ListWithholdingCertificates func(context.Context, *withholdinpb.ListWithholdingCertificatesRequest) (*withholdinpb.ListWithholdingCertificatesResponse, error)
}
