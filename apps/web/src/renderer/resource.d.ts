import type { SortOrder } from "@/components/data-table";
export declare const DEFAULT_PAGE_SIZE = 10;
/**
 * Frozen list-endpoint rule (I-010-001 v0.2.0 · A-001 F-001): `table.props.dataSource`
 * must be a single-slash same-origin absolute path — starts with one `/`, no `//`
 * (no protocol-relative host), no scheme, no whitespace, backslash, `?` or `#`.
 * Query strings are appended by `buildResourceQuery`, never authored in dataSource;
 * fragments are never allowed. Validated before any (auth) fetch is attempted.
 * Mirrors `DataRef.url`'s `^/(?!/)[^\s\\]*$` but additionally rejects `?`/`#`.
 */
export declare const DATASOURCE_URL_PATTERN: RegExp;
/** A schema-driven resource row: any plain JSON object (no field whitelist). */
export interface ResourceItem {
    [key: string]: unknown;
}
/** Unified list envelope frozen across resources: `{items,total,page,pageSize}`. */
export interface ResourceList {
    items: ResourceItem[];
    total: number;
    page: number;
    pageSize: number;
}
/** A generic resource list query. */
export interface ResourceQuery {
    q?: string;
    sort?: string;
    order?: SortOrder;
    page?: number;
    pageSize?: number;
    /** Schema-driven table filters (table node props.filters): field to value. */
    filters?: Record<string, string>;
}
/**
 * Shared query for display components (statCard/chart) that read the first
 * item + envelope of a list endpoint. The pageSize keeps the historical RPC
 * shape (one generous bucket); wallet-ensure uses the exact same query so its
 * existence probe coalesces with the statCards' request (per-page fetch cache).
 */
export declare const DISPLAY_LIST_QUERY: ResourceQuery;
/** One field-level validation failure (GOAL-014 D-002 §2.1). */
export interface FieldError {
    field: string;
    reason: string;
}
/** A resource API failure carrying the frozen envelope `{error, message}`. */
export declare class ResourceApiError extends Error {
    readonly code: string;
    readonly status: number;
    /** VP-007 S4: stable catalog key when the server cataloged this error. */
    readonly messageKey?: string;
    /** VP-007 S4: interpolation params for messageKey. */
    readonly params?: Record<string, unknown>;
    /** GOAL-014: field-level validation failures, when the server attached them. */
    readonly fieldErrors: FieldError[];
    /** R1 correlation identifier for support/operator lookup. */
    readonly correlationId?: string;
    constructor(status: number, code: string, message: string, messageKey?: string, params?: Record<string, unknown>, fieldErrors?: FieldError[], correlationId?: string);
}
/** W19: missing self-wallet is an empty surface, not a hard page error. */
export declare function isWalletNotFoundError(err: unknown): boolean;
export declare const EMPTY_RESOURCE_LIST: ResourceList;
export type ResourceCreateBody = ResourceItem;
export type ResourcePatch = Partial<ResourceItem>;
/** F-001: single executable rule for a table list endpoint (fail-closed on any violation). */
export declare function isValidDataSource(url: string): boolean;
/** Reads the frozen error envelope from a non-OK response. */
export declare function readResourceApiError(response: Response, label: string): Promise<ResourceApiError>;
/** Serializes a ResourceQuery into a URL query string (query-serialization). */
/**
 * v2.9 dataSource params resolution (ADR-0039, capability data.route-binding).
 * Literal scalars pass through; whole `$context.route.query.*` /
 * `$context.route.params.*` bindings are resolved from the route snapshot;
 * a missing key (or absent route) is a tombstone — the parameter is dropped
 * (ADR-0010 semantics, same as null values). Returns the merged query string
 * (empty when nothing resolves).
 */
export declare function resolveDataParamsQuery(params: Record<string, unknown> | undefined, route: {
    query?: Record<string, string>;
    params?: Record<string, string>;
}): string;
export declare function buildResourceQuery(query: ResourceQuery): string;
/**
 * Maps an unknown list payload to the unified ResourceList envelope, fail-closed
 * on shape drift. Items are arbitrary JSON objects (F-002 rowKey is enforced at
 * the table surface, not here).
 */
export declare function parseResourceList(value: unknown): ResourceList;
/**
 * Builds the query-string portion of a resource list request — the exact
 * construction `fetchResourceList` sends on the wire. `extraQuery` carries
 * v2.9 ADR-0039 dataSource params that already resolved to literal key=value
 * pairs and is merged over the standard q/sort/order/page/pageSize query.
 */
export declare function resourceListQueryString(query: ResourceQuery, extraQuery?: string): string;
/**
 * The final URL for a resource list request (F-001-validated baseURL + query
 * string). Shared by `fetchResourceList` and the renderer's per-page fetch
 * cache so the cache key and the wire request can never drift apart.
 */
export declare function resourceListURL(baseURL: string, query: ResourceQuery, extraQuery?: string): string;
/**
 * Fetches a page of a schema-driven resource list (request-construction).
 *
 * `extraQuery` carries v2.9 ADR-0039 dataSource params that already resolved
 * to literal key=value pairs (e.g. `dictKey=order_status` from a
 * `$context.route.query.*` binding). It is merged with the standard
 * q/sort/order/page/pageSize query; baseURL itself must stay a bare
 * single-slash same-origin path (F-001).
 */
export declare function fetchResourceList(fetcher: typeof fetch, baseURL: string, query: ResourceQuery, extraQuery?: string): Promise<ResourceList>;
/** Creates a resource row via POST; returns the 201 row. */
export declare function createResource(fetcher: typeof fetch, baseURL: string, body: ResourceCreateBody): Promise<ResourceItem>;
export declare function updateResource(fetcher: typeof fetch, baseURL: string, id: string, patch: ResourcePatch): Promise<ResourceItem>;
export declare function deleteResource(fetcher: typeof fetch, baseURL: string, id: string): Promise<void>;
