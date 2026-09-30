package block

import (
	"strings"
	"testing"

	"github.com/erniealice/espyna-golang/consumer/compose"
)

// UI-CP-10: the charge policy unit is opt-in. Without the option the unit set is
// unchanged (so service-admin / school-admin compositions stay byte-identical);
// with it, exactly one unit is appended.
func TestChargePoliciesAreOptIn(t *testing.T) {
	base := AllUnits(nil, nil)
	for _, u := range base {
		if u.Key == "ledger.charge_policy" {
			t.Fatal("charge policy unit present without WithChargePolicies")
		}
	}
	off := AllUnits(nil, nil, WithChargePolicies(false))
	if len(off) != len(base) {
		t.Fatalf("WithChargePolicies(false) changed the unit set: %d vs %d", len(off), len(base))
	}
	on := AllUnits(nil, nil, WithChargePolicies(true))
	if len(on) != len(base)+1 {
		t.Fatalf("WithChargePolicies(true) added %d units, want 1", len(on)-len(base))
	}
	last := on[len(on)-1]
	if last.Key != "ledger.charge_policy" {
		t.Fatalf("appended unit = %q", last.Key)
	}
	for i := range base {
		if base[i].Key != on[i].Key {
			t.Fatalf("default unit order changed at %d: %s vs %s", i, base[i].Key, on[i].Key)
		}
	}
	if last.Nav.Permission != "charge_policy:list" || len(last.Nav.Items) != 1 || last.Nav.Items[0].Route != "charge_policy.list" {
		t.Errorf("nav contribution = %+v", last.Nav)
	}
}

func TestLegacyChargePolicyOptionIsNeverImpliedByEnableAll(t *testing.T) {
	cfg := &blockConfig{enableAll: true}
	if cfg.wantChargePolicy() {
		t.Error("enableAll must not enable charge policies")
	}
	WithChargePolicy()(cfg)
	if !cfg.wantChargePolicy() {
		t.Error("WithChargePolicy did not enable the module")
	}
}

// C13: an option-gated module's closures are in RequireFor — enabling the charge
// policy pages with any closure unbound refuses boot instead of degrading at runtime.
func TestRequireForCoversChargePolicy(t *testing.T) {
	cfg := &blockConfig{}
	WithChargePolicy()(cfg)
	err := (&UseCases{}).RequireFor(cfg)
	if err == nil {
		t.Fatal("RequireFor accepted an empty UseCases with the charge policy module enabled")
	}
	for _, name := range []string{
		"UseCases.ChargePolicy.CreateChargePolicy",
		"UseCases.ChargePolicy.GetChargePolicyListPageData",
		"UseCases.ChargePolicy.ApproveChargePolicyVersion",
		"UseCases.ChargePolicy.DeleteChargePolicyPosting",
	} {
		if !strings.Contains(err.Error(), name) {
			t.Errorf("RequireFor error does not name %s: %v", name, err)
		}
	}
	// Module off: no charge policy requirement.
	if err := (&UseCases{}).RequireFor(&blockConfig{}); err != nil && strings.Contains(err.Error(), "ChargePolicy") {
		t.Errorf("charge policy closures required while the module is off: %v", err)
	}
}

// R4 M1: compose-v2 Mount does not run RequireFor, so the unit itself must fail
// closed when a required closure is unbound.
func TestChargePolicyUnitMountFailsClosed(t *testing.T) {
	units := AllUnits(&UseCases{}, nil, WithChargePolicies(true))
	last := units[len(units)-1]
	if last.Key != "ledger.charge_policy" {
		t.Fatalf("appended unit = %q", last.Key)
	}
	if err := last.Mount(&compose.MountContext{}); err == nil {
		t.Fatal("Mount accepted an empty UseCases")
	}
	nilUnit := ChargePolicyUnit(nil, nil)
	if err := nilUnit.Mount(&compose.MountContext{}); err == nil {
		t.Fatal("Mount accepted nil UseCases")
	}
}
