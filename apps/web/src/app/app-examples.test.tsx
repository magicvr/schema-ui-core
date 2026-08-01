// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { validateAppManifest } from "@/protocol/app-manifest";

// D-003 (GOAL-003): the 5 hand-written EXAMPLE_PAGES are migrated to Schema
// documents that render through the default path
// (`page.route → page.schemaUrl → loadPageDocument → RenderPage`). These tests
// inject schema fixtures via App's schemaFetcher and assert:
//   (a) example routes render Schema content, not the hand-written surfaces;
//   (b) the hand-written EXAMPLE_PAGES are no longer part of the render path;
//   (c) a missing / invalid schema fails closed with the unified error surface.

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
      {
        pageId: "form-with-reactions",
        title: "Form with reactions",
        schemaUrl: "/schema/form-with-reactions",
        route: "/form-with-reactions",
      },
    ],
    navigation: {
      sidebar: [
        { pageRef: "data-table", label: "Data table" },
        { pageRef: "search-form-table", label: "Search + table" },
        { pageRef: "list-edit-lifecycle", label: "List + edit" },
        { pageRef: "form-controls", label: "Form controls" },
        { pageRef: "form-with-reactions", label: "Form with reactions" },
      ],
    },
  });
}

// Structurally valid page document (meta + section + text), mirroring the
// shape the pinned page/node schemas accept.
function schemaDocument(pageId: string, title: string, text: string) {
  return {
    meta: {
      pageId,
      title,
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest", "app.navigation"],
    },
    body: {
      type: "section",
      children: [{ type: "text", props: { text } }],
    },
  };
}

// Fetcher keyed by resolved schemaUrl pathname; a missing key 404s so the
// loader maps it to PAGE_NOT_FOUND (fail-closed), exactly like the API does.
function schemaFetcher(documents: Record<string, unknown>) {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    const document = documents[pathname];
    if (document === undefined) {
      return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), {
        status: 404,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response(JSON.stringify(document), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;
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
  documents: Record<string, unknown>,
): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <App
        manifest={exampleManifest()}
        navigationContext={{}}
        schemaFetcher={schemaFetcher(documents)}
      />,
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

describe("schema-driven default path (GOAL-003)", () => {
  it("renders a migrated example route from its Schema document", async () => {
    const container = await renderApp("/data-table", {
      "/schema/data-table": schemaDocument("data-table", "Data table", "Schema-driven records"),
    });
    expect(container.querySelector("h1")?.textContent).toContain("Data table");
    expect(container.textContent).toContain("Schema-driven records");
  });

  it("does not render the hand-written example surface on its route", async () => {
    const container = await renderApp("/form-with-reactions", {
      "/schema/form-with-reactions": schemaDocument(
        "form-with-reactions",
        "Form with reactions",
        "Schema-driven reactive form",
      ),
    });
    // The hand-written surface markers (Context snapshot / Name (input)) are gone.
    expect(container.textContent).toContain("Schema-driven reactive form");
    expect(container.textContent).not.toContain("Context snapshot");
    expect(container.textContent).not.toContain("evaluateExpression");
  });

  it("fails closed with the unified error when the schema document is missing", async () => {
    const container = await renderApp("/search-form-table", {});
    expect(container.textContent).toContain("PAGE_NOT_FOUND");
    expect(container.textContent).toContain("/schema/search-form-table");
  });

  it("fails closed with the unified error for an invalid schema document", async () => {
    const container = await renderApp("/form-controls", {
      "/schema/form-controls": {
        meta: { pageId: "form-controls", title: "Form controls" },
        body: { type: "section" },
      },
    });
    expect(container.textContent).toContain("PAGE_SCHEMA_INVALID");
  });
});
