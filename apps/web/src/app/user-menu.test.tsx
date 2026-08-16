// @vitest-environment jsdom
//
// T-01 (GOAL-013 D-002): the topbar user dropdown folds 个人中心 / 设置 /
// (我的钱包 when present) / 退出登录 into one trigger for every breakpoint,
// and the mobile drawer no longer repeats the user chain.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

function userMenuManifest(): AppManifest {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "um", name: "UserMenu", homePageRef: "home" },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/" },
      { pageId: "account", title: "Account", schemaUrl: "/schema/account", route: "/account" },
      { pageId: "settings", title: "Settings", schemaUrl: "/schema/settings", route: "/settings" },
    ],
    navigation: {
      user: [
        { pageRef: "account", label: "Account", icon: "user" },
        { pageRef: "settings", label: "Settings", icon: "settings" },
      ],
    },
  });
}

function schemaDocument(pageId: string) {
  return {
    meta: { pageId, title: pageId, protocolVersion: "2.7", requiredCapabilities: ["app.manifest", "app.navigation"] },
    body: { type: "section", children: [{ type: "text", props: { text: pageId } }] },
  };
}

function schemaFetcher() {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    const pageId = pathname.split("/").pop() ?? "home";
    return new Response(JSON.stringify(schemaDocument(pageId)), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
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
          manifest={userMenuManifest()}
          schemaFetcher={schemaFetcher()}
          currentUser={{ id: "usr-admin", name: "Alice" }}
          onLogout={vi.fn()}
        />
      </I18nProvider>,
    );
  });
  return container;
}

describe("T-01 user dropdown (GOAL-013 D-002)", () => {
  it("renders one avatar+name trigger and folds user nav + signout into the menu", async () => {
    const container = await renderApp("/");
    const trigger = container.querySelector('button[aria-label="User menu"]');
    expect(trigger).not.toBeNull();
    // Trigger shows the display name.
    expect(container.textContent).toContain("Alice");
    // Desktop horizontal user nav is gone: no separate "User navigation" nav.
    expect(container.querySelector('nav[aria-label="User navigation"]')).toBeNull();
    // The menu is closed initially.
    expect(container.querySelector('[role="menu"]')).toBeNull();

    await act(async () => {
      (trigger as HTMLButtonElement).click();
    });
    const menu = container.querySelector('[role="menu"]');
    expect(menu).not.toBeNull();
    const items = Array.from(menu!.querySelectorAll('[role="menuitem"]')).map((el) =>
      el.textContent?.trim(),
    );
    // projection.user declaration order + divider + signout.
    expect(items).toEqual(["Account", "Settings", "Sign out"]);
  });

  it("navigates from a menu item and closes the menu", async () => {
    const container = await renderApp("/");
    await act(async () => {
      (container.querySelector('button[aria-label="User menu"]') as HTMLButtonElement).click();
    });
    const accountItem = Array.from(container.querySelectorAll('[role="menuitem"]')).find(
      (el) => el.textContent?.trim() === "Account",
    );
    await act(async () => {
      (accountItem as HTMLButtonElement).click();
    });
    expect(window.location.pathname).toBe("/account");
    // Menu closed after navigation.
    expect(container.querySelector('[role="menu"]')).toBeNull();
  });

  it("keeps the mobile drawer free of the user chain (D-002 §3)", async () => {
    const container = await renderApp("/");
    await act(async () => {
      (container.querySelector('button[aria-label="Open navigation menu"]') as HTMLButtonElement).click();
    });
    const drawer = container.querySelector('nav[aria-label="Mobile navigation"]');
    expect(drawer).not.toBeNull();
    expect(drawer!.textContent).not.toContain("Account");
    expect(drawer!.textContent).not.toContain("Settings");
  });

  it("closes on Escape", async () => {
    const container = await renderApp("/");
    await act(async () => {
      (container.querySelector('button[aria-label="User menu"]') as HTMLButtonElement).click();
    });
    expect(container.querySelector('[role="menu"]')).not.toBeNull();
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape" }));
    });
    expect(container.querySelector('[role="menu"]')).toBeNull();
  });
});
