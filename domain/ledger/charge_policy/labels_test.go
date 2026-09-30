package charge_policy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// The Go json tags must byte-match the Lyngua key tree: decoding the general and
// leasing charge_policy.json files with unknown fields disallowed proves that
// every Lyngua key is consumed by a Go field.
func TestLabelsMatchLynguaKeyTree(t *testing.T) {
	for _, rel := range []string{
		"general/charge_policy.json",
		"leasing/charge_policy.json",
	} {
		path := filepath.Join("..", "..", "..", "..", "lyngua", "translations", "en", rel)
		raw, err := os.ReadFile(path)
		if err != nil {
			t.Skipf("lyngua file not reachable (%s): %v", rel, err)
		}
		var doc struct {
			ChargePolicy Labels `json:"charge_policy"`
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&doc); err != nil {
			t.Errorf("%s does not decode into Labels: %v", rel, err)
		}
	}
}

// C14 / M3: the English defaults equal the general Lyngua value of EVERY string
// (and string-map) field — a Go default without a general-tier key is a defect.
func TestDefaultLabelsEqualGeneralLyngua(t *testing.T) {
	raw, err := os.ReadFile(filepath.Join("..", "..", "..", "..", "lyngua", "translations", "en", "general", "charge_policy.json"))
	if err != nil {
		t.Skipf("lyngua file not reachable: %v", err)
	}
	var doc struct {
		ChargePolicy Labels `json:"charge_policy"`
	}
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatal(err)
	}
	var walk func(path string, def, got reflect.Value)
	walk = func(path string, def, got reflect.Value) {
		switch def.Kind() {
		case reflect.Struct:
			for i := 0; i < def.NumField(); i++ {
				walk(path+"."+def.Type().Field(i).Name, def.Field(i), got.Field(i))
			}
		case reflect.String:
			if def.String() == "" {
				t.Errorf("%s: empty Go default", path)
			}
			if def.String() != got.String() {
				t.Errorf("%s: default %q != general lyngua %q", path, def.String(), got.String())
			}
		case reflect.Map:
			if def.Len() != got.Len() {
				t.Errorf("%s: default has %d keys, lyngua %d", path, def.Len(), got.Len())
			}
			for _, k := range def.MapKeys() {
				g := got.MapIndex(k)
				if !g.IsValid() || g.String() != def.MapIndex(k).String() {
					t.Errorf("%s[%v]: default %q != lyngua %q", path, k, def.MapIndex(k), g)
				}
			}
		}
	}
	walk("Labels", reflect.ValueOf(DefaultLabels()), reflect.ValueOf(doc.ChargePolicy))
}

// The 13 use-case error codes each map to their own general key and message.
func TestEveryErrorCodeHasALabel(t *testing.T) {
	l := DefaultLabels()
	codes := []string{ErrKindValidation, ErrKindNotFound, ErrKindNotDraft, ErrKindDraftExists, ErrKindNoApprovedVersion,
		ErrKindRetired, ErrKindSelfApproval, ErrKindChecklistFailed, ErrKindUnsupportedCombination, ErrKindInUse,
		ErrKindTransactionRequired, ErrKindTaxTreatmentNotFound, ErrKindUnverifiable}
	seen := map[string]string{}
	for _, c := range codes {
		m := l.ErrorMessage(c)
		if m == "" || m == l.Errors.Generic {
			t.Errorf("code %q has no dedicated message", c)
		}
		if prev, dup := seen[m]; dup {
			t.Errorf("codes %q and %q share the message %q", prev, c, m)
		}
		seen[m] = c
	}
	if l.ErrorMessage("nonsense") != l.Errors.Generic {
		t.Error("unknown code must map to the generic message")
	}
}

type codedErr string

func (e codedErr) Error() string     { return "x: " + string(e) }
func (e codedErr) ErrorCode() string { return string(e) }

// Errors are classified through ErrorCode() only (no espyna sentinels).
func TestErrorKindUsesErrorCode(t *testing.T) {
	if got := ErrorKind(fmt.Errorf("wrapped: %w", codedErr("self_approval"))); got != ErrKindSelfApproval {
		t.Errorf("wrapped coded error = %q", got)
	}
	if got := ErrorKind(errors.New("boom")); got != ErrKindUnknown {
		t.Errorf("uncoded error = %q", got)
	}
	if ErrorKind(nil) != ErrKindNone {
		t.Error("nil error must classify as none")
	}
}

func TestFormatSubstitutesPositionalArgs(t *testing.T) {
	if got := Format("Back to {0} · v{1}", "Utility", "3"); got != "Back to Utility · v3" {
		t.Errorf("Format = %q", got)
	}
}

func TestSummarizeBuckets(t *testing.T) {
	// covered through the list tests; here only the ordering invariant.
	s := Summarize(nil, nil)
	if s.Bucket != "draft" {
		t.Errorf("empty non-retired summary bucket = %q", s.Bucket)
	}
}
