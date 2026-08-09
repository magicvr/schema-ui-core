import {
  Fragment,
  createContext,
  useCallback,
  useContext,
  useEffect,
  useMemo,
  useRef,
  useState,
  type ComponentType,
  type ReactNode,
} from "react";

import type { NavigationContext } from "@/protocol/app-manifest";
import { applyComponentFormat } from "@/protocol/conformance/component-format";
import { resolveAsyncDisplayState } from "@/components/ui/async-state";
import { Card, CardContent } from "@/components/ui/card";
import { Skeleton } from "@/components/ui/skeleton";
import {
  constructRequest,
  normalizeSelection,
} from "@/protocol/conformance/request-construction";
import {
  uploadFilesWithFetch,
  type UploadActionResult,
  type UploadableFile,
} from "@/protocol/conformance/upload-orchestration";
import { ConfirmDialog } from "@/renderer/confirm";
import { FormControls } from "@/renderer/form-controls.tsx";
import {
  FORM_RECORD_LOAD_CAPABILITY,
  coerceFieldValue,
  type FormControlField,
} from "@/renderer/form-controls";
import { ModalHost } from "@/renderer/modal";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
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
  resolveFullFormReactions,
  type FormControlStateMap,
  type ReactionError,
} from "@/renderer/reactions";
import {
  gateRenderFormFields,
  parseRenderNode,
  resolveFormReactions,
  resolveResponsePath,
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
    onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
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
  /** VP-007 S4: catalog key for the frontend localization floor. */
  messageKey?: string;
  /** VP-007 S4: interpolation params for messageKey. */
  params?: Record<string, unknown>;
}

export interface SchemaCrudConfirm {
  actionRef: string;
  actionKey: string;
  row: Record<string, unknown>;
  requestMapping?: Record<string, unknown>;
  message: string;
  /** Batch trigger confirm (ADR-0022 D4/D5): carries the selection snapshot. */
  batch?: { tableId: string; selection: TableSelection };
  /** Toolbar item batchMapping carried through confirm. */
  batchMapping?: Record<string, unknown>;
}

/** Selection snapshot (ADR-0022 D3): ordered keys + count. */
export interface TableSelection {
  keys: unknown[];
  count: number;
}

export type ActionResult =
  | { ok: true }
  | { ok: false; code: string; message: string; messageKey?: string; params?: Record<string, unknown> };

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
  /** Dispatches a batch toolbar trigger (ADR-0022): gate → confirm → request. */
  invokeBatchAction: (item: Record<string, unknown>, tableId: string) => void;
  /** Upload control transport (ADR-0012): resolves action/actionRef → validates → uploads. */
  uploadFiles: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
  /** Per-table selection state (keys + count; normalized by the table). */
  selection: (tableId: string) => TableSelection | undefined;
  setSelection: (tableId: string, keys: unknown[]) => void;
  clearSelection: (tableId: string) => void;
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

type Translator = (key: string, params?: MessageParams, literalFallback?: string) => string;

function successMessageFor(method: unknown, t: Translator): string {
  switch (method) {
    case "POST":
      return t("feedback.itemCreated");
    case "PATCH":
      return t("feedback.itemUpdated");
    case "DELETE":
      return t("feedback.itemDeleted");
    default:
      return t("feedback.actionCompleted");
  }
}

/** Batch triggers process a selection: delete-shaped URLs get the plural message. */
function batchSuccessMessageFor(action: JsonRecord, t: Translator): string {
  const url = stringOf(action.url);
  if (url.endsWith("/batch-delete")) {
    return t("feedback.itemsDeleted");
  }
  return successMessageFor(action.method, t);
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
    return {
      ok: false,
      code: apiError.code,
      message: apiError.message,
      ...(apiError.messageKey === undefined ? {} : { messageKey: apiError.messageKey }),
      ...(apiError.params === undefined ? {} : { params: apiError.params }),
    };
  }
  // Branding refresh after settings PATCH lives in the App/host layer (A-006 R-002),
  // not in this generic request executor — keep Renderer free of product endpoints.
  return { ok: true };
}

/**
 * Executes one batch toolbar trigger (ADR-0022 D5d): gate → normalize
 * selection → construct batchRequest → fetch → reload (clears selection).
 */
async function runBatchRequest(
  document: RenderPageDocument,
  context: Record<string, unknown>,
  fetcher: typeof fetch,
  actionRef: string,
  item: Record<string, unknown>,
  selection: TableSelection,
): Promise<ActionResult> {
  const action = actionOf(document, actionRef);
  if (action === undefined) {
    return { ok: false, code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page` };
  }
  if (action.type !== "request") {
    return { ok: false, code: "ACTION_NOT_REQUEST", message: `action "${actionRef}" is not a request action` };
  }
  const targetId = stringOf(item.key) !== "" ? stringOf(item.key) : actionRef;
  // ADR-0022 D5d: unmarked triggers do not participate in the permission
  // cascade ("未声明 intent 的入口仍只适用本地 permissions"); only gate when
  // the page declares a permission entry for this target.
  const hasPermissionEntry = evaluatePermissionTargets(
    document as unknown as JsonRecord,
    context as NavigationContext,
  ).some((entry) => entry.targetId === targetId);
  if (hasPermissionEntry) {
    const gate = executeAction(
      document as unknown as JsonRecord,
      {
        targetId,
        visible: true,
        disabled: false,
        requiresSelection: item.requiresSelection === true,
      },
      context as NavigationContext,
    );
    if (gate.outcome !== "EXECUTED") {
      return {
        ok: false,
        code: gate.reason ?? gate.outcome,
        message: `batch action "${actionRef}" was not executed (${gate.reason ?? gate.outcome})`,
      };
    }
  }
  const batchMapping = isRecord(item.batchMapping) ? item.batchMapping : undefined;
  const constructed = constructRequest({
    kind: "batchRequest",
    action,
    batchMapping,
    selection: { keys: selection.keys, count: selection.count },
  });
  if (!constructed.ok) {
    return { ok: false, code: constructed.code, message: `batch request construction failed (${constructed.path})` };
  }
  const request = constructed.request;
  if (request === undefined) {
    return { ok: false, code: "NO_REQUEST", message: "batch action produced no request" };
  }
  const body = request.body === null || request.body === undefined ? undefined : JSON.stringify(request.body);
  const response = await fetcher(request.url, {
    method: request.method,
    headers: { "Content-Type": "application/json" },
    ...(body === undefined ? {} : { body }),
  });
  if (!response.ok) {
    const apiError = await readResourceApiError(response, actionRef);
    return {
      ok: false,
      code: apiError.code,
      message: apiError.message,
      ...(apiError.messageKey === undefined ? {} : { messageKey: apiError.messageKey }),
      ...(apiError.params === undefined ? {} : { params: apiError.params }),
    };
  }
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
  const [selections, setSelections] = useState<Record<string, unknown[]>>({});
  const [reloadToken, setReloadToken] = useState(0);
  const [activeModal, setActiveModal] = useState<{
    actionRef: string;
    row: Record<string, unknown> | null;
    title: string;
  } | null>(null);
  const [pendingConfirm, setPendingConfirm] = useState<SchemaCrudConfirm | null>(null);
  const [feedback, setFeedback] = useState<SchemaCrudFeedback | null>(null);
  const [fetcher, setFetcher] = useState<typeof fetch>(() => initialFetcher ?? globalThis.fetch);
  const t = useTranslate();

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
  const selection = useCallback(
    (id: string): TableSelection | undefined => {
      const keys = selections[id];
      if (keys === undefined) {
        return undefined;
      }
      return { keys, count: keys.length };
    },
    [selections],
  );
  const setSelection = useCallback((id: string, keys: unknown[]) => {
    setSelections((prev) => ({ ...prev, [id]: normalizeSelection(keys).keys }));
  }, []);
  const clearSelection = useCallback((id: string) => {
    setSelections((prev) => {
      if (!(id in prev)) {
        return prev;
      }
      const next = { ...prev };
      delete next[id];
      return next;
    });
  }, []);
  // ADR-0022 D2: any data reload success clears every table selection.
  const reloadList = useCallback(() => {
    setSelections({});
    setReloadToken((token) => token + 1);
  }, []);

  const openModal = useCallback(
    (actionRef: string, row: Record<string, unknown> | null, title: string) => {
      // Prefill comes from activeModal.row (modalRow) — do not select the list row
      // (selectRow drives the recordView Drawer and would open it under Edit/New).
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
        setFeedback({ kind: "success", message: successMessageFor(actionOf(document, actionRef)?.method, t) });
        reloadList();
        setSelectedRow(null);
      } else {
        setFeedback(errorFeedback(result));
      }
      return result;
    },
    [runRequestCallback, document, reloadList, t],
  );

  const invokeAction = useCallback(
    (item: Record<string, unknown>, row: Record<string, unknown> | null) => {
      // actionButton nodes carry `props.actionId`; other entries use `actionRef`.
      const actionRef =
        stringOf(item.actionRef) !== "" ? stringOf(item.actionRef) : stringOf(item.actionId);
      const action = actionOf(document, actionRef);
      if (actionRef === "" || action === undefined) {
        setFeedback({ kind: "error", code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page` });
        return;
      }
      // Row actions carry the row in modal/confirm/request payloads. Do NOT call
      // setSelectedRow here — that opens the recordView Drawer and is wrong for
      // Edit / Delete / any toolbar-row action (user gap 2026-08-09).
      if (action.type === "modal") {
        setActiveModal({
          actionRef,
          row,
          title: resolveTextProp(item, "labelKey", "label", t, t("feedback.action")),
        });
        return;
      }
      const gateTargetId = stringOf(item.key) !== "" ? stringOf(item.key) : actionRef;
      const requestMapping = isRecord(item.requestMapping) ? item.requestMapping : undefined;
      const confirmMessage = resolveTextProp(item, "confirmKey", "confirm", t, "");
      if (confirmMessage !== "") {
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

  // ADR-0022 D4: batch toolbar trigger — confirm first, then run the batch.
  const invokeBatchAction = useCallback((item: Record<string, unknown>, tableId: string) => {
      const actionRef = stringOf(item.actionRef);
      const action = actionOf(document, actionRef);
      if (actionRef === "" || action === undefined) {
        setFeedback({ kind: "error", code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page` });
        return;
      }
      const current = selection(tableId);
      if (current === undefined || current.count === 0) {
        setFeedback({ kind: "error", code: "EMPTY_SELECTION", message: "select at least one row first" });
        return;
      }
      const confirmMessage = resolveTextProp(item, "confirmKey", "confirm", t, "");
      if (confirmMessage !== "") {
        setPendingConfirm({
          actionRef,
          actionKey: stringOf(item.key) !== "" ? stringOf(item.key) : actionRef,
          row: {},
          message: confirmMessage,
          batch: { tableId, selection: current },
          ...(isRecord(item.batchMapping) ? { batchMapping: item.batchMapping } : {}),
        });
        return;
      }
      void runBatchRequest(document, context, fetcher, actionRef, item, current).then((result) => {
        if (result.ok) {
          setFeedback({ kind: "success", message: batchSuccessMessageFor(action, t) });
          reloadList();
        } else {
          setFeedback(errorFeedback(result));
        }
      });
    },
    [document, context, fetcher, selection, reloadList, t],
  );

  // ADR-0012 upload transport: resolve action (actionRef → page action, else
  // direct URL), validate + upload through the shared orchestrator.
  const uploadFiles = useCallback(
    async (field: FormControlField, files: UploadableFile[]): Promise<unknown> => {
      let action: UploadActionResult;
      if (typeof field.actionRef === "string" && field.actionRef !== "") {
        const referenced = actionOf(document, field.actionRef);
        if (referenced === undefined || referenced.type !== "upload") {
          throw new Error(`upload action "${field.actionRef}" is not a type=upload action`);
        }
        action = referenced as unknown as UploadActionResult;
      } else if (typeof field.action === "string" && field.action !== "") {
        // Direct-URL mode: constraints come from the field props (registry).
        action = {
          url: field.action,
          accept: field.accept,
          maxSize: field.maxSize,
          multiple: field.multiple === true,
        };
      } else {
        throw new Error("upload field requires action or actionRef");
      }
      const result = await uploadFilesWithFetch(action, files, fetcher);
      if (!result.ok) {
        throw new Error(`${result.code} (file ${result.fileIndex})`);
      }
      return result.fieldValue;
    },
    [document, fetcher],
  );

  const resolveConfirm = useCallback(async (confirmed: boolean) => {
      if (pendingConfirm === null) {
        return;
      }
      const { actionRef, actionKey, row, requestMapping, batch } = pendingConfirm;
      setPendingConfirm(null);
      if (!confirmed) {
        return;
      }
      if (batch !== undefined) {
        const item = { actionRef, key: actionKey, batchMapping: pendingConfirm.batchMapping };
        const result = await runBatchRequest(document, context, fetcher, actionRef, item, batch.selection);
        if (result.ok) {
          setFeedback({ kind: "success", message: batchSuccessMessageFor(actionOf(document, actionRef) ?? {}, t) });
          reloadList();
        } else {
          setFeedback(errorFeedback(result));
        }
        return;
      }
      await runRowAction(actionRef, { row, requestMapping, gateTargetId: actionKey, confirmed: true });
    },
    [pendingConfirm, runRowAction, document, context, fetcher, reloadList, t],
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
        setFeedback({ kind: "success", message: successMessageFor(actionOf(document, submitAction)?.method, t) });
        reloadList();
        setActiveModal(null);
        setSelectedRow(null);
      }
      return result;
    },
    [runRequestCallback, document, activeModal, reloadList, t],
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
      selection,
      setSelection,
      clearSelection,
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
      invokeBatchAction,
      uploadFiles,
      submitForm,
      searchFormSubmit,
      effectivePermission,
    }),
    [
      selectedRow,
      tableQuery,
      setTableQuery,
      selection,
      setSelection,
      clearSelection,
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
      invokeBatchAction,
      uploadFiles,
      submitForm,
      searchFormSubmit,
      effectivePermission,
    ],
  );

  return <SchemaCrudContext.Provider value={value}>{children}</SchemaCrudContext.Provider>;
}

function errorFeedback(result: {
  code: string;
  message: string;
  messageKey?: string;
  params?: Record<string, unknown>;
}): SchemaCrudFeedback {
  return {
    kind: "error",
    code: result.code,
    message: result.message,
    ...(result.messageKey === undefined ? {} : { messageKey: result.messageKey }),
    ...(result.params === undefined ? {} : { params: result.params }),
  };
}

function FeedbackRegion({ feedback }: { feedback: SchemaCrudFeedback }) {
  const t = useTranslate();
  // VP-007 S4 frontend floor: render the catalog entry by key/params when the
  // catalog has it (current locale → en-US); otherwise the server message.
  let text = feedback.message;
  if (feedback.messageKey !== undefined && feedback.messageKey !== "") {
    const localized = t(feedback.messageKey, feedback.params as MessageParams | undefined);
    if (localized !== feedback.messageKey) {
      text = localized;
    }
  }
  return (
    <div
      role={feedback.kind === "error" ? "alert" : "status"}
      className={`mb-4 rounded-md border px-3 py-2 text-sm ${
        feedback.kind === "error"
          ? "border-destructive/50 bg-destructive/10 text-destructive"
          : "border-success/50 bg-success/10 text-success"
      }`}
    >
      {feedback.code !== undefined && feedback.code !== "" ? `${feedback.code}: ` : ""}
      {text}
    </div>
  );
}

type RecordSourcePrefillState =
  | { status: "idle" }
  | { status: "loading" }
  | { status: "error"; message: string }
  | { status: "ready"; values: Record<string, unknown> };

function hasRequiredCapability(metaValue: unknown, capability: string): boolean {
  return (
    isRecord(metaValue) &&
    Array.isArray(metaValue.requiredCapabilities) &&
    (metaValue.requiredCapabilities as unknown[]).includes(capability)
  );
}

/**
 * ADR-0021 `form.props.recordSource` prefill (S6): loads the record from a
 * detail GET and maps it to field ids via `responseMapping` (field → dot-path).
 * Fails closed on missing capability / invalid shape / network error — never
 * renders an editable blank form that could overwrite the record.
 */
function useRecordSourcePrefill(
  node: RenderFormNode,
  metaValue: unknown,
  crud: SchemaCrudValue | null,
): RecordSourcePrefillState {
  const recordSource = node.props.recordSource;
  const [state, setState] = useState<RecordSourcePrefillState>({ status: "idle" });
  useEffect(() => {
    if (recordSource === undefined) {
      setState({ status: "idle" });
      return;
    }
    if (node.props.mode === "search") {
      setState({ status: "error", message: "form.recordSource is forbidden on search-mode forms" });
      return;
    }
    if (!hasRequiredCapability(metaValue, FORM_RECORD_LOAD_CAPABILITY)) {
      setState({
        status: "error",
        message: `form.recordSource requires capability "${FORM_RECORD_LOAD_CAPABILITY}" in meta.requiredCapabilities`,
      });
      return;
    }
    const constructed = constructRequest({
      kind: "recordSource",
      recordSource,
      route: { params: {}, query: {} },
      baseURL: "",
    });
    if (!constructed.ok) {
      setState({
        status: "error",
        message: `recordSource construction failed (${constructed.path}): ${constructed.code}`,
      });
      return;
    }
    const url = constructed.request?.url;
    if (typeof url !== "string") {
      setState({ status: "error", message: "recordSource produced no request URL" });
      return;
    }
    let cancelled = false;
    setState({ status: "loading" });
    const fetcher = crud?.fetcher ?? globalThis.fetch;
    fetcher(url)
      .then(async (response) => {
        if (cancelled) {
          return;
        }
        if (!response.ok) {
          const apiError = await readResourceApiError(response, "recordSource");
          if (!cancelled) {
            setState({ status: "error", message: `${apiError.code}: ${apiError.message}` });
          }
          return;
        }
        const record: unknown = await response.json();
        const raw = isRecord(record) ? record : {};
        const mapping = isRecord(recordSource.responseMapping)
          ? recordSource.responseMapping
          : {};
        const values: Record<string, unknown> = {};
        for (const [fieldId, pathExpr] of Object.entries(mapping)) {
          if (typeof pathExpr !== "string") {
            continue;
          }
          const resolved = resolveResponsePath(raw, pathExpr);
          if (resolved !== undefined) {
            values[fieldId] = resolved;
          }
        }
        if (!cancelled) {
          setState({ status: "ready", values });
        }
      })
      .catch((error: unknown) => {
        if (cancelled) {
          return;
        }
        setState({
          status: "error",
          message: error instanceof Error ? error.message : String(error),
        });
      });
    return () => {
      cancelled = true;
    };
  }, [recordSource, node.props.mode, metaValue, crud?.fetcher, crud?.reloadToken]);
  return state;
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
  const crud = useSchemaCrud();
  const t = useTranslate();
  const prefill = useRecordSourcePrefill(node, metaValue, crud);

  if (prefill.status === "loading") {
    return (
      <div
        role="status"
        aria-label={t("feedback.loading")}
        className="space-y-2 rounded-md border border-border bg-card p-4"
      >
        <Skeleton className="h-4 w-40" />
        <Skeleton className="h-9 w-full" />
        <Skeleton className="h-9 w-full" />
      </div>
    );
  }
  if (prefill.status === "error") {
    return (
      <p role="alert" className="text-sm text-destructive">
        {prefill.message}
      </p>
    );
  }
  return (
    <FormInner
      // Remount on reload so recordSource forms re-initialize from the fresh
      // record after a save / batch reload (no stale typed values).
      key={node.props.recordSource !== undefined ? crud?.reloadToken : undefined}
      node={node}
      metaValue={metaValue}
      context={context}
      formComponent={formComponent}
      prefillValues={prefill.status === "ready" ? prefill.values : null}
    />
  );
}

function FormInner({
  node,
  metaValue,
  context,
  formComponent,
  prefillValues,
}: {
  node: RenderFormNode;
  metaValue: unknown;
  context: Record<string, unknown>;
  formComponent?: RendererComponentProps["formComponent"];
  prefillValues: Record<string, unknown> | null;
}) {
  const Component = formComponent ?? FormControls;
  const crud = useSchemaCrud();
  const t = useTranslate();
  const modalRow = crud?.modalRow ?? null;
  const gate = gateRenderFormFields(metaValue, node.props.fields, "fields");
  const [values, setValues] = useState<Record<string, unknown>>(() => {
    const initial: Record<string, unknown> = {};
    for (const raw of node.props.fields) {
      if (!isRecord(raw) || typeof raw.id !== "string") {
        continue;
      }
      // Precedence: modal row (edit-in-modal) → recordSource prefill → empty.
      // Always coerce through the field wire kind so row arrays (e.g. roles[])
      // become textarea strings or checkboxGroup string[] as appropriate (A-006 R-004).
      const modalValue =
        modalRow !== null && modalRow[raw.id] !== undefined ? modalRow[raw.id] : undefined;
      const prefillValue =
        prefillValues !== null && prefillValues[raw.id] !== undefined
          ? prefillValues[raw.id]
          : undefined;
      const fromRow = modalValue !== undefined ? modalValue : prefillValue;
      initial[raw.id] = coerceFieldValue(raw as unknown as FormControlField, fromRow);
    }
    return initial;
  });
  // Baseline snapshot for the full $deps engine (02 §14: baseline = the
  // field's value at first mount unless explicitly overridden).
  const mountBaselines = useRef(values);
  // Full engine (upstream per-field reactions) runs on every value change and
  // converges; the frozen $context engine remains the fallback.
  const fullReaction = useMemo(
    () => resolveFullFormReactions(node.props.fields, values, mountBaselines.current),
    [node.props.fields, values],
  );
  const reaction = fullReaction.usesFullEngine
    ? fullReaction
    : resolveFormReactions(node, context);
  const reactionState: FormControlStateMap = fullReaction.usesFullEngine
    ? fullReaction.state
    : reaction.state;
  const reactionErrors: ReactionError[] = fullReaction.usesFullEngine
    ? fullReaction.errors
    : reaction.errors;

  useEffect(() => {
    if (!fullReaction.usesFullEngine) {
      return;
    }
    const entries = Object.entries(fullReaction.values);
    if (entries.length === 0) {
      return;
    }
    setValues((prev) => {
      const next = { ...prev };
      let changed = false;
      for (const [key, value] of entries) {
        if (!Object.is(prev[key], value)) {
          next[key] = value;
          changed = true;
        }
      }
      return changed ? next : prev;
    });
  }, [fullReaction]);

  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<{
    code: string;
    message: string;
    messageKey?: string;
    params?: Record<string, unknown>;
  } | null>(null);

  const isSearch = node.props.mode === "search";
  const submitAction = node.props.submitAction;
  const canSubmit = isSearch || typeof submitAction === "string";
  const hasBlockingErrors = gate.errors.length > 0 || reactionErrors.length > 0;
  // Permission gate: a default-mode form with a declared submit permission is
  // read-only when the `${formId}:submit` target is denied — the viewer sees
  // current values but cannot save (backend `settings.write` stays the hard
  // gate). Absent a declared target the effective permission defaults to true.
  const canEdit = crud === null ? true : crud.effectivePermission(`${node.id ?? "unnamed"}:submit`);

  const visibleFields = gate.fields.filter(
    (raw) => reactionState[raw.id]?.visible !== false,
  );

  const fieldDisabled = (id: string) =>
    reactionState[id]?.disabled === true || (!isSearch && !canEdit);

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
      setFormError({
        code: result.code,
        message: result.message,
        ...(result.messageKey === undefined ? {} : { messageKey: result.messageKey }),
        ...(result.params === undefined ? {} : { params: result.params }),
      });
    }
    setSubmitting(false);
  };

  const title = resolveTextProp(
    node.props as unknown as Record<string, unknown>,
    "titleKey",
    "title",
    t,
    "",
  );

  return (
    <form className="space-y-3" onSubmit={(event) => {
      event.preventDefault();
      void handleSubmit();
    }}>
      {title !== "" ? (
        <h2 className="text-lg font-semibold tracking-tight text-foreground">{title}</h2>
      ) : null}
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
        onUpload={crud?.uploadFiles}
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
          {formError.code}:{" "}
          {formError.messageKey !== undefined && formError.messageKey !== ""
            ? t(formError.messageKey, formError.params as MessageParams | undefined)
            : formError.message}
        </p>
      ) : null}
      {canSubmit ? (
        <button
          type="submit"
          disabled={submitting || hasBlockingErrors || (!isSearch && !canEdit)}
          className="h-9 rounded-md bg-primary px-3 text-sm font-medium text-primary-foreground transition-opacity hover:opacity-90 disabled:opacity-50"
        >
          {submitting
            ? t("feedback.submitting")
            : resolveTextProp(
                node.props as unknown as Record<string, unknown>,
                "submitLabelKey",
                "submitLabel",
                t,
                isSearch ? t("feedback.search") : t("feedback.submit"),
              )}
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
  // Mobile: always single column; md+: honor schema columns (fluid via CSS + minmax(0,1fr)).
  const desktopCols =
    columns !== undefined && columns > 1
      ? Math.min(Math.max(Math.floor(columns), 1), 6)
      : undefined;
  return (
    <div
      className="grid w-full min-w-0 gap-4 grid-cols-1"
      data-grid-responsive="true"
      data-schema-grid-cols={desktopCols !== undefined ? String(desktopCols) : undefined}
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
  const t = useTranslate();
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
        <p className="text-sm text-muted-foreground">{t("feedback.tabsNoChildren")}</p>
      )}
    </div>
  );
}

function TextView({ node }: { node: RenderTextNode }) {
  const t = useTranslate();
  return (
    <p className="text-sm text-foreground">
      {resolveTextProp(node.props as unknown as Record<string, unknown>, "textKey", "text", t, "")}
    </p>
  );
}

/**
 * recordView surface (D-004 §5 / S2):
 * - Empty: inline placeholder
 * - With record: desktop right Drawer + mobile full-height Sheet (not centered Modal)
 * Static `props.record` still uses the same chrome so fixtures remain observable.
 */
function RecordView({ node }: { node: RenderRecordViewNode }) {
  const crud = useSchemaCrud();
  const t = useTranslate();
  const staticRecord = node.props?.record;
  const hasStatic = isRecord(staticRecord);
  const record = hasStatic ? staticRecord : (crud?.selectedRow ?? null);
  const entries = isRecord(record) ? Object.entries(record) : [];
  if (entries.length === 0) {
    return (
      <p
        data-record-view="empty"
        className="text-sm text-muted-foreground"
      >
        {t("feedback.selectRecordToView")}
      </p>
    );
  }

  const canClose = !hasStatic && crud !== null;
  const onClose = () => {
    if (canClose) {
      crud.selectRow(null);
    }
  };

  return (
    <>
      {/* Dimmer only when selection-driven (static fixtures stay non-modal). */}
      {canClose ? (
        <div
          data-record-view="backdrop"
          className="fixed inset-0 z-40 bg-overlay"
          aria-hidden="true"
          onClick={onClose}
        />
      ) : null}
      <aside
        data-record-view="panel"
        data-record-view-mode={canClose ? "drawer" : "panel"}
        role="dialog"
        aria-modal={canClose ? true : undefined}
        aria-label={t("feedback.recordDetails")}
        className={
          canClose
            ? // Desktop (md+): right Drawer; mobile (<768, D-004): full-height Sheet
              "fixed inset-y-0 right-0 z-50 flex w-full max-w-md flex-col border-l border-border bg-card shadow-lg max-md:inset-x-0 max-md:top-auto max-md:h-[min(92vh,100%)] max-md:rounded-t-xl max-md:border-l-0 max-md:border-t md:max-w-md"
            : "w-full max-w-md rounded-lg border border-border bg-card shadow-sm"
        }
      >
        <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3">
          <h2 className="text-sm font-semibold tracking-tight text-foreground">
            {t("feedback.recordDetails")}
          </h2>
          {canClose ? (
            <button
              type="button"
              aria-label={t("feedback.closeRecordDetails")}
              onClick={onClose}
              className="inline-flex size-8 items-center justify-center rounded-md text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
            >
              ×
            </button>
          ) : null}
        </div>
        <dl className="flex-1 space-y-3 overflow-y-auto p-4">
          {entries.map(([key, value]) => (
            <div key={key} className="grid gap-0.5 text-sm sm:grid-cols-[8rem_1fr] sm:gap-3">
              <dt className="text-xs font-medium uppercase tracking-[0.08em] text-muted-foreground">
                {key}
              </dt>
              <dd className="break-words text-foreground">{String(value)}</dd>
            </div>
          ))}
        </dl>
      </aside>
    </>
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
    // A-002 finding 3: clear any stale error before a refetch starts, and on
    // success — otherwise a failed fetch followed by a successful reload
    // (same dataSource, new reloadToken) keeps showing the old error because
    // resolveAsyncDisplayState prefers `error` over `ready`.
    setError(null);
    fetchResourceList(fetcher ?? fetch, dataSource, { page: 1, pageSize: 100 })
      .then((next) => {
        if (!cancelled) {
          setList(next);
          setError(null);
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
  const t = useTranslate();
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
  const displayState = resolveAsyncDisplayState({ loading: list === null, error });
  if (displayState === "error") {
    return <p role="alert" className="text-sm text-destructive">statCard data failed to load: {error}</p>;
  }
  if (displayState === "loading") {
    return (
      <div role="status" aria-label={t("feedback.loadingStatCard")} className="space-y-2 rounded-md border border-border bg-card p-4">
        <Skeleton className="h-3 w-16" />
        <Skeleton className="h-7 w-24" />
      </div>
    );
  }
  const readyList = list!;
  const first = readyList.items[0];
  // Registry: valueField = "指定从 API 响应中取哪个字段作为展示值". The list
  // envelope (total/page/pageSize) is part of the response, so envelope fields
  // resolve too; row fields take precedence for display values.
  const raw =
    valueField !== undefined
      ? (first !== undefined && first[valueField] !== undefined
          ? first[valueField]
          : valueField === "total"
            ? readyList.total
            : valueField === "page"
              ? readyList.page
              : valueField === "pageSize"
                ? readyList.pageSize
                : 0)
      : 0;
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
  const label = resolveTextProp(
    node.props as unknown as Record<string, unknown>,
    "labelKey",
    "label",
    t,
    valueField ?? t("feedback.value"),
  );
  const unit = node.props?.unit;
  return (
    <Card data-display-surface="statCard" className="shadow-sm">
      <CardContent className="p-4">
        <p className="text-xs font-medium text-muted-foreground">{label}</p>
        <p className="mt-1 text-2xl font-semibold tracking-tight text-foreground">
          {String(formatResult.value)}
          {unit !== undefined && unit !== "" ? (
            <span className="ml-1 text-sm font-normal text-muted-foreground">{unit}</span>
          ) : null}
        </p>
      </CardContent>
    </Card>
  );
}

function ChartView({ node }: { node: RenderChartNode }) {
  const crud = useSchemaCrud();
  const t = useTranslate();
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
  const chartDisplayState = resolveAsyncDisplayState({ loading: list === null, error });
  if (chartDisplayState === "error") {
    return <p role="alert" className="text-sm text-destructive">chart data failed to load: {error}</p>;
  }
  if (chartDisplayState === "loading") {
    return (
      <div role="status" aria-label={t("feedback.loadingChart")} className="space-y-2 rounded-md border border-border bg-card p-4">
        <Skeleton className="h-40 w-full" />
      </div>
    );
  }
  const points = list!.items
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
                  stroke={`var(--color-chart-${(index % 5) + 1})`}
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
  const t = useTranslate();
  const crud = useSchemaCrud();
  const gate = tableActionGate(node.props ?? {}, context);
  if (!gate.visible) {
    return null;
  }
  // ADR-0023 D4b: actionButton permission-intent targets are gated by their
  // `props.key` (fallback node id); a denied target renders the button disabled.
  const targetId =
    typeof node.props?.key === "string" && node.props.key !== ""
      ? node.props.key
      : node.id;
  const canAct =
    crud === null || targetId === undefined
      ? true
      : crud.effectivePermission(targetId);
  return (
    <div className="space-y-1">
      <button
        type="button"
        disabled={gate.disabled || !canAct}
        onClick={() => onAction?.(node)}
        className="h-9 rounded-md border border-input bg-background px-3 text-sm transition-colors focus-visible:ring-2 focus-visible:ring-ring disabled:opacity-50"
      >
        {resolveTextProp(
          node.props as unknown as Record<string, unknown>,
          "labelKey",
          "label",
          t,
          t("feedback.action"),
        )}
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
  // VP-007 S3: actionButton nodes are first-class page actions in the default
  // app path — dispatch through the frozen Schema CRUD executor (gate →
  // confirm → request) unless the host overrides onAction.
  const resolvedOnAction = onAction ?? ((node: RenderActionButtonNode) => {
    crud.invokeAction(node.props as unknown as Record<string, unknown>, null);
  });
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
        onAction: resolvedOnAction,
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
            onAction: resolvedOnAction,
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
