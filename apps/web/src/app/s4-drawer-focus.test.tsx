// @vitest-environment jsdom
//
// F-002 跨模块 UI 可访问性下限 · 移动抽屉焦点可复跑断言（S0 D-003 §8）：
// 抽屉打开焦点进入、Escape 关闭、关闭后焦点恢复到触发元素。补齐 S1 close-path
// 中「移动抽屉对应可复跑断言」——此前仅有实现（App.tsx）与结构检查（shell.test.ts）。
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest, type NavigationContext } from "@/protocol/app-manifest";

function drawerManifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "integration", name: "Integration", homePageRef: "home" },
    pages: [
      { pageId: "home", title: "Home", schemaUrl: "/schema/home", route: "/home" },
      { pageId: "catalog", title: "Catalog", schemaUrl: "/schema/catalog", route: "/catalog" },
    ],
    navigation: {
      top: [{ pageRef: "home", label: "Home" }],
      sidebar: [{ pageRef: "catalog", label: "Catalog" }],
    },
  });
}

function drawerSchemaFetcher() {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    const documents: Record<string, unknown> = {
      "/schema/home": {
        meta: { pageId: "home", title: "Home", protocolVersion: "2.7", requiredCapabilities: [] },
        body: { type: "section", children: [{ type: "text", props: { text: "Home body" } }] },
      },
      "/schema/catalog": {
        meta: { pageId: "catalog", title: "Catalog", protocolVersion: "2.7", requiredCapabilities: [] },
        body: { type: "section", children: [{ type: "text", props: { text: "Catalog body" } }] },
      },
    };
    if (documents[pathname]) {
      return new Response(JSON.stringify(documents[pathname]), { status: 200 });
    }
    return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), { status: 404 });
  }) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
});

async function renderDrawer(path: string, navigationContext?: NavigationContext) {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <App
          manifest={drawerManifest()}
          navigationContext={navigationContext}
          schemaFetcher={drawerSchemaFetcher()}
        />
      </I18nProvider>,
    );
  });
  return container;
}

function hamburgerButton(container: HTMLElement): HTMLButtonElement {
  const button = Array.from(container.querySelectorAll("button")).find(
    (b) => b.getAttribute("aria-label")?.includes("Open") || b.getAttribute("aria-expanded") !== null,
  );
  if (!button) {
    throw new Error("hamburger button not found");
  }
  return button;
}

describe("S4 · mobile drawer focus management (F-002)", () => {
  it("moves focus into the drawer on open and restores it on Escape-close", async () => {
    const container = await renderDrawer("/home", { user: { roles: ["admin"] }, features: {} });
    const hamburger = hamburgerButton(container);
    hamburger.focus();

    // Open the drawer.
    await act(async () => {
      hamburger.click();
    });
    const drawer = container.querySelector('nav[aria-label="Mobile navigation"]') as HTMLElement | null;
    expect(drawer).not.toBeNull();
    // Focus enters the drawer (first focusable = close button).
    expect(drawer?.contains(document.activeElement)).toBe(true);

    // Escape closes the drawer.
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape", bubbles: true }));
    });
    expect(container.querySelector('nav[aria-label="Mobile navigation"]')).toBeNull();

    // Focus restored to the hamburger trigger.
    expect(document.activeElement).toBe(hamburger);
  });

  it("traps Tab focus within the drawer", async () => {
    const container = await renderDrawer("/home", { user: { roles: ["admin"] }, features: {} });
    const hamburger = hamburgerButton(container);
    await act(async () => {
      hamburger.click();
    });
    const drawer = container.querySelector('nav[aria-label="Mobile navigation"]') as HTMLElement;
    const focusables = Array.from(
      drawer.querySelectorAll<HTMLElement>(
        'a[href], button:not([disabled]), [tabindex]:not([tabindex="-1"])',
      ),
    );
    expect(focusables.length).toBeGreaterThanOrEqual(1);
    const first = focusables[0];
    const last = focusables[focusables.length - 1];

    last.focus();
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Tab", bubbles: true }));
    });
    expect(document.activeElement).toBe(first);

    first.focus();
    await act(async () => {
      document.dispatchEvent(
        new KeyboardEvent("keydown", { key: "Tab", shiftKey: true, bubbles: true }),
      );
    });
    expect(document.activeElement).toBe(last);
  });
});
