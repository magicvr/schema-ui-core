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

describe("LocaleSwitcher", () => {
  it("renders auto + both supported locales with catalog labels", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    const select = container.querySelector<HTMLSelectElement>("select");
    expect(select).not.toBeNull();
    const options = Array.from(select?.options ?? []);
    expect(options.map((option) => option.value)).toEqual(["auto", "zh-CN", "en-US"]);
    expect(options[0].text).toBe("Auto (follow browser)");
    expect(options[1].text).toBe("简体中文");
    expect(options[2].text).toBe("English");
    expect(select?.getAttribute("aria-label")).toBe("Language");
  });

  it("shows translated labels under zh-CN", async () => {
    const container = await renderSwitcher({ browserLanguages: ["zh-CN"] });
    const select = container.querySelector<HTMLSelectElement>("select");
    expect(select?.getAttribute("aria-label")).toBe("语言");
    expect(select?.options[0].text).toBe("自动（跟随浏览器）");
  });

  it("switching persists to localStorage and takes effect", async () => {
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    const select = container.querySelector<HTMLSelectElement>("select");
    expect(select?.value).toBe("auto");

    await act(async () => {
      select!.value = "zh-CN";
      select!.dispatchEvent(new Event("change", { bubbles: true }));
    });

    expect(localStorage.getItem(LOCALE_STORAGE_KEY)).toBe("zh-CN");
    expect(select?.value).toBe("zh-CN");
    expect(select?.options[0].text).toBe("自动（跟随浏览器）"); // labels re-render via zh-CN? option 0 is "auto" label — same in both; aria-label is the zh label:
    expect(select?.getAttribute("aria-label")).toBe("语言");
  });

  it("reachable without any settings permission — no auth prop required", async () => {
    // The switcher only needs I18nProvider; no user/permission inputs exist.
    const container = await renderSwitcher({ browserLanguages: ["en-US"] });
    expect(container.querySelector("select")).not.toBeNull();
  });
});
