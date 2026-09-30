package action

import (
	"net/http"
	"net/url"
	"testing"

	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"

	ds "github.com/erniealice/fycha-golang/domain/ledger/document_series"
	"github.com/erniealice/fycha-golang/domain/ledger/document_series/internal/dstest"
)

func newDeps(f *dstest.Fake) *Deps {
	return &Deps{Routes: ds.DefaultRoutes(), Labels: ds.DefaultLabels(), CommonLabels: dstest.Common(), UseCases: f.UseCases()}
}

func formBody(kv ...string) string {
	v := url.Values{}
	for i := 0; i+1 < len(kv); i += 2 {
		v.Set(kv[i], kv[i+1])
	}
	return v.Encode()
}

func TestEveryWriteIsPermissionGatedAndMakesNoCall(t *testing.T) {
	f := seedFake()
	deps := newDeps(f)
	none := dstest.Ctx("document_series:list")
	for name, run := range map[string]func() int{
		"add": func() int {
			return NewAddAction(deps).Handle(none, dstest.Request(http.MethodPost, "/a", formBody("code", "X"))).StatusCode
		},
		"edit": func() int {
			return NewEditAction(deps).Handle(none, dstest.Request(http.MethodPost, "/e", formBody("issuer_name", "X"), "id", "s1")).StatusCode
		},
		"retire": func() int {
			return NewRetireAction(deps).Handle(none, dstest.Request(http.MethodPost, "/r", "", "id", "s1")).StatusCode
		},
	} {
		if got := run(); got != http.StatusUnprocessableEntity {
			t.Errorf("%s without permission = %d, want 422 (HTMXError)", name, got)
		}
	}
	if f.MutationCalls() != 0 || f.Calls["Read"] != 0 {
		t.Errorf("a denied request reached the use cases: %v", f.Calls)
	}
}

func seedFake() *dstest.Fake {
	f := dstest.New()
	f.Add("s1", "SOA", false, 4)
	f.Add("s2", "OLD", true, 9)
	f.Add("s3", "NEW", false, 1) // nothing issued yet: prefix/padding still editable
	return f
}

func TestAddCreatesAnActiveSeriesWithOnlyNoneFiscalReset(t *testing.T) {
	f := seedFake()
	res := NewAddAction(newDeps(f)).Handle(dstest.Ctx("document_series:create"), dstest.Request(http.MethodPost, "/a",
		formBody("code", "soa", "name", "Statements", "issuer_name", "Acme", "document_kind", "RECOVERY_DOCUMENT", "prefix", "SOA-", "number_padding", "6", "next_number", "1")))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("create status = %d", res.StatusCode)
	}
	got := f.LastCreate
	if got == nil || got.GetCode() != "SOA" || got.GetDocumentKind() != enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT {
		t.Fatalf("create request = %v", got)
	}
	if got.GetFiscalReset() != pb.DocumentSeriesFiscalReset_DOCUMENT_SERIES_FISCAL_RESET_NONE {
		t.Error("only NONE fiscal reset is implemented; the form must never send YEARLY")
	}
}

func TestAddRejectsInvalidFormWithoutCallingTheUseCase(t *testing.T) {
	for name, body := range map[string]string{
		"no code":    formBody("issuer_name", "A", "document_kind", "INVOICE", "number_padding", "6"),
		"no issuer":  formBody("code", "A", "document_kind", "INVOICE", "number_padding", "6"),
		"bad kind":   formBody("code", "A", "issuer_name", "A", "document_kind", "NOPE", "number_padding", "6"),
		"padding 0":  formBody("code", "A", "issuer_name", "A", "document_kind", "INVOICE", "number_padding", "0"),
		"padding 19": formBody("code", "A", "issuer_name", "A", "document_kind", "INVOICE", "number_padding", "19"),
		"next 0":     formBody("code", "A", "issuer_name", "A", "document_kind", "INVOICE", "number_padding", "6", "next_number", "0"),
	} {
		f := seedFake()
		res := NewAddAction(newDeps(f)).Handle(dstest.Ctx("document_series:create"), dstest.Request(http.MethodPost, "/a", body))
		if res.StatusCode != http.StatusUnprocessableEntity || f.MutationCalls() != 0 {
			t.Errorf("%s: status %d, calls %v", name, res.StatusCode, f.Calls)
		}
	}
}

// Refusal codes map to their own Lyngua messages (C14).
func TestUseCaseRefusalsMapToTheirLabels(t *testing.T) {
	l := ds.DefaultLabels()
	cases := []struct {
		err  error
		want string
	}{
		{dstest.ErrCodeTaken, l.Errors.CodeTaken},
		{dstest.ErrRetired, l.Errors.Retired},
		{dstest.ErrNotFound, l.Errors.NotFound},
	}
	for _, c := range cases {
		f := seedFake()
		f.CreateErr = c.err
		res := NewAddAction(newDeps(f)).Handle(dstest.Ctx("document_series:create"), dstest.Request(http.MethodPost, "/a",
			formBody("code", "A", "issuer_name", "A", "document_kind", "INVOICE", "number_padding", "6")))
		if got := res.Headers["HX-Error-Message"]; got != c.want {
			t.Errorf("%v -> %q, want %q", c.err, got, c.want)
		}
	}
}

func TestEditSendsOnlyMutableFieldsAndNeverStatus(t *testing.T) {
	f := seedFake()
	res := NewEditAction(newDeps(f)).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/e",
		formBody("issuer_name", "New Issuer", "name", "N", "number_padding", "8", "code", "HACK", "next_number", "999"), "id", "s3"))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("edit status = %d", res.StatusCode)
	}
	u := f.LastUpdate
	if u.GetId() != "s3" || u.GetIssuerName() != "New Issuer" || u.GetNumberPadding() != 8 {
		t.Errorf("update = %v", u)
	}
	if u.GetCode() != "" || u.GetNextNumber() != 0 || u.GetStatus() != pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_UNSPECIFIED {
		t.Errorf("edit must not send code/next_number/status: %v", u)
	}
}

func TestEditDrawerRefusesARetiredSeries(t *testing.T) {
	f := seedFake()
	res := NewEditAction(newDeps(f)).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodGet, "/e", "", "id", "s2"))
	if res.Headers["HX-Error-Message"] != ds.DefaultLabels().Errors.Retired {
		t.Errorf("retired edit drawer = %v", res.Headers)
	}
}

func TestRetireIsAStatusOnlyUpdate(t *testing.T) {
	f := seedFake()
	res := NewRetireAction(newDeps(f)).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/r", "", "id", "s1"))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("retire status = %d", res.StatusCode)
	}
	u := f.LastUpdate
	if u.GetStatus() != pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED || u.GetNumberPadding() != 6 {
		t.Errorf("retire update = %v", u)
	}
	// R5 m5: status-only — the descriptive fields are not resent, so a
	// concurrent edit of them cannot be overwritten.
	if u.Name != nil || u.GetIssuerName() != "" || u.Prefix != nil || u.BranchCode != nil || u.IssuerTaxId != nil {
		t.Errorf("retire resent descriptive fields: %v", u)
	}
	// Unknown id refuses with not_found and writes nothing.
	f2 := seedFake()
	res = NewRetireAction(newDeps(f2)).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/r", "", "id", "nope"))
	if res.Headers["HX-Error-Message"] != ds.DefaultLabels().Errors.NotFound || f2.Calls["Update"] != 0 {
		t.Errorf("unknown id: %v calls %v", res.Headers, f2.Calls)
	}
}

func TestUnboundUseCasesAnswer503(t *testing.T) {
	deps := newDeps(seedFake())
	deps.UseCases = &ds.UseCases{}
	if got := NewRetireAction(deps).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/r", "", "id", "s1")).StatusCode; got != http.StatusServiceUnavailable {
		t.Errorf("retire with no use cases = %d", got)
	}
}

// Once numbers were issued (next_number > 1) prefix and padding are locked: the
// view keeps the stored values whatever the form posts (R5 m5).
func TestEditKeepsPrefixAndPaddingOnceNumbersWereIssued(t *testing.T) {
	f := seedFake() // s1: next_number 4, prefix "SOA-", padding 6
	res := NewEditAction(newDeps(f)).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/e",
		formBody("issuer_name", "New Issuer", "prefix", "HACK-", "number_padding", "12"), "id", "s1"))
	if res.StatusCode != http.StatusOK {
		t.Fatalf("edit status = %d", res.StatusCode)
	}
	u := f.LastUpdate
	if u.GetPrefix() != "SOA-" || u.GetNumberPadding() != 6 || u.GetIssuerName() != "New Issuer" {
		t.Errorf("locked edit update = %v", u)
	}
	// A disabled padding input posts nothing: still accepted for a locked series.
	res = NewEditAction(newDeps(seedFake())).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/e",
		formBody("issuer_name", "X"), "id", "s1"))
	if res.StatusCode != http.StatusOK {
		t.Errorf("locked edit without prefix/padding = %d", res.StatusCode)
	}
}

// A strict-gate denial after the view's permission check maps to the shared
// PermissionDenied message, not the generic/unknown error (R5 M1, C20).
func TestStrictGateDenialMapsToPermissionDenied(t *testing.T) {
	want := dstest.Common().Errors.PermissionDenied
	f := seedFake()
	f.CreateErr = dstest.ErrPermissionDenied
	res := NewAddAction(newDeps(f)).Handle(dstest.Ctx("document_series:create"), dstest.Request(http.MethodPost, "/a",
		formBody("code", "A", "issuer_name", "A", "document_kind", "INVOICE", "number_padding", "6")))
	if got := res.Headers["HX-Error-Message"]; got != want || want == "" {
		t.Errorf("denial -> %q, want %q", got, want)
	}
}

// A numbering_locked refusal from the update use case surfaces its own message
// on the edit POST, not the generic error (R5 B2, C2/C31).
func TestEditMapsNumberingLockedRefusal(t *testing.T) {
	f := seedFake()
	f.UpdateErr = dstest.ErrNumberingLocked
	res := NewEditAction(newDeps(f)).Handle(dstest.Ctx("document_series:update"), dstest.Request(http.MethodPost, "/e",
		formBody("issuer_name", "X"), "id", "s1"))
	want := ds.DefaultLabels().Errors.NumberingLocked
	if got := res.Headers["HX-Error-Message"]; got != want || want == "" {
		t.Errorf("numbering_locked -> %q, want %q", got, want)
	}
}
