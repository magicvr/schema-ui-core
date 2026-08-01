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

/**
 * Default schema-driven table surface (R1 · GOAL-004 / D-004).
 *
 * Renders a whitelisted table node's `props.columns` over its `props.dataSource`
 * (reuse of the demo `GET /api/records` D-DATA contract). Provides loading /
 * error / empty states and column sort, mirroring the hand-written example's
 * surface without owning data logic. Fails closed when the node declares no
 * columns or no data source.
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
  const [query, setQuery] = useState<RecordsQuery>({ page: 1, pageSize: 10 });
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
  }, [fetcher, dataSource, query]);

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

  const dataColumns: DataTableColumn<RecordItem>[] = columns.map((column) => ({
    key: column.field,
    label: column.label ?? column.field,
    sortable: column.sortable === true,
  }));

  return (
    <div className="space-y-2">
      <DataTable
        columns={dataColumns}
        rows={list?.items ?? []}
        rowKey={(row) => row.id}
        sort={sort}
        onSortChange={onSortChange}
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
