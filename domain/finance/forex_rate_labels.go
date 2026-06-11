package finance

// ---------------------------------------------------------------------------
// ForexRateLabels
// Lyngua root key: "forexRate"
// ---------------------------------------------------------------------------

// ForexRateLabels holds all translatable strings for the Forex Rate read-only views.
// Forex rates are read-only in the UI; rows are appended only via RecordOperatorRate.
type ForexRateLabels struct {
	Page    ForexRatePageLabels   `json:"page"`
	Columns ForexRateColumnLabels `json:"columns"`
	Actions ForexRateActionLabels `json:"actions"`
	Empty   ForexRateEmptyLabels  `json:"empty"`
	Detail  ForexRateDetailLabels `json:"detail"`
}

// ForexRatePageLabels holds page heading strings.
type ForexRatePageLabels struct {
	HeadingActive     string `json:"headingActive"`
	CaptionActive     string `json:"captionActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	CaptionSuperseded string `json:"captionSuperseded"`
}

// ForexRateColumnLabels holds table column headers.
type ForexRateColumnLabels struct {
	FromCurrency   string `json:"fromCurrency"`
	ToCurrency     string `json:"toCurrency"`
	RateMicroUnits string `json:"rateMicroUnits"`
	Source         string `json:"source"`
	EffectiveFrom  string `json:"effectiveFrom"`
	Status         string `json:"status"`
}

// ForexRateActionLabels holds action button labels.
type ForexRateActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// ForexRateEmptyLabels holds empty-state strings.
type ForexRateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// ForexRateDetailLabels holds detail-page field labels.
type ForexRateDetailLabels struct {
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

// DefaultForexRateLabels returns ForexRateLabels with sensible English defaults.
func DefaultForexRateLabels() ForexRateLabels {
	return ForexRateLabels{
		Page: ForexRatePageLabels{
			HeadingActive:     "Forex Rates",
			CaptionActive:     "Active foreign exchange rates used during billing",
			HeadingSuperseded: "Superseded Forex Rates",
			CaptionSuperseded: "Historical foreign exchange rates (superseded by newer versions)",
		},
		Columns: ForexRateColumnLabels{
			FromCurrency:   "From",
			ToCurrency:     "To",
			RateMicroUnits: "Rate (micro-units)",
			Source:         "Source",
			EffectiveFrom:  "Effective From",
			Status:         "Status",
		},
		Actions: ForexRateActionLabels{
			View:         "View",
			NoPermission: "You do not have permission to view forex rates",
		},
		Empty: ForexRateEmptyLabels{
			Title:   "No forex rates found",
			Message: "Forex rates are recorded automatically when revenue is recognized with a foreign currency.",
		},
		Detail: ForexRateDetailLabels{
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
