# fycha-golang

Accounting and finance domain package for Ichizen OS. Owns **eight esqyma proto domains** -- asset, expenditure, finance, funding, ledger, payroll, tax, treasury -- making it the largest multi-domain package in the monorepo.

**Module path:** `github.com/erniealice/fycha-golang`

## Domain ownership

fycha maps to eight `proto/v1/domain/` directories in esqyma. Three domains are split with other packages:

| esqyma domain | fycha owns | other package owns |
|---|---|---|
| `asset` | all entities (asset, asset_category, asset_revaluation, depreciation, depreciation_run) | -- |
| `expenditure` | prepayment | centymo: core expenditure entities |
| `finance` | forex_rate | -- |
| `funding` | all entities (fund, fund_allocation, fund_transaction) | -- |
| `ledger` | all entities (account, journal_entry, journal_line, equity_account, equity_transaction, equity_dashboard, recurring_journal_template, fiscal_period) | -- |
| `payroll` | all entities (payroll_run, payroll_remittance, payroll_employee, payroll_settings, payroll_dashboard) | -- |
| `tax` | tax_rate | entydad: tax_registration |
| `treasury` | loan, loan_payment, petty_cash_fund, replenishment, voucher, withholding_certificate | centymo: collection, disbursement, advance |

Reports are a **service-driven surface** under `service/report/` -- they are not an esqyma entity domain and are charter-exempt from entity placement rules.

## Package structure (Option B)

Under Option B the ENTITY is the contract package. Each `domain/<d>/<e>/` directory is one esqyma entity. The domain facade (`domain/<d>/<d>.go`) re-exports entity-local types as Go type aliases so consumers never change their import paths.

```
fycha-golang/
  placement_test.go            # B-STRICT placement gate -- the ONLY test file at root
  assets.go                    # //go:embed assets -- static CSS/JS embed (root stub)
  go.mod / go.sum
  assets/{css,js}/             # embedded stylesheets + filter JS
  domain/
    asset/                     # package asset -- facade for the asset domain
      asset.go                 # facade: type AssetLabels = asset.Labels, etc. (aliases only)
      asset_module.go          # NewAssetModule() assembler
      depreciation_run_module.go # NewDepreciationRunModule() assembler
      asset/                   # entity: esqyma asset/asset
        descriptor.go          # compose.Unit descriptor (Describe())
        labels.go              # Labels struct
        routes.go              # Routes struct + DefaultRoutes()
        embed.go               # template embed.FS
        list/page.go           # list page handler
        detail/                # detail page + tabs + deps + attachment
        action/                # CRUD action handlers
        dashboard/page.go      # asset dashboard handler
        form/                  # form.Data structs
        templates/             # HTML templates
      asset_category/          # entity: esqyma asset/asset_category
        descriptor.go, action/, policies/, templates/
      asset_revaluation/       # entity: esqyma asset/asset_revaluation
        labels.go
      depreciation_run/        # entity: esqyma asset/depreciation_run
        descriptor.go, detail/, list/, shared/, templates/
      lapsing_schedule/        # legacyAllow -- fycha view over depreciation schedule
        descriptor.go, list/, templates/
      depreciation_policies/   # legacyAllow -- per-category policy rollup view
        labels.go
    expenditure/               # package expenditure -- facade for the expenditure domain
      expenditure.go           # facade
      prepayment_module.go     # NewPrepaymentModule() assembler
      prepayment/              # entity: esqyma expenditure/prepayment
        descriptor.go, templates/
    finance/                   # package finance -- facade for the finance domain
      finance.go               # facade
      forex_rate_module.go     # NewForexRateModule() assembler
      forex_rate/              # entity: esqyma finance/forex_rate
        descriptor.go, list/, templates/
    funding/                   # package funding -- facade (funding_module.go + labels.go)
      funding_module.go        # NewFundingModule() assembler (doubles as facade)
      labels.go                # FundingFormLabels
      fund/                    # entity: esqyma funding/fund (labels.go only)
      fund_allocation/         # entity: esqyma funding/fund_allocation (labels.go only)
      fund_transaction/        # entity: esqyma funding/fund_transaction (labels.go only)
      funding/                 # legacyAllow -- aggregate funding view
        descriptor.go
        allocation/, card/, draw/, labels/, settlement/, source/, transfer/, templates/
    ledger/                    # package ledger -- facade for the ledger domain
      ledger.go                # facade: type AccountLabels = account.Labels, etc.
      ledger_module.go         # NewLedgerModule() assembler (multi-entity: account + journal + fiscal_period + recurring_template + settings)
      equity_module.go         # NewEquityModule() assembler
      account/                 # entity: esqyma ledger/account (labels.go, routes.go -- no descriptor)
      journal/                 # legacyAllow -- esqyma entity is journal_entry/journal_line
      fiscal_period/           # entity: esqyma ledger/fiscal_period (labels.go, routes.go)
      recurring_template/      # legacyAllow -- esqyma entity is recurring_journal_template
      equity/                  # legacyAllow -- composite view (equity_account/equity_transaction/equity_dashboard)
        descriptor.go
        capitalaccounts/, dashboard/, equitytransactions/, templates/
      ledger/                  # legacyAllow -- ledger overview view
        action/, dashboard/, detail/, fiscal/, form/, journal/, journal_action/, journal_detail/, list/, recurring/, reports/, settings/, templates/
      seeder/                  # legacyAllow -- COA seeder helper (default_coa.go, default_permissions.go, default_role_permissions.go)
    payroll/                   # package payroll -- facade for the payroll domain
      payroll.go               # facade
      labels.go, routes.go     # domain-level label/route types
      run_module.go            # NewRunModule() assembler
      remittance_module.go     # NewRemittanceModule() assembler
      payrolldashboard_module.go  # NewPayrollDashboardModule() assembler
      payrollemployee_module.go   # NewPayrollEmployeeModule() assembler
      payrollsettings_module.go   # NewPayrollSettingsModule() assembler
      payrolldashboard/        # domain-view: payrolldashboard (prefixed -- valid R2')
        descriptor.go, templates/
      payrollemployee/         # domain-view: payrollemployee (prefixed -- valid R2')
        descriptor.go, list/, templates/
      payrollsettings/         # domain-view: payrollsettings (prefixed -- valid R2')
        descriptor.go, templates/
      remittance/              # legacyAllow -- esqyma entity is payroll_remittance
        descriptor.go, list/, templates/
      run/                     # legacyAllow -- esqyma entity is payroll_run
        descriptor.go, list/, templates/
    tax/                       # package tax -- facade for the tax domain
      tax.go                   # facade
      tax_rate_module.go       # NewTaxRateModule() assembler
      tax_rate/                # entity: esqyma tax/tax_rate
        descriptor.go, list/, templates/
    treasury/                  # package treasury -- facade for the treasury domain
      treasury.go              # facade: type LoanLabels = loan.Labels, etc.
      routes.go                # domain-level route helpers
      loan_module.go           # NewLoanModule() assembler
      petty_cash_module.go     # NewPettyCashModule() assembler
      withholding_certificate_module.go  # NewWithholdingCertificateModule() assembler
      loan/                    # entity: esqyma treasury/loan
        descriptor.go, dashboard/, loanlist/, loanpayments/, templates/
      loan_payment/            # entity: esqyma treasury/loan_payment (labels.go, routes.go)
      petty_cash/              # legacyAllow -- aggregate view (petty_cash_fund/replenishment/voucher)
        descriptor.go, templates/
      withholding_certificate/ # entity: esqyma treasury/withholding_certificate
        descriptor.go, action/, form/, list/, templates/
  service/
    report/                    # report routes/labels/filter/toolbar facade (charter-exempt)
      labels.go, routes.go, filter.go, toolbar.go
      views/                   # one folder per report view
        module.go, embed.go, placeholder.go
        balance_sheet/         # Balance Sheet
        cash_book/             # Cash Book (+ export handler)
        cash_flow/             # Cash Flow Statement
        collection_summary_report/  # Collection Summary
        cost_of_sales/         # Cost of Sales
        dashboard/             # Reports Dashboard
        disbursement_report/   # Disbursement Report
        equity_changes/        # Statement of Changes in Equity
        expenditure_report/    # Expenditure Report
        expenses/              # Expenses List
        financial/             # Financial Statements (income_statement, balance_sheet sub-views)
        gross_profit/          # Gross Profit
        income_statement/      # Income Statement
        net_profit/            # Net Profit
        payables_aging/        # Payables Aging (simple)
        payables_aging_report/ # Payables Aging (parameterized)
        receivables_aging_report/  # Receivables Aging
        revenue/               # Revenue List
        revenue_report/        # Revenue Report
        templates/             # shared report HTML templates
    document/                  # document generation service
      document_service.go      # DocumentService (storage + template orchestration)
  services/
    doctemplate/               # DOCX template engine (placeholder substitution, XML processing)
      engine.go, placeholder.go, xmlprocessor.go, docx.go, testdata/
    pdfconv/                   # PDF converter (centymo imports this for invoice generation)
      converter.go
  block/
    block.go                   # Block() constructor -- pyeza.AppOption entry point
    options.go                 # BlockOption / blockConfig / WithX / wantX accessors
    usecases.go                # *UseCases typed wiring contract + RequireFor + MustValidate
    catalog.go                 # composition-v2 unit binders (AllUnits / per-entity XxxUnit)
    asset.go                   # asset-module wiring (wireAssetModule + proto<->record translators)
    callbacks.go               # lapsing-schedule + revaluation callback helpers
    dashboard_wiring.go        # nil-safe dashboard closure wiring (ledger/equity/payroll/loan)
    helpers.go                 # workspace / functional-currency helpers
    infra.go                   # infra deps type
    mustvalidate_test.go       # MustValidate fail-closed wiring tests
    usecases_test.go           # RequireFor completeness tests
  tests/
    e2e/                       # Playwright E2E specs
```

## Placement gate (`placement_test.go`)

fycha carries a **B-STRICT** placement gate (v2, Option B). `legacyAllow` holds 14 dated residuals (all EXPIRES 2026-07-15) pending dir renames or restructures. The target state is empty (STRICT).

| Rule | What it checks |
|------|----------------|
| **R1** Empty root | No package `.go` files at module root -- only `_test.go` permitted (excused: `assets.go` embed stub) |
| **R2** Canonical dirs | Every first-level dir is an allowed infra surface; every `domain/<d>` is an esqyma proto domain |
| **R2'** Entity dirs | Every `domain/<d>/<child>/` DIR is an esqyma entity of domain `<d>`, `shared`, or a domain-view (name starts with `<d>`) |
| **R3'** Entity contract | No real `*Labels`/`*Routes` type declaration at the domain root -- only alias re-exports (`type X = pkg.Y`) are allowed |
| **R4** No god-files | No `.go` file (excl. `_test.go`) may exceed 1,200 lines |
| **R5** Facade exists | A facade `domain/<d>/<d>.go` must exist for every domain dir with >=1 entity subdir |
| **R6** No cycles | Enforced by `lint-no-domain-cycles.sh` (external, go-list based) |

`crossCutting = false` -- the domain variant applies. esqyma's `proto/v1/domain/` is located at test time so the rules never drift from the live proto tree.

Current `legacyAllow` residuals (all EXPIRES 2026-07-15):

- **R1:** `assets.go` (embed FS stub at root), `docs` (planning markdown)
- **R2':** `domain/asset/lapsing_schedule`, `domain/asset/depreciation_policies`, `domain/funding/funding`, `domain/ledger/equity`, `domain/ledger/journal`, `domain/ledger/ledger`, `domain/ledger/recurring_template`, `domain/ledger/seeder`, `domain/payroll/remittance`, `domain/payroll/run`, `domain/treasury/petty_cash`
- **R5:** `domain/funding/funding.go` (facade is `funding_module.go` + `labels.go`, not the conventional `funding.go`)

## Fail-closed wiring (`block/usecases.go`)

`*UseCases` is the typed wiring contract between service-admin's composition layer and fycha's view modules. `RequireFor(cfg)` lists every missing REQUIRED closure for the enabled modules. `MustValidate(cfg)` adds fail-closed posture:

- **dev/test** (`testing.Testing()` true or `FYCHA_BLOCK_STRICT` truthy): PANIC with the full field list -- uncatchable-by-accident, stack-traced, fails CI loudly.
- **prod**: `log.Printf("FATAL: ...")` at the seam AND returns the error -- `Block()` propagates -- `NewServiceAdmin` halts boot.

OPTIONAL closures (the four dashboard closures, the dashboard-only Loans/Equity/Payroll modules, the TODO-stub Cash/Expenses/Financial modules, `Workspace.Read`, the FiscalPeriod mutators, and `Revaluation`) are never flagged -- they degrade gracefully to empty-state or disabled CTA.

The `UseCases` struct is organized by domain:

| Group | Struct | Closures |
|---|---|---|
| Workspace | `WorkspaceUseCases` | `Read` (functional currency lookup) |
| Asset | `AssetUseCases` | `GetListPageData`, `Create`, `Read`, `Update`, `SetActive`, `Category.ListWithPolicyRollup` |
| DepRun | `DepreciationRunUseCases` | `ListCandidates`, `Generate`, `List`, `Read`, `ListEntries` |
| Revaluation | `RevaluationUseCases` | `Revalue`, `Preview` |
| Ledger | `LedgerUseCases` | `Account.{GetListPageData,Create,Read,Update,Delete}`, `JournalEntry.{GetListPageData,Create,Read,Update,Delete,Post,Reverse}` |
| FiscalPeriod | `FiscalPeriodUseCases` | `GetListPageData`, `Create`, `Close` |
| Tax | `TaxUseCases` | `ListTaxRates` |
| Finance | `FinanceUseCases` | `ListForexRates` |
| Treasury | `TreasuryUseCases` | `ListWithholdingCertificates` |
| Reports | `ReportsUseCases` | 5 sub-groups: `ARAging`, `APAging`, `GrossCashFlow`, `DomainSpecific`, `Statements` |
| Dashboards | (top-level closures) | `GetLedgerDashboardPageData`, `GetEquityDashboardPageData`, `GetPayrollDashboardPageData`, `GetLoanDashboardPageData` |

## Block -- the composition entry point

`block.Block()` is fycha's Lego entry: it returns a `pyeza.AppOption` that registers the selected modules. Consumer apps import it and (optionally) alias:

```go
import fychablock "github.com/erniealice/fycha-golang/block"

fychablock.Block()                  // all modules
fychablock.Block(                   // selective
    fychablock.WithReports(),
    fychablock.WithLedger(),
    fychablock.WithUseCases(uc),    // required -- the typed wiring contract
)
```

Twelve module toggles: `WithReports`, `WithAsset`, `WithLedger`, `WithLoans`, `WithEquity`, `WithPayroll`, `WithCash`, `WithExpenses`, `WithFinancial`, `WithTaxRate`, `WithForexRate`, `WithWithholdingCertificate`. When no module-toggling option is passed, all modules are enabled (`enableAll`). Non-module options (`WithUseCases`, `WithAssetDepreciationRunURL`) never flip `enableAll` off.

## Reports

Report views live under `service/report/views/<report>/` (19 report views + shared templates) and consume service-driven typed closures via `useCases.Reports.<Group>.<UseCase>`. The five report groups -- AR aging, AP aging, gross-cashflow, domain-specific (revenue/expenditure/disbursement), and counterparty statements -- were migrated out of the retired `fycha.DataSource` duck interface into proto-shaped closures during Wave B (2026-05-20/21). Reports query **operational entities, never GL journals** (monetary amounts are centavos -- display / 100).

## Document service

`service/document/` orchestrates document template processing with storage I/O. It delegates core processing to `services/doctemplate` (a pure bytes-in/bytes-out DOCX template engine with XML/placeholder substitution) and adds the storage layer via a `StorageReadWriter` interface. `services/pdfconv` provides PDF conversion; centymo imports `pdfconv` for invoice generation.

## Private services

`services/doctemplate` and `services/pdfconv` are chartered private helpers under `services/` (an allowed first-level directory). They are not exported as a separate module.

- `services/doctemplate` -- DOCX template engine: reads a `.docx` template, substitutes placeholders via XML processing, and returns the result as bytes. Has its own test suite + testdata fixtures.
- `services/pdfconv` -- PDF converter: wraps document-to-PDF conversion. Imported by centymo for invoice PDF generation.

## COA seeder

`domain/ledger/seeder/` provides default Chart of Accounts seed data (`default_coa.go`), default permissions (`default_permissions.go`), and default role-permission mappings (`default_role_permissions.go`). This is a legacyAllow residual pending relocation to `service/` or `scripts/`.

## Dependencies

- `github.com/erniealice/pyeza-golang` -- UI framework (view system, template engine, types, compose units)
- `github.com/erniealice/esqyma` -- proto schemas (asset, expenditure, finance, funding, ledger, payroll, tax, treasury domains + reporting services)
- `github.com/erniealice/lyngua` -- translation/i18n (labels + routes per business type)
- `github.com/erniealice/espyna-golang` -- `reference.Checker` (in-use gates for deletable entities)
- `github.com/erniealice/hybra-golang` -- cross-cutting helpers (indirect)
- `github.com/beevik/etree` -- XML tree processing (used by `services/doctemplate`)

## Role in the monorepo

fycha sits in the domain layer above pyeza and espyna. Consumer apps (e.g., `apps/service-admin`) call `block.Block()` to mount the accounting domain, supplying a `*UseCases` via `block.WithUseCases(...)`. The typed contract ensures any drift between espyna and fycha is a compile error, not a silent nil.

See `docs/wiki/articles/vertical-slices.md` for the full entity trace and `docs/wiki/articles/package-map.md` for the monorepo dependency graph.

## Verification

```bash
cd packages/fycha-golang
go build ./...                                                   # clean compile
go vet ./...
go test -run Placement ./...                                     # B-STRICT gate
go test ./block/ -run 'MustValidate|RequireFor'                  # fail-closed wiring
bash ../../docs/orchestrate/20260610-package-cleanup/lint-no-domain-cycles.sh fycha   # R6
```
