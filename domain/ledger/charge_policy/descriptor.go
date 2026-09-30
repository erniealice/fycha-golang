package charge_policy

import "github.com/erniealice/espyna-golang/consumer/compose"

// Describe returns the compose descriptor for the charge policy pages.
// The unit is mounted by fycha's ChargePolicyUnit only when the app opts in
// (block.WithChargePolicies); it is not part of the default fycha unit set.
func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "ledger.charge_policy",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "charge_policy"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "charge_policy.json", Key: "charge_policy"},
		LabelName: "ChargePolicyLabels",
		Templates: TemplatesFS,
		Nav: compose.NavContrib{
			Permission: "charge_policy:list",
			Items: []compose.NavItem{
				{Key: "charge-policies", Route: "charge_policy.list", Params: map[string]string{"status": "active"},
					Label: "Charge Policies", Icon: "icon-shield-check", Permission: "charge_policy:list",
					LabelKey: "charge_policies_label", IconKey: "charge_policies_icon"},
			},
		},
	}
}
