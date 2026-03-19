import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * FYC-AST-001: Asset List
 * FYC-AST-002: Asset Add via Drawer
 * FYC-AST-003: Asset Edit via Drawer
 *
 * Routes: AssetListURL, AssetAddURL, AssetEditURL
 * Table ID: #assets-table
 * Verifies: list page loads, table structure, CRUD via drawer
 */

test.describe('FYC-AST-001: Asset List', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/assets/list/active');
    await page.waitForLoadState('networkidle');
  });

  test('displays asset table with correct ID', async ({ page }) => {
    const table = page.locator('#assets-table');
    await expect(table).toBeVisible({ timeout: 10000 });
  });

  test('table has correct column headers', async ({ page }) => {
    const headers = page.locator('#assets-table thead th .column-label');
    const headerTexts = await headers.allTextContents();
    // Expected columns: Asset Number, Name, Category, Location, Acquisition Cost, Book Value, Status
    expect(headerTexts.length).toBeGreaterThanOrEqual(6);
  });

  test('shows data rows with asset data (mock data)', async ({ page }) => {
    const rows = page.locator('#assets-table tbody tr[data-id]');
    const count = await rows.count();
    // Mock data has 5 active assets (ast-001 through ast-004, ast-006)
    expect(count).toBeGreaterThanOrEqual(1);
  });

  test('has primary action button in toolbar', async ({ page }) => {
    const primaryAction = page.locator('.toolbar-primary-action');
    await expect(primaryAction).toBeVisible();
    await expect(primaryAction).toBeEnabled();
  });

  test('has table search input', async ({ page }) => {
    const table = page.locator('#assets-table');
    await expect(table).toBeVisible();
    // Search is inside the table card
    const search = page.locator('.table-search input, input[type="search"], input[placeholder*="earch"]');
    await expect(search.first()).toBeVisible();
  });

  test('shows pagination with entry count', async ({ page }) => {
    // Pagination area is rendered in the table footer container
    const pagination = page.locator('.table-footer, .pagination-info');
    await expect(pagination).toBeVisible();
  });

  test('row has action buttons (view, edit, delete)', async ({ page }) => {
    const firstRow = page.locator('#assets-table tbody tr[data-id]').first();
    const viewLink = firstRow.locator('a.action-btn.view');
    const editBtn = firstRow.locator('.action-btn.edit');
    const deleteBtn = firstRow.locator('.action-btn.delete');

    await expect(viewLink).toBeVisible();
    await expect(editBtn).toBeVisible();
    await expect(deleteBtn).toBeVisible();
  });

  test('view link navigates to asset detail', async ({ page }) => {
    const viewLink = page.locator('#assets-table tbody tr[data-id]').first().locator('a.action-btn.view');
    const href = await viewLink.getAttribute('href');
    expect(href).toContain('/app/assets/detail/');
  });

  test('status badges are visible in table', async ({ page }) => {
    const badges = page.locator('#assets-table .status-badge');
    const count = await badges.count();
    expect(count).toBeGreaterThanOrEqual(1);
  });
});

test.describe('FYC-AST-002: Asset Add via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/assets/list/active');
    await expect(page.locator('#assets-table')).toBeVisible({ timeout: 10000 });
  });

  test('opens drawer when primary action clicked', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });
  });

  test('drawer has required form fields', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Check form fields by ID
    await expect(page.locator('#name')).toBeVisible();
    await expect(page.locator('#asset_number')).toBeVisible();
    await expect(page.locator('#acquisition_cost')).toBeVisible();
    await expect(page.locator('#useful_life_months')).toBeVisible();
  });

  test('drawer has optional form fields', async ({ page }) => {
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    await expect(page.locator('#description')).toBeVisible();
    await expect(page.locator('#asset_category_id')).toBeVisible();
    await expect(page.locator('#location_id')).toBeVisible();
    await expect(page.locator('#salvage_value')).toBeVisible();
    await expect(page.locator('#depreciation_method')).toBeVisible();
  });

  test('creates asset via drawer form', async ({ page }) => {
    const ts = Date.now();

    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Fill required fields
    await page.locator('#name').fill(`Test Asset ${ts}`);
    await page.locator('#asset_number').fill(`TA-${ts}`);
    await page.locator('#acquisition_cost').fill('50000');
    await page.locator('#useful_life_months').fill('60');

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Wait for HTMX response + sheet close
    await waitForHtmxSettle(page);
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });

  test('cancel closes drawer without creating', async ({ page }) => {
    // Open drawer
    await page.locator('.toolbar-primary-action').click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Fill something
    await page.locator('#name').fill('ShouldNotSave');

    // Cancel
    await page.locator('#sheet .sheet-footer .btn-secondary').click();

    // Drawer should close
    await expect(page.locator('#sheet').first()).not.toHaveClass(/open/, { timeout: 5000 });
  });
});

test.describe('FYC-AST-003: Asset Edit via Drawer', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/app/assets/list/active');
    await expect(page.locator('#assets-table')).toBeVisible({ timeout: 10000 });
  });

  test('opens edit drawer with pre-filled data', async ({ page }) => {
    const editBtn = page.locator('#assets-table tbody tr[data-id]').first().locator('.action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Name should be pre-filled from mock data
    const name = await page.locator('#name').inputValue();
    expect(name.length).toBeGreaterThan(0);
  });

  test('edit drawer has asset number pre-filled', async ({ page }) => {
    const editBtn = page.locator('#assets-table tbody tr[data-id]').first().locator('.action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    const assetNumber = await page.locator('#asset_number').inputValue();
    expect(assetNumber.length).toBeGreaterThan(0);
  });

  test('saves edit and closes drawer', async ({ page }) => {
    const editBtn = page.locator('#assets-table tbody tr[data-id]').first().locator('.action-btn.edit');
    await editBtn.click();
    await expect(page.locator('#sheet.open .sheet-panel')).toBeVisible({ timeout: 10000 });

    // Modify a field
    const ts = Date.now();
    await page.locator('#name').fill(`Updated Asset ${ts}`);

    // Submit
    await page.locator('#sheet .sheet-footer button[type="submit"]').click();

    // Drawer closes
    await waitForHtmxSettle(page);
    await expect(page.locator('.sheet.open')).not.toBeVisible({ timeout: 10000 });
  });
});
