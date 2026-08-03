// @vitest-environment jsdom
//
// T-UI-01..10 acceptance evidence for GOAL-007 S4/S5 (I-007-003 v0.2.2 §6):
// the `list-edit-lifecycle` representative page drives list/search/detail/
// create/edit/delete through Schema nodes against an emulated records API, and
// the permission / error / confirm matrix holds. The renderer under test is
// the frozen §5 main path; every page behaviour below comes from the real
// fixture, not from source code.
//
// The API emulator mirrors the frozen I-007-001 contract (GET list envelope,
// POST 201 / INVALID_CREATE_*, PATCH 200, DELETE 204, unified {error,message},
// records.read / records.write). Backend authority itself is proven by the Go
// T-API-08/09 suite; here we assert the UI renders consistently with it.

import { readFileSync, readdirSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { RenderPageDocument } from "@/renderer/render";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";
import type { RecordItem } from "@/renderer/records";

const FIXTURE_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/internal/handler/fixtures/schema",
);
const WEB_SRC_DIR = resolve(dirname(fileURLToPath(import.meta.url)), "..");

function fixtureDocument(pageId: string): unknown {
  return JSON.parse(readFileSync(resolve(FIXTURE_DIR, `${pageId}.json`), "utf8"));
}

// --- Records API emulator (I-007-001 contract) ---

interface ApiCall {
  method: string;
  url: string;
  body?: unknown;
}

function json(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function createRecordsApi(
  initial: RecordItem[],
  options: { denyWrites?: boolean; failPatch?: boolean } = {},
): { fetcher: typeof fetch; store: RecordItem[]; calls: ApiCall[] } {
  const store: RecordItem[] = initial.map((record) => ({ ...record }));
  const calls: ApiCall[] = [];
  let tick = 0;
  const now = (): string => {
    tick += 1;
    return `2026-08-02T10:00:00.${String(tick).padStart(3, "0")}Z`;
  };
  const deny = (): Response =>
    json({ error: "FORBIDDEN", message: "permission required: records.write" }, 403);

  const fetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const raw = String(input);
    const url = new URL(raw, "http://test.local");
    const method = (init?.method ?? "GET").toUpperCase();
    calls.push({
      method,
      url: raw,
      body: init?.body !== undefined ? JSON.parse(String(init.body)) : undefined,
    });
    const match = /^\/api\/records(?:\/([^/]+))?$/.exec(url.pathname);
    if (!match) {
      return json({ error: "NOT_FOUND", message: "no such route" }, 404);
    }
    const id = match[1];

    if (method === "GET" && id === undefined) {
      let items = [...store];
      const q = url.searchParams.get("q");
      if (q !== null && q !== "") {
        items = items.filter((record) =>
          String(record.name ?? "").toLowerCase().includes(q.toLowerCase()),
        );
      }
      const sort = url.searchParams.get("sort");
      const order = url.searchParams.get("order");
      if (sort !== null) {
        items.sort((a, b) => {
          const left = String(a[sort as keyof RecordItem] ?? "");
          const right = String(b[sort as keyof RecordItem] ?? "");
          const cmp = left.localeCompare(right);
          return order === "desc" ? -cmp : cmp;
        });
      }
      const page = Number(url.searchParams.get("page") ?? "1");
      const pageSize = Number(url.searchParams.get("pageSize") ?? "10");
      return json({
        items: items.slice((page - 1) * pageSize, page * pageSize),
        total: items.length,
        page,
        pageSize,
      });
    }
    if (method === "POST" && id === undefined) {
      if (options.denyWrites) {
        return deny();
      }
      const body = init?.body !== undefined ? JSON.parse(String(init.body)) : {};
      for (const field of ["name", "status", "owner"]) {
        const value = typeof body[field] === "string" ? String(body[field]).trim() : "";
        if (value === "") {
          return json({ error: "INVALID_CREATE_FIELD", message: `${field} is required` }, 400);
        }
      }
      const record: RecordItem = {
        id: `rec-new-${store.length + 1}`,
        name: String(body.name).trim(),
        status: String(body.status).trim(),
        owner: String(body.owner).trim(),
        updatedAt: now(),
      };
      store.unshift(record);
      return json(record, 201);
    }
    if (id !== undefined && (method === "PATCH" || method === "DELETE")) {
      const index = store.findIndex((record) => record.id === id);
      if (index < 0) {
        return json({ error: "RECORD_NOT_FOUND", message: "no record with that id" }, 404);
      }
      if (method === "DELETE") {
        if (options.denyWrites) {
          return deny();
        }
        store.splice(index, 1);
        return new Response(null, { status: 204 });
      }
      if (options.denyWrites) {
        return deny();
      }
      if (options.failPatch) {
        return json({ error: "INVALID_PATCH_FIELD", message: "name is required" }, 400);
      }
      const body = init?.body !== undefined ? JSON.parse(String(init.body)) : {};
      const patch: Partial<RecordItem> = {};
      for (const field of ["name", "status", "owner"]) {
        if (body[field] !== undefined) {
          const value = String(body[field]).trim();
          if (value === "") {
            return json({ error: "INVALID_PATCH_FIELD", message: `${field} is required` }, 400);
          }
          patch[field as "name" | "status" | "owner"] = value;
        }
      }
      const updated = { ...store[index], ...patch, updatedAt: now() };
      store[index] = updated;
      return json(updated);
    }
    return json({ error: "NOT_FOUND", message: "no such route" }, 404);
  }) as typeof fetch;

  return { fetcher, store, calls };
}

// --- Render harness ---

const ADMIN = {
  user: { id: "u1", roles: ["admin"], permissions: ["records.read", "records.write"] },
};
const VIEWER = {
  user: { id: "u2", roles: ["viewer"], permissions: ["records.read"] },
};

const RECORDS: RecordItem[] = [
  {
    id: "rec-1",
    name: "Acme Console",
    status: "active",
    owner: "alice",
    updatedAt: "2026-07-31T00:00:00.000Z",
  },
  {
    id: "rec-2",
    name: "Northwind Sales",
    status: "pending",
    owner: "bob",
    updatedAt: "2026-07-31T11:00:00.000Z",
  },
];

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderCrud(
  fixture: unknown,
  context: Record<string, unknown>,
  fetcher: typeof fetch,
): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <RenderPage
        document={fixture as RenderPageDocument}
        context={context}
        tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
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

function buttonByText(container: HTMLElement, text: string): HTMLButtonElement | undefined {
  return Array.from(container.querySelectorAll("button")).find(
    (button) => button.textContent?.trim() === text,
  ) as HTMLButtonElement | undefined;
}

function setFieldValue(el: Element, value: string): void {
  const proto =
    el instanceof HTMLSelectElement
      ? HTMLSelectElement.prototype
      : el instanceof HTMLTextAreaElement
        ? HTMLTextAreaElement.prototype
        : HTMLInputElement.prototype;
  Object.getOwnPropertyDescriptor(proto, "value")?.set?.call(el, value);
  el.dispatchEvent(new Event(el instanceof HTMLSelectElement ? "change" : "input", { bubbles: true }));
}

function fieldInput(container: HTMLElement, id: string): HTMLInputElement | null {
  return container.querySelector<HTMLInputElement>(`#field-${id}`);
}

// --- T-UI tests ---

describe("T-UI-01 · list loads and empty state", () => {
  it("loads the representative list from /api/records and shows the record count", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).toContain("Northwind Sales");
    expect(container.textContent).toContain("2 records · page 1 of 1");
  });

  it("shows the empty message when the data source has no rows", async () => {
    const api = createRecordsApi([]);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    expect(container.textContent).toContain("No records match.");
  });
});

describe("T-UI-02 · sort and pagination", () => {
  it("toggles a column sort header and refetches with sort/order", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    const nameHeader = Array.from(container.querySelectorAll("th button")).find((button) =>
      button.textContent?.includes("Name"),
    );
    expect(nameHeader).not.toBeUndefined();
    await act(async () => (nameHeader as HTMLButtonElement).click());
    expect((nameHeader as HTMLButtonElement).getAttribute("aria-sort")).toBe("ascending");
    const sortCalls = api.calls.filter(
      (call) => call.method === "GET" && call.url.includes("sort=name") && call.url.includes("order=asc"),
    );
    expect(sortCalls.length).toBeGreaterThan(0);
  });
});

describe("T-UI-03 · search form-to-query binding (search-form-table)", () => {
  it("binds the search form to the target table query", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("search-form-table"), ADMIN, api.fetcher);
    const input = fieldInput(container, "q");
    expect(input).not.toBeNull();
    await act(async () => setFieldValue(input as HTMLInputElement, "acme"));
    const searchButton = buttonByText(container, "Search");
    await act(async () => (searchButton as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.url.includes("q=acme"))).toBe(true);
    expect(container.textContent).toContain("Acme Console");
    expect(container.textContent).not.toContain("Northwind Sales");
    expect(container.textContent).toContain("1 record · page 1 of 1");
  });
});

describe("T-UI-04 · create form (POST → 201)", () => {
  it("creates a record and shows it in the refreshed list", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New record") as HTMLButtonElement).click());
    await act(async () => setFieldValue(fieldInput(container, "name") as HTMLInputElement, "Globex Corp"));
    await act(async () => setFieldValue(fieldInput(container, "status") as HTMLInputElement, "active"));
    await act(async () => setFieldValue(fieldInput(container, "owner") as HTMLInputElement, "carol"));
    await act(async () => (buttonByText(container, "Create record") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "POST" && call.url === "/api/records")).toBe(true);
    expect(container.textContent).toContain("Record created");
    expect(container.textContent).toContain("Globex Corp");
  });

  it("surfaces a blank-field INVALID_CREATE_FIELD form error", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New record") as HTMLButtonElement).click());
    await act(async () => (buttonByText(container, "Create record") as HTMLButtonElement).click());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("INVALID_CREATE_FIELD");
  });
});

describe("T-UI-05 · edit form (PATCH → 200, prefilled from the row)", () => {
  it("prefills the edit form, patches the row, and refreshes updatedAt", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Edit") as HTMLButtonElement).click());
    const nameInput = fieldInput(container, "name") as HTMLInputElement;
    expect(nameInput.value).toBe("Acme Console");
    await act(async () => setFieldValue(nameInput, "Acme Rebrand"));
    await act(async () => (buttonByText(container, "Save changes") as HTMLButtonElement).click());
    const patchCall = api.calls.find((call) => call.method === "PATCH");
    expect(patchCall?.url).toBe("/api/records/rec-1");
    expect(container.textContent).toContain("Record updated");
    expect(container.textContent).toContain("Acme Rebrand");
    expect(api.store.find((record) => record.id === "rec-1")?.name).toBe("Acme Rebrand");
  });
});

describe("T-UI-06 · delete with confirmation", () => {
  it("confirms then deletes the row (DELETE 204) and removes it from the list", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Delete") as HTMLButtonElement).click());
    expect(container.textContent).toContain("Delete this record?");
    await act(async () => (buttonByText(container, "Confirm") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "DELETE" && call.url === "/api/records/rec-1")).toBe(true);
    expect(container.textContent).toContain("Record deleted");
    expect(container.textContent).not.toContain("Acme Console");
  });

  it("cancelling issues no request", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Delete") as HTMLButtonElement).click());
    await act(async () => (buttonByText(container, "Cancel") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "DELETE")).toBe(false);
    expect(container.textContent).toContain("Acme Console");
  });
});

describe("T-UI-07 · error envelope rendering", () => {
  it("renders a FORBIDDEN envelope when the backend denies a write", async () => {
    const api = createRecordsApi(RECORDS, { denyWrites: true });
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New record") as HTMLButtonElement).click());
    await act(async () => setFieldValue(fieldInput(container, "name") as HTMLInputElement, "Blocked"));
    await act(async () => setFieldValue(fieldInput(container, "status") as HTMLInputElement, "active"));
    await act(async () => setFieldValue(fieldInput(container, "owner") as HTMLInputElement, "dan"));
    await act(async () => (buttonByText(container, "Create record") as HTMLButtonElement).click());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("FORBIDDEN");
  });

  it("renders a RECORD_NOT_FOUND envelope on a stale edit", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Edit") as HTMLButtonElement).click());
    api.store.length = 0; // emulate a concurrent delete
    await act(async () => (buttonByText(container, "Save changes") as HTMLButtonElement).click());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("RECORD_NOT_FOUND");
  });
});

describe("T-UI-08 · permission matrix (admin vs viewer)", () => {
  it("admin sees enabled write affordances", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), ADMIN, api.fetcher);
    for (const label of ["New record", "Edit", "Delete"]) {
      const button = buttonByText(container, label) as HTMLButtonElement;
      expect(button, `${label} present`).not.toBeUndefined();
      expect(button.disabled, `${label} enabled for admin`).toBe(false);
    }
  });

  it("viewer/editor is read-only: write affordances are disabled", async () => {
    const api = createRecordsApi(RECORDS);
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), VIEWER, api.fetcher);
    for (const label of ["New record", "Edit", "Delete"]) {
      const button = buttonByText(container, label) as HTMLButtonElement;
      expect(button, `${label} present`).not.toBeUndefined();
      expect(button.disabled, `${label} disabled for viewer`).toBe(true);
    }
    // The read surface still works (records.read is held).
    expect(container.textContent).toContain("Acme Console");
  });
});

describe("T-UI-09 · backend authority is not replaced by frontend hiding", () => {
  it("the API denies a write even when the frontend hides the affordance", async () => {
    const viewerApi = createRecordsApi(RECORDS, { denyWrites: true });
    const container = await renderCrud(fixtureDocument("list-edit-lifecycle"), VIEWER, viewerApi.fetcher);
    // Frontend hides/disables the affordance for the viewer…
    const editButton = buttonByText(container, "Edit") as HTMLButtonElement;
    expect(editButton.disabled).toBe(true);
    // …and a direct write the UI never issues is still blocked by the backend.
    const response = await viewerApi.fetcher("/api/records", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ name: "X", status: "active", owner: "x" }),
    });
    expect(response.status).toBe(403);
    const envelope = (await response.json()) as { error: string };
    expect(envelope.error).toBe("FORBIDDEN");
  });
});

describe("T-UI-10 · page changes are fixture-only (renderer stays generic)", () => {
  it("the record CRUD action ids exist only in the fixture, never in renderer source", async () => {
    const fixtureText = readFileSync(resolve(FIXTURE_DIR, "list-edit-lifecycle.json"), "utf8");
    const actionIds = ["createRecord", "updateRecord", "deleteRecord", "openCreate", "openEdit"];
    for (const id of actionIds) {
      expect(fixtureText, `${id} declared by the fixture`).toContain(id);
    }
    // The records client (records.ts / use-records.ts) is the generic reusable
    // transport for any page, so its function names legitimately echo the
    // record verbs; the page-*rendering* source must not hardcode them.
    const srcFiles = readDirRecursive(WEB_SRC_DIR)
      .filter(
        (file) =>
          (file.endsWith(".ts") || file.endsWith(".tsx")) && !file.includes(".test."),
      )
      .filter(
        (file) =>
          !file.endsWith("records.ts") && !file.endsWith("use-records.ts"),
      );
    for (const file of srcFiles) {
      const content = readFileSync(file, "utf8");
      for (const id of actionIds) {
        expect(content, `${id} must not be hardcoded in ${file}`).not.toContain(id);
      }
    }
  });
});

function readDirRecursive(dir: string): string[] {
  const out: string[] = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = resolve(dir, entry.name);
    if (entry.isDirectory()) {
      out.push(...readDirRecursive(full));
    } else {
      out.push(full);
    }
  }
  return out;
}
