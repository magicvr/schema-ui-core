// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { RenderTableNode } from "@/renderer/render.types";
import {
  SchemaTable,
  pagerPages,
  rowActionDisabled,
  schemaTableColumns,
  schemaTableDataSource,
  schemaTableFilters,
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

const SAMPLE_ROWS = {
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

function rowsFetcher(status = 200) {
  return (async (input: RequestInfo | URL) => {
    const url = new URL(String(input), "http://test.local");
    const q = url.searchParams.get("q");
    const items = q
      ? SAMPLE_ROWS.items.filter((item) =>
          item.name.toLowerCase().includes(q.toLowerCase()),
        )
      : SAMPLE_ROWS.items;
    if (status !== 200) {
      return new Response(JSON.stringify({ error: "resource down" }), { status });
    }
    return new Response(
      JSON.stringify({ ...SAMPLE_ROWS, items, total: items.length }),
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
  return { type: "table", id: "schema-table", props: props as RenderTableNode["props"] };
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

  it("carries the column truncate flag into the shipped table cell (W4 · GOAL-005)", async () => {
    const truncateColumns = [
      { field: "name", label: "Name", sortable: true },
      { field: "permissions", label: "Permissions", truncate: true },
    ];
    const fetcher = itemsFetcher([
      { id: "role-1", name: "Admin", permissions: ["users.read", "roles.write"] },
    ]);
    const container = await renderTable(
      tableNode({ columns: truncateColumns, dataSource: "/api/roles" }),
      fetcher,
    );
    const cell = Array.from(
      container.querySelectorAll('[data-table-cell="truncated"]'),
    ).find((entry) => entry.getAttribute("title") === "users.read,roles.write");
    expect(cell).not.toBeUndefined();
  });

  it("fails closed (null) when the table node has no data source (no /api/users fallback)", () => {
    const node = tableNode({ columns: COLUMNS });
    expect(schemaTableDataSource(node)).toBeNull();
  });

  it("fails closed (null) on a non-single-slash data source (F-001)", () => {
    const node = tableNode({ columns: COLUMNS, dataSource: "//evil.example/x" });
    expect(schemaTableDataSource(node)).toBeNull();
  });

  it("renders rows from the injected data source", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      rowsFetcher(),
    );
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).toContain("Northwind Sales");
    expect(container.textContent).toContain("2 items · page 1 of 1");
  });

  it("formats currency columns from cent values (W16-F04)", async () => {
    const currencyColumns = [
      { field: "amount", label: "Amount", format: "currency" as const },
    ];
    const fetcher = itemsFetcher([{ id: "wallet-1", amount: 123456 }]);
    const container = await renderTable(
      tableNode({ columns: currencyColumns, dataSource: "/api/wallet" }),
      fetcher,
    );
    expect(container.textContent).toContain("1,234.56");
  });

  it("fails closed when the table node declares no columns", async () => {
    const container = await renderTable(tableNode({}), rowsFetcher());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "table node requires a columns array",
    );
  });

  it("fails closed without fetching when the data source is missing (F-001)", async () => {
    const fetcher = vi.fn(rowsFetcher());
    const container = await renderTable(tableNode({ columns: COLUMNS }), fetcher);
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "valid dataSource",
    );
    expect(fetcher).not.toHaveBeenCalled();
  });

  it("surfaces a fail-closed error when the data source request fails", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      rowsFetcher(500),
    );
    expect(container.textContent).toContain("resource fetch failed");
  });

  it("retries the resource fetch from the table error state (W15-F02)", async () => {
    let calls = 0;
    const fetcher = (async () => {
      calls += 1;
      if (calls === 1) {
        return new Response(JSON.stringify({ error: "resource down" }), { status: 500 });
      }
      return new Response(JSON.stringify(SAMPLE_ROWS), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }) as typeof fetch;
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      fetcher,
    );
    expect(container.querySelector("[data-table-retry]")).not.toBeNull();
    await act(async () => {
      (container.querySelector("[data-table-retry]") as HTMLButtonElement).click();
    });
    expect(calls).toBe(2);
    expect(container.textContent).toContain("Acme Console");
  });

  it("toggles column sort and marks the active column", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      rowsFetcher(),
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

// --- Title / filters / pager (account page sessions-table surface) ---

const MANY_ROWS = Array.from({ length: 25 }, (_, index) => ({
  id: "rec-" + String(index + 1),
  name: "Item " + String(index + 1),
  status: index % 2 === 0 ? "active" : "revoked",
}));

/** Server-side pagination emulator that also honours a status filter param. */
function pagedFetcher(calls: string[] = []) {
  return (async (input: RequestInfo | URL) => {
    const url = new URL(String(input), "http://test.local");
    calls.push(url.search);
    const page = Number(url.searchParams.get("page") ?? "1");
    const pageSize = Number(url.searchParams.get("pageSize") ?? "10");
    const status = url.searchParams.get("status");
    const items =
      status === null
        ? MANY_ROWS
        : MANY_ROWS.filter((row) => row.status === status);
    const start = (page - 1) * pageSize;
    return new Response(
      JSON.stringify({
        items: items.slice(start, start + pageSize),
        total: items.length,
        page,
        pageSize,
      }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  }) as typeof fetch;
}

describe("pagerPages", () => {
  it("lists every page when the run is short", () => {
    expect(pagerPages(1, 1)).toEqual([1]);
    expect(pagerPages(3, 7)).toEqual([1, 2, 3, 4, 5, 6, 7]);
  });

  it("windows long page runs with gap markers", () => {
    expect(pagerPages(1, 10)).toEqual([1, 2, "gap", 10]);
    expect(pagerPages(5, 10)).toEqual([1, "gap", 4, 5, 6, "gap", 10]);
    expect(pagerPages(10, 10)).toEqual([1, "gap", 9, 10]);
  });

  it("clamps out-of-range current pages", () => {
    expect(pagerPages(0, 10)).toEqual([1, 2, "gap", 10]);
    expect(pagerPages(99, 10)).toEqual([1, "gap", 9, 10]);
  });
});

describe("SchemaTable title / filters / pager", () => {
  it("renders a title heading from titleKey/title", async () => {
    const container = await renderTable(
      tableNode({
        columns: COLUMNS,
        dataSource: "/api/users",
        title: "Signed-in sessions",
        titleKey: "schema.account.session.title",
      }),
      rowsFetcher(),
    );
    const heading = container.querySelector("h2");
    expect(heading?.textContent).toBe("Signed-in sessions");
  });

  it("parses only well-formed select filters (fail-closed on malformed entries)", () => {
    const node = tableNode({
      columns: COLUMNS,
      dataSource: "/api/users",
      filters: [
        { field: "status", type: "select", options: [{ value: "active" }] },
        { field: "bogus", type: "checkbox" },
        { field: "" },
        { field: "noType" },
        42,
        { field: "badOptions", type: "select", options: [42, { value: "ok" }] },
      ],
    });
    const filters = schemaTableFilters(node);
    expect(filters.map((filter) => filter.field)).toEqual([
      "status",
      "badOptions",
    ]);
    expect(filters[0].options.map((option) => option.value)).toEqual(["active"]);
    expect(filters[1].options.map((option) => option.value)).toEqual(["ok"]);
  });

  it("renders a status filter and refetches with the status param on change", async () => {
    const calls: string[] = [];
    const node = tableNode({
      columns: COLUMNS,
      dataSource: "/api/users",
      filters: [
        {
          field: "status",
          type: "select",
          labelKey: "schema.account.filter.status",
          options: [
            { value: "", labelKey: "schema.account.filter.status.all" },
            { value: "active", labelKey: "schema.account.filter.status.active" },
            { value: "revoked", labelKey: "schema.account.filter.status.revoked" },
          ],
        },
      ],
    });
    const container = await renderTable(node, pagedFetcher(calls));
    const select = container.querySelector(
      "[data-table-filters] select",
    ) as HTMLSelectElement;
    expect(select).not.toBeNull();
    expect(select.value).toBe("");
    await act(async () => {
      select.value = "active";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    expect(calls.some((query) => query.includes("status=active"))).toBe(true);
    expect(container.querySelectorAll("tbody tr").length).toBe(10);
    // Filter change resets to page 1 (no page=2 request for the filtered list).
    expect(
      calls.some(
        (query) => query.includes("status=active") && query.includes("page=2"),
      ),
    ).toBe(false);
  });

  it("turns pages with the pager controls", async () => {
    const calls: string[] = [];
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      pagedFetcher(calls),
    );
    expect(container.textContent).toContain("25 items · page 1 of 3");
    const prev = container.querySelector(
      '[aria-label="Previous page"]',
    ) as HTMLButtonElement;
    const next = container.querySelector(
      '[aria-label="Next page"]',
    ) as HTMLButtonElement;
    expect(prev.disabled).toBe(true);
    expect(next.disabled).toBe(false);
    await act(async () => next.click());
    expect(calls.some((query) => query.includes("page=2"))).toBe(true);
    expect(container.textContent).toContain("25 items · page 2 of 3");
    // The current page button carries aria-current="page".
    const current = container.querySelector('[aria-current="page"]');
    expect(current?.textContent).toBe("2");
    // Jump to the last page via the pager number button.
    const pageButtons = Array.from(container.querySelectorAll("nav button"));
    await act(async () =>
      (pageButtons[pageButtons.length - 1] as HTMLButtonElement).click(),
    );
    expect(calls.some((query) => query.includes("page=3"))).toBe(true);
    expect(
      (container.querySelector('[aria-label="Next page"]') as HTMLButtonElement)
        .disabled,
    ).toBe(true);
  });

  it("keeps the current rows rendered while paginating (no skeleton swap)", async () => {
    let resolvePage2: (value: Response) => void = () => {};
    const calls: string[] = [];
    const controlled = (async (input: RequestInfo | URL) => {
      const url = new URL(String(input), "http://test.local");
      calls.push(url.search);
      const page = Number(url.searchParams.get("page") ?? "1");
      if (page === 1) {
        return new Response(
          JSON.stringify({ items: MANY_ROWS.slice(0, 10), total: 25, page: 1, pageSize: 10 }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      }
      return new Promise<Response>((resolve) => {
        resolvePage2 = resolve;
      });
    }) as typeof fetch;

    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      controlled,
    );
    expect(container.textContent).toContain("Item 1");
    const next = container.querySelector(
      '[aria-label="Next page"]',
    ) as HTMLButtonElement;
    await act(async () => next.click());
    // While page 2 is in flight the old rows stay rendered: no loading
    // skeleton swap, so the list height (and the scroll anchor) is stable.
    expect(container.querySelector('[data-table-presentation="loading"]')).toBeNull();
    expect(container.textContent).toContain("Item 1");
    await act(async () => {
      resolvePage2(
        new Response(
          JSON.stringify({ items: MANY_ROWS.slice(10, 20), total: 25, page: 2, pageSize: 10 }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        ),
      );
    });
    expect(container.textContent).toContain("Item 11");
    // Page-1-only rows gone. Row cells are rendered as "Item N<status>", so a
    // page-1 row starts with "Item 1a" (active) or "Item 1r" (revoked) while
    // page-2 rows start with "Item 1<digit>".
    const rowTexts = Array.from(container.querySelectorAll("tbody tr")).map(
      (row) => row.textContent ?? "",
    );
    expect(rowTexts.length).toBe(10);
    expect(rowTexts.every((text) => !/^Item 1[ar]/.test(text))).toBe(true);
    expect(calls.some((query) => query.includes("page=2"))).toBe(true);
  });

  it("hides the pager when everything fits on one page", async () => {
    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      rowsFetcher(),
    );
    expect(container.querySelector("nav")).toBeNull();
  });

  // W11 · U-06: pageSize switcher + go-to-page control.
  it("switches the per-page size (resets to page 1) and jumps to a page", async () => {
    const calls: string[] = [];
    const controlled = (async (input: RequestInfo | URL) => {
      const url = new URL(String(input), "http://test.local");
      calls.push(url.search);
      const pageSize = Number(url.searchParams.get("pageSize") ?? "10");
      const page = Number(url.searchParams.get("page") ?? "1");
      return new Response(
        JSON.stringify({
          items: MANY_ROWS.slice((page - 1) * pageSize, page * pageSize),
          total: 25,
          page,
          pageSize,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }) as typeof fetch;

    const container = await renderTable(
      tableNode({ columns: COLUMNS, dataSource: "/api/users" }),
      controlled,
    );
    // Page-size switch → page=1 + pageSize=20.
    const sizeSelect = container.querySelector<HTMLSelectElement>('[aria-label="Rows per page"]');
    // Go-to-page jump first: page 3 is valid at the default pageSize of 10
    // (25 items → 3 pages); after switching to 20 per page only 2 pages exist.
    const goToForm = container.querySelector<HTMLFormElement>('[aria-label="Go to page"]');
    expect(goToForm).not.toBeNull();
    const input = goToForm!.querySelector<HTMLInputElement>('input[type="number"]');
    expect(input).not.toBeNull();
    await act(async () => {
      Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set?.call(
        input,
        "3",
      );
      input!.dispatchEvent(new Event("input", { bubbles: true }));
      goToForm!.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true }));
    });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(calls.some((query) => query.includes("page=3"))).toBe(true);

    // Page-size switch → page=1 (default, omitted) + pageSize=20; the jump to
    // page 3 is no longer valid at 20 per page (2 pages), so no page=3+20 call.
    await act(async () => {
      Object.getOwnPropertyDescriptor(window.HTMLSelectElement.prototype, "value")?.set?.call(
        sizeSelect,
        "20",
      );
      sizeSelect!.dispatchEvent(new Event("change", { bubbles: true }));
    });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 0));
    });
    expect(calls.some((query) => query.includes("pageSize=20"))).toBe(true);
    expect(calls.some((query) => query.includes("page=2") && query.includes("pageSize=20"))).toBe(false);
  });

  // W11 · U-05 fix: the overflow menu is portaled to document.body and fixed to
  // the viewport, so the table's overflow-x-auto container cannot clip it.
  it("portals the overflow menu to document.body (not clipped by the table)", async () => {
    const node = tableNode({
      columns: COLUMNS,
      dataSource: "/api/users",
      actions: [
        { key: "edit", label: "Edit", actionRef: "edit", requestMapping: { path: { id: "$row.id" } } },
        { key: "delete", label: "Delete", actionRef: "delete", requestMapping: { path: { id: "$row.id" } } },
        { key: "enable", label: "Enable", actionRef: "enable", requestMapping: { path: { id: "$row.id" } } },
        { key: "disable", label: "Disable", actionRef: "disable", requestMapping: { path: { id: "$row.id" } } },
      ],
    });
    const container = await renderTable(node, rowsFetcher());
    // Edit + Delete stay inline (MAX_INLINE_ROW_ACTIONS = 2)…
    expect(container.textContent).toContain("Edit");
    expect(container.textContent).toContain("Delete");
    // …Enable/Disable live in the "⋯" menu, rendered OUTSIDE the table DOM.
    const trigger = container.querySelector<HTMLButtonElement>(
      '[data-row-actions-menu] button[aria-label="More actions"]',
    );
    expect(trigger).not.toBeNull();
    await act(async () => trigger!.click());
    const menu = document.body.querySelector<HTMLElement>('[role="menu"]');
    expect(menu).not.toBeNull();
    expect(menu!.textContent).toContain("Enable");
    expect(menu!.textContent).toContain("Disable");
    // The menu is not a descendant of the table container.
    expect(container.querySelector('[role="menu"]')).toBeNull();
    // Fixed viewport positioning is applied (jsdom rects are zero, so the
    // menu anchors at the trigger's bottom + 4 = 4).
    expect(menu!.style.position).toBe("fixed");
    expect(Number(menu!.style.zIndex)).toBeGreaterThanOrEqual(50);
    // Escape closes it again.
    await act(async () => {
      document.dispatchEvent(new KeyboardEvent("keydown", { key: "Escape", bubbles: true }));
    });
    expect(document.body.querySelector('[role="menu"]')).toBeNull();
  });

  // W23 (GOAL-034): scroll closes the menu only when the trigger actually
  // MOVED. A scroll that leaves the anchor untouched — e.g. a dialog-close
  // focus-restore scrolling a row into view, or an unrelated inner container
  // scroll — must not rip the freshly opened menu away (e2e flake).
  it("keeps the portaled menu open on a scroll that does not move the trigger, and closes when it moves", async () => {
    const node = tableNode({
      columns: COLUMNS,
      dataSource: "/api/users",
      actions: [
        { key: "edit", label: "Edit", actionRef: "edit" },
        { key: "delete", label: "Delete", actionRef: "delete" },
        { key: "enable", label: "Enable", actionRef: "enable" },
      ],
    });
    const container = await renderTable(node, rowsFetcher());
    const trigger = container.querySelector<HTMLButtonElement>(
      '[data-row-actions-menu] button[aria-label="More actions"]',
    );
    expect(trigger).not.toBeNull();
    // Pin a stable rect so the open-time anchor is deterministic in jsdom.
    const stableRect = {
      top: 100,
      bottom: 132,
      left: 40,
      right: 72,
      width: 32,
      height: 32,
      x: 40,
      y: 100,
      toJSON(): Record<string, number> {
        return {
          top: this.top,
          bottom: this.bottom,
          left: this.left,
          right: this.right,
          width: this.width,
          height: this.height,
          x: this.x,
          y: this.y,
        };
      },
    };
    Object.defineProperty(trigger!, "getBoundingClientRect", {
      configurable: true,
      value: () => stableRect,
    });
    await act(async () => trigger!.click());
    expect(document.body.querySelector('[role="menu"]')).not.toBeNull();

    // A scroll that does not move the trigger keeps the menu open.
    await act(async () => {
      document.dispatchEvent(new Event("scroll", { bubbles: true }));
    });
    expect(document.body.querySelector('[role="menu"]')).not.toBeNull();

    // Once the trigger actually moves, the anchor is invalidated → close.
    stableRect.top = 140;
    stableRect.bottom = 172;
    await act(async () => {
      document.dispatchEvent(new Event("scroll", { bubbles: true }));
    });
    expect(document.body.querySelector('[role="menu"]')).toBeNull();

    // Outside pointerdown still closes (unchanged contract).
    await act(async () => trigger!.click());
    expect(document.body.querySelector('[role="menu"]')).not.toBeNull();
    await act(async () => {
      document.dispatchEvent(new MouseEvent("pointerdown", { bubbles: true }));
    });
    expect(document.body.querySelector('[role="menu"]')).toBeNull();
  });
});
