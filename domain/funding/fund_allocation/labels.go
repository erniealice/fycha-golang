// Package fundallocation holds label structs for the FundAllocation entity
// of the funding domain.
// Lyngua root key: "funding"
package fundallocation

// FundingAllocationFormLabels holds field/button labels for the allocation drawer.
type FundingAllocationFormLabels struct {
	AllocatedLimit string `json:"allocatedLimit"`
	Mode           string `json:"mode"`
	ModeHardLimit  string `json:"modeHardLimit"`
	ModeSoftLimit  string `json:"modeSoftLimit"`
}
