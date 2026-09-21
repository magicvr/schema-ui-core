import test from "node:test";
import assert from "node:assert/strict";
import { npmjsPackageName } from "./npmjs-package-name.mjs";

test("maps internal short package names to the public schema-ui scope", () => {
  assert.equal(npmjsPackageName("@schema-ui/lib"), "@magicvr/schema-ui-lib");
  assert.equal(npmjsPackageName("@schema-ui/protocol"), "@magicvr/schema-ui-protocol");
});

test("preserves already-public names and normalizes pnpm tarball names", () => {
  assert.equal(npmjsPackageName("@magicvr/schema-ui-shell"), "@magicvr/schema-ui-shell");
  assert.equal(npmjsPackageName("schema-ui-renderer-0.3.14"), "@magicvr/schema-ui-renderer");
  assert.equal(npmjsPackageName("magicvr-schema-ui-theme-0.1.4"), "@magicvr/schema-ui-theme");
});

test("supports an explicitly selected publishing scope", () => {
  assert.equal(npmjsPackageName("@schema-ui/ui", "@acme"), "@acme/schema-ui-ui");
});

test("rejects a name outside the six published packages", () => {
  assert.throws(() => npmjsPackageName("@schema-ui/unknown"), /unsupported schema-ui package name/);
});
