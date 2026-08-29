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
export type PageSchemaErrorCode = "PAGE_LOAD_FAILED" | "PAGE_NOT_FOUND" | "PAGE_PARSE_FAILED" | "PAGE_SCHEMA_INVALID" | "PAGE_ID_MISMATCH";
export interface PageSchemaValidationIssue {
    path: string;
    message: string;
    keyword?: string;
}
/** Unified, observable error for the schema loading + validation pipeline. */
export declare class PageSchemaError extends Error {
    readonly code: PageSchemaErrorCode;
    readonly url: string;
    readonly issues?: PageSchemaValidationIssue[];
    constructor(code: PageSchemaErrorCode, url: string, message: string, issues?: PageSchemaValidationIssue[]);
}
export interface LoadPageOptions {
    /** Origin/base used to resolve relative schemaUrl values. Defaults to `location.origin`. */
    baseURL?: string;
    /** Injectable fetch; defaults to `globalThis.fetch`. */
    fetcher?: typeof fetch;
    /**
     * Optional in-memory document cache keyed by resolved schemaUrl (owned by
     * the App shell). A hit skips the network fetch AND the D-VAL structural
     * validation (the stored document was already validated at load time, and
     * its meta.pageId was verified against the manifest page). Intentionally
     * opt-in: tests and callers that must observe every load pass none.
     */
    cache?: Map<string, unknown>;
}
/**
 * Load and validate a page document for the given manifest page and route
 * params. Resolves the schemaUrl (expanding `{param}` placeholders), fetches
 * it, enforces structural validation, and verifies the document's `meta.pageId`
 * matches the manifest page. Returns the parsed page document on success.
 */
export declare function loadPageDocument(page: PageEntry, params: Record<string, string>, options?: LoadPageOptions): Promise<unknown>;
