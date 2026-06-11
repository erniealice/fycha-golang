package treasury

// ---------------------------------------------------------------------------
// PettyCashLabels (Cash — Petty Cash)
// ---------------------------------------------------------------------------

// PettyCashLabels holds all translatable strings for the Petty Cash module.
type PettyCashLabels struct {
	Page    PettyCashPageLabels   `json:"page"`
	Buttons PettyCashButtonLabels `json:"buttons"`
	Columns PettyCashColumnLabels `json:"columns"`
	Status  PettyCashStatusLabels `json:"status"`
	Empty   PettyCashEmptyLabels  `json:"empty"`
	Form    PettyCashFormLabels   `json:"form"`
	Actions PettyCashActionLabels `json:"actions"`
}

type PettyCashPageLabels struct {
	RegisterHeading          string `json:"registerHeading"`
	RegisterCaption          string `json:"registerCaption"`
	ReplenishmentsHeading    string `json:"replenishmentsHeading"`
	ReplenishmentsCaption    string `json:"replenishmentsCaption"`
	CustodianBalancesHeading string `json:"custodianBalancesHeading"`
	CustodianBalancesCaption string `json:"custodianBalancesCaption"`
}

type PettyCashButtonLabels struct {
	AddFund   string `json:"addFund"`
	Replenish string `json:"replenish"`
}

type PettyCashColumnLabels struct {
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

type PettyCashStatusLabels struct {
	Active   string `json:"active"`
	Inactive string `json:"inactive"`
}

type PettyCashEmptyLabels struct {
	RegisterTitle         string `json:"registerTitle"`
	RegisterMessage       string `json:"registerMessage"`
	ReplenishmentsTitle   string `json:"replenishmentsTitle"`
	ReplenishmentsMessage string `json:"replenishmentsMessage"`
	CustodianTitle        string `json:"custodianTitle"`
	CustodianMessage      string `json:"custodianMessage"`
}

type PettyCashFormLabels struct {
	Name                  string `json:"name"`
	NamePlaceholder       string `json:"namePlaceholder"`
	AuthorizedAmount      string `json:"authorizedAmount"`
	AuthorizedPlaceholder string `json:"authorizedPlaceholder"`
	CustodianID           string `json:"custodianId"`
	LocationID            string `json:"locationId"`
}

type PettyCashActionLabels struct {
	View          string `json:"view"`
	Replenish     string `json:"replenish"`
	Delete        string `json:"delete"`
	NoPermission  string `json:"noPermission"`
	ConfirmDelete string `json:"confirmDelete"`
}

// DefaultPettyCashLabels returns PettyCashLabels with hardcoded English defaults.
func DefaultPettyCashLabels() PettyCashLabels {
	return PettyCashLabels{
		Page: PettyCashPageLabels{
			RegisterHeading:          "Petty Cash Register",
			RegisterCaption:          "Manage petty cash funds across locations and custodians",
			ReplenishmentsHeading:    "Petty Cash Replenishments",
			ReplenishmentsCaption:    "Track fund replenishments and reimbursements",
			CustodianBalancesHeading: "Custodian Balances",
			CustodianBalancesCaption: "Current balance summary by custodian",
		},
		Buttons: PettyCashButtonLabels{
			AddFund:   "Add Fund",
			Replenish: "Replenish",
		},
		Columns: PettyCashColumnLabels{
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
		Status: PettyCashStatusLabels{
			Active:   "Active",
			Inactive: "Inactive",
		},
		Empty: PettyCashEmptyLabels{
			RegisterTitle:         "No petty cash funds",
			RegisterMessage:       "Set up petty cash funds for each location or department.",
			ReplenishmentsTitle:   "No replenishments",
			ReplenishmentsMessage: "Replenishment records will appear here when funds are restocked.",
			CustodianTitle:        "No custodian data",
			CustodianMessage:      "Assign custodians to funds to see balance summaries here.",
		},
		Form: PettyCashFormLabels{
			Name:                  "Fund Name",
			NamePlaceholder:       "e.g. Main Office Petty Cash",
			AuthorizedAmount:      "Authorized Amount",
			AuthorizedPlaceholder: "0.00",
			CustodianID:           "Custodian",
			LocationID:            "Location",
		},
		Actions: PettyCashActionLabels{
			View:          "View",
			Replenish:     "Replenish",
			Delete:        "Delete",
			NoPermission:  "No permission",
			ConfirmDelete: "Are you sure you want to delete this petty cash fund? This action cannot be undone.",
		},
	}
}
