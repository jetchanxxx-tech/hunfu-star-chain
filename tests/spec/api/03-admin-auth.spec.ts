import { test, expect } from '@playwright/test';

const BASE = '/api/v1/admin/login';

test.describe('Admin Login API', () => {
  test('login with admin/admin123 succeeds', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'admin', password: 'admin123' },
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('token');
    expect(body).toHaveProperty('username', 'admin');
    expect(body).toHaveProperty('role', 'super_admin');
  });

  test('login with steward01/steward123 succeeds', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'steward01', password: 'steward123' },
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('role', 'steward');
  });

  test('login with doctor01/doctor123 succeeds', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'doctor01', password: 'doctor123' },
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('role', 'doctor');
  });

  test('login with operator01/operator123 succeeds', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'operator01', password: 'operator123' },
    });
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(body).toHaveProperty('role', 'operator');
  });

  test('login with wrong password returns 401', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'admin', password: 'wrongpass' },
    });
    expect(res.status()).toBe(401);
    const body = await res.json();
    expect(body).toHaveProperty('error');
    expect(body.error).toContain('用户名或密码错误');
  });

  test('login with nonexistent user returns 401', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'nonexistent_user', password: 'x' },
    });
    expect(res.status()).toBe(401);
  });

  test('login without password returns error', async ({ request }) => {
    const res = await request.post(BASE, {
      data: { username: 'admin' },
    });
    expect(res.status()).toBeGreaterThanOrEqual(400);
    expect(res.status()).toBeLessThanOrEqual(401);
  });
});
