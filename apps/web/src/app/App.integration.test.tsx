// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
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
): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <App manifest={testManifest()} navigationContext={navigationContext} />,
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

    const catalogLink = container.querySelector('a[href="/catalog"]');
    expect(catalogLink).not.toBeNull();
    await act(async () => {
      catalogLink?.dispatchEvent(new MouseEvent("click", { bubbles: true, button: 0 }));
    });
    expect(window.location.pathname).toBe("/catalog");
    expect(container.querySelector("h1")?.textContent).toBe("Catalog");

    window.history.pushState({}, "", "/home");
    await act(async () => {
      window.dispatchEvent(new PopStateEvent("popstate"));
    });
    expect(container.querySelector("h1")?.textContent).toBe("Home");
  });

  it("renders a fail-closed fallback and returns to the manifest home route", async () => {
    const container = await renderApp("/unknown");
    expect(container.textContent).toContain("Page not found");
    expect(container.textContent).toContain("/unknown");

    const homeButton = Array.from(container.querySelectorAll("button")).find((button) =>
      button.textContent?.includes("Return to home"),
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
        <App
          manifest={testManifest()}
          navigationContext={{}}
          accountError={new Error("account unavailable")}
        />,
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
      healthyRoot.render(<App manifest={testManifest()} navigationContext={{}} />);
    });
    expect(healthy.textContent).not.toContain("Account session failed to load");
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
