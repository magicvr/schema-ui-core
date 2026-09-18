/**
 * GOAL-008 (R6 A-002 F-005) · type-check evidence convention guard.
 *
 * WHY THIS EXISTS
 * ---------------
 * `apps/web/tsconfig.json` is a solution-style config: it declares no `files`
 * and no `include`, only `references` to `tsconfig.app.json` /
 * `tsconfig.node.json`. In non-build mode TypeScript compiles exactly what the
 * config selects, so running `tsc` against it compiles an EMPTY program and
 * exits 0 — it reports "success" because it checked nothing at all.
 *
 * That is not a config bug (solution-style + project references is the normal
 * Vite/TS layout, and `tsc -b` walks the references correctly). The bug was
 * choosing the wrong command: several stage records across this repo cited a
 * bare `tsc --noEmit` as "type check passed" when the command could never fail.
 * See GOAL-008 D-001 / E-001 for the injected-error proof
 * (`tsc --noEmit` → exit 0, `tsc -b` → TS2322, exit 2).
 *
 * WHAT THIS GUARD PINS
 * --------------------
 * 1. The root config stays solution-style (the precondition that makes the
 *    bare form vacuous — change it and you must revisit this convention);
 * 2. `npm run typecheck` exists and is a real checking form;
 * 3. no executable surface (package scripts, CI workflows, repo/app scripts)
 *    invokes `tsc` in a non-checking form;
 * 4. the README keeps documenting the convention.
 *
 * WHAT IT DELIBERATELY DOES NOT DO
 * --------------------------------
 * It does not scan governance prose. Historical stage records legitimately
 * quote the wrong command while explaining it, and history must not be
 * rewritten by a linter. Prose discipline is held by the README convention
 * plus the fact that the correct command is now the easiest one to run.
 */

import { readFileSync, readdirSync, statSync } from "node:fs";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const HERE = dirname(fileURLToPath(import.meta.url));
/** apps/web — this guard lives at src/, so the package root is one level up. */
const WEB_ROOT = resolve(HERE, "..");
/** Repository root (apps/web → apps → repo). */
const REPO_ROOT = resolve(WEB_ROOT, "../..");

const ROOT_TSCONFIG = join(WEB_ROOT, "tsconfig.json");
const PACKAGE_JSON = join(WEB_ROOT, "package.json");
const README = join(WEB_ROOT, "README.md");

interface TsConfigShape {
  files?: unknown;
  include?: unknown;
  references?: unknown;
}

/**
 * True when the config selects at least one source file of its own. A
 * solution-style config selects none, which is what makes a non-build `tsc`
 * run vacuous against it.
 */
function selectsOwnSources(config: TsConfigShape): boolean {
  const files = Array.isArray(config.files) ? config.files : [];
  if (files.length > 0) {
    return true;
  }
  const include = config.include;
  if (typeof include === "string") {
    return include.trim() !== "";
  }
  if (Array.isArray(include)) {
    return include.length > 0;
  }
  return false;
}

/**
 * A `tsc` invocation is a real check only when it walks project references
 * (`-b` / `--build`) or targets an explicit project (`-p` / `--project`).
 * Everything else — `tsc`, `tsc --noEmit`, `tsc --noEmit --pretty false` —
 * compiles only what the (empty) root config selects.
 *
 * The flag must be a standalone token; trailing shell/JSON punctuation
 * (quote, comma, semicolon, pipe, bracket) legitimately ends it.
 */
const CHECKING_FLAG = /(^|[\s"'(=[])(-b|--build|-p|--project)([\s"',;|&)\]]|$)/;

/** `tsc` as a command word — not `tsc-built`, not `@types/tsc-foo`. */
const TSC_COMMAND = /(?<![A-Za-z0-9_$.-])tsc(?![A-Za-z0-9_$-])/g;

/**
 * Context that puts `tsc` in command position: line start, a shell separator,
 * a package runner, or the opening of an argv element (`"tsc"`, `['tsc'`).
 * A bare mention inside prose or a log message ("… = tsc 全产物") is not an
 * invocation and must not be flagged.
 */
const PRECEDING_COMMAND =
  /(^|&&|\|\||;|\||npx(\.cmd)?|pnpm|yarn|bunx|npm\s+exec(\s+--)?|["'[(])\s*$/;

/**
 * Reports non-checking `tsc` usages on one line. Flags may sit anywhere after
 * the command (same string, or a later element of a spawn argument array), so
 * the whole remainder is inspected rather than a truncated segment.
 */
function nonCheckingTscOnLine(line: string): string[] {
  const offenders: string[] = [];
  TSC_COMMAND.lastIndex = 0;
  let match: RegExpExecArray | null;
  while ((match = TSC_COMMAND.exec(line)) !== null) {
    if (!PRECEDING_COMMAND.test(line.slice(0, match.index))) {
      continue; // prose mention, not an invocation
    }
    const remainder = line.slice(match.index + match[0].length);
    if (!CHECKING_FLAG.test(remainder)) {
      offenders.push(line.trim());
    }
  }
  return offenders;
}

/** Lines that only talk about a command (comments) are not invocations. */
function isCommentLine(line: string): boolean {
  const trimmed = line.trim();
  return (
    trimmed === "" ||
    trimmed.startsWith("//") ||
    trimmed.startsWith("#") ||
    trimmed.startsWith("*") ||
    trimmed.startsWith("/*")
  );
}

function walkFiles(dir: string, out: string[]): void {
  let entries: string[];
  try {
    entries = readdirSync(dir);
  } catch {
    return;
  }
  for (const entry of entries) {
    if (entry === "node_modules" || entry === "dist" || entry === "dist-lib") {
      continue;
    }
    const abs = join(dir, entry);
    if (statSync(abs).isDirectory()) {
      walkFiles(abs, out);
    } else {
      out.push(abs);
    }
  }
}

/** Executable surfaces that could carry a verification command. */
function executableSurfaces(): string[] {
  const surfaces = [PACKAGE_JSON];

  const workflows = join(REPO_ROOT, ".github/workflows");
  for (const entry of readdirSync(workflows)) {
    if (entry.endsWith(".yml") || entry.endsWith(".yaml")) {
      surfaces.push(join(workflows, entry));
    }
  }

  for (const scriptsDir of [join(REPO_ROOT, "scripts"), join(WEB_ROOT, "scripts")]) {
    const files: string[] = [];
    walkFiles(scriptsDir, files);
    surfaces.push(...files.filter((file) => /\.(mjs|cjs|js|ts|sh|ps1)$/.test(file)));
  }

  return surfaces;
}

describe("GOAL-008 · type-check evidence convention", () => {
  const rootConfig = JSON.parse(readFileSync(ROOT_TSCONFIG, "utf8")) as TsConfigShape;
  const pkg = JSON.parse(readFileSync(PACKAGE_JSON, "utf8")) as {
    scripts?: Record<string, string>;
  };
  const scripts = pkg.scripts ?? {};
  const rootSelectsSources = selectsOwnSources(rootConfig);

  it("keeps the root tsconfig solution-style (the documented precondition)", () => {
    // If this ever fails, the package gained its own source selection and the
    // "bare tsc is vacuous" premise no longer holds — revisit GOAL-008 D-001
    // and this guard instead of deleting the assertion.
    expect(
      rootSelectsSources,
      "apps/web/tsconfig.json now selects sources; revisit the type-check convention",
    ).toBe(false);
    expect(Array.isArray(rootConfig.references)).toBe(true);
    expect((rootConfig.references as unknown[]).length).toBeGreaterThan(0);
  });

  it("exposes a typecheck script that actually checks", () => {
    const typecheck = scripts.typecheck;
    expect(typecheck, "package.json is missing a typecheck script").toBeTypeOf("string");
    expect(
      CHECKING_FLAG.test(typecheck),
      `typecheck must use -b or -p (got: ${typecheck})`,
    ).toBe(true);
  });

  // The e2e specs live in their own project (e2e/tsconfig.json) which is NOT in
  // the root config's `references` graph, so `tsc -b` alone never checks them —
  // the same vacuity one level down. `typecheck` must reach them explicitly.
  it("type-checks the e2e project too (it is outside the references graph)", () => {
    const typecheck = scripts.typecheck ?? "";
    expect(
      /-p\s+e2e\/tsconfig\.json|--project\s+e2e\/tsconfig\.json/.test(typecheck),
      `typecheck must also check the e2e project (got: ${typecheck})`,
    ).toBe(true);

    const e2eConfigPath = join(WEB_ROOT, "e2e/tsconfig.json");
    const e2eConfig = JSON.parse(readFileSync(e2eConfigPath, "utf8")) as TsConfigShape;
    // If e2e ever joins the references graph this assertion is the reminder to
    // simplify the script rather than check it twice.
    const referenced = (rootConfig.references as Array<{ path?: string }> | undefined) ?? [];
    const inGraph = referenced.some((ref) => (ref.path ?? "").includes("e2e"));
    expect(inGraph, "e2e project unexpectedly joined the references graph").toBe(false);
    expect(selectsOwnSources(e2eConfig)).toBe(true);
  });

  it("keeps the build script type-checked", () => {
    const build = scripts.build ?? "";
    expect(build).toContain("tsc");
    expect(
      CHECKING_FLAG.test(build),
      `build must run a checking tsc form (got: ${build})`,
    ).toBe(true);
  });

  it("uses no non-checking tsc invocation in any executable surface", () => {
    const offenders: string[] = [];
    for (const surface of executableSurfaces()) {
      const lines = readFileSync(surface, "utf8").split(/\r?\n/);
      lines.forEach((line, index) => {
        if (isCommentLine(line)) {
          return;
        }
        for (const bad of nonCheckingTscOnLine(line)) {
          offenders.push(
            `${surface.replaceAll("\\", "/").replace(REPO_ROOT.replaceAll("\\", "/") + "/", "")}:${index + 1}: ${bad}`,
          );
        }
      });
    }
    expect(
      offenders,
      "non-checking tsc invocation(s) found — use `tsc -b` or `tsc -p <config>`",
    ).toEqual([]);
  });

  it("keeps the convention documented in the web README", () => {
    const readme = readFileSync(README, "utf8");
    expect(readme).toContain("npm run typecheck");
    expect(readme).toContain("tsc -b");
    // The warning itself is the teaching moment; losing it is a regression.
    expect(readme).toContain("tsc --noEmit");
  });
});
