/**
 * Browser-safe structural validation for vendored page/node/action/reaction
 * schemas. Mirrors `protocol/conformance/schema-validate.ts` but imports the
 * pinned `docs/schemas/*.json` at build time (Vite `@schemas` alias) instead of
 * reading them from disk, so the runtime loader can enforce D-VAL in the
 * browser. The schema set is identical, so runtime and test-time validators
 * stay aligned and neither redefines upstream node/page semantics.
 */
export type RuntimeSchemaKind = "node" | "page" | "action" | "reaction";
export interface RuntimeSchemaValidationResult {
    ok: boolean;
    errors: Array<{
        path: string;
        message: string;
        keyword?: string;
    }>;
}
/**
 * Structural validation of a fetched page document against the pinned page/node
 * schemas. `ok: false` means the document must fail closed and never reach the
 * renderer.
 */
export declare function validatePageDocument(document: unknown): RuntimeSchemaValidationResult;
