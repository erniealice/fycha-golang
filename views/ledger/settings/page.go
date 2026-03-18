// Package settings provides the ledger settings views, including the
// Account Templates page (/app/ledger/settings/account-templates).
//
// The Account Templates page shows the default Philippine service business
// Chart of Accounts (from seeder.DefaultCoA) as a set of template cards.
// Users can preview the template and apply it via an HTMX POST action.
package settings

import (
	"context"
	"fmt"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	"github.com/erniealice/fycha-golang/seeder"
)

// ---------------------------------------------------------------------------
// Template data types
// ---------------------------------------------------------------------------

// CoATemplate describes one pre-built Chart of Accounts template.
type CoATemplate struct {
	ID           string // e.g. "service-ph"
	Icon         string // icon name, e.g. "icon-scissors"
	Name         string // display name
	Description  string // 2-3 sentence description
	AccountCount int    // total accounts in the template
	AssetCount   int
	LiabilityCount int
	EquityCount  int
	RevenueCount int
	ExpenseCount int
	ApplyURL     string // POST action URL
	PreviewURL   string // GET action URL (returns dialog content)
	IsApplied    bool   // whether this template has already been applied
}

// PreviewEntry is a simplified account entry for the preview dialog.
type PreviewEntry struct {
	Code    string
	Name    string
	Element string
	Class   string
	IsGroup bool
	Level   int
}

// ---------------------------------------------------------------------------
// View dependencies
// ---------------------------------------------------------------------------

// Deps holds view dependencies for the account templates settings page.
type Deps struct {
	Routes       fycha.AccountRoutes
	Labels       fycha.AccountLabels
	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels
	// ApplyTemplateURL is the POST action URL for the seeder.
	// Defaults to "/action/ledger/settings/account-templates/apply" if empty.
	ApplyTemplateURL string
	PreviewTemplateURL string
}

// PageData holds the template data for the account templates settings page.
type PageData struct {
	types.PageData
	ContentTemplate string
	Templates       []CoATemplate
	// Summary info about the current CoA state (for the info alert)
	CurrentAccountCount int
	HasExistingAccounts bool
}

// PreviewPageData holds the template data for the preview dialog partial.
type PreviewPageData struct {
	types.PageData
	TemplateName string
	AccountCount int
	Entries      []PreviewEntry
	ApplyURL     string
}

// ---------------------------------------------------------------------------
// Default URLs
// ---------------------------------------------------------------------------

const (
	defaultApplyURL   = "/action/ledger/settings/account-templates/apply"
	defaultPreviewURL = "/action/ledger/settings/account-templates/preview"
)

// ---------------------------------------------------------------------------
// Views
// ---------------------------------------------------------------------------

// NewView creates the account templates settings page (full page).
func NewView(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		applyURL := deps.ApplyTemplateURL
		if applyURL == "" {
			applyURL = defaultApplyURL
		}
		previewURL := deps.PreviewTemplateURL
		if previewURL == "" {
			previewURL = defaultPreviewURL
		}

		templates := buildTemplates(applyURL, previewURL)

		pageData := &PageData{
			PageData: types.PageData{
				CacheVersion:   viewCtx.CacheVersion,
				Title:          "Account Templates",
				CurrentPath:    viewCtx.CurrentPath,
				ActiveNav:      deps.Routes.ActiveNav,
				ActiveSubNav:   "account-templates",
				HeaderTitle:    "Account Templates",
				HeaderSubtitle: "Pre-built Chart of Accounts for your business type",
				HeaderIcon:     "icon-layout",
				CommonLabels:   deps.CommonLabels,
			},
			ContentTemplate:     "account-templates-content",
			Templates:           templates,
			CurrentAccountCount: 0, // Phase 3: replace with DB query
			HasExistingAccounts: false,
		}

		return view.OK("account-templates", pageData)
	})
}

// NewPreviewAction creates the template preview action (GET — returns dialog body partial).
func NewPreviewAction(deps *Deps) view.View {
	return view.ViewFunc(func(ctx context.Context, viewCtx *view.ViewContext) view.ViewResult {
		templateID := viewCtx.Request.URL.Query().Get("template_id")
		if templateID == "" {
			templateID = "service-ph"
		}

		applyURL := deps.ApplyTemplateURL
		if applyURL == "" {
			applyURL = defaultApplyURL
		}

		entries := buildPreviewEntries()
		pageData := &PreviewPageData{
			PageData:     types.PageData{CommonLabels: deps.CommonLabels},
			TemplateName: templateName(templateID),
			AccountCount: len(entries),
			Entries:      entries,
			ApplyURL:     fmt.Sprintf("%s?template_id=%s", applyURL, templateID),
		}

		return view.OK("account-templates-preview", pageData)
	})
}

// ---------------------------------------------------------------------------
// Template builder
// ---------------------------------------------------------------------------

func buildTemplates(applyURL, previewURL string) []CoATemplate {
	// Count accounts from the seeder for the service business template
	accounts := seeder.DefaultCoA()
	var assets, liabilities, equity, revenue, expenses int
	for _, a := range accounts {
		switch a.Element.String() {
		case "ACCOUNT_ELEMENT_ASSET":
			assets++
		case "ACCOUNT_ELEMENT_LIABILITY":
			liabilities++
		case "ACCOUNT_ELEMENT_EQUITY":
			equity++
		case "ACCOUNT_ELEMENT_REVENUE":
			revenue++
		case "ACCOUNT_ELEMENT_EXPENSE":
			expenses++
		}
	}
	total := len(accounts)

	return []CoATemplate{
		{
			ID:             "service-ph",
			Icon:           "icon-scissors",
			Name:           "Service Business (Salon / Spa)",
			Description:    "Standard CoA for Philippine service businesses: salons, spas, clinics, and consulting firms. Includes payroll liabilities (SSS, PhilHealth, Pag-IBIG), VAT accounts, and service-specific revenue accounts.",
			AccountCount:   total,
			AssetCount:     assets,
			LiabilityCount: liabilities,
			EquityCount:    equity,
			RevenueCount:   revenue,
			ExpenseCount:   expenses,
			ApplyURL:       fmt.Sprintf("%s?template_id=service-ph", applyURL),
			PreviewURL:     fmt.Sprintf("%s?template_id=service-ph", previewURL),
			IsApplied:      false, // Phase 3: check if already seeded
		},
		{
			ID:             "retail-ph",
			Icon:           "icon-shopping-bag",
			Name:           "Retail Business",
			Description:    "CoA for Philippine retail businesses with inventory and COGS accounts. Includes supplier payables, merchandise inventory, and product sales revenue streams.",
			AccountCount:   0, // coming soon
			ApplyURL:       "",
			PreviewURL:     "",
			IsApplied:      false,
		},
		{
			ID:             "professional-ph",
			Icon:           "icon-briefcase",
			Name:           "Professional Services",
			Description:    "Streamlined CoA for consultants, accountants, lawyers, and other professional services firms. Minimal inventory, focus on receivables and project billing.",
			AccountCount:   0, // coming soon
			ApplyURL:       "",
			PreviewURL:     "",
			IsApplied:      false,
		},
	}
}

func templateName(id string) string {
	switch id {
	case "service-ph":
		return "Service Business (Salon / Spa)"
	case "retail-ph":
		return "Retail Business"
	case "professional-ph":
		return "Professional Services"
	default:
		return "Standard Template"
	}
}

// ---------------------------------------------------------------------------
// Preview entries builder
// ---------------------------------------------------------------------------

// buildPreviewEntries builds a flat, hierarchical list of accounts from the
// seeder data, suitable for rendering as a read-only tree in the preview dialog.
// The seeder data is flat (no explicit group rows), so we synthesize group rows
// for each classification grouping.
func buildPreviewEntries() []PreviewEntry {
	accounts := seeder.DefaultCoA()

	// We emit synthetic element-level headers + classification headers
	// then the leaf accounts under each.
	var entries []PreviewEntry

	type elementGroup struct {
		element string
		label   string
		code    string
	}
	elements := []elementGroup{
		{"ACCOUNT_ELEMENT_ASSET", "ASSETS", "1000"},
		{"ACCOUNT_ELEMENT_LIABILITY", "LIABILITIES", "2000"},
		{"ACCOUNT_ELEMENT_EQUITY", "EQUITY", "3000"},
		{"ACCOUNT_ELEMENT_REVENUE", "REVENUE", "4000"},
		{"ACCOUNT_ELEMENT_EXPENSE", "EXPENSES", "5000"},
	}

	for _, eg := range elements {
		// Find accounts for this element
		var elementAccounts []seeder.DefaultCoAEntry
		for _, a := range accounts {
			if a.Element.String() == eg.element {
				elementAccounts = append(elementAccounts, a)
			}
		}
		if len(elementAccounts) == 0 {
			continue
		}

		// Emit element group header
		entries = append(entries, PreviewEntry{
			Code:    eg.code,
			Name:    eg.label,
			Element: elementLabel(eg.element),
			IsGroup: true,
			Level:   0,
		})

		// Group by classification
		seen := map[string]bool{}
		for _, a := range elementAccounts {
			classLabel := classificationLabel(a.Classification.String())
			if !seen[classLabel] {
				seen[classLabel] = true
				entries = append(entries, PreviewEntry{
					Code:    "",
					Name:    classLabel,
					Element: elementLabel(eg.element),
					Class:   classLabel,
					IsGroup: true,
					Level:   1,
				})
			}
			entries = append(entries, PreviewEntry{
				Code:    a.Code,
				Name:    a.Name,
				Element: elementLabel(eg.element),
				Class:   classificationLabel(a.Classification.String()),
				IsGroup: false,
				Level:   2,
			})
		}
	}

	return entries
}

func elementLabel(pbElement string) string {
	switch pbElement {
	case "ACCOUNT_ELEMENT_ASSET":
		return "Asset"
	case "ACCOUNT_ELEMENT_LIABILITY":
		return "Liability"
	case "ACCOUNT_ELEMENT_EQUITY":
		return "Equity"
	case "ACCOUNT_ELEMENT_REVENUE":
		return "Revenue"
	case "ACCOUNT_ELEMENT_EXPENSE":
		return "Expense"
	default:
		return pbElement
	}
}

func classificationLabel(pbClass string) string {
	switch pbClass {
	case "ACCOUNT_CLASSIFICATION_CURRENT_ASSET":
		return "Current Assets"
	case "ACCOUNT_CLASSIFICATION_NON_CURRENT_ASSET":
		return "Non-Current Assets"
	case "ACCOUNT_CLASSIFICATION_CURRENT_LIABILITY":
		return "Current Liabilities"
	case "ACCOUNT_CLASSIFICATION_NON_CURRENT_LIABILITY":
		return "Non-Current Liabilities"
	case "ACCOUNT_CLASSIFICATION_EQUITY":
		return "Equity"
	case "ACCOUNT_CLASSIFICATION_OPERATING_REVENUE":
		return "Operating Revenue"
	case "ACCOUNT_CLASSIFICATION_OTHER_INCOME":
		return "Other Income"
	case "ACCOUNT_CLASSIFICATION_COST_OF_SALES":
		return "Cost of Sales"
	case "ACCOUNT_CLASSIFICATION_OPERATING_EXPENSE":
		return "Operating Expenses"
	case "ACCOUNT_CLASSIFICATION_FINANCE_COST":
		return "Finance Costs"
	default:
		return pbClass
	}
}
