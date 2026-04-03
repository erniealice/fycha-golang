package recurring

import (
	"context"
	"log"

	recurringpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/recurring_journal_template"
	lynguaV1 "github.com/erniealice/lyngua/golang/v1"
	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
)

// ---------------------------------------------------------------------------
// View dependencies + page data
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the recurring journal template list page.
type Deps struct {
	Routes       fycha.LedgerSettingsRoutes
	Labels       fycha.RecurringTemplateLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Recurring template use cases
	GetRecurringTemplateList func(ctx context.Context) ([]*recurringpb.RecurringJournalTemplate, error)
}

// PageData holds the data for the recurring templates list page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Table           *types.TableConfig
}

// RecurringTemplateRow is the view-model for a single recurring template row.
type RecurringTemplateRow struct {
	ID          string
	Name        string
	Description string
	Frequency   string // "daily", "weekly", "monthly", "quarterly", "yearly"
	NextRunDate string
	Active      bool
}

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the recurring templates list view (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		templates := fetchTemplates(ctx, deps)

		perms := view.GetUserPermissions(ctx)
		tableConfig := buildTableConfig(deps, templates, perms)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          deps.Labels.Page.Heading,
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "recurring-templates",
				HeaderTitle:    deps.Labels.Page.Heading,
				HeaderSubtitle: deps.Labels.Page.Caption,
				HeaderIcon:     "icon-repeat",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate: "recurring-templates-content",
			Table:           tableConfig,
		}

		// KB help content
		if viewCtx.Translations != nil {
			if provider, ok := viewCtx.Translations.(*lynguaV1.TranslationProvider); ok {
				if kb, _ := provider.LoadKBIfExists(viewCtx.Lang, viewCtx.BusinessType, "ledger-recurring"); kb != nil {
					pageData.HasHelp = true
					pageData.HelpContent = kb.Body
				}
			}
		}

		return view.OK("recurring-templates", pageData)
	})
}

// ---------------------------------------------------------------------------
// Data fetcher
// ---------------------------------------------------------------------------

// fetchTemplates calls the use case and converts the response to view-model rows.
// Falls back to mock data when no use case is wired (placeholder mode).
func fetchTemplates(ctx context.Context, deps *Deps) []RecurringTemplateRow {
	if deps.GetRecurringTemplateList == nil {
		return mockTemplates()
	}

	templates, err := deps.GetRecurringTemplateList(ctx)
	if err != nil {
		log.Printf("GetRecurringTemplateList error: %v", err)
		return mockTemplates()
	}

	rows := make([]RecurringTemplateRow, 0, len(templates))
	for _, t := range templates {
		rows = append(rows, protoToRow(t))
	}
	return rows
}

// protoToRow converts a proto RecurringJournalTemplate to a view-model row.
func protoToRow(t *recurringpb.RecurringJournalTemplate) RecurringTemplateRow {
	return RecurringTemplateRow{
		ID:          t.GetId(),
		Name:        t.GetName(),
		Description: t.GetDescription(),
		Frequency:   frequencyString(t.GetFrequency()),
		NextRunDate: t.GetNextRunDateString(),
		Active:      t.GetActive(),
	}
}

// ---------------------------------------------------------------------------
// Table builder
// ---------------------------------------------------------------------------

func buildTableConfig(deps *Deps, templates []RecurringTemplateRow, perms *types.UserPermissions) *types.TableConfig {
	l := deps.Labels
	columns := templateColumns(l)
	rows := buildTableRows(templates, l, perms)
	types.ApplyColumnStyles(columns, rows)

	tableConfig := &types.TableConfig{
		ID:                "recurring-templates-table",
		Columns:           columns,
		Rows:              rows,
		ShowSearch:        true,
		ShowActions:       true,
		ShowExport:        false,
		ShowEntries:       true,
		DefaultSortColumn: "name",
		Labels:            deps.TableLabels,
		EmptyState: types.TableEmptyState{
			Title:   l.Empty.Title,
			Message: l.Empty.Message,
		},
		PrimaryAction: &types.PrimaryAction{
			Label:           l.Buttons.AddTemplate,
			ActionURL:       "#",
			Icon:            "icon-plus",
			Disabled:        !perms.Can("recurring_template", "create"),
			DisabledTooltip: l.Actions.NoPermission,
		},
	}
	types.ApplyTableSettings(tableConfig)

	return tableConfig
}

func templateColumns(l fycha.RecurringTemplateLabels) []types.TableColumn {
	return []types.TableColumn{
		{Key: "name", Label: l.Columns.Name, Sortable: false},
		{Key: "frequency", Label: l.Columns.Frequency, Sortable: false, Width: "120px"},
		{Key: "next_run", Label: l.Columns.NextRun, Sortable: false, Width: "120px"},
		{Key: "status", Label: l.Columns.Status, Sortable: false, Width: "110px"},
	}
}

func buildTableRows(templates []RecurringTemplateRow, l fycha.RecurringTemplateLabels, perms *types.UserPermissions) []types.TableRow {
	rows := []types.TableRow{}
	for _, t := range templates {
		canUpdate := perms.Can("recurring_template", "update")

		actions := []types.TableAction{
			{
				Type:            "edit",
				Label:           l.Actions.Edit,
				Action:          "edit",
				URL:             "#",
				DrawerTitle:     l.Actions.Edit,
				Disabled:        !canUpdate,
				DisabledTooltip: l.Actions.NoPermission,
			},
		}

		freqVariant := frequencyBadgeVariant(t.Frequency)
		freqLabel := frequencyDisplayLabel(l, t.Frequency)

		statusVariant := "success"
		statusLabel := l.Status.Active
		if !t.Active {
			statusVariant = "muted"
			statusLabel = l.Status.Inactive
		}

		row := types.TableRow{
			ID: t.ID,
			Cells: []types.TableCell{
				{Type: "text", Value: t.Name},
				{Type: "badge", Value: freqLabel, Variant: freqVariant},
				{Type: "text", Value: t.NextRunDate},
				{Type: "badge", Value: statusLabel, Variant: statusVariant},
			},
			Actions: actions,
		}
		rows = append(rows, row)
	}
	return rows
}

// ---------------------------------------------------------------------------
// Mock data (placeholder until DB is wired)
// ---------------------------------------------------------------------------

func mockTemplates() []RecurringTemplateRow {
	return []RecurringTemplateRow{
		{ID: "rt-01", Name: "Monthly Depreciation", Description: "Depreciation for IT equipment and furniture", Frequency: "monthly", NextRunDate: "2026-04-01", Active: true},
		{ID: "rt-02", Name: "Rent Accrual", Description: "Monthly office rent accrual", Frequency: "monthly", NextRunDate: "2026-04-01", Active: true},
		{ID: "rt-03", Name: "Insurance Amortization", Description: "Annual insurance policy spread monthly", Frequency: "monthly", NextRunDate: "2026-04-01", Active: true},
		{ID: "rt-04", Name: "Quarterly Tax Provision", Description: "Estimated income tax provision", Frequency: "quarterly", NextRunDate: "2026-06-30", Active: true},
	}
}

// ---------------------------------------------------------------------------
// Proto enum → display string converters
// ---------------------------------------------------------------------------

func frequencyString(f recurringpb.RecurrenceFrequency) string {
	switch f {
	case recurringpb.RecurrenceFrequency_RECURRENCE_FREQUENCY_DAILY:
		return "daily"
	case recurringpb.RecurrenceFrequency_RECURRENCE_FREQUENCY_WEEKLY:
		return "weekly"
	case recurringpb.RecurrenceFrequency_RECURRENCE_FREQUENCY_MONTHLY:
		return "monthly"
	case recurringpb.RecurrenceFrequency_RECURRENCE_FREQUENCY_QUARTERLY:
		return "quarterly"
	case recurringpb.RecurrenceFrequency_RECURRENCE_FREQUENCY_YEARLY:
		return "yearly"
	default:
		return "monthly"
	}
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

func frequencyBadgeVariant(frequency string) string {
	switch frequency {
	case "daily":
		return "amber"
	case "weekly":
		return "navy"
	case "monthly":
		return "info"
	case "quarterly":
		return "sage"
	case "yearly":
		return "terracotta"
	default:
		return "default"
	}
}

func frequencyDisplayLabel(l fycha.RecurringTemplateLabels, frequency string) string {
	switch frequency {
	case "daily":
		return l.Frequency.Daily
	case "weekly":
		return l.Frequency.Weekly
	case "monthly":
		return l.Frequency.Monthly
	case "quarterly":
		return l.Frequency.Quarterly
	case "yearly":
		return l.Frequency.Yearly
	default:
		return frequency
	}
}
