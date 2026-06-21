package funding

import (
	"github.com/erniealice/espyna-golang/consumer/compose"

	fundinglabels "github.com/erniealice/fycha-golang/domain/funding/funding/labels"
)

// Describe returns the funding view module's compose.Unit. This unit is a
// view-module contributor: it supplies an AppEntry, sidebar Items, and
// lyngua-overlayable FundingFormLabels for the "funding" app. Entity-level
// routes (fund, fund_allocation, fund_transaction) live in the centymo
// descriptor units; this unit owns the view-module wiring via its Mount
// closure (set by the block catalog).
func Describe() compose.Unit {
	l := fundinglabels.DefaultFundingFormLabels()
	return compose.Unit{
		Key:       "funding.funding",
		Templates: TemplatesFS,
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "funding.json", Key: "funding"},
		LabelName: "FundingFormLabels",
		Nav: compose.NavContrib{
			Permission: "fund:list",
			AppEntry: &compose.AppEntry{
				Key:        "funding",
				Route:      "fund.source.list",
				Label:      "Funding",
				Icon:       "icon-wallet",
				Permission: "fund:list|fund_allocation:list",
			},
			// NOTE: nav items (sources-all / cards-* / activity-*) removed —
			// they reference route keys "fund.source.list", "fund.card.list",
			// and "fund_transaction.list", but the fund / fund_card /
			// fund_transaction entities are scaffold-only (labels.go, no
			// routes.go / descriptor / mounted unit) in BOTH fycha and centymo.
			// Within compose, phase-3 nav validation runs against the per-engine
			// RouteMap, so these unmounted-entity references fail-closed the
			// entire fycha engine. Restore once the fund* entities ship route
			// structs and are mounted in their owning engine's AllUnits.
			Items: []compose.NavItem{},
		},
	}
}
