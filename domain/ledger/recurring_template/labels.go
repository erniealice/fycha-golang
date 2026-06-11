package recurring_template

// ---------------------------------------------------------------------------
// RecurringTemplate labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the recurring journal template module.
type Labels struct {
	Page      PageLabels      `json:"page"`
	Buttons   ButtonLabels    `json:"buttons"`
	Columns   ColumnLabels    `json:"columns"`
	Frequency FrequencyLabels `json:"frequency"`
	Status    StatusLabels    `json:"status"`
	Empty     EmptyLabels     `json:"empty"`
	Actions   ActionLabels    `json:"actions"`
}

type PageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type ButtonLabels struct {
	AddTemplate string `json:"addTemplate"`
}

type ColumnLabels struct {
	Name      string `json:"name"`
	Frequency string `json:"frequency"`
	NextRun   string `json:"nextRun"`
	Status    string `json:"status"`
}

type FrequencyLabels struct {
	Daily     string `json:"daily"`
	Weekly    string `json:"weekly"`
	Monthly   string `json:"monthly"`
	Quarterly string `json:"quarterly"`
	Yearly    string `json:"yearly"`
}

type StatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type ActionLabels struct {
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	Pause        string `json:"pause"`
	Resume       string `json:"resume"`
	NoPermission string `json:"noPermission"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			Heading: "Recurring Journal Entries",
			Caption: "Automated entries for depreciation, amortization, accruals",
		},
		Buttons: ButtonLabels{
			AddTemplate: "Add Recurring Entry",
		},
		Columns: ColumnLabels{
			Name:      "Name",
			Frequency: "Frequency",
			NextRun:   "Next Run",
			Status:    "Status",
		},
		Frequency: FrequencyLabels{
			Daily:     "Daily",
			Weekly:    "Weekly",
			Monthly:   "Monthly",
			Quarterly: "Quarterly",
			Yearly:    "Yearly",
		},
		Status: StatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		Empty: EmptyLabels{
			Title:   "No recurring templates",
			Message: "Add your first recurring journal entry template to automate periodic entries.",
		},
		Actions: ActionLabels{
			Edit:         "Edit",
			Delete:       "Delete",
			Pause:        "Pause",
			Resume:       "Resume",
			NoPermission: "No permission",
		},
	}
}
