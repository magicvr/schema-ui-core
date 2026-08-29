/**
 * Full protocol expression engine (schema-ui-docs@2.7.0 · docs/02-reaction-expression.md).
 *
 * Whitelisted grammar (no eval / new Function):
 *   - namespaces: $deps.<field>, $self, $context.user.* / $context.features.*
 *   - operators:  ==  !=  >  >=  <  <=  contains  &&  ||  !
 *   - grouping:   ( )
 *   - literals:   'str' | "str" | number (int/float/exponent) | true | false | null
 *
 * Semantics (ADR-0016 / §7-§14):
 *   - strict typing, no coercion: 1 == 1.0 true; 1 == '1' false; true == 1 false
 *   - contains: array left operand, strict element equality; non-array → false
 *   - string ordering by Unicode code points (not UTF-16 code units)
 *   - undefined values compare as false (missing deps never throw)
 */
export type ExprValue = string | number | boolean | null | ExprValue[] | {
    [key: string]: ExprValue;
} | undefined;
export interface ReactionEnv {
    /** $deps.<field> values (form field snapshot). */
    deps?: Record<string, unknown>;
    /** $self (current field value). */
    self?: unknown;
    /** $context.user.* / $context.features.* snapshots. */
    context?: Record<string, unknown>;
}
export type ParseError = {
    ok: false;
    code: "SYNTAX";
    message: string;
} | {
    ok: false;
    code: "FORBIDDEN_VARIABLE";
    message: string;
} | {
    ok: false;
    code: "UNSUPPORTED_OPERATOR";
    message: string;
};
export type ParsedExpr = {
    kind: "literal";
    value: ExprValue;
} | {
    kind: "var";
    path: string[];
} | {
    kind: "self";
} | {
    kind: "not";
    operand: ParsedExpr;
} | {
    kind: "and";
    left: ParsedExpr;
    right: ParsedExpr;
} | {
    kind: "or";
    left: ParsedExpr;
    right: ParsedExpr;
} | {
    kind: "compare";
    op: "==" | "!=" | ">" | ">=" | "<" | "<=" | "contains";
    left: ParsedExpr;
    right: ParsedExpr;
};
/** Parses a whitelisted expression; returns the AST or a typed error. */
export declare function parseExpression(source: string): ParsedExpr | ParseError;
/** Deep equality (arrays/objects by structure; scalars by strict identity). */
export declare function deepEqual(left: unknown, right: unknown): boolean;
/** Evaluates a whitelisted expression to a boolean; invalid input → false. */
export declare function evaluateFullExpression(source: string, env: ReactionEnv): boolean;
/** True when `source` parses under the full grammar. */
export declare function isValidFullExpression(source: string): boolean;
/** Extracts the $deps.<root> field names referenced by an expression. */
export declare function expressionDependencyFields(source: string): string[];
