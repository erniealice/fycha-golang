package journal_action

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	jepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
	consumer "github.com/erniealice/espyna-golang/consumer"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Form data types
// ---------------------------------------------------------------------------

// FormData is the template data for the journal entry drawer form.
type FormData struct {
	FormAction  string
	IsEdit      bool
	ID          string
	Date        string
	Description string
	Notes       string
	Lines       []FormLine
	Labels      fycha.JournalFormLabels
	CommonLabels any
}

// FormLine represents one editable journal line in the form.
// The account selector submits account_id[N] (hidden) alongside debit[N], credit[N], memo[N].
// Consumer apps wire account lookup to populate AccountCode and AccountName for display.
type FormLine struct {
	Index       int    // 1-based line number for display
	AccountID   string // selected account ID (stored as hidden input)
	AccountCode string // display code (e.g. "1110")
	AccountName string // display name (e.g. "Cash on Hand")
	Debit       string
	Credit      string
	Memo        string
}

// ParsedLine holds one parsed journal line from the form submission.
// Consumer apps create JournalLine protos from these after creating the JournalEntry.
type ParsedLine struct {
	AccountID string
	Debit     float64
	Credit    float64
	Memo      string
	Order     int32
}

// ---------------------------------------------------------------------------
// Deps
// ---------------------------------------------------------------------------

// Deps holds dependencies for journal action handlers.
type Deps struct {
	Routes fycha.JournalRoutes
	Labels fycha.JournalLabels

	// Journal use cases
	CreateJournalEntry          func(ctx context.Context, req *jepb.CreateJournalEntryRequest) (*jepb.CreateJournalEntryResponse, error)
	ReadJournalEntry            func(ctx context.Context, req *jepb.ReadJournalEntryRequest) (*jepb.ReadJournalEntryResponse, error)
	UpdateJournalEntry          func(ctx context.Context, req *jepb.UpdateJournalEntryRequest) (*jepb.UpdateJournalEntryResponse, error)
	DeleteJournalEntry          func(ctx context.Context, req *jepb.DeleteJournalEntryRequest) (*jepb.DeleteJournalEntryResponse, error)
	PostJournalEntry            func(ctx context.Context, req *jepb.PostJournalEntryRequest) (*jepb.PostJournalEntryResponse, error)
	ReverseJournalEntry         func(ctx context.Context, req *jepb.ReverseJournalEntryRequest) (*jepb.ReverseJournalEntryResponse, error)
	GetJournalEntryItemPageData func(ctx context.Context, req *jepb.GetJournalEntryItemPageDataRequest) (*jepb.GetJournalEntryItemPageDataResponse, error)

	// Optional hook called after CreateJournalEntry to persist journal lines.
	// Receives the new journal entry ID and the parsed lines from the form.
	// Consumer apps wire this to CreateJournalLine for each ParsedLine.
	CreateLines func(ctx context.Context, journalEntryID string, lines []ParsedLine) error
}

// ---------------------------------------------------------------------------
// Add action (GET = empty form, POST = create draft or create+post)
// ---------------------------------------------------------------------------

// NewAddAction creates the journal entry add action.
// GET returns the empty sheet form.
// POST with submit=draft creates a draft; POST with submit=post creates then immediately posts.
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("journal", "post_manual") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("journal-drawer-form", newEmptyForm(deps))
		}

		// POST -- parse form and create
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
		}

		entry, parsedLines, err := parseJournalForm(viewCtx.Request)
		if err != nil {
			return fycha.HTMXError(err.Error())
		}

		submitAction := viewCtx.Request.FormValue("submit")
		entry.Status = jepb.JournalEntryStatus_JOURNAL_ENTRY_STATUS_DRAFT

		if deps.CreateJournalEntry == nil {
			log.Printf("CreateJournalEntry use case not wired")
			return fycha.HTMXSuccess("journals-table")
		}

		resp, err := deps.CreateJournalEntry(ctx, &jepb.CreateJournalEntryRequest{Data: entry})
		if err != nil {
			log.Printf("CreateJournalEntry error: %v", err)
			return fycha.HTMXError(deps.Labels.Actions.SaveError)
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := deps.Labels.Actions.SaveError
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		newID := ""
		if len(resp.GetData()) > 0 {
			newID = resp.GetData()[0].GetId()
		}

		// Persist lines via optional hook
		if newID != "" && deps.CreateLines != nil && len(parsedLines) > 0 {
			if lineErr := deps.CreateLines(ctx, newID, parsedLines); lineErr != nil {
				log.Printf("CreateLines error for %s: %v", newID, lineErr)
				// Non-fatal: entry was created, continue
			}
		}

		// If "post" was clicked, post the newly created entry
		if submitAction == "post" && newID != "" && deps.PostJournalEntry != nil {
			postedByCreate := consumer.ExtractUserIDFromContext(ctx)
			postResp, postErr := deps.PostJournalEntry(ctx, &jepb.PostJournalEntryRequest{
				JournalEntryId: newID,
				PostedBy:       postedByCreate,
			})
			if postErr != nil {
				log.Printf("PostJournalEntry error after create for %s: %v", newID, postErr)
			} else if postResp != nil && !postResp.GetSuccess() {
				log.Printf("PostJournalEntry response not success for %s", newID)
			}
		}

		return fycha.HTMXSuccess("journals-table")
	})
}

// ---------------------------------------------------------------------------
// Edit action (GET = form with existing data, POST = update)
// ---------------------------------------------------------------------------

// NewEditAction creates the journal entry edit action (draft entries only).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("journal", "update") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			formData := loadEditFormData(ctx, deps, id)
			return view.OK("journal-drawer-form", formData)
		}

		// POST -- update draft
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
		}

		entry, _, err := parseJournalForm(viewCtx.Request)
		if err != nil {
			return fycha.HTMXError(err.Error())
		}
		entry.Id = id

		if deps.UpdateJournalEntry == nil {
			log.Printf("UpdateJournalEntry use case not wired")
			return fycha.HTMXSuccess("journals-table")
		}

		resp, err := deps.UpdateJournalEntry(ctx, &jepb.UpdateJournalEntryRequest{Data: entry})
		if err != nil {
			log.Printf("UpdateJournalEntry error for %s: %v", id, err)
			return fycha.HTMXError(deps.Labels.Actions.SaveError)
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := deps.Labels.Actions.SaveError
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("journals-table")
	})
}

// ---------------------------------------------------------------------------
// Post action (POST only)
// ---------------------------------------------------------------------------

// NewPostAction creates the journal post action (transitions draft -> posted).
func NewPostAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("journal", "post_manual") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return fycha.HTMXError("Journal entry ID is required")
		}

		if deps.PostJournalEntry == nil {
			log.Printf("PostJournalEntry use case not wired")
			return fycha.HTMXSuccess("journals-table")
		}

		// Extract the authenticated user ID from context.
		postedBy := consumer.ExtractUserIDFromContext(ctx)

		resp, err := deps.PostJournalEntry(ctx, &jepb.PostJournalEntryRequest{
			JournalEntryId: id,
			PostedBy:       postedBy,
		})
		if err != nil {
			log.Printf("PostJournalEntry error for %s: %v", id, err)
			return fycha.HTMXError(deps.Labels.Actions.PostError)
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := deps.Labels.Actions.PostError
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("journals-table")
	})
}

// ---------------------------------------------------------------------------
// Reverse action (POST only)
// ---------------------------------------------------------------------------

// NewReverseAction creates the journal reverse action (posted -> reversed, new reversal entry created).
func NewReverseAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("journal", "post_manual") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return fycha.HTMXError("Journal entry ID is required")
		}

		if deps.ReverseJournalEntry == nil {
			log.Printf("ReverseJournalEntry use case not wired")
			return fycha.HTMXSuccess("journals-table")
		}

		resp, err := deps.ReverseJournalEntry(ctx, &jepb.ReverseJournalEntryRequest{
			JournalEntryId: id,
		})
		if err != nil {
			log.Printf("ReverseJournalEntry error for %s: %v", id, err)
			return fycha.HTMXError(deps.Labels.Actions.ReverseError)
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := deps.Labels.Actions.ReverseError
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("journals-table")
	})
}

// ---------------------------------------------------------------------------
// Delete action (POST only)
// ---------------------------------------------------------------------------

// NewDeleteAction creates the journal delete action (draft entries only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("journal", "delete") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			return fycha.HTMXError("Journal entry ID is required")
		}

		if deps.DeleteJournalEntry == nil {
			log.Printf("DeleteJournalEntry use case not wired")
			return fycha.HTMXSuccess("journals-table")
		}

		resp, err := deps.DeleteJournalEntry(ctx, &jepb.DeleteJournalEntryRequest{
			Data: &jepb.JournalEntry{Id: id},
		})
		if err != nil {
			log.Printf("DeleteJournalEntry error for %s: %v", id, err)
			return fycha.HTMXError(deps.Labels.Actions.DeleteError)
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := deps.Labels.Actions.DeleteError
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("journals-table")
	})
}

// ---------------------------------------------------------------------------
// Form helpers
// ---------------------------------------------------------------------------

func newEmptyForm(deps *Deps) *FormData {
	today := time.Now().Format("2006-01-02")
	return &FormData{
		FormAction:   deps.Routes.AddURL,
		IsEdit:       false,
		Date:         today,
		Labels:       deps.Labels.Form,
		CommonLabels: nil,
		Lines: []FormLine{
			{Index: 1},
			{Index: 2},
		},
	}
}

func loadEditFormData(ctx context.Context, deps *Deps, id string) *FormData {
	base := &FormData{
		FormAction:   deps.Routes.EditURL,
		IsEdit:       true,
		ID:           id,
		Labels:       deps.Labels.Form,
		CommonLabels: nil,
		Lines:        []FormLine{{Index: 1}, {Index: 2}},
	}

	if deps.GetJournalEntryItemPageData == nil {
		return base
	}

	resp, err := deps.GetJournalEntryItemPageData(ctx, &jepb.GetJournalEntryItemPageDataRequest{
		JournalEntryId: id,
	})
	if err != nil {
		log.Printf("GetJournalEntryItemPageData error for edit form %s: %v", id, err)
		return base
	}
	if resp == nil || !resp.GetSuccess() || resp.GetJournalEntry() == nil {
		return base
	}

	e := resp.GetJournalEntry()
	dateStr := e.GetEntryDateString()
	notes := e.GetNotes()

	return &FormData{
		FormAction:   deps.Routes.EditURL,
		IsEdit:       true,
		ID:           id,
		Date:         dateStr,
		Description:  e.GetDescription(),
		Notes:        notes,
		Lines:        []FormLine{{Index: 1}, {Index: 2}}, // lines loaded separately by consumer
		Labels:       deps.Labels.Form,
		CommonLabels: nil,
	}
}

// ParseJournalFormLines parses the line items from the form.
// Exported so consumer apps can call it when they need to persist lines
// via their own JournalLine service.
//
// Line fields: account_id[N], debit[N], credit[N], memo[N] (N is 1-based).
func ParseJournalFormLines(r *http.Request) []ParsedLine {
	var lines []ParsedLine
	order := int32(1)
	for i := 1; i <= 50; i++ {
		key := fmt.Sprintf("%d", i)
		accountID := r.FormValue("account_id[" + key + "]")
		if accountID == "" {
			continue
		}

		debit := parseAmount(r.FormValue("debit[" + key + "]"))
		credit := parseAmount(r.FormValue("credit[" + key + "]"))
		memo := r.FormValue("memo[" + key + "]")

		if debit == 0 && credit == 0 {
			continue
		}

		lines = append(lines, ParsedLine{
			AccountID: accountID,
			Debit:     debit,
			Credit:    credit,
			Memo:      memo,
			Order:     order,
		})
		order++
	}
	return lines
}

// parseJournalForm parses the multipart form into a JournalEntry proto + parsed lines.
func parseJournalForm(r *http.Request) (*jepb.JournalEntry, []ParsedLine, error) {
	desc := r.FormValue("description")
	if strings.TrimSpace(desc) == "" {
		return nil, nil, fmt.Errorf("description is required")
	}

	dateStr := r.FormValue("date")
	notes := r.FormValue("notes")

	lines := ParseJournalFormLines(r)
	if len(lines) < 2 {
		return nil, nil, fmt.Errorf("at least 2 journal lines are required")
	}

	// Validate balance
	var totalDebit, totalCredit float64
	for _, l := range lines {
		totalDebit += l.Debit
		totalCredit += l.Credit
	}
	diff := totalDebit - totalCredit
	if diff < 0 {
		diff = -diff
	}
	if diff > 0.005 {
		return nil, nil, fmt.Errorf(
			"journal entry is unbalanced: total debits %.2f != total credits %.2f",
			totalDebit, totalCredit,
		)
	}

	entry := &jepb.JournalEntry{
		Description:     desc,
		EntryDateString: &dateStr,
		TotalDebit:      totalDebit,
		TotalCredit:     totalCredit,
	}
	if notes != "" {
		entry.Notes = &notes
	}

	return entry, lines, nil
}

func parseAmount(s string) float64 {
	if s == "" {
		return 0
	}
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimPrefix(s, "\u20b1")
	s = strings.TrimPrefix(s, "$")
	f, err := strconv.ParseFloat(strings.TrimSpace(s), 64)
	if err != nil {
		return 0
	}
	return f
}
