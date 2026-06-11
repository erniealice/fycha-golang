package depreciationrun

// ---------------------------------------------------------------------------
// Depreciation Run labels (Surfaces A, B, C, D, F) + Revaluation labels (Surface E)
// Lyngua root key: "depreciationRun" / "assetRevaluation" / "depreciationPolicies"
// Naming: depreciationRun / DepreciationRun / depreciation_run / depreciation-run everywhere
// except user-visible VALUES supplied by lyngua (e.g. "Lapsing Schedule").
// ---------------------------------------------------------------------------

// Labels is the top-level struct for all Depreciation Run copy.
type Labels struct {
	AppLabel string `json:"appLabel"`
	// Surface A — per-asset drawer
	RunForm FormLabels `json:"runForm"`
	// Surface B — Lapsing Schedule list page (workspace overview)
	LapsingSchedule LapsingScheduleLabels `json:"lapsingSchedule"`
	// Surface C — per-category / per-policy drawer
	CategoryRunForm CategoryFormLabels `json:"categoryRunForm"`
	// Surface D — run history list + detail
	List   ListLabels   `json:"list"`
	Detail DetailLabels `json:"detail"`
	// Status badges shared across Surfaces B and D
	StatusBadges StatusBadgeLabels `json:"statusBadges"`
	// Scope kind display labels
	ScopeKind ScopeKindLabels `json:"scopeKind"`
	// Entry outcome labels
	EntryOutcome EntryOutcomeLabels `json:"entryOutcome"`
	// Cross-cutting toast labels
	Toast ToastLabels `json:"toast"`
	// Errors
	Errors ErrorLabels `json:"errors"`
	// Cross-cutting asset-edit field labels
	AssetEditDepreciationFieldsLockedWarning  string `json:"assetEditDepreciationFieldsLockedWarning"`
	AssetEditUnitsOfProductionDisabledTooltip string `json:"assetEditUnitsOfProductionDisabledTooltip"`
}

// FormLabels holds copy for the Surface A per-asset drawer.
type FormLabels struct {
	Title            string `json:"title"`
	SubtitleTemplate string `json:"subtitleTemplate"`
	// AsOfDate input
	AsOfDateLabel string `json:"asOfDateLabel"`
	AsOfDateHint  string `json:"asOfDateHint"`
	// Pending-periods table column headers
	ColPeriod         string `json:"colPeriod"`
	ColAmount         string `json:"colAmount"`
	ColAccumulated    string `json:"colAccumulated"`
	ColBookValueAfter string `json:"colBookValueAfter"`
	// Generate button (label and count-template variant)
	GenerateLabel             string `json:"generateLabel"`
	GenerateWithCountTemplate string `json:"generateWithCountTemplate"`
	CancelLabel               string `json:"cancelLabel"`
	// Blocker chip labels
	BlockerNotInService     string `json:"blockerNotInService"`
	BlockerFullyDepreciated string `json:"blockerFullyDepreciated"`
	BlockerMissingMethod    string `json:"blockerMissingMethod"`
	BlockerUnitsRequired    string `json:"blockerUnitsRequired"`
	// UoP-specific blocker messaging (rendered as a translated message + link in the drawer)
	BlockerUnitsRequiredMessage string `json:"blockerUnitsRequiredMessage"`
	BlockerUnitsRequiredLink    string `json:"blockerUnitsRequiredLink"`
	// Empty state (no pending periods)
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// LapsingScheduleLabels holds copy for the Surface B workspace
// lapsing-schedule list page.
type LapsingScheduleLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// AsOfDate input
	AsOfDateLabel string `json:"asOfDateLabel"`
	// Table column labels
	Columns LapsingScheduleColumnLabels `json:"columns"`
	// Status badge variants
	StatusUpToDate                string `json:"statusUpToDate"`
	StatusNPeriodsPendingTemplate string `json:"statusNPeriodsPendingTemplate"`
	StatusNotStarted              string `json:"statusNotStarted"`
	StatusBlockedTemplate         string `json:"statusBlockedTemplate"`
	// BlockedPrefix is the human-readable prefix for blocked-status badges
	// when a specific BlockerLabel is provided (e.g. "Blocked: Units required").
	// The trailing space is intentional — it is concatenated with the reason string.
	BlockedPrefix string `json:"blockedPrefix"`
	// Bulk action labels
	BulkRunForSelected    string `json:"bulkRunForSelected"`
	BulkRunForAllMatching string `json:"bulkRunForAllMatching"`
	// Filter chip labels
	FilterCategory string `json:"filterCategory"`
	FilterPolicy   string `json:"filterPolicy"`
	FilterStatus   string `json:"filterStatus"`
	FilterCurrency string `json:"filterCurrency"`
	// Empty state
	EmptyTitle   string `json:"emptyTitle"`
	EmptyMessage string `json:"emptyMessage"`
}

// LapsingScheduleColumnLabels holds column headers for Surface B.
type LapsingScheduleColumnLabels struct {
	Asset             string `json:"asset"`
	Category          string `json:"category"`
	Policy            string `json:"policy"`
	Currency          string `json:"currency"`
	CurrentBookValue  string `json:"currentBookValue"`
	LastPostedPeriod  string `json:"lastPostedPeriod"`
	NextPendingPeriod string `json:"nextPendingPeriod"`
	Pending           string `json:"pending"`
	NextAmount        string `json:"nextAmount"`
	Status            string `json:"status"`
	Actions           string `json:"actions"`
}

// CategoryFormLabels holds copy for the Surface C per-category /
// per-policy drawer. The same drawer serves both entry points; only the breadcrumb differs.
type CategoryFormLabels struct {
	// Category breadcrumb variant
	TitleCategory string `json:"titleCategory"`
	// Policy breadcrumb variant
	TitlePolicy      string `json:"titlePolicy"`
	SubtitleTemplate string `json:"subtitleTemplate"`
	// Per-asset row column headers
	ColAsset    string `json:"colAsset"`
	ColMethod   string `json:"colMethod"`
	ColPending  string `json:"colPending"`
	ColAmount   string `json:"colAmount"`
	ColBlockers string `json:"colBlockers"`
	// Submit and cancel
	SubmitLabel string `json:"submitLabel"`
	CancelLabel string `json:"cancelLabel"`
}

// ListLabels holds copy for the Surface D run-history list page.
type ListLabels struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	// Table column labels
	Columns ListColumnLabels `json:"columns"`
	// Status filter chip labels
	FilterPending  string `json:"filterPending"`
	FilterComplete string `json:"filterComplete"`
	FilterFailed   string `json:"filterFailed"`
	// Stale-pending warning
	StalePendingWarning string `json:"stalePendingWarning"`
	// Per-status empty states
	Empty ListEmptyLabels `json:"empty"`
}

// ListColumnLabels holds column headers for Surface D list.
type ListColumnLabels struct {
	ID          string `json:"id"`
	Scope       string `json:"scope"`
	AsOfDate    string `json:"asOfDate"`
	Initiator   string `json:"initiator"`
	InitiatedAt string `json:"initiatedAt"`
	Status      string `json:"status"`
	Created     string `json:"created"`
	Skipped     string `json:"skipped"`
	Errored     string `json:"errored"`
	Actions     string `json:"actions"`
}

// ListEmptyLabels holds per-status empty-state copy for Surface D list.
type ListEmptyLabels struct {
	Pending  ListEmptyStateLabels `json:"pending"`
	Complete ListEmptyStateLabels `json:"complete"`
	Failed   ListEmptyStateLabels `json:"failed"`
}

// ListEmptyStateLabels holds title + message for one empty-state variant.
type ListEmptyStateLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DetailLabels holds copy for the Surface D run-detail page.
type DetailLabels struct {
	Title   string          `json:"title"`
	Tabs    DetailTabLabels `json:"tabs"`
	Summary SummaryLabels   `json:"summary"`
}

// DetailTabLabels holds tab labels for the Surface D detail page.
type DetailTabLabels struct {
	Summary      string `json:"summary"`
	Selections   string `json:"selections"`
	Results      string `json:"results"`
	Transactions string `json:"transactions"`
	History      string `json:"history"`
}

// SummaryLabels holds stat-card labels for the Surface D summary tab.
type SummaryLabels struct {
	Scope                   string `json:"scope"`
	AsOfDate                string `json:"asOfDate"`
	Initiator               string `json:"initiator"`
	InitiatedAt             string `json:"initiatedAt"`
	CompletedAt             string `json:"completedAt"`
	Status                  string `json:"status"`
	Created                 string `json:"created"`
	Skipped                 string `json:"skipped"`
	Errored                 string `json:"errored"`
	Totals                  string `json:"totals"`
	PossiblyInterruptedNote string `json:"possiblyInterruptedNote"`
}

// StatusBadgeLabels holds display labels for each run status value.
type StatusBadgeLabels struct {
	Pending             string `json:"pending"`
	Complete            string `json:"complete"`
	Failed              string `json:"failed"`
	PossiblyInterrupted string `json:"possiblyInterrupted"`
}

// ScopeKindLabels holds display labels for each scope kind enum value.
type ScopeKindLabels struct {
	Asset     string `json:"asset"`
	Category  string `json:"category"`
	Policy    string `json:"policy"`
	Workspace string `json:"workspace"`
}

// EntryOutcomeLabels holds display labels for per-entry outcome values.
type EntryOutcomeLabels struct {
	Created string `json:"created"`
	Skipped string `json:"skipped"`
	Errored string `json:"errored"`
}

// ToastLabels holds toast message templates for all Depreciation Run surfaces.
type ToastLabels struct {
	// SuccessTemplate supports {{.Created}}/{{.Skipped}}/{{.Errored}} placeholders.
	SuccessTemplate string `json:"successTemplate"`
	// SkippedTemplate is shown when created_count=0 and skipped_count>0.
	SkippedTemplate string `json:"skippedTemplate"`
	// ErroredTemplate is shown when errored_count>0.
	ErroredTemplate string `json:"erroredTemplate"`
	// ViewRunLink is the link label used on single-run toasts.
	ViewRunLink string `json:"viewRunLink"`
}

// ErrorLabels holds error message strings for the depreciation-run module.
type ErrorLabels struct {
	RunForSelectedCapExceededError string `json:"runForSelectedCapExceededError"`
	PermissionDenied               string `json:"permissionDenied"`
	UseCaseUnavailable             string `json:"useCaseUnavailable"`
	FormParseFailed                string `json:"formParseFailed"`
	GenerateFailed                 string `json:"generateFailed"`
	InvalidSelection               string `json:"invalidSelection"`
	IdempotencyConflict            string `json:"idempotencyConflict"`
	WorkspaceMismatch              string `json:"workspaceMismatch"`
}

// DefaultLabels returns Labels with sensible English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		AppLabel: "Lapsing Schedule",
		RunForm: FormLabels{
			Title:                       "Run Depreciation",
			SubtitleTemplate:            "Post depreciation for {{.AssetName}} through {{.AsOfDate}}",
			AsOfDateLabel:               "As of date",
			AsOfDateHint:                "Periods up to and including this date will be posted.",
			ColPeriod:                   "Period",
			ColAmount:                   "Amount",
			ColAccumulated:              "Accumulated",
			ColBookValueAfter:           "Book value after",
			GenerateLabel:               "Generate",
			GenerateWithCountTemplate:   "Generate ({{.Count}})",
			CancelLabel:                 "Cancel",
			BlockerNotInService:         "Not in service",
			BlockerFullyDepreciated:     "Fully depreciated",
			BlockerMissingMethod:        "Missing depreciation method",
			BlockerUnitsRequired:        "Units of Production not yet supported",
			BlockerUnitsRequiredMessage: "Units of Production depreciation requires per-period units input. See the future-release plan.",
			BlockerUnitsRequiredLink:    "Open future-release plan",
			EmptyTitle:                  "No pending periods",
			EmptyMessage:                "All periods up to the selected date have been posted.",
		},
		LapsingSchedule: LapsingScheduleLabels{
			Title:         "Lapsing Schedule",
			Subtitle:      "In-service assets with pending depreciation periods",
			AsOfDateLabel: "As of date",
			Columns: LapsingScheduleColumnLabels{
				Asset:             "Asset",
				Category:          "Category",
				Policy:            "Policy",
				Currency:          "Currency",
				CurrentBookValue:  "Book Value",
				LastPostedPeriod:  "Last posted",
				NextPendingPeriod: "Next pending",
				Pending:           "Pending",
				NextAmount:        "Next amount",
				Status:            "Status",
				Actions:           "Actions",
			},
			StatusUpToDate:                "Up to date",
			StatusNPeriodsPendingTemplate: "{{.Count}} periods pending",
			StatusNotStarted:              "Not started",
			StatusBlockedTemplate:         "Blocked: {{.Reason}}",
			BlockedPrefix:                 "Blocked: ",
			BulkRunForSelected:            "Run for selected",
			BulkRunForAllMatching:         "Run for all matching",
			FilterCategory:                "Category",
			FilterPolicy:                  "Policy",
			FilterStatus:                  "Status",
			FilterCurrency:                "Currency",
			EmptyTitle:                    "No assets in service",
			EmptyMessage:                  "Add an asset to get started.",
		},
		CategoryRunForm: CategoryFormLabels{
			TitleCategory:    "Run depreciation for category",
			TitlePolicy:      "Run depreciation for policy",
			SubtitleTemplate: "{{.Count}} assets eligible",
			ColAsset:         "Asset",
			ColMethod:        "Method",
			ColPending:       "Pending",
			ColAmount:        "Amount",
			ColBlockers:      "Blockers",
			SubmitLabel:      "Run depreciation",
			CancelLabel:      "Cancel",
		},
		List: ListLabels{
			Title:    "Depreciation Runs",
			Subtitle: "History of depreciation run batches",
			Columns: ListColumnLabels{
				ID:          "Run ID",
				Scope:       "Scope",
				AsOfDate:    "As of date",
				Initiator:   "Initiator",
				InitiatedAt: "Initiated",
				Status:      "Status",
				Created:     "Created",
				Skipped:     "Skipped",
				Errored:     "Errored",
				Actions:     "Actions",
			},
			FilterPending:       "Pending",
			FilterComplete:      "Complete",
			FilterFailed:        "Failed",
			StalePendingWarning: "This run has been pending for an unusually long time and may have been interrupted.",
			Empty: ListEmptyLabels{
				Pending: ListEmptyStateLabels{
					Title:   "No pending runs",
					Message: "There are no depreciation runs currently in progress.",
				},
				Complete: ListEmptyStateLabels{
					Title:   "No completed runs",
					Message: "No depreciation runs have completed yet.",
				},
				Failed: ListEmptyStateLabels{
					Title:   "No failed runs",
					Message: "No depreciation runs have failed.",
				},
			},
		},
		Detail: DetailLabels{
			Title: "Depreciation Run",
			Tabs: DetailTabLabels{
				Summary:      "Summary",
				Selections:   "Selections",
				Results:      "Results",
				Transactions: "Transactions",
				History:      "History",
			},
			Summary: SummaryLabels{
				Scope:                   "Scope",
				AsOfDate:                "As of date",
				Initiator:               "Initiator",
				InitiatedAt:             "Initiated",
				CompletedAt:             "Completed",
				Status:                  "Status",
				Created:                 "Created",
				Skipped:                 "Skipped",
				Errored:                 "Errored",
				Totals:                  "Totals",
				PossiblyInterruptedNote: "This run may have been interrupted before completing. Some periods may be missing.",
			},
		},
		StatusBadges: StatusBadgeLabels{
			Pending:             "Pending",
			Complete:            "Complete",
			Failed:              "Failed",
			PossiblyInterrupted: "Possibly interrupted",
		},
		ScopeKind: ScopeKindLabels{
			Asset:     "Asset",
			Category:  "Category",
			Policy:    "Policy",
			Workspace: "Workspace",
		},
		EntryOutcome: EntryOutcomeLabels{
			Created: "Created",
			Skipped: "Skipped",
			Errored: "Errored",
		},
		Toast: ToastLabels{
			SuccessTemplate: "{{.Created}} periods posted, {{.Skipped}} skipped, {{.Errored}} errored",
			SkippedTemplate: "{{.Skipped}} periods already posted (skipped)",
			ErroredTemplate: "{{.Errored}} periods failed to post",
			ViewRunLink:     "View run",
		},
		Errors: ErrorLabels{
			RunForSelectedCapExceededError: "Batch cap exceeded — maximum 500 assets per run. Narrow the filter to run the rest.",
			PermissionDenied:               "You do not have permission to run depreciation.",
			UseCaseUnavailable:             "Service unavailable. Please try again.",
			FormParseFailed:                "Form data could not be read.",
			GenerateFailed:                 "Failed to record the depreciation run.",
			InvalidSelection:               "One or more selected assets are invalid.",
			IdempotencyConflict:            "Depreciation for one or more periods has already been posted.",
			WorkspaceMismatch:              "Selected assets belong to a different workspace.",
		},
		AssetEditDepreciationFieldsLockedWarning:  "Posted depreciation exists for this asset. Changing depreciation configuration requires a Useful Life Change action (not yet available — see Run history for posted periods).",
		AssetEditUnitsOfProductionDisabledTooltip: "Units of Production depreciation is not yet supported.",
	}
}
