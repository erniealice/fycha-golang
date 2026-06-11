package treasury

// ---------------------------------------------------------------------------
// WithholdingCertificateLabels
// Lyngua root key: "withholdingCertificate"
// ---------------------------------------------------------------------------

// WithholdingCertificateLabels holds all translatable strings for the
// Withholding Certificate CRUD views.
type WithholdingCertificateLabels struct {
	Page    WithholdingCertificatePageLabels   `json:"page"`
	Columns WithholdingCertificateColumnLabels `json:"columns"`
	Buttons WithholdingCertificateButtonLabels `json:"buttons"`
	Actions WithholdingCertificateActionLabels `json:"actions"`
	Empty   WithholdingCertificateEmptyLabels  `json:"empty"`
	Fields  WithholdingCertificateFieldLabels  `json:"fields"`
}

// WithholdingCertificatePageLabels holds page heading strings.
type WithholdingCertificatePageLabels struct {
	HeadingActive string `json:"headingActive"`
	CaptionActive string `json:"captionActive"`
	HeadingVoided string `json:"headingVoided"`
	CaptionVoided string `json:"captionVoided"`
}

// WithholdingCertificateColumnLabels holds table column headers.
type WithholdingCertificateColumnLabels struct {
	CertificateNumber  string `json:"certificateNumber"`
	RevenueID          string `json:"revenueId"`
	PeriodYear         string `json:"periodYear"`
	PeriodQuarter      string `json:"periodQuarter"`
	WhtAmountCertified string `json:"whtAmountCertified"`
	Status             string `json:"status"`
	DateIssued         string `json:"dateIssued"`
}

// WithholdingCertificateButtonLabels holds button text.
type WithholdingCertificateButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
	Void   string `json:"void"`
}

// WithholdingCertificateActionLabels holds action dropdown labels.
type WithholdingCertificateActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	NoPermission string `json:"noPermission"`
}

// WithholdingCertificateEmptyLabels holds empty-state strings.
type WithholdingCertificateEmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// WithholdingCertificateFieldLabels holds drawer form field labels.
type WithholdingCertificateFieldLabels struct {
	CertificateNumber  string `json:"certificateNumber"`
	RevenueID          string `json:"revenueId"`
	TaxAuthorityID     string `json:"taxAuthorityId"`
	PayorTin           string `json:"payorTin"`
	PayorName          string `json:"payorName"`
	PayeeTin           string `json:"payeeTin"`
	PayeeName          string `json:"payeeName"`
	PeriodYear         string `json:"periodYear"`
	PeriodQuarter      string `json:"periodQuarter"`
	WhtAmountCertified string `json:"whtAmountCertified"`
	Status             string `json:"status"`
	DateIssued         string `json:"dateIssued"`
	Notes              string `json:"notes"`
}

// DefaultWithholdingCertificateLabels returns WithholdingCertificateLabels with
// sensible English defaults.
func DefaultWithholdingCertificateLabels() WithholdingCertificateLabels {
	return WithholdingCertificateLabels{
		Page: WithholdingCertificatePageLabels{
			HeadingActive: "Withholding Certificates",
			CaptionActive: "BIR Form 2307 withholding tax certificates",
			HeadingVoided: "Voided Certificates",
			CaptionVoided: "Voided withholding tax certificates",
		},
		Columns: WithholdingCertificateColumnLabels{
			CertificateNumber:  "Certificate No.",
			RevenueID:          "Invoice",
			PeriodYear:         "Period Year",
			PeriodQuarter:      "Quarter",
			WhtAmountCertified: "WHT Certified",
			Status:             "Status",
			DateIssued:         "Date Issued",
		},
		Buttons: WithholdingCertificateButtonLabels{
			Add:    "Add Certificate",
			Edit:   "Edit",
			Delete: "Delete",
			Void:   "Void",
		},
		Actions: WithholdingCertificateActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			NoPermission: "You do not have permission to manage withholding certificates",
		},
		Empty: WithholdingCertificateEmptyLabels{
			Title:   "No withholding certificates",
			Message: "Add a withholding certificate received from a customer to record creditable WHT.",
		},
		Fields: WithholdingCertificateFieldLabels{
			CertificateNumber:  "Certificate Number",
			RevenueID:          "Invoice",
			TaxAuthorityID:     "Tax Authority",
			PayorTin:           "Payor TIN",
			PayorName:          "Payor Name",
			PayeeTin:           "Payee TIN",
			PayeeName:          "Payee Name",
			PeriodYear:         "Period Year",
			PeriodQuarter:      "Quarter",
			WhtAmountCertified: "WHT Amount Certified",
			Status:             "Status",
			DateIssued:         "Date Issued",
			Notes:              "Notes",
		},
	}
}
