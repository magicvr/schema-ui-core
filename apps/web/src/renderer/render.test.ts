import { describe, expect, it } from "vitest";

import {
  collectFieldIds,
  gateRenderFormFields,
  isWhitelistedNodeType,
  parseRecordViewFields,
  parseRenderNode,
  resolveActionGate,
  resolveFormReactions,
  resolveResponsePath,
  tableActionGate,
  type RenderFormNode,
  type RenderRecordViewNode,
} from "@/renderer/render.types";

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

  it("accepts the full registry display types (I-PROTO-FULL-001)", () => {
    expect(isWhitelistedNodeType("statCard")).toBe(true);
    expect(isWhitelistedNodeType("chart")).toBe(true);
  });

  it("rejects everything outside the registry renderer whitelist", () => {
    expect(isWhitelistedNodeType("modal")).toBe(false);
    expect(isWhitelistedNodeType("upload")).toBe(false);
    expect(isWhitelistedNodeType("slider")).toBe(false);
  });

  // GOAL-018: custom nodes are whitelisted and dispatch to the registry.
  it("accepts the custom node type (GOAL-018)", () => {
    expect(isWhitelistedNodeType("custom")).toBe(true);
    expect(parseRenderNode({ type: "custom", component: "mfa-manager" }, "body")).toMatchObject({
      type: "custom",
      component: "mfa-manager",
    });
    expect(parseRenderNode({ type: "custom", props: {} }, "body")).toMatchObject({
      code: "RENDER_INVALID_BODY",
    });
  });
});

describe("parseRenderNode", () => {
  it("fails closed on non-object bodies and unknown types", () => {
    expect(parseRenderNode(null, "body")).toMatchObject({ code: "RENDER_INVALID_BODY" });
    expect(parseRenderNode({ type: "slider" }, "body")).toMatchObject({
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
    // Shape-only test: parseRenderNode normalizes regardless of dataSource
    // validity; runtime schemaTableDataSource rejects non-rooted paths.
    const table = parseRenderNode({ type: "table", props: { dataSource: "/api/users" } }, "body");
    expect(table).toMatchObject({ type: "table" });
  });

  it("preserves textKey on text nodes and labelKey on statCard nodes (F-01)", () => {
    const textNode = parseRenderNode(
      { type: "text", props: { text: "Fallback", textKey: "schema.dashboard.text.intro" } },
      "body",
    );
    expect(textNode).toMatchObject({
      type: "text",
      props: { text: "Fallback", textKey: "schema.dashboard.text.intro" },
    });
    const statCard = parseRenderNode(
      {
        type: "statCard",
        id: "total",
        props: { label: "Users", labelKey: "schema.dashboard.statCard.users", format: "plain", valueField: "total", dataSource: "/api/users" },
      },
      "body",
    );
    expect(statCard).toMatchObject({
      type: "statCard",
      props: { label: "Users", labelKey: "schema.dashboard.statCard.users", valueField: "total", dataSource: "/api/users" },
    });
  });

  it("normalizes statCard and chart nodes with their registry props", () => {
    const statCard = parseRenderNode(
      {
        type: "statCard",
        id: "total",
        props: { label: "Total", format: "currency", valueField: "amount", dataSource: "/api/users" },
      },
      "body",
    );
    expect(statCard).toMatchObject({
      type: "statCard",
      id: "total",
      props: { label: "Total", format: "currency", valueField: "amount", dataSource: "/api/users" },
    });
    const chart = parseRenderNode(
      { type: "chart", props: { chartType: "bar", xField: "month", yField: "count", dataSource: "/api/users" } },
      "body",
    );
    expect(chart).toMatchObject({
      type: "chart",
      props: { chartType: "bar", xField: "month", yField: "count", dataSource: "/api/users" },
    });
  });

  it("keeps recordView title/titleKey/fields and drops malformed field rows", () => {
    const node = parseRenderNode(
      {
        type: "recordView",
        id: "user-detail",
        props: {
          title: "User details",
          titleKey: "schema.users.detail.title",
          fields: [
            { key: "username", label: "Username", labelKey: "schema.users.column.username" },
            { key: "", label: "skip" },
            { label: "no-key" },
          ],
        },
      },
      "body",
    ) as RenderRecordViewNode;
    expect(node).toMatchObject({
      type: "recordView",
      id: "user-detail",
      props: {
        title: "User details",
        titleKey: "schema.users.detail.title",
        fields: [
          { key: "username", label: "Username", labelKey: "schema.users.column.username" },
        ],
      },
    });
    expect(parseRecordViewFields([{ key: "id" }, null, "x"])).toEqual([{ key: "id" }]);
  });

  it("drops non-registry props from statCard/chart nodes (fail-closed on shape)", () => {
    const statCard = parseRenderNode(
      { type: "statCard", props: { label: 42, bogus: "x", valueField: "v", dataSource: "/api/users" } },
      "body",
    );
    expect(statCard).toMatchObject({ type: "statCard", props: { valueField: "v", dataSource: "/api/users" } });
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

describe("resolveResponsePath (form.recordSource.responseMapping · S6)", () => {
  it("resolves a nested dot-path and an identity path", () => {
    const record = { id: "default", customer: { name: "Acme" } };
    expect(resolveResponsePath(record, "customer.name")).toBe("Acme");
    expect(resolveResponsePath(record, "id")).toBe("default");
  });

  it("returns undefined for a missing segment or a non-object intermediate", () => {
    const record = { customer: { name: "Acme" }, tags: ["a"] };
    expect(resolveResponsePath(record, "customer.missing")).toBeUndefined();
    expect(resolveResponsePath(record, "customer.name.deep")).toBeUndefined();
    expect(resolveResponsePath(record, "tags.length")).toBeUndefined();
  });

  it("returns undefined for non-object inputs or an empty path", () => {
    expect(resolveResponsePath(null, "a")).toBeUndefined();
    expect(resolveResponsePath("x", "a")).toBeUndefined();
    expect(resolveResponsePath({ a: 1 }, "")).toBeUndefined();
  });
});

describe("GOAL-014 · constraint passthrough + fieldErrors echo (A-003 F-001/F-002)", () => {
  it("gateRenderFormFields preserves submit-time constraints (F-001)", () => {
    const result = gateRenderFormFields(
      { protocolVersion: "2.7", requiredCapabilities: ["form.controls.extended"] },
      [
        {
          id: "name",
          type: "input",
          required: true,
          pattern: "^[a-z]+$",
          minLength: 2,
          maxLength: 10,
        },
      ],
      "fields",
    );
    expect(result.errors).toEqual([]);
    const field = result.fields[0];
    expect(field.required).toBe(true);
    expect(field.pattern).toBe("^[a-z]+$");
    expect(field.minLength).toBe(2);
    expect(field.maxLength).toBe(10);
  });

  it("preserves numeric and ISO-date bounds for their registry control types", () => {
    const result = gateRenderFormFields(
      { protocolVersion: "2.7", requiredCapabilities: [] },
      [
        { id: "quantity", type: "inputNumber", min: 1, max: 100 },
        { id: "expiresAt", type: "datePicker", min: "2001-09-09", max: "2099-12-31" },
      ],
      "fields",
    );
    expect(result.errors).toEqual([]);
    expect(result.fields).toEqual([
      expect.objectContaining({ id: "quantity", min: 1, max: 100 }),
      expect.objectContaining({ id: "expiresAt", min: "2001-09-09", max: "2099-12-31" }),
    ]);
  });
});

describe("gateRenderFormFields · afterComponent (W17)", () => {
  it("preserves afterComponent on input fields", () => {
    const result = gateRenderFormFields(
      { protocolVersion: "2.7", requiredCapabilities: [] },
      [{ id: "cron", type: "input", afterComponent: "cron-preview" }],
      "fields",
    );
    expect(result.errors).toEqual([]);
    expect(result.fields[0]?.afterComponent).toBe("cron-preview");
  });
});

describe("gateRenderFormFields · readOnly (ADR-0040)", () => {
  it("passes readOnly through and gates it on 2.9 + capability", () => {
    const meta = { protocolVersion: "2.9", requiredCapabilities: ["form.controls.readonly"] };
    const ok = gateRenderFormFields(meta, [{ id: "dictKey", type: "input", readOnly: true }], "fields");
    expect(ok.errors).toEqual([]);
    expect(ok.fields[0]?.readOnly).toBe(true);
  });

  it("rejects readOnly fields when the capability is missing", () => {
    const meta = { protocolVersion: "2.9", requiredCapabilities: [] };
    const bad = gateRenderFormFields(meta, [{ id: "dictKey", type: "input", readOnly: true }], "fields");
    expect(bad.fields.length).toBe(0);
    expect(bad.errors.some((e) => e.code === "FORM_CAPABILITY_REQUIRED")).toBe(true);
  });
});
