// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider, useI18n } from "./runtime";
import { TIMEZONE_STORAGE_KEY } from "./timezone";

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

/** Probe harness exposing the I18n state the tests need. */
function Harness() {
  const i18n = useI18n();
  return (
    <div
      data-testid="harness"
      data-timezone={i18n.timezone}
      data-preference={i18n.timezonePreference}
      data-date={i18n.formatDate(new Date("2026-08-09T03:00:00.000Z"))}
    >
      <button
        type="button"
        data-action="set-shanghai"
        onClick={() => i18n.setTimezonePreference("Asia/Shanghai")}
      />
      <button type="button" data-action="set-auto" onClick={() => i18n.setTimezonePreference("auto")} />
    </div>
  );
}

const PROBE = () => "America/New_York";

// ── L1/L3/L4 resolution wiring (C1 consumers) ────────────────────────────────

describe("I18nProvider timezone wiring", () => {
  it("L4: no inputs → effective timezone is auto", () => {
    const container = mount(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={() => ""} storedTimezone="auto" siteTimezone="auto">
        <Harness />
      </I18nProvider>,
    );
    const probe = container.querySelector("[data-testid='harness']");
    expect(probe?.getAttribute("data-timezone")).toBe("auto");
    expect(probe?.getAttribute("data-preference")).toBe("auto");
  });

  it("L3: site default applies when no user override and probe is empty", () => {
    const container = mount(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={() => ""} storedTimezone="auto" siteTimezone="Europe/London">
        <Harness />
      </I18nProvider>,
    );
    const probe = container.querySelector("[data-testid='harness']");
    expect(probe?.getAttribute("data-timezone")).toBe("Europe/London");
  });

  it("L1: user override wins over site default and probe (C2 persistence)", () => {
    const container = mount(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={PROBE} storedTimezone="Asia/Shanghai" siteTimezone="Europe/London">
        <Harness />
      </I18nProvider>,
    );
    const probe = container.querySelector("[data-testid='harness']");
    expect(probe?.getAttribute("data-timezone")).toBe("Asia/Shanghai");
    // The storedTimezone prop only seeds state; persistence is asserted in
    // the setTimezonePreference round-trip below.
    expect(probe?.getAttribute("data-preference")).toBe("Asia/Shanghai");
  });

  it("setTimezonePreference persists via the single channel and flips formatting (C2/C4)", () => {
    const container = mount(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={PROBE} storedTimezone="auto">
        <Harness />
      </I18nProvider>,
    );
    const probe = container.querySelector("[data-testid='harness']");
    // L2 probe applies while the preference is auto.
    const probeDate = probe?.getAttribute("data-date") ?? "";

    const shanghaiButton = container.querySelector<HTMLButtonElement>("[data-action='set-shanghai']");
    act(() => shanghaiButton?.click());
    expect(probe?.getAttribute("data-timezone")).toBe("Asia/Shanghai");
    expect(localStorage.getItem(TIMEZONE_STORAGE_KEY)).toBe("Asia/Shanghai");
    const shanghaiDate = probe?.getAttribute("data-date") ?? "";
    expect(shanghaiDate).not.toBe(probeDate);

    const autoButton = container.querySelector<HTMLButtonElement>("[data-action='set-auto']");
    act(() => autoButton?.click());
    expect(probe?.getAttribute("data-timezone")).toBe("America/New_York"); // L2 probe again
    expect(localStorage.getItem(TIMEZONE_STORAGE_KEY)).toBeNull();
  });

  it("explicit options.timeZone still overrides the effective timezone (C4 contract)", () => {
    const container = mount(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={PROBE} storedTimezone="Asia/Shanghai">
        <ProbeWithOverride />
      </I18nProvider>,
    );
    const probe = container.querySelector("[data-testid='override']");
    // 2026-08-09T03:00:00Z in UTC renders "3:00 AM" — in Asia/Shanghai it is
    // 11:00 — the override must win (UTC render contains "3:00").
    expect(probe?.getAttribute("data-date")).toContain("3:00");
  });

  it("C3: systemDefaultUrl fetch supplies siteTimezone (L3)", async () => {
    vi.stubGlobal(
      "fetch",
      vi.fn().mockResolvedValue(
        new Response(
          JSON.stringify({ defaultLocale: "zh-CN", siteTimezone: "Europe/London" }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        ),
      ),
    );
    const container = mount(
      <I18nProvider
        browserLanguages={["en-US"]}
        detectTimezone={() => ""}
        storedTimezone="auto"
        systemDefaultUrl="/api/branding"
      >
        <Harness />
      </I18nProvider>,
    );
    // Flush the fetch promise so the site default reaches state (C3).
    await act(async () => {});
    const probe = container.querySelector("[data-testid='harness']");
    expect(probe?.getAttribute("data-timezone")).toBe("Europe/London");
  });
});

function ProbeWithOverride() {
  const { formatDate } = useI18n();
  return (
    <div
      data-testid="override"
      data-date={formatDate(new Date("2026-08-09T03:00:00.000Z"), { timeZone: "UTC" })}
    />
  );
}