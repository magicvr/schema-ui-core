import { mkdtempSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

import { defineConfig } from "@playwright/test";

// R6 minimal browser matrix (I-008-005): one Chromium project that boots both
// services and exercises the critical path (shell + account context via the
// Web /api proxy to the Go API). This is the minimum evidence the user chose
// over accepting a platform residual.
//
// WEB_PORT defaults to 5173 (CI/Linux product port). On Windows hosts where
// 5173 falls in a Hyper-V excluded range, set WEB_PORT=9999 (or any free port)
// for local runs without changing the committed default.
//
// Each Playwright run gets a fresh SQLite file so seedRBAC / seedRecords are
// deterministic and parallel browser specs do not fight a developer DB.
const webPort = Number(process.env.WEB_PORT || 5173);
const webOrigin = `http://127.0.0.1:${webPort}`;
const e2eDbPath = join(mkdtempSync(join(tmpdir(), "schema-ui-e2e-")), "e2e.db");

export default defineConfig({
  testDir: "./e2e",
  timeout: 60_000,
  retries: 0,
  // SQLite API is single-writer; keep browser specs serial against one server.
  workers: 1,
  fullyParallel: false,
  reporter: [["list"]],
  use: {
    baseURL: webOrigin,
    trace: "retain-on-failure",
  },
  projects: [{ name: "chromium", use: { browserName: "chromium" } }],
  webServer: [
    {
      command: "go run ./cmd/server",
      cwd: "../api",
      url: "http://127.0.0.1:8080/healthz",
      // Never reuse: a leftover developer server may point at a non-seeded DB
      // without menu_list_edit_lifecycle grants.
      reuseExistingServer: false,
      timeout: 60_000,
      env: {
        ...process.env,
        DB_PATH: e2eDbPath,
        ADMIN_INITIAL_PASSWORD: "admin",
        APP_ENV: "development",
      },
    },
    {
      command: "npm run dev",
      cwd: ".",
      url: `${webOrigin}/`,
      reuseExistingServer: false,
      timeout: 30_000,
      env: {
        ...process.env,
        WEB_PORT: String(webPort),
      },
    },
  ],
});
