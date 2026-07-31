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
  // jsdom cannot resolve relative URLs through fetch; stub the records API
  // used by the R5 example pages (one row so the Edit/Delete gates render).
  globalThis.fetch = (async (input: RequestInfo | URL, _init?: RequestInit) => {
    const url = String(input);
    if (url.startsWith("/api/records/")) {
      // PATCH / DELETE detail routes are not exercised by these surface tests.
      return new Response(null, { status: 204 });
    }
    return new Response(
      JSON.stringify({
        items: [
          {
            id: "rec-1",
            name: "Acme Console",
            status: "active",
            owner: "alice",
            updatedAt: "2026-07-31T00:00:00Z",
          },
        ],
        total: 1,
        page: 1,
        pageSize: 10,
      }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  }) as typeof fetch;
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
      {
        pageId: "list-edit-lifecycle",
        title: "List + edit lifecycle",
        schemaUrl: "/schema/list-edit-lifecycle",
        route: "/list-edit-lifecycle",
      },
      {
        pageId: "form-controls",
        title: "Form controls",
        schemaUrl: "/schema/form-controls",
        route: "/form-controls",
      },
    ],
    navigation: {
      sidebar: [
        { pageRef: "data-table", label: "Data table" },
        { pageRef: "search-form-table", label: "Search + table" },
        { pageRef: "list-edit-lifecycle", label: "List + edit" },
        { pageRef: "form-controls", label: "Form controls" },
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

  it("renders the list-edit-lifecycle example surface with Edit/Delete gates", async () => {
    const container = await renderApp("/list-edit-lifecycle");
    expect(container.querySelector("h1")?.textContent).toContain("List + edit lifecycle");
    // Row action buttons are present (permission gate passes for the dev admin).
    expect(container.textContent).toContain("Edit");
    expect(container.textContent).toContain("Delete");
  });

  it("renders the form-controls example surface with the whitelist gate", async () => {
    const container = await renderApp("/form-controls");
    expect(container.querySelector("h1")?.textContent).toContain("Form controls");
    // The capability gate passes for the 2.7 + extended/advanced meta.
    expect(container.textContent).toContain("Capability gate");
  });

  it("keeps the manifest fallback for non-example pages", async () => {
    const container = await renderApp("/overview");
    expect(container.querySelector("h1")?.textContent).toContain("Overview");
  });
});
