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
  const parsed = items.map((item, index) => {
    if (!isRecord(item)) {
      return fail("parseRecordList", `items[${index}] expected an object`);
    }
    return {
      id: requireString(item.id, `items[${index}].id`),
      name: requireString(item.name, `items[${index}].name`),
      status: requireString(item.status, `items[${index}].status`),
      owner: requireString(item.owner, `items[${index}].owner`),
      updatedAt: requireString(item.updatedAt, `items[${index}].updatedAt`),
    };
  });
  return {
    items: parsed,
    total: requireNumber(value.total, "total"),
    page: requireNumber(value.page, "page"),
    pageSize: requireNumber(value.pageSize, "pageSize"),
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
    throw new Error(`records fetch failed: HTTP ${response.status}`);
  }
  return parseRecordList(await response.json());
}

export interface RecordPatch {
  name?: string;
  status?: string;
  owner?: string;
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
    throw new Error(`record update failed: HTTP ${response.status}`);
  }
  const value: unknown = await response.json();
  if (!isRecord(value)) {
    return fail("updateRecord", "expected an object response");
  }
  return {
    id: requireString(value.id, "id"),
    name: requireString(value.name, "name"),
    status: requireString(value.status, "status"),
    owner: requireString(value.owner, "owner"),
    updatedAt: requireString(value.updatedAt, "updatedAt"),
  };
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
    throw new Error(`record delete failed: HTTP ${response.status}`);
  }
}
