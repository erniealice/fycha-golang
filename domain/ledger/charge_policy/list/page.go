// Package list renders the Charge Policies list (Active · Draft · Retired).
package list

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	espynahttp "github.com/erniealice/espyna-golang/contrib/http"
	"github.com/erniealice/espyna-golang/shared/tableparams"
	versionpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/charge_policy_version"
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	cp "github.com/erniealice/fycha-golang/domain/ledger/charge_policy"
)

// Deps holds the list view dependencies.
type Deps struct {
	Routes       cp.Routes
	Labels       cp.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *cp.UseCases
}

// PageData is the list page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	StatusTabs      []pyeza.TabItem
	ActiveStatus    string
}

// NewView renders the full page (or the HTMX content partial).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "list") {
			return view.Forbidden("charge_policy:list")
		}
		if deps.UseCases == nil || deps.UseCases.GetChargePolicyListPageData == nil {
			return view.Error(fmt.Errorf("%s", deps.Labels.Errors.Unavailable))
		}
		status := cp.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, tabs, err := build(ctx, deps, viewCtx, status, perms)
		if err != nil {
			return view.Error(err)
		}
		title := statusTitle(deps.Labels, status)
		return view.OK("charge-policy-list", &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    title,
				HeaderSubtitle: deps.Labels.Page.Subtitle,
				HeaderIcon:     "icon-shield-check",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "charge-policy-list-content",
			Table:           table,
			StatusTabs:      tabs,
			ActiveStatus:    status,
		})
	})
}

// NewTableView renders only the table card (HTMX refresh).
func NewTableView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("charge_policy", "list") {
			return view.Forbidden("charge_policy:list")
		}
		if deps.UseCases == nil || deps.UseCases.GetChargePolicyListPageData == nil {
			return view.Error(fmt.Errorf("%s", deps.Labels.Errors.Unavailable))
		}
		status := cp.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, _, err := build(ctx, deps, viewCtx, status, perms)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("table-card", table)
	})
}

func statusTitle(l cp.Labels, status string) string {
	switch status {
	case "draft":
		return l.Page.TitleDraft
	case "retired":
		return l.Page.TitleRetired
	default:
		return l.Page.TitleActive
	}
}

func statusTabLabel(l cp.Labels, status string) string {
	switch status {
	case "draft":
		return l.Tabs.Draft
	case "retired":
		return l.Tabs.Retired
	default:
		return l.Tabs.Active
	}
}

func emptyTitle(l cp.Labels, status string) string {
	switch status {
	case "draft":
		return l.Empty.DraftTitle
	case "retired":
		return l.Empty.RetiredTitle
	default:
		return l.Empty.ActiveTitle
	}
}

func columnsFor(l cp.Labels, status string) []types.TableColumn {
	switch status {
	case "draft":
		return []types.TableColumn{
			{Key: "code", Label: l.Columns.Code, WidthClass: "col-3xl"},
			{Key: "name", Label: l.Columns.Name},
			{Key: "draft_version", Label: l.Columns.DraftVersion, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
			{Key: "prepared_by", Label: l.Columns.PreparedBy, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
			{Key: "modified", Label: l.Columns.Modified, WidthClass: "col-2xl"},
		}
	case "retired":
		return []types.TableColumn{
			{Key: "code", Label: l.Columns.Code, WidthClass: "col-3xl"},
			{Key: "name", Label: l.Columns.Name},
			{Key: "last_version", Label: l.Columns.LastVersion, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
			{Key: "retired_on", Label: l.Columns.RetiredOn, WidthClass: "col-2xl"},
			{Key: "retired_by", Label: l.Columns.RetiredBy, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		}
	default:
		return []types.TableColumn{
			{Key: "code", Label: l.Columns.Code, WidthClass: "col-3xl"},
			{Key: "name", Label: l.Columns.Name},
			{Key: "accounting_role", Label: l.Columns.AccountingRole, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
			{Key: "document_kind", Label: l.Columns.DocumentKind, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
			{Key: "current_version", Label: l.Columns.CurrentVersion, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
			{Key: "in_use", Label: l.Columns.InUse, NoSort: true, NoFilter: true, WidthClass: "col-xl"},
			{Key: "modified", Label: l.Columns.Modified, WidthClass: "col-2xl"},
		}
	}
}

func build(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, status string, perms *types.UserPermissions) (*types.TableConfig, []pyeza.TabItem, error) {
	l := deps.Labels
	columns := columnsFor(l, status)
	p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "code", "asc")
	if err != nil {
		return nil, nil, err
	}

	policies, err := deps.UseCases.ListPolicies(ctx)
	if err != nil {
		log.Printf("charge_policy list: ListPolicies: %v", err)
		return nil, nil, err
	}

	summaries := make([]cp.Summary, 0, len(policies))
	counts := map[string]int{}
	for _, pol := range policies {
		if pol == nil {
			continue
		}
		versions := versionsOf(ctx, deps, pol.GetId())
		s := cp.Summarize(pol, versions)
		counts[s.Bucket]++
		if s.Bucket == status {
			summaries = append(summaries, s)
		}
	}
	summaries = applyQuery(summaries, p)

	inUse := map[string]bool{}
	if status == "active" && deps.UseCases.GetChargePolicyInUseIds != nil {
		ids := make([]string, 0, len(summaries))
		for _, s := range summaries {
			ids = append(ids, s.Policy.GetId())
		}
		if m, err := deps.UseCases.InUseIDs(ctx, ids); err == nil && m != nil {
			inUse = m
		}
	}
	// Draft rows need the in-use flag for the Delete disabled state as well.
	if status == "draft" && deps.UseCases.GetChargePolicyInUseIds != nil {
		ids := make([]string, 0, len(summaries))
		for _, s := range summaries {
			ids = append(ids, s.Policy.GetId())
		}
		if m, err := deps.UseCases.InUseIDs(ctx, ids); err == nil && m != nil {
			inUse = m
		}
	}

	rows := buildRows(ctx, deps, summaries, status, inUse, perms)
	types.ApplyColumnStyles(columns, rows)

	bulk := pyeza.MapBulkConfig(deps.CommonLabels)
	bulk.Actions = nil
	if status == "draft" {
		bulk.Actions = []types.BulkAction{{
			Key:              "delete",
			Label:            deps.CommonLabels.Bulk.Delete,
			Icon:             "icon-trash-2",
			Variant:          "danger",
			Endpoint:         deps.Routes.BulkDeleteURL,
			ConfirmTitle:     l.Confirm.DeleteTitle,
			ConfirmMessage:   l.Confirm.DeleteMsg,
			RequiresDataAttr: "deletable",
		}}
	}

	table := &types.TableConfig{
		ID:                   "charge-policies-table",
		RefreshURL:           route.ResolveURL(deps.Routes.TableURL, "status", status),
		Columns:              columns,
		Rows:                 rows,
		ShowSearch:           true,
		ShowActions:          true,
		ShowFilters:          false,
		ShowSort:             true,
		ShowColumns:          true,
		ShowExport:           false,
		ShowDensity:          true,
		ShowEntries:          true,
		DefaultSortColumn:    "code",
		DefaultSortDirection: "asc",
		Labels:               deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   emptyTitle(l, status),
			Message: l.Empty.Message,
		},
	}
	if len(bulk.Actions) > 0 {
		table.BulkActions = &bulk
	}
	// "Add" is offered on Active and Draft; disabled (never hidden) without the permission.
	if status != "retired" {
		table.PrimaryAction = &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			TestID:          "charge-policy-add",
			Disabled:        !perms.Can("charge_policy", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "charge_policy:create"),
		}
	}
	types.ApplyTableSettings(table)

	tabs := make([]pyeza.TabItem, 0, len(cp.Statuses))
	for _, st := range cp.Statuses {
		tabs = append(tabs, pyeza.TabItem{
			Key:   st,
			Label: statusTabLabel(l, st),
			Href:  route.ResolveURL(deps.Routes.ListURL, "status", st),
			Count: counts[st],
		})
	}
	return table, tabs, nil
}

func versionsOf(ctx context.Context, deps *Deps, policyID string) []*versionpb.ChargePolicyVersion {
	if deps.UseCases.ListChargePolicyVersions == nil {
		return nil
	}
	vs, err := deps.UseCases.ListVersions(ctx, policyID)
	if err != nil {
		log.Printf("charge_policy list: ListVersions(%s): %v", policyID, err)
		return nil
	}
	return vs
}

// applyQuery applies search and sort in memory (the port returns the whole
// workspace set; a workspace holds a handful of policies).
func applyQuery(in []cp.Summary, p tableparams.TableQueryParams) []cp.Summary {
	out := in
	if q := strings.ToLower(strings.TrimSpace(p.Search)); q != "" {
		out = out[:0:0]
		for _, s := range in {
			if strings.Contains(strings.ToLower(s.Policy.GetCode()), q) || strings.Contains(strings.ToLower(s.Policy.GetName()), q) {
				out = append(out, s)
			}
		}
	}
	key := func(s cp.Summary) string {
		switch p.SortColumn {
		case "name":
			return strings.ToLower(s.Policy.GetName())
		case "modified":
			return fmt.Sprintf("%020d", s.Policy.GetDateModified())
		case "retired_on":
			return fmt.Sprintf("%020d", s.Policy.GetRetiredAt())
		default:
			return strings.ToLower(s.Policy.GetCode())
		}
	}
	desc := p.SortDir == "desc"
	sort.SliceStable(out, func(i, j int) bool {
		a, b := key(out[i]), key(out[j])
		if a == b {
			return out[i].Policy.GetId() < out[j].Policy.GetId()
		}
		if desc {
			return a > b
		}
		return a < b
	})
	return out
}

func buildRows(ctx context.Context, deps *Deps, summaries []cp.Summary, status string, inUse map[string]bool, perms *types.UserPermissions) []types.TableRow {
	l := deps.Labels
	rows := make([]types.TableRow, 0, len(summaries))
	for _, s := range summaries {
		id := s.Policy.GetId()
		detailURL := route.ResolveURL(deps.Routes.DetailURL, "id", id)
		codeCell := types.TableCell{Type: "link", Value: s.Policy.GetCode(), Href: detailURL}
		nameCell := types.TableCell{Type: "text", Value: s.Policy.GetName()}
		modified := types.TableCell{Type: "text", Value: cp.FormatMillis(ctx, s.Policy.DateModified)}

		var cells []types.TableCell
		switch status {
		case "draft":
			cells = []types.TableCell{codeCell, nameCell,
				{Type: "text", Value: versionText(s.Draft)},
				{Type: "text", Value: s.Draft.GetPreparedBy()},
				modified}
		case "retired":
			cells = []types.TableCell{codeCell, nameCell,
				{Type: "text", Value: versionText(s.Last)},
				{Type: "text", Value: cp.FormatMillis(ctx, s.Policy.RetiredAt)},
				{Type: "text", Value: s.Policy.GetRetiredBy()}}
		default:
			role, doc := currentClassification(ctx, deps, s)
			inUseText := l.Enums.No
			if inUse[id] {
				inUseText = l.Enums.Yes
			}
			roleCell := types.TableCell{Type: "text", Value: ""}
			if role != "" {
				roleCell = types.TableCell{Type: "badge", Value: role, Variant: "info"}
			}
			cells = []types.TableCell{codeCell, nameCell, roleCell,
				{Type: "text", Value: doc},
				{Type: "text", Value: versionText(s.Current)},
				{Type: "text", Value: inUseText},
				modified}
		}

		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"testid":    "charge-policy-row-" + id,
				"name":      s.Policy.GetName(),
				"status":    status,
				"deletable": strconv.FormatBool(status == "draft" && !inUse[id]),
			},
			Actions: rowActions(deps, s, status, inUse[id], perms),
		})
	}
	return rows
}

func versionText(v *versionpb.ChargePolicyVersion) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("v%d", v.GetVersionNumber())
}

// currentClassification returns the role badge text and the document text of
// the current approved version (document from its components).
func currentClassification(ctx context.Context, deps *Deps, s cp.Summary) (role, doc string) {
	if s.Current == nil {
		return "", ""
	}
	l := deps.Labels
	if s.Current.AccountingRole != nil {
		role = l.AccountingRoleLabel(s.Current.GetAccountingRole())
	}
	if deps.UseCases.ReadChargePolicyVersion == nil {
		return role, ""
	}
	d, err := deps.UseCases.ReadVersion(ctx, s.Current.GetId())
	if err != nil || d == nil {
		return role, ""
	}
	seen := map[enumspb.ChargeDocumentKind]bool{}
	var labels []string
	for _, c := range d.Components {
		if !seen[c.GetDocumentKind()] {
			seen[c.GetDocumentKind()] = true
			labels = append(labels, l.DocumentKindLabel(c.GetDocumentKind()))
		}
	}
	return role, strings.Join(labels, " + ")
}

func rowActions(deps *Deps, s cp.Summary, status string, isInUse bool, perms *types.UserPermissions) []types.TableAction {
	l := deps.Labels
	cl := deps.CommonLabels
	id := s.Policy.GetId()
	missing := func(code string) string { return fmt.Sprintf(cl.Errors.MissingPermission, code) }

	view := types.TableAction{Type: "view", Label: l.Buttons.View, Action: "view",
		Href: route.ResolveURL(deps.Routes.DetailURL, "id", id)}
	edit := types.TableAction{Type: "edit", Label: l.Buttons.Edit, Action: "edit",
		URL: route.ResolveURL(deps.Routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
		Disabled: !perms.Can("charge_policy", "update"), DisabledTooltip: missing("charge_policy:update")}

	switch status {
	case "retired":
		return []types.TableAction{view}
	case "draft":
		del := types.TableAction{Type: "delete", Label: l.Buttons.Delete, Action: "delete",
			URL: deps.Routes.DeleteURL, ItemName: s.Policy.GetName(),
			ConfirmTitle: l.Confirm.DeleteTitle, ConfirmMessage: l.Confirm.DeleteMsg}
		if isInUse {
			del.Disabled, del.DisabledTooltip = true, l.Errors.InUse
		} else if !perms.Can("charge_policy", "delete") {
			del.Disabled, del.DisabledTooltip = true, missing("charge_policy:delete")
		}
		return []types.TableAction{view, edit, del}
	}

	newVersion := types.TableAction{Type: "clone", Label: l.Buttons.NewVersion, Action: "clone",
		URL: route.ResolveURL(deps.Routes.VersionAddURL, "id", id), DrawerTitle: l.Buttons.NewVersion}
	switch {
	case s.Draft != nil:
		newVersion.Disabled, newVersion.DisabledTooltip = true, l.Errors.DraftExists
	case !perms.Can("charge_policy", "update"):
		newVersion.Disabled, newVersion.DisabledTooltip = true, missing("charge_policy:update")
	}
	retire := types.TableAction{Type: "deactivate", Label: l.Buttons.Retire, Action: "deactivate",
		URL: route.ResolveURL(deps.Routes.RetireURL, "id", id), ItemName: s.Policy.GetName(),
		ConfirmTitle: l.Confirm.RetireTitle, ConfirmMessage: l.Confirm.RetireMsg,
		Disabled: !perms.Can("charge_policy", "retire"), DisabledTooltip: missing("charge_policy:retire")}
	return []types.TableAction{view, edit, newVersion, retire}
}
