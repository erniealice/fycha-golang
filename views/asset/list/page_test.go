package list

// Phase 5 (H5) — soft-delete gate tests for buildTableRows.
//
// Verifies that:
//   - An asset with transactions (inUseIDs[id] = true) renders data-deletable="false"
//     and has its delete action disabled with the CannotDeleteInUse tooltip.
//   - An asset with no transactions (inUseIDs[id] = false) renders data-deletable="true"
//     and has its delete action enabled.
//   - An asset that belongs to a different workspace (not in inUseIDs) is treated
//     as not in-use (map lookup returns false for missing key — workspace scoping
//     is enforced by the checker itself; the list page honours the result faithfully).

import (
	"testing"

	"github.com/erniealice/pyeza-golang/types"

	fycha "github.com/erniealice/fycha-golang"
)

// testPermsWithDelete returns a UserPermissions that allows all asset actions.
func testPermsWithDelete() *types.UserPermissions {
	return types.NewUserPermissions([]string{"asset:create", "asset:read", "asset:update", "asset:delete"})
}

func testAssetLabels() fycha.AssetLabels {
	l := fycha.DefaultAssetLabels()
	l.Actions.CannotDeleteInUse = "Cannot delete: asset has posted transactions."
	return l
}

func testRoutes() fycha.AssetRoutes {
	return fycha.AssetRoutes{
		DetailURL:    "/assets/{id}",
		EditURL:      "/assets/{id}/edit",
		DeleteURL:    "/assets/delete",
		SetStatusURL: "/assets/set-status",
	}
}

// TestBuildTableRows_InUseAsset_DataDeletableFalse verifies that an asset with
// at least one asset_transaction row (inUseIDs[id]=true) gets data-deletable="false"
// and a disabled delete action with the CannotDeleteInUse tooltip.
func TestBuildTableRows_InUseAsset_DataDeletableFalse(t *testing.T) {
	t.Parallel()

	assets := []AssetRow{
		{ID: "asset-with-tx", Name: "Laptop", AssetNumber: "FA-001", Active: true},
	}
	inUseIDs := map[string]bool{"asset-with-tx": true}

	rows := buildTableRows(assets, testAssetLabels(), testRoutes(), testPermsWithDelete(), "active", inUseIDs)

	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]

	// data-deletable attribute must be "false"
	if got := row.DataAttrs["deletable"]; got != "false" {
		t.Errorf("data-deletable = %q, want %q", got, "false")
	}

	// Delete action must be disabled
	var deleteAction *types.TableAction
	for i := range row.Actions {
		if row.Actions[i].Action == "delete" {
			deleteAction = &row.Actions[i]
			break
		}
	}
	if deleteAction == nil {
		t.Fatal("delete action not found in row actions")
	}
	if !deleteAction.Disabled {
		t.Error("delete action should be disabled for in-use asset")
	}
	wantTooltip := testAssetLabels().Actions.CannotDeleteInUse
	if deleteAction.DisabledTooltip != wantTooltip {
		t.Errorf("DisabledTooltip = %q, want %q", deleteAction.DisabledTooltip, wantTooltip)
	}
}

// TestBuildTableRows_NotInUseAsset_DataDeletableTrue verifies that an asset with
// no asset_transaction rows (not in inUseIDs) gets data-deletable="true" and an
// enabled delete action.
func TestBuildTableRows_NotInUseAsset_DataDeletableTrue(t *testing.T) {
	t.Parallel()

	assets := []AssetRow{
		{ID: "asset-no-tx", Name: "Chair", AssetNumber: "FA-002", Active: true},
	}
	inUseIDs := map[string]bool{} // empty — no assets in use

	rows := buildTableRows(assets, testAssetLabels(), testRoutes(), testPermsWithDelete(), "active", inUseIDs)

	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]

	// data-deletable attribute must be "true"
	if got := row.DataAttrs["deletable"]; got != "true" {
		t.Errorf("data-deletable = %q, want %q", got, "true")
	}

	// Delete action must be enabled
	var deleteAction *types.TableAction
	for i := range row.Actions {
		if row.Actions[i].Action == "delete" {
			deleteAction = &row.Actions[i]
			break
		}
	}
	if deleteAction == nil {
		t.Fatal("delete action not found in row actions")
	}
	if deleteAction.Disabled {
		t.Errorf("delete action should be enabled for asset with no transactions; DisabledTooltip=%q", deleteAction.DisabledTooltip)
	}
}

// TestBuildTableRows_CrossWorkspaceAsset_TreatedAsNotInUse verifies that the
// workspace scoping contract is honoured: an asset whose ID does not appear in
// inUseIDs (because the checker filtered it out by workspace) is treated as
// not in-use. The list page page.go trusts the checker's result faithfully.
func TestBuildTableRows_CrossWorkspaceAsset_TreatedAsNotInUse(t *testing.T) {
	t.Parallel()

	// Simulate: workspace-A checker returned only workspace-A asset IDs.
	// asset-other-ws belongs to workspace-B and is absent from the map.
	assets := []AssetRow{
		{ID: "asset-other-ws", Name: "Server", AssetNumber: "FA-003", Active: true},
	}
	inUseIDs := map[string]bool{
		// Only workspace-A asset IDs would appear here — workspace-B ID absent.
		"asset-workspace-a": true,
	}

	rows := buildTableRows(assets, testAssetLabels(), testRoutes(), testPermsWithDelete(), "active", inUseIDs)

	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]

	// asset-other-ws is not in the inUseIDs map → treated as not in-use
	if got := row.DataAttrs["deletable"]; got != "true" {
		t.Errorf("data-deletable = %q, want %q (cross-workspace asset absent from map = not in-use)", got, "true")
	}
}

// TestBuildTableRows_NilInUseMap_AllDeletable verifies nil-safety: when
// inUseIDs is nil (checker unavailable), all assets render as deletable.
func TestBuildTableRows_NilInUseMap_AllDeletable(t *testing.T) {
	t.Parallel()

	assets := []AssetRow{
		{ID: "asset-001", Name: "Desk", AssetNumber: "FA-004", Active: true},
		{ID: "asset-002", Name: "Monitor", AssetNumber: "FA-005", Active: true},
	}

	rows := buildTableRows(assets, testAssetLabels(), testRoutes(), testPermsWithDelete(), "active", nil)

	if len(rows) != 2 {
		t.Fatalf("got %d rows, want 2", len(rows))
	}
	for _, row := range rows {
		if got := row.DataAttrs["deletable"]; got != "true" {
			t.Errorf("row %s: data-deletable = %q, want %q (nil map = all deletable)", row.ID, got, "true")
		}
	}
}

// TestBuildTableRows_InUseNoPermission_DisabledWithInUseTooltip verifies that
// when an asset is both in-use AND the user lacks delete permission, the
// CannotDeleteInUse tooltip takes priority (more informative than NoPermission).
func TestBuildTableRows_InUseNoPermission_DisabledWithInUseTooltip(t *testing.T) {
	t.Parallel()

	assets := []AssetRow{
		{ID: "asset-in-use", Name: "Printer", AssetNumber: "FA-006", Active: true},
	}
	inUseIDs := map[string]bool{"asset-in-use": true}

	// No delete permission
	noDeletePerms := types.NewUserPermissions([]string{"asset:create", "asset:read", "asset:update"})

	rows := buildTableRows(assets, testAssetLabels(), testRoutes(), noDeletePerms, "active", inUseIDs)

	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]

	var deleteAction *types.TableAction
	for i := range row.Actions {
		if row.Actions[i].Action == "delete" {
			deleteAction = &row.Actions[i]
			break
		}
	}
	if deleteAction == nil {
		t.Fatal("delete action not found")
	}
	if !deleteAction.Disabled {
		t.Error("delete action should be disabled")
	}
	wantTooltip := testAssetLabels().Actions.CannotDeleteInUse
	if deleteAction.DisabledTooltip != wantTooltip {
		t.Errorf("DisabledTooltip = %q, want CannotDeleteInUse tooltip %q", deleteAction.DisabledTooltip, wantTooltip)
	}
}

// TestBuildTableRows_PermDeniedNotInUse_DisabledWithNoPermissionTooltip verifies
// that when the asset is NOT in-use but the user lacks delete permission, the
// NoPermission tooltip is shown.
func TestBuildTableRows_PermDeniedNotInUse_DisabledWithNoPermissionTooltip(t *testing.T) {
	t.Parallel()

	assets := []AssetRow{
		{ID: "asset-ok", Name: "Whiteboard", AssetNumber: "FA-007", Active: true},
	}
	inUseIDs := map[string]bool{} // not in use

	noDeletePerms := types.NewUserPermissions([]string{"asset:create", "asset:read", "asset:update"})
	l := testAssetLabels()

	rows := buildTableRows(assets, l, testRoutes(), noDeletePerms, "active", inUseIDs)

	if len(rows) != 1 {
		t.Fatalf("got %d rows, want 1", len(rows))
	}
	row := rows[0]

	var deleteAction *types.TableAction
	for i := range row.Actions {
		if row.Actions[i].Action == "delete" {
			deleteAction = &row.Actions[i]
			break
		}
	}
	if deleteAction == nil {
		t.Fatal("delete action not found")
	}
	if !deleteAction.Disabled {
		t.Error("delete action should be disabled (no delete permission)")
	}
	if deleteAction.DisabledTooltip != l.Actions.NoPermission {
		t.Errorf("DisabledTooltip = %q, want NoPermission tooltip %q", deleteAction.DisabledTooltip, l.Actions.NoPermission)
	}
}

