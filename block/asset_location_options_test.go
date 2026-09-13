package block

import (
	"context"
	"errors"
	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	locationpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/entity/location"
	assetform "github.com/erniealice/fycha-golang/domain/asset/asset/form"
	"testing"
)

func TestAssetLocationListRequestHasNoStatusFilter(t *testing.T) {
	ctx := context.WithValue(context.Background(), struct{}{}, "workspace scope marker")
	calls := 0
	uc := &UseCases{GetLocationListPageData: func(got context.Context, req *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error) {
		calls++
		if got != ctx {
			t.Fatal("request context lost")
		}
		if req.GetFilters() != nil {
			t.Fatal("location picker must not filter status")
		}
		if req.GetPagination().GetOffset().GetPage() != int32(calls) {
			t.Fatal("pagination not advanced")
		}
		if calls == 1 {
			return &locationpb.GetLocationListPageDataResponse{LocationList: []*locationpb.Location{{Id: "active-id", Name: "A Active", Active: true}}, Pagination: &commonpb.PaginationResponse{HasNext: true}}, nil
		}
		return &locationpb.GetLocationListPageDataResponse{LocationList: []*locationpb.Location{{Id: "inactive-id", Name: "Z Inactive", Active: false}}}, nil
	}}
	options, err := assetLocationOptionsLoader(uc)(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if calls != 2 || len(options) != 2 || options[0].Value != "active-id" || options[1].Value != "inactive-id" || options[1].Label != "Z Inactive" {
		t.Fatalf("wrong options: %+v", options)
	}
	failure := errors.New("location list failed")
	uc.GetLocationListPageData = func(context.Context, *locationpb.GetLocationListPageDataRequest) (*locationpb.GetLocationListPageDataResponse, error) {
		return nil, failure
	}
	if _, err := assetLocationOptionsLoader(uc)(ctx); !errors.Is(err, failure) {
		t.Fatalf("failure hidden: %v", err)
	}
}

func TestAssetLocationRecordMapping(t *testing.T) {
	record := &assetform.Record{ID: "asset-id", LocationID: "inactive-location-id"}
	proto := recordToAsset(record)
	if proto.GetLocationId() != record.LocationID {
		t.Fatal("location ID lost before persistence")
	}
	if assetToRecord(proto).LocationID != record.LocationID {
		t.Fatal("location ID lost after reading")
	}
}
