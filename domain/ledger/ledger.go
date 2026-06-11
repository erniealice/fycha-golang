package ledger

// ledger.go — ledger-domain facade.
//
// Re-exports every ledger entity package's Routes/Labels with the entity
// prefix so consumers (block/) keep writing ledger.DefaultAccountRoutes(),
// ledger.DefaultLedgerStatementRoutes(), ledger.DefaultEquityLabels(), etc.
//
// Note: the LedgerStatement routes live in the account entity package
// (account.StatementRoutes); LedgerSettings routes live in the
// recurring_template entity package (recurring_template.SettingsRoutes).

import (
	account "github.com/erniealice/fycha-golang/domain/ledger/account"
	equity "github.com/erniealice/fycha-golang/domain/ledger/equity"
	fiscalperiod "github.com/erniealice/fycha-golang/domain/ledger/fiscal_period"
	journal "github.com/erniealice/fycha-golang/domain/ledger/journal"
	recurringtemplate "github.com/erniealice/fycha-golang/domain/ledger/recurring_template"
)

// ---------------------------------------------------------------------------
// Account (domain/ledger/account)
// ---------------------------------------------------------------------------

type AccountRoutes = account.Routes
type AccountLabels = account.Labels

func DefaultAccountRoutes() AccountRoutes { return account.DefaultRoutes() }
func DefaultAccountLabels() AccountLabels { return account.DefaultLabels() }

// ---------------------------------------------------------------------------
// LedgerStatement (domain/ledger/account — StatementRoutes)
// ---------------------------------------------------------------------------

type LedgerStatementRoutes = account.StatementRoutes

func DefaultLedgerStatementRoutes() LedgerStatementRoutes { return account.DefaultStatementRoutes() }

// ---------------------------------------------------------------------------
// Journal (domain/ledger/journal)
// ---------------------------------------------------------------------------

type JournalRoutes = journal.Routes
type JournalLabels = journal.Labels

func DefaultJournalRoutes() JournalRoutes { return journal.DefaultRoutes() }
func DefaultJournalLabels() JournalLabels { return journal.DefaultLabels() }

// ---------------------------------------------------------------------------
// FiscalPeriod (domain/ledger/fiscal_period)
// ---------------------------------------------------------------------------

type FiscalPeriodRoutes = fiscalperiod.Routes
type FiscalPeriodLabels = fiscalperiod.Labels

func DefaultFiscalPeriodRoutes() FiscalPeriodRoutes { return fiscalperiod.DefaultRoutes() }
func DefaultFiscalPeriodLabels() FiscalPeriodLabels { return fiscalperiod.DefaultLabels() }

// ---------------------------------------------------------------------------
// LedgerSettings + RecurringTemplate (domain/ledger/recurring_template)
// ---------------------------------------------------------------------------

type LedgerSettingsRoutes = recurringtemplate.SettingsRoutes
type RecurringTemplateLabels = recurringtemplate.Labels

func DefaultLedgerSettingsRoutes() LedgerSettingsRoutes {
	return recurringtemplate.DefaultSettingsRoutes()
}
func DefaultRecurringTemplateLabels() RecurringTemplateLabels {
	return recurringtemplate.DefaultLabels()
}

// ---------------------------------------------------------------------------
// Equity (domain/ledger/equity)
// ---------------------------------------------------------------------------

type EquityRoutes = equity.Routes
type EquityLabels = equity.Labels

func DefaultEquityRoutes() EquityRoutes { return equity.DefaultRoutes() }
func DefaultEquityLabels() EquityLabels { return equity.DefaultLabels() }
