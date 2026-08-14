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

  constructor(
    status: number,
    code: string,
    message: string,
    messageKey?: string,
    params?: Record<string, unknown>,
    fieldErrors?: FieldError[],
  ) {
    super(message);
    this.name = "ResourceApiError";
    this.status = status;
    this.code = code;
    this.messageKey = messageKey;
    this.params = params;
    this.fieldErrors = fieldErrors ?? [];
  }
}

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
  const suffix =
    envelope.message === "" ? "" : envelope.message.startsWith(":") ? envelope.message : `: ${envelope.message}`;
  return new ResourceApiError(
    response.status,
    envelope.code,
    `${label} failed: HTTP ${response.status}${suffix}`,
    envelope.messageKey,
    envelope.params,
    envelope.fieldErrors,
  );
}

/** Serializes a ResourceQuery into a URL query string (query-serialization). */
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

/** Fetches a page of a schema-driven resource list (request-construction). */
export async function fetchResourceList(
  fetcher: typeof fetch,
  baseURL: string,
  query: ResourceQuery,
): Promise<ResourceList> {
  // F-001: validate BEFORE touching the (auth) fetcher so an invalid dataSource
  // never reaches Bearer-attaching transport.
  if (!isValidDataSource(baseURL)) {
    throw new Error(
      `invalid dataSource "${baseURL}": expected a single-slash same-origin path (no //, scheme, whitespace, backslash, ? or #)`,
    );
  }
  const queryString = buildResourceQuery(query);
  const url = queryString === "" ? baseURL : `${baseURL}?${queryString}`;
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
