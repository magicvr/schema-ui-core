/**
 * Browser-safe structural validation for vendored page/node/action/reaction
 * schemas. Mirrors `protocol/conformance/schema-validate.ts` but imports the
 * pinned `docs/schemas/*.json` at build time (Vite `@schemas` alias) instead of
 * reading them from disk, so the runtime loader can enforce D-VAL in the
 * browser. The schema set is identical, so runtime and test-time validators
 * stay aligned and neither redefines upstream node/page semantics.
 */

import Ajv, { type ErrorObject, type ValidateFunction } from "ajv";

import actionSchema from "@schemas/action.schema.json";
import nodeSchema from "@schemas/node.schema.json";
import pageSchema from "@schemas/page.schema.json";
import reactionSchema from "@schemas/reaction.schema.json";

export type RuntimeSchemaKind = "node" | "page" | "action" | "reaction";

export interface RuntimeSchemaValidationResult {
  ok: boolean;
  errors: Array<{ path: string; message: string; keyword?: string }>;
}

let cached: Record<RuntimeSchemaKind, ValidateFunction> | null = null;

function buildValidators(): Record<RuntimeSchemaKind, ValidateFunction> {
  const ajv = new Ajv({
    allErrors: true,
    strict: false,
    validateSchema: false,
  });

  // Register by both $id and the relative filenames used in $ref so cross-schema
  // references (page -> node -> reaction) resolve exactly like schema-validate.ts.
  // GOAL-018 local extension: custom nodes carry a top-level component key.
  // The extension is applied at validation time so the upstream-pinned
  // node.schema.json artifact (I-PROTO-004, schema-ui-docs@2.9.0) stays
  // byte-identical — this is NOT an upstream protocol change.
  // The upstream node.schema.json types do not know the local extension;
  // cast through unknown to keep the spread clean (local-only addition).
  const extendedNodeSchema = {
    ...nodeSchema,
    properties: {
      ...(nodeSchema.properties ?? {}),
      component: {
        type: "string",
        description:
          "GOAL-018 local extension: custom node component key (renderer custom-components registry)",
      },
    },
  } as unknown as typeof nodeSchema;
  ajv.addSchema(extendedNodeSchema);
  ajv.addSchema(extendedNodeSchema, "node.schema.json");
  ajv.addSchema(pageSchema);
  ajv.addSchema(pageSchema, "page.schema.json");
  ajv.addSchema(actionSchema);
  ajv.addSchema(actionSchema, "action.schema.json");
  ajv.addSchema(reactionSchema);
  ajv.addSchema(reactionSchema, "reaction.schema.json");

  return {
    node: ajv.getSchema(nodeSchema.$id as string) ?? ajv.compile(nodeSchema),
    page: ajv.getSchema(pageSchema.$id as string) ?? ajv.compile(pageSchema),
    action: ajv.getSchema(actionSchema.$id as string) ?? ajv.compile(actionSchema),
    reaction: ajv.getSchema(reactionSchema.$id as string) ?? ajv.compile(reactionSchema),
  };
}

function getValidators(): Record<RuntimeSchemaKind, ValidateFunction> {
  if (!cached) {
    cached = buildValidators();
  }
  return cached;
}

function mapErrors(
  errors: ErrorObject[] | null | undefined,
): RuntimeSchemaValidationResult["errors"] {
  if (!errors) {
    return [];
  }
  return errors.map((error) => ({
    path: error.instancePath || "/",
    message: error.message ?? "invalid",
    keyword: error.keyword,
  }));
}

/**
 * Structural validation of a fetched page document against the pinned page/node
 * schemas. `ok: false` means the document must fail closed and never reach the
 * renderer.
 */
export function validatePageDocument(
  document: unknown,
): RuntimeSchemaValidationResult {
  const validate = getValidators().page;
  const ok = validate(document) as boolean;
  return {
    ok,
    errors: ok ? [] : mapErrors(validate.errors),
  };
}
