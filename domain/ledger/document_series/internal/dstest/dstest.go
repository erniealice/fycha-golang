// Package dstest is an in-memory fake of the document series ports plus request
// helpers, shared by the view tests of the document series pages. Test-only.
package dstest

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"

	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ds "github.com/erniealice/fycha-golang/domain/ledger/document_series"
)

// Coded is a fake use-case refusal: it exposes ErrorCode() like espyna's Error.
type Coded struct{ code string }

func (e *Coded) Error() string     { return "document_series: " + e.code }
func (e *Coded) ErrorCode() string { return e.code }

// Refusals returned by the fake.
var (
	ErrCodeTaken error = &Coded{"code_taken"}
	ErrRetired   error = &Coded{"retired"}
	ErrNotFound  error = &Coded{"not_found"}
	// ErrNumberingLocked is the prefix/padding edit refusal once numbers were issued.
	ErrNumberingLocked error = &Coded{"numbering_locked"}
	// ErrPermissionDenied is the strict-gate denial (actiongate `permission_denied`).
	ErrPermissionDenied error = &Coded{"permission_denied"}
)

// Fake is the in-memory store behind the ports.
type Fake struct {
	Series     []*pb.DocumentSeries
	CreateErr  error
	UpdateErr  error
	Calls      map[string]int
	LastCreate *pb.DocumentSeries
	LastUpdate *pb.DocumentSeries
}

// New returns an empty fake.
func New() *Fake { return &Fake{Calls: map[string]int{}} }

// MutationCalls counts every write call made through the ports.
func (f *Fake) MutationCalls() int { return f.Calls["Create"] + f.Calls["Update"] }

// Add registers a series.
func (f *Fake) Add(id, code string, retired bool, next int64) *pb.DocumentSeries {
	st := pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_ACTIVE
	if retired {
		st = pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED
	}
	prefix := code + "-"
	name := "Series " + code
	s := &pb.DocumentSeries{Id: id, Code: code, Name: &name, IssuerName: "Acme Realty", Prefix: &prefix,
		NextNumber: next, NumberPadding: 6, Status: st, Active: !retired}
	f.Series = append(f.Series, s)
	return s
}

// UseCases returns the closures backed by the fake.
func (f *Fake) UseCases() *ds.UseCases {
	c := func(n string) { f.Calls[n]++ }
	return &ds.UseCases{
		CreateDocumentSeries: func(ctx context.Context, req *pb.CreateDocumentSeriesRequest) (*pb.CreateDocumentSeriesResponse, error) {
			c("Create")
			if f.CreateErr != nil {
				return nil, f.CreateErr
			}
			f.LastCreate = req.GetData()
			row := *req.GetData()
			row.Id = "new-" + row.Code
			f.Series = append(f.Series, &row)
			return &pb.CreateDocumentSeriesResponse{Data: []*pb.DocumentSeries{&row}, Success: true}, nil
		},
		ReadDocumentSeries: func(ctx context.Context, req *pb.ReadDocumentSeriesRequest) (*pb.ReadDocumentSeriesResponse, error) {
			c("Read")
			for _, s := range f.Series {
				if s.GetId() == req.GetData().GetId() {
					return &pb.ReadDocumentSeriesResponse{Data: []*pb.DocumentSeries{s}, Success: true}, nil
				}
			}
			return &pb.ReadDocumentSeriesResponse{Success: true}, nil
		},
		UpdateDocumentSeries: func(ctx context.Context, req *pb.UpdateDocumentSeriesRequest) (*pb.UpdateDocumentSeriesResponse, error) {
			c("Update")
			if f.UpdateErr != nil {
				return nil, f.UpdateErr
			}
			f.LastUpdate = req.GetData()
			for _, s := range f.Series {
				if s.GetId() == req.GetData().GetId() && req.GetData().GetStatus() == pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED {
					s.Status = pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED
				}
			}
			return &pb.UpdateDocumentSeriesResponse{Success: true}, nil
		},
		GetDocumentSeriesListPageData: func(ctx context.Context, req *pb.GetDocumentSeriesListPageDataRequest) (*pb.GetDocumentSeriesListPageDataResponse, error) {
			c("List")
			return &pb.GetDocumentSeriesListPageDataResponse{DocumentSeriesList: f.Series, Success: true}, nil
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
