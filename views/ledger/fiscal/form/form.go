package form

import (
	fycha "github.com/erniealice/fycha-golang"
)

// Data is the template data for the fiscal period add drawer form.
type Data struct {
	FormAction   string
	Labels       fycha.FiscalPeriodFormLabels
	CommonLabels any
}
