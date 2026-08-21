import { describe, expect, it } from "vitest";
import { existsSync, readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";

describe("theme FOUC bootstrap CSP compatibility (W8 F-002)", () => {
  const root = fileURLToPath(new URL("../../", import.meta.url));
  const indexHtml = readFileSync(`${root}index.html`, "utf8");
  const themeInit = `${root}public/theme-init.js`;

  it("references the external theme-init.js script", () => {
    expect(indexHtml).toContain('<script src="/theme-init.js"></script>');
  });

  it("does not contain an inline bootstrap script block", () => {
    // The old inline FOUC script was blocked by production CSP script-src 'self'.
    expect(indexHtml).not.toContain('localStorage.getItem("theme")');
  });

  it("serves the external bootstrap file from public", () => {
    expect(existsSync(themeInit)).toBe(true);
    const script = readFileSync(themeInit, "utf8");
    expect(script).toContain('localStorage.getItem("theme")');
    expect(script).toContain("prefers-color-scheme: dark");
  });
});