// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { RenderTableNode } from "@/renderer/render";
import {
  SchemaTable,
  rowActionDisabled,
  schemaTableColumns,
  schemaTableDataSource,
  schemaTableRowKey,
} from "@/renderer/schema-table";

describe("rowActionDisabled", () => {
  it("evaluates a structured row field equality condition", () => {
    const action = { disabledWhen: { field: "editable", equals: false } };
    expect(rowActionDisabled(action, { id: "role-admin", editable: false })).toBe(true);
    expect(rowActionDisabled(action, { id: "role-ops", editable: true })).toBe(false);
  });

  it("fails closed on malformed conditions", () => {
    expect(rowActionDisabled({ disabledWhen: "editable" }, { id: "role-admin" })).toBe(true);
    expect(rowActionDisabled({ disabledWhen: { field: "" } }, { id: "role-admin" })).toBe(true);
    expect(
      rowActionDisabled({ disabledWhen: { field: "editable", equals: false } }, { id: "role-admin" }),
    ).toBe(true);
  });
});

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

function recordsFetcher(status = 200) {
  return (async (input: RequestInfo | URL) => {
    const url = new URL(String(input), "http://test.local");
    const q = url.searchParams.get("q");
    const items = q
      ? RECORDS.items.filter((record) =>
          record.name.toLowerCase().includes(q.toLowerCase()),
        )
      : RECORDS.items;
    if (status !== 200) {
      return new Response(JSON.stringify({ error: "records down" }), { status });
    }
    return new Response(
      JSON.stringify({ ...RECORDS, items, total: items.length }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  }) as typeof fetch;
}

/** A fetcher that returns arbitrary resource rows from a fixed envelope. */
function itemsFetcher(items: Array<Record<string, unknown>>) {
  return (async () =>
    new Response(
      JSON.stringify({ items, total: items.length, page: 1, pageSize: 10 }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    )) as typeof fetch;
}

function tableNode(props: Record<string, unknown>): RenderTableNode {
  return { type: "table", id: "records-table", props: props as RenderTableNode["props"] };
}

const COLUMNS = [
  { field: "name", label: "Name", sortable: true },
  { field: "status", label: "Status", sortable: true },
];

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderTable(node: RenderTableNode, fetcher: typeof fetch) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<SchemaTable node={node} fetcher={fetcher} />);
  });
  return container;
}

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

describe("SchemaTable (R1 list-data injection)", () => {
  it("reads column specs and the data source from the table node props", () => {
    const node = tableNode({ columns: COLUMNS, dataSource: "/api/users" });
    expect(schemaTableColumns(node).map((column) => column.field)).toEqual([
      "name",
      "status",
    ]);
    expect(schemaTableDataSource(node)).toBe("/api/users");
  });

  it("fails closed (null) when the table node has no data source (no /api/users fallback)", () => {
    const node = tableNode({ columns: COLUMNS });
    expect(schemaTableDataSource(node)).toBeNull();
  });

  it("fails closed (null) on a non-single-slash data source (F-001)", () => {
    const node = tableNode({ columns: COLUMNS, dataSource: "//evil.example/x" });
    expect(schemaTableDataSource(node)).toBeNull();
  });

  it("renders records from the injected data source", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      recordsFetcher(),
    );
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).toContain("Northwind Sales");
    expect(container.textContent).toContain("2 items · page 1 of 1");
  });

  it("fails closed when the table node declares no columns", async () => {
    const container = await renderTable(tableNode({}), recordsFetcher());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "table node requires a columns array",
    );
  });

  it("fails closed without fetching when the data source is missing (F-001)", async () => {
    const fetcher = vi.fn(recordsFetcher());
    const container = await renderTable(tableNode({ columns: COLUMNS }), fetcher);
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "valid dataSource",
    );
    expect(fetcher).not.toHaveBeenCalled();
  });

  it("surfaces a fail-closed error when the data source request fails", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      recordsFetcher(500),
    );
    expect(container.textContent).toContain("resource fetch failed");
  });

  it("toggles column sort and marks the active column", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      recordsFetcher(),
    );
    const nameHeader = Array.from(container.querySelectorAll("th button")).find((button) =>
      button.textContent?.includes("Name"),
    );
    expect(nameHeader).not.toBeUndefined();
    await act(async () => (nameHeader as HTMLButtonElement).click());
    expect((nameHeader as HTMLButtonElement).getAttribute("aria-sort")).toBe(
      "ascending",
    );
  });
});

describe("SchemaTable rowKey (F-002 · I-010-001 v0.2.0 §3)", () => {
  it("defaults the row key field to id", () => {
    expect(schemaTableRowKey(tableNode({}))).toBe("id");
  });

  it("reads a non-id rowKey from the table node props", () => {
    expect(schemaTableRowKey(tableNode({ rowKey: "sku" }))).toBe("sku");
  });

  it("renders a new entity with a non-id row key (positive)", async () => {
    const node = tableNode({
      columns: [{ field: "title", label: "Title" }],
      dataSource: "/api/catalog",
      rowKey: "sku",
    });
    const container = await renderTable(
      node,
      itemsFetcher([
        { sku: "S-1", title: "Widget", price: 19 },
        { sku: "S-2", title: "Gadget", price: 29 },
      ]),
    );
    expect(container.textContent).toContain("Widget");
    expect(container.textContent).toContain("Gadget");
    expect(container.textContent).toContain("2 items · page 1 of 1");
  });

  it("fails closed on a missing row key", async () => {
    const node = tableNode({
      columns: [{ field: "title", label: "Title" }],
      dataSource: "/api/catalog",
      rowKey: "sku",
    });
    const container = await renderTable(
      node,
      itemsFetcher([{ title: "Widget" }, { sku: "S-2", title: "Gadget" }]),
    );
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      'no valid "sku" key',
    );
    expect(container.querySelectorAll("tbody tr").length).toBe(0);
  });

  it("fails closed on duplicate row keys", async () => {
    const node = tableNode({
      columns: [{ field: "title", label: "Title" }],
      dataSource: "/api/catalog",
      rowKey: "sku",
    });
    const container = await renderTable(
      node,
      itemsFetcher([
        { sku: "S-1", title: "Widget" },
        { sku: "S-1", title: "Gadget" },
      ]),
    );
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      'duplicate row key "S-1"',
    );
    expect(container.querySelectorAll("tbody tr").length).toBe(0);
  });

  it("fails closed on a non-scalar row key (wrong type)", async () => {
    const node = tableNode({
      columns: [{ field: "title", label: "Title" }],
      dataSource: "/api/catalog",
      rowKey: "sku",
    });
    const container = await renderTable(
      node,
      itemsFetcher([{ sku: {}, title: "Widget" }]),
    );
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      'no valid "sku" key',
    );
    expect(container.querySelectorAll("tbody tr").length).toBe(0);
  });
});
