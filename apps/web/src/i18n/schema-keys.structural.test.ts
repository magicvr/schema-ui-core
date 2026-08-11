/**
 * S2 structural completeness (GOAL-003 · C3/F-V029 denominator).
 *
 * Asserts over the REAL shipped schema documents + catalogs:
 * 1. Every user-visible text prop (label/text/submitLabel/confirm/placeholder)
 *    in the frozen 12-page denominator has a *Key counterpart.
 * 2. Every *Key used in the denominator exists in BOTH catalogs (zh-CN + en-US)
 *    — the maintainable translation surface is total for the denominator.
 * 3. The en-US catalog is the canonical baseline (identical key sets).
 */

import { readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const __dir = dirname(fileURLToPath(import.meta.url));
const MODULES = resolve(__dir, "../../../api/internal/modules");

const SCHEMA_FILES = [
  "dev/examples/schema/overview.json",
  "dev/examples/schema/admin-list-batch.json",
  "dev/examples/schema/data-display.json",
  "dev/examples/schema/data-table.json",
  "dev/examples/schema/search-form-table.json",
  "dev/examples/schema/form-controls.json",
  "dev/examples/schema/form-with-reactions.json",
  "dev/examples/schema/form-with-upload.json",
  "users/schema/users.json",
  "roles/schema/roles.json",
  "settings/schema/settings.json",
  "activity/schema/activity.json",
];

const TEXT_PROPS = ["label", "text", "content", "submitLabel", "confirm", "placeholder"];
const KEY_PROPS = ["labelKey", "textKey", "contentKey", "submitLabelKey", "confirmKey", "placeholderKey"];

function loadCatalogs() {
  const en = JSON.parse(readFileSync(join(__dir, "messages/en-US.json"), "utf8")) as Record<string, string>;
  const zh = JSON.parse(readFileSync(join(__dir, "messages/zh-CN.json"), "utf8")) as Record<string, string>;
  return { en, zh };
}

function collectTextAndKeys(doc: unknown): Array<{ path: string; text: string; key?: string }> {
  const out: Array<{ path: string; text: string; key?: string }> = [];
  const visit = (node: unknown, p: string) => {
    if (node === null || node === undefined) return;
    if (Array.isArray(node)) {
      node.forEach((c, i) => visit(c, `${p}[${i}]`));
      return;
    }
    if (typeof node === "object") {
      const record = node as Record<string, unknown>;
      for (const [k, v] of Object.entries(record)) {
        const childPath = `${p}.${k}`;
        visit(v, childPath);
        if (typeof v === "string" && TEXT_PROPS.includes(k)) {
          const key = KEY_PROPS.map((kp) => record[kp]).find((entry): entry is string => typeof entry === "string" && entry !== "");
          out.push({ path: childPath, text: v, key });
        }
      }
    }
  };
  visit(doc, "root");
  return out;
}

describe("S2 · F-V029 denominator schema key completeness", () => {
  const { en, zh } = loadCatalogs();

  it("every user-visible text in the 12 schema documents has a *Key", () => {
    const missing: string[] = [];
    for (const rel of SCHEMA_FILES) {
      const doc = JSON.parse(readFileSync(resolve(MODULES, rel), "utf8"));
      for (const entry of collectTextAndKeys(doc)) {
        if (entry.key === undefined) {
          missing.push(`${rel} ${entry.path} = "${entry.text}"`);
        }
      }
    }
    expect(missing).toEqual([]);
  });

  it("every *Key used in the denominator exists in both catalogs", () => {
    const unknown: string[] = [];
    for (const rel of SCHEMA_FILES) {
      const doc = JSON.parse(readFileSync(resolve(MODULES, rel), "utf8"));
      for (const entry of collectTextAndKeys(doc)) {
        if (entry.key === undefined) continue;
        if (!(entry.key in en)) unknown.push(`${rel} ${entry.key} missing in en-US`);
        if (!(entry.key in zh)) unknown.push(`${rel} ${entry.key} missing in zh-CN`);
      }
    }
    expect(unknown).toEqual([]);
  });

  it("manifest titleKey/labelKey used by the served manifests exist in both catalogs", () => {
    const manifestFiles = [
      "apps/api/internal/manifest/app-manifest.json",
      "apps/api/internal/modules/users/manifest/fragment.json",
      "apps/api/internal/modules/roles/manifest/fragment.json",
      "apps/api/internal/modules/settings/manifest/fragment.json",
      "apps/api/internal/modules/activity/manifest/fragment.json",
    ];
    const unknown: string[] = [];
    const keys: string[] = [];
    const collectKeys = (node: unknown) => {
      if (node === null || typeof node !== "object") return;
      if (Array.isArray(node)) {
        node.forEach(collectKeys);
        return;
      }
      const record = node as Record<string, unknown>;
      if (typeof record.titleKey === "string") keys.push(record.titleKey);
      if (typeof record.labelKey === "string") keys.push(record.labelKey);
      Object.values(record).forEach(collectKeys);
    };
    for (const rel of manifestFiles) {
      const doc = JSON.parse(readFileSync(resolve(__dir, "../../../..", rel), "utf8"));
      collectKeys(doc);
    }
    for (const key of keys) {
      if (!(key in en)) unknown.push(`${key} missing in en-US`);
      if (!(key in zh)) unknown.push(`${key} missing in zh-CN`);
    }
    expect(unknown).toEqual([]);
    expect(keys.length).toBeGreaterThan(0);
  });

  it("catalogs have identical key sets (zh-CN complete vs en-US baseline)", () => {
    expect(Object.keys(zh).sort()).toEqual(Object.keys(en).sort());
  });
});
