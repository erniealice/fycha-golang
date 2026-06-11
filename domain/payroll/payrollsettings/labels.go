// ---------------------------------------------------------------------------
// Payroll Settings labels
// ---------------------------------------------------------------------------

package payrollsettings

// Labels holds labels for Payroll Settings pages.
type Labels struct {
	GovRates   GovRatesLabels   `json:"govRates"`
	PayPeriods PayPeriodsLabels `json:"payPeriods"`
}

type GovRatesLabels struct {
	Page   GovRatesPageLabels   `json:"page"`
	Agency GovRatesAgencyLabels `json:"agency"`
}

type GovRatesPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

type GovRatesAgencyLabels struct {
	SSS            string `json:"sss"`
	PhilHealth     string `json:"philHealth"`
	PagIBIG        string `json:"pagIbig"`
	BIRWithholding string `json:"birWithholding"`
}

type PayPeriodsLabels struct {
	Page PayPeriodsPageLabels `json:"page"`
}

type PayPeriodsPageLabels struct {
	Heading string `json:"heading"`
	Caption string `json:"caption"`
}

// DefaultLabels returns Labels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files.
func DefaultLabels() Labels {
	return Labels{
		GovRates: GovRatesLabels{
			Page: GovRatesPageLabels{
				Heading: "Government Contribution Rates",
				Caption: "Philippine mandatory contribution rates — SSS, PhilHealth, Pag-IBIG, BIR",
			},
			Agency: GovRatesAgencyLabels{
				SSS:            "SSS (Social Security System)",
				PhilHealth:     "PhilHealth",
				PagIBIG:        "Pag-IBIG (HDMF)",
				BIRWithholding: "BIR Withholding Tax",
			},
		},
		PayPeriods: PayPeriodsLabels{
			Page: PayPeriodsPageLabels{
				Heading: "Pay Period Settings",
				Caption: "Configure payroll cut-off dates and pay schedules",
			},
		},
	}
}
