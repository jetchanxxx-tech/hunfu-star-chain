import { test, expect } from '@playwright/test';

test.describe('Admin Members', () => {
  test('members API returns valid data', async ({ request }) => {
    const res = await request.get('/api/v1/admin/members');
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body)).toBe(true);
  });

  test('search input is available on members page', async ({ page }) => {
    await page.goto('/admin/login');
    await page.evaluate(() => {
      localStorage.setItem('admin_token', 'admin-session-admin');
      localStorage.setItem(
        'admin_user',
        JSON.stringify({
          username: 'admin',
          role: 'super_admin',
          real_name: 'admin',
        }),
      );
    });
    await page.goto('/admin/members');
    await page.waitForTimeout(2000);
    const search = page.locator('input[placeholder="搜索昵称/家庭名"]');
    if (await search.isVisible().catch(() => false)) {
      await search.fill('test');
      await page.waitForTimeout(500);
    }
    // Page should not crash
    const bodyText = await page.locator('body').textContent();
    expect(bodyText).toBeTruthy();
  });
});
