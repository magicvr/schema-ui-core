// @vitest-environment jsdom
//
// GOAL-005 R4 · the shared job-result download decision.
//
// Both export entry points (the users-page batch-export component and the jobs
// page's row action) delegate to this module, so these cases are the single
// place where "which bytes, under which name" is decided — and the reason
// I-038-013 (the two entry points naming the same export differently) cannot
// recur.

import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import {
  downloadJobResultDocument,
  jobResultFileName,
  jobResultHasCSV,
} from "@/lib/job-result-download";

class RecordingBlob extends Blob {
  readonly parts: string;
  readonly type: string;
  constructor(parts: BlobPart[], options?: BlobPropertyBag) {
    super(parts, options);
    this.parts = parts.map((part) => String(part)).join("");
    this.type = options?.type ?? "";
  }
}

let blobs: RecordingBlob[] = [];
let downloads: string[] = [];

beforeEach(() => {
  blobs = [];
  downloads = [];
  vi.stubGlobal("Blob", RecordingBlob);
  Object.defineProperty(URL, "createObjectURL", {
    configurable: true,
    writable: true,
    value: (blob: Blob) => {
      blobs.push(blob as RecordingBlob);
      return "blob:mock";
    },
  });
  Object.defineProperty(URL, "revokeObjectURL", {
    configurable: true,
    writable: true,
    value: () => {},
  });
  vi.spyOn(HTMLAnchorElement.prototype, "click").mockImplementation(function (
    this: HTMLAnchorElement,
  ) {
    downloads.push(this.download);
  });
});

afterEach(() => {
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
});

const CSV_DOCUMENT = {
  resource: "users",
  rowCount: 2,
  fileName: "users-selection.csv",
  csv: "id,name\nusr-1,Ada\n",
};

describe("job result download", () => {
  it("recognizes a CSV-carrying result envelope", () => {
    expect(jobResultHasCSV(CSV_DOCUMENT)).toBe(true);
    expect(jobResultHasCSV({ csv: "" })).toBe(false);
    expect(jobResultHasCSV({ csv: 42 })).toBe(false);
    expect(jobResultHasCSV(null)).toBe(false);
    expect(jobResultHasCSV("csv")).toBe(false);
  });

  it("prefers the server's own filename and falls back only when it is missing", () => {
    expect(jobResultFileName(CSV_DOCUMENT, "fallback.csv")).toBe("users-selection.csv");
    expect(jobResultFileName({ fileName: "  " }, "fallback.csv")).toBe("fallback.csv");
    expect(jobResultFileName({}, "fallback.csv")).toBe("fallback.csv");
  });

  it("downloads the CSV bytes as text/csv when the envelope carries them", () => {
    const filename = downloadJobResultDocument(CSV_DOCUMENT, "fallback.csv");
    expect(filename).toBe("users-selection.csv");
    expect(downloads).toEqual(["users-selection.csv"]);
    expect(blobs[0]?.parts).toBe("id,name\nusr-1,Ada\n");
    expect(blobs[0]?.type).toContain("text/csv");
  });

  it("downloads the raw JSON document for any non-export job", () => {
    const payload = { ok: true };
    const filename = downloadJobResultDocument(payload, "job-1.json");
    expect(filename).toBe("job-1.json");
    expect(downloads).toEqual(["job-1.json"]);
    expect(JSON.parse(blobs[0]?.parts ?? "null")).toEqual({ ok: true });
    expect(blobs[0]?.type).toContain("application/json");
  });
});
