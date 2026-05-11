// Package list provides the Withholding Certificate list view.
package list

import (
	"context"
	"fmt"
	"log"

	withholdingcertificatepb "github.com/erniealice/esqyma/pkg/schema/v1/domain/treasury/withholding_certificate"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/route"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the withholding certificate list page.
type Deps struct {
	Routes       fycha.WithholdingCertificateRoutes
	Labels       fycha.WithholdingCertificateLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Withholding certificate use cases
	ListWithholdingCertificates func(ctx context.Context, req *withholdingcertificatepb.ListWithholdingCertificatesRequest) (*withholdingcertificatepb.ListWithholdingCertificatesResponse, error)
}

// PageData holds the data for the withholding certificate list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	ActiveStatus    string
	StatusTabs      []pyeza.TabItem
	Table           *types.TableConfig
	Labels          fycha.WithholdingCertificateLabels
}

// WithholdingCertificateRow is the view-model for a single withholding certificate row.
type WithholdingCertificateRow struct {
	ID                 string
	CertificateNumber  string
	RevenueID          string
	PeriodYear         string
	PeriodQuarter      string
	WhtAmountCertified string
	Status             string
	DateIssued         string
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the withholding certificate list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		status := viewCtx.Request.PathValue("status")
		if status == "" {
			status = "active"
		}

		rows := fetchCertificates(ctx, deps, status)
		perms := view.GetUserPermissions(ctx)
		statusTabs := buildStatusTabs(deps)
		tableConfig := buildTableConfig(deps, rows, perms)

		heading, caption := headingForStatus(deps.Labels, status)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "withholding-certificates",
				HeaderTitle:    heading,
				HeaderSubtitle: caption,
				HeaderIcon:     "icon-file-text",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "withholding-certificate-list-content",
			ActiveStatus:    status,
			StatusTabs:      statusTabs,
			Table:           tableConfig,
			Labels:          deps.Labels,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "withholding_certificate"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("withholding-certificate-list", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

func fetchCertificates(ctx context.Context, deps *Deps, status string) []WithholdingCertificateRow {
	if deps.ListWithholdingCertificates == nil {
		return []WithholdingCertificateRow{}
	}

	resp, err := deps.ListWithholdingCertificates(ctx, &withholdingcertificatepb.ListWithholdingCertificatesRequest{})
	if err != nil {
		log.Printf("ListWithholdingCertificates error: %v", err)
		return []WithholdingCertificateRow{}
	}
	if resp == nil {
		return []WithholdingCertificateRow{}
	}

	rows := make([]WithholdingCertificateRow, 0)
	for _, wc := range resp.GetData() {
		row := protoToRow(wc)
		if status == "active" && row.Status != "active" {
			continue
		}
		if status == "voided" && row.Status != "voided" {
			continue
		}
		rows = append(rows, row)
	}
	return rows
}

func protoToRow(wc *withholdingcertificatepb.WithholdingCertificate) WithholdingCertificateRow {
	return WithholdingCertificateRow{
		ID:                 wc.GetId(),
		CertificateNumber:  wc.GetCertificateNumber(),
		RevenueID:          wc.GetRevenueId(),
		PeriodYear:         wc.GetCertificatePeriod(), // CertificatePeriod encodes period info
		PeriodQuarter:      "",                        // Derived from CertificatePeriod in Phase 5
		WhtAmountCertified: fmt.Sprintf("%.2f", float64(wc.GetActualAmount())/100.0),
		Status:             statusString(wc.GetStatus()),
		DateIssued:         wc.GetIssuedDate(),
	}
}

// ---------------------------------------------------------------------------
// Tab builder
// ---------------------------------------------------------------------------

func buildStatusTabs(deps *Deps) []pyeza.TabItem {
	base := deps.Routes.ListURL
	return []pyeza.TabItem{
		{Key: "active", Label: "Active", Href: route.ResolveURL(base, "status", "active"), Icon: ""},
		{Key: "voided", Label: "Voided", Href: route.ResolveURL(base, "status", "voided"), Icon: ""},
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, rows []WithholdingCertificateRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := []types.TableColumn{
		{Key: "certificate_number", Label: l.Columns.CertificateNumber, WidthClass: "col-3xl"},
		{Key: "revenue_id", Label: l.Columns.RevenueID, WidthClass: "col-3xl"},
		{Key: "period_year", Label: l.Columns.PeriodYear, WidthClass: "col-xl"},
		{Key: "period_quarter", Label: l.Columns.PeriodQuarter, WidthClass: "col-xl", Align: "center"},
		{Key: "wht_amount", Label: l.Columns.WhtAmountCertified, WidthClass: "col-3xl", Align: "right"},
		{Key: "date_issued", Label: l.Columns.DateIssued, WidthClass: "col-3xl"},
	}

	tableRows := []types.TableRow{}
	canView := perms.Can("withholding_certificate", "read")
	canCreate := perms.Can("withholding_certificate", "create")
	canDelete := perms.Can("withholding_certificate", "delete")
	for _, r := range rows {
		actions := []types.TableAction{
			{
				Type:            "view",
				Label:           l.Actions.View,
				Action:          "view",
				Href:            deps.Routes.DetailFor(r.ID),
				Disabled:        !canView,
				DisabledTooltip: l.Actions.NoPermission,
			},
			{
				Type:            "edit",
				Label:           l.Actions.Edit,
				Action:          "edit",
				Href:            deps.Routes.EditFor(r.ID),
				Disabled:        !canCreate,
				DisabledTooltip: l.Actions.NoPermission,
			},
			{
				Type:            "delete",
				Label:           l.Actions.Delete,
				Action:          "delete",
				Href:            deps.Routes.DeleteURL,
				Disabled:        !canDelete,
				DisabledTooltip: l.Actions.NoPermission,
			},
		}

		tableRows = append(tableRows, types.TableRow{
			ID:   r.ID,
			Href: deps.Routes.DetailFor(r.ID),
			Cells: []types.TableCell{
				{Type: "text", Value: r.CertificateNumber},
				{Type: "text", Value: r.RevenueID},
				{Type: "text", Value: r.PeriodYear},
				{Type: "text", Value: r.PeriodQuarter},
				{Type: "money", Value: r.WhtAmountCertified},
				{Type: "text", Value: r.DateIssued},
			},
			Actions: actions,
		})
	}

	types.ApplyColumnStyles(columns, tableRows)

	tableConfig := &types.TableConfig{
		ID:                "withholding-certs-table",
		Columns:           columns,
		Rows:              tableRows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowEntries:       true,
		DefaultSortColumn: "certificate_number",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.Add,
			ActionURL:       deps.Routes.AddURL,
			Icon:            "icon-plus",
			Disabled:        !canCreate,
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func headingForStatus(l fycha.WithholdingCertificateLabels, status string) (heading, caption string) {
	switch status {
	case "voided":
		return l.Page.HeadingVoided, l.Page.CaptionVoided
	default:
		return l.Page.HeadingActive, l.Page.CaptionActive
	}
}

func statusString(s withholdingcertificatepb.WithholdingCertificateStatus) string {
	switch s {
	case withholdingcertificatepb.WithholdingCertificateStatus_WITHHOLDING_CERTIFICATE_STATUS_RECEIVED:
		return "active"
	case withholdingcertificatepb.WithholdingCertificateStatus_WITHHOLDING_CERTIFICATE_STATUS_RECEIVED_WITH_VARIANCE:
		return "active"
	case withholdingcertificatepb.WithholdingCertificateStatus_WITHHOLDING_CERTIFICATE_STATUS_PENDING_RECEIPT:
		return "active"
	case withholdingcertificatepb.WithholdingCertificateStatus_WITHHOLDING_CERTIFICATE_STATUS_VOIDED:
		return "voided"
	case withholdingcertificatepb.WithholdingCertificateStatus_WITHHOLDING_CERTIFICATE_STATUS_REJECTED:
		return "voided"
	default:
		return "active"
	}
}
