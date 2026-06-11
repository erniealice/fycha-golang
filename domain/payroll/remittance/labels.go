// ---------------------------------------------------------------------------
// Payroll Remittance labels
// ---------------------------------------------------------------------------

package remittance

// Labels holds all translatable strings for the Payroll Remittance sub-module.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Tabs    TabLabels    `json:"tabs"`
	Columns ColumnLabels `json:"columns"`
	Types   TypeLabels   `json:"types"`
	Empty   EmptyLabels  `json:"empty"`
}

type PageLabels struct {
	HeadingPending  string `json:"headingPending"`
	SubtitlePending string `json:"subtitlePending"`
	HeadingFiled    string `json:"headingFiled"`
	SubtitleFiled   string `json:"subtitleFiled"`
	HeadingPaid     string `json:"headingPaid"`
	SubtitlePaid    string `json:"subtitlePaid"`
}

type TabLabels struct {
	Pending string `json:"pending"`
	Filed   string `json:"filed"`
	Paid    string `json:"paid"`
}

type ColumnLabels struct {
	RemittanceType  string `json:"remittanceType"`
	Amount          string `json:"amount"`
	DueDate         string `json:"dueDate"`
	Status          string `json:"status"`
	FiledAt         string `json:"filedAt"`
	ReferenceNumber string `json:"referenceNumber"`
}

type TypeLabels struct {
	SSS            string `json:"sss"`
	PhilHealth     string `json:"philHealth"`
	PagIBIG        string `json:"pagIbig"`
	BIRWithholding string `json:"birWithholding"`
}

type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingPending:  "Pending Remittances",
			SubtitlePending: "Government contributions due for filing and payment",
			HeadingFiled:    "Filed Remittances",
			SubtitleFiled:   "Remittances filed with the government agency",
			HeadingPaid:     "Paid Remittances",
			SubtitlePaid:    "Remittances confirmed paid to the government agency",
		},
		Tabs: TabLabels{
			Pending: "Pending",
			Filed:   "Filed",
			Paid:    "Paid",
		},
		Columns: ColumnLabels{
			RemittanceType:  "Agency",
			Amount:          "Amount",
			DueDate:         "Due Date",
			Status:          "Status",
			FiledAt:         "Filed At",
			ReferenceNumber: "Reference #",
		},
		Types: TypeLabels{
			SSS:            "SSS",
			PhilHealth:     "PhilHealth",
			PagIBIG:        "Pag-IBIG",
			BIRWithholding: "BIR Withholding",
		},
		Empty: EmptyLabels{
			Title:   "No remittances found",
			Message: "Government contribution remittances will appear here once payroll runs are processed.",
		},
	}
}
