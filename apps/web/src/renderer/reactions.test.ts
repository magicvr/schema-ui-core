import { describe, expect, it } from "vitest";

import {
  evaluateReactions,
  parseAndEvaluateReactions,
  parseReactionRule,
  type ReactionRule,
} from "@/renderer/reactions";

const FIELD_IDS = ["name", "kind", "notes", "enabled", "approval"];

function adminContext() {
  return { user: { roles: ["admin"] }, features: { audit: true } };
}

describe("parseReactionRule", () => {
  it("parses a valid reaction rule", () => {
    const result = parseReactionRule(
      {
        id: "r1",
        when: '$context.features.audit == true',
        apply: [{ fieldId: "notes", visible: true }],
      },
      "reactions[0]",
    );
    expect(result).toEqual({
      id: "r1",
      when: "$context.features.audit == true",
      apply: [{ fieldId: "notes", visible: true }],
    });
  });

  it("rejects non-object rules and missing id/when", () => {
    expect(parseReactionRule("x", "p")).toMatchObject({ code: "REACTION_APPLY_INVALID" });
    expect(parseReactionRule({ id: "r", apply: [] }, "p")).toMatchObject({
      code: "REACTION_APPLY_INVALID",
    });
  });

  it("rejects expressions outside the frozen grammar", () => {
    const result = parseReactionRule(
      { id: "r", when: "$deps.admin == true", apply: [{ fieldId: "x", visible: true }] },
      "p",
    );
    expect(result).toMatchObject({ code: "REACTION_EXPRESSION_INVALID" });
  });

  it("rejects apply entries without fieldId or with non-boolean flags", () => {
    expect(
      parseReactionRule(
        { id: "r", when: "$context.features.a == true", apply: [{ visible: true }] },
        "p",
      ),
    ).toMatchObject({ code: "REACTION_APPLY_INVALID" });
    expect(
      parseReactionRule(
        { id: "r", when: "$context.features.a == true", apply: [{ fieldId: "x", visible: "yes" }] },
        "p",
      ),
    ).toMatchObject({ code: "REACTION_APPLY_INVALID" });
  });
});

describe("evaluateReactions", () => {
  it("keeps every field at its default state when no rule holds", () => {
    const result = evaluateReactions([], adminContext(), FIELD_IDS);
    expect(result.errors).toEqual([]);
    expect(result.state).toEqual({
      name: { visible: true, disabled: false },
      kind: { visible: true, disabled: false },
      notes: { visible: true, disabled: false },
      enabled: { visible: true, disabled: false },
      approval: { visible: true, disabled: false },
    });
  });

  it("applies a rule when the expression holds", () => {
    const rules: ReactionRule[] = [
      {
        id: "r1",
        when: '$context.features.audit == true',
        apply: [{ fieldId: "notes", visible: false }],
      },
    ];
    const result = evaluateReactions(rules, adminContext(), FIELD_IDS);
    expect(result.state.notes).toEqual({ visible: false, disabled: false });
    expect(result.state.name).toEqual({ visible: true, disabled: false });
  });

  it("skips rules whose expression is false", () => {
    const rules: ReactionRule[] = [
      {
        id: "r1",
        when: '$context.user.roles contains "manager"',
        apply: [{ fieldId: "approval", visible: true, disabled: true }],
      },
    ];
    const result = evaluateReactions(rules, adminContext(), FIELD_IDS);
    expect(result.state.approval).toEqual({ visible: true, disabled: false });
  });

  it("fails closed on an unknown apply field", () => {
    const rules: ReactionRule[] = [
      {
        id: "r1",
        when: '$context.features.audit == true',
        apply: [{ fieldId: "ghost", visible: false }],
      },
    ];
    const result = evaluateReactions(rules, adminContext(), FIELD_IDS);
    expect(result.errors).toMatchObject([
      { code: "REACTION_APPLY_FIELD_UNKNOWN", path: "reactions[r1].apply.ghost" },
    ]);
    // Unknown field keeps its default; no throw.
    expect(result.state).not.toHaveProperty("ghost");
  });

  it("combines multiple applies on the same field", () => {
    const rules: ReactionRule[] = [
      {
        id: "r1",
        when: '$context.features.audit == true',
        apply: [{ fieldId: "approval", visible: false, disabled: true }],
      },
    ];
    const result = evaluateReactions(rules, adminContext(), FIELD_IDS);
    expect(result.state.approval).toEqual({ visible: false, disabled: true });
  });
});

describe("parseAndEvaluateReactions", () => {
  it("parses raw rules and applies them", () => {
    const result = parseAndEvaluateReactions(
      [
        {
          id: "r1",
          when: '$context.features.audit == true',
          apply: [{ fieldId: "notes", visible: false }],
        },
      ],
      adminContext(),
      FIELD_IDS,
    );
    expect(result.errors).toEqual([]);
    expect(result.state.notes!.visible).toBe(false);
  });

  it("reports parse errors without throwing and keeps other fields intact", () => {
    const result = parseAndEvaluateReactions(
      [
        { id: "bad", when: "$deps.x == true", apply: [{ fieldId: "notes", visible: false }] },
        {
          id: "good",
          when: '$context.features.audit == true',
          apply: [{ fieldId: "kind", visible: false }],
        },
      ],
      adminContext(),
      FIELD_IDS,
    );
    expect(result.errors).toMatchObject([
      { code: "REACTION_EXPRESSION_INVALID" },
    ]);
    expect(result.state.notes!.visible).toBe(true);
    expect(result.state.kind!.visible).toBe(false);
  });

  it("returns default state for non-array input", () => {
    const result = parseAndEvaluateReactions(undefined, adminContext(), FIELD_IDS);
    expect(result.errors).toEqual([]);
    expect(result.state.name).toEqual({ visible: true, disabled: false });
  });
});
