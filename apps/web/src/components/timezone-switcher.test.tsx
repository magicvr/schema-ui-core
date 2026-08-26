// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { TimezoneSwitcher } from "./timezone-switcher";
import { I18nProvider } from "@/i18n/runtime";
import { TIMEZONE_STORAGE_KEY } from "@/i18n/timezone";

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
});

function mountSwitcher() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  act(() => {
    root.render(
      <I18nProvider browserLanguages={["en-US"]} detectTimezone={() => "America/New_York"}>
        <TimezoneSwitcher />
      </I18nProvider>,
    );
  });
  activeRoots.push({ root, container });
  return container;
}

describe("TimezoneSwitcher", () => {
  it("opens a menu with auto + the common IANA set and marks the current one", () => {
    const container = mountSwitcher();
    const trigger = container.querySelector<HTMLButtonElement>("button[aria-haspopup='menu']");
    expect(trigger?.getAttribute("aria-label")).toBe("Timezone");

    act(() => trigger?.click());
    const menu = container.querySelector<HTMLElement>("[role='menu']");
    expect(menu).not.toBeNull();
    const items = Array.from(menu?.querySelectorAll<HTMLButtonElement>("[role='menuitemradio']") ?? []);
    const labels = items.map((item) => item.textContent?.trim());
    expect(labels).toContain("Auto (system)");
    expect(labels).toContain("Asia/Shanghai");
    expect(labels).toContain("UTC");
    // auto is the initial preference → checked.
    expect(items[0]?.getAttribute("aria-checked")).toBe("true");
  });

  it("selecting a timezone persists to localStorage and closes the menu", () => {
    const container = mountSwitcher();
    const trigger = container.querySelector<HTMLButtonElement>("button[aria-haspopup='menu']");
    act(() => trigger?.click());
    const menu = container.querySelector<HTMLElement>("[role='menu']");
    const shanghai = Array.from(
      menu?.querySelectorAll<HTMLButtonElement>("[role='menuitemradio']") ?? [],
    ).find((item) => item.textContent?.includes("Asia/Shanghai"));

    act(() => shanghai?.click());
    expect(localStorage.getItem(TIMEZONE_STORAGE_KEY)).toBe("Asia/Shanghai");
    expect(container.querySelector<HTMLElement>("[role='menu']")).toBeNull();
  });

  it("selecting auto removes the stored key", () => {
    localStorage.setItem(TIMEZONE_STORAGE_KEY, "Asia/Shanghai");
    const container = mountSwitcher();
    const trigger = container.querySelector<HTMLButtonElement>("button[aria-haspopup='menu']");
    act(() => trigger?.click());
    const menu = container.querySelector<HTMLElement>("[role='menu']");
    const auto = Array.from(
      menu?.querySelectorAll<HTMLButtonElement>("[role='menuitemradio']") ?? [],
    ).find((item) => item.textContent?.includes("Auto (system)"));

    // The stored preference wins → Asia/Shanghai is the checked item now.
    const shanghaiItem = Array.from(
      menu?.querySelectorAll<HTMLButtonElement>("[role='menuitemradio']") ?? [],
    ).find((item) => item.textContent?.includes("Asia/Shanghai"));
    expect(shanghaiItem?.getAttribute("aria-checked")).toBe("true");

    act(() => auto?.click());
    expect(localStorage.getItem(TIMEZONE_STORAGE_KEY)).toBeNull();
  });

  it("closes on outside pointerdown (parity with the locale switcher)", () => {
    const container = mountSwitcher();
    const trigger = container.querySelector<HTMLButtonElement>("button[aria-haspopup='menu']");
    act(() => trigger?.click());
    expect(container.querySelector("[role='menu']")).not.toBeNull();
    act(() => {
      document.dispatchEvent(new MouseEvent("pointerdown", { bubbles: true }));
    });
    expect(container.querySelector("[role='menu']")).toBeNull();
  });
});