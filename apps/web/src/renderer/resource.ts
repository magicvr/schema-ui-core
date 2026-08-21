import type { SortOrder } from "@/components/data-table";

export const DEFAULT_PAGE_SIZE = 10;

/**
 * Frozen list-endpoint rule (I-010-001 v0.2.0 · A-001 F-001): `table.props.dataSource`
 * must be a single-slash same-origin absolute path — starts with one `/`, no `//`
 * (no protocol-relative host), no scheme, no whitespace, backslash, `?` or `#`.
 * Query strings are appended by `buildResourceQuery`, never authored in dataSource;
 * fragments are never allowed. Validated before any (auth) fetch is attempted.
 * Mirrors `DataRef.url`'s `^/(?!/)[^\s\\]*$` but additionally rejects `?`/`#`.
 */
export const DATASOURCE_URL_PATTERN = /^\/(?!\/)[^\s\\?#]*$/;

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

/** One field-level validation failure (GOAL-014 D-002 §2.1). */
export interface FieldError {
  field: string;
  reason: string;
}

/** A resource API failure carrying the frozen envelope `{error, message}`. */
export class ResourceApiError extends Error {
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

  constructor(
    status: number,
    code: string,
    message: string,
    messageKey?: string,
    params?: Record<string, unknown>,
    fieldErrors?: FieldError[],
    correlationId?: string,
  ) {
    super(message);
    this.name = "ResourceApiError";
    this.status = status;
    this.code = code;
    this.messageKey = messageKey;
    this.params = params;
    this.fieldErrors = fieldErrors ?? [];
    this.correlationId = correlationId;
  }
}

/** W19: missing self-wallet is an empty surface, not a hard page error. */
export function isWalletNotFoundError(err: unknown): boolean {
  return err instanceof ResourceApiError && err.code === "WALLET_NOT_FOUND";
}

export const EMPTY_RESOURCE_LIST: ResourceList = { items: [], total: 0, page: 1, pageSize: 100 };

export type ResourceCreateBody = ResourceItem;
export type ResourcePatch = Partial<ResourceItem>;

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function fail(label: string, detail: string): never {
  throw new Error(`${label}: ${detail}`);
}

function requireNumber(value: unknown, label: string): number {
  if (typeof value !== "number" || !Number.isFinite(value)) {
    return fail(label, "expected a finite number");
  }
  return value;
}

/** F-001: single executable rule for a table list endpoint (fail-closed on any violation). */
export function isValidDataSource(url: string): boolean {
  return typeof url === "string" && url !== "" && DATASOURCE_URL_PATTERN.test(url);
}

/** Parses one resource row: any plain JSON object (no per-field whitelist). */
function parseResourceItem(value: unknown, label: string): ResourceItem {
  if (!isRecord(value)) {
    return fail(label, "expected an object");
  }
  return value;
}

interface Envelope {
  code: string;
  message: string;
  messageKey?: string;
  params?: Record<string, unknown>;
  fieldErrors?: FieldError[];
  correlationId?: string;
}

function parseFieldErrors(value: unknown): FieldError[] {
  if (!Array.isArray(value)) return [];
  const out: FieldError[] = [];
  for (const raw of value) {
    if (!isRecord(raw)) continue;
    const field = typeof raw.field === "string" ? raw.field : "";
    const reason = typeof raw.reason === "string" ? raw.reason : "";
    if (field !== "" || reason !== "") {
      out.push({ field, reason });
    }
  }
  return out;
}

/** Reads the frozen error envelope from a non-OK response, defaulting safely. */
async function readEnvelope(response: Response): Promise<Envelope> {
  try {
    const value: unknown = await response.json();
    if (isRecord(value)) {
      return {
        code: typeof value.error === "string" && value.error !== "" ? value.error : "UNKNOWN",
        message: typeof value.message === "string" ? value.message : "",
        ...(typeof value.messageKey === "string" && value.messageKey !== ""
          ? { messageKey: value.messageKey }
          : {}),
        ...(isRecord(value.params) ? { params: value.params } : {}),
        ...(value.fieldErrors !== undefined ? { fieldErrors: parseFieldErrors(value.fieldErrors) } : {}),
        ...(typeof value.correlation_id === "string" && value.correlation_id !== ""
          ? { correlationId: value.correlation_id }
          : typeof value.correlationId === "string" && value.correlationId !== ""
            ? { correlationId: value.correlationId }
            : {}),
      };
    }
  } catch {
    // non-JSON body
  }
  return { code: "UNKNOWN", message: "" };
}

/** Reads the frozen error envelope from a non-OK response. */
export async function readResourceApiError(response: Response, label: string): Promise<ResourceApiError> {
  const envelope = await readEnvelope(response);
  const correlationId = envelope.correlationId ?? response.headers.get("X-Request-ID") ?? undefined;
  const suffix =
    envelope.message === "" ? "" : envelope.message.startsWith(":") ? envelope.message : `: ${envelope.message}`;
  const correlationSuffix = correlationId === undefined ? "" : ` (request ${correlationId})`;
  return new ResourceApiError(
    response.status,
    envelope.code,
    `${label} failed: HTTP ${response.status}${suffix}${correlationSuffix}`,
    envelope.messageKey,
    envelope.params,
    envelope.fieldErrors,
    correlationId,
  );
}

/** Serializes a ResourceQuery into a URL query string (query-serialization). */

/**
 * v2.9 dataSource params resolution (ADR-0039, capability data.route-binding).
 * Literal scalars pass through; whole `$context.route.query.*` /
 * `$context.route.params.*` bindings are resolved from the route snapshot;
 * a missing key (or absent route) is a tombstone — the parameter is dropped
 * (ADR-0010 semantics, same as null values). Returns the merged query string
 * (empty when nothing resolves).
 */
export function resolveDataParamsQuery(
  params: Record<string, unknown> | undefined,
  route: { query?: Record<string, string>; params?: Record<string, string> },
): string {
  const out = new URLSearchParams();
  if (params === undefined) {
    return "";
  }
  for (const [key, value] of Object.entries(params)) {
    if (typeof value === "string" && value.startsWith("$context.route.")) {
      if (value.startsWith("$context.route.query.")) {
        const routeKey = value.slice("$context.route.query.".length);
        const resolved = route.query?.[routeKey];
        if (resolved === undefined || resolved === null) {
          continue; // tombstone: drop the parameter
        }
        out.set(key, String(resolved));
        continue;
      }
      if (value.startsWith("$context.route.params.")) {
        const routeKey = value.slice("$context.route.params.".length);
        const resolved = route.params?.[routeKey];
        if (resolved === undefined || resolved === null) {
          continue; // tombstone: drop the parameter
        }
        out.set(key, String(resolved));
        continue;
      }
      // Unknown $context.route.* shape: fail closed by dropping the param.
      continue;
    }
    if (value === null || value === undefined) {
      continue; // tombstone
    }
    out.set(key, String(value));
  }
  return out.toString();
}

export function buildResourceQuery(query: ResourceQuery): string {
  const params = new URLSearchParams();
  if (query.q !== undefined && query.q.trim() !== "") {
    params.set("q", query.q.trim());
  }
  if (query.sort !== undefined) {
    params.set("sort", query.sort);
  }
  if (query.order !== undefined) {
    params.set("order", query.order);
  }
  if (query.page !== undefined && query.page > 1) {
    params.set("page", String(query.page));
  }
  if (query.pageSize !== undefined && query.pageSize !== DEFAULT_PAGE_SIZE) {
    params.set("pageSize", String(query.pageSize));
  }
  if (query.filters !== undefined) {
    for (const [field, value] of Object.entries(query.filters)) {
      // An empty filter value means "all" — omit the parameter entirely so
      // clearing a filter never sends a stale/meaningless value.
      if (field !== "" && value !== "") {
        params.set(field, value);
      }
    }
  }
  return params.toString();
}

/**
 * Maps an unknown list payload to the unified ResourceList envelope, fail-closed
 * on shape drift. Items are arbitrary JSON objects (F-002 rowKey is enforced at
 * the table surface, not here).
 */
export function parseResourceList(value: unknown): ResourceList {
  if (!isRecord(value)) {
    return fail("parseResourceList", "expected an object");
  }
  const items = value.items;
  if (!Array.isArray(items)) {
    return fail("parseResourceList", "expected items array");
  }
  const parsed = items.map((item, index) => parseResourceItem(item, `parseResourceList.items[${index}]`));
  return {
    items: parsed,
    total: requireNumber(value.total, "parseResourceList.total"),
    page: requireNumber(value.page, "parseResourceList.page"),
    pageSize: requireNumber(value.pageSize, "parseResourceList.pageSize"),
  };
}

/**
 * Fetches a page of a schema-driven resource list (request-construction).
 *
 * `extraQuery` carries v2.9 ADR-0039 dataSource params that already resolved
 * to literal key=value pairs (e.g. `dictKey=order_status` from a
 * `$context.route.query.*` binding). It is merged with the standard
 * q/sort/order/page/pageSize query; baseURL itself must stay a bare
 * single-slash same-origin path (F-001).
 */
export async function fetchResourceList(
  fetcher: typeof fetch,
  baseURL: string,
  query: ResourceQuery,
  extraQuery?: string,
): Promise<ResourceList> {
  // F-001: validate BEFORE touching the (auth) fetcher so an invalid dataSource
  // never reaches Bearer-attaching transport.
  if (!isValidDataSource(baseURL)) {
    throw new Error(
      `invalid dataSource "${baseURL}": expected a single-slash same-origin path (no //, scheme, whitespace, backslash, ? or #)`,
    );
  }
  const params = new URLSearchParams(buildResourceQuery(query));
  if (extraQuery !== undefined && extraQuery !== "") {
    for (const [key, value] of new URLSearchParams(extraQuery).entries()) {
      params.set(key, value);
    }
  }
  const url = params.size === 0 ? baseURL : `${baseURL}?${params.toString()}`;
  const response = await fetcher(url);
  if (!response.ok) {
    throw await readResourceApiError(response, "resource fetch");
  }
  return parseResourceList(await response.json());
}

/** Creates a resource row via POST; returns the 201 row. */
export async function createResource(
  fetcher: typeof fetch,
  baseURL: string,
  body: ResourceCreateBody,
): Promise<ResourceItem> {
  const response = await fetcher(baseURL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!response.ok) {
    throw await readResourceApiError(response, "resource create");
  }
  return parseResourceItem(await response.json(), "createResource");
}

export async function updateResource(
  fetcher: typeof fetch,
  baseURL: string,
  id: string,
  patch: ResourcePatch,
): Promise<ResourceItem> {
  const response = await fetcher(`${baseURL}/${encodeURIComponent(id)}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patch),
  });
  if (!response.ok) {
    throw await readResourceApiError(response, "resource update");
  }
  return parseResourceItem(await response.json(), "updateResource");
}

export async function deleteResource(
  fetcher: typeof fetch,
  baseURL: string,
  id: string,
): Promise<void> {
  const response = await fetcher(`${baseURL}/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
  if (!response.ok && response.status !== 204) {
    throw await readResourceApiError(response, "resource delete");
  }
}
