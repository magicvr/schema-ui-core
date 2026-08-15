import { describe, expect, it } from "vitest";

import { formatDisplayTime } from "@/lib/datetime";

function localForm(date: Date): string {
  const pad = (value: number) => String(value).padStart(2, "0");
  return (
    date.getFullYear() +
    "-" +
    pad(date.getMonth() + 1) +
    "-" +
    pad(date.getDate()) +
    " " +
    pad(date.getHours()) +
    ":" +
    pad(date.getMinutes())
  );
}

describe("formatDisplayTime", () => {
  it("formats UTC ISO timestamps to local YYYY-MM-DD HH:mm", () => {
    const input = "2026-08-01T18:02:44.000Z";
    expect(formatDisplayTime(input)).toBe(localForm(new Date(input)));
    expect(formatDisplayTime(input)).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/);
  });

  it("accepts timestamps with and without fractional seconds and offsets", () => {
    expect(formatDisplayTime("2026-08-01T18:02:44Z")).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/);
    expect(formatDisplayTime("2026-08-01T18:02:44.000Z")).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/);
    expect(formatDisplayTime("2026-08-01T18:02:44+08:00")).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/);
  });

  it("leaves non-timestamp values untouched (null)", () => {
    expect(formatDisplayTime("Acme Console")).toBeNull();
    expect(formatDisplayTime("2026-08-01")).toBeNull();
    expect(formatDisplayTime("18:02:44")).toBeNull();
  });

  it("leaves non-string and malformed values untouched (null)", () => {
    expect(formatDisplayTime(12345)).toBeNull();
    expect(formatDisplayTime(null)).toBeNull();
    expect(formatDisplayTime(undefined)).toBeNull();
    expect(formatDisplayTime(true)).toBeNull();
    expect(formatDisplayTime("not-a-dateT01:02:03Z")).toBeNull();
    expect(formatDisplayTime("2026-13-45T99:99:99Z")).toBeNull();
  });
});