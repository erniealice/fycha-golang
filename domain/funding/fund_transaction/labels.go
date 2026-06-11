// Package fundtransaction holds label structs for the FundTransaction entity
// of the funding domain.
// Lyngua root key: "funding"
package fundtransaction

// FundingDrawFormLabels holds field/button labels for the draw (charge) drawer.
type FundingDrawFormLabels struct {
	Amount      string `json:"amount"`
	Description string `json:"description"`
	Submit      string `json:"submit"`
}

// FundingSettlementFormLabels holds field/button labels for the settlement drawer.
type FundingSettlementFormLabels struct {
	Amount string `json:"amount"`
	Submit string `json:"submit"`
}

// FundingTransferFormLabels holds field/button labels for the transfer drawer.
type FundingTransferFormLabels struct {
	DestinationFundID string `json:"destinationFundId"`
	Amount            string `json:"amount"`
	Submit            string `json:"submit"`
}
