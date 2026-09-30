// Package cptest is an in-memory fake of the charge policy ports plus request
// helpers, shared by the view tests of the charge policy pages. Test-only.
package cptest

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	componentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_component"
	postingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_posting"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
)

// Coded is a fake use-case refusal: it exposes ErrorCode() like espyna's Error.
type Coded struct{ code string }

func (e *Coded) Error() string     { return "charge_policy: " + e.code }
func (e *Coded) ErrorCode() string { return e.code }

// Errors returned by the fake (classified by cp.ErrorKind through ErrorCode()).
var (
	ErrNotFound     error = &Coded{"not_found"}
	ErrNotDraft     error = &Coded{"not_draft"}
	ErrSelfApproval error = &Coded{"self_approval"}
	ErrChecklist    error = &Coded{"checklist_failed"}
	// ErrPermissionDenied is the strict-gate denial (actiongate `permission_denied`).
	ErrPermissionDenied error = &Coded{"permission_denied"}
)

// Fake is the in-memory store behind the ports.
type Fake struct {
	Policies   []*policypb.ChargePolicy
	Versions   map[string][]*versionpb.ChargePolicyVersion // by policy id
	Components map[string][]*componentpb.ChargePolicyComponent
	Postings   map[string][]*postingpb.ChargePolicyPosting
	Accounts   []*accountpb.Account
	InUse      map[string]bool

	Validation  *versionpb.ChargePolicyApprovalChecklist
	ApproveErr  error
	Calls       map[string]int
	LastVersion *versionpb.ChargePolicyVersion
}

// New returns an empty fake.
func New() *Fake {
	return &Fake{
		Versions:   map[string][]*versionpb.ChargePolicyVersion{},
		Components: map[string][]*componentpb.ChargePolicyComponent{},
		Postings:   map[string][]*postingpb.ChargePolicyPosting{},
		InUse:      map[string]bool{},
		Calls:      map[string]int{},
	}
}

// MutationCalls counts every write call made through the ports.
func (f *Fake) MutationCalls() int {
	n := 0
	for k, v := range f.Calls {
		switch k {
		case "CreatePolicy", "UpdatePolicy", "RetirePolicy", "DeletePolicy", "CreateDraftVersion",
			"UpdateVersion", "DeleteVersion", "CreateComponent", "UpdateComponent", "DeleteComponent",
			"CreatePosting", "UpdatePosting", "DeletePosting", "Approve":
			n += v
		}
	}
	return n
}

// AddPolicy registers a policy and its versions.
func (f *Fake) AddPolicy(id, code, name string, retired bool, versions ...*versionpb.ChargePolicyVersion) *policypb.ChargePolicy {
	st := enumspb.ChargePolicyStatus_CHARGE_POLICY_STATUS_ACTIVE
	if retired {
		st = enumspb.ChargePolicyStatus_CHARGE_POLICY_STATUS_RETIRED
	}
	p := &policypb.ChargePolicy{Id: id, Code: code, Name: name, Status: st, Active: true}
	f.Policies = append(f.Policies, p)
	for _, v := range versions {
		v.ChargePolicyId = id
		f.Versions[id] = append(f.Versions[id], v)
	}
	return p
}

// Version builds a version row.
func Version(id string, number int32, status enumspb.ChargePolicyVersionStatus, preparedBy string) *versionpb.ChargePolicyVersion {
	pb := preparedBy
	return &versionpb.ChargePolicyVersion{Id: id, VersionNumber: number, Status: status, PreparedBy: &pb, Active: true}
}

// Draft / Approved shortcuts.
func Draft(id string, number int32, preparedBy string) *versionpb.ChargePolicyVersion {
	return Version(id, number, enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_DRAFT, preparedBy)
}

func Approved(id string, number int32) *versionpb.ChargePolicyVersion {
	return Version(id, number, enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_APPROVED, "u-prep")
}

func (f *Fake) findVersion(vid string) *versionpb.ChargePolicyVersion {
	for _, vs := range f.Versions {
		for _, v := range vs {
			if v.GetId() == vid {
				return v
			}
		}
	}
	return nil
}

// UseCases exposes the fake as the proto-closure set the views consume.
func (f *Fake) UseCases() *cp.UseCases {
	c := func(name string) { f.Calls[name]++ }
	return &cp.UseCases{
		GetChargePolicyListPageData: func(ctx context.Context, req *policypb.GetChargePolicyListPageDataRequest) (*policypb.GetChargePolicyListPageDataResponse, error) {
			c("ListPolicies")
			return &policypb.GetChargePolicyListPageDataResponse{ChargePolicyList: f.Policies, Success: true}, nil
		},
		ReadChargePolicy: func(ctx context.Context, req *policypb.ReadChargePolicyRequest) (*policypb.ReadChargePolicyResponse, error) {
			c("ReadPolicy")
			id := req.GetData().GetId()
			for _, p := range f.Policies {
				if p.GetId() == id {
					return &policypb.ReadChargePolicyResponse{Data: []*policypb.ChargePolicy{p}, Versions: f.Versions[id], Success: true}, nil
				}
			}
			return nil, ErrNotFound
		},
		CreateChargePolicy: func(ctx context.Context, req *policypb.CreateChargePolicyRequest) (*policypb.CreateChargePolicyResponse, error) {
			c("CreatePolicy")
			p := &policypb.ChargePolicy{Id: "p-new", Code: req.GetData().GetCode(), Name: req.GetData().GetName()}
			return &policypb.CreateChargePolicyResponse{Data: []*policypb.ChargePolicy{p}, Draft: Draft("v-new", 1, "u-actor"), Success: true}, nil
		},
		UpdateChargePolicy: func(ctx context.Context, req *policypb.UpdateChargePolicyRequest) (*policypb.UpdateChargePolicyResponse, error) {
			c("UpdatePolicy")
			return &policypb.UpdateChargePolicyResponse{Success: true}, nil
		},
		RetireChargePolicy: func(ctx context.Context, req *policypb.RetireChargePolicyRequest) (*policypb.RetireChargePolicyResponse, error) {
			c("RetirePolicy")
			return &policypb.RetireChargePolicyResponse{Success: true}, nil
		},
		DeleteChargePolicy: func(ctx context.Context, req *policypb.DeleteChargePolicyRequest) (*policypb.DeleteChargePolicyResponse, error) {
			c("DeletePolicy")
			return &policypb.DeleteChargePolicyResponse{Success: true}, nil
		},
		GetChargePolicyInUseIds: func(ctx context.Context, req *policypb.GetChargePolicyInUseIdsRequest) (*policypb.GetChargePolicyInUseIdsResponse, error) {
			c("InUseIDs")
			resp := &policypb.GetChargePolicyInUseIdsResponse{Success: true}
			for id, in := range f.InUse {
				if in {
					resp.InUseChargePolicyIds = append(resp.InUseChargePolicyIds, id)
				}
			}
			return resp, nil
		},
		ListChargePolicyVersions: func(ctx context.Context, req *versionpb.ListChargePolicyVersionsRequest) (*versionpb.ListChargePolicyVersionsResponse, error) {
			c("ListVersions")
			return &versionpb.ListChargePolicyVersionsResponse{Data: f.Versions[req.GetChargePolicyId()], Success: true}, nil
		},
		CreateDraftChargePolicyVersion: func(ctx context.Context, req *versionpb.CreateDraftChargePolicyVersionRequest) (*versionpb.CreateDraftChargePolicyVersionResponse, error) {
			c("CreateDraftVersion")
			return &versionpb.CreateDraftChargePolicyVersionResponse{Draft: Draft("v-clone", 9, "u-actor"), Success: true}, nil
		},
		ReadChargePolicyVersion: func(ctx context.Context, req *versionpb.ReadChargePolicyVersionRequest) (*versionpb.ReadChargePolicyVersionResponse, error) {
			c("ReadVersion")
			vid := req.GetData().GetId()
			v := f.findVersion(vid)
			if v == nil {
				return nil, ErrNotFound
			}
			return &versionpb.ReadChargePolicyVersionResponse{Data: []*versionpb.ChargePolicyVersion{v}, Components: f.Components[vid], Postings: f.Postings[vid], Success: true}, nil
		},
		UpdateChargePolicyVersion: func(ctx context.Context, req *versionpb.UpdateChargePolicyVersionRequest) (*versionpb.UpdateChargePolicyVersionResponse, error) {
			c("UpdateVersion")
			f.LastVersion = req.GetData()
			return &versionpb.UpdateChargePolicyVersionResponse{Success: true}, nil
		},
		DeleteChargePolicyVersion: func(ctx context.Context, req *versionpb.DeleteChargePolicyVersionRequest) (*versionpb.DeleteChargePolicyVersionResponse, error) {
			c("DeleteVersion")
			return &versionpb.DeleteChargePolicyVersionResponse{Success: true}, nil
		},
		CreateChargePolicyComponent: func(ctx context.Context, req *componentpb.CreateChargePolicyComponentRequest) (*componentpb.CreateChargePolicyComponentResponse, error) {
			c("CreateComponent")
			return &componentpb.CreateChargePolicyComponentResponse{Success: true}, nil
		},
		UpdateChargePolicyComponent: func(ctx context.Context, req *componentpb.UpdateChargePolicyComponentRequest) (*componentpb.UpdateChargePolicyComponentResponse, error) {
			c("UpdateComponent")
			return &componentpb.UpdateChargePolicyComponentResponse{Success: true}, nil
		},
		DeleteChargePolicyComponent: func(ctx context.Context, req *componentpb.DeleteChargePolicyComponentRequest) (*componentpb.DeleteChargePolicyComponentResponse, error) {
			c("DeleteComponent")
			return &componentpb.DeleteChargePolicyComponentResponse{Success: true}, nil
		},
		CreateChargePolicyPosting: func(ctx context.Context, req *postingpb.CreateChargePolicyPostingRequest) (*postingpb.CreateChargePolicyPostingResponse, error) {
			c("CreatePosting")
			return &postingpb.CreateChargePolicyPostingResponse{Success: true}, nil
		},
		UpdateChargePolicyPosting: func(ctx context.Context, req *postingpb.UpdateChargePolicyPostingRequest) (*postingpb.UpdateChargePolicyPostingResponse, error) {
			c("UpdatePosting")
			return &postingpb.UpdateChargePolicyPostingResponse{Success: true}, nil
		},
		DeleteChargePolicyPosting: func(ctx context.Context, req *postingpb.DeleteChargePolicyPostingRequest) (*postingpb.DeleteChargePolicyPostingResponse, error) {
			c("DeletePosting")
			return &postingpb.DeleteChargePolicyPostingResponse{Success: true}, nil
		},
		ValidateChargePolicyVersionForApproval: func(ctx context.Context, req *versionpb.ValidateChargePolicyVersionForApprovalRequest) (*versionpb.ValidateChargePolicyVersionForApprovalResponse, error) {
			c("Validate")
			return &versionpb.ValidateChargePolicyVersionForApprovalResponse{Checklist: f.Validation, Success: true}, nil
		},
		ApproveChargePolicyVersion: func(ctx context.Context, req *versionpb.ApproveChargePolicyVersionRequest) (*versionpb.ApproveChargePolicyVersionResponse, error) {
			c("Approve")
			if f.ApproveErr != nil {
				return nil, f.ApproveErr
			}
			return &versionpb.ApproveChargePolicyVersionResponse{Checklist: f.Validation, Success: true}, nil
		},
		GetAccountListPageData: func(ctx context.Context, req *accountpb.GetAccountListPageDataRequest) (*accountpb.GetAccountListPageDataResponse, error) {
			c("ListAccounts")
			return &accountpb.GetAccountListPageDataResponse{AccountList: f.Accounts, Success: true}, nil
		},
	}
}

// Common returns CommonLabels with the message templates the views format.
func Common() pyeza.CommonLabels {
	var cl pyeza.CommonLabels
	cl.Errors.MissingPermission = "Missing permission: %s"
	cl.Errors.PermissionDenied = "Permission denied"
	return cl
}

// Ctx returns a context carrying the given permission codes.
func Ctx(codes ...string) context.Context {
	return view.WithUserPermissions(context.Background(), types.NewUserPermissions(codes))
}

// Request builds a view context for a request with path values.
func Request(method, target string, form string, pathValues ...string) *view.ViewContext {
	var req *http.Request
	if form != "" {
		req = httptest.NewRequest(method, target, strings.NewReader(form))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req = httptest.NewRequest(method, target, nil)
	}
	for i := 0; i+1 < len(pathValues); i += 2 {
		req.SetPathValue(pathValues[i], pathValues[i+1])
	}
	return &view.ViewContext{Request: req, CurrentPath: req.URL.Path}
}

// Errf is a tiny helper for tests that need an ad-hoc error.
func Errf(format string, a ...any) error { return fmt.Errorf(format, a...) }
