package block

// dashboard_wiring.go wires dashboard use cases from the espyna UseCases
// aggregate into fycha module ModuleDeps callbacks.
//
// Since the dashboard use-case request/response types live in espyna's
// internal packages (unreachable from fycha), we use reflection to:
//  1. Dereference the use-case pointer field (e.g. useCases.Ledger.Dashboard)
//  2. Build the Execute request via reflect.New + field-name assignment
//  3. Call Execute reflectively
//  4. Copy matching fields from the response to the view-layer Response type
//
// All helpers are nil-safe: if the Dashboard field is nil the callback is
// left unset and the dashboard view renders empty state (its existing behaviour).

import (
	"context"
	"reflect"
	"time"

	consumer "github.com/erniealice/espyna-golang/consumer"

	ledgerdashboard "github.com/erniealice/fycha-golang/views/ledger/dashboard"
	equitydashboard "github.com/erniealice/fycha-golang/views/equity/dashboard"
	payrolldashboard "github.com/erniealice/fycha-golang/views/payroll/dashboard"
	loansdashboard "github.com/erniealice/fycha-golang/views/loans/dashboard"

	ledgermod "github.com/erniealice/fycha-golang/views/ledger"
	equitymod "github.com/erniealice/fycha-golang/views/equity"
	payrollmod "github.com/erniealice/fycha-golang/views/payroll"
	loansmod "github.com/erniealice/fycha-golang/views/loans"

	equitytransactionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/equity_transaction"
	journalentrypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
	loanpaymentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/loan_payment"
	payrollremittancepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_remittance"
	payrollrunpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/payroll/payroll_run"
)

// callExecute calls a use-case's Execute method via reflection.
// The useCase value must be a non-nil pointer to a use-case struct.
// workspaceID and now are set on the request struct by field name.
// Returns the response (as reflect.Value) and an error.
func callExecute(useCase reflect.Value, ctx context.Context, workspaceID string, now time.Time) (reflect.Value, error) {
	m := useCase.MethodByName("Execute")
	if !m.IsValid() {
		return reflect.Value{}, nil
	}
	// Execute signature: func(ctx context.Context, req *Request) (*Response, error)
	// Build request: reflect.New gives us a *Request
	reqType := m.Type().In(1).Elem() // *Request → Request
	reqPtr := reflect.New(reqType)
	if f := reqPtr.Elem().FieldByName("WorkspaceID"); f.IsValid() && f.CanSet() {
		f.SetString(workspaceID)
	}
	if f := reqPtr.Elem().FieldByName("Now"); f.IsValid() && f.CanSet() {
		f.Set(reflect.ValueOf(now))
	}
	results := m.Call([]reflect.Value{reflect.ValueOf(ctx), reqPtr})
	if len(results) < 2 {
		return reflect.Value{}, nil
	}
	if !results[1].IsNil() {
		return reflect.Value{}, results[1].Interface().(error)
	}
	resp := results[0]
	if resp.Kind() == reflect.Ptr && !resp.IsNil() {
		return resp.Elem(), nil
	}
	return resp, nil
}

// int64Field reads an int64 field by name from a reflect.Value (struct).
func int64Field(v reflect.Value, name string) int64 {
	if !v.IsValid() {
		return 0
	}
	f := v.FieldByName(name)
	if !f.IsValid() {
		return 0
	}
	return f.Int()
}

// mapStringInt64Field reads a map[string]int64 field by name.
func mapStringInt64Field(v reflect.Value, name string) map[string]int64 {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	return f.Interface().(map[string]int64)
}

// protoSliceField reads a proto slice field by name and type-asserts it.
func protoSliceField[T any](v reflect.Value, name string) []T {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	if s, ok := f.Interface().([]T); ok {
		return s
	}
	return nil
}

// float64SliceField reads a []float64 field by name.
func float64SliceField(v reflect.Value, name string) []float64 {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	if s, ok := f.Interface().([]float64); ok {
		return s
	}
	return nil
}

// stringSliceField reads a []string field by name.
func stringSliceField(v reflect.Value, name string) []string {
	if !v.IsValid() {
		return nil
	}
	f := v.FieldByName(name)
	if !f.IsValid() || f.IsNil() {
		return nil
	}
	if s, ok := f.Interface().([]string); ok {
		return s
	}
	return nil
}

// ---------------------------------------------------------------------------
// Ledger dashboard wiring
// ---------------------------------------------------------------------------

// wireLedgerDashboard sets ledgerDeps.GetLedgerDashboardPageData if
// useCases.Ledger.Dashboard is non-nil.
func wireLedgerDashboard(deps *ledgermod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Ledger == nil || useCases.Ledger.Dashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Ledger.Dashboard)
	deps.GetLedgerDashboardPageData = func(ctx context.Context, req *ledgerdashboard.Request) (*ledgerdashboard.Response, error) {
		workspaceID := ""
		if req != nil {
			workspaceID = req.WorkspaceID
		}
		resp, err := callExecute(uc, ctx, workspaceID, time.Now())
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		return &ledgerdashboard.Response{
			TotalAssets:      int64Field(resp, "TotalAssets"),
			TotalLiabilities: int64Field(resp, "TotalLiabilities"),
			TotalEquity:      int64Field(resp, "TotalEquity"),
			NetIncomeMTD:     int64Field(resp, "NetIncomeMTD"),
			UnpostedJournals: int64Field(resp, "UnpostedJournals"),
			BalanceByType:    mapStringInt64Field(resp, "BalanceByType"),
			UnpostedTop:      protoSliceField[*journalentrypb.JournalEntry](resp, "UnpostedTop"),
			RecentEntries:    protoSliceField[*journalentrypb.JournalEntry](resp, "RecentEntries"),
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Equity dashboard wiring
// ---------------------------------------------------------------------------

// wireEquityDashboard sets equityDeps.GetEquityDashboardPageData if
// useCases.Ledger.EquityDashboard is non-nil.
func wireEquityDashboard(deps *equitymod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Ledger == nil || useCases.Ledger.EquityDashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Ledger.EquityDashboard)
	deps.GetEquityDashboardPageData = func(ctx context.Context, req *equitydashboard.Request) (*equitydashboard.Response, error) {
		workspaceID := ""
		if req != nil {
			workspaceID = req.WorkspaceID
		}
		resp, err := callExecute(uc, ctx, workspaceID, time.Now())
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Copy TopContributors: []EquityAccountSlice → []EquityAccountRow
		var topContribs []equitydashboard.EquityAccountRow
		if f := resp.FieldByName("TopContributors"); f.IsValid() && !f.IsNil() {
			for i := 0; i < f.Len(); i++ {
				s := f.Index(i)
				topContribs = append(topContribs, equitydashboard.EquityAccountRow{
					ID:          s.FieldByName("ID").String(),
					Name:        s.FieldByName("Name").String(),
					OwnerName:   s.FieldByName("OwnerName").String(),
					AccountType: s.FieldByName("AccountType").String(),
					Balance:     s.FieldByName("Balance").Int(),
				})
			}
		}
		return &equitydashboard.Response{
			TotalContributed: int64Field(resp, "TotalContributed"),
			ActiveOwners:     int64Field(resp, "ActiveOwners"),
			DistributionsYTD: int64Field(resp, "DistributionsYTD"),
			NetMovementYTD:   int64Field(resp, "NetMovementYTD"),
			ByTypeYTD:        mapStringInt64Field(resp, "ByTypeYTD"),
			TopContributors:  topContribs,
			Recent:           protoSliceField[*equitytransactionpb.EquityTransaction](resp, "Recent"),
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Payroll dashboard wiring
// ---------------------------------------------------------------------------

// wirePayrollDashboard sets payrollDeps.GetPayrollDashboardPageData if
// useCases.Payroll.Dashboard is non-nil.
func wirePayrollDashboard(deps *payrollmod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Payroll == nil || useCases.Payroll.Dashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Payroll.Dashboard)
	deps.GetPayrollDashboardPageData = func(ctx context.Context, req *payrolldashboard.Request) (*payrolldashboard.Response, error) {
		workspaceID := ""
		if req != nil {
			workspaceID = req.WorkspaceID
		}
		resp, err := callExecute(uc, ctx, workspaceID, time.Now())
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var status string
		var employees int32
		var totalGross, remittances int64
		if stats := resp.FieldByName("Stats"); stats.IsValid() {
			status = stats.FieldByName("CurrentRunStatus").String()
			employees = int32(stats.FieldByName("EmployeesInCurrent").Int())
			totalGross = stats.FieldByName("TotalGrossMTD").Int()
			remittances = stats.FieldByName("RemittancesDue30Cnt").Int()
		}
		// Extract LatestRun single pointer.
		var latestRunProto *payrollrunpb.PayrollRun
		if f := resp.FieldByName("LatestRun"); f.IsValid() && !f.IsNil() {
			if v, ok := f.Interface().(*payrollrunpb.PayrollRun); ok {
				latestRunProto = v
			}
		}
		return &payrolldashboard.Response{
			CurrentRunStatus:    status,
			EmployeesInCurrent:  employees,
			TotalGrossMTD:       totalGross,
			RemittancesDue30Cnt: remittances,
			LatestRun:           latestRunProto,
			RecentRuns:          protoSliceField[*payrollrunpb.PayrollRun](resp, "RecentRuns"),
			UpcomingDeadlines:   protoSliceField[*payrollremittancepb.PayrollRemittance](resp, "UpcomingDeadlines"),
			GrossTrendLabels:    stringSliceField(resp, "GrossTrendLabels"),
			GrossTrendValues:    float64SliceField(resp, "GrossTrendValues"),
		}, nil
	}
}

// ---------------------------------------------------------------------------
// Loans dashboard wiring
// ---------------------------------------------------------------------------

// wireLoansDashboard sets loansDeps.GetLoanDashboardPageData if
// useCases.Treasury.LoanDashboard is non-nil.
func wireLoansDashboard(deps *loansmod.ModuleDeps, useCases *consumer.UseCases) {
	if useCases == nil || useCases.Treasury == nil || useCases.Treasury.LoanDashboard == nil {
		return
	}
	uc := reflect.ValueOf(useCases.Treasury.LoanDashboard)
	deps.GetLoanDashboardPageData = func(ctx context.Context, req *loansdashboard.Request) (*loansdashboard.Response, error) {
		workspaceID := ""
		if req != nil {
			workspaceID = req.WorkspaceID
		}
		resp, err := callExecute(uc, ctx, workspaceID, time.Now())
		if err != nil || !resp.IsValid() {
			return nil, err
		}
		// Stats sub-struct
		var totalOut, interestYTD, paymentsDue30, defaultedCount int64
		if stats := resp.FieldByName("Stats"); stats.IsValid() {
			totalOut = stats.FieldByName("TotalOutstanding").Int()
			interestYTD = stats.FieldByName("InterestYTD").Int()
			paymentsDue30 = stats.FieldByName("PaymentsDue30").Int()
			defaultedCount = stats.FieldByName("DefaultedCount").Int()
		}
		// TopLoans: []LoanSlice → []LoanRow
		var topLoans []loansdashboard.LoanRow
		if f := resp.FieldByName("TopLoans"); f.IsValid() && !f.IsNil() {
			for i := 0; i < f.Len(); i++ {
				s := f.Index(i)
				topLoans = append(topLoans, loansdashboard.LoanRow{
					ID:               s.FieldByName("ID").String(),
					LoanNumber:       s.FieldByName("LoanNumber").String(),
					LenderName:       s.FieldByName("LenderName").String(),
					RemainingBalance: s.FieldByName("RemainingBalance").Int(),
					PrincipalAmount:  s.FieldByName("PrincipalAmount").Int(),
					Status:           s.FieldByName("Status").String(),
				})
			}
		}
		return &loansdashboard.Response{
			TotalOutstanding: totalOut,
			InterestYTD:      interestYTD,
			PaymentsDue30:    paymentsDue30,
			DefaultedCount:   defaultedCount,
			TrendLabels:      stringSliceField(resp, "TrendLabels"),
			TrendValues:      float64SliceField(resp, "TrendValues"),
			TopLoans:         topLoans,
			RecentPayments:   protoSliceField[*loanpaymentpb.LoanPayment](resp, "RecentPayments"),
		}, nil
	}
}
