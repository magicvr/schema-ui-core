import { useState, type FormEvent } from "react";

import { Button } from "@/components/ui/button";
import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";
import { useRecords } from "@/renderer/use-records";
import type { RecordItem } from "@/renderer/records";

const COLUMNS: DataTableColumn<RecordItem>[] = [
  { key: "id", label: "ID" },
  { key: "name", label: "Name", sortable: true },
  { key: "status", label: "Status", sortable: true },
  { key: "owner", label: "Owner", sortable: true },
  { key: "updatedAt", label: "Updated", sortable: true },
];

export function SearchFormTablePage() {
  const { list, loading, error, query, setQuery } = useRecords();
  const [draft, setDraft] = useState("");

  const sort: SortState | undefined =
    query.sort !== undefined && query.order !== undefined
      ? { field: query.sort, order: query.order }
      : undefined;

  const onSortChange = (next: SortState) => {
    setQuery({ ...query, sort: next.field, order: next.order, page: 1 });
  };

  const onSubmit = (event: FormEvent) => {
    event.preventDefault();
    setQuery({ ...query, q: draft.trim(), page: 1 });
  };

  return (
    <section className="space-y-6" aria-labelledby="search-table-title">
      <div className="space-y-2">
        <h1 id="search-table-title" className="text-3xl font-semibold tracking-tight">
          Search + table
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
          A search form filters the record list (D-TABLE <code>search-table</code>);
          column headers declare sort (D-TABLE <code>table-sort</code>).
        </p>
      </div>
      <form onSubmit={onSubmit} className="flex max-w-md items-center gap-2">
        <label className="sr-only" htmlFor="search-input">
          Search records
        </label>
        <input
          id="search-input"
          type="search"
          value={draft}
          onChange={(event) => setDraft(event.target.value)}
          placeholder="Search name, status, or owner…"
          className="h-9 flex-1 rounded-md border border-input bg-background px-3 text-sm outline-none transition-colors focus-visible:ring-2 focus-visible:ring-ring"
        />
        <Button type="submit" variant="secondary">
          Search
        </Button>
      </form>
      <DataTable
        columns={COLUMNS}
        rows={list?.items ?? []}
        rowKey={(row) => row.id}
        sort={sort}
        onSortChange={onSortChange}
        loading={loading}
        error={error}
        emptyMessage="No records match."
        caption="Search results"
      />
      {list !== null ? (
        <p className="text-xs text-muted-foreground">
          {list.total} record{list.total === 1 ? "" : "s"}
        </p>
      ) : null}
    </section>
  );
}
