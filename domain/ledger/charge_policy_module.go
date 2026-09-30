package ledger

import (
	"context"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	chargepolicy "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
	cpaction "github.com/erniealice/fycha-golang/domain/ledger/charge_policy/action"
	cpdetail "github.com/erniealice/fycha-golang/domain/ledger/charge_policy/detail"
	cplist "github.com/erniealice/fycha-golang/domain/ledger/charge_policy/list"
	cpversion "github.com/erniealice/fycha-golang/domain/ledger/charge_policy/version"
)

// ChargePolicyModuleDeps holds all dependencies for the charge policy module.
// UseCases is nil when the espyna charge policy use cases are not wired; the
// module then fails closed (routes answer 503, never fake data).
type ChargePolicyModuleDeps struct {
	Routes       chargepolicy.Routes
	Labels       chargepolicy.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *chargepolicy.UseCases

	// Evidence attachments (version approval tab).
	NewAttachmentID  func() string
	UploadFile       func(ctx context.Context, bucket, key string, content []byte, contentType string) error
	ListAttachments  func(ctx context.Context, moduleKey, foreignKey string) (*attachmentpb.ListAttachmentsResponse, error)
	CreateAttachment func(ctx context.Context, req *attachmentpb.CreateAttachmentRequest) (*attachmentpb.CreateAttachmentResponse, error)
	DeleteAttachment func(ctx context.Context, req *attachmentpb.DeleteAttachmentRequest) (*attachmentpb.DeleteAttachmentResponse, error)

	// Audit history (optional; nil renders the empty history tab).
	auditlog.AuditOps
}

// ChargePolicyModule holds all constructed charge policy views.
type ChargePolicyModule struct {
	routes chargepolicy.Routes

	List, Table, Detail, TabAction view.View
	Add, Edit, Delete, BulkDelete  view.View
	Retire                         view.View

	VersionDetail, VersionTabAction view.View
	VersionAdd, VersionEdit         view.View
	VersionDelete, VersionApprove   view.View

	ComponentAdd, ComponentEdit, ComponentDelete view.View
	PostingAdd, PostingEdit, PostingDelete       view.View

	AttachmentUpload, AttachmentDelete view.View
}

// NewChargePolicyModule wires every charge policy view.
func NewChargePolicyModule(deps *ChargePolicyModuleDeps) *ChargePolicyModule {
	if deps == nil {
		deps = &ChargePolicyModuleDeps{}
	}
	ucs := deps.UseCases
	if ucs == nil {
		ucs = &chargepolicy.UseCases{} // fail closed: every feature nil
	}

	listDeps := &cplist.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs}
	detailDeps := &cpdetail.Deps{Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs, AuditOps: deps.AuditOps}
	versionDeps := &cpversion.Deps{
		Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, TableLabels: deps.TableLabels, UseCases: ucs,
		AttachmentOps: attachment.AttachmentOps{
			UploadFile: deps.UploadFile, ListAttachments: deps.ListAttachments,
			CreateAttachment: deps.CreateAttachment, DeleteAttachment: deps.DeleteAttachment,
			NewAttachmentID: deps.NewAttachmentID,
		},
		AuditOps: deps.AuditOps,
	}
	actionDeps := &cpaction.Deps{
		Routes: deps.Routes, Labels: deps.Labels, CommonLabels: deps.CommonLabels, UseCases: ucs,
		AttachmentConfig: func() *attachment.Config { return cpversion.AttachmentConfig(versionDeps) },
	}

	return &ChargePolicyModule{
		routes: deps.Routes,

		List:      cplist.NewView(listDeps),
		Table:     cplist.NewTableView(listDeps),
		Detail:    cpdetail.NewView(detailDeps),
		TabAction: cpdetail.NewTabAction(detailDeps),

		Add:        cpaction.NewAddAction(actionDeps),
		Edit:       cpaction.NewEditAction(actionDeps),
		Delete:     cpaction.NewDeleteAction(actionDeps),
		BulkDelete: cpaction.NewBulkDeleteAction(actionDeps),
		Retire:     cpaction.NewRetireAction(actionDeps),

		VersionDetail:    cpversion.NewView(versionDeps),
		VersionTabAction: cpversion.NewTabAction(versionDeps),
		VersionAdd:       cpaction.NewVersionAddAction(actionDeps),
		VersionEdit:      cpaction.NewVersionEditAction(actionDeps),
		VersionDelete:    cpaction.NewVersionDeleteAction(actionDeps),
		VersionApprove:   cpaction.NewVersionApproveAction(actionDeps),

		ComponentAdd:    cpaction.NewComponentAddAction(actionDeps),
		ComponentEdit:   cpaction.NewComponentEditAction(actionDeps),
		ComponentDelete: cpaction.NewComponentDeleteAction(actionDeps),
		PostingAdd:      cpaction.NewPostingAddAction(actionDeps),
		PostingEdit:     cpaction.NewPostingEditAction(actionDeps),
		PostingDelete:   cpaction.NewPostingDeleteAction(actionDeps),

		AttachmentUpload: cpaction.NewAttachmentUploadAction(actionDeps),
		AttachmentDelete: cpaction.NewAttachmentDeleteAction(actionDeps),
	}
}

// RegisterRoutes registers every charge policy route.
func (m *ChargePolicyModule) RegisterRoutes(r view.RouteRegistrar) {
	rt := m.routes
	// Pages.
	r.GET(rt.ListURL, m.List)
	r.GET(rt.TableURL, m.Table)
	r.GET(rt.DetailURL, m.Detail)
	r.GET(rt.TabActionURL, m.TabAction)
	r.GET(rt.VersionDetailURL, m.VersionDetail)
	r.GET(rt.VersionTabActionURL, m.VersionTabAction)

	// Policy header.
	r.GET(rt.AddURL, m.Add)
	r.POST(rt.AddURL, m.Add)
	r.GET(rt.EditURL, m.Edit)
	r.POST(rt.EditURL, m.Edit)
	r.POST(rt.DeleteURL, m.Delete)
	r.POST(rt.BulkDeleteURL, m.BulkDelete)
	r.POST(rt.RetireURL, m.Retire)

	// Versions.
	r.GET(rt.VersionAddURL, m.VersionAdd)
	r.POST(rt.VersionAddURL, m.VersionAdd)
	r.GET(rt.VersionEditURL, m.VersionEdit)
	r.POST(rt.VersionEditURL, m.VersionEdit)
	r.POST(rt.VersionDeleteURL, m.VersionDelete)
	r.POST(rt.VersionApproveURL, m.VersionApprove)

	// Components / postings.
	r.GET(rt.ComponentAddURL, m.ComponentAdd)
	r.POST(rt.ComponentAddURL, m.ComponentAdd)
	r.GET(rt.ComponentEditURL, m.ComponentEdit)
	r.POST(rt.ComponentEditURL, m.ComponentEdit)
	r.POST(rt.ComponentDeleteURL, m.ComponentDelete)
	r.GET(rt.PostingAddURL, m.PostingAdd)
	r.POST(rt.PostingAddURL, m.PostingAdd)
	r.GET(rt.PostingEditURL, m.PostingEdit)
	r.POST(rt.PostingEditURL, m.PostingEdit)
	r.POST(rt.PostingDeleteURL, m.PostingDelete)

	// Evidence attachments.
	r.GET(rt.AttachmentUploadURL, m.AttachmentUpload)
	r.POST(rt.AttachmentUploadURL, m.AttachmentUpload)
	r.POST(rt.AttachmentDeleteURL, m.AttachmentDelete)
}
