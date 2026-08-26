// @vitest-environment jsdom

import { describe, expect, it, vi } from "vitest";

import {
  AUTO_TIMEZONE,
  detectBrowserTimezone,
  isValidIanaTimeZone,
  normalizeTimezonePreference,
  readStoredTimezone,
  resolveEffectiveTimezone,
  TIMEZONE_STORAGE_KEY,
  writeStoredTimezone,
  type TimezoneResolutionInput,
} from "./timezone";

const noZone = () => "";

function resolve(partial: Partial<TimezoneResolutionInput>): string {
  return resolveEffectiveTimezone({
    stored: null,
    siteDefault: null,
    detect: noZone,
    ...partial,
  });
}

// ── isValidIanaTimeZone ───────────────────────────────────────────────────────

describe("isValidIanaTimeZone", () => {
  it("accepts real IANA names (incl. UTC)", () => {
    expect(isValidIanaTimeZone("Asia/Shanghai")).toBe(true);
    expect(isValidIanaTimeZone("America/New_York")).toBe(true);
    expect(isValidIanaTimeZone("UTC")).toBe(true);
  });

  it("rejects empty, non-string and unknown names", () => {
    expect(isValidIanaTimeZone("")).toBe(false);
    expect(isValidIanaTimeZone("   ")).toBe(false);
    expect(isValidIanaTimeZone(null)).toBe(false);
    expect(isValidIanaTimeZone(undefined)).toBe(false);
    expect(isValidIanaTimeZone("Foo/Bar")).toBe(false);
  });
});

// ── normalizeTimezonePreference ───────────────────────────────────────────────

describe("normalizeTimezonePreference", () => {
  it("maps invalid / empty / auto values to auto", () => {
    expect(normalizeTimezonePreference(null)).toBe(AUTO_TIMEZONE);
    expect(normalizeTimezonePreference(undefined)).toBe(AUTO_TIMEZONE);
    expect(normalizeTimezonePreference("")).toBe(AUTO_TIMEZONE);
    expect(normalizeTimezonePreference("auto")).toBe(AUTO_TIMEZONE);
    expect(normalizeTimezonePreference("AUTO")).toBe(AUTO_TIMEZONE);
    expect(normalizeTimezonePreference("Foo/Bar")).toBe(AUTO_TIMEZONE);
  });

  it("keeps valid IANA names (trimmed)", () => {
    expect(normalizeTimezonePreference("Asia/Shanghai")).toBe("Asia/Shanghai");
    expect(normalizeTimezonePreference("  Europe/London  ")).toBe("Europe/London");
  });
});

// ── resolveEffectiveTimezone · priority L1..L4 ────────────────────────────────

describe("resolveEffectiveTimezone", () => {
  it("L1: user override wins over detection and site default", () => {
    const out = resolve({
      stored: "Asia/Shanghai",
      siteDefault: "Europe/London",
      detect: () => "America/New_York",
    });
    expect(out).toBe("Asia/Shanghai");
  });

  it("L1: 'auto' stored means no override (falls through)", () => {
    const out = resolve({ stored: "auto", detect: () => "America/New_York" });
    expect(out).toBe("America/New_York");
  });

  it("L1: invalid stored value is ignored (falls through)", () => {
    const out = resolve({ stored: "Foo/Bar", detect: () => "America/New_York" });
    expect(out).toBe("America/New_York");
  });

  it("L2: session probe is used when no override", () => {
    const out = resolve({ detect: () => "America/New_York" });
    expect(out).toBe("America/New_York");
  });

  it("L2: empty or invalid probe is skipped", () => {
    expect(resolve({ detect: () => "" })).toBe(AUTO_TIMEZONE);
    expect(resolve({ detect: () => "Foo/Bar" })).toBe(AUTO_TIMEZONE);
  });

  it("L3: site default is used when probe is unavailable", () => {
    const out = resolve({ detect: noZone, siteDefault: "Europe/London" });
    expect(out).toBe("Europe/London");
  });

  it("L3: 'auto' / empty site default is skipped", () => {
    expect(resolve({ detect: noZone, siteDefault: "auto" })).toBe(AUTO_TIMEZONE);
    expect(resolve({ detect: noZone, siteDefault: "" })).toBe(AUTO_TIMEZONE);
  });

  it("L4: embedded default 'auto' when everything is unset", () => {
    expect(resolve({})).toBe(AUTO_TIMEZONE);
  });

  it("never throws on hostile inputs", () => {
    expect(() =>
      resolve({ stored: "Foo/Bar", siteDefault: "Foo/Baz", detect: () => "Foo/Qux" }),
    ).not.toThrow();
  });
});

// ── detectBrowserTimezone ─────────────────────────────────────────────────────

describe("detectBrowserTimezone", () => {
  it("returns a valid IANA name or empty string — never throws", () => {
    expect(() => detectBrowserTimezone()).not.toThrow();
    const zone = detectBrowserTimezone();
    expect(typeof zone).toBe("string");
    if (zone !== "") {
      expect(isValidIanaTimeZone(zone)).toBe(true);
    }
  });
});

// ── read / write stored preference ────────────────────────────────────────────

describe("timezone storage", () => {
  it("round-trips via the single channel and auto removes the key", () => {
    const get = vi.spyOn(Storage.prototype, "getItem");
    const set = vi.spyOn(Storage.prototype, "setItem");
    const remove = vi.spyOn(Storage.prototype, "removeItem");

    writeStoredTimezone("Asia/Shanghai");
    expect(set).toHaveBeenCalledWith(TIMEZONE_STORAGE_KEY, "Asia/Shanghai");
    expect(readStoredTimezone()).toBe("Asia/Shanghai");

    writeStoredTimezone(AUTO_TIMEZONE);
    expect(remove).toHaveBeenCalledWith(TIMEZONE_STORAGE_KEY);

    get.mockRestore();
    set.mockRestore();
    remove.mockRestore();
  });
});