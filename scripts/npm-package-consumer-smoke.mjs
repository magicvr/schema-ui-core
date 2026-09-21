#!/usr/bin/env node
/**
 * Verify the six built package faces from an isolated consumer node_modules.
 * Run after build-lib-packages.mjs and rewrite-lib-aliases.mjs.
 */
import { cpSync, existsSync, mkdtempSync, mkdirSync, readFileSync, readdirSync, rmSync, writeFileSync } from "node:fs";
import { execFileSync } from "node:child_process";
import os from "node:os";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { npmjsPackageName } from "./npmjs-package-name.mjs";

const here = path.dirname(fileURLToPath(import.meta.url));
const root = path.resolve(here, "..");
const webRoot = path.join(root, "apps/web");
const distRoot = path.join(webRoot, "dist-lib/@schema-ui");
const consumerRoot = mkdtempSync(path.join(os.tmpdir(), "schema-ui-npm-consumer-"));

try {
  const stageRoot = path.join(consumerRoot, "stage");
  const packedRoot = path.join(consumerRoot, "packed");
  mkdirSync(stageRoot);
  mkdirSync(packedRoot);

  for (const name of ["protocol", "lib", "renderer", "ui", "shell", "theme"]) {
    const source = path.join(distRoot, name);
    if (!existsSync(path.join(source, "package.json"))) {
      throw new Error(`missing built package: ${source}`);
    }
    const packageDir = path.join(stageRoot, name);
    cpSync(source, packageDir, { recursive: true });
    const manifestPath = path.join(packageDir, "package.json");
    const manifest = JSON.parse(readFileSync(manifestPath, "utf8"));
    manifest.name = npmjsPackageName(manifest.name);
    manifest.publishConfig = { access: "public" };
    writeFileSync(manifestPath, JSON.stringify(manifest, null, 2));
    execFileSync("npm", ["pack", "--silent", "--pack-destination", packedRoot], {
      cwd: packageDir,
      stdio: "ignore",
      shell: process.platform === "win32",
    });
  }

  const tarballs = readdirSync(packedRoot)
    .filter((name) => name.endsWith(".tgz"))
    .map((name) => path.join(packedRoot, name));
  if (tarballs.length !== 6) throw new Error(`expected six packed npm packages, got ${tarballs.length}`);
  writeFileSync(path.join(consumerRoot, "package.json"), JSON.stringify({ private: true, type: "module" }));
  execFileSync(
    "npm",
    ["install", "--no-audit", "--no-fund", ...tarballs, "react@19", "react-dom@19"],
    { cwd: consumerRoot, stdio: "inherit", shell: process.platform === "win32" },
  );

  const smokePath = path.join(consumerRoot, "smoke.mjs");
  writeFileSync(
    smokePath,
    `const imports = [
  "@magicvr/schema-ui-protocol",
  "@magicvr/schema-ui-lib",
  "@magicvr/schema-ui-renderer",
  "@magicvr/schema-ui-ui",
  "@magicvr/schema-ui-shell",
  "@magicvr/schema-ui-theme",
  "@magicvr/schema-ui-protocol/protocol/app-manifest.js",
  "@magicvr/schema-ui-protocol/protocol/host-support.js",
  "@magicvr/schema-ui-protocol/protocol/conformance/runtime-schema-validate.js",
  "@magicvr/schema-ui-protocol/schemas/node.schema.json",
  "@magicvr/schema-ui-lib/lib/utils.js",
  "@magicvr/schema-ui-ui/lib/utils.js",
  "@magicvr/schema-ui-ui/components/ui/button.js",
];

for (const specifier of imports) {
  const attributes = specifier.endsWith(".json") ? { with: { type: "json" } } : undefined;
  await import(specifier, attributes);
  console.log("PASS " + specifier);
}
`,
  );
  execFileSync(process.execPath, [smokePath], { cwd: consumerRoot, stdio: "inherit" });
} finally {
  const resolvedConsumerRoot = path.resolve(consumerRoot);
  if (
    path.dirname(resolvedConsumerRoot) !== path.resolve(os.tmpdir()) ||
    !path.basename(resolvedConsumerRoot).startsWith("schema-ui-npm-consumer-")
  ) {
    throw new Error(`refusing to remove unexpected consumer path: ${resolvedConsumerRoot}`);
  }
  rmSync(resolvedConsumerRoot, { recursive: true, force: true });
}
