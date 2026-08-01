import { describe, expect, it, vi } from "vitest";

import type { PageEntry } from "@/protocol/app-manifest";
import { validatePageDocument } from "@/protocol/conformance/runtime-schema-validate";
import {
  loadPageDocument,
  PageSchemaError,
  type PageSchemaErrorCode,
} from "@/protocol/load-page";

const BASE = "https://example.test";

// Structurally valid page document (mirrors the seeded Go fixture
// apps/api/internal/handler/fixtures/schema/overview.json so the loader test
// exercises the same shape the endpoint serves).
const VALID_DOCUMENT = {
  meta: {
    pageId: "overview",
    title: "Overview",
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation"],
  },
  body: {
    type: "section",
    props: { title: "Overview" },
    children: [{ type: "text", props: { text: "Welcome" } }],
  },
};

// Structurally invalid: meta is missing the required protocolVersion.
const INVALID_DOCUMENT = {
  meta: { pageId: "overview", title: "Overview" },
  body: { type: "section" },
};

const OVERVIEW_PAGE: PageEntry = {
  pageId: "overview",
  title: "Overview",
  schemaUrl: "/api/schema/overview",
  route: "/overview",
};

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function expectErrorCode(
  promise: Promise<unknown>,
  code: PageSchemaErrorCode,
): Promise<PageSchemaError> {
  return promise.then(
    () => {
      throw new Error(`expected loadPageDocument to reject with ${code}`);
    },
    (error: unknown) => {
      expect(error).toBeInstanceOf(PageSchemaError);
      const schemaError = error as PageSchemaError;
      expect(schemaError.code).toBe(code);
      return schemaError;
    },
  );
}

describe("validatePageDocument", () => {
  it("accepts a structurally valid page document", () => {
    expect(validatePageDocument(VALID_DOCUMENT).ok).toBe(true);
  });

  it("fails closed on a document missing required meta fields", () => {
    const result = validatePageDocument(INVALID_DOCUMENT);
    expect(result.ok).toBe(false);
    expect(result.errors.length).toBeGreaterThan(0);
  });
});

describe("loadPageDocument", () => {
  it("loads and returns a valid page document", async () => {
    const fetcher = vi.fn(async () => jsonResponse(VALID_DOCUMENT));
    const loaded = await loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher });
    expect(loaded).toEqual(VALID_DOCUMENT);
    expect(fetcher).toHaveBeenCalledWith("https://example.test/api/schema/overview");
  });

  it("expands route-parameter placeholders in schemaUrl", async () => {
    const detailPage: PageEntry = {
      pageId: "order-detail",
      title: "Order detail",
      schemaUrl: "/api/schema/orders/{id}",
      route: "/orders/{id}",
    };
    const fetcher = vi.fn(async () =>
      jsonResponse({
        meta: { pageId: "order-detail", title: "Order detail", protocolVersion: "2.7" },
        body: { type: "section" },
      }),
    );
    await loadPageDocument(detailPage, { id: "o-42" }, { baseURL: BASE, fetcher });
    expect(fetcher).toHaveBeenCalledWith("https://example.test/api/schema/orders/o-42");
  });

  it("maps a 404 response to PAGE_NOT_FOUND", async () => {
    const fetcher = vi.fn(async () => jsonResponse({ error: "SCHEMA_NOT_FOUND" }, 404));
    const error = await expectErrorCode(
      loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher }),
      "PAGE_NOT_FOUND",
    );
    expect(error.url).toBe("https://example.test/api/schema/overview");
  });

  it("maps a 5xx response to PAGE_LOAD_FAILED", async () => {
    const fetcher = vi.fn(async () => new Response("boom", { status: 500 }));
    await expectErrorCode(
      loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher }),
      "PAGE_LOAD_FAILED",
    );
  });

  it("maps a network failure to PAGE_LOAD_FAILED", async () => {
    const fetcher = vi.fn(async () => {
      throw new TypeError("fetch failed");
    });
    await expectErrorCode(
      loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher }),
      "PAGE_LOAD_FAILED",
    );
  });

  it("maps a non-JSON body to PAGE_PARSE_FAILED", async () => {
    const fetcher = vi.fn(async () => new Response("<html>oops</html>", { status: 200 }));
    await expectErrorCode(
      loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher }),
      "PAGE_PARSE_FAILED",
    );
  });

  it("fails closed on a structurally invalid document with issues", async () => {
    const fetcher = vi.fn(async () => jsonResponse(INVALID_DOCUMENT));
    const error = await expectErrorCode(
      loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher }),
      "PAGE_SCHEMA_INVALID",
    );
    expect(error.issues?.length ?? 0).toBeGreaterThan(0);
  });

  it("rejects a document whose meta.pageId differs from the manifest pageId", async () => {
    const fetcher = vi.fn(async () =>
      jsonResponse({
        meta: { pageId: "catalog", title: "Overview", protocolVersion: "2.7" },
        body: { type: "section" },
      }),
    );
    const error = await expectErrorCode(
      loadPageDocument(OVERVIEW_PAGE, {}, { baseURL: BASE, fetcher }),
      "PAGE_ID_MISMATCH",
    );
    expect(error.message).toContain("catalog");
  });
});
