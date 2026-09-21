import { describe, expect, it } from "vitest";

import {
  createSavedViewRecord,
  normalizeSavedViewState,
  readSavedViews,
  savedViewStorageKey,
  updateSavedViewRecord,
  writeSavedViews,
  type SavedViewRecord,
  type SavedViewStorageLike,
} from "@/renderer/saved-views";

const config = {
  columnFields: ["name", "status", "owner"],
  sortableFields: ["name", "updatedAt"],
  filterFields: ["status", "owner"],
} as const;

function memoryStorage(initial?: string): SavedViewStorageLike & { value: string | null } {
  let value = initial ?? null;
  return {
    get value() {
      return value;
    },
    getItem: () => value,
    setItem: (_key, next) => {
      value = next;
    },
  };
}

function view(id = "view-1"): SavedViewRecord {
  const result = createSavedViewRecord(
    "Operations",
    { q: "alice", filters: { status: "active" }, sort: "name", order: "asc", pageSize: 20 },
    ["name", "status"],
    config,
    "2026-09-17T00:00:00.000Z",
    id,
  );
  if (!result.ok) throw new Error(result.message);
  return result.view;
}

describe("Saved View storage contract", () => {
  it("encodes the user/page/table storage namespace", () => {
    expect(savedViewStorageKey("user/a", "admin users", "table.1")).toBe(
      "schema-ui.saved-views.v1.user%2Fa.admin%20users.table%2E1",
    );
    expect(savedViewStorageKey("a.b", "c", "d")).not.toBe(savedViewStorageKey("a", "b", "c.d"));
  });

  it("projects only the query and column allowlist", () => {
    const result = normalizeSavedViewState(
      {
        q: "  alice ",
        filters: { status: " active " },
        sort: "name",
        order: "desc",
        page: 99,
        pageSize: 50,
      },
      ["name", "owner"],
      config,
    );
    expect(result).toEqual({
      ok: false,
      code: "SAVED_VIEW_STATE_INVALID",
      message: "Saved View query contains an unknown field.",
    });
    const valid = normalizeSavedViewState(
      { q: "  alice ", filters: { status: " active " }, sort: "name", order: "desc", pageSize: 50 },
      ["name", "owner"],
      config,
    );
    expect(valid).toEqual({
      ok: true,
      state: {
        query: { q: "alice", filters: { status: "active" }, sort: "name", order: "desc", pageSize: 50 },
        visibleColumns: ["name", "owner"],
      },
    });
  });

  it("round-trips views and the active view id", () => {
    const storage = memoryStorage();
    const first = view("view-1");
    const second = view("view-2");
    expect(writeSavedViews(storage, "key", [first, second], config, "view-2")).toEqual({ ok: true });
    expect(readSavedViews(storage, "key", config)).toEqual({
      ok: true,
      views: [first, second],
      activeViewId: "view-2",
      droppedCount: 0,
    });
  });

  it("drops malformed and duplicate records without applying them", () => {
    const storage = memoryStorage(
      JSON.stringify({
        version: 1,
        views: [
          view("view-1"),
          view("view-1"),
          { ...view("view-bad"), visibleColumns: ["not-in-schema"] },
        ],
        activeViewId: "view-bad",
      }),
    );
    expect(readSavedViews(storage, "key", config)).toEqual({
      ok: true,
      views: [view("view-1")],
      droppedCount: 2,
    });
  });

  it("fails closed for malformed documents and storage failures", () => {
    expect(readSavedViews(memoryStorage("not-json"), "key", config)).toMatchObject({
      ok: false,
      code: "SAVED_VIEW_STORAGE_INVALID",
    });
    const broken: SavedViewStorageLike = {
      getItem: () => {
        throw new Error("denied");
      },
      setItem: () => {
        throw new Error("denied");
      },
    };
    expect(readSavedViews(broken, "key", config)).toMatchObject({
      ok: false,
      code: "SAVED_VIEW_STORAGE_UNAVAILABLE",
    });
    expect(writeSavedViews(broken, "key", [view()], config)).toMatchObject({
      ok: false,
      code: "SAVED_VIEW_STORAGE_WRITE_FAILED",
    });
  });

  it("rejects unknown document fields and an active pointer outside the stored set", () => {
    const stored = memoryStorage(
      JSON.stringify({ version: 1, views: [view()], activeViewId: "missing", debug: true }),
    );
    expect(readSavedViews(stored, "key", config)).toMatchObject({
      ok: false,
      code: "SAVED_VIEW_STORAGE_INVALID",
    });
    expect(writeSavedViews(memoryStorage(), "key", [view()], config, "missing")).toMatchObject({
      ok: false,
      code: "SAVED_VIEW_STATE_INVALID",
    });
  });

  it("updates a view without changing its identity or creation timestamp", () => {
    const existing = view();
    const result = updateSavedViewRecord(
      existing,
      { q: "bob" },
      ["name", "owner"],
      config,
      "2026-09-18T00:00:00.000Z",
    );
    expect(result).toEqual({
      ok: true,
      view: {
        ...existing,
        query: { q: "bob" },
        visibleColumns: ["name", "owner"],
        updatedAt: "2026-09-18T00:00:00.000Z",
      },
    });
  });
});
