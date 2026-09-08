// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { defaultBranding } from "@/app/branding";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest } from "@/protocol/app-manifest";

const GROUP_KEY = "manifest.nav.group.identityAccess";
const GROUP_STORAGE_KEY = "schema-ui:nav-groups:v1";

function manifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "nav-groups", name: "Navigation groups", homePageRef: "dashboard" },
    pages: [
      { pageId: "dashboard", title: "Dashboard", schemaUrl: "/schema/dashboard", route: "/dashboard" },
      { pageId: "users", title: "Users", schemaUrl: "/schema/users", route: "/users" },
      {
        pageId: "users-invites",
        title: "Invitations",
        schemaUrl: "/schema/users-invites",
        route: "/users-invites",
      },
    ],
    navigation: {
      sidebar: [
        {
          label: "Workspace · WORKSPACE",
          labelKey: "manifest.nav.group.workspace",
          items: [{ pageRef: "dashboard", label: "Dashboard · 01" }],
        },
        {
          label: "Identity & access · IAM",
          labelKey: GROUP_KEY,
          items: [{ pageRef: "users", label: "Users · SQL" }],
        },
      ],
    },
  });
}

function schemaDocument(pageId: string, title: string) {
  return {
    meta: {
      pageId,
      title,
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
    },
    body: { type: "section", children: [{ type: "text", props: { text: `${title} body` } }] },
  };
}

function schemaFetcher(input: RequestInfo | URL): Promise<Response> {
  const pathname = new URL(String(input), "http://test.local").pathname;
  const documents: Record<string, unknown> = {
    "/schema/dashboard": schemaDocument("dashboard", "Dashboard"),
    "/schema/users": schemaDocument("users", "Users"),
    "/schema/users-invites": schemaDocument("users-invites", "Invitations"),
  };
  const document = documents[pathname];
  return Promise.resolve(
    document === undefined
      ? new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), { status: 404 })
      : new Response(JSON.stringify(document), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        }),
  );
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  window.sessionStorage.clear();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
  window.sessionStorage.clear();
});

async function renderApp(path: string): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <App
          manifest={manifest()}
          schemaFetcher={schemaFetcher as typeof fetch}
          branding={defaultBranding()}
        />
      </I18nProvider>,
    );
  });
  return container;
}

describe("R3 collapsible navigation groups", () => {
  it("toggles with Enter/Space and persists the session preference", async () => {
    const container = await renderApp("/users");
    const groupButton = container.querySelector<HTMLButtonElement>(
      `section[data-navigation-group="${GROUP_KEY}"] button`,
    );
    expect(groupButton).not.toBeNull();
    expect(groupButton?.className).toContain("hover:bg-accent/60");
    expect(groupButton?.closest("section")?.className).toContain("space-y-1");
    expect(container.querySelector('[id^="navigation-group-"]')?.className).toContain("border-l");
    expect(groupButton?.textContent).toContain("IAM");
    expect(container.querySelector('a[href="/users"]')?.textContent).toContain("SQL");
    expect(container.querySelector('[data-navigation-secondary="SQL"]')).not.toBeNull();
    const activeMarker = container.querySelector<HTMLElement>('[data-navigation-active-marker="active"]');
    expect(activeMarker).not.toBeNull();
    expect(activeMarker?.className).toContain("bg-primary");
    expect(container.querySelector(".animate-ping")).toBeNull();
    expect(container.querySelector('aside[data-shell-region="sidenav"] > div')?.className).toContain("space-y-2");
    expect(groupButton?.getAttribute("aria-expanded")).toBe("true");
    expect(container.querySelector('a[href="/users"]')).not.toBeNull();

    await act(async () => {
      groupButton?.dispatchEvent(new KeyboardEvent("keydown", { key: " ", bubbles: true }));
    });
    expect(groupButton?.getAttribute("aria-expanded")).toBe("false");
    expect(container.querySelector('a[href="/users"]')).toBeNull();
    expect(JSON.parse(window.sessionStorage.getItem(GROUP_STORAGE_KEY) ?? "{}")[GROUP_KEY]).toBe(false);

    await act(async () => {
      groupButton?.dispatchEvent(new KeyboardEvent("keydown", { key: "Enter", bubbles: true }));
    });
    expect(groupButton?.getAttribute("aria-expanded")).toBe("true");
    expect(container.querySelector('a[href="/users"]')).not.toBeNull();
  });

  it("defaults an inactive group to closed", async () => {
    const container = await renderApp("/dashboard");
    const groupButton = container.querySelector<HTMLButtonElement>(
      `section[data-navigation-group="${GROUP_KEY}"] button`,
    );
    expect(groupButton?.getAttribute("aria-expanded")).toBe("false");
    expect(container.querySelector('a[href="/users"]')).toBeNull();
  });

  it("auto-expands a collapsed group for an inner-page deep link", async () => {
    window.sessionStorage.setItem(GROUP_STORAGE_KEY, JSON.stringify({ [GROUP_KEY]: false }));
    const container = await renderApp("/users-invites");
    const groupButton = container.querySelector<HTMLButtonElement>(
      `section[data-navigation-group="${GROUP_KEY}"] button`,
    );
    expect(groupButton?.getAttribute("aria-expanded")).toBe("true");
    expect(container.querySelector('a[href="/users"]')?.getAttribute("aria-current")).toBe("page");
  });

  it("falls back to the default closed state when session storage is malformed", async () => {
    window.sessionStorage.setItem(GROUP_STORAGE_KEY, "not-json");
    const container = await renderApp("/dashboard");
    const groupButton = container.querySelector<HTMLButtonElement>(
      `section[data-navigation-group="${GROUP_KEY}"] button`,
    );
    expect(groupButton?.getAttribute("aria-expanded")).toBe("false");
  });
});
