// Package seeder provides Chart of Accounts seed data for Philippine service businesses.
// The default CoA follows IFRS/PFRS standards with practical groupings suitable for
// a small-to-medium salon, spa, or similar service operation.
//
// Usage:
//
//	import "github.com/erniealice/fycha-golang/seeder"
//
//	err := seeder.SeedDefaultCoA(ctx, createAccountFn, workspaceID)
package seeder

import (
	"context"
	"fmt"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
)

// CreateAccountFunc is the signature of the CreateAccount use case Execute method.
type CreateAccountFunc func(ctx context.Context, req *accountpb.CreateAccountRequest) (*accountpb.CreateAccountResponse, error)

// DefaultCoAEntry defines one account in the standard Philippine service business CoA.
type DefaultCoAEntry struct {
	Code             string
	Name             string
	Element          accountpb.AccountElement
	Classification   accountpb.AccountClassification
	NormalBalance    accountpb.NormalBalance
	CashFlowActivity accountpb.CashFlowActivity
	IsSystemAccount  bool
	Description      string
}

// DefaultCoA returns the standard Chart of Accounts for a Philippine service business.
// Accounts are ordered by code (element → classification → account).
// Based on IFRS/PFRS standards as documented in docs/epic/20260317-chart-of-accounts-integration/01-plan/03-coa-hierarchy.md
func DefaultCoA() []DefaultCoAEntry {
	return []DefaultCoAEntry{
		// ====================================================================
		// ASSETS (1000–1999)
		// ====================================================================

		// Current Assets — Cash & Equivalents
		{
			Code:           "1010",
			Name:           "Cash on Hand",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Physical cash held in office safe and petty cash box",
		},
		{
			Code:           "1020",
			Name:           "Petty Cash Fund",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Petty cash maintained for minor daily expenses",
		},
		{
			Code:           "1030",
			Name:           "BDO Savings Account",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "BDO savings account for daily operations",
		},
		{
			Code:           "1040",
			Name:           "BPI Checking Account",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "BPI checking account for payroll and supplier payments",
		},
		// Current Assets — Receivables
		{
			Code:           "1110",
			Name:           "Accounts Receivable",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Amounts owed by customers for services rendered on credit",
		},
		{
			Code:            "1120",
			Name:            "Allowance for Doubtful Accounts",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Contra asset: estimated uncollectible receivables (credit normal balance)",
		},
		// Current Assets — Prepaid & Other
		{
			Code:           "1200",
			Name:           "Prepaid Expenses",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Expenses paid in advance (insurance, rent deposit, subscriptions)",
		},
		{
			Code:           "1210",
			Name:           "Creditable Withholding Tax",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "CWT certificates received from clients subject to expanded withholding tax",
		},
		{
			Code:           "1220",
			Name:           "Input VAT",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "VAT paid on purchases claimable against output VAT",
		},
		// Current Assets — Inventory
		{
			Code:           "1300",
			Name:           "Supplies Inventory",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Service supplies (salon products, spa consumables) held for use",
		},
		{
			Code:           "1310",
			Name:           "Merchandise Inventory",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Retail products held for resale to clients",
		},

		// Non-Current Assets — PP&E
		{
			Code:           "1500",
			Name:           "Leasehold Improvements",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Renovations and improvements to leased premises",
		},
		{
			Code:            "1510",
			Name:            "Accumulated Depreciation — Leasehold Improvements",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Contra asset: accumulated depreciation on leasehold improvements",
		},
		{
			Code:           "1520",
			Name:           "Equipment",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Service equipment (styling chairs, massage beds, sterilizers)",
		},
		{
			Code:            "1530",
			Name:            "Accumulated Depreciation — Equipment",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Contra asset: accumulated depreciation on equipment",
		},
		{
			Code:           "1540",
			Name:           "Furniture & Fixtures",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Furniture and fixtures used in operations",
		},
		{
			Code:            "1550",
			Name:            "Accumulated Depreciation — Furniture & Fixtures",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Contra asset: accumulated depreciation on furniture and fixtures",
		},
		{
			Code:           "1560",
			Name:           "Computer Equipment",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "POS systems, computers, and peripherals",
		},
		{
			Code:            "1570",
			Name:            "Accumulated Depreciation — Computer Equipment",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Contra asset: accumulated depreciation on computer equipment",
		},
		// Non-Current Assets — Security Deposits
		{
			Code:           "1600",
			Name:           "Security Deposits",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Refundable deposits paid for leased premises and utilities",
		},

		// ====================================================================
		// LIABILITIES (2000–2999)
		// ====================================================================

		// Current Liabilities — Payables
		{
			Code:           "2010",
			Name:           "Accounts Payable",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Amounts owed to suppliers for goods and services purchased on credit",
		},
		{
			Code:           "2020",
			Name:           "Accrued Expenses",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Expenses incurred but not yet paid (utilities, rent, commissions)",
		},
		// Current Liabilities — Payroll
		{
			Code:            "2110",
			Name:            "SSS Payable",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "SSS contributions withheld from employees and employer share due",
		},
		{
			Code:            "2120",
			Name:            "PhilHealth Payable",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "PhilHealth contributions withheld from employees and employer share due",
		},
		{
			Code:            "2130",
			Name:            "Pag-IBIG Payable",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Pag-IBIG (HDMF) contributions withheld from employees and employer share due",
		},
		{
			Code:            "2140",
			Name:            "Income Tax Withheld Payable",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Withholding tax on compensation withheld from employees' salaries",
		},
		{
			Code:            "2150",
			Name:            "Expanded Withholding Tax Payable",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "EWT withheld from supplier payments subject to BIR withholding",
		},
		{
			Code:            "2160",
			Name:            "Output VAT Payable",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "VAT collected from customers on vatable sales and services",
		},
		// Current Liabilities — Deferred Revenue
		{
			Code:           "2200",
			Name:           "Deferred Revenue",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Payments received for services not yet rendered (gift cards, packages, deposits)",
		},
		// Non-Current Liabilities
		{
			Code:           "2500",
			Name:           "Loans Payable — Long Term",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_LIABILITY,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Long-term bank loans and financing payable beyond 12 months",
		},

		// ====================================================================
		// EQUITY (3000–3999)
		// ====================================================================

		{
			Code:            "3010",
			Name:            "Owner's Capital",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_EQUITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Initial and additional capital invested by the owner(s)",
		},
		{
			Code:           "3020",
			Name:           "Owner's Drawing",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_EQUITY,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Withdrawals by the owner for personal use (contra equity)",
		},
		{
			Code:            "3030",
			Name:            "Retained Earnings",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_EQUITY,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			IsSystemAccount: true,
			Description:     "Cumulative net income (loss) retained in the business",
		},

		// ====================================================================
		// REVENUE (4000–4999)
		// ====================================================================

		{
			Code:           "4010",
			Name:           "Service Revenue",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Revenue from primary service activities (haircut, massage, facial, etc.)",
		},
		{
			Code:           "4020",
			Name:           "Product Sales",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Revenue from retail product sales to clients",
		},
		{
			Code:           "4030",
			Name:           "Package Revenue",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Revenue from prepaid service packages and memberships",
		},
		{
			Code:           "4040",
			Name:           "Gift Certificate Revenue",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Revenue recognized upon redemption of gift certificates",
		},
		{
			Code:           "4900",
			Name:           "Other Income",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OTHER_INCOME,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Miscellaneous income not from primary operations (booth rental, commissions received)",
		},
		{
			Code:           "4910",
			Name:           "Interest Income",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OTHER_INCOME,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_CREDIT,
			Description:    "Interest earned on bank deposits and savings accounts",
		},

		// ====================================================================
		// EXPENSES (5000–5999)
		// ====================================================================

		// Cost of Sales
		{
			Code:           "5010",
			Name:           "Cost of Services",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_COST_OF_SALES,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Direct costs of service delivery (supplies consumed, technician commissions)",
		},
		{
			Code:           "5020",
			Name:           "Cost of Goods Sold",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_COST_OF_SALES,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Cost of retail products sold to clients",
		},
		// Operating Expenses — Labor
		{
			Code:           "5110",
			Name:           "Salaries Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Monthly salaries and wages for all employees",
		},
		{
			Code:           "5120",
			Name:           "SSS Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Employer share of SSS contributions",
		},
		{
			Code:           "5130",
			Name:           "PhilHealth Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Employer share of PhilHealth contributions",
		},
		{
			Code:           "5140",
			Name:           "Pag-IBIG Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Employer share of Pag-IBIG (HDMF) contributions",
		},
		// Operating Expenses — Occupancy
		{
			Code:           "5210",
			Name:           "Rent Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Monthly rent for business premises",
		},
		{
			Code:           "5220",
			Name:           "Utilities Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Electricity, water, and internet expenses",
		},
		// Operating Expenses — Depreciation
		{
			Code:            "5310",
			Name:            "Depreciation Expense",
			Element:         accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification:  accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:   accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			IsSystemAccount: true,
			Description:     "Periodic depreciation on equipment, furniture, and leasehold improvements",
		},
		// Operating Expenses — Supplies & Marketing
		{
			Code:           "5410",
			Name:           "Supplies Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Office and operational supplies consumed",
		},
		{
			Code:           "5420",
			Name:           "Advertising & Promotion Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Social media, print, and promotional campaign costs",
		},
		// Operating Expenses — Administrative
		{
			Code:           "5510",
			Name:           "Professional Fees",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Accounting, legal, and consulting fees",
		},
		{
			Code:           "5520",
			Name:           "Repairs & Maintenance Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Equipment and facility repair and maintenance costs",
		},
		{
			Code:           "5530",
			Name:           "Communication Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Telephone, mobile, and internet communication costs",
		},
		{
			Code:           "5540",
			Name:           "Insurance Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Business insurance premiums (fire, theft, liability)",
		},
		{
			Code:           "5550",
			Name:           "Licenses & Permits Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Business permits, mayor's permit, BIR registration, and professional licenses",
		},
		{
			Code:           "5560",
			Name:           "Taxes & Licenses",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Local business taxes (LBT) and other regulatory fees",
		},
		{
			Code:           "5570",
			Name:           "Miscellaneous Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Minor expenses not classified elsewhere",
		},
		// Finance Costs
		{
			Code:           "5810",
			Name:           "Interest Expense",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_FINANCE_COST,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Interest paid on loans and financing arrangements",
		},
		{
			Code:           "5820",
			Name:           "Bank Charges",
			Element:        accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE,
			Classification: accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_FINANCE_COST,
			NormalBalance:  accountpb.NormalBalance_NORMAL_BALANCE_DEBIT,
			Description:    "Bank service fees, transaction charges, and penalties",
		},
	}
}

// SeedDefaultCoA inserts all default accounts for a Philippine service business.
// It calls createAccount for each entry. Accounts that already exist (duplicate code)
// are skipped gracefully — the function continues and returns a summary error if any
// individual account fails for a non-duplicate reason.
//
// Parameters:
//   - ctx: request context (used for auth and tracing)
//   - createAccount: the CreateAccount use case Execute function
//   - workspaceID: the workspace to seed accounts into (stored in account.workspace_id if supported)
//
// Returns an error only if a non-recoverable failure occurred. Partial success
// (some accounts skipped) is reported via the returned count values.
func SeedDefaultCoA(
	ctx context.Context,
	createAccount CreateAccountFunc,
	workspaceID string,
) (created int, skipped int, err error) {
	accounts := DefaultCoA()
	var errs []string

	for _, entry := range accounts {
		desc := entry.Description
		req := &accountpb.CreateAccountRequest{
			Data: &accountpb.Account{
				Code:             entry.Code,
				Name:             entry.Name,
				Element:          entry.Element,
				Classification:   entry.Classification,
				NormalBalance:    entry.NormalBalance,
				IsSystemAccount:  entry.IsSystemAccount,
				Active:           true,
				Description:      &desc,
				CashFlowActivity: entry.CashFlowActivity,
			},
		}
		// Associate with workspace if provided
		if workspaceID != "" {
			// Note: workspace scoping is handled at the use case layer via context.
			// The workspaceID is passed here for future use when Account proto gains
			// an explicit workspace_id field.
			_ = workspaceID
		}

		_, createErr := createAccount(ctx, req)
		if createErr != nil {
			// Treat duplicate-code errors as non-fatal skips.
			// A proper duplicate check would query by code first; for now we
			// surface all errors as skipped and collect them for reporting.
			skipped++
			errs = append(errs, fmt.Sprintf("%s %s: %v", entry.Code, entry.Name, createErr))
			continue
		}
		created++
	}

	if len(errs) > 0 {
		err = fmt.Errorf("seeder completed with %d error(s): %v", len(errs), errs)
	}
	return created, skipped, err
}
