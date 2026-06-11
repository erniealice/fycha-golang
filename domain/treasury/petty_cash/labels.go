package petty_cash

// ---------------------------------------------------------------------------
// PettyCash labels (Cash — Petty Cash)
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Petty Cash module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Buttons ButtonLabels `json:"buttons"`
	Columns ColumnLabels `json:"columns"`
	Status  StatusLabels `json:"status"`
	Empty   EmptyLabels  `json:"empty"`
	Form    FormLabels   `json:"form"`
	Actions ActionLabels `json:"actions"`
}

type PageLabels struct {
	RegisterHeading          string `json:"registerHeading"`
	RegisterCaption          string `json:"registerCaption"`
	ReplenishmentsHeading    string `json:"replenishmentsHeading"`
	ReplenishmentsCaption    string `json:"replenishmentsCaption"`
	CustodianBalancesHeading string `json:"custodianBalancesHeading"`
	CustodianBalancesCaption string `json:"custodianBalancesCaption"`
}

type ButtonLabels struct {
	AddFund   string `json:"addFund"`
	Replenish string `json:"replenish"`
}

type ColumnLabels struct {
	// Register columns
	Name             string `json:"name"`
	AuthorizedAmount string `json:"authorizedAmount"`
	CurrentBalance   string `json:"currentBalance"`
	Custodian        string `json:"custodian"`
	Location         string `json:"location"`
	Status           string `json:"status"`
	// Replenishment columns
	Fund   string `json:"fund"`
	Amount string `json:"amount"`
	Date   string `json:"date"`
	Notes  string `json:"notes"`
	// Custodian balance columns
	TotalFunds   string `json:"totalFunds"`
	TotalBalance string `json:"totalBalance"`
}

type StatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type EmptyLabels struct {
	RegisterTitle         string `json:"registerTitle"`
	RegisterMessage       string `json:"registerMessage"`
	ReplenishmentsTitle   string `json:"replenishmentsTitle"`
	ReplenishmentsMessage string `json:"replenishmentsMessage"`
	CustodianTitle        string `json:"custodianTitle"`
	CustodianMessage      string `json:"custodianMessage"`
}

type FormLabels struct {
	Name                  string `json:"name"`
	NamePlaceholder       string `json:"namePlaceholder"`
	AuthorizedAmount      string `json:"authorizedAmount"`
	AuthorizedPlaceholder string `json:"authorizedPlaceholder"`
	CustodianID           string `json:"custodianId"`
	LocationID            string `json:"locationId"`
}

type ActionLabels struct {
	View          string `json:"view"`
	Replenish     string `json:"replenish"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			RegisterHeading:          "Petty Cash Register",
			RegisterCaption:          "Manage petty cash funds across locations and custodians",
			ReplenishmentsHeading:    "Petty Cash Replenishments",
			ReplenishmentsCaption:    "Track fund replenishments and reimbursements",
			CustodianBalancesHeading: "Custodian Balances",
			CustodianBalancesCaption: "Current balance summary by custodian",
		},
		Buttons: ButtonLabels{
			AddFund:   "Add Fund",
			Replenish: "Replenish",
		},
		Columns: ColumnLabels{
			Name:             "Fund Name",
			AuthorizedAmount: "Authorized Amount",
			CurrentBalance:   "Current Balance",
			Custodian:        "Custodian",
			Location:         "Location",
			Status:           "Status",
			Fund:             "Fund",
			Amount:           "Amount",
			Date:             "Date",
			Notes:            "Notes",
			TotalFunds:       "Total Funds",
			TotalBalance:     "Total Balance",
		},
		Status: StatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		Empty: EmptyLabels{
			RegisterTitle:         "No petty cash funds",
			RegisterMessage:       "Set up petty cash funds for each location or department.",
			ReplenishmentsTitle:   "No replenishments",
			ReplenishmentsMessage: "Replenishment records will appear here when funds are restocked.",
			CustodianTitle:        "No custodian data",
			CustodianMessage:      "Assign custodians to funds to see balance summaries here.",
		},
		Form: FormLabels{
			Name:                  "Fund Name",
			NamePlaceholder:       "e.g. Main Office Petty Cash",
			AuthorizedAmount:      "Authorized Amount",
			AuthorizedPlaceholder: "0.00",
			CustodianID:           "Custodian",
			LocationID:            "Location",
		},
		Actions: ActionLabels{
			View:          "View",
			Replenish:     "Replenish",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this petty cash fund? This action cannot be undone.",
		},
	}
}
