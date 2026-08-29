/**
 * request-construction fixture adapter (schema-ui-docs v2.7.0).
 *
 * Builds HTTP request / navigation / modal outcomes from declarative mappings.
 * Batch kinds are out of scope for MVP callers (Q1); this module still implements
 * non-batch kinds used by stage3 execution.
 */
export type RequestConstructionResult = {
    ok: true;
    request?: {
        method: string;
        url: string;
        body: unknown;
        headers?: Record<string, string>;
    };
    navigation?: {
        url: string;
    };
    modalOpen?: {
        modalId: string;
    };
    resolvedBase?: string;
    /** Batch triggers: selection snapshot after a successful reload (ADR-0022 D2). */
    selectionAfterSuccessReload?: {
        keys: unknown[];
        count: number;
    };
} | {
    ok: false;
    code: string;
    path: string;
};
/** D3 invariants: scalar keys only, dedupe preserving order, count = keys.length. */
export declare function normalizeSelection(keys: unknown[]): {
    keys: unknown[];
    count: number;
};
/**
 * Run one request-construction fixture case.
 * Batch kinds return a structured error; stage3 excludes them via Q1.
 */
export declare function constructRequest(input: Record<string, unknown>): RequestConstructionResult;
