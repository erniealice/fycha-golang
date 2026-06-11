package tax

// tax.go — tax-domain facade.
//
// Re-exports the tax_rate entity package's Routes/Labels with the entity
// prefix so consumers (block/) keep writing tax.DefaultTaxRateRoutes(),
// tax.DefaultTaxRateLabels(), etc.

import (
	taxrate "github.com/erniealice/fycha-golang/domain/tax/tax_rate"
)

// ---------------------------------------------------------------------------
// TaxRate (domain/tax/tax_rate)
// ---------------------------------------------------------------------------

type TaxRateRoutes = taxrate.Routes
type TaxRateLabels = taxrate.Labels

func DefaultTaxRateRoutes() TaxRateRoutes { return taxrate.DefaultRoutes() }
func DefaultTaxRateLabels() TaxRateLabels { return taxrate.DefaultLabels() }
