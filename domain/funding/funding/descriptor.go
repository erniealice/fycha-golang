package funding

import (
	"github.com/erniealice/pyeza-golang/compose"

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
			Items: []compose.NavItem{
				{Key: "sources-all", Route: "fund.source.list", Label: "All Sources", Icon: "icon-list", Permission: "fund:list"},
				{Key: "cards-active", Route: "fund.card.list", Query: "?status=active", Label: "Active", Icon: "icon-credit-card", Permission: "fund_allocation:list"},
				{Key: "cards-suspended", Route: "fund.card.list", Query: "?status=suspended", Label: "Suspended", Icon: "icon-pause-circle", Permission: "fund_allocation:list"},
				{Key: "cards-closed", Route: "fund.card.list", Query: "?status=closed", Label: "Closed", Icon: "icon-x-circle", Permission: "fund_allocation:list"},
				{Key: "activity-posted", Route: "fund_transaction.list", Params: map[string]string{"status": "posted"}, Label: "Posted", Icon: "icon-check-circle", Permission: "fund_transaction:list"},
				{Key: "activity-draft", Route: "fund_transaction.list", Params: map[string]string{"status": "draft"}, Label: "Draft", Icon: "icon-edit", Permission: "fund_transaction:list"},
				{Key: "activity-voided", Route: "fund_transaction.list", Params: map[string]string{"status": "voided"}, Label: "Voided", Icon: "icon-x-circle", Permission: "fund_transaction:list"},
			},
		},
	}
}
