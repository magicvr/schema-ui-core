/**
 * User-scoped Saved View persistence for schema tables (VP-037 R2).
 *
 * This module deliberately owns the storage contract rather than letting a
 * table write arbitrary JSON to localStorage. Every read and write projects
 * through the current table's allowlist so malformed, stale, or permission-
 * unsafe state cannot reach the resource query or column renderer.
 */

export const SAVED_VIEW_STORAGE_VERSION = 1 as const;
export const SAVED_VIEW_MAX_VIEWS = 50;
export const SAVED_VIEW_MAX_NAME_LENGTH = 80;
export const SAVED_VIEW_MAX_QUERY_LENGTH = 500;
export const SAVED_VIEW_MAX_FILTER_LENGTH = 200;
export const SAVED_VIEW_MAX_STORAGE_BYTES = 256 * 1024;

const STORAGE_KEY_PREFIX = "schema-ui.saved-views.v1";
const ALLOWED_PAGE_SIZES = new Set([10, 20, 50, 100]);
const ID_PATTERN = /^[A-Za-z0-9._:-]+$/;

export interface SavedViewStorageLike {
  getItem(key: string): string | null;
  setItem(key: string, value: string): void;
}

export interface SavedViewTableConfig {
  columnFields: readonly string[];
  sortableFields: readonly string[];
  filterFields: readonly string[];
}

export interface SavedViewQuery {
  q?: string;
  filters?: Record<string, string>;
  sort?: string;
  order?: "asc" | "desc";
  pageSize?: number;
}

export interface SavedViewState {
  query: SavedViewQuery;
  visibleColumns: string[];
}

export interface SavedViewRecord extends SavedViewState {
  id: string;
  name: string;
  createdAt: string;
  updatedAt: string;
}

export interface SavedViewStorageDocument {
  version: typeof SAVED_VIEW_STORAGE_VERSION;
  views: SavedViewRecord[];
  activeViewId?: string;
}

export type SavedViewErrorCode =
  | "SAVED_VIEW_STORAGE_UNAVAILABLE"
  | "SAVED_VIEW_STORAGE_INVALID"
  | "SAVED_VIEW_STORAGE_WRITE_FAILED"
  | "SAVED_VIEW_NAME_INVALID"
  | "SAVED_VIEW_STATE_INVALID"
  | "SAVED_VIEW_LIMIT_REACHED";

export interface SavedViewFailure {
  ok: false;
  code: SavedViewErrorCode;
  message: string;
}

export interface SavedViewReadSuccess {
  ok: true;
  views: SavedViewRecord[];
  activeViewId?: string;
  /** Number of records dropped because they failed the current allowlist. */
  droppedCount: number;
}

export type SavedViewReadResult = SavedViewReadSuccess | SavedViewFailure;
export type SavedViewWriteResult = { ok: true } | SavedViewFailure;
export type SavedViewRecordResult =
  | { ok: true; view: SavedViewRecord }
  | SavedViewFailure;

function isSavedViewFailure(value: unknown): value is SavedViewFailure {
  return isRecord(value) && value.ok === false;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function uniqueFields(fields: readonly string[]): string[] {
  return [...new Set(fields.filter((field) => typeof field === "string" && field !== ""))];
}

function isAllowedField(field: string, fields: readonly string[]): boolean {
  return fields.includes(field);
}

function normalizeName(name: unknown): string | SavedViewFailure {
  if (typeof name !== "string") {
    return { ok: false, code: "SAVED_VIEW_NAME_INVALID", message: "Saved View name must be text." };
  }
  const normalized = name.trim();
  if (
    normalized === "" ||
    normalized.length > SAVED_VIEW_MAX_NAME_LENGTH ||
    /[\u0000-\u001f\u007f]/.test(normalized)
  ) {
    return {
      ok: false,
      code: "SAVED_VIEW_NAME_INVALID",
      message: `Saved View name must be 1-${SAVED_VIEW_MAX_NAME_LENGTH} visible characters.`,
    };
  }
  return normalized;
}

function normalizeQuery(
  raw: unknown,
  config: SavedViewTableConfig,
): SavedViewQuery | SavedViewFailure {
  if (!isRecord(raw)) {
    return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View query is not an object." };
  }
  const allowedKeys = new Set(["q", "filters", "sort", "order", "pageSize"]);
  if (Object.keys(raw).some((key) => !allowedKeys.has(key))) {
    return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View query contains an unknown field." };
  }
  const query: SavedViewQuery = {};
  if (raw.q !== undefined) {
    if (typeof raw.q !== "string") {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View query text is invalid." };
    }
    const q = raw.q.trim();
    if (q.length > SAVED_VIEW_MAX_QUERY_LENGTH) {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View query text is too long." };
    }
    if (q !== "") query.q = q;
  }
  if (raw.filters !== undefined) {
    if (!isRecord(raw.filters)) {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View filters are invalid." };
    }
    const filters: Record<string, string> = {};
    for (const [field, rawValue] of Object.entries(raw.filters)) {
      if (!isAllowedField(field, config.filterFields) || typeof rawValue !== "string") {
        return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View contains an unsupported filter." };
      }
      const value = rawValue.trim();
      if (value.length > SAVED_VIEW_MAX_FILTER_LENGTH) {
        return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View filter value is too long." };
      }
      if (value !== "") filters[field] = value;
    }
    if (Object.keys(filters).length > 0) query.filters = filters;
  }
  if (raw.sort !== undefined) {
    if (typeof raw.sort !== "string" || !isAllowedField(raw.sort, config.sortableFields)) {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View sort field is not allowed." };
    }
    query.sort = raw.sort;
  }
  if (raw.order !== undefined) {
    if (raw.order !== "asc" && raw.order !== "desc") {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View sort order is invalid." };
    }
    if (query.sort === undefined) {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View sort order has no sort field." };
    }
    query.order = raw.order;
  }
  if (raw.pageSize !== undefined) {
    if (typeof raw.pageSize !== "number" || !Number.isInteger(raw.pageSize) || !ALLOWED_PAGE_SIZES.has(raw.pageSize)) {
      return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View page size is invalid." };
    }
    query.pageSize = raw.pageSize;
  }
  return query;
}

function normalizeVisibleColumns(
  raw: unknown,
  config: SavedViewTableConfig,
): string[] | SavedViewFailure {
  if (!Array.isArray(raw)) {
    return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View columns are invalid." };
  }
  const columnFields = uniqueFields(config.columnFields);
  const visible = raw.filter((field): field is string => typeof field === "string");
  if (
    visible.length !== raw.length ||
    new Set(visible).size !== visible.length ||
    visible.some((field) => !columnFields.includes(field)) ||
    (columnFields.length > 0 && visible.length === 0)
  ) {
    return { ok: false, code: "SAVED_VIEW_STATE_INVALID", message: "Saved View columns are outside the current Schema." };
  }
  return visible;
}

export function normalizeSavedViewState(
  query: unknown,
  visibleColumns: unknown,
  config: SavedViewTableConfig,
): { ok: true; state: SavedViewState } | SavedViewFailure {
  const normalizedQuery = normalizeQuery(query, config);
  if (isSavedViewFailure(normalizedQuery)) {
    return normalizedQuery;
  }
  const normalizedColumns = normalizeVisibleColumns(visibleColumns, config);
  if (isSavedViewFailure(normalizedColumns)) {
    return normalizedColumns;
  }
  return {
    ok: true,
    state: { query: normalizedQuery, visibleColumns: normalizedColumns },
  };
}

function normalizeTimestamp(value: unknown): string | null {
  return typeof value === "string" && value !== "" && value.length <= 64 ? value : null;
}

function normalizeRecord(
  raw: unknown,
  config: SavedViewTableConfig,
): SavedViewRecord | null {
  if (!isRecord(raw)) return null;
  const allowedKeys = new Set(["id", "name", "query", "visibleColumns", "createdAt", "updatedAt"]);
  if (Object.keys(raw).some((key) => !allowedKeys.has(key))) return null;
  if (
    typeof raw.id !== "string" ||
    raw.id === "" ||
    raw.id.length > 120 ||
    !ID_PATTERN.test(raw.id)
  ) {
    return null;
  }
  const name = normalizeName(raw.name);
  const createdAt = normalizeTimestamp(raw.createdAt);
  const updatedAt = normalizeTimestamp(raw.updatedAt);
  if (typeof name !== "string" || createdAt === null || updatedAt === null) return null;
  const state = normalizeSavedViewState(raw.query, raw.visibleColumns, config);
  if (!state.ok) return null;
  return { id: raw.id, name, ...state.state, createdAt, updatedAt };
}

function failure(code: SavedViewErrorCode, message: string): SavedViewFailure {
  return { ok: false, code, message };
}

/** Builds the per-user/per-page/per-table key required by D-003. */
export function savedViewStorageKey(userId: string, pageId: string, tableId: string): string {
  const encodedSegments = [userId, pageId, tableId]
    // encodeURIComponent leaves `.` untouched; encode it explicitly so the
    // dot separator cannot become ambiguous if an ID alphabet expands later.
    .map((segment) => encodeURIComponent(segment).replaceAll(".", "%2E"));
  return [STORAGE_KEY_PREFIX, ...encodedSegments].join(".");
}

/** Returns browser storage without turning an unavailable storage API into a fake success. */
export function getBrowserSavedViewStorage(): SavedViewStorageLike | null {
  if (typeof window === "undefined") return null;
  try {
    return window.localStorage;
  } catch {
    return null;
  }
}

export function readSavedViews(
  storage: SavedViewStorageLike | null,
  key: string,
  config: SavedViewTableConfig,
): SavedViewReadResult {
  if (storage === null) {
    return failure("SAVED_VIEW_STORAGE_UNAVAILABLE", "Browser storage is unavailable.");
  }
  let raw: string | null;
  try {
    raw = storage.getItem(key);
  } catch {
    return failure("SAVED_VIEW_STORAGE_UNAVAILABLE", "Browser storage could not be read.");
  }
  if (raw === null || raw === "") {
    return { ok: true, views: [], droppedCount: 0 };
  }
  if (raw.length > SAVED_VIEW_MAX_STORAGE_BYTES) {
    return failure("SAVED_VIEW_STORAGE_INVALID", "Saved View storage exceeds the safety limit.");
  }
  let parsed: unknown;
  try {
    parsed = JSON.parse(raw);
  } catch {
    return failure("SAVED_VIEW_STORAGE_INVALID", "Saved View storage is not valid JSON.");
  }
  if (!isRecord(parsed) || parsed.version !== SAVED_VIEW_STORAGE_VERSION || !Array.isArray(parsed.views)) {
    return failure("SAVED_VIEW_STORAGE_INVALID", "Saved View storage has an unsupported shape.");
  }
  const allowedDocumentKeys = new Set(["version", "views", "activeViewId"]);
  if (Object.keys(parsed).some((key) => !allowedDocumentKeys.has(key))) {
    return failure("SAVED_VIEW_STORAGE_INVALID", "Saved View storage contains an unknown field.");
  }
  if (parsed.activeViewId !== undefined && typeof parsed.activeViewId !== "string") {
    return failure("SAVED_VIEW_STORAGE_INVALID", "Saved View active state is invalid.");
  }
  if (parsed.views.length > SAVED_VIEW_MAX_VIEWS) {
    return failure("SAVED_VIEW_STORAGE_INVALID", "Saved View storage contains too many views.");
  }
  const views: SavedViewRecord[] = [];
  let droppedCount = 0;
  const seen = new Set<string>();
  for (const entry of parsed.views) {
    const view = normalizeRecord(entry, config);
    if (view === null || seen.has(view.id)) {
      droppedCount += 1;
      continue;
    }
    seen.add(view.id);
    views.push(view);
  }
  const activeViewId =
    typeof parsed.activeViewId === "string" && views.some((view) => view.id === parsed.activeViewId)
      ? parsed.activeViewId
      : undefined;
  return { ok: true, views, ...(activeViewId === undefined ? {} : { activeViewId }), droppedCount };
}

export function writeSavedViews(
  storage: SavedViewStorageLike | null,
  key: string,
  views: readonly SavedViewRecord[],
  config: SavedViewTableConfig,
  activeViewId?: string,
): SavedViewWriteResult {
  if (storage === null) {
    return failure("SAVED_VIEW_STORAGE_UNAVAILABLE", "Browser storage is unavailable.");
  }
  if (views.length > SAVED_VIEW_MAX_VIEWS) {
    return failure("SAVED_VIEW_LIMIT_REACHED", `A maximum of ${SAVED_VIEW_MAX_VIEWS} Saved Views is supported.`);
  }
  const normalizedViews: SavedViewRecord[] = [];
  const seen = new Set<string>();
  for (const entry of views) {
    const view = normalizeRecord(entry, config);
    if (view === null || seen.has(view.id)) {
      return failure("SAVED_VIEW_STATE_INVALID", "Saved View state no longer matches the current Schema.");
    }
    seen.add(view.id);
    normalizedViews.push(view);
  }
  if (activeViewId !== undefined && !seen.has(activeViewId)) {
    return failure("SAVED_VIEW_STATE_INVALID", "The active Saved View is no longer available.");
  }
  const document: SavedViewStorageDocument = {
    version: SAVED_VIEW_STORAGE_VERSION,
    views: normalizedViews,
    ...(activeViewId === undefined ? {} : { activeViewId }),
  };
  let serialized: string;
  try {
    serialized = JSON.stringify(document);
  } catch {
    return failure("SAVED_VIEW_STATE_INVALID", "Saved View state could not be serialized.");
  }
  if (serialized.length > SAVED_VIEW_MAX_STORAGE_BYTES) {
    return failure("SAVED_VIEW_LIMIT_REACHED", "Saved View storage exceeds the safety limit.");
  }
  try {
    storage.setItem(key, serialized);
  } catch {
    return failure("SAVED_VIEW_STORAGE_WRITE_FAILED", "Browser storage could not be written.");
  }
  return { ok: true };
}

function createId(): string {
  try {
    const candidate = globalThis.crypto?.randomUUID?.();
    if (typeof candidate === "string") return `view-${candidate}`;
  } catch {
    // Use the bounded timestamp fallback below.
  }
  return `view-${Date.now()}-${Math.random().toString(36).slice(2, 10)}`;
}

export function createSavedViewRecord(
  name: unknown,
  query: unknown,
  visibleColumns: unknown,
  config: SavedViewTableConfig,
  now = new Date().toISOString(),
  id = createId(),
): SavedViewRecordResult {
  const normalizedName = normalizeName(name);
  if (typeof normalizedName !== "string") return normalizedName;
  const state = normalizeSavedViewState(query, visibleColumns, config);
  if (!state.ok) return state;
  if (typeof id !== "string" || id === "" || id.length > 120 || !ID_PATTERN.test(id)) {
    return failure("SAVED_VIEW_STATE_INVALID", "Saved View id is invalid.");
  }
  return {
    ok: true,
    view: {
      id,
      name: normalizedName,
      ...state.state,
      createdAt: now,
      updatedAt: now,
    },
  };
}

export function updateSavedViewRecord(
  existing: SavedViewRecord,
  query: unknown,
  visibleColumns: unknown,
  config: SavedViewTableConfig,
  now = new Date().toISOString(),
): SavedViewRecordResult {
  const state = normalizeSavedViewState(query, visibleColumns, config);
  if (!state.ok) return state;
  const name = normalizeName(existing.name);
  if (typeof name !== "string") return name;
  if (typeof existing.id !== "string" || existing.id === "" || !ID_PATTERN.test(existing.id)) {
    return failure("SAVED_VIEW_STATE_INVALID", "Saved View id is invalid.");
  }
  const createdAt = normalizeTimestamp(existing.createdAt);
  if (createdAt === null) return failure("SAVED_VIEW_STATE_INVALID", "Saved View creation time is invalid.");
  return {
    ok: true,
    view: { id: existing.id, name, ...state.state, createdAt, updatedAt: now },
  };
}
