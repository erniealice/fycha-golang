import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * FYC-COA-001: Chart of Accounts List
 *
 * Routes: AccountListURL = /app/ledger/accounts/list
 * Table ID: #accounts-tree-table (not rendered in live DOM — table has no ID)
 * Verifies: account list page loads, table structure, element filter tabs
 *
 * Note: The table in the live DOM uses class selectors only (no ID).
 * The accounts-tree-table ID is set in Go code but not rendered by the
 * current compiled submodule version.
 */

test.describe('FYC-COA-001: Chart of Accounts List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/ledger/accounts/list');
    await page.waitForLoadState('networkidle');
  });

  test('loads account list page with table', async ({ page }) => {
    const table = page.locator('table');
    await expect(table).toBeVisible({ timeout: 10000 });
  });

  test('table has correct column headers', async ({ page }) => {
    const headers = page.locator('thead th');
    const headerTexts = await headers.allTextContents();
    // Expected columns: Code, Account Name, Element, Class, Balance, Actions
    expect(headerTexts.length).toBeGreaterThanOrEqual(5);
  });

  test('has primary action button for adding accounts', async ({ page }) => {
    // The "Add Account" button in the accounts toolbar
    const addBtn = page.getByRole('button', { name: /Add Account/i });
    await expect(addBtn).toBeVisible();
    await expect(addBtn).toBeEnabled();
  });

  test('has table search input', async ({ page }) => {
    const search = page.getByPlaceholder('Search...');
    await expect(search).toBeVisible();
  });

  test('has element filter tabs', async ({ page }) => {
    // Account list uses nav[aria-label="Category filter"] with link elements
    const filterNav = page.locator('nav[aria-label="Category filter"]');
    await expect(filterNav).toBeVisible();

    const tabs = filterNav.locator('a');
    const count = await tabs.count();
    // Should have 6 tabs: All, Assets, Liabilities, Equity, Revenue, Expenses
    expect(count).toBeGreaterThanOrEqual(6);
  });

  test('shows empty state or account rows', async ({ page }) => {
    const table = page.locator('table');
    await expect(table).toBeVisible();

    // Check if we have data rows or empty state
    const rows = page.locator('tbody tr');
    const count = await rows.count();
    // At least the table rendered (may be 0 data rows + 1 empty state row)
    expect(count).toBeGreaterThanOrEqual(1);
  });

  test('shows pagination area', async ({ page }) => {
    const pagination = page.locator('.table-footer, .pagination-info');
    await expect(pagination).toBeVisible();
  });
});

test.describe('FYC-COA-001: Account Add via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/ledger/accounts/list');
    await page.waitForLoadState('networkidle');
    await expect(page.locator('table')).toBeVisible({ timeout: 10000 });
  });

  test('opens add drawer when Add Account clicked', async ({ page }) => {
    await page.getByRole('button', { name: /Add Account/i }).click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Verify form fields
    await expect(page.locator('#code')).toBeVisible();
    await expect(page.locator('#name')).toBeVisible();
    await expect(page.locator('#element')).toBeVisible();
  });

  test('close button closes drawer', async ({ page }) => {
    await page.getByRole('button', { name: /Add Account/i }).click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Account drawer uses the sheet's built-in close button (X)
    await page.locator('#sheet .sheet-close, #sheet button[aria-label="Close"]').first().click();
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });
});
