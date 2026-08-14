import { describe, expect, it, vi } from "vitest";

import {
  buildResourceQuery,
  deleteResource,
  fetchResourceList,
  isValidDataSource,
  parseResourceList,
  resolveDataParamsQuery,
  updateResource,
  type ResourceList,
} from "@/renderer/resource";

describe("buildResourceQuery (query-serialization)", () => {
  it("omits empty query", () => {
    expect(buildResourceQuery({})).toBe("");
  });

  it("serializes search and sort/order", () => {
    expect(buildResourceQuery({ q: "alice", sort: "name", order: "asc" })).toBe(
      "q=alice&sort=name&order=asc",
    );
  });

  it("omits default page and pageSize", () => {
    expect(buildResourceQuery({ page: 1, pageSize: 10 })).toBe("");
  });

  it("serializes non-default pagination", () => {
    expect(buildResourceQuery({ page: 2, pageSize: 25 })).toBe("page=2&pageSize=25");
  });

  it("trims blank search", () => {
    expect(buildResourceQuery({ q: "   " })).toBe("");
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

describe("parseResourceList (response-mapping)", () => {
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
    const list = parseResourceList(value);
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
    const list = parseResourceList(value);
    expect(list.items[0]).toEqual({ sku: "S-1", title: "Widget", price: 19 });
  });

  it("fails closed on a missing items array", () => {
    expect(() => parseResourceList({ total: 1, page: 1, pageSize: 10 })).toThrow();
  });

  it("fails closed on a non-object item", () => {
    expect(() =>
      parseResourceList({
        items: [null],
        total: 1,
        page: 1,
        pageSize: 10,
      }),
    ).toThrow();
  });

  it("fails closed on a non-object payload", () => {
    expect(() => parseResourceList(null)).toThrow();
    expect(() => parseResourceList([])).toThrow();
  });
});

describe("resolveDataParamsQuery (v2.9 ADR-0039 route bindings)", () => {
  it("passes literal scalars through and drops nulls", () => {
    expect(resolveDataParamsQuery({ status: "paid", flag: true, n: 3, gone: null }, { query: {}, params: {} })).toBe("status=paid&flag=true&n=3");
  });

  it("resolves whole $context.route.query.* bindings", () => {
    expect(
      resolveDataParamsQuery(
        { dictKey: "$context.route.query.dictKey", status: "paid" },
        { query: { dictKey: "ORDER_STATUS" }, params: {} },
      ),
    ).toBe("dictKey=ORDER_STATUS&status=paid");
  });

  it("resolves $context.route.params.* bindings", () => {
    expect(
      resolveDataParamsQuery(
        { type: "$context.route.params.type" },
        { query: {}, params: { type: "admin" } },
      ),
    ).toBe("type=admin");
  });

  it("tombstones a missing route key (ADR-0010)", () => {
    expect(
      resolveDataParamsQuery(
        { dictKey: "$context.route.query.dictKey" },
        { query: {}, params: {} },
      ),
    ).toBe("");
  });

  it("tombstones when no route snapshot is provided", () => {
    expect(
      resolveDataParamsQuery(
        { dictKey: "$context.route.query.dictKey" },
        { query: undefined, params: undefined },
      ),
    ).toBe("");
  });

  it("drops unknown $context.route.* shapes fail-closed", () => {
    expect(
      resolveDataParamsQuery(
        { k: "$context.route.session.id" },
        { query: { session: "x" }, params: {} },
      ),
    ).toBe("");
  });
});

describe("fetchResourceList (request-construction)", () => {
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
    const list: ResourceList = await fetchResourceList(
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
      fetchResourceList(fetcher as unknown as typeof fetch, "/api/users", {}),
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
        fetchResourceList(fetcher as unknown as typeof fetch, bad, {}),
      ).rejects.toThrow(/invalid dataSource/);
    }
    expect(fetcher).not.toHaveBeenCalled();
  });
});

describe("updateResource (PATCH)", () => {
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
    const updated = await updateResource(
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
      updateResource(fetcher as unknown as typeof fetch, "/api/users", "usr-3", { name: "" }),
    ).rejects.toThrow("HTTP 400");
  });
});

describe("deleteResource (DELETE)", () => {
  it("sends a DELETE and resolves on 204", async () => {
    const fetcher = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      expect(String(input)).toBe("/api/users/usr-3");
      expect(init?.method).toBe("DELETE");
      return new Response(null, { status: 204 });
    });
    await expect(
      deleteResource(fetcher as unknown as typeof fetch, "/api/users", "usr-3"),
    ).resolves.toBeUndefined();
  });

  it("throws on a non-OK response", async () => {
    const fetcher = vi.fn(async () => new Response("bad", { status: 404 }));
    await expect(
      deleteResource(fetcher as unknown as typeof fetch, "/api/users", "usr-999"),
    ).rejects.toThrow("HTTP 404");
  });
});

describe("fetchResourceList with the system-monitoring status envelope (S-03 · GOAL-009 A-003 F-001)", () => {
  // The status endpoint serves a single-row list envelope so statCard
  // dataSource loading (fetchResourceList) can bind every valueField.
  it("parses the status row and reads every statCard valueField", async () => {
    const statusRow = {
      status: "ok",
      ready: true,
      version: "dev",
      commit: "abc",
      uptimeSeconds: 42,
      moduleCount: 12,
      modules: ["core.auth-session", "admin.users"],
      dbSizeBytes: 1024,
    };
    const fetcher = vi.fn(async () =>
      new Response(
        JSON.stringify({ items: [statusRow], total: 1, page: 1, pageSize: 1 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    );
    const list = await fetchResourceList(fetcher as typeof fetch, "/api/system-monitoring/status", {});
    expect(list.total).toBe(1);
    expect(list.items).toHaveLength(1);
    const row = list.items[0] as Record<string, unknown>;
    for (const field of ["status", "ready", "version", "commit", "uptimeSeconds", "moduleCount", "modules", "dbSizeBytes"]) {
      expect(row[field]).toEqual(statusRow[field as keyof typeof statusRow]);
    }
  });

  it("rejects a flat (non-envelope) status body fail-closed", async () => {
    const fetcher = vi.fn(async () =>
      new Response(
        JSON.stringify({ status: "ok", ready: true, uptimeSeconds: 1 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      ),
    );
    await expect(
      fetchResourceList(fetcher as typeof fetch, "/api/system-monitoring/status", {}),
    ).rejects.toThrow();
  });
});

describe("fetchResourceList extraQuery (v2.9 ADR-0039)", () => {
  it("merges extraQuery params with the standard query", async () => {
    const seen: string[] = [];
    const fetcher = (async (input: RequestInfo | URL) => {
      seen.push(String(input));
      return new Response(
        JSON.stringify({ items: [], total: 0, page: 1, pageSize: 10 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }) as typeof fetch;
    const list = await fetchResourceList(fetcher, "/api/entries", { q: "paid" }, "dictKey=order_status");
    expect(list.total).toBe(0);
    expect(seen[0]).toContain("q=paid");
    expect(seen[0]).toContain("dictKey=order_status");
  });

  it("keeps baseURL bare — F-001 rejects a ?-carrying dataSource even with extraQuery", async () => {
    const fetcher = (async () =>
      new Response(
        JSON.stringify({ items: [], total: 0, page: 1, pageSize: 10 }),
        { status: 200 },
      )
    ) as typeof fetch;
    await expect(
      fetchResourceList(fetcher, "/api/entries?x=1", {}, "dictKey=order_status"),
    ).rejects.toThrow(/invalid dataSource/);
  });
});
