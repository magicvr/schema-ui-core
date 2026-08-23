import { existsSync } from "node:fs";
import { spawnSync } from "node:child_process";
import { resolve } from "node:path";

// W24 (GOAL-035): store-contract validation with fail-fast diagnosis.
//
// Runs after the webServer entries are up and before any spec: verify the API
// actually operates on the store the harness declared (playwright.config.ts).
// Historically a gitignored apps/api/configs/.env silently redirected the API
// to the developer's shared Postgres — no sqlite file, no seed, every login
// 401 (W23 N-001). The suite must never reach the specs in that state.
const apiDir = resolve(process.cwd(), "../api");
const CONTRACT_TIMEOUT_MS = 60_000;

function delay(ms: number): Promise<void> {
  return new Promise((resolveDelay) => setTimeout(resolveDelay, ms));
}

function runHelper(args: string[]): { status: number; stderr: string } {
  const result = spawnSync("go", ["run", "./cmd/e2e-pgset", ...args], {
    cwd: apiDir,
    env: process.env,
    encoding: "utf8",
  });
  return { status: result.status ?? -1, stderr: (result.stderr ?? "").trim() };
}

export default async function globalSetup(): Promise<void> {
  const dialect = process.env.E2E_DB_DIALECT ?? "sqlite";
  const deadline = Date.now() + CONTRACT_TIMEOUT_MS;

  if (dialect === "postgres") {
    const name = process.env.E2E_PG_NAME;
    if (!name) {
      throw new Error(
        "[e2e] postgres dialect requires E2E_PG_NAME (set by playwright.config.ts provisioning); did the config load early?",
      );
    }
    for (;;) {
      const check = runHelper(["verify", name]);
      if (check.status === 0) {
        return;
      }
      if (Date.now() >= deadline) {
        throw new Error(
          `[e2e] store contract violated: the API never migrated scratch database "${name}" ` +
            `(e2e-pgset verify: ${check.stderr || "unknown"}). The postgres run must use the ` +
            `dedicated provisioned DB_NAME — a shared developer database would silently break fresh-seed ` +
            `logins (W23 N-001 root cause). See GOAL-035 D-001.`,
        );
      }
      await delay(2000);
    }
  }

  const dbPath = process.env.E2E_DB_PATH;
  if (!dbPath) {
    throw new Error("[e2e] sqlite dialect requires E2E_DB_PATH (set by playwright.config.ts)");
  }
  for (;;) {
    if (existsSync(dbPath)) {
      return;
    }
    if (Date.now() >= deadline) {
      throw new Error(
        `[e2e] store contract violated: no SQLite database file appeared at "${dbPath}" after API boot. ` +
          `The API was redirected away from the harness DB_PATH — check apps/api/configs/.env and ` +
          `CONFIG_ENV_FILE for a DB_DIALECT/DB_* override (W23 N-001 root cause: shared Postgres instead ` +
          `of the isolated store). See GOAL-035 D-001.`,
      );
    }
    await delay(2000);
  }
}