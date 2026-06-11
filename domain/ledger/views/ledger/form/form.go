package form

import (
	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	pyeza "github.com/erniealice/pyeza-golang/types"

	ledger "github.com/erniealice/fycha-golang/domain/ledger"
)

// Data is the template data for the account drawer form.
type Data struct {
	FormAction    string
	WorkspaceID    string // injected by C1: populated by ViewAdapter.injectWorkspaceID for action_workspace_guard
	IsEdit        bool
	ID            string
	Code          string
	Name          string
	Element       string
	Class         string
	ParentCode    string
	Group         string
	IsGroup       bool
	Active        bool
	Description   string
	CashFlowClass string
	Labels        ledger.AccountFormLabels
	CommonLabels  any

	// Option lists for select elements (value/label pairs)
	ElementOptions  []pyeza.SelectOption
	ClassOptions    []pyeza.SelectOption
	CashFlowOptions []pyeza.SelectOption
}

// ---------------------------------------------------------------------------
// Proto ↔ form string converters (pure helpers, no deps)
// ---------------------------------------------------------------------------

func parseElement(s string) accountpb.AccountElement {
	switch s {
	case "asset":
		return accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET
	case "liability":
		return accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY
	case "equity":
		return accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY
	case "revenue":
		return accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE
	case "expense":
		return accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE
	default:
		return accountpb.AccountElement_ACCOUNT_ELEMENT_UNSPECIFIED
	}
}

func parseClassification(s string) accountpb.AccountClassification {
	switch s {
	case "current_asset":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET
	case "non_current_asset":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET
	case "current_liability":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY
	case "non_current_liability":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_LIABILITY
	case "equity":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_EQUITY
	case "operating_revenue":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE
	case "other_income":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OTHER_INCOME
	case "cost_of_sales":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_COST_OF_SALES
	case "operating_expense":
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE
	default:
		return accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_UNSPECIFIED
	}
}

func parseCashFlow(s string) accountpb.CashFlowActivity {
	switch s {
	case "operating":
		return accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_OPERATING
	case "investing":
		return accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_INVESTING
	case "financing":
		return accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_FINANCING
	case "":
		return accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_NONE
	default:
		return accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_UNSPECIFIED
	}
}

// parseNormalBalance derives normal balance from the element (accounting rule).
func parseNormalBalance(e accountpb.AccountElement) accountpb.NormalBalance {
	switch e {
	case accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET,
		accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE:
		return accountpb.NormalBalance_NORMAL_BALANCE_DEBIT
	case accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY,
		accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY,
		accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE:
		return accountpb.NormalBalance_NORMAL_BALANCE_CREDIT
	default:
		return accountpb.NormalBalance_NORMAL_BALANCE_UNSPECIFIED
	}
}

func elementStringFromProto(e accountpb.AccountElement) string {
	switch e {
	case accountpb.AccountElement_ACCOUNT_ELEMENT_ASSET:
		return "asset"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_LIABILITY:
		return "liability"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_EQUITY:
		return "equity"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_REVENUE:
		return "revenue"
	case accountpb.AccountElement_ACCOUNT_ELEMENT_EXPENSE:
		return "expense"
	default:
		return ""
	}
}

func classStringFromProto(c accountpb.AccountClassification) string {
	switch c {
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_ASSET:
		return "current_asset"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET:
		return "non_current_asset"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY:
		return "current_liability"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_NON_CURRENT_LIABILITY:
		return "non_current_liability"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_EQUITY:
		return "equity"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_REVENUE:
		return "operating_revenue"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OTHER_INCOME:
		return "other_income"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_COST_OF_SALES:
		return "cost_of_sales"
	case accountpb.AccountClassification_ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE:
		return "operating_expense"
	default:
		return ""
	}
}

func cashFlowStringFromProto(c accountpb.CashFlowActivity) string {
	switch c {
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_OPERATING:
		return "operating"
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_INVESTING:
		return "investing"
	case accountpb.CashFlowActivity_CASH_FLOW_ACTIVITY_FINANCING:
		return "financing"
	default:
		return ""
	}
}

// ---------------------------------------------------------------------------
// Option list helpers (pure functions, no deps)
// ---------------------------------------------------------------------------

func ElementOptions(l ledger.AccountFormLabels) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: "asset", Label: l.ElementAsset},
		{Value: "liability", Label: l.ElementLiability},
		{Value: "equity", Label: l.ElementEquity},
		{Value: "revenue", Label: l.ElementRevenue},
		{Value: "expense", Label: l.ElementExpense},
	}
}

func ClassOptions(element string, l ledger.AccountFormLabels) []pyeza.SelectOption {
	switch element {
	case "asset":
		return []pyeza.SelectOption{
			{Value: "current_asset", Label: l.ClassCurrentAsset},
			{Value: "non_current_asset", Label: l.ClassNonCurrentAsset},
		}
	case "liability":
		return []pyeza.SelectOption{
			{Value: "current_liability", Label: l.ClassCurrentLiability},
			{Value: "non_current_liability", Label: l.ClassNonCurrentLiability},
		}
	case "equity":
		return []pyeza.SelectOption{
			{Value: "equity", Label: l.ClassEquity},
		}
	case "revenue":
		return []pyeza.SelectOption{
			{Value: "operating_revenue", Label: l.ClassOperatingRevenue},
			{Value: "other_income", Label: l.ClassOtherIncome},
		}
	case "expense":
		return []pyeza.SelectOption{
			{Value: "cost_of_sales", Label: l.ClassCostOfSales},
			{Value: "operating_expense", Label: l.ClassOperatingExpense},
		}
	default:
		// Return all classes when element is not yet selected
		return []pyeza.SelectOption{
			{Value: "current_asset", Label: l.ClassCurrentAsset},
			{Value: "non_current_asset", Label: l.ClassNonCurrentAsset},
			{Value: "current_liability", Label: l.ClassCurrentLiability},
			{Value: "non_current_liability", Label: l.ClassNonCurrentLiability},
			{Value: "equity", Label: l.ClassEquity},
			{Value: "operating_revenue", Label: l.ClassOperatingRevenue},
			{Value: "other_income", Label: l.ClassOtherIncome},
			{Value: "cost_of_sales", Label: l.ClassCostOfSales},
			{Value: "operating_expense", Label: l.ClassOperatingExpense},
		}
	}
}

func CashFlowOptions(l ledger.AccountFormLabels) []pyeza.SelectOption {
	return []pyeza.SelectOption{
		{Value: "", Label: l.CashFlowNone},
		{Value: "operating", Label: l.CashFlowOperating},
		{Value: "investing", Label: l.CashFlowInvesting},
		{Value: "financing", Label: l.CashFlowFinancing},
	}
}

// Exported parsing helpers used by action/ (keeping them lowercase but accessible).
// They live here but are called via form package by action/action.go.
// ParseElement converts string to proto AccountElement enum.
func ParseElement(s string) accountpb.AccountElement {
	return parseElement(s)
}

// ParseClassification converts string to proto AccountClassification enum.
func ParseClassification(s string) accountpb.AccountClassification {
	return parseClassification(s)
}

// ParseCashFlow converts string to proto CashFlowActivity enum.
func ParseCashFlow(s string) accountpb.CashFlowActivity {
	return parseCashFlow(s)
}

// ParseNormalBalance derives normal balance from AccountElement.
func ParseNormalBalance(e accountpb.AccountElement) accountpb.NormalBalance {
	return parseNormalBalance(e)
}

// ElementStringFromProto converts proto AccountElement enum to string.
func ElementStringFromProto(e accountpb.AccountElement) string {
	return elementStringFromProto(e)
}

// ClassStringFromProto converts proto AccountClassification enum to string.
func ClassStringFromProto(c accountpb.AccountClassification) string {
	return classStringFromProto(c)
}

// CashFlowStringFromProto converts proto CashFlowActivity enum to string.
func CashFlowStringFromProto(c accountpb.CashFlowActivity) string {
	return cashFlowStringFromProto(c)
}
