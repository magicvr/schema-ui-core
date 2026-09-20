// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { formatDate } from "./format";
import { I18nProvider, useI18n } from "./runtime";

/**
 * workspace-040 R3-B · VP-020 round-trip matrix (I-040-004; carrier ruled by
 * I-041-007): the Go unit tests pin the wire shape, the parser compatibility
 * matrix and microsecond-exact round-tripping; these component/unit tests pin the
 * DISPLAY round-trip — the stored/transmitted value is a UTC instant that never
 * moves, while the rendered wall clock follows the effective timezone
 * (L1 user override → L2 session probe → L3 site default → L4 auto).
 *
 * The API input path is date-only (`datePicker` → `YYYY-MM-DD`, asserted in
 * `renderer/form-controls.test.ts`) and the server rejects zoneless wire input
 * (Root D-005), so no zoneless local value can reach storage from here.
 */

const WIRE = "2026-09-20T12:57:15.900000Z";

const WIRE_INSTANTS = [
  "2026-09-20T12:57:15.900000Z", // trailing-zero microseconds (a width-losing renderer drops to ".9")
  "2026-01-15T00:00:00.000000Z",
  "2026-03-08T07:30:00.000001Z", // US DST spring-forward day
  "2026-11-01T05:30:00.999999Z", // US DST fall-back day
  "1965-03-04T05:06:07.654321Z", // negative Unix epoch is a legal instant (Root D-015)
] as const;

const ZONES = ["UTC", "Asia/Shanghai", "America/New_York", "Asia/Kathmandu"] as const;

interface WallClock {
  year: number;
  month: number;
  day: number;
  hour: number;
  minute: number;
  second: number;
}

const WALL_CLOCK_FORMAT: Intl.DateTimeFormatOptions = {
  hour12: false,
  year: "numeric",
  month: "2-digit",
  day: "2-digit",
  hour: "2-digit",
  minute: "2-digit",
  second: "2-digit",
};

/** The wall clock an Intl-based renderer shows for `instantMs` in `timeZone`. */
function wallClockIn(instantMs: number, timeZone: string): WallClock {
  const parts = new Intl.DateTimeFormat("en-US", { ...WALL_CLOCK_FORMAT, timeZone })
    .formatToParts(new Date(instantMs))
    .filter((part) => part.type !== "literal");
  const read = (type: Intl.DateTimeFormatPartTypes) => Number(parts.find((p) => p.type === type)?.value);
  return {
    year: read("year"),
    month: read("month"),
    day: read("day"),
    hour: read("hour") % 24,
    minute: read("minute"),
    second: read("second"),
  };
}

/** The zone's offset (ms) at `instantMs`, at second granularity. */
function zoneOffsetMs(instantMs: number, timeZone: string): number {
  const base = Math.floor(instantMs / 1000) * 1000;
  const wall = wallClockIn(base, timeZone);
  return Date.UTC(wall.year, wall.month - 1, wall.day, wall.hour, wall.minute, wall.second) - base;
}

/** Recovers the instant a user in `timeZone` means by that wall clock. */
function instantFromWallClock(wall: WallClock, timeZone: string): number {
  const asUtc = Date.UTC(wall.year, wall.month - 1, wall.day, wall.hour, wall.minute, wall.second);
  let instant = asUtc;
  for (let i = 0; i < 3; i += 1) {
    const next = asUtc - zoneOffsetMs(instant, timeZone);
    if (next === instant) {
      return instant;
    }
    instant = next;
  }
  return instant;
}

function canonicalWire(instantMs: number): string {
  return new Date(instantMs).toISOString().replace(/\.\d{3}Z$/, ".000000Z");
}

function wireAtSecond(wire: string): string {
  return `${wire.slice(0, 19)}.000000Z`;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
  vi.unstubAllGlobals();
  vi.restoreAllMocks();
});

function mount(children: React.ReactNode) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  act(() => {
    root.render(children);
  });
  activeRoots.push({ root, container });
  return container;
}

/** Probe exposing the display values the round-trip assertions need. */
function Harness({ value }: { value: string }) {
  const i18n = useI18n();
  return (
    <div
      data-testid="harness"
      data-locale={i18n.locale}
      data-timezone={i18n.timezone}
      data-display={i18n.formatDate(value)}
      data-display-session={i18n.formatDate(value, { timeZone: "America/New_York" })}
      data-display-utc={i18n.formatDate(value, { timeZone: "UTC" })}
    />
  );
}

describe("VP-020 UTC storage ↔ session timezone display round-trip", () => {
  it("renders the wire instant in the effective (user) timezone, not in the probe or UTC", () => {
    const container = mount(
      <I18nProvider
        browserLanguages={["en-US"]}
        detectTimezone={() => "America/New_York"}
        storedTimezone="Asia/Shanghai"
        siteTimezone="Europe/London"
      >
        <Harness value={WIRE} />
      </I18nProvider>,
    );
    const probe = container.querySelector("[data-testid='harness']");
    expect(probe?.getAttribute("data-timezone")).toBe("Asia/Shanghai");

    const locale = (probe?.getAttribute("data-locale") ?? "en-US") as "en-US";
    const display = probe?.getAttribute("data-display") ?? "";
    // The displayed value is exactly the production renderer bound to the
    // resolved session zone…
    expect(display).toBe(formatDate(WIRE, locale, { timeZone: "Asia/Shanghai" }));
    // …and it is neither the session-probe rendering nor the UTC rendering, so
    // the assertion would fail if the provider resolved the wrong layer.
    expect(display).not.toBe(probe?.getAttribute("data-display-session"));
    expect(display).not.toBe(probe?.getAttribute("data-display-utc"));
    // A whole-hour zone shift is visible in the rendered wall clock.
    expect(display).not.toBe("");
  });

  it("does not mutate the instant it renders", () => {
    const instant = new Date(WIRE);
    const before = instant.toISOString();
    mount(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={() => "UTC"} storedTimezone="Asia/Shanghai">
        <Harness value={WIRE} />
      </I18nProvider>,
    );
    expect(instant.toISOString()).toBe(before);
    expect(formatDate(instant, "en-US", { timeZone: "Asia/Shanghai" })).toBe(
      formatDate(WIRE, "en-US", { timeZone: "Asia/Shanghai" }),
    );
  });

  it.each(ZONES)("display → instant → wire round-trips in %s", (zone) => {
    for (const wire of WIRE_INSTANTS) {
      const instantMs = Date.parse(wire);
      expect(Number.isFinite(instantMs)).toBe(true);
      const wall = wallClockIn(instantMs, zone);
      const back = instantFromWallClock(wall, zone);
      // The human-readable display resolves to seconds, so the recovered wire
      // value is exact at second granularity (microseconds stay storage-only).
      expect(canonicalWire(back)).toBe(wireAtSecond(wire));
    }
  });

  it.each([
    ["UTC", "12:57"],
    ["Asia/Shanghai", "20:57"],
    ["America/New_York", "08:57"],
    ["Asia/Kathmandu", "18:42"],
  ])("the local wall-clock helper agrees with the production formatter in %s (%s)", (zone, hhmm) => {
    // The round-trip helpers below are test-local; tie them to the production
    // renderer so a helper that silently disagrees with the shipped display is
    // caught here (independent audit A-002 F-I-005).
    expect(formatDate(WIRE, "zh-CN", { timeZone: zone })).toContain(hhmm);
    const wall = wallClockIn(Date.parse(WIRE), zone);
    expect(`${String(wall.hour).padStart(2, "0")}:${String(wall.minute).padStart(2, "0")}`).toBe(hhmm);
  });

  it("covers a non-hour zone offset (+05:45) and both US DST edges", () => {
    // Asia/Kathmandu is UTC+05:45: a naive whole-hour assumption would shift it.
    const kathmandu = wallClockIn(Date.parse("2026-09-20T12:57:15.900000Z"), "Asia/Kathmandu");
    expect(kathmandu.hour).toBe(18);
    expect(kathmandu.minute).toBe(42);

    // Spring forward (2026-03-08): 07:30Z is 02:30 EST / 03:30 EDT — the
    // recovered instant must not drift by the hour that does not exist.
    expect(canonicalWire(instantFromWallClock(wallClockIn(Date.parse("2026-03-08T07:30:00.000001Z"), "America/New_York"), "America/New_York")))
      .toBe("2026-03-08T07:30:00.000000Z");
    // Fall back (2026-11-01) is unambiguous at 05:30Z.
    expect(canonicalWire(instantFromWallClock(wallClockIn(Date.parse("2026-11-01T05:30:00.999999Z"), "America/New_York"), "America/New_York")))
      .toBe("2026-11-01T05:30:00.000000Z");
  });

  it("shows the same instant differently per session timezone without touching the stored value", () => {
    const stored = WIRE;
    const shanghai = formatDate(stored, "en-US", { timeZone: "Asia/Shanghai" });
    const utc = formatDate(stored, "en-US", { timeZone: "UTC" });
    expect(shanghai).not.toBe(utc);
    // The wire value is the same string before and after every rendering: the
    // session timezone is a presentation concern only.
    expect(stored).toBe(WIRE);
    expect(formatDate(stored, "en-US", { timeZone: "Asia/Shanghai" })).toBe(shanghai);
  });
});
