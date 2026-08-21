import { describe, expect, it } from "vitest";

import {
  DEFAULT_LOCALE,
  isSupportedLocale,
  normalizeLocaleCandidate,
  normalizePreference,
  resolveLocale,
} from "./locale";

// ── normalizeLocaleCandidate ─────────────────────────────────────────────────

describe("normalizeLocaleCandidate", () => {
  it("accepts exact supported tags", () => {
    expect(normalizeLocaleCandidate("zh-CN")).toBe("zh-CN");
    expect(normalizeLocaleCandidate("en-US")).toBe("en-US");
  });

  it("accepts case and separator variants", () => {
    expect(normalizeLocaleCandidate("zh-cn")).toBe("zh-CN");
    expect(normalizeLocaleCandidate("zh_CN")).toBe("zh-CN");
    expect(normalizeLocaleCandidate("en-us")).toBe("en-US");
    expect(normalizeLocaleCandidate("EN")).toBe("en-US");
  });

  it("maps language-only prefixes to the canonical supported locale", () => {
    expect(normalizeLocaleCandidate("zh")).toBe("zh-CN");
    expect(normalizeLocaleCandidate("en")).toBe("en-US");
    expect(normalizeLocaleCandidate("en-GB")).toBe("en-US");
    expect(normalizeLocaleCandidate("en-AU")).toBe("en-US");
  });

  it("rejects unsupported and empty candidates (incl. auto)", () => {
    expect(normalizeLocaleCandidate("auto")).toBeNull();
    expect(normalizeLocaleCandidate("fr-FR")).toBeNull();
    expect(normalizeLocaleCandidate("ja")).toBeNull();
    expect(normalizeLocaleCandidate("")).toBeNull();
    expect(normalizeLocaleCandidate(null)).toBeNull();
    expect(normalizeLocaleCandidate(undefined)).toBeNull();
  });
});

// ── isSupportedLocale ────────────────────────────────────────────────────────

describe("isSupportedLocale", () => {
  it("accepts only the exact supported tags", () => {
    expect(isSupportedLocale("zh-CN")).toBe(true);
    expect(isSupportedLocale("en-US")).toBe(true);
    expect(isSupportedLocale("zh-cn")).toBe(false);
    expect(isSupportedLocale("fr")).toBe(false);
    expect(isSupportedLocale(null)).toBe(false);
  });
});

// ── resolveLocale: frozen priority ───────────────────────────────────────────
// user explicit → system default (non-auto) → browser preference → en-US

describe("resolveLocale priority (D-002 §I-L10N-002)", () => {
  const browser = ["fr-FR", "en-US", "zh-CN"] as const;

  it("explicit user choice beats system default and browser", () => {
    expect(
      resolveLocale({ stored: "zh-CN", systemDefault: "en-US", browserLanguages: browser }),
    ).toBe("zh-CN");
    expect(
      resolveLocale({ stored: "en-US", systemDefault: "zh-CN", browserLanguages: browser }),
    ).toBe("en-US");
  });

  it("system default (non-auto) beats browser preference", () => {
    expect(
      resolveLocale({ stored: null, systemDefault: "zh-CN", browserLanguages: browser }),
    ).toBe("zh-CN");
    expect(
      resolveLocale({ stored: null, systemDefault: "en-US", browserLanguages: browser }),
    ).toBe("en-US");
  });

  it("auto system default defers to browser preference", () => {
    expect(
      resolveLocale({ stored: null, systemDefault: "auto", browserLanguages: ["zh-CN", "en-US"] }),
    ).toBe("zh-CN");
  });

  it("browser preference picks the first supported language", () => {
    expect(
      resolveLocale({ stored: null, systemDefault: null, browserLanguages: ["fr-FR", "zh-CN", "en-US"] }),
    ).toBe("zh-CN");
    expect(
      resolveLocale({ stored: null, systemDefault: null, browserLanguages: ["fr-FR", "en-GB"] }),
    ).toBe("en-US");
  });

  it("falls back to en-US when nothing matches", () => {
    expect(resolveLocale({ stored: null, systemDefault: null, browserLanguages: [] })).toBe(
      DEFAULT_LOCALE,
    );
    expect(
      resolveLocale({ stored: "fr-FR", systemDefault: "ja-JP", browserLanguages: ["de-DE"] }),
    ).toBe(DEFAULT_LOCALE);
  });

  it("applies the same priority pre- and post-login (no auth dependency)", () => {
    // The resolver has no auth input at all; the same function serves both.
    const input = { stored: "zh-CN", systemDefault: "en-US", browserLanguages: ["en-US"] };
    expect(resolveLocale(input)).toBe("zh-CN");
  });
});

// ── normalizePreference ──────────────────────────────────────────────────────

describe("normalizePreference", () => {
  it("maps supported values to themselves and everything else to auto", () => {
    expect(normalizePreference("zh-CN")).toBe("zh-CN");
    expect(normalizePreference("en-US")).toBe("en-US");
    expect(normalizePreference(null)).toBe("auto");
    expect(normalizePreference(undefined)).toBe("auto");
    expect(normalizePreference("bogus")).toBe("auto");
    expect(normalizePreference("")).toBe("auto");
  });
});
