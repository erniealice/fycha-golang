package journal_detail

import (
	"context"
	"fmt"
	"log"

	jepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
	jlpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_line"
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
	GetJournalEntryItemPageData func(ctx context.Context, req *jepb.GetJournalEntryItemPageDataRequest) (*jepb.GetJournalEntryItemPageDataResponse, error)
}

// PageData holds the data for the journal detail page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Labels          fycha.JournalLabels
	// Identity
	ID          string
	EntryNumber string
	Description string
	EntryDate   string
	Status      string
	StatusVariant string
	SourceType  string
	SourceID    string
	SourceHref  string // link to originating transaction (if auto-generated)
	Notes       string
	// Posting info
	PostedBy   string
	PostedAt   string
	// Reversal info
	ReversedBy      string
	ReversedAt      string
	ReversalEntryID string
	ReversalHref    string
	// Audit
	Created      string
	LastModified string
	// Totals
	TotalDebit  string
	TotalCredit string
	Difference  string
	IsBalanced  bool
	// Lines table
	LinesTable *types.TableConfig
	// Action URLs + permissions
	EditURL        string
	PostURL        string
	ReverseURL     string
	DeleteURL      string
	CanEdit        bool
	CanPost        bool
	CanReverse     bool
	CanDelete      bool
}

// LineRow is the view-model for a single journal line.
type LineRow struct {
	AccountID   string
	AccountCode string
	AccountName string
	Description string
	Debit       string
	Credit      string
	LineOrder   int
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the journal detail view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		id := viewCtx.Request.PathValue("id")
		perms := view.GetUserPermissions(ctx)
		pageData := buildPageData(ctx, deps, id, viewCtx, perms)
		return view.OK("journal-detail", pageData)
	})
}

// ---------------------------------------------------------------------------
// Page data builder
// ---------------------------------------------------------------------------

func buildPageData(ctx context.Context, deps *Deps, id string, viewCtx *view.ViewContext, perms *types.UserPermissions) *PageData {
	vm := fetchEntry(ctx, deps, id)

	linesTable := buildLinesTable(vm.Lines, deps.Labels, deps.TableLabels)

	postURL := route.ResolveURL(deps.Routes.PostURL, "id", id)
	reverseURL := route.ResolveURL(deps.Routes.ReverseURL, "id", id)
	editURL := route.ResolveURL(deps.Routes.EditURL, "id", id)

	return &PageData{
		PageData: types.PageData{
			CacheVersion:   viewCtx.CacheVersion,
			Title:          fmt.Sprintf("Journal Entry %s", vm.EntryNumber),
			CurrentPath:    viewCtx.CurrentPath,
			ActiveNav:      deps.Routes.ActiveNav,
			ActiveSubNav:   subNavForStatus(vm.Status),
			HeaderTitle:    fmt.Sprintf("Journal Entry %s", vm.EntryNumber),
			HeaderSubtitle: vm.Description,
			HeaderIcon:     "icon-file-text",
			CommonLabels:   deps.CommonLabels,
		},
		ContentTemplate: "journal-detail-content",
		Labels:          deps.Labels,
		ID:              id,
		EntryNumber:     vm.EntryNumber,
		Description:     vm.Description,
		EntryDate:       vm.EntryDate,
		Status:          vm.Status,
		StatusVariant:   statusBadgeVariant(vm.Status),
		SourceType:      vm.SourceType,
		SourceID:        vm.SourceID,
		SourceHref:      vm.SourceHref,
		Notes:           vm.Notes,
		PostedBy:        vm.PostedBy,
		PostedAt:        vm.PostedAt,
		ReversedBy:      vm.ReversedBy,
		ReversedAt:      vm.ReversedAt,
		ReversalEntryID: vm.ReversalEntryID,
		ReversalHref:    route.ResolveURL(deps.Routes.DetailURL, "id", vm.ReversalEntryID),
		Created:         vm.Created,
		LastModified:    vm.LastModified,
		TotalDebit:      vm.TotalDebit,
		TotalCredit:     vm.TotalCredit,
		Difference:      vm.Difference,
		IsBalanced:      vm.IsBalanced,
		LinesTable:      linesTable,
		EditURL:         editURL,
		PostURL:         postURL,
		ReverseURL:      reverseURL,
		DeleteURL:       deps.Routes.DeleteURL,
		CanEdit:         perms.Can("journal", "update") && vm.Status == "draft",
		CanPost:         perms.Can("journal", "post_manual") && vm.Status == "draft",
		CanReverse:      perms.Can("journal", "post_manual") && vm.Status == "posted",
		CanDelete:       perms.Can("journal", "delete") && vm.Status == "draft",
	}
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

// journalViewModel is the normalised view-model for a single journal entry.
type journalViewModel struct {
	EntryNumber     string
	Description     string
	EntryDate       string
	Status          string
	SourceType      string
	SourceID        string
	SourceHref      string
	Notes           string
	PostedBy        string
	PostedAt        string
	ReversedBy      string
	ReversedAt      string
	ReversalEntryID string
	Created         string
	LastModified    string
	TotalDebit      string
	TotalCredit     string
	Difference      string
	IsBalanced      bool
	Lines           []LineRow
}

func placeholder() journalViewModel {
	return journalViewModel{
		EntryNumber: "\u2014",
		Description: "Journal entry not found",
		Status:      "draft",
		TotalDebit:  "\u20b10.00",
		TotalCredit: "\u20b10.00",
		Difference:  "\u20b10.00",
		IsBalanced:  false,
	}
}

func fetchEntry(ctx context.Context, deps *Deps, id string) journalViewModel {
	if deps.GetJournalEntryItemPageData == nil {
		return placeholder()
	}

	resp, err := deps.GetJournalEntryItemPageData(ctx, &jepb.GetJournalEntryItemPageDataRequest{
		JournalEntryId: id,
	})
	if err != nil {
		log.Printf("GetJournalEntryItemPageData error for %s: %v", id, err)
		return placeholder()
	}
	if resp == nil || !resp.GetSuccess() || resp.GetJournalEntry() == nil {
		return placeholder()
	}

	return protoToViewModel(resp.GetJournalEntry())
}

func protoToViewModel(e *jepb.JournalEntry) journalViewModel {
	status := statusString(e.GetStatus())
	totalDebit := e.GetTotalDebit()
	totalCredit := e.GetTotalCredit()
	diff := totalDebit - totalCredit
	if diff < 0 {
		diff = -diff
	}

	vm := journalViewModel{
		EntryNumber:     e.GetEntryNumber(),
		Description:     e.GetDescription(),
		EntryDate:       e.GetEntryDateString(),
		Status:          status,
		SourceType:      sourceTypeLabel(e.GetSourceType()),
		SourceID:        e.GetSourceId(),
		Notes:           e.GetNotes(),
		PostedBy:        e.GetPostedBy(),
		PostedAt:        e.GetPostedAtString(),
		ReversedBy:      e.GetReversedBy(),
		ReversedAt:      e.GetReversedAtString(),
		ReversalEntryID: e.GetReversalEntryId(),
		Created:         e.GetDateCreatedString(),
		LastModified:    e.GetDateModifiedString(),
		TotalDebit:      fmt.Sprintf("\u20b1%.2f", totalDebit),
		TotalCredit:     fmt.Sprintf("\u20b1%.2f", totalCredit),
		Difference:      fmt.Sprintf("\u20b1%.2f", diff),
		IsBalanced:      diff < 0.005, // allow for float precision
	}

	return vm
}

// ---------------------------------------------------------------------------
// Lines table builder
// ---------------------------------------------------------------------------

func buildLinesTable(lines []LineRow, labels fycha.JournalLabels, tableLabels types.TableLabels) *types.TableConfig {
	columns := []types.TableColumn{
		{Key: "account_code", Label: labels.Lines.AccountCode, Width: "120px"},
		{Key: "account_name", Label: labels.Lines.AccountName},
		{Key: "memo", Label: labels.Lines.Memo},
		{Key: "debit", Label: labels.Lines.Debit, Width: "130px", Align: "right"},
		{Key: "credit", Label: labels.Lines.Credit, Width: "130px", Align: "right"},
	}

	rows := make([]types.TableRow, len(lines))
	for i, l := range lines {
		rows[i] = types.TableRow{
			ID: fmt.Sprintf("line-%d", i+1),
			Cells: []types.TableCell{
				{Type: "text", Value: l.AccountCode},
				{Type: "text", Value: l.AccountName},
				{Type: "text", Value: l.Description},
				{Type: "text", Value: l.Debit},
				{Type: "text", Value: l.Credit},
			},
		}
	}

	types.ApplyColumnStyles(columns, rows)

	cfg := &types.TableConfig{
		ID:         "journal-lines-table",
		Minimal:    true,
		Columns:    columns,
		Rows:       rows,
		ShowSearch: false,
		Labels:     tableLabels,
		EmptyState: types.TableEmptyState{
			Title:   labels.Lines.EmptyTitle,
			Message: labels.Lines.EmptyMessage,
		},
	}
	types.ApplyTableSettings(cfg)
	return cfg
}

// ---------------------------------------------------------------------------
// Proto enum converters
// ---------------------------------------------------------------------------

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

// LinesToViewModels converts journal line protos to LineRow view-models.
// Called by the detail view when lines are embedded in the JournalEntry response.
// accountCodeByID maps account_id -> code+name for display; pass nil to show ID as fallback.
func LinesToViewModels(lines []*jlpb.JournalLine, accountCodeByID map[string]string, accountNameByID map[string]string) []LineRow {
	rows := make([]LineRow, len(lines))
	for i, l := range lines {
		code := accountCodeByID[l.GetAccountId()]
		name := accountNameByID[l.GetAccountId()]
		if code == "" {
			code = l.GetAccountId()
		}

		debit := ""
		credit := ""
		if l.GetDebitAmount() > 0 {
			debit = fmt.Sprintf("%.2f", l.GetDebitAmount())
		}
		if l.GetCreditAmount() > 0 {
			credit = fmt.Sprintf("%.2f", l.GetCreditAmount())
		}

		rows[i] = LineRow{
			AccountID:   l.GetAccountId(),
			AccountCode: code,
			AccountName: name,
			Description: l.GetDescription(),
			Debit:       debit,
			Credit:      credit,
			LineOrder:   int(l.GetLineOrder()),
		}
	}
	return rows
}
