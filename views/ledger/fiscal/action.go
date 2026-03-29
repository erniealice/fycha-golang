package fiscal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// Action form data
// ---------------------------------------------------------------------------

// AddFormData is the template data for the fiscal period add drawer form.
type AddFormData struct {
	FormAction   string
	Labels       fycha.FiscalPeriodFormLabels
	CommonLabels any
}

// ---------------------------------------------------------------------------
// Action deps
// ---------------------------------------------------------------------------

// ActionDeps holds dependencies for fiscal period action handlers.
type ActionDeps struct {
	Routes fycha.FiscalPeriodRoutes
	Labels fycha.FiscalPeriodLabels

	// Use cases (nil-safe — falls back to mock success when not wired)
	CreateFiscalPeriod func(ctx context.Context, req *fiscalperiodpb.CreateFiscalPeriodRequest) (*fiscalperiodpb.CreateFiscalPeriodResponse, error)
	CloseFiscalPeriod  func(ctx context.Context, req *fiscalperiodpb.CloseFiscalPeriodRequest) (*fiscalperiodpb.CloseFiscalPeriodResponse, error)
}

// ---------------------------------------------------------------------------
// Add action (GET = form, POST = create)
// ---------------------------------------------------------------------------

// NewAddAction creates the fiscal period add action.
func NewAddAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("fiscal_period", "create") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		if viewCtx.Request.Method == http.MethodGet {
			return view.OK("fiscal-period-drawer-form", &AddFormData{
				FormAction:   deps.Routes.AddURL,
				Labels:       deps.Labels.Form,
				CommonLabels: nil, // injected by ViewAdapter
			})
		}

		// POST — create fiscal period
		if err := viewCtx.Request.ParseForm(); err != nil {
			return fycha.HTMXError(deps.Labels.Actions.NoPermission)
		}

		if deps.CreateFiscalPeriod == nil {
			log.Printf("CreateFiscalPeriod use case not wired")
			return fycha.HTMXSuccess("fiscal-periods-table")
		}

		name := viewCtx.Request.FormValue("name")
		periodNumberStr := viewCtx.Request.FormValue("period_number")
		fiscalYearStr := viewCtx.Request.FormValue("fiscal_year")
		startDate := viewCtx.Request.FormValue("start_date")
		endDate := viewCtx.Request.FormValue("end_date")

		periodNumber, _ := strconv.ParseInt(periodNumberStr, 10, 32)
		fiscalYear, _ := strconv.ParseInt(fiscalYearStr, 10, 32)

		req := &fiscalperiodpb.CreateFiscalPeriodRequest{
			Data: &fiscalperiodpb.FiscalPeriod{
				Name:            name,
				PeriodNumber:    int32(periodNumber),
				FiscalYear:      int32(fiscalYear),
				StartDateString: &startDate,
				EndDateString:   &endDate,
				Status:          fiscalperiodpb.FiscalPeriodStatus_FISCAL_PERIOD_STATUS_OPEN,
			},
		}

		resp, err := deps.CreateFiscalPeriod(ctx, req)
		if err != nil {
			log.Printf("CreateFiscalPeriod error: %v", err)
			return fycha.HTMXError("Failed to create fiscal period")
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := "Failed to create fiscal period"
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("fiscal-periods-table")
	})
}

// ---------------------------------------------------------------------------
// Close action (POST only)
// ---------------------------------------------------------------------------

// NewCloseAction creates the fiscal period close action.
func NewCloseAction(deps *ActionDeps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		perms := view.GetUserPermissions(ctx)
		if !perms.Can("fiscal_period", "update") {
			return view.Error(fmt.Errorf("permission denied"))
		}

		id := viewCtx.Request.PathValue("id")
		if id == "" {
			return fycha.HTMXError("Fiscal period ID is required")
		}

		if deps.CloseFiscalPeriod == nil {
			log.Printf("CloseFiscalPeriod use case not wired")
			return fycha.HTMXSuccess("fiscal-periods-table")
		}

		// Extract the current user ID from context (set by session middleware as "uid")
		closedBy, _ := ctx.Value("uid").(string)

		resp, err := deps.CloseFiscalPeriod(ctx, &fiscalperiodpb.CloseFiscalPeriodRequest{
			FiscalPeriodId: id,
			ClosedBy:       closedBy,
		})
		if err != nil {
			log.Printf("CloseFiscalPeriod error for %s: %v", id, err)
			return fycha.HTMXError("Failed to close fiscal period")
		}
		if resp == nil || !resp.GetSuccess() {
			errMsg := "Failed to close fiscal period"
			if resp.GetError() != nil {
				errMsg = resp.GetError().GetMessage()
			}
			return fycha.HTMXError(errMsg)
		}

		return fycha.HTMXSuccess("fiscal-periods-table")
	})
}
