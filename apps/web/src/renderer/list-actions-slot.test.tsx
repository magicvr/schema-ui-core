// @vitest-environment jsdom
//
// GOAL-045 W33 (D-001) · the list page-actions slot.
//
// A custom node may declare `props.slot = "list-page-actions"` + `props.targetTable`
// to be rendered inside THAT table's page-actions row, left segment — instead of
// in document flow, which is where a section child lands (the pre-W33 position
// that put the control above the filter panel).
//
// Three properties are pinned here, each guarding a distinct failure mode:
//   1. a slot-declared node really renders inside the target table's left
//      segment (placement is the whole point of the feature);
//   2. a node WITHOUT the declaration keeps the old in-flow behaviour — this is
//      a compatible extension, not a silent relocation of every custom node;
//   3. a slot declaration whose target table does not exist falls back to
//      in-flow rendering rather than disappearing (D-001 §3 fail-open): a
//      layout capability must never turn a schema typo into a missing control.

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import type { RenderPageDocument } from "@/renderer/render.types";
import { registerCustomComponent, resetCustomComponentsForTests } from "@/renderer/custom-components";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";
import "@/components/jobs-batch-export";

const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");
const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
  resetCustomComponentsForTests();
  registerCustomComponent("w33-probe", () => <span data-w33-probe="true">probe</span>);
  registerCustomComponent("w33-other", () => <span data-w33-other="true">other</span>);
});

afterEach(async () => {
  resetCustomComponentsForTests();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

function documentWith(nodes: Array<Record<string, unknown>>, tableId = "w33-table"): RenderPageDocument {
  return {
    meta: {
      pageId: "w33",
      title: "W33",
      protocolVersion: "2.7",
      requiredCapabilities: ["app.manifest"],
    },
    body: { type: "section", id: "w33", children: [...nodes, {
      type: "table",
      id: tableId,
      props: { columns: [{ field: "id", label: "ID" }], dataSource: "/api/w33" },
    }] },
  } as unknown as RenderPageDocument;
}

async function renderDocument(pageDocument: RenderPageDocument): Promise<HTMLDivElement> {
  const fetcher = (async () =>
    new Response(JSON.stringify({ items: [{ id: "1" }], total: 1, page: 1, pageSize: 20 }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    })) as typeof fetch;
  const container = globalThis.document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <RenderPage
          document={pageDocument}
          context={{}}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </I18nProvider>,
    );
  });
  await act(async () => {
    await Promise.resolve();
  });
  return container;
}

describe("W33 · list page-actions slot", () => {
  it("renders a slot-declared node inside the target table's left segment", async () => {
    const container = await renderDocument(
      documentWith([
        {
          type: "custom",
          id: "slotted",
          component: "w33-probe",
          props: { targetTable: "w33-table", slot: "list-page-actions" },
        },
      ]),
    );
    const left = container.querySelector("[data-list-page-actions-left]");
    expect(left, "the target table must expose the left segment").not.toBeNull();
    expect(left!.querySelector("[data-w33-probe]"), "the node must land inside the left segment").not.toBeNull();
    // …and NOT in document flow at the section level.
    const row = container.querySelector("[data-list-page-actions]");
    expect(row, "the page-actions row must exist").not.toBeNull();
    expect(row!.firstElementChild).toBe(left);
  });

  it("leaves a node without the declaration in document flow", async () => {
    const container = await renderDocument(
      documentWith([{ type: "custom", id: "in-flow", component: "w33-other", props: {} }]),
    );
    expect(container.querySelector("[data-w33-other]")).not.toBeNull();
    // No slot consumer => no left segment (a table nobody slots into gains no
    // empty segment).
    expect(container.querySelector("[data-list-page-actions-left]")).toBeNull();
  });

  it("falls back to document flow when the target table does not exist", async () => {
    const container = await renderDocument(
      documentWith([
        {
          type: "custom",
          id: "orphan",
          component: "w33-probe",
          props: { targetTable: "no-such-table", slot: "list-page-actions" },
        },
      ]),
    );
    expect(container.querySelector("[data-w33-probe]"), "fail-open: the control must stay visible").not.toBeNull();
    expect(container.querySelector("[data-list-page-actions-left]")).toBeNull();
  });

  it("treats an unknown slot value as no slot at all", async () => {
    const container = await renderDocument(
      documentWith([
        {
          type: "custom",
          id: "typo",
          component: "w33-probe",
          props: { targetTable: "w33-table", slot: "list-page-action" },
        },
      ]),
    );
    expect(container.querySelector("[data-w33-probe]")).not.toBeNull();
    expect(container.querySelector("[data-list-page-actions-left]")).toBeNull();
  });

  it("keeps two different tables' column state apart across a page switch", async () => {
    // Regression lock for the bug this feature uncovered in browser E2E: two
    // page documents can align their table node at the same child index, and a
    // reuse of the mounted SchemaTable carried the FIRST table's visible columns
    // into the second (the users table came back showing only the columns the
    // users and roles schemas have in common). The renderer now keys a table
    // node by its id, so a different table mounts a different instance.
    const docFor = (pageId: string, tableId: string, middleColumn: string): RenderPageDocument =>
      ({
        meta: { pageId, title: pageId, protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
        body: {
          type: "section",
          id: pageId,
          children: [
            {
              type: "custom",
              id: `${tableId}-trigger`,
              component: "w33-probe",
              props: { targetTable: tableId, slot: "list-page-actions" },
            },
            {
              type: "table",
              id: tableId,
              props: {
                columns: [
                  { field: "id", label: "ID" },
                  { field: middleColumn.toLowerCase(), label: middleColumn },
                  { field: "updatedAt", label: "Updated" },
                ],
                dataSource: `/api/${pageId}`,
              },
            },
          ],
        },
      }) as unknown as RenderPageDocument;

    const container = globalThis.document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    const fetcher = (async () =>
      new Response(JSON.stringify({ items: [{ id: "1" }], total: 1, page: 1, pageSize: 20 }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      })) as typeof fetch;
    const render = async (doc: RenderPageDocument) => {
      await act(async () => {
        root.render(
          <I18nProvider stored="en-US">
            <RenderPage
              document={doc}
              context={{}}
              dataFetcher={fetcher}
              tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
            />
          </I18nProvider>,
        );
      });
      await act(async () => {
        await Promise.resolve();
      });
    };
    const headers = () =>
      Array.from(container.querySelectorAll("thead th")).map((cell) => (cell.textContent ?? "").trim());

    await render(docFor("users", "users-table", "Username"));
    expect(headers()).toEqual(["ID", "Username", "Updated"]);
    await render(docFor("roles", "roles-table", "Key"));
    expect(headers()).toEqual(["ID", "Key", "Updated"]);
    // Back to users: the column set must be intact, never the intersection of
    // the two schemas (["ID", "Updated"]).
    await render(docFor("users", "users-table", "Username"));
    expect(headers()).toEqual(["ID", "Username", "Updated"]);
  });

  it("places the shipped users and roles triggers through the slot", () => {
    // The two shipped pages are the actual consumers; assert their declarations
    // so a future edit cannot quietly move the trigger back above the filters.
    for (const [file, nodeId, tableId, resource] of [
      ["users/schema/users.json", "users-batch-export", "users-table", "users"],
      ["roles/schema/roles.json", "roles-batch-export", "roles-table", "roles"],
    ] as const) {
      const doc = JSON.parse(readFileSync(resolve(MODULES, file), "utf8")) as {
        body: { children: Array<{ id?: string; props?: Record<string, unknown> }> };
      };
      const node = doc.body.children.find((child) => child.id === nodeId);
      expect(node, `${file} must declare ${nodeId}`).toBeDefined();
      expect(node!.props?.slot).toBe("list-page-actions");
      expect(node!.props?.targetTable).toBe(tableId);
      expect(node!.props?.resource).toBe(resource);
    }
  });
});
