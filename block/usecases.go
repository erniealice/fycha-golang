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

	// Dashboard closures — typed; service-admin builds these from espyna's
	// internal use cases and maps their internal response types to fycha's
	// view-layer types. This removes the need for reflection in dashboard_wiring.go.
	GetLedgerDashboardPageData  func(ctx context.Context, req *ledgerdashboard.Request) (*ledgerdashboard.Response, error)
	GetEquityDashboardPageData  func(ctx context.Context, req *equitydashboard.Request) (*equitydashboard.Response, error)
	GetPayrollDashboardPageData func(ctx context.Context, req *payrolldashboard.Request) (*payrolldashboard.Response, error)
	GetLoanDashboardPageData    func(ctx context.Context, req *loansdashboard.Request) (*loansdashboard.Response, error)
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
