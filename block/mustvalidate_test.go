package block

import (
	"context"
	"strings"
	"testing"
	"time"

	apagingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ap_aging"
	aragingpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/ar_aging"
	dspb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/domain_specific"
	gcfpb "github.com/erniealice/esqyma/pkg/schema/v1/service/reporting/gross_cashflow"
)

// ---------------------------------------------------------------------------
// MustValidate — FAIL-CLOSED wiring guard (architecture-roast burn #1).
//
// RequireFor returns an error; MustValidate adds the posture: in dev/test
// (testing.Testing() is true here) a missing REQUIRED closure PANICS — loud,
// stack-traced, uncatchable-by-accident — so a nil-closure wiring gap can never
// be silently dropped into an empty-state render. OPTIONAL nils never trip it.
//
// fycha had NO RequireFor before this wave; the skeleton authored alongside this
// test makes each report's report/list/get closure, the asset/ledger CRUD set,
// and the read-only list modules' single closure REQUIRED, while the four
// dashboard closures + the dashboard-only / TODO-stub modules stay OPTIONAL.
// ---------------------------------------------------------------------------

// wireReportsRequired sets every closure RequireFor checks for the Reports
// module: the AR/AP aging, gross-profit, and domain-specific report/list closures.
func wireReportsRequired(uc *UseCases) {
	ar := &uc.Reports.ARAging
	ar.GetReceivablesAgingReport = func(context.Context, *aragingpb.GetReceivablesAgingRequest) (*aragingpb.GetReceivablesAgingResponse, error) {
		return nil, nil
	}
	ar.GetCollectionSummaryReport = func(context.Context, *aragingpb.GetCollectionSummaryRequest) (*aragingpb.GetCollectionSummaryResponse, error) {
		return nil, nil
	}
	uc.Reports.APAging.GetPayablesAgingReport = func(context.Context, *apagingpb.GetPayablesAgingRequest) (*apagingpb.GetPayablesAgingResponse, error) {
		return nil, nil
	}
	uc.Reports.GrossCashFlow.GetGrossProfitReport = func(context.Context, *gcfpb.GetGrossProfitRequest) (*gcfpb.GetGrossProfitResponse, error) {
		return nil, nil
	}
	ds := &uc.Reports.DomainSpecific
	ds.GetRevenueReport = func(context.Context, *dspb.GetRevenueReportRequest) (*dspb.GetRevenueReportResponse, error) {
		return nil, nil
	}
	ds.GetExpenditureReport = func(context.Context, *dspb.GetExpenditureReportRequest) (*dspb.GetExpenditureReportResponse, error) {
		return nil, nil
	}
	ds.GetDisbursementReport = func(context.Context, *dspb.GetDisbursementReportRequest) (*dspb.GetDisbursementReportResponse, error) {
		return nil, nil
	}
	ds.ListRevenue = func(context.Context, *time.Time, *time.Time) ([]map[string]any, error) { return nil, nil }
	ds.ListExpenses = func(context.Context, *time.Time, *time.Time) ([]map[string]any, error) { return nil, nil }
}

// TestMustValidate_NilRequiredClosure_Panics is the core burn-#1 proof: with the
// Reports module enabled but one REQUIRED closure (ListExpenses) left nil,
// MustValidate must PANIC under test — not return an empty render, not silently
// degrade. This is the loud failure the prior bare nil-check path lacked.
func TestMustValidate_NilRequiredClosure_Panics(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	wireReportsRequired(uc)
	uc.Reports.DomainSpecific.ListExpenses = nil // drop exactly one REQUIRED closure

	defer func() {
		r := recover()
		if r == nil {
			t.Fatal("MustValidate(Reports enabled, ListExpenses nil) should PANIC in dev/test, but did not")
		}
		msg, _ := r.(string)
		if !strings.Contains(msg, "ListExpenses") {
			t.Fatalf("panic message should name the missing field; got %q", msg)
		}
	}()

	// Should not reach the next line — MustValidate panics first.
	_ = uc.MustValidate(&blockConfig{reports: true})
	t.Fatal("MustValidate returned instead of panicking on a nil REQUIRED closure")
}

// TestMustValidate_EmptyUseCases_EnableAll_Panics: a fully empty UseCases with
// every module enabled (the "permanently nil dashboard" trap) must panic loudly
// in dev/test rather than register a wall of empty views.
func TestMustValidate_EmptyUseCases_EnableAll_Panics(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	defer func() {
		if recover() == nil {
			t.Fatal("MustValidate(empty UseCases, enableAll) should PANIC in dev/test")
		}
	}()
	_ = uc.MustValidate(&blockConfig{enableAll: true})
	t.Fatal("MustValidate returned instead of panicking on an empty enableAll wiring")
}

// TestMustValidate_NilOptionalClosure_OK proves the required-vs-optional
// discrimination survives the fail-closed wrapper: the dashboard-only OPTIONAL
// modules (Loans/Equity/Payroll — not in RequireFor) with their dashboard
// closures left nil must pass MustValidate with NO panic and NO error.
func TestMustValidate_NilOptionalClosure_OK(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	// Dashboard-only modules enabled, their (optional) dashboard closures nil.
	// Also leave Workspace.Read / Revaluation / FiscalPeriod mutators nil.
	cfg := &blockConfig{loans: true, equity: true, payroll: true}

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustValidate(optional nil closures) must NOT panic; panicked with %v", r)
		}
	}()
	if err := uc.MustValidate(cfg); err != nil {
		t.Fatalf("MustValidate(optional nil closures) should be nil, got %v", err)
	}
}

// TestMustValidate_FullyWired_OK: a completely wired REQUIRED set passes with no
// panic and no error (happy path — guard is silent when wiring is complete).
func TestMustValidate_FullyWired_OK(t *testing.T) {
	t.Parallel()

	uc := &UseCases{}
	wireReportsRequired(uc)

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("MustValidate(fully wired Reports) must NOT panic; panicked with %v", r)
		}
	}()
	if err := uc.MustValidate(&blockConfig{reports: true}); err != nil {
		t.Fatalf("MustValidate(fully wired Reports) should be nil, got %v", err)
	}
}
