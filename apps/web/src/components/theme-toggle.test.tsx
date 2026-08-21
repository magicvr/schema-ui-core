// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { ThemeToggle } from "@/components/theme-toggle";
import { I18nProvider } from "@/i18n/runtime";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  document.documentElement.classList.remove("dark");
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  document.documentElement.classList.remove("dark");
});

async function renderToggle(browserLanguages: readonly string[]) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider browserLanguages={browserLanguages}>
        <ThemeToggle />
      </I18nProvider>,
    );
  });
  return container;
}

describe("ThemeToggle (W8 visual unification · GOAL-009)", () => {
  it("is a ghost icon button matching the header icon row (no form select/outline)", async () => {
    const container = await renderToggle(["en-US"]);
    const button = container.querySelector<HTMLButtonElement>('button[aria-label="Toggle color theme"]');
    expect(button).not.toBeNull();
    expect(button?.className).toContain("size-9");
    expect(button?.className).toContain("hover:bg-accent");
    expect(button?.className).not.toContain("border"); // no outline/form styling
    expect(button?.querySelector("svg")).not.toBeNull(); // Sun/Moon icon
  });

  it("tooltip and aria-label follow the active locale (no hardcoded English)", async () => {
    const en = await renderToggle(["en-US"]);
    const enButton = en.querySelector<HTMLButtonElement>('button[aria-label="Toggle color theme"]');
    expect(enButton?.getAttribute("title")).toBe("Toggle color theme");

    const zh = await renderToggle(["zh-CN"]);
    const zhButton = zh.querySelector<HTMLButtonElement>('button[aria-label="切换明暗主题"]');
    expect(zhButton).not.toBeNull();
    expect(zhButton?.getAttribute("title")).toBe("切换明暗主题");
  });

  it("clicking toggles the dark class on the root", async () => {
    const container = await renderToggle(["en-US"]);
    const button = container.querySelector<HTMLButtonElement>('button[aria-label="Toggle color theme"]');
    expect(document.documentElement.classList.contains("dark")).toBe(false);
    await act(async () => button?.click());
    expect(document.documentElement.classList.contains("dark")).toBe(true);
    await act(async () => button?.click());
    expect(document.documentElement.classList.contains("dark")).toBe(false);
  });
});
