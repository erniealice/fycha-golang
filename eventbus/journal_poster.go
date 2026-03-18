package eventbus

// journal_poster.go — Phase 9 auto-posting skeleton.
//
// JournalPoster registers handlers for all known operational event types and
// maps each event to the correct double-entry accounting rule. In Phase 9 all
// handlers are stubs (they log and return nil). Full implementation wires the
// CreateJournalEntry use case once the JournalEntry repository is available in
// the operational context.
//
// Accounting rules documented per event type:
//
//   "revenue.completed"
//     DR  Accounts Receivable (or Cash if collected immediately)
//     CR  Revenue
//
//   "collection.received"
//     DR  Cash / Bank
//     CR  Accounts Receivable
//
//   "expenditure.approved"
//     DR  Expense Account (by category)
//     CR  Accounts Payable (or Cash if paid immediately)
//
//   "disbursement.paid"
//     DR  Accounts Payable
//     CR  Cash / Bank
//
//   "asset.acquired"
//     DR  Fixed Asset Account
//     CR  Cash / Accounts Payable
//
//   "asset.deprecated"
//     DR  Depreciation Expense
//     CR  Accumulated Depreciation
//
//   "prepayment.created"
//     DR  Prepaid Expense (asset)
//     CR  Cash / Accounts Payable
//
//   "prepayment.amortized"
//     DR  Expense Account
//     CR  Prepaid Expense (asset)
//
//   "loan.received"
//     DR  Cash / Bank
//     CR  Loan Payable (liability)
//
//   "loan.payment"
//     DR  Loan Payable (principal portion)
//     DR  Interest Expense (interest portion)
//     CR  Cash / Bank
//
//   "equity.contribution"
//     DR  Cash / Bank
//     CR  Owner's Capital / Paid-In Capital
//
//   "equity.withdrawal"
//     DR  Owner's Draw / Drawings
//     CR  Cash / Bank
//
//   "payroll.posted"
//     DR  Salary Expense (gross pay)
//     CR  Cash / Bank (net pay)
//     CR  Government Payables (SSS, PhilHealth, Pag-IBIG, withholding tax)
//
//   "petty_cash.replenished"
//     DR  Petty Cash Expense accounts (per voucher categories)
//     CR  Cash / Bank (reimbursement amount)

import (
	"context"
	"log"
)

// EventTypeRevenuCompleted fires when a revenue record is marked completed/collected.
const EventTypeRevenueCompleted = "revenue.completed"

// EventTypeCollectionReceived fires when a cash collection is recorded.
const EventTypeCollectionReceived = "collection.received"

// EventTypeExpenditureApproved fires when an expenditure is approved for payment.
const EventTypeExpenditureApproved = "expenditure.approved"

// EventTypeDisbursementPaid fires when a disbursement is marked as paid.
const EventTypeDisbursementPaid = "disbursement.paid"

// EventTypeAssetAcquired fires when a fixed asset acquisition is recorded.
const EventTypeAssetAcquired = "asset.acquired"

// EventTypeAssetDepreciated fires when depreciation is run on an asset.
const EventTypeAssetDepreciated = "asset.deprecated"

// EventTypePrepaymentCreated fires when a prepaid expense is recorded.
const EventTypePrepaymentCreated = "prepayment.created"

// EventTypePrepaymentAmortized fires when a prepayment amortization period is posted.
const EventTypePrepaymentAmortized = "prepayment.amortized"

// EventTypeLoanReceived fires when a loan drawdown is recorded.
const EventTypeLoanReceived = "loan.received"

// EventTypeLoanPayment fires when a loan repayment is made.
const EventTypeLoanPayment = "loan.payment"

// EventTypeEquityContribution fires when an equity capital injection is recorded.
const EventTypeEquityContribution = "equity.contribution"

// EventTypeEquityWithdrawal fires when an equity draw/withdrawal is recorded.
const EventTypeEquityWithdrawal = "equity.withdrawal"

// EventTypePayrollPosted fires when a payroll run is posted.
const EventTypePayrollPosted = "payroll.posted"

// EventTypePettyCashReplenished fires when a petty cash fund is replenished.
const EventTypePettyCashReplenished = "petty_cash.replenished"

// JournalPoster handles all accounting event types and auto-posts journal entries.
//
// In Phase 9 all handlers are stubs. Implement CreateJournalEntry in each
// handler once the use case is injected via JournalPosterDeps.
type JournalPoster struct {
	// TODO Phase 9 full implementation: inject CreateJournalEntry use case.
	// CreateJournalEntry func(ctx context.Context, req *journalentrypb.CreateJournalEntryRequest) (*journalentrypb.CreateJournalEntryResponse, error)
}

// NewJournalPoster creates a JournalPoster.
// In Phase 9 it carries no dependencies; wiring the use case is future work.
func NewJournalPoster() *JournalPoster {
	return &JournalPoster{}
}

// RegisterAll subscribes all 14 handler stubs to the provided EventBus.
func (p *JournalPoster) RegisterAll(bus EventBus) {
	bus.Subscribe(EventTypeRevenueCompleted, p.handleRevenueCompleted)
	bus.Subscribe(EventTypeCollectionReceived, p.handleCollectionReceived)
	bus.Subscribe(EventTypeExpenditureApproved, p.handleExpenditureApproved)
	bus.Subscribe(EventTypeDisbursementPaid, p.handleDisbursementPaid)
	bus.Subscribe(EventTypeAssetAcquired, p.handleAssetAcquired)
	bus.Subscribe(EventTypeAssetDepreciated, p.handleAssetDepreciated)
	bus.Subscribe(EventTypePrepaymentCreated, p.handlePrepaymentCreated)
	bus.Subscribe(EventTypePrepaymentAmortized, p.handlePrepaymentAmortized)
	bus.Subscribe(EventTypeLoanReceived, p.handleLoanReceived)
	bus.Subscribe(EventTypeLoanPayment, p.handleLoanPayment)
	bus.Subscribe(EventTypeEquityContribution, p.handleEquityContribution)
	bus.Subscribe(EventTypeEquityWithdrawal, p.handleEquityWithdrawal)
	bus.Subscribe(EventTypePayrollPosted, p.handlePayrollPosted)
	bus.Subscribe(EventTypePettyCashReplenished, p.handlePettyCashReplenished)
}

// handleRevenueCompleted auto-posts:
//   DR  Accounts Receivable (or Cash if payload["collected"] == true)
//   CR  Revenue
//
// Expected payload keys:
//   "amount"         float64   — gross revenue amount
//   "revenue_id"     string    — source revenue ID (for journal source link)
//   "collected"      bool      — true if cash was collected immediately
//   "account_ar"     string    — account ID for Accounts Receivable
//   "account_cash"   string    — account ID for Cash (used when collected == true)
//   "account_rev"    string    — account ID for Revenue
func (p *JournalPoster) handleRevenueCompleted(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: revenue.completed source=%s — DR AR/Cash, CR Revenue", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_ar (or account_cash if collected), amount
	//   line 2: credit account_rev, amount
	// then call p.CreateJournalEntry(ctx, req)
	return nil
}

// handleCollectionReceived auto-posts:
//   DR  Cash / Bank
//   CR  Accounts Receivable
//
// Expected payload keys:
//   "amount"         float64   — collection amount
//   "collection_id"  string    — source collection ID
//   "account_cash"   string    — account ID for Cash / Bank
//   "account_ar"     string    — account ID for Accounts Receivable
func (p *JournalPoster) handleCollectionReceived(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: collection.received source=%s — DR Cash, CR AR", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_cash, amount
	//   line 2: credit account_ar, amount
	return nil
}

// handleExpenditureApproved auto-posts:
//   DR  Expense Account (by category/type)
//   CR  Accounts Payable (or Cash if paid immediately)
//
// Expected payload keys:
//   "amount"             float64   — expenditure amount
//   "expenditure_id"     string    — source expenditure ID
//   "paid_immediately"   bool      — true if cash was disbursed immediately
//   "account_expense"    string    — account ID for the expense category
//   "account_ap"         string    — account ID for Accounts Payable
//   "account_cash"       string    — account ID for Cash (when paid_immediately)
func (p *JournalPoster) handleExpenditureApproved(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: expenditure.approved source=%s — DR Expense, CR AP/Cash", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_expense, amount
	//   line 2: credit account_ap (or account_cash if paid_immediately), amount
	return nil
}

// handleDisbursementPaid auto-posts:
//   DR  Accounts Payable
//   CR  Cash / Bank
//
// Expected payload keys:
//   "amount"            float64   — disbursement amount
//   "disbursement_id"   string    — source disbursement ID
//   "account_ap"        string    — account ID for Accounts Payable
//   "account_cash"      string    — account ID for Cash / Bank
func (p *JournalPoster) handleDisbursementPaid(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: disbursement.paid source=%s — DR AP, CR Cash", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_ap, amount
	//   line 2: credit account_cash, amount
	return nil
}

// handleAssetAcquired auto-posts:
//   DR  Fixed Asset Account
//   CR  Cash / Accounts Payable
//
// Expected payload keys:
//   "amount"          float64   — acquisition cost
//   "asset_id"        string    — source asset ID
//   "paid_in_cash"    bool      — true if paid in cash, false if financed (AP)
//   "account_asset"   string    — account ID for the Fixed Asset
//   "account_cash"    string    — account ID for Cash (when paid_in_cash)
//   "account_ap"      string    — account ID for Accounts Payable (when financed)
func (p *JournalPoster) handleAssetAcquired(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: asset.acquired source=%s — DR Fixed Asset, CR Cash/AP", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_asset, amount
	//   line 2: credit account_cash (or account_ap if financed), amount
	return nil
}

// handleAssetDepreciated auto-posts:
//   DR  Depreciation Expense
//   CR  Accumulated Depreciation
//
// Expected payload keys:
//   "amount"                    float64   — depreciation amount for the period
//   "asset_id"                  string    — source asset ID
//   "account_depreciation_exp"  string    — account ID for Depreciation Expense
//   "account_accum_dep"         string    — account ID for Accumulated Depreciation
func (p *JournalPoster) handleAssetDepreciated(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: asset.deprecated source=%s — DR Depreciation Exp, CR Accum Dep", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_depreciation_exp, amount
	//   line 2: credit account_accum_dep, amount
	return nil
}

// handlePrepaymentCreated auto-posts:
//   DR  Prepaid Expense (asset)
//   CR  Cash / Accounts Payable
//
// Expected payload keys:
//   "amount"           float64   — prepayment amount
//   "prepayment_id"    string    — source prepayment ID
//   "paid_in_cash"     bool      — true if paid in cash, false if on account
//   "account_prepaid"  string    — account ID for Prepaid Expense (asset)
//   "account_cash"     string    — account ID for Cash
//   "account_ap"       string    — account ID for Accounts Payable
func (p *JournalPoster) handlePrepaymentCreated(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: prepayment.created source=%s — DR Prepaid Expense, CR Cash/AP", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_prepaid, amount
	//   line 2: credit account_cash (or account_ap), amount
	return nil
}

// handlePrepaymentAmortized auto-posts:
//   DR  Expense Account (matching prepayment category)
//   CR  Prepaid Expense (asset)
//
// Expected payload keys:
//   "amount"           float64   — amortization amount for the period
//   "prepayment_id"    string    — source prepayment ID
//   "account_expense"  string    — account ID for the appropriate Expense
//   "account_prepaid"  string    — account ID for Prepaid Expense (asset)
func (p *JournalPoster) handlePrepaymentAmortized(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: prepayment.amortized source=%s — DR Expense, CR Prepaid Expense", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_expense, amount
	//   line 2: credit account_prepaid, amount
	return nil
}

// handleLoanReceived auto-posts:
//   DR  Cash / Bank
//   CR  Loan Payable (liability)
//
// Expected payload keys:
//   "amount"           float64   — loan principal received
//   "loan_id"          string    — source loan ID
//   "account_cash"     string    — account ID for Cash / Bank
//   "account_loan"     string    — account ID for Loan Payable
func (p *JournalPoster) handleLoanReceived(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: loan.received source=%s — DR Cash, CR Loan Payable", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_cash, amount
	//   line 2: credit account_loan, amount
	return nil
}

// handleLoanPayment auto-posts:
//   DR  Loan Payable (principal portion)
//   DR  Interest Expense (interest portion)
//   CR  Cash / Bank (total payment)
//
// Expected payload keys:
//   "total_amount"        float64   — total payment made
//   "principal_amount"    float64   — principal portion of the payment
//   "interest_amount"     float64   — interest portion of the payment
//   "loan_payment_id"     string    — source loan payment ID
//   "account_loan"        string    — account ID for Loan Payable
//   "account_interest"    string    — account ID for Interest Expense
//   "account_cash"        string    — account ID for Cash / Bank
func (p *JournalPoster) handleLoanPayment(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: loan.payment source=%s — DR Loan Payable + DR Interest Exp, CR Cash", event.SourceID)
	// TODO: build CreateJournalEntryRequest with three lines:
	//   line 1: debit account_loan, principal_amount
	//   line 2: debit account_interest, interest_amount
	//   line 3: credit account_cash, total_amount
	return nil
}

// handleEquityContribution auto-posts:
//   DR  Cash / Bank
//   CR  Owner's Capital / Paid-In Capital
//
// Expected payload keys:
//   "amount"                float64   — equity contribution amount
//   "equity_transaction_id" string    — source equity transaction ID
//   "account_cash"          string    — account ID for Cash / Bank
//   "account_capital"       string    — account ID for Owner's Capital
func (p *JournalPoster) handleEquityContribution(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: equity.contribution source=%s — DR Cash, CR Owner's Capital", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_cash, amount
	//   line 2: credit account_capital, amount
	return nil
}

// handleEquityWithdrawal auto-posts:
//   DR  Owner's Draw / Drawings
//   CR  Cash / Bank
//
// Expected payload keys:
//   "amount"                float64   — withdrawal amount
//   "equity_transaction_id" string    — source equity transaction ID
//   "account_drawings"      string    — account ID for Owner's Draw/Drawings
//   "account_cash"          string    — account ID for Cash / Bank
func (p *JournalPoster) handleEquityWithdrawal(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: equity.withdrawal source=%s — DR Owner's Draw, CR Cash", event.SourceID)
	// TODO: build CreateJournalEntryRequest with two lines:
	//   line 1: debit account_drawings, amount
	//   line 2: credit account_cash, amount
	return nil
}

// handlePayrollPosted auto-posts:
//   DR  Salary Expense (gross pay)
//   CR  Cash / Bank (net pay disbursed)
//   CR  Government Payables (SSS, PhilHealth, Pag-IBIG, withholding tax)
//
// Expected payload keys:
//   "gross_pay"             float64   — total gross payroll amount
//   "net_pay"               float64   — net pay disbursed to employees
//   "gov_contributions"     float64   — total government contribution withholdings
//   "payroll_run_id"        string    — source payroll run ID
//   "account_salary_exp"    string    — account ID for Salary Expense
//   "account_cash"          string    — account ID for Cash / Bank
//   "account_gov_payables"  string    — account ID for Government Payables
func (p *JournalPoster) handlePayrollPosted(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: payroll.posted source=%s — DR Salary Exp, CR Cash + CR Gov Payables", event.SourceID)
	// TODO: build CreateJournalEntryRequest with three lines:
	//   line 1: debit account_salary_exp, gross_pay
	//   line 2: credit account_cash, net_pay
	//   line 3: credit account_gov_payables, gov_contributions
	// Note: gross_pay == net_pay + gov_contributions (entries must balance)
	return nil
}

// handlePettyCashReplenished auto-posts:
//   DR  Petty Cash Expense accounts (per voucher categories)
//   CR  Cash / Bank (replenishment amount)
//
// Expected payload keys:
//   "replenishment_amount"   float64          — total replenishment disbursed
//   "petty_cash_fund_id"     string           — source petty cash fund ID
//   "expense_lines"          []map[string]any — per-category expense breakdown:
//                              "account_id" string  — expense account ID
//                              "amount"     float64 — amount for this category
//   "account_cash"           string           — account ID for Cash / Bank
func (p *JournalPoster) handlePettyCashReplenished(_ context.Context, event Event) error {
	log.Printf("[eventbus] stub: petty_cash.replenished source=%s — DR Expense accounts, CR Cash", event.SourceID)
	// TODO: build CreateJournalEntryRequest with N+1 lines:
	//   for each entry in expense_lines:
	//     line N: debit account_id, amount
	//   last line: credit account_cash, replenishment_amount
	// Note: sum of expense_lines amounts must equal replenishment_amount
	return nil
}
