import { readFileSync } from "node:fs";

import { describe, expect, it, vi } from "vitest";

import {
  APP_MANIFEST_PROTOCOL_VERSION,
  DEFAULT_MANIFEST_PATH,
  ManifestError,
  type AppManifest,
  type PageEntry,
  evaluateExpression,
  loadAppManifest,
  matchRoute,
  pageIdMatches,
  resolveInitialRoute,
  resolveLogoUrl,
  resolveSchemaUrl,
  validateAppManifest,
} from "@/protocol/app-manifest";

const pages: PageEntry[] = [
  {
    pageId: "home",
    title: "Home",
    schemaUrl: "/schema/home",
    route: "/home",
  },
  {
    pageId: "orders",
    title: "Orders",
    schemaUrl: "/schema/orders",
    route: "/orders",
  },
  {
    pageId: "orders-new",
    title: "New Order",
    schemaUrl: "/schema/orders-new",
    route: "/orders/new",
  },
  {
    pageId: "orders-detail",
    title: "Order Detail",
    schemaUrl: "/schema/orders/{id}",
    route: "/orders/{id}",
  },
];

const adminManifestFixtureBytes = readFileSync(
  new URL("../test-fixtures/app-manifest.admin-dogfood.json", import.meta.url),
);

function manifest(overrides: Record<string, unknown> = {}): Record<string, unknown> {
  return {
    protocolVersion: APP_MANIFEST_PROTOCOL_VERSION,
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: {
      appId: "demo-app",
      name: "Demo",
      homePageRef: "home",
    },
    pages,
    navigation: {
      top: [{ pageRef: "home", label: "Home" }],
      sidebar: [{ pageRef: "orders", label: "Orders" }],
    },
    ...overrides,
  };
}

function thrown(action: () => unknown): unknown {
  try {
    action();
  } catch (error) {
    return error;
  }
  throw new Error("Expected action to throw.");
}

function expectCode(action: () => unknown, code: string) {
  expect(thrown(action)).toMatchObject({ code });
}

describe("app manifest validation", () => {
  it("accepts the pinned default manifest shape", () => {
    const result = validateAppManifest(manifest());
    expect(result.protocolVersion).toBe(APP_MANIFEST_PROTOCOL_VERSION);
    expect(result.pages).toHaveLength(4);
    expect(result.navigation?.sidebar).toHaveLength(1);
  });

  it("rejects unknown fields and missing navigation capability", () => {
    expectCode(
      () => validateAppManifest({ ...manifest(), extra: true }),
      "UNKNOWN_MANIFEST_FIELD",
    );
    expectCode(
      () =>
        validateAppManifest({
          ...manifest(),
          requiredCapabilities: ["app.manifest"],
        }),
      "CAPABILITY_REQUIRED",
    );
  });

  it("rejects unsupported versions, unsafe home routes, and empty-page refs", () => {
    expectCode(
      () => validateAppManifest({ ...manifest(), protocolVersion: "2.4" }),
      "PROTOCOL_VERSION_TOO_LOW",
    );
    expectCode(
      () => validateAppManifest({ ...manifest(), protocolVersion: "2.6" }),
      "UNSUPPORTED_PROTOCOL_VERSION",
    );
    expectCode(
      () =>
        validateAppManifest({
          ...manifest(),
          app: { appId: "demo-app", name: "Demo", homePageRef: "orders-detail" },
        }),
      "MANIFEST_HOME_ROUTE_PARAMETRIC",
    );
    expectCode(
      () =>
        validateAppManifest({
          ...manifest(),
          pages: [],
          app: { appId: "demo-app", name: "Demo", homePageRef: "home" },
          navigation: undefined,
          requiredCapabilities: ["app.manifest"],
        }),
      "PAGE_REF_WITH_EMPTY_PAGES",
    );
  });

  it("rejects navigation mutexes, nested groups, and unknown page refs", () => {
    expectCode(
      () =>
        validateAppManifest({
          ...manifest(),
          navigation: { sidebar: [{ pageRef: "home", url: "/home", label: "Home" }] },
        }),
      "NAV_LINK_MUTEX",
    );
    expectCode(
      () =>
        validateAppManifest({
          ...manifest(),
          navigation: {
            sidebar: [
              {
                label: "Outer",
                items: [{ label: "Inner", items: [{ pageRef: "home", label: "Home" }] }],
              },
            ],
          },
        }),
      "NAV_GROUP_NESTED",
    );
    expectCode(
      () =>
        validateAppManifest({
          ...manifest(),
          navigation: { sidebar: [{ pageRef: "missing", label: "Missing" }] },
        }),
      "NAV_PAGE_REF_UNKNOWN",
    );
  });
});

describe("D4a route and path semantics", () => {
  it("prefers literals, then route length, then declaration order", () => {
    expect(matchRoute(pages, "/orders/new")?.page.pageId).toBe("orders-new");
    expect(matchRoute(pages, "/orders/42")?.params).toEqual({ id: "42" });

    const sameShape: PageEntry[] = [
      { pageId: "shorter", title: "S", schemaUrl: "/s", route: "/x/{id}" },
      { pageId: "longer", title: "L", schemaUrl: "/s", route: "/x/{oid}" },
    ];
    expect(matchRoute(sameShape, "/x/1")?.page.pageId).toBe("longer");

    const declarationOrder: PageEntry[] = [
      { pageId: "first", title: "F", schemaUrl: "/s", route: "/{a}/{b}" },
      { pageId: "second", title: "S", schemaUrl: "/s", route: "/{c}/{d}" },
    ];
    expect(matchRoute(declarationOrder, "/x/y")?.page.pageId).toBe("first");
  });

  it("decodes each segment, keeps plus literal, and rejects bad escapes", () => {
    const encoded: PageEntry = {
      pageId: "item",
      title: "Item",
      schemaUrl: "/schema/item/{id}",
      route: "/items/{id}",
    };
    expect(matchRoute([encoded], "/items/a%2Bb")?.params).toEqual({ id: "a+b" });
    expect(matchRoute([encoded], "/items/a+b")?.params).toEqual({ id: "a+b" });
    expect(matchRoute([encoded], "/items/%ZZ")).toBeUndefined();
    expect(matchRoute([encoded], "/items/a/")).toBeUndefined();
  });

  it("resolves deep links before home and ignores query for matching", () => {
    const validated = validateAppManifest(manifest()) as AppManifest;
    expect(resolveInitialRoute(validated, "/")?.path).toBe("/home");
    expect(resolveInitialRoute(validated, "/orders/new?source=nav")).toMatchObject({
      path: "/orders/new",
      source: "deepLink",
      query: { source: "nav" },
    });
    expect(resolveInitialRoute(validated, "/missing")).toBeUndefined();
  });

  it("joins API paths and preserves HTTPS logos", () => {
    expect(resolveSchemaUrl("https://api.example.com/v1", "/schema/orders", {})).toBe(
      "https://api.example.com/v1/schema/orders",
    );
    expect(
      resolveSchemaUrl("https://api.example.com", "/schema/orders/{id}", { id: "42" }),
    ).toBe("https://api.example.com/schema/orders/42");
    expect(resolveLogoUrl("https://api.example.com/v1", "/assets/logo.svg")).toBe(
      "https://api.example.com/v1/assets/logo.svg",
    );
    expect(resolveLogoUrl("https://api.example.com", "https://cdn.example.com/logo.svg")).toBe(
      "https://cdn.example.com/logo.svg",
    );
    expectCode(
      () => resolveSchemaUrl("https://api.example.com", "/schema/orders/{id}", {}),
      "MISSING_PATH_BINDING",
    );
    expectCode(
      () => resolveLogoUrl("https://api.example.com", "http://cdn.example.com/logo.svg"),
      "INVALID_LOGO_URL",
    );
  });
});

describe("manifest loading and expression boundaries", () => {
  it("loads and validates one snapshot from the default fetch path", async () => {
    let requested: RequestInfo | URL | undefined;
    const loaded = await loadAppManifest({
      fetcher: async (input) => {
        requested = input;
        return new Response(JSON.stringify(manifest()), { status: 200 });
      },
    });
    expect(requested).toBe("/.well-known/schema-ui/app-manifest.json");
    expect(loaded.app.appId).toBe("demo-app");
  });

  it("validates the admin-profile test fixture through the loader", async () => {
    let requested: RequestInfo | URL | undefined;
    const loaded = await loadAppManifest({
      fetcher: async (input) => {
        requested = input;
        return new Response(adminManifestFixtureBytes, { status: 200 });
      },
    });
    expect(requested).toBe(DEFAULT_MANIFEST_PATH);
    expect(loaded.app.appId).toBe("schema-ui-core");
    // F-01 (GOAL-003): the production home is now the dashboard; the R2 wave
    // added dashboard/account/notifications pages to the fixture union.
    expect(loaded.app.homePageRef).toBe("dashboard");
    expect(loaded.pages.map((page) => page.pageId)).toEqual([
      "admin-list-batch",
      "data-display",
      "overview",
      "data-table",
      "search-form-table",
      "form-controls",
      "form-with-reactions",
      "form-with-upload",
      "users",
      "roles",
      "settings",
      "activity",
      "dashboard",
      "account",
      "notifications",
      "file-library",
      "data-dictionary",
      "dictionary-entries",
      "system-monitoring",
      "scheduled-tasks",
      "task-runs",
      "recycle-bin",
      "data-permission",
      "wallet",
      "wallet-entries",
      "my-wallet",
    ]);
  });

  it("warns when a real development response lacks the API source marker", async () => {
    const warning = vi.spyOn(console, "warn").mockImplementation(() => undefined);
    await loadAppManifest({
      fetcher: async () =>
        ({
          ok: true,
          status: 200,
          url: "http://127.0.0.1:5173/.well-known/schema-ui/app-manifest.json",
          headers: new Headers(),
          arrayBuffer: async () =>
            new TextEncoder().encode(JSON.stringify(manifest())).buffer as ArrayBuffer,
        }) as Response,
    });
    expect(warning).toHaveBeenCalledWith(
      expect.stringContaining("/.well-known/schema-ui/app-manifest.json"),
    );
    warning.mockRestore();
  });

  it("fails closed when the manifest endpoint fails", async () => {
    await expect(
      loadAppManifest({
        fetcher: async () => new Response("missing", { status: 503 }),
      }),
    ).rejects.toMatchObject({ code: "MANIFEST_LOAD_FAILED" });
  });

  it("supports only context expressions and matches their values", () => {
    expect(
      evaluateExpression('$context.user.roles contains "admin"', {
        user: { roles: ["admin"] },
      }),
    ).toBe(true);
    expect(
      evaluateExpression("$context.features.beta == true", {
        features: { beta: true },
      }),
    ).toBe(true);
    expect(evaluateExpression("$deps.admin == true", {})).toBe(false);
    expect(ManifestError).toBeDefined();
  });

  it("checks a renderer page id against the manifest entry", () => {
    expect(pageIdMatches(pages[0], "home")).toBe(true);
    expectCode(() => pageIdMatches(pages[0], "other"), "MANIFEST_PAGE_ID_MISMATCH");
  });
});