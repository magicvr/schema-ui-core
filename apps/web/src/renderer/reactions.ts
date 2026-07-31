import { evaluateExpression, isValidExpression } from "@/protocol/app-manifest";

/**
 * R5 D-EXPR reaction surface (frozen Q: reactions operate on the $context
 * namespace only; no field-value triggers).
 *
 * A reaction is a `{ when, apply }` rule: `when` is a frozen $context
 * expression (evaluateExpression grammar), and `apply` turns a form control
 * state on/off when the expression holds. This is the renderer-level
 * equivalent of the navigation visibleWhen gate, applied to form fields.
 *
 * Fails closed: an unparseable `when` expression or an unknown apply target
 * keeps the field at its default state instead of mutating it silently.
 */

export interface ReactionApply {
  fieldId: string;
  /** visible/disabled toggles are explicit booleans (no implicit flip). */
  visible?: boolean;
  disabled?: boolean;
}

export interface ReactionRule {
  id: string;
  when: string;
  apply: ReactionApply[];
}

export type ReactionErrorCode =
  | "REACTION_EXPRESSION_INVALID"
  | "REACTION_APPLY_FIELD_UNKNOWN"
  | "REACTION_APPLY_INVALID";

export interface ReactionError {
  code: ReactionErrorCode;
  path: string;
  message: string;
}

export interface FormControlState {
  visible: boolean;
  disabled: boolean;
}

export type FormControlStateMap = Record<string, FormControlState>;

export interface ReactionEvaluation {
  state: FormControlStateMap;
  errors: ReactionError[];
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

function initialState(fieldIds: string[]): FormControlStateMap {
  const state: FormControlStateMap = {};
  for (const fieldId of fieldIds) {
    state[fieldId] = { visible: true, disabled: false };
  }
  return state;
}

/** Parses a raw reaction rule, fail-closed on malformed shapes. */
export function parseReactionRule(value: unknown, path: string): ReactionRule | ReactionError {
  if (!isRecord(value)) {
    return { code: "REACTION_APPLY_INVALID", path, message: "reaction must be an object" };
  }
  if (typeof value.id !== "string" || typeof value.when !== "string") {
    return { code: "REACTION_APPLY_INVALID", path, message: "reaction requires string id and when" };
  }
  if (!isValidExpression(value.when)) {
    return { code: "REACTION_EXPRESSION_INVALID", path: `${path}.when`, message: value.when };
  }
  if (!Array.isArray(value.apply) || value.apply.length === 0) {
    return { code: "REACTION_APPLY_INVALID", path: `${path}.apply`, message: "apply must be a non-empty array" };
  }
  const apply: ReactionApply[] = [];
  for (const [index, entry] of value.apply.entries()) {
    if (!isRecord(entry) || typeof entry.fieldId !== "string") {
      return {
        code: "REACTION_APPLY_INVALID",
        path: `${path}.apply[${index}]`,
        message: "apply entry requires a fieldId",
      };
    }
    if (entry.visible !== undefined && typeof entry.visible !== "boolean") {
      return {
        code: "REACTION_APPLY_INVALID",
        path: `${path}.apply[${index}].visible`,
        message: "visible must be a boolean",
      };
    }
    if (entry.disabled !== undefined && typeof entry.disabled !== "boolean") {
      return {
        code: "REACTION_APPLY_INVALID",
        path: `${path}.apply[${index}].disabled`,
        message: "disabled must be a boolean",
      };
    }
    apply.push({
      fieldId: entry.fieldId,
      ...(entry.visible === undefined ? {} : { visible: entry.visible }),
      ...(entry.disabled === undefined ? {} : { disabled: entry.disabled }),
    });
  }
  return { id: value.id, when: value.when, apply };
}

/**
 * Evaluates a reaction rule list against the frozen $context snapshot.
 * Unknown apply fieldIds fail closed (keep default) and are reported.
 */
export function evaluateReactions(
  rules: ReactionRule[],
  context: Record<string, unknown>,
  fieldIds: string[],
): ReactionEvaluation {
  const state = initialState(fieldIds);
  const errors: ReactionError[] = [];
  const known = new Set(fieldIds);

  for (const rule of rules) {
    const holds = evaluateExpression(rule.when, context);
    if (!holds) {
      continue;
    }
    for (const apply of rule.apply) {
      if (!known.has(apply.fieldId)) {
        errors.push({
          code: "REACTION_APPLY_FIELD_UNKNOWN",
          path: `reactions[${rule.id}].apply.${apply.fieldId}`,
          message: `unknown field: ${apply.fieldId}`,
        });
        continue;
      }
      if (apply.visible !== undefined) {
        state[apply.fieldId]!.visible = apply.visible;
      }
      if (apply.disabled !== undefined) {
        state[apply.fieldId]!.disabled = apply.disabled;
      }
    }
  }
  return { state, errors };
}

/** Validates raw rules and returns parsed rules + fail-closed errors. */
export function parseAndEvaluateReactions(
  rawRules: unknown,
  context: Record<string, unknown>,
  fieldIds: string[],
): ReactionEvaluation {
  const parsed: ReactionRule[] = [];
  const errors: ReactionError[] = [];
  if (!Array.isArray(rawRules)) {
    return { state: initialState(fieldIds), errors };
  }
  for (const [index, raw] of rawRules.entries()) {
    const result = parseReactionRule(raw, `reactions[${index}]`);
    if ("code" in result) {
      errors.push(result);
      continue;
    }
    parsed.push(result);
  }
  const evaluation = evaluateReactions(parsed, context, fieldIds);
  return { state: evaluation.state, errors: [...errors, ...evaluation.errors] };
}
