import { test, expect } from '@playwright/test';
import { waitForHtmxSettle } from '../helpers/htmx';

/**
 * FYC-NAV-001: Page Header Updates on Navigation
 *
 * Verifies that the <h1> heading changes when switching between pages
 * via sidebar links. Tests structural change only — does not assert
 * specific heading text.
 */

test.describe('FYC-NAV-001: Page Header Updates on Navigation', () => {
  test('heading changes when switching financial report pages', async ({ page }) => {
    // Navigate to revenue report
    await page.goto('/app/financial/revenue');
    const h1 = page.locator('h1').first();
    await expect(h1).toBeVisible();
    const revenueHeading = await h1.textContent();

    // Navigate to cost-of-sales via sidebar link
    await page.locator('nav a[href="/app/financial/cost-of-sales"]').click();
    await waitForHtmxSettle(page);

    // Heading should be different from the revenue heading
    const cosHeading = await h1.textContent();
    expect(cosHeading).not.toEqual(revenueHeading);
  });

  test('heading changes between assets dashboard and list', async ({ page }) => {
    await page.goto('/app/assets/dashboard');
    const h1 = page.locator('h1').first();
    await expect(h1).toBeVisible();
    const dashHeading = await h1.textContent();

    // Navigate to assets list via sidebar link
    await page.locator('nav a[href*="/assets/list"]').click();
    await waitForHtmxSettle(page);

    // Heading should be different from the dashboard heading
    const listHeading = await h1.textContent();
    expect(listHeading).not.toEqual(dashHeading);
  });
});
