package block

import (
	"context"

	"github.com/erniealice/espyna-golang/ports"
	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
)

// Infra carries the subset of AppContext that view modules need beyond
// the typed UseCases: attachment ops and reference checker. Built once by
// service-admin and passed into each catalog binder.
//
// Fycha modules do not need a raw DB handle — all data access goes through
// typed espyna use cases on UseCases (hexagonal invariant, no-direct-sql rule).
type Infra struct {
	UploadFile       func(context.Context, string, string, []byte, string) error
	ListAttachments  func(context.Context, string, string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(context.Context, *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(context.Context, *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
	NewAttachmentID  func() string
	RefChecker       ports.Checker

	// AssetDepreciationRunURL is the optional URL for the asset depreciation-run
	// toast / redirect on the asset module (Surface A). Empty means no redirect.
	AssetDepreciationRunURL string
}
