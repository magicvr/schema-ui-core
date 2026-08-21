/**
 * Runtime page-schema loader for the schema-driven render path (R1 · GOAL-002).
 *
 * Resolves a manifest `PageEntry`'s schemaUrl (with route-parameter expansion),
 * fetches the page document, forces structural validation (D-VAL) against the
 * pinned page/node schemas, and returns the parsed document or throws a unified
 * `PageSchemaError`. The fetcher is injectable so tests exercise network,
 * parse, and validation failure paths without a server. This module does NOT
 * switch the app's default render branch (that is GOAL-003).
 */

import type { PageEntry } from "@/protocol/app-manifest";
import { resolveSchemaUrl } from "@/protocol/app-manifest";
import { validatePageDocument } from "@/protocol/conformance/runtime-schema-validate";
import { withTimeout } from "@/lib/fetch-timeout";

export type PageSchemaErrorCode =
  | "PAGE_LOAD_FAILED"
  | "PAGE_NOT_FOUND"
  | "PAGE_PARSE_FAILED"
  | "PAGE_SCHEMA_INVALID"
  | "PAGE_ID_MISMATCH";

export interface PageSchemaValidationIssue {
  path: string;
  message: string;
  keyword?: string;
}

/** Unified, observable error for the schema loading + validation pipeline. */
export class PageSchemaError extends Error {
  readonly code: PageSchemaErrorCode;
  readonly url: string;
  readonly issues?: PageSchemaValidationIssue[];

  constructor(
    code: PageSchemaErrorCode,
    url: string,
    message: string,
    issues?: PageSchemaValidationIssue[],
  ) {
    super(message);
    this.name = "PageSchemaError";
    this.code = code;
    this.url = url;
    this.issues = issues;
  }
}

export interface LoadPageOptions {
  /** Origin/base used to resolve relative schemaUrl values. Defaults to `location.origin`. */
  baseURL?: string;
  /** Injectable fetch; defaults to `globalThis.fetch`. */
  fetcher?: typeof fetch;
}

function defaultBaseURL(): string {
  return typeof globalThis.location !== "undefined" ? globalThis.location.origin : "";
}

/**
 * Load and validate a page document for the given manifest page and route
 * params. Resolves the schemaUrl (expanding `{param}` placeholders), fetches
 * it, enforces structural validation, and verifies the document's `meta.pageId`
 * matches the manifest page. Returns the parsed page document on success.
 */
export async function loadPageDocument(
  page: PageEntry,
  params: Record<string, string>,
  options: LoadPageOptions = {},
): Promise<unknown> {
  const baseURL = options.baseURL ?? defaultBaseURL();
  // W10 F-002: the default transport is timeout-bounded; an injected fetcher
  // (tests) is honored as-is so failure paths stay deterministic.
  const fetcher = options.fetcher ?? withTimeout();
  const url = resolveSchemaUrl(baseURL, page.schemaUrl, params);

  if (typeof fetcher !== "function") {
    throw new PageSchemaError("PAGE_LOAD_FAILED", url, "Fetch is unavailable.");
  }

  let response: Response;
  try {
    response = await fetcher(url);
  } catch (error) {
    throw new PageSchemaError(
      "PAGE_LOAD_FAILED",
      url,
      error instanceof Error ? error.message : "Network request failed.",
    );
  }

  if (!response.ok) {
    throw new PageSchemaError(
      response.status === 404 ? "PAGE_NOT_FOUND" : "PAGE_LOAD_FAILED",
      url,
      `Page document request failed with HTTP ${response.status}.`,
    );
  }

  let document: unknown;
  try {
    document = await response.json();
  } catch {
    throw new PageSchemaError("PAGE_PARSE_FAILED", url, "Response body is not valid JSON.");
  }

  const validation = validatePageDocument(document);
  if (!validation.ok) {
    throw new PageSchemaError(
      "PAGE_SCHEMA_INVALID",
      url,
      "Page document failed structural validation.",
      validation.errors,
    );
  }

  const meta = (document as { meta?: { pageId?: unknown } } | null)?.meta;
  if (typeof meta?.pageId === "string" && meta.pageId !== page.pageId) {
    throw new PageSchemaError(
      "PAGE_ID_MISMATCH",
      url,
      `Page document meta.pageId (${meta.pageId}) does not match manifest pageId (${page.pageId}).`,
    );
  }

  return document;
}
