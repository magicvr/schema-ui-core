// test:e2e:postgres — run the browser E2E suite against the postgres dialect.
//
// The harness provisions a dedicated scratch database per run (create → run →
// drop) using apps/api/cmd/e2e-pgset; connection details come from the process
// env or apps/api/configs/.env (DB_* keys), exactly like the API server itself.
// Requires a reachable PostgreSQL with CREATEDB rights for that user.
import { spawnSync } from "node:child_process";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const cli = resolve(
  here,
  "..",
  "node_modules",
  ".bin",
  process.platform === "win32" ? "playwright.cmd" : "playwright",
);
const result = spawnSync(cli, ["test", ...process.argv.slice(2)], {
  stdio: "inherit",
  shell: process.platform === "win32",
  env: { ...process.env, DB_DIALECT: "postgres" },
});
process.exit(result.status ?? 1);