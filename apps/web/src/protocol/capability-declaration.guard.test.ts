// @vitest-environment node
//
// C-005 (GOAL-041 S2/S3) enforcement guard: the v2.9 capability pair
// (data.route-binding / form.controls.readonly) must never be declared on a
// page without actual usage — this is the drift class C-005 addressed
// (digitaloffer pages declared both without any $context.route.* / readOnly
// usage; user P-004 decision 2026-09-06 removed those declarations).
//
// Scope note: legacy capabilities (permissions.inheritance, actions.row.request,
// form.controls.extended, …) are declared conservatively across many pages;
// upstream page.schema.json only constrains usage→declaration (one-directional),
// so conservative declaration is legal there. A full declaration↔usage audit of
// the legacy set is a NON-BLOCKING follow-up (recorded in GOAL-041 E-005), NOT
// part of C-005's scope. This guard intentionally covers only the v2.9 pair.

import { existsSync, readdirSync, readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const WEB_SRC = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");

// v2.9 capability → usage-marker predicate over the raw document text.
const V29_MARKERS: Record<string, (text: string) => boolean> = {
  "data.route-binding": (text) => /\$context\.route\./.test(text),
  "form.controls.readonly": (text) => /"readOnly"\s*:\s*true/.test(text),
};

function walkSchemaFiles(dir: string, out: string[]): void {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const abs = join(dir, entry.name);
    if (entry.isDirectory()) {
      walkSchemaFiles(abs, out);
    } else if (entry.isFile() && /\\schema\\/.test(abs) && entry.name.endsWith(".json")) {
      out.push(abs);
    }
  }
}

function collectSchemaFiles(): string[] {
  const out: string[] = [];
  for (const moduleName of readdirSync(MODULES)) {
    const moduleRoot = join(MODULES, moduleName);
    if (!existsSync(moduleRoot)) {
      continue;
    }
    walkSchemaFiles(moduleRoot, out);
  }
  return out;
}

describe("C-005 · v2.9 capabilities are declared only when actually used", () => {
  const files = collectSchemaFiles();
  expect(files.length).toBeGreaterThanOrEqual(30);

  for (const file of files) {
    it(`${file.replace(/\\/g, "/").replace(WEB_SRC.replace(/\\/g, "/") + "/../", "")}`, () => {
      const text = readFileSync(file, "utf8");
      const document = JSON.parse(text) as {
        meta?: { pageId?: string; requiredCapabilities?: string[] };
      };
      const declared = document.meta?.requiredCapabilities ?? [];
      const violations: string[] = [];
      for (const capability of declared) {
        const marker = V29_MARKERS[capability];
        if (marker === undefined) {
          continue; // not part of the v2.9 pair (legacy set: out of scope)
        }
        if (!marker(text)) {
          violations.push(
            `declared capability "${capability}" has no usage marker in the document`,
          );
        }
      }
      expect(violations, `${document.meta?.pageId ?? file}: ${violations.join("; ")}`).toEqual([]);
    });
  }
});
