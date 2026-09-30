package document_series

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
)

func lyngua(rel string) string {
	return filepath.Join("..", "..", "..", "..", "lyngua", "translations", "en", rel)
}

// The Go json tags must byte-match the Lyngua key tree: decoding the general and
// leasing document_series.json files with unknown fields disallowed proves that
// every Lyngua key is consumed by a Go field.
func TestLabelsMatchLynguaKeyTree(t *testing.T) {
	for _, rel := range []string{"general/document_series.json", "leasing/document_series.json"} {
		raw, err := os.ReadFile(lyngua(rel))
		if err != nil {
			t.Skipf("lyngua file not reachable (%s): %v", rel, err)
		}
		var doc struct {
			DocumentSeries Labels `json:"document_series"`
		}
		dec := json.NewDecoder(bytes.NewReader(raw))
		dec.DisallowUnknownFields()
		if err := dec.Decode(&doc); err != nil {
			t.Errorf("%s does not decode into Labels: %v", rel, err)
		}
	}
}

// C14: the English defaults equal the general Lyngua value of EVERY string field.
func TestDefaultLabelsEqualGeneralLyngua(t *testing.T) {
	raw, err := os.ReadFile(lyngua("general/document_series.json"))
	if err != nil {
		t.Skipf("lyngua file not reachable: %v", err)
	}
	var doc struct {
		DocumentSeries Labels `json:"document_series"`
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
		}
	}
	walk("Labels", reflect.ValueOf(DefaultLabels()), reflect.ValueOf(doc.DocumentSeries))
}

// Every use-case error code maps to its own message.
func TestEveryErrorCodeHasALabel(t *testing.T) {
	l := DefaultLabels()
	seen := map[string]string{}
	for _, c := range []string{ErrKindValidation, ErrKindNotFound, ErrKindCodeTaken, ErrKindRetired, ErrKindTransactionRequired, ErrKindNumberingLocked} {
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

func TestErrorKindUsesErrorCode(t *testing.T) {
	if got := ErrorKind(fmt.Errorf("wrapped: %w", codedErr("code_taken"))); got != ErrKindCodeTaken {
		t.Errorf("wrapped coded error = %q", got)
	}
	if got := ErrorKind(errors.New("boom")); got != ErrKindUnknown {
		t.Errorf("uncoded error = %q", got)
	}
	if ErrorKind(nil) != ErrKindNone {
		t.Error("nil error must classify as none")
	}
}

func TestNextDocumentNumberPadsAndPrefixes(t *testing.T) {
	prefix := "SOA-"
	s := &pb.DocumentSeries{Prefix: &prefix, NextNumber: 42, NumberPadding: 6}
	if got := NextDocumentNumber(s); got != "SOA-000042" {
		t.Errorf("NextDocumentNumber = %q", got)
	}
	if got := Format("Next document: {0}", "SOA-000042"); got != "Next document: SOA-000042" {
		t.Errorf("Format = %q", got)
	}
}

func TestListSeriesFailsClosedPastPageBound(t *testing.T) {
	u := &UseCases{GetDocumentSeriesListPageData: func(_ context.Context, _ *pb.GetDocumentSeriesListPageDataRequest) (*pb.GetDocumentSeriesListPageDataResponse, error) {
		return &pb.GetDocumentSeriesListPageDataResponse{Pagination: &commonpb.PaginationResponse{HasNext: true}}, nil
	}}
	if _, err := u.ListSeries(context.Background()); !errors.Is(err, ErrTooManyPages) {
		t.Fatalf("err = %v, want ErrTooManyPages", err)
	}
}
