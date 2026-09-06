// @vitest-environment node
//
// F1 (GOAL-042 W30): full capability declaration↔usage audit guard. Every
// page schema's declared `meta.requiredCapabilities` must correspond to actual
// usage in the document, except the envelope capabilities (app.manifest /
// app.navigation) and the host-level ones (host.*). This extends the C-005
// guard (GOAL-041 S2/S3, which covered only the v2.9 pair) to the whole
// non-exempt capability set: a declared-but-unused capability is either an
// aspirational lie or a drift risk.
//
// Markers are intentionally textual (scan the raw JSON) so the guard stays
// robust to schema shape changes; false negatives are impossible (a marker
// missing for a genuinely used capability fails the test and forces an update
// of this mapping table). Pages whose usage happens inside a custom component
// (not textually visible) must register intent in the INTENT_OVERRIDES map
// below instead of silently keeping a dangling declaration.

import { existsSync, readdirSync, readFileSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const WEB_SRC = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");

// Envelope / host-level capabilities that may be declared without a
// document-local usage marker.
const EXEMPT = new Set([
  "app.manifest",
  "app.navigation",
  "host.bootstrap",
  "host.failure-recovery",
  "host.conformance-claim",
]);

// Capability → usage-marker predicate over the raw document text.
const MARKERS: Record<string, (text: string) => boolean> = {
  "data.route-binding": (text) => /\$context\.route\./.test(text),
  "form.controls.readonly": (text) => /"readOnly"\s*:\s*true/.test(text),
  "form.controls.extended": (text) =>
    /"type"\s*:\s*"(textarea|switch|checkbox|radio)"/.test(text) ||
    /"mode"\s*:\s*"multiple"/.test(text),
  "form.controls.advanced": (text) =>
    /"type"\s*:\s*"(cascader|checkboxGroup|richText|password)"/.test(text) ||
    /"defaultValue"\s*:/.test(text),
  "table.sort": (text) => /"sortable"|"sortField"|"defaultSort"/.test(text),
  "permissions.inheritance": (text) => /"permissionCascade"|"permissionIntent"/.test(text),
  "record.view.load": (text) => /"type"\s*:\s*"recordView"/.test(text),
  "form.record.load": (text) => /"recordSource"/.test(text),
  "table.selection": (text) => /"requiresSelection"/.test(text),
  "actions.upload": (text) => /"type"\s*:\s*"upload"/.test(text),
  "actions.row.request": (text) => /"requestMapping"/.test(text),
  "actions.page.trigger": (text) =>
    /"actionRef"/.test(text) || /"type"\s*:\s*"actionButton"/.test(text),
  "actions.row.navigate": (text) => /"navigateMapping"/.test(text),
  "actions.batch.request": (text) => /\/batch-delete/.test(text),
};

// Documented-intent overrides: pageId → capabilities whose usage lives inside
// a custom component / host surface and cannot be textually marked. Only pages
// that genuinely consume the capability outside the visible document may be
// listed here; each entry must name the consuming surface.
const INTENT_OVERRIDES: Record<string, string[]> = {
  // notifications: actions.page.trigger fires from the Host notification bell
  // (markAllRead) and the notification-center custom component.
  notifications: ["actions.page.trigger"],
};

function walkSchemaFiles(dir: string, out: string[]): void {
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const abs = join(dir, entry.name);
    if (entry.isDirectory()) {
      walkSchemaFiles(abs, out);
    } else if (entry.isFile() && abs.replace(/\\/g, "/").includes("/schema/") && entry.name.endsWith(".json")) {
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

describe("F1 · declared capabilities are actually used (or exempt / intent-registered)", () => {
  const files = collectSchemaFiles();
  expect(files.length).toBeGreaterThanOrEqual(30);

  for (const file of files) {
    it(`${file.replace(/\\/g, "/").replace(WEB_SRC.replace(/\\/g, "/") + "/../", "")}`, () => {
      const text = readFileSync(file, "utf8");
      const document = JSON.parse(text) as {
        meta?: { pageId?: string; requiredCapabilities?: string[] };
      };
      const pageId = document.meta?.pageId ?? file;
      const declared = document.meta?.requiredCapabilities ?? [];
      const violations: string[] = [];
      for (const capability of declared) {
        if (EXEMPT.has(capability)) {
          continue;
        }
        const marker = MARKERS[capability];
        if (marker === undefined) {
          violations.push(`no usage marker defined for capability "${capability}"`);
          continue;
        }
        if (marker(text)) {
          continue;
        }
        const intent = INTENT_OVERRIDES[pageId];
        if (intent !== undefined && intent.includes(capability)) {
          continue; // documented custom-surface consumption
        }
        violations.push(
          `declared capability "${capability}" has no usage marker in the document`,
        );
      }
      expect(violations, `${pageId}: ${violations.join("; ")}`).toEqual([]);
    });
  }
});

