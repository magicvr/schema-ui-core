import { describe, expect, it } from "vitest";

import {
  aggregateSearchableItems,
  createManifestSearchProvider,
  filterSearchableItems,
  normalizeSearchText,
  SEARCHABLE_PROVIDER_VERSION,
  type SearchableProviderContext,
} from "@/app/searchable";
import { validateAppManifest } from "@/protocol/app-manifest";

const translate = (key: string, _params?: Record<string, string | number>, fallback?: string) =>
  fallback ?? key;

function contextFor(
  manifest: ReturnType<typeof validateAppManifest>,
  documents: Record<string, unknown>,
  overrides: Partial<SearchableProviderContext> = {},
): SearchableProviderContext {
  return {
    manifest,
    navigationContext: { user: { permissions: ["users.write"] }, features: { menu_users: true } },
    currentPath: "/home",
    t: translate,
    loadPage: async (page) => {
      const document = documents[page.pageId];
      if (document === undefined) {
        throw new Error("missing document");
      }
      return document;
    },
    ...overrides,
  };
}

function manifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "searchable", name: "Searchable", homePageRef: "home" },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
      { pageId: "users", title: "Users", schemaUrl: "/schema/users", route: "/users" },
      {
        pageId: "users-invites",
        title: "Invitations",
        schemaUrl: "/schema/users-invites",
        route: "/users-invites",
      },
      {
        pageId: "dictionary-entries",
        title: "Dictionary entries",
        schemaUrl: "/schema/dictionary-entries",
        route: "/dictionary-entries/{dictKey}",
      },
      {
        pageId: "notifications",
        title: "Notifications",
        schemaUrl: "/schema/notifications",
        route: "/notifications",
      },
    ],
    navigation: {
      top: [{ pageRef: "home", label: "Home" }],
      sidebar: [
        {
          label: "Identity",
          items: [
            {
              pageRef: "users",
              label: "Users",
              visibleWhen: { when: "$context.features.menu_users == true" },
            },
            { pageRef: "users-invites", label: "Invitations" },
            { pageRef: "dictionary-entries", label: "Entries" },
          ],
        },
      ],
    },
  });
}

function pageDocument(pageId: string, body: Record<string, unknown>, actions?: Record<string, unknown>) {
  return {
    meta: {
      pageId,
      title: pageId,
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
    },
    ...(actions === undefined ? {} : { actions }),
    body,
  };
}

describe("SearchableItem matching and aggregation", () => {
  it("normalizes case, accents, and whitespace", () => {
    expect(normalizeSearchText("  Crème   Brûlée ")).toBe("creme brulee");
  });

  it("scores exact/prefix/keyword matches and caps deterministically", () => {
    const items = [
      { id: "z", kind: "page" as const, label: "User settings", rank: 2 },
      { id: "a", kind: "page" as const, label: "Users", rank: 1 },
      { id: "b", kind: "action" as const, label: "Create user", keywords: ["users"], rank: 3 },
      { id: "c", kind: "page" as const, label: "Account", keywords: ["users"], rank: 4 },
    ];
    expect(filterSearchableItems(items, " users ", 3).map((item) => item.id)).toEqual(["a", "b", "c"]);
    expect(filterSearchableItems(items, "", 2).map((item) => item.id)).toEqual(["a", "z"]);
    expect(filterSearchableItems(items, "users", 0)).toEqual([]);
  });

  it("removes both sides of duplicate ids and reports provider failures", async () => {
    const base = contextFor(manifest(), {});
    const result = await aggregateSearchableItems(
      [
        {
          id: "b-provider",
          version: SEARCHABLE_PROVIDER_VERSION,
          getItems: () => ({ items: [{ id: "duplicate", kind: "page", label: "B" }] }),
        },
        {
          id: "a-provider",
          version: SEARCHABLE_PROVIDER_VERSION,
          getItems: () => ({ items: [{ id: "duplicate", kind: "page", label: "A" }] }),
        },
        {
          id: "failing",
          version: SEARCHABLE_PROVIDER_VERSION,
          getItems: () => {
            throw new Error("boom");
          },
        },
      ],
      base,
    );
    expect(result.items).toEqual([]);
    expect(result.errors.map((error) => error.code)).toEqual(["DUPLICATE_ITEM_ID", "PROVIDER_FAILED"]);
  });
});

describe("manifest SearchableProvider", () => {
  it("uses projected visible navigation, excludes dynamic/unlinked pages, and indexes direct actions", async () => {
    const appManifest = manifest();
    const documents = {
      home: pageDocument("home", { type: "section", children: [] }),
      users: pageDocument(
        "users",
        {
          type: "section",
          children: [
            {
              type: "table",
              id: "users-table",
              props: {
                toolbar: [
                  { key: "create", label: "Create user", actionRef: "openCreate" },
                  { key: "batch", label: "Batch", actionRef: "batch", requiresSelection: true },
                ],
                actions: [{ key: "row", label: "Row action", actionRef: "row" }],
              },
            },
            {
              type: "actionButton",
              id: "open-users-help",
              props: { label: "Help", actionId: "openHelp" },
            },
          ],
        },
        {
          openCreate: { type: "modal", content: { type: "section", children: [] } },
          openHelp: { type: "navigate", url: "/home" },
          batch: { type: "request", method: "POST", url: "/api/batch" },
          row: { type: "request", method: "POST", url: "/api/row/{id}" },
        },
      ),
      "users-invites": pageDocument("users-invites", { type: "section", children: [] }),
      notifications: pageDocument("notifications", { type: "section", children: [] }),
    };
    const provider = createManifestSearchProvider();
    const result = await provider.getItems(
      contextFor(appManifest, documents, { includeShellNotifications: true }),
    );

    expect(result.errors ?? []).toEqual([]);
    expect(result.items.filter((item) => item.kind !== "action").map((item) => item.id)).toEqual([
      "page:home",
      "page:users",
      "page:users-invites",
      "page:notifications",
    ]);
    expect(result.items.filter((item) => item.kind === "action").map((item) => item.id)).toEqual([
      "action:users:create",
      "action:users:open-users-help",
    ]);
    expect(result.items.some((item) => item.id === "page:dictionary-entries")).toBe(false);
  });

  it("fails closed for a denied mounted action while keeping its visible page", async () => {
    const appManifest = manifest();
    const denied = pageDocument(
      "users",
      {
        type: "section",
        children: [
          {
            type: "table",
            id: "users-table",
            permissionCascade: { keys: ["edit"] },
            permissions: { edit: '$context.user.permissions contains "users.write"' },
            props: {
              toolbar: [
                { key: "create", label: "Create user", actionRef: "openCreate", permissionIntent: "edit" },
              ],
            },
          },
        ],
      },
      { openCreate: { type: "modal", content: { type: "section", children: [] } } },
    );
    const result = await createManifestSearchProvider().getItems(
      contextFor(appManifest, { users: denied }, {
        navigationContext: { user: { permissions: [] }, features: { menu_users: true } },
      }),
    );
    expect(result.items.some((item) => item.id === "page:users")).toBe(true);
    expect(result.items.some((item) => item.id === "action:users:create")).toBe(false);
  });
});
