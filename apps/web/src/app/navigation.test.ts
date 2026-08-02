import { readFileSync } from "node:fs";
import { describe, expect, it } from "vitest";

import { projectNavigation } from "@/app/navigation";
import {
  type AppManifest,
  loadAppManifest,
  validateAppManifest,
} from "@/protocol/app-manifest";

function testManifest(): AppManifest {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: {
      appId: "admin",
      name: "Admin",
      homePageRef: "home",
    },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/s/home", route: "/home" },
      { pageId: "orders", title: "Orders", schemaUrl: "/s/orders", route: "/orders" },
      {
        pageId: "orders-detail",
        title: "Order detail",
        schemaUrl: "/s/orders/{id}",
        route: "/orders/{id}",
      },
    ],
    navigation: {
      top: [{ pageRef: "home", label: "Home", icon: "home" }],
      sidebar: [
        { pageRef: "orders", label: "Orders", icon: "orders" },
        {
          label: "Admin",
          items: [
            {
              pageRef: "orders-detail",
              label: "Details",
              permissions: { view: '$context.user.roles contains "admin"' },
            },
          ],
        },
      ],
      user: [
        { url: "/profile", label: "Profile" },
        { url: "/profile/settings", label: "Settings" },
      ],
    },
  });
}

describe("navigation projection", () => {
  it("projects top/sidebar/user slots in declaration order", () => {
    const result = projectNavigation(testManifest(), "/orders", {
      user: { roles: ["admin"] },
      features: {},
    });
    expect(result.top.map((item) => item.label)).toEqual(["Home"]);
    expect(result.sidebar.map((item) => item.label)).toEqual(["Orders", "Admin"]);
    expect(result.user.map((item) => item.label)).toEqual(["Profile", "Settings"]);
    expect(result.sidebar[0]).toMatchObject({ type: "link", active: true });
    expect(result.sidebar[1]).toMatchObject({
      type: "group",
      items: [{ label: "Details", active: false }],
    });
    if (result.sidebar[1]?.type === "group") {
      expect(result.sidebar[1].items[0]).not.toHaveProperty("href");
    }
  });

  it("filters permission-gated links and prunes empty groups", () => {
    const result = projectNavigation(testManifest(), "/orders", {
      user: { roles: ["viewer"] },
      features: {},
    });
    expect(result.sidebar.map((item) => item.label)).toEqual(["Orders"]);
  });

  it("uses D4a for pageRef active state and exact URL active state", () => {
    const result = projectNavigation(testManifest(), "/orders/42?tab=events", {
      user: { roles: ["admin"] },
      features: {},
    });
    expect(result.sidebar[0]).toMatchObject({ active: false });
    expect(result.sidebar[1]).toMatchObject({
      type: "group",
      items: [{ label: "Details", active: true, href: "/orders/42" }],
    });

    const urlResult = projectNavigation(testManifest(), "/profile?tab=security", {
      user: {},
      features: {},
    });
    expect(urlResult.user[0]).toMatchObject({ active: true });
    expect(urlResult.user[1]).toMatchObject({ active: false });
  });
});

// GOAL-006 S5 · the real checked-in manifest gates list-edit-lifecycle on
// $context.features.menu_list_edit_lifecycle (V-MENU-03/04/05/06).
const checkedInManifestBytes = readFileSync(
  new URL("../../public/.well-known/schema-ui/app-manifest.json", import.meta.url),
);

describe("GOAL-006 S5 · list-edit-lifecycle menu projection", () => {
  async function realManifest(): Promise<AppManifest> {
    return loadAppManifest({
      fetcher: async () => new Response(checkedInManifestBytes, { status: 200 }),
    });
  }

  async function examplesLabels(features: Record<string, boolean>): Promise<string[]> {
    const manifest = await realManifest();
    const sidebar = projectNavigation(manifest, "/", { user: { roles: [] }, features }).sidebar;
    const group = sidebar.find((item) => item.type === "group" && item.label === "Examples");
    return group?.type === "group" ? group.items.map((child) => child.label) : [];
  }

  it("shows List + edit for admin and hides it for viewer, keeping the group (V-MENU-04/06)", async () => {
    expect(await examplesLabels({ menu_list_edit_lifecycle: true })).toEqual([
      "Data table",
      "Search + table",
      "List + edit",
      "Form controls",
      "Form with reactions",
    ]);
    // viewer: the child is hidden, declaration order and other children remain
    // (the Examples group is not pruned).
    expect(await examplesLabels({ menu_list_edit_lifecycle: false })).toEqual([
      "Data table",
      "Search + table",
      "Form controls",
      "Form with reactions",
    ]);
  });

  it("fails closed when the feature is missing or falsy (V-MENU-05)", async () => {
    expect(await examplesLabels({})).not.toContain("List + edit");
    expect(await examplesLabels({ menu_list_edit_lifecycle: false })).not.toContain("List + edit");
    // wrong type also denies: the renderer never surfaces the child.
    expect(await examplesLabels({ menu_list_edit_lifecycle: "yes" as unknown as boolean })).not.toContain(
      "List + edit",
    );
  });
});
