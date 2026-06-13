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
		Nav: compose.NavContrib{
			Permission: "withholding_certificate:list",
			Items: []compose.NavItem{
				{Key: "withholding-certs", Route: "withholding_certificate.list", Params: map[string]string{"status": "active"}, Label: "Withholding Certificates", Icon: "icon-file-text", Permission: "withholding_certificate:list"},
			},
		},
	}
}
