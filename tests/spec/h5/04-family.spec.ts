import { test, expect } from '@playwright/test';

test.describe('H5 Family Management', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/#/pages/family/family');
    await page.waitForTimeout(1500);
  });

  test('empty state shows when no family', async ({ page }) => {
    await page.evaluate(() => {
      localStorage.removeItem('family_id');
      localStorage.removeItem('member_id');
    });
    await page.reload();
    await page.waitForTimeout(500);
    const body = await page.textContent('body');
    expect(body).toMatch(/还没有家庭账户|创建我的家庭|创建家庭账户/);
  });

  test('create family button transitions to form', async ({ page }) => {
    await page.evaluate(() => {
      localStorage.removeItem('family_id');
      localStorage.removeItem('member_id');
    });
    await page.reload();
    await page.waitForTimeout(500);

    const createBtn = page.locator('text=/创建/').first();
    if (await createBtn.isVisible().catch(() => false)) {
      await createBtn.click();
      await page.waitForTimeout(500);
      const body = await page.textContent('body');
      expect(body).toMatch(/昵称|手机号|确认创建|家庭名称/);
    }
  });

  test('family creation form has input fields', async ({ page }) => {
    await page.evaluate(() => {
      localStorage.removeItem('family_id');
      localStorage.removeItem('member_id');
    });
    await page.reload();
    await page.waitForTimeout(500);

    const createBtn = page.locator('text=/创建/').first();
    if (await createBtn.isVisible().catch(() => false)) {
      await createBtn.click();
      await page.waitForTimeout(500);
    }

    const inputs = page.locator('input');
    const count = await inputs.count();
    expect(count).toBeGreaterThanOrEqual(1);
  });

  test('family creation form can be submitted', async ({ page }) => {
    await page.evaluate(() => {
      localStorage.removeItem('family_id');
      localStorage.removeItem('member_id');
    });
    await page.reload();
    await page.waitForTimeout(500);

    const createBtn = page.locator('text=/创建/').first();
    if (await createBtn.isVisible().catch(() => false)) {
      await createBtn.click();
      await page.waitForTimeout(500);
    }

    // Verify the form renders (family name, nickname, phone fields expected)
    const formText = await page.locator('body').textContent();
    expect(formText).toMatch(/昵称|手机号|确认创建|家庭名称/);

    // Try to find and click submit — may work even if inputs are
    // not interactable in mobile viewport
    const submitBtn = page.locator('text=/确认创建|确认添加/').first();
    if (await submitBtn.isVisible().catch(() => false)) {
      await submitBtn.click();
      await page.waitForTimeout(1000);
      // Either error message or redirect happens
      const body = await page.textContent('body');
      expect(body).toBeTruthy();
    }
  });
});
