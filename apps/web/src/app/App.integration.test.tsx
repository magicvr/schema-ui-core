// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import {
  CONFIG_CHANGED_HEADER,
  createConfigAwareFetcher,
  SETTINGS_BRANDING_NAMESPACE,
} from "@/app/config-events";
import { ManifestFailure } from "@/app/ManifestFailure";
import {
  ManifestError,
  type NavigationContext,
  validateAppManifest,
} from "@/protocol/app-manifest";

function testManifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "integration", name: "Integration", homePageRef: "home" },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
      {
        pageId: "catalog",
        title: "Catalog",
        schemaUrl: "/schema/catalog",
        route: "/catalog",
      },
      {
        pageId: "catalog-detail",
        title: "Catalog detail",
        schemaUrl: "/schema/catalog/{id}",
        route: "/catalog/{id}",
      },
    ],
    navigation: {
      top: [{ pageRef: "home", label: "Home" }],
      sidebar: [
        { pageRef: "catalog", label: "Catalog" },
        {
          label: "Operations",
          items: [
            {
              pageRef: "catalog-detail",
              label: "Detail",
              permissions: { view: '$context.user.roles contains "admin"' },
            },
          ],
        },
      ],
    },
  });
}

// Structurally valid page document, mirroring the pinned page/node schema shape.
function schemaDocument(pageId: string, title: string, text: string) {
  return {
    meta: {
      pageId,
      title,
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
    },
    body: {
      type: "section",
      children: [{ type: "text", props: { text } }],
    },
  };
}

// Keyed by resolved schemaUrl pathname; a missing key 404s (fail-closed).
const DEFAULT_DOCUMENTS: Record<string, unknown> = {
  "/schema/home": schemaDocument("home", "Home", "Schema home body"),
  "/schema/catalog": schemaDocument("catalog", "Catalog", "Schema catalog body"),
  "/schema/catalog/42": schemaDocument("catalog-detail", "Catalog detail", "Schema detail body"),
};

function schemaFetcher(documents: Record<string, unknown>) {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    const document = documents[pathname];
    if (document === undefined) {
      return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), {
        status: 404,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response(JSON.stringify(document), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderApp(
  path: string,
  navigationContext?: NavigationContext,
  documents: Record<string, unknown> = DEFAULT_DOCUMENTS,
  resourceFetcher?: typeof fetch,
  manifest = testManifest(),
): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <App
          manifest={manifest}
          navigationContext={navigationContext}
          schemaFetcher={schemaFetcher(documents)}
          resourceFetcher={resourceFetcher}
        />
      </I18nProvider>,
    );
  });
  return container;
}

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
});

describe("App shell integration", () => {
  it("redirects root to home, navigates with history, and handles popstate", async () => {
    const container = await renderApp("/");
    expect(window.location.pathname).toBe("/home");
    expect(container.querySelector("h1")?.textContent).toBe("Home");
    expect(container.textContent).toContain("Schema home body");

    const catalogLink = container.querySelector('a[href="/catalog"]');
    expect(catalogLink).not.toBeNull();
    await act(async () => {
      catalogLink?.dispatchEvent(new MouseEvent("click", { bubbles: true, button: 0 }));
    });
    expect(window.location.pathname).toBe("/catalog");
    expect(container.querySelector("h1")?.textContent).toBe("Catalog");
    expect(container.textContent).toContain("Schema catalog body");

    window.history.pushState({}, "", "/home");
    await act(async () => {
      window.dispatchEvent(new PopStateEvent("popstate"));
    });
    expect(container.querySelector("h1")?.textContent).toBe("Home");
    expect(container.textContent).toContain("Schema home body");
  });


  // GOAL-015 semantic breadcrumbs (user ruling 2026-08-14): the trail is
  // 首页 => 一级页 => ... => n级内页 — the home page (manifest homePageRef)
  // always leads, then nav group labels, then the current page. It is
  // hierarchy, not visit history: a deep link shows the same trail.
  it("shows the semantic breadcrumb trail (首页 => 组 => 内页) on nested pages", async () => {
    const container = await renderApp("/");
    // Home is the root — single level, no trail UI, no back button.
    expect(container.querySelector('nav[aria-label="Breadcrumb"]')).toBeNull();

    // Deep-link straight into the nested page (no prior visits).
    window.history.pushState({}, "", "/catalog/42");
    await act(async () => {
      window.dispatchEvent(new PopStateEvent("popstate"));
    });
    expect(container.querySelector("h1")?.textContent).toBe("Catalog detail");

    // Semantic trail: Home (root) / Operations (group label) / Catalog detail.
    const breadcrumb = container.querySelector('nav[aria-label="Breadcrumb"]');
    expect(breadcrumb).not.toBeNull();
    expect(breadcrumb?.textContent).toContain("Home");
    expect(breadcrumb?.textContent).toContain("Operations");
    expect(breadcrumb?.textContent).toContain("Catalog detail");
    // No routed ancestor besides home (the group label is not a page) → the
    // compact back button stays hidden; 首页 itself is the clickable root.
    expect(breadcrumb?.querySelector('button[aria-label]')).toBeNull();
  });

  it("opens the Telegram operator page through an app-internal navigate action", async () => {
    const manifest = validateAppManifest({
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
      app: { appId: "integration", name: "Integration", homePageRef: "home" },
      pages: [
        { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
        {
          pageId: "telegram-settings",
          title: "Telegram channel",
          schemaUrl: "/schema/telegram-settings",
          route: "/telegram-settings",
        },
        {
          pageId: "telegram-operator",
          title: "Operator conversations",
          schemaUrl: "/schema/telegram-operator",
          route: "/telegram-settings/operator",
        },
      ],
      navigation: {
        top: [{ pageRef: "home", label: "Home" }],
        sidebar: [{ pageRef: "telegram-settings", label: "Telegram channel" }],
      },
    });
    const documents: Record<string, unknown> = {
      "/schema/home": schemaDocument("home", "Home", "Schema home body"),
      "/schema/telegram-settings": {
        meta: {
          pageId: "telegram-settings",
          title: "Telegram channel",
          protocolVersion: "2.7",
          requiredCapabilities: ["app.manifest", "app.navigation", "actions.page.trigger"],
        },
        actions: {
          openTelegramOperator: { type: "navigate", url: "/telegram-settings/operator" },
        },
        body: {
          type: "section",
          children: [{ type: "actionButton", props: { label: "Open operator conversations", actionId: "openTelegramOperator" } }],
        },
      },
      "/schema/telegram-operator": schemaDocument(
        "telegram-operator",
        "Operator conversations",
        "Operator chat body",
      ),
    };
    const container = await renderApp("/telegram-settings", {}, documents, undefined, manifest);
    const entryButton = Array.from(container.querySelectorAll("button")).find((button) =>
      button.textContent?.trim() === "Open operator conversations",
    );
    expect(entryButton).not.toBeUndefined();

    await act(async () => {
      entryButton!.click();
      await Promise.resolve();
    });

    expect(window.location.pathname).toBe("/telegram-settings/operator");
    expect(container.querySelector("h1")?.textContent).toBe("Operator conversations");
    expect(container.textContent).toContain("Operator chat body");
    const breadcrumb = container.querySelector('nav[aria-label="Breadcrumb"]');
    expect(breadcrumb?.textContent).toContain("Telegram channel");
    expect(breadcrumb?.textContent).toContain("Operator conversations");

    const backButton = breadcrumb?.querySelector('button[aria-label]') as HTMLButtonElement | null;
    expect(backButton).not.toBeNull();
    await act(async () => backButton?.click());
    expect(window.location.pathname).toBe("/telegram-settings");
  });


  // GOAL-015 F-001 (grok audit): a schema-driven navigate action (type
  // navigate + row navigateMapping, e.g. the dictionary types page 「条目」
  // row action) must navigate session-internally through the host's
  // onNavigate — the entries page then keeps its breadcrumb trail + back
  // button (no full reload that would reset the visit stack).
  it("keeps the breadcrumb trail when a row navigate action opens the inner page", async () => {
    const manifest = validateAppManifest({
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
      app: { appId: "integration", name: "Integration", homePageRef: "home" },
      pages: [
        { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
        {
          pageId: "data-dictionary",
          title: "Data dictionary",
          schemaUrl: "/schema/dictionary",
          route: "/data-dictionary",
        },
        {
          pageId: "dictionary-entries",
          title: "Dictionary entries",
          schemaUrl: "/schema/dictionary-entries",
          route: "/dictionary-entries/{dictKey}",
        },
      ],
      navigation: {
        top: [{ pageRef: "home", label: "Home" }],
        sidebar: [{ pageRef: "data-dictionary", label: "Data dictionary" }],
      },
    });
    const dictionaryDoc = {
      meta: {
        pageId: "data-dictionary",
        title: "Data dictionary",
        protocolVersion: "2.7",
        requiredCapabilities: ["app.manifest", "app.navigation", "actions.row.navigate"],
      },
      actions: {
        openEntries: { type: "navigate", url: "/dictionary-entries/{dictKey}" },
      },
      body: {
        type: "section",
        children: [
          {
            type: "table",
            id: "dict-types-table",
            data: { source: "api", url: "/api/data-dictionary/types" },
            props: {
              rowKey: "id",
              columns: [
                { field: "key", label: "Type key" },
                { field: "name", label: "Type name" },
              ],
              actions: [
                {
                  key: "entries",
                  label: "Entries",
                  actionRef: "openEntries",
                  navigateMapping: {
                    path: { dictKey: "$row.key" },
                    query: { dictTypeName: "$row.name" },
                  },
                },
              ],
            },
          },
        ],
      },
    };
    const entriesDoc = {
      meta: {
        pageId: "dictionary-entries",
        title: "Dictionary entries",
        protocolVersion: "2.7",
        requiredCapabilities: ["app.manifest", "app.navigation"],
      },
      body: {
        type: "section",
        children: [{ type: "text", props: { text: "Entries body" } }],
      },
    };
    const documents: Record<string, unknown> = {
      "/schema/home": schemaDocument("home", "Home", "Schema home body"),
      "/schema/dictionary": dictionaryDoc,
      "/schema/dictionary-entries": entriesDoc,
    };
    const resourceFetcher = (async () =>
      new Response(
        JSON.stringify({
          items: [{ id: "order_status", key: "order_status", name: "Order status" }],
          total: 1,
          page: 1,
          pageSize: 10,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      )
    ) as typeof fetch;
    const container = await renderApp("/", {}, documents, resourceFetcher, manifest);

    // Home → Data dictionary via the sidebar.
    await act(async () => {
      container.querySelector('a[href="/data-dictionary"]')?.dispatchEvent(
        new MouseEvent("click", { bubbles: true, button: 0 }),
      );
    });
    expect(container.querySelector("h1")?.textContent).toBe("Data dictionary");

    // Click the row 「条目」 action: session-internal navigate with the row
    // dictKey binding — the URL carries the query and the visit stack keeps
    // Home + Data dictionary as ancestors.
    await act(async () => {
      (Array.from(container.querySelectorAll("button")).find((button) =>
        button.textContent?.trim() === "Entries",
      ) as HTMLButtonElement).click();
    });
    expect(window.location.pathname).toBe("/dictionary-entries/order_status");
    expect(window.location.search).toBe("?dictTypeName=Order%20status");
    // Semantic trail 首页 => 一级页 => 内页 (independent of how the page was
    // reached): Home / Data dictionary / Dictionary entries.
    const breadcrumb = container.querySelector('nav[aria-label="Breadcrumb"]');
    expect(breadcrumb).not.toBeNull();
    expect(breadcrumb?.textContent).toContain("Home");
    expect(breadcrumb?.textContent).toContain("Data dictionary");
    expect(breadcrumb?.textContent).toContain("Dictionary entries");
    // The compact circular back button (semantic parent = data-dictionary).
    const backButton = breadcrumb?.querySelector('button[aria-label]') as HTMLButtonElement | null;
    expect(backButton).not.toBeNull();
    await act(async () => backButton?.click());
    expect(window.location.pathname).toBe("/data-dictionary");
  });


  // GOAL-015 F-007(b) (user ruling): the entries inner page requires the
  // type context — its route is /dictionary-entries/{dictKey}, so a deep
  // link without the param fails closed at the router (HOST_ROUTE_NOT_FOUND)
  // instead of showing an unfiltered list / broken create form.
  it("fails closed when the entries inner page is deep-linked without its dictKey param", async () => {
    const manifest = validateAppManifest({
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
      app: { appId: "integration", name: "Integration", homePageRef: "home" },
      pages: [
        { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
        {
          pageId: "dictionary-entries",
          title: "Dictionary entries",
          schemaUrl: "/schema/dictionary-entries",
          route: "/dictionary-entries/{dictKey}",
        },
      ],
      navigation: { top: [{ pageRef: "home", label: "Home" }] },
    });
    const container = await renderApp("/dictionary-entries", {}, DEFAULT_DOCUMENTS, undefined, manifest);
    expect(container.textContent).toContain("Page not found");
    expect(container.querySelector("h1")?.textContent).not.toBe("Dictionary entries");
  });

  it("renders a fail-closed fallback and returns to the manifest home route", async () => {
    const container = await renderApp("/unknown");
    // HOST_ROUTE_NOT_FOUND global failure surface (ADR-0036 D3/D3a).
    expect(container.textContent).toContain("Page not found");

    const homeButton = Array.from(container.querySelectorAll("button")).find((button) =>
      button.textContent?.includes("Return home"),
    );
    expect(homeButton).not.toBeUndefined();
    await act(async () => homeButton?.click());
    expect(window.location.pathname).toBe("/home");
    expect(container.querySelector("h1")?.textContent).toBe("Home");
  });

  it("uses the injected boot context and resolves a parametric page link", async () => {
    const container = await renderApp("/catalog/42", {
      user: { roles: ["admin"] },
      features: {},
    });
    expect(container.querySelector('a[href="/catalog/42"]')).not.toBeNull();
    expect(container.textContent).toContain("Detail");
    expect(container.textContent).toContain("Schema detail body");

    const viewer = await renderApp("/catalog", {
      user: { roles: ["viewer"] },
      features: {},
    });
    expect(viewer.textContent).not.toContain("Detail");
  });

  it("surfaces a non-blocking notice when the account session fails to load", async () => {
    window.history.replaceState({}, "", "/home");
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider>
          <App
            manifest={testManifest()}
            navigationContext={{}}
            accountError={new Error("account unavailable")}
          />
        </I18nProvider>,
      );
    });
    // Fail-closed: the shell still renders, but the failure is observable.
    expect(container.textContent).toContain("Account session failed to load");
    expect(container.querySelector("h1")?.textContent).toBe("Home");

    const healthy = document.createElement("div");
    document.body.appendChild(healthy);
    const healthyRoot = createRoot(healthy);
    activeRoots.push({ root: healthyRoot, container: healthy });
    await act(async () => {
      healthyRoot.render(
        <I18nProvider>
          <App manifest={testManifest()} navigationContext={{}} />
        </I18nProvider>,
      );
    });
    expect(healthy.textContent).not.toContain("Account session failed to load");
  });

  it("reloads branding when a generic resource response carries a config-change header", async () => {
    const originalFetch = globalThis.fetch;
    let siteTitle = "Before change";
    globalThis.fetch = (async (input: RequestInfo | URL) => {
      if (String(input) === "/api/branding") {
        return new Response(JSON.stringify({ siteTitle, logoUrl: "" }), { status: 200 });
      }
      return new Response("not found", { status: 404 });
    }) as typeof fetch;

    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    const resourceFetcher = createConfigAwareFetcher(async () =>
      new Response("{}", {
        status: 200,
        headers: { [CONFIG_CHANGED_HEADER]: SETTINGS_BRANDING_NAMESPACE },
      }),
    );

    try {
      await act(async () => {
        root.render(
          <I18nProvider>
            <App
              manifest={testManifest()}
              navigationContext={{}}
              schemaFetcher={schemaFetcher(DEFAULT_DOCUMENTS)}
              resourceFetcher={resourceFetcher}
            />
          </I18nProvider>,
        );
        await new Promise((resolve) => setTimeout(resolve, 0));
      });
      expect(container.textContent).toContain("Before change");

      siteTitle = "After change";
      await act(async () => {
        await resourceFetcher("/api/settings", { method: "PATCH" });
        await new Promise((resolve) => setTimeout(resolve, 0));
      });
      expect(container.textContent).toContain("After change");
      expect(document.title).toBe("After change");
    } finally {
      globalThis.fetch = originalFetch;
    }
  });

  it("fails closed with the unified schema error when the page document is missing", async () => {
    const container = await renderApp("/catalog", {}, {});
    // The shell header still renders the manifest title; the page body surfaces
    // the loader's PAGE_NOT_FOUND instead of a placeholder or hand-written page.
    expect(container.querySelector("h1")?.textContent).toBe("Catalog");
    expect(container.textContent).toContain("PAGE_NOT_FOUND");
    expect(container.textContent).toContain("/schema/catalog");
  });
});

describe("manifest failure surface", () => {
  it("shows the stable error code and retry affordance", async () => {
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <ManifestFailure
          error={new ManifestError("MANIFEST_LOAD_FAILED", "/manifest", "HTTP 503")}
        />,
      );
    });
    expect(container.textContent).toContain("MANIFEST_LOAD_FAILED");
    expect(container.textContent).toContain("HTTP 503");
    expect(container.textContent).toContain("Retry manifest");
  });
});
