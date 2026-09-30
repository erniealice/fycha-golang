// Package detail renders the charge policy detail page
// (info · versions · usage · audit-history).
package detail

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/erniealice/hybra-golang/views/auditlog"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
)

// Deps holds the detail view dependencies.
type Deps struct {
	Routes       cp.Routes
	Labels       cp.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *cp.UseCases

	auditlog.AuditOps
}

// PageData is the detail page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          cp.Labels
	ActiveTab       string
	TabItems        []pyeza.TabItem

	ID            string
	Code          string
	Name          string
	Description   string
	StatusLabel   string
	StatusVariant string
	CurrentLabel  string
	CreatedDate   string
	ModifiedDate  string
	RetiredDate   string
	IsRetired     bool

	// Actions (rendered disabled — never hidden — when not allowed).
	EditURL           string
	CanEdit           bool
	RetireURL         string
	CanRetire         bool
	RetireDisabled    string // tooltip when disabled
	NewVersionURL     string
	NewVersionOff     bool
	NewVersionTip     string
	MissingPermEdit   string
	MissingPermRetire string

	VersionsTable *types.TableConfig

	InUse       bool
	UsageTables []UsageTable

	AuditEntries    []auditlog.AuditEntryView
	AuditHasNext    bool
	AuditNextCursor string
	AuditHistoryURL string
}

// UsageTable is one read-only "where used" section.
type UsageTable struct {
	Title string
	Table *types.TableConfig
}

// NewView renders the full detail page.
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "read") {
			return view.Forbidden("charge_policy:read")
		}
		id := viewCtx.Request.PathValue("id")
		tab := normalizeTab(viewCtx.Request.URL.Query().Get("tab"))
		pd, status, err := buildPageData(ctx, deps, viewCtx, id, tab, perms)
		if err != nil {
			return failure(status, err)
		}
		return view.OK("charge-policy-detail", pd)
	})
}

// NewTabAction serves the HTMX tab partials.
func NewTabAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "read") {
			return view.Forbidden("charge_policy:read")
		}
		id := viewCtx.Request.PathValue("id")
		tab := normalizeTab(viewCtx.Request.PathValue("tab"))
		pd, status, err := buildPageData(ctx, deps, viewCtx, id, tab, perms)
		if err != nil {
			return failure(status, err)
		}
		name := "charge-policy-tab-" + tab
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
	case "info", "versions", "usage", "audit-history":
		return t
	default:
		return "info"
	}
}

func buildPageData(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, id, tab string, perms *types.UserPermissions) (*PageData, int, error) {
	l := deps.Labels
	if deps.UseCases == nil || deps.UseCases.ReadChargePolicy == nil {
		return nil, http.StatusServiceUnavailable, fmt.Errorf("%s", l.Errors.Unavailable)
	}
	detail, err := deps.UseCases.ReadPolicy(ctx, id)
	if err != nil || detail == nil || detail.Policy == nil {
		if err != nil && cp.ErrorKind(err) != cp.ErrKindNotFound {
			log.Printf("charge_policy detail: ReadPolicy(%s): %v", id, err)
			return nil, http.StatusInternalServerError, err
		}
		return nil, http.StatusNotFound, fmt.Errorf("%s", l.Errors.NotFound)
	}
	s := cp.Summarize(detail.Policy, detail.Versions)

	cl := deps.CommonLabels
	pd := &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          s.Policy.GetName(),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   deps.Routes.ActiveSubNav,
			HeaderTitle:    s.Policy.GetName(),
			HeaderSubtitle: s.Policy.GetCode(),
			HeaderIcon:     "icon-shield-check",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "charge-policy-detail-content",
		Labels:          l,
		ActiveTab:       tab,
		ID:              id,
		Code:            s.Policy.GetCode(),
		Name:            s.Policy.GetName(),
		Description:     s.Policy.GetDescription(),
		CreatedDate:     cp.FormatMillis(ctx, s.Policy.DateCreated),
		ModifiedDate:    cp.FormatMillis(ctx, s.Policy.DateModified),
		RetiredDate:     cp.FormatMillis(ctx, s.Policy.RetiredAt),
		IsRetired:       s.IsRetired(),
		EditURL:         route.ResolveURL(deps.Routes.EditURL, "id", id),
		CanEdit:         perms.Can("charge_policy", "update"),
		RetireURL:       route.ResolveURL(deps.Routes.RetireURL, "id", id),
		CanRetire:       perms.Can("charge_policy", "retire") && !s.IsRetired(),
		NewVersionURL:   route.ResolveURL(deps.Routes.VersionAddURL, "id", id),

		MissingPermEdit:   fmt.Sprintf(cl.Errors.MissingPermission, "charge_policy:update"),
		MissingPermRetire: fmt.Sprintf(cl.Errors.MissingPermission, "charge_policy:retire"),
	}
	switch s.Bucket {
	case "retired":
		pd.StatusLabel, pd.StatusVariant = l.Enums.PolicyStatusRetired, "muted"
	case "draft":
		pd.StatusLabel, pd.StatusVariant = l.Enums.PolicyStatusDraft, "warning"
	default:
		pd.StatusLabel, pd.StatusVariant = l.Enums.PolicyStatusActive, "success"
	}
	if s.Current != nil {
		pd.CurrentLabel = cp.Format(l.Detail.CurrentVersion, fmt.Sprintf("%d", s.Current.GetVersionNumber()))
	} else {
		pd.CurrentLabel = l.Detail.NoCurrentVersion
	}
	switch {
	case s.IsRetired():
		pd.RetireDisabled = l.Errors.Retired
	case !perms.Can("charge_policy", "retire"):
		pd.RetireDisabled = pd.MissingPermRetire
	}
	switch {
	case s.IsRetired():
		pd.NewVersionOff, pd.NewVersionTip = true, l.Errors.Retired
	case s.Draft != nil:
		pd.NewVersionOff, pd.NewVersionTip = true, l.Errors.DraftExists
	case !perms.Can("charge_policy", "update"):
		pd.NewVersionOff, pd.NewVersionTip = true, fmt.Sprintf(cl.Errors.MissingPermission, "charge_policy:update")
	}

	pd.TabItems = tabItems(deps, id, len(s.Versions))

	switch tab {
	case "versions":
		pd.VersionsTable = buildVersionsTable(ctx, deps, s, perms, pd)
	case "usage":
		buildUsage(ctx, deps, pd, id)
	case "info":
		if deps.UseCases.GetChargePolicyInUseIds != nil {
			if m, err := deps.UseCases.InUseIDs(ctx, []string{id}); err == nil {
				pd.InUse = m[id]
			}
		}
	case "audit-history":
		if deps.ListAuditHistory != nil {
			cursor := viewCtx.Request.URL.Query().Get("cursor")
			resp, err := deps.ListAuditHistory(ctx, &auditlog.ListAuditRequest{
				EntityType: "charge_policy", EntityID: id, Limit: 20, CursorToken: cursor,
			})
			if err != nil {
				log.Printf("charge_policy detail: audit history: %v", err)
			}
			if resp != nil {
				pd.AuditEntries = resp.Entries
				pd.AuditHasNext = resp.HasNext
				pd.AuditNextCursor = resp.NextCursor
			}
		}
		pd.AuditHistoryURL = route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "") + "audit-history"
	}
	return pd, 0, nil
}

func tabItems(deps *Deps, id string, versionCount int) []pyeza.TabItem {
	l := deps.Labels
	base := route.ResolveURL(deps.Routes.DetailURL, "id", id)
	action := route.ResolveURL(deps.Routes.TabActionURL, "id", id, "tab", "")
	return []pyeza.TabItem{
		{Key: "info", Label: l.Tabs.Info, Href: base + "?tab=info", HxGet: action + "info", Icon: "icon-info"},
		{Key: "versions", Label: l.Tabs.Versions, Href: base + "?tab=versions", HxGet: action + "versions", Icon: "icon-layers", Count: versionCount},
		{Key: "usage", Label: l.Tabs.Usage, Href: base + "?tab=usage", HxGet: action + "usage", Icon: "icon-link"},
		{Key: "audit-history", Label: l.Tabs.History, Href: base + "?tab=audit-history", HxGet: action + "audit-history", Icon: "icon-clock"},
	}
}

func buildVersionsTable(ctx context.Context, deps *Deps, s cp.Summary, perms *types.UserPermissions, pd *PageData) *types.TableConfig {
	l := deps.Labels
	id := s.Policy.GetId()
	columns := []types.TableColumn{
		{Key: "version", Label: l.Columns.Version, NoSort: true, NoFilter: true, WidthClass: "col-lg"},
		{Key: "status", Label: l.Columns.Status, NoSort: true, NoFilter: true, WidthClass: "col-xl"},
		{Key: "accounting_role", Label: l.Columns.AccountingRole, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "document_kind", Label: l.Columns.DocumentKind, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "prepared_by", Label: l.Columns.PreparedBy, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "approved_by", Label: l.Columns.ApprovedBy, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "approved_on", Label: l.Columns.ApprovedOn, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
	}
	rows := make([]types.TableRow, 0, len(s.Versions))
	for _, v := range s.Versions {
		statusLabel, variant := l.VersionStatusLabel(v.GetStatus())
		role := ""
		if v.AccountingRole != nil {
			role = l.AccountingRoleLabel(v.GetAccountingRole())
		}
		doc := documentText(ctx, deps, v)
		href := route.ResolveURL(deps.Routes.VersionDetailURL, "id", id, "vid", v.GetId())
		rows = append(rows, types.TableRow{
			ID:   v.GetId(),
			Href: href,
			DataAttrs: map[string]string{
				"testid": "charge-policy-version-row-" + v.GetId(),
			},
			Cells: []types.TableCell{
				{Type: "link", Value: fmt.Sprintf("v%d", v.GetVersionNumber()), Href: href},
				{Type: "badge", Value: statusLabel, Variant: variant},
				{Type: "text", Value: role},
				{Type: "text", Value: doc},
				{Type: "text", Value: v.GetPreparedBy()},
				{Type: "text", Value: v.GetApprovedBy()},
				{Type: "text", Value: cp.FormatMillis(ctx, v.ApprovedAt)},
			},
			Actions: []types.TableAction{{Type: "view", Label: l.Buttons.View, Action: "view", Href: href}},
		})
	}
	types.ApplyColumnStyles(columns, rows)
	table := &types.TableConfig{
		ID:          "charge-policy-versions-table",
		Columns:     columns,
		Rows:        rows,
		ShowSearch:  false,
		ShowActions: true,
		ShowEntries: true,
		Labels:      deps.TableLabels,
		EmptyState:  types.TableEmptyState{Title: l.Tabs.Versions, Message: l.Detail.NoCurrentVersion},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.NewVersion,
			ActionURL:       pd.NewVersionURL,
			Icon:            "icon-plus",
			TestID:          "charge-policy-new-version",
			Disabled:        pd.NewVersionOff,
			DisabledTooltip: pd.NewVersionTip,
		},
	}
	types.ApplyTableSettings(table)
	return table
}

func documentText(ctx context.Context, deps *Deps, v *versionpb.ChargePolicyVersion) string {
	if deps.UseCases.ReadChargePolicyVersion == nil {
		return ""
	}
	d, err := deps.UseCases.ReadVersion(ctx, v.GetId())
	if err != nil || d == nil {
		return ""
	}
	seen := map[enumspb.ChargeDocumentKind]bool{}
	out := ""
	for _, c := range d.Components {
		if seen[c.GetDocumentKind()] {
			continue
		}
		seen[c.GetDocumentKind()] = true
		if out != "" {
			out += " + "
		}
		out += deps.Labels.DocumentKindLabel(c.GetDocumentKind())
	}
	return out
}

func buildUsage(ctx context.Context, deps *Deps, pd *PageData, id string) {
	if deps.UseCases.GetChargePolicyInUseIds != nil {
		if m, err := deps.UseCases.InUseIDs(ctx, []string{id}); err == nil {
			pd.InUse = m[id]
		}
	}
}
