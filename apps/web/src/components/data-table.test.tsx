// @vitest-environment jsdom

import { useState } from "react";
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";

interface Row {
  id: string;
  name: string;
}

const columns: DataTableColumn<Row>[] = [
  { key: "id", label: "ID" },
  { key: "name", label: "Name", sortable: true },
];

const rows: Row[] = [
  { id: "rec-1", name: "Acme Console" },
  { id: "rec-2", name: "Northwind Sales" },
];

function rowKey(row: Row): string {
  return row.id;
}

interface LongRow {
  id: string;
  permissions: string[];
}

const longColumns: DataTableColumn<LongRow>[] = [
  { key: "id", label: "ID" },
  { key: "permissions", label: "Permissions", truncate: true },
];

const longRows: LongRow[] = [
  {
    id: "role-1",
    permissions: ["users.read", "users.write", "roles.read", "roles.write"],
  },
];

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderTable(children: React.ReactNode): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => root.render(children));
  return container;
}

function buttonText(container: HTMLElement, label: string): HTMLButtonElement | null {
  const buttons = Array.from(container.querySelectorAll("button"));
  return buttons.find((button) => button.textContent?.includes(label)) ?? null;
}

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

describe("DataTable", () => {
  it("renders headers and cell values", async () => {
    const container = await renderTable(
      <DataTable columns={columns} rows={rows} rowKey={rowKey} caption="Example rows" />,
    );
    expect(container.textContent).toContain("ID");
    expect(container.textContent).toContain("Name");
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).toContain("Northwind Sales");
  });

  it("toggles sort asc → desc on a sortable header", async () => {
    function Harness() {
      const [sort, setSort] = useState<SortState | undefined>(undefined);
      return (
        <DataTable columns={columns} rows={rows} rowKey={rowKey} sort={sort} onSortChange={(next) => setSort(next ?? undefined)} />
      );
    }
    const container = await renderTable(<Harness />);
    const nameHeader = buttonText(container, "Name");
    expect(nameHeader).not.toBeNull();
    await act(async () => nameHeader!.click());
    expect(buttonText(container, "Name")?.getAttribute("aria-sort")).toBe("ascending");
    await act(async () => nameHeader!.click());
    expect(buttonText(container, "Name")?.getAttribute("aria-sort")).toBe("descending");
  });

  it("marks the active sort column as aria-sort ascending", async () => {
    const container = await renderTable(
      <DataTable
        columns={columns}
        rows={rows}
        rowKey={rowKey}
        sort={{ field: "name", order: "asc" }}
      />,
    );
    const button = buttonText(container, "Name");
    expect(button?.getAttribute("aria-sort")).toBe("ascending");
  });

  it("shows the empty message and a loading row", async () => {
    const emptyContainer = await renderTable(
      <DataTable columns={columns} rows={[]} rowKey={rowKey} emptyMessage="No rows" />,
    );
    expect(emptyContainer.textContent).toContain("No rows");

    const loadingContainer = await renderTable(
      <DataTable columns={columns} rows={[]} rowKey={rowKey} loading />,
    );
    expect(loadingContainer.querySelector('[role="status"]')).not.toBeNull();
    expect(loadingContainer.querySelectorAll(".animate-pulse").length).toBeGreaterThan(0);
  });

  it("renders an error row when fetch fails", async () => {
    const container = await renderTable(
      <DataTable
        columns={columns}
        rows={[]}
        rowKey={rowKey}
        error="resource fetch failed: HTTP 500"
      />,
    );
    expect(container.textContent).toContain("resource fetch failed: HTTP 500");
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "resource fetch failed: HTTP 500",
    );
  });

  it("exposes a retry control on the error path (W15-F02)", async () => {
    const onRetry = vi.fn();
    const container = await renderTable(
      <DataTable
        columns={columns}
        rows={[]}
        rowKey={rowKey}
        error="resource fetch failed: HTTP 500"
        onRetry={onRetry}
      />,
    );
    const retry = container.querySelector("[data-table-retry]") as HTMLButtonElement | null;
    expect(retry).not.toBeNull();
    retry?.click();
    expect(onRetry).toHaveBeenCalledTimes(1);
  });

  it("renders a custom cell renderer", async () => {
    const customColumns: DataTableColumn<Row>[] = [
      { key: "id", label: "ID" },
      {
        key: "name",
        label: "Name",
        render: (row) => <strong>{String(row.name)}!</strong>,
      },
    ];
    const container = await renderTable(
      <DataTable columns={customColumns} rows={rows} rowKey={rowKey} />,
    );
    expect(container.textContent).toContain("Acme Console!");
  });

  it("does not fire onRowClick when an action button is clicked", async () => {
    const clicks: string[] = [];
    const actionColumns: DataTableColumn<Row>[] = [
      { key: "id", label: "ID" },
      { key: "name", label: "Name" },
      {
        key: "actions",
        label: "",
        render: (row) => (
          <button type="button" data-testid={`edit-${row.id}`}>
            Edit
          </button>
        ),
      },
    ];
    const container = await renderTable(
      <DataTable
        columns={actionColumns}
        rows={rows}
        rowKey={rowKey}
        onRowClick={(row) => {
          clicks.push(row.id);
        }}
      />,
    );
    const edit = container.querySelector<HTMLButtonElement>('[data-testid="edit-rec-1"]');
    expect(edit).not.toBeNull();
    await act(async () => {
      edit!.click();
    });
    expect(clicks).toEqual([]);
    // Data cell still selects
    const nameCell = Array.from(container.querySelectorAll("td")).find((td) =>
      td.textContent?.includes("Acme Console"),
    );
    expect(nameCell).not.toBeUndefined();
    await act(async () => {
      nameCell!.click();
    });
    expect(clicks).toEqual(["rec-1"]);
  });

  it("ships dual-end presentation: dense desktop table + mobile card list (D-004 / S2)", async () => {
    const container = await renderTable(
      <DataTable columns={columns} rows={rows} rowKey={rowKey} caption="Example rows" />,
    );
    const root = container.querySelector('[data-table-presentation="dual-end"]');
    expect(root).not.toBeNull();
    const desktop = container.querySelector('[data-table-presentation="desktop-table"]');
    expect(desktop).not.toBeNull();
    expect(desktop?.className).toMatch(/hidden/);
    expect(desktop?.className).toMatch(/md:block/);
    expect(desktop?.querySelector("table")).not.toBeNull();

    const mobile = container.querySelector('[data-table-presentation="mobile-cards"]');
    expect(mobile).not.toBeNull();
    expect(mobile?.getAttribute("aria-label")).toBe("Mobile card list");
    expect(mobile?.className).toMatch(/md:hidden/);
    expect(mobile?.querySelectorAll("li").length).toBe(2);
    expect(mobile?.textContent).toContain("Acme Console");
    expect(mobile?.textContent).toContain("Northwind Sales");
  });

  it("truncates a truncate column with a full-text title affordance (W4 · GOAL-005)", async () => {
    const container = await renderTable(
      <DataTable columns={longColumns} rows={longRows} rowKey={(row) => row.id} />,
    );
    const full = "users.read,users.write,roles.read,roles.write";
    const cell = Array.from(
      container.querySelectorAll('[data-table-cell="truncated"]'),
    ).find((entry) => entry.getAttribute("title") === full);
    expect(cell).not.toBeUndefined();
    expect(cell?.className).toMatch(/truncate/);
    expect(cell?.className).toMatch(/max-w-\[16rem\]/);
    expect(cell?.getAttribute("title")).toBe(full);
    expect(cell?.textContent).toBe(full);
    // W4 · GOAL-005 / A-003 F-3: the width cap must live on the td too —
    // table-layout auto sizes a column by its cell max-content, so the
    // inner span's max-width alone does not clamp the column.
    const td = cell?.closest("td");
    expect(td?.className).toMatch(/max-w-\[16rem\]/);
  });

  it("truncates every string cell by default with a full-text title (universal ellipsis)", async () => {
    const container = await renderTable(
      <DataTable columns={columns} rows={rows} rowKey={rowKey} />,
    );
    const cells = Array.from(container.querySelectorAll('[data-table-cell="truncated"]'));
    expect(cells.length).toBeGreaterThan(0);
    // Non-truncate columns cap at the universal 20rem; truncate columns keep 16rem.
    const nameCell = cells.find((cell) => cell.getAttribute("title") === "Acme Console");
    expect(nameCell).not.toBeUndefined();
    expect(nameCell?.className).toMatch(/truncate/);
    expect(nameCell?.closest("td")?.className).toMatch(/max-w-\[20rem\]/);
    expect(container.textContent).toContain("Acme Console");
  });

  it("renders a muted placeholder for null / undefined / empty cells", async () => {
    const sparse = [
      { id: "r1", name: "With value", note: null },
      { id: "r2", name: "", note: "has note" },
      { id: "r3" },
    ];
    const container = await renderTable(
      <DataTable
        columns={[
          { key: "id", label: "ID" },
          { key: "name", label: "Name" },
          { key: "note", label: "Note" },
        ]}
        rows={sparse}
        rowKey={(row) => row.id}
      />,
    );
    const rowsEl = Array.from(container.querySelectorAll("tbody tr"));
    expect(rowsEl[0]?.textContent).toContain("—");
    expect(rowsEl[1]?.textContent).toContain("—");
    expect(rowsEl[2]?.textContent).toContain("—");
  });
});
