// routes.go defines configurable route structs for fycha ledger journal views.
package journal

const (
	// Ledger — Journal Entries
	ListURL             = "/ledger/journals/list/{status}"
	DetailURL           = "/ledger/journals/detail/{id}"
	TabActionURL        = "/action/ledger/journal/{id}/tab/{tab}"
	AttachmentUploadURL = "/action/ledger/journal/{id}/attachments/upload"
	AttachmentDeleteURL = "/action/ledger/journal/{id}/attachments/delete"
	AddURL              = "/action/ledger/journal/add"
	EditURL             = "/action/ledger/journal/edit/{id}"
	PostURL             = "/action/ledger/journal/post/{id}"
	ReverseURL          = "/action/ledger/journal/reverse/{id}"
	DeleteURL           = "/action/ledger/journal/delete"
)

// ---------------------------------------------------------------------------
// Routes
// ---------------------------------------------------------------------------

// Routes holds route paths for Journal Entry views.
type Routes struct {
	ActiveNav    string `json:"active_nav"`
	ActiveSubNav string `json:"active_sub_nav"`
	ListURL      string `json:"list_url"`
	DetailURL    string `json:"detail_url"`
	TabActionURL string `json:"tab_action_url"`
	AddURL       string `json:"add_url"`
	EditURL      string `json:"edit_url"`
	PostURL      string `json:"post_url"`
	ReverseURL   string `json:"reverse_url"`
	DeleteURL    string `json:"delete_url"`

	// Attachment routes
	AttachmentUploadURL string `json:"attachment_upload_url"`
	AttachmentDeleteURL string `json:"attachment_delete_url"`
}

func DefaultRoutes() Routes {
	return Routes{
		ActiveNav:    "ledger",
		ActiveSubNav: "journals-draft",
		ListURL:      ListURL,
		DetailURL:    DetailURL,
		TabActionURL: TabActionURL,
		AddURL:       AddURL,
		EditURL:      EditURL,
		PostURL:      PostURL,
		ReverseURL:   ReverseURL,
		DeleteURL:    DeleteURL,

		AttachmentUploadURL: AttachmentUploadURL,
		AttachmentDeleteURL: AttachmentDeleteURL,
	}
}

func (r Routes) RouteMap() map[string]string {
	return map[string]string{
		"ledger.journal.list":    r.ListURL,
		"ledger.journal.detail":  r.DetailURL,
		"ledger.journal.add":     r.AddURL,
		"ledger.journal.edit":    r.EditURL,
		"ledger.journal.post":    r.PostURL,
		"ledger.journal.reverse": r.ReverseURL,
		"ledger.journal.delete":  r.DeleteURL,

		"ledger.journal.attachment.upload": r.AttachmentUploadURL,
		"ledger.journal.attachment.delete": r.AttachmentDeleteURL,
	}
}
