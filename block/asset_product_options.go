package block

import (
	"context"
	"sort"
	"strings"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	productpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/product/product"
	"github.com/erniealice/pyeza-golang/types"
)

func assetProductOptionsLoader(uc *UseCases) func(context.Context) ([]types.SelectOption, error) {
	if uc == nil || uc.ListAssetProducts == nil {
		return nil
	}
	return func(ctx context.Context) ([]types.SelectOption, error) {
		var options []types.SelectOption
		seen := map[string]bool{}
		for _, active := range []bool{true, false} {
			for page := int32(1); ; page++ {
				resp, err := uc.ListAssetProducts(ctx, &productpb.ListProductsRequest{
					Filters:    &commonpb.FilterRequest{Filters: []*commonpb.TypedFilter{{Field: "active", FilterType: &commonpb.TypedFilter_BooleanFilter{BooleanFilter: &commonpb.BooleanFilter{Value: active}}}}},
					Pagination: &commonpb.PaginationRequest{Limit: 100, Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}}},
				})
				if err != nil {
					return nil, err
				}
				rows := resp.GetData()
				for _, p := range rows {
					if p == nil || seen[p.GetId()] {
						continue
					}
					seen[p.GetId()] = true
					options = append(options, types.SelectOption{Value: p.GetId(), Label: p.GetName()})
				}
				if len(rows) < 100 {
					break
				}
			}
		}
		sort.SliceStable(options, func(i, j int) bool { return strings.ToLower(options[i].Label) < strings.ToLower(options[j].Label) })
		return options, nil
	}
}
