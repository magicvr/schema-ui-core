// @vitest-environment jsdom
//
// Cross-boundary regression for the R1 representative pages (GOAL-004):
// reads the actual Go-embedded core and module-owned page fixtures, validates them through the
// browser loader (`validatePageDocument` / `loadPageDocument`), and renders
// them via `RenderPage`. This closes the GOAL-002 F-001/F-002 follow-up
// ("actual endpoint fixture through the loader") and proves the migrated
// list/form/composite pages render on the schema-driven main path.

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { PageEntry } from "@/protocol/app-manifest";
import { validatePageDocument } from "@/protocol/conformance/runtime-schema-validate";
import { loadPageDocument } from "@/protocol/load-page";
import type { RenderPageDocument } from "@/renderer/render";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";

const CORE_FIXTURE_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/internal/modules/schemarender/schema",
);
const MODULE_FIXTURE_DIRS: Record<string, string> = {
  settings: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/settings/schema",
  ),
  activity: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/activity/schema",
  ),
  users: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/users/schema",
  ),
  roles: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/roles/schema",
  ),
};

const MIGRATED_PAGE_IDS = [
  "data-table",
  "search-form-table",
  "form-controls",
  "form-with-reactions",
  "users",
  "roles",
  "settings",
  "activity",
];

function fixtureDocument(pageId: string): unknown {
  const directory = MODULE_FIXTURE_DIRS[pageId] ?? CORE_FIXTURE_DIR;
  return JSON.parse(readFileSync(resolve(directory, `${pageId}.json`), "utf8"));
}

/** Fetcher that serves the real fixture document for a schemaUrl pathname. */
function fixtureFetcher(pageId: string): typeof fetch {
  const document = fixtureDocument(pageId);
  return (async () =>
    new Response(JSON.stringify(document), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    })) as typeof fetch;
}

const USERS = {
  items: [
    {
      id: "usr-1",
      username: "alice",
      name: "Alice",
      roles: ["admin"],
      updatedAt: "2026-08-03T00:00:00.000Z",
    },
  ],
  total: 1,
  page: 1,
  pageSize: 10,
};

function resourceFetcher(): typeof fetch {
  return (async () =>
    new Response(JSON.stringify(USERS), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    })) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderDocument(
  pageDoc: RenderPageDocument,
  context: Record<string, unknown>,
): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <RenderPage
        document={pageDoc}
        context={context}
        tableRenderer={(node) => <SchemaTable node={node} fetcher={resourceFetcher()} />}
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
});

describe("migrated representative pages (GOAL-004)", () => {
  it("every migrated fixture passes structural validation and loads through the loader", async () => {
    for (const pageId of MIGRATED_PAGE_IDS) {
      const document = fixtureDocument(pageId);
      const validation = validatePageDocument(document);
      expect(validation.ok, `${pageId} must pass D-VAL`).toBe(true);

      const page = {
        pageId,
        title: pageId,
        schemaUrl: `/api/schema/${pageId}`,
        route: `/${pageId}`,
      } as PageEntry;
      const loaded = await loadPageDocument(page, {}, { fetcher: fixtureFetcher(pageId) });
      const meta = (loaded as { meta?: { pageId?: unknown } }).meta;
      expect(meta?.pageId).toBe(pageId);
    }
  });

  it("renders the list page (data-table) with an injected table surface", async () => {
    const container = await renderDocument(
      fixtureDocument("data-table") as RenderPageDocument,
      {},
    );
    expect(container.textContent).toContain("alice");
    expect(container.textContent).toContain("Name");
  });

  it("renders the search + table page structure (form + table)", async () => {
    const container = await renderDocument(
      fixtureDocument("search-form-table") as RenderPageDocument,
      {},
    );
    expect(container.textContent).toContain("Search roles");
    expect(container.textContent).toContain("Alice");
  });

  it("renders the form-controls page with whitelisted controls", async () => {
    const container = await renderDocument(
      fixtureDocument("form-controls") as RenderPageDocument,
      {},
    );
    expect(container.textContent).toContain("Name (input)");
    expect(container.textContent).toContain("Kind (select)");
    expect(container.textContent).toContain("Audit note (textarea)");
  });

  it("applies $context reactions on the form-with-reactions page", async () => {
    const pageDoc = fixtureDocument("form-with-reactions") as RenderPageDocument;
    const admin = await renderDocument(pageDoc, {
      user: { roles: ["admin"] },
      features: { audit: true },
    });
    expect(admin.textContent).toContain("Approval (switch)");

    const viewer = await renderDocument(pageDoc, {
      user: { roles: ["viewer"] },
      features: { audit: false },
    });
    expect(viewer.textContent).not.toContain("Approval (switch)");
    const auditNote = Array.from(viewer.querySelectorAll("textarea")).find(
      (area) => (area as HTMLTextAreaElement).id === "field-auditNote",
    );
    expect((auditNote as HTMLTextAreaElement | undefined)?.disabled).toBe(true);
  });

  it("renders the CRUD lifecycle page (table + toolbar + row actions + recordView)", async () => {
    const admin = {
      user: { id: "u1", roles: ["admin"], permissions: ["users.read", "users.write"] },
    };
    const container = await renderDocument(
      fixtureDocument("users") as RenderPageDocument,
      admin,
    );
    // S4 structure (I-007-003 §2.1): a table surface with a create toolbar
    // trigger, row edit/delete actions, and a selected-row recordView.
    expect(container.textContent).toContain("New user");
    expect(container.textContent).toContain("Edit");
    expect(container.textContent).toContain("Delete");
    expect(container.textContent).toContain("alice");
    // recordView renders the selected-row copy once a row is selected.
    expect(container.textContent).toContain("Select a record to view details.");
    const editRow = Array.from(container.querySelectorAll("button")).find(
      (button) => button.textContent === "Edit",
    );
    expect(editRow).not.toBeUndefined();
    await act(async () => (editRow as HTMLButtonElement).click());
    expect(container.textContent).toContain("Save changes");
    expect(container.textContent).toContain("alice");
  });

  it("fails closed on an unknown node type on a representative path", async () => {
    const pageDoc = {
      meta: {
        pageId: "search-form-table",
        title: "Search + table",
        protocolVersion: "2.7",
        requiredCapabilities: ["app.manifest"],
      },
      body: { type: "chart", id: "x", props: {} },
    } as unknown as RenderPageDocument;
    const container = await renderDocument(pageDoc, {});
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "outside the §5 renderer whitelist",
    );
  });
});
