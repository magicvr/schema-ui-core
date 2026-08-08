import type { ReactNode } from "react";

import { resolveAsyncDisplayState } from "@/components/ui/async-state";
import { Skeleton } from "@/components/ui/skeleton";
import { cn } from "@/lib/utils";

export type SortOrder = "asc" | "desc";

export interface SortState {
  field: string;
  order: SortOrder;
}

export interface DataTableColumn<T> {
  key: string;
  /** Cell content or header node (checkboxes render in headers for selection). */
  label: ReactNode;
  sortable?: boolean;
  render?: (row: T) => ReactNode;
}

export interface DataTableProps<T> {
  columns: DataTableColumn<T>[];
  rows: T[];
  rowKey: (row: T) => string;
  sort?: SortState;
  onSortChange?: (sort: SortState) => void;
  loading?: boolean;
  error?: string | null;
  emptyMessage?: string;
  caption?: string;
  /** Invoked when a data row is clicked (row selection, S4 · GOAL-007). */
  onRowClick?: (row: T) => void;
  /** Row key of the currently selected row (highlight), S4 · GOAL-007. */
  selectedKey?: string;
}

function cellContent<T>(
  column: DataTableColumn<T>,
  row: T,
): ReactNode {
  if (column.render !== undefined) {
    return column.render(row);
  }
  const value = (row as Record<string, unknown>)[column.key];
  if (value === undefined || value === null) {
    return <span className="text-muted-foreground">—</span>;
  }
  return String(value);
}

export function DataTable<T>({
  columns,
  rows,
  rowKey,
  sort,
  onSortChange,
  loading = false,
  error = null,
  emptyMessage = "No rows.",
  caption,
  onRowClick,
  selectedKey,
}: DataTableProps<T>) {
  const toggleSort = (column: DataTableColumn<T>) => {
    if (!column.sortable || onSortChange === undefined) {
      return;
    }
    const nextOrder: SortOrder =
      sort?.field === column.key && sort.order === "asc" ? "desc" : "asc";
    onSortChange({ field: column.key, order: nextOrder });
  };

  const arrowFor = (column: DataTableColumn<T>): string => {
    if (sort?.field !== column.key) {
      return "";
    }
    return sort.order === "asc" ? " ↑" : " ↓";
  };

  return (
    <div className="overflow-hidden border border-border">
      <table className="w-full border-collapse text-sm">
        {caption ? <caption className="sr-only">{caption}</caption> : null}
        <thead>
          <tr className="border-b border-border bg-muted/50">
            {columns.map((column) => {
              const isActive = sort?.field === column.key;
              return (
                <th
                  key={column.key}
                  className="px-4 py-3 text-left text-xs font-semibold uppercase tracking-[0.12em] text-muted-foreground"
                  scope="col"
                >
                  {column.sortable ? (
                    <button
                      type="button"
                      onClick={() => toggleSort(column)}
                      aria-sort={
                        isActive ? (sort.order === "asc" ? "ascending" : "descending") : undefined
                      }
                      className={cn(
                        "inline-flex items-center gap-1 uppercase tracking-[0.12em]",
                        isActive ? "text-foreground" : "hover:text-foreground",
                      )}
                    >
                      {column.label}
                      <span aria-hidden="true" className="text-[10px]">
                        {arrowFor(column)}
                      </span>
                    </button>
                  ) : (
                    column.label
                  )}
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody>
          {(() => {
            const state = resolveAsyncDisplayState({ loading, error, isEmpty: rows.length === 0 });
            if (state === "loading") {
              return (
                <tr>
                  <td colSpan={columns.length} className="px-4 py-6">
                    <div role="status" aria-label="Loading" className="space-y-2">
                      <Skeleton className="h-4 w-full" />
                      <Skeleton className="h-4 w-full" />
                      <Skeleton className="h-4 w-3/4" />
                    </div>
                  </td>
                </tr>
              );
            }
            if (state === "error") {
              return (
                <tr>
                  <td
                    colSpan={columns.length}
                    role="alert"
                    className="px-4 py-6 text-sm text-destructive"
                  >
                    {error}
                  </td>
                </tr>
              );
            }
            if (state === "empty") {
              return (
                <tr>
                  <td colSpan={columns.length} className="px-4 py-6 text-sm text-muted-foreground">
                    {emptyMessage}
                  </td>
                </tr>
              );
            }
            return rows.map((row) => {
              const key = rowKey(row);
              const selected = selectedKey !== undefined && selectedKey === key;
              return (
                <tr
                  key={key}
                  onClick={onRowClick === undefined ? undefined : () => onRowClick(row)}
                  aria-selected={onRowClick === undefined ? undefined : selected}
                  className={cn(
                    "border-b border-border last:border-b-0",
                    onRowClick === undefined
                      ? ""
                      : "cursor-pointer transition-colors hover:bg-accent",
                    selected ? "bg-accent/60" : "",
                  )}
                >
                  {columns.map((column) => (
                    <td key={column.key} className="px-4 py-3 align-middle">
                      {cellContent(column, row)}
                    </td>
                  ))}
                </tr>
              );
            });
          })()}
        </tbody>
      </table>
    </div>
  );
}
