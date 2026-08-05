// @vitest-environment jsdom

import { useState } from "react";
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

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
        <DataTable columns={columns} rows={rows} rowKey={rowKey} sort={sort} onSortChange={setSort} />
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
    expect(loadingContainer.textContent).toContain("Loading…");
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
});
