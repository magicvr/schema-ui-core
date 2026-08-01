// @vitest-environment jsdom
//
// Full-path integration evidence for the R1 representative pages (GOAL-004):
// uses the real app manifest (`apps/web/public/.well-known/schema-ui/`) and the
// real Go-embedded page fixtures (`apps/api/internal/handler/fixtures/schema/`)
// through the App's schema-driven default path, with the records API surface
// injected. Asserts "改 Schema 即可出现页面" holds on the main path and that
// unknown / illegal inputs fail closed with observable errors.

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const MANIFEST_PATH = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../public/.well-known/schema-ui/app-manifest.json",
);
const FIXTURE_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/internal/handler/fixtures/schema",
);

const MIGRATED_PAGE_IDS = [
  "overview",
  "catalog",
  "data-table",
  "search-form-table",
  "list-edit-lifecycle",
  "form-controls",
  "form-with-reactions",
];

const RECORDS = {
  items: [
    {
      id: "rec-1",
      name: "Acme Console",
      status: "active",
      owner: "alice",
      updatedAt: "2026-07-31T00:00:00Z",
    },
    {
      id: "rec-2",
      name: "Northwind Sales",
      status: "pending",
      owner: "bob",
      updatedAt: "2026-07-31T11:00:00Z",
    },
  ],
  total: 2,
  page: 1,
  pageSize: 10,
};

function fixtureDocument(pageId: string): unknown {
  return JSON.parse(readFileSync(resolve(FIXTURE_DIR, `${pageId}.json`), "utf8"));
}

function manifest(): AppManifest {
  return validateAppManifest(
    JSON.parse(readFileSync(MANIFEST_PATH, "utf8")),
  );
}

/** Serves /api/schema/* from the real fixtures and /api/records from the demo envelope. */
function combinedFetcher(
  fixtures: Record<string, unknown>,
  recordsStatus = 200,
): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    if (pathname.startsWith("/api/schema/")) {
      const pageId = pathname.slice("/api/schema/".length);
      const document = fixtures[pageId];
      if (document === undefined) {
        return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), { status: 404 });
      }
      return new Response(JSON.stringify(document), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/records")) {
      if (recordsStatus !== 200) {
        return new Response(JSON.stringify({ error: "records down" }), { status: recordsStatus });
      }
      return new Response(JSON.stringify(RECORDS), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    return new Response("not found", { status: 404 });
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
  context: Record<string, unknown>,
  fixtures: Record<string, unknown>,
  recordsStatus = 200,
): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const fetcher = combinedFetcher(fixtures, recordsStatus);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <App
        manifest={manifest()}
        navigationContext={context}
        schemaFetcher={fetcher}
        recordsFetcher={fetcher}
      />,
    );
  });
  return container;
}

function realFixtures(): Record<string, unknown> {
  const fixtures: Record<string, unknown> = {};
  for (const pageId of MIGRATED_PAGE_IDS) {
    fixtures[pageId] = fixtureDocument(pageId);
  }
  return fixtures;
}

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
});

describe("representative pages through the real manifest (GOAL-004)", () => {
  it("renders a migrated list page with records via the default path", async () => {
    const container = await renderApp("/data-table", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Data table");
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).toContain("Northwind Sales");
    expect(container.textContent).toContain("2 records · page 1 of 1");
  });

  it("renders the search + table page structure", async () => {
    const container = await renderApp("/search-form-table", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Search + table");
    expect(container.textContent).toContain("Search records");
    expect(container.textContent).toContain("Acme Console");
  });

  it("renders a migrated form page", async () => {
    const container = await renderApp("/form-controls", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Form controls");
    expect(container.textContent).toContain("Name (input)");
    expect(container.textContent).toContain("Kind (select)");
  });

  it("renders the composite/detail page (tabs + recordView + form)", async () => {
    const container = await renderApp("/list-edit-lifecycle", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("List + edit lifecycle");
    expect(container.textContent).toContain("Detail");
    expect(container.textContent).toContain("Edit");
    expect(container.textContent).toContain("Acme Console");
    // The Detail tab is active by default; switch to Edit to reveal the form.
    const editTab = Array.from(container.querySelectorAll('[role="tab"]')).find(
      (tab) => tab.textContent === "Edit",
    );
    expect(editTab).not.toBeUndefined();
    await act(async () => (editTab as HTMLElement).click());
    expect(container.textContent).toContain("Status");
    expect(container.textContent).toContain("Name");
  });

  it("applies $context reactions on the migrated reactive form", async () => {
    const admin = await renderApp("/form-with-reactions", { user: { roles: ["admin"] } }, realFixtures());
    expect(admin.textContent).toContain("Approval (switch)");

    const viewer = await renderApp("/form-with-reactions", { user: { roles: ["viewer"] } }, realFixtures());
    expect(viewer.textContent).not.toContain("Approval (switch)");
  });

  it("fails closed when a representative page document has an unknown node type", async () => {
    const fixtures = realFixtures();
    fixtures["search-form-table"] = {
      meta: {
        pageId: "search-form-table",
        title: "Search + table",
        protocolVersion: "2.7",
        requiredCapabilities: ["app.manifest"],
      },
      body: { type: "chart", id: "x", props: {} },
    };
    const container = await renderApp("/search-form-table", {}, fixtures);
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "outside the §5 renderer whitelist",
    );
  });

  it("fails closed when the records data source is unreachable on a list page", async () => {
    const container = await renderApp("/data-table", {}, realFixtures(), 500);
    expect(container.textContent).toContain("records fetch failed");
  });
});
