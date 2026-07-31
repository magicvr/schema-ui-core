import { describe, expect, it } from "vitest";

import {
  collectFieldIds,
  gateAction,
  isWhitelistedNodeType,
  parseRenderNode,
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
  it("accepts the three whitelisted node types", () => {
    expect(isWhitelistedNodeType("form")).toBe(true);
    expect(isWhitelistedNodeType("section")).toBe(true);
    expect(isWhitelistedNodeType("table")).toBe(true);
  });

  it("rejects everything outside the §5 renderer whitelist", () => {
    expect(isWhitelistedNodeType("chart")).toBe(false);
    expect(isWhitelistedNodeType("modal")).toBe(false);
    expect(isWhitelistedNodeType("grid")).toBe(false);
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

describe("gateAction / tableActionGate", () => {
  it("accepts booleans and evaluates expressions against $context", () => {
    expect(gateAction(true, CONTEXT)).toBe(true);
    expect(gateAction(false, CONTEXT)).toBe(false);
    expect(gateAction('$context.user.roles contains "admin"', CONTEXT)).toBe(true);
    expect(gateAction('$context.user.roles contains "guest"', CONTEXT)).toBe(false);
  });

  it("returns the default for invalid expressions (fail-closed)", () => {
    expect(gateAction("$deps.x == true", CONTEXT)).toBe(true);
    expect(gateAction("$deps.x == true", CONTEXT, false)).toBe(false);
    expect(gateAction(42, CONTEXT)).toBe(true);
  });

  it("gates visible and disabled independently", () => {
    expect(
      tableActionGate(
        {
          visibleWhen: '$context.features.audit == true',
          disabledWhen: '$context.user.roles contains "admin"',
        },
        CONTEXT,
      ),
    ).toEqual({ visible: true, disabled: true });
    expect(
      tableActionGate(
        { visibleWhen: '$context.features.audit == false' },
        CONTEXT,
      ),
    ).toEqual({ visible: false, disabled: false });
  });
});
