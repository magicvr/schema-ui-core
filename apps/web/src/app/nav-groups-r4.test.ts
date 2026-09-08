import { describe, expect, it } from "vitest";

import { projectNavigation } from "@/app/navigation";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const GROUP_KEYS = {
  workspace: "manifest.nav.group.workspace",
  identity: "manifest.nav.group.identityAccess",
  content: "manifest.nav.group.contentData",
  operations: "manifest.nav.group.operations",
  communications: "manifest.nav.group.communications",
  commerce: "manifest.nav.group.commerce",
} as const;

type RouteCase = {
  path: string;
  groupKey: string;
  pageRef: string;
};

const routeCases: RouteCase[] = [
  { path: "/dashboard", groupKey: GROUP_KEYS.workspace, pageRef: "dashboard" },
  { path: "/users", groupKey: GROUP_KEYS.identity, pageRef: "users" },
  { path: "/users-invites", groupKey: GROUP_KEYS.identity, pageRef: "users" },
  { path: "/roles", groupKey: GROUP_KEYS.identity, pageRef: "roles" },
  { path: "/data-permission", groupKey: GROUP_KEYS.identity, pageRef: "data-permission" },
  { path: "/file-library", groupKey: GROUP_KEYS.content, pageRef: "file-library" },
  { path: "/data-dictionary", groupKey: GROUP_KEYS.content, pageRef: "data-dictionary" },
  { path: "/dictionary-entries/colors", groupKey: GROUP_KEYS.content, pageRef: "data-dictionary" },
  { path: "/activity", groupKey: GROUP_KEYS.operations, pageRef: "activity" },
  { path: "/system-monitoring", groupKey: GROUP_KEYS.operations, pageRef: "system-monitoring" },
  { path: "/scheduled-tasks", groupKey: GROUP_KEYS.operations, pageRef: "scheduled-tasks" },
  { path: "/task-runs", groupKey: GROUP_KEYS.operations, pageRef: "scheduled-tasks" },
  { path: "/recycle-bin", groupKey: GROUP_KEYS.operations, pageRef: "recycle-bin" },
  { path: "/mail", groupKey: GROUP_KEYS.communications, pageRef: "mail" },
  { path: "/mail-outbox", groupKey: GROUP_KEYS.communications, pageRef: "mail-outbox" },
  { path: "/telegram-settings", groupKey: GROUP_KEYS.communications, pageRef: "telegram-settings" },
  {
    path: "/telegram-settings/operator",
    groupKey: GROUP_KEYS.communications,
    pageRef: "telegram-settings",
  },
  { path: "/wallet", groupKey: GROUP_KEYS.commerce, pageRef: "wallet" },
  { path: "/wallet-entries/account-1", groupKey: GROUP_KEYS.commerce, pageRef: "wallet" },
  { path: "/wallet-vouchers", groupKey: GROUP_KEYS.commerce, pageRef: "wallet-vouchers" },
  { path: "/digitaloffer-offers", groupKey: GROUP_KEYS.commerce, pageRef: "digitaloffer-offers" },
  {
    path: "/digitaloffer-entitlements",
    groupKey: GROUP_KEYS.commerce,
    pageRef: "digitaloffer-entitlements",
  },
  {
    path: "/digitaloffer-purchases",
    groupKey: GROUP_KEYS.commerce,
    pageRef: "digitaloffer-purchases",
  },
];

function routeMatrixManifest(): AppManifest {
  const pages = [
    ["dashboard", "/dashboard"],
    ["users", "/users"],
    ["users-invites", "/users-invites"],
    ["roles", "/roles"],
    ["data-permission", "/data-permission"],
    ["file-library", "/file-library"],
    ["data-dictionary", "/data-dictionary"],
    ["dictionary-entries", "/dictionary-entries/{dictKey}"],
    ["activity", "/activity"],
    ["system-monitoring", "/system-monitoring"],
    ["scheduled-tasks", "/scheduled-tasks"],
    ["task-runs", "/task-runs"],
    ["recycle-bin", "/recycle-bin"],
    ["mail", "/mail"],
    ["mail-outbox", "/mail-outbox"],
    ["telegram-settings", "/telegram-settings"],
    ["telegram-operator", "/telegram-settings/operator"],
    ["wallet", "/wallet"],
    ["wallet-entries", "/wallet-entries/{id}"],
    ["wallet-vouchers", "/wallet-vouchers"],
    ["digitaloffer-offers", "/digitaloffer-offers"],
    ["digitaloffer-entitlements", "/digitaloffer-entitlements"],
    ["digitaloffer-purchases", "/digitaloffer-purchases"],
    ["account", "/account"],
    ["settings", "/settings"],
    ["my-wallet", "/my-wallet"],
  ].map(([pageId, route]) => ({
    pageId,
    title: pageId,
    schemaUrl: `/schema/${pageId}`,
    route,
  }));
  const link = (pageRef: string) => ({ pageRef, label: pageRef });
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "r4-route-matrix", name: "R4 route matrix", homePageRef: "dashboard" },
    pages,
    navigation: {
      sidebar: [
        {
          label: "Workspace · WORKSPACE",
          labelKey: GROUP_KEYS.workspace,
          items: [{ ...link("dashboard"), label: "dashboard · 01" }],
        },
        {
          label: "Identity & access",
          labelKey: GROUP_KEYS.identity,
          items: [link("users"), link("roles"), link("data-permission")],
        },
        {
          label: "Content & data",
          labelKey: GROUP_KEYS.content,
          items: [link("file-library"), link("data-dictionary")],
        },
        {
          label: "Operations",
          labelKey: GROUP_KEYS.operations,
          items: [link("activity"), link("system-monitoring"), link("scheduled-tasks"), link("recycle-bin")],
        },
        {
          label: "Communications",
          labelKey: GROUP_KEYS.communications,
          items: [link("mail"), link("mail-outbox"), link("telegram-settings")],
        },
        {
          label: "Commerce",
          labelKey: GROUP_KEYS.commerce,
          items: [
            link("wallet"),
            link("wallet-vouchers"),
            link("digitaloffer-offers"),
            link("digitaloffer-entitlements"),
            link("digitaloffer-purchases"),
          ],
        },
      ],
      user: [link("account"), link("my-wallet"), link("settings")],
    },
  });
}

describe("R4 grouped navigation route matrix", () => {
  it("activates every grouped NodeID and its parent group for direct/deep URLs", () => {
    const manifest = routeMatrixManifest();
    for (const testCase of routeCases) {
      const projection = projectNavigation(manifest, testCase.path, { features: {} });
      const group = projection.sidebar.find(
        (item) => item.type === "group" && item.key === testCase.groupKey,
      );
      expect(group, `${testCase.path} group`).toMatchObject({ type: "group", active: true });
      if (group?.type === "group") {
        expect(group.items.find((item) => item.pageRef === testCase.pageRef), testCase.path).toMatchObject({
          active: true,
        });
      }
    }
  });

  it("places Dashboard in Workspace and keeps user-slot links outside groups", () => {
    const manifest = routeMatrixManifest();
    const dashboard = projectNavigation(manifest, "/dashboard", { features: {} });
    expect(dashboard.sidebar[0]).toMatchObject({
      type: "group",
      key: GROUP_KEYS.workspace,
      secondary: "WORKSPACE",
      active: true,
      items: [{ pageRef: "dashboard", secondary: "01", active: true }],
    });
    expect(dashboard.sidebar.filter((item) => item.type === "group")).toHaveLength(6);

    const settings = projectNavigation(manifest, "/settings", { features: {} });
    expect(settings.user.find((item) => item.type === "link" && item.pageRef === "settings")).toMatchObject({
      active: true,
    });
    expect(settings.sidebar.some((item) => item.type === "link" && item.pageRef === "settings")).toBe(false);
  });
});
