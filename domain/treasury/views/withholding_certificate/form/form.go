// Package form provides the data shape and label builder for the
// Withholding Certificate drawer form.
package form

import (
	treasury "github.com/erniealice/fycha-golang/domain/treasury"
)

// Labels holds the translatable strings used by the drawer template.
type Labels struct {
	DrawerTitleAdd     string
	DrawerTitleEdit    string
	CertificateNumber  string
	RevenueID          string
	TaxAuthorityID     string
	PayorTin           string
	PayorName          string
	PayeeTin           string
	PayeeName          string
	PeriodYear         string
	PeriodQuarter      string
	WhtAmountCertified string
	Status             string
	DateIssued         string
	Notes              string
}

// BuildLabels converts WithholdingCertificateLabels to the flat Labels struct.
func BuildLabels(l treasury.WithholdingCertificateLabels) Labels {
	return Labels{
		DrawerTitleAdd:     l.Buttons.Add,
		DrawerTitleEdit:    l.Buttons.Edit,
		CertificateNumber:  l.Fields.CertificateNumber,
		RevenueID:          l.Fields.RevenueID,
		TaxAuthorityID:     l.Fields.TaxAuthorityID,
		PayorTin:           l.Fields.PayorTin,
		PayorName:          l.Fields.PayorName,
		PayeeTin:           l.Fields.PayeeTin,
		PayeeName:          l.Fields.PayeeName,
		PeriodYear:         l.Fields.PeriodYear,
		PeriodQuarter:      l.Fields.PeriodQuarter,
		WhtAmountCertified: l.Fields.WhtAmountCertified,
		Status:             l.Fields.Status,
		DateIssued:         l.Fields.DateIssued,
		Notes:              l.Fields.Notes,
	}
}

// Data is the template data for the withholding certificate drawer form.
type Data struct {
	FormAction string
	WorkspaceID string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit     bool
	ID         string
	// Pre-populated fields
	RevenueID         string
	CertificateNumber string
	TaxAuthorityID    string
	PayorTin          string
	PayorName         string
	PayeeTin          string
	PayeeName         string
	PeriodYear        string
	PeriodQuarter     string
	// WhtAmountCertified is the centavo amount expressed as a display string
	// (e.g. "1500.00"). Converted to int64 centavos on POST.
	WhtAmountCertified string
	// ExpectedAmount is pre-filled from revenue.wht_amount_expected when launched
	// from the "Add Withholding Certificate" CTA on the revenue detail drawer.
	ExpectedAmount string
	StatusOptions  []StatusOption
	DateIssued     string
	Notes          string

	Labels       Labels
	CommonLabels any
}

// StatusOption is a single entry in the status <select>.
type StatusOption struct {
	Value    string
	Label    string
	Selected bool
	// Description must be present (even if empty) because form-group.html's
	// select branch reads data-description="{{.Description}}" on every option.
	Description string
}

// BuildStatusOptions constructs the status options for the select element.
func BuildStatusOptions(current string) []StatusOption {
	opts := []struct{ v, l string }{
		{"WITHHOLDING_CERTIFICATE_STATUS_PENDING_RECEIPT", "Pending Receipt"},
		{"WITHHOLDING_CERTIFICATE_STATUS_RECEIVED", "Received"},
		{"WITHHOLDING_CERTIFICATE_STATUS_RECEIVED_WITH_VARIANCE", "Received (Variance)"},
		{"WITHHOLDING_CERTIFICATE_STATUS_VOIDED", "Voided"},
		{"WITHHOLDING_CERTIFICATE_STATUS_REJECTED", "Rejected"},
	}
	out := make([]StatusOption, 0, len(opts))
	for _, o := range opts {
		out = append(out, StatusOption{
			Value:    o.v,
			Label:    o.l,
			Selected: o.v == current,
		})
	}
	return out
}
