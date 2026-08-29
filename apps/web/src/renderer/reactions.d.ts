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
export type ReactionErrorCode = "REACTION_EXPRESSION_INVALID" | "REACTION_APPLY_FIELD_UNKNOWN" | "REACTION_APPLY_INVALID";
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
/** Parses a raw reaction rule, fail-closed on malformed shapes. */
export declare function parseReactionRule(value: unknown, path: string): ReactionRule | ReactionError;
/**
 * Evaluates a reaction rule list against the frozen $context snapshot.
 * Unknown apply fieldIds fail closed (keep default) and are reported.
 */
export declare function evaluateReactions(rules: ReactionRule[], context: Record<string, unknown>, fieldIds: string[]): ReactionEvaluation;
/** Validates raw rules and returns parsed rules + fail-closed errors. */
export declare function parseAndEvaluateReactions(rawRules: unknown, context: Record<string, unknown>, fieldIds: string[]): ReactionEvaluation;
export interface FullReactionResult {
    /** True when the form declares upstream-shaped per-field reactions ($deps). */
    usesFullEngine: boolean;
    /** Per-field control state after convergence (visible/disabled). */
    state: FormControlStateMap;
    /** Value commits to merge into the form values (last-wins, convergent). */
    values: Record<string, unknown>;
    errors: ReactionError[];
}
/**
 * Resolves the full multi-round $deps reaction engine over a form's fields.
 *
 * Upstream shape (02-reaction-expression.md): each field node carries
 * `reactions: [{ when, fulfill, otherwise }]`. Returns the convergent control
 * state + value commits; malformed rules fail closed (reported, not applied).
 */
export declare function resolveFullFormReactions(rawFields: unknown, values: Record<string, unknown>, baselines: Record<string, unknown>): FullReactionResult;
