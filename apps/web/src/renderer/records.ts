import type { SortOrder } from "@/components/data-table";

export const DEFAULT_RECORDS_URL = "/api/records";
export const DEFAULT_PAGE_SIZE = 10;

export interface RecordItem {
  id: string;
  name: string;
  status: string;
  owner: string;
  updatedAt: string;
}

export interface RecordList {
  items: RecordItem[];
  total: number;
  page: number;
  pageSize: number;
}

export interface RecordsQuery {
  q?: string;
  sort?: string;
  order?: SortOrder;
  page?: number;
  pageSize?: number;
}

/** A records API failure carrying the frozen envelope `{error, message}`. */
export class RecordApiError extends Error {
  readonly code: string;
  readonly status: number;

  constructor(status: number, code: string, message: string) {
    super(message);
    this.name = "RecordApiError";
    this.status = status;
    this.code = code;
  }
}

export interface RecordCreateBody {
  name: string;
  status: string;
  owner: string;
}

export interface RecordPatch {
  name?: string;
  status?: string;
  owner?: string;
}

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function fail(label: string, detail: string): never {
  throw new Error(`${label}: ${detail}`);
}

function requireString(value: unknown, label: string): string {
  if (typeof value !== "string") {
    return fail(label, "expected a string");
  }
  return value;
}

function requireNumber(value: unknown, label: string): number {
  if (typeof value !== "number" || !Number.isFinite(value)) {
    return fail(label, "expected a finite number");
  }
  return value;
}

function parseRecordItem(value: unknown, label: string): RecordItem {
  if (!isRecord(value)) {
    return fail(label, "expected an object");
  }
  return {
    id: requireString(value.id, `${label}.id`),
    name: requireString(value.name, `${label}.name`),
    status: requireString(value.status, `${label}.status`),
    owner: requireString(value.owner, `${label}.owner`),
    updatedAt: requireString(value.updatedAt, `${label}.updatedAt`),
  };
}

/** Reads the frozen error envelope from a non-OK response, defaulting safely. */
async function readEnvelope(response: Response): Promise<{ code: string; message: string }> {
  try {
    const value: unknown = await response.json();
    if (isRecord(value)) {
      return {
        code: typeof value.error === "string" && value.error !== "" ? value.error : "UNKNOWN",
        message: typeof value.message === "string" ? value.message : "",
      };
    }
  } catch {
    // non-JSON body
  }
  return { code: "UNKNOWN", message: "" };
}

/** Reads the frozen error envelope from a non-OK response into a RecordApiError. */
export async function readRecordApiError(response: Response, label: string): Promise<RecordApiError> {
  const envelope = await readEnvelope(response);
  const suffix =
    envelope.message === "" ? "" : envelope.message.startsWith(":") ? envelope.message : `: ${envelope.message}`;
  return new RecordApiError(
    response.status,
    envelope.code,
    `${label} failed: HTTP ${response.status}${suffix}`,
  );
}

/** Serializes a RecordsQuery into a URL query string (query-serialization). */
export function buildRecordsQuery(query: RecordsQuery): string {
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

/** Maps an unknown list payload to RecordList, fail-closed on shape drift. */
export function parseRecordList(value: unknown): RecordList {
  if (!isRecord(value)) {
    return fail("parseRecordList", "expected an object");
  }
  const items = value.items;
  if (!Array.isArray(items)) {
    return fail("parseRecordList", "expected items array");
  }
  const parsed = items.map((item, index) => parseRecordItem(item, `parseRecordList.items[${index}]`));
  return {
    items: parsed,
    total: requireNumber(value.total, "parseRecordList.total"),
    page: requireNumber(value.page, "parseRecordList.page"),
    pageSize: requireNumber(value.pageSize, "parseRecordList.pageSize"),
  };
}

/** Fetches a page of records from the Go list API (request-construction). */
export async function fetchRecords(
  fetcher: typeof fetch,
  baseURL: string,
  query: RecordsQuery,
): Promise<RecordList> {
  const queryString = buildRecordsQuery(query);
  const url = queryString === "" ? baseURL : `${baseURL}?${queryString}`;
  const response = await fetcher(url);
  if (!response.ok) {
    throw await readRecordApiError(response, "records fetch");
  }
  return parseRecordList(await response.json());
}

/** Creates a record via POST /api/records (records.write); returns the 201 record. */
export async function createRecord(
  fetcher: typeof fetch,
  baseURL: string,
  body: RecordCreateBody,
): Promise<RecordItem> {
  const response = await fetcher(baseURL, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!response.ok) {
    throw await readRecordApiError(response, "record create");
  }
  return parseRecordItem(await response.json(), "createRecord");
}

export async function updateRecord(
  fetcher: typeof fetch,
  baseURL: string,
  id: string,
  patch: RecordPatch,
): Promise<RecordItem> {
  const response = await fetcher(`${baseURL}/${encodeURIComponent(id)}`, {
    method: "PATCH",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(patch),
  });
  if (!response.ok) {
    throw await readRecordApiError(response, "record update");
  }
  return parseRecordItem(await response.json(), "updateRecord");
}

export async function deleteRecord(
  fetcher: typeof fetch,
  baseURL: string,
  id: string,
): Promise<void> {
  const response = await fetcher(`${baseURL}/${encodeURIComponent(id)}`, {
    method: "DELETE",
  });
  if (!response.ok && response.status !== 204) {
    throw await readRecordApiError(response, "record delete");
  }
}
