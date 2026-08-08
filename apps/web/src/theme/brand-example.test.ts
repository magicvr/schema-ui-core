/**
 * Fork brand example structural check (S5 · GOAL-005).
 *
 * `brand.example.css` is documentation-as-code: it must only override CSS
 * custom properties that the base `index.css` Token system already declares
 * (D-002/D-003). This test parses both files' `:root`/`.dark` custom
 * property declarations and asserts the example is a strict subset of the
 * base — i.e. copying the example into a fork can never introduce an unknown
 * token name that the Renderer/components do not already consume.
 */

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

const HERE = dirname(fileURLToPath(import.meta.url));
const INDEX_CSS_PATH = resolve(HERE, "../index.css");
const BRAND_EXAMPLE_PATH = resolve(HERE, "./brand.example.css");

/** Extracts declared `--custom-property` names from a CSS source string. */
function declaredCustomProperties(css: string): Set<string> {
  const names = new Set<string>();
  const pattern = /(--[a-zA-Z0-9-]+)\s*:/g;
  let match: RegExpExecArray | null;
  while ((match = pattern.exec(css)) !== null) {
    names.add(match[1]);
  }
  return names;
}

describe("fork brand example (brand.example.css)", () => {
  const baseTokens = declaredCustomProperties(readFileSync(INDEX_CSS_PATH, "utf8"));
  const exampleTokens = declaredCustomProperties(readFileSync(BRAND_EXAMPLE_PATH, "utf8"));

  it("declares at least one override (not an empty/no-op example)", () => {
    expect(exampleTokens.size).toBeGreaterThan(0);
  });

  it("only overrides token names that the base index.css already declares", () => {
    const unknown = [...exampleTokens].filter((name) => !baseTokens.has(name));
    expect(unknown).toEqual([]);
  });

  it("overrides the brand-relevant token families (primary / chart / radius)", () => {
    const overriddenFamilies = new Set(
      [...exampleTokens].map((name) => name.replace(/-\d+$/, "")),
    );
    expect(overriddenFamilies.has("--primary")).toBe(true);
    expect(overriddenFamilies.has("--chart")).toBe(true);
    expect(overriddenFamilies.has("--radius")).toBe(true);
  });
});
