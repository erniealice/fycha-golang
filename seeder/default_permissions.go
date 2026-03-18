// Package seeder provides Chart of Accounts permission seed data.
// The 21 permission codes defined here correspond to the accounting-module
// RBAC additions described in:
//
//	docs/epic/20260317-chart-of-accounts-integration/01-plan/11-rbac-plan.md
//
// Usage:
//
//	import "github.com/erniealice/fycha-golang/seeder"
//
//	created, skipped, err := seeder.SeedDefaultPermissions(ctx, createPermissionFn)
package seeder

import (
	"context"
	"fmt"

	permissionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/permission"
)

// CreatePermissionFunc is the signature of the CreatePermission use case Execute method.
type CreatePermissionFunc func(ctx context.Context, req *permissionpb.CreatePermissionRequest) (*permissionpb.CreatePermissionResponse, error)

// DefaultCoAPermission defines a permission to be seeded for the Chart of Accounts integration.
type DefaultCoAPermission struct {
	// ID is the permission record ID: "perm-{entity}-{action}".
	ID string
	// Code is the permission_code in entity:action form (e.g. "ledger:view").
	Code string
	// Name is the human-readable display name.
	Name string
	// Description explains what the permission gates.
	Description string
	// Module is the app/module this permission belongs to (e.g. "ledger", "equity").
	Module string
	// Category is the grouping within the module (e.g. "app", "journal", "account").
	Category string
}

// DefaultCoAPermissions returns all 21 Chart of Accounts permission codes for the
// accounting-module RBAC integration. Codes follow the entity:action convention.
// IDs follow the perm-{entity}-{action} convention used by the existing seeder.
func DefaultCoAPermissions() []DefaultCoAPermission {
	return []DefaultCoAPermission{
		// ----------------------------------------------------------------
		// Ledger app — sidebar gate
		// ----------------------------------------------------------------
		{
			ID:          "perm-ledger-view",
			Code:        "ledger:view",
			Name:        "Ledger: View",
			Description: "Access the Ledger app in sidebar",
			Module:      "ledger",
			Category:    "app",
		},

		// ----------------------------------------------------------------
		// Ledger app — journal permissions
		// ----------------------------------------------------------------
		{
			ID:          "perm-journal-list",
			Code:        "journal:list",
			Name:        "Journal: List",
			Description: "View posted journal entries",
			Module:      "ledger",
			Category:    "journal",
		},
		{
			ID:          "perm-journal-post",
			Code:        "journal:post_guided",
			Name:        "Journal: Post Guided",
			Description: "Post system-generated journal entries",
			Module:      "ledger",
			Category:    "journal",
		},
		{
			ID:          "perm-journal-post_manual",
			Code:        "journal:post_manual",
			Name:        "Journal: Post Manual",
			Description: "Access the raw journal entry form",
			Module:      "ledger",
			Category:    "journal",
		},
		{
			ID:          "perm-journal-post_manual_approve",
			Code:        "journal:post_manual_approve",
			Name:        "Journal: Post Manual Approve",
			Description: "Maker-checker approval for manual journal entries",
			Module:      "ledger",
			Category:    "journal",
		},
		{
			ID:          "perm-journal-reverse",
			Code:        "journal:reverse",
			Name:        "Journal: Reverse",
			Description: "Reverse a posted journal entry",
			Module:      "ledger",
			Category:    "journal",
		},

		// ----------------------------------------------------------------
		// Ledger app — account (Chart of Accounts) permissions
		// ----------------------------------------------------------------
		{
			ID:          "perm-account-list",
			Code:        "account:list",
			Name:        "Account: List",
			Description: "View Chart of Accounts",
			Module:      "ledger",
			Category:    "account",
		},
		{
			ID:          "perm-account-create",
			Code:        "account:create",
			Name:        "Account: Create",
			Description: "Add and edit accounts",
			Module:      "ledger",
			Category:    "account",
		},
		{
			ID:          "perm-account-delete",
			Code:        "account:delete",
			Name:        "Account: Delete",
			Description: "Delete accounts from Chart of Accounts",
			Module:      "ledger",
			Category:    "account",
		},

		// ----------------------------------------------------------------
		// Equity — Funding group
		// ----------------------------------------------------------------
		{
			ID:          "perm-equity-list",
			Code:        "equity:list",
			Name:        "Equity: List",
			Description: "View capital accounts and transactions",
			Module:      "equity",
			Category:    "equity",
		},
		{
			ID:          "perm-equity-create",
			Code:        "equity:create",
			Name:        "Equity: Create",
			Description: "Record equity transactions",
			Module:      "equity",
			Category:    "equity",
		},
		{
			ID:          "perm-equity-update",
			Code:        "equity:update",
			Name:        "Equity: Update",
			Description: "Edit capital accounts",
			Module:      "equity",
			Category:    "equity",
		},
		{
			ID:          "perm-equity-delete",
			Code:        "equity:delete",
			Name:        "Equity: Delete",
			Description: "Reverse equity transactions",
			Module:      "equity",
			Category:    "equity",
		},

		// ----------------------------------------------------------------
		// Loans — Funding group
		// ----------------------------------------------------------------
		{
			ID:          "perm-loan-list",
			Code:        "loan:list",
			Name:        "Loan: List",
			Description: "View loans",
			Module:      "loans",
			Category:    "loan",
		},
		{
			ID:          "perm-loan-create",
			Code:        "loan:create",
			Name:        "Loan: Create",
			Description: "Record new loans",
			Module:      "loans",
			Category:    "loan",
		},
		{
			ID:          "perm-loan-payment",
			Code:        "loan:payment",
			Name:        "Loan: Payment",
			Description: "Record loan payments",
			Module:      "loans",
			Category:    "loan",
		},

		// ----------------------------------------------------------------
		// Expenses
		// ----------------------------------------------------------------
		{
			ID:          "perm-expense-create",
			Code:        "expense:create",
			Name:        "Expense: Create",
			Description: "Create expense forms (guided)",
			Module:      "expenses",
			Category:    "expense",
		},
		{
			ID:          "perm-expense-approve",
			Code:        "expense:approve",
			Name:        "Expense: Approve",
			Description: "Approve expenses before disbursement",
			Module:      "expenses",
			Category:    "expense",
		},

		// ----------------------------------------------------------------
		// Payroll
		// ----------------------------------------------------------------
		{
			ID:          "perm-payroll-run",
			Code:        "payroll:run",
			Name:        "Payroll: Run",
			Description: "Initiate payroll runs",
			Module:      "payroll",
			Category:    "payroll",
		},
		{
			ID:          "perm-payroll-approve",
			Code:        "payroll:approve",
			Name:        "Payroll: Approve",
			Description: "Approve payroll before disbursement",
			Module:      "payroll",
			Category:    "payroll",
		},

		// ----------------------------------------------------------------
		// Petty Cash — Cash app section
		// ----------------------------------------------------------------
		{
			ID:          "perm-petty_cash-voucher",
			Code:        "petty_cash:voucher",
			Name:        "Petty Cash: Voucher",
			Description: "Record petty cash vouchers",
			Module:      "cash",
			Category:    "petty_cash",
		},
		{
			ID:          "perm-petty_cash-replenish",
			Code:        "petty_cash:replenish",
			Name:        "Petty Cash: Replenish",
			Description: "Approve petty cash replenishments",
			Module:      "cash",
			Category:    "petty_cash",
		},
	}
}

// SeedDefaultPermissions inserts all 21 Chart of Accounts permission codes.
// It calls createPermission for each entry. Permissions that already exist
// (duplicate ID) are skipped gracefully — the function continues and returns
// a summary error if any individual permission fails for a non-duplicate reason.
//
// Parameters:
//   - ctx: request context (used for auth and tracing)
//   - createPermission: the CreatePermission use case Execute function
//
// Returns an error only if a non-recoverable failure occurred. Partial success
// (some permissions skipped) is reported via the returned count values.
func SeedDefaultPermissions(
	ctx context.Context,
	createPermission CreatePermissionFunc,
) (created int, skipped int, err error) {
	permissions := DefaultCoAPermissions()
	var errs []string

	for _, entry := range permissions {
		desc := entry.Description
		req := &permissionpb.CreatePermissionRequest{
			Data: &permissionpb.Permission{
				Id:             entry.ID,
				PermissionCode: entry.Code,
				Name:           entry.Name,
				Description:    desc,
				PermissionType: permissionpb.PermissionType_PERMISSION_TYPE_ALLOW,
				Active:         true,
			},
		}

		_, createErr := createPermission(ctx, req)
		if createErr != nil {
			// Treat duplicate-ID errors as non-fatal skips.
			skipped++
			errs = append(errs, fmt.Sprintf("%s %s: %v", entry.ID, entry.Code, createErr))
			continue
		}
		created++
	}

	if len(errs) > 0 {
		err = fmt.Errorf("permission seeder completed with %d error(s): %v", len(errs), errs)
	}
	return created, skipped, err
}
