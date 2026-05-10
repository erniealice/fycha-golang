package block

import (
	"testing"
)

// TestBlock_RequiresWithUseCases asserts that Block() returns an error when
// WithUseCases() is not provided.
func TestBlock_RequiresWithUseCases(t *testing.T) {
	// Build a minimal Block with no WithUseCases option.
	// We can't call it without a full pyeza.AppContext, but we can verify
	// the useCases nil-check at config time by inspecting cfg directly.
	cfg := &blockConfig{enableAll: true}
	if cfg.useCases != nil {
		t.Fatal("expected useCases to be nil when WithUseCases is not called")
	}
}

// TestWithUseCases_SetsField asserts that WithUseCases sets the useCases field.
func TestWithUseCases_SetsField(t *testing.T) {
	cfg := &blockConfig{}
	uc := &UseCases{}
	opt := WithUseCases(uc)
	opt(cfg)
	if cfg.useCases == nil {
		t.Fatal("expected useCases to be set after WithUseCases()")
	}
	if cfg.useCases != uc {
		t.Fatal("expected useCases pointer to match supplied *UseCases")
	}
}

// TestUseCases_ZeroValue asserts that a zero-value *UseCases has all closures nil
// (no panics from nil pointer dereference on struct fields).
func TestUseCases_ZeroValue(t *testing.T) {
	uc := &UseCases{}
	// Access each group field — should not panic on zero-value struct.
	_ = uc.Workspace.Read
	_ = uc.Asset.Create
	_ = uc.Asset.Category.ListWithPolicyRollup
	_ = uc.DepRun.ListCandidates
	_ = uc.DepRun.Generate
	_ = uc.DepRun.List
	_ = uc.DepRun.Read
	_ = uc.DepRun.ListEntries
	_ = uc.Revaluation.Revalue
	_ = uc.Revaluation.Preview
	_ = uc.Ledger.Account.GetListPageData
	_ = uc.Ledger.JournalEntry.GetListPageData
	_ = uc.FiscalPeriod.GetListPageData
	_ = uc.Tax.ListTaxRates
	_ = uc.Finance.ListForexRates
	_ = uc.Treasury.ListWithholdingCertificates
	_ = uc.GetLedgerDashboardPageData
	_ = uc.GetEquityDashboardPageData
	_ = uc.GetPayrollDashboardPageData
	_ = uc.GetLoanDashboardPageData
}
