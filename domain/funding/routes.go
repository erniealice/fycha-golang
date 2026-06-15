// routes.go defines configurable route structs for fycha funding domain views.
//
// FundRoutes holds the five route URLs for the Fund primitive
// (cross-workspace fund sources + cards). The compose engine in service-admin
// uses DefaultFundRoutes() to register sidebar nav entries and populate the
// route map used by template {{route}} / {{routeWith}} functions.
//
// Post-P12 (2026-05-22) the /app/ prefix was removed from all fund routes —
// workspace_path middleware dispatches /w/{slug}/funding/sources to the bare
// /funding/sources handler.
package funding

// FundRoutes holds routes for the Fund primitive (cross-workspace fund sources + cards).
type FundRoutes struct {
	SourcesURL        string
	SourceDetailURL   string
	CardsURL          string
	CardDetailURL     string
	ReconciliationURL string
}

// DefaultFundRoutes returns the default Fund route URLs.
func DefaultFundRoutes() FundRoutes {
	return FundRoutes{
		SourcesURL:        "/funding/sources",
		SourceDetailURL:   "/funding/sources/{fund_id}",
		CardsURL:          "/funding/cards",
		CardDetailURL:     "/funding/cards/{allocation_id}",
		ReconciliationURL: "/funding/reports/gl-reconciliation",
	}
}

// RouteMap returns dot-notation keys for template route lookups.
func (r FundRoutes) RouteMap() map[string]string {
	return map[string]string{
		"fund.source.list":    r.SourcesURL,
		"fund.source.detail":  r.SourceDetailURL,
		"fund.card.list":      r.CardsURL,
		"fund.card.detail":    r.CardDetailURL,
		"fund.reconciliation": r.ReconciliationURL,
	}
}
