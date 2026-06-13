package withholding_certificate

import "github.com/erniealice/pyeza-golang/compose"

func Describe() compose.Unit {
	r := DefaultRoutes()
	l := DefaultLabels()
	return compose.Unit{
		Key:       "treasury.withholding_certificate",
		Routes:    &r,
		RouteJSON: compose.JSONBinding{File: "route.json", Key: "withholding_certificate"},
		Labels:    &l,
		LabelJSON: compose.JSONBinding{File: "withholding_certificate.json", Key: ""},
		LabelName: "WithholdingCertificateLabels",
		Templates: TemplatesFS,
	}
}
