package ledger

import (
	"context"
	"net/http"

	pyeza "github.com/erniealice/pyeza-golang"
	"github.com/erniealice/pyeza-golang/types"
	"github.com/erniealice/pyeza-golang/view"

	fycha "github.com/erniealice/fycha-golang"
	accountaction "github.com/erniealice/fycha-golang/views/ledger/action"
	accountdetail "github.com/erniealice/fycha-golang/views/ledger/detail"
	fiscalview "github.com/erniealice/fycha-golang/views/ledger/fiscal"
	journalview "github.com/erniealice/fycha-golang/views/ledger/journal"
	journalactionview "github.com/erniealice/fycha-golang/views/ledger/journal_action"
	journaldetailview "github.com/erniealice/fycha-golang/views/ledger/journal_detail"
	accountlist "github.com/erniealice/fycha-golang/views/ledger/list"
	recurringview "github.com/erniealice/fycha-golang/views/ledger/recurring"
	ledgerreports "github.com/erniealice/fycha-golang/views/ledger/reports"
	ledgersettings "github.com/erniealice/fycha-golang/views/ledger/settings"

	accountpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/account"
	fiscalperiodpb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/fiscal_period"
	journalentrypb "github.com/erniealice/esqyma/pkg/schema/v1/domain/ledger/journal_entry"
)

// ModuleDeps holds all dependencies for the ledger module.
// Phase 2: Account views wired to real use cases.
// Phase 3: Journal Entry + FiscalPeriod views wired; GL and Trial Balance with mock data.
type ModuleDeps struct {
	// Account routes
	Routes          fycha.AccountRoutes
	StatementRoutes fycha.LedgerStatementRoutes

	// Journal + FiscalPeriod routes (Phase 3)
	JournalRoutes      fycha.JournalRoutes
	FiscalPeriodRoutes fycha.FiscalPeriodRoutes

	CommonLabels pyeza.CommonLabels
	TableLabels  types.TableLabels

	// Account labels
	Labels fycha.AccountLabels

	// Journal labels (Phase 3)
	JournalLabels fycha.JournalLabels

	// FiscalPeriod labels (Phase 3)
	FiscalPeriodLabels fycha.FiscalPeriodLabels

	// Account use cases
	GetAccountListPageData func(ctx context.Context, req *accountpb.GetAccountListPageDataRequest) (*accountpb.GetAccountListPageDataResponse, error)
	CreateAccount          func(ctx context.Context, req *accountpb.CreateAccountRequest) (*accountpb.CreateAccountResponse, error)
	ReadAccount            func(ctx context.Context, req *accountpb.ReadAccountRequest) (*accountpb.ReadAccountResponse, error)
	UpdateAccount          func(ctx context.Context, req *accountpb.UpdateAccountRequest) (*accountpb.UpdateAccountResponse, error)
	DeleteAccount          func(ctx context.Context, req *accountpb.DeleteAccountRequest) (*accountpb.DeleteAccountResponse, error)

	// Ledger statement use cases (Phase 3: nil → mock data; Phase 4: wire real DB queries)
	GetGeneralLedger func(ctx context.Context, accountID, startDate, endDate string) (*ledgerreports.GLAccountSection, error)
	GetTrialBalance  func(ctx context.Context, asOfDate string) ([]ledgerreports.TBAccountRow, error)

	// Journal Entry use cases (Phase 3)
	GetJournalEntryListPageData func(ctx context.Context, req *journalentrypb.GetJournalEntryListPageDataRequest) (*journalentrypb.GetJournalEntryListPageDataResponse, error)
	GetJournalEntryItemPageData func(ctx context.Context, req *journalentrypb.GetJournalEntryItemPageDataRequest) (*journalentrypb.GetJournalEntryItemPageDataResponse, error)
	CreateJournalEntry          func(ctx context.Context, req *journalentrypb.CreateJournalEntryRequest) (*journalentrypb.CreateJournalEntryResponse, error)
	ReadJournalEntry            func(ctx context.Context, req *journalentrypb.ReadJournalEntryRequest) (*journalentrypb.ReadJournalEntryResponse, error)
	UpdateJournalEntry          func(ctx context.Context, req *journalentrypb.UpdateJournalEntryRequest) (*journalentrypb.UpdateJournalEntryResponse, error)
	DeleteJournalEntry          func(ctx context.Context, req *journalentrypb.DeleteJournalEntryRequest) (*journalentrypb.DeleteJournalEntryResponse, error)
	PostJournalEntry            func(ctx context.Context, req *journalentrypb.PostJournalEntryRequest) (*journalentrypb.PostJournalEntryResponse, error)
	ReverseJournalEntry         func(ctx context.Context, req *journalentrypb.ReverseJournalEntryRequest) (*journalentrypb.ReverseJournalEntryResponse, error)

	// FiscalPeriod use cases (Phase 3; nil-safe — falls back to mock data)
	GetFiscalPeriodListPageData func(ctx context.Context) ([]*fiscalperiodpb.FiscalPeriod, error)
	CreateFiscalPeriod          func(ctx context.Context, req *fiscalperiodpb.CreateFiscalPeriodRequest) (*fiscalperiodpb.CreateFiscalPeriodResponse, error)
	CloseFiscalPeriod           func(ctx context.Context, req *fiscalperiodpb.CloseFiscalPeriodRequest) (*fiscalperiodpb.CloseFiscalPeriodResponse, error)

	// Ledger settings routes + labels (Phase 4: RecurringTemplates + BadDebtPolicy)
	LedgerSettingsRoutes    fycha.LedgerSettingsRoutes
	RecurringTemplateLabels fycha.RecurringTemplateLabels
}

// Module holds all constructed ledger views.
type Module struct {
	routes          fycha.AccountRoutes
	statementRoutes fycha.LedgerStatementRoutes
	journalRoutes   fycha.JournalRoutes
	fiscalRoutes    fycha.FiscalPeriodRoutes

	// Account CRUD
	AccountList      view.View
	AccountDetail    view.View
	AccountTabAction view.View
	AccountAdd       view.View
	AccountEdit      view.View
	AccountDelete    view.View

	// Account search (JSON endpoint for journal form autocomplete)
	accountSearchHandler http.HandlerFunc

	// Account settings
	AccountTemplates        view.View
	AccountTemplatesPreview view.View
	AccountTemplatesApply   view.View

	// Ledger statements (Phase 3: mock data; Phase 4: real DB)
	GeneralLedger view.View
	TrialBalance  view.View

	// Journal Entry views (Phase 3)
	JournalList    view.View
	JournalDetail  view.View
	JournalAdd     view.View
	JournalEdit    view.View
	JournalPost    view.View
	JournalReverse view.View
	JournalDelete  view.View

	// FiscalPeriod views (Phase 3)
	FiscalPeriodList  view.View
	FiscalPeriodAdd   view.View
	FiscalPeriodClose view.View

	// Ledger settings views (Phase 4)
	RecurringTemplates view.View
	BadDebtPolicy      view.View
}

// NewModule creates a ledger module with Account views, GL/TB reports, Journal Entry,
// and FiscalPeriod views wired.
func NewModule(deps *ModuleDeps) *Module {
	// Default statement routes if not provided
	statementRoutes := deps.StatementRoutes
	if statementRoutes.ActiveNav == "" {
		statementRoutes = fycha.DefaultLedgerStatementRoutes()
	}

	listDeps := &accountlist.ListViewDeps{
		Routes:                 deps.Routes,
		Labels:                 deps.Labels,
		CommonLabels:           deps.CommonLabels,
		TableLabels:            deps.TableLabels,
		GetAccountListPageData: deps.GetAccountListPageData,
	}
	detailDeps := &accountdetail.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
		ReadAccount:  deps.ReadAccount,
	}
	actionDeps := &accountaction.Deps{
		Routes:        deps.Routes,
		Labels:        deps.Labels,
		CreateAccount: deps.CreateAccount,
		ReadAccount:   deps.ReadAccount,
		UpdateAccount: deps.UpdateAccount,
		DeleteAccount: deps.DeleteAccount,
	}
	settingsDeps := &ledgersettings.Deps{
		Routes:       deps.Routes,
		Labels:       deps.Labels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
	}

	glDeps := &ledgerreports.GeneralLedgerDeps{
		Routes:           statementRoutes,
		Labels:           deps.Labels,
		CommonLabels:     deps.CommonLabels,
		TableLabels:      deps.TableLabels,
		GetGeneralLedger: deps.GetGeneralLedger,
	}
	tbDeps := &ledgerreports.TrialBalanceDeps{
		Routes:          statementRoutes,
		Labels:          deps.Labels,
		CommonLabels:    deps.CommonLabels,
		TableLabels:     deps.TableLabels,
		GetTrialBalance: deps.GetTrialBalance,
	}

	journalListDeps := &journalview.Deps{
		Routes:                      deps.JournalRoutes,
		Labels:                      deps.JournalLabels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		GetJournalEntryListPageData: deps.GetJournalEntryListPageData,
	}
	journalDetailDeps := &journaldetailview.Deps{
		Routes:                      deps.JournalRoutes,
		Labels:                      deps.JournalLabels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		GetJournalEntryItemPageData: deps.GetJournalEntryItemPageData,
	}
	journalActionDeps := &journalactionview.Deps{
		Routes:                      deps.JournalRoutes,
		Labels:                      deps.JournalLabels,
		CreateJournalEntry:          deps.CreateJournalEntry,
		ReadJournalEntry:            deps.ReadJournalEntry,
		UpdateJournalEntry:          deps.UpdateJournalEntry,
		DeleteJournalEntry:          deps.DeleteJournalEntry,
		PostJournalEntry:            deps.PostJournalEntry,
		ReverseJournalEntry:         deps.ReverseJournalEntry,
		GetJournalEntryItemPageData: deps.GetJournalEntryItemPageData,
	}

	fiscalDeps := &fiscalview.Deps{
		Routes:                      deps.FiscalPeriodRoutes,
		Labels:                      deps.FiscalPeriodLabels,
		CommonLabels:                deps.CommonLabels,
		TableLabels:                 deps.TableLabels,
		GetFiscalPeriodListPageData: deps.GetFiscalPeriodListPageData,
	}

	fiscalActionDeps := &fiscalview.ActionDeps{
		Routes:             deps.FiscalPeriodRoutes,
		Labels:             deps.FiscalPeriodLabels,
		CreateFiscalPeriod: deps.CreateFiscalPeriod,
		CloseFiscalPeriod:  deps.CloseFiscalPeriod,
	}

	// Ledger settings routes: prefer provided, fall back to defaults
	settingsRoutes := deps.LedgerSettingsRoutes
	if settingsRoutes.ActiveNav == "" {
		settingsRoutes = fycha.DefaultLedgerSettingsRoutes()
	}

	recurringDeps := &recurringview.Deps{
		Routes:       settingsRoutes,
		Labels:       deps.RecurringTemplateLabels,
		CommonLabels: deps.CommonLabels,
		TableLabels:  deps.TableLabels,
		// GetRecurringTemplateList: nil — falls back to mock data until DB is wired
	}

	accountSearchDeps := &accountaction.AccountSearchDeps{
		GetAccountListPageData: deps.GetAccountListPageData,
	}

	return &Module{
		routes:          deps.Routes,
		statementRoutes: statementRoutes,
		journalRoutes:   deps.JournalRoutes,
		fiscalRoutes:    deps.FiscalPeriodRoutes,

		accountSearchHandler:    accountaction.NewSearchAccountsHandler(accountSearchDeps),
		AccountList:             accountlist.NewView(listDeps),
		AccountDetail:           accountdetail.NewView(detailDeps),
		AccountTabAction:        accountdetail.NewTabAction(detailDeps),
		AccountAdd:              accountaction.NewAddAction(actionDeps),
		AccountEdit:             accountaction.NewEditAction(actionDeps),
		AccountDelete:           accountaction.NewDeleteAction(actionDeps),
		AccountTemplates:        ledgersettings.NewView(settingsDeps),
		AccountTemplatesPreview: ledgersettings.NewPreviewAction(settingsDeps),
		AccountTemplatesApply:   accountaction.NewApplyTemplateAction(actionDeps),
		GeneralLedger:           ledgerreports.NewGeneralLedgerView(glDeps),
		TrialBalance:            ledgerreports.NewTrialBalanceView(tbDeps),

		JournalList:    journalview.NewView(journalListDeps),
		JournalDetail:  journaldetailview.NewView(journalDetailDeps),
		JournalAdd:     journalactionview.NewAddAction(journalActionDeps),
		JournalEdit:    journalactionview.NewEditAction(journalActionDeps),
		JournalPost:    journalactionview.NewPostAction(journalActionDeps),
		JournalReverse: journalactionview.NewReverseAction(journalActionDeps),
		JournalDelete:  journalactionview.NewDeleteAction(journalActionDeps),

		FiscalPeriodList:  fiscalview.NewView(fiscalDeps),
		FiscalPeriodAdd:   fiscalview.NewAddAction(fiscalActionDeps),
		FiscalPeriodClose: fiscalview.NewCloseAction(fiscalActionDeps),

		RecurringTemplates: recurringview.NewView(recurringDeps),
		BadDebtPolicy:      badDebtPolicyView(deps.CommonLabels),
	}
}

// routeRegistrarFull extends view.RouteRegistrar with HandleFunc support for raw
// http.HandlerFunc routes (e.g. JSON search endpoints).
type routeRegistrarFull interface {
	view.RouteRegistrar
	HandleFunc(method, path string, handler http.HandlerFunc, middlewares ...string)
}

// RegisterRoutes registers all ledger routes with the given route registrar.
func (m *Module) RegisterRoutes(r view.RouteRegistrar) {
	// Accounts — Phase 2: real views
	r.GET(m.routes.ListURL, m.AccountList)
	r.GET(m.routes.DetailURL, m.AccountDetail)
	r.GET("/action/ledger/accounts/{id}/tab/{tab}", m.AccountTabAction)
	r.GET(m.routes.AddURL, m.AccountAdd)
	r.POST(m.routes.AddURL, m.AccountAdd)
	r.GET(m.routes.EditURL, m.AccountEdit)
	r.POST(m.routes.EditURL, m.AccountEdit)
	r.POST(m.routes.DeleteURL, m.AccountDelete)

	// Account search — JSON endpoint for journal form autocomplete
	// GET /action/ledger/accounts/search?q=<term>
	if m.accountSearchHandler != nil {
		if full, ok := r.(routeRegistrarFull); ok {
			full.HandleFunc("GET", "/action/ledger/accounts/search", m.accountSearchHandler)
		}
	}

	// Journals — Phase 3: real views
	r.GET(m.journalRoutes.ListURL, m.JournalList)
	r.GET(m.journalRoutes.DetailURL, m.JournalDetail)
	// Journal actions — Phase 3: add / edit / post / reverse / delete
	r.GET(m.journalRoutes.AddURL, m.JournalAdd)
	r.POST(m.journalRoutes.AddURL, m.JournalAdd)
	r.GET(m.journalRoutes.EditURL, m.JournalEdit)
	r.POST(m.journalRoutes.EditURL, m.JournalEdit)
	r.POST(m.journalRoutes.PostURL, m.JournalPost)
	r.POST(m.journalRoutes.ReverseURL, m.JournalReverse)
	r.POST(m.journalRoutes.DeleteURL, m.JournalDelete)

	// Reports — Phase 3: real views with mock data
	r.GET(m.statementRoutes.GeneralLedgerURL, m.GeneralLedger)
	r.GET(m.statementRoutes.TrialBalanceURL, m.TrialBalance)

	// Settings — Account Templates: real view
	r.GET(m.routes.TemplatesURL, m.AccountTemplates)
	r.GET("/action/ledger/settings/account-templates/preview", m.AccountTemplatesPreview)
	r.POST("/action/ledger/settings/account-templates/apply", m.AccountTemplatesApply)

	// Settings — Phase 3: FiscalPeriod wired
	r.GET(m.fiscalRoutes.ListURL, m.FiscalPeriodList)
	r.GET(m.fiscalRoutes.AddURL, m.FiscalPeriodAdd)
	r.POST(m.fiscalRoutes.AddURL, m.FiscalPeriodAdd)
	r.POST(m.fiscalRoutes.CloseURL, m.FiscalPeriodClose)

	// Settings — Phase 4: RecurringTemplates + BadDebtPolicy wired
	r.GET(fycha.RecurringTemplatesURL, m.RecurringTemplates)
	r.GET(fycha.BadDebtPolicyURL, m.BadDebtPolicy)
}

// badDebtPolicyView returns a view that renders the bad-debt-policy template.
// The template is a coming-soon placeholder; it uses CommonLabels for icon injection.
func badDebtPolicyView(commonLabels pyeza.CommonLabels) view.View {
	return view.ViewFunc(func(_ context.Context, viewCtx *view.ViewContext) view.ViewResult {
		templateName := "bad-debt-policy"
		if viewCtx.IsHTMX {
			templateName = "bad-debt-policy-content"
		}
		return view.OK(templateName, &types.PageData{
			CacheVersion: viewCtx.CacheVersion,
			Title:        "Bad Debt Policy",
			CurrentPath:  viewCtx.CurrentPath,
			ActiveNav:    "ledger",
			ActiveSubNav: "bad-debt-policy",
			HeaderTitle:  "Bad Debt Policy",
			HeaderIcon:   "icon-alert-triangle",
			CommonLabels: commonLabels,
		})
	})
}
