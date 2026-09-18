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
 *    invokes `tsc` in a non-checking form — where "checking" means `-b`, or a
 *    `-p` target whose config selects sources of its own (GOAL-010);
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
import { dirname, join, relative, resolve } from "node:path";
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
 * (`-b` / `--build`) or targets an explicit project (`-p` / `--project`) whose
 * config selects sources of its own.
 *
 * Targeting a project is NOT sufficient by itself (GOAL-010, hardening
 * GOAL-008 A-002 F-001): `tsc --noEmit -p tsconfig.json` names the root
 * solution-style config, which selects no files, so it compiles an empty
 * program and exits 0 — the same vacuity in a `-p` wrapper. The `-p` target is
 * therefore resolved and inspected; an unreadable or unparsable target fails
 * closed (treated as non-checking) rather than passing silently.
 *
 * Everything else — `tsc`, `tsc --noEmit`, `tsc --noEmit --pretty false`,
 * `tsc -p <solution-style config>` — is vacuous.
 */
/** Marks a shell command boundary (`&&`, `||`, `;`, `|`) between argv tokens. */
const BOUNDARY = "\u0000boundary";
/** Characters that end an argv token. */
const ARG_SEPARATOR = /[\s,[\](){}]/;
/** A quote preceded by one of these starts a value; otherwise it closes one. */
const OPENS_QUOTE = /[\s,=([{;:]/;

/**
 * Split one command remainder into argv-ish tokens. The remainder starts right
 * after the `tsc` token, i.e. mid-line, so quotes are resolved by position: a
 * quote that begins a value (after a separator) opens a quoted region, while a
 * quote glued to the previous character closes the surrounding literal and is
 * dropped. Backtick regions keep `${…}` interpolations — which may contain `"`,
 * `)` or `]` — in one piece. Unquoted shell separators become `BOUNDARY` tokens
 * so a following command's flags are never attributed to this one.
 */
function tokenizeArgs(segment: string): string[] {
  const tokens: string[] = [];
  let current = "";
  let quote: string | null = null;
  const flush = (): void => {
    if (current !== "") {
      tokens.push(current);
      current = "";
    }
  };
  for (let index = 0; index < segment.length; index += 1) {
    const char = segment[index];
    const previous = index === 0 ? undefined : segment[index - 1];
    if (quote !== null) {
      if (char === quote) {
        quote = null;
      } else {
        current += char;
      }
      continue;
    }
    if (char === '"' || char === "'" || char === "`") {
      if (previous !== undefined && OPENS_QUOTE.test(previous)) {
        quote = char;
      }
      continue;
    }
    if (char === "&" || char === "|" || char === ";") {
      flush();
      if (tokens[tokens.length - 1] !== BOUNDARY) {
        tokens.push(BOUNDARY);
      }
      continue;
    }
    if (ARG_SEPARATOR.test(char)) {
      flush();
      continue;
    }
    current += char;
  }
  flush();
  return tokens;
}

/** Drop punctuation left over from JSON/`=` wrapping (`x.json",`). */
const cleanArg = (value: string): string =>
  value.replace(/^["'`,;]+/, "").replace(/["'`,;]+$/, "");

/** `-p x` / `-p=x` / `--project x` / `--project=x` targets of one invocation. */
function projectTargetsFrom(tokens: string[]): string[] {
  const targets: string[] = [];
  tokens.forEach((token, index) => {
    if (token === "-p" || token === "--project") {
      const value = tokens[index + 1];
      if (value !== undefined && value !== "" && value !== BOUNDARY) {
        targets.push(cleanArg(value));
      }
      return;
    }
    if (token.startsWith("--project=") || token.startsWith("-p=")) {
      targets.push(cleanArg(token.slice(token.indexOf("=") + 1)));
    }
  });
  return targets;
}

/** Tokens of one `tsc` invocation, cut at the next shell separator. */
function tscCallTokens(remainder: string): string[] {
  const tokens = tokenizeArgs(remainder);
  const boundary = tokens.indexOf(BOUNDARY);
  return boundary === -1 ? tokens : tokens.slice(0, boundary);
}

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
 * A `-p` target counts as a real project only when its config selects its own
 * sources. Directories resolve the way `tsc` resolves them (`./tsconfig.json`).
 * Anything unreadable or unparsable returns false so the caller stays
 * fail-closed.
 */
function projectTargetSelectsSources(target: string): boolean {
  const cleaned = target.trim().replace(/^["'`]+|["'`]+$/g, "");
  if (cleaned === "") {
    return false;
  }
  const resolved = resolve(WEB_ROOT, cleaned);
  let configPath = resolved;
  try {
    if (statSync(resolved).isDirectory()) {
      configPath = join(resolved, "tsconfig.json");
    }
  } catch {
    return false; // does not exist — fail closed
  }
  try {
    return selectsOwnSources(JSON.parse(readFileSync(configPath, "utf8")) as TsConfigShape);
  } catch {
    return false; // not a readable JSON config — fail closed
  }
}

const escapeRegExp = (value: string): string => value.replace(/[.*+?^${}()|[\]\\]/g, "\\$&");
/** `${expr}` / `$VAR` placeholders, as they appear after regex escaping. */
const INTERPOLATION = /\\\$\\\{[^}]*\\\}|\\\$[A-Za-z_][A-Za-z0-9_]*/g;

/** Every tsconfig in the package (top level and nested), ignoring build output. */
function allTsConfigs(): string[] {
  const found: string[] = [];
  walkFiles(WEB_ROOT, found);
  return found.filter(
    (file) => /(^|[\\/])tsconfig[^\\/]*\.json$/.test(file) && !file.includes("node_modules"),
  );
}

/**
 * `scripts/build-lib-packages.mjs` legitimately builds one project per lib
 * package by interpolating the package name into the config path, so the target
 * is not a literal we can stat. Turn the interpolation into a wildcard, then
 * require that the pattern matches at least one config and that EVERY config it
 * can denote selects sources — fail closed otherwise.
 */
function dynamicTargetSelectsSources(target: string): boolean {
  const normalized = target.trim().replace(/^["'`]+|["'`]+$/g, "").replaceAll("\\", "/");
  const escaped = escapeRegExp(normalized);
  const wildcarded = escaped.replace(INTERPOLATION, "[^/]*");
  if (wildcarded === escaped) {
    return false; // no interpolation to resolve — stay fail-closed
  }
  const pattern = new RegExp(`^${wildcarded}$`);
  const matches = allTsConfigs()
    .map((file) => relative(WEB_ROOT, file).replaceAll("\\", "/"))
    .filter((relativePath) => pattern.test(relativePath));
  return (
    matches.length > 0 &&
    matches.every((relativePath) => projectTargetSelectsSources(relativePath))
  );
}

/** Resolve a `-p` target, literal or interpolated. */
function projectTargetIsRealCheck(target: string): boolean {
  return /[$`]/.test(target)
    ? dynamicTargetSelectsSources(target)
    : projectTargetSelectsSources(target);
}

/**
 * True when the command remainder after `tsc` describes a call that cannot fail
 * a type check. `-b` is sufficient; otherwise at least one `-p` target must be
 * a config that selects its own sources.
 */
function tscCallChecksNothing(remainder: string): boolean {
  const tokens = tscCallTokens(remainder);
  if (tokens.includes("-b") || tokens.includes("--build")) {
    return false;
  }
  return !projectTargetsFrom(tokens).some(projectTargetIsRealCheck);
}

/**
 * GOAL-010 · regression table. Every row is a command form that has either been
 * seen in this repo or is one flag away from the vacuous forms that were.
 * `vacuous: true` means the guard must flag it.
 */
const COMMAND_FORMS: Array<{ command: string; vacuous: boolean; note: string }> = [
  { command: "tsc --noEmit", vacuous: true, note: "bare: compiles the empty root program" },
  { command: "npx tsc --noEmit --pretty false", vacuous: true, note: "flags add no files" },
  { command: "tsc --noEmit -p tsconfig.json", vacuous: true, note: "GOAL-010: -p at the solution-style root" },
  { command: "tsc --noEmit --project tsconfig.json", vacuous: true, note: "long form of the same mistake" },
  { command: "tsc --noEmit --project=tsconfig.json", vacuous: true, note: "= form of the same mistake" },
  { command: "tsc --noEmit -p .", vacuous: true, note: "directory resolves to the root config" },
  { command: "tsc --noEmit -p tsconfig.missing.json", vacuous: true, note: "unreadable target fails closed" },
  { command: "tsc --noEmit -p", vacuous: true, note: "missing target fails closed" },
  { command: "npm exec -- tsc --noEmit -p tsconfig.json", vacuous: true, note: "runner prefix does not help" },
  { command: "tsc -p tsconfig.app.json --noEmit", vacuous: false, note: "app config selects src" },
  { command: "tsc -p e2e/tsconfig.json", vacuous: false, note: "e2e config selects ./**/*.ts" },
  { command: "tsc -b", vacuous: false, note: "build mode walks references" },
  { command: "tsc --build", vacuous: false, note: "long form" },
  { command: "tsc -b && tsc -p e2e/tsconfig.json", vacuous: false, note: "the typecheck script" },
  { command: "tsc --noEmit -p tsconfig.json && tsc -p e2e/tsconfig.json", vacuous: true, note: "first call still vacuous" },
  { command: "验证方式：tsc 全产物检查", vacuous: false, note: "prose mention is not an invocation" },
];

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
    if (tscCallChecksNothing(remainder)) {
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
      nonCheckingTscOnLine(typecheck ?? ""),
      `typecheck must run a checking tsc form (got: ${typecheck})`,
    ).toEqual([]);
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
      nonCheckingTscOnLine(build),
      `build must run a checking tsc form (got: ${build})`,
    ).toEqual([]);
  });

  // GOAL-010 · the `-p` target itself has to select sources. Without this the
  // guard accepts `tsc --noEmit -p tsconfig.json`, which compiles the empty root
  // program and exits 0 — GOAL-008 A-002 F-001.
  it("rejects every vacuous command form, including -p at the root config", () => {
    const wrong = COMMAND_FORMS.filter(
      (form) => nonCheckingTscOnLine(form.command).length > 0 !== form.vacuous,
    ).map((form) => `${form.command} -> expected vacuous=${form.vacuous} (${form.note})`);
    expect(wrong, "command forms judged incorrectly").toEqual([]);
    // The table must keep covering both directions.
    expect(COMMAND_FORMS.some((form) => form.vacuous)).toBe(true);
    expect(COMMAND_FORMS.some((form) => !form.vacuous)).toBe(true);
  });

  it("resolves -p targets by config content, not by the flag being present", () => {
    expect(projectTargetSelectsSources("tsconfig.json")).toBe(false);
    expect(projectTargetSelectsSources(".")).toBe(false);
    expect(projectTargetSelectsSources("./tsconfig.json")).toBe(false);
    expect(projectTargetSelectsSources("tsconfig.app.json")).toBe(true);
    expect(projectTargetSelectsSources("e2e/tsconfig.json")).toBe(true);
    expect(projectTargetSelectsSources("tsconfig.nope.json")).toBe(false);
    expect(projectTargetSelectsSources("")).toBe(false);
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
