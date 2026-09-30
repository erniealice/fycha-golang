// Package list renders the Document Series settings list (Active · Retired).
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
	pb "github.com/erniealice/esqyma/pkg/schema/v1/domain/revenue/document_series"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	ds "github.com/erniealice/fycha-golang/domain/ledger/document_series"
)

// Deps holds the list view dependencies.
type Deps struct {
	Routes       ds.Routes
	Labels       ds.Labels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	UseCases     *ds.UseCases
}

// PageData is the list page view model.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
	StatusTabs      []pyeza.TabItem
	ActiveStatus    string
}

func available(deps *Deps) bool {
	return deps.UseCases != nil && deps.UseCases.GetDocumentSeriesListPageData != nil
}

// NewView renders the full page (or the HTMX content partial).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("document_series", "list") {
			return view.Forbidden("document_series:list")
		}
		if !available(deps) {
			return view.Error(fmt.Errorf("%s", deps.Labels.Errors.Unavailable))
		}
		status := ds.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, tabs, err := build(ctx, deps, viewCtx, status, perms)
		if err != nil {
			return view.Error(err)
		}
		title := deps.Labels.Page.TitleActive
		if status == "retired" {
			title = deps.Labels.Page.TitleRetired
		}
		return view.OK("document-series-list", &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          title,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   deps.Routes.ActiveSubNav,
				HeaderTitle:    title,
				HeaderSubtitle: deps.Labels.Page.Subtitle,
				HeaderIcon:     "icon-hash",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "document-series-list-content",
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
		if !perms.Can("document_series", "list") {
			return view.Forbidden("document_series:list")
		}
		if !available(deps) {
			return view.Error(fmt.Errorf("%s", deps.Labels.Errors.Unavailable))
		}
		status := ds.NormalizeStatus(viewCtx.Request.PathValue("status"))
		table, _, err := build(ctx, deps, viewCtx, status, perms)
		if err != nil {
			return view.Error(err)
		}
		return view.OK("table-card", table)
	})
}

func columnsFor(l ds.Labels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "code", Label: l.Columns.Code, WidthClass: "col-3xl"},
		{Key: "name", Label: l.Columns.Name},
		{Key: "issuer", Label: l.Columns.Issuer, WidthClass: "col-3xl"},
		{Key: "document_kind", Label: l.Columns.DocumentKind, NoSort: true, NoFilter: true, WidthClass: "col-3xl"},
		{Key: "prefix", Label: l.Columns.Prefix, WidthClass: "col-xl"},
		{Key: "next_number", Label: l.Columns.NextNumber, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
		{Key: "fiscal_reset", Label: l.Columns.FiscalReset, NoSort: true, NoFilter: true, WidthClass: "col-2xl"},
		{Key: "status", Label: l.Columns.Status, NoSort: true, NoFilter: true, WidthClass: "col-xl"},
	}
}

func build(ctx context.Context, deps *Deps, viewCtx *view.ViewContext, status string, perms *types.UserPermissions) (*types.TableConfig, []pyeza.TabItem, error) {
	l := deps.Labels
	columns := columnsFor(l)
	p, err := espynahttp.ParseTableParamsWithFilters(viewCtx.Request, types.SortableKeys(columns), types.FilterableKeys(columns), "code", "asc")
	if err != nil {
		return nil, nil, err
	}

	all, err := deps.UseCases.ListSeries(ctx)
	if err != nil {
		log.Printf("document_series list: ListSeries: %v", err)
		return nil, nil, err
	}
	counts := map[string]int{}
	var shown []*pb.DocumentSeries
	for _, s := range all {
		if s == nil {
			continue
		}
		bucket := "active"
		if ds.IsRetired(s) {
			bucket = "retired"
		}
		counts[bucket]++
		if bucket == status {
			shown = append(shown, s)
		}
	}
	shown = applyQuery(shown, p)

	rows := buildRows(deps, shown, status, perms)
	types.ApplyColumnStyles(columns, rows)

	emptyTitle, emptyMsg := l.Empty.ActiveTitle, l.Empty.ActiveMessage
	if status == "retired" {
		emptyTitle, emptyMsg = l.Empty.RetiredTitle, ""
	}
	table := &types.TableConfig{
		ID:                   "document-series-table",
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
		EmptyState:           types.TableEmptyState{Title: emptyTitle, Message: emptyMsg},
	}
	if status == "active" {
		table.PrimaryAction = &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			TestID:          "document-series-add",
			Disabled:        !perms.Can("document_series", "create"),
			DisabledTooltip: fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "document_series:create"),
		}
	}
	types.ApplyTableSettings(table)

	tabs := make([]pyeza.TabItem, 0, len(ds.Statuses))
	for _, st := range ds.Statuses {
		label := l.Tabs.Active
		if st == "retired" {
			label = l.Tabs.Retired
		}
		tabs = append(tabs, pyeza.TabItem{
			Key: st, Label: label, Href: route.ResolveURL(deps.Routes.ListURL, "status", st), Count: counts[st],
		})
	}
	return table, tabs, nil
}

// applyQuery applies search and sort in memory (a workspace holds a handful of
// series; the walk above already read every page).
func applyQuery(in []*pb.DocumentSeries, p tableparams.TableQueryParams) []*pb.DocumentSeries {
	out := in
	if q := strings.ToLower(strings.TrimSpace(p.Search)); q != "" {
		out = out[:0:0]
		for _, s := range in {
			if strings.Contains(strings.ToLower(s.GetCode()), q) || strings.Contains(strings.ToLower(s.GetName()), q) ||
				strings.Contains(strings.ToLower(s.GetIssuerName()), q) || strings.Contains(strings.ToLower(s.GetPrefix()), q) {
				out = append(out, s)
			}
		}
	}
	key := func(s *pb.DocumentSeries) string {
		switch p.SortColumn {
		case "name":
			return strings.ToLower(s.GetName())
		case "issuer":
			return strings.ToLower(s.GetIssuerName())
		case "prefix":
			return strings.ToLower(s.GetPrefix())
		default:
			return strings.ToLower(s.GetCode())
		}
	}
	desc := p.SortDir == "desc"
	sort.SliceStable(out, func(i, j int) bool {
		a, b := key(out[i]), key(out[j])
		if a == b {
			return out[i].GetId() < out[j].GetId()
		}
		if desc {
			return a > b
		}
		return a < b
	})
	return out
}

func buildRows(deps *Deps, series []*pb.DocumentSeries, status string, perms *types.UserPermissions) []types.TableRow {
	l := deps.Labels
	rows := make([]types.TableRow, 0, len(series))
	for _, s := range series {
		id := s.GetId()
		statusText, variant := l.Enums.StatusActive, "success"
		if status == "retired" {
			statusText, variant = l.Enums.StatusRetired, "default"
		}
		cells := []types.TableCell{
			{Type: "text", Value: s.GetCode()},
			{Type: "text", Value: s.GetName()},
			{Type: "text", Value: s.GetIssuerName()},
			{Type: "text", Value: l.DocumentKindLabel(s.GetDocumentKind())},
			{Type: "text", Value: s.GetPrefix()},
			{Type: "text", Value: strconv.FormatInt(s.GetNextNumber(), 10)},
			{Type: "text", Value: l.FiscalResetLabel(s.GetFiscalReset())},
			{Type: "badge", Value: statusText, Variant: variant},
		}
		rows = append(rows, types.TableRow{
			ID:    id,
			Cells: cells,
			DataAttrs: map[string]string{
				"testid": "document-series-row-" + id,
				"name":   s.GetCode(),
				"status": status,
			},
			Actions: rowActions(deps, s, status, perms),
		})
	}
	return rows
}

func rowActions(deps *Deps, s *pb.DocumentSeries, status string, perms *types.UserPermissions) []types.TableAction {
	if status == "retired" {
		return nil
	}
	l := deps.Labels
	id := s.GetId()
	missing := fmt.Sprintf(deps.CommonLabels.Errors.MissingPermission, "document_series:update")
	can := perms.Can("document_series", "update")
	edit := types.TableAction{Type: "edit", Label: l.Buttons.Edit, Action: "edit",
		URL: route.ResolveURL(deps.Routes.EditURL, "id", id), DrawerTitle: l.Buttons.Edit,
		TestID: "document-series-edit-" + id, Disabled: !can, DisabledTooltip: missing}
	retire := types.TableAction{Type: "deactivate", Label: l.Buttons.Retire, Action: "deactivate",
		URL: route.ResolveURL(deps.Routes.RetireURL, "id", id), ItemName: s.GetCode(),
		ConfirmTitle: l.Confirm.RetireTitle, ConfirmMessage: l.Confirm.RetireMsg,
		TestID: "document-series-retire-" + id, Disabled: !can, DisabledTooltip: missing}
	return []types.TableAction{edit, retire}
}
