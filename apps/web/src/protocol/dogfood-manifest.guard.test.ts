// @vitest-environment node
//
// C-006 (GOAL-041 S2): the static dogfood manifests are test fixtures, NOT
// runtime authorities. They hand-copy pageId/route/schemaUrl from the API
// module fragments, so they can drift. This guard re-derives the module union
// from `apps/api/modules/**/manifest/fragment.json` (recursive — nested
// modules like channel/telegram and dev/examples included) and asserts every
// dogfood page entry matches the union exactly (same pageId, route, schemaUrl).
// Renaming/deleting a page in a module without touching the dogfood fixture
// now fails the build instead of silently diverging.

import { existsSync, readFileSync, readdirSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const WEB_SRC = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");
const FIXTURES = resolve(WEB_SRC, "test-fixtures");

interface FragmentPage {
  pageId: string;
  title?: string;
  schemaUrl: string;
  route: string;
}

interface Fragment {
  protocolVersion?: string;
  pages?: FragmentPage[];
}

interface ManifestFixture {
  protocolVersion?: string;
  pages?: Array<{
    pageId?: string;
    schemaUrl?: string;
    route?: string;
  }>;
}

function walkFiles(dir: string, out: string[]): void {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const abs = join(dir, entry.name);
    if (entry.isDirectory()) {
      walkFiles(abs, out);
    } else if (entry.isFile()) {
      out.push(abs);
    }
  }
}

function modulePageUnion(): Map<string, { route: string; schemaUrl: string }> {
  const fragments: string[] = [];
  walkFiles(MODULES, fragments);
  const union = new Map<string, { route: string; schemaUrl: string }>();
  for (const file of fragments) {
    // Normalize separators so the suffix match works on Windows (backslash).
    const normalized = file.replace(/\\/g, "/");
    if (!normalized.endsWith("manifest/fragment.json")) {
      continue;
    }
    const raw = readFileSync(file, "utf8");
    if (raw.trim() === "") {
      continue;
    }
    const fragment = JSON.parse(raw) as Fragment;
    for (const page of fragment.pages ?? []) {
      union.set(page.pageId, { route: page.route, schemaUrl: page.schemaUrl });
    }
  }
  return union;
}

describe("C-006 · dogfood manifests stay consistent with the module union", () => {
  const union = modulePageUnion();

  it("derives a non-trivial page union from all module fragments", () => {
    expect(union.size).toBeGreaterThanOrEqual(30);
  });

  for (const name of ["app-manifest.mvp-dogfood.json", "app-manifest.admin-dogfood.json"]) {
    it(`${name} pages are a consistent subset of the module union`, () => {
      const path = resolve(FIXTURES, name);
      expect(existsSync(path), `${name} exists`).toBe(true);
      const manifest = JSON.parse(readFileSync(path, "utf8")) as ManifestFixture;
      const pages = manifest.pages ?? [];
      expect(pages.length).toBeGreaterThan(0);
      const problems: string[] = [];
      for (const page of pages) {
        if (typeof page.pageId !== "string" || typeof page.route !== "string" || typeof page.schemaUrl !== "string") {
          problems.push(`${name}: page missing pageId/route/schemaUrl: ${JSON.stringify(page)}`);
          continue;
        }
        const expected = union.get(page.pageId);
        if (expected === undefined) {
          problems.push(`${name}: pageId "${page.pageId}" not present in any module fragment`);
        } else {
          if (expected.route !== page.route) {
            problems.push(
              `${name}: "${page.pageId}" route mismatch (fixture "${page.route}" vs module "${expected.route}")`,
            );
          }
          if (expected.schemaUrl !== page.schemaUrl) {
            problems.push(
              `${name}: "${page.pageId}" schemaUrl mismatch (fixture "${page.schemaUrl}" vs module "${expected.schemaUrl}")`,
            );
          }
        }
      }
      expect(problems).toEqual([]);
    });
  }
});
