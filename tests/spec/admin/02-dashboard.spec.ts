import { test, expect } from '@playwright/test';

test.describe('Admin Dashboard', () => {
  test('dashboard API returns valid stats', async ({ request }) => {
    const res = await request.get('/api/v1/admin/dashboard');
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('total_members');
    expect(typeof body.total_members).toBe('number');
  });

  test('dashboard page handles errors without crash', async ({ page }) => {
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
    await page.goto('/admin/dashboard');
    await page.waitForTimeout(2000);

    // Simulate API failure
    await page.route('**/api/v1/admin/dashboard', (route) =>
      route.fulfill({ status: 500, body: '{"error":"db error"}' }),
    );
    await page.reload();
    await page.waitForTimeout(1000);
    // Page should not crash
    const bodyText = await page.locator('body').textContent();
    expect(bodyText).toBeTruthy();
  });
});
