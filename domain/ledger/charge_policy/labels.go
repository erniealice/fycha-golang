package charge_policy

import "strings"

// Labels holds every translatable string of the charge policy pages.
// Lyngua file: general/charge_policy.json, root key "charge_policy" (loaded with
// an empty binding key elsewhere in fycha; here the descriptor binds
// Key "charge_policy"). JSON tags byte-match the lyngua key tree.
//
// Every field has a general-tier Lyngua key (TestDefaultLabelsEqualGeneralLyngua
// compares every string field); the Go defaults are the English fallback.
type Labels struct {
	Page    PageLabels    `json:"page"`
	Tabs    TabLabels     `json:"tabs"`
	Columns ColumnLabels  `json:"columns"`
	Buttons ButtonLabels  `json:"buttons"`
	Form    FormLabels    `json:"form"`
	Enums   EnumLabels    `json:"enums"`
	Detail  DetailLabels  `json:"detail"`
	Confirm ConfirmLabels `json:"confirm"`
	Empty   EmptyLabels   `json:"empty"`
	Errors  ErrorLabels   `json:"errors"`
}

type PageLabels struct {
	Title        string `json:"title"`
	TitleActive  string `json:"title_active"`
	TitleDraft   string `json:"title_draft"`
	TitleRetired string `json:"title_retired"`
	Subtitle     string `json:"subtitle"`
}

type TabLabels struct {
	Active         string `json:"active"`
	Draft          string `json:"draft"`
	Retired        string `json:"retired"`
	Info           string `json:"info"`
	Versions       string `json:"versions"`
	Usage          string `json:"usage"`
	History        string `json:"history"`
	Classification string `json:"classification"`
	Components     string `json:"components"`
	Postings       string `json:"postings"`
	Approval       string `json:"approval"`
}

type ColumnLabels struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	AccountingRole   string `json:"accounting_role"`
	DocumentKind     string `json:"document_kind"`
	CurrentVersion   string `json:"current_version"`
	InUse            string `json:"in_use"`
	Modified         string `json:"modified"`
	DraftVersion     string `json:"draft_version"`
	PreparedBy       string `json:"prepared_by"`
	RetiredOn        string `json:"retired_on"`
	RetiredBy        string `json:"retired_by"`
	Version          string `json:"version"`
	Status           string `json:"status"`
	ApprovedBy       string `json:"approved_by"`
	ApprovedOn       string `json:"approved_on"`
	ComponentRole    string `json:"component_role"`
	BookPresentation string `json:"book_presentation"`
	Event            string `json:"event"`
	PostingRole      string `json:"posting_role"`
	Account          string `json:"account"`
	LastVersion      string `json:"last_version"`
	Kind             string `json:"kind"`
	UsageDetail      string `json:"usage_detail"`
	Description      string `json:"description"`
	Rule             string `json:"rule"`
	Passed           string `json:"passed"`
	SelfApproved     string `json:"self_approved"`
}

type ButtonLabels struct {
	Add          string `json:"add"`
	Edit         string `json:"edit"`
	Retire       string `json:"retire"`
	Delete       string `json:"delete"`
	NewVersion   string `json:"new_version"`
	Approve      string `json:"approve"`
	AddComponent string `json:"add_component"`
	AddPosting   string `json:"add_posting"`
	View         string `json:"view"`
	EditVersion  string `json:"edit_version"`
	DeleteDraft  string `json:"delete_draft"`
	BackToPolicy string `json:"back_to_policy"`
}

type FormLabels struct {
	CodeLabel             string `json:"code_label"`
	CodePlaceholder       string `json:"code_placeholder"`
	NameLabel             string `json:"name_label"`
	DescriptionLabel      string `json:"description_label"`
	TemplateLabel         string `json:"template_label"`
	TemplatePlaceholder   string `json:"template_placeholder"`
	AccountingRoleLabel   string `json:"accounting_role_label"`
	TaxPositionLabel      string `json:"tax_position_label"`
	TaxTreatmentLabel     string `json:"tax_treatment_label"`
	AssessmentScopeLabel  string `json:"assessment_scope_label"`
	AssessmentNoteLabel   string `json:"assessment_note_label"`
	BookPresentationLabel string `json:"book_presentation_label"`
	ComponentRoleLabel    string `json:"component_role_label"`
	DocumentKindLabel     string `json:"document_kind_label"`
	EventLabel            string `json:"event_label"`
	PostingRoleLabel      string `json:"posting_role_label"`
	AccountLabel          string `json:"account_label"`
	AccountPlaceholder    string `json:"account_placeholder"`
	Unspecified           string `json:"unspecified"`
}

type EnumLabels struct {
	AccountingRolePrincipal string `json:"accounting_role_principal"`
	AccountingRoleAgent     string `json:"accounting_role_agent"`
	DocumentKindInvoice     string `json:"document_kind_invoice"`
	DocumentKindRecovery    string `json:"document_kind_recovery_document"`
	ComponentRoleOwnRevenue string `json:"component_role_own_revenue"`
	ComponentRoleRecovery   string `json:"component_role_recovery_cost"`
	ComponentRoleFee        string `json:"component_role_fee"`
	VersionStatusDraft      string `json:"version_status_draft"`
	VersionStatusApproved   string `json:"version_status_approved"`
	VersionStatusSuperseded string `json:"version_status_superseded"`

	BookPresentationRevenue  string `json:"book_presentation_revenue"`
	BookPresentationExcluded string `json:"book_presentation_excluded_reimbursement"`
	TaxPositionOwnSupply     string `json:"tax_position_own_supply"`
	TaxPositionExcluded      string `json:"tax_position_excluded_reimbursement"`
	EventAllocation          string `json:"event_allocation"`
	EventIssue               string `json:"event_issue"`
	EventApplication         string `json:"event_application"`
	EventReversal            string `json:"event_reversal"`
	PostingRoleReceivable    string `json:"posting_role_receivable"`
	PostingRoleClearing      string `json:"posting_role_clearing"`
	PostingRoleRevenue       string `json:"posting_role_revenue"`
	PostingRoleRefund        string `json:"posting_role_refund_liability"`
	PostingRoleExpense       string `json:"posting_role_expense"`
	PostingRoleCash          string `json:"posting_role_cash"`
	PolicyStatusActive       string `json:"policy_status_active"`
	PolicyStatusRetired      string `json:"policy_status_retired"`
	PolicyStatusDraft        string `json:"policy_status_draft"`
	Yes                      string `json:"yes"`
	No                       string `json:"no"`
}

type DetailLabels struct {
	CurrentVersion string `json:"current_version"` // "Current: v{0}"
	LockedNotice   string `json:"locked_notice"`
	ChecklistTitle string `json:"checklist_title"`
	MissingMapping string `json:"missing_mapping"`
	BackToPolicy   string `json:"back_to_policy"` // "Back to {0}"

	NoCurrentVersion   string            `json:"no_current_version"`
	ClonedFrom         string            `json:"cloned_from"`
	SelfApprovedNote   string            `json:"self_approved_note"`
	ChecklistPassed    string            `json:"checklist_passed"`
	ChecklistFailed    string            `json:"checklist_failed"`
	UnsupportedTitle   string            `json:"unsupported_title"`
	NoDescription      string            `json:"no_description"`
	VersionTitle       string            `json:"version_title"`
	ApprovedBy         string            `json:"approved_by"`
	ApprovedAt         string            `json:"approved_at"`
	Created            string            `json:"created"`
	UsageEmpty         string            `json:"usage_empty"`
	UsageInUse         string            `json:"usage_in_use"`
	UsageNotInUse      string            `json:"usage_not_in_use"`
	OpenDraftExists    string            `json:"open_draft_exists"`
	NewVersionMessage  string            `json:"new_version_message"`
	EvidenceTitle      string            `json:"evidence_title"`
	ChecklistCodeLabel map[string]string `json:"checklist_codes"`
	ReasonLabel        map[string]string `json:"reasons"`
}

type ConfirmLabels struct {
	RetireTitle  string `json:"retire_title"`
	RetireMsg    string `json:"retire_msg"`
	ApproveTitle string `json:"approve_title"`
	ApproveMsg   string `json:"approve_msg"`

	DeleteTitle string `json:"delete_title"`
	DeleteMsg   string `json:"delete_msg"`
}

type EmptyLabels struct {
	ActiveTitle  string `json:"active_title"`
	DraftTitle   string `json:"draft_title"`
	RetiredTitle string `json:"retired_title"`

	Message string `json:"message"`
}

// ErrorLabels: the first 13 fields are the charge_policy use-case error codes
// (`ErrorCode()`), keyed `charge_policy.errors.<code>` and shared with espyna's
// Translator (values equal espyna's defaults); the rest are view-only messages.
type ErrorLabels struct {
	Validation             string `json:"validation"`
	NotFound               string `json:"not_found"`
	NotDraft               string `json:"not_draft"`
	DraftExists            string `json:"draft_exists"`
	NoApprovedVersion      string `json:"no_approved_version"`
	Retired                string `json:"retired"`
	SelfApproval           string `json:"self_approval"`
	ChecklistFailed        string `json:"checklist_failed"`
	UnsupportedCombination string `json:"unsupported_combination"`
	InUse                  string `json:"in_use"`
	TransactionRequired    string `json:"transaction_required"`
	TaxTreatmentNotFound   string `json:"tax_treatment_not_found"`
	Unverifiable           string `json:"unverifiable"`

	FormInvalid        string `json:"form_invalid"`
	Unavailable        string `json:"unavailable"`
	Generic            string `json:"generic"`
	DeleteDisabledUsed string `json:"delete_disabled_used"`
	DraftRequired      string `json:"draft_required"`
}

// DefaultLabels returns the English defaults (identical to lyngua
// general/charge_policy.json for the keys it carries).
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Title:        "Charge Policies",
			TitleActive:  "Active Charge Policies",
			TitleDraft:   "Draft Charge Policies",
			TitleRetired: "Retired Charge Policies",
			Subtitle:     "How each kind of charge is accounted for, which document it goes on, and which accounts it uses",
		},
		Tabs: TabLabels{
			Active: "Active", Draft: "Draft", Retired: "Retired",
			Info: "Info", Versions: "Versions", Usage: "Where used", History: "History",
			Classification: "Classification", Components: "Documents", Postings: "Accounts",
			Approval: "Evidence & approval",
		},
		Columns: ColumnLabels{
			Code: "Code", Name: "Name", AccountingRole: "Role", DocumentKind: "Document",
			CurrentVersion: "Current version", InUse: "In use", Modified: "Modified",
			DraftVersion: "Draft version", PreparedBy: "Prepared by", RetiredOn: "Retired on",
			RetiredBy: "Retired by", Version: "Version", Status: "Status", ApprovedBy: "Approved by",
			ApprovedOn: "Approved on", ComponentRole: "Component", BookPresentation: "Presentation",
			Event: "Event", PostingRole: "Posting role", Account: "Account",
			LastVersion: "Last version", Kind: "Kind", UsageDetail: "Detail", Description: "Description",
			Rule: "Check", Passed: "Result", SelfApproved: "Self-approved",
		},
		Buttons: ButtonLabels{
			Add: "Add charge policy", Edit: "Edit", Retire: "Retire", Delete: "Delete",
			NewVersion: "New draft version", Approve: "Approve version",
			AddComponent: "Add document rule", AddPosting: "Add account mapping",
			View: "View", EditVersion: "Edit", DeleteDraft: "Delete draft", BackToPolicy: "Back",
		},
		Form: FormLabels{
			CodeLabel: "Code", CodePlaceholder: "e.g. UTILITY_RECOVERY", NameLabel: "Name",
			DescriptionLabel: "Description", TemplateLabel: "Start from", TemplatePlaceholder: "Blank",
			AccountingRoleLabel: "Accounting role", TaxPositionLabel: "Tax position",
			TaxTreatmentLabel: "Tax treatment", AssessmentScopeLabel: "Assessment scope",
			AssessmentNoteLabel:   "Assessment note",
			BookPresentationLabel: "Book presentation", ComponentRoleLabel: "Component",
			DocumentKindLabel: "Document", EventLabel: "Event", PostingRoleLabel: "Posting role",
			AccountLabel: "Account", AccountPlaceholder: "Select an account", Unspecified: "Not set",
		},
		Enums: EnumLabels{
			AccountingRolePrincipal: "Principal (own supply)",
			AccountingRoleAgent:     "Agent (on behalf of another)",
			DocumentKindInvoice:     "Invoice",
			DocumentKindRecovery:    "Recovery document",
			ComponentRoleOwnRevenue: "Own revenue",
			ComponentRoleRecovery:   "Recovered cost",
			ComponentRoleFee:        "Fee",
			VersionStatusDraft:      "Draft",
			VersionStatusApproved:   "Approved",
			VersionStatusSuperseded: "Superseded",

			BookPresentationRevenue:  "Revenue",
			BookPresentationExcluded: "Excluded reimbursement",
			TaxPositionOwnSupply:     "Own supply",
			TaxPositionExcluded:      "Excluded reimbursement",
			EventAllocation:          "Allocation",
			EventIssue:               "Issue",
			EventApplication:         "Application",
			EventReversal:            "Reversal",
			PostingRoleReceivable:    "Receivable",
			PostingRoleClearing:      "Clearing",
			PostingRoleRevenue:       "Revenue",
			PostingRoleRefund:        "Refund liability",
			PostingRoleExpense:       "Expense",
			PostingRoleCash:          "Cash",
			PolicyStatusActive:       "Active",
			PolicyStatusRetired:      "Retired",
			PolicyStatusDraft:        "Draft",
			Yes:                      "Yes",
			No:                       "No",
		},
		Detail: DetailLabels{
			CurrentVersion: "Current: v{0}",
			LockedNotice:   "Approved versions cannot be changed. Create a new draft to amend.",
			ChecklistTitle: "Before approval",
			MissingMapping: "Missing",
			BackToPolicy:   "Back to {0}",

			NoCurrentVersion:  "No approved version yet",
			ClonedFrom:        "Cloned from",
			SelfApprovedNote:  "Approved by its preparer",
			ChecklistPassed:   "Passed",
			ChecklistFailed:   "Not met",
			UnsupportedTitle:  "This combination cannot be approved yet",
			NoDescription:     "No description",
			VersionTitle:      "{0} · v{1}",
			ApprovedBy:        "Approved by",
			ApprovedAt:        "Approved on",
			Created:           "Created",
			UsageEmpty:        "Not used by any package line or agreement yet.",
			UsageInUse:        "This policy is in use, so it cannot be deleted.",
			UsageNotInUse:     "This policy is not in use.",
			OpenDraftExists:   "A draft version is already open.",
			NewVersionMessage: "Create a new draft version by copying the current approved version. Approved versions stay unchanged.",
			EvidenceTitle:     "Assessment evidence",
			ChecklistCodeLabel: map[string]string{
				"assessment_scope":               "Assessment scope is filled in",
				"posting:ISSUE:RECEIVABLE":       "Issue: receivable account mapped",
				"posting:ISSUE:CLEARING":         "Issue: clearing account mapped",
				"posting:APPLICATION:CASH":       "Application: cash account mapped",
				"posting:APPLICATION:RECEIVABLE": "Application: receivable account mapped",
			},
			ReasonLabel: map[string]string{
				"principal_excluded":                      "A principal cannot use an excluded-reimbursement tax position.",
				"accounting_role_not_agent":               "Only the agent accounting role is supported.",
				"own_supply":                              "An own-supply tax position is not supported yet.",
				"tax_position_not_excluded_reimbursement": "Only the excluded-reimbursement tax position is supported.",
				"component_count_not_one":                 "Exactly one document rule is required.",
				"fee":                                     "A fee component is not supported yet.",
				"own_revenue":                             "An own-revenue component is not supported yet.",
				"component_role_not_recovery_cost":        "The component must be a recovered cost.",
				"invoice":                                 "An invoice document is not supported yet.",
				"document_kind_not_recovery_document":     "The document must be a recovery document.",
				"component_book_presentation_not_excluded_reimbursement": "The component must be presented as an excluded reimbursement.",
			},
		},
		Confirm: ConfirmLabels{
			RetireTitle:  "Retire charge policy?",
			RetireMsg:    "It can no longer be chosen for new lines. Signed agreements keep their version.",
			ApproveTitle: "Approve this version?",
			ApproveMsg:   "Approved versions are permanent and will be used for new agreements.",
			DeleteTitle:  "Delete?",
			DeleteMsg:    "This cannot be undone.",
		},
		Empty: EmptyLabels{
			ActiveTitle: "No active charge policies", DraftTitle: "No drafts",
			RetiredTitle: "Nothing retired",
			Message:      "Add a charge policy to describe how a kind of charge is accounted for.",
		},
		Errors: ErrorLabels{
			Validation:             "The charge policy request is invalid",
			NotFound:               "Charge policy record not found",
			NotDraft:               "Only a draft version can be changed",
			DraftExists:            "The policy already has a draft version",
			NoApprovedVersion:      "The policy has no approved version",
			Retired:                "The charge policy is retired",
			SelfApproval:           "A draft you prepared or edited needs the approve-own permission to approve",
			ChecklistFailed:        "The approval checklist has failed items",
			UnsupportedCombination: "This classification is not supported yet",
			InUse:                  "The charge policy was approved or is referenced and cannot be deleted",
			TransactionRequired:    "A database transaction is required for this operation",
			TaxTreatmentNotFound:   "The tax treatment does not exist",
			Unverifiable:           "The request cannot be verified because a dependency is not wired",

			FormInvalid:        "Check the highlighted fields and try again.",
			Unavailable:        "Charge policies are not available.",
			Generic:            "Something went wrong. Please try again.",
			DeleteDisabledUsed: "In use — retire it instead.",
			DraftRequired:      "Only draft versions can be changed.",
		},
	}
}

// Format substitutes {0}, {1} placeholders (Lyngua positional convention).
func Format(tmpl string, args ...string) string {
	for i, a := range args {
		tmpl = strings.ReplaceAll(tmpl, "{"+itoa(i)+"}", a)
	}
	return tmpl
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	s := ""
	for i > 0 {
		s = string(rune('0'+i%10)) + s
		i /= 10
	}
	return s
}

// ErrorMessage maps a use-case error code (ErrorCode()) to its Lyngua message.
func (l Labels) ErrorMessage(code string) string {
	switch code {
	case ErrKindValidation:
		return l.Errors.Validation
	case ErrKindNotFound:
		return l.Errors.NotFound
	case ErrKindNotDraft:
		return l.Errors.NotDraft
	case ErrKindDraftExists:
		return l.Errors.DraftExists
	case ErrKindNoApprovedVersion:
		return l.Errors.NoApprovedVersion
	case ErrKindRetired:
		return l.Errors.Retired
	case ErrKindSelfApproval:
		return l.Errors.SelfApproval
	case ErrKindChecklistFailed:
		return l.Errors.ChecklistFailed
	case ErrKindUnsupportedCombination:
		return l.Errors.UnsupportedCombination
	case ErrKindInUse:
		return l.Errors.InUse
	case ErrKindTransactionRequired:
		return l.Errors.TransactionRequired
	case ErrKindTaxTreatmentNotFound:
		return l.Errors.TaxTreatmentNotFound
	case ErrKindUnverifiable:
		return l.Errors.Unverifiable
	default:
		return l.Errors.Generic
	}
}
