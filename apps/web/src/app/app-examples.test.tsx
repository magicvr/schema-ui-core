// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { validateAppManifest } from "@/protocol/app-manifest";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

function exampleManifest() {
  return validateAppManifest({
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
    app: { appId: "examples", name: "Examples", homePageRef: "overview" },
    pages: [
      { pageId: "overview", title: "Overview", schemaUrl: "/schema/overview", route: "/overview" },
      { pageId: "data-table", title: "Data table", schemaUrl: "/schema/data-table", route: "/data-table" },
      {
        pageId: "search-form-table",
        title: "Search + table",
        schemaUrl: "/schema/search-form-table",
        route: "/search-form-table",
      },
    ],
    navigation: {
      sidebar: [
        { pageRef: "data-table", label: "Data table" },
        { pageRef: "search-form-table", label: "Search + table" },
      ],
    },
  });
}

async function renderApp(path: string): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<App manifest={exampleManifest()} navigationContext={{}} />);
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

describe("R5 example pages in the shell", () => {
  it("renders the data-table example surface on its route", async () => {
    const container = await renderApp("/data-table");
    expect(container.querySelector("h1")?.textContent).toContain("Data table");
  });

  it("renders the search-form-table example surface on its route", async () => {
    const container = await renderApp("/search-form-table");
    expect(container.querySelector("h1")?.textContent).toContain("Search + table");
  });

  it("keeps the manifest fallback for non-example pages", async () => {
    const container = await renderApp("/overview");
    expect(container.querySelector("h1")?.textContent).toContain("Overview");
  });
});
