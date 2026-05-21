import { test, expect } from '@playwright/test';

test.describe('Admin System', () => {
  test('users API returns array with valid fields', async ({ request }) => {
    const res = await request.get('/api/v1/admin/users');
    expect(res.status()).toBe(200);
    const body = await res.json();
    expect(Array.isArray(body)).toBe(true);
    if (body.length > 0) {
      const user = body[0];
      const nameKey = user.username !== undefined ? 'username' : 'Username';
      const roleKey = user.role !== undefined ? 'role' : 'Role';
      expect(user).toHaveProperty(nameKey);
      expect(user).toHaveProperty(roleKey);
    }
  });
});
