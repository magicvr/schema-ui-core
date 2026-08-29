import { type ComponentType, type ReactNode } from "react";
import { type UploadableFile } from "@/protocol/conformance/upload-orchestration";
import { type FormControlField } from "@/renderer/form-controls.types";
import { type ResourceList, type ResourceQuery } from "@/renderer/resource";
import { type RenderActionButtonNode, type RenderChartNode, type RenderFormNode, type RenderPageDocument, type RenderStatCardNode, type RenderTableNode } from "@/renderer/render.types";
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
    batch?: {
        tableId: string;
        selection: TableSelection;
    };
    /** Toolbar item batchMapping carried through confirm. */
    batchMapping?: Record<string, unknown>;
}
/** Selection snapshot (ADR-0022 D3): ordered keys + count. */
export interface TableSelection {
    keys: unknown[];
    count: number;
}
export type ActionResult = {
    ok: true;
    fieldErrors?: Array<{
        field: string;
        reason: string;
        rowNumber?: number;
    }>;
    message?: string;
    messageKey?: string;
} | {
    ok: false;
    code: string;
    message: string;
    messageKey?: string;
    params?: Record<string, unknown>;
    /** GOAL-014 D-002 §2: server field-level validation failures. */
    fieldErrors?: Array<{
        field: string;
        reason: string;
        rowNumber?: number;
    }>;
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
    fetchList: (dataSource: string, query: ResourceQuery, extraQuery?: string, transport?: typeof fetch) => Promise<ResourceList>;
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
    activeModal: {
        actionRef: string;
        row: Record<string, unknown> | null;
        title: string;
    } | null;
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
    runRowAction: (actionRef: string, opts: RunRequestOptions) => Promise<ActionResult>;
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
    route: {
        query: Record<string, string>;
        params: Record<string, string>;
    };
}
export interface RunRequestOptions {
    row?: Record<string, unknown> | null;
    formValues?: Record<string, unknown>;
    requestMapping?: Record<string, unknown>;
    gateTargetId?: string;
    confirmed?: boolean;
}
/** Page CRUD context consumed by custom components (rendered nodes). */
export declare const SchemaCrudContext: import("react").Context<SchemaCrudValue | null>;
/** Reads the page-level Schema CRUD provider (null when rendered bare). */
export declare function useSchemaCrud(): SchemaCrudValue | null;
export declare function RenderPage({ document, context, tableRenderer, dataFetcher, onAction, onNavigate, formComponent, }: RendererComponentProps & {
    dataFetcher?: typeof fetch;
}): import("react").JSX.Element;
