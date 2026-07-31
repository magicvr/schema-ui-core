/**
 * table-sort fixture adapter (ADR-0027 three-state sort + L2 validate).
 */

import { serializeQuery, type QuerySource } from "./query-serialize";

interface Column {
  field: string;
  label?: string;
  sortable?: boolean;
  sortField?: string;
}

interface TableSpec {
  columns: Column[];
  pagination?: { mode?: string; pageSize?: number };
  defaultSort?: { field: string; order: "asc" | "desc" };
}

interface TableState {
  filters: Record<string, unknown>;
  page: number;
  pageSize: number;
  sort: string | null;
}

function parseVersion(version: string): [number, number] | null {
  const m = /^(\d+)\.(\d+)$/.exec(version);
  if (!m) {
    return null;
  }
  return [Number(m[1]), Number(m[2])];
}

function versionGte(version: string, min: string): boolean {
  const a = parseVersion(version);
  const b = parseVersion(min);
  if (!a || !b) {
    return false;
  }
  if (a[0] !== b[0]) {
    return a[0] > b[0];
  }
  return a[1] >= b[1];
}

function sortKeyOf(column: Column): string {
  return column.sortField ?? column.field;
}

function knownSortKeys(table: TableSpec): Set<string> {
  const keys = new Set<string>();
  for (const column of table.columns) {
    if (column.sortable) {
      keys.add(sortKeyOf(column));
    }
  }
  return keys;
}

function buildUrl(baseUrl: string, state: TableState): string {
  const pairs: Array<[string, unknown]> = [];
  for (const [key, value] of Object.entries(state.filters)) {
    pairs.push([key, value]);
  }
  pairs.push(["page", state.page]);
  pairs.push(["pageSize", state.pageSize]);
  if (state.sort) {
    pairs.push(["sort", state.sort]);
  }
  const sources: QuerySource[] = [pairs];
  const result = serializeQuery(baseUrl, sources);
  if (!result.ok) {
    throw new Error(result.code);
  }
  return result.url;
}

function validateSortTable(
  meta: { protocolVersion?: string; requiredCapabilities?: string[] },
  table: TableSpec,
): Record<string, unknown> | null {
  const version = meta.protocolVersion ?? "0.0";
  const caps = meta.requiredCapabilities ?? [];
  const usesSort =
    table.columns.some((c) => c.sortable) || table.defaultSort !== undefined;
  if (!usesSort) {
    return null;
  }
  if (!versionGte(version, "2.5")) {
    return {
      ok: false,
      code: "PROTOCOL_VERSION_TOO_LOW",
      path: "meta.protocolVersion",
      errors: [{ code: "PROTOCOL_VERSION_TOO_LOW", path: "meta.protocolVersion" }],
    };
  }
  if (!caps.includes("table.sort")) {
    return {
      ok: false,
      code: "CAPABILITY_REQUIRED",
      path: "meta.requiredCapabilities",
      errors: [{ code: "CAPABILITY_REQUIRED", path: "meta.requiredCapabilities" }],
    };
  }
  const seen = new Set<string>();
  for (let i = 0; i < table.columns.length; i++) {
    const column = table.columns[i]!;
    if (!column.sortable) {
      continue;
    }
    const key = sortKeyOf(column);
    if (seen.has(key)) {
      return {
        ok: false,
        code: "SORT_KEY_DUPLICATE",
        path: `columns[${i}]`,
        errors: [{ code: "SORT_KEY_DUPLICATE", path: `columns[${i}]` }],
      };
    }
    seen.add(key);
  }
  return null;
}

export function runTableSort(input: Record<string, unknown>): Record<string, unknown> {
  const meta = (input.meta as { protocolVersion?: string; requiredCapabilities?: string[] }) ?? {};
  const table = input.table as TableSpec;
  const baseUrl = (input.baseUrl as string) ?? "/orders";

  if (input.operation === "validate") {
    const err = validateSortTable(meta, table);
    if (err) {
      return err;
    }
    return { ok: true };
  }

  // v2.4 compat: sort interaction disabled
  const protocolSortInteraction = versionGte(meta.protocolVersion ?? "0.0", "2.5");

  let state: TableState = input.state
    ? {
        filters: { ...((input.state as TableState).filters ?? {}) },
        page: (input.state as TableState).page ?? 1,
        pageSize: (input.state as TableState).pageSize ?? table.pagination?.pageSize ?? 20,
        sort: (input.state as TableState).sort ?? null,
      }
    : {
        filters: {},
        page: 1,
        pageSize: table.pagination?.pageSize ?? 20,
        sort: null,
      };

  let selection = input.selection
    ? {
        keys: [...((input.selection as { keys: unknown[] }).keys ?? [])],
        count: (input.selection as { count: number }).count ?? 0,
      }
    : undefined;

  const event = input.event as Record<string, unknown> | null | undefined;
  const keys = knownSortKeys(table);

  // Unknown current sort key blocks request
  if (state.sort) {
    const field = state.sort.split(":")[0] ?? "";
    if (!keys.has(field)) {
      return {
        ok: false,
        code: "TABLE_SORT_FIELD_UNKNOWN",
        requestEmitted: false,
        state,
      };
    }
  }

  let requestEmitted = false;

  if (!protocolSortInteraction) {
    if (event?.type === "clickSort") {
      return {
        ok: true,
        protocolSortInteraction: false,
        state,
        url: buildUrl(baseUrl, state),
      };
    }
  }

  if (event?.type === "init") {
    if (table.defaultSort) {
      state = {
        ...state,
        sort: `${table.defaultSort.field}:${table.defaultSort.order}`,
        page: 1,
      };
      requestEmitted = true;
    }
  } else if (event?.type === "clickSort") {
    const columnField = event.field as string;
    const column = table.columns.find((c) => c.field === columnField);
    if (!column || !column.sortable) {
      requestEmitted = true;
      return {
        ok: true,
        requestEmitted,
        state,
        url: buildUrl(baseUrl, state),
        ...(selection ? { selection } : {}),
      };
    }
    const key = sortKeyOf(column);
    const current = state.sort;
    let next: string | null;
    if (current === `${key}:asc`) {
      next = `${key}:desc`;
    } else if (current === `${key}:desc`) {
      next = null;
    } else {
      next = `${key}:asc`;
    }
    state = { ...state, sort: next, page: 1 };
    requestEmitted = true;
    if (selection) {
      selection = { keys: [], count: 0 };
    }
  } else if (event?.type === "submitSearch") {
    state = {
      ...state,
      filters: { ...((event.filters as Record<string, unknown>) ?? {}) },
      page: 1,
    };
    requestEmitted = true;
  } else if (event == null) {
    // no event — just evaluate current state (unknown sort already handled)
    requestEmitted = false;
  }

  const result: Record<string, unknown> = {
    ok: true,
    requestEmitted,
    state,
    url: buildUrl(baseUrl, state),
  };
  if (selection) {
    result.selection = selection;
  }
  if (!protocolSortInteraction) {
    result.protocolSortInteraction = false;
  }
  return result;
}
