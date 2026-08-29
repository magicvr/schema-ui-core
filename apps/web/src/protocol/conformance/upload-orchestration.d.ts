/**
 * Upload orchestration (schema-ui-docs@2.7.0 · docs/07-actions-contract.md §7
 * + ADR-0012). I-PROTO-FULL-001 D-UPLOAD include.
 *
 * Client contract:
 *   - one multipart request per file (fieldName default "file"; method POST
 *     default, PUT allowed)
 *   - constraints (multiple / maxSize / accept) are checked BEFORE any
 *     request; any violation rejects the whole batch atomically
 *   - accept tokens: ".ext" matches the file extension case-insensitively,
 *     "type/*" matches the MIME main type, other tokens match MIME exactly
 *     (case-insensitive)
 *   - serial uploads in selection order; any failure stops the batch and no
 *     partial field value is committed
 *   - field value: url (priority) or id; multiple → array in selection order
 *   - retryPolicy idempotent → Idempotency-Key header "{invocationId}:{index}"
 */
export interface UploadFile {
    name: string;
    type: string;
    size: number;
    contentId: string;
}
/** Real-transport file: an UploadFile plus the bytes to send (browser File). */
export interface UploadableFile extends UploadFile {
    blob: Blob;
}
export interface UploadActionResult {
    url?: string;
    id?: string;
    method?: string;
    fieldName?: string;
    accept?: string;
    maxSize?: number;
    multiple?: boolean;
    retryPolicy?: string;
}
export interface UploadRequestPart {
    name: string;
    fileName: string;
    contentId: string;
}
export interface UploadRequest {
    method: string;
    url: string;
    part: UploadRequestPart;
    headers?: Record<string, string>;
}
export type UploadResult = {
    type: "success";
    response: {
        url?: string;
        id?: string;
    };
} | {
    type: "failure";
    status?: number;
};
export interface UploadBatchResultOk {
    ok: true;
    requests: UploadRequest[];
    fieldValue: string | string[] | null;
}
export interface UploadBatchResultError {
    ok: false;
    code: string;
    fileIndex: number;
    requests: UploadRequest[];
    fieldValue: null;
}
export type UploadBatchResult = UploadBatchResultOk | UploadBatchResultError;
/** accept token matching (ADR-0012 D2): ".ext", "type/*", exact MIME. */
export declare function acceptTokenMatches(token: string, file: UploadFile): boolean;
/** Client-side constraints (ADR-0012 D2); returns the first violation. */
export declare function validateUploadSelection(action: UploadActionResult, files: UploadFile[]): {
    ok: true;
} | {
    ok: false;
    code: string;
    fileIndex: number;
};
/** Builds the per-file multipart request descriptors (pre-request). */
export declare function buildUploadRequests(action: UploadActionResult, files: UploadFile[], invocationId?: string): UploadBatchResultError | {
    ok: true;
    requests: UploadRequest[];
};
/**
 * Fixture-driven orchestration: replays pre-recorded transport results
 * (`conformance/fixtures/uploads/cases.json` contract).
 */
export declare function runUploadBatch(input: {
    action: UploadActionResult;
    files: UploadFile[];
    results?: UploadResult[];
    invocationId?: string;
}): UploadBatchResult;
/**
 * Real transport orchestration used by the Renderer's upload control: the
 * same constraints and request shape, executed against a live fetch with
 * multipart FormData carrying the actual file bytes (ADR-0012 D4).
 */
export declare function uploadFilesWithFetch(action: UploadActionResult, files: UploadableFile[], fetcher: typeof fetch, invocationId?: string): Promise<UploadBatchResult>;
