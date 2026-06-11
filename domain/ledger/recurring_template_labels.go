package ledger

// ---------------------------------------------------------------------------
// RecurringTemplate labels
// ---------------------------------------------------------------------------

// RecurringTemplateLabels holds all translatable strings for the recurring journal template module.
type RecurringTemplateLabels struct {
	Page      RecurringTemplatePageLabels      `json:"page"`
	Buttons   RecurringTemplateButtonLabels    `json:"buttons"`
	Columns   RecurringTemplateColumnLabels    `json:"columns"`
	Frequency RecurringTemplateFrequencyLabels `json:"frequency"`
	Status    RecurringTemplateStatusLabels    `json:"status"`
	Empty     RecurringTemplateEmptyLabels     `json:"empty"`
	Actions   RecurringTemplateActionLabels    `json:"actions"`
}

type RecurringTemplatePageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type RecurringTemplateButtonLabels struct {
	AddTemplate string `json:"addTemplate"`
}

type RecurringTemplateColumnLabels struct {
	Name      string `json:"name"`
	Frequency string `json:"frequency"`
	NextRun   string `json:"nextRun"`
	Status    string `json:"status"`
}

type RecurringTemplateFrequencyLabels struct {
	Daily     string `json:"daily"`
	Weekly    string `json:"weekly"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Yearly    string `json:"yearly"`
}

type RecurringTemplateStatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type RecurringTemplateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type RecurringTemplateActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Pause        string `json:"pause"`
	Resume       string `json:"resume"`
	NoPermission string `json:"noPermission"`
}

// DefaultRecurringTemplateLabels returns RecurringTemplateLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultRecurringTemplateLabels() RecurringTemplateLabels {
	return RecurringTemplateLabels{
		Page: RecurringTemplatePageLabels{
			Heading: "Recurring Journal Entries",
			Caption: "Automated entries for depreciation, amortization, accruals",
		},
		Buttons: RecurringTemplateButtonLabels{
			AddTemplate: "Add Recurring Entry",
		},
		Columns: RecurringTemplateColumnLabels{
			Name:      "Name",
			Frequency: "Frequency",
			NextRun:   "Next Run",
			Status:    "Status",
		},
		Frequency: RecurringTemplateFrequencyLabels{
			Daily:     "Daily",
			Weekly:    "Weekly",
			Monthly:   "Monthly",
			Quarterly: "Quarterly",
			Yearly:    "Yearly",
		},
		Status: RecurringTemplateStatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		Empty: RecurringTemplateEmptyLabels{
			Title:   "No recurring templates",
			Message: "Add your first recurring journal entry template to automate periodic entries.",
		},
		Actions: RecurringTemplateActionLabels{
			Edit:         "Edit",
			Delete:       "Delete",
			Pause:        "Pause",
			Resume:       "Resume",
			NoPermission: "No permission",
		},
	}
}
