import { describe, expect, it, vi } from "vitest";

import {
  buildRecordsQuery,
  fetchRecords,
  parseRecordList,
  type RecordList,
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

  it("fails closed on a missing items array", () => {
    expect(() => parseRecordList({ total: 1, page: 1, pageSize: 10 })).toThrow();
  });

  it("fails closed on a non-string field", () => {
    expect(() =>
      parseRecordList({
        items: [{ id: 1, name: "x", status: "active", owner: "a", updatedAt: "t" }],
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
    const list: RecordList = await fetchRecords(
      fetcher as unknown as typeof fetch,
      "/api/records",
      { q: "acme" },
    );
    expect(fetcher).toHaveBeenCalledWith("/api/records?q=acme");
    expect(list.total).toBe(1);
  });

  it("throws on a non-OK response", async () => {
    const fetcher = vi.fn(async () => new Response("nope", { status: 500 }));
    await expect(
      fetchRecords(fetcher as unknown as typeof fetch, "/api/records", {}),
    ).rejects.toThrow("HTTP 500");
  });
});
