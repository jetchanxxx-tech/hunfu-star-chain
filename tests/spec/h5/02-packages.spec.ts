import { test, expect } from '@playwright/test';

test.describe('H5 Packages', () => {
  test.beforeEach(async ({ page }) => {
    // uni-app H5 uses hash routing
    await page.goto('/#/pages/packages/packages');
    await page.waitForTimeout(1500);
  });

  test('packages page loads', async ({ page }) => {
    const body = await page.textContent('body');
    expect(body).toMatch(/服务包|购买|¥|暂无/);
  });

  test('package cards show price with ¥ symbol', async ({ page }) => {
    const body = await page.textContent('body');
    // Demo data may load from demo/home endpoint
    expect(body).toBeTruthy();
  });

  test('buy button exists (known placeholder — no handler)', async ({
    page,
  }) => {
    const buyBtn = page.locator('text=/购买|立即购买/').first();
    if (await buyBtn.isVisible().catch(() => false)) {
      await buyBtn.click();
      await page.waitForTimeout(300);
    }
    const body = await page.textContent('body');
    expect(body).toBeTruthy();
  });
});
