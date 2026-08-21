import { describe, expect, it } from "vitest";

import { formatDate, formatNumber } from "./format";

// ── formatDate ───────────────────────────────────────────────────────────────

describe("formatDate", () => {
  const sample = new Date("2026-08-09T03:00:00.000Z");

  it("formats with the zh-CN locale (locale-aware, no custom template)", () => {
    const out = formatDate(sample, "zh-CN");
    // Intl zh-CN medium date + short time — must contain the year and day.
    expect(out).toContain("2026");
    expect(out).toContain("8");
  });

  it("formats with the en-US locale", () => {
    const out = formatDate(sample, "en-US");
    expect(out).toContain("2026");
    expect(out.length).toBeGreaterThan(0);
  });

  it("applies an explicit IANA timezone when provided", () => {
    const shanghai = formatDate(sample, "en-US", { timeZone: "Asia/Shanghai" });
    const utc = formatDate(sample, "en-US", { timeZone: "UTC" });
    expect(shanghai).not.toBe(utc);
  });

  it("degrades safely on an invalid timezone instead of throwing", () => {
    const out = formatDate(sample, "en-US", { timeZone: "Foo/Bar" });
    expect(out.length).toBeGreaterThan(0);
    expect(out).toContain("2026");
  });

  it("returns an empty string for invalid date input", () => {
    expect(formatDate("not-a-date", "en-US")).toBe("");
    expect(formatDate(Number.NaN, "en-US")).toBe("");
  });

  it("accepts ISO strings and unix timestamps", () => {
    expect(formatDate("2026-08-09T03:00:00.000Z", "en-US")).toContain("2026");
    expect(formatDate(1782846000000, "en-US")).toContain("2026");
  });
});

// ── formatNumber ─────────────────────────────────────────────────────────────

describe("formatNumber", () => {
  it("formats with locale group separators", () => {
    const en = formatNumber(1234567.5, "en-US");
    const zh = formatNumber(1234567.5, "zh-CN");
    expect(en).toContain("1,234,567");
    expect(zh.length).toBeGreaterThan(0);
  });

  it("respects Intl.NumberFormatOptions", () => {
    expect(formatNumber(0.25, "en-US", { style: "percent" })).toBe("25%");
  });

  it("returns an empty string for invalid input", () => {
    expect(formatNumber(Number.NaN, "en-US")).toBe("");
    expect(formatNumber(Number.POSITIVE_INFINITY, "en-US")).toBe("");
  });
});
