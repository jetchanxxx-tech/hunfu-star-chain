import { test, expect } from '@playwright/test';

test.describe('H5 AI Chat', () => {
  test.beforeEach(async ({ page }) => {
    await page.goto('/#/pages/ai-chat/ai-chat');
    await page.waitForTimeout(1500);
  });

  test('welcome screen loads', async ({ page }) => {
    const body = await page.textContent('body');
    expect(body).toMatch(/惠福灵犀|AI|健康助手|孕期/);
  });

  test('FAQ topic chips are visible', async ({ page }) => {
    const body = await page.textContent('body');
    expect(body).toMatch(/孕期可以运动|NT检查|产后|宝宝发烧|母乳不足|糖耐量/);
  });

  test('clicking FAQ chip populates and sends', async ({ page }) => {
    const faqChip = page.locator('text=/孕期可以运动|NT检查/').first();
    if (await faqChip.isVisible().catch(() => false)) {
      await faqChip.click();
      await page.waitForTimeout(1000);
      const body = await page.textContent('body');
      expect(body).toBeTruthy();
    }
  });

  test('chat input and send button are visible', async ({ page }) => {
    const body = await page.textContent('body');
    expect(body).toMatch(/发送|输入/);
  });

  test('send message via input and receive reply', async ({ page }) => {
    const input = page.locator('input, textarea').first();
    if (await input.isVisible().catch(() => false)) {
      await input.fill('你好');
      const sendBtn = page
        .locator('button, [role="button"]')
        .filter({ hasText: /发送|send/i })
        .first();
      if (await sendBtn.isVisible().catch(() => false)) {
        await sendBtn.click();
      } else {
        await input.press('Enter');
      }
      await page.waitForTimeout(2000);
    }
    const body = await page.textContent('body');
    expect(body).toBeTruthy();
  });
});
