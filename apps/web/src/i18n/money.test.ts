import { describe, expect, it } from "vitest";

import { formatMoney, parseLocalizedMoney, parseLocalizedNumber, defaultCurrencyFor, normalizeCurrencyCode, resolveEffectiveCurrency, DEFAULT_CURRENCY } from "./money";

// ── formatMoney (C1 · §3.1) ───────────────────────────────────────────────────

describe("formatMoney", () => {
  it("formats zh-CN CNY from minor units (symbol/position/fraction via Intl)", () => {
    const out = formatMoney(12345, "zh-CN", { currency: "CNY" });
    expect(out).toContain("¥");
    expect(out).toContain("123.45");
  });

  it("formats en-US USD from minor units", () => {
    const out = formatMoney(12345, "en-US", { currency: "USD" });
    expect(out).toBe("$123.45");
  });

  it("uses the locale's embedded default currency when none is given (§4.3)", () => {
    expect(formatMoney(12345, "zh-CN")).toContain("¥");
    expect(formatMoney(12345, "en-US")).toBe("$123.45");
  });

  it("respects explicit non-default currencies", () => {
    const out = formatMoney(12345, "en-US", { currency: "CNY" });
    expect(out).toContain("¥"); // en-US CNY renders "CN¥123.45"
  });

  it("returns '' for invalid input (fail-safe, never throws)", () => {
    expect(formatMoney(Number.NaN, "en-US")).toBe("");
    expect(formatMoney(Number.POSITIVE_INFINITY, "en-US")).toBe("");
    expect(formatMoney(123, "en-US", { currency: "NOT-A-CURRENCY" })).toBe("");
    expect(formatMoney(123, "en-US", { currency: "" })).toBe("");
  });

  it("supports non-default minor units (incl. exponent 0 like JPY)", () => {
    // 12345 minor units of a 0-decimal currency → "¥12,345"
    expect(formatMoney(12345, "en-US", { currency: "JPY", minorUnits: 0 })).toBe("¥12,345");
    expect(parseLocalizedMoney("¥12,345", "en-US", { currency: "JPY", minorUnits: 0 })).toBe(12345);
  });

  it("R3 F-002: site default currency wins over the per-locale map, explicit wins over site", () => {
    // en-US with site default CNY → ¥ (site over locale map which would be USD)
    expect(formatMoney(12345, "en-US", { siteDefaultCurrency: "CNY" })).toContain("¥");
    // explicit currency always wins over the site default
    expect(formatMoney(12345, "en-US", { siteDefaultCurrency: "CNY", currency: "USD" })).toBe("$123.45");
    // unset site default → embedded map
    expect(formatMoney(12345, "en-US", { siteDefaultCurrency: "" })).toBe("$123.45");
    // invalid site default is ignored (per-locale map applies)
    expect(formatMoney(12345, "en-US", { siteDefaultCurrency: "nope" })).toBe("$123.45");
  });
});

// ── defaultCurrencyFor (C2 · §4.3) ────────────────────────────────────────────

describe("defaultCurrencyFor", () => {
  it("maps zh-CN → CNY and en-US → USD", () => {
    expect(defaultCurrencyFor("zh-CN")).toBe("CNY");
    expect(defaultCurrencyFor("en-US")).toBe("USD");
  });

  it("falls back to USD for unknown locales and never throws", () => {
    expect(defaultCurrencyFor("fr-FR")).toBe(DEFAULT_CURRENCY);
    expect(defaultCurrencyFor("")).toBe(DEFAULT_CURRENCY);
  });
});

// ── resolveEffectiveCurrency (R3 F-002 · site default channel) ───────────────

describe("resolveEffectiveCurrency", () => {
  it("site default wins when set; map applies otherwise", () => {
    expect(resolveEffectiveCurrency("en-US", "CNY")).toBe("CNY");
    expect(resolveEffectiveCurrency("zh-CN", "USD")).toBe("USD");
    expect(resolveEffectiveCurrency("en-US", null)).toBe("USD");
    expect(resolveEffectiveCurrency("en-US", undefined)).toBe("USD");
    expect(resolveEffectiveCurrency("en-US", "")).toBe("USD");
    expect(resolveEffectiveCurrency("zh-CN", "CNY")).toBe("CNY");
  });

  it("ignores invalid site values (falls back to the map)", () => {
    expect(resolveEffectiveCurrency("en-US", "nope")).toBe("USD");
  });
});

// ── normalizeCurrencyCode ─────────────────────────────────────────────────────

describe("normalizeCurrencyCode", () => {
  it("uppercases 3-letter codes and rejects the rest", () => {
    expect(normalizeCurrencyCode("cny")).toBe("CNY");
    expect(normalizeCurrencyCode(" USD ")).toBe("USD");
    expect(normalizeCurrencyCode("USDX")).toBeNull();
    expect(normalizeCurrencyCode("US")).toBeNull();
    expect(normalizeCurrencyCode("")).toBeNull();
    expect(normalizeCurrencyCode(null)).toBeNull();
  });
});

// ── parseLocalizedMoney (C3 · §3.2) ───────────────────────────────────────────

describe("parseLocalizedMoney", () => {
  it("parses en-US currency strings into minor units", () => {
    expect(parseLocalizedMoney("$1,234.56", "en-US", { currency: "USD" })).toBe(123456);
    expect(parseLocalizedMoney("$123.45", "en-US", { currency: "USD" })).toBe(12345);
    // site default currency participates in symbol stripping too — en-US
    // renders CNY with the "CN¥" prefix.
    expect(parseLocalizedMoney("CN¥123.45", "en-US", { siteDefaultCurrency: "CNY" })).toBe(12345);
  });

  it("parses zh-CN currency strings", () => {
    expect(parseLocalizedMoney("¥123.45", "zh-CN", { currency: "CNY" })).toBe(12345);
    expect(parseLocalizedMoney("¥12,345.67", "zh-CN", { currency: "CNY" })).toBe(1234567);
  });

  it("tolerates the bare ISO code and negative amounts", () => {
    expect(parseLocalizedMoney("CNY 123.45", "zh-CN", { currency: "CNY" })).toBe(12345);
    expect(parseLocalizedMoney("-$123.45", "en-US", { currency: "USD" })).toBe(-12345);
  });

  it("applies the embedded default currency when none is given", () => {
    expect(parseLocalizedMoney("¥123.45", "zh-CN")).toBe(12345);
    expect(parseLocalizedMoney("$123.45", "en-US")).toBe(12345);
  });

  it("returns null for unparsable input (never submits raw strings)", () => {
    expect(parseLocalizedMoney("", "en-US")).toBeNull();
    expect(parseLocalizedMoney("abc", "en-US")).toBeNull();
    expect(parseLocalizedMoney("$", "en-US")).toBeNull();
    expect(parseLocalizedMoney("$1.2.3", "en-US")).toBeNull();
  });

  it("documents the grouping-position tolerance (not validated in R3)", () => {
    // "12,34.5" — naive separator stripping accepts odd grouping; position
    // correctness is out of the R3 verification scope (see money.ts note).
    expect(parseLocalizedMoney("12,34.5", "en-US", { currency: "USD" })).toBe(123450);
  });

  it("supports non-default minor units", () => {
    // zh-CN renders JPY with the "JP¥" prefix — the round-trip drives the symbol.
    const rendered = formatMoney(1234, "zh-CN", { currency: "JPY", minorUnits: 0 });
    expect(rendered).toContain("JP¥");
    expect(parseLocalizedMoney(rendered, "zh-CN", { currency: "JPY", minorUnits: 0 })).toBe(1234);
  });
});

// ── parseLocalizedNumber (C3 · §3.2) ──────────────────────────────────────────

describe("parseLocalizedNumber", () => {
  it("normalizes locale grouping into machine numbers", () => {
    expect(parseLocalizedNumber("1,234.5", "en-US")).toBe(1234.5);
    expect(parseLocalizedNumber("1,234,567", "en-US")).toBe(1234567);
  });

  it("returns null for unparsable input", () => {
    expect(parseLocalizedNumber("", "en-US")).toBeNull();
    expect(parseLocalizedNumber("n/a", "en-US")).toBeNull();
    expect(parseLocalizedNumber("1.2.3", "en-US")).toBeNull();
  });
});

// ── C5 · bidirectional consistency (same contract both ways) ──────────────────

describe("format ↔ parse round-trip (C5 · §3 bidirectionality)", () => {
  it("en-US: format → parse returns the original minor units", () => {
    for (const minor of [0, 1, 12345, 123456, 99999999]) {
      const rendered = formatMoney(minor, "en-US", { currency: "USD" });
      expect(parseLocalizedMoney(rendered, "en-US", { currency: "USD" })).toBe(minor);
    }
  });

  it("zh-CN: format → parse returns the original minor units", () => {
    for (const minor of [0, 1, 12345, 123456, 99999999]) {
      const rendered = formatMoney(minor, "zh-CN", { currency: "CNY" });
      expect(parseLocalizedMoney(rendered, "zh-CN", { currency: "CNY" })).toBe(minor);
    }
  });
});