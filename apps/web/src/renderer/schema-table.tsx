import { useEffect, useMemo, useState } from "react";

import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";
import {
  fetchRecords,
  isValidDataSource,
  type RecordsQuery,
  type ResourceItem,
  type ResourceList,
} from "@/renderer/records";
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
  /** Injectable fetch (defaults to `globalThis.fetch`); tests inject records). */
  fetcher?: typeof fetch;
}

export interface SchemaTableColumnSpec {
  field: string;
  label?: string;
  sortable?: boolean;
}

function stringOf(value: unknown): string {
  return typeof value === "string" ? value : "";
}

/** Spreads a typed row into a plain object the generic action executor accepts. */
function rowAsRecord(row: ResourceItem): Record<string, unknown> {
  return { ...row };
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
 * Resolves the table node's list endpoint (F-001). Returns null when absent or
 * not a single-slash same-origin path; the table then fails closed and never
 * fetches (the `/api/records` fallback was removed in GOAL-010 S3).
 */
export function schemaTableDataSource(node: RenderTableNode): string | null {
  const raw = node.props?.dataSource;
  return typeof raw === "string" && isValidDataSource(raw) ? raw : null;
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
  const rowKeyField = schemaTableRowKey(node);
  const crud = useSchemaCrud();
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
  const [localQuery, setLocalQuery] = useState<RecordsQuery>({ page: 1, pageSize: 10 });
  const query = providerQuery ?? localQuery;
  const setQuery = (next: RecordsQuery) => {
    if (crud !== null) {
      crud.setTableQuery(tableId, next);
    } else {
      setLocalQuery(next);
    }
  };

  const [list, setList] = useState<ResourceList | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

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
    fetchRecords(fetcher ?? fetch, dataSource, query)
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
  }, [fetcher, dataSource, query, crud?.reloadToken]);

  // F-002: validate row keys on every fetched page; invalid → fail closed.
  const keyCheck = useMemo(
    () => (list === null ? null : checkRowKeys(list.items, rowKeyField)),
    [list, rowKeyField],
  );

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

  const dataColumns: DataTableColumn<ResourceItem>[] = [
    ...columns.map((column) => ({
      key: column.field,
      label: column.label ?? column.field,
      sortable: column.sortable === true,
    })),
    ...(rowActions.length > 0
      ? [
          {
            key: "actions",
            label: "",
            render: (row: ResourceItem) => (
              <div className="flex justify-end gap-2">
                {rowActions.map((action) => {
                  const key = stringOf(action.key) !== "" ? stringOf(action.key) : stringOf(action.actionRef);
                  const permitted = crud?.effectivePermission(key) ?? true;
                  return (
                    <button
                      key={key}
                      type="button"
                      disabled={!permitted}
                      onClick={() => crud?.invokeAction(action, rowAsRecord(row))}
                      className="h-8 rounded-md border border-input bg-background px-2.5 text-xs text-muted-foreground transition-colors hover:bg-accent hover:text-foreground disabled:opacity-50"
                    >
                      {stringOf(action.label) ?? key}
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
    <div className="space-y-2">
      {toolbar.length > 0 ? (
        <div className="flex items-center justify-end">
          {toolbar.map((trigger) => {
            const key = stringOf(trigger.key) !== "" ? stringOf(trigger.key) : stringOf(trigger.actionRef);
            const permitted = crud?.effectivePermission(key) ?? true;
            return (
              <button
                key={key}
                type="button"
                disabled={!permitted}
                onClick={() => crud?.invokeAction(trigger, null)}
                className="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-50"
              >
                {stringOf(trigger.label) ?? key}
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
        emptyMessage="No records match."
        caption="Schema-driven records"
      />
      {list !== null ? (
        <p className="text-xs text-muted-foreground">
          {list.total} record{list.total === 1 ? "" : "s"} · page {list.page} of{" "}
          {Math.max(1, Math.ceil(list.total / list.pageSize))}
        </p>
      ) : null}
    </div>
  );
}
