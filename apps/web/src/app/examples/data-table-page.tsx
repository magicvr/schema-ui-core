import { useRecords } from "@/renderer/use-records";
import type { RecordItem } from "@/renderer/records";
import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";

const COLUMNS: DataTableColumn<RecordItem>[] = [
  { key: "id", label: "ID" },
  { key: "name", label: "Name", sortable: true },
  { key: "status", label: "Status", sortable: true },
  { key: "owner", label: "Owner", sortable: true },
  { key: "updatedAt", label: "Updated", sortable: true },
];

export function DataTablePage() {
  const { list, loading, error, query, setQuery } = useRecords();
  const sort: SortState | undefined =
    query.sort !== undefined && query.order !== undefined
      ? { field: query.sort, order: query.order }
      : undefined;

  const onSortChange = (next: SortState) => {
    setQuery({ ...query, sort: next.field, order: next.order, page: 1 });
  };

  return (
    <section className="space-y-6" aria-labelledby="data-table-title">
      <div className="space-y-2">
        <h1 id="data-table-title" className="text-3xl font-semibold tracking-tight">
          Data table
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
          Sortable records from the Go list API (<code>/api/records</code>). This
          example exercises D-DATA list/detail and D-TABLE sort declarations.
        </p>
      </div>
      <DataTable
        columns={COLUMNS}
        rows={list?.items ?? []}
        rowKey={(row) => row.id}
        sort={sort}
        onSortChange={onSortChange}
        loading={loading}
        error={error}
        emptyMessage="No records match."
        caption="Example records"
      />
      {list !== null ? (
        <p className="text-xs text-muted-foreground">
          {list.total} record{list.total === 1 ? "" : "s"} · page {list.page} of{" "}
          {Math.max(1, Math.ceil(list.total / list.pageSize))}
        </p>
      ) : null}
    </section>
  );
}
