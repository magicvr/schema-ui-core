import { readFileSync, readdirSync, statSync } from "node:fs";
import { dirname, join, resolve, sep } from "node:path";
import { fileURLToPath } from "node:url";

import { describe, expect, it } from "vitest";

import { createManifestSearchProvider, type SearchableProviderContext } from "@/app/searchable";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const ROOT = resolve(dirname(fileURLToPath(import.meta.url)), "../../../..");
const MODULE_ROOT = join(ROOT, "apps/api/modules");

const BASE_PAGES = ["dashboard", "users", "roles", "account", "notifications", "users-invites"];
const ADMIN_PAGES = [
  "settings",
  "mail",
  "mail-outbox",
  "activity",
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
  "wallet-vouchers",
];
const DEMO_PAGES = [
  "overview",
  "data-table",
  "admin-list-batch",
  "data-display",
  "search-form-table",
  "form-controls",
  "form-with-reactions",
  "form-with-upload",
];
const OPTIONAL_PAGES = ["telegram-settings", "telegram-operator", "digitaloffer-offers", "digitaloffer-entitlements", "digitaloffer-purchases"];

const MENU_FEATURES: Record<string, string> = {
  users: "menu_users",
  roles: "menu_roles",
  dashboard: "menu_dashboard",
  account: "menu_account",
  activity: "menu_activity",
  settings: "menu_settings",
  mail: "menu_mail",
  "mail-outbox": "menu_mail_outbox",
  "file-library": "menu_files",
  "data-dictionary": "menu_dictionary",
  "system-monitoring": "menu_monitoring",
  "scheduled-tasks": "menu_scheduled_tasks",
  "recycle-bin": "menu_recycle_bin",
  "data-permission": "menu_data_permission",
  wallet: "menu_wallet",
  "wallet-vouchers": "menu_wallet_vouchers",
  "my-wallet": "menu_wallet_self",
  "telegram-settings": "menu_telegram",
  "digitaloffer-offers": "menu_digitaloffer_offers",
  "digitaloffer-entitlements": "menu_digitaloffer_entitlements",
  "digitaloffer-purchases": "menu_digitaloffer_purchases",
};

function pageEntry(pageId: string) {
  const routes: Record<string, string> = {
    "dictionary-entries": "/dictionary-entries/{dictKey}",
    "wallet-entries": "/wallet-entries/{id}",
    "telegram-operator": "/telegram-settings/operator",
  };
  return {
    pageId,
    title: pageId,
    schemaUrl: `/api/schema/${pageId}`,
    route: routes[pageId] ?? `/${pageId}`,
  };
}

function navLink(pageId: string) {
  const feature = MENU_FEATURES[pageId];
  return {
    pageRef: pageId,
    label: pageId,
    ...(feature === undefined ? {} : { visibleWhen: { when: `$context.features.${feature} == true` } }),
  };
}

function profileManifest(profile: "mvp" | "admin" | "demo" | "custom"): AppManifest {
  const pages =
    profile === "mvp"
      ? [...BASE_PAGES]
      : profile === "admin"
        ? [...BASE_PAGES, ...ADMIN_PAGES]
        : profile === "demo"
          ? [...BASE_PAGES, ...DEMO_PAGES]
          : [...BASE_PAGES, ...ADMIN_PAGES, ...OPTIONAL_PAGES];
  const sidebarPageIds =
    profile === "mvp"
      ? ["dashboard", "users", "roles"]
      : profile === "admin"
        ? [
            "dashboard",
            "users",
            "roles",
            "activity",
            "mail",
            "mail-outbox",
            "file-library",
            "data-dictionary",
            "system-monitoring",
            "scheduled-tasks",
            "recycle-bin",
            "data-permission",
            "wallet",
            "wallet-vouchers",
          ]
        : profile === "demo"
          ? ["dashboard", "users", "roles", ...DEMO_PAGES]
          : [
              "dashboard",
              "users",
              "roles",
              "activity",
              "mail",
              "mail-outbox",
              "file-library",
              "data-dictionary",
              "system-monitoring",
              "scheduled-tasks",
              "recycle-bin",
              "data-permission",
              "wallet",
              "wallet-vouchers",
              "telegram-settings",
              "digitaloffer-offers",
              "digitaloffer-entitlements",
              "digitaloffer-purchases",
            ];
  const userPageIds = profile === "mvp" || profile === "demo" ? ["account"] : ["account", "my-wallet", "settings"];
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: `profile-${profile}`, name: profile, homePageRef: profile === "demo" ? "overview" : "dashboard" },
    pages: pages.map(pageEntry),
    navigation: {
      sidebar: sidebarPageIds.map(navLink),
      user: userPageIds.map(navLink),
    },
  });
}

function schemaDocuments(): Record<string, unknown> {
  const result: Record<string, unknown> = {};
  const walk = (directory: string) => {
    for (const entry of readdirSync(directory)) {
      const path = join(directory, entry);
      if (statSync(path).isDirectory()) {
        walk(path);
        continue;
      }
      if (!path.endsWith(".json") || !path.split(sep).includes("schema")) {
        continue;
      }
      const parsed = JSON.parse(readFileSync(path, "utf8")) as { meta?: { pageId?: unknown } };
      if (typeof parsed.meta?.pageId === "string") {
        result[parsed.meta.pageId] = parsed;
      }
    }
  };
  walk(MODULE_ROOT);
  return result;
}

function providerContext(manifest: AppManifest, documents: Record<string, unknown>): SearchableProviderContext {
  return {
    manifest,
    navigationContext: {
      user: {
        id: "profile-matrix-admin",
        roles: ["admin", "editor", "viewer"],
        permissions: [
          "users.read", "users.write", "users.invite", "roles.read", "roles.write", "roles.assign",
          "settings.read", "settings.write", "operations.read", "files.read", "files.write", "files.delete",
          "dictionary.read", "dictionary.write", "monitoring.read", "tasks.read", "tasks.write",
          "recycle.read", "recycle.write", "data-permission.read", "data-permission.write",
          "wallet.read", "wallet.write", "wallet.adjust", "wallet.voucher.issue", "data.export", "data.import",
          "users.enable", "users.disable", "telegram.operator.read", "telegram.operator.write",
          "digitaloffer.read", "digitaloffer.offer.manage", "digitaloffer.entitlement.void",
        ],
      },
      features: Object.fromEntries(Object.values(MENU_FEATURES).map((feature) => [feature, true])),
    },
    currentPath: manifest.app.homePageRef === "overview" ? "/overview" : "/dashboard",
    t: (key, _params, fallback) => fallback ?? key,
    loadPage: async (page) => {
      const document = documents[page.pageId];
      if (document === undefined) {
        throw new Error(`schema missing for ${page.pageId}`);
      }
      return document;
    },
    includeShellNotifications: true,
  };
}

describe("VP-036 profile denominator matrix", () => {
  const documents = schemaDocuments();
  const cases = [
    { name: "mvp" as const, pages: 5, actions: 6 },
    { name: "admin" as const, pages: 18, actions: 16 },
    { name: "demo" as const, pages: 13, actions: 6 },
    { name: "custom" as const, pages: 22, actions: 18 },
  ];

  for (const testCase of cases) {
    it(`${testCase.name} exposes only the corrected page/action denominator`, async () => {
      const result = await createManifestSearchProvider().getItems(
        providerContext(profileManifest(testCase.name), documents),
      );
      const pageItems = result.items.filter((item) => item.kind === "page" || item.kind === "navigation");
      const actionItems = result.items.filter((item) => item.kind === "action");
      expect(pageItems).toHaveLength(testCase.pages);
      expect(actionItems).toHaveLength(testCase.actions);
      expect(pageItems.some((item) => item.id.includes("dictionary-entries"))).toBe(false);
      expect(pageItems.some((item) => item.id.includes("wallet-entries"))).toBe(false);
      expect(actionItems.every((item) => item.action?.trigger.requiresSelection !== true)).toBe(true);
      expect(result.errors ?? []).toEqual([]);
    });
  }
});
