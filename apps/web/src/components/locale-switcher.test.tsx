// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { LocaleSwitcher } from "@/components/locale-switcher";
import { I18nProvider, LOCALE_STORAGE_KEY } from "@/i18n/runtime";

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

async function renderSwitcher(props: { browserLanguages?: readonly string[] }) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider {...props}>
        <LocaleSwitcher />
      </I18nProvider>,
    );
  });
  return container;
}

describe("LocaleSwitcher (header dropdown)", () => {
  it("renders a compact icon trigger (no form select)", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    expect(container.querySelector("select")).toBeNull();
    const trigger = container.querySelector<HTMLButtonElement>('button[aria-label="Language"]');
    expect(trigger).not.toBeNull();
    expect(trigger?.getAttribute("aria-haspopup")).toBe("menu");
    expect(trigger?.getAttribute("aria-expanded")).toBe("false");
    // No menu content before opening.
    expect(container.querySelector('[role="menu"]')).toBeNull();
  });

  it("opens a dropdown with auto + both supported locales and a checkmark on the current one", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    const trigger = container.querySelector<HTMLButtonElement>('button[aria-label="Language"]');
    await act(async () => trigger?.click());
    const menu = container.querySelector('[role="menu"]');
    expect(menu).not.toBeNull();
    const items = Array.from(menu?.querySelectorAll('[role="menuitemradio"]') ?? []);
    expect(items.map((item) => item.textContent?.trim())).toEqual([
      "Auto (follow browser)",
      "简体中文",
      "English",
    ]);
    // auto is the current preference → checkmarked.
    expect(items[0]?.getAttribute("aria-checked")).toBe("true");
    expect(items[0]?.querySelector("svg")).not.toBeNull();
    expect(items[1]?.getAttribute("aria-checked")).toBe("false");
  });

  it("selecting a locale persists to localStorage and closes the menu", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    const trigger = container.querySelector<HTMLButtonElement>('button[aria-label="Language"]');
    await act(async () => trigger?.click());
    const zh = Array.from(container.querySelectorAll('[role="menuitemradio"]')).find(
      (item) => item.textContent?.trim() === "简体中文",
    );
    await act(async () => (zh as HTMLButtonElement).click());
    expect(localStorage.getItem(LOCALE_STORAGE_KEY)).toBe("zh-CN");
    expect(container.querySelector('[role="menu"]')).toBeNull(); // closed after select
    // Reopen: zh-CN now carries the checkmark.
    await act(async () => trigger?.click());
    const items = Array.from(container.querySelectorAll('[role="menuitemradio"]'));
    expect(items[1]?.getAttribute("aria-checked")).toBe("true");
  });

  it("closes on Escape and on outside pointerdown", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    const trigger = container.querySelector<HTMLButtonElement>('button[aria-label="Language"]');
    await act(async () => trigger?.click());
    expect(container.querySelector('[role="menu"]')).not.toBeNull();
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));
    });
    expect(container.querySelector('[role="menu"]')).toBeNull();
    await act(async () => trigger?.click());
    await act(async () => {
      document.dispatchEvent(new Event("pointerdown", { bubbles: true }));
    });
    expect(container.querySelector('[role="menu"]')).toBeNull();
  });

  it("shows translated labels under zh-CN", async () => {
    const container = await renderSwitcher({ browserLanguages: ["zh-CN"] });
    const trigger = container.querySelector<HTMLButtonElement>('button[aria-label="语言"]');
    expect(trigger).not.toBeNull();
    await act(async () => trigger?.click());
    const items = Array.from(container.querySelectorAll('[role="menuitemradio"]'));
    expect(items.map((item) => item.textContent?.trim())).toEqual([
      "自动（跟随浏览器）",
      "简体中文",
      "English",
    ]);
  });

  it("reachable without any settings permission — no auth prop required", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    expect(container.querySelector('button[aria-label="Language"]')).not.toBeNull();
  });
});
