// routes.go defines configurable route structs for fycha ledger recurring_template views.
package recurring_template

const (
	// Ledger — Recurring Templates
	RecurringTemplatesURL = "/ledger/settings/recurring"
)

// ---------------------------------------------------------------------------
// SettingsRoutes
// ---------------------------------------------------------------------------

// SettingsRoutes holds route paths for ledger settings views
// (Bad Debt Policy, Recurring Templates).
type SettingsRoutes struct {
	ActiveNav             string `json:"active_nav"`
	BadDebtPolicyURL      string `json:"bad_debt_policy_url"`
	RecurringTemplatesURL string `json:"recurring_templates_url"`
}

func DefaultSettingsRoutes() SettingsRoutes {
	return SettingsRoutes{
		ActiveNav:             "ledger",
		BadDebtPolicyURL:      "/ledger/settings/bad-debt-policy",
		RecurringTemplatesURL: RecurringTemplatesURL,
	}
}

func (r SettingsRoutes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.settings.bad_debt_policy":     r.BadDebtPolicyURL,
		"ledger.settings.recurring_templates": r.RecurringTemplatesURL,
	}
}
