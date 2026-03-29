package action

import (
	"context"
	"fmt"
	"log"
	"net/http"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	"github.com/erniealice/fycha-golang/seeder"
)

// FormData is the template data for the account drawer form.
type FormData struct {
	FormAction    string
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
	Labels        fycha.AccountFormLabels
	CommonLabels  any

	// Option lists for select elements (value/label pairs)
	ElementOptions  []SelectOption
	ClassOptions    []SelectOption
	CashFlowOptions []SelectOption
}

// SelectOption holds a select option value/label pair.
type SelectOption struct {
	Value    string
	Label    string
	Selected bool
}

// Deps holds dependencies for account action handlers.
type Deps struct {
	Routes fycha.AccountRoutes
	Labels fycha.AccountLabels

	// Account use cases
	CreateAccount func(ctx context.Context, req *accountpb.CreateAccountRequest) (*accountpb.CreateAccountResponse, error)
	ReadAccount   func(ctx context.Context, req *accountpb.ReadAccountRequest) (*accountpb.ReadAccountResponse, error)
	UpdateAccount func(ctx context.Context, req *accountpb.UpdateAccountRequest) (*accountpb.UpdateAccountResponse, error)
	DeleteAccount func(ctx context.Context, req *accountpb.DeleteAccountRequest) (*accountpb.DeleteAccountResponse, error)
}

// NewAddAction creates the account add action (GET = form, POST = create).
func NewAddAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("account", "create") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("account-drawer-form", &FormData{
				FormAction:      deps.Routes.AddURL,
				Active:          true,
				Labels:          deps.Labels.Form,
				CommonLabels:    nil, // injected by ViewAdapter
				ElementOptions:  elementOptions(deps.Labels.Form),
				ClassOptions:    classOptions("", deps.Labels.Form),
				CashFlowOptions: cashFlowOptions(deps.Labels.Form),
			})
		}

		// POST -- create account
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
		}

		if deps.CreateAccount == nil {
			log.Printf("CreateAccount use case not wired")
			return fycha.HTMXSuccess("accounts-tree-table")
		}

		element := parseElement(viewCtx.Request.FormValue("element"))
		classification := parseClassification(viewCtx.Request.FormValue("class"))
		cashFlow := parseCashFlow(viewCtx.Request.FormValue("cash_flow_class"))
		normalBal := parseNormalBalance(element)
		desc := viewCtx.Request.FormValue("description")
		active := viewCtx.Request.FormValue("active") == "true" || viewCtx.Request.FormValue("active") == "on"

		req := &accountpb.CreateAccountRequest{
			Data: &accountpb.Account{
				Code:             viewCtx.Request.FormValue("code"),
				Name:             viewCtx.Request.FormValue("name"),
				Element:          element,
				Classification:   classification,
				CashFlowActivity: cashFlow,
				NormalBalance:    normalBal,
				Active:           active,
			},
		}
		if desc != "" {
			req.Data.Description = &desc
		}

		resp, err := deps.CreateAccount(ctx, req)
		if err != nil {
			log.Printf("CreateAccount error: %v", err)
			return fycha.HTMXError("Failed to save account")
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := "Failed to save account"
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("accounts-tree-table")
	})
}

// NewEditAction creates the account edit action (GET = form, POST = update).
func NewEditAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("account", "update") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.PathValue("id")

		if viewCtx.Request.Method == http.MethodGet {
			// Load existing account to pre-populate the form
			formData := loadEditFormData(ctx, deps, id)
			return view.OK("account-drawer-form", formData)
		}

		// POST -- update account
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
		}

		if deps.UpdateAccount == nil {
			log.Printf("UpdateAccount use case not wired")
			return fycha.HTMXSuccess("accounts-tree-table")
		}

		element := parseElement(viewCtx.Request.FormValue("element"))
		classification := parseClassification(viewCtx.Request.FormValue("class"))
		cashFlow := parseCashFlow(viewCtx.Request.FormValue("cash_flow_class"))
		normalBal := parseNormalBalance(element)
		desc := viewCtx.Request.FormValue("description")
		active := viewCtx.Request.FormValue("active") == "true" || viewCtx.Request.FormValue("active") == "on"

		req := &accountpb.UpdateAccountRequest{
			Data: &accountpb.Account{
				Id:               id,
				Code:             viewCtx.Request.FormValue("code"),
				Name:             viewCtx.Request.FormValue("name"),
				Element:          element,
				Classification:   classification,
				CashFlowActivity: cashFlow,
				NormalBalance:    normalBal,
				Active:           active,
			},
		}
		if desc != "" {
			req.Data.Description = &desc
		}

		resp, err := deps.UpdateAccount(ctx, req)
		if err != nil {
			log.Printf("UpdateAccount error for %s: %v", id, err)
			return fycha.HTMXError("Failed to save account")
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := "Failed to save account"
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("accounts-tree-table")
	})
}

// NewDeleteAction creates the account delete action (POST only).
func NewDeleteAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("account", "delete") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.URL.Query().Get("id")
		if id == "" {
			return fycha.HTMXError("Account ID is required")
		}

		if deps.DeleteAccount == nil {
			log.Printf("DeleteAccount use case not wired")
			return fycha.HTMXSuccess("accounts-tree-table")
		}

		resp, err := deps.DeleteAccount(ctx, &accountpb.DeleteAccountRequest{
			Data: &accountpb.Account{Id: id},
		})
		if err != nil {
			log.Printf("DeleteAccount error for %s: %v", id, err)
			return fycha.HTMXError("Failed to delete account")
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := "Failed to delete account"
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("accounts-tree-table")
	})
}

// ---------------------------------------------------------------------------
// Edit form loader
// ---------------------------------------------------------------------------

// loadEditFormData fetches the account by ID and populates FormData for edit.
// Falls back to empty form if ReadAccount is nil or fails.
func loadEditFormData(ctx context.Context, deps *Deps, id string) *FormData {
	base := &FormData{
		FormAction:      deps.Routes.EditURL,
		IsEdit:          true,
		ID:              id,
		Active:          true,
		Labels:          deps.Labels.Form,
		CommonLabels:    nil,
		ElementOptions:  elementOptions(deps.Labels.Form),
		ClassOptions:    classOptions("", deps.Labels.Form),
		CashFlowOptions: cashFlowOptions(deps.Labels.Form),
	}

	if deps.ReadAccount == nil {
		return base
	}

	resp, err := deps.ReadAccount(ctx, &accountpb.ReadAccountRequest{
		Data: &accountpb.Account{Id: id},
	})
	if err != nil {
		log.Printf("ReadAccount error for edit form %s: %v", id, err)
		return base
	}
	if resp == nil || !resp.GetSuccess() || len(resp.GetData()) == 0 {
		return base
	}

	a := resp.GetData()[0]
	element := elementStringFromProto(a.GetElement())
	class := classStringFromProto(a.GetClassification())
	cashFlow := cashFlowStringFromProto(a.GetCashFlowActivity())

	return &FormData{
		FormAction:      deps.Routes.EditURL,
		IsEdit:          true,
		ID:              id,
		Code:            a.GetCode(),
		Name:            a.GetName(),
		Element:         element,
		Class:           class,
		Description:     a.GetDescription(),
		CashFlowClass:   cashFlow,
		Active:          a.GetActive(),
		Labels:          deps.Labels.Form,
		CommonLabels:    nil,
		ElementOptions:  elementOptions(deps.Labels.Form),
		ClassOptions:    classOptions(element, deps.Labels.Form),
		CashFlowOptions: cashFlowOptions(deps.Labels.Form),
	}
}

// ---------------------------------------------------------------------------
// Proto ↔ form string converters
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
// Option list helpers
// ---------------------------------------------------------------------------

func elementOptions(l fycha.AccountFormLabels) []SelectOption {
	return []SelectOption{
		{Value: "asset", Label: l.ElementAsset},
		{Value: "liability", Label: l.ElementLiability},
		{Value: "equity", Label: l.ElementEquity},
		{Value: "revenue", Label: l.ElementRevenue},
		{Value: "expense", Label: l.ElementExpense},
	}
}

func classOptions(element string, l fycha.AccountFormLabels) []SelectOption {
	switch element {
	case "asset":
		return []SelectOption{
			{Value: "current_asset", Label: l.ClassCurrentAsset},
			{Value: "non_current_asset", Label: l.ClassNonCurrentAsset},
		}
	case "liability":
		return []SelectOption{
			{Value: "current_liability", Label: l.ClassCurrentLiability},
			{Value: "non_current_liability", Label: l.ClassNonCurrentLiability},
		}
	case "equity":
		return []SelectOption{
			{Value: "equity", Label: l.ClassEquity},
		}
	case "revenue":
		return []SelectOption{
			{Value: "operating_revenue", Label: l.ClassOperatingRevenue},
			{Value: "other_income", Label: l.ClassOtherIncome},
		}
	case "expense":
		return []SelectOption{
			{Value: "cost_of_sales", Label: l.ClassCostOfSales},
			{Value: "operating_expense", Label: l.ClassOperatingExpense},
		}
	default:
		// Return all classes when element is not yet selected
		return []SelectOption{
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

func cashFlowOptions(l fycha.AccountFormLabels) []SelectOption {
	return []SelectOption{
		{Value: "", Label: l.CashFlowNone},
		{Value: "operating", Label: l.CashFlowOperating},
		{Value: "investing", Label: l.CashFlowInvesting},
		{Value: "financing", Label: l.CashFlowFinancing},
	}
}

// ---------------------------------------------------------------------------
// Account Template seeder action
// ---------------------------------------------------------------------------

// NewApplyTemplateAction creates the action handler that seeds default CoA accounts.
//
// POST /action/ledger/settings/account-templates/apply?template_id=service-ph
//
// Currently only "service-ph" is implemented (Philippine service business from
// seeder.DefaultCoA). Other template IDs return a "not yet available" error.
//
// The handler calls the CreateAccount use case for each account in the template.
// Accounts with duplicate codes are silently skipped (idempotent).
// Returns HTMX success trigger to refresh the page on success.
func NewApplyTemplateAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("account", "create") {
			return fycha.HTMXError("You do not have permission to apply account templates.")
		}

		templateID := viewCtx.Request.URL.Query().Get("template_id")
		if templateID == "" {
			templateID = "service-ph"
		}

		if templateID != "service-ph" {
			return fycha.HTMXError(
				fmt.Sprintf("Template %q is not yet available. Only 'service-ph' is currently supported.", templateID),
			)
		}

		if deps.CreateAccount == nil {
			// No use case wired — log and return success for dev/demo mode
			log.Printf("ApplyTemplate: CreateAccount use case not wired, skipping seeder")
			return fycha.HTMXSuccess("account-templates-content")
		}

		created, skipped, err := seeder.SeedDefaultCoA(ctx, deps.CreateAccount, "")
		if err != nil {
			// Partial success is non-fatal — the seeder collects errors for individual
			// accounts but continues. Log the aggregate error and return success.
			log.Printf("ApplyTemplate seeder completed with errors: %v", err)
		}

		log.Printf("ApplyTemplate: created=%d skipped=%d template=%s", created, skipped, templateID)
		return fycha.HTMXSuccess("account-templates-content")
	})
}
