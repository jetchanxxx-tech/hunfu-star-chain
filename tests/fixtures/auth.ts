import { Page } from '@playwright/test';

const ADMIN_TOKENS: Record<string, { token: string; role: string }> = {
  admin: { token: 'admin-session-admin', role: 'super_admin' },
  steward01: { token: 'admin-session-steward01', role: 'steward' },
  doctor01: { token: 'admin-session-doctor01', role: 'doctor' },
  operator01: { token: 'admin-session-operator01', role: 'operator' },
};

export async function setAdminAuth(
  page: Page,
  username: keyof typeof ADMIN_TOKENS = 'admin',
) {
  const info = ADMIN_TOKENS[username];
  // Navigate to login page first (guest route, no redirect), then set token
  await page.goto('/admin/login');
  await page.evaluate(
    ({ token, role, username }) => {
      localStorage.setItem('admin_token', token);
      localStorage.setItem(
        'admin_user',
        JSON.stringify({ username, role, real_name: username }),
      );
    },
    { ...info, username },
  );
}

export async function setH5Auth(
  page: Page,
  memberId: string,
  familyId: string,
) {
  await page.goto('/');
  await page.evaluate(
    ({ mid, fid }) => {
      localStorage.setItem('member_id', mid);
      localStorage.setItem('family_id', fid);
      localStorage.setItem('token', 'mock-h5-token');
    },
    { mid: memberId, fid: familyId },
  );
}

export async function loginAsAdmin(
  page: Page,
  username: string,
  password: string,
) {
  await page.goto('/admin/login');
  await page.waitForTimeout(500);
  await page.locator('.el-input__inner').first().fill(username);
  await page.locator('input[type="password"]').fill(password);
  await page.locator('button:has-text("登录")').click();
}
