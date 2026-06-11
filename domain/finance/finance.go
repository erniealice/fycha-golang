package finance

// finance.go — finance-domain facade.
//
// Re-exports the forex_rate entity package's Routes/Labels with the entity
// prefix so consumers (block/) keep writing finance.DefaultForexRateRoutes(),
// finance.DefaultForexRateLabels(), etc.

import (
	forexrate "github.com/erniealice/fycha-golang/domain/finance/forex_rate"
)

// ---------------------------------------------------------------------------
// ForexRate (domain/finance/forex_rate)
// ---------------------------------------------------------------------------

type ForexRateRoutes = forexrate.Routes
type ForexRateLabels = forexrate.Labels

func DefaultForexRateRoutes() ForexRateRoutes { return forexrate.DefaultRoutes() }
func DefaultForexRateLabels() ForexRateLabels { return forexrate.DefaultLabels() }
