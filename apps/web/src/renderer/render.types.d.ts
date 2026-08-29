import { type FormControlField, type FormControlGateError } from "@/renderer/form-controls.types";
import { type FormControlStateMap, type ReactionError } from "@/renderer/reactions";
/**
 * D-COMP page renderer (frozen §5 whitelist + I-PROTO-FULL-001 full registry
 * surface; resolve R4 F-002).
 *
 * The renderer walks a page document's node tree and dispatches whitelisted
 * node types to the components in this directory. It is deliberately minimal:
 * only the node types surfaced by the example pages are handled, and any
 * other type fails closed instead of rendering a silent fallback.
 *
 * A page document looks like the example pages' `PAGE_DOCUMENT`:
 *   { meta: { protocolVersion, requiredCapabilities }, body: Node }
 *
 * Node types supported (registry surface):
 *   - layout:  grid / section / tabs
 *   - data/action: text / table / recordView / actionButton / statCard / chart
 *   - form:    form → FormControls (with reactions applied to field state)
 * The form control whitelist itself is enforced by D-FORM
 * (isWhitelistedFormControl / checkFormCapabilities).
 */
export type RenderNodeType = "form" | "section" | "table" | "grid" | "tabs" | "text" | "recordView" | "actionButton" | "statCard" | "chart" | "custom";
export interface RenderMeta {
    protocolVersion: string;
    requiredCapabilities: string[];
}
export interface ReactionState {
    /** fieldId → control state after reactions. */
    state: FormControlStateMap;
    errors: ReactionError[];
}
export interface RenderFormNode {
    type: "form";
    id?: string;
    props: {
        fields: Array<Record<string, unknown>>;
        reactions?: unknown;
        submitLabel?: string;
        /** S2 (VP-007): i18n key resolved before `submitLabel` (local doc convention). */
        submitLabelKey?: string;
        /** Default-mode submit: the top-level action id to run on submit (S4). */
        submitAction?: string;
        /** Search-mode form: binds its fields to the target table's query (S4). */
        mode?: "default" | "search";
        targetTable?: string;
        /** Section heading rendered above the fields (registry form `title`/`titleKey`). */
        title?: string;
        /** i18n key resolved before `title`. */
        titleKey?: string;
        /**
         * ADR-0021 edit-form record GET prefill (registry since 2.1): loads current
         * values from a detail GET and initializes the fields via `responseMapping`
         * (field id → dot-path in the response). Requires capability
         * `form.record.load`; search-mode forms forbid it.
         */
        recordSource?: Record<string, unknown>;
        /** GOAL-014 D-002 §4: responsive form column count (>1 enables the grid). */
        columns?: number;
    };
    children?: RenderNode[];
}
export interface RenderSectionNode {
    type: "section";
    id?: string;
    props?: Record<string, unknown>;
    children: RenderNode[];
}
export interface RenderGridNode {
    type: "grid";
    id?: string;
    props?: Record<string, unknown>;
    children?: RenderNode[];
}
export interface RenderTabsNode {
    type: "tabs";
    id?: string;
    props?: Record<string, unknown>;
    children?: RenderNode[];
}
export interface RenderTextNode {
    type: "text";
    id?: string;
    props?: {
        text?: string;
        /** i18n key resolved before `text` (F-01 · GOAL-003 A-003 F-001). */
        textKey?: string;
    };
    children?: RenderNode[];
}
/**
 * v2.9 DataRef subset consumed by data nodes (ADR-0039): source:api url +
 * params (literal scalars or whole $context.route.query.* / params.* bindings).
 */
export interface RenderDataRef {
    url: string;
    params?: Record<string, unknown>;
}
export interface RenderTableNode {
    type: "table";
    id?: string;
    /** v2.9 node-level DataRef (source:api). Preferred over props.dataSource. */
    data?: RenderDataRef;
    props: {
        columns?: Array<Record<string, unknown>>;
        actions?: Array<Record<string, unknown>>;
        toolbar?: Array<Record<string, unknown>>;
        dataSource?: string;
        /** Table heading (literal fallback for titleKey). */
        title?: string;
        /** i18n key resolved before title. */
        titleKey?: string;
        /** Direct field name of each row's unique key (F-002 · I-010-001 v0.2.0 §3; default "id"). */
        rowKey?: string;
        /** ADR-0022 multi-select model (registry; mode: multiple only). */
        selection?: {
            mode?: string;
        };
        /** Schema-driven filters (select-only; see schemaTableFilters). */
        filters?: Array<Record<string, unknown>>;
    };
    children?: RenderNode[];
}
/** Read-only field row on a recordView (registry `fields[]` since 2.4). */
export interface RenderRecordViewField {
    key: string;
    label?: string;
    labelKey?: string;
}
export interface RenderRecordViewNode {
    type: "recordView";
    id?: string;
    props?: {
        record?: Record<string, unknown>;
        /** Detail panel heading (registry `title` since 2.4). */
        title?: string;
        /** i18n key resolved before `title`. */
        titleKey?: string;
        /** Declared display fields; when set, only these rows render. */
        fields?: RenderRecordViewField[];
    };
    children?: RenderNode[];
}
export interface RenderActionButtonNode {
    type: "actionButton";
    id?: string;
    props?: {
        label?: string;
        /** i18n key resolved before `label` (S3). */
        labelKey?: string;
        actionId?: string;
        visibleWhen?: unknown;
        disabledWhen?: unknown;
        /** Permission-intent key (ADR-0023 D4b mount); gates the button target. */
        permissionIntent?: string;
        /** Target id for the permission target / action gate (falls back to node id). */
        key?: string;
        /** Confirm message (shown before executing the referenced action). */
        confirm?: string;
        /** i18n key resolved before `confirm`. */
        confirmKey?: string;
    };
    children?: RenderNode[];
}
export interface RenderStatCardNode {
    type: "statCard";
    id?: string;
    /** v2.9 node-level DataRef (source:api). Preferred over props.dataSource. */
    data?: RenderDataRef;
    props?: {
        label?: string;
        /** i18n key resolved before `label` (F-01 · GOAL-003 A-003 F-001). */
        labelKey?: string;
        unit?: string;
        /** plain | currency | percent (registry enum). */
        format?: string;
        /** Field of the dataSource rows to display (registry, required since 0.2). */
        valueField?: string;
        /** Single-slash same-origin data path (same invariant as table.dataSource). */
        dataSource?: string;
    };
    children?: RenderNode[];
}
export interface RenderChartNode {
    type: "chart";
    id?: string;
    /** v2.9 node-level DataRef (source:api). Preferred over props.dataSource. */
    data?: RenderDataRef;
    props?: {
        /** line | bar | pie (registry enum, required). */
        chartType?: string;
        xField?: string;
        yField?: string;
        /** Single-slash same-origin data path (same invariant as table.dataSource). */
        dataSource?: string;
    };
    children?: RenderNode[];
}
export interface RenderCustomNode {
    type: "custom";
    id?: string;
    /** Registered component key (GOAL-018: renderer custom-component registry). */
    component: string;
    props?: Record<string, unknown>;
    children?: RenderNode[];
}
export type RenderNode = RenderFormNode | RenderSectionNode | RenderGridNode | RenderTabsNode | RenderTextNode | RenderTableNode | RenderRecordViewNode | RenderActionButtonNode | RenderStatCardNode | RenderChartNode | RenderCustomNode;
/** Page-level action table entry (registry: modal | request | navigate). */
export interface RenderPageAction {
    type: string;
    /** Modal content; required when type=modal. */
    content?: RenderNode;
    /** Request/navigate target URL. */
    url?: string;
    /** HTTP method for request-shaped actions. */
    method?: string;
    /** Permission intent name, gated against page-level permissionCascade. */
    permissionIntent?: string;
}
export interface RenderPageDocument {
    meta: RenderMeta;
    body: RenderNode;
    /** Page-level action table referenced by actionRef / actionId. */
    actions?: Record<string, RenderPageAction>;
}
export type RenderErrorCode = "RENDER_UNKNOWN_NODE_TYPE" | "RENDER_INVALID_BODY" | "RENDER_META_INVALID" | "RENDER_FORM_FIELD_INVALID";
export interface RenderError {
    code: RenderErrorCode;
    path: string;
    message: string;
}
export type ActionGateErrorCode = "ACTION_GATE_EXPRESSION_INVALID";
export interface ActionGateError {
    code: ActionGateErrorCode;
    path: string;
    message: string;
}
export type ActionGateResult = {
    kind: "ok";
    value: boolean;
} | {
    kind: "error";
    error: ActionGateError;
};
export interface RenderFormFieldGate {
    /** Fields that passed the type whitelist and the version/capability gate. */
    fields: FormControlField[];
    /** Deterministic errors for rejected fields and gate failures. */
    errors: FormControlGateError[];
}
/** Keeps only well-formed recordView field rows (key required). */
export declare function parseRecordViewFields(raw: unknown): RenderRecordViewField[];
export declare function isWhitelistedNodeType(type: string): type is RenderNodeType;
/**
 * Resolves a dot-path on a record for `form.recordSource.responseMapping`
 * (e.g. `"customer.name"` → record.customer.name). Missing segments or a
 * non-object intermediate return `undefined` (field keeps its default/empty).
 */
export declare function resolveResponsePath(record: unknown, path: string): unknown;
/** Normalizes an unknown body value into a typed RenderNode, fail-closed. */
export declare function parseRenderNode(value: unknown, path: string): RenderNode | RenderError;
/** Collects all form fieldIds from a page document (for reaction resolution). */
export declare function collectFieldIds(node: RenderNode, into?: string[]): string[];
/** Resolves the reaction state for a form node's fields. */
export declare function resolveFormReactions(form: RenderFormNode, context: Record<string, unknown>): ReactionState;
/**
 * Resolves a table action gate expression.
 *
 * Distinguishes an absent property (→ `absentDefault`) from an explicit but
 * invalid expression (→ fail-closed `error`, no silent default). Valid $context
 * expressions are evaluated against the frozen engine; booleans pass through.
 */
export declare function resolveActionGate(expression: unknown, context: Record<string, unknown>, absentDefault: boolean, path: string): ActionGateResult;
/**
 * Evaluates a whitelisted table action's visibility/disabled gate.
 * Explicit invalid expressions fail closed (hidden / disabled) and are reported.
 */
export declare function tableActionGate(action: Record<string, unknown>, context: Record<string, unknown>): {
    visible: boolean;
    disabled: boolean;
    errors: ActionGateError[];
};
/**
 * Parses and gates a form node's raw fields against the D-FORM §5 whitelist
 * and the page meta version/capability rules (F-002 / F-003).
 *
 * Returns only the fields that pass every gate; rejected fields produce
 * deterministic errors instead of being silently rendered (or silently dropped).
 */
export declare function gateRenderFormFields(metaValue: unknown, rawFields: unknown, path: string): RenderFormFieldGate;
