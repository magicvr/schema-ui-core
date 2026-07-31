import { useMemo, useState, type FormEvent } from "react";

import { Button } from "@/components/ui/button";
import { DataTable, type DataTableColumn, type SortState } from "@/components/data-table";
import { FormControls } from "@/renderer/form-controls.tsx";
import type { FormControlField } from "@/renderer/form-controls";
import { runRowAction } from "@/renderer/row-action";
import {
  deleteRecord,
  updateRecord,
  type RecordItem,
  type RecordPatch,
} from "@/renderer/records";
import { useRecords } from "@/renderer/use-records";

const STATUS_OPTIONS = [
  { value: "active", label: "Active" },
  { value: "pending", label: "Pending" },
  { value: "archived", label: "Archived" },
];

const EDIT_FIELDS: FormControlField[] = [
  { id: "name", label: "Name", type: "input" },
  { id: "status", label: "Status", type: "select", options: STATUS_OPTIONS },
  { id: "owner", label: "Owner", type: "input" },
];

// Page document consumed by the R4 executeAction gate (D-ACT row actions).
const PAGE_DOCUMENT = {
  meta: {
    protocolVersion: "2.7",
    requiredCapabilities: ["app.manifest", "app.navigation", "permissions.inheritance"],
  },
  body: {
    type: "table",
    props: {
      columns: [{ field: "name" }],
      actions: [
        { key: "edit", label: "Edit", permissionIntent: "edit" },
        { key: "delete", label: "Delete", permissionIntent: "delete" },
      ],
    },
  },
} as const;

function actionGate(targetId: string, context: Record<string, unknown>) {
  return runRowAction({
    page: PAGE_DOCUMENT as unknown as Record<string, unknown>,
    targetId,
    context,
  });
}

function formatUpdatedAt(value: string): string {
  const date = new Date(value);
  if (Number.isNaN(date.getTime())) {
    return value;
  }
  return date.toLocaleDateString(undefined, {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

export function ListEditLifecyclePage() {
  const { list, loading, error, query, setQuery } = useRecords();
  const [editId, setEditId] = useState<string | null>(null);
  const [editValues, setEditValues] = useState<Record<string, unknown>>({});
  const [patchError, setPatchError] = useState<string | null>(null);
  const [confirmDeleteId, setConfirmDeleteId] = useState<string | null>(null);
  const [deleteError, setDeleteError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);

  // R4 D-PERM context: the dev session carries admin + editor roles.
  const context = useMemo(() => ({ user: { roles: ["admin"] }, features: {} }), []);
  const editGate = actionGate("edit", context);
  const deleteGate = actionGate("delete", context);

  const sort: SortState | undefined =
    query.sort !== undefined && query.order !== undefined
      ? { field: query.sort, order: query.order }
      : undefined;

  const onSortChange = (next: SortState) => {
    setQuery({ ...query, sort: next.field, order: next.order, page: 1 });
  };

  const startEdit = (record: RecordItem) => {
    setPatchError(null);
    setEditValues({ name: record.name, status: record.status, owner: record.owner });
    setEditId(record.id);
  };

  const submitEdit = async (event: FormEvent) => {
    event.preventDefault();
    if (editId === null) {
      return;
    }
    setPatchError(null);
    const patch: RecordPatch = {
      name: typeof editValues.name === "string" ? editValues.name : undefined,
      status: typeof editValues.status === "string" ? editValues.status : undefined,
      owner: typeof editValues.owner === "string" ? editValues.owner : undefined,
    };
    try {
      await updateRecord(fetch, "/api/records", editId, patch);
      setEditId(null);
      setNotice(`Updated ${editId}`);
      // Force a refresh of the current page.
      setQuery({ ...query });
    } catch (err) {
      setPatchError(err instanceof Error ? err.message : String(err));
    }
  };

  const requestDelete = (record: RecordItem) => {
    setDeleteError(null);
    setConfirmDeleteId(record.id);
  };

  const confirmDelete = async () => {
    if (confirmDeleteId === null) {
      return;
    }
    setDeleteError(null);
    try {
      await deleteRecord(fetch, "/api/records", confirmDeleteId);
      setConfirmDeleteId(null);
      setNotice(`Deleted ${confirmDeleteId}`);
      setQuery({ ...query });
    } catch (err) {
      setDeleteError(err instanceof Error ? err.message : String(err));
    }
  };

  const cancelDelete = () => setConfirmDeleteId(null);

  const columns: DataTableColumn<RecordItem>[] = [
    { key: "id", label: "ID" },
    { key: "name", label: "Name", sortable: true },
    {
      key: "status",
      label: "Status",
      sortable: true,
      render: (row) => <span>{row.status}</span>,
    },
    { key: "owner", label: "Owner", sortable: true },
    { key: "updatedAt", label: "Updated", render: (row) => formatUpdatedAt(row.updatedAt) },
    {
      key: "actions",
      label: "",
      render: (row) => (
        <div className="flex items-center gap-2">
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={editGate.outcome !== "EXECUTED"}
            onClick={() => startEdit(row)}
          >
            Edit
          </Button>
          <Button
            type="button"
            variant="outline"
            size="sm"
            disabled={deleteGate.outcome !== "EXECUTED"}
            onClick={() => requestDelete(row)}
          >
            Delete
          </Button>
        </div>
      ),
    },
  ];

  const activeRecord = list?.items.find((item) => item.id === editId);

  return (
    <section className="space-y-6" aria-labelledby="list-edit-title">
      <div className="space-y-2">
        <h1 id="list-edit-title" className="text-3xl font-semibold tracking-tight">
          List + edit lifecycle
        </h1>
        <p className="max-w-2xl text-sm leading-6 text-muted-foreground">
          Row actions gate through the R4 permission engine (D-ACT); the edit form
          uses whitelisted form controls (D-FORM) and writes back with PATCH.
        </p>
      </div>

      {notice !== null ? (
        <p role="status" className="text-sm text-foreground">
          {notice}
        </p>
      ) : null}

      {editId !== null && activeRecord !== undefined ? (
        <form onSubmit={submitEdit} className="space-y-4 border border-border bg-card p-4">
          <div className="flex items-center justify-between gap-2">
            <h2 className="text-sm font-semibold">Edit {editId}</h2>
            <Button type="button" variant="ghost" size="sm" onClick={() => setEditId(null)}>
              Cancel
            </Button>
          </div>
          <FormControls
            fields={EDIT_FIELDS}
            values={editValues}
            onChange={(id, value) => setEditValues((prev) => ({ ...prev, [id]: value }))}
          />
          {patchError !== null ? (
            <p role="alert" className="text-sm text-destructive">
              {patchError}
            </p>
          ) : null}
          <Button type="submit" size="sm">
            Save
          </Button>
        </form>
      ) : null}

      <DataTable
        columns={columns}
        rows={list?.items ?? []}
        rowKey={(row) => row.id}
        sort={sort}
        onSortChange={onSortChange}
        loading={loading}
        error={error}
        emptyMessage="No records match."
        caption="Editable records"
      />

      {confirmDeleteId !== null ? (
        <div
          role="dialog"
          aria-modal="true"
          aria-labelledby="delete-confirm-title"
          className="border border-destructive/50 bg-background p-4"
        >
          <h2 id="delete-confirm-title" className="text-sm font-semibold">
            Delete {confirmDeleteId}?
          </h2>
          <p className="mt-1 text-sm text-muted-foreground">
            This is an example D-ACT confirm gate. Deleting removes the record.
          </p>
          {deleteError !== null ? (
            <p role="alert" className="mt-2 text-sm text-destructive">
              {deleteError}
            </p>
          ) : null}
          <div className="mt-3 flex items-center gap-2">
            <Button type="button" variant="outline" size="sm" onClick={cancelDelete}>
              Cancel
            </Button>
            <Button type="button" size="sm" onClick={confirmDelete}>
              Confirm delete
            </Button>
          </div>
        </div>
      ) : null}
    </section>
  );
}
