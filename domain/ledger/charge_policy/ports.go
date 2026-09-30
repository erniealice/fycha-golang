package charge_policy

import (
	"context"
	"errors"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	componentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_component"
	postingpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_posting"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
)

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `charge_policy.errors.<code>` (Labels.ErrorMessage). The codes are
// the espyna charge_policy use-case codes (progress-f2-chargepolicy.md).
const (
	ErrKindNone                   = ""
	ErrKindValidation             = "validation"
	ErrKindNotFound               = "not_found"
	ErrKindNotDraft               = "not_draft"
	ErrKindDraftExists            = "draft_exists"
	ErrKindNoApprovedVersion      = "no_approved_version"
	ErrKindRetired                = "retired"
	ErrKindSelfApproval           = "self_approval"
	ErrKindChecklistFailed        = "checklist_failed"
	ErrKindUnsupportedCombination = "unsupported_combination"
	ErrKindInUse                  = "in_use"
	ErrKindTransactionRequired    = "transaction_required"
	ErrKindTaxTreatmentNotFound   = "tax_treatment_not_found"
	ErrKindUnverifiable           = "unverifiable"
	// ErrKindPermissionDenied is the strict-gate denial (actiongate `permission_denied`);
	// views map it to the shared PermissionDenied message, not the generic error.
	ErrKindPermissionDenied = "permission_denied"
	ErrKindUnknown          = "unknown"
)

// ErrorKind classifies a use-case error through its `ErrorCode()` (no import of
// espyna). Unknown / uncoded errors classify as ErrKindUnknown.
func ErrorKind(err error) string {
	if err == nil {
		return ErrKindNone
	}
	var coded interface{ ErrorCode() string }
	if errors.As(err, &coded) {
		if c := coded.ErrorCode(); c != "" {
			return c
		}
	}
	return ErrKindUnknown
}

// PolicyDetail is a policy with all its versions (any status), newest first.
type PolicyDetail struct {
	Policy   *policypb.ChargePolicy
	Versions []*versionpb.ChargePolicyVersion
}

// VersionDetail is a version with its component and posting rows.
type VersionDetail struct {
	Version    *versionpb.ChargePolicyVersion
	Components []*componentpb.ChargePolicyComponent
	Postings   []*postingpb.ChargePolicyPosting
}

// UseCases is the set of espyna charge policy use-case closures the views
// consume, exactly the proto request/response shape of
// `uc.Ledger.ChargePolicy.<Field>.Execute` (fiscal_period precedent). The block
// binds them; a nil closure disables the feature that needs it and the view
// fails closed. Permission gates live in the views AND in the use cases.
type UseCases struct {
	CreateChargePolicy          func(context.Context, *policypb.CreateChargePolicyRequest) (*policypb.CreateChargePolicyResponse, error)
	ReadChargePolicy            func(context.Context, *policypb.ReadChargePolicyRequest) (*policypb.ReadChargePolicyResponse, error)
	UpdateChargePolicy          func(context.Context, *policypb.UpdateChargePolicyRequest) (*policypb.UpdateChargePolicyResponse, error)
	GetChargePolicyListPageData func(context.Context, *policypb.GetChargePolicyListPageDataRequest) (*policypb.GetChargePolicyListPageDataResponse, error)
	RetireChargePolicy          func(context.Context, *policypb.RetireChargePolicyRequest) (*policypb.RetireChargePolicyResponse, error)
	DeleteChargePolicy          func(context.Context, *policypb.DeleteChargePolicyRequest) (*policypb.DeleteChargePolicyResponse, error)
	GetChargePolicyInUseIds     func(context.Context, *policypb.GetChargePolicyInUseIdsRequest) (*policypb.GetChargePolicyInUseIdsResponse, error)

	CreateDraftChargePolicyVersion         func(context.Context, *versionpb.CreateDraftChargePolicyVersionRequest) (*versionpb.CreateDraftChargePolicyVersionResponse, error)
	ReadChargePolicyVersion                func(context.Context, *versionpb.ReadChargePolicyVersionRequest) (*versionpb.ReadChargePolicyVersionResponse, error)
	ListChargePolicyVersions               func(context.Context, *versionpb.ListChargePolicyVersionsRequest) (*versionpb.ListChargePolicyVersionsResponse, error)
	UpdateChargePolicyVersion              func(context.Context, *versionpb.UpdateChargePolicyVersionRequest) (*versionpb.UpdateChargePolicyVersionResponse, error)
	DeleteChargePolicyVersion              func(context.Context, *versionpb.DeleteChargePolicyVersionRequest) (*versionpb.DeleteChargePolicyVersionResponse, error)
	ValidateChargePolicyVersionForApproval func(context.Context, *versionpb.ValidateChargePolicyVersionForApprovalRequest) (*versionpb.ValidateChargePolicyVersionForApprovalResponse, error)
	ApproveChargePolicyVersion             func(context.Context, *versionpb.ApproveChargePolicyVersionRequest) (*versionpb.ApproveChargePolicyVersionResponse, error)

	CreateChargePolicyComponent func(context.Context, *componentpb.CreateChargePolicyComponentRequest) (*componentpb.CreateChargePolicyComponentResponse, error)
	UpdateChargePolicyComponent func(context.Context, *componentpb.UpdateChargePolicyComponentRequest) (*componentpb.UpdateChargePolicyComponentResponse, error)
	DeleteChargePolicyComponent func(context.Context, *componentpb.DeleteChargePolicyComponentRequest) (*componentpb.DeleteChargePolicyComponentResponse, error)
	CreateChargePolicyPosting   func(context.Context, *postingpb.CreateChargePolicyPostingRequest) (*postingpb.CreateChargePolicyPostingResponse, error)
	UpdateChargePolicyPosting   func(context.Context, *postingpb.UpdateChargePolicyPostingRequest) (*postingpb.UpdateChargePolicyPostingResponse, error)
	DeleteChargePolicyPosting   func(context.Context, *postingpb.DeleteChargePolicyPostingRequest) (*postingpb.DeleteChargePolicyPostingResponse, error)

	// Chart-of-accounts feeder for the posting drawer's account picker.
	GetAccountListPageData func(context.Context, *accountpb.GetAccountListPageDataRequest) (*accountpb.GetAccountListPageDataResponse, error)
}

// ---- view-side conveniences over the proto closures ------------------------
// These build presentation shapes from proto responses. Callers check the
// backing closure for nil first (the view fails closed).

const pageLimit = int32(100)

func pageRequest(page int32) *commonpb.PaginationRequest {
	return &commonpb.PaginationRequest{
		Limit:  pageLimit,
		Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}},
	}
}

// pageAll walks every page of a paginated list (no silent 100-row cap, C11).
func pageAll[T any](fetch func(*commonpb.PaginationRequest) ([]T, *commonpb.PaginationResponse, error)) ([]T, error) {
	var out []T
	for page := int32(1); page <= 1000; page++ {
		items, pag, err := fetch(pageRequest(page))
		if err != nil {
			return nil, err
		}
		out = append(out, items...)
		if pag != nil {
			if !pag.GetHasNext() {
				return out, nil
			}
		} else if int32(len(items)) < pageLimit {
			return out, nil
		}
	}
	return out, nil
}

// ListPolicies returns every policy of the workspace (paged through).
func (u *UseCases) ListPolicies(ctx context.Context) ([]*policypb.ChargePolicy, error) {
	return pageAll(func(p *commonpb.PaginationRequest) ([]*policypb.ChargePolicy, *commonpb.PaginationResponse, error) {
		resp, err := u.GetChargePolicyListPageData(ctx, &policypb.GetChargePolicyListPageDataRequest{Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetChargePolicyList(), resp.GetPagination(), nil
	})
}

// ReadPolicy returns the policy with all its versions.
func (u *UseCases) ReadPolicy(ctx context.Context, id string) (*PolicyDetail, error) {
	resp, err := u.ReadChargePolicy(ctx, &policypb.ReadChargePolicyRequest{Data: &policypb.ChargePolicy{Id: id}})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 {
		return nil, nil
	}
	return &PolicyDetail{Policy: resp.GetData()[0], Versions: resp.GetVersions()}, nil
}

// CreatePolicy creates the policy and its draft v1.
func (u *UseCases) CreatePolicy(ctx context.Context, code, name, description string) (*policypb.ChargePolicy, *versionpb.ChargePolicyVersion, error) {
	data := &policypb.ChargePolicy{Code: code, Name: name}
	if description != "" {
		data.Description = &description
	}
	resp, err := u.CreateChargePolicy(ctx, &policypb.CreateChargePolicyRequest{Data: data})
	if err != nil {
		return nil, nil, err
	}
	var pol *policypb.ChargePolicy
	if len(resp.GetData()) > 0 {
		pol = resp.GetData()[0]
	}
	return pol, resp.GetDraft(), nil
}

// UpdatePolicy updates name and description.
func (u *UseCases) UpdatePolicy(ctx context.Context, id, name, description string) error {
	data := &policypb.ChargePolicy{Id: id, Name: name}
	if description != "" {
		data.Description = &description
	}
	_, err := u.UpdateChargePolicy(ctx, &policypb.UpdateChargePolicyRequest{Data: data})
	return err
}

// RetirePolicy retires the policy.
func (u *UseCases) RetirePolicy(ctx context.Context, id string) error {
	_, err := u.RetireChargePolicy(ctx, &policypb.RetireChargePolicyRequest{ChargePolicyId: id})
	return err
}

// DeletePolicy deletes a never-approved, unreferenced policy.
func (u *UseCases) DeletePolicy(ctx context.Context, id string) error {
	_, err := u.DeleteChargePolicy(ctx, &policypb.DeleteChargePolicyRequest{Data: &policypb.ChargePolicy{Id: id}})
	return err
}

// InUseIDs returns the ids among ids that cannot be deleted.
func (u *UseCases) InUseIDs(ctx context.Context, ids []string) (map[string]bool, error) {
	resp, err := u.GetChargePolicyInUseIds(ctx, &policypb.GetChargePolicyInUseIdsRequest{ChargePolicyIds: ids})
	if err != nil {
		return nil, err
	}
	out := make(map[string]bool, len(resp.GetInUseChargePolicyIds()))
	for _, id := range resp.GetInUseChargePolicyIds() {
		out[id] = true
	}
	return out, nil
}

// ListVersions returns every version of a policy.
func (u *UseCases) ListVersions(ctx context.Context, policyID string) ([]*versionpb.ChargePolicyVersion, error) {
	return pageAll(func(p *commonpb.PaginationRequest) ([]*versionpb.ChargePolicyVersion, *commonpb.PaginationResponse, error) {
		resp, err := u.ListChargePolicyVersions(ctx, &versionpb.ListChargePolicyVersionsRequest{ChargePolicyId: &policyID, Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetData(), resp.GetPagination(), nil
	})
}

// CreateDraftVersion clones the latest approved version into a draft.
func (u *UseCases) CreateDraftVersion(ctx context.Context, policyID string) (*versionpb.ChargePolicyVersion, error) {
	resp, err := u.CreateDraftChargePolicyVersion(ctx, &versionpb.CreateDraftChargePolicyVersionRequest{ChargePolicyId: policyID})
	if err != nil {
		return nil, err
	}
	return resp.GetDraft(), nil
}

// ReadVersion returns a version with its component and posting rows.
func (u *UseCases) ReadVersion(ctx context.Context, versionID string) (*VersionDetail, error) {
	resp, err := u.ReadChargePolicyVersion(ctx, &versionpb.ReadChargePolicyVersionRequest{Data: &versionpb.ChargePolicyVersion{Id: versionID}})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 {
		return nil, nil
	}
	return &VersionDetail{Version: resp.GetData()[0], Components: resp.GetComponents(), Postings: resp.GetPostings()}, nil
}

// UpdateVersion updates a draft version's classification fields.
func (u *UseCases) UpdateVersion(ctx context.Context, v *versionpb.ChargePolicyVersion) error {
	_, err := u.UpdateChargePolicyVersion(ctx, &versionpb.UpdateChargePolicyVersionRequest{Data: v})
	return err
}

// DeleteVersion deletes a draft version.
func (u *UseCases) DeleteVersion(ctx context.Context, versionID string) error {
	_, err := u.DeleteChargePolicyVersion(ctx, &versionpb.DeleteChargePolicyVersionRequest{Data: &versionpb.ChargePolicyVersion{Id: versionID}})
	return err
}

// CreateComponent adds a component row to a draft version.
func (u *UseCases) CreateComponent(ctx context.Context, c *componentpb.ChargePolicyComponent) error {
	_, err := u.CreateChargePolicyComponent(ctx, &componentpb.CreateChargePolicyComponentRequest{Data: c})
	return err
}

// UpdateComponent updates a component row.
func (u *UseCases) UpdateComponent(ctx context.Context, c *componentpb.ChargePolicyComponent) error {
	_, err := u.UpdateChargePolicyComponent(ctx, &componentpb.UpdateChargePolicyComponentRequest{Data: c})
	return err
}

// DeleteComponent deletes a component row.
func (u *UseCases) DeleteComponent(ctx context.Context, id string) error {
	_, err := u.DeleteChargePolicyComponent(ctx, &componentpb.DeleteChargePolicyComponentRequest{Data: &componentpb.ChargePolicyComponent{Id: id}})
	return err
}

// CreatePosting adds a posting row to a draft version.
func (u *UseCases) CreatePosting(ctx context.Context, p *postingpb.ChargePolicyPosting) error {
	_, err := u.CreateChargePolicyPosting(ctx, &postingpb.CreateChargePolicyPostingRequest{Data: p})
	return err
}

// UpdatePosting updates a posting row.
func (u *UseCases) UpdatePosting(ctx context.Context, p *postingpb.ChargePolicyPosting) error {
	_, err := u.UpdateChargePolicyPosting(ctx, &postingpb.UpdateChargePolicyPostingRequest{Data: p})
	return err
}

// DeletePosting deletes a posting row.
func (u *UseCases) DeletePosting(ctx context.Context, id string) error {
	_, err := u.DeleteChargePolicyPosting(ctx, &postingpb.DeleteChargePolicyPostingRequest{Data: &postingpb.ChargePolicyPosting{Id: id}})
	return err
}

// Validate returns the approval-readiness report of a draft version.
func (u *UseCases) Validate(ctx context.Context, versionID string) (*versionpb.ChargePolicyApprovalChecklist, error) {
	resp, err := u.ValidateChargePolicyVersionForApproval(ctx, &versionpb.ValidateChargePolicyVersionForApprovalRequest{ChargePolicyVersionId: versionID})
	if err != nil {
		return nil, err
	}
	return resp.GetChecklist(), nil
}

// Approve approves the version. On a checklist / unsupported-combination
// refusal the coded error may carry the per-item report via `Checklist()`; it
// is returned alongside the error.
func (u *UseCases) Approve(ctx context.Context, versionID string) (*versionpb.ChargePolicyApprovalChecklist, error) {
	resp, err := u.ApproveChargePolicyVersion(ctx, &versionpb.ApproveChargePolicyVersionRequest{ChargePolicyVersionId: versionID})
	if err != nil {
		var withList interface {
			Checklist() *versionpb.ChargePolicyApprovalChecklist
		}
		if errors.As(err, &withList) {
			return withList.Checklist(), err
		}
		return nil, err
	}
	return resp.GetChecklist(), nil
}

// ListAccounts returns every account of the chart (paged through).
func (u *UseCases) ListAccounts(ctx context.Context) ([]*accountpb.Account, error) {
	return pageAll(func(p *commonpb.PaginationRequest) ([]*accountpb.Account, *commonpb.PaginationResponse, error) {
		resp, err := u.GetAccountListPageData(ctx, &accountpb.GetAccountListPageDataRequest{Pagination: p})
		if err != nil {
			return nil, nil, err
		}
		return resp.GetAccountList(), resp.GetPagination(), nil
	})
}
