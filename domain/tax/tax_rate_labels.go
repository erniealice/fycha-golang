package tax

// ---------------------------------------------------------------------------
// TaxRateLabels
// Lyngua root key: "taxRate"
// ---------------------------------------------------------------------------

// TaxRateLabels holds all translatable strings for the Tax Rate read-only views.
// Tax rates are read-only in the UI; supersession is via admin SQL recipe.
type TaxRateLabels struct {
	Page    TaxRatePageLabels   `json:"page"`
	Columns TaxRateColumnLabels `json:"columns"`
	Actions TaxRateActionLabels `json:"actions"`
	Empty   TaxRateEmptyLabels  `json:"empty"`
	Detail  TaxRateDetailLabels `json:"detail"`
}

// TaxRatePageLabels holds page heading strings.
type TaxRatePageLabels struct {
	HeadingActive     string `json:"headingActive"`
	CaptionActive     string `json:"captionActive"`
	HeadingSuperseded string `json:"headingSuperseded"`
	CaptionSuperseded string `json:"captionSuperseded"`
}

// TaxRateColumnLabels holds table column headers.
type TaxRateColumnLabels struct {
	Jurisdiction  string `json:"jurisdiction"`
	AuthorityCode string `json:"authorityCode"`
	Kind          string `json:"kind"`
	TreatmentCode string `json:"treatmentCode"`
	Direction     string `json:"direction"`
	RateBps       string `json:"rateBps"`
	EffectiveFrom string `json:"effectiveFrom"`
	Status        string `json:"status"`
}

// TaxRateActionLabels holds action button labels.
type TaxRateActionLabels struct {
	View         string `json:"view"`
	NoPermission string `json:"noPermission"`
}

// TaxRateEmptyLabels holds empty-state strings.
type TaxRateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// TaxRateDetailLabels holds detail-page field labels.
type TaxRateDetailLabels struct {
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

// DefaultTaxRateLabels returns TaxRateLabels with sensible English defaults.
func DefaultTaxRateLabels() TaxRateLabels {
	return TaxRateLabels{
		Page: TaxRatePageLabels{
			HeadingActive:     "Tax Rates",
			CaptionActive:     "Active tax rates applied during computation",
			HeadingSuperseded: "Superseded Tax Rates",
			CaptionSuperseded: "Historical tax rates (superseded by newer versions)",
		},
		Columns: TaxRateColumnLabels{
			Jurisdiction:  "Jurisdiction",
			AuthorityCode: "Authority",
			Kind:          "Kind",
			TreatmentCode: "Treatment",
			Direction:     "Direction",
			RateBps:       "Rate (bps)",
			EffectiveFrom: "Effective From",
			Status:        "Status",
		},
		Actions: TaxRateActionLabels{
			View:         "View",
			NoPermission: "You do not have permission to view tax rates",
		},
		Empty: TaxRateEmptyLabels{
			Title:   "No tax rates found",
			Message: "Tax rates are added via the seed CSVs and superseded via admin SQL recipe.",
		},
		Detail: TaxRateDetailLabels{
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
