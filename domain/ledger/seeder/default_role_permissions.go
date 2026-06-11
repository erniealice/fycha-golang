package seeder

// RolePermissionMapping defines which permission codes belong to a default role.
// This is data only — actual role-permission assignment uses the existing
// role-permission use cases and the seed-rbac seeder in
// apps/service-admin/cmd/seed-rbac/main.go.
type RolePermissionMapping struct {
	// RoleID is the seeded role identifier (e.g. "role-admin").
	RoleID string
	// RoleName is the human-readable role name.
	RoleName string
	// PermissionCodes is the list of permission_code values assigned to this role.
	PermissionCodes []string
}

// allCoAPermCodes is the full set of 21 permission codes for convenient reuse.
var allCoAPermCodes = []string{
	// Ledger
	"ledger:view",
	"journal:list",
	"journal:post_guided",
	"journal:post_manual",
	"journal:post_manual_approve",
	"journal:reverse",
	"account:list",
	"account:create",
	"account:delete",
	// Equity
	"equity:list",
	"equity:create",
	"equity:update",
	"equity:delete",
	// Loans
	"loan:list",
	"loan:create",
	"loan:payment",
	// Expenses
	"expense:create",
	"expense:approve",
	// Payroll
	"payroll:run",
	"payroll:approve",
	// Petty Cash
	"petty_cash:voucher",
	"petty_cash:replenish",
}

// DefaultCoARolePermissions returns the role-to-permission-code mappings for all
// default roles that interact with Chart of Accounts features.
//
// These mappings reflect the agreed design from:
//
//	docs/epic/20260317-chart-of-accounts-integration/01-plan/11-rbac-plan.md
//
// Assignment logic is in apps/service-admin/cmd/seed-rbac/main.go — this
// function provides the data definitions only.
func DefaultCoARolePermissions() []RolePermissionMapping {
	return []RolePermissionMapping{
		{
			// Owner / Admin: all 21 CoA permissions.
			RoleID:          "role-admin",
			RoleName:        "Admin",
			PermissionCodes: allCoAPermCodes,
		},
		{
			// Accountant: ledger management + post guided JEs + read-only loans/equity/payroll.
			// Cannot post manual JEs or approve payroll (those are CFO/Owner actions).
			RoleID:   "role-accountant",
			RoleName: "Accountant",
			PermissionCodes: []string{
				"ledger:view",
				"journal:list",
				"journal:post_guided",
				"account:list",
				"account:create",
				"equity:list",
				"loan:list",
				"petty_cash:replenish",
			},
		},
		{
			// Bookkeeper: view ledger and journals, post guided JEs — day-to-day maker role.
			// Cannot manage accounts or access funding/payroll.
			RoleID:   "role-bookkeeper",
			RoleName: "Bookkeeper",
			PermissionCodes: []string{
				"ledger:view",
				"journal:list",
				"journal:post_guided",
				"petty_cash:voucher",
			},
		},
		{
			// Manager: approve expenses and payroll; approve petty cash replenishments.
			RoleID:   "role-operations-staff",
			RoleName: "Operations Staff",
			PermissionCodes: []string{
				"expense:approve",
				"payroll:approve",
				"petty_cash:replenish",
			},
		},
		{
			// Frontline Staff: create expense forms and record petty cash vouchers.
			RoleID:   "role-people-manager",
			RoleName: "People Manager",
			PermissionCodes: []string{
				"expense:create",
				"petty_cash:voucher",
			},
		},
	}
}
