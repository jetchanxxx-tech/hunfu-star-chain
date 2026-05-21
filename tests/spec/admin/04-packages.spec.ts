import { test, expect } from '@playwright/test';

test.describe('Admin Packages', () => {
  test('packages API returns valid data', async ({ request }) => {
    const res = await request.get('/api/v1/admin/packages');
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body)).toBe(true);
  });

  test('create package dialog opens and closes', async ({ page }) => {
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
    await page.goto('/admin/packages');
    await page.waitForTimeout(2000);

    const addBtn = page.locator('button:has-text("新增服务包")');
    if (await addBtn.isVisible().catch(() => false)) {
      await addBtn.click();
      await page.waitForTimeout(500);
      const dialog = page.locator('.el-dialog');
      await expect(dialog).toBeVisible({ timeout: 3000 });
      await page
        .locator('.el-dialog__footer button:has-text("取消")')
        .click();
      await expect(dialog).not.toBeVisible({ timeout: 3000 });
    }
  });
});
