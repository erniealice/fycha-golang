package block

import (
	"context"
	"errors"
	"fmt"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	"testing"
)

func TestAssetProductOptionsPaginationAndErrors(t *testing.T) {
	type scopeKey struct{}
	ctx := context.WithValue(context.Background(), scopeKey{}, "scope")
	var calls []string
	uc := &UseCases{ListAssetProducts: func(got context.Context, req *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
		if got.Value(scopeKey{}) != "scope" {
			t.Fatal("context lost")
		}
		active := req.Filters.Filters[0].GetBooleanFilter().Value
		page := req.Pagination.GetOffset().Page
		calls = append(calls, fmt.Sprintf("%v/%d", active, page))
		if req.Pagination.Limit != 100 {
			t.Fatalf("limit=%d", req.Pagination.Limit)
		}
		var rows []*productpb.Product
		if active && page == 1 {
			for i := 0; i < 100; i++ {
				rows = append(rows, &productpb.Product{Id: fmt.Sprintf("p%d", i), Name: fmt.Sprintf("Product %03d", i)})
			}
		}
		if active && page == 2 {
			rows = []*productpb.Product{{Id: "last", Name: "Zulu"}}
		}
		if !active {
			rows = []*productpb.Product{nil, {Id: "p0", Name: "Duplicate"}, {Id: "inactive", Name: "alpha"}}
		}
		return &productpb.ListProductsResponse{Data: rows}, nil
	}}
	options, err := assetProductOptionsLoader(uc)(ctx)
	if err != nil || len(options) != 102 || options[0].Value != "inactive" {
		t.Fatalf("options=%d err=%v", len(options), err)
	}
	if fmt.Sprint(calls) != "[true/1 true/2 false/1]" {
		t.Fatalf("calls=%v", calls)
	}
	failure := errors.New("scoped list denied")
	uc.ListAssetProducts = func(context.Context, *productpb.ListProductsRequest) (*productpb.ListProductsResponse, error) {
		return nil, failure
	}
	if _, err := assetProductOptionsLoader(uc)(ctx); !errors.Is(err, failure) {
		t.Fatalf("error=%v", err)
	}
	if assetProductOptionsLoader(nil) != nil || assetProductOptionsLoader(&UseCases{}) != nil {
		t.Fatal("missing dependency did not stay optional")
	}
}
