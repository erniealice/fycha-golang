// Package version renders the charge policy version page
// (classification · components · postings · approval · audit-history).
package version

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	attachmentpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/document/attachment"
	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	"github.com/erniealice/hybra-golang/views/attachment"
	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
)

// Deps holds the version page dependencies.
type Deps struct {
	Routes       cp.Routes
	Labels       cp.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *cp.UseCases

	attachment.AttachmentOps
	auditlog.AuditOps
}

// requiredPostings are the posting pairs the approval checklist demands.
var requiredPostings = []struct {
	Event enumspb.ChargePostingEvent
	Role  enumspb.ChargePostingRole
}{
	{enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ISSUE, enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_RECEIVABLE},
	{enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ISSUE, enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_CLEARING},
	{enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_APPLICATION, enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_CASH},
	{enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_APPLICATION, enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_RECEIVABLE},
}

// ChecklistRow is one rendered approval-checklist line.
type ChecklistRow struct {
	Text   string
	Passed bool
	Detail string
}

// Field is one classification info-grid entry.
type Field struct {
	Label  string
	Value  string
	TestID string
}

// PageData is the version page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          cp.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	PolicyID      string
	PolicyName    string
	PolicyURL     string
	BackLabel     string
	VersionID     string
	VersionNumber int32
	IsDraft       bool
	StatusLabel   string
	StatusVariant string
	ShowLocked    bool
	SelfApproved  bool

	// Version-level actions (rendered disabled, never hidden).
	CanEdit       bool
	EditURL       string
	EditTip       string
	DeleteURL     string
	CanDelete     bool
	DeleteTip     string
	MissingUpdate string

	Fields          []Field
	ComponentsTable *types.TableConfig
	PostingsTable   *types.TableConfig

	// Approval.
	Validation      *versionpb.ChargePolicyApprovalChecklist
	Checklist       []ChecklistRow
	ChecklistReady  bool // validation ran
	Reasons         []string
	ReasonsText     string
	ApproveURL      string
	ApproveOff      bool
	ApproveTip      string
	ApprovedBy      string
	ApprovedOn      string
	AttachmentTable *types.TableConfig

	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// NewView renders the full version page.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "read") {
			return view.Forbidden("charge_policy:read")
		}
		tab := normalizeTab(viewCtx.Request.URL.Query().Get("tab"))
		pd, status, err := buildPageData(ctx, deps, viewCtx, tab, perms)
		if err != nil {
			return failure(status, err)
		}
		return view.OK("charge-policy-version-detail", pd)
	})
}

// NewTabAction serves the HTMX tab partials.
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "read") {
			return view.Forbidden("charge_policy:read")
		}
		rawTab := viewCtx.Request.PathValue("tab")
		tab := normalizeTab(rawTab)
		if rawTab == "attachments" {
			tab = "approval" // the attachments table refresh target
		}
		pd, status, err := buildPageData(ctx, deps, viewCtx, tab, perms)
		if err != nil {
			return failure(status, err)
		}
		// Table-only refresh (row delete / drawer save): return the table card.
		if rawTab == "attachments" && pd.AttachmentTable != nil {
			return view.OK("table-card", pd.AttachmentTable)
		}
		if viewCtx.Request.URL.Query().Get("table") == "1" {
			switch {
			case tab == "components" && pd.ComponentsTable != nil:
				return view.OK("table-card", pd.ComponentsTable)
			case tab == "postings" && pd.PostingsTable != nil:
				return view.OK("table-card", pd.PostingsTable)
			}
		}
		name := "charge-policy-version-tab-" + tab
		if tab == "audit-history" {
			name = "audit-history-tab"
		}
		return view.OK(name, pd)
	})
}

func failure(status int, err error) view.ViewResult {
	r := view.Error(err)
	if status != 0 {
		r.StatusCode = status
	}
	return r
}

func normalizeTab(t string) string {
	switch t {
	case "classification", "components", "postings", "approval", "audit-history":
		return t
	default:
		return "classification"
	}
}

// AttachmentConfig builds the evidence-attachment config (shared with the
// upload/delete action handlers). The explicit Policy allows PDF and images.
func AttachmentConfig(deps *Deps) *attachment.Config {
	return &attachment.Config{
		EntityType:         "charge_policy_version",
		BucketName:         "attachments",
		RefreshURL:         deps.Routes.VersionTabActionURL,
		UploadURL:          deps.Routes.AttachmentUploadURL,
		DeleteURL:          deps.Routes.AttachmentDeleteURL,
		PrimaryIDPathParam: "vid",
		Labels:             attachment.DefaultLabels(),
		CommonLabels:       deps.CommonLabels,
		TableLabels:        deps.TableLabels,
		NewID:              deps.NewAttachmentID,
		UploadFile:         deps.UploadFile,
		ListAttachments:    deps.ListAttachments,
		CreateAttachment:   deps.CreateAttachment,
		DeleteAttachment:   deps.DeleteAttachment,
		Policy: &attachment.Policy{
			AllowedContentTypes: []string{"application/pdf", "image/png", "image/jpeg"},
			AllowedExtensions:   []string{".pdf", ".png", ".jpg", ".jpeg"},
			MaxFileCount:        20,
		},
	}
}

func buildPageData(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, tab string, perms *types.UserPermissions) (*PageData, int, error) {
	l := deps.Labels
	cl := deps.CommonLabels
	if deps.UseCases == nil || deps.UseCases.ReadChargePolicyVersion == nil || deps.UseCases.ReadChargePolicy == nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("%s", l.Errors.Unavailable)
	}
	policyID := viewCtx.Request.PathValue("id")
	vid := viewCtx.Request.PathValue("vid")

	vd, err := deps.UseCases.ReadVersion(ctx, vid)
	if err != nil || vd == nil || vd.Version == nil || vd.Version.GetChargePolicyId() != policyID {
		if err != nil && cp.ErrorKind(err) != cp.ErrKindNotFound {
			log.Printf("charge_policy version: ReadVersion(%s): %v", vid, err)
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, fmt.Errorf("%s", l.Errors.NotFound)
	}
	pol, err := deps.UseCases.ReadPolicy(ctx, policyID)
	if err != nil || pol == nil || pol.Policy == nil {
		return nil, http.StatusNotFound, fmt.Errorf("%s", l.Errors.NotFound)
	}
	v := vd.Version
	statusLabel, variant := l.VersionStatusLabel(v.GetStatus())
	title := cp.Format(l.Detail.VersionTitle, pol.Policy.GetName(), fmt.Sprintf("%d", v.GetVersionNumber()))
	draft := cp.IsDraft(v)
	retired := pol.Policy.GetStatus() == enumspb.ChargePolicyStatus_CHARGE_POLICY_STATUS_RETIRED

	pd := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          title,
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    title,
			HeaderSubtitle: pol.Policy.GetCode(),
			HeaderIcon:     "icon-shield-check",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "charge-policy-version-detail-content",
		Labels:          l,
		ActiveTab:       tab,
		PolicyID:        policyID,
		PolicyName:      pol.Policy.GetName(),
		PolicyURL:       route.ResolveURL(deps.Routes.DetailURL, "id", policyID),
		BackLabel:       cp.Format(l.Detail.BackToPolicy, pol.Policy.GetName()),
		VersionID:       vid,
		VersionNumber:   v.GetVersionNumber(),
		IsDraft:         draft,
		StatusLabel:     statusLabel,
		StatusVariant:   variant,
		ShowLocked:      !draft,
		SelfApproved:    v.GetSelfApproved(),
		MissingUpdate:   fmt.Sprintf(cl.Errors.MissingPermission, "charge_policy:update"),
		ApproveURL:      route.ResolveURL(deps.Routes.VersionApproveURL, "id", policyID, "vid", vid),
		EditURL:         route.ResolveURL(deps.Routes.VersionEditURL, "id", policyID, "vid", vid),
		DeleteURL:       route.ResolveURL(deps.Routes.VersionDeleteURL, "id", policyID, "vid", vid),
		ApprovedBy:      v.GetApprovedBy(),
		ApprovedOn:      cp.FormatMillis(ctx, v.ApprovedAt),
	}
	switch {
	case !draft:
		pd.EditTip, pd.DeleteTip = l.Errors.NotDraft, l.Errors.NotDraft
	case !perms.Can("charge_policy", "update"):
		pd.EditTip, pd.DeleteTip = pd.MissingUpdate, pd.MissingUpdate
	default:
		pd.CanEdit = true
	}
	pd.CanDelete = draft && perms.Can("charge_policy", "delete")
	if draft && !perms.Can("charge_policy", "delete") {
		pd.DeleteTip = fmt.Sprintf(cl.Errors.MissingPermission, "charge_policy:delete")
	}

	pd.TabItems = tabItems(deps, policyID, vid, len(vd.Components), len(vd.Postings))
	pd.Fields = classificationFields(deps, v)

	// The checklist card is shown on every tab while the version is a draft.
	if draft && deps.UseCases.ValidateChargePolicyVersionForApproval != nil {
		if val, err := deps.UseCases.Validate(ctx, vid); err != nil {
			log.Printf("charge_policy version: Validate(%s): %v", vid, err)
		} else if val != nil {
			pd.Validation = val
			pd.ChecklistReady = true
			pd.Checklist = checklistRows(l, val)
			for _, r := range val.GetReasons() {
				pd.Reasons = append(pd.Reasons, reasonText(l, r))
			}
		}
	}
	pd.ReasonsText = strings.Join(pd.Reasons, " ")
	pd.ApproveOff, pd.ApproveTip = approveState(ctx, deps, perms, v, pd.Validation, draft, retired)

	switch tab {
	case "components":
		pd.ComponentsTable = buildComponentsTable(deps, pd, vd, perms)
	case "postings":
		pd.PostingsTable = buildPostingsTable(ctx, deps, pd, vd, perms)
	case "approval":
		if deps.ListAttachments != nil {
			cfg := AttachmentConfig(deps)
			resp, err := deps.ListAttachments(ctx, cfg.EntityType, vid)
			if err != nil {
				log.Printf("charge_policy version: list attachments: %v", err)
			}
			var items []*attachmentpb.Attachment
			if resp != nil {
				items = resp.GetData()
			}
			pd.AttachmentTable = attachment.BuildTable(items, cfg, policyID, "vid", vid)
			if pd.AttachmentTable != nil && pd.AttachmentTable.PrimaryAction != nil && !pd.CanEdit {
				pd.AttachmentTable.PrimaryAction.Disabled = true
				pd.AttachmentTable.PrimaryAction.DisabledTooltip = pd.EditTip
			}
		}
	case "audit-history":
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			resp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
				EntityType: "charge_policy_version", EntityID: vid, Limit: 20, CursorToken: cursor,
			})
			if err != nil {
				log.Printf("charge_policy version: audit history: %v", err)
			}
			if resp != nil {
				pd.AuditEntries = resp.Entries
				pd.AuditHasNext = resp.HasNext
				pd.AuditNextCursor = resp.NextCursor
			}
		}
		pd.AuditHistoryURL = route.ResolveURL(deps.Routes.VersionTabActionURL, "id", policyID, "vid", vid, "tab", "") + "audit-history"
	}
	return pd, 0, nil
}

// approveState decides whether the Approve button is disabled and why. The
// authoritative gate is the use case (strict permission, SoD, checklist, S1
// matrix); this only mirrors it so the operator sees the reason.
func approveState(ctx context.Context, deps *Deps, perms *types.UserPermissions, v interface {
	GetPreparedBy() string
}, val *versionpb.ChargePolicyApprovalChecklist, draft, retired bool) (off bool, tip string) {
	l := deps.Labels
	switch {
	case !draft:
		return true, l.Errors.NotDraft
	case retired:
		return true, l.Errors.Retired
	case !perms.Can("charge_policy", "approve"):
		return true, fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "charge_policy:approve")
	case v.GetPreparedBy() != "" && v.GetPreparedBy() == cp.ActorID(ctx) && !perms.Can("charge_policy", "approve_own"):
		return true, l.Errors.SelfApproval
	case deps.UseCases == nil || deps.UseCases.ApproveChargePolicyVersion == nil:
		return true, l.Errors.Unavailable
	case val == nil:
		return true, l.Errors.ChecklistFailed
	case !val.GetSupported():
		return true, l.Errors.UnsupportedCombination
	case !val.GetPassed():
		return true, l.Errors.ChecklistFailed
	}
	return false, ""
}

func tabItems(deps *Deps, policyID, vid string, components, postings int) []pyeza.TabItem {
	l := deps.Labels
	base := route.ResolveURL(deps.Routes.VersionDetailURL, "id", policyID, "vid", vid)
	action := route.ResolveURL(deps.Routes.VersionTabActionURL, "id", policyID, "vid", vid, "tab", "")
	return []pyeza.TabItem{
		{Key: "classification", Label: l.Tabs.Classification, Href: base + "?tab=classification", HxGet: action + "classification", Icon: "icon-info"},
		{Key: "components", Label: l.Tabs.Components, Href: base + "?tab=components", HxGet: action + "components", Icon: "icon-file-text", Count: components},
		{Key: "postings", Label: l.Tabs.Postings, Href: base + "?tab=postings", HxGet: action + "postings", Icon: "icon-book", Count: postings},
		{Key: "approval", Label: l.Tabs.Approval, Href: base + "?tab=approval", HxGet: action + "approval", Icon: "icon-check-circle"},
		{Key: "audit-history", Label: l.Tabs.History, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

func orDash(s string) string {
	if strings.TrimSpace(s) == "" {
		return "—"
	}
	return s
}

func classificationFields(deps *Deps, v interface {
	GetAccountingRole() enumspb.AccountingRole
	GetBookPresentation() enumspb.BookPresentation
	GetTaxPosition() enumspb.TaxPosition
	GetTaxTreatmentId() string
	GetAssessmentScope() string
	GetAssessmentNote() string
	GetClonedFromVersionId() string
}) []Field {
	l := deps.Labels
	return []Field{
		{Label: l.Form.AccountingRoleLabel, Value: orDash(l.AccountingRoleLabel(v.GetAccountingRole())), TestID: "cpv-info-accounting-role"},
		{Label: l.Form.BookPresentationLabel, Value: orDash(l.BookPresentationLabel(v.GetBookPresentation())), TestID: "cpv-info-book-presentation"},
		{Label: l.Form.TaxPositionLabel, Value: orDash(l.TaxPositionLabel(v.GetTaxPosition())), TestID: "cpv-info-tax-position"},
		{Label: l.Form.TaxTreatmentLabel, Value: orDash(v.GetTaxTreatmentId()), TestID: "cpv-info-tax-treatment"},
		{Label: l.Form.AssessmentScopeLabel, Value: orDash(v.GetAssessmentScope()), TestID: "cpv-info-assessment-scope"},
		{Label: l.Form.AssessmentNoteLabel, Value: orDash(v.GetAssessmentNote()), TestID: "cpv-info-assessment-note"},
		{Label: l.Detail.ClonedFrom, Value: orDash(v.GetClonedFromVersionId()), TestID: "cpv-info-cloned-from"},
	}
}

func checklistRows(l cp.Labels, val *versionpb.ChargePolicyApprovalChecklist) []ChecklistRow {
	rows := make([]ChecklistRow, 0, len(val.GetItems()))
	for _, it := range val.GetItems() {
		rows = append(rows, ChecklistRow{Text: checklistText(l, it.GetCode()), Passed: it.GetPassed(), Detail: it.GetDetail()})
	}
	return rows
}

// checklistText renders a checklist code (assessment_scope, posting:E:R,
// posting_account:E:R) as operator text.
func checklistText(l cp.Labels, code string) string {
	if t, ok := l.Detail.ChecklistCodeLabel[code]; ok {
		return t
	}
	parts := strings.Split(code, ":")
	if len(parts) == 3 && (parts[0] == "posting" || parts[0] == "posting_account") {
		ev := enumspb.ChargePostingEvent(enumspb.ChargePostingEvent_value["CHARGE_POSTING_EVENT_"+parts[1]])
		ro := enumspb.ChargePostingRole(enumspb.ChargePostingRole_value["CHARGE_POSTING_ROLE_"+parts[2]])
		text := l.EventLabel(ev) + ": " + l.PostingRoleLabel(ro)
		if parts[0] == "posting_account" {
			return text + " (" + l.Columns.Account + ")"
		}
		return text
	}
	return code
}

func reasonText(l cp.Labels, code string) string {
	if t, ok := l.Detail.ReasonLabel[code]; ok {
		return t
	}
	return l.Errors.UnsupportedCombination
}

func buildComponentsTable(deps *Deps, pd *PageData, vd *cp.VersionDetail, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "component_role", Label: l.Columns.ComponentRole, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "document_kind", Label: l.Columns.DocumentKind, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "book_presentation", Label: l.Columns.BookPresentation, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
	}
	rows := make([]types.TableRow, 0, len(vd.Components))
	for _, c := range vd.Components {
		id := c.GetId()
		edit := types.TableAction{Type: "edit", Label: l.Buttons.Edit, Action: "edit",
			URL:         route.ResolveURL(deps.Routes.ComponentEditURL, "id", pd.PolicyID, "vid", pd.VersionID, "cid", id),
			DrawerTitle: l.Buttons.Edit, Disabled: !pd.CanEdit, DisabledTooltip: pd.EditTip}
		del := types.TableAction{Type: "delete", Label: l.Buttons.Delete, Action: "delete",
			URL:      route.ResolveURL(deps.Routes.ComponentDeleteURL, "id", pd.PolicyID, "vid", pd.VersionID),
			ItemName: l.ComponentRoleLabel(c.GetComponentRole()), ConfirmTitle: l.Confirm.DeleteTitle, ConfirmMessage: l.Confirm.DeleteMsg,
			Disabled: !pd.CanEdit, DisabledTooltip: pd.EditTip}
		rows = append(rows, types.TableRow{
			ID:        id,
			DataAttrs: map[string]string{"testid": "charge-policy-component-row-" + id},
			Cells: []types.TableCell{
				{Type: "text", Value: l.ComponentRoleLabel(c.GetComponentRole())},
				{Type: "text", Value: l.DocumentKindLabel(c.GetDocumentKind())},
				{Type: "text", Value: l.BookPresentationLabel(c.GetBookPresentation())},
			},
			Actions: []types.TableAction{edit, del},
		})
	}
	types.ApplyColumnStyles(columns, rows)
	table := &types.TableConfig{
		ID: "charge-policy-components-table", Columns: columns, Rows: rows,
		RefreshURL:  route.ResolveURL(deps.Routes.VersionTabActionURL, "id", pd.PolicyID, "vid", pd.VersionID, "tab", "components") + "?table=1",
		ShowActions: true, ShowEntries: true, Labels: deps.TableLabels,
		EmptyState: types.TableEmptyState{Title: l.Tabs.Components},
		PrimaryAction: &types.PrimaryAction{
			Label: l.Buttons.AddComponent, ActionURL: route.ResolveURL(deps.Routes.ComponentAddURL, "id", pd.PolicyID, "vid", pd.VersionID),
			Icon: "icon-plus", TestID: "charge-policy-component-add",
			Disabled: !pd.CanEdit, DisabledTooltip: pd.EditTip,
		},
	}
	types.ApplyTableSettings(table)
	return table
}

func buildPostingsTable(ctx context.Context, deps *Deps, pd *PageData, vd *cp.VersionDetail, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	accounts := accountNames(ctx, deps)
	columns := []types.TableColumn{
		{Key: "event", Label: l.Columns.Event, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "posting_role", Label: l.Columns.PostingRole, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "account", Label: l.Columns.Account, NoSort: true, NoFilter: true},
	}
	type key struct {
		e enumspb.ChargePostingEvent
		r enumspb.ChargePostingRole
	}
	have := map[key]bool{}
	rows := make([]types.TableRow, 0, len(vd.Postings)+len(requiredPostings))
	for _, p := range vd.Postings {
		have[key{p.GetEvent(), p.GetPostingRole()}] = true
		id := p.GetId()
		name := accounts[p.GetAccountId()]
		if name == "" {
			name = p.GetAccountId()
		}
		edit := types.TableAction{Type: "edit", Label: l.Buttons.Edit, Action: "edit",
			URL:         route.ResolveURL(deps.Routes.PostingEditURL, "id", pd.PolicyID, "vid", pd.VersionID, "pid", id),
			DrawerTitle: l.Buttons.Edit, Disabled: !pd.CanEdit, DisabledTooltip: pd.EditTip}
		del := types.TableAction{Type: "delete", Label: l.Buttons.Delete, Action: "delete",
			URL:          route.ResolveURL(deps.Routes.PostingDeleteURL, "id", pd.PolicyID, "vid", pd.VersionID),
			ItemName:     l.EventLabel(p.GetEvent()) + " / " + l.PostingRoleLabel(p.GetPostingRole()),
			ConfirmTitle: l.Confirm.DeleteTitle, ConfirmMessage: l.Confirm.DeleteMsg,
			Disabled: !pd.CanEdit, DisabledTooltip: pd.EditTip}
		rows = append(rows, types.TableRow{
			ID:        id,
			DataAttrs: map[string]string{"testid": "charge-policy-posting-row-" + id},
			Cells: []types.TableCell{
				{Type: "text", Value: l.EventLabel(p.GetEvent())},
				{Type: "text", Value: l.PostingRoleLabel(p.GetPostingRole())},
				{Type: "text", Value: name},
			},
			Actions: []types.TableAction{edit, del},
		})
	}
	// A missing required pair shows as a red "Missing" row (draft only).
	if pd.IsDraft {
		for _, rp := range requiredPostings {
			if have[key{rp.Event, rp.Role}] {
				continue
			}
			rid := "missing-" + rp.Event.String() + "-" + rp.Role.String()
			rows = append(rows, types.TableRow{
				ID:        rid,
				DataAttrs: map[string]string{"testid": "charge-policy-posting-missing-" + strings.ToLower(strings.TrimPrefix(rp.Event.String(), "CHARGE_POSTING_EVENT_")) + "-" + strings.ToLower(strings.TrimPrefix(rp.Role.String(), "CHARGE_POSTING_ROLE_"))},
				Cells: []types.TableCell{
					{Type: "text", Value: l.EventLabel(rp.Event)},
					{Type: "text", Value: l.PostingRoleLabel(rp.Role)},
					{Type: "badge", Value: l.Detail.MissingMapping, Variant: "danger"},
				},
			})
		}
	}
	types.ApplyColumnStyles(columns, rows)
	table := &types.TableConfig{
		ID: "charge-policy-postings-table", Columns: columns, Rows: rows,
		RefreshURL:  route.ResolveURL(deps.Routes.VersionTabActionURL, "id", pd.PolicyID, "vid", pd.VersionID, "tab", "postings") + "?table=1",
		ShowActions: true, ShowEntries: true, Labels: deps.TableLabels,
		EmptyState: types.TableEmptyState{Title: l.Tabs.Postings},
		PrimaryAction: &types.PrimaryAction{
			Label: l.Buttons.AddPosting, ActionURL: route.ResolveURL(deps.Routes.PostingAddURL, "id", pd.PolicyID, "vid", pd.VersionID),
			Icon: "icon-plus", TestID: "charge-policy-posting-add",
			Disabled: !pd.CanEdit, DisabledTooltip: pd.EditTip,
		},
	}
	types.ApplyTableSettings(table)
	return table
}

// accountNames maps account id → "code · name" for display and pickers.
func accountNames(ctx context.Context, deps *Deps) map[string]string {
	out := map[string]string{}
	if deps.UseCases == nil || deps.UseCases.GetAccountListPageData == nil {
		return out
	}
	accs, err := deps.UseCases.ListAccounts(ctx)
	if err != nil {
		log.Printf("charge_policy version: ListAccounts: %v", err)
		return out
	}
	for _, a := range accs {
		out[a.GetId()] = AccountLabel(a)
	}
	return out
}

// AccountLabel renders "code · name" for an account.
func AccountLabel(a *accountpb.Account) string {
	if a.GetCode() != "" {
		return a.GetCode() + " · " + a.GetName()
	}
	return a.GetName()
}
