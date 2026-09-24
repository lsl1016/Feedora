import { defineConfig } from '@playwright/test';

export default defineConfig({
  testDir: './tests/e2e',
  timeout: 30_000,
  fullyParallel: false,
  retries: 1,
  workers: 1,
  reporter: [['list']],
  use: {
    headless: true,
    trace: 'retain-on-failure',
    viewport: { width: 1600, height: 1000 },
  },
  webServer: [
    {
      command: 'pnpm dev:web',
      url: 'http://127.0.0.1:5173',
      reuseExistingServer: false,
      timeout: 120_000,
    },
    {
      command: 'pnpm dev:desktop-ui',
      url: 'http://127.0.0.1:5174',
      reuseExistingServer: false,
      timeout: 120_000,
    },
  ],
});
