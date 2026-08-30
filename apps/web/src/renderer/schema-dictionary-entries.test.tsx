// @vitest-environment jsdom
//
// GOAL-015 · 数据字典内页（按类型过滤）+ 面包屑层级导航 — S3 验证证据。
//
// T-DE-01..T-DE-05 drive the REAL module-owned dictionary-entries page schema
// (apps/api/modules/datadictionary/schema/dictionary-entries.json)
// through the frozen §5 main render path, emulating the dictionary API:
//   - T-DE-01: v2.9 data.route-binding (ADR-0039) — table dataSource params
//     bind dictKey from $context.route.query.dictKey and the list request
//     carries ?dictKey=<type> (server-side filter).
//   - T-DE-02: missing route key is a tombstone — no dictKey param is sent (the
//     router already fails closed on deep links without the path param, F-007b).
//   - T-DE-03: create modal — the readOnly dictKey field is seeded by the Host
//     from the route query (ADR-0040: value source in modal scenarios) and the
//     POST body still carries the type key.
//   - T-DE-04: edit modal — dictKey is readOnly (row value) and dictTypeName is
//     displayed read-only; the PATCH body keeps the type key.
//   - T-DE-05: a readOnly field cannot be edited by the user.

import { readFileSync } from "node:fs";
import { resolve } from "node:path";
import { dirname } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { RenderPageDocument } from "@/renderer/render.types";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";
import type { ResourceItem } from "@/renderer/resource";

const ENTRIES_SCHEMA_PATH = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/modules/datadictionary/schema/dictionary-entries.json",
);

function entriesDocument(): unknown {
  return JSON.parse(readFileSync(ENTRIES_SCHEMA_PATH, "utf8"));
}

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

const TYPE_NAMES: Record<string, string> = {
  order_status: "Order status",
  other_status: "Other status",
};

function createEntriesApi(
  initial: ResourceItem[],
): { fetcher: typeof fetch; store: ResourceItem[]; calls: ApiCall[] } {
  const store: ResourceItem[] = initial.map((entry) => ({ ...entry }));
  const calls: ApiCall[] = [];
  let tick = 0;
  const now = (): string => {
    tick += 1;
    return "2026-08-14T10:00:00." + String(tick).padStart(3, "0") + "Z";
  };

  const fetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const raw = String(input);
    const url = new URL(raw, "http://test.local");
    const method = (init?.method ?? "GET").toUpperCase();
    calls.push({
      method,
      url: raw,
      body: init?.body !== undefined ? JSON.parse(String(init.body)) : undefined,
    });
    const match = /^\/api\/data-dictionary\/entries(?:\/([^/]+))?$/.exec(url.pathname);
    if (!match) {
      return json({ error: "NOT_FOUND", message: "no such route" }, 404);
    }
    const id = match[1];

    if (method === "GET" && id === undefined) {
      // Server-side dictKey filter (ExtraQuery) + q composition.
      let items = [...store];
      const dictKey = url.searchParams.get("dictKey");
      if (dictKey !== null && dictKey !== "") {
        items = items.filter((entry) => entry.dictKey === dictKey);
      }
      const q = url.searchParams.get("q");
      if (q !== null && q !== "") {
        items = items.filter((entry) =>
          String(entry.label ?? "").toLowerCase().includes(q.toLowerCase()),
        );
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
      const body = init?.body !== undefined ? JSON.parse(String(init.body)) : {};
      for (const field of ["dictKey", "entryKey", "label"]) {
        const value = typeof body[field] === "string" ? String(body[field]).trim() : "";
        if (value === "") {
          return json({ error: "INVALID_CREATE_FIELD", message: field + " is required" }, 400);
        }
      }
      const entry: ResourceItem = {
        id: "entry-new-" + (store.length + 1),
        dictKey: String(body.dictKey).trim(),
        dictTypeName: TYPE_NAMES[String(body.dictKey).trim()] ?? String(body.dictKey).trim(),
        entryKey: String(body.entryKey).trim(),
        label: String(body.label).trim(),
        enabled: body.enabled === true,
        sort: typeof body.sort === "number" ? body.sort : 0,
        remark: typeof body.remark === "string" ? body.remark : "",
        updatedAt: now(),
      };
      store.unshift(entry);
      return json(entry, 201);
    }
    if (id !== undefined && (method === "PATCH" || method === "DELETE")) {
      const index = store.findIndex((entry) => entry.id === id);
      if (index < 0) {
        return json({ error: "ENTRY_NOT_FOUND", message: "no entry with that id" }, 404);
      }
      if (method === "DELETE") {
        store.splice(index, 1);
        return new Response(null, { status: 204 });
      }
      const body = init?.body !== undefined ? JSON.parse(String(init.body)) : {};
      const updated = {
        ...store[index],
        label: body.label === undefined ? store[index]?.label : String(body.label).trim(),
        enabled: body.enabled === undefined ? store[index]?.enabled : body.enabled === true,
        sort: body.sort === undefined ? store[index]?.sort : body.sort,
        remark: body.remark === undefined ? store[index]?.remark : String(body.remark),
        updatedAt: now(),
      };
      store[index] = updated;
      return json(updated);
    }
    return json({ error: "NOT_FOUND", message: "no such route" }, 404);
  }) as typeof fetch;

  return { fetcher, store, calls };
}

// F-007(b): the entries page route is /dictionary-entries/{dictKey} — the
// render context carries the resolved path params + query (App injects
// route: {params, query}); deep links without the param fail closed at the
// router level (HOST_ROUTE_NOT_FOUND), so the table binding always resolves.
function dictContext(
  params: Record<string, string>,
  query: Record<string, string> = {},
): Record<string, unknown> {
  return {
    user: { id: "u1", roles: ["admin"], permissions: ["dictionary.read", "dictionary.write"] },
    route: { params, query },
  };
}

const ENTRIES: ResourceItem[] = [
  {
    id: "ent-1",
    dictKey: "order_status",
    dictTypeName: "Order status",
    entryKey: "PAID",
    label: "Paid",
    enabled: true,
    sort: 0,
    updatedAt: "2026-08-14T08:00:00.000Z",
  },
  {
    id: "ent-2",
    dictKey: "order_status",
    dictTypeName: "Order status",
    entryKey: "CANCELLED",
    label: "Cancelled",
    enabled: true,
    sort: 1,
    updatedAt: "2026-08-14T08:05:00.000Z",
  },
  {
    id: "ent-3",
    dictKey: "other_status",
    dictTypeName: "Other status",
    entryKey: "DONE",
    label: "Done",
    enabled: true,
    sort: 0,
    updatedAt: "2026-08-14T08:10:00.000Z",
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
  return container.querySelector<HTMLInputElement>("#field-" + id);
}

describe("T-DE · GOAL-015 dictionary entries inner page (v2.9)", () => {
  it("T-DE-01: table dataSource binds dictKey from the route path param (ADR-0039)", async () => {
    const api = createEntriesApi(ENTRIES);
    const container = await renderCrud(
      entriesDocument(),
      dictContext({ dictKey: "order_status" }),
      api.fetcher,
    );

    expect(api.calls.some((call) => call.method === "GET" && call.url.includes("dictKey=order_status"))).toBe(true);
    expect(container.textContent).toContain("Paid");
    expect(container.textContent).toContain("Cancelled");
    expect(container.textContent).not.toContain("Done");
    expect(container.textContent).toContain("Order status");
  });

  it("T-DE-02: missing route key is a tombstone — no dictKey param is sent", async () => {
    const api = createEntriesApi(ENTRIES);
    const container = await renderCrud(entriesDocument(), dictContext({}), api.fetcher);

    const firstList = api.calls.find((call) => call.method === "GET");
    expect(firstList?.url).toBe("/api/data-dictionary/entries");
    expect(firstList?.url).not.toContain("dictKey");
    expect(container.textContent).toContain("Paid");
    expect(container.textContent).toContain("Done");
  });

  it("T-DE-03: create modal seeds the readOnly dictKey from the route and POSTs the type key", async () => {
    const api = createEntriesApi(ENTRIES);
    const container = await renderCrud(
      entriesDocument(),
      dictContext({ dictKey: "order_status" }, { dictTypeName: "Order status" }),
      api.fetcher,
    );

    await act(async () => (buttonByText(container, "New entry") as HTMLButtonElement).click());
    const dictKeyInput = fieldInput(container, "dictKey") as HTMLInputElement;
    expect(dictKeyInput).not.toBeNull();
    expect(dictKeyInput.readOnly).toBe(true);
    expect(dictKeyInput.value).toBe("order_status");
    // F-006: the create form also shows the type NAME read-only (seeded from
    // the navigate query, which now carries dictTypeName=$row.name).
    const typeNameInput = fieldInput(container, "dictTypeName") as HTMLInputElement;
    expect(typeNameInput).not.toBeNull();
    expect(typeNameInput.readOnly).toBe(true);
    expect(typeNameInput.value).toBe("Order status");

    await act(async () => setFieldValue(fieldInput(container, "entryKey") as HTMLInputElement, "REFUNDED"));
    await act(async () => setFieldValue(fieldInput(container, "label") as HTMLInputElement, "Refunded"));
    await act(async () => (buttonByText(container, "Create entry") as HTMLButtonElement).click());

    const post = api.calls.find((call) => call.method === "POST");
    expect(post?.url).toBe("/api/data-dictionary/entries");
    expect(post?.body).toMatchObject({
      dictKey: "order_status",
      entryKey: "REFUNDED",
      label: "Refunded",
    });
    expect(container.textContent).toContain("Item created");
  });

  it("T-DE-04: edit form shows dictKey + dictTypeName read-only and PATCHes the type key", async () => {
    const api = createEntriesApi(ENTRIES);
    const container = await renderCrud(
      entriesDocument(),
      dictContext({ dictKey: "order_status" }),
      api.fetcher,
    );

    await act(async () => (buttonByText(container, "Edit") as HTMLButtonElement).click());
    const dictKeyInput = fieldInput(container, "dictKey") as HTMLInputElement;
    expect(dictKeyInput).not.toBeNull();
    expect(dictKeyInput.readOnly).toBe(true);
    expect(dictKeyInput.value).toBe("order_status");
    const typeNameInput = fieldInput(container, "dictTypeName") as HTMLInputElement;
    expect(typeNameInput).not.toBeNull();
    expect(typeNameInput.readOnly).toBe(true);
    expect(typeNameInput.value).toBe("Order status");

    const labelInput = fieldInput(container, "label") as HTMLInputElement;
    await act(async () => setFieldValue(labelInput, "Paid (updated)"));
    await act(async () => (buttonByText(container, "Save changes") as HTMLButtonElement).click());

    const patch = api.calls.find((call) => call.method === "PATCH");
    expect(patch?.url).toBe("/api/data-dictionary/entries/ent-1");
    expect(patch?.body).toMatchObject({ dictKey: "order_status", label: "Paid (updated)" });
    expect(container.textContent).toContain("Item updated");
  });

  it("T-DE-05: a readOnly field cannot be edited by the user", async () => {
    const api = createEntriesApi(ENTRIES);
    const container = await renderCrud(
      entriesDocument(),
      dictContext({ dictKey: "order_status" }),
      api.fetcher,
    );

    await act(async () => (buttonByText(container, "New entry") as HTMLButtonElement).click());
    const dictKeyInput = fieldInput(container, "dictKey") as HTMLInputElement;
    await act(async () => setFieldValue(dictKeyInput, "hacked"));
    expect(dictKeyInput.value).toBe("order_status");
  });
});
