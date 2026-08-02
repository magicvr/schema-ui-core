import { useEffect, useState } from "react";

import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";
import {
  DEFAULT_RECORDS_URL,
  fetchRecords,
  type RecordList,
  type RecordItem,
  type RecordsQuery,
} from "@/renderer/records";
import type { RenderTableNode } from "@/renderer/render";
import { useSchemaCrud } from "@/renderer/render.tsx";

/**
 * Default schema-driven table surface (R1 · GOAL-004 / D-004) + S4 CRUD wiring.
 *
 * Renders a whitelisted table node's `props.columns` over its `props.dataSource`
 * (reuse of the demo `GET /api/records` D-DATA contract). Provides loading /
 * error / empty states and column sort, mirroring the hand-written example's
 * surface without owning data logic. Fails closed when the node declares no
 * columns or no data source.
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

type JsonRecord = Record<string, unknown>;

function stringOf(value: unknown): string {
  return typeof value === "string" ? value : "";
}

/** Spreads a typed row into a plain object the generic action executor accepts. */
function rowAsRecord(row: RecordItem): JsonRecord {
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

/** Resolves the table node's data source URL, defaulting to the demo API. */
export function schemaTableDataSource(node: RenderTableNode): string {
  return typeof node.props?.dataSource === "string" && node.props.dataSource !== ""
    ? node.props.dataSource
    : DEFAULT_RECORDS_URL;
}

export function SchemaTable({ node, fetcher }: SchemaTableProps) {
  const columns = schemaTableColumns(node);
  const dataSource = schemaTableDataSource(node);
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

  const [list, setList] = useState<RecordList | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
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

  if (columns.length === 0) {
    return (
      <p role="alert" className="text-sm text-destructive">
        table node requires a columns array
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

  const onRowClick = (row: RecordItem) => {
    crud?.selectRow(rowAsRecord(row));
  };

  const dataColumns: DataTableColumn<RecordItem>[] = [
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
            render: (row: RecordItem) => (
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
        rowKey={(row) => row.id}
        sort={sort}
        onSortChange={onSortChange}
        onRowClick={crud !== null ? onRowClick : undefined}
        selectedKey={crud?.selectedRow?.id as string | undefined}
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
