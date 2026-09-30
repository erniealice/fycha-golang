package charge_policy

import (
	"context"
	"errors"
	"testing"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
)

// C11: lists shown in views walk every page (no silent 100-row cap).
func TestListPoliciesPagesThroughEveryRow(t *testing.T) {
	calls := 0
	u := &UseCases{GetChargePolicyListPageData: func(ctx context.Context, req *policypb.GetChargePolicyListPageDataRequest) (*policypb.GetChargePolicyListPageDataResponse, error) {
		calls++
		page := req.GetPagination().GetOffset().GetPage()
		resp := &policypb.GetChargePolicyListPageDataResponse{
			ChargePolicyList: []*policypb.ChargePolicy{{Id: string(rune('a' + page))}},
			Pagination:       &commonpb.PaginationResponse{HasNext: page < 3},
		}
		return resp, nil
	}}
	got, err := u.ListPolicies(context.Background())
	if err != nil || len(got) != 3 || calls != 3 {
		t.Fatalf("got %d rows in %d calls, err %v", len(got), calls, err)
	}
}

type checklistErr struct{}

func (checklistErr) Error() string     { return "charge_policy: checklist_failed" }
func (checklistErr) ErrorCode() string { return "checklist_failed" }
func (checklistErr) Checklist() *versionpb.ChargePolicyApprovalChecklist {
	return &versionpb.ChargePolicyApprovalChecklist{Supported: true}
}

// A failed approval hands the per-item report back through the error's Checklist().
func TestApproveReturnsChecklistFromError(t *testing.T) {
	u := &UseCases{ApproveChargePolicyVersion: func(context.Context, *versionpb.ApproveChargePolicyVersionRequest) (*versionpb.ApproveChargePolicyVersionResponse, error) {
		return nil, checklistErr{}
	}}
	cl, err := u.Approve(context.Background(), "v1")
	if err == nil || cl == nil || !cl.GetSupported() {
		t.Fatalf("checklist %v err %v", cl, err)
	}
	if ErrorKind(err) != ErrKindChecklistFailed {
		t.Errorf("kind %q", ErrorKind(err))
	}
	u.ApproveChargePolicyVersion = func(context.Context, *versionpb.ApproveChargePolicyVersionRequest) (*versionpb.ApproveChargePolicyVersionResponse, error) {
		return nil, errors.New("plain")
	}
	if cl, _ := u.Approve(context.Background(), "v1"); cl != nil {
		t.Error("uncoded error must not invent a checklist")
	}
}
