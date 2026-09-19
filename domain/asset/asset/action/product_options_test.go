package action

import (
	"context"
	assetform "github.com/erniealice/fycha-golang/domain/asset/asset/form"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestAssetProductDrawerSelection(t *testing.T) {
	deps := testDeps()
	deps.LoadProductOptions = func(context.Context) ([]types.SelectOption, error) {
		return []types.SelectOption{{Value: "product", Label: "Sample Product"}}, nil
	}
	deps.ReadAsset = func(context.Context, string) (*assetform.Record, error) {
		return &assetform.Record{ID: "asset", ProductID: "product"}, nil
	}
	req := httptest.NewRequest(http.MethodGet, "/action/asset/edit", nil)
	req.SetPathValue("id", "asset")
	result := NewEditAction(deps).Handle(ctxWithPerms("asset:update"), &view.ViewContext{Request: req})
	d := result.Data.(*assetform.Data)
	if !d.ShowProduct || !d.ProductLocked || d.ProductID != "product" || len(d.ProductOptions) != 1 {
		t.Fatalf("wrong product form: %+v", d)
	}
	saved := ""
	deps.CreateAsset = func(_ context.Context, r *assetform.Record) error { saved = r.ProductID; return nil }
	values := url.Values{"name": {"Sample Asset"}, "product_id": {"product"}, "acquisition_cost": {"100"}}
	req = httptest.NewRequest(http.MethodPost, "/action/asset/add", strings.NewReader(values.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	NewAddAction(deps).Handle(ctxWithPerms("asset:create"), &view.ViewContext{Request: req})
	if saved != "product" {
		t.Fatalf("product not forwarded: %s", saved)
	}
}
