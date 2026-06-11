# fycha-golang

Accounting domain package for Ichizen OS. Ships the view + wiring layer for the
finance verticals — fixed assets, the general ledger (chart of accounts,
journals, fiscal periods, equity), treasury (loans, petty cash, withholding),
funding, payroll, tax/forex, prepayments, and the P&L / aging / cash-flow report
family — for Go + HTMX admin apps built on pyeza-golang.

**Module path:** `github.com/erniealice/fycha-golang`

**Dependencies:**
- `github.com/erniealice/pyeza-golang` — UI framework (views, types, templates)
- `github.com/erniealice/esqyma` — protobuf schemas (the domain source of truth)
- `github.com/erniealice/lyngua` — translation provider (labels + routes per business type)
- `github.com/erniealice/espyna-golang` — `reference.Checker` (in-use gates)

## Architecture — Option B (domain-first, entity-as-unit)

fycha is laid out **by domain, then by entity** — mirroring esqyma's proto
domains (`packages/esqyma/proto/v1/domain/`). The **entity** is the package unit:
each esqyma entity gets its own folder holding its contract types
(`Labels`/`Routes`) and views. A hand-written **facade** per domain re-exports
those entity-local types under their public names so consumers keep writing
`asset.AssetLabels`, `ledger.DefaultAccountRoutes()`, etc.

The layout is enforced by `placement_test.go` (the v2 B-STRICT gate, adopted
byte-identical from
`docs/orchestrate/20260610-package-cleanup/placement_test.template.go`) and the
`lint-no-domain-cycles.sh` R6 cycle lint. The root carries **no package code** —
only `placement_test.go` and the `assets.go` embed stub.

```
packages/fycha-golang/
  go.mod / go.sum
  placement_test.go        — B-STRICT placement gate (package fycha_test)
  assets.go                — //go:embed assets — static CSS/JS embed (legacy root stub)
  assets/{css,js}/         — embedded stylesheets + filter JS

  domain/<domain>/         — one dir per esqyma proto domain
    <domain>.go            — FACADE: type-alias + value re-exports of the entity packages
    <entity>_module.go     — hoisted module assembler (NewXxxModule + ModuleDeps)
    <entity>/              — ENTITY package: Labels/Routes + views (list/detail/form/action/…)
    <domain>view/          — disambiguated domain-level view (name prefixed with the domain)

  service/                 — service surfaces (cross-entity, not a single esqyma entity)
    report/                — report routes/labels/filter/toolbar facade
      views/<report>/      — one folder per report (balance_sheet, cash_book, gross_profit, …)
    document/              — document_service (PDF/template rendering)

  block/                   — the Lego ASSEMBLER (composition entry point)
  tests/                   — Playwright e2e specs
```

### domain/ layout (realized)

| domain | entity folders | domain-view folders | facade |
|---|---|---|---|
| `asset` | `asset`, `asset_category`, `asset_revaluation`, `depreciation_run`, `lapsing_schedule`*, `depreciation_policies`* | — | `asset.go` |
| `expenditure` | `prepayment` | — | `expenditure.go` |
| `finance` | `forex_rate` | — | `finance.go` |
| `funding` | `fund`, `fund_allocation`, `fund_transaction`, `funding`* | — | `funding_module.go` + `labels.go`* |
| `ledger` | `account`, `fiscal_period`, `journal`*, `equity`*, `ledger`*, `recurring_template`*, `seeder`* | — | `ledger.go` |
| `payroll` | `remittance`*, `run`* | `payrolldashboard`, `payrollemployee`, `payrollsettings` | `payroll.go` |
| `tax` | `tax_rate` | — | `tax.go` |
| `treasury` | `loan`, `loan_payment`, `withholding_certificate`, `petty_cash`* | — | `treasury.go` |

`*` = a residual whose folder name is not a 1:1 esqyma entity (legacy short
name or a cross-entity view aggregate); excused in `placement_test.go`'s
`legacyAllow` with a dated `EXPIRES 2026-07-15` stamp until it is renamed to a
domain-view, split per esqyma entity, or hand-folded into its facade.

### Facade pattern

Each `domain/<d>/<d>.go` is a thin facade: `type AssetLabels = asset.Labels`
aliases plus by-value re-exports of route-URL consts and `DefaultXxx()`
constructors. Entity packages import **nothing** from their domain facade; only
the facade and the hoisted `<entity>_module.go` assemblers import the entity
packages. This breaks the import cycle (`views import root for types`) and keeps
each `domain/<d>/` subtree a shallow child→parent DAG (no intra-domain cycle —
enforced by the R6 lint).

## Block — the composition entry point

`block.Block()` is fycha's Lego entry: it returns a `pyeza.AppOption` that
registers the selected modules. Consumer apps import it and (optionally) alias:

```go
import fychablock "github.com/erniealice/fycha-golang/block"

fychablock.Block()                  // all modules
fychablock.Block(                   // selective
    fychablock.WithReports(),
    fychablock.WithLedger(),
    fychablock.WithUseCases(uc),    // required — the typed wiring contract
)
```

`block/` lives in a sub-package (not the root) to avoid the Go import cycle
(`domain/<d>/<e>/` views import the domain facade for types, so `Block()` cannot
sit in a package those views import). Companion files:

- `options.go` — `BlockOption` / `blockConfig` / `WithX` / `wantX` accessors
- `usecases.go` — the typed `UseCases` wiring contract + **`RequireFor` / `MustValidate`**
- `block.go` — `Block()` assembler (loads routes/labels, registers modules)
- `asset.go` / `callbacks.go` — asset-module wiring + proto↔record translators
- `dashboard_wiring.go` — nil-safe dashboard-closure wiring helpers
- `helpers.go` — workspace / functional-currency helpers

### Typed wiring contract + fail-closed gate

`block.UseCases` declares exactly what fycha needs from the host app
(service-admin's `buildFychaUseCases` constructs it from espyna's consumer
container). `Block()` calls **`UseCases.MustValidate(cfg)`** at entry:

- **`RequireFor`** asserts every REQUIRED closure of each *enabled* module is
  non-nil — each report's report/list/get closure, the asset/ledger CRUD + list
  + page-data closures, and the read-only list modules' single closure. Required
  status lives entirely in `RequireFor`: a field is REQUIRED iff it is asserted
  inside an enabled `if cfg.wantXxx()` block.
- **`MustValidate`** adds the fail-closed posture (mirrors service-admin's
  `AUTHZ_ENFORCE` boot-guard): in dev/test (`go test`, or `FYCHA_BLOCK_STRICT`
  truthy) a missing REQUIRED closure **panics** with the field list; in prod it
  logs a screaming `FATAL:` and returns the error so boot halts. A nil closure
  can never silently render an empty feature.
- **OPTIONAL** (never asserted): the four dashboard closures, the dashboard-only
  Loans/Equity/Payroll modules, the TODO-stub Cash/Expenses/Financial modules,
  `Workspace.Read`, the FiscalPeriod mutators, and `Revaluation` — these
  legitimately degrade to empty-state / disabled-CTA when nil.

## Reports

Report views live under `service/report/views/<report>/` and consume
service-driven typed closures via `useCases.Reports.<Group>.<UseCase>` (AR/AP
aging, gross-cashflow, domain-specific revenue/expenditure/disbursement,
counterparty statements). They query **operational entities, never GL journals**
(monetary amounts are centavos — display ÷100).

## Verification

```bash
cd packages/fycha-golang
go build ./...                                                   # clean
go vet ./...
go test -run Placement ./...                                     # B-STRICT gate
go test ./block/ -run 'MustValidate|RequireFor'                 # fail-closed wiring
bash ../../docs/orchestrate/20260610-package-cleanup/lint-no-domain-cycles.sh fycha   # R6
```
