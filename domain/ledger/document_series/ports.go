package document_series

import (
	"context"
	"errors"
	"fmt"

	commonpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/common"
	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
)

// Error codes carried by a use-case refusal (`ErrorCode()`); each maps to the
// Lyngua key `document_series.errors.<code>` (Labels.ErrorMessage). They are the
// espyna document_series use-case codes.
const (
	ErrKindNone                = ""
	ErrKindValidation          = "validation"
	ErrKindNotFound            = "not_found"
	ErrKindCodeTaken           = "code_taken"
	ErrKindRetired             = "retired"
	ErrKindTransactionRequired = "transaction_required"
	ErrKindNumberingLocked     = "numbering_locked"
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

// UseCases is the set of espyna document series use-case closures the views
// consume, exactly the proto request/response shape of
// `uc.Revenue.DocumentSeries.<Field>.Execute`. The block binds them and refuses
// boot when one is nil (RequireDocumentSeries / Mount).
type UseCases struct {
	CreateDocumentSeries          func(context.Context, *pb.CreateDocumentSeriesRequest) (*pb.CreateDocumentSeriesResponse, error)
	ReadDocumentSeries            func(context.Context, *pb.ReadDocumentSeriesRequest) (*pb.ReadDocumentSeriesResponse, error)
	UpdateDocumentSeries          func(context.Context, *pb.UpdateDocumentSeriesRequest) (*pb.UpdateDocumentSeriesResponse, error)
	GetDocumentSeriesListPageData func(context.Context, *pb.GetDocumentSeriesListPageDataRequest) (*pb.GetDocumentSeriesListPageDataResponse, error)
}

const pageLimit = int32(100)
const maxPages = 1000

// ErrTooManyPages is returned when a list walk passes the page bound: the view
// fails closed instead of silently showing a truncated list.
var ErrTooManyPages = errors.New("document_series: list exceeds the page bound")

func pageRequest(page int32) *commonpb.PaginationRequest {
	return &commonpb.PaginationRequest{
		Limit:  pageLimit,
		Method: &commonpb.PaginationRequest_Offset{Offset: &commonpb.OffsetPagination{Page: page}},
	}
}

// ListSeries returns every series of the workspace (paged through, C11).
func (u *UseCases) ListSeries(ctx context.Context) ([]*pb.DocumentSeries, error) {
	var out []*pb.DocumentSeries
	for page := int32(1); page <= maxPages; page++ {
		resp, err := u.GetDocumentSeriesListPageData(ctx, &pb.GetDocumentSeriesListPageDataRequest{Pagination: pageRequest(page)})
		if err != nil {
			return nil, err
		}
		items := resp.GetDocumentSeriesList()
		out = append(out, items...)
		if pag := resp.GetPagination(); pag != nil {
			if !pag.GetHasNext() {
				return out, nil
			}
		} else if int32(len(items)) < pageLimit {
			return out, nil
		}
	}
	return nil, ErrTooManyPages
}

// ReadSeries returns one series, nil when absent.
func (u *UseCases) ReadSeries(ctx context.Context, id string) (*pb.DocumentSeries, error) {
	resp, err := u.ReadDocumentSeries(ctx, &pb.ReadDocumentSeriesRequest{Data: &pb.DocumentSeries{Id: id}})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 {
		return nil, nil
	}
	return resp.GetData()[0], nil
}

// Create creates an ACTIVE series.
func (u *UseCases) Create(ctx context.Context, s *pb.DocumentSeries) (*pb.DocumentSeries, error) {
	resp, err := u.CreateDocumentSeries(ctx, &pb.CreateDocumentSeriesRequest{Data: s})
	if err != nil {
		return nil, err
	}
	if len(resp.GetData()) == 0 {
		return nil, fmt.Errorf("document_series: create returned no row")
	}
	return resp.GetData()[0], nil
}

// Update edits the mutable fields (name, issuer, prefix, branch, padding).
func (u *UseCases) Update(ctx context.Context, s *pb.DocumentSeries) error {
	_, err := u.UpdateDocumentSeries(ctx, &pb.UpdateDocumentSeriesRequest{Data: s})
	return err
}

// Retire retires a series as a status-only update: id + RETIRED, plus the
// current number_padding only because the use case requires it (>= 1) on every
// update. Name, issuer, prefix and branch are NOT resent, so a concurrent edit
// of them between the read and the update cannot be overwritten (R5 m5); the
// use case applies the change under its FOR UPDATE row lock.
func (u *UseCases) Retire(ctx context.Context, id string) error {
	cur, err := u.ReadSeries(ctx, id)
	if err != nil {
		return err
	}
	if cur == nil {
		return errNotFound
	}
	return u.Update(ctx, &pb.DocumentSeries{
		Id: cur.GetId(), NumberPadding: cur.GetNumberPadding(),
		Status: pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED,
	})
}

type codedError string

func (e codedError) Error() string     { return "document_series: " + string(e) }
func (e codedError) ErrorCode() string { return string(e) }

const errNotFound = codedError(ErrKindNotFound)

// IsRetired reports whether the series is retired.
func IsRetired(s *pb.DocumentSeries) bool {
	return s.GetStatus() == pb.DocumentSeriesStatus_DOCUMENT_SERIES_STATUS_RETIRED
}
