/**
 * S2/S3 visual fidelity tests (workspace-006 / D-004).
 * Drive real shipped modules — not re-implementations of the surface.
 */
// @vitest-environment jsdom

import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { DataTable } from "@/components/data-table";
import { FormControls } from "@/renderer/form-controls.tsx";
import { RenderPage } from "@/renderer/render.tsx";
import type { RenderPageDocument } from "@/renderer/render";
import { SchemaTable } from "@/renderer/schema-table";

const __dir = dirname(fileURLToPath(import.meta.url));
const formControlsSource = readFileSync(join(__dir, "form-controls.tsx"), "utf-8");
const renderSource = readFileSync(join(__dir, "render.tsx"), "utf-8");

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

async function mount(node: React.ReactNode): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(node);
  });
  return container;
}

describe("S2 form-controls design-system consumption", () => {
  it("imports Input / Label / Textarea primitives from components/ui", () => {
    expect(formControlsSource).toContain('from "@/components/ui/input"');
    expect(formControlsSource).toContain('from "@/components/ui/label"');
    expect(formControlsSource).toContain('from "@/components/ui/textarea"');
    expect(formControlsSource).toContain('data-form-controls="design-system"');
  });

  it("renders input + textarea through real FormControls using design primitives", async () => {
    const container = await mount(
      <FormControls
        fields={[
          { id: "name", label: "Name", type: "input" },
          { id: "notes", label: "Notes", type: "textarea" },
        ]}
        values={{ name: "Ada", notes: "hi" }}
        onChange={() => undefined}
      />,
    );
    expect(container.querySelector('[data-form-controls="design-system"]')).not.toBeNull();
    const name = container.querySelector<HTMLInputElement>("#field-name");
    expect(name).not.toBeNull();
    expect(name?.value).toBe("Ada");
    // Input primitive signature
    expect(name?.className).toMatch(/shadow-sm/);
    const notes = container.querySelector<HTMLTextAreaElement>("#field-notes");
    expect(notes).not.toBeNull();
    expect(notes?.value).toBe("hi");
  });
});

describe("S2 recordView Drawer/Sheet presentation", () => {
  it("implements Drawer/Sheet chrome (not centered Modal-only detail)", () => {
    expect(renderSource).toContain('data-record-view="panel"');
    expect(renderSource).toContain("Record details");
    expect(renderSource).toContain("Close record details");
    expect(renderSource).toMatch(/fixed inset-y-0 right-0/);
    expect(renderSource).toMatch(/role="dialog"/);
    // D-004 mobile band uses md (768), not max-sm (640) alone
    expect(renderSource).toMatch(/max-md:/);
  });

  it("renders static recordView with dialog panel chrome via RenderPage", async () => {
    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "recordView",
        props: { record: { id: "rec-9", status: "active", name: "Ada" } },
      },
    };
    const container = await mount(<RenderPage document={pageDoc} context={{}} />);
    const panel = container.querySelector('[data-record-view="panel"]');
    expect(panel).not.toBeNull();
    expect(panel?.getAttribute("role")).toBe("dialog");
    expect(panel?.getAttribute("aria-label")).toBe("Record details");
    expect(container.textContent).toContain("Ada");
    expect(container.textContent).toContain("active");
    // Static fixtures do not mount the selection backdrop
    expect(container.querySelector('[data-record-view="backdrop"]')).toBeNull();
  });

  it("does not open recordView drawer when row Edit action is clicked", async () => {
    const items = [
      { id: "u-1", name: "Ada Lovelace", role: "admin" },
      { id: "u-2", name: "Grace Hopper", role: "editor" },
    ];
    const fetcher = (async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.startsWith("/api/items")) {
        return new Response(
          JSON.stringify({ items, total: items.length, page: 1, pageSize: 10 }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      }
      return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
    }) as typeof fetch;

    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "grid",
        props: { columns: 1 },
        children: [
          {
            type: "table",
            id: "items",
            props: {
              dataSource: "/api/items",
              rowKey: "id",
              columns: [
                { field: "name", label: "Name" },
                { field: "role", label: "Role" },
              ],
              actions: [
                {
                  key: "edit",
                  label: "Edit",
                  actionRef: "editItem",
                },
              ],
            },
          },
          { type: "recordView", id: "detail", props: {} },
        ],
      },
      actions: {
        editItem: {
          type: "modal",
          content: { type: "text", props: { text: "edit-modal-body" } },
        },
      },
    };

    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <RenderPage
          document={pageDoc}
          context={{}}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />,
      );
    });
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    const editButtons = Array.from(container.querySelectorAll("button")).filter(
      (button) => button.textContent?.trim() === "Edit",
    );
    expect(editButtons.length).toBeGreaterThan(0);
    await act(async () => {
      editButtons[0]!.click();
    });
    // Must NOT open selection-driven drawer; modal path may open for edit
    expect(container.querySelector('[data-record-view="backdrop"]')).toBeNull();
    expect(container.querySelector('[data-record-view-mode="drawer"]')).toBeNull();
  });

  it("opens selection-driven Drawer with backdrop and closes via selectRow(null)", async () => {
    const items = [
      { id: "u-1", name: "Ada Lovelace", role: "admin" },
      { id: "u-2", name: "Grace Hopper", role: "editor" },
    ];
    const fetcher = (async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.startsWith("/api/items")) {
        return new Response(
          JSON.stringify({ items, total: items.length, page: 1, pageSize: 10 }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      }
      return new Response(JSON.stringify({ error: "NOT_FOUND" }), { status: 404 });
    }) as typeof fetch;

    const pageDoc: RenderPageDocument = {
      meta: { protocolVersion: "2.7", requiredCapabilities: ["app.manifest"] },
      body: {
        type: "grid",
        props: { columns: 1 },
        children: [
          {
            type: "table",
            id: "items",
            props: {
              dataSource: "/api/items",
              rowKey: "id",
              columns: [
                { field: "name", label: "Name" },
                { field: "role", label: "Role" },
              ],
            },
          },
          { type: "recordView", id: "detail", props: {} },
        ],
      },
    };

    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <RenderPage
          document={pageDoc}
          context={{}}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />,
      );
    });
    // Allow list fetch to settle
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    expect(container.textContent).toContain("Ada Lovelace");
    expect(container.querySelector('[data-record-view="panel"]')).toBeNull();

    // Click a desktop table row (first data row with name)
    const nameCells = Array.from(container.querySelectorAll("td")).filter((td) =>
      td.textContent?.includes("Ada Lovelace"),
    );
    expect(nameCells.length).toBeGreaterThan(0);
    await act(async () => {
      nameCells[0]!.closest("tr")?.click();
    });

    const panel = container.querySelector('[data-record-view="panel"]');
    expect(panel).not.toBeNull();
    expect(panel?.getAttribute("data-record-view-mode")).toBe("drawer");
    expect(panel?.getAttribute("aria-modal")).toBe("true");
    expect(container.querySelector('[data-record-view="backdrop"]')).not.toBeNull();
    expect(container.textContent).toContain("admin");

    const close = container.querySelector<HTMLButtonElement>(
      'button[aria-label="Close record details"]',
    );
    expect(close).not.toBeNull();
    await act(async () => {
      close!.click();
    });
    expect(container.querySelector('[data-record-view="panel"]')).toBeNull();
    expect(container.querySelector('[data-record-view="backdrop"]')).toBeNull();
  });
});

describe("S2 dual-end DataTable is the shipped list surface", () => {
  it("SchemaTable imports DataTable that exposes dual-end markers", async () => {
    const container = await mount(
      <DataTable
        columns={[
          { key: "id", label: "ID" },
          { key: "name", label: "Name" },
        ]}
        rows={[
          { id: "1", name: "One" },
          { id: "2", name: "Two" },
        ]}
        rowKey={(row) => row.id}
      />,
    );
    expect(container.querySelector('[data-table-presentation="dual-end"]')).not.toBeNull();
    expect(container.querySelector('[data-table-presentation="mobile-cards"]')).not.toBeNull();
    expect(container.querySelector('[data-table-presentation="desktop-table"] table')).not.toBeNull();
  });
});
