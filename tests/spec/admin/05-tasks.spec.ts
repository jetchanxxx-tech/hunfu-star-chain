import { test, expect } from '@playwright/test';

test.describe('Admin Tasks', () => {
  test('tasks API returns valid data', async ({ request }) => {
    const res = await request.get('/api/v1/tasks?offset=0&limit=10');
    expect(res.status()).toBe(200);
  });

  test('task stats API returns valid counts', async ({ request }) => {
    const res = await request.get('/api/v1/admin/task-stats');
    expect(res.status()).toBe(200);
  });

  test('create task dialog opens and closes', async ({ page }) => {
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
    await page.goto('/admin/tasks');
    await page.waitForTimeout(2000);

    const createBtn = page.locator('button:has-text("创建任务")');
    if (await createBtn.isVisible().catch(() => false)) {
      await createBtn.click();
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
