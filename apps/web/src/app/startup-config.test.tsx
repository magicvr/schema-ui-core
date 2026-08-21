// @vitest-environment jsdom

/**
 * S3 startup configuration tests (GOAL-004 · C2/C3/C5):
 * - fetchBranding parses the VP-007 extended public startup payload.
 * - applyDocumentBranding applies faviconUrl (fallback logoUrl) + title.
 * - The shell applies light/dark logo variants and the system default theme
 *   when branding loads; the login page applies the favicon + default theme.
 * - I18nProvider resolves the system default locale from /api/branding.
 * - Settings page renders the four-category surface with permission gating.
 */

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { LoginPage } from "@/app/LoginPage";
import {
  applyDocumentBranding,
  defaultBranding,
  fetchBranding,
  isSafeBrandingUrl,
  type Branding,
} from "@/app/branding";
import { I18nProvider, LOCALE_STORAGE_KEY, useI18n } from "@/i18n/runtime";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const __dir = dirname(fileURLToPath(import.meta.url));
const MANIFEST_PATH = resolve(__dir, "../test-fixtures/app-manifest.admin.json");
const SETTINGS_SCHEMA_PATH = resolve(
  __dir,
  "../../../api/internal/modules/settings/schema/settings.json",
);

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

const STARTUP_BODY: Record<string, unknown> = {
  siteTitle: "Acme Admin",
  logoUrl: "/assets/logo.svg",
  logoUrlLight: "/assets/logo-light.svg",
  logoUrlDark: "/assets/logo-dark.svg",
  faviconUrl: "/favicon.ico",
  defaultLocale: "zh-CN",
  supportedLocales: ["zh-CN", "en-US"],
  siteTimezone: "Asia/Shanghai",
  defaultTheme: "dark",
};

const SETTINGS_ROW: Record<string, unknown> = {
  id: "default",
  siteTitle: "Acme Admin",
  logoUrl: "/assets/logo.svg",
  logoUrlLight: "/assets/logo-light.svg",
  logoUrlDark: "/assets/logo-dark.svg",
  faviconUrl: "/favicon.ico",
  defaultLocale: "zh-CN",
  siteTimezone: "Asia/Shanghai",
  defaultTheme: "dark",
  updatedAt: "2026-08-09T00:00:00.000Z",
};

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  window.matchMedia = ((query: string) => ({
    matches: false,
    media: query,
    onchange: null,
    addListener: () => undefined,
    removeListener: () => undefined,
    addEventListener: () => undefined,
    removeEventListener: () => undefined,
    dispatchEvent: () => false,
  })) as typeof window.matchMedia;
  localStorage.clear();
  document.title = "";
  document.documentElement.lang = "";
  document.head.querySelectorAll("link[data-schema-ui-branding]").forEach((link) => link.remove());
  document.documentElement.classList.remove("dark");
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
});

function adminManifest(): AppManifest {
  return validateAppManifest(JSON.parse(readFileSync(MANIFEST_PATH, "utf8")));
}

function settingsSchemaDocument(): Record<string, unknown> {
  return JSON.parse(readFileSync(SETTINGS_SCHEMA_PATH, "utf8"));
}

function adminContext() {
  return {
    user: { id: "u1", roles: ["admin"], permissions: ["settings.read", "settings.write"] },
    features: { menu_settings: true },
  };
}

function fetcherFor(): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    if (pathname === "/api/branding") {
      return new Response(JSON.stringify(STARTUP_BODY), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname === "/api/settings/default") {
      return new Response(JSON.stringify(SETTINGS_ROW), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname === "/api/captcha/settings") {
      return new Response(JSON.stringify({ enabled: "true" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/settings")) {
      return new Response(JSON.stringify({ items: [SETTINGS_ROW], total: 1, page: 1, pageSize: 10 }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/schema/settings")) {
      return new Response(JSON.stringify(settingsSchemaDocument()), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
  }) as typeof fetch;
}

async function renderAt(path: string, element: React.ReactElement): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(element);
  });
  return container;
}

// ── C2 · fetchBranding extended payload ───────────────────────────────────────

describe("S3 · fetchBranding startup payload", () => {
  it("parses the VP-007 extended fields and defaults the legacy shape", async () => {
    globalThis.fetch = (async () =>
      new Response(JSON.stringify(STARTUP_BODY), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      })) as typeof fetch;
    const branding = await fetchBranding();
    expect(branding.siteTitle).toBe("Acme Admin");
    expect(branding.logoUrlLight).toBe("/assets/logo-light.svg");
    expect(branding.logoUrlDark).toBe("/assets/logo-dark.svg");
    expect(branding.faviconUrl).toBe("/favicon.ico");
    expect(branding.defaultLocale).toBe("zh-CN");
    expect(branding.supportedLocales).toEqual(["zh-CN", "en-US"]);
    expect(branding.siteTimezone).toBe("Asia/Shanghai");
    expect(branding.defaultTheme).toBe("dark");
  });

  it("falls back to safe defaults on failure or missing fields", async () => {
    globalThis.fetch = (async () => new Response("boom", { status: 503 })) as typeof fetch;
    const failed = await fetchBranding();
    expect(failed.siteTitle).toBe("Schema UI Core");
    expect(failed.defaultLocale).toBe("auto");
    expect(failed.defaultTheme).toBe("auto");
    expect(failed.supportedLocales).toEqual([]);

    globalThis.fetch = (async () =>
      new Response(JSON.stringify({ siteTitle: "T" }), { status: 200 })) as typeof fetch;
    const sparse = await fetchBranding();
    expect(sparse.defaultLocale).toBe("auto");
    expect(sparse.siteTimezone).toBe("auto");
  });

  it("applyDocumentBranding sets the title and favicon from faviconUrl", () => {
    const branding: Branding = { ...defaultBranding(), siteTitle: "Acme", faviconUrl: "/favicon.ico" };
    applyDocumentBranding(branding);
    expect(document.title).toBe("Acme");
    const link = document.querySelector<HTMLLinkElement>("link[rel='icon'][data-schema-ui-branding]");
    expect(link?.href).toContain("/favicon.ico");
  });

  it("applyDocumentBranding falls back to logoUrl for the favicon and clears on empty", () => {
    applyDocumentBranding({ ...defaultBranding(), logoUrl: "/assets/logo.svg" });
    const link = document.querySelector<HTMLLinkElement>("link[rel='icon'][data-schema-ui-branding]");
    expect(link?.href).toContain("/assets/logo.svg");
    applyDocumentBranding(defaultBranding());
    expect(document.querySelector("link[rel='icon'][data-schema-ui-branding]")).toBeNull();
  });

  it("drops unsafe branding URLs instead of writing them to img/link sinks", async () => {
    expect(isSafeBrandingUrl("/assets/logo.svg")).toBe(true);
    expect(isSafeBrandingUrl("https://cdn.example/logo.png")).toBe(true);
    expect(isSafeBrandingUrl("javascript:alert(1)")).toBe(false);
    expect(isSafeBrandingUrl("//evil.example/x")).toBe(false);
    expect(isSafeBrandingUrl("/\\evil.example/x")).toBe(false);
    expect(isSafeBrandingUrl("data:image/svg+xml,<svg>")).toBe(false);

    globalThis.fetch = (async () =>
      new Response(
        JSON.stringify({
          siteTitle: "X",
          logoUrl: "javascript:alert(1)",
          logoUrlLight: "//evil.example/l",
          logoUrlDark: "/\\evil.example/d",
          faviconUrl: "data:text/html,x",
        }),
        { status: 200 },
      )) as typeof fetch;
    const branding = await fetchBranding();
    expect(branding.logoUrl).toBe("");
    expect(branding.logoUrlLight).toBe("");
    expect(branding.logoUrlDark).toBe("");
    expect(branding.faviconUrl).toBe("");

    applyDocumentBranding({ ...defaultBranding(), faviconUrl: "javascript:alert(1)" });
    expect(document.querySelector("link[rel='icon'][data-schema-ui-branding]")).toBeNull();
  });
});

// ── C3 · projection: shell logos, favicon, system default theme/locale ────────

describe("S3 · shell/login projection of startup config", () => {
  it("applies the system default theme when the user has no explicit choice", async () => {
    globalThis.fetch = (async () =>
      new Response(JSON.stringify(STARTUP_BODY), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      })) as typeof fetch;
    const container = await renderAt(
      "/",
      <I18nProvider>
        <LoginPage onLogin={async () => undefined} />
      </I18nProvider>,
    );
    // Branding fetch resolves async: defaultTheme dark must flip the root.
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    expect(document.documentElement.classList.contains("dark")).toBe(true);
    expect(document.title).toBe("Acme Admin");
    container.remove();
  });

  it("keeps the explicit user theme over the system default", async () => {
    localStorage.setItem("theme", "light");
    globalThis.fetch = (async () =>
      new Response(JSON.stringify(STARTUP_BODY), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      })) as typeof fetch;
    const container = await renderAt(
      "/",
      <I18nProvider>
        <LoginPage onLogin={async () => undefined} />
      </I18nProvider>,
    );
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    expect(document.documentElement.classList.contains("dark")).toBe(false);
    container.remove();
  });

  it("renders light/dark logo variants in the shell and applies the favicon", async () => {
    globalThis.fetch = (async () =>
      new Response(JSON.stringify(STARTUP_BODY), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      })) as typeof fetch;
    const container = await renderAt(
      "/settings",
      <I18nProvider>
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor()}
          resourceFetcher={fetcherFor()}
        />
      </I18nProvider>,
    );
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const lightImg = container.querySelector<HTMLImageElement>("img[src='/assets/logo-light.svg']");
    const darkImg = container.querySelector<HTMLImageElement>("img[src='/assets/logo-dark.svg']");
    expect(lightImg).not.toBeNull();
    expect(darkImg).not.toBeNull();
    expect(lightImg?.className).toContain("dark:hidden");
    expect(darkImg?.className).toContain("dark:block");
    const favicon = document.querySelector<HTMLLinkElement>("link[rel='icon'][data-schema-ui-branding]");
    expect(favicon?.href).toContain("/favicon.ico");
  });
});

// ── I18nProvider system default from /api/branding ────────────────────────────

describe("S3 · provider system default locale", () => {
  it("re-resolves the locale from the startup payload when no explicit choice", async () => {
    globalThis.fetch = (async (input: RequestInfo | URL) => {
      if (String(input) === "/api/branding") {
        return new Response(JSON.stringify(STARTUP_BODY), { status: 200 });
      }
      return new Response("not found", { status: 404 });
    }) as typeof fetch;

    function Harness() {
      return (
        <I18nProvider systemDefaultUrl="/api/branding" browserLanguages={["en-US"]}>
          <LocaleProbe />
        </I18nProvider>
      );
    }
    function LocaleProbe() {
      const { locale, t } = useI18n();
      return (
        <div data-locale={locale} data-text={t("locale.switcher.label")}>
          {locale}
        </div>
      );
    }
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(<Harness />);
    });
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const probe = container.querySelector("[data-locale]");
    expect(probe?.getAttribute("data-locale")).toBe("zh-CN");
    expect(probe?.getAttribute("data-text")).toBe("语言");
  });

  it("explicit user choice still wins over the fetched system default", async () => {
    localStorage.setItem(LOCALE_STORAGE_KEY, "en-US");
    globalThis.fetch = (async () =>
      new Response(JSON.stringify(STARTUP_BODY), { status: 200 })) as typeof fetch;
    function Probe() {
      const { locale } = useI18n();
      return <div data-locale={locale} />;
    }
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider systemDefaultUrl="/api/branding" browserLanguages={["zh-CN"]}>
          <Probe />
        </I18nProvider>,
      );
    });
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    expect(container.querySelector("[data-locale]")?.getAttribute("data-locale")).toBe("en-US");
  });
});

// ── C4/C5 · settings page four-category surface + permission gating ───────────

describe("S3 · settings page four-category surface", () => {
  // W13 T-01: the settings page switches by functional unit through tabs
  // (same shape as the account page); Restore defaults stays reachable
  // outside the tabs.
  function clickTab(container: HTMLElement, label: string): void {
    const tab = [...container.querySelectorAll("button[role=tab]")].find((button) =>
      button.textContent?.includes(label),
    );
    expect(tab, "tab " + label + " present").not.toBeUndefined();
    (tab as HTMLButtonElement).click();
  }

  it("renders functional-unit tabs + prefilled forms per tab + restore defaults (zh-CN)", async () => {
    const container = await renderAt(
      "/settings",
      <I18nProvider stored="zh-CN" browserLanguages={["en-US"]}>
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={fetcherFor()}
          resourceFetcher={fetcherFor()}
        />
      </I18nProvider>,
    );
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const text = container.textContent ?? "";
    // The five functional-unit tab labels are always rendered.
    expect(text).toContain("常规");
    expect(text).toContain("品牌");
    expect(text).toContain("本地化");
    expect(text).toContain("外观");
    expect(text).toContain("安全");
    // Restore defaults stays outside the tabs (any tab can reach it).
    expect(text).toContain("恢复默认");

    // First tab (常规) is active by default: siteTitle form prefilled.
    const titleInput = container.querySelector<HTMLInputElement>("#field-siteTitle");
    expect(titleInput?.value).toBe("Acme Admin");

    // 品牌 tab → upload fields surface.
    await act(async () => {
      clickTab(container, "品牌");
    });
    expect(container.textContent).toContain("浅色主题 Logo");
    expect(container.textContent).toContain("Favicon");

    // 本地化 tab → locale + timezone prefilled.
    await act(async () => {
      clickTab(container, "本地化");
    });
    const localeSelect = container.querySelector<HTMLSelectElement>("#field-defaultLocale");
    expect(localeSelect?.value).toBe("zh-CN");
    expect(container.textContent).toContain("时区");

    // 外观 tab → theme prefilled.
    await act(async () => {
      clickTab(container, "外观");
    });
    const themeSelect = container.querySelector<HTMLSelectElement>("#field-defaultTheme");
    expect(themeSelect?.value).toBe("dark");

    // 安全 tab → captcha select surfaces.
    await act(async () => {
      clickTab(container, "安全");
    });
    expect(container.textContent).toContain("登录时要求验证码");
  });

  it("edits the inline General form and saves through the real PATCH action", async () => {
    const patched: Array<{ url: string; body: unknown }> = [];
    const base = fetcherFor();
    const tracking: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (init?.method === "PATCH" && url.startsWith("/api/settings/")) {
        patched.push({ url, body: JSON.parse(String(init.body)) });
        return new Response(JSON.stringify({ ...SETTINGS_ROW, ...JSON.parse(String(init.body)) }), {
          status: 200,
          headers: { "X-Schema-UI-Config-Changed": "settings.branding", "Content-Type": "application/json" },
        });
      }
      return base(input, init);
    }) as typeof fetch;
    const container = await renderAt(
      "/settings",
      <I18nProvider stored="zh-CN" browserLanguages={["en-US"]}>
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={tracking}
          resourceFetcher={tracking}
        />
      </I18nProvider>,
    );
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const titleInput = container.querySelector<HTMLInputElement>("#field-siteTitle");
    expect(titleInput).not.toBeNull();
    await act(async () => {
      // React tracks controlled-input values via a value tracker; use the
      // native prototype setter so onChange observes the new value (same
      // pattern as schema-crud.test.tsx setFieldValue).
      Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set?.call(
        titleInput,
        "Renamed",
      );
      titleInput!.dispatchEvent(new Event("input", { bubbles: true }));
    });
    const form = container.querySelector("form");
    expect(form).not.toBeNull();
    await act(async () => {
      form?.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
    });
    expect(patched.length).toBe(1);
    expect(patched[0]!.url).toContain("/api/settings/default");
    expect((patched[0]!.body as Record<string, unknown>).siteTitle).toBe("Renamed");
  });

  it("gates the category actions behind settings.write (no write → disabled)", async () => {
    const viewerContext = {
      user: { id: "v1", roles: ["viewer"], permissions: ["settings.read"] },
      features: { menu_settings: true },
    };
    const container = await renderAt(
      "/settings",
      <I18nProvider>
        <App
          manifest={adminManifest()}
          navigationContext={viewerContext}
          schemaFetcher={fetcherFor()}
          resourceFetcher={fetcherFor()}
        />
      </I18nProvider>,
    );
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    // W13 T-01: the Restore defaults action sits outside the tabs and is
    // always visible; it is gated like every category action. The provider
    // re-resolves the locale from /api/branding (defaultLocale zh-CN), so
    // match labels in either locale.
    const actionLabels = ["Restore defaults", "恢复默认"];
    const saveLabels = ["Save settings", "保存设置"];
    const tabLabels = [["General", "常规"], ["Branding", "品牌"], ["Localization", "本地化"], ["Appearance", "外观"], ["Security", "安全"]];
    const restoreButtons = [...container.querySelectorAll("button")].filter((button) =>
      actionLabels.some((label) => button.textContent?.includes(label)),
    );
    expect(restoreButtons.length).toBe(1);
    expect((restoreButtons[0] as HTMLButtonElement).disabled).toBe(true);
    // Each functional-unit tab's save button is disabled for the viewer.
    for (const [tabLabel, zhTabLabel] of tabLabels) {
      const tab = [...container.querySelectorAll("button[role=tab]")].find((button) =>
        button.textContent?.includes(tabLabel) || button.textContent?.includes(zhTabLabel),
      );
      expect(tab, "tab " + tabLabel).not.toBeUndefined();
      await act(async () => {
        (tab as HTMLButtonElement).click();
      });
      const saveButtons = [...container.querySelectorAll("button")].filter((button) =>
        saveLabels.some((label) => button.textContent?.includes(label)),
      );
      expect(saveButtons.length).toBe(1);
      expect((saveButtons[0] as HTMLButtonElement).disabled).toBe(true);
    }
    // Read-only viewer still sees the current values on the General tab,
    // with fields disabled (switch back from the last visited tab).
    await act(async () => {
      const generalTab = [...container.querySelectorAll("button[role=tab]")].find((button) =>
        button.textContent?.includes("General") || button.textContent?.includes("常规"),
      )!;
      (generalTab as HTMLButtonElement).click();
    });
    const titleInput = container.querySelector<HTMLInputElement>("#field-siteTitle");
    expect(titleInput?.value).toBe("Acme Admin");
    expect(titleInput?.disabled).toBe(true);
  });

  it("restore defaults runs the reset request after confirmation", async () => {
    const resets: string[] = [];
    const base = fetcherFor();
    const tracking: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (init?.method === "POST" && url.endsWith("/reset")) {
        resets.push(url);
        return new Response(JSON.stringify({ ...SETTINGS_ROW, defaultLocale: "auto" }), {
          status: 200,
          headers: { "X-Schema-UI-Config-Changed": "settings.branding", "Content-Type": "application/json" },
        });
      }
      return base(input, init);
    }) as typeof fetch;
    const container = await renderAt(
      "/settings",
      <I18nProvider>
        <App
          manifest={adminManifest()}
          navigationContext={adminContext()}
          schemaFetcher={tracking}
          resourceFetcher={tracking}
        />
      </I18nProvider>,
    );
    await act(async () => {
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    const resetButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("Restore defaults"),
    )!;
    await act(async () => {
      resetButton.click();
    });
    expect(container.textContent).toContain("Restore all settings to their defaults?");
    const confirmButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("Confirm"),
    )!;
    await act(async () => {
      confirmButton.click();
      await new Promise((resolvePromise) => setTimeout(resolvePromise, 0));
    });
    expect(resets).toHaveLength(1);
    expect(resets[0]).toContain("/api/settings/default/reset");
  });
});
