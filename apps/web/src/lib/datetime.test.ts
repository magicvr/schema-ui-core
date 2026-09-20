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

  // workspace-040 R3-A fixture sync (inventory §Web consumer): the API wire
  // value became the canonical fixed-6 shape, so the display formatter must
  // render it, and must still render the legacy 3-digit shape it replaced.
  it("renders the canonical fixed-6 wire shape and the legacy 3-digit shape", () => {
    const fixed6 = "2026-08-01T18:02:44.123456Z";
    const legacy3 = "2026-08-01T18:02:44.123Z";
    expect(formatDisplayTime(fixed6)).toBe(localForm(new Date(fixed6)));
    expect(formatDisplayTime(legacy3)).toBe(localForm(new Date(legacy3)));
    // Same instant, different fraction width -> identical rendering.
    expect(formatDisplayTime(fixed6)).toBe(formatDisplayTime(legacy3));
    expect(formatDisplayTime(fixed6)).toMatch(/^\d{4}-\d{2}-\d{2} \d{2}:\d{2}$/);
  });

  it("renders a non-zero-offset fixed-6 value at the same instant as its UTC form", () => {
    expect(formatDisplayTime("2026-08-02T02:02:44.123456+08:00")).toBe(
      formatDisplayTime("2026-08-01T18:02:44.123456Z"),
    );
  });

  it("leaves non-timestamp values untouched (null)", () => {
    expect(formatDisplayTime("Acme Console")).toBeNull();
    expect(formatDisplayTime("2026-08-01")).toBeNull();
    expect(formatDisplayTime("18:02:44")).toBeNull();
    // Zoneless fixed-6 is ambiguous local time: never guessed, always raw text.
    expect(formatDisplayTime("2026-08-01T18:02:44.123456")).toBeNull();
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