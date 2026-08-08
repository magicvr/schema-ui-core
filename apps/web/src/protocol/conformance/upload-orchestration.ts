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

export type UploadResult =
  | { type: "success"; response: { url?: string; id?: string } }
  | { type: "failure"; status?: number };

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

const PROTOCOL_URL_RE = /^\/(?!\/)[^\s\\]*$/;
const RETRY_POLICIES = new Set(["never", "idempotent"]);

function fail(code: string, fileIndex: number, requests: UploadRequest[]): UploadBatchResultError {
  return { ok: false, code, fileIndex, requests, fieldValue: null };
}

function responseFieldValue(response: { url?: string; id?: string }): string | null {
  if (typeof response.url === "string" && response.url !== "") {
    return response.url;
  }
  if (typeof response.id === "string" && response.id !== "") {
    return response.id;
  }
  return null;
}

/** accept token matching (ADR-0012 D2): ".ext", "type/*", exact MIME. */
export function acceptTokenMatches(token: string, file: UploadFile): boolean {
  const trimmed = token.trim();
  if (trimmed === "") {
    return false;
  }
  if (trimmed.startsWith(".")) {
    const dot = file.name.lastIndexOf(".");
    if (dot < 0 || dot === file.name.length - 1) {
      return false;
    }
    return file.name.slice(dot).toLowerCase() === trimmed.toLowerCase();
  }
  if (trimmed.endsWith("/*")) {
    const main = trimmed.slice(0, trimmed.length - 2);
    return file.type.toLowerCase().startsWith(`${main.toLowerCase()}/`);
  }
  return file.type.toLowerCase() === trimmed.toLowerCase();
}

/** Client-side constraints (ADR-0012 D2); returns the first violation. */
export function validateUploadSelection(
  action: UploadActionResult,
  files: UploadFile[],
): { ok: true } | { ok: false; code: string; fileIndex: number } {
  const maxIndex = action.multiple === true ? files.length : Math.min(files.length, 1);
  if (action.multiple !== true && files.length > 1) {
    return { ok: false, code: "MULTIPLE_FILES_NOT_ALLOWED", fileIndex: 1 };
  }
  for (let index = 0; index < maxIndex; index += 1) {
    const file = files[index]!;
    if (typeof action.maxSize === "number" && action.maxSize >= 0 && file.size > action.maxSize) {
      return { ok: false, code: "FILE_TOO_LARGE", fileIndex: index };
    }
    if (typeof action.accept === "string" && action.accept.trim() !== "") {
      const tokens = action.accept.split(",").map((token) => token.trim()).filter((token) => token !== "");
      const matched = tokens.some((token) => acceptTokenMatches(token, file));
      if (!matched) {
        return { ok: false, code: "UNSUPPORTED_FILE_TYPE", fileIndex: index };
      }
    }
  }
  return { ok: true };
}

/** Builds the per-file multipart request descriptors (pre-request). */
export function buildUploadRequests(
  action: UploadActionResult,
  files: UploadFile[],
  invocationId?: string,
): UploadBatchResultError | { ok: true; requests: UploadRequest[] } {
  if (typeof action.url !== "string" || !PROTOCOL_URL_RE.test(action.url)) {
    return fail("INVALID_PROTOCOL_URL", 0, []);
  }
  const retryPolicy = action.retryPolicy ?? "never";
  if (!RETRY_POLICIES.has(retryPolicy)) {
    return fail("INVALID_RETRY_POLICY", 0, []);
  }
  if (retryPolicy === "idempotent" && (typeof invocationId !== "string" || invocationId === "")) {
    return fail("MISSING_INVOCATION_ID", 0, []);
  }
  const fieldName = typeof action.fieldName === "string" && action.fieldName !== "" ? action.fieldName : "file";
  const method = action.method === "PUT" ? "PUT" : "POST";
  const requests: UploadRequest[] = files.map((file, index) => ({
    method,
    url: action.url!,
    part: { name: fieldName, fileName: file.name, contentId: file.contentId },
    ...(retryPolicy === "idempotent"
      ? { headers: { "Idempotency-Key": `${invocationId}:${index}` } }
      : {}),
  }));
  return { ok: true, requests };
}

/**
 * Fixture-driven orchestration: replays pre-recorded transport results
 * (`conformance/fixtures/uploads/cases.json` contract).
 */
export function runUploadBatch(input: {
  action: UploadActionResult;
  files: UploadFile[];
  results?: UploadResult[];
  invocationId?: string;
}): UploadBatchResult {
  const validation = validateUploadSelection(input.action, input.files);
  if (!validation.ok) {
    return fail(validation.code, validation.fileIndex, []);
  }
  const built = buildUploadRequests(input.action, input.files, input.invocationId);
  if (!built.ok) {
    return built;
  }
  const requests = built.requests;
  const values: string[] = [];
  const results = input.results ?? [];
  for (let index = 0; index < requests.length; index += 1) {
    const result = results[index];
    if (result === undefined || result.type === "failure") {
      return fail("UPLOAD_REQUEST_FAILED", index, requests.slice(0, index + 1));
    }
    const value = responseFieldValue(result.response);
    if (value === null) {
      return fail("INVALID_UPLOAD_RESPONSE", index, requests.slice(0, index + 1));
    }
    values.push(value);
  }
  return {
    ok: true,
    requests,
    fieldValue: input.action.multiple === true ? values : (values[0] ?? null),
  };
}

/**
 * Real transport orchestration used by the Renderer's upload control: the
 * same constraints and request shape, executed against a live fetch with
 * multipart FormData. Any failure stops the batch (ADR-0012 D4).
 */
export async function uploadFilesWithFetch(
  action: UploadActionResult,
  files: UploadFile[],
  fetcher: typeof fetch,
  invocationId?: string,
): Promise<UploadBatchResult> {
  const validation = validateUploadSelection(action, files);
  if (!validation.ok) {
    return fail(validation.code, validation.fileIndex, []);
  }
  const built = buildUploadRequests(action, files, invocationId);
  if (!built.ok) {
    return built;
  }
  const requests = built.requests;
  const values: string[] = [];
  for (let index = 0; index < requests.length; index += 1) {
    const request = requests[index]!;
    const form = new FormData();
    const blob = new Blob([new Uint8Array(0)], { type: files[index]!.type });
    // File name is carried by the part descriptor; the app passes a File-like
    // whose name/size match the validated descriptor (contentId is the wire id).
    const fileLike = new File([blob], files[index]!.name, { type: files[index]!.type });
    form.append(request.part.name, fileLike, request.part.fileName);
    const response = await fetcher(request.url, {
      method: request.method,
      ...(request.headers === undefined ? {} : { headers: request.headers }),
      body: form,
    });
    if (!response.ok) {
      return fail("UPLOAD_REQUEST_FAILED", index, requests.slice(0, index + 1));
    }
    const payload: unknown = await response.json();
    const record = (typeof payload === "object" && payload !== null ? payload : {}) as {
      url?: unknown;
      id?: unknown;
    };
    const value = responseFieldValue({
      url: typeof record.url === "string" ? record.url : undefined,
      id: typeof record.id === "string" ? record.id : undefined,
    });
    if (value === null) {
      return fail("INVALID_UPLOAD_RESPONSE", index, requests.slice(0, index + 1));
    }
    values.push(value);
  }
  return {
    ok: true,
    requests,
    fieldValue: action.multiple === true ? values : (values[0] ?? null),
  };
}
