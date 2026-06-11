// Package withholding_certificate provides labels, routes, and views for the
// Withholding Certificate entity.
// Lyngua root key: "withholdingCertificate"
package withholding_certificate

// ---------------------------------------------------------------------------
// Labels
// ---------------------------------------------------------------------------

// Labels holds all translatable strings for the Withholding Certificate CRUD views.
type Labels struct {
	Page    PageLabels   `json:"page"`
	Columns ColumnLabels `json:"columns"`
	Buttons ButtonLabels `json:"buttons"`
	Actions ActionLabels `json:"actions"`
	Empty   EmptyLabels  `json:"empty"`
	Fields  FieldLabels  `json:"fields"`
}

// PageLabels holds page heading strings.
type PageLabels struct {
	HeadingActive string `json:"headingActive"`
	CaptionActive string `json:"captionActive"`
	HeadingVoided string `json:"headingVoided"`
	CaptionVoided string `json:"captionVoided"`
}

// ColumnLabels holds table column headers.
type ColumnLabels struct {
	CertificateNumber  string `json:"certificateNumber"`
	RevenueID          string `json:"revenueId"`
	PeriodYear         string `json:"periodYear"`
	PeriodQuarter      string `json:"periodQuarter"`
	WhtAmountCertified string `json:"whtAmountCertified"`
	Status             string `json:"status"`
	DateIssued         string `json:"dateIssued"`
}

// ButtonLabels holds button text.
type ButtonLabels struct {
	Add    string `json:"add"`
	Edit   string `json:"edit"`
	Delete string `json:"delete"`
	Void   string `json:"void"`
}

// ActionLabels holds action dropdown labels.
type ActionLabels struct {
	View         string `json:"view"`
	Edit         string `json:"edit"`
	Delete       string `json:"delete"`
	NoPermission string `json:"noPermission"`
}

// EmptyLabels holds empty-state strings.
type EmptyLabels struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

// FieldLabels holds drawer form field labels.
type FieldLabels struct {
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

// DefaultLabels returns Labels with sensible English defaults.
func DefaultLabels() Labels {
	return Labels{
		Page: PageLabels{
			HeadingActive: "Withholding Certificates",
			CaptionActive: "BIR Form 2307 withholding tax certificates",
			HeadingVoided: "Voided Certificates",
			CaptionVoided: "Voided withholding tax certificates",
		},
		Columns: ColumnLabels{
			CertificateNumber:  "Certificate No.",
			RevenueID:          "Invoice",
			PeriodYear:         "Period Year",
			PeriodQuarter:      "Quarter",
			WhtAmountCertified: "WHT Certified",
			Status:             "Status",
			DateIssued:         "Date Issued",
		},
		Buttons: ButtonLabels{
			Add:    "Add Certificate",
			Edit:   "Edit",
			Delete: "Delete",
			Void:   "Void",
		},
		Actions: ActionLabels{
			View:         "View",
			Edit:         "Edit",
			Delete:       "Delete",
			NoPermission: "You do not have permission to manage withholding certificates",
		},
		Empty: EmptyLabels{
			Title:   "No withholding certificates",
			Message: "Add a withholding certificate received from a customer to record creditable WHT.",
		},
		Fields: FieldLabels{
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
