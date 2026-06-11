package ledger

// ---------------------------------------------------------------------------
// FiscalPeriod labels
// ---------------------------------------------------------------------------

// FiscalPeriodLabels holds all translatable strings for the fiscal period module.
type FiscalPeriodLabels struct {
	Page    FiscalPeriodPageLabels   `json:"page"`
	Buttons FiscalPeriodButtonLabels `json:"buttons"`
	Columns FiscalPeriodColumnLabels `json:"columns"`
	Status  FiscalPeriodStatusLabels `json:"status"`
	Empty   FiscalPeriodEmptyLabels  `json:"empty"`
	Actions FiscalPeriodActionLabels `json:"actions"`
	Form    FiscalPeriodFormLabels   `json:"form"`
}

// FiscalPeriodFormLabels holds field-level labels for the fiscal period add/edit form.
type FiscalPeriodFormLabels struct {
	Name         string `json:"name"`
	PeriodNumber string `json:"period_number"`
	FiscalYear   string `json:"fiscal_year"`
	StartDate    string `json:"start_date"`
	EndDate      string `json:"end_date"`
	Status       string `json:"status"`
	// Info popover text
	PeriodNumberInfo string `json:"periodNumberInfo"`
	FiscalYearInfo   string `json:"fiscalYearInfo"`
	StartDateInfo    string `json:"startDateInfo"`
	EndDateInfo      string `json:"endDateInfo"`
}

type FiscalPeriodPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type FiscalPeriodButtonLabels struct {
	AddPeriod    string `json:"addPeriod"`
	CloseYearEnd string `json:"closeYearEnd"`
}

type FiscalPeriodColumnLabels struct {
	Period    string `json:"period"`
	Year      string `json:"year"`
	StartDate string `json:"startDate"`
	EndDate   string `json:"endDate"`
	Status    string `json:"status"`
}

type FiscalPeriodStatusLabels struct {
	Open   string `json:"open"`
	Closed string `json:"closed"`
	Locked string `json:"locked"`
}

type FiscalPeriodEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

type FiscalPeriodActionLabels struct {
	Close        string `json:"close"`
	NoPermission string `json:"noPermission"`
	// Confirm messages
	ConfirmClose string `json:"confirmClose"`
}

// DefaultFiscalPeriodLabels returns FiscalPeriodLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultFiscalPeriodLabels() FiscalPeriodLabels {
	return FiscalPeriodLabels{
		Page: FiscalPeriodPageLabels{
			Heading: "Fiscal Periods",
			Caption: "Manage accounting periods and year-end close",
		},
		Buttons: FiscalPeriodButtonLabels{
			AddPeriod:    "Add Period",
			CloseYearEnd: "Close Year-End",
		},
		Columns: FiscalPeriodColumnLabels{
			Period:    "Period",
			Year:      "Year",
			StartDate: "Start Date",
			EndDate:   "End Date",
			Status:    "Status",
		},
		Status: FiscalPeriodStatusLabels{
			Open:   "Open",
			Closed: "Closed",
			Locked: "Locked",
		},
		Empty: FiscalPeriodEmptyLabels{
			Title:   "No fiscal periods found",
			Message: "Add your first fiscal period to start tracking accounting periods.",
		},
		Actions: FiscalPeriodActionLabels{
			Close:        "Close",
			NoPermission: "No permission",
			ConfirmClose: "Are you sure you want to close %s? This will prevent new journal entries from being posted to this period.",
		},
		Form: FiscalPeriodFormLabels{
			Name:             "Name",
			PeriodNumber:     "Period Number",
			FiscalYear:       "Fiscal Year",
			StartDate:        "Start Date",
			EndDate:          "End Date",
			Status:           "Status",
			PeriodNumberInfo: "Sequential number of this period within the fiscal year (e.g. 1 for January).",
			FiscalYearInfo:   "The fiscal year this period belongs to (e.g. 2026).",
			StartDateInfo:    "First day of this accounting period. Journal entries dated on or after this date fall within the period.",
			EndDateInfo:      "Last day of this accounting period. Journal entries dated on or before this date fall within the period.",
		},
	}
}
