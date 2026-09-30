package charge_policy

import (
	"context"
	"sort"
	"time"

	"github.com/erniealice/espyna-golang/shared/identity"
	policypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pyezatypes "github.com/erniealice/pyeza-golang/types"
)

// Summary is the per-policy view model shared by the list and detail pages.
// The Draft list is COMPUTED (UI-CP-6): a policy with no APPROVED version that
// is not retired.
type Summary struct {
	Policy   *policypb.ChargePolicy
	Versions []*versionpb.ChargePolicyVersion // version_number descending
	Current  *versionpb.ChargePolicyVersion   // newest APPROVED
	Draft    *versionpb.ChargePolicyVersion   // the open DRAFT, if any
	Last     *versionpb.ChargePolicyVersion   // highest version_number
	Bucket   string                           // active | draft | retired
}

// Summarize builds a Summary; versions may be in any order.
func Summarize(p *policypb.ChargePolicy, versions []*versionpb.ChargePolicyVersion) Summary {
	vs := make([]*versionpb.ChargePolicyVersion, 0, len(versions))
	for _, v := range versions {
		if v != nil {
			vs = append(vs, v)
		}
	}
	sort.SliceStable(vs, func(i, j int) bool { return vs[i].GetVersionNumber() > vs[j].GetVersionNumber() })
	s := Summary{Policy: p, Versions: vs}
	for _, v := range vs {
		if s.Last == nil {
			s.Last = v
		}
		switch v.GetStatus() {
		case enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_APPROVED:
			if s.Current == nil {
				s.Current = v
			}
		case enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_DRAFT:
			if s.Draft == nil {
				s.Draft = v
			}
		}
	}
	switch {
	case p.GetStatus() == enumspb.ChargePolicyStatus_CHARGE_POLICY_STATUS_RETIRED:
		s.Bucket = "retired"
	case s.Current != nil:
		s.Bucket = "active"
	default:
		s.Bucket = "draft"
	}
	return s
}

// IsRetired reports whether the policy is retired.
func (s Summary) IsRetired() bool {
	return s.Policy.GetStatus() == enumspb.ChargePolicyStatus_CHARGE_POLICY_STATUS_RETIRED
}

// FormatMillis renders an epoch-millisecond timestamp in the workspace time zone.
func FormatMillis(ctx context.Context, ms *int64) string {
	if ms == nil || *ms == 0 {
		return ""
	}
	return pyezatypes.FormatInTZ(time.UnixMilli(*ms), pyezatypes.LocationFromContext(ctx), "2006-01-02")
}

// IsDraft reports whether the version is an open draft.
func IsDraft(v *versionpb.ChargePolicyVersion) bool {
	return v.GetStatus() == enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_DRAFT
}

// ActorID returns the authenticated user id, or "" when the request carries no
// identity (never panics).
func ActorID(ctx context.Context) string {
	if id, ok := identity.FromContext(ctx); ok && id != nil {
		return id.UserID
	}
	return ""
}
