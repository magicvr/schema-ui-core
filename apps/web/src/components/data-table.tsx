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

function labelText(label: ReactNode): string {
  if (typeof label === "string" || typeof label === "number") {
    return String(label);
  }
  return "";
}

/**
 * Columns eligible for the mobile card surface. Action / selection chrome stays
 * out of the primary title/secondary stack and is rendered as trailing controls.
 */
function contentColumns<T>(columns: DataTableColumn<T>[]): DataTableColumn<T>[] {
  return columns.filter((column) => column.key !== "__selection" && column.key !== "actions");
}

function actionColumn<T>(columns: DataTableColumn<T>[]): DataTableColumn<T> | undefined {
  return columns.find((column) => column.key === "actions");
}

function MobileCardList<T>({
  columns,
  rows,
  rowKey,
  onRowClick,
  selectedKey,
}: {
  columns: DataTableColumn<T>[];
  rows: T[];
  rowKey: (row: T) => string;
  onRowClick?: (row: T) => void;
  selectedKey?: string;
}) {
  const fields = contentColumns(columns);
  const actions = actionColumn(columns);
  const titleColumn = fields[0];
  const secondaryColumns = fields.slice(1, 3);

  return (
    <ul
      data-table-presentation="mobile-cards"
      className="space-y-2 md:hidden"
      aria-label="Mobile card list"
    >
      {rows.map((row) => {
        const key = rowKey(row);
        const selected = selectedKey !== undefined && selectedKey === key;
        return (
          <li key={key}>
            <div
              role={onRowClick === undefined ? undefined : "button"}
              tabIndex={onRowClick === undefined ? undefined : 0}
              aria-selected={onRowClick === undefined ? undefined : selected}
              onClick={onRowClick === undefined ? undefined : () => onRowClick(row)}
              onKeyDown={
                onRowClick === undefined
                  ? undefined
                  : (event) => {
                      if (event.key === "Enter" || event.key === " ") {
                        event.preventDefault();
                        onRowClick(row);
                      }
                    }
              }
              className={cn(
                "rounded-lg border border-border bg-card p-3 text-left shadow-sm transition-colors",
                onRowClick === undefined ? "" : "cursor-pointer hover:bg-accent/40",
                selected ? "border-primary/40 bg-accent/50 ring-1 ring-primary/20" : "",
              )}
            >
              <div className="flex items-start justify-between gap-3">
                <div className="min-w-0 flex-1 space-y-1">
                  {titleColumn !== undefined ? (
                    <p className="truncate text-sm font-semibold text-foreground">
                      {cellContent(titleColumn, row)}
                    </p>
                  ) : null}
                  {secondaryColumns.map((column) => (
                    <p
                      key={column.key}
                      className="truncate text-xs text-muted-foreground"
                    >
                      {labelText(column.label) !== "" ? (
                        <span className="mr-1 font-medium text-muted-foreground/80">
                          {labelText(column.label)}:
                        </span>
                      ) : null}
                      {cellContent(column, row)}
                    </p>
                  ))}
                </div>
                {actions !== undefined ? (
                  <div
                    className="shrink-0"
                    onClick={(event) => event.stopPropagation()}
                    onKeyDown={(event) => event.stopPropagation()}
                  >
                    {actions.render !== undefined
                      ? actions.render(row)
                      : cellContent(actions, row)}
                  </div>
                ) : (
                  <span
                    aria-hidden="true"
                    className="shrink-0 px-1 text-sm text-muted-foreground"
                  >
                    ⋯
                  </span>
                )}
              </div>
            </div>
          </li>
        );
      })}
    </ul>
  );
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

  const state = resolveAsyncDisplayState({ loading, error, isEmpty: rows.length === 0 });

  if (state === "loading") {
    return (
      <div className="space-y-2" data-table-presentation="loading">
        <div role="status" aria-label="Loading" className="space-y-2 rounded-md border border-border p-4">
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-full" />
          <Skeleton className="h-4 w-3/4" />
        </div>
      </div>
    );
  }

  if (state === "error") {
    return (
      <p role="alert" className="rounded-md border border-destructive/30 bg-destructive/5 px-4 py-6 text-sm text-destructive">
        {error}
      </p>
    );
  }

  if (state === "empty") {
    return (
      <p className="rounded-md border border-border bg-card px-4 py-6 text-sm text-muted-foreground">
        {emptyMessage}
      </p>
    );
  }

  return (
    <div className="w-full min-w-0 space-y-0" data-table-presentation="dual-end">
      {/* Desktop / tablet dense table (D-004 §4): hidden below md; scrolls within viewport */}
      <div
        data-table-presentation="desktop-table"
        className="hidden w-full min-w-0 overflow-x-auto rounded-md border border-border md:block"
      >
        <table className="w-full min-w-[32rem] border-collapse text-sm">
          {caption ? <caption className="sr-only">{caption}</caption> : null}
          <thead>
            <tr className="border-b border-border bg-muted/40">
              {columns.map((column) => {
                const isActive = sort?.field === column.key;
                return (
                  <th
                    key={column.key}
                    className="px-3 py-2 text-left text-[11px] font-semibold uppercase tracking-[0.12em] text-muted-foreground"
                    scope="col"
                  >
                    {column.sortable ? (
                      <button
                        type="button"
                        onClick={() => toggleSort(column)}
                        aria-sort={
                          isActive
                            ? sort.order === "asc"
                              ? "ascending"
                              : "descending"
                            : undefined
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
            {rows.map((row) => {
              const key = rowKey(row);
              const selected = selectedKey !== undefined && selectedKey === key;
              return (
                <tr
                  key={key}
                  onClick={
                    onRowClick === undefined
                      ? undefined
                      : (event) => {
                          // Action/selection chrome must not select the row
                          // (selecting opens recordView drawer — bad UX for Edit/Delete).
                          const target = event.target as HTMLElement | null;
                          if (
                            target?.closest(
                              "button, a, input, select, textarea, label, [data-row-click-ignore]",
                            )
                          ) {
                            return;
                          }
                          onRowClick(row);
                        }
                  }
                  aria-selected={onRowClick === undefined ? undefined : selected}
                  className={cn(
                    "border-b border-border last:border-b-0",
                    onRowClick === undefined
                      ? ""
                      : "cursor-pointer transition-colors hover:bg-accent/50",
                    selected ? "bg-accent/60" : "",
                  )}
                >
                  {columns.map((column) => {
                    const interactive =
                      column.key === "actions" || column.key === "__selection";
                    return (
                      <td
                        key={column.key}
                        data-row-click-ignore={interactive ? "true" : undefined}
                        className="px-3 py-2 align-middle text-sm"
                        onClick={
                          interactive
                            ? (event) => event.stopPropagation()
                            : undefined
                        }
                      >
                        {cellContent(column, row)}
                      </td>
                    );
                  })}
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>

      {/* Mobile card list (D-004 §4): visible only below md */}
      <MobileCardList
        columns={columns}
        rows={rows}
        rowKey={rowKey}
        onRowClick={onRowClick}
        selectedKey={selectedKey}
      />
    </div>
  );
}
