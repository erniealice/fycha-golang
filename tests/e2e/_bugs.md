# Fycha E2E Test Bugs

Discovered during E2E test writing on 2026-03-19.

---

## BUG-001: Income Statement renders "Page content not available"

**Route:** `/app/reports/income-statement`
**Severity:** P0 (page completely broken)
**Status:** Open

**Description:** Navigating to the Income Statement page shows "Page content not available" instead of the financial statement table. The route IS registered in `apps/service-admin/internal/presentation/financial/module.go` (line 55) and the sidebar link works, but the template content fails to render.

**Root cause (suspected):** The fycha-golang submodule reference in service-admin is out of date. The compiled binary does not include the `income-statement-content` template that the view handler returns via `view.OK("income-statement-content", pageData)`. The pyeza `renderContent` function cannot find the template and falls back to the "Page content not available" message.

**Affected tests:** 5 tests skipped in `reports-navigation.spec.ts` (FYC-RPT-002)

**Fix:** Update the fycha-golang submodule ref in the service-admin go.work / go.mod, rebuild, and verify the template renders.

---

## BUG-002: Balance Sheet renders "Page content not available"

**Route:** `/app/reports/balance-sheet`
**Severity:** P0 (page completely broken)
**Status:** Open

**Description:** Same issue as BUG-001. The Balance Sheet page shows "Page content not available". The route is registered at `module.go` line 56 but the `balance-sheet-content` template is not found at runtime.

**Root cause:** Same as BUG-001 — submodule version mismatch.

**Affected tests:** 4 tests skipped in `reports-navigation.spec.ts` (FYC-RPT-003)

**Fix:** Same as BUG-001.

---

## BUG-003: Cash Flow table missing `id="cash-flow-table"` attribute

**Route:** `/app/reports/cash-flow`
**Severity:** P2 (cosmetic / selector mismatch)
**Status:** Open

**Description:** The Cash Flow page renders correctly with full data, but the `<table>` element does not have the `id="cash-flow-table"` attribute that exists in the source template at `views/reports/templates/cash-flow.html:57`. The rendered HTML has `<table class="financial-statement-table fs-collapsible">` without an ID.

**Root cause (suspected):** Same submodule version issue. The compiled binary uses a template version that predates the ID addition. The class `financial-statement-table` is present and usable as a selector workaround.

**Workaround in tests:** Use `.financial-statement-table` class selector instead of `#cash-flow-table`.

**Affected tests:** All FYC-RPT-004 tests use the class workaround (all pass).

---

## BUG-004: Account drawer form has no footer/cancel button

**Route:** `/app/ledger/accounts/list` -> Add Account drawer
**Severity:** P3 (UX gap)
**Status:** Open

**Description:** The account drawer form (`account-drawer-form.html`) renders without a `sheet-form-footer` template. Unlike other drawer forms (suppliers, assets), there are no Save/Cancel buttons at the bottom. The only way to close the drawer is the X close button in the sheet header.

**Root cause:** The account drawer template at `views/ledger/templates/account-drawer-form.html` does not include `{{template "sheet-form-footer" ...}}` — it ends with just `</form>`. This appears to be an oversight in the template.

**Workaround in tests:** Use `#sheet .sheet-close` or `button[aria-label="Close"]` to close the drawer.

**Affected tests:** `accounts-list.spec.ts` "close button closes drawer" uses the sheet close button.

---

## Summary

| Bug | Route | Severity | Tests Affected | Status |
|-----|-------|----------|----------------|--------|
| BUG-001 | `/app/reports/income-statement` | P0 | 5 skipped | Open |
| BUG-002 | `/app/reports/balance-sheet` | P0 | 4 skipped | Open |
| BUG-003 | `/app/reports/cash-flow` | P2 | Workaround in place | Open |
| BUG-004 | `/app/ledger/accounts/list` (drawer) | P3 | Workaround in place | Open |
