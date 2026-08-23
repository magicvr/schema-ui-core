import { mkdtempSync, writeFileSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { tmpdir } from "node:os";
import { join, resolve } from "node:path";

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
// Each Playwright run gets a fresh store so seedRBAC is deterministic and
// parallel browser specs do not fight a developer DB.
const webPort = Number(process.env.WEB_PORT || 25173);
const webOrigin = `http://127.0.0.1:${webPort}`;
const e2eDbPath = join(mkdtempSync(join(tmpdir(), "schema-ui-e2e-")), "e2e.db");

// W24 (GOAL-035 / workspace-010): the browser E2E suite is a DUAL-DIALECT
// acceptance surface 閳?it must run once per store dialect, not merely "not
// get redirected" (that was W23 N-001: a gitignored apps/api/configs/.env
// with DB_DIALECT=postgres silently redirected the API to the developer's
// shared database and every fresh-seed login 401'd).
//
// The harness DECLARES the dialect contract:
//   - DB_DIALECT=sqlite (default): the API gets an isolated temp SQLite at
//     DB_PATH; configs/.env can never override an already-set process env.
//   - DB_DIALECT=postgres (explicit opt-in): the harness provisions a
//     DEDICATED scratch database via cmd/e2e-pgset (create 閳?run 閳?drop,
//     same pattern as internal/pgtest and the CI api-postgres job) and passes
//     its name as DB_NAME; credentials come from process env or configs/.env.
//
// e2e/global-setup.ts then VALIDATES the contract after boot (sqlite: DB file
// appeared; postgres: scratch schema_migrations exists) and fails fast with a
// diagnosis instead of letting specs die with cryptic 401s halfway through.
const dbDialect = (process.env.DB_DIALECT ?? "sqlite").trim().toLowerCase();
if (dbDialect !== "sqlite" && dbDialect !== "postgres") {
  throw new Error(`DB_DIALECT must be sqlite or postgres for browser E2E (got ${dbDialect || "empty"})`);
}
const apiDir = resolve(process.cwd(), "../api");

let e2ePgName: string | null = null;
if (dbDialect === "postgres") {
  // Playwright evaluates this config more than once per real run (main +
  // worker config serialization, in DIFFERENT processes). Guard on E2E_PG_NAME
  // so a later evaluation inherits the first one's scratch database instead
  // of provisioning a second one.
  if (process.env.E2E_PG_NAME !== undefined && process.env.E2E_PG_NAME !== "") {
    e2ePgName = process.env.E2E_PG_NAME;
  } else {
    e2ePgName = `schema_ui_e2e_${Date.now().toString(36)}${Math.random().toString(36).slice(2, 8)}`;
    const create = spawnSync("go", ["run", "./cmd/e2e-pgset", "create", e2ePgName], {
      cwd: apiDir,
      env: process.env,
      encoding: "utf8",
    });
    if (create.status !== 0) {
      throw new Error(
        `[e2e] postgres scratch database provisioning failed for ${e2ePgName}: ` +
          `${(create.stderr ?? create.stdout ?? "").trim() || "go run exited " + create.status}. ` +
          `Check DB_* credentials (env or apps/api/configs/.env) and CREATEDB rights.`,
      );
    }
    process.env.E2E_PG_NAME = e2ePgName;
  }
}

// Carried to e2e/global-setup.ts / global-teardown.ts (spawned workers inherit).
process.env.E2E_DB_DIALECT = dbDialect;
process.env.E2E_DB_PATH = e2eDbPath;
if (e2ePgName !== null) {
  process.env.E2E_PG_NAME = e2ePgName;
}

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
  globalSetup: "./e2e/global-setup",
  globalTeardown: "./e2e/global-teardown",
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
        // The harness contract wins over configs/.env (process env first).
        DB_DIALECT: dbDialect,
        ...(e2ePgName !== null ? { DB_NAME: e2ePgName } : {}),
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




