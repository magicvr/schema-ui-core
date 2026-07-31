import { describe, expect, it } from "vitest";

import {
  collectFieldIds,
  isWhitelistedNodeType,
  parseRenderNode,
  resolveActionGate,
  resolveFormReactions,
  tableActionGate,
  type RenderFormNode,
} from "@/renderer/render";

const CONTEXT = { user: { roles: ["admin"] }, features: { audit: true } };

function formNode(): RenderFormNode {
  return parseRenderNode(
    {
      type: "form",
      id: "f1",
      props: {
        fields: [
          { id: "name", type: "input" },
          { id: "notes", type: "textarea" },
        ],
        reactions: [
          {
            id: "r1",
            when: "$context.features.audit == true",
            apply: [{ fieldId: "notes", visible: false }],
          },
        ],
      },
    },
    "body",
  ) as RenderFormNode;
}

describe("isWhitelistedNodeType", () => {
  it("accepts the frozen §5 whitelisted node types", () => {
    expect(isWhitelistedNodeType("form")).toBe(true);
    expect(isWhitelistedNodeType("section")).toBe(true);
    expect(isWhitelistedNodeType("table")).toBe(true);
    expect(isWhitelistedNodeType("grid")).toBe(true);
    expect(isWhitelistedNodeType("tabs")).toBe(true);
    expect(isWhitelistedNodeType("text")).toBe(true);
    expect(isWhitelistedNodeType("recordView")).toBe(true);
    expect(isWhitelistedNodeType("actionButton")).toBe(true);
  });

  it("rejects everything outside the frozen §5 renderer whitelist", () => {
    expect(isWhitelistedNodeType("chart")).toBe(false);
    expect(isWhitelistedNodeType("modal")).toBe(false);
    expect(isWhitelistedNodeType("upload")).toBe(false);
  });
});

describe("parseRenderNode", () => {
  it("fails closed on non-object bodies and unknown types", () => {
    expect(parseRenderNode(null, "body")).toMatchObject({ code: "RENDER_INVALID_BODY" });
    expect(parseRenderNode({ type: "chart" }, "body")).toMatchObject({
      code: "RENDER_UNKNOWN_NODE_TYPE",
    });
  });

  it("requires a fields array on form nodes", () => {
    expect(parseRenderNode({ type: "form", props: {} }, "body")).toMatchObject({
      code: "RENDER_FORM_FIELD_INVALID",
    });
  });

  it("normalizes form/section/table nodes", () => {
    expect(formNode()).toMatchObject({ type: "form", id: "f1" });
    const section = parseRenderNode(
      { type: "section", children: [{ type: "form", props: { fields: [] } }] },
      "body",
    );
    expect(section).toMatchObject({ type: "section" });
    const table = parseRenderNode({ type: "table", props: { dataSource: "records" } }, "body");
    expect(table).toMatchObject({ type: "table" });
  });
});

describe("collectFieldIds", () => {
  it("collects field ids from nested forms", () => {
    const section = parseRenderNode(
      { type: "section", children: [formNode()] },
      "body",
    ) as { type: "section"; children: RenderFormNode[] };
    expect(collectFieldIds(section)).toEqual(["name", "notes"]);
  });
});

describe("resolveFormReactions", () => {
  it("applies the frozen $context reaction to field state", () => {
    const result = resolveFormReactions(formNode(), CONTEXT);
    expect(result.errors).toEqual([]);
    expect(result.state.notes).toEqual({ visible: false, disabled: false });
    expect(result.state.name).toEqual({ visible: true, disabled: false });
  });

  it("keeps defaults when the expression is false", () => {
    const form = parseRenderNode(
      {
        type: "form",
        props: {
          fields: [{ id: "approval", type: "input" }],
          reactions: [
            {
              id: "r1",
              when: '$context.user.roles contains "manager"',
              apply: [{ fieldId: "approval", visible: false }],
            },
          ],
        },
      },
      "body",
    ) as RenderFormNode;
    const result = resolveFormReactions(form, CONTEXT);
    expect(result.state.approval).toEqual({ visible: true, disabled: false });
  });
});

describe("resolveActionGate", () => {
  it("accepts booleans and evaluates expressions against $context", () => {
    expect(resolveActionGate(true, CONTEXT, true, "visibleWhen")).toEqual({
      kind: "ok",
      value: true,
    });
    expect(resolveActionGate(false, CONTEXT, true, "visibleWhen")).toEqual({
      kind: "ok",
      value: false,
    });
    expect(resolveActionGate('$context.user.roles contains "admin"', CONTEXT, true, "visibleWhen")).toEqual({
      kind: "ok",
      value: true,
    });
    expect(resolveActionGate('$context.user.roles contains "guest"', CONTEXT, true, "visibleWhen")).toEqual({
      kind: "ok",
      value: false,
    });
  });

  it("uses the absent default only when the property is missing, not on invalid expressions", () => {
    // Absent → default applies.
    expect(resolveActionGate(undefined, CONTEXT, true, "visibleWhen")).toEqual({
      kind: "ok",
      value: true,
    });
    expect(resolveActionGate(null, CONTEXT, true, "visibleWhen")).toEqual({
      kind: "ok",
      value: true,
    });
    // Explicit invalid string → error (fail closed), NOT the default.
    const invalid = resolveActionGate("$deps.x == true", CONTEXT, true, "visibleWhen");
    expect(invalid).toEqual({
      kind: "error",
      error: { code: "ACTION_GATE_EXPRESSION_INVALID", path: "visibleWhen", message: "$deps.x == true" },
    });
    // Explicit non-expression value → error (fail closed).
    expect(resolveActionGate(42, CONTEXT, true, "visibleWhen")).toMatchObject({ kind: "error" });
  });
});

describe("tableActionGate", () => {
  it("gates visible and disabled independently", () => {
    expect(
      tableActionGate(
        {
          visibleWhen: '$context.features.audit == true',
          disabledWhen: '$context.user.roles contains "admin"',
        },
        CONTEXT,
      ),
    ).toEqual({ visible: true, disabled: true, errors: [] });
    expect(
      tableActionGate(
        { visibleWhen: '$context.features.audit == false' },
        CONTEXT,
      ),
    ).toEqual({ visible: false, disabled: false, errors: [] });
  });

  it("fails closed with a checkable error on an invalid visibleWhen", () => {
    const result = tableActionGate(
      { visibleWhen: "$deps.x == true" },
      CONTEXT,
    );
    expect(result.visible).toBe(false);
    expect(result.errors).toEqual([
      { code: "ACTION_GATE_EXPRESSION_INVALID", path: "visibleWhen", message: "$deps.x == true" },
    ]);
  });

  it("fails closed with a checkable error on an invalid disabledWhen", () => {
    const result = tableActionGate(
      { visibleWhen: true, disabledWhen: "$deps.y == 1" },
      CONTEXT,
    );
    expect(result.disabled).toBe(true);
    expect(result.errors).toEqual([
      { code: "ACTION_GATE_EXPRESSION_INVALID", path: "disabledWhen", message: "$deps.y == 1" },
    ]);
  });

  it("defaults absent gates without errors", () => {
    expect(tableActionGate({}, CONTEXT)).toEqual({ visible: true, disabled: false, errors: [] });
  });
});
