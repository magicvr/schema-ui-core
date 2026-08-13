// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import {
  applyLocaleToDocument,
  I18nProvider,
  LOCALE_STORAGE_KEY,
  readStoredLocale,
  useI18n,
  writeStoredLocale,
} from "./runtime";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
  document.documentElement.lang = "";
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
});

function Harness(props: {
  systemDefault?: string | null;
  stored?: string | null;
  browserLanguages?: readonly string[];
}) {
  void props;
  const i18n = useI18n();
  return (
    <div
      data-testid="harness"
      data-locale={i18n.locale}
      data-preference={i18n.preference}
      data-translated={i18n.t("locale.switcher.label")}
      data-date={i18n.formatDate(new Date("2026-08-09T03:00:00.000Z"))}
      data-number={i18n.formatNumber(1234567.5)}
    >
      <button
        type="button"
        data-action="set-zh"
        onClick={() => i18n.setPreference("zh-CN")}
      />
      <button
        type="button"
        data-action="set-auto"
        onClick={() => i18n.setPreference("auto")}
      />
    </div>
  );
}

async function renderHarness(props: {
  systemDefault?: string | null;
  stored?: string | null;
  browserLanguages?: readonly string[];
}) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider {...props}>
        <Harness {...props} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("I18nProvider", () => {
  it("resolves from browser preference when nothing is stored (auto)", async () => {
    const container = await renderHarness({ browserLanguages: ["zh-CN", "en-US"] });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    expect(harness?.dataset.locale).toBe("zh-CN");
    expect(harness?.dataset.preference).toBe("auto");
  });

  it("applies explicit stored preference over browser languages", async () => {
    const container = await renderHarness({
      stored: "en-US",
      browserLanguages: ["zh-CN"],
    });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    expect(harness?.dataset.locale).toBe("en-US");
  });

  it("applies system default when no explicit choice (priority order)", async () => {
    const container = await renderHarness({
      systemDefault: "zh-CN",
      browserLanguages: ["en-US"],
    });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    expect(harness?.dataset.locale).toBe("zh-CN");
  });

  it("falls back to en-US safe default", async () => {
    const container = await renderHarness({ browserLanguages: [] });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    expect(harness?.dataset.locale).toBe("en-US");
  });

  it("translates and formats with the effective locale", async () => {
    const container = await renderHarness({
      stored: "zh-CN",
      browserLanguages: ["en-US"],
    });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    expect(harness?.dataset.translated).toBe("语言");
    expect(harness?.dataset.number).toContain("1,234,567");
  });

  it("persists the switch to localStorage and re-resolves immediately", async () => {
    const container = await renderHarness({ browserLanguages: ["en-US"] });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    expect(harness?.dataset.locale).toBe("en-US");

    const setZh = container.querySelector<HTMLButtonElement>("[data-action='set-zh']");
    await act(async () => setZh?.click());

    expect(harness?.dataset.locale).toBe("zh-CN");
    expect(harness?.dataset.preference).toBe("zh-CN");
    expect(localStorage.getItem(LOCALE_STORAGE_KEY)).toBe("zh-CN");
    expect(harness?.dataset.translated).toBe("语言");
  });

  it("switching back to auto removes the stored key (single channel)", async () => {
    const container = await renderHarness({ browserLanguages: ["zh-CN"] });
    const harness = container.querySelector<HTMLElement>("[data-testid='harness']");
    const setZh = container.querySelector<HTMLButtonElement>("[data-action='set-zh']");
    await act(async () => setZh?.click());
    expect(localStorage.getItem(LOCALE_STORAGE_KEY)).toBe("zh-CN");

    const setAuto = container.querySelector<HTMLButtonElement>("[data-action='set-auto']");
    await act(async () => setAuto?.click());
    expect(localStorage.getItem(LOCALE_STORAGE_KEY)).toBeNull();
    expect(harness?.dataset.preference).toBe("auto");
    // Resolves back through browser preference.
    expect(harness?.dataset.locale).toBe("zh-CN");
  });
});

// ── storage helpers + document lang ──────────────────────────────────────────

describe("storage helpers and document lang", () => {
  it("readStoredLocale / writeStoredLocale round-trip", () => {
    expect(readStoredLocale()).toBeNull();
    writeStoredLocale("zh-CN");
    expect(readStoredLocale()).toBe("zh-CN");
    writeStoredLocale("auto");
    expect(readStoredLocale()).toBeNull();
  });

  it("readStoredLocale / writeStoredLocale swallow disabled-storage throws", () => {
    const boom = () => {
      throw new Error("SecurityError");
    };
    const original = window.localStorage;
    Object.defineProperty(window, "localStorage", {
      configurable: true,
      value: {
        getItem: boom,
        setItem: boom,
        removeItem: boom,
        clear: boom,
        key: () => null,
        length: 0,
      },
    });
    try {
      expect(readStoredLocale()).toBeNull();
      expect(() => writeStoredLocale("zh-CN")).not.toThrow();
      expect(() => writeStoredLocale("auto")).not.toThrow();
    } finally {
      Object.defineProperty(window, "localStorage", { configurable: true, value: original });
    }
  });

  it("applyLocaleToDocument sets <html lang>", () => {
    expect(document.documentElement.lang).not.toBe("zh-CN");
    applyLocaleToDocument("zh-CN");
    expect(document.documentElement.lang).toBe("zh-CN");
    applyLocaleToDocument("en-US");
    expect(document.documentElement.lang).toBe("en-US");
  });

  it("provider applies <html lang> on mount and on switch", async () => {
    const container = await renderHarness({ browserLanguages: ["zh-CN"] });
    expect(document.documentElement.lang).toBe("zh-CN");
    const setZh = container.querySelector<HTMLButtonElement>("[data-action='set-zh']");
    await act(async () => setZh?.click());
    // Stays zh-CN (explicit now); switch to en-US via stored seam below.
    expect(document.documentElement.lang).toBe("zh-CN");
  });
});
