import { test, expect } from '@playwright/test';

test.describe('API Health', () => {
  test('GET /api/health returns 200 ok', async ({ request }) => {
    const res = await request.get('/api/health');
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('status', 'ok');
    expect(body).toHaveProperty('mysql');
  });

  test('health endpoint works without auth header', async ({ request }) => {
    const res = await request.get('/api/health', {
      headers: {},
    });
    expect(res.status()).toBe(200);
  });

  test('health response time is reasonable', async ({ request }) => {
    const start = Date.now();
    await request.get('/api/health');
    const elapsed = Date.now() - start;
    expect(elapsed).toBeLessThan(10000);
  });
});
