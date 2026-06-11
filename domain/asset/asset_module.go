package asset

import (
	"context"

	"github.com/erniealice/hybra-golang/views/attachment"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	entityasset "github.com/erniealice/fycha-golang/domain/asset/asset"
	assetaction "github.com/erniealice/fycha-golang/domain/asset/asset/action"
	assetdashboard "github.com/erniealice/fycha-golang/domain/asset/asset/dashboard"
	assetdetail "github.com/erniealice/fycha-golang/domain/asset/asset/detail"
	assetform "github.com/erniealice/fycha-golang/domain/asset/asset/form"
	assetlist "github.com/erniealice/fycha-golang/domain/asset/asset/list"
	assetrevaluation "github.com/erniealice/fycha-golang/domain/asset/asset_revaluation"
	depreciationrun "github.com/erniealice/fycha-golang/domain/asset/depreciation_run"
)

// AssetModuleDeps holds all dependencies for the asset module.
type AssetModuleDeps struct {
	Routes       entityasset.Routes
	CommonLabels pyeza.CommonLabels
	Labels       entityasset.Labels
	TableLabels  types.TableLabels

	// Depreciation Run + Revaluation labels (Surface A / E)
	DepreciationRunLabels  depreciationrun.Labels
	AssetRevaluationLabels assetrevaluation.Labels

	// CRUD operations (wired from block.go via raw SQL)
	CreateAsset func(ctx context.Context, asset *assetform.Record) error
	ReadAsset   func(ctx context.Context, id string) (*assetform.Record, error)
	UpdateAsset func(ctx context.Context, asset *assetform.Record) error
	DeleteAsset func(ctx context.Context, id string) error
	SetActive   func(ctx context.Context, id string, active bool) error
	ListAssets  func(ctx context.Context, status string) ([]assetlist.AssetRow, error)
	NewID       func() string

	// Depreciation-fields lock check (wired when espyna use cases are available)
	DepreciationFieldsLockedFn func(ctx context.Context, assetID string) (bool, error)

	// GetAssetInUseIDs returns a map of asset IDs that have any asset_transaction
	// row. Wired from block.go via the reference checker. Used for both the list
	// page (data-deletable attribute) and the delete handler (H5 server-side gate).
	// Nil = skip the check (mock build or use cases not yet wired).
	GetAssetInUseIDs func(ctx context.Context, ids []string) (map[string]bool, error)

	// Depreciation run use-case wrappers (Surface A)
	ListDepreciationCandidates func(ctx context.Context, assetID, asOfDate string) ([]assetaction.DepreciationCandidate, error)
	GenerateDepreciationRun    func(ctx context.Context, req assetaction.DepreciationRunRequest) (*assetaction.DepreciationRunResult, error)
	// AssetDepreciationRunURL is the run-detail URL for toast link plumbing.
	// Injected by block.go via WithAssetDepreciationRunURL.
	AssetDepreciationRunURL string
	// DepreciationRunRoutes for resolving run-detail links
	DepreciationRunRoutes depreciationrun.Routes

	// Revaluation use-case wrappers (Surface E)
	RevalueAsset       func(ctx context.Context, req assetaction.RevaluationRequest) (*assetaction.RevaluationResult, error)
	PreviewRevaluation func(ctx context.Context, assetID string, newFairValue int64) (*assetaction.RevaluationPreview, error)

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
}

// AssetModule holds all constructed asset views.
type AssetModule struct {
	routes           entityasset.Routes
	Dashboard        view.View
	List             view.View
	Table            view.View
	Detail           view.View
	TabAction        view.View
	Add              view.View
	Edit             view.View
	Delete           view.View
	BulkDelete       view.View
	SetStatus        view.View
	BulkSetStatus    view.View
	AttachmentUpload view.View
	AttachmentDelete view.View
	// Surface A — per-asset depreciation-run drawer
	DepreciationRun view.View
	// Surface E — per-asset revaluation drawer + preview endpoint
	Revaluation        view.View
	RevaluationPreview view.View
}

// NewAssetModule creates an asset module with all views wired.
func NewAssetModule(deps *AssetModuleDeps) *AssetModule {
	listDeps := &assetlist.ListViewDeps{
		Routes:           deps.Routes,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		ListAssets:       deps.ListAssets,
		GetAssetInUseIDs: deps.GetAssetInUseIDs,
	}
	actionDeps := &assetaction.Deps{
		Routes:                     deps.Routes,
		Labels:                     deps.Labels,
		CreateAsset:                deps.CreateAsset,
		ReadAsset:                  deps.ReadAsset,
		UpdateAsset:                deps.UpdateAsset,
		DeleteAsset:                deps.DeleteAsset,
		SetActive:                  deps.SetActive,
		NewID:                      deps.NewID,
		DepreciationFieldsLockedFn: deps.DepreciationFieldsLockedFn,
		GetAssetInUseIDs:           deps.GetAssetInUseIDs,
	}
	detailDeps := &assetdetail.DetailViewDeps{
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		Routes:                 deps.Routes,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		DepreciationRunLabels:  deps.DepreciationRunLabels,
		AssetRevaluationLabels: deps.AssetRevaluationLabels,
	}
	dashboardDeps := &assetdashboard.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
	}

	depRunDeps := &assetaction.DepreciationRunDeps{
		Routes:                     deps.Routes,
		DepreciationRunRoutes:      deps.DepreciationRunRoutes,
		Labels:                     deps.DepreciationRunLabels,
		ListDepreciationCandidates: deps.ListDepreciationCandidates,
		GenerateDepreciationRun:    deps.GenerateDepreciationRun,
		AssetDepreciationRunURL:    deps.AssetDepreciationRunURL,
	}

	revalDeps := &assetaction.RevaluationDeps{
		Routes:             deps.Routes,
		Labels:             deps.AssetRevaluationLabels,
		RevalueAsset:       deps.RevalueAsset,
		PreviewRevaluation: deps.PreviewRevaluation,
	}

	return &AssetModule{
		routes:             deps.Routes,
		Dashboard:          assetdashboard.NewView(dashboardDeps),
		List:               assetlist.NewView(listDeps),
		Table:              assetlist.NewTableView(listDeps),
		Detail:             assetdetail.NewView(detailDeps),
		TabAction:          assetdetail.NewTabAction(detailDeps),
		Add:                assetaction.NewAddAction(actionDeps),
		Edit:               assetaction.NewEditAction(actionDeps),
		Delete:             assetaction.NewDeleteAction(actionDeps),
		BulkDelete:         assetaction.NewBulkDeleteAction(actionDeps),
		SetStatus:          assetaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:      assetaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload:   assetdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete:   assetdetail.NewAttachmentDeleteAction(detailDeps),
		DepreciationRun:    assetaction.NewDepreciationRunAction(depRunDeps),
		Revaluation:        assetaction.NewRevaluationAction(revalDeps),
		RevaluationPreview: assetaction.NewRevaluationPreviewAction(revalDeps),
	}
}

// RegisterRoutes registers all asset routes with the given route registrar.
func (m *AssetModule) RegisterRoutes(r view.RouteRegistrar) {
	r.GET(m.routes.DashboardURL, m.Dashboard)
	r.GET(m.routes.ListURL, m.List)
	r.GET(m.routes.TableURL, m.Table)
	r.GET(m.routes.DetailURL, m.Detail)
	r.GET(m.routes.TabActionURL, m.TabAction)
	r.GET(m.routes.AddURL, m.Add)
	r.POST(m.routes.AddURL, m.Add)
	r.GET(m.routes.EditURL, m.Edit)
	r.POST(m.routes.EditURL, m.Edit)
	r.POST(m.routes.DeleteURL, m.Delete)
	r.POST(m.routes.BulkDeleteURL, m.BulkDelete)
	r.POST(m.routes.SetStatusURL, m.SetStatus)
	r.POST(m.routes.BulkSetStatusURL, m.BulkSetStatus)
	// Attachments
	if m.AttachmentUpload != nil {
		r.GET(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentUploadURL, m.AttachmentUpload)
		r.POST(m.routes.AttachmentDeleteURL, m.AttachmentDelete)
	}
	// Surface A — per-asset depreciation-run drawer
	if m.DepreciationRun != nil {
		r.GET(m.routes.DepreciationRunURL, m.DepreciationRun)
		r.POST(m.routes.DepreciationRunURL, m.DepreciationRun)
	}
	// Surface E — per-asset revaluation drawer + preview
	if m.Revaluation != nil {
		r.GET(m.routes.RevaluationURL, m.Revaluation)
		r.POST(m.routes.RevaluationURL, m.Revaluation)
	}
	if m.RevaluationPreview != nil {
		r.POST(m.routes.RevaluationPreviewURL, m.RevaluationPreview)
	}
}
