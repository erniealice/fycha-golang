// Package funding holds the funding-domain facade: re-exports the funding
// drawer-form label types (owned by the fundinglabels leaf package) so that
// consumers keep writing funding.FundingFormLabels,
// funding.DefaultFundingFormLabels(), etc.
//
// The concrete definitions live in domain/funding/funding/labels
// (package fundinglabels) to keep the entity DAG acyclic: the facade and the
// funding leaf views both import the leaf, never each other.
package funding

import (
	fundinglabels "github.com/erniealice/fycha-golang/domain/funding/funding/labels"
)

type FundingFormLabels = fundinglabels.FundingFormLabels
type FundingAllocationFormLabels = fundinglabels.FundingAllocationFormLabels
type FundingDrawFormLabels = fundinglabels.FundingDrawFormLabels
type FundingSettlementFormLabels = fundinglabels.FundingSettlementFormLabels
type FundingTransferFormLabels = fundinglabels.FundingTransferFormLabels
type FundingSourceListLabels = fundinglabels.FundingSourceListLabels
type FundingSourceKindLabels = fundinglabels.FundingSourceKindLabels
type FundingSourceStatusLabels = fundinglabels.FundingSourceStatusLabels

// DefaultFundingFormLabels returns FundingFormLabels with hardcoded English defaults.
func DefaultFundingFormLabels() FundingFormLabels { return fundinglabels.DefaultFundingFormLabels() }
