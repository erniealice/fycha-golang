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
	"github.com/erniealice/fycha-golang/views/ledger/form"
)

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
			// 2026-05-14 permission-gates P3: error-shape fix.
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("account-drawer-form", &form.Data{
				FormAction:      deps.Routes.AddURL,
				Active:          true,
				Labels:          deps.Labels.Form,
				CommonLabels:    nil, // injected by ViewAdapter
				ElementOptions:  form.ElementOptions(deps.Labels.Form),
				ClassOptions:    form.ClassOptions("", deps.Labels.Form),
				CashFlowOptions: form.CashFlowOptions(deps.Labels.Form),
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

		element := form.ParseElement(viewCtx.Request.FormValue("element"))
		classification := form.ParseClassification(viewCtx.Request.FormValue("class"))
		cashFlow := form.ParseCashFlow(viewCtx.Request.FormValue("cash_flow_class"))
		normalBal := form.ParseNormalBalance(element)
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
			// 2026-05-14 permission-gates P3: error-shape fix.
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
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

		element := form.ParseElement(viewCtx.Request.FormValue("element"))
		classification := form.ParseClassification(viewCtx.Request.FormValue("class"))
		cashFlow := form.ParseCashFlow(viewCtx.Request.FormValue("cash_flow_class"))
		normalBal := form.ParseNormalBalance(element)
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
			// 2026-05-14 permission-gates P3: error-shape fix.
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
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

// loadEditFormData fetches the account by ID and populates form.Data for edit.
// Falls back to empty form if ReadAccount is nil or fails.
func loadEditFormData(ctx context.Context, deps *Deps, id string) *form.Data {
	base := &form.Data{
		FormAction:      deps.Routes.EditURL,
		IsEdit:          true,
		ID:              id,
		Active:          true,
		Labels:          deps.Labels.Form,
		CommonLabels:    nil,
		ElementOptions:  form.ElementOptions(deps.Labels.Form),
		ClassOptions:    form.ClassOptions("", deps.Labels.Form),
		CashFlowOptions: form.CashFlowOptions(deps.Labels.Form),
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
	element := form.ElementStringFromProto(a.GetElement())
	class := form.ClassStringFromProto(a.GetClassification())
	cashFlow := form.CashFlowStringFromProto(a.GetCashFlowActivity())

	return &form.Data{
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
		ElementOptions:  form.ElementOptions(deps.Labels.Form),
		ClassOptions:    form.ClassOptions(element, deps.Labels.Form),
		CashFlowOptions: form.CashFlowOptions(deps.Labels.Form),
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
