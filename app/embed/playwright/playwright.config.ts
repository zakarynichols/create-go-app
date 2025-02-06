import type { PlaywrightTestConfig } from "@playwright/test";

const E2E_HOST = process.env.NODE_CLIENT_HOST;

/**
 * Read environment variables from file.
 * https://github.com/motdotla/dotenv
 */
// require('dotenv').config();

/**
 * See https://playwright.dev/docs/test-configuration.
 */
const config: PlaywrightTestConfig = {
  testDir: "./dist",
  use: {
    ignoreHTTPSErrors: true,
    baseURL: `http://${E2E_HOST}:7777/`,
    headless: true,
    browserName: "chromium",
    viewport: { width: 1280, height: 720 },
    screenshot: "only-on-failure",
    video: "retain-on-failure",
  },
  // webServer: {
  //   command: "npm run start",
  //   port: 3000,
  //   timeout: 120 * 1000,
  //   reuseExistingServer: !process.env.CI,
  // },
};

export default config;
