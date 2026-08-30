/**
 * W15 F-006 (GOAL-016 A-001): the canonical test-fixture root is
 * `apps/api/modules` — `apps/api/internal/modules` was retired and does not
 * exist on disk. This guard makes the root a CI-time fact:
 *   1. the canonical root exists relative to this package;
 *   2. no file under apps/web/src references the retired root any more.
 * A future test that hard-codes the old path fails here before it can fail
 * with a confusing ENOENT cascade.
 */

import { existsSync, readFileSync, readdirSync, statSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const SRC = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const CANONICAL_ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");
// This file is the sentinel itself: it must reference the retired path in its
// own description, so it is the one file exempted from the scan.
const SELF = fileURLToPath(import.meta.url).replaceAll("\\", "/");

function walk(dir: string, out: string[]): void {
  for (const entry of readdirSync(dir)) {
    const abs = join(dir, entry);
    if (statSync(abs).isDirectory()) {
      walk(abs, out);
    } else if (
      (abs.endsWith(".ts") || abs.endsWith(".tsx") || abs.endsWith(".js") || abs.endsWith(".jsx")) &&
      abs.replaceAll("\\", "/") !== SELF
    ) {
      out.push(abs);
    }
  }
}

describe("W15 F-006 · canonical fixture root", () => {
  it("apps/api/modules exists relative to web src", () => {
    expect(existsSync(CANONICAL_ROOT), `missing canonical fixture root ${CANONICAL_ROOT}`).toBe(true);
  });

  it("no source file references the retired apps/api/internal/modules root", () => {
    const files: string[] = [];
    walk(SRC, files);
    const offenders: string[] = [];
    for (const file of files) {
      const text = readFileSync(file, "utf8");
      if (text.includes("api/internal/modules")) {
        offenders.push(file.replaceAll("\\", "/").replace(SRC.replaceAll("\\", "/") + "/", ""));
      }
    }
    expect(offenders, "files still referencing apps/api/internal/modules").toEqual([]);
  });
});