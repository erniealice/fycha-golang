package action

import (
	"context"
	"errors"
	assetform "github.com/erniealice/fycha-golang/domain/asset/asset/form"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAssetLocationOptionsAllStatuses(t *testing.T) {
	deps := testDeps()
	deps.LoadLocationOptions = func(context.Context) ([]types.SelectOption, error) {
		return []types.SelectOption{{Value: "active", Label: "Active location"}, {Value: "inactive", Label: "Inactive location"}}, nil
	}
	deps.ReadAsset = func(context.Context, string) (*assetform.Record, error) {
		return &assetform.Record{ID: "asset", LocationID: "inactive"}, nil
	}
	for _, edit := range []bool{false, true} {
		req := httptest.NewRequest(http.MethodGet, "/action/asset/add", nil)
		req.SetPathValue("id", "asset")
		handler := NewAddAction(deps)
		if edit {
			handler = NewEditAction(deps)
		}
		result := handler.Handle(ctxWithPerms("asset:create", "asset:update"), &view.ViewContext{Request: req})
		if result.Error != nil {
			t.Fatal(result.Error)
		}
		data := result.Data.(*assetform.Data)
		if len(data.LocationOptions) != 2 {
			t.Fatal("locations lost")
		}
		if edit && data.LocationID != "inactive" {
			t.Fatal("existing inactive location lost")
		}
	}
}

func TestAssetLocationLoadFailure(t *testing.T) {
	deps := testDeps()
	failure := errors.New("lookup failed")
	deps.LoadLocationOptions = func(context.Context) ([]types.SelectOption, error) { return nil, failure }
	result := NewAddAction(deps).Handle(ctxWithPerms("asset:create"), &view.ViewContext{Request: httptest.NewRequest(http.MethodGet, "/action/asset/add", nil)})
	if !errors.Is(result.Error, failure) {
		t.Fatal("lookup failure hidden")
	}
}

func TestAssetSaveSelectedLocation(t *testing.T) {
	for _, edit := range []bool{false, true} {
		deps := testDeps()
		saved := ""
		save := func(_ context.Context, r *assetform.Record) error { saved = r.LocationID; return nil }
		deps.CreateAsset = save
		deps.UpdateAsset = save
		data := url.Values{"name": {"Asset test"}, "location_id": {"chosen-location-id"}, "acquisition_cost": {"100"}}
		req := httptest.NewRequest(http.MethodPost, "/action/asset/add", strings.NewReader(data.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.SetPathValue("id", "asset")
		handler := NewAddAction(deps)
		if edit {
			handler = NewEditAction(deps)
		}
		result := handler.Handle(ctxWithPerms("asset:create", "asset:update"), &view.ViewContext{Request: req})
		if result.Error != nil {
			t.Fatal(result.Error)
		}
		if saved != "chosen-location-id" {
			t.Fatalf("saved location_id=%q", saved)
		}
	}
}
