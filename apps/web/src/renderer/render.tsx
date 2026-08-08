import {
  Fragment,
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useState,
  type ComponentType,
  type ReactNode,
} from "react";

import type { NavigationContext } from "@/protocol/app-manifest";
import { applyComponentFormat } from "@/protocol/conformance/component-format";
import { constructRequest } from "@/protocol/conformance/request-construction";
import { ConfirmDialog } from "@/renderer/confirm";
import { FormControls } from "@/renderer/form-controls.tsx";
import { coerceFieldValue, type FormControlField } from "@/renderer/form-controls";
import { ModalHost } from "@/renderer/modal";
import {
  executeAction,
  evaluatePermissionTargets,
} from "@/renderer/permissions";
import {
  fetchResourceList,
  isValidDataSource,
  readResourceApiError,
  type ResourceList,
  type ResourceQuery,
} from "@/renderer/resource";
import {
  gateRenderFormFields,
  parseRenderNode,
  resolveFormReactions,
  tableActionGate,
  type RenderActionButtonNode,
  type RenderChartNode,
  type RenderFormNode,
  type RenderGridNode,
  type RenderNode,
  type RenderPageDocument,
  type RenderRecordViewNode,
  type RenderSectionNode,
  type RenderStatCardNode,
  type RenderTabsNode,
  type RenderTableNode,
  type RenderTextNode,
} from "@/renderer/render";

/**
 * R5 D-COMP minimal Renderer (resolve R4 F-002) + S4 Schema CRUD (GOAL-007).
 *
 * Dispatch layer: parses a page document, applies the frozen $context
 * reaction engine to form field state, and renders whitelisted node types
 * through the components in this directory. Unknown node types fail closed.
 *
 * S4 one-time completion (I-007-003 v0.2.2 §9): a SchemaCrudProvider owns the
 * cross-node state a Schema-driven CRUD page needs — selected row (feeds
 * recordView + edit-form prefill), per-table query (search form-to-query
 * binding), a reload token, the active modal, the pending delete confirm, and
 * a generic action executor that gates through the frozen `executeAction`
 * engine and constructs requests with the pinned conformance constructor
 * (`request-construction.ts`). Every page-level behaviour stays fixture-driven;
 * after this completion page behavior remains schema-owned; core fixtures and
 * module-owned schema packages can add pages without Renderer changes.
 *
 * Scope: the frozen §5 node whitelist — layout (grid/section/tabs),
 * data/action (text/table/recordView/actionButton) and form. The form control
 * whitelist itself is enforced by D-FORM (isWhitelistedFormControl /
 * checkFormCapabilities) via gateRenderFormFields. The default app path wires
 * a schema-driven table surface (SchemaTable, GOAL-004) as `tableRenderer`;
 * a table node dispatched without one fails closed with an observable note.
 *
 * A-002 F-002-002 (GOAL-009 S1): form submission is blocked while any
 * gate/reaction error is present — the submit button is disabled and
 * handleSubmit re-rejects before any request can be constructed.
 */

export interface RendererComponentProps {
  document: RenderPageDocument;
  context: Record<string, unknown>;
  /** Renders a table node; provided by the example page that owns data. */
  tableRenderer?: (node: RenderTableNode) => ReactNode;
  /** Renders statCard/chart nodes (supportsData display, registry); defaults to built-ins. */
  dataRenderer?: (node: RenderStatCardNode | RenderChartNode) => ReactNode;
  /** Invoked when an actionButton node is activated. */
  onAction?: (node: RenderActionButtonNode) => void;
  /** Overrides the default FormControls component (keeps field wiring local). */
  formComponent?: ComponentType<{
    fields: FormControlField[];
    values: Record<string, unknown>;
    onChange: (id: string, value: unknown) => void;
    fieldDisabled?: (id: string) => boolean;
  }>;
}

type JsonRecord = Record<string, unknown>;

function isRecord(value: unknown): value is JsonRecord {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function stringOf(value: unknown): string {
  return typeof value === "string" ? value : "";
}

// --- Schema CRUD context (S4 · I-007-003 §9) ---

export interface SchemaCrudFeedback {
  kind: "success" | "error";
  message: string;
  code?: string;
}

export interface SchemaCrudConfirm {
  actionRef: string;
  actionKey: string;
  row: Record<string, unknown>;
  requestMapping?: Record<string, unknown>;
  message: string;
}

export type ActionResult = { ok: true } | { ok: false; code: string; message: string };

export interface SchemaCrudValue {
  selectedRow: Record<string, unknown> | null;
  selectRow: (row: Record<string, unknown> | null) => void;
  tableQuery: (id: string) => ResourceQuery | undefined;
  setTableQuery: (id: string, query: ResourceQuery) => void;
  reloadToken: number;
  reloadList: () => void;
  activeModal: { actionRef: string; row: Record<string, unknown> | null; title: string } | null;
  modalRow: Record<string, unknown> | null;
  openModal: (actionRef: string, row: Record<string, unknown> | null, title: string) => void;
  closeModal: () => void;
  pendingConfirm: SchemaCrudConfirm | null;
  requestConfirm: (confirm: SchemaCrudConfirm) => void;
  resolveConfirm: (confirmed: boolean) => Promise<void>;
  feedback: SchemaCrudFeedback | null;
  registerFetcher: (fetcher: typeof fetch) => void;
  /** The currently registered transport (globalThis.fetch until injected). */
  fetcher: typeof fetch;
  /** Runs a request action end-to-end (gate → construct → fetch → feedback/reload). */
  runRowAction: (
    actionRef: string,
    opts: RunRequestOptions,
  ) => Promise<ActionResult>;
  /** Dispatches a toolbar/row action entry: modal open, confirm, or request. */
  invokeAction: (item: Record<string, unknown>, row: Record<string, unknown> | null) => void;
  /** Submits a default-mode form against its `submitAction`. */
  submitForm: (form: RenderFormNode, values: Record<string, unknown>) => Promise<ActionResult>;
  /** Binds a search-mode form's fields to its target table query. */
  searchFormSubmit: (form: RenderFormNode, values: Record<string, unknown>) => void;
  effectivePermission: (targetId: string) => boolean;
}

export interface RunRequestOptions {
  row?: Record<string, unknown> | null;
  formValues?: Record<string, unknown>;
  requestMapping?: Record<string, unknown>;
  gateTargetId?: string;
  confirmed?: boolean;
}

const SchemaCrudContext = createContext<SchemaCrudValue | null>(null);

/** Reads the page-level Schema CRUD provider (null when rendered bare). */
export function useSchemaCrud(): SchemaCrudValue | null {
  return useContext(SchemaCrudContext);
}

function pageActions(document: RenderPageDocument): Record<string, unknown> {
  const doc = document as unknown as JsonRecord;
  return isRecord(doc.actions) ? doc.actions : {};
}

function actionOf(document: RenderPageDocument, actionRef: string): JsonRecord | undefined {
  const action = pageActions(document)[actionRef];
  return isRecord(action) ? action : undefined;
}

function successMessageFor(method: unknown): string {
  switch (method) {
    case "POST":
      return "Item created";
    case "PATCH":
      return "Item updated";
    case "DELETE":
      return "Item deleted";
    default:
      return "Action completed";
  }
}

/**
 * Executes one request action against the frozen permission/confirm gate and
 * the pinned request-construction adapter (I-007-003 §9.1/§9.1a).
 *
 * - `formValues` → formAction construction (create/edit form submits).
 * - `row` + `requestMapping` → rowAction construction (delete).
 * - `{id}` slots on a formAction URL are resolved from the captured row context
 *   (bounded extension, §9.1a); row actions resolve `$row.*` in the constructor.
 */
async function runRequest(
  document: RenderPageDocument,
  context: Record<string, unknown>,
  fetcher: typeof fetch,
  actionRef: string,
  opts: RunRequestOptions,
): Promise<ActionResult> {
  const action = actionOf(document, actionRef);
  if (action === undefined) {
    return { ok: false, code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page` };
  }
  if (action.type !== "request") {
    return { ok: false, code: "ACTION_NOT_REQUEST", message: `action "${actionRef}" is not a request action` };
  }
  if (opts.gateTargetId !== undefined) {
    const gate = executeAction(
      document as unknown as JsonRecord,
      {
        targetId: opts.gateTargetId,
        visible: true,
        ...(opts.confirmed !== undefined ? { confirm: true, confirmed: opts.confirmed } : {}),
        disabled: false,
      },
      context as NavigationContext,
    );
    if (gate.outcome !== "EXECUTED") {
      return {
        ok: false,
        code: gate.reason ?? gate.outcome,
        message: `action "${actionRef}" was not executed (${gate.reason ?? gate.outcome})`,
      };
    }
  }

  const formProvided = opts.formValues !== undefined;
  const rowProvided = opts.row !== undefined && opts.row !== null;
  const kind = formProvided ? "formAction" : rowProvided ? "rowAction" : "pageTriggerRequest";
  const input: JsonRecord = { kind, action, baseURL: "" };
  if (formProvided) {
    input.formValues = opts.formValues;
  }
  if (rowProvided) {
    input.row = opts.row;
    if (opts.requestMapping !== undefined) {
      input.requestMapping = opts.requestMapping;
    }
  }
  const constructed = constructRequest(input);
  if (!constructed.ok) {
    return { ok: false, code: constructed.code, message: `request construction failed (${constructed.path})` };
  }
  const request = constructed.request;
  if (request === undefined) {
    return { ok: false, code: "NO_REQUEST", message: "action produced no request" };
  }

  // Bounded formAction `{id}` slot binding (I-007-003 §9.1a): a default form
  // submit targeting a path slot resolves it from the captured row context.
  let url = request.url;
  const rowId = opts.row?.id;
  if (formProvided && rowProvided && url.includes("{id}") && typeof rowId === "string") {
    url = url.replaceAll("{id}", encodeURIComponent(rowId));
  }

  const body = request.body === null || request.body === undefined ? undefined : JSON.stringify(request.body);
  const response = await fetcher(url, {
    method: request.method,
    headers: { "Content-Type": "application/json" },
    ...(body === undefined ? {} : { body }),
  });
  if (!response.ok) {
    const apiError = await readResourceApiError(response, actionRef);
    return { ok: false, code: apiError.code, message: apiError.message };
  }
  // Branding refresh after settings PATCH lives in the App/host layer (A-006 R-002),
  // not in this generic request executor — keep Renderer free of product endpoints.
  return { ok: true };
}

function SchemaCrudProvider({
  document,
  context,
  children,
  initialFetcher,
}: {
  document: RenderPageDocument;
  context: Record<string, unknown>;
  children: ReactNode;
  initialFetcher?: typeof fetch;
}) {
  const [selectedRow, setSelectedRow] = useState<Record<string, unknown> | null>(null);
  const [queries, setQueries] = useState<Record<string, ResourceQuery>>({});
  const [reloadToken, setReloadToken] = useState(0);
  const [activeModal, setActiveModal] = useState<{
    actionRef: string;
    row: Record<string, unknown> | null;
    title: string;
  } | null>(null);
  const [pendingConfirm, setPendingConfirm] = useState<SchemaCrudConfirm | null>(null);
  const [feedback, setFeedback] = useState<SchemaCrudFeedback | null>(null);
  const [fetcher, setFetcher] = useState<typeof fetch>(() => initialFetcher ?? globalThis.fetch);

  const permissionTargets = useMemo(
    () => evaluatePermissionTargets(document as unknown as JsonRecord, context as NavigationContext),
    [document, context],
  );

  const effectivePermission = useCallback(
    (targetId: string) => {
      const entry = permissionTargets.find((target) => target.targetId === targetId);
      // Absent target = no declared permission (engine default is allow).
      return entry === undefined ? true : entry.effectivePermission;
    },
    [permissionTargets],
  );

  const registerFetcher = useCallback((next: typeof fetch) => {
    // Keep the first injected transport; the fixture-level fetcher identity in
    // tests can change every render, and returning `prev` bails out of state
    // updates (no re-render loop). Production fetchers are stable anyway.
    setFetcher((prev: typeof fetch) => (prev !== globalThis.fetch ? prev : next));
  }, []);

  const tableQuery = useCallback((id: string) => queries[id], [queries]);
  const setTableQuery = useCallback((id: string, query: ResourceQuery) => {
    setQueries((prev) => ({ ...prev, [id]: query }));
  }, []);
  const reloadList = useCallback(() => setReloadToken((token) => token + 1), []);

  const openModal = useCallback(
    (actionRef: string, row: Record<string, unknown> | null, title: string) => {
      setSelectedRow(row);
      setActiveModal({ actionRef, row, title });
    },
    [],
  );
  const closeModal = useCallback(() => setActiveModal(null), []);

  const runRequestCallback = useCallback(
    (actionRef: string, opts: RunRequestOptions) =>
      runRequest(document, context, fetcher, actionRef, opts),
    [document, context, fetcher],
  );

  const runRowAction = useCallback(
    async (actionRef: string, opts: RunRequestOptions): Promise<ActionResult> => {
      const result = await runRequestCallback(actionRef, opts);
      if (result.ok) {
        setFeedback({ kind: "success", message: successMessageFor(actionOf(document, actionRef)?.method) });
        reloadList();
        setSelectedRow(null);
      } else {
        setFeedback({ kind: "error", code: result.code, message: result.message });
      }
      return result;
    },
    [runRequestCallback, document, reloadList],
  );

  const invokeAction = useCallback(
    (item: Record<string, unknown>, row: Record<string, unknown> | null) => {
      const actionRef = stringOf(item.actionRef);
      const action = actionOf(document, actionRef);
      if (actionRef === "" || action === undefined) {
        setFeedback({ kind: "error", code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page` });
        return;
      }
      setSelectedRow(row);
      if (action.type === "modal") {
        setActiveModal({ actionRef, row, title: stringOf(item.label) ?? "Action" });
        return;
      }
      const gateTargetId = stringOf(item.key) !== "" ? stringOf(item.key) : actionRef;
      const requestMapping = isRecord(item.requestMapping) ? item.requestMapping : undefined;
      const confirmMessage = typeof item.confirm === "string" && item.confirm !== "" ? item.confirm : undefined;
      if (confirmMessage !== undefined) {
        setPendingConfirm({
          actionRef,
          actionKey: gateTargetId,
          row: row ?? {},
          requestMapping,
          message: confirmMessage,
        });
        return;
      }
      void runRowAction(actionRef, {
        row: row ?? undefined,
        requestMapping,
        gateTargetId,
      });
    },
    [document, runRowAction],
  );

  const resolveConfirm = useCallback(
    async (confirmed: boolean) => {
      if (pendingConfirm === null) {
        return;
      }
      const { actionRef, actionKey, row, requestMapping } = pendingConfirm;
      setPendingConfirm(null);
      if (!confirmed) {
        return;
      }
      await runRowAction(actionRef, { row, requestMapping, gateTargetId: actionKey, confirmed: true });
    },
    [pendingConfirm, runRowAction],
  );

  const submitForm = useCallback(
    async (form: RenderFormNode, values: Record<string, unknown>): Promise<ActionResult> => {
      const submitAction = form.props.submitAction;
      if (typeof submitAction !== "string") {
        return { ok: false, code: "NO_SUBMIT_ACTION", message: "form has no submitAction" };
      }
      const rawStringFields = new Set(
        form.props.fields
          .filter((field) => isRecord(field) && field.type === "password")
          .map((field) => (isRecord(field) && typeof field.id === "string" ? field.id : ""))
          .filter((field) => field !== ""),
      );
      // Text fields keep the historical trim rule; secret controls preserve
      // the exact string so the browser and API hash the same password bytes.
      const prepared: Record<string, unknown> = {};
      for (const [key, value] of Object.entries(values)) {
        prepared[key] =
          typeof value === "string" && !rawStringFields.has(key) ? value.trim() : value;
      }
      // dateRangePicker binds two independent fields (registry startField /
      // endField): expand the {start,end} pair into the two output keys and
      // drop the range control's own id from the submit projection.
      for (const raw of form.props.fields) {
        if (!isRecord(raw) || raw.type !== "dateRangePicker") {
          continue;
        }
        const rangeId = typeof raw.id === "string" ? raw.id : "";
        const startField = typeof raw.startField === "string" ? raw.startField : "";
        const endField = typeof raw.endField === "string" ? raw.endField : "";
        if (rangeId === "" || startField === "" || endField === "") {
          continue;
        }
        const pair = prepared[rangeId];
        if (isRecord(pair)) {
          prepared[startField] = typeof pair.start === "string" ? pair.start.trim() : "";
          prepared[endField] = typeof pair.end === "string" ? pair.end.trim() : "";
        } else {
          prepared[startField] = "";
          prepared[endField] = "";
        }
        delete prepared[rangeId];
      }
      const gateTargetId = `${form.id ?? "unnamed"}:submit`;
      const result = await runRequestCallback(submitAction, {
        formValues: prepared,
        row: activeModal?.row ?? null,
        gateTargetId,
      });
      if (result.ok) {
        setFeedback({ kind: "success", message: successMessageFor(actionOf(document, submitAction)?.method) });
        reloadList();
        setActiveModal(null);
        setSelectedRow(null);
      }
      return result;
    },
    [runRequestCallback, document, activeModal, reloadList],
  );

  const searchFormSubmit = useCallback((form: RenderFormNode, values: Record<string, unknown>) => {
    const targetTable = form.props.targetTable;
    if (typeof targetTable !== "string" || targetTable === "") {
      return;
    }
    const raw = values.q;
    const q = typeof raw === "string" ? raw.trim() : "";
    setQueries((prev) => {
      const current = prev[targetTable] ?? { page: 1, pageSize: 10 };
      return { ...prev, [targetTable]: { ...current, page: 1, ...(q === "" ? {} : { q }) } };
    });
  }, []);

  const value = useMemo<SchemaCrudValue>(
    () => ({
      selectedRow,
      selectRow: setSelectedRow,
      tableQuery,
      setTableQuery,
      reloadToken,
      reloadList,
      activeModal,
      modalRow: activeModal?.row ?? null,
      openModal,
      closeModal,
      pendingConfirm,
      requestConfirm: setPendingConfirm,
      resolveConfirm,
      feedback,
      registerFetcher,
      fetcher,
      runRowAction,
      invokeAction,
      submitForm,
      searchFormSubmit,
      effectivePermission,
    }),
    [
      selectedRow,
      tableQuery,
      setTableQuery,
      reloadToken,
      reloadList,
      activeModal,
      openModal,
      closeModal,
      pendingConfirm,
      resolveConfirm,
      feedback,
      registerFetcher,
      fetcher,
      runRowAction,
      invokeAction,
      submitForm,
      searchFormSubmit,
      effectivePermission,
    ],
  );

  return <SchemaCrudContext.Provider value={value}>{children}</SchemaCrudContext.Provider>;
}

function FeedbackRegion({ feedback }: { feedback: SchemaCrudFeedback }) {
  return (
    <div
      role={feedback.kind === "error" ? "alert" : "status"}
      className={`mb-4 rounded-md border px-3 py-2 text-sm ${
        feedback.kind === "error"
          ? "border-destructive/50 bg-destructive/10 text-destructive"
          : "border-emerald-500/50 bg-emerald-500/10 text-emerald-700"
      }`}
    >
      {feedback.code !== undefined && feedback.code !== "" ? `${feedback.code}: ` : ""}
      {feedback.message}
    </div>
  );
}

function FormView({
  node,
  metaValue,
  context,
  formComponent,
}: {
  node: RenderFormNode;
  metaValue: unknown;
  context: Record<string, unknown>;
  formComponent?: RendererComponentProps["formComponent"];
}) {
  const Component = formComponent ?? FormControls;
  const reaction = resolveFormReactions(node, context);
  const gate = gateRenderFormFields(metaValue, node.props.fields, "fields");
  const crud = useSchemaCrud();
  const modalRow = crud?.modalRow ?? null;
  const [values, setValues] = useState<Record<string, unknown>>(() => {
    const initial: Record<string, unknown> = {};
    for (const raw of node.props.fields) {
      if (!isRecord(raw) || typeof raw.id !== "string") {
        continue;
      }
      // Always coerce through the field wire kind so row arrays (e.g. roles[])
      // become textarea strings or checkboxGroup string[] as appropriate (A-006 R-004).
      const fromRow =
        modalRow !== null && modalRow[raw.id] !== undefined ? modalRow[raw.id] : undefined;
      initial[raw.id] = coerceFieldValue(raw as unknown as FormControlField, fromRow);
    }
    return initial;
  });
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<{ code: string; message: string } | null>(null);

  const isSearch = node.props.mode === "search";
  const submitAction = node.props.submitAction;
  const canSubmit = isSearch || typeof submitAction === "string";
  const hasBlockingErrors = gate.errors.length > 0 || reaction.errors.length > 0;

  const visibleFields = gate.fields.filter(
    (raw) => reaction.state[raw.id]?.visible !== false,
  );

  const fieldDisabled = (id: string) => reaction.state[id]?.disabled === true;

  const handleSubmit = async () => {
    if (crud === null) {
      return;
    }
    if (hasBlockingErrors) {
      return;
    }
    if (isSearch) {
      crud.searchFormSubmit(node, values);
      return;
    }
    if (typeof submitAction !== "string") {
      return;
    }
    setSubmitting(true);
    setFormError(null);
    const result = await crud.submitForm(node, values);
    if (!result.ok) {
      setFormError({ code: result.code, message: result.message });
    }
    setSubmitting(false);
  };

  return (
    <form className="space-y-3" onSubmit={(event) => {
      event.preventDefault();
      void handleSubmit();
    }}>
      {gate.errors.length > 0 ? (
        <ul role="alert" className="space-y-1 text-sm text-destructive">
          {gate.errors.map((error, index) => (
            <li key={index}>
              {error.code}: {error.message}
            </li>
          ))}
        </ul>
      ) : null}
      <Component
        fields={visibleFields}
        values={values}
        onChange={(id, value) => setValues((prev) => ({ ...prev, [id]: value }))}
        fieldDisabled={fieldDisabled}
      />
      {reaction.errors.length > 0 ? (
        <ul role="alert" className="space-y-1 text-sm text-destructive">
          {reaction.errors.map((error, index) => (
            <li key={index}>
              {error.code}: {error.message}
            </li>
          ))}
        </ul>
      ) : null}
      {formError !== null ? (
        <p role="alert" className="text-sm text-destructive">
          {formError.code}: {formError.message}
        </p>
      ) : null}
      {canSubmit ? (
        <button
          type="submit"
          disabled={submitting || hasBlockingErrors}
          className="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {submitting ? "Submitting…" : (node.props.submitLabel ?? (isSearch ? "Search" : "Submit"))}
        </button>
      ) : null}
    </form>
  );
}

function SectionView({
  node,
  metaValue,
  context,
  tableRenderer,
  onAction,
  formComponent,
}: {
  node: RenderSectionNode;
  metaValue: unknown;
  context: Record<string, unknown>;
  tableRenderer?: RendererComponentProps["tableRenderer"];
  onAction?: RendererComponentProps["onAction"];
  formComponent?: RendererComponentProps["formComponent"];
}) {
  return (
    <div className="space-y-6">
      {node.children.map((child, index) => (
        <Fragment key={index}>
          {dispatchNode({
            node: child,
            path: `body.children[${index}]`,
            metaValue,
            context,
            tableRenderer,
            onAction,
            formComponent,
          })}
        </Fragment>
      ))}
    </div>
  );
}

function GridView({
  node,
  metaValue,
  context,
  tableRenderer,
  onAction,
  formComponent,
}: {
  node: RenderGridNode;
  metaValue: unknown;
  context: Record<string, unknown>;
  tableRenderer?: RendererComponentProps["tableRenderer"];
  onAction?: RendererComponentProps["onAction"];
  formComponent?: RendererComponentProps["formComponent"];
}) {
  const columns =
    isRecord(node.props) && typeof node.props.columns === "number"
      ? node.props.columns
      : undefined;
  return (
    <div
      className="grid gap-4"
      style={columns !== undefined ? { gridTemplateColumns: `repeat(${columns}, minmax(0, 1fr))` } : undefined}
    >
      {(node.children ?? []).map((child, index) => (
        <Fragment key={index}>
          {dispatchNode({
            node: child,
            path: `body.children[${index}]`,
            metaValue,
            context,
            tableRenderer,
            onAction,
            formComponent,
          })}
        </Fragment>
      ))}
    </div>
  );
}

function TabsView({
  node,
  metaValue,
  context,
  tableRenderer,
  onAction,
  formComponent,
}: {
  node: RenderTabsNode;
  metaValue: unknown;
  context: Record<string, unknown>;
  tableRenderer?: RendererComponentProps["tableRenderer"];
  onAction?: RendererComponentProps["onAction"];
  formComponent?: RendererComponentProps["formComponent"];
}) {
  const children = node.children ?? [];
  const [active, setActive] = useState(0);
  const current = children[Math.min(active, children.length - 1)];
  return (
    <div className="space-y-3">
      {children.length > 1 ? (
        <div role="tablist" className="flex flex-wrap items-center gap-1 border-b border-border">
          {children.map((child, index) => {
            const rawProps = (child as unknown as { props?: unknown }).props;
            const label =
              isRecord(rawProps) && typeof rawProps.label === "string"
                ? rawProps.label
                : `Tab ${index + 1}`;
            return (
              <button
                key={index}
                type="button"
                role="tab"
                aria-selected={index === active}
                onClick={() => setActive(index)}
                className={`rounded-t px-3 py-1.5 text-sm ${
                  index === active
                    ? "border-b-2 border-primary font-medium text-foreground"
                    : "text-muted-foreground"
                }`}
              >
                {label}
              </button>
            );
          })}
        </div>
      ) : null}
      {current !== undefined ? (
        dispatchNode({
          node: current,
          path: "body.children",
          metaValue,
          context,
          tableRenderer,
          onAction,
          formComponent,
        })
      ) : (
        <p className="text-sm text-muted-foreground">tabs node has no children</p>
      )}
    </div>
  );
}

function TextView({ node }: { node: RenderTextNode }) {
  return <p className="text-sm text-foreground">{node.props?.text ?? ""}</p>;
}

function RecordView({ node }: { node: RenderRecordViewNode }) {
  const crud = useSchemaCrud();
  const staticRecord = node.props?.record;
  const record = isRecord(staticRecord) ? staticRecord : (crud?.selectedRow ?? null);
  const entries = isRecord(record) ? Object.entries(record) : [];
  if (entries.length === 0) {
    return <p className="text-sm text-muted-foreground">Select a record to view details.</p>;
  }
  return (
    <dl className="space-y-1 rounded-md border border-border bg-card p-4">
      {entries.map(([key, value]) => (
        <div key={key} className="flex gap-4 text-sm">
          <dt className="w-32 shrink-0 text-muted-foreground">{key}</dt>
          <dd className="text-foreground">{String(value)}</dd>
        </div>
      ))}
    </dl>
  );
}

/** Shared data fetch for statCard/chart (supportsData components, registry). */
function useDisplayData(
  dataSource: string | null,
  fetcher: typeof fetch,
): { list: ResourceList | null; error: string | null } {
  const crud = useSchemaCrud();
  useEffect(() => {
    if (crud !== null && fetcher !== undefined) {
      crud.registerFetcher(fetcher);
    }
  }, [crud, fetcher]);
  const [list, setList] = useState<ResourceList | null>(null);
  const [error, setError] = useState<string | null>(null);
  useEffect(() => {
    if (dataSource === null) {
      setList(null);
      setError(null);
      return;
    }
    let cancelled = false;
    fetchResourceList(fetcher ?? fetch, dataSource, { page: 1, pageSize: 100 })
      .then((next) => {
        if (!cancelled) {
          setList(next);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          setError(err instanceof Error ? err.message : String(err));
        }
      });
    return () => {
      cancelled = true;
    };
  }, [fetcher, dataSource, crud?.reloadToken]);
  return { list, error };
}

function StatCardView({ node }: { node: RenderStatCardNode }) {
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;
  const valueField = node.props?.valueField;
  const format = node.props?.format ?? "plain";
  const dataSource =
    typeof node.props?.dataSource === "string" && isValidDataSource(node.props.dataSource)
      ? node.props.dataSource
      : null;
  const { list, error } = useDisplayData(dataSource, fetcher);

  if (dataSource === null) {
    return (
      <p role="alert" className="text-sm text-destructive">
        statCard node requires a valid dataSource (single-slash same-origin path)
      </p>
    );
  }
  if (error !== null) {
    return <p role="alert" className="text-sm text-destructive">statCard data failed to load: {error}</p>;
  }
  if (list === null) {
    return <p className="text-sm text-muted-foreground">Loading statCard…</p>;
  }
  const first = list.items[0];
  const raw = valueField !== undefined && first !== undefined ? first[valueField] : 0;
  // Registry statCard format enum: plain | currency | percent.
  if (format !== "plain" && format !== "currency" && format !== "percent") {
    return (
      <p role="alert" className="text-sm text-destructive">
        statCard format "{format}" is outside the registry enum (plain/currency/percent)
      </p>
    );
  }
  const formatResult = format === "plain" ? { ok: true as const, value: raw } : applyComponentFormat(format, raw);
  if (!formatResult.ok) {
    return (
      <p role="alert" className="text-sm text-destructive">
        statCard format "{format}" rejects the value type of field "{valueField}" ({formatResult.code})
      </p>
    );
  }
  const label = node.props?.label ?? valueField ?? "Value";
  const unit = node.props?.unit;
  return (
    <div className="rounded-md border border-border bg-card p-4">
      <p className="text-xs font-medium text-muted-foreground">{label}</p>
      <p className="mt-1 text-2xl font-semibold text-foreground">
        {String(formatResult.value)}
        {unit !== undefined && unit !== "" ? <span className="ml-1 text-sm text-muted-foreground">{unit}</span> : null}
      </p>
    </div>
  );
}

function ChartView({ node }: { node: RenderChartNode }) {
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;
  const chartType = node.props?.chartType;
  const xField = node.props?.xField;
  const yField = node.props?.yField;
  const dataSource =
    typeof node.props?.dataSource === "string" && isValidDataSource(node.props.dataSource)
      ? node.props.dataSource
      : null;
  const { list, error } = useDisplayData(dataSource, fetcher);

  const missingProps = chartType === undefined || xField === undefined || yField === undefined;
  if (dataSource === null || missingProps) {
    return (
      <p role="alert" className="text-sm text-destructive">
        chart node requires chartType / xField / yField and a valid dataSource
      </p>
    );
  }
  if (error !== null) {
    return <p role="alert" className="text-sm text-destructive">chart data failed to load: {error}</p>;
  }
  if (list === null) {
    return <p className="text-sm text-muted-foreground">Loading chart…</p>;
  }
  const points = list.items
    .map((row) => ({
      x: String(row[xField] ?? ""),
      y: typeof row[yField] === "number" ? row[yField] : Number(row[yField]),
    }))
    .filter((point) => Number.isFinite(point.y) && point.x !== "");
  if (points.length === 0) {
    return <p className="text-sm text-muted-foreground">chart has no plottable data points</p>;
  }
  const maxY = Math.max(...points.map((point) => point.y), 1);
  const width = 320;
  const height = 160;
  const pad = 8;
  const innerW = width - pad * 2;
  const innerH = height - pad * 2;
  return (
    <div className="rounded-md border border-border bg-card p-4">
      <svg
        role="img"
        aria-label={`${chartType} chart (${xField} / ${yField})`}
        viewBox={`0 0 ${width} ${height}`}
        className="h-40 w-full max-w-md"
      >
        {chartType === "line" ? (
          <polyline
            fill="none"
            stroke="currentColor"
            strokeWidth={2}
            points={points
              .map(
                (point, index) =>
                  `${pad + (index / Math.max(points.length - 1, 1)) * innerW},${pad + innerH - (point.y / maxY) * innerH}`,
              )
              .join(" ")}
          />
        ) : chartType === "bar" ? (
          points.map((point, index) => {
            const barWidth = innerW / points.length;
            const barHeight = (point.y / maxY) * innerH;
            return (
              <rect
                key={index}
                x={pad + index * barWidth + barWidth * 0.2}
                y={pad + innerH - barHeight}
                width={barWidth * 0.6}
                height={barHeight}
                fill="currentColor"
              />
            );
          })
        ) : chartType === "pie" ? (
          // Donut slices via stroke-dasharray (no trig needed; ≤4 slices fits
          // the example data and keeps the renderer dependency-free).
          (() => {
            const total = points.reduce((sum, point) => sum + point.y, 0) || 1;
            const radius = 56;
            const circumference = 2 * Math.PI * radius;
            let offset = 0;
            return points.map((point, index) => {
              const fraction = point.y / total;
              const dash = fraction * circumference;
              const slice = (
                <circle
                  key={index}
                  cx={width / 2}
                  cy={height / 2}
                  r={radius}
                  fill="none"
                  strokeWidth={28}
                  stroke={`hsl(${(index * 137.5) % 360} 65% 55%)`}
                  strokeDasharray={`${dash} ${circumference - dash}`}
                  strokeDashoffset={-offset}
                />
              );
              offset += dash;
              return slice;
            });
          })()
        ) : null}
      </svg>
      <ul className="mt-2 space-y-0.5 text-xs text-muted-foreground">
        {points.map((point, index) => (
          <li key={index}>
            {point.x}: {point.y}
          </li>
        ))}
      </ul>
    </div>
  );
}

function ActionButtonView({
  node,
  context,
  onAction,
}: {
  node: RenderActionButtonNode;
  context: Record<string, unknown>;
  onAction?: RendererComponentProps["onAction"];
}) {
  const gate = tableActionGate(node.props ?? {}, context);
  if (!gate.visible) {
    return null;
  }
  return (
    <div className="space-y-1">
      <button
        type="button"
        disabled={gate.disabled}
        onClick={() => onAction?.(node)}
        className="h-9 rounded-md border border-input bg-background px-3 text-sm transition-colors focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50"
      >
        {node.props?.label ?? node.props?.actionId ?? "Action"}
      </button>
      {gate.errors.length > 0 ? (
        <ul role="alert" className="space-y-1 text-sm text-destructive">
          {gate.errors.map((error, index) => (
            <li key={index}>
              {error.code}: {error.message}
            </li>
          ))}
        </ul>
      ) : null}
    </div>
  );
}

function dispatchNode({
  node,
  path,
  metaValue,
  context,
  tableRenderer,
  onAction,
  formComponent,
}: {
  node: unknown;
  path: string;
  metaValue: unknown;
  context: Record<string, unknown>;
  tableRenderer?: RendererComponentProps["tableRenderer"];
  onAction?: RendererComponentProps["onAction"];
  formComponent?: RendererComponentProps["formComponent"];
}): ReactNode {
  const parsed = parseRenderNode(node, path);
  if ("code" in parsed) {
    return (
      <p key={path} role="alert" className="text-sm text-destructive">
        {parsed.message}
      </p>
    );
  }
  return dispatchParsedNode({
    node: parsed,
    metaValue,
    context,
    tableRenderer,
    onAction,
    formComponent,
  });
}

function dispatchParsedNode({
  node,
  metaValue,
  context,
  tableRenderer,
  onAction,
  formComponent,
}: {
  node: RenderNode;
  metaValue: unknown;
  context: Record<string, unknown>;
  tableRenderer?: RendererComponentProps["tableRenderer"];
  onAction?: RendererComponentProps["onAction"];
  formComponent?: RendererComponentProps["formComponent"];
}): ReactNode {
  switch (node.type) {
    case "form":
      return <FormView node={node} metaValue={metaValue} context={context} formComponent={formComponent} />;
    case "section":
      return (
        <SectionView
          node={node}
          metaValue={metaValue}
          context={context}
          tableRenderer={tableRenderer}
          onAction={onAction}
          formComponent={formComponent}
        />
      );
    case "grid":
      return (
        <GridView
          node={node}
          metaValue={metaValue}
          context={context}
          tableRenderer={tableRenderer}
          onAction={onAction}
          formComponent={formComponent}
        />
      );
    case "tabs":
      return (
        <TabsView
          node={node}
          metaValue={metaValue}
          context={context}
          tableRenderer={tableRenderer}
          onAction={onAction}
          formComponent={formComponent}
        />
      );
    case "text":
      return <TextView node={node} />;
    case "recordView":
      return <RecordView node={node} />;
    case "actionButton":
      return <ActionButtonView node={node} context={context} onAction={onAction} />;
    case "statCard":
      return <StatCardView node={node} />;
    case "chart":
      return <ChartView node={node} />;
    case "table":
      return (
        tableRenderer?.(node) ?? (
          <p className="text-sm text-muted-foreground">
            table node rendered without a tableRenderer (the app wires SchemaTable)
          </p>
        )
      );
  }
}

function RenderPageSurface({
  document,
  context,
  tableRenderer,
  onAction,
  formComponent,
}: RendererComponentProps) {
  const crud = useSchemaCrud()!;
  const modalAction =
    crud.activeModal !== null ? actionOf(document, crud.activeModal.actionRef) : undefined;
  const modalContent =
    crud.activeModal !== null &&
    modalAction !== undefined &&
    modalAction.type === "modal" &&
    isRecord(modalAction.content)
      ? modalAction.content
      : undefined;

  return (
    <Fragment>
      {crud.feedback !== null ? <FeedbackRegion feedback={crud.feedback} /> : null}
      {dispatchNode({
        node: document.body,
        path: "body",
        metaValue: document.meta,
        context,
        tableRenderer,
        onAction,
        formComponent,
      })}
      {modalContent !== undefined ? (
        <ModalHost
          key={`${crud.activeModal!.actionRef}:${isRecord(crud.activeModal!.row) ? stringOf(crud.activeModal!.row.id) : "new"}`}
          title={crud.activeModal!.title}
          onClose={crud.closeModal}
        >
          {dispatchNode({
            node: modalContent,
            path: `actions.${crud.activeModal!.actionRef}.content`,
            metaValue: document.meta,
            context,
            tableRenderer,
            onAction,
            formComponent,
          })}
        </ModalHost>
      ) : null}
      {crud.pendingConfirm !== null ? (
        <ConfirmDialog
          message={crud.pendingConfirm.message}
          onConfirm={() => void crud.resolveConfirm(true)}
          onCancel={() => void crud.resolveConfirm(false)}
        />
      ) : null}
    </Fragment>
  );
}

export function RenderPage({
  document,
  context,
  tableRenderer,
  dataFetcher,
  onAction,
  formComponent,
}: RendererComponentProps & { dataFetcher?: typeof fetch }) {
  return (
    <SchemaCrudProvider document={document} context={context} initialFetcher={dataFetcher}>
      <RenderPageSurface
        document={document}
        context={context}
        tableRenderer={tableRenderer}
        onAction={onAction}
        formComponent={formComponent}
      />
    </SchemaCrudProvider>
  );
}
