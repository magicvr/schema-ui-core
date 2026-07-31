/**
 * search-table fixture adapter — four-layer query merge + selection events.
 */

import { serializeQuery, type QuerySource } from "./query-serialize";

export interface TableState {
  filters: Record<string, unknown>;
  page: number;
  pageSize: number;
  sort: string | null;
}

export interface SelectionState {
  keys: unknown[];
  count: number;
}

function isScalarKey(value: unknown): boolean {
  return (
    typeof value === "string" ||
    typeof value === "number" ||
    typeof value === "boolean"
  );
}

function dedupeKeys(keys: unknown[]): unknown[] {
  const seen = new Set<string>();
  const out: unknown[] = [];
  for (const key of keys) {
    if (!isScalarKey(key)) {
      continue;
    }
    const token = `${typeof key}:${String(key)}`;
    if (seen.has(token)) {
      continue;
    }
    seen.add(token);
    out.push(key);
  }
  return out;
}

function buildUrl(
  baseUrl: string,
  staticParams: Record<string, unknown>,
  state: TableState,
): string {
  const sources: QuerySource[] = [];
  const staticPairs: Array<[string, unknown]> = Object.entries(staticParams);
  if (staticPairs.length > 0) {
    sources.push(staticPairs);
  }
  const statePairs: Array<[string, unknown]> = [];
  for (const [key, value] of Object.entries(state.filters)) {
    statePairs.push([key, value]);
  }
  statePairs.push(["page", state.page]);
  statePairs.push(["pageSize", state.pageSize]);
  if (state.sort !== null && state.sort !== undefined) {
    statePairs.push(["sort", state.sort]);
  }
  sources.push(statePairs);
  const result = serializeQuery(baseUrl, sources);
  if (!result.ok) {
    throw new Error(result.code);
  }
  return result.url;
}

export function runSearchTable(input: Record<string, unknown>): Record<string, unknown> {
  const baseUrl = input.baseUrl as string;
  const staticParams = (input.staticParams as Record<string, unknown>) ?? {};
  let state: TableState = {
    filters: { ...((input.state as TableState)?.filters ?? {}) },
    page: (input.state as TableState)?.page ?? 1,
    pageSize: (input.state as TableState)?.pageSize ?? 20,
    sort: (input.state as TableState)?.sort ?? null,
  };
  let selection: SelectionState | undefined = input.selection
    ? {
        keys: [...((input.selection as SelectionState).keys ?? [])],
        count: (input.selection as SelectionState).count ?? 0,
      }
    : undefined;

  const event = input.event as Record<string, unknown> | null | undefined;
  const selectionEvent = input.selectionEvent as Record<string, unknown> | undefined;

  if (event) {
    switch (event.type) {
      case "submitSearch": {
        state = {
          ...state,
          filters: { ...((event.filters as Record<string, unknown>) ?? {}) },
          page: 1,
        };
        if (selection) {
          selection = { keys: [], count: 0 };
        }
        break;
      }
      case "clearSearch": {
        state = { ...state, filters: {}, page: 1 };
        break;
      }
      case "changePage": {
        state = { ...state, page: event.page as number };
        if (selection) {
          selection = { keys: [], count: 0 };
        }
        break;
      }
      case "changeSort": {
        state = { ...state, sort: event.sort as string, page: 1 };
        break;
      }
      default:
        break;
    }
  }

  if (selectionEvent?.type === "setKeys") {
    const keys = dedupeKeys((selectionEvent.keys as unknown[]) ?? []);
    selection = { keys, count: keys.length };
  }

  const url = buildUrl(baseUrl, staticParams, state);
  const result: Record<string, unknown> = { state, url };
  if (selection !== undefined || input.selection !== undefined || selectionEvent) {
    result.selection = selection ?? { keys: [], count: 0 };
  }
  return result;
}
