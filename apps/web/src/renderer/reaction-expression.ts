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

export type ExprValue =
  | string
  | number
  | boolean
  | null
  | ExprValue[]
  | { [key: string]: ExprValue }
  | undefined;

export interface ReactionEnv {
  /** $deps.<field> values (form field snapshot). */
  deps?: Record<string, unknown>;
  /** $self (current field value). */
  self?: unknown;
  /** $context.user.* / $context.features.* snapshots. */
  context?: Record<string, unknown>;
}

export type ParseError =
  | { ok: false; code: "SYNTAX"; message: string }
  | { ok: false; code: "FORBIDDEN_VARIABLE"; message: string }
  | { ok: false; code: "UNSUPPORTED_OPERATOR"; message: string };

export type ParsedExpr =
  | { kind: "literal"; value: ExprValue }
  | { kind: "var"; path: string[] }
  | { kind: "self" }
  | { kind: "not"; operand: ParsedExpr }
  | { kind: "and"; left: ParsedExpr; right: ParsedExpr }
  | { kind: "or"; left: ParsedExpr; right: ParsedExpr }
  | {
      kind: "compare";
      op: "==" | "!=" | ">" | ">=" | "<" | "<=" | "contains";
      left: ParsedExpr;
      right: ParsedExpr;
    };

// --- Tokenizer ---

type Token =
  | { type: "ident"; value: string }
  | { type: "string"; value: string }
  | { type: "number"; value: number }
  | { type: "bool"; value: boolean }
  | { type: "null" }
  | { type: "op"; value: string }
  | { type: "lparen" }
  | { type: "rparen" }
  | { type: "eof" };

const COMPARISON_OPS = new Set(["==", "!=", ">", ">=", "<", "<=", "contains"]);

function tokenize(source: string): Token[] | ParseError {
  const tokens: Token[] = [];
  let index = 0;
  const fail = (message: string): ParseError => ({ ok: false, code: "SYNTAX", message });

  while (index < source.length) {
    const char = source[index]!;
    if (/\s/.test(char)) {
      index += 1;
      continue;
    }
    if (char === "(") {
      tokens.push({ type: "lparen" });
      index += 1;
      continue;
    }
    if (char === ")") {
      tokens.push({ type: "rparen" });
      index += 1;
      continue;
    }
    if (char === "'" || char === '"') {
      const quote = char;
      let value = "";
      let closed = false;
      index += 1;
      while (index < source.length) {
        const current = source[index]!;
        if (current === "\\" && index + 1 < source.length) {
          const next = source[index + 1]!;
          if (next === quote || next === "\\") {
            value += next;
            index += 2;
            continue;
          }
          value += next;
          index += 2;
          continue;
        }
        if (current === quote) {
          closed = true;
          index += 1;
          break;
        }
        value += current;
        index += 1;
      }
      if (!closed) {
        return fail(`unterminated string literal at ${index}`);
      }
      tokens.push({ type: "string", value });
      continue;
    }
    if (/[0-9]/.test(char) || (char === "." && /[0-9]/.test(source[index + 1] ?? ""))) {
      const match = /^(\d+(?:\.\d+)?(?:[eE][+-]?\d+)?)/.exec(source.slice(index));
      if (match === null) {
        return fail(`invalid number literal at ${index}`);
      }
      tokens.push({ type: "number", value: Number(match[1]) });
      index += match[1].length;
      continue;
    }
    const two = source.slice(index, index + 2);
    if (two === "==" || two === "!=" || two === ">=" || two === "<=" || two === "&&" || two === "||") {
      tokens.push({ type: "op", value: two });
      index += 2;
      continue;
    }
    if (char === ">" || char === "<" || char === "!") {
      tokens.push({ type: "op", value: char });
      index += 1;
      continue;
    }
    if (/[a-zA-Z_$]/.test(char)) {
      const match = /^[a-zA-Z_$][a-zA-Z0-9_.$]*/.exec(source.slice(index));
      if (match === null) {
        return fail(`invalid identifier at ${index}`);
      }
      const word = match[0];
      if (word === "true" || word === "false") {
        tokens.push({ type: "bool", value: word === "true" });
      } else if (word === "null") {
        tokens.push({ type: "null" });
      } else if (word === "contains") {
        tokens.push({ type: "op", value: "contains" });
      } else {
        tokens.push({ type: "ident", value: word });
      }
      index += word.length;
      continue;
    }
    return fail(`unexpected character "${char}" at ${index}`);
  }
  tokens.push({ type: "eof" });
  return tokens;
}

// --- Parser (recursive descent, comparison non-chainable per §3) ---

class Parser {
  private readonly tokens: Token[];
  private position = 0;

  constructor(tokens: Token[]) {
    this.tokens = tokens;
  }

  private peek(): Token {
    return this.tokens[this.position]!;
  }

  private next(): Token {
    const token = this.tokens[this.position]!;
    this.position += 1;
    return token;
  }

  parse(): ParsedExpr | ParseError {
    const expr = this.parseOr();
    if ("ok" in expr && !expr.ok) {
      return expr;
    }
    if (this.peek().type !== "eof") {
      return { ok: false, code: "SYNTAX", message: `unexpected token after expression (position ${this.position})` };
    }
    return expr as ParsedExpr;
  }

  private parseOr(): ParsedExpr | ParseError {
    let left = this.parseAnd();
    if ("ok" in left && !left.ok) {
      return left;
    }
    while (true) {
      const token = this.peek();
      if (token.type !== "op" || token.value !== "||") {
        break;
      }
      this.next();
      const right = this.parseAnd();
      if ("ok" in right && !right.ok) {
        return right;
      }
      left = { kind: "or", left: left as ParsedExpr, right: right as ParsedExpr };
    }
    return left;
  }

  private parseAnd(): ParsedExpr | ParseError {
    let left = this.parseUnary();
    if ("ok" in left && !left.ok) {
      return left;
    }
    while (true) {
      const token = this.peek();
      if (token.type !== "op" || token.value !== "&&") {
        break;
      }
      this.next();
      const right = this.parseUnary();
      if ("ok" in right && !right.ok) {
        return right;
      }
      left = { kind: "and", left: left as ParsedExpr, right: right as ParsedExpr };
    }
    return left;
  }

  private parseUnary(): ParsedExpr | ParseError {
    const token = this.peek();
    if (token.type === "op" && token.value === "!") {
      this.next();
      const operand = this.parseUnary();
      if ("ok" in operand && !operand.ok) {
        return operand;
      }
      return { kind: "not", operand: operand as ParsedExpr };
    }
    if (token.type === "lparen") {
      this.next();
      const inner = this.parseOr();
      if ("ok" in inner && !inner.ok) {
        return inner;
      }
      if (this.peek().type !== "rparen") {
        return { ok: false, code: "SYNTAX", message: "missing closing parenthesis" };
      }
      this.next();
      return inner as ParsedExpr;
    }
    return this.parseComparison();
  }

  private parseComparison(): ParsedExpr | ParseError {
    const left = this.parseOperand();
    if ("ok" in left && !left.ok) {
      return left;
    }
    const token = this.peek();
    if (token.type === "op" && COMPARISON_OPS.has(token.value)) {
      this.next();
      const right = this.parseOperand();
      if ("ok" in right && !right.ok) {
        return right;
      }
      return {
        kind: "compare",
        op: token.value as "==" | "!=" | ">" | ">=" | "<" | "<=" | "contains",
        left: left as ParsedExpr,
        right: right as ParsedExpr,
      };
    }
    // Comparisons are not chainable; a bare operand in boolean position is
    // invalid (the grammar requires an explicit comparison or ! grouping).
    return { ok: false, code: "SYNTAX", message: `expected a comparison operator, got ${token.type}` };
  }

  private parseOperand(): ParsedExpr | ParseError {
    const token = this.next();
    switch (token.type) {
      case "string":
      case "number":
      case "bool":
        return { kind: "literal", value: token.value };
      case "null":
        return { kind: "literal", value: null };
      case "ident": {
        const variable = token.value;
        if (variable === "$self") {
          return { kind: "self" };
        }
        if (variable.startsWith("$deps.")) {
          const path = variable.slice("$deps.".length).split(".");
          if (path.some((part) => part === "")) {
            return { ok: false, code: "FORBIDDEN_VARIABLE", message: `invalid $deps path "${variable}"` };
          }
          return { kind: "var", path: ["$deps", ...path] };
        }
        if (variable.startsWith("$context.")) {
          const rest = variable.slice("$context.".length);
          const [root, ...tail] = rest.split(".");
          if (root !== "user" && root !== "features") {
            return {
              ok: false,
              code: "FORBIDDEN_VARIABLE",
              message: `$context.${root} is outside the whitelist (user/features)`,
            };
          }
          return { kind: "var", path: ["$context", root, ...tail] };
        }
        return { ok: false, code: "FORBIDDEN_VARIABLE", message: `unknown variable "${variable}"` };
      }
      default:
        return { ok: false, code: "SYNTAX", message: `unexpected token in operand position` };
    }
  }
}

/** Parses a whitelisted expression; returns the AST or a typed error. */
export function parseExpression(source: string): ParsedExpr | ParseError {
  const trimmed = source.trim();
  if (trimmed === "") {
    return { ok: false, code: "SYNTAX", message: "empty expression" };
  }
  const tokens = tokenize(trimmed);
  if (!Array.isArray(tokens)) {
    return tokens;
  }
  return new Parser(tokens).parse();
}

// --- Evaluation ---

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/** Deep equality (arrays/objects by structure; scalars by strict identity). */
export function deepEqual(left: unknown, right: unknown): boolean {
  if (Object.is(left, right)) {
    return true;
  }
  if (Array.isArray(left) && Array.isArray(right)) {
    return left.length === right.length && left.every((entry, index) => deepEqual(entry, right[index]));
  }
  if (isRecord(left) && isRecord(right)) {
    const leftKeys = Object.keys(left).sort();
    const rightKeys = Object.keys(right).sort();
    return (
      leftKeys.length === rightKeys.length &&
      leftKeys.every((key, index) => key === rightKeys[index] && deepEqual(left[key], right[key]))
    );
  }
  return false;
}

/** Unicode code-point ordering for strings (fixture: U+10000 > U+E000). */
function compareCodePoints(left: string, right: string): number {
  const a = [...left];
  const b = [...right];
  const length = Math.min(a.length, b.length);
  for (let index = 0; index < length; index += 1) {
    const ca = a[index]!.codePointAt(0)!;
    const cb = b[index]!.codePointAt(0)!;
    if (ca !== cb) {
      return ca < cb ? -1 : 1;
    }
  }
  return a.length === b.length ? 0 : a.length < b.length ? -1 : 1;
}

function resolveVar(path: string[], env: ReactionEnv): unknown {
  if (path[0] === "$deps") {
    return env.deps?.[path.slice(1).join(".")];
  }
  // $context.user.<tail>
  let current: unknown = env.context?.[path[1] ?? ""];
  for (const part of path.slice(2)) {
    if (!isRecord(current)) {
      return undefined;
    }
    current = current[part];
  }
  return current;
}

function evaluateNode(node: ParsedExpr, env: ReactionEnv): unknown {
  switch (node.kind) {
    case "literal":
      return node.value;
    case "var":
      return resolveVar(node.path, env);
    case "self":
      return env.self;
    case "not":
      return !evaluateNode(node.operand, env);
    case "and":
      return evaluateNode(node.left, env) === true && evaluateNode(node.right, env) === true;
    case "or":
      return evaluateNode(node.left, env) === true || evaluateNode(node.right, env) === true;
    case "compare": {
      const left = evaluateNode(node.left, env);
      const right = evaluateNode(node.right, env);
      return compareValues(node.op, left, right);
    }
  }
}

/** Scalar comparison with strict typing (ADR-0016); undefined → false. */
function compareValues(
  op: "==" | "!=" | ">" | ">=" | "<" | "<=" | "contains",
  left: unknown,
  right: unknown,
): boolean {
  switch (op) {
    case "==":
      return deepEqual(left, right);
    case "!=":
      return !deepEqual(left, right);
    case "contains":
      return Array.isArray(left) ? left.some((entry) => deepEqual(entry, right)) : false;
    case ">":
    case ">=":
    case "<":
    case "<=": {
      if (left === undefined || right === undefined) {
        return false;
      }
      if (typeof left === "number" && typeof right === "number") {
        return op === ">" ? left > right : op === ">=" ? left >= right : op === "<" ? left < right : left <= right;
      }
      if (typeof left === "string" && typeof right === "string") {
        const order = compareCodePoints(left, right);
        return op === ">" ? order > 0 : op === ">=" ? order >= 0 : op === "<" ? order < 0 : order <= 0;
      }
      return false;
    }
  }
}

/** Evaluates a whitelisted expression to a boolean; invalid input → false. */
export function evaluateFullExpression(source: string, env: ReactionEnv): boolean {
  const parsed = parseExpression(source);
  if ("ok" in parsed && !parsed.ok) {
    return false;
  }
  return evaluateNode(parsed as ParsedExpr, env) === true;
}

/** True when `source` parses under the full grammar. */
export function isValidFullExpression(source: string): boolean {
  const parsed = parseExpression(source);
  return !("ok" in parsed && !parsed.ok);
}

/** Extracts the $deps.<root> field names referenced by an expression. */
export function expressionDependencyFields(source: string): string[] {
  const parsed = parseExpression(source);
  if ("ok" in parsed && !parsed.ok) {
    return [];
  }
  const fields: string[] = [];
  const visit = (node: ParsedExpr): void => {
    switch (node.kind) {
      case "var":
        if (node.path[0] === "$deps" && node.path.length >= 2) {
          fields.push(node.path[1]!);
        }
        return;
      case "not":
        visit(node.operand);
        return;
      case "and":
      case "or":
        visit(node.left);
        visit(node.right);
        return;
      case "compare":
        visit(node.left);
        visit(node.right);
        return;
      default:
        return;
    }
  };
  visit(parsed as ParsedExpr);
  return [...new Set(fields)];
}
