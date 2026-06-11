// Package fund holds label structs for the Fund entity of the funding domain.
// Lyngua root key: "funding"
package fund

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
