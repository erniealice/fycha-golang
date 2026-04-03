package action

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// testDeps returns a Deps with all labels populated and mock CRUD functions.
func testDeps() *Deps {
	return &Deps{
		Routes: fycha.AssetRoutes{
			AddURL:    "/assets/add",
			EditURL:   "/assets/{id}/edit",
			DeleteURL: "/assets/delete",
		},
		Labels: fycha.AssetLabels{
			Form: fycha.AssetFormLabels{
				Name: "Name",
			},
			Actions: fycha.AssetActionLabels{
				InvalidFormData:     "Invalid form data",
				IDRequired:          "ID is required",
				NoIDsProvided:       "No IDs provided",
				InvalidStatus:       "Invalid status",
				InvalidTargetStatus: "Invalid target status",
			},
		},
		CreateAsset: func(ctx context.Context, asset *AssetRecord) error {
			return nil
		},
		ReadAsset: func(ctx context.Context, id string) (*AssetRecord, error) {
			return &AssetRecord{
				ID:               id,
				Name:             "Test Asset",
				AssetNumber:      "FA-001",
				AcquisitionCost:  50000,
				SalvageValue:     5000,
				UsefulLifeMonths: 60,
			}, nil
		},
		UpdateAsset: func(ctx context.Context, asset *AssetRecord) error {
			return nil
		},
		DeleteAsset: func(ctx context.Context, id string) error {
			return nil
		},
		SetActive: func(ctx context.Context, id string, active bool) error {
			return nil
		},
		NewID: func() string { return "test-id-123" },
	}
}

// ctxWithPerms returns a context with the given permission codes.
func ctxWithPerms(codes ...string) context.Context {
	perms := types.NewUserPermissions(codes)
	return view.WithUserPermissions(context.Background(), perms)
}

// ctxNoPerms returns a context with an empty permission set (no permissions granted).
func ctxNoPerms() context.Context {
	perms := types.NewUserPermissions([]string{})
	return view.WithUserPermissions(context.Background(), perms)
}

// ----- NewAddAction tests -----

func TestNewAddAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewAddAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/add", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxNoPerms(), viewCtx)

	if result.Error == nil {
		t.Fatal("expected error for permission denied")
	}
	if result.Error.Error() != "permission denied" {
		t.Errorf("error = %q, want %q", result.Error.Error(), "permission denied")
	}
	if result.StatusCode != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusInternalServerError)
	}
}

func TestNewAddAction_MissingName(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewAddAction(deps)

	// POST with no name field
	form := url.Values{}
	form.Set("acquisition_cost", "10000")
	req := httptest.NewRequest(http.MethodPost, "/assets/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:create"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Name is required" {
		t.Errorf("HX-Error-Message = %q, want %q", result.Headers["HX-Error-Message"], "Name is required")
	}
}

func TestNewAddAction_EmptyStringName(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewAddAction(deps)

	form := url.Values{}
	form.Set("name", "")
	req := httptest.NewRequest(http.MethodPost, "/assets/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:create"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d for empty name", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestNewAddAction_InvalidNumericFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		acquisitionCost string
		salvageValue    string
		usefulLife      string
	}{
		{
			name:            "non-numeric acquisition cost",
			acquisitionCost: "not-a-number",
			salvageValue:    "5000",
			usefulLife:      "60",
		},
		{
			name:            "non-numeric salvage value",
			acquisitionCost: "50000",
			salvageValue:    "abc",
			usefulLife:      "60",
		},
		{
			name:            "non-numeric useful life",
			acquisitionCost: "50000",
			salvageValue:    "5000",
			usefulLife:      "five-years",
		},
		{
			name:            "negative acquisition cost",
			acquisitionCost: "-10000",
			salvageValue:    "5000",
			usefulLife:      "60",
		},
		{
			name:            "negative salvage value",
			acquisitionCost: "50000",
			salvageValue:    "-5000",
			usefulLife:      "60",
		},
		{
			name:            "negative useful life",
			acquisitionCost: "50000",
			salvageValue:    "5000",
			usefulLife:      "-60",
		},
		{
			name:            "all empty numeric fields",
			acquisitionCost: "",
			salvageValue:    "",
			usefulLife:      "",
		},
		{
			name:            "float overflow",
			acquisitionCost: "99999999999999999999999999999999999999999999999999",
			salvageValue:    "5000",
			usefulLife:      "60",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := testDeps()
			var createdRecord *AssetRecord
			deps.CreateAsset = func(ctx context.Context, asset *AssetRecord) error {
				createdRecord = asset
				return nil
			}

			v := NewAddAction(deps)

			form := url.Values{}
			form.Set("name", "Test Asset")
			form.Set("acquisition_cost", tt.acquisitionCost)
			form.Set("salvage_value", tt.salvageValue)
			form.Set("useful_life_months", tt.usefulLife)

			req := httptest.NewRequest(http.MethodPost, "/assets/add", strings.NewReader(form.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			viewCtx := &view.ViewContext{Request: req}

			result := v.Handle(ctxWithPerms("asset:create"), viewCtx)

			// strconv.ParseFloat/Atoi silently return 0 for invalid input,
			// so the handler should succeed (create with zero values).
			// This test documents that behavior.
			if result.StatusCode != http.StatusOK {
				t.Errorf("status = %d, want %d (handler should tolerate bad numeric input)",
					result.StatusCode, http.StatusOK)
			}
			if createdRecord == nil {
				t.Fatal("expected CreateAsset to be called")
			}
		})
	}
}

func TestNewAddAction_CreateAssetError(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	deps.CreateAsset = func(ctx context.Context, asset *AssetRecord) error {
		return errors.New("database connection refused")
	}

	v := NewAddAction(deps)

	form := url.Values{}
	form.Set("name", "Test Asset")
	req := httptest.NewRequest(http.MethodPost, "/assets/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:create"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Failed to create asset" {
		t.Errorf("HX-Error-Message = %q, want %q",
			result.Headers["HX-Error-Message"], "Failed to create asset")
	}
}

func TestNewAddAction_GetReturnsForm(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewAddAction(deps)

	req := httptest.NewRequest(http.MethodGet, "/assets/add", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:create"), viewCtx)

	if result.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusOK)
	}
	if result.Template != "asset-drawer-form" {
		t.Errorf("template = %q, want %q", result.Template, "asset-drawer-form")
	}
}

// ----- NewEditAction tests -----

func TestNewEditAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewEditAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/123/edit", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxNoPerms(), viewCtx)

	if result.Error == nil {
		t.Fatal("expected error for permission denied")
	}
	if result.Error.Error() != "permission denied" {
		t.Errorf("error = %q, want %q", result.Error.Error(), "permission denied")
	}
}

func TestNewEditAction_MissingName(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewEditAction(deps)

	form := url.Values{}
	form.Set("acquisition_cost", "10000")

	req := httptest.NewRequest(http.MethodPost, "/assets/123/edit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "123")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:update"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestNewEditAction_ReadAssetError(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	deps.ReadAsset = func(ctx context.Context, id string) (*AssetRecord, error) {
		return nil, errors.New("asset not found")
	}

	v := NewEditAction(deps)

	req := httptest.NewRequest(http.MethodGet, "/assets/999/edit", nil)
	req.SetPathValue("id", "999")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:update"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != "Failed to read asset" {
		t.Errorf("HX-Error-Message = %q, want %q",
			result.Headers["HX-Error-Message"], "Failed to read asset")
	}
}

func TestNewEditAction_UpdateAssetError(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	deps.UpdateAsset = func(ctx context.Context, asset *AssetRecord) error {
		return errors.New("update failed")
	}

	v := NewEditAction(deps)

	form := url.Values{}
	form.Set("name", "Updated Asset")

	req := httptest.NewRequest(http.MethodPost, "/assets/123/edit", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetPathValue("id", "123")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:update"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ----- NewDeleteAction tests -----

func TestNewDeleteAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewDeleteAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/delete?id=123", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxNoPerms(), viewCtx)

	if result.Error == nil {
		t.Fatal("expected error for permission denied")
	}
}

func TestNewDeleteAction_MissingID(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewDeleteAction(deps)

	// No id in query or form
	req := httptest.NewRequest(http.MethodPost, "/assets/delete", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:delete"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
	if result.Headers["HX-Error-Message"] != deps.Labels.Actions.IDRequired {
		t.Errorf("HX-Error-Message = %q, want %q",
			result.Headers["HX-Error-Message"], deps.Labels.Actions.IDRequired)
	}
}

func TestNewDeleteAction_EmptyIDQueryParam(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewDeleteAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/delete?id=", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:delete"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d for empty id", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestNewDeleteAction_DeleteError(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	deps.DeleteAsset = func(ctx context.Context, id string) error {
		return errors.New("delete failed: foreign key constraint")
	}

	v := NewDeleteAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/delete?id=123", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:delete"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

func TestNewDeleteAction_IDFromFormBody(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	var deletedID string
	deps.DeleteAsset = func(ctx context.Context, id string) error {
		deletedID = id
		return nil
	}

	v := NewDeleteAction(deps)

	form := url.Values{}
	form.Set("id", "asset-456")
	req := httptest.NewRequest(http.MethodPost, "/assets/delete", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:delete"), viewCtx)

	if result.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusOK)
	}
	if deletedID != "asset-456" {
		t.Errorf("deleted ID = %q, want %q", deletedID, "asset-456")
	}
}

// ----- NewSetStatusAction tests -----

func TestNewSetStatusAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewSetStatusAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/set-status?id=123&status=active", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxNoPerms(), viewCtx)

	if result.Error == nil {
		t.Fatal("expected error for permission denied")
	}
}

func TestNewSetStatusAction_InvalidStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		status string
	}{
		{name: "empty status", status: ""},
		{name: "unknown status", status: "pending"},
		{name: "uppercase ACTIVE", status: "ACTIVE"},
		{name: "sql injection", status: "active'; DROP TABLE --"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			deps := testDeps()
			v := NewSetStatusAction(deps)

			req := httptest.NewRequest(http.MethodPost,
				"/assets/set-status?id=123&status="+url.QueryEscape(tt.status), nil)
			viewCtx := &view.ViewContext{Request: req}

			result := v.Handle(ctxWithPerms("asset:update"), viewCtx)

			if result.StatusCode != http.StatusUnprocessableEntity {
				t.Errorf("status = %d, want %d for invalid status %q",
					result.StatusCode, http.StatusUnprocessableEntity, tt.status)
			}
		})
	}
}

func TestNewSetStatusAction_MissingID(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewSetStatusAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/set-status?status=active", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:update"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ----- NewBulkDeleteAction tests -----

func TestNewBulkDeleteAction_PermissionDenied(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewBulkDeleteAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/bulk-delete", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxNoPerms(), viewCtx)

	if result.Error == nil {
		t.Fatal("expected error for permission denied")
	}
}

func TestNewBulkDeleteAction_NoIDs(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewBulkDeleteAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/bulk-delete", strings.NewReader(""))
	req.Header.Set("Content-Type", "multipart/form-data; boundary=----")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:delete"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ----- NewBulkSetStatusAction tests -----

func TestNewBulkSetStatusAction_InvalidTargetStatus(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	v := NewBulkSetStatusAction(deps)

	form := url.Values{}
	form.Set("id", "asset-1")
	form.Set("target_status", "bogus")

	req := httptest.NewRequest(http.MethodPost, "/assets/bulk-set-status", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:update"), viewCtx)

	if result.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("status = %d, want %d", result.StatusCode, http.StatusUnprocessableEntity)
	}
}

// ----- NilDeps edge cases -----

func TestNewAddAction_NilCreateAsset(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	deps.CreateAsset = nil // no create wired

	v := NewAddAction(deps)

	form := url.Values{}
	form.Set("name", "Test Asset")
	req := httptest.NewRequest(http.MethodPost, "/assets/add", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:create"), viewCtx)

	// Should succeed even without CreateAsset wired (mock mode)
	if result.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d (nil CreateAsset should not error)", result.StatusCode, http.StatusOK)
	}
}

func TestNewDeleteAction_NilDeleteAsset(t *testing.T) {
	t.Parallel()

	deps := testDeps()
	deps.DeleteAsset = nil // no delete wired

	v := NewDeleteAction(deps)

	req := httptest.NewRequest(http.MethodPost, "/assets/delete?id=123", nil)
	viewCtx := &view.ViewContext{Request: req}

	result := v.Handle(ctxWithPerms("asset:delete"), viewCtx)

	// Should succeed in mock mode
	if result.StatusCode != http.StatusOK {
		t.Errorf("status = %d, want %d (nil DeleteAsset should not error)", result.StatusCode, http.StatusOK)
	}
}
