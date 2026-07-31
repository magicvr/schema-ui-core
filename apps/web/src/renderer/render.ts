import { evaluateExpression, isValidExpression } from "@/protocol/app-manifest";

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
 * Node types supported (whitelist):
 *   - form        → FormControls (with reactions applied to field state)
 *   - section     → container passthrough
 *   - table       → DataTable (whitelisted; the example page owns data wiring)
 * The form control whitelist itself is enforced by D-FORM
 * (isWhitelistedFormControl / checkFormCapabilities).
 */

export type RenderNodeType = "form" | "section" | "table";

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
  };
  children?: RenderNode[];
}

export interface RenderSectionNode {
  type: "section";
  id?: string;
  props?: Record<string, unknown>;
  children: RenderNode[];
}

export interface RenderTableNode {
  type: "table";
  id?: string;
  props: {
    columns?: Array<Record<string, unknown>>;
    actions?: Array<Record<string, unknown>>;
    dataSource?: string;
  };
  children?: RenderNode[];
}

export type RenderNode = RenderFormNode | RenderSectionNode | RenderTableNode;

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

const WHITELISTED_NODE_TYPES = new Set<RenderNodeType>(["form", "section", "table"]);

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
      },
      ...(value.children === undefined ? {} : { children: value.children }),
    } as RenderFormNode;
  }
  if (value.type === "section") {
    const section: RenderSectionNode = {
      type: "section",
      children: Array.isArray(value.children)
        ? (value.children as RenderNode[])
        : [],
    };
    if (typeof value.id === "string") {
      section.id = value.id;
    }
    if (isRecord(value.props)) {
      section.props = value.props;
    }
    return section;
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

/** Applies a permission gate to a row action via the frozen $context engine. */
export function gateAction(
  expression: unknown,
  context: Record<string, unknown>,
  defaultValue = true,
): boolean {
  if (typeof expression === "boolean") {
    return expression;
  }
  if (typeof expression === "string") {
    if (!isValidExpression(expression)) {
      return defaultValue;
    }
    return evaluateExpression(expression, context);
  }
  return defaultValue;
}

/** Evaluates a whitelisted table action's visibility/disabled gate. */
export function tableActionGate(
  action: Record<string, unknown>,
  context: Record<string, unknown>,
): { visible: boolean; disabled: boolean } {
  const visible = gateAction(action.visibleWhen, context);
  const disabled = gateAction(action.disabledWhen, context, false);
  return { visible, disabled };
}
