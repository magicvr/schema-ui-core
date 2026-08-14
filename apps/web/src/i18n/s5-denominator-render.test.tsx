// @vitest-environment jsdom

/**
 * S5 · F-V029 denominator runtime bilingual render (GOAL-006 A-001 F-002/F-003).
 *
 * Drives real shipped App + catalogs against real module schema fixtures so
 * structural key completeness is not silently treated as runtime bilingual
 * render for the remaining pageId union (esp. roles + example pages).
 * Also asserts mvp Profile boundary: settings/activity pageIds absent and
 * /settings is not a reachable schema surface.
 */

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import { resetMissingTranslationReports } from "@/i18n/catalog";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const __dir = dirname(fileURLToPath(import.meta.url));
const ADMIN_MANIFEST = resolve(__dir, "../test-fixtures/app-manifest.admin.json");
const MVP_MANIFEST = resolve(__dir, "../test-fixtures/app-manifest.mvp.json");
const MODULES = resolve(__dir, "../../../api/internal/modules");

const SCHEMA_PATHS: Record<string, string> = {
  overview: resolve(MODULES, "dev/examples/schema/overview.json"),
  "admin-list-batch": resolve(MODULES, "dev/examples/schema/admin-list-batch.json"),
  "data-display": resolve(MODULES, "dev/examples/schema/data-display.json"),
  "data-table": resolve(MODULES, "dev/examples/schema/data-table.json"),
  "search-form-table": resolve(MODULES, "dev/examples/schema/search-form-table.json"),
  "form-controls": resolve(MODULES, "dev/examples/schema/form-controls.json"),
  "form-with-reactions": resolve(MODULES, "dev/examples/schema/form-with-reactions.json"),
  "form-with-upload": resolve(MODULES, "dev/examples/schema/form-with-upload.json"),
  users: resolve(MODULES, "users/schema/users.json"),
  roles: resolve(MODULES, "roles/schema/roles.json"),
  settings: resolve(MODULES, "settings/schema/settings.json"),
  activity: resolve(MODULES, "activity/schema/activity.json"),
  dashboard: resolve(MODULES, "dashboard/schema/dashboard.json"),
  account: resolve(MODULES, "account/schema/account.json"),
  notifications: resolve(MODULES, "notifications/schema/notifications.json"),
  "file-library": resolve(MODULES, "filelibrary/schema/file-library.json"),
  "data-dictionary": resolve(MODULES, "datadictionary/schema/data-dictionary.json"),
  "dictionary-entries": resolve(MODULES, "datadictionary/schema/dictionary-entries.json"),
  "system-monitoring": resolve(MODULES, "systemmonitoring/schema/system-monitoring.json"),
  "scheduled-tasks": resolve(MODULES, "scheduledtasks/schema/scheduled-tasks.json"),
  "recycle-bin": resolve(MODULES, "recyclebin/schema/recycle-bin.json"),
  "task-runs": resolve(MODULES, "scheduledtasks/schema/task-runs.json"),
};

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
  resetMissingTranslationReports();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
  resetMissingTranslationReports();
});

function loadManifest(path: string): AppManifest {
  return validateAppManifest(JSON.parse(readFileSync(path, "utf8")));
}

function loadSchema(pageId: string): Record<string, unknown> {
  return JSON.parse(readFileSync(SCHEMA_PATHS[pageId], "utf8"));
}

function allDocuments(): Record<string, unknown> {
  const docs: Record<string, unknown> = {};
  for (const pageId of Object.keys(SCHEMA_PATHS)) {
    docs[pageId] = loadSchema(pageId);
  }
  return docs;
}

const ROLES_LIST = {
  items: [
    {
      id: "role-1",
      key: "admin",
      name: "Administrator",
      permissions: ["users.read"],
      menuItems: ["users"],
      isSystem: true,
      userCount: 1,
      updatedAt: "2026-08-03T00:00:00.000Z",
    },
  ],
  total: 1,
  page: 1,
  pageSize: 10,
};

const EMPTY_LIST = { items: [], total: 0, page: 1, pageSize: 10 };

function fetcherFor(documents: Record<string, unknown>): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    if (pathname.startsWith("/api/schema/")) {
      const pageId = pathname.slice("/api/schema/".length);
      const document = documents[pageId];
      if (document === undefined) {
        return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), { status: 404 });
      }
      return new Response(JSON.stringify(document), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/roles")) {
      return new Response(JSON.stringify(ROLES_LIST), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/users") || pathname.startsWith("/api/operations")) {
      return new Response(JSON.stringify(EMPTY_LIST), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname === "/api/settings/default") {
      return new Response(
        JSON.stringify({
          id: "default",
          siteTitle: "Schema UI Core",
          logoUrl: "",
          logoUrlLight: "",
          logoUrlDark: "",
          faviconUrl: "",
          defaultLocale: "auto",
          siteTimezone: "auto",
          defaultTheme: "auto",
          updatedAt: "2026-08-09T00:00:00.000Z",
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    if (pathname === "/api/branding") {
      return new Response(
        JSON.stringify({
          siteTitle: "Schema UI Core",
          logoUrl: "",
          defaultLocale: "auto",
          supportedLocales: ["zh-CN", "en-US"],
          siteTimezone: "auto",
          defaultTheme: "auto",
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
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
  await act(async () => {
    await new Promise((r) => setTimeout(r, 0));
  });
  return container;
}

function shell(children: React.ReactNode, locale: "zh-CN" | "en-US") {
  return (
    <I18nProvider stored={locale} browserLanguages={["en-US"]}>
      {children}
    </I18nProvider>
  );
}

function adminContext(features: Record<string, boolean> = {}) {
  return {
    user: {
      id: "u1",
      roles: ["admin"],
      permissions: [
        "users.read",
        "users.write",
        "roles.read",
        "roles.write",
        "settings.read",
        "settings.write",
        "operations.read",
      ],
    },
    features: {
      menu_users: true,
      menu_roles: true,
      menu_activity: true,
      menu_settings: true,
      ...features,
    },
  };
}

// ── F-002 · roles + example page runtime bilingual render ───────────────────

describe("S5 · denominator runtime bilingual render (F-002)", () => {
  it("renders roles schema labels via labelKey under zh-CN and en-US", async () => {
    const documents = allDocuments();
    const manifest = loadManifest(ADMIN_MANIFEST);
    const fetchImpl = fetcherFor(documents);

    const zh = await renderAt(
      "/roles",
      shell(
        <App
          manifest={manifest}
          navigationContext={adminContext()}
          schemaFetcher={fetchImpl}
          resourceFetcher={fetchImpl}
        />,
        "zh-CN",
      ),
    );
    const zhText = zh.textContent ?? "";
    expect(zh.querySelector("h1")?.textContent).toContain("角色");
    expect(zhText).toContain("新建角色");
    expect(zhText).toContain("名称");
    expect(zhText).toContain("键");

    const en = await renderAt(
      "/roles",
      shell(
        <App
          manifest={manifest}
          navigationContext={adminContext()}
          schemaFetcher={fetchImpl}
          resourceFetcher={fetchImpl}
        />,
        "en-US",
      ),
    );
    const enText = en.textContent ?? "";
    expect(en.querySelector("h1")?.textContent).toContain("Roles");
    expect(enText).toContain("New role");
    expect(enText).toContain("Name");
    expect(enText).toContain("Key");
  });

  it("resolves titleKey for remaining union pages under both locales", async () => {
    const documents = allDocuments();
    const manifest = loadManifest(ADMIN_MANIFEST);
    const fetchImpl = fetcherFor(documents);

    const cases: Array<{ route: string; en: string; zh: string }> = [
      { route: "/overview", en: "Overview", zh: "总览" },
      { route: "/data-table", en: "Data table", zh: "数据表格" },
      { route: "/data-display", en: "Data display", zh: "数据展示" },
      { route: "/search-form-table", en: "Search + table", zh: "搜索 + 表格" },
      { route: "/form-controls", en: "Form controls", zh: "表单控件" },
      { route: "/form-with-reactions", en: "Form with reactions", zh: "联动表单" },
      { route: "/form-with-upload", en: "Form with upload", zh: "上传表单" },
      { route: "/admin-list-batch", en: "List + batch", zh: "列表 + 批量" },
      { route: "/activity", en: "Activity", zh: "操作日志" },
    ];

    for (const entry of cases) {
      const en = await renderAt(
        entry.route,
        shell(
          <App
            manifest={manifest}
            navigationContext={adminContext()}
            schemaFetcher={fetchImpl}
            resourceFetcher={fetchImpl}
          />,
          "en-US",
        ),
      );
      expect(en.querySelector("h1")?.textContent ?? "", `en ${entry.route}`).toContain(entry.en);

      const zh = await renderAt(
        entry.route,
        shell(
          <App
            manifest={manifest}
            navigationContext={adminContext()}
            schemaFetcher={fetchImpl}
            resourceFetcher={fetchImpl}
          />,
          "zh-CN",
        ),
      );
      expect(zh.querySelector("h1")?.textContent ?? "", `zh ${entry.route}`).toContain(entry.zh);
    }
  });

  it("renders activity intro/columns under zh-CN (schema body, not only titleKey)", async () => {
    const documents = allDocuments();
    const manifest = loadManifest(ADMIN_MANIFEST);
    const fetchImpl = fetcherFor(documents);
    const container = await renderAt(
      "/activity",
      shell(
        <App
          manifest={manifest}
          navigationContext={adminContext()}
          schemaFetcher={fetchImpl}
          resourceFetcher={fetchImpl}
        />,
        "zh-CN",
      ),
    );
    const text = container.textContent ?? "";
    expect(text).toMatch(/操作日志|只读/);
    expect(text).toMatch(/事件|时间|操作/);
  });
});

// ── F-003 · mvp Profile boundary ────────────────────────────────────────────

describe("S5 · mvp Profile boundary (F-003)", () => {
  it("mvp manifest excludes settings/activity pageIds; admin includes them", () => {
    const mvp = loadManifest(MVP_MANIFEST);
    const admin = loadManifest(ADMIN_MANIFEST);
    const mvpIds = mvp.pages.map((p) => p.pageId);
    const adminIds = admin.pages.map((p) => p.pageId);
    expect(mvpIds).not.toContain("settings");
    expect(mvpIds).not.toContain("activity");
    expect(mvpIds).toContain("users");
    expect(mvpIds).toContain("roles");
    expect(adminIds).toContain("settings");
    expect(adminIds).toContain("activity");
  });

  it("mvp /settings is not a schema-driven page; admin /settings renders bilingual surface", async () => {
    const documents = allDocuments();
    const fetchImpl = fetcherFor(documents);

    const mvp = await renderAt(
      "/settings",
      shell(
        <App
          manifest={loadManifest(MVP_MANIFEST)}
          navigationContext={adminContext({ menu_settings: true })}
          schemaFetcher={fetchImpl}
          resourceFetcher={fetchImpl}
        />,
        "zh-CN",
      ),
    );
    // No settings page in mvp manifest → route fallback / page not found chrome.
    expect(mvp.textContent ?? "").not.toContain("站点标题");
    expect(mvp.textContent ?? "").toMatch(/找不到|not found|Page not found|页面/i);

    const admin = await renderAt(
      "/settings",
      shell(
        <App
          manifest={loadManifest(ADMIN_MANIFEST)}
          navigationContext={adminContext()}
          schemaFetcher={fetchImpl}
          resourceFetcher={fetchImpl}
        />,
        "zh-CN",
      ),
    );
    expect(admin.querySelector("h1")?.textContent).toContain("设置");
    expect(admin.textContent ?? "").toMatch(/常规|站点标题/);
  });
});