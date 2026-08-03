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
 * R5 D-COMP minimal page renderer (frozen §5 whitelist; resolve R4 F-002).
 *
 * The renderer walks a page document's node tree and dispatches whitelisted
 * node types to the components in this directory. It is deliberately minimal:
 * only the node types actually surfaced by the R5 example pages are handled,
 * and any other type fails closed instead of rendering a silent fallback.
 *
 * A page document looks like the example pages' `PAGE_DOCUMENT`:
 *   { meta: { protocolVersion, requiredCapabilities }, body: Node }
 *
 * Node types supported (frozen §5 whitelist):
 *   - layout:  grid / section / tabs
 *   - data/action: text / table / recordView / actionButton
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
  | "actionButton";

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
    /** Default-mode submit: the top-level action id to run on submit (S4). */
    submitAction?: string;
    /** Search-mode form: binds its fields to the target table's query (S4). */
    mode?: "default" | "search";
    targetTable?: string;
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
  };
  children?: RenderNode[];
}

export interface RenderTableNode {
  type: "table";
  id?: string;
  props: {
    columns?: Array<Record<string, unknown>>;
    actions?: Array<Record<string, unknown>>;
    toolbar?: Array<Record<string, unknown>>;
    dataSource?: string;
    /** Direct field name of each row's unique key (F-002 · I-010-001 v0.2.0 §3; default "id"). */
    rowKey?: string;
  };
  children?: RenderNode[];
}

export interface RenderRecordViewNode {
  type: "recordView";
  id?: string;
  props?: {
    record?: Record<string, unknown>;
  };
  children?: RenderNode[];
}

export interface RenderActionButtonNode {
  type: "actionButton";
  id?: string;
  props?: {
    label?: string;
    actionId?: string;
    visibleWhen?: unknown;
    disabledWhen?: unknown;
  };
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
  | RenderActionButtonNode;

export interface RenderPageDocument {
  meta: RenderMeta;
  body: RenderNode;
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
]);

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

export function isWhitelistedNodeType(type: string): type is RenderNodeType {
  return WHITELISTED_NODE_TYPES.has(type as RenderNodeType);
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
      message: `node type "${value.type}" is outside the §5 renderer whitelist`,
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
        ...(typeof value.props.submitAction === "string"
          ? { submitAction: value.props.submitAction }
          : {}),
        ...(value.props.mode === "search" ? { mode: "search" as const } : {}),
        ...(typeof value.props.targetTable === "string"
          ? { targetTable: value.props.targetTable }
          : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderFormNode;
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
    return {
      type: "text",
      ...(value.id === undefined ? {} : { id: value.id }),
      props: {
        ...(isRecord(value.props) && typeof value.props.text === "string"
          ? { text: value.props.text }
          : {}),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderTextNode;
  }
  if (value.type === "recordView") {
    return {
      type: "recordView",
      ...(value.id === undefined ? {} : { id: value.id }),
      props: {
        ...(isRecord(value.props) && isRecord(value.props.record)
          ? { record: value.props.record }
          : {}),
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
        ...(typeof props.actionId === "string" ? { actionId: props.actionId } : {}),
        ...(props.visibleWhen === undefined ? {} : { visibleWhen: props.visibleWhen }),
        ...(props.disabledWhen === undefined ? {} : { disabledWhen: props.disabledWhen }),
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderActionButtonNode;
  }
  return {
    type: "table",
    ...(value.id === undefined ? {} : { id: value.id }),
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
      ...(entry.mode === "multiple" ? { mode: "multiple" } : {}),
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
