import { describe, expect, it, vi } from "vitest";

import {
  buildRecordsQuery,
  deleteRecord,
  fetchRecords,
  isValidDataSource,
  parseRecordList,
  updateRecord,
  type ResourceList,
} from "@/renderer/records";

describe("buildRecordsQuery (query-serialization)", () => {
  it("omits empty query", () => {
    expect(buildRecordsQuery({})).toBe("");
  });

  it("serializes search and sort/order", () => {
    expect(buildRecordsQuery({ q: "alice", sort: "name", order: "asc" })).toBe(
      "q=alice&sort=name&order=asc",
    );
  });

  it("omits default page and pageSize", () => {
    expect(buildRecordsQuery({ page: 1, pageSize: 10 })).toBe("");
  });

  it("serializes non-default pagination", () => {
    expect(buildRecordsQuery({ page: 2, pageSize: 25 })).toBe("page=2&pageSize=25");
  });

  it("trims blank search", () => {
    expect(buildRecordsQuery({ q: "   " })).toBe("");
  });
});

describe("isValidDataSource (F-001 · I-010-001 v0.2.0 §2)", () => {
  it("accepts single-slash same-origin paths", () => {
    expect(isValidDataSource("/api/users")).toBe(true);
    expect(isValidDataSource("/api/catalog")).toBe(true);
    expect(isValidDataSource("/")).toBe(true);
  });

  it("rejects protocol-relative and absolute URLs", () => {
    expect(isValidDataSource("//evil.example/api/users")).toBe(false);
    expect(isValidDataSource("http://evil.example/api/users")).toBe(false);
    expect(isValidDataSource("https://evil.example/api/users")).toBe(false);
    expect(isValidDataSource("javascript:alert(1)")).toBe(false);
  });

  it("rejects relative (non-rooted) paths", () => {
    expect(isValidDataSource("api/users")).toBe(false);
    expect(isValidDataSource("records")).toBe(false);
  });

  it("rejects whitespace, backslash, query and fragment", () => {
    expect(isValidDataSource("/api/rec ords")).toBe(false);
    expect(isValidDataSource("/api\\records")).toBe(false);
    expect(isValidDataSource("/api/users?q=x")).toBe(false);
    expect(isValidDataSource("/api/users#frag")).toBe(false);
  });

  it("rejects empty and non-string input", () => {
    expect(isValidDataSource("")).toBe(false);
  });
});

describe("parseRecordList (response-mapping)", () => {
  it("maps an envelope", () => {
    const value = {
      items: [
        {
          id: "rec-1",
          name: "Acme Console",
          status: "active",
          owner: "alice",
          updatedAt: "2026-07-31T00:00:00Z",
        },
      ],
      total: 1,
      page: 1,
      pageSize: 10,
    };
    const list = parseRecordList(value);
    expect(list.items[0].id).toBe("rec-1");
    expect(list.total).toBe(1);
  });

  it("accepts arbitrary object rows (no five-field whitelist)", () => {
    const value = {
      items: [{ sku: "S-1", title: "Widget", price: 19 }],
      total: 1,
      page: 1,
      pageSize: 10,
    };
    const list = parseRecordList(value);
    expect(list.items[0]).toEqual({ sku: "S-1", title: "Widget", price: 19 });
  });

  it("fails closed on a missing items array", () => {
    expect(() => parseRecordList({ total: 1, page: 1, pageSize: 10 })).toThrow();
  });

  it("fails closed on a non-object item", () => {
    expect(() =>
      parseRecordList({
        items: [null],
        total: 1,
        page: 1,
        pageSize: 10,
      }),
    ).toThrow();
  });

  it("fails closed on a non-object payload", () => {
    expect(() => parseRecordList(null)).toThrow();
    expect(() => parseRecordList([])).toThrow();
  });
});

describe("fetchRecords (request-construction)", () => {
  it("builds the URL and maps the response", async () => {
    const fetcher = vi.fn(async (_input: RequestInfo | URL) => {
      return new Response(
        JSON.stringify({
          items: [
            {
              id: "rec-1",
              name: "Acme Console",
              status: "active",
              owner: "alice",
              updatedAt: "2026-07-31T00:00:00Z",
            },
          ],
          total: 1,
          page: 1,
          pageSize: 10,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    });
    const list: ResourceList = await fetchRecords(
      fetcher as unknown as typeof fetch,
      "/api/users",
      { q: "acme" },
    );
    expect(fetcher).toHaveBeenCalledWith("/api/users?q=acme");
    expect(list.total).toBe(1);
  });

  it("throws on a non-OK response", async () => {
    const fetcher = vi.fn(async () => new Response("nope", { status: 500 }));
    await expect(
      fetchRecords(fetcher as unknown as typeof fetch, "/api/users", {}),
    ).rejects.toThrow("HTTP 500");
  });

  it("rejects an invalid dataSource before touching the fetcher (F-001)", async () => {
    const fetcher = vi.fn(async () => new Response("{}", { status: 200 }));
    for (const bad of [
      "//evil.example/api/users",
      "http://evil.example/api/users",
      "api/users",
      "/api/rec ords",
      "/api/users?q=x",
      "",
    ]) {
      await expect(
        fetchRecords(fetcher as unknown as typeof fetch, bad, {}),
      ).rejects.toThrow(/invalid dataSource/);
    }
    expect(fetcher).not.toHaveBeenCalled();
  });
});

describe("updateRecord (PATCH)", () => {
  it("sends a PATCH and maps the response", async () => {
    const fetcher = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      expect(String(input)).toBe("/api/users/usr-3");
      expect(init?.method).toBe("PATCH");
      return new Response(
        JSON.stringify({
          id: "usr-3",
          name: "Hooli Rebrand",
          status: "archived",
          owner: "carol",
          updatedAt: "2026-07-31T00:00:00Z",
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    });
    const updated = await updateRecord(
      fetcher as unknown as typeof fetch,
      "/api/users",
      "usr-3",
      { name: "Hooli Rebrand", status: "archived" },
    );
    expect(updated.name).toBe("Hooli Rebrand");
    expect(updated.status).toBe("archived");
  });

  it("throws on a non-OK response", async () => {
    const fetcher = vi.fn(async () => new Response("bad", { status: 400 }));
    await expect(
      updateRecord(fetcher as unknown as typeof fetch, "/api/users", "usr-3", { name: "" }),
    ).rejects.toThrow("HTTP 400");
  });
});

describe("deleteRecord (DELETE)", () => {
  it("sends a DELETE and resolves on 204", async () => {
    const fetcher = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      expect(String(input)).toBe("/api/users/usr-3");
      expect(init?.method).toBe("DELETE");
      return new Response(null, { status: 204 });
    });
    await expect(
      deleteRecord(fetcher as unknown as typeof fetch, "/api/users", "usr-3"),
    ).resolves.toBeUndefined();
  });

  it("throws on a non-OK response", async () => {
    const fetcher = vi.fn(async () => new Response("bad", { status: 404 }));
    await expect(
      deleteRecord(fetcher as unknown as typeof fetch, "/api/users", "usr-999"),
    ).rejects.toThrow("HTTP 404");
  });
});
