package block

import (
	"context"

	equitydashboardview "github.com/erniealice/fycha-golang/domain/ledger/equity/dashboard"
	ledgerdashboardview "github.com/erniealice/fycha-golang/domain/ledger/ledger/dashboard"
	payrolldashboardview "github.com/erniealice/fycha-golang/domain/payroll/payrolldashboard"
	loansdashboardview "github.com/erniealice/fycha-golang/domain/treasury/loan/dashboard"

	"github.com/erniealice/espyna-golang/consumer"
	consumerapp "github.com/erniealice/espyna-golang/consumer/app"
	"github.com/erniealice/espyna-golang/reference"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	equitydashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/equity"
	ledgerdashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/ledger"
	payrolldashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/payroll"
	treasurydashpb "github.com/erniealice/esqyma/pkg/schema/v1/service/dashboard/treasury"
)

// fychaEngineBlock returns a consumerapp.AppOption that registers all fycha domain
// modules via the compose engine (replaces legacy fychaBlock).
func EngineBlock(depRunURL string) consumerapp.AppOption {
	return func(ctx *consumerapp.AppContext) error {
		uc, err := consumerapp.RequireUseCases(ctx, "fychaEngineBlock")
		if err != nil {
			return err
		}
		adapted := buildFychaUseCases(uc)

		infra := &Infra{AssetDepreciationRunURL: depRunURL}
		infra.UploadFile, _ = ctx.UploadFile.(func(context.Context, string, string, []byte, string) error)
		infra.ListAttachments, _ = ctx.ListAttachments.(func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error))
		infra.CreateAttachment, _ = ctx.CreateAttachment.(func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error))
		infra.DeleteAttachment, _ = ctx.DeleteAttachment.(func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error))
		infra.NewAttachmentID, _ = ctx.NewAttachmentID.(func() string)
		if ctx.RefChecker != nil {
			if rc, ok := ctx.RefChecker.(reference.Checker); ok {
				infra.RefChecker = rc
			}
		}

		units := AllUnits(adapted, infra)
		return consumerapp.AssembleEngineBlock("fycha", units, ctx)
	}
}

// ---------------------------------------------------------------------------
// fycha adapter
// ---------------------------------------------------------------------------

// buildFychaUseCases maps espyna's *consumer.UseCases to fycha block's
// typed shape. All sub-group wiring is nil-safe.
func buildFychaUseCases(uc *consumer.UseCases) *UseCases {
	result := &UseCases{}

	// Workspace.Read — for functional currency lookup
	if uc.Entity != nil && uc.Entity.Workspace != nil {
		result.Workspace.Read = uc.Entity.Workspace.ReadWorkspace.Execute
	}

	// Asset domain
	if uc.Asset != nil {
		if uc.Asset.Asset != nil {
			result.Asset.GetListPageData = uc.Asset.Asset.GetAssetListPageData.Execute
			result.Asset.Create = uc.Asset.Asset.CreateAsset.Execute
			result.Asset.Read = uc.Asset.Asset.ReadAsset.Execute
			result.Asset.Update = uc.Asset.Asset.UpdateAsset.Execute
			result.Asset.SetActive = uc.Asset.Asset.SetAssetActive.Execute
		}
		if uc.Asset.AssetCategory != nil {
			result.Asset.Category.ListWithPolicyRollup = uc.Asset.AssetCategory.ListAssetCategoriesWithPolicyRollup.Execute
		}
		if uc.Asset.DepreciationRun != nil {
			result.DepRun.ListCandidates = uc.Asset.DepreciationRun.ListDepreciationCandidates.Execute
			result.DepRun.Generate = uc.Asset.DepreciationRun.GenerateDepreciationRun.Execute
			result.DepRun.List = uc.Asset.DepreciationRun.ListDepreciationRuns.Execute
			result.DepRun.Read = uc.Asset.DepreciationRun.ReadDepreciationRun.Execute
			result.DepRun.ListEntries = uc.Asset.DepreciationRun.ListDepreciationRunEntries.Execute
		}
		if uc.Asset.AssetRevaluation != nil {
			result.Revaluation.Revalue = uc.Asset.AssetRevaluation.RevalueAsset.Execute
			result.Revaluation.Preview = uc.Asset.AssetRevaluation.PreviewRevaluation.Execute
		}
	}

	// Ledger domain
	if uc.Ledger != nil {
		if uc.Ledger.Account != nil {
			result.Ledger.Account.GetListPageData = uc.Ledger.Account.GetAccountListPageData.Execute
			result.Ledger.Account.Create = uc.Ledger.Account.CreateAccount.Execute
			result.Ledger.Account.Read = uc.Ledger.Account.ReadAccount.Execute
			result.Ledger.Account.Update = uc.Ledger.Account.UpdateAccount.Execute
			result.Ledger.Account.Delete = uc.Ledger.Account.DeleteAccount.Execute
		}
		if uc.Ledger.JournalEntry != nil {
			result.Ledger.JournalEntry.GetListPageData = uc.Ledger.JournalEntry.GetJournalEntryListPageData.Execute
			result.Ledger.JournalEntry.Create = uc.Ledger.JournalEntry.CreateJournalEntry.Execute
			result.Ledger.JournalEntry.Read = uc.Ledger.JournalEntry.ReadJournalEntry.Execute
			result.Ledger.JournalEntry.Update = uc.Ledger.JournalEntry.UpdateJournalEntry.Execute
			result.Ledger.JournalEntry.Delete = uc.Ledger.JournalEntry.DeleteJournalEntry.Execute
			result.Ledger.JournalEntry.Post = uc.Ledger.JournalEntry.PostJournalEntry.Execute
			result.Ledger.JournalEntry.Reverse = uc.Ledger.JournalEntry.ReverseJournalEntry.Execute
		}
		if uc.Ledger.FiscalPeriod != nil {
			result.FiscalPeriod.GetListPageData = uc.Ledger.FiscalPeriod.GetFiscalPeriodListPageData.Execute
			result.FiscalPeriod.Create = uc.Ledger.FiscalPeriod.CreateFiscalPeriod.Execute
			result.FiscalPeriod.Close = uc.Ledger.FiscalPeriod.CloseFiscalPeriod.Execute
		}
	}

	// Tax domain
	if uc.Tax != nil && uc.Tax.TaxRate != nil {
		result.Tax.ListTaxRates = uc.Tax.TaxRate.ListTaxRates.Execute
	}

	// Finance domain (forex rates)
	if uc.Finance != nil && uc.Finance.ForexRate != nil {
		result.Finance.ListForexRates = uc.Finance.ForexRate.ListForexRates.Execute
	}

	// Treasury (withholding certificates)
	if uc.Treasury != nil && uc.Treasury.WithholdingCertificate != nil {
		result.Treasury.ListWithholdingCertificates = uc.Treasury.WithholdingCertificate.ListWithholdingCertificates.Execute
	}

	// Funding domain (cross-workspace fund sources + cards + transactions)
	if uc.Funding != nil {
		if uc.Funding.Fund != nil {
			if uc.Funding.Fund.Read != nil {
				result.Funding.ReadFund = uc.Funding.Fund.Read.Execute
			}
			if uc.Funding.Fund.List != nil {
				result.Funding.ListFunds = uc.Funding.Fund.List.Execute
			}
		}
		if uc.Funding.FundAllocation != nil {
			if uc.Funding.FundAllocation.Read != nil {
				result.Funding.ReadAllocation = uc.Funding.FundAllocation.Read.Execute
			}
			if uc.Funding.FundAllocation.List != nil {
				result.Funding.ListAllocations = uc.Funding.FundAllocation.List.Execute
			}
		}
		if uc.Funding.FundTransaction != nil {
			if uc.Funding.FundTransaction.List != nil {
				result.Funding.ListTransactions = uc.Funding.FundTransaction.List.Execute
			}
		}
	}

	// 20260520-21 Wave B P1.E.1-P1.E.5 — service-driven reporting closures.
	// Captured directly off the typed `Reporting.<Group>` sub-aggregates;
	// available whenever the wrapping service umbrella is non-nil (the use
	// cases are built nil-Reporter on mock builds but still construct).
	//
	// Each `<UseCase>.Execute` method binding here is the integration
	// point between service-admin and fycha's report views — those views
	// no longer reach through the legacy `fycha.DataSource` duck
	// interface, they consume the closures below.
	if uc.Service != nil && uc.Service.Reporting != nil {
		// P1.E.1 — AR aging
		if ar := uc.Service.Reporting.ARAging; ar != nil {
			if ar.GetReceivablesAgingReport != nil {
				result.Reports.ARAging.GetReceivablesAgingReport = ar.GetReceivablesAgingReport.Execute
			}
			if ar.GetCollectionSummaryReport != nil {
				result.Reports.ARAging.GetCollectionSummaryReport = ar.GetCollectionSummaryReport.Execute
			}
		}
		// P1.E.2 — AP aging
		if ap := uc.Service.Reporting.APAging; ap != nil {
			if ap.GetPayablesAgingReport != nil {
				result.Reports.APAging.GetPayablesAgingReport = ap.GetPayablesAgingReport.Execute
			}
			if ap.GetSimplePayablesAgingReport != nil {
				result.Reports.APAging.GetSimplePayablesAgingReport = ap.GetSimplePayablesAgingReport.Execute
			}
		}
		// P1.E.3 — Gross/CashFlow
		if gcf := uc.Service.Reporting.GrossCashFlow; gcf != nil {
			if gcf.GetGrossProfitReport != nil {
				result.Reports.GrossCashFlow.GetGrossProfitReport = gcf.GetGrossProfitReport.Execute
			}
			if gcf.GetCashBookReport != nil {
				result.Reports.GrossCashFlow.GetCashBookReport = gcf.GetCashBookReport.Execute
			}
		}
		// P1.E.4 — Statements (fycha consumes typed closures for future
		// migration; entydad uses the entydad-block adapter below for the
		// map-shaped balance shim).
		if st := uc.Service.Reporting.Statements; st != nil {
			if st.GetClientStatement != nil {
				result.Reports.Statements.GetClientStatement = st.GetClientStatement.Execute
			}
			if st.GetSupplierStatement != nil {
				result.Reports.Statements.GetSupplierStatement = st.GetSupplierStatement.Execute
			}
			if st.ListClientBalances != nil {
				result.Reports.Statements.ListClientBalances = st.ListClientBalances.Execute
			}
			if st.ListSupplierBalances != nil {
				result.Reports.Statements.ListSupplierBalances = st.ListSupplierBalances.Execute
			}
		}
		// P1.E.5 — DomainSpecific
		if ds := uc.Service.Reporting.DomainSpecific; ds != nil {
			if ds.GetRevenueReport != nil {
				result.Reports.DomainSpecific.GetRevenueReport = ds.GetRevenueReport.Execute
			}
			if ds.GetExpenditureReport != nil {
				result.Reports.DomainSpecific.GetExpenditureReport = ds.GetExpenditureReport.Execute
			}
			if ds.GetDisbursementReport != nil {
				result.Reports.DomainSpecific.GetDisbursementReport = ds.GetDisbursementReport.Execute
			}
			if ds.ListRevenue != nil {
				result.Reports.DomainSpecific.ListRevenue = ds.ListRevenue.Execute
			}
			if ds.ListExpenses != nil {
				result.Reports.DomainSpecific.ListExpenses = ds.ListExpenses.Execute
			}
		}
	}

	// Wave B P1.C.3 — ledger dashboard rewired to service-driven path.
	// Q-SDM-DASHBOARD-DOWNSTREAM (LOCKED 2026-05-20): same-commit rewire from
	// `uc.Ledger.Dashboard.Execute` (entity-layer, RETIRED) to
	// `uc.Service.Dashboard.Ledger.GetLedgerDashboard.Execute` (service-
	// layer). The proto type relocated from `domain/ledger/dashboard` to
	// `service/dashboard/ledger`; the view-Response field names stay
	// identical so only the import path on `ledgerdashpb` and the use-case
	// chain change. The ledgerdashboardview Response shape is the view-side
	// contract owned by fycha.
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Ledger != nil && uc.Service.Dashboard.Ledger.GetLedgerDashboard != nil {
		ledgerDash := uc.Service.Dashboard.Ledger.GetLedgerDashboard
		result.GetLedgerDashboardPageData = func(ctx context.Context, req *ledgerdashboardview.Request) (*ledgerdashboardview.Response, error) {
			wsID := req.WorkspaceID
			if wsID == "" {
				wsID = consumer.GetWorkspaceIDFromContext(ctx)
			}
			resp, err := ledgerDash.Execute(ctx, &ledgerdashpb.GetLedgerDashboardRequest{
				WorkspaceId: wsID,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			return &ledgerdashboardview.Response{
				TotalAssets:      resp.GetStats().GetTotalAssets(),
				TotalLiabilities: resp.GetStats().GetTotalLiabilities(),
				TotalEquity:      resp.GetStats().GetTotalEquity(),
				NetIncomeMTD:     resp.GetStats().GetNetIncomeMtd(),
				UnpostedJournals: resp.GetStats().GetUnpostedJournals(),
				BalanceByType:    resp.GetBalanceByType(),
				UnpostedTop:      resp.GetUnpostedTop(),
				RecentEntries:    resp.GetRecentEntries(),
			}, nil
		}
	}

	// Equity dashboard relocated to service.Dashboard.Equity per Wave B P1.C.4
	// (Q-SDM-DASHBOARD-LAYOUT / Q-SDM-DASHBOARD-DOWNSTREAM, 2026-05-20).
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Equity != nil && uc.Service.Dashboard.Equity.GetEquityDashboard != nil {
		equityDash := uc.Service.Dashboard.Equity.GetEquityDashboard
		result.GetEquityDashboardPageData = func(ctx context.Context, req *equitydashboardview.Request) (*equitydashboardview.Response, error) {
			wsID := req.WorkspaceID
			if wsID == "" {
				wsID = consumer.GetWorkspaceIDFromContext(ctx)
			}
			resp, err := equityDash.Execute(ctx, &equitydashpb.GetEquityDashboardRequest{
				WorkspaceId: wsID,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := &equitydashboardview.Response{
				TotalContributed: resp.GetStats().GetTotalContributed(),
				ActiveOwners:     resp.GetStats().GetActiveOwners(),
				DistributionsYTD: resp.GetStats().GetDistributionsYtd(),
				NetMovementYTD:   resp.GetStats().GetNetMovementYtd(),
				ByTypeYTD:        resp.GetByTypeYtd(),
				Recent:           resp.GetRecent(),
			}
			for _, c := range resp.GetTopContributors() {
				out.TopContributors = append(out.TopContributors, equitydashboardview.EquityAccountRow{
					ID:          c.GetId(),
					Name:        c.GetName(),
					OwnerName:   c.GetOwnerName(),
					AccountType: c.GetAccountType(),
					Balance:     c.GetBalance(),
				})
			}
			return out, nil
		}
	}

	// Wave B P1.C.6 — payroll dashboard rewired to service-driven path.
	// Q-SDM-DASHBOARD-DOWNSTREAM (LOCKED 2026-05-20): same-commit rewire from
	// `uc.Payroll.Dashboard.Execute` (entity-layer) to
	// `uc.Service.Dashboard.Payroll.GetPayrollDashboard.Execute` (service-
	// layer). The proto type relocated from `domain/payroll/dashboard` to
	// `service/dashboard/payroll`; the field/getter names stay identical so
	// only the import path on `payrolldashpb` and the use-case chain change.
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Payroll != nil && uc.Service.Dashboard.Payroll.GetPayrollDashboard != nil {
		result.GetPayrollDashboardPageData = func(ctx context.Context, req *payrolldashboardview.Request) (*payrolldashboardview.Response, error) {
			wsID := req.WorkspaceID
			if wsID == "" {
				wsID = consumer.GetWorkspaceIDFromContext(ctx)
			}
			resp, err := uc.Service.Dashboard.Payroll.GetPayrollDashboard.Execute(ctx, &payrolldashpb.GetPayrollDashboardRequest{
				WorkspaceId: wsID,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			return &payrolldashboardview.Response{
				CurrentRunStatus:    resp.GetStats().GetCurrentRunStatus(),
				EmployeesInCurrent:  resp.GetStats().GetEmployeesInCurrent(),
				TotalGrossMTD:       resp.GetStats().GetTotalGrossMtd(),
				RemittancesDue30Cnt: resp.GetStats().GetRemittancesDue30Cnt(),
				LatestRun:           resp.GetLatestRun(),
				RecentRuns:          resp.GetRecentRuns(),
				UpcomingDeadlines:   resp.GetUpcomingDeadlines(),
				GrossTrendLabels:    resp.GetGrossTrendLabels(),
				GrossTrendValues:    resp.GetGrossTrendValues(),
			}, nil
		}
	}

	// Wave B P1.C.5 — Treasury (unified Loan+Cash) rewired to service-driven
	// path. Q-SDM-DASHBOARD-DOWNSTREAM + Q-SDM-DASHBOARD-COUNT: same-commit
	// rewire from `uc.Treasury.LoanDashboard.Execute` (entity-layer, RETIRED)
	// to `uc.Service.Dashboard.Treasury.Loan.GetLoanDashboard.Execute` (service-
	// layer). The proto type relocated from `domain/treasury/dashboard` to
	// `service/dashboard/treasury` (single proto bundles Loan + Cash messages).
	if uc.Service != nil && uc.Service.Dashboard != nil && uc.Service.Dashboard.Treasury != nil && uc.Service.Dashboard.Treasury.Loan != nil && uc.Service.Dashboard.Treasury.Loan.GetLoanDashboard != nil {
		loanDash := uc.Service.Dashboard.Treasury.Loan.GetLoanDashboard
		result.GetLoanDashboardPageData = func(ctx context.Context, req *loansdashboardview.Request) (*loansdashboardview.Response, error) {
			wsID := req.WorkspaceID
			if wsID == "" {
				wsID = consumer.GetWorkspaceIDFromContext(ctx)
			}
			resp, err := loanDash.Execute(ctx, &treasurydashpb.GetLoanDashboardRequest{
				WorkspaceId: wsID,
			})
			if err != nil {
				return nil, err
			}
			if resp == nil {
				return nil, nil
			}
			out := &loansdashboardview.Response{
				TotalOutstanding: resp.GetStats().GetTotalOutstanding(),
				InterestYTD:      resp.GetStats().GetInterestYtd(),
				PaymentsDue30:    resp.GetStats().GetPaymentsDue30(),
				DefaultedCount:   resp.GetStats().GetDefaultedCount(),
				TrendLabels:      resp.GetTrendLabels(),
				TrendValues:      resp.GetTrendValues(),
				RecentPayments:   resp.GetRecentPayments(),
			}
			for _, l := range resp.GetTopLoans() {
				out.TopLoans = append(out.TopLoans, loansdashboardview.LoanRow{
					ID:               l.GetId(),
					LoanNumber:       l.GetLoanNumber(),
					LenderName:       l.GetLenderName(),
					RemainingBalance: l.GetRemainingBalance(),
					PrincipalAmount:  l.GetPrincipalAmount(),
					Status:           l.GetStatus(),
				})
			}
			return out, nil
		}
	}

	return result
}
