/**
 * Job result download (GOAL-005 R4 · VP-038 result center).
 *
 * A succeeded job's result document is served by the shared management-scope
 * route `GET /api/jobs/{id}/result` (R2). For the batch-export kind the document
 * is a JSON envelope `{resource, rowCount, fileName, csv}`; the CSV inside it is
 * what the operator asked for. Any other kind keeps the server's own attachment
 * semantics and downloads as JSON.
 *
 * This module owns that decision ONCE: both consumers — the R3 batch-export
 * component on the users page and the R4 result-center row action in the
 * renderer — call `downloadJobResultDocument`, so a CSV-rendering rule can never
 * diverge between the two entry points (the exact class of drift I-038-013 is
 * about).
 */

/** The subset of the result envelope this module needs. */
export interface JobResultDocument {
  csv?: unknown;
  fileName?: unknown;
  rowCount?: unknown;
  resource?: unknown;
}

function isRecord(value: unknown): value is Record<string, unknown> {
  return typeof value === "object" && value !== null && !Array.isArray(value);
}

/**
 * Reports whether the payload is a batch-export result carrying CSV bytes.
 * Anything else (including a malformed payload) is NOT treated as CSV — the
 * caller falls back to the JSON document rather than downloading an empty file.
 */
export function jobResultHasCSV(payload: unknown): payload is JobResultDocument & { csv: string } {
  return isRecord(payload) && typeof payload.csv === "string" && payload.csv !== "";
}

/**
 * The download filename: the document's own `fileName` (server-derived from the
 * exported resource, and the SAME string the synchronous export uses, so the
 * two export entry points cannot name the same selection differently), else the
 * caller's fallback.
 */
export function jobResultFileName(payload: unknown, fallback: string): string {
  if (isRecord(payload) && typeof payload.fileName === "string" && payload.fileName.trim() !== "") {
    return payload.fileName;
  }
  return fallback;
}

/** Triggers a browser download for an already-fetched blob. */
export function triggerBlobDownload(blob: Blob, filename: string): void {
  const objectUrl = URL.createObjectURL(blob);
  const anchor = document.createElement("a");
  anchor.href = objectUrl;
  anchor.download = filename;
  document.body.appendChild(anchor);
  anchor.click();
  anchor.remove();
  URL.revokeObjectURL(objectUrl);
}

/**
 * Downloads a fetched job result: the CSV bytes when the document carries them,
 * otherwise the raw JSON document. Returns the filename actually used.
 */
export function downloadJobResultDocument(payload: unknown, fallbackName: string): string {
  const filename = jobResultFileName(payload, fallbackName);
  const blob = jobResultHasCSV(payload)
    ? new Blob([payload.csv], { type: "text/csv;charset=utf-8" })
    : new Blob([JSON.stringify(payload)], { type: "application/json" });
  triggerBlobDownload(blob, filename);
  return filename;
}
