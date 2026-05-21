import { test, expect } from '@playwright/test';
import { setH5Auth } from '../../fixtures/auth';

test.describe('H5 Home Page (Demo Mode)', () => {
  test('demo home page loads without login', async ({ page }) => {
    await page.goto('/');
    await page.waitForTimeout(1000);
    // Verify the page loaded with content
    await expect(page.locator('text=/你好|星球居民/')).toBeVisible({
      timeout: 5000,
    });
  });

  test('greeting shows demo user or resident', async ({ page }) => {
    await page.goto('/');
    await page.waitForTimeout(1000);
    const body = await page.textContent('body');
    expect(body).toMatch(/你好|星球居民|关键节点|最近报告/);
  });

  test('quick action items are visible', async ({ page }) => {
    await page.goto('/');
    await page.waitForTimeout(1000);
    const body = await page.textContent('body');
    expect(body).toMatch(/时光轴|同心圆|服务包|灵犀问答/);
  });

  test('key nodes section renders', async ({ page }) => {
    await page.goto('/');
    await page.waitForTimeout(1000);
    const body = await page.textContent('body');
    expect(body).toMatch(/关键节点|暂无数据/);
  });

  test('recent reports section renders', async ({ page }) => {
    await page.goto('/');
    await page.waitForTimeout(1000);
    const body = await page.textContent('body');
    expect(body).toMatch(/最近报告|暂无数据/);
  });

  test('home page with H5 auth', async ({ page }) => {
    await setH5Auth(page, '1', 'test-family-id');
    await page.goto('/');
    await page.waitForTimeout(1000);
    // Should load without crashing
    const body = await page.textContent('body');
    expect(body).toBeTruthy();
  });
});
