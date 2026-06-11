// Package funding holds label structs and route constants for the
// esqyma funding domain: fund, fund_allocation, fund_transaction.
package funding

// ---------------------------------------------------------------------------
// FundingFormLabels (Funding > drawer forms)
// Lyngua root key: "funding"
// ---------------------------------------------------------------------------

// FundingFormLabels holds all translatable strings for the four funding
// drawer forms: allocation, draw (charge), settlement, and transfer.
// The lyngua root key is "funding"; JSON lives in
// packages/lyngua/translations/en/{general,professional}/funding.json.
type FundingFormLabels struct {
	Allocation FundingAllocationFormLabels `json:"allocation"`
	Draw       FundingDrawFormLabels       `json:"draw"`
	Settlement FundingSettlementFormLabels `json:"settlement"`
	Transfer   FundingTransferFormLabels   `json:"transfer"`
	Source     FundingSourceListLabels     `json:"source"`
}

// FundingAllocationFormLabels holds field/button labels for the allocation drawer.
type FundingAllocationFormLabels struct {
	AllocatedLimit string `json:"allocatedLimit"`
	Mode           string `json:"mode"`
	ModeHardLimit  string `json:"modeHardLimit"`
	ModeSoftLimit  string `json:"modeSoftLimit"`
}

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

// FundingSourceListLabels holds page-level strings for the fund source list view.
type FundingSourceListLabels struct {
	Title    string                    `json:"title"`
	Subtitle string                    `json:"subtitle"`
	Kind     FundingSourceKindLabels   `json:"kind"`
	Status   FundingSourceStatusLabels `json:"status"`
}

// FundingSourceKindLabels maps FundKind enum values to display strings for the
// fund source list. Keys mirror the proto FundKind enum.
type FundingSourceKindLabels struct {
	CashOnHand  string `json:"cashOnHand"`
	BankAccount string `json:"bankAccount"`
	PettyCash   string `json:"pettyCash"`
	CreditCard  string `json:"creditCard"`
	CreditLine  string `json:"creditLine"`
	PrepaidCard string `json:"prepaidCard"`
	MobileMoney string `json:"mobileMoney"`
	Unknown     string `json:"unknown"`
}

// FundingSourceStatusLabels maps FundStatus enum values to display strings for
// the fund source list. Keys mirror the proto FundStatus enum.
type FundingSourceStatusLabels struct {
	Draft     string `json:"draft"`
	Active    string `json:"active"`
	Suspended string `json:"suspended"`
	Archived  string `json:"archived"`
	Unknown   string `json:"unknown"`
}

// DefaultFundingFormLabels returns FundingFormLabels with hardcoded English defaults.
// Consumer apps should override these via lyngua JSON files
// (packages/lyngua/translations/en/{general,professional}/funding.json).
func DefaultFundingFormLabels() FundingFormLabels {
	return FundingFormLabels{
		Allocation: FundingAllocationFormLabels{
			AllocatedLimit: "Allocated Limit",
			Mode:           "Mode",
			ModeHardLimit:  "Hard Limit",
			ModeSoftLimit:  "Soft Limit",
		},
		Draw: FundingDrawFormLabels{
			Amount:      "Amount",
			Description: "Description",
			Submit:      "Charge",
		},
		Settlement: FundingSettlementFormLabels{
			Amount: "Settlement Amount",
			Submit: "Settle",
		},
		Transfer: FundingTransferFormLabels{
			DestinationFundID: "Destination Fund ID",
			Amount:            "Amount",
			Submit:            "Transfer",
		},
		Source: FundingSourceListLabels{
			Title:    "Fund Sources",
			Subtitle: "Funds you own and share with workspaces",
			Kind: FundingSourceKindLabels{
				CashOnHand:  "Cash on Hand",
				BankAccount: "Bank Account",
				PettyCash:   "Petty Cash",
				CreditCard:  "Credit Card",
				CreditLine:  "Credit Line",
				PrepaidCard: "Prepaid Card",
				MobileMoney: "Mobile Money",
				Unknown:     "Unknown",
			},
			Status: FundingSourceStatusLabels{
				Draft:     "Draft",
				Active:    "Active",
				Suspended: "Suspended",
				Archived:  "Archived",
				Unknown:   "Unknown",
			},
		},
	}
}
