import { mkdtempSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";

import { defineConfig } from "@playwright/test";

// R6 minimal browser matrix (I-008-005): one Chromium project that boots both
// services and exercises the critical path (shell + account context via the
// Web /api proxy to the Go API). This is the minimum evidence the user chose
// over accepting a platform residual.
//
// The module set follows the repo's YAML-only contract (T-06 / GOAL-013):
// APP_PROFILE is no longer read by the API, so the harness writes a small
// overlay config (app.profile) into a temp file and points CONFIG_FILE at it.
// mvp/admin are the production profiles; demo (W2, GOAL-003 / workspace-010)
// is the non-production demonstration profile (mvp capability + dev.examples).
const appProfile = (process.env.APP_PROFILE || "mvp").trim().toLowerCase();
if (appProfile !== "mvp" && appProfile !== "admin" && appProfile !== "demo") {
  throw new Error(`APP_PROFILE must be mvp, admin or demo for browser E2E (got ${appProfile || "empty"})`);
}
const e2eConfigPath = join(mkdtempSync(join(tmpdir(), "schema-ui-e2e-cfg-")), "config.yaml");
writeFileSync(
  e2eConfigPath,
  `app:
  env: development
  profile: ${appProfile}
`,
  "utf8",
);

// WEB_PORT defaults to 25173 (>25000) so local Windows runs stay outside the
// Hyper-V excluded ranges; WEB_PORT still overrides when another port is needed.
//
// Each Playwright run gets a fresh SQLite file so seedRBAC is deterministic and
// parallel browser specs do not fight a developer DB.
const webPort = Number(process.env.WEB_PORT || 25173);
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
      url: "http://127.0.0.1:25080/readyz",
      // Never reuse: a leftover developer server may point at a non-seeded DB
      // without menu_users grants.
      reuseExistingServer: false,
      timeout: 60_000,
      env: {
        ...process.env,
        CONFIG_FILE: e2eConfigPath,
        DB_PATH: e2eDbPath,
        // W23 (GOAL-034 D-001): PIN the dialect so a gitignored local
        // apps/api/configs/.env (DB_DIALECT=postgres, created 2026-08-21)
        // can never silently redirect the e2e API away from this isolated
        // temp SQLite. Without the pin, fresh-seed admin/admin login 401s
        // against the developer's Postgres and every auth-gated spec fails
        // ("login stays on /" was this isolation drift, not a routing bug).
        DB_DIALECT: "sqlite",
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
