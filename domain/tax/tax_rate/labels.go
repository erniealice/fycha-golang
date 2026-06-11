package tax_rate

// ---------------------------------------------------------------------------
// Labels
// Lyngua root key: "taxRate"
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Tax Rate read-only views.
// Tax rates are read-only in the UI; supersession is via admin SQL recipe.
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
	Jurisdiction  string `json:"jurisdiction"`
	AuthorityCode string `json:"authorityCode"`
	Kind          string `json:"kind"`
	TreatmentCode string `json:"treatmentCode"`
	Direction     string `json:"direction"`
	RateBps       string `json:"rateBps"`
	EffectiveFrom string `json:"effectiveFrom"`
	Status        string `json:"status"`
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
	Title         string `json:"title"`
	Jurisdiction  string `json:"jurisdiction"`
	AuthorityCode string `json:"authorityCode"`
	Kind          string `json:"kind"`
	TreatmentCode string `json:"treatmentCode"`
	Direction     string `json:"direction"`
	RateBps       string `json:"rateBps"`
	EffectiveFrom string `json:"effectiveFrom"`
	EffectiveTo   string `json:"effectiveTo"`
	Status        string `json:"status"`
	SupersedesID  string `json:"supersedesId"`
	WorkspaceID   string `json:"workspaceId"`
}

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingActive:     "Tax Rates",
			CaptionActive:     "Active tax rates applied during computation",
			HeadingSuperseded: "Superseded Tax Rates",
			CaptionSuperseded: "Historical tax rates (superseded by newer versions)",
		},
		Columns: ColumnLabels{
			Jurisdiction:  "Jurisdiction",
			AuthorityCode: "Authority",
			Kind:          "Kind",
			TreatmentCode: "Treatment",
			Direction:     "Direction",
			RateBps:       "Rate (bps)",
			EffectiveFrom: "Effective From",
			Status:        "Status",
		},
		Actions: ActionLabels{
			View:         "View",
			NoPermission: "You do not have permission to view tax rates",
		},
		Empty: EmptyLabels{
			Title:   "No tax rates found",
			Message: "Tax rates are added via the seed CSVs and superseded via admin SQL recipe.",
		},
		Detail: DetailLabels{
			Title:         "Tax Rate Detail",
			Jurisdiction:  "Jurisdiction",
			AuthorityCode: "Authority Code",
			Kind:          "Kind",
			TreatmentCode: "Treatment",
			Direction:     "Direction",
			RateBps:       "Rate (basis points)",
			EffectiveFrom: "Effective From",
			EffectiveTo:   "Effective To",
			Status:        "Status",
			SupersedesID:  "Supersedes",
			WorkspaceID:   "Workspace",
		},
	}
}
