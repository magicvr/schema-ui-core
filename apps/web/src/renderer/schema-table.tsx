import { useEffect, useMemo, useState } from "react";

import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";
import { resolveTextProp } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
import {
  fetchResourceList,
  isValidDataSource,
  resolveDataParamsQuery,
  type ResourceQuery,
  type ResourceItem,
  type ResourceList,
} from "@/renderer/resource";
import type { RenderTableNode } from "@/renderer/render";
import { useSchemaCrud } from "@/renderer/render.tsx";

/**
 * Default schema-driven table surface (R1 · GOAL-004 / D-004) + S4 CRUD wiring
 * + A-002 generic adapter (GOAL-010 S3).
 *
 * Renders a whitelisted table node's `props.columns` over its `props.dataSource`.
 * Fails closed when the node declares no columns, no (valid) data source, or a
 * response that violates the frozen rowKey invariant (F-001 / F-002).
 *
 * F-001 (I-010-001 v0.2.0 §2): `dataSource` must be a single-slash same-origin
 * path; invalid/absent sources render an observable fail-closed state and never
 * reach the (auth) fetcher.
 *
 * F-002 (I-010-001 v0.2.0 §3): `props.rowKey` (default `"id"`) is the direct
 * field name of each row's unique, non-empty scalar key (string / finite
 * number). A missing / empty / non-scalar / duplicate key stops rendering the
 * table data and forbids row actions and selection.
 *
 * S4 (GOAL-007 · I-007-003 v0.2.2 §9): when rendered inside a `RenderPage`
 * (which always wraps content in the SchemaCrudProvider), the table surfaces
 * `props.toolbar` (create trigger) and `props.actions` (row edit/delete),
 * participates in row selection for the page's recordView / edit-form prefill,
 * reads its query from the shared per-table query state (so the search-form
 * binding and the post-write reload both reach it), and registers its injected
 * fetcher so modal form submits / row actions use the same transport. All of
 * this is fixture-driven: no record-specific logic lives here.
 */

export interface SchemaTableProps {
  node: RenderTableNode;
  /** Injectable fetch (defaults to `globalThis.fetch`). */
  fetcher?: typeof fetch;
}

export interface SchemaTableColumnSpec {
  field: string;
  label?: string;
  sortable?: boolean;
  /** W4 · GOAL-005: single-line truncate + title affordance for long values. */
  truncate?: boolean;
}

function stringOf(value: unknown): string {
  return typeof value === "string" ? value : "";
}

/** Spreads a typed row into a plain object the generic action executor accepts. */
function rowAsRecord(row: ResourceItem): Record<string, unknown> {
  return { ...row };
}

/** Structured row predicate used by Schema actions; malformed predicates fail closed. */
export function rowActionDisabled(action: unknown, row: ResourceItem): boolean {
  if (typeof action !== "object" || action === null || Array.isArray(action)) {
    return true;
  }
  const condition = (action as Record<string, unknown>).disabledWhen;
  if (condition === undefined) {
    return false;
  }
  if (typeof condition !== "object" || condition === null || Array.isArray(condition)) {
    return true;
  }
  const field = (condition as Record<string, unknown>).field;
  if (typeof field !== "string" || field === "" || !("equals" in condition)) {
    return true;
  }
  if (!(field in row)) {
    return true;
  }
  return row[field] === (condition as Record<string, unknown>).equals;
}

function isColumnSpec(value: unknown): value is SchemaTableColumnSpec {
  return (
    typeof value === "object" &&
    value !== null &&
    typeof (value as Record<string, unknown>).field === "string"
  );
}

/** Extracts the table node's column specs, fail-closed on malformed entries. */
export function schemaTableColumns(node: RenderTableNode): SchemaTableColumnSpec[] {
  const raw = node.props?.columns;
  return Array.isArray(raw)
    ? (raw as unknown[]).filter((entry): entry is SchemaTableColumnSpec =>
        isColumnSpec(entry),
      )
    : [];
}

/**
 * Resolves the table node's list endpoint (F-001). v2.9 prefers the node-level
 * DataRef (ADR-0039, source:api url); the legacy props.dataSource string stays
 * supported. Returns null when absent or not a single-slash same-origin path;
 * the table then fails closed and never fetches (the fixture fallback was
 * removed in GOAL-010 S3).
 */
export function schemaTableDataSource(node: RenderTableNode): string | null {
  const dataUrl = node.data?.url;
  if (typeof dataUrl === "string" && isValidDataSource(dataUrl)) {
    return dataUrl;
  }
  const raw = node.props?.dataSource;
  return typeof raw === "string" && isValidDataSource(raw) ? raw : null;
}

/** v2.9 node-level DataRef params (ADR-0039 route bindings), else empty. */
export function schemaTableDataParams(node: RenderTableNode): Record<string, unknown> | undefined {
  return node.data?.params;
}

/** F-002: the direct field name used as each row's unique key (default "id"). */
export function schemaTableRowKey(node: RenderTableNode): string {
  const raw = node.props?.rowKey;
  return typeof raw === "string" && raw !== "" ? raw : "id";
}

/** F-002: a row key must be a non-empty string or a finite number (JSON scalar). */
function scalarRowKey(value: unknown): string | null {
  if (typeof value === "string" && value !== "") {
    return value;
  }
  if (typeof value === "number" && Number.isFinite(value)) {
    return String(value);
  }
  return null;
}

type RowKeyCheck = { ok: true } | { ok: false; message: string };

/** F-002: every row must carry a non-empty, unique, scalar key in `field`. */
function checkRowKeys(items: ResourceItem[], field: string): RowKeyCheck {
  const seen = new Set<string>();
  for (const [index, row] of items.entries()) {
    const key = scalarRowKey(row[field]);
    if (key === null) {
      return {
        ok: false,
        message: `row ${index} has no valid "${field}" key (expected a non-empty string or finite number)`,
      };
    }
    if (seen.has(key)) {
      return { ok: false, message: `duplicate row key "${key}" for "${field}"` };
    }
    seen.add(key);
  }
  return { ok: true };
}

export function SchemaTable({ node, fetcher }: SchemaTableProps) {
  const columns = schemaTableColumns(node);
  const dataSource = schemaTableDataSource(node);
  const dataParams = schemaTableDataParams(node);
  const rowKeyField = schemaTableRowKey(node);
  const crud = useSchemaCrud();
  const t = useTranslate();
  const tableId = node.id ?? "default";
  const rowActions = Array.isArray(node.props?.actions) ? node.props.actions : [];
  const toolbar = Array.isArray(node.props?.toolbar) ? node.props.toolbar : [];

  // Register the injected transport with the page's Schema CRUD provider so
  // modal form submits and row actions share the same fetcher (S4).
  useEffect(() => {
    if (crud !== null && fetcher !== undefined) {
      crud.registerFetcher(fetcher);
    }
  }, [crud, fetcher]);

  const providerQuery = crud?.tableQuery(tableId);
  const [localQuery, setLocalQuery] = useState<ResourceQuery>({ page: 1, pageSize: 10 });
  const query = providerQuery ?? localQuery;
  const setQuery = (next: ResourceQuery) => {
    if (crud !== null) {
      crud.setTableQuery(tableId, next);
    } else {
      setLocalQuery(next);
    }
  };

  const [list, setList] = useState<ResourceList | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  // v2.9 ADR-0039: route snapshot for dataSource params bindings. Prefers the
  // provider's route context (App injects route: {params, query}); hostless
  // renders fall back to the location query.
  const routeSnapshot = useMemo(() => {
    if (crud !== null) {
      return crud.route.query;
    }
    const query: Record<string, string> = {};
    if (typeof window !== "undefined") {
      for (const [key, value] of new URLSearchParams(window.location.search).entries()) {
        query[key] = value;
      }
    }
    return query;
  }, [crud]);

  useEffect(() => {
    // F-001: absent/invalid dataSource fails closed — never fetch.
    if (dataSource === null) {
      setList(null);
      setLoading(false);
      setError(null);
      return;
    }
    let cancelled = false;
    setLoading(true);
    setError(null);
    const paramsQuery = resolveDataParamsQuery(dataParams, {
      query: routeSnapshot,
      params: crud !== null ? crud.route.params : {},
    });
    fetchResourceList(fetcher ?? fetch, dataSource, query, paramsQuery)
      .then((next) => {
        if (!cancelled) {
          setList(next);
          setLoading(false);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : String(err));
          setLoading(false);
        }
      });
    return () => {
      cancelled = true;
    };
  }, [fetcher, dataSource, dataParams, routeSnapshot, query, crud?.reloadToken]);

  // F-002: validate row keys on every fetched page; invalid → fail closed.
  const keyCheck = useMemo(
    () => (list === null ? null : checkRowKeys(list.items, rowKeyField)),
    [list, rowKeyField],
  );

  // --- Multi-select (ADR-0022 D2, I-PROTO-FULL-001 include) ---
  const selectionEnabled =
    typeof node.props?.selection === "object" &&
    node.props?.selection !== null &&
    (node.props.selection as Record<string, unknown>).mode === "multiple";

  // ADR-0022 D2: filter/page/sort changes clear the page selection.
  useEffect(() => {
    if (selectionEnabled && crud !== null) {
      crud.clearSelection(tableId);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [query.page, query.pageSize, query.sort, query.order, query.q]);

  if (columns.length === 0) {
    return (
      <p role="alert" className="text-sm text-destructive">
        table node requires a columns array
      </p>
    );
  }

  if (dataSource === null) {
    return (
      <p role="alert" className="text-sm text-destructive">
        table node requires a valid dataSource (single-slash same-origin path)
      </p>
    );
  }

  if (keyCheck !== null && !keyCheck.ok) {
    return (
      <p role="alert" className="text-sm text-destructive">
        {keyCheck.message}
      </p>
    );
  }

  const sort: SortState | undefined =
    query.sort !== undefined && query.order !== undefined
      ? { field: query.sort, order: query.order }
      : undefined;

  const onSortChange = (next: SortState) => {
    setQuery({ ...query, sort: next.field, order: next.order, page: 1 });
  };

  const onRowClick = (row: ResourceItem) => {
    crud?.selectRow(rowAsRecord(row));
  };

  const selectedRow = crud?.selectedRow;
  const selectedKey =
    selectedRow === null || selectedRow === undefined
      ? undefined
      : (scalarRowKey(selectedRow[rowKeyField]) ?? undefined);

  // --- Multi-select (ADR-0022 D2, I-PROTO-FULL-001 include) ---
  const currentSelection = selectionEnabled ? crud?.selection(tableId) : undefined;
  const selectedKeys = new Set(
    (currentSelection?.keys ?? []).map((key) => `${typeof key}:${String(key)}`),
  );

  const rowKeyToken = (row: ResourceItem): string | null => {
    const key = scalarRowKey(row[rowKeyField]);
    return key === null ? null : `${typeof row[rowKeyField]}:${key}`;
  };
  const selectionTokenOf = (token: string): unknown => {
    const colon = token.indexOf(":");
    const kind = token.slice(0, colon);
    const raw = token.slice(colon + 1);
    return kind === "number" ? Number(raw) : kind === "boolean" ? raw === "true" : raw;
  };

  const toggleRowSelection = (row: ResourceItem) => {
    const token = rowKeyToken(row);
    if (token === null || crud === null) {
      return;
    }
    const next = new Set(selectedKeys);
    if (next.has(token)) {
      next.delete(token);
    } else {
      next.add(token);
    }
    crud.setSelection(tableId, [...next].map(selectionTokenOf));
  };

  const allPageSelected =
    list !== null && keyCheck !== null && keyCheck.ok && list.items.length > 0
      ? list.items.every((row) => {
          const token = rowKeyToken(row);
          return token !== null && selectedKeys.has(token);
        })
      : false;
  const toggleAllPage = () => {
    if (crud === null || list === null || keyCheck === null || !keyCheck.ok) {
      return;
    }
    const next = new Set(selectedKeys);
    if (allPageSelected) {
      for (const row of list.items) {
        const token = rowKeyToken(row);
        if (token !== null) {
          next.delete(token);
        }
      }
    } else {
      for (const row of list.items) {
        const token = rowKeyToken(row);
        if (token !== null) {
          next.add(token);
        }
      }
    }
    crud.setSelection(tableId, [...next].map(selectionTokenOf));
  };

  const dataColumns: DataTableColumn<ResourceItem>[] = [
    ...(selectionEnabled
      ? [
          {
            key: "__selection",
            label: (
              <input
                type="checkbox"
                aria-label={t("feedback.selectAllOnPage")}
                checked={allPageSelected}
                disabled={list === null || keyCheck === null || !keyCheck.ok}
                onChange={toggleAllPage}
              />
            ),
            render: (row: ResourceItem) => {
              const token = rowKeyToken(row);
              return (
                <input
                  type="checkbox"
                  aria-label={t("feedback.selectRow")}
                  checked={token !== null && selectedKeys.has(token)}
                  disabled={token === null}
                  onClick={(event) => event.stopPropagation()}
                  onChange={() => toggleRowSelection(row)}
                />
              );
            },
          },
        ]
      : []),
    ...columns.map((column) => ({
      key: column.field,
      label: resolveTextProp(
        column as unknown as Record<string, unknown>,
        "labelKey",
        "label",
        t,
        column.field,
      ),
      sortable: column.sortable === true,
      truncate: column.truncate === true,
    })),
    ...(rowActions.length > 0
      ? [
          {
            key: "actions",
            label: "",
            render: (row: ResourceItem) => (
              <div
                className="flex justify-end gap-2"
                data-row-click-ignore="true"
                onClick={(event) => event.stopPropagation()}
              >
                {rowActions.map((action) => {
                  const key = stringOf(action.key) !== "" ? stringOf(action.key) : stringOf(action.actionRef);
                  const permitted = crud?.effectivePermission(key) ?? true;
                  const disabled = !permitted || rowActionDisabled(action, row);
                  return (
                    <button
                      key={key}
                      type="button"
                      disabled={disabled}
                      onClick={(event) => {
                        event.stopPropagation();
                        crud?.invokeAction(action, rowAsRecord(row));
                      }}
                      className="h-8 rounded-md border border-input bg-background px-2.5 text-xs font-medium text-muted-foreground shadow-sm transition-colors hover:bg-accent hover:text-foreground disabled:opacity-50"
                    >
                      {resolveTextProp(
                        action as unknown as Record<string, unknown>,
                        "labelKey",
                        "label",
                        t,
                        key,
                      )}
                    </button>
                  );
                })}
              </div>
            ),
          },
        ]
      : []),
  ];

  return (
    <div className="w-full min-w-0 space-y-2">
      {toolbar.length > 0 ? (
        <div className="flex flex-wrap items-center justify-end gap-2">
          {toolbar.map((trigger) => {
            const key = stringOf(trigger.key) !== "" ? stringOf(trigger.key) : stringOf(trigger.actionRef);
            const permitted = crud?.effectivePermission(key) ?? true;
            const isBatch =
              typeof trigger === "object" &&
              trigger !== null &&
              !Array.isArray(trigger) &&
              ((trigger as Record<string, unknown>).batchMapping !== undefined ||
                (trigger as Record<string, unknown>).requiresSelection === true);
            const requiresSelection = trigger.requiresSelection === true;
            const selectionDisabled =
              requiresSelection && (currentSelection === undefined || currentSelection.count === 0);
            const disabled = !permitted || selectionDisabled;
            return (
              <button
                key={key}
                type="button"
                disabled={disabled}
                onClick={() =>
                  isBatch && selectionEnabled
                    ? crud?.invokeBatchAction(trigger, tableId)
                    : crud?.invokeAction(trigger, null)
                }
                className="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground shadow-sm transition-opacity hover:opacity-90 disabled:opacity-50"
              >
                {resolveTextProp(
                  trigger as unknown as Record<string, unknown>,
                  "labelKey",
                  "label",
                  t,
                  key,
                )}
              </button>
            );
          })}
        </div>
      ) : null}
      <DataTable
        columns={dataColumns}
        rows={list?.items ?? []}
        rowKey={(row) => String(row[rowKeyField])}
        sort={sort}
        onSortChange={onSortChange}
        onRowClick={crud !== null ? onRowClick : undefined}
        selectedKey={selectedKey}
        loading={loading}
        error={error}
        emptyMessage={t("feedback.noItemsMatch")}
        caption={t("feedback.schemaDrivenItems")}
      />
      {list !== null ? (
        <p className="text-xs text-muted-foreground">
          {list.total} {list.total === 1 ? t("feedback.item") : t("feedback.items")} ·{" "}
          {t("feedback.pageOf", {
            page: String(list.page),
            total: String(Math.max(1, Math.ceil(list.total / list.pageSize))),
          })}
        </p>
      ) : null}
    </div>
  );
}
