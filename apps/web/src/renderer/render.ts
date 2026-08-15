import { evaluateExpression, isValidExpression } from "@/protocol/app-manifest";

import {
  checkFormCapabilitiesRaw,
  isWhitelistedFormControl,
  type FormControlField,
  type FormControlGateError,
} from "@/renderer/form-controls";
import {
  parseAndEvaluateReactions,
  type FormControlStateMap,
  type ReactionError,
} from "@/renderer/reactions";

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

export type RenderNodeType =
  | "form"
  | "section"
  | "table"
  | "grid"
  | "tabs"
  | "text"
  | "recordView"
  | "actionButton"
  | "statCard"
  | "chart"
  | "custom";

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
    /** Direct field name of each row's unique key (F-002 · I-010-001 v0.2.0 §3; default "id"). */
    rowKey?: string;
    /** ADR-0022 multi-select model (registry; mode: multiple only). */
    selection?: { mode?: string };
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

export type RenderNode =
  | RenderFormNode
  | RenderSectionNode
  | RenderGridNode
  | RenderTabsNode
  | RenderTextNode
  | RenderTableNode
  | RenderRecordViewNode
  | RenderActionButtonNode
  | RenderStatCardNode
  | RenderChartNode
  | RenderCustomNode;

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

export type RenderErrorCode =
  | "RENDER_UNKNOWN_NODE_TYPE"
  | "RENDER_INVALID_BODY"
  | "RENDER_META_INVALID"
  | "RENDER_FORM_FIELD_INVALID";

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

export type ActionGateResult =
  | { kind: "ok"; value: boolean }
  | { kind: "error"; error: ActionGateError };

export interface RenderFormFieldGate {
  /** Fields that passed the type whitelist and the version/capability gate. */
  fields: FormControlField[];
  /** Deterministic errors for rejected fields and gate failures. */
  errors: FormControlGateError[];
}

const WHITELISTED_NODE_TYPES = new Set<RenderNodeType>([
  "form",
  "section",
  "table",
  "grid",
  "tabs",
  "text",
  "recordView",
  "actionButton",
  "statCard",
  "chart",
  // GOAL-018: custom nodes dispatch to the registered component map.
  "custom",
]);

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/** Keeps only well-formed recordView field rows (key required). */
export function parseRecordViewFields(raw: unknown): RenderRecordViewField[] {
  if (!Array.isArray(raw)) {
    return [];
  }
  const fields: RenderRecordViewField[] = [];
  const seen = new Set<string>();
  for (const item of raw) {
    if (!isRecord(item) || typeof item.key !== "string" || item.key === "") {
      continue;
    }
    if (seen.has(item.key)) continue;
    seen.add(item.key);
    fields.push({
      key: item.key,
      ...(typeof item.label === "string" ? { label: item.label } : {}),
      ...(typeof item.labelKey === "string" ? { labelKey: item.labelKey } : {}),
    });
  }
  return fields;
}

export function isWhitelistedNodeType(type: string): type is RenderNodeType {
  return WHITELISTED_NODE_TYPES.has(type as RenderNodeType);
}

/**
 * Resolves a dot-path on a record for `form.recordSource.responseMapping`
 * (e.g. `"customer.name"` → record.customer.name). Missing segments or a
 * non-object intermediate return `undefined` (field keeps its default/empty).
 */
export function resolveResponsePath(record: unknown, path: string): unknown {
  if (!isRecord(record) || path === "") {
    return undefined;
  }
  let value: unknown = record;
  for (const segment of path.split(".")) {
    if (isRecord(value)) {
      value = value[segment];
    } else {
      return undefined;
    }
  }
  return value;
}

/** Parses a raw node-level DataRef (v2.9, source:api url + params). */
function parseRenderDataRef(value: unknown): RenderDataRef | undefined {
  if (!isRecord(value)) {
    return undefined;
  }
  const url = value.url;
  if (typeof url !== "string" || url === "") {
    return undefined;
  }
  return {
    url,
    ...(isRecord(value.params) ? { params: value.params as Record<string, unknown> } : {}),
  };
}

/** Normalizes an unknown body value into a typed RenderNode, fail-closed. */
export function parseRenderNode(value: unknown, path: string): RenderNode | RenderError {
  if (!isRecord(value) || typeof value.type !== "string") {
    return { code: "RENDER_INVALID_BODY", path, message: "expected a node object with a type" };
  }
  if (!isWhitelistedNodeType(value.type)) {
    return {
      code: "RENDER_UNKNOWN_NODE_TYPE",
      path,
      message: `node type "${value.type}" is outside the registry renderer whitelist`,
    };
  }
  if (value.type === "form") {
    if (!isRecord(value.props) || !Array.isArray(value.props.fields)) {
      return {
        code: "RENDER_FORM_FIELD_INVALID",
        path: `${path}.props.fields`,
        message: "form nodes require a fields array",
      };
    }
    return {
      type: "form",
      ...(value.id === undefined ? {} : { id: value.id }),
      props: {
        fields: value.props.fields,
        ...(value.props.reactions === undefined ? {} : { reactions: value.props.reactions }),
        ...(value.props.submitLabel === undefined
          ? {}
          : { submitLabel: value.props.submitLabel }),
        ...(typeof value.props.submitLabelKey === "string"
          ? { submitLabelKey: value.props.submitLabelKey }
          : {}),
        ...(typeof value.props.submitAction === "string"
          ? { submitAction: value.props.submitAction }
          : {}),
        ...(value.props.mode === "search" ? { mode: "search" as const } : {}),
        ...(typeof value.props.targetTable === "string"
          ? { targetTable: value.props.targetTable }
          : {}),
        ...(typeof value.props.title === "string" ? { title: value.props.title } : {}),
        ...(typeof value.props.titleKey === "string" ? { titleKey: value.props.titleKey } : {}),
        ...(isRecord(value.props.recordSource)
          ? { recordSource: value.props.recordSource }
          : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderFormNode;
  }
  if (value.type === "custom") {
    if (typeof value.component !== "string" || value.component === "") {
      return {
        code: "RENDER_INVALID_BODY",
        path: `${path}.component`,
        message: "custom nodes require a component key",
      };
    }
    return {
      type: "custom",
      ...(value.id === undefined ? {} : { id: value.id as string }),
      component: value.component,
      ...(isRecord(value.props) ? { props: value.props as Record<string, unknown> } : {}),
    };
  }
  if (value.type === "section" || value.type === "grid" || value.type === "tabs") {
    return {
      type: value.type,
      ...(value.id === undefined ? {} : { id: value.id }),
      ...(isRecord(value.props) ? { props: value.props } : {}),
      ...(Array.isArray(value.children) ? { children: value.children } : {}),
    } as RenderNode;
  }
  if (value.type === "text") {
    const textProps = isRecord(value.props) ? value.props : {};
    return {
      type: "text",
      ...(value.id === undefined ? {} : { id: value.id }),
      props: {
        ...(typeof textProps.text === "string" ? { text: textProps.text } : {}),
        // F-01 (GOAL-003 A-003 F-001): preserve textKey so the default render
        // path localizes text nodes (was dropped before).
        ...(typeof textProps.textKey === "string" ? { textKey: textProps.textKey } : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderTextNode;
  }
  if (value.type === "recordView") {
    const props = isRecord(value.props) ? value.props : {};
    const fields = parseRecordViewFields(props.fields);
    return {
      type: "recordView",
      ...(value.id === undefined ? {} : { id: value.id }),
      props: {
        ...(isRecord(props.record) ? { record: props.record } : {}),
        ...(typeof props.title === "string" ? { title: props.title } : {}),
        ...(typeof props.titleKey === "string" ? { titleKey: props.titleKey } : {}),
        ...(fields !== undefined ? { fields } : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderRecordViewNode;
  }
  if (value.type === "actionButton") {
    const props = isRecord(value.props) ? value.props : {};
    return {
      type: "actionButton",
      ...(value.id === undefined ? {} : { id: value.id }),
      props: {
        ...(typeof props.label === "string" ? { label: props.label } : {}),
        ...(typeof props.labelKey === "string" ? { labelKey: props.labelKey } : {}),
        ...(typeof props.actionId === "string" ? { actionId: props.actionId } : {}),
        ...(typeof props.permissionIntent === "string"
          ? { permissionIntent: props.permissionIntent }
          : {}),
        ...(typeof props.key === "string" ? { key: props.key } : {}),
        ...(typeof props.confirm === "string" ? { confirm: props.confirm } : {}),
        ...(typeof props.confirmKey === "string" ? { confirmKey: props.confirmKey } : {}),
        ...(props.visibleWhen === undefined ? {} : { visibleWhen: props.visibleWhen }),
        ...(props.disabledWhen === undefined ? {} : { disabledWhen: props.disabledWhen }),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderActionButtonNode;
  }
  if (value.type === "statCard") {
    const props = isRecord(value.props) ? value.props : {};
    return {
      type: "statCard",
      ...(value.id === undefined ? {} : { id: value.id }),
      ...(parseRenderDataRef(value.data) === undefined
        ? {}
        : { data: parseRenderDataRef(value.data) }),
      props: {
        ...(typeof props.label === "string" ? { label: props.label } : {}),
        // F-01 (GOAL-003 A-003 F-001): preserve labelKey for runtime i18n.
        ...(typeof props.labelKey === "string" ? { labelKey: props.labelKey } : {}),
        ...(typeof props.unit === "string" ? { unit: props.unit } : {}),
        ...(typeof props.format === "string" ? { format: props.format } : {}),
        ...(typeof props.valueField === "string" ? { valueField: props.valueField } : {}),
        ...(typeof props.dataSource === "string" ? { dataSource: props.dataSource } : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderStatCardNode;
  }
  if (value.type === "chart") {
    const props = isRecord(value.props) ? value.props : {};
    return {
      type: "chart",
      ...(value.id === undefined ? {} : { id: value.id }),
      ...(parseRenderDataRef(value.data) === undefined
        ? {}
        : { data: parseRenderDataRef(value.data) }),
      props: {
        ...(typeof props.chartType === "string" ? { chartType: props.chartType } : {}),
        ...(typeof props.xField === "string" ? { xField: props.xField } : {}),
        ...(typeof props.yField === "string" ? { yField: props.yField } : {}),
        ...(typeof props.dataSource === "string" ? { dataSource: props.dataSource } : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderChartNode;
  }
  return {
    type: "table",
    ...(value.id === undefined ? {} : { id: value.id }),
    ...(parseRenderDataRef(value.data) === undefined
      ? {}
      : { data: parseRenderDataRef(value.data) }),
    props: {
      ...(isRecord(value.props) ? value.props : {}),
    },
    ...(value.children === undefined ? {} : { children: value.children }),
  } as RenderTableNode;
}

/** Collects all form fieldIds from a page document (for reaction resolution). */
export function collectFieldIds(node: RenderNode, into: string[] = []): string[] {
  if (node.type === "form") {
    for (const raw of node.props.fields) {
      if (isRecord(raw) && typeof raw.id === "string") {
        into.push(raw.id);
      }
    }
  }
  for (const child of node.children ?? []) {
    collectFieldIds(child, into);
  }
  return into;
}

/** Resolves the reaction state for a form node's fields. */
export function resolveFormReactions(
  form: RenderFormNode,
  context: Record<string, unknown>,
): ReactionState {
  const fieldIds = form.props.fields
    .map((raw) => (isRecord(raw) && typeof raw.id === "string" ? raw.id : null))
    .filter((id): id is string => id !== null);
  return parseAndEvaluateReactions(form.props.reactions, context, fieldIds);
}

/**
 * Resolves a table action gate expression.
 *
 * Distinguishes an absent property (→ `absentDefault`) from an explicit but
 * invalid expression (→ fail-closed `error`, no silent default). Valid $context
 * expressions are evaluated against the frozen engine; booleans pass through.
 */
export function resolveActionGate(
  expression: unknown,
  context: Record<string, unknown>,
  absentDefault: boolean,
  path: string,
): ActionGateResult {
  if (expression === undefined || expression === null) {
    return { kind: "ok", value: absentDefault };
  }
  if (typeof expression === "boolean") {
    return { kind: "ok", value: expression };
  }
  if (typeof expression === "string") {
    if (isValidExpression(expression)) {
      return { kind: "ok", value: evaluateExpression(expression, context) };
    }
    return {
      kind: "error",
      error: { code: "ACTION_GATE_EXPRESSION_INVALID", path, message: expression },
    };
  }
  return {
    kind: "error",
    error: {
      code: "ACTION_GATE_EXPRESSION_INVALID",
      path,
      message: `expected a boolean or $context expression, got ${typeof expression}`,
    },
  };
}

/**
 * Evaluates a whitelisted table action's visibility/disabled gate.
 * Explicit invalid expressions fail closed (hidden / disabled) and are reported.
 */
export function tableActionGate(
  action: Record<string, unknown>,
  context: Record<string, unknown>,
): { visible: boolean; disabled: boolean; errors: ActionGateError[] } {
  const visible = resolveActionGate(action.visibleWhen, context, true, "visibleWhen");
  const disabled = resolveActionGate(action.disabledWhen, context, false, "disabledWhen");
  const errors = [visible, disabled]
    .filter((result): result is { kind: "error"; error: ActionGateError } => result.kind === "error")
    .map((result) => result.error);
  return {
    visible: visible.kind === "ok" ? visible.value : false,
    disabled: disabled.kind === "ok" ? disabled.value : true,
    errors,
  };
}

/**
 * Parses and gates a form node's raw fields against the D-FORM §5 whitelist
 * and the page meta version/capability rules (F-002 / F-003).
 *
 * Returns only the fields that pass every gate; rejected fields produce
 * deterministic errors instead of being silently rendered (or silently dropped).
 */
export function gateRenderFormFields(
  metaValue: unknown,
  rawFields: unknown,
  path: string,
): RenderFormFieldGate {
  const errors: FormControlGateError[] = [];
  if (!isRecord(metaValue)) {
    return {
      fields: [],
      errors: [{ code: "FORM_META_INVALID", path: "meta", message: "page meta must be an object" }],
    };
  }
  const raw = Array.isArray(rawFields) ? rawFields : null;
  if (raw === null) {
    return {
      fields: [],
      errors: [
        { code: "FORM_FIELDS_INVALID", path, message: "form node requires a fields array" },
      ],
    };
  }

  const fields: FormControlField[] = [];
  for (const [index, entry] of raw.entries()) {
    if (!isRecord(entry) || typeof entry.id !== "string" || typeof entry.type !== "string") {
      errors.push({
        code: "FORM_FIELD_INVALID",
        path: `${path}[${index}]`,
        message: "field requires string id and type",
      });
      continue;
    }
    if (!isWhitelistedFormControl(entry.type)) {
      errors.push({
        code: "FORM_TYPE_NOT_WHITELISTED",
        path: `${path}[${index}].type`,
        message: `field type "${entry.type}" is outside the frozen §5 whitelist`,
      });
      continue;
    }
    fields.push({
      id: entry.id,
      type: entry.type,
      ...(typeof entry.label === "string" ? { label: entry.label } : {}),
      ...(typeof entry.labelKey === "string" ? { labelKey: entry.labelKey } : {}),
      ...(typeof entry.placeholder === "string" ? { placeholder: entry.placeholder } : {}),
      ...(typeof entry.placeholderKey === "string" ? { placeholderKey: entry.placeholderKey } : {}),
      ...(entry.mode === "multiple" ? { mode: "multiple" } : {}),
      ...(typeof entry.startField === "string" ? { startField: entry.startField } : {}),
      ...(typeof entry.endField === "string" ? { endField: entry.endField } : {}),
      ...(typeof entry.min === "number" ? { min: entry.min } : {}),
      ...(typeof entry.max === "number" ? { max: entry.max } : {}),
      ...(typeof entry.step === "number" ? { step: entry.step } : {}),
      ...(typeof entry.precision === "number" ? { precision: entry.precision } : {}),
      // GOAL-014 D-002 §3: pass through submit-time validation constraints
      // (A-003 F-001: without this they were dropped at parse time).
      ...(entry.required === true ? { required: true } : {}),
      // ADR-0040 (since 2.9): readOnly fields render non-editable but keep
      // their value in the submit projection (gated by checkFormCapabilities).
      ...(entry.readOnly === true ? { readOnly: true } : {}),
      ...(typeof entry.pattern === "string" ? { pattern: entry.pattern } : {}),
      ...(typeof entry.minLength === "number" ? { minLength: entry.minLength } : {}),
      ...(typeof entry.maxLength === "number" ? { maxLength: entry.maxLength } : {}),
      ...(typeof entry.format === "string" ? { format: entry.format } : {}),
      ...(typeof entry.action === "string" ? { action: entry.action } : {}),
      ...(typeof entry.actionRef === "string" ? { actionRef: entry.actionRef } : {}),
      ...(typeof entry.accept === "string" ? { accept: entry.accept } : {}),
      ...(typeof entry.maxSize === "number" ? { maxSize: entry.maxSize } : {}),
      ...(entry.multiple === true ? { multiple: true } : {}),
      ...(Array.isArray(entry.options)
        ? {
            options: entry.options
              .filter((option): option is Record<string, unknown> => isRecord(option))
              .filter(
                (option): option is Record<string, unknown> & { value: string } =>
                  typeof option.value === "string",
              )
              .map((option) => ({
                value: option.value,
                ...(typeof option.label === "string" ? { label: option.label } : {}),
                ...(typeof option.labelKey === "string" ? { labelKey: option.labelKey } : {}),
              })),
          }
        : {}),
      ...(entry.defaultValue !== undefined ? { defaultValue: entry.defaultValue } : {}),
    });
  }

  const validFields: FormControlField[] = [];
  for (const field of fields) {
    const gate = checkFormCapabilitiesRaw(metaValue, [field]);
    if (gate.length === 0) {
      validFields.push(field);
    } else {
      errors.push(...gate.map((error) => ({ ...error, path: `fields[${field.id}]` })));
    }
  }
  return { fields: validFields, errors };
}
