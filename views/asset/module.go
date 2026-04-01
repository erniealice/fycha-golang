package asset

import (
	"context"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	fycha "github.com/erniealice/fycha-golang"
	assetaction "github.com/erniealice/fycha-golang/views/asset/action"
	assetdashboard "github.com/erniealice/fycha-golang/views/asset/dashboard"
	assetdetail "github.com/erniealice/fycha-golang/views/asset/detail"
	assetlist "github.com/erniealice/fycha-golang/views/asset/list"
)

// ModuleDeps holds all dependencies for the asset module.
type ModuleDeps struct {
	Routes       fycha.AssetRoutes
	CommonLabels pyeza.CommonLabels
	Labels       fycha.AssetLabels
	TableLabels  types.TableLabels

	// CRUD operations (wired from block.go via raw SQL)
	CreateAsset  func(ctx context.Context, asset *assetaction.AssetRecord) error
	ReadAsset    func(ctx context.Context, id string) (*assetaction.AssetRecord, error)
	UpdateAsset  func(ctx context.Context, asset *assetaction.AssetRecord) error
	DeleteAsset  func(ctx context.Context, id string) error
	SetActive    func(ctx context.Context, id string, active bool) error
	ListAssets   func(ctx context.Context, status string) ([]assetlist.AssetRow, error)
	NewID        func() string

	// Attachment operations
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)
}

// Module holds all constructed asset views.
type Module struct {
	routes           fycha.AssetRoutes
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
}

// NewModule creates an asset module with all views wired.
func NewModule(deps *ModuleDeps) *Module {
	listDeps := &assetlist.ListViewDeps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
		ListAssets:   deps.ListAssets,
	}
	actionDeps := &assetaction.Deps{
		Routes:      deps.Routes,
		Labels:      deps.Labels,
		CreateAsset: deps.CreateAsset,
		ReadAsset:   deps.ReadAsset,
		UpdateAsset: deps.UpdateAsset,
		DeleteAsset: deps.DeleteAsset,
		SetActive:   deps.SetActive,
		NewID:       deps.NewID,
	}
	detailDeps := &assetdetail.DetailViewDeps{
		AttachmentOps: attachment.AttachmentOps{
			UploadFile:       deps.UploadFile,
			ListAttachments:  deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment,
			DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID:  deps.NewID,
		},
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}
	dashboardDeps := &assetdashboard.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
	}

	return &Module{
		routes:           deps.Routes,
		Dashboard:        assetdashboard.NewView(dashboardDeps),
		List:             assetlist.NewView(listDeps),
		Table:            assetlist.NewTableView(listDeps),
		Detail:           assetdetail.NewView(detailDeps),
		TabAction:        assetdetail.NewTabAction(detailDeps),
		Add:              assetaction.NewAddAction(actionDeps),
		Edit:             assetaction.NewEditAction(actionDeps),
		Delete:           assetaction.NewDeleteAction(actionDeps),
		BulkDelete:       assetaction.NewBulkDeleteAction(actionDeps),
		SetStatus:        assetaction.NewSetStatusAction(actionDeps),
		BulkSetStatus:    assetaction.NewBulkSetStatusAction(actionDeps),
		AttachmentUpload: assetdetail.NewAttachmentUploadAction(detailDeps),
		AttachmentDelete: assetdetail.NewAttachmentDeleteAction(detailDeps),
	}
}

// RegisterRoutes registers all asset routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
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
}
