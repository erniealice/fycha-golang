package charge_policy

import (
	enumspb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/enums"
)

// Option is one <select> option (value = canonical proto enum name).
type Option struct {
	Value    string
	Label    string
	Selected bool
	// Description feeds the pyeza select's per-option hint (optional).
	Description string
}

// AccountingRoleLabel renders the accounting role.
func (l Labels) AccountingRoleLabel(v enumspb.AccountingRole) string {
	switch v {
	case enumspb.AccountingRole_ACCOUNTING_ROLE_PRINCIPAL:
		return l.Enums.AccountingRolePrincipal
	case enumspb.AccountingRole_ACCOUNTING_ROLE_AGENT:
		return l.Enums.AccountingRoleAgent
	}
	return ""
}

// BookPresentationLabel renders the book presentation.
func (l Labels) BookPresentationLabel(v enumspb.BookPresentation) string {
	switch v {
	case enumspb.BookPresentation_BOOK_PRESENTATION_REVENUE:
		return l.Enums.BookPresentationRevenue
	case enumspb.BookPresentation_BOOK_PRESENTATION_EXCLUDED_REIMBURSEMENT:
		return l.Enums.BookPresentationExcluded
	}
	return ""
}

// TaxPositionLabel renders the tax position.
func (l Labels) TaxPositionLabel(v enumspb.TaxPosition) string {
	switch v {
	case enumspb.TaxPosition_TAX_POSITION_OWN_SUPPLY:
		return l.Enums.TaxPositionOwnSupply
	case enumspb.TaxPosition_TAX_POSITION_EXCLUDED_REIMBURSEMENT:
		return l.Enums.TaxPositionExcluded
	}
	return ""
}

// ComponentRoleLabel renders the component role.
func (l Labels) ComponentRoleLabel(v enumspb.ChargeComponentRole) string {
	switch v {
	case enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_OWN_REVENUE:
		return l.Enums.ComponentRoleOwnRevenue
	case enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_RECOVERY_COST:
		return l.Enums.ComponentRoleRecovery
	case enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_FEE:
		return l.Enums.ComponentRoleFee
	}
	return ""
}

// DocumentKindLabel renders the document kind.
func (l Labels) DocumentKindLabel(v enumspb.ChargeDocumentKind) string {
	switch v {
	case enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_INVOICE:
		return l.Enums.DocumentKindInvoice
	case enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT:
		return l.Enums.DocumentKindRecovery
	}
	return ""
}

// EventLabel renders the posting event.
func (l Labels) EventLabel(v enumspb.ChargePostingEvent) string {
	switch v {
	case enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ALLOCATION:
		return l.Enums.EventAllocation
	case enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ISSUE:
		return l.Enums.EventIssue
	case enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_APPLICATION:
		return l.Enums.EventApplication
	case enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_REVERSAL:
		return l.Enums.EventReversal
	}
	return ""
}

// PostingRoleLabel renders the posting role.
func (l Labels) PostingRoleLabel(v enumspb.ChargePostingRole) string {
	switch v {
	case enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_RECEIVABLE:
		return l.Enums.PostingRoleReceivable
	case enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_CLEARING:
		return l.Enums.PostingRoleClearing
	case enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_REVENUE:
		return l.Enums.PostingRoleRevenue
	case enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_REFUND_LIABILITY:
		return l.Enums.PostingRoleRefund
	case enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_EXPENSE:
		return l.Enums.PostingRoleExpense
	case enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_CASH:
		return l.Enums.PostingRoleCash
	}
	return ""
}

// VersionStatusLabel renders a version status; the variant is the badge variant.
func (l Labels) VersionStatusLabel(v enumspb.ChargePolicyVersionStatus) (label, variant string) {
	switch v {
	case enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_DRAFT:
		return l.Enums.VersionStatusDraft, "warning"
	case enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_APPROVED:
		return l.Enums.VersionStatusApproved, "success"
	case enumspb.ChargePolicyVersionStatus_CHARGE_POLICY_VERSION_STATUS_SUPERSEDED:
		return l.Enums.VersionStatusSuperseded, "muted"
	}
	return "", "default"
}

// Option builders (value = proto enum name; parsed back with the *_value maps).

func (l Labels) AccountingRoleOptions(selected enumspb.AccountingRole) []Option {
	vs := []enumspb.AccountingRole{enumspb.AccountingRole_ACCOUNTING_ROLE_PRINCIPAL, enumspb.AccountingRole_ACCOUNTING_ROLE_AGENT}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.AccountingRoleLabel(v), Selected: v == selected})
	}
	return out
}

func (l Labels) BookPresentationOptions(selected enumspb.BookPresentation) []Option {
	vs := []enumspb.BookPresentation{enumspb.BookPresentation_BOOK_PRESENTATION_REVENUE, enumspb.BookPresentation_BOOK_PRESENTATION_EXCLUDED_REIMBURSEMENT}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.BookPresentationLabel(v), Selected: v == selected})
	}
	return out
}

func (l Labels) TaxPositionOptions(selected enumspb.TaxPosition) []Option {
	vs := []enumspb.TaxPosition{enumspb.TaxPosition_TAX_POSITION_OWN_SUPPLY, enumspb.TaxPosition_TAX_POSITION_EXCLUDED_REIMBURSEMENT}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.TaxPositionLabel(v), Selected: v == selected})
	}
	return out
}

func (l Labels) ComponentRoleOptions(selected enumspb.ChargeComponentRole) []Option {
	vs := []enumspb.ChargeComponentRole{
		enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_OWN_REVENUE,
		enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_RECOVERY_COST,
		enumspb.ChargeComponentRole_CHARGE_COMPONENT_ROLE_FEE,
	}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.ComponentRoleLabel(v), Selected: v == selected})
	}
	return out
}

func (l Labels) DocumentKindOptions(selected enumspb.ChargeDocumentKind) []Option {
	vs := []enumspb.ChargeDocumentKind{
		enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_INVOICE,
		enumspb.ChargeDocumentKind_CHARGE_DOCUMENT_KIND_RECOVERY_DOCUMENT,
	}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.DocumentKindLabel(v), Selected: v == selected})
	}
	return out
}

func (l Labels) EventOptions(selected enumspb.ChargePostingEvent) []Option {
	vs := []enumspb.ChargePostingEvent{
		enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ALLOCATION,
		enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_ISSUE,
		enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_APPLICATION,
		enumspb.ChargePostingEvent_CHARGE_POSTING_EVENT_REVERSAL,
	}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.EventLabel(v), Selected: v == selected})
	}
	return out
}

func (l Labels) PostingRoleOptions(selected enumspb.ChargePostingRole) []Option {
	vs := []enumspb.ChargePostingRole{
		enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_RECEIVABLE,
		enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_CLEARING,
		enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_REVENUE,
		enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_REFUND_LIABILITY,
		enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_EXPENSE,
		enumspb.ChargePostingRole_CHARGE_POSTING_ROLE_CASH,
	}
	out := make([]Option, 0, len(vs))
	for _, v := range vs {
		out = append(out, Option{Value: v.String(), Label: l.PostingRoleLabel(v), Selected: v == selected})
	}
	return out
}
