// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import type { RenderPageDocument } from "@/renderer/render.types";
import { RenderPage } from "@/renderer/render.tsx";
import { savedViewStorageKey } from "@/renderer/saved-views";
import { SchemaTable } from "@/renderer/schema-table";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

const documentFixture = {
  meta: {
    pageId: "saved-view-users",
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest"],
  },
  body: {
    type: "table",
    id: "users-table",
    props: {
      dataSource: "/api/users",
      columns: [
        { field: "name", label: "Name", sortable: true },
        { field: "status", label: "Status", sortable: true },
      ],
      filters: [
        {
          field: "status",
          type: "select",
          label: "Status",
          options: [
            { value: "", label: "All" },
            { value: "active", label: "Active" },
          ],
        },
      ],
    },
  },
} as unknown as RenderPageDocument;

function fetcherFor(calls: string[]): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const raw = String(input);
    calls.push(raw);
    const url = new URL(raw, "http://test.local");
    return new Response(
      JSON.stringify({
        items: [{ id: "u-1", name: "Alice", status: url.searchParams.get("status") ?? "active" }],
        total: 1,
        page: Number(url.searchParams.get("page") ?? "1"),
        pageSize: Number(url.searchParams.get("pageSize") ?? "10"),
      }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  }) as typeof fetch;
}

async function renderFixture(calls: string[], userId = "user-1") {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  const resourceFetcher = fetcherFor(calls);
  await act(async () => {
    root.render(
      <RenderPage
        document={documentFixture}
        context={{ user: { id: userId } }}
        tableRenderer={(node) => <SchemaTable node={node} fetcher={resourceFetcher} />}
      />,
    );
  });
  return container;
}

function setSelectValue(select: HTMLSelectElement, value: string): void {
  Object.getOwnPropertyDescriptor(HTMLSelectElement.prototype, "value")?.set?.call(select, value);
  select.dispatchEvent(new Event("change", { bubbles: true }));
}

async function saveNamedView(container: HTMLDivElement, name: string): Promise<void> {
  await act(async () =>
    container.querySelector<HTMLButtonElement>('[data-saved-view-action="save"]')?.click(),
  );
  const nameInput = container.querySelector<HTMLInputElement>("#users-table-saved-view-name");
  expect(nameInput).not.toBeNull();
  await act(async () => {
    Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, "value")?.set?.call(nameInput, name);
    nameInput?.dispatchEvent(new Event("input", { bubbles: true }));
  });
  await act(async () =>
    container.querySelector<HTMLButtonElement>('[data-saved-view-action="confirm-save"]')?.click(),
  );
}

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  localStorage.clear();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  localStorage.clear();
});

describe("Saved View table integration", () => {
  it("saves the allowlisted state and restores it when selected", async () => {
    const calls: string[] = [];
    const container = await renderFixture(calls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    const tableFilter = container.querySelector<HTMLSelectElement>("[data-table-filters] select");
    expect(tableFilter).not.toBeNull();
    await act(async () => setSelectValue(tableFilter!, "active"));

    await saveNamedView(container, "Active users");

    const key = savedViewStorageKey("user-1", "saved-view-users", "users-table");
    const stored = JSON.parse(localStorage.getItem(key) ?? "null") as {
      version: number;
      views: Array<{ name: string; query: Record<string, unknown>; visibleColumns: string[] }>;
      activeViewId?: string;
    };
    expect(stored.version).toBe(1);
    expect(stored.views).toHaveLength(1);
    expect(stored.views[0]).toMatchObject({
      name: "Active users",
      query: { filters: { status: "active" } },
      visibleColumns: ["name", "status"],
    });
    expect(stored.activeViewId).toBeTruthy();

    const reloadCalls: string[] = [];
    await renderFixture(reloadCalls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });
    expect(reloadCalls.some((url) => new URL(url, "http://test.local").searchParams.get("status") === "active")).toBe(true);

    await act(async () => setSelectValue(tableFilter!, ""));
    const savedViewSelect = container.querySelector<HTMLSelectElement>("[data-saved-view-select]");
    expect(savedViewSelect).not.toBeNull();
    await act(async () => setSelectValue(savedViewSelect!, stored.activeViewId!));
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });
    expect(calls.some((url) => new URL(url, "http://test.local").searchParams.get("status") === "active")).toBe(true);
  });

  it("updates and deletes the selected view without leaking stale state", async () => {
    const calls: string[] = [];
    const container = await renderFixture(calls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    const tableFilter = container.querySelector<HTMLSelectElement>("[data-table-filters] select");
    expect(tableFilter).not.toBeNull();
    await act(async () => setSelectValue(tableFilter!, "active"));
    await saveNamedView(container, "Active users");

    const key = savedViewStorageKey("user-1", "saved-view-users", "users-table");
    const initial = JSON.parse(localStorage.getItem(key) ?? "null") as {
      views: Array<{ id: string; query: Record<string, unknown> }>;
      activeViewId?: string;
    };
    expect(initial.views).toHaveLength(1);
    expect(initial.views[0].query).toMatchObject({ filters: { status: "active" } });

    await act(async () => setSelectValue(tableFilter!, ""));
    await act(async () =>
      container.querySelector<HTMLButtonElement>('[data-saved-view-action="update"]')?.click(),
    );
    const updated = JSON.parse(localStorage.getItem(key) ?? "null") as {
      views: Array<{ id: string; query: Record<string, unknown> }>;
      activeViewId?: string;
    };
    expect(updated.views[0].id).toBe(initial.views[0].id);
    expect(updated.views[0].query).toEqual({ pageSize: 10 });
    expect(updated.activeViewId).toBe(initial.activeViewId);

    const confirm = vi.spyOn(window, "confirm").mockReturnValue(true);
    await act(async () =>
      container.querySelector<HTMLButtonElement>('[data-saved-view-action="delete"]')?.click(),
    );
    confirm.mockRestore();
    expect(JSON.parse(localStorage.getItem(key) ?? "null")).toMatchObject({ views: [] });
    expect((JSON.parse(localStorage.getItem(key) ?? "null") as { activeViewId?: string }).activeViewId).toBeUndefined();
  });

  it("shows an observable error and leaves the table query untouched for invalid stored views", async () => {
    const key = savedViewStorageKey("user-1", "saved-view-users", "users-table");
    localStorage.setItem(
      key,
      JSON.stringify({
        version: 1,
        views: [
          {
            id: "bad-view",
            name: "Bad view",
            query: { filters: { unknown: "secret" } },
            visibleColumns: ["name", "status"],
            createdAt: "2026-09-17T00:00:00.000Z",
            updatedAt: "2026-09-17T00:00:00.000Z",
          },
        ],
        activeViewId: "bad-view",
      }),
    );
    const calls: string[] = [];
    const container = await renderFixture(calls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    expect(container.querySelector('[data-feedback-toast="error"]')).not.toBeNull();
    expect(container.querySelector<HTMLSelectElement>("[data-saved-view-select]")?.value).toBe("");
    expect(calls.some((url) => new URL(url, "http://test.local").searchParams.has("unknown"))).toBe(false);
  });

  it("persists column visibility and keeps separate user namespaces", async () => {
    const userACalls: string[] = [];
    const userBCalls: string[] = [];
    const userA = await renderFixture(userACalls, "user-a");
    const userB = await renderFixture(userBCalls, "user-b");
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    const userAColumns = userA.querySelectorAll<HTMLInputElement>('[data-saved-views] input[type="checkbox"]');
    expect(userAColumns).toHaveLength(2);
    await act(async () => userAColumns[1].click());
    await saveNamedView(userA, "Names only");
    await saveNamedView(userB, "All users");

    const userAKey = savedViewStorageKey("user-a", "saved-view-users", "users-table");
    const userBKey = savedViewStorageKey("user-b", "saved-view-users", "users-table");
    const storedA = JSON.parse(localStorage.getItem(userAKey) ?? "null") as {
      views: Array<{ visibleColumns: string[] }>;
      activeViewId?: string;
    };
    const storedB = JSON.parse(localStorage.getItem(userBKey) ?? "null") as {
      views: Array<{ visibleColumns: string[] }>;
      activeViewId?: string;
    };
    expect(userAKey).not.toBe(userBKey);
    expect(storedA.views[0].visibleColumns).toEqual(["name"]);
    expect(storedB.views[0].visibleColumns).toEqual(["name", "status"]);

    const userASelect = userA.querySelector<HTMLSelectElement>("[data-saved-view-select]");
    expect(userASelect).not.toBeNull();
    await act(async () => setSelectValue(userASelect!, ""));
    await act(async () => setSelectValue(userASelect!, storedA.activeViewId!));
    expect(userA.querySelectorAll<HTMLInputElement>('[data-saved-views] input[type="checkbox"]')[1].checked).toBe(false);
  });

  it("surfaces a write failure without claiming the view was saved", async () => {
    const calls: string[] = [];
    const container = await renderFixture(calls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });
    const key = savedViewStorageKey("user-1", "saved-view-users", "users-table");
    const originalSetItem = Storage.prototype.setItem;
    const setItem = vi.spyOn(Storage.prototype, "setItem").mockImplementation(function (this: Storage, storageKey: string, value: string) {
      if (storageKey === key) {
        throw new Error("quota");
      }
      return originalSetItem.call(this, storageKey, value);
    });
    try {
      await saveNamedView(container, "Will fail");
      expect(container.querySelector('[data-feedback-toast="error"]')).not.toBeNull();
      expect(localStorage.getItem(key)).toBeNull();
    } finally {
      setItem.mockRestore();
    }
  });

  // R6 C7 (user 2026-09-18, item 3): the "add view" form was surfacing in the
  // page-level actions row, far from the "保存为视图" button that opened it. It
  // must render inside the view surface, directly beneath the tab group.
  it("renders the new-view form inside the view surface, next to its trigger", async () => {
    const calls: string[] = [];
    const container = await renderFixture(calls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    const viewSurface = container.querySelector('[data-saved-views-surface="true"]');
    const tabs = container.querySelector('[data-saved-view-tabs="true"]');
    expect(viewSurface).not.toBeNull();
    expect(tabs).not.toBeNull();
    expect(container.querySelector("#users-table-saved-view-name")).toBeNull();

    await act(async () =>
      container.querySelector<HTMLButtonElement>('[data-saved-view-action="save"]')?.click(),
    );

    const nameInput = container.querySelector("#users-table-saved-view-name");
    const management = container.querySelector('[data-saved-view-management="true"]');
    const pageActions = container.querySelector('[data-list-page-actions]');
    expect(nameInput).not.toBeNull();
    expect(management).not.toBeNull();
    // The form lives inside the view surface, after the tab group.
    expect(viewSurface?.contains(management as Node)).toBe(true);
    expect(management?.contains(nameInput as Node)).toBe(true);
    expect(tabs?.compareDocumentPosition(management as Node)).toBe(
      Node.DOCUMENT_POSITION_FOLLOWING,
    );
    // The page actions row must not own the view form.
    expect(pageActions?.contains(management as Node) ?? false).toBe(false);
  });

  // R6 C7 (user 2026-09-18, item 1): the reference page pairs "列配置" and
  // "保存视图" with a glyph as visual emphasis.
  it("carries reference-page icons on the column and save-view triggers", async () => {
    const calls: string[] = [];
    const container = await renderFixture(calls);
    await act(async () => {
      await Promise.resolve();
      await Promise.resolve();
    });

    const columnsTrigger = container.querySelector('[data-saved-view-columns-trigger="true"]');
    const saveTrigger = container.querySelector('[data-saved-view-action="save"]');
    expect(columnsTrigger).not.toBeNull();
    expect(saveTrigger).not.toBeNull();
    expect(columnsTrigger?.querySelector("svg")).not.toBeNull();
    expect(saveTrigger?.querySelector("svg")).not.toBeNull();
    // The glyph must not replace the accessible label.
    expect(columnsTrigger?.textContent).toContain("Columns");
    expect(saveTrigger?.textContent).toContain("Save as view");
    expect(columnsTrigger?.querySelector("svg")?.getAttribute("aria-hidden")).toBe("true");
    expect(saveTrigger?.querySelector("svg")?.getAttribute("aria-hidden")).toBe("true");
  });
});
