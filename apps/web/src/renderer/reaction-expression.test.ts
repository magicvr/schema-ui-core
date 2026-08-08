import { describe, expect, it } from "vitest";

import {
  evaluateFullExpression,
  expressionDependencyFields,
  isValidFullExpression,
  parseExpression,
} from "@/renderer/reaction-expression";
import {
  runReactionEngine,
  runReactionEngineDetailed,
} from "@/renderer/reaction-engine";

const deps = (values: Record<string, unknown>) => ({ deps: values, context: {} });

describe("full expression grammar (02-reaction-expression.md)", () => {
  it("evaluates strict comparisons without coercion (ADR-0016)", () => {
    expect(evaluateFullExpression("$deps.count == 1.0", deps({ count: 1 }))).toBe(true);
    expect(evaluateFullExpression("$deps.count == '1'", deps({ count: 1 }))).toBe(false);
    expect(evaluateFullExpression("$deps.enabled == 1", deps({ enabled: true }))).toBe(false);
    expect(evaluateFullExpression("$deps.count > 0", deps({ count: 1 }))).toBe(true);
    expect(evaluateFullExpression("$deps.count != 2", deps({ count: 1 }))).toBe(true);
    expect(evaluateFullExpression("$deps.count >= 1", deps({ count: 1 }))).toBe(true);
    expect(evaluateFullExpression("$deps.count < 2", deps({ count: 1 }))).toBe(true);
    expect(evaluateFullExpression("$deps.count <= 1", deps({ count: 1 }))).toBe(true);
  });

  it("supports exponent and float literals", () => {
    expect(evaluateFullExpression("$deps.count == 1e0", deps({ count: 1 }))).toBe(true);
    expect(evaluateFullExpression("$deps.count == 1e+21", deps({ count: 1e21 }))).toBe(true);
  });

  it("compares strings by Unicode code points", () => {
    expect(evaluateFullExpression("$deps.astral > '\\ue000'", deps({ astral: "𐀀" }))).toBe(true);
    expect(evaluateFullExpression("$deps.name > 'b'", deps({ name: "a" }))).toBe(false);
  });

  it("contains uses strict element equality on arrays", () => {
    expect(evaluateFullExpression("$deps.roles contains 'admin'", deps({ roles: ["admin", 1] }))).toBe(true);
    expect(evaluateFullExpression("$deps.roles contains 1", deps({ roles: ["admin", 1] }))).toBe(true);
    expect(evaluateFullExpression("$deps.roles contains '1'", deps({ roles: ["admin", 1] }))).toBe(false);
    expect(evaluateFullExpression("$deps.roles contains 'x'", deps({ roles: [] }))).toBe(false);
    // Non-array left operand short-circuits to false.
    expect(evaluateFullExpression("$deps.roles contains 'x'", deps({ roles: "admin" }))).toBe(false);
  });

  it("supports && || ! and grouping", () => {
    const env = deps({ age: 20, hasLicense: true, status: "draft" });
    expect(evaluateFullExpression("$deps.age >= 18 && $deps.hasLicense == true", env)).toBe(true);
    expect(evaluateFullExpression("$deps.age < 18 || $deps.hasLicense == true", env)).toBe(true);
    expect(evaluateFullExpression("!($deps.status == 'draft')", env)).toBe(false);
    expect(evaluateFullExpression("($deps.age >= 18 && $deps.hasLicense == true) && $deps.status != 'x'", env)).toBe(true);
  });

  it("evaluates $context.user / $context.features paths", () => {
    const env = { deps: {}, context: { user: { roles: ["admin", "editor"] }, features: { newDashboard: true } } };
    expect(evaluateFullExpression("$context.user.roles contains 'admin'", env)).toBe(true);
    expect(evaluateFullExpression("$context.features.newDashboard == true", env)).toBe(true);
    expect(evaluateFullExpression("$context.features.missing == true", env)).toBe(false);
    // Deep path on undefined short-circuits to false.
    expect(evaluateFullExpression("$context.user.profile.admin == true", env)).toBe(false);
  });

  it("supports null literals and string escapes", () => {
    expect(evaluateFullExpression("$deps.result == null", deps({ result: null }))).toBe(true);
    expect(evaluateFullExpression("$deps.result != null", deps({ result: "x" }))).toBe(true);
    expect(evaluateFullExpression("$deps.note == 'it\\'s'", deps({ note: "it's" }))).toBe(true);
    expect(evaluateFullExpression('$deps.note == "double"', deps({ note: "double" }))).toBe(true);
  });

  it("fails closed on invalid syntax and forbidden variables", () => {
    expect(evaluateFullExpression("", deps({}))).toBe(false);
    expect(evaluateFullExpression("$deps.x", deps({ x: true }))).toBe(false); // no bare operand
    expect(evaluateFullExpression("$deps.x == true == false", deps({ x: true }))).toBe(false); // non-chainable
    expect(evaluateFullExpression("$row.x == 1", deps({}))).toBe(false); // $row not in form scope
    expect(evaluateFullExpression("window.alert()", deps({}))).toBe(false);
    expect(evaluateFullExpression("$deps.x + 1 == 2", deps({ x: 1 }))).toBe(false); // no arithmetic
  });

  it("parses but reports typed errors", () => {
    const bad = parseExpression("$row.x == 1");
    expect(bad).toMatchObject({ ok: false, code: "FORBIDDEN_VARIABLE" });
    const syntax = parseExpression("$deps.x ===");
    expect(syntax).toMatchObject({ ok: false, code: "SYNTAX" });
  });

  it("extracts $deps dependency fields", () => {
    expect(expressionDependencyFields("$deps.orderType == 'x' && $deps.customerLevel == 'vip'")).toEqual([
      "orderType",
      "customerLevel",
    ]);
    expect(expressionDependencyFields("$context.user.roles contains 'admin'")).toEqual([]);
  });

  it("validates expressions under the full grammar", () => {
    expect(isValidFullExpression("$deps.status == 'closed'")).toBe(true);
    expect(isValidFullExpression("$deps.status ===")).toBe(false);
  });
});

describe("reaction engine unit edges (02 §14)", () => {
  it("stops immediately when nothing commits", () => {
    const result = runReactionEngine({
      initialValues: { a: 1 },
      fields: [{ field: "a", reactions: [{ when: "$deps.a == 2", fulfill: { value: 3 } }] }],
    });
    expect(result).toEqual({
      ok: true,
      values: { a: 1 },
      rounds: [{ round: 1, snapshot: { a: 1 }, observations: {}, commits: [] }],
      warnings: [],
    });
  });

  it("deep-equal writes do not schedule a next round", () => {
    const result = runReactionEngine({
      initialValues: { status: "closed" },
      fields: [{ field: "status", reactions: [{ when: "$deps.status == 'closed'", fulfill: { value: "closed" } }] }],
    });
    expect(result.ok).toBe(true);
    if (result.ok) {
      expect(result.rounds).toHaveLength(1);
      expect(result.rounds[0]!.commits).toEqual([]);
    }
  });

  it("applies visible/disabled state through the detailed runner", () => {
    const detailed = runReactionEngineDetailed({
      initialValues: { trigger: true, note: "" },
      fields: [
        {
          field: "note",
          reactions: [
            { when: "$deps.trigger == true", fulfill: { visible: false, value: "on" } },
          ],
        },
      ],
    });
    expect(detailed.result.ok).toBe(true);
    expect(detailed.fieldStates.note).toMatchObject({ visible: false });
    expect(detailed.fieldStates.note?.disabled).toBeUndefined();
  });
});
