import { Page, expect } from '@playwright/test';

export async function waitForTableLoad(page: Page) {
  await page.locator('.el-table').waitFor({ state: 'visible' });
  await page
    .locator('.el-loading-mask')
    .waitFor({ state: 'hidden', timeout: 10000 })
    .catch(() => {});
}

export async function waitForDialog(page: Page, title?: string) {
  const dialog = page.locator('.el-dialog');
  await dialog.waitFor({ state: 'visible' });
  if (title) {
    await expect(dialog.locator('.el-dialog__title')).toContainText(title);
  }
}

export async function closeDialog(page: Page) {
  await page.locator('.el-dialog__footer button:has-text("取消")').click();
  await page.locator('.el-dialog').waitFor({ state: 'hidden' });
}

export async function expectSuccessMessage(page: Page, text?: string) {
  const msg = page.locator('.el-message--success');
  await msg.waitFor({ state: 'visible', timeout: 5000 });
  if (text) await expect(msg).toContainText(text);
}

export async function expectErrorMessage(page: Page, text?: string) {
  const msg = page.locator('.el-message--error');
  await msg.waitFor({ state: 'visible', timeout: 5000 });
  if (text) await expect(msg).toContainText(text);
}
