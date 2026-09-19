package block

import (
	"context"
	"github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/espyna-golang/consumer/compose"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	asset "github.com/erniealice/fycha-golang/domain/asset/asset"
	assetform "github.com/erniealice/fycha-golang/domain/asset/asset/form"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"
	"net/http/httptest"
	"reflect"
	"testing"
)

type commercialRegistrar struct{ get map[string]view.View }

func (r *commercialRegistrar) GET(path string, v view.View, _ ...string) { r.get[path] = v }
func (r *commercialRegistrar) POST(string, view.View, ...string)         {}

func TestAssetProductSelectionOptions(t *testing.T) {
	for _, enabled := range []bool{false, true} {
		for _, available := range []bool{false, true} {
			raw := &consumer.UseCases{}
			if available {
				// Allocate the exported aggregate's concrete use-case chain without importing
				// Espyna internals. Execute is not called: binding presence is the contract.
				v := reflect.ValueOf(raw).Elem()
				for _, name := range []string{"Product", "Product", "ListProducts"} {
					f := v.FieldByName(name)
					f.Set(reflect.New(f.Type().Elem()))
					v = f.Elem()
				}
			}
			cfg := engineConfig{}
			WithAssetProductSelection(enabled)(&cfg)
			adapted := &UseCases{}
			bindAssetProducts(adapted, raw, cfg.assetProductSelection)
			want := enabled && available
			if (adapted.ListAssetProducts != nil) != want {
				t.Fatalf("enabled=%v available=%v binding mismatch", enabled, available)
			}
			// Exercise the downstream real mount with a typed use-case stub; scope and
			// permissions are not replaced by the capability.
			if want {
				adapted.ListAssetProducts = func(context.Context, *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
					return &productpb.ListProductsResponse{Data: []*productpb.Product{{Id: "product-1", Name: "Sample Product"}}}, nil
				}
			}
			for _, vertical := range []string{"general", "leasing"} {
				r := &commercialRegistrar{get: map[string]view.View{}}
				u := AssetUnit(adapted, &Infra{})
				if err := u.Mount(&compose.MountContext{Routes: r}); err != nil {
					t.Fatal(err)
				}
				path := asset.DefaultRoutes().AddURL
				ctx := view.WithUserPermissions(context.Background(), types.NewUserPermissions([]string{"asset:create"}))
				result := r.get[path].Handle(ctx, &view.ViewContext{Request: httptest.NewRequest("GET", path, nil), BusinessType: vertical})
				d, ok := result.Data.(*assetform.Data)
				if !ok || d.ShowProduct != want {
					t.Fatalf("vertical=%s want=%v result=%+v", vertical, want, result)
				}
				if want && len(d.ProductOptions) != 1 {
					t.Fatalf("options=%+v", d.ProductOptions)
				}
			}
		}
	}
	bindAssetProducts(nil, nil, true)
	bindAssetProducts(&UseCases{}, nil, true)
}
