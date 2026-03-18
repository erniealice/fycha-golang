package journal

import (
	"context"
	"fmt"
	"log"

	jepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies.
type Deps struct {
	Routes       fycha.JournalRoutes
	Labels       fycha.JournalLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Journal use cases
	GetJournalEntryListPageData func(ctx context.Context, req *jepb.GetJournalEntryListPageDataRequest) (*jepb.GetJournalEntryListPageDataResponse, error)
}

// PageData holds the data for the journal list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveStatus    string          // "draft", "posted", "reversed"
	Table           *types.TableConfig
}

// JournalRow is the view-model for a single journal entry row.
type JournalRow struct {
	ID          string
	EntryNumber string
	Description string
	EntryDate   string
	Status      string // "draft", "posted", "reversed"
	SourceType  string // human-readable source type label
	SourceID    string // FK to originating entity (may be empty)
	TotalDebit  string
	TotalCredit string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the journal list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "draft"
		}

		entries := fetchEntries(ctx, deps)
		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, status, entries, perms)

		heading, subtitle := headingForStatus(deps.Labels, status)
		activeSubNav := subNavForStatus(status)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   activeSubNav,
				HeaderTitle:    heading,
				HeaderSubtitle: subtitle,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "journal-list-content",
			ActiveStatus:    status,
			Table:           tableConfig,
		}

		return view.OK("journal-list", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchEntries(ctx context.Context, deps *Deps) []JournalRow {
	if deps.GetJournalEntryListPageData == nil {
		return []JournalRow{}
	}

	resp, err := deps.GetJournalEntryListPageData(ctx, &jepb.GetJournalEntryListPageDataRequest{})
	if err != nil {
		log.Printf("GetJournalEntryListPageData error: %v", err)
		return []JournalRow{}
	}
	if resp == nil || !resp.GetSuccess() {
		return []JournalRow{}
	}

	rows := make([]JournalRow, 0, len(resp.GetJournalEntryList()))
	for _, e := range resp.GetJournalEntryList() {
		rows = append(rows, protoToRow(e))
	}
	return rows
}

func protoToRow(e *jepb.JournalEntry) JournalRow {
	return JournalRow{
		ID:          e.GetId(),
		EntryNumber: e.GetEntryNumber(),
		Description: e.GetDescription(),
		EntryDate:   e.GetEntryDateString(),
		Status:      statusString(e.GetStatus()),
		SourceType:  sourceTypeLabel(e.GetSourceType()),
		SourceID:    e.GetSourceId(),
		TotalDebit:  fmt.Sprintf("%.2f", e.GetTotalDebit()),
		TotalCredit: fmt.Sprintf("%.2f", e.GetTotalCredit()),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, status string, entries []JournalRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := journalColumns(l)
	rows := buildTableRows(entries, status, l, deps.Routes, perms)
	types.ApplyColumnStyles(columns, rows)

	canCreate := perms.Can("journal", "post_manual")

	tableConfig := &types.TableConfig{
		ID:                "journals-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowExport:        true,
		ShowEntries:       true,
		DefaultSortColumn: "entry_number",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.NewEntry,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !canCreate,
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func journalColumns(l fycha.JournalLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "entry_number", Label: l.Columns.EntryNumber, Sortable: true, Width: "110px"},
		{Key: "date", Label: l.Columns.Date, Sortable: true, Width: "100px"},
		{Key: "description", Label: l.Columns.Description, Sortable: false},
		{Key: "amount", Label: l.Columns.Amount, Sortable: false, Width: "140px", Align: "right"},
		{Key: "source", Label: l.Columns.Source, Sortable: false, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: false, Width: "100px"},
	}
}

func buildTableRows(entries []JournalRow, status string, l fycha.JournalLabels, routes fycha.JournalRoutes, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, je := range entries {
		if je.Status != status {
			continue
		}

		detailHref := route.ResolveURL(routes.DetailURL, "id", je.ID)
		canUpdate := perms.Can("journal", "update")
		canDelete := perms.Can("journal", "delete")
		canPost := perms.Can("journal", "post_manual")

		var actions []types.TableAction
		switch je.Status {
		case "draft":
			actions = append(actions,
				types.TableAction{Type: "view", Label: l.Actions.View, Action: "view", Href: detailHref},
				types.TableAction{
					Type:            "edit",
					Label:           l.Actions.Edit,
					Action:          "edit",
					URL:             route.ResolveURL(routes.EditURL, "id", je.ID),
					DrawerTitle:     l.Actions.Edit,
					Disabled:        !canUpdate,
					DisabledTooltip: l.Actions.NoPermission,
				},
				types.TableAction{
					Type:            "action",
					Label:           l.Actions.Post,
					Action:          "post",
					URL:             route.ResolveURL(routes.PostURL, "id", je.ID),
					Disabled:        !canPost,
					DisabledTooltip: l.Actions.NoPermission,
				},
				types.TableAction{
					Type:            "delete",
					Label:           l.Actions.Delete,
					Action:          "delete",
					URL:             routes.DeleteURL,
					ItemName:        je.EntryNumber,
					Disabled:        !canDelete,
					DisabledTooltip: l.Actions.NoPermission,
				},
			)
		case "posted":
			actions = append(actions,
				types.TableAction{Type: "view", Label: l.Actions.View, Action: "view", Href: detailHref},
				types.TableAction{
					Type:            "action",
					Label:           l.Actions.Reverse,
					Action:          "reverse",
					URL:             route.ResolveURL(routes.ReverseURL, "id", je.ID),
					Disabled:        !canPost,
					DisabledTooltip: l.Actions.NoPermission,
				},
			)
		case "reversed":
			actions = append(actions,
				types.TableAction{Type: "view", Label: l.Actions.View, Action: "view", Href: detailHref},
			)
		}

		statusVariant := statusBadgeVariant(je.Status)
		statusLabel := statusLabel(l, je.Status)

		row := types.TableRow{
			ID:   je.ID,
			Href: detailHref,
			Cells: []types.TableCell{
				{Type: "link", Value: je.EntryNumber, Href: detailHref},
				{Type: "text", Value: je.EntryDate},
				{Type: "text", Value: je.Description},
				{Type: "text", Value: je.TotalDebit},
				{Type: "text", Value: je.SourceType},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			Actions: actions,
		}
		rows = append(rows, row)
	}
	return rows
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func headingForStatus(l fycha.JournalLabels, status string) (heading, subtitle string) {
	switch status {
	case "draft":
		return l.Page.HeadingDraft, l.Page.SubtitleDraft
	case "posted":
		return l.Page.HeadingPosted, l.Page.SubtitlePosted
	case "reversed":
		return l.Page.HeadingReversed, l.Page.SubtitleReversed
	default:
		return l.Page.HeadingDraft, l.Page.SubtitleDraft
	}
}

func subNavForStatus(status string) string {
	switch status {
	case "posted":
		return "journals-posted"
	case "reversed":
		return "journals-reversed"
	default:
		return "journals-draft"
	}
}

func statusString(s jepb.JournalEntryStatus) string {
	switch s {
	case jepb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_DRAFT:
		return "draft"
	case jepb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_POSTED:
		return "posted"
	case jepb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_REVERSED:
		return "reversed"
	default:
		return "draft"
	}
}

func statusBadgeVariant(status string) string {
	switch status {
	case "draft":
		return "warning"
	case "posted":
		return "success"
	case "reversed":
		return "muted"
	default:
		return "default"
	}
}

func statusLabel(l fycha.JournalLabels, status string) string {
	switch status {
	case "draft":
		return l.Tabs.Draft
	case "posted":
		return l.Tabs.Posted
	case "reversed":
		return l.Tabs.Reversed
	default:
		return status
	}
}

func sourceTypeLabel(t jepb.JournalSourceType) string {
	switch t {
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_MANUAL:
		return "Manual"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_REVENUE:
		return "Revenue"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_EXPENDITURE:
		return "Expenditure"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_COLLECTION:
		return "Collection"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_DISBURSEMENT:
		return "Disbursement"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_DEPRECIATION:
		return "Depreciation"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_ASSET_ACQUISITION:
		return "Asset Acquisition"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_ASSET_DISPOSAL:
		return "Asset Disposal"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_PREPAYMENT:
		return "Prepayment"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_PREPAYMENT_AMORTIZATION:
		return "Prepayment Amortization"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_LOAN_RECEIPT:
		return "Loan Receipt"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_LOAN_PAYMENT:
		return "Loan Payment"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_PETTY_CASH_REPLENISHMENT:
		return "Petty Cash Replenishment"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_BAD_DEBT_PROVISION:
		return "Bad Debt Provision"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_DEFERRED_REVENUE:
		return "Deferred Revenue"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_EQUITY_CONTRIBUTION:
		return "Equity Contribution"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_EQUITY_WITHDRAWAL:
		return "Equity Withdrawal"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_EQUITY_DISTRIBUTION:
		return "Equity Distribution"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_YEAR_END_CLOSE:
		return "Year-End Close"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_RECURRING:
		return "Recurring"
	case jepb.JournalSourceType_JOURNAL_SOURCE_TYPE_PAYROLL:
		return "Payroll"
	default:
		return "Manual"
	}
}
