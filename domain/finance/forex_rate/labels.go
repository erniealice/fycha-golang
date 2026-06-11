package forex_rate

// ---------------------------------------------------------------------------
// Labels
// Lyngua root key: "forexRate"
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Forex Rate read-only views.
// Forex rates are read-only in the UI; rows are appended only via RecordOperatorRate.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Columns ColumnLabels `json:"columns"`
	Actions ActionLabels `json:"actions"`
	Empty   EmptyLabels  `json:"empty"`
	Detail  DetailLabels `json:"detail"`
}

// PageLabels holds page heading strings.
type PageLabels struct {
	HeadingActive     string `json:"headingActive"`
	CaptionActive     string `json:"captionActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	CaptionSuperseded string `json:"captionSuperseded"`
}

// ColumnLabels holds table column headers.
type ColumnLabels struct {
	FromCurrency   string `json:"fromCurrency"`
	ToCurrency     string `json:"toCurrency"`
	RateMicroUnits string `json:"rateMicroUnits"`
	Source         string `json:"source"`
	EffectiveFrom  string `json:"effectiveFrom"`
	Status         string `json:"status"`
}

// ActionLabels holds action button labels.
type ActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// EmptyLabels holds empty-state strings.
type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DetailLabels holds detail-page field labels.
type DetailLabels struct {
	Title           string `json:"title"`
	FromCurrency    string `json:"fromCurrency"`
	ToCurrency      string `json:"toCurrency"`
	RateMicroUnits  string `json:"rateMicroUnits"`
	Source          string `json:"source"`
	EffectiveFrom   string `json:"effectiveFrom"`
	EffectiveTo     string `json:"effectiveTo"`
	Status          string `json:"status"`
	SupersedesID    string `json:"supersedesId"`
	WorkspaceID     string `json:"workspaceId"`
	CreatedByUserID string `json:"createdByUserId"`
	Notes           string `json:"notes"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingActive:     "Forex Rates",
			CaptionActive:     "Active foreign exchange rates used during billing",
			HeadingSuperseded: "Superseded Forex Rates",
			CaptionSuperseded: "Historical foreign exchange rates (superseded by newer versions)",
		},
		Columns: ColumnLabels{
			FromCurrency:   "From",
			ToCurrency:     "To",
			RateMicroUnits: "Rate (micro-units)",
			Source:         "Source",
			EffectiveFrom:  "Effective From",
			Status:         "Status",
		},
		Actions: ActionLabels{
			View:         "View",
			NoPermission: "You do not have permission to view forex rates",
		},
		Empty: EmptyLabels{
			Title:   "No forex rates found",
			Message: "Forex rates are recorded automatically when revenue is recognized with a foreign currency.",
		},
		Detail: DetailLabels{
			Title:           "Forex Rate Detail",
			FromCurrency:    "From Currency",
			ToCurrency:      "To Currency",
			RateMicroUnits:  "Rate (micro-units)",
			Source:          "Source",
			EffectiveFrom:   "Effective From",
			EffectiveTo:     "Effective To",
			Status:          "Status",
			SupersedesID:    "Supersedes",
			WorkspaceID:     "Workspace",
			CreatedByUserID: "Created By",
			Notes:           "Notes",
		},
	}
}
