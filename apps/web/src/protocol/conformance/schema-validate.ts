/**
 * Structural validation entry for vendored node/page/action/reaction schemas.
 * Uses Ajv draft-07 against pinned schema-ui-docs@2.7.0 artifacts in docs/schemas/.
 */

import { readFileSync } from "node:fs";
import { dirname, join } from "node:path";
import { fileURLToPath } from "node:url";

import Ajv, { type ErrorObject, type ValidateFunction } from "ajv";

export type SchemaKind = "node" | "page" | "action" | "reaction";

export interface SchemaValidationResult {
  ok: boolean;
  errors: Array<{ path: string; message: string; keyword?: string }>;
}

const __dirname = dirname(fileURLToPath(import.meta.url));
const SCHEMAS_DIR = join(__dirname, "../../../../../docs/schemas");

function loadSchema(name: string): Record<string, unknown> {
  const bytes = readFileSync(join(SCHEMAS_DIR, name), "utf8");
  return JSON.parse(bytes) as Record<string, unknown>;
}

let cached: {
  validators: Record<SchemaKind, ValidateFunction>;
} | null = null;

function buildValidators(): Record<SchemaKind, ValidateFunction> {
  const ajv = new Ajv({
    allErrors: true,
    strict: false,
    validateSchema: false,
  });

  const node = loadSchema("node.schema.json");
  const page = loadSchema("page.schema.json");
  const action = loadSchema("action.schema.json");
  const reaction = loadSchema("reaction.schema.json");

  // Register by both $id and relative filenames used in $ref.
  ajv.addSchema(node);
  ajv.addSchema(node, "node.schema.json");
  ajv.addSchema(page);
  ajv.addSchema(page, "page.schema.json");
  ajv.addSchema(action);
  ajv.addSchema(action, "action.schema.json");
  ajv.addSchema(reaction);
  ajv.addSchema(reaction, "reaction.schema.json");

  return {
    node: ajv.getSchema(node.$id as string) ?? ajv.compile(node),
    page: ajv.getSchema(page.$id as string) ?? ajv.compile(page),
    action: ajv.getSchema(action.$id as string) ?? ajv.compile(action),
    reaction: ajv.getSchema(reaction.$id as string) ?? ajv.compile(reaction),
  };
}

function getValidators(): Record<SchemaKind, ValidateFunction> {
  if (!cached) {
    cached = { validators: buildValidators() };
  }
  return cached.validators;
}

function mapErrors(errors: ErrorObject[] | null | undefined): SchemaValidationResult["errors"] {
  if (!errors) {
    return [];
  }
  return errors.map((error) => ({
    path: error.instancePath || "/",
    message: error.message ?? "invalid",
    keyword: error.keyword,
  }));
}

export function validateAgainstSchema(
  kind: SchemaKind,
  document: unknown,
): SchemaValidationResult {
  const validate = getValidators()[kind];
  const ok = validate(document) as boolean;
  return {
    ok,
    errors: ok ? [] : mapErrors(validate.errors),
  };
}

/** Minimal valid page document using §5 whitelist types (for structural smoke). */
export function sampleWhitelistedPage(): Record<string, unknown> {
  return {
    meta: {
      pageId: "sample-form",
      title: "Sample Form",
      protocolVersion: "2.7",
      requiredCapabilities: ["form.controls.extended", "form.controls.advanced"],
    },
    body: {
      type: "form",
      props: {
        fields: [
          { type: "input", field: "name", label: "Name" },
          { type: "textarea", field: "notes", label: "Notes" },
        ],
      },
      children: [
        {
          type: "section",
          props: { title: "Details" },
          children: [
            { type: "text", props: { content: "Hello" } },
            {
              type: "table",
              props: {
                columns: [{ field: "id", label: "ID" }],
              },
            },
            {
              type: "actionButton",
              props: { label: "Save" },
            },
            {
              type: "recordView",
              props: { title: "Record" },
            },
          ],
        },
        {
          type: "grid",
          children: [{ type: "text", props: { content: "Cell" } }],
        },
        {
          type: "tabs",
          props: {
            items: [{ key: "a", label: "A" }],
          },
          children: [{ type: "text", props: { content: "Tab A" } }],
        },
      ],
    },
  };
}
