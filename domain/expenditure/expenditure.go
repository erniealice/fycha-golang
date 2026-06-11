package expenditure

// expenditure.go — expenditure-domain facade.
//
// Re-exports the prepayment entity package's Routes/Labels with the entity
// prefix. block builds the prepayment module via its own ModuleDeps; these
// aliases are provided for symmetry and any future consumer references.

import (
	prepayment "github.com/erniealice/fycha-golang/domain/expenditure/prepayment"
)

// ---------------------------------------------------------------------------
// Prepayment (domain/expenditure/prepayment)
// ---------------------------------------------------------------------------

type PrepaymentRoutes = prepayment.Routes
type PrepaymentLabels = prepayment.Labels

func DefaultPrepaymentRoutes() PrepaymentRoutes { return prepayment.DefaultRoutes() }
func DefaultPrepaymentLabels() PrepaymentLabels { return prepayment.DefaultLabels() }
