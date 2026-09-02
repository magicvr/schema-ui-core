import { RotateCcw, Search, X } from "lucide-react";
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
import { formatDisplayTime } from "@/lib/datetime";
import {
  constructRequest,
  normalizeSelection,
  type RequestConstructionResult,
} from "@/protocol/conformance/request-construction";
import {
  uploadFilesWithFetch,
  type UploadActionResult,
  type UploadableFile,
} from "@/protocol/conformance/upload-orchestration";
import { ConfirmDialog } from "@/renderer/confirm";
import { FormControls } from "@/renderer/form-controls.tsx";
import { getCustomComponent } from "@/renderer/custom-components";
import {
  FORM_RECORD_LOAD_CAPABILITY,
  coerceFieldValue,
  validateFieldValues,
  type FormControlField,
} from "@/renderer/form-controls.types";
import { optionList } from "@/renderer/form-controls.tsx";
import { ModalHost } from "@/renderer/modal";
import { resolveTextProp, type MessageParams } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
import {
  executeAction,
  evaluatePermissionTargets,
  validatePermissions,
} from "@/renderer/permissions";
import {
  DISPLAY_LIST_QUERY,
  EMPTY_RESOURCE_LIST,
  fetchResourceList,
  isValidDataSource,
  isWalletNotFoundError,
  readResourceApiError,
  resolveDataParamsQuery,
  resourceListURL,
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
} from "@/renderer/render.types";

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
  /**
   * Session-internal navigation hook (ADR-0021 navigate actions; GOAL-015
   * F-001): the host pushes the target onto its own history/visit stack so
   * breadcrumbs survive. Falls back to window.location.assign when absent.
   */
  onNavigate?: (url: string) => void;
  /** Overrides the default FormControls component (keeps field wiring local). */
  formComponent?: ComponentType<{
    fields: FormControlField[];
    values: Record<string, unknown>;
    onChange: (id: string, value: unknown) => void;
    fieldDisabled?: (id: string) => boolean;
    onUpload?: (field: FormControlField, files: UploadableFile[]) => Promise<unknown>;
    /** W11 · U-01/U-02: auth-aware transport for dynamic option sources. */
    fetcher?: typeof fetch;
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
  | { ok: true; fieldErrors?: Array<{ field: string; reason: string; rowNumber?: number }>; message?: string; messageKey?: string; data?: unknown }
  | {
      ok: false;
      code: string;
      message: string;
      messageKey?: string;
      params?: Record<string, unknown>;
      /** GOAL-014 D-002 §2: server field-level validation failures. */
      fieldErrors?: Array<{ field: string; reason: string; rowNumber?: number }>;
    };

export interface SchemaCrudValue {
  selectedRow: Record<string, unknown> | null;
  selectRow: (row: Record<string, unknown> | null) => void;
  tableQuery: (id: string) => ResourceQuery | undefined;
  setTableQuery: (id: string, query: ResourceQuery) => void;
  reloadToken: number;
  reloadList: () => void;
  /**
   * Fetches a resource list through the page-level in-flight coalescer:
   * simultaneous consumers of the same URL share ONE network request (three
   * statCards + the wallet-ensure probe on one "我的钱包" visit merge into a
   * single GET /me). Requests are otherwise never memoized — every query,
   * reset or reload refetches — and reloadList drops the in-flight map so a
   * reload issued during a slow fetch starts its own fresh request.
   * `transport` is the caller's own transport (a directly injected fixture or
   * the auth fetcher); defaults to the provider's registered fetcher.
   */
  fetchList: (
    dataSource: string,
    query: ResourceQuery,
    extraQuery?: string,
    transport?: typeof fetch,
  ) => Promise<ResourceList>;
  /**
   * Targeted display-data refresh (W25): bumps the per-URL refresh token for
   * the standard display query of `dataSource` (statCard/chart consume it),
   * so a consumer can refetch ONE surface without a full-page reload wave.
   * Only applies to display nodes without route-param bindings (the standard
   * DISPLAY_LIST_QUERY shape); tables/misc surfaces keep their data until a
   * manual reload.
   */
  refreshList: (dataSource: string) => void;
  /** Current refresh token for a display dataSource (0 when untouched). */
  listRefreshToken: (dataSource: string) => number;
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
  /**
   * Current route snapshot (from the render context; App injects
   * route: {params, query}). Used for ADR-0039 dataSource bindings and
   * create-modal readOnly seeding instead of reading window.location.
   */
  route: { query: Record<string, string>; params: Record<string, string> };
}

export interface RunRequestOptions {
  row?: Record<string, unknown> | null;
  formValues?: Record<string, unknown>;
  requestMapping?: Record<string, unknown>;
  gateTargetId?: string;
  confirmed?: boolean;
}

/** Page CRUD context consumed by custom components (rendered nodes). */
export const SchemaCrudContext = createContext<SchemaCrudValue | null>(null);

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

/**
 * F-02 (GOAL-004 D-002 §5): whitelisted custom-action handlers. Export
 * handlers fetch the CSV with the authed transport and trigger a browser
 * download (blob + anchor); the filename derives from the handler name.
 */
const CUSTOM_HANDLER_URLS: Record<string, string> = {
  "export.users": "/api/export/users",
  "export.roles": "/api/export/roles",
  // S-02 (GOAL-007 D-002 §5): file-library row download. The {id} slot is
  // resolved from the captured row context (bounded binding, same posture as
  // the formAction {id} slots, I-007-003 §9.1a).
  "library.download": "/api/library/files/{id}/download",
  // W16-F02 (GOAL-026): preview opens the same authed URL in a new tab;
  // copyLink writes the URL to the clipboard.
  "library.preview": "/api/library/files/{id}/download",
  "library.copyLink": "/api/library/files/{id}/download",
};

async function runCustomAction(
  action: JsonRecord,
  fetcher: typeof fetch,
  row?: Record<string, unknown> | null,
): Promise<ActionResult> {
  const handler = stringOf(action.handler);
  let url = CUSTOM_HANDLER_URLS[handler];
  if (url === undefined) {
    return { ok: false, code: "CUSTOM_HANDLER_NOT_FOUND", message: "custom handler not whitelisted: " + handler, messageKey: "error.customHandlerNotFound", params: { handler } };
  }
  // Row-scoped custom handlers resolve the {id} slot from the captured row
  // context (S-02, GOAL-007 D-002 §5). A missing row id on a templated
  // handler fails closed.
  const rowId = row?.id;
  if (url.includes("{id}")) {
    if (typeof rowId !== "string" || rowId === "") {
      return { ok: false, code: "CUSTOM_HANDLER_MISSING_ROW_ID", message: "custom handler requires a row id: " + handler, messageKey: "error.customHandlerMissingRowId", params: { handler } };
    }
    url = url.replaceAll("{id}", encodeURIComponent(rowId));
  }
  // W16-F02: preview/copy fetch the file through the authed transport as a
  // blob, then open/copy the blob object URL. This keeps content accessible
  // without leaking the Bearer token into a raw download URL (the download
  // endpoint itself is attachment + bearer-gated, so a bare window.open 401s).
  // W18: copy the origin-absolute download path (session-gated). A blob:
  // URL is not pasteable outside this page (GOAL-024 A-007 F-001).
  if (handler === "library.copyLink") {
    const absolute = new URL(url, window.location.origin).href;
    try {
      await navigator.clipboard.writeText(absolute);
      return { ok: true };
    } catch {
      return { ok: false, code: "CLIPBOARD_UNAVAILABLE", message: "clipboard is unavailable", messageKey: "error.clipboardUnavailable" };
    }
  }
  if (handler === "library.preview") {
    // Open synchronously in the user gesture, then navigate the blank tab
    // after the authed blob fetch — a post-await window.open is blocked.
    const previewWindow = window.open("about:blank", "_blank");
    let blob: Blob;
    try {
      const response = await fetcher(url, { method: "GET", headers: { Accept: "*/*" } });
      if (!response.ok) {
        previewWindow?.close();
        const apiError = await readResourceApiError(response, handler);
        return { ok: false, code: apiError.code, message: apiError.message };
      }
      blob = await response.blob();
    } catch (error) {
      previewWindow?.close();
      return { ok: false, code: "REQUEST_FAILED", message: requestFailedMessage(error) };
    }
    const objectUrl = URL.createObjectURL(blob);
    if (previewWindow === null || previewWindow.closed) {
      URL.revokeObjectURL(objectUrl);
      return { ok: false, code: "POPUP_BLOCKED", message: "preview window was blocked", messageKey: "error.popupBlocked" };
    }
    // A-003 recommended F-003: cut the opener link as soon as we hold the
    // direct reference — if the preview template ever embeds untrusted
    // content, the opened document can no longer reach back into this app.
    previewWindow.opener = null;
        // W7 F-010: do not navigate the preview tab directly to a blob: URL — that
    // strips the server Content-Disposition/CSP sandbox headers. Embed the blob
    // in a sandboxed (no allow-scripts, no allow-same-origin) iframe instead.
    const previewDocument = previewWindow.document;
    previewDocument.open();
    previewDocument.write(
      "<!doctype html><meta charset=\"utf-8\"><title>Preview</title>" +
        "<style>html,body{margin:0;height:100%}iframe{width:100vw;height:100vh;border:0;display:block}</style>" +
        '<iframe sandbox="" src="' + objectUrl + '"></iframe>',
    );
    previewDocument.close();
    window.setTimeout(() => URL.revokeObjectURL(objectUrl), 60_000);
    return { ok: true };
  }
  let response: Response;
  try {
    response = await fetcher(url, { method: "GET", headers: { "Content-Type": "application/json" } });
  } catch (error) {
    return { ok: false, code: "REQUEST_FAILED", message: requestFailedMessage(error) };
  }
  if (!response.ok) {
    const apiError = await readResourceApiError(response, handler);
    return { ok: false, code: apiError.code, message: apiError.message };
  }
  const blob = await response.blob();
  // The download filename prefers the row's stored name, scrubbed with the
  // server allowlist shape ([A-Za-z0-9._-], trimmed, capped — A-003 F-001);
  // the export fallback keeps the historical "<handler>.csv" shape.
  const rowName = row?.name;
  const rawName = typeof rowName === "string" ? rowName.trim() : "";
  const filename =
    rawName !== ""
      ? sanitizeClientFilename(rawName)
      : handler.split(".").pop() + ".csv";
  triggerBlobDownload(blob, filename);
  return { ok: true };
}

/**
 * S-02 (GOAL-007 A-003 F-001): client download names mirror the server
 * allowlist — path separators, controls and everything outside [A-Za-z0-9._-]
 * become underscores, leading/trailing separators are trimmed, and the result
 * is capped so a path-like or control-bearing stored name can never reach the
 * local filesystem unscrubbed.
 */
function sanitizeClientFilename(name: string): string {
  let out = name.replace(/[^A-Za-z0-9._-]/g, "_").replace(/^[._-]+|[._-]+$/g, "");
  if (out === "") {
    out = "download";
  }
  if (out.length > 100) {
    out = out.slice(0, 100);
  }
  return out;
}

/** Triggers a browser download from a fetched blob (F-02 local extension). */
function triggerBlobDownload(blob: Blob, filename: string): void {
  const objectUrl = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = objectUrl;
  anchor.download = filename;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(objectUrl);
}

async function runRequest(
  document: RenderPageDocument,
  context: Record<string, unknown>,
  fetcher: typeof fetch,
  actionRef: string,
  opts: RunRequestOptions,
): Promise<ActionResult> {
  const action = actionOf(document, actionRef);
  if (action === undefined) {
    return { ok: false, code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page`, messageKey: "error.actionNotFound", params: { action: actionRef } };
  }
  // F-02 (GOAL-004 D-002 §5): local custom-action dispatch — the protocol's
  // CustomAction extension point (action.schema.json): a schema action may
  // reference a whitelisted handler name; the renderer resolves it locally.
  // Unknown handler names fail closed (CUSTOM_HANDLER_NOT_FOUND).
  if (action.type === "custom") {
    return runCustomAction(action, fetcher, opts.row ?? null);
  }
  if (action.type !== "request") {
    return { ok: false, code: "ACTION_NOT_REQUEST", message: `action "${actionRef}" is not a request action`, messageKey: "error.actionNotRequest", params: { action: actionRef } };
  }
  if (opts.gateTargetId !== undefined) {
    // Absent target = no declared permission entry (engine default is allow):
    // only gate when the page actually declares a permission for this target,
    // matching the batch path (ADR-0022 D5d) and effectivePermission (C7).
    const hasPermissionEntry = evaluatePermissionTargets(
      document as unknown as JsonRecord,
      context as NavigationContext,
    ).some((entry) => entry.targetId === opts.gateTargetId);
    if (hasPermissionEntry) {
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
          messageKey: "error.actionNotExecuted",
          params: { action: actionRef, reason: gate.reason ?? gate.outcome },
        };
      }
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
  // W4 P1-1: constructRequest can throw on malformed schema values (e.g.
  // serializeQueryValue rejects non-scalars). The UI call sites must never let
  // that escape as an unhandled rejection / white screen — return a failed
  // ActionResult so the caller surfaces it through the normal feedback channel.
  let constructed: RequestConstructionResult;
  try {
    constructed = constructRequest(input);
  } catch (error) {
    return {
      ok: false,
      code: "REQUEST_CONSTRUCTION_FAILED",
      message: error instanceof Error ? error.message : "request construction failed",
    };
  }
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
  let response: Response;
  try {
    response = await fetcher(url, {
      method: request.method,
      headers: { "Content-Type": "application/json" },
      ...(body === undefined ? {} : { body }),
    });
  } catch (error) {
    // Network-level failure (offline, server down, CORS): surface as an action
    // result so every caller shows feedback instead of an unhandled rejection.
    return { ok: false, code: "REQUEST_FAILED", message: requestFailedMessage(error) };
  }
  if (!response.ok) {
    const apiError = await readResourceApiError(response, actionRef);
    return {
      ok: false,
      code: apiError.code,
      message: apiError.message,
      ...(apiError.messageKey === undefined ? {} : { messageKey: apiError.messageKey }),
      ...(apiError.params === undefined ? {} : { params: apiError.params }),
      ...(apiError.fieldErrors.length > 0 ? { fieldErrors: apiError.fieldErrors } : {}),
    };
  }


  // Branding refresh after settings PATCH lives in the App/host layer (A-006 R-002),
  // not in this generic request executor — keep Renderer free of product endpoints.
  // W16-F03: a 200 import response may still carry `fieldErrors` (partial
  // failure) — surface them to the form instead of pretending full success.
  let successFieldErrors: Array<{ field: string; reason: string }> | undefined;
  let parsedData: unknown = undefined;
  try {
    const text = await response.text();
    if (text !== "") {
      const parsed = JSON.parse(text) as { fieldErrors?: unknown };
      parsedData = parsed;
      if (Array.isArray(parsed.fieldErrors)) {
        const cleaned = parsed.fieldErrors
          .filter((entry): entry is { field?: unknown; reason?: unknown; rowNumber?: unknown } =>
            typeof entry === "object" && entry !== null && !Array.isArray(entry),
          )
          .map((entry) => ({
            field: typeof entry.field === "string" ? entry.field : "",
            reason: typeof entry.reason === "string" ? entry.reason : "",
            ...(typeof entry.rowNumber === "number" ? { rowNumber: entry.rowNumber } : {}),
          }))
          .filter((entry) => entry.field !== "" || entry.reason !== "");
        if (cleaned.length > 0) {
          successFieldErrors = cleaned;
        }
      }
    }
  } catch {
    // non-JSON success body (e.g. 204/CSV) — nothing to surface
  }
  return successFieldErrors !== undefined
    ? { ok: true, fieldErrors: successFieldErrors, data: parsedData }
    : { ok: true, data: parsedData };
}

function requestFailedMessage(error: unknown): string {
  if (error instanceof Error && error.message !== "") {
    return error.message;
  }
  return "request failed (network error)";
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
    return { ok: false, code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page`, messageKey: "error.actionNotFound", params: { action: actionRef } };
  }
  if (action.type !== "request") {
    return { ok: false, code: "ACTION_NOT_REQUEST", message: `action "${actionRef}" is not a request action`, messageKey: "error.actionNotRequest", params: { action: actionRef } };
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
        messageKey: "error.batchActionNotExecuted",
        params: { action: actionRef, reason: gate.reason ?? gate.outcome },
      };
    }
  }
  const batchMapping = isRecord(item.batchMapping) ? item.batchMapping : undefined;
  // W4 P1-1: constructRequest can throw on malformed schema values; never let
  // that escape as an unhandled rejection — return a failed result instead.
  let constructed: RequestConstructionResult;
  try {
    constructed = constructRequest({
      kind: "batchRequest",
      action,
      batchMapping,
      selection: { keys: selection.keys, count: selection.count },
    });
  } catch (error) {
    return {
      ok: false,
      code: "REQUEST_CONSTRUCTION_FAILED",
      message: error instanceof Error ? error.message : "batch request construction failed",
    };
  }
  if (!constructed.ok) {
    return { ok: false, code: constructed.code, message: `batch request construction failed (${constructed.path})` };
  }
  const request = constructed.request;
  if (request === undefined) {
    return { ok: false, code: "NO_REQUEST", message: "batch action produced no request" };
  }
  const body = request.body === null || request.body === undefined ? undefined : JSON.stringify(request.body);
  let response: Response;
  try {
    response = await fetcher(request.url, {
      method: request.method,
      headers: { "Content-Type": "application/json" },
      ...(body === undefined ? {} : { body }),
    });
  } catch (error) {
    return { ok: false, code: "REQUEST_FAILED", message: requestFailedMessage(error) };
  }
  if (!response.ok) {
    const apiError = await readResourceApiError(response, actionRef);
    return {
      ok: false,
      code: apiError.code,
      message: apiError.message,
      ...(apiError.messageKey === undefined ? {} : { messageKey: apiError.messageKey }),
      ...(apiError.params === undefined ? {} : { params: apiError.params }),
      ...(apiError.fieldErrors.length > 0 ? { fieldErrors: apiError.fieldErrors } : {}),
    };
  }
  return { ok: true };
}

function SchemaCrudProvider({
  document,
  context,
  children,
  initialFetcher,
  onNavigate,
}: {
  document: RenderPageDocument;
  context: Record<string, unknown>;
  children: ReactNode;
  initialFetcher?: typeof fetch;
  /** Session-internal navigation (breadcrumb-preserving); see RendererComponentProps. */
  onNavigate?: (url: string) => void;
}) {
  const [selectedRow, setSelectedRow] = useState<Record<string, unknown> | null>(null);
  const [queries, setQueries] = useState<Record<string, ResourceQuery>>({});
  const [selections, setSelections] = useState<Record<string, unknown[]>>({});
  const [reloadToken, setReloadToken] = useState(0);
  // Page-level in-flight coalescing (per wire URL): simultaneous consumers of
  // the same dataSource — three statCards + the wallet-ensure probe on GET
  // /api/wallet/me — share ONE network request. There is deliberately no
  // long-lived result memo: every query change / reset / reload must observe
  // fresh data, so only requests still in flight are merged. Scoped to this
  // provider (one page view); a navigation remount starts fresh.
  const listInFlight = useRef<Map<string, Promise<ResourceList>>>(new Map());
  const [activeModal, setActiveModal] = useState<{
    actionRef: string;
    row: Record<string, unknown> | null;
    title: string;
  } | null>(null);
  const [pendingConfirm, setPendingConfirm] = useState<SchemaCrudConfirm | null>(null);
  const [feedback, setFeedback] = useState<SchemaCrudFeedback | null>(null);
  const [fetcher, setFetcher] = useState<typeof fetch>(() => initialFetcher ?? globalThis.fetch);
  const t = useTranslate();

  // Route snapshot from the render context (App injects route: {params, query});
  // hostless renderers (tests) fall back to an empty snapshot.
  const route = useMemo(() => {
    const raw = isRecord(context.route) ? context.route : {};
    const query: Record<string, string> = {};
    if (isRecord(raw.query)) {
      for (const [key, value] of Object.entries(raw.query)) {
        if (typeof value === "string") query[key] = value;
      }
    }
    const params: Record<string, string> = {};
    if (isRecord(raw.params)) {
      for (const [key, value] of Object.entries(raw.params)) {
        if (typeof value === "string") params[key] = value;
      }
    }
    return { query, params };
  }, [context]);

  const permissionTargets = useMemo(
    () => evaluatePermissionTargets(document as unknown as JsonRecord, context as NavigationContext),
    [document, context],
  );

  // W9 A-005 R-F-001: the L2 permission validator now runs in the production
  // render path. A malformed permission structure surfaces at load time
  // (console.error with codes + paths) and fails CLOSED — while invalid, every
  // registered permission target is denied instead of trusting a structure the
  // validator rejected. Unmarked targets keep the protocol's default-allow.
  const permissionStructureInvalid = useMemo(() => {
    const l2Errors = validatePermissions(document as unknown as JsonRecord);
    if (l2Errors.length > 0) {
      console.error(
        "[schema-ui] permission L2 validation failed; gated targets are denied",
        l2Errors,
      );
    }
    return l2Errors.length > 0;
  }, [document]);

  const effectivePermission = useCallback(
    (targetId: string) => {
      if (permissionStructureInvalid) {
        return false;
      }
      const entry = permissionTargets.find((target) => target.targetId === targetId);
      // Absent target = no declared permission (engine default is allow).
      return entry === undefined ? true : entry.effectivePermission;
    },
    [permissionTargets, permissionStructureInvalid],
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
  // ADR-0022 D2: any data reload success clears every table selection. The
  // in-flight map is dropped too, so a reload issued DURING a slow fetch
  // starts its own fresh request instead of joining the (pre-mutation) one.
  const reloadList = useCallback(() => {
    setSelections({});
    listInFlight.current = new Map();
    setReloadToken((token) => token + 1);
  }, []);

  // A transport swap must never join requests that started through the
  // previous transport (auth identity / fixture skew).
  useEffect(() => {
    listInFlight.current = new Map();
  }, [fetcher]);

  // Per-URL refresh tokens for display data (W25): statCard/chart consume the
  // standard DISPLAY_LIST_QUERY shape, so a targeted refresh — e.g.
  // monitoring-auto-refresh refetching only its /status cards — bumps the
  // token for that exact key without a full-page reload wave.
  const [listRefreshTokens, setListRefreshTokens] = useState<Record<string, number>>({});
  const refreshList = useCallback((dataSource: string) => {
    const key = resourceListURL(dataSource, DISPLAY_LIST_QUERY, undefined);
    // Symmetric with reloadList (A-001 F-003, independent): drop any in-flight
    // request for the targeted URL so a refresh issued DURING a slow fetch
    // starts its own request instead of joining the pre-refresh one.
    listInFlight.current.delete(key);
    setListRefreshTokens((prev) => ({ ...prev, [key]: (prev[key] ?? 0) + 1 }));
  }, []);
  const listRefreshToken = useCallback(
    (dataSource: string) =>
      listRefreshTokens[resourceListURL(dataSource, DISPLAY_LIST_QUERY, undefined)] ?? 0,
    [listRefreshTokens],
  );

  const fetchList = useCallback(
    (
      dataSource: string,
      query: ResourceQuery,
      extraQuery?: string,
      transport?: typeof fetch,
    ): Promise<ResourceList> => {
      const key = resourceListURL(dataSource, query, extraQuery);
      const inFlight = listInFlight.current.get(key);
      if (inFlight !== undefined) {
        return inFlight;
      }
      const promise = fetchResourceList(transport ?? fetcher, dataSource, query, extraQuery).finally(
        () => {
          listInFlight.current.delete(key);
        },
      );
      listInFlight.current.set(key, promise);
      return promise;
    },
    [fetcher],
  );

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
        setFeedback({ kind: "error", code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page`, messageKey: "error.actionNotFound", params: { action: actionRef } });
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
      // ADR-0021 navigate actions (S-01 · GOAL-008 A-003 F-002): the host
      // executes top-level type:navigate actions referenced by rows/toolbars.
      // Only single-slash same-origin application-route paths are accepted
      // (the protocol's NavigateAction shape); anything else fails closed.
      //
      // GOAL-015 D-002 §3.2: a row action may carry navigateMapping to bind
      // path/query parameters from the row (e.g. dictKey) — the request
      // constructor's rowNavigate branch builds the final URL; without a
      // mapping the plain navigate path is used unchanged.
      if (action.type === "navigate") {
        const url = stringOf((action as JsonRecord).url);
        if (url === "" || !/^\/(?!\/)[^\s\\?#]*$/.test(url)) {
          setFeedback({
            kind: "error",
            code: "INVALID_NAVIGATE_URL",
            message: `${actionRef} has an invalid url`,
            messageKey: "error.invalidNavigateUrl",
            params: { action: actionRef },
          });
          return;
        }
        const navigateMapping = isRecord(item.navigateMapping) ? item.navigateMapping : undefined;
        if (navigateMapping !== undefined) {
          const constructed = constructRequest({
            kind: "rowNavigate",
            action: action as JsonRecord,
            navigateMapping,
            row: row ?? {},
          });
          if (!constructed.ok) {
            setFeedback({
              kind: "error",
              code: constructed.code,
              message: `row navigation construction failed (${constructed.path})`,
              messageKey: "error.rowNavigationFailed",
              params: { path: constructed.path },
            });
            return;
          }
          // rowNavigate returns navigation.url (not request.url).
          const target = constructed.navigation?.url ?? constructed.request?.url;
          if (target === undefined) {
            setFeedback({ kind: "error", code: "INVALID_NAVIGATE_URL", message: "navigate mapping produced no url", messageKey: "error.navigateNoUrl" });
            return;
          }
          // GOAL-015 F-001: navigate through the host's session-internal
          // onNavigate (pushState + visit stack) so the entries page keeps its
          // breadcrumb trail; location.assign (full reload) is only the fallback
          // when no host hook is wired (e.g. embedded renderer tests).
          if (onNavigate !== undefined) {
            onNavigate(target);
          } else {
            window.location.assign(target);
          }
          return;
        }
        if (onNavigate !== undefined) {
          onNavigate(url);
        } else {
          window.location.assign(url);
        }
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
      }).catch((error: unknown) => {
        setFeedback(errorFeedback({ code: "ROW_ACTION_FAILED", message: String(error) }));
      });
    },
    [document, runRowAction],
  );

  // ADR-0022 D4: batch toolbar trigger — confirm first, then run the batch.
  const invokeBatchAction = useCallback((item: Record<string, unknown>, tableId: string) => {
      const actionRef = stringOf(item.actionRef);
      const action = actionOf(document, actionRef);
      if (actionRef === "" || action === undefined) {
        setFeedback({ kind: "error", code: "ACTION_NOT_FOUND", message: `action "${actionRef}" is not defined on this page`, messageKey: "error.actionNotFound", params: { action: actionRef } });
        return;
      }
      const current = selection(tableId);
      if (current === undefined || current.count === 0) {
        setFeedback({ kind: "error", code: "EMPTY_SELECTION", message: "select at least one row first", messageKey: "feedback.selectRowFirst" });
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
      // W4 P1-1: runBatchRequest resolves (never rejects) for construction and
      // network failures, but keep a catch so any unforeseen rejection still
      // surfaces as feedback instead of an unhandled promise rejection.
      void runBatchRequest(document, context, fetcher, actionRef, item, current)
        .then((result) => {
          if (result.ok) {
            setFeedback({ kind: "success", message: batchSuccessMessageFor(action, t) });
            reloadList();
          } else {
            setFeedback(errorFeedback(result));
          }
        })
        .catch((error: unknown) => {
          setFeedback(errorFeedback({ code: "BATCH_FAILED", message: String(error) }));
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
          throw new Error(t("error.uploadNotUploadAction", { action: field.actionRef }));
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
        throw new Error(t("error.uploadRequiresAction"));
      }
      const result = await uploadFilesWithFetch(action, files, fetcher);
      if (!result.ok) {
        throw new Error(`${result.code} (file ${result.fileIndex})`);
      }
      return result.fieldValue;
    },
    [document, fetcher, t],
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
      if (result.ok && (result.fieldErrors === undefined || result.fieldErrors.length === 0)) {
        // VP-029 判据 #2/#5（A-005 F-002 → A-008；声明载体修正见 E-007）：一次性
        // 明文 CSV 导出是**声明驱动**——只有提交表单节点的 props 声明 downloadCsv
        // 时，才在同一手势导出 items[].code；其它表单响应携带同名字段不再被启发
        // 式误触发下载。声明放在表单节点 props（node schema 的业务级参数区，upstream
        // pin 允许任意业务键）而不是 action.onSuccess：后者受 pinned action.schema
        // OutcomeBehavior 严格结构约束（additionalProperties: false），放在那里会
        // 使页面文档 D-VAL 失败（PAGE_SCHEMA_INVALID，用户可见「页面 Schema 错误」）。
        const formProps = isRecord(form.props) ? form.props : undefined;
        const downloadDecl = isRecord(formProps?.downloadCsv)
          ? (formProps.downloadCsv as Record<string, unknown>)
          : undefined;
        const columns = Array.isArray(downloadDecl?.columns)
          ? (downloadDecl.columns as unknown[]).filter(
              (column): column is string => typeof column === "string" && column !== "",
            )
          : [];
        const fileNameTemplate =
          typeof downloadDecl?.fileName === "string" && downloadDecl.fileName !== ""
            ? downloadDecl.fileName
            : "";
        const csvHeaders: Record<string, string> = {
          code: "Code",
          codePrefix: "Prefix",
          batchId: "BatchId",
          amount: "Amount",
          currency: "Currency",
          createdAt: "CreatedAt",
        };
        if (
          downloadDecl !== undefined &&
          columns.includes("code") &&
          fileNameTemplate !== "" &&
          isRecord(result.data) &&
          Array.isArray(result.data.items)
        ) {
          const vouchersWithCode = result.data.items.filter(
            (it): it is Record<string, unknown> => isRecord(it) && typeof it.code === "string",
          );
          if (vouchersWithCode.length > 0) {
            const lines = [columns.map((column) => csvHeaders[column] ?? column).join(",")];
            for (const v of vouchersWithCode) {
              const cells = columns.map((column) => {
                if (column === "amount") {
                  // E-008: wire amounts are CNY min units (分); the export and
                  // the table column (format: currency) both present yuan with
                  // two decimals (mirrors the W16-F04 display convention).
                  const cents = Number(v[column] ?? 0);
                  return Number.isFinite(cents) ? (cents / 100).toFixed(2) : "0.00";
                }
                return `"${String(v[column] ?? "")}"`;
              });
              lines.push(cells.join(","));
            }
            const csvBlob = new Blob([lines.join("\n")], { type: "text/csv;charset=utf-8" });
            const batchId =
              typeof prepared.batchId === "string" && prepared.batchId !== ""
                ? prepared.batchId
                : typeof vouchersWithCode[0]?.batchId === "string" && vouchersWithCode[0].batchId !== ""
                  ? (vouchersWithCode[0].batchId as string)
                  : "vouchers";
            triggerBlobDownload(csvBlob, sanitizeClientFilename(fileNameTemplate.replace("{batchId}", batchId)));
          }
        }
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
    // T-02 (GOAL-013 D-003): search forms may declare extra fields beyond the
    // keyword box (selects for discrete states). Every non-q field is bound
    // to the target table query filters (serialized by buildResourceQuery
    // into URL params) — this is the ExtraQuery pipe the backend whitelists.
    // Field ids owned by this form are first removed from the previous
    // filters so clearing a select actually clears the parameter (same
    // discipline as the q reset, C6). Table-level props.filters selects live
    // in the same query filters map and are preserved because they are not
    // owned by this form.
    const owned = new Set<string>();
    for (const rawField of form.props.fields) {
      if (typeof rawField.id === "string" && rawField.id !== "" && rawField.id !== "q") {
        owned.add(rawField.id);
      }
    }
    setQueries((prev) => {
      const current = prev[targetTable] ?? { page: 1, pageSize: 10 };
      const filters: Record<string, string> = { ...(current.filters ?? {}) };
      for (const id of owned) {
        delete filters[id];
      }
      for (const rawField of form.props.fields) {
        const id = rawField.id;
        if (typeof id !== "string" || id === "" || id === "q") {
          continue;
        }
        const value = values[id];
        const text =
          typeof value === "string"
            ? value.trim()
            : value === undefined || value === null
              ? ""
              : String(value);
        if (text !== "") {
          filters[id] = text;
        }
      }
      const next: ResourceQuery = { ...current, page: 1, q };
      if (Object.keys(filters).length > 0) {
        next.filters = filters;
      } else {
        // All owned conditions cleared: drop the inherited filters entirely
        // (the spread above would otherwise keep stale params — A-003 reset).
        delete next.filters;
      }
      return { ...prev, [targetTable]: next };
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
      fetchList,
      refreshList,
      listRefreshToken,
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
      route,
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
      fetchList,
      refreshList,
      listRefreshToken,
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
      route,
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
  const [dismissed, setDismissed] = useState(false);
  // W11 · U-03: toasts auto-dismiss; the parent remounts this component per
  // feedback occurrence, so the timer always starts fresh.
  useEffect(() => {
    setDismissed(false);
    if (feedback.kind === "error") {
      return;
    }
    const timer = window.setTimeout(() => setDismissed(true), FEEDBACK_TOAST_MS);
    return () => window.clearTimeout(timer);
  }, [feedback]);
  if (dismissed) {
    return null;
  }
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
      data-feedback-toast={feedback.kind}
      className={`fixed right-4 top-16 z-50 flex max-w-sm items-start justify-between gap-3 rounded-md border px-3 py-2 text-sm shadow-lg ${
        feedback.kind === "error"
          ? "border-destructive/50 bg-destructive/10 text-destructive"
          : "border-success/50 bg-success/10 text-success"
      }`}
    >
      <span data-feedback-code={feedback.code ?? ""} title={feedback.code ?? undefined}>
        {text}
      </span>
      <button
        type="button"
        aria-label={t("feedback.cancel")}
        onClick={() => setDismissed(true)}
        className="shrink-0 text-xs opacity-70 transition-opacity hover:opacity-100"
      >
        {"×"}
      </button>
    </div>
  );
}

/** W11 · U-03: auto-dismiss window for operation feedback toasts (ms). */
const FEEDBACK_TOAST_MS = 4000;

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
  context?: Record<string, unknown>,
): RecordSourcePrefillState {
  const recordSource = node.props.recordSource;
  // A-002 F-001: start a recordSource form in `loading` (not `idle`) so the
  // first commit renders the skeleton, never a blank editable form that could
  // overwrite the record while the prefill GET is pending.
  const [state, setState] = useState<RecordSourcePrefillState>(() =>
    recordSource !== undefined ? { status: "loading" } : { status: "idle" },
  );
  // W11 F-012: the route is a prefill input (query/params flow into
  // recordSource construction). The route OBJECT identity is unstable (App
  // rebuilds the render context each render), so the effect depends on this
  // serialized key: it re-runs only when query/params actually change, not
  // on unrelated parent renders.
  const routeKey = JSON.stringify(crud?.route ?? null);
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
    // Route bindings ($context.route.query/params) flow in via the render
    // context so deep links and query strings reach recordSource mapping (C8).
    const route = isRecord(context?.route)
      ? (context.route as { params?: Record<string, unknown>; query?: Record<string, unknown> })
      : undefined;
    // W4 P1-1: constructRequest throws on malformed recordSource values (e.g.
    // serializeQueryValue rejects non-scalars). This effect runs in a React
    // useEffect — an uncaught throw would unmount the whole tree (no
    // ErrorBoundary). Trap it into the error state instead.
    let constructed: RequestConstructionResult;
    try {
      constructed = constructRequest({
        kind: "recordSource",
        recordSource,
        route: route ?? { params: {}, query: {} },
        baseURL: "",
      });
    } catch (error) {
      setState({
        status: "error",
        message: `recordSource construction failed: ${error instanceof Error ? error.message : "unknown error"}`,
      });
      return;
    }
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
  // W11 F-012: the route is a prefill input. The route OBJECT identity is
  // unstable (App rebuilds the render context each render), so depend on a
  // serialized key — the effect re-runs only when query/params actually
  // change, not on unrelated parent renders. The previous deps omitted the
  // route entirely: a same-page query change left the PREVIOUS record's
  // values pre-filled and a save could write to the wrong row.
  }, [recordSource, node.props.mode, metaValue, crud?.fetcher, crud?.reloadToken, routeKey]);
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
  const prefill = useRecordSourcePrefill(node, metaValue, crud, context);

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
    // GOAL-015 / ADR-0040: a create-modal (no row, no recordSource) may have
    // its readOnly fields seeded by the Host from the current route — the
    // v2.9 route-filter scenario ("值由 Host 在 modal 场景提供") wires the
    // inner-page type key into the new-entry form this way. The provider's
    // route snapshot covers both path params (dictKey) and the query
    // (dictTypeName); hostless renders fall back to the location query.
    const routeValues: Record<string, string> = {};
    if (crud !== null) {
      for (const [key, value] of Object.entries(crud.route.params)) {
        routeValues[key] = value;
      }
      for (const [key, value] of Object.entries(crud.route.query)) {
        routeValues[key] = value;
      }
    } else if (typeof window !== "undefined") {
      for (const [key, value] of new URLSearchParams(window.location.search).entries()) {
        routeValues[key] = value;
      }
    }
    for (const raw of node.props.fields) {
      if (!isRecord(raw) || typeof raw.id !== "string") {
        continue;
      }
      // Precedence: modal row (edit-in-modal) → recordSource prefill → route
      // query (create-modal readOnly fields) → empty. Always coerce through
      // the field wire kind so row arrays (e.g. roles[]) become textarea
      // strings or checkboxGroup string[] as appropriate (A-006 R-004).
      const modalValue =
        modalRow !== null && modalRow[raw.id] !== undefined ? modalRow[raw.id] : undefined;
      const prefillValue =
        prefillValues !== null && prefillValues[raw.id] !== undefined
          ? prefillValues[raw.id]
          : undefined;
      const routeValue =
        raw.readOnly === true && routeValues[raw.id] !== undefined
          ? routeValues[raw.id]
          : undefined;
      const fromRow =
        modalValue !== undefined ? modalValue : prefillValue !== undefined ? prefillValue : routeValue;
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
  // GOAL-014 D-002 §3: submit-time validation + server fieldErrors echo,
  // keyed by field id for inline display.
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});
  // W16-F03: row-level import failure detail (200 response fieldErrors).
  const [importErrorRows, setImportErrorRows] = useState<
    Array<{ rowNumber?: number; field: string; reason: string }>
  >([]);

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
    // GOAL-014 D-002 §3.1: submit-time client validation — failures block the
    // request and inline on their fields.
    const validation = validateFieldValues(visibleFields, values);
    if (validation.length > 0) {
      const byField: Record<string, string> = {};
      for (const err of validation) {
        const localized = t(err.messageKey);
        byField[err.field] = localized === err.messageKey ? err.message : localized;
      }
      setFieldErrors(byField);
      setFormError(null);
      return;
    }
    // W16-F04: retain a second confirmation for risky wallet adjustments.
    if (node.id === "adjust-wallet-form") {
      const delta = values.amountDelta;
      const risky =
        typeof delta === "number" && (delta < 0 || Math.abs(delta) > 100000);
      if (risky && !window.confirm(t("schema.wallet.adjustConfirm"))) {
        return;
      }
    }
    setFieldErrors({});
    setImportErrorRows([]);
    setSubmitting(true);
    setFormError(null);
    try {
      const result = await crud.submitForm(node, values);
      if (!result.ok) {
        // GOAL-014 D-002 §2: echo server fieldErrors onto the matching inputs;
        // unmatched fields fall back to the form-level alert.
        const byField: Record<string, string> = {};
        for (const fe of result.fieldErrors ?? []) {
          byField[fe.field] = fe.reason;
        }
        setFieldErrors(byField);
        setFormError({
          code: result.code,
          message: result.message,
          ...(result.messageKey === undefined ? {} : { messageKey: result.messageKey }),
          ...(result.params === undefined ? {} : { params: result.params }),
        });
      } else if (result.ok && (result.fieldErrors ?? []).length > 0) {
        // W16-F03: a 200 import response with fieldErrors is a partial
        // failure — keep the modal open and show row-level errors.
        const rows = (result.fieldErrors ?? []).map((fe) => ({
          rowNumber: fe.rowNumber,
          field: fe.field,
          reason: fe.reason,
        }));
        setImportErrorRows(rows);
        setFormError({
          code: "IMPORT_HAS_ERRORS",
          message: t("importErrors.title"),
        });
        return;
      }
    } catch (error) {
      // Defensive: a throwing submit (unexpected fetch/transport failure) must
      // never leave the button stuck in its disabled Submitting state (C5).
      setFieldErrors({});
      setFormError({ code: "REQUEST_FAILED", message: requestFailedMessage(error) });
    } finally {
      setSubmitting(false);
    }
  };

  const title = resolveTextProp(
    node.props as unknown as Record<string, unknown>,
    "titleKey",
    "title",
    t,
    "",
  );

  // A-003 (audit response): search-mode toolbar — active filter chips with
  // per-condition clear + one-click clear-all; the reset button restores the
  // default values and re-runs the search immediately.
  //
  // T-07 (GOAL-014): chips reflect the SUBMITTED query (the target table's
  // query state), not the local draft values — a keyword that was typed but
  // not submitted (or a select that was changed but not applied) never shows
  // as an active filter. This keeps the chips truthful: they describe what
  // the list is actually filtered by.
  const targetTable =
    isSearch && typeof node.props.targetTable === "string" ? node.props.targetTable : "";
  const submittedQuery =
    isSearch && crud !== null && targetTable !== "" ? crud.tableQuery(targetTable) : undefined;
  const activeFilters = isSearch
    ? visibleFields
        .map((field) => {
          const raw = field.id === "q" ? submittedQuery?.q : submittedQuery?.filters?.[field.id];
          if (raw === undefined || raw === null || raw === "") {
            return null;
          }
          const text = typeof raw === "string" ? raw.trim() : String(raw);
          if (text === "") {
            return null;
          }
          const valueLabel =
            field.type === "select"
              ? (optionList(field, t).find((option) => option.value === text)?.label ?? text)
              : text;
          return { field, value: text, valueLabel };
        })
        .filter((entry): entry is { field: FormControlField; value: string; valueLabel: string } => entry !== null)
    : [];
  // T-07 (GOAL-014): non-text filter controls (selects, date pickers, …)
  // apply IMMEDIATELY on change — the list re-filters as the condition
  // changes. The keyword box keeps the submit-on-search contract (its value
  // only joins the query when the paired search button is pressed).
  const handleFieldChange = (id: string, value: unknown) => {
    const next = { ...values, [id]: value };
    setValues(next);
    if (!isSearch || crud === null) {
      return;
    }
    const field = visibleFields.find((entry) => entry.id === id);
    if (field !== undefined && field.type !== "input") {
      crud.searchFormSubmit(node, next);
    }
  };

  const resetValues = () => {
    const cleared: Record<string, unknown> = {};
    for (const field of visibleFields) {
      cleared[field.id] =
        field.defaultValue !== undefined ? field.defaultValue : "";
    }
    setValues(cleared);
    if (crud !== null) {
      crud.searchFormSubmit(node, cleared);
    }
  };
  const removeFilter = (field: FormControlField) => {
    const next = { ...values, [field.id]: "" };
    setValues(next);
    if (crud !== null) {
      crud.searchFormSubmit(node, next);
    }
  };

  return (
    <form
      className={isSearch ? "space-y-3.5 rounded-xl border border-border/70 bg-card/85 p-4 shadow-[0_1px_3px_0_rgba(0,0,0,0.03),0_1px_2px_-1px_rgba(0,0,0,0.03)] dark:border-border/60 dark:bg-card/70 dark:shadow-[0_1px_3px_0_rgba(0,0,0,0.2)]" : "space-y-3"}
      onSubmit={(event) => {
        event.preventDefault();
        void handleSubmit();
      }}
    >
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
        onChange={handleFieldChange}
        fieldDisabled={fieldDisabled}
        onUpload={crud?.uploadFiles}
        fieldErrors={fieldErrors}
        // W11 · U-01/U-02: auth-aware transport for dynamic option sources.
        fetcher={crud?.fetcher}
        columns={
          isRecord(node.props) && typeof node.props.columns === "number"
            ? node.props.columns
            : undefined
        }
        searchMode={isSearch}
        // A-003 pairing rule (user 2026-08-16): the search button belongs
        // beside its keyword input — one button per text input, adjacent.
        searchButtonSlot={
          isSearch ? (
            <button
              type="submit"
              disabled={hasBlockingErrors}
              className="inline-flex h-9 shrink-0 cursor-pointer items-center justify-center gap-1.5 rounded-l-none rounded-r-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-2xs transition-all duration-150 hover:bg-primary/90 hover:shadow-xs active:scale-[0.98] disabled:pointer-events-none disabled:opacity-50 -ml-px"
            >
              <Search aria-hidden="true" className="size-3.5 stroke-[2.2]" />
              {resolveTextProp(
                node.props as unknown as Record<string, unknown>,
                "submitLabelKey",
                "submitLabel",
                t,
                t("feedback.search"),
              )}
            </button>
          ) : undefined
        }
        actionSlot={
          isSearch ? (
            <div className="flex items-end">
              <button
                type="button"
                onClick={resetValues}
                className="inline-flex h-9 shrink-0 cursor-pointer items-center justify-center gap-1.5 rounded-md border border-input/80 bg-background px-3.5 text-sm font-medium text-muted-foreground shadow-2xs transition-all duration-150 hover:border-muted-foreground/30 hover:bg-accent/40 hover:text-foreground hover:shadow-xs active:scale-[0.98]"
              >
                <RotateCcw aria-hidden="true" className="size-3.5" />
                {t("feedback.reset")}
              </button>
            </div>
          ) : undefined
        }
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
      {importErrorRows.length > 0 ? (
        <ul role="alert" className="max-h-40 space-y-1 overflow-y-auto rounded-md border border-destructive/30 bg-destructive/5 p-2 text-xs" data-import-error-rows>
          {importErrorRows.map((row, index) => (
            <li key={index} className="flex gap-2">
              <span className="shrink-0 font-mono text-muted-foreground">
                {row.rowNumber === undefined ? "—" : `#${row.rowNumber}`}
              </span>
              <span className="shrink-0 font-medium">{row.field}</span>
              <span className="text-destructive">{row.reason}</span>
            </li>
          ))}
        </ul>
      ) : null}
      {activeFilters.length > 0 ? (
        <div className="flex flex-wrap items-center gap-2 border-t border-border/50 pt-2.5 text-xs" data-filter-chips>
          {activeFilters.map(({ field, valueLabel }) => {
            const fieldLabel = resolveTextProp(
              field as unknown as Record<string, unknown>,
              "labelKey",
              "label",
              t,
              field.id,
            );
            return (
              <span
                key={field.id}
                className="inline-flex h-6.5 items-center gap-1.5 rounded-md border border-border/70 bg-muted/40 px-2.5 text-xs font-normal text-foreground/90 shadow-2xs transition-all hover:border-border dark:bg-muted/30"
              >
                <span className="text-muted-foreground/80">{fieldLabel}:</span>
                <span className="font-medium text-foreground">{valueLabel}</span>
                <button
                  type="button"
                  aria-label={t("feedback.removeFilter") + ": " + fieldLabel}
                  className="ml-0.5 inline-flex size-3.5 cursor-pointer items-center justify-center rounded-full text-muted-foreground/60 transition-colors hover:bg-muted-foreground/20 hover:text-foreground"
                  onClick={() => removeFilter(field)}
                >
                  <X aria-hidden="true" className="size-3 stroke-[2.5]" />
                </button>
              </span>
            );
          })}
          <button
            type="button"
            className="ml-1 cursor-pointer text-xs font-medium text-muted-foreground/70 underline-offset-4 transition-colors hover:text-destructive hover:underline"
            onClick={resetValues}
          >
            {t("feedback.clearFilters")}
          </button>
        </div>
      ) : null}
      {canSubmit && !isSearch ? (
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
  const baseId = typeof node.id === "string" && node.id !== "" ? node.id : "tabs";
  const focusTab = (from: number, direction: -1 | 0 | 1) => {
    if (children.length === 0) {
      return;
    }
    const next = (from + direction + children.length) % children.length;
    setActive(next);
    const button = document.getElementById(`tab-${baseId}-${next}`);
    button?.focus();
  };
  return (
    <div className="space-y-3">
      {children.length > 1 ? (
        <div role="tablist" aria-label={t("feedback.tabs")} className="flex flex-wrap items-center gap-1 border-b border-border">
          {children.map((child, index) => {
            const rawProps = (child as unknown as { props?: unknown }).props;
            // T-03 (GOAL-013 D-004): resolve labelKey before the literal
            // label (registry convention), so tabs stay bilingual.
            const label = resolveTextProp(
              isRecord(rawProps) ? rawProps : {},
              "labelKey",
              "label",
              t,
              `Tab ${index + 1}`,
            );
            const panelId = `tab-panel-${baseId}-${index}`;
            return (
              <button
                key={index}
                id={`tab-${baseId}-${index}`}
                type="button"
                role="tab"
                aria-selected={index === active}
                aria-controls={panelId}
                tabIndex={index === active ? 0 : -1}
                onClick={() => setActive(index)}
                onKeyDown={(event) => {
                  if (event.key === "ArrowRight") {
                    event.preventDefault();
                    focusTab(index, 1);
                  } else if (event.key === "ArrowLeft") {
                    event.preventDefault();
                    focusTab(index, -1);
                  } else if (event.key === "Home") {
                    event.preventDefault();
                    focusTab(0, 0);
                  } else if (event.key === "End") {
                    event.preventDefault();
                    focusTab(children.length - 1, 0);
                  }
                }}
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
        <div
          id={`tab-panel-${baseId}-${active}`}
          role="tabpanel"
          aria-labelledby={`tab-${baseId}-${active}`}
          tabIndex={0}
          className="outline-none"
        >
          {dispatchNode({
            node: current,
            path: "body.children",
            metaValue,
            context,
            tableRenderer,
            onAction,
            formComponent,
          })}
        </div>
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
function formatRecordViewValue(value: unknown): string {
  if (Array.isArray(value)) {
    return value.join(", ");
  }
  if (value === undefined || value === null) {
    return "—";
  }
  // ISO timestamps render as local "YYYY-MM-DD HH:mm" (display formatting).
  return formatDisplayTime(value) ?? String(value);
}

function RecordView({ node }: { node: RenderRecordViewNode }) {
  const crud = useSchemaCrud();
  const t = useTranslate();
  const staticRecord = node.props?.record;
  const hasStatic = isRecord(staticRecord);
  const record = hasStatic ? staticRecord : (crud?.selectedRow ?? null);
  const declaredFields = node.props?.fields ?? [];
  const rows =
    record === null
      ? []
      : declaredFields.length > 0
        ? declaredFields.map((field) => ({
            key: field.key,
            label: resolveTextProp(
              field as unknown as Record<string, unknown>,
              "labelKey",
              "label",
              t,
              field.key,
            ),
            value: record[field.key],
          }))
        : Object.entries(record).map(([key, value]) => ({ key, label: key, value }));
  if (rows.length === 0) {
    return (
      <p
        data-record-view="empty"
        className="text-sm text-muted-foreground"
      >
        {t("feedback.selectRecordToView")}
      </p>
    );
  }

  const title = resolveTextProp(
    node.props as unknown as Record<string, unknown>,
    "titleKey",
    "title",
    t,
    t("feedback.recordDetails"),
  );

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
        aria-label={title}
        className={
          canClose
            ? // Desktop (md+): right Drawer; mobile (<768, D-004): full-height Sheet
              "fixed inset-y-0 right-0 z-50 flex w-full max-w-md flex-col border-l border-border bg-card shadow-lg max-md:inset-x-0 max-md:top-auto max-md:h-[min(92vh,100%)] max-md:rounded-t-xl max-md:border-l-0 max-md:border-t md:max-w-md"
            : "w-full max-w-md rounded-lg border border-border bg-card shadow-sm"
        }
      >
        <div className="flex items-center justify-between gap-3 border-b border-border px-4 py-3">
          <h2 className="text-sm font-semibold tracking-tight text-foreground">
            {title}
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
          {rows.map((row) => (
            <div key={row.key} className="grid gap-0.5 text-sm sm:grid-cols-[8rem_minmax(0,1fr)] sm:gap-3">
              <dt className="text-xs font-medium text-muted-foreground">
                {row.label}
              </dt>
              <dd className="break-words text-foreground">
                {formatRecordViewValue(row.value)}
              </dd>
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
  params?: Record<string, unknown>,
): { list: ResourceList | null; error: string | null } {
  const crud = useSchemaCrud();
  useEffect(() => {
    if (crud !== null && fetcher !== undefined) {
      crud.registerFetcher(fetcher);
    }
  }, [crud, fetcher]);
  const [list, setList] = useState<ResourceList | null>(null);
  const [error, setError] = useState<string | null>(null);
  // v2.9 ADR-0039: route snapshot for dataSource params bindings. Prefers the
  // provider's route context; hostless renders fall back to the location query.
  const routeSnapshot = useMemo(() => {
    if (crud !== null) {
      return crud.route;
    }
    const query: Record<string, string> = {};
    if (typeof window !== "undefined") {
      for (const [key, value] of new URLSearchParams(window.location.search).entries()) {
        query[key] = value;
      }
    }
    return { query, params: {} };
  }, [crud]);
  // W25: targeted refresh token for this dataSource (refreshList): a change
  // re-runs the fetch effect for THIS URL only, without a full page reload.
  const targetedRefreshToken =
    crud === null || dataSource === null ? 0 : crud.listRefreshToken(dataSource);
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
    const paramsQuery = resolveDataParamsQuery(params, routeSnapshot);
    const request =
      crud !== null
        ? crud.fetchList(dataSource, DISPLAY_LIST_QUERY, paramsQuery, fetcher)
        : fetchResourceList(fetcher ?? fetch, dataSource, DISPLAY_LIST_QUERY, paramsQuery);
    request
      .then((next) => {
        if (!cancelled) {
          setList(next);
          setError(null);
        }
      })
      .catch((err: unknown) => {
        if (!cancelled) {
          if (isWalletNotFoundError(err)) {
            setList(EMPTY_RESOURCE_LIST);
            setError(null);
            return;
          }
          setError(err instanceof Error ? err.message : String(err));
        }
      });
    return () => {
      cancelled = true;
    };
  }, [fetcher, dataSource, params, routeSnapshot, crud?.reloadToken, targetedRefreshToken]);
  return { list, error };
}

function StatCardView({ node }: { node: RenderStatCardNode }) {
  const crud = useSchemaCrud();
  const t = useTranslate();
  const fetcher = crud?.fetcher ?? globalThis.fetch;
  const valueField = node.props?.valueField;
  const format = node.props?.format ?? "plain";
  // v2.9 ADR-0039: node-level DataRef preferred over legacy props.dataSource.
  const dataSource =
    typeof node.data?.url === "string" && isValidDataSource(node.data.url)
      ? node.data.url
      : typeof node.props?.dataSource === "string" && isValidDataSource(node.props.dataSource)
        ? node.props.dataSource
        : null;
  const { list, error } = useDisplayData(dataSource, fetcher, node.data?.params);

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
  // v2.9 ADR-0039: node-level DataRef preferred over legacy props.dataSource.
  const dataSource =
    typeof node.data?.url === "string" && isValidDataSource(node.data.url)
      ? node.data.url
      : typeof node.props?.dataSource === "string" && isValidDataSource(node.props.dataSource)
        ? node.props.dataSource
        : null;
  const { list, error } = useDisplayData(dataSource, fetcher, node.data?.params);

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

/**
 * F-008 (grok audit): gates a data node that uses v2.9 ADR-0039 route-bound
 * params ($context.route.* bindings in node.data.params) on the page meta —
 * protocolVersion >= 2.9 AND the data.route-binding capability must be
 * declared, else the node fails closed with an observable error (same
 * discipline as checkFormCapabilities for form fields).
 */
function gateDataRouteBinding(
  metaValue: unknown,
  node: RenderTableNode | RenderStatCardNode | RenderChartNode,
): string | null {
  const params = node.data?.params;
  if (params === undefined) {
    return null;
  }
  const usesRouteBinding = Object.values(params).some(
    (value) => typeof value === "string" && value.startsWith("$context.route."),
  );
  if (!usesRouteBinding) {
    return null;
  }
  const meta = isRecord(metaValue) ? metaValue : {};
  const version = typeof meta.protocolVersion === "string" ? meta.protocolVersion : "";
  const versionMatch = /^(\d+)\.(\d+)$/.exec(version);
  const versionOk =
    versionMatch !== null &&
    (Number(versionMatch[1]) > 2 ||
      (Number(versionMatch[1]) === 2 && Number(versionMatch[2]) >= 9));
  const capabilities = Array.isArray(meta.requiredCapabilities)
    ? (meta.requiredCapabilities as unknown[]).filter(
        (c): c is string => typeof c === "string",
      )
    : [];
  if (!versionOk || !capabilities.includes("data.route-binding")) {
    return "dataSource params use $context.route.* bindings which require protocol >= 2.9 and the data.route-binding capability";
  }
  return null;
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
    case "custom": {
      // GOAL-018: custom nodes dispatch to the module-level registry; an
      // unregistered component renders a safe fallback (never crashes).
      const Custom = getCustomComponent(node.component);
      if (Custom === null) {
        return (
          <p className="text-sm text-muted-foreground">
            unknown custom component: {node.component}
          </p>
        );
      }
      return <Custom node={node} context={context} children={node.children} />;
    }
    case "recordView":
      return <RecordView node={node} />;
    case "actionButton":
      return <ActionButtonView node={node} context={context} onAction={onAction} />;
    case "statCard":
    case "chart":
    case "table": {
      const routeGate = gateDataRouteBinding(metaValue, node);
      if (routeGate !== null) {
        return (
          <p key={node.id ?? node.type} role="alert" className="text-sm text-destructive">
            {routeGate}
          </p>
        );
      }
      if (node.type === "statCard") {
        return <StatCardView node={node} />;
      }
      if (node.type === "chart") {
        return <ChartView node={node} />;
      }
      return (
        tableRenderer?.(node) ?? (
          <p className="text-sm text-muted-foreground">
            table node rendered without a tableRenderer (the app wires SchemaTable)
          </p>
        )
      );
    }
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
            // The triggering row rides the modal context so custom content can
            // act on it (workspace-019 resend dialog consumes modalRow).
            context: { ...context, modalRow: crud.activeModal?.row ?? null },
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
  onNavigate,
  formComponent,
}: RendererComponentProps & { dataFetcher?: typeof fetch }) {
  return (
    <SchemaCrudProvider
      document={document}
      context={context}
      initialFetcher={dataFetcher}
      onNavigate={onNavigate}
    >
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
