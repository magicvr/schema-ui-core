import { useEffect, useMemo, useRef, useState } from "react";

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
  /** Column width hint (px number or CSS length). */
  width?: number | string;
  /** Minimum column width (px number or CSS length). */
  minWidth?: number | string;
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

/** One select option of a schema-driven table filter. */
export interface SchemaTableFilterOptionSpec {
  value: string;
  label?: string;
  labelKey?: string;
}

/** A schema-driven table filter (table node `props.filters`; select-only). */
export interface SchemaTableFilterSpec {
  field: string;
  type: "select";
  label?: string;
  labelKey?: string;
  options: SchemaTableFilterOptionSpec[];
}

function isFilterRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function filterOptionOf(value: unknown): SchemaTableFilterOptionSpec | null {
  if (!isFilterRecord(value) || typeof value.value !== "string") {
    return null;
  }
  return {
    value: value.value,
    ...(typeof value.label === "string" ? { label: value.label } : {}),
    ...(typeof value.labelKey === "string" ? { labelKey: value.labelKey } : {}),
  };
}

/**
 * Extracts the table node's filter specs (fail-closed on malformed entries).
 * Only `select` filters are supported; other/unknown types are dropped so a
 * schema typo never renders a broken control.
 */
export function schemaTableFilters(node: RenderTableNode): SchemaTableFilterSpec[] {
  const raw = node.props?.filters;
  if (!Array.isArray(raw)) {
    return [];
  }
  const out: SchemaTableFilterSpec[] = [];
  for (const entry of raw) {
    if (!isFilterRecord(entry)) {
      continue;
    }
    const field = stringOf(entry.field);
    if (field === "") {
      continue;
    }
    if (entry.type !== "select") {
      continue;
    }
    const options = (Array.isArray(entry.options) ? entry.options : [])
      .map(filterOptionOf)
      .filter((option): option is SchemaTableFilterOptionSpec => option !== null);
    out.push({
      field,
      type: "select",
      ...(typeof entry.label === "string" ? { label: entry.label } : {}),
      ...(typeof entry.labelKey === "string" ? { labelKey: entry.labelKey } : {}),
      options,
    });
  }
  return out;
}

/**
 * Windowed pager page list: 1 … current±1 … total, with gap markers for
 * long page runs (stable shape for tests and a11y).
 */
export function pagerPages(current: number, total: number): Array<number | "gap"> {
  const pages = Math.max(1, Math.floor(total));
  const now = Math.min(Math.max(1, Math.floor(current)), pages);
  if (pages <= 7) {
    return Array.from({ length: pages }, (_, index) => index + 1);
  }
  const out: Array<number | "gap"> = [1];
  if (now > 3) {
    out.push("gap");
  }
  const start = Math.max(2, now - 1);
  const end = Math.min(pages - 1, now + 1);
  for (let page = start; page <= end; page += 1) {
    out.push(page);
  }
  if (now < pages - 2) {
    out.push("gap");
  }
  out.push(pages);
  return out;
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

/**
 * W11 · U-05: actions beyond the primary two collapse into a "⋯ More" menu
 * (one open row at a time; outside click / Escape closes). Keeps dense tables
 * scannable and reduces accidental taps on destructive row actions.
 */
const MAX_INLINE_ROW_ACTIONS = 2;

function RowActionsMenu({
  row,
  actions,
  crud,
  t,
  rowKeyField,
  disabledWhen,
}: {
  row: ResourceItem;
  actions: Array<Record<string, unknown>>;
  crud: ReturnType<typeof useSchemaCrud>;
  t: ReturnType<typeof useTranslate>;
  rowKeyField: string;
  disabledWhen: (action: unknown, row: ResourceItem) => boolean;
}) {
  const [open, setOpen] = useState(false);
  const rootRef = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (!open) {
      return;
    }
    const onPointerDown = (event: PointerEvent) => {
      if (rootRef.current !== null && !rootRef.current.contains(event.target as Node)) {
        setOpen(false);
      }
    };
    const onKeyDown = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
      }
    };
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    return () => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    };
  }, [open]);

  const rowKey = scalarRowKey(row[rowKeyField]) ?? "row";
  return (
    <div ref={rootRef} className="relative" data-row-actions-menu={rowKey}>
      <button
        type="button"
        aria-expanded={open}
        aria-label={t("feedback.moreActions")}
        className="flex h-7 w-7 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
        onClick={(event) => {
          event.stopPropagation();
          setOpen((value) => !value);
        }}
      >
        {"⋯"}
      </button>
      {open ? (
        <div
          role="menu"
          aria-label={t("feedback.moreActions")}
          className="absolute right-0 top-8 z-30 min-w-36 rounded-md border border-border bg-card py-1 shadow-lg"
        >
          {actions.map((action) => {
            const key = stringOf(action.key) !== "" ? stringOf(action.key) : stringOf(action.actionRef);
            const permitted = crud?.effectivePermission(key) ?? true;
            const disabled = !permitted || disabledWhen(action, row);
            return (
              <button
                key={key}
                type="button"
                role="menuitem"
                disabled={disabled}
                onClick={(event) => {
                  event.stopPropagation();
                  setOpen(false);
                  crud?.invokeAction(action, rowAsRecord(row));
                }}
                className="block w-full px-3 py-1.5 text-left text-xs text-foreground transition-colors hover:bg-accent disabled:opacity-50"
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
      ) : null}
    </div>
  );
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
  const filters = schemaTableFilters(node);
  const title = resolveTextProp(
    node.props as unknown as Record<string, unknown>,
    "titleKey",
    "title",
    t,
    "",
  );

  // Register the injected transport with the page's Schema CRUD provider so
  // modal form submits and row actions share the same fetcher (S4).
  useEffect(() => {
    if (crud !== null && fetcher !== undefined) {
      crud.registerFetcher(fetcher);
    }
  }, [crud, fetcher]);

  // W11 · U-06: go-to-page input ref (submit reads it; no controlled state).
  const goToPageRef = useRef<HTMLInputElement>(null);
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
    // Keep the current rows rendered while refetching (page/filter/sort
    // changes): swapping rows in place keeps the list height and the browser
    // scroll anchor stable, so paginating never yanks the viewport back to
    // the top. The loading skeleton only appears on the first load (no rows
    // yet). (GOAL-011 W10)
    if (list === null) {
      setLoading(true);
    }
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
  }, [query.page, query.pageSize, query.sort, query.order, query.q, query.filters]);

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
      ...(column.width !== undefined ? { width: column.width } : {}),
      ...(column.minWidth !== undefined ? { minWidth: column.minWidth } : {}),
    })),
    ...(rowActions.length > 0
      ? [
          {
            key: "actions",
            label: "",
            render: (row: ResourceItem) => {
              const primary = rowActions.slice(0, MAX_INLINE_ROW_ACTIONS);
              const overflow = rowActions.slice(MAX_INLINE_ROW_ACTIONS);
              return (
                <div
                  className="flex items-center justify-end gap-1"
                  data-row-click-ignore="true"
                  onClick={(event) => event.stopPropagation()}
                >
                  {primary.map((action) => {
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
                        className="rounded px-2 py-1 text-xs font-medium text-muted-foreground transition-colors hover:bg-accent hover:text-foreground disabled:opacity-50"
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
                  {overflow.length > 0 ? (
                    <RowActionsMenu
                      row={row}
                      actions={overflow}
                      crud={crud}
                      t={t}
                      rowKeyField={rowKeyField}
                      disabledWhen={rowActionDisabled}
                    />
                  ) : null}
                </div>
              );
            },
          },
        ]
      : []),
  ];

  return (
    <div className="w-full min-w-0 space-y-2">
      {title !== "" ? (
        <h2 className="text-lg font-semibold tracking-tight text-foreground">{title}</h2>
      ) : null}
      {filters.length > 0 ? (
        <div className="flex flex-wrap items-center gap-4" data-table-filters>
          {filters.map((filter) => {
            const label = resolveTextProp(
              filter as unknown as Record<string, unknown>,
              "labelKey",
              "label",
              t,
              filter.field,
            );
            const value = query.filters?.[filter.field] ?? "";
            return (
              <label
                key={filter.field}
                className="flex items-center gap-2 text-sm text-muted-foreground"
              >
                <span>{label}</span>
                <select
                  value={value}
                  onChange={(event) =>
                    setQuery({
                      ...query,
                      filters: { ...(query.filters ?? {}), [filter.field]: event.target.value },
                      page: 1,
                    })
                  }
                  className="h-8 rounded-md border border-input bg-background px-2 text-sm"
                >
                  {filter.options.map((option) => (
                    <option key={option.value} value={option.value}>
                      {resolveTextProp(
                        option as unknown as Record<string, unknown>,
                        "labelKey",
                        "label",
                        t,
                        option.value,
                      )}
                    </option>
                  ))}
                </select>
              </label>
            );
          })}
        </div>
      ) : null}
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
        <div className="flex flex-wrap items-center justify-between gap-2">
          <div className="flex flex-wrap items-center gap-3">
            <p className="pl-0.5 text-xs text-muted-foreground">
              {list.total} {list.total === 1 ? t("feedback.item") : t("feedback.items")} ·{" "}
              {t("feedback.pageOf", {
                page: String(list.page),
                total: String(Math.max(1, Math.ceil(list.total / list.pageSize))),
              })}
            </p>
            {/* W11 · U-06: per-page size switcher (resets to page 1). */}
            <label className="flex items-center gap-1.5 text-xs text-muted-foreground">
              <span>{t("feedback.pageSize")}</span>
              <select
                aria-label={t("feedback.pageSize")}
                value={String(query.pageSize ?? 10)}
                onChange={(event) =>
                  setQuery({ ...query, pageSize: Number(event.target.value), page: 1 })
                }
                className="h-7 rounded-md border border-input bg-background px-1.5 text-xs scheme-light dark:scheme-dark"
              >
                {[10, 20, 50, 100].map((size) => (
                  <option key={size} value={String(size)}>
                    {size}
                  </option>
                ))}
              </select>
            </label>
          </div>
          {Math.max(1, Math.ceil(list.total / list.pageSize)) > 1 ? (
            <nav aria-label={t("feedback.pagination")} className="flex items-center gap-1">
              <button
                type="button"
                disabled={list.page <= 1}
                aria-label={t("feedback.previousPage")}
                onClick={() => setQuery({ ...query, page: list.page - 1 })}
                className="flex h-7 w-7 items-center justify-center rounded-md border border-input bg-background text-sm text-muted-foreground shadow-sm transition-colors hover:bg-accent hover:text-foreground disabled:opacity-40"
              >
                {"‹"}
              </button>
              {pagerPages(list.page, Math.max(1, Math.ceil(list.total / list.pageSize))).map(
                (page, index) =>
                  page === "gap" ? (
                    <span key={"gap-" + String(index)} className="px-1 text-xs text-muted-foreground">
                      {"…"}
                    </span>
                  ) : (
                    <button
                      key={page}
                      type="button"
                      disabled={page === list.page}
                      aria-current={page === list.page ? "page" : undefined}
                      aria-label={t("feedback.pageNumber", { page: String(page) })}
                      onClick={() => setQuery({ ...query, page })}
                      className="flex h-7 w-7 items-center justify-center rounded-md border border-input bg-background text-sm shadow-sm transition-colors hover:bg-accent hover:text-foreground disabled:cursor-default disabled:border-primary/40 disabled:bg-primary/10 disabled:text-foreground"
                    >
                      {page}
                    </button>
                  ),
              )}
              <button
                type="button"
                disabled={list.page >= Math.max(1, Math.ceil(list.total / list.pageSize))}
                aria-label={t("feedback.nextPage")}
                onClick={() => setQuery({ ...query, page: list.page + 1 })}
                className="flex h-7 w-7 items-center justify-center rounded-md border border-input bg-background text-sm text-muted-foreground shadow-sm transition-colors hover:bg-accent hover:text-foreground disabled:opacity-40"
              >
                {"›"}
              </button>
            </nav>
          ) : null}
          {/* W11 · U-06: quick jump to a specific page. */}
          {Math.max(1, Math.ceil(list.total / list.pageSize)) > 1 ? (
            <form
              aria-label={t("feedback.goToPage")}
              className="flex items-center gap-1.5"
              onSubmit={(event) => {
                event.preventDefault();
                const target = Number(goToPageRef.current?.value ?? "");
                const pages = Math.max(1, Math.ceil(list.total / list.pageSize));
                if (Number.isFinite(target) && target >= 1 && target <= pages) {
                  setQuery({ ...query, page: Math.floor(target) });
                }
              }}
            >
              <label className="text-xs text-muted-foreground">{t("feedback.goToPage")}</label>
              <input
                ref={goToPageRef}
                type="number"
                min={1}
                max={Math.max(1, Math.ceil(list.total / list.pageSize))}
                defaultValue=""
                placeholder={String(list.page)}
                aria-label={t("feedback.goToPage")}
                className="h-7 w-16 rounded-md border border-input bg-background px-1.5 text-xs"
              />
              <button
                type="submit"
                className="h-7 rounded-md border border-input bg-background px-2 text-xs text-muted-foreground shadow-sm transition-colors hover:bg-accent hover:text-foreground"
              >
                {t("feedback.search")}
              </button>
            </form>
          ) : null}
        </div>
      ) : null}
    </div>
  );
}
