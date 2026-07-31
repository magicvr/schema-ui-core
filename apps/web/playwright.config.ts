import { defineConfig } from "@playwright/test";

// R6 minimal browser matrix (I-008-005): one Chromium project that boots both
// services and exercises the critical path (shell + account context via the
// Web /api proxy to the Go API). This is the minimum evidence the user chose
// over accepting a platform residual.
export default defineConfig({
  testDir: "./e2e",
  timeout: 60_000,
  retries: 0,
  reporter: [["list"]],
  use: {
    baseURL: "http://127.0.0.1:5173",
    trace: "retain-on-failure",
  },
  projects: [{ name: "chromium", use: { browserName: "chromium" } }],
  webServer: [
    {
      command: "go run ./cmd/server",
      cwd: "../api",
      url: "http://127.0.0.1:8080/healthz",
      reuseExistingServer: true,
      timeout: 30_000,
    },
    {
      command: "npm run dev",
      cwd: ".",
      url: "http://127.0.0.1:5173/",
      reuseExistingServer: true,
      timeout: 30_000,
    },
  ],
});
