import { spawnSync } from "node:child_process";
import { resolve } from "node:path";

// W24 (GOAL-035): drop the provisioned postgres scratch database after the run
// (best-effort; a leftover can be removed with `go run ./cmd/e2e-pgset drop <name>`).
// No-op for sqlite (the temp directory is OS-managed).
const apiDir = resolve(process.cwd(), "../api");

export default async function globalTeardown(): Promise<void> {
  const name = process.env.E2E_PG_NAME;
  if (!name) {
    return;
  }
  // Stdio inherited so a failed drop is VISIBLE (leftovers can then be cleaned
  // with `go run ./cmd/e2e-pgset list`/`drop`); WITH (FORCE) terminates the
  // API's still-open connections so the drop cannot be silently blocked.
  spawnSync("go", ["run", "./cmd/e2e-pgset", "drop", name], {
    cwd: apiDir,
    env: process.env,
    stdio: "inherit",
    encoding: "utf8",
  });
}