import test from "node:test";
import assert from "node:assert/strict";
import { normalizePackageEsmImports } from "./package-esm-compat.mjs";

test("normalizes protocol self-subpaths and packaged JSON schema imports", () => {
  const output = normalizePackageEsmImports([
    'import nodeSchema from "@schemas/node.schema.json";',
    'import { validate } from "@schema-ui/protocol/protocol/conformance/runtime-schema-validate";',
    'import hostSupport from "./host-support.json";',
  ].join("\n"));

  assert.match(output, /@magicvr\/schema-ui-protocol\/schemas\/node\.schema\.json" with \{ type: "json" \}/);
  assert.match(output, /@magicvr\/schema-ui-protocol\/protocol\/conformance\/runtime-schema-validate\.js/);
  assert.match(output, /\.\/host-support\.json" with \{ type: "json" \}/);
});

test("rewrites renderer internal package imports to public .js subpaths", () => {
  const output = normalizePackageEsmImports([
    'import { Card } from "@schema-ui/ui/components/ui/card";',
    'import { cn } from "@schema-ui/lib/lib/utils";',
    'import { parse } from "@magicvr/schema-ui-protocol/protocol/parse";',
  ].join("\n"));

  assert.match(output, /@magicvr\/schema-ui-ui\/components\/ui\/card\.js/);
  assert.match(output, /@magicvr\/schema-ui-lib\/lib\/utils\.js/);
  assert.match(output, /@magicvr\/schema-ui-protocol\/protocol\/parse\.js/);
});

test("adds JSON import attributes and keeps finalization idempotent", () => {
  const source = [
    'import enUS from "./messages/en-US.json";',
    'export { default as zhCN } from "./messages/zh-CN.json";',
    'import alreadyValid from "./messages/fr-FR.json" with { type: "json" };',
  ].join("\n");
  const once = normalizePackageEsmImports(source);

  assert.match(once, /\.\/messages\/en-US\.json" with \{ type: "json" \}/);
  assert.match(once, /\.\/messages\/zh-CN\.json" with \{ type: "json" \}/);
  assert.match(once, /\.\/messages\/fr-FR\.json" with \{ type: "json" \}/);
  assert.equal(normalizePackageEsmImports(once), once);
});

test("leaves package roots and third-party imports unchanged", () => {
  const source = 'import React from "react"; import protocol from "@schema-ui/protocol";';
  assert.equal(
    normalizePackageEsmImports(source),
    'import React from "react"; import protocol from "@magicvr/schema-ui-protocol";',
  );
});
