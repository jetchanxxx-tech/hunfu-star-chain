import { defineConfig, devices } from '@playwright/test';

export default defineConfig({
  testDir: './spec',
  fullyParallel: false,
  retries: 1,
  workers: 1,
  reporter: [
    ['list'],
    ['html', { outputFolder: '../playwright-report' }],
  ],
  timeout: 30000,
  expect: { timeout: 10000 },
  use: {
    baseURL: 'https://huifu.pangu-cloud.com',
    trace: 'on-first-retry',
    screenshot: 'only-on-failure',
  },
  projects: [
    {
      name: 'admin-desktop',
      use: {
        ...devices['Desktop Chrome'],
        viewport: { width: 1440, height: 900 },
      },
      testMatch: 'admin/**/*.spec.ts',
    },
    {
      name: 'h5-mobile',
      use: {
        ...devices['iPhone 13'],
      },
      testMatch: 'h5/**/*.spec.ts',
    },
    {
      name: 'api',
      use: {
        ...devices['Desktop Chrome'],
      },
      testMatch: 'api/**/*.spec.ts',
    },
  ],
});
