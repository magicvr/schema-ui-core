// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { RenderTableNode } from "@/renderer/render";
import {
  SchemaTable,
  schemaTableColumns,
  schemaTableDataSource,
} from "@/renderer/schema-table";

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
    const node = tableNode({ columns: COLUMNS, dataSource: "/api/records" });
    expect(schemaTableColumns(node).map((column) => column.field)).toEqual([
      "name",
      "status",
    ]);
    expect(schemaTableDataSource(node)).toBe("/api/records");
  });

  it("defaults the data source to the demo records API", () => {
    const node = tableNode({ columns: COLUMNS });
    expect(schemaTableDataSource(node)).toBe("/api/records");
  });

  it("renders records from the injected data source", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/records" }),
      recordsFetcher(),
    );
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).toContain("Northwind Sales");
    expect(container.textContent).toContain("2 records · page 1 of 1");
  });

  it("fails closed when the table node declares no columns", async () => {
    const container = await renderTable(tableNode({}), recordsFetcher());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "table node requires a columns array",
    );
  });

  it("surfaces a fail-closed error when the data source request fails", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/records" }),
      recordsFetcher(500),
    );
    expect(container.textContent).toContain("records fetch failed");
  });

  it("toggles column sort and marks the active column", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/records" }),
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
