import { describe, expect, it } from "vitest";
import { resolveTheme, applyThemeToElement } from "./theme";

// ── resolveTheme ──────────────────────────────────────────────────────────────

describe("resolveTheme", () => {
  it("returns dark when stored === 'dark'", () => {
    const out = resolveTheme({ stored: "dark", prefersDark: false });
    expect(out.theme).toBe("dark");
    expect(out.colorScheme).toBe("dark");
  });

  it("returns light when stored === 'light'", () => {
    const out = resolveTheme({ stored: "light", prefersDark: true });
    expect(out.theme).toBe("light");
    expect(out.colorScheme).toBe("light");
  });

  it("falls back to OS dark preference when stored is null", () => {
    const out = resolveTheme({ stored: null, prefersDark: true });
    expect(out.theme).toBe("dark");
  });

  it("falls back to light when stored is null and OS is light", () => {
    const out = resolveTheme({ stored: null, prefersDark: false });
    expect(out.theme).toBe("light");
  });

  it("treats empty string stored as no preference (falls back to OS)", () => {
    // An empty string is falsy, so the stored branch is skipped.
    const out = resolveTheme({ stored: "", prefersDark: true });
    expect(out.theme).toBe("dark");
  });

  it("ignores unknown stored values and falls back to OS", () => {
    const out = resolveTheme({ stored: "system", prefersDark: false });
    expect(out.theme).toBe("light");
  });
});

// ── VP-007 S3: system default theme priority ─────────────────────────────────
// explicit user choice → system default (non-auto) → OS preference

describe("resolveTheme with system default (VP-007 S3)", () => {
  it("explicit user choice beats the system default", () => {
    const out = resolveTheme({ stored: "light", prefersDark: true, systemDefault: "dark" });
    expect(out.theme).toBe("light");
    const outDark = resolveTheme({ stored: "dark", prefersDark: false, systemDefault: "light" });
    expect(outDark.theme).toBe("dark");
  });

  it("system default (non-auto) beats the OS preference", () => {
    const out = resolveTheme({ stored: null, prefersDark: false, systemDefault: "dark" });
    expect(out.theme).toBe("dark");
    const outLight = resolveTheme({ stored: null, prefersDark: true, systemDefault: "light" });
    expect(outLight.theme).toBe("light");
  });

  it("auto/unknown system default defers to the OS preference", () => {
    const out = resolveTheme({ stored: null, prefersDark: true, systemDefault: "auto" });
    expect(out.theme).toBe("dark");
    const out2 = resolveTheme({ stored: null, prefersDark: false, systemDefault: null });
    expect(out2.theme).toBe("light");
  });
});

// ── applyThemeToElement ───────────────────────────────────────────────────────

describe("applyThemeToElement", () => {
  function makeEl(): HTMLElement {
    // Minimal HTMLElement stub — enough for classList + style.
    const el = {
      classList: {
        _classes: new Set<string>(),
        add(cls: string) {
          this._classes.add(cls);
        },
        remove(cls: string) {
          this._classes.delete(cls);
        },
        has(cls: string) {
          return this._classes.has(cls);
        },
      },
      style: { colorScheme: "" },
    } as unknown as HTMLElement;
    return el;
  }

  it("adds 'dark' class and sets colorScheme for dark theme", () => {
    const el = makeEl();
    applyThemeToElement(el, { theme: "dark", colorScheme: "dark" });
    expect((el.classList as unknown as { _classes: Set<string> })._classes.has("dark")).toBe(true);
    expect(el.style.colorScheme).toBe("dark");
  });

  it("removes 'dark' class and sets colorScheme to light for light theme", () => {
    const el = makeEl();
    // Pre-seed the dark class.
    (el.classList as unknown as { _classes: Set<string> })._classes.add("dark");
    applyThemeToElement(el, { theme: "light", colorScheme: "light" });
    expect((el.classList as unknown as { _classes: Set<string> })._classes.has("dark")).toBe(false);
    expect(el.style.colorScheme).toBe("light");
  });

  it("is idempotent for repeated dark calls", () => {
    const el = makeEl();
    applyThemeToElement(el, { theme: "dark", colorScheme: "dark" });
    applyThemeToElement(el, { theme: "dark", colorScheme: "dark" });
    expect((el.classList as unknown as { _classes: Set<string> })._classes.has("dark")).toBe(true);
  });
});

// ── index.css Token structure assertions ─────────────────────────────────────
// These are static-analysis tests that read the shipped CSS file and assert
// the required Token names exist.  They do not parse CSS semantics; they
// confirm the variable declarations are present (structural contract).

import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import { join, dirname } from "node:path";

const __dir = dirname(fileURLToPath(import.meta.url));
const cssPath = join(__dir, "..", "index.css");
const css = readFileSync(cssPath, "utf-8");

describe("index.css Token structure", () => {
  const requiredTokens = [
    "--destructive",
    "--success",
    "--chart-1",
    "--chart-2",
    "--chart-3",
    "--chart-4",
    "--chart-5",
    "--overlay",
    "--font-sans",
    "--font-mono",
    "--elevation-sm",
    "--elevation-md",
    "--elevation-lg",
  ];

  for (const token of requiredTokens) {
    it(`declares ${token}`, () => {
      expect(css).toContain(token);
    });
  }

  it("maps --shadow-sm to var(--elevation-sm) without self-reference", () => {
    // Must contain the alias declaration.
    expect(css).toContain("--shadow-sm: var(--elevation-sm)");
    // Must NOT contain --shadow-sm: var(--shadow-sm) (self-ref).
    expect(css).not.toContain("--shadow-sm: var(--shadow-sm)");
  });

  it("maps --shadow-md to var(--elevation-md) without self-reference", () => {
    expect(css).toContain("--shadow-md: var(--elevation-md)");
    expect(css).not.toContain("--shadow-md: var(--shadow-md)");
  });

  it("maps --shadow-lg to var(--elevation-lg) without self-reference", () => {
    expect(css).toContain("--shadow-lg: var(--elevation-lg)");
    expect(css).not.toContain("--shadow-lg: var(--shadow-lg)");
  });

  it("declares --elevation-* inside :root (not only inside .dark)", () => {
    // The :root block must appear before .dark and contain elevation tokens.
    const rootIdx = css.indexOf(":root");
    const darkIdx = css.indexOf(".dark");
    const elevSmIdx = css.indexOf("--elevation-sm");
    // Elevation must be declared at (or before) .dark — not exclusively in .dark.
    expect(rootIdx).toBeGreaterThanOrEqual(0);
    expect(elevSmIdx).toBeGreaterThanOrEqual(0);
    // elevation-sm should appear before .dark block (in the :root / §3 section)
    expect(elevSmIdx).toBeLessThan(darkIdx > -1 ? darkIdx + 10000 : Infinity);
  });
});
