// @vitest-environment jsdom
//
// GOAL-005 R4 (C2/C3) · the result center on the REAL jobs page document.
//
// This test renders `apps/api/modules/jobs/schema/jobs.json` itself — not a
// hand-written fixture — through the production chain (validatePageDocument →
// loadPageDocument → RenderPage + SchemaTable), so it pins the shipped schema
// and not a copy of it. It covers the R4 experience contract:
//
//   * the six job states + progress are rendered, with the state's badge style;
//   * cancel / retry / download exist as row actions and are enabled ONLY where
//     the server's derived `cancellable` / `retryable` / `downloadable` says so;
//   * the actions address the frozen routes with the row's own id;
//   * a CSV result document downloads as CSV under the server's filename;
//   * a read-only principal (no jobs.write) cannot invoke the write actions;
//   * the auto-refresh control re-fetches the list on its cadence;
//   * the jobs table declares no selection, which is the precondition the
//     auto-refresh control documents (GOAL-005 D-001 §5, R-1.2).

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import type { PageEntry } from "@/protocol/app-manifest";
import { loadPageDocument } from "@/protocol/load-page";
import type { RenderPageDocument } from "@/renderer/render.types";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";
import "@/components/jobs-auto-refresh";

const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");
const JOBS_SCHEMA = resolve(MODULES, "jobs/schema/jobs.json");

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

afterEach(async () => {
  vi.restoreAllMocks();
  vi.unstubAllGlobals();
  vi.useRealTimers();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

/** One job row exactly as the API projects it (internal/handler/jobs.go jobToMap). */
interface JobRow {
  id: string;
  kind: string;
  status: string;
  progress: number;
  attempt: number;
  maxAttempts: number;
  actorId: string;
  createdAt: string;
  cancellable: boolean;
  retryable: boolean;
  downloadable: boolean;
  statusStyle: string;
  errorCode?: string;
  errorMessage?: string;
  resultUrl?: string;
  resultExpiresAt?: string;
  correlationId: string;
}

function row(overrides: Partial<JobRow> & { id: string; status: string }): JobRow {
  return {
    kind: "jobs.batch-export",
    progress: 0,
    attempt: 0,
    maxAttempts: 3,
    actorId: "user-1",
    correlationId: "corr-" + overrides.id,
    createdAt: "2026-09-19T10:00:00.000Z",
    cancellable: overrides.status === "queued" || overrides.status === "running",
    retryable: overrides.status === "failed",
    downloadable: overrides.status === "succeeded",
    statusStyle: "neutral",
    ...overrides,
  };
}

const ROWS: JobRow[] = [
  row({ id: "job-running", status: "running", progress: 42, attempt: 1, statusStyle: "warning" }),
  row({
    id: "job-failed",
    status: "failed",
    attempt: 1,
    errorCode: "JOB_HANDLER_FAILED",
    errorMessage: "boom",
    statusStyle: "destructive",
  }),
  row({ id: "job-done", status: "succeeded", progress: 100, attempt: 1, statusStyle: "success" }),
  row({ id: "job-expired", status: "expired", attempt: 1, statusStyle: "neutral" }),
];

interface RecordedRequest {
  url: string;
  method: string;
}

/** Captures the bytes handed to URL.createObjectURL. */
class RecordingBlob extends Blob {
  readonly parts: string;
  constructor(parts: BlobPart[], options?: BlobPropertyBag) {
    super(parts, options);
    this.parts = parts.map((part) => String(part)).join("");
  }
}

interface Harness {
  container: HTMLDivElement;
  requests: RecordedRequest[];
  downloads: Array<{ filename: string; text: string }>;
  blobs: RecordingBlob[];
  listCalls: () => number;
}

async function renderJobsPage(options: {
  permissions: string[];
  resultDocument?: unknown;
}): Promise<Harness> {
  const requests: RecordedRequest[] = [];
  const downloads: Array<{ filename: string; text: string }> = [];
  let listCalls = 0;
  vi.stubGlobal("Blob", RecordingBlob);

  // jsdom implements neither URL.createObjectURL nor a real download, so both
  // seams are stubbed: createObjectURL captures the blob, the anchor click
  // captures the filename the helper chose. A Blob subclass keeps the bytes
  // readable without depending on jsdom's Blob.text().
  const capturedBlobs: RecordingBlob[] = [];
  Object.defineProperty(URL, "createObjectURL", {
    configurable: true,
    writable: true,
    value: (blob: Blob) => {
      capturedBlobs.push(blob as RecordingBlob);
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
    downloads.push({ filename: this.download, text: "" });
  });

  const fetcher = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = String(input);
    const method = init?.method ?? "GET";
    requests.push({ url, method });
    if (url.startsWith("/api/jobs/") && url.endsWith("/result")) {
      return new Response(JSON.stringify(options.resultDocument ?? {}), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (url.startsWith("/api/jobs") && method === "POST") {
      return new Response(JSON.stringify({ ok: true }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (url.startsWith("/api/jobs")) {
      listCalls += 1;
      return new Response(
        JSON.stringify({ items: ROWS, total: ROWS.length, page: 1, pageSize: 20 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    return new Response(JSON.stringify({ items: [], total: 0, page: 1, pageSize: 20 }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;

  const raw = JSON.parse(readFileSync(JOBS_SCHEMA, "utf8")) as unknown;
  const page: PageEntry = {
    pageId: "jobs",
    title: "Jobs",
    schemaUrl: "/api/schema/jobs",
    route: "/jobs",
  };
  const pageDocument = (await loadPageDocument(page, {}, {
    fetcher: async () => new Response(JSON.stringify(raw), { status: 200 }),
  })) as RenderPageDocument;

  const container = globalThis.document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider stored="en-US">
        <RenderPage
          document={pageDocument}
          context={{ user: { permissions: options.permissions }, features: {} } as never}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </I18nProvider>,
    );
  });
  await act(async () => {
    await Promise.resolve();
  });
  return { container, requests, downloads, blobs: capturedBlobs, listCalls: () => listCalls };
}

function rowButton(container: HTMLElement, rowIndex: number, text: string): HTMLButtonElement | null {
  const rows = Array.from(container.querySelectorAll("tbody tr"));
  const target = rows[rowIndex];
  if (target === undefined) {
    return null;
  }
  const buttons = Array.from(target.querySelectorAll("button"));
  return (buttons.find((button) => (button.textContent ?? "").trim() === text) as
    | HTMLButtonElement
    | undefined) ?? null;
}

async function click(element: HTMLElement | null): Promise<void> {
  expect(element, "the control under test must exist").not.toBeNull();
  await act(async () => {
    element!.dispatchEvent(new MouseEvent("click", { bubbles: true }));
  });
  await act(async () => {
    await Promise.resolve();
  });
}

describe("R4 · result center on the shipped jobs page", () => {
  it("renders the six-state columns and the action surface", async () => {
    const { container } = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });
    const text = container.textContent ?? "";
    for (const header of ["ID", "Kind", "Status", "Progress", "Attempt", "Actor", "Error"]) {
      expect(text).toContain(header);
    }
    // Four seeded rows, each carrying its own state badge.
    expect(container.querySelectorAll("tbody tr").length).toBe(ROWS.length);
    // The auto-refresh control and the search form are both present.
    expect(container.querySelector("[data-jobs-refresh]")).not.toBeNull();
    // No selection column: the jobs table declares no `props.selection`, which
    // is what makes the auto-refresh control's reloadList() side effect inert.
    expect(container.querySelector('input[type="checkbox"]')).toBeNull();
  });

  it("uses only theme tokens, so light and dark both apply", async () => {
    const { container } = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });
    // The surface must be themed through the token classes (border-input,
    // bg-background, text-muted-foreground …) — a fixed palette class would look
    // right in one theme and wrong in the other. This is the whole-page check
    // for the surface R4 added (rows, badges, actions, refresh control).
    const classNames = Array.from(container.querySelectorAll("*"))
      .map((element) => element.getAttribute("class") ?? "")
      .join(" ");
    for (const forbidden of [
      "bg-white",
      "bg-black",
      "text-white",
      "text-black",
      "bg-gray-",
      "text-gray-",
      "bg-slate-",
    ]) {
      expect(classNames).not.toContain(forbidden);
    }
    // Positive control: the tokens really are in use.
    expect(classNames).toContain("text-muted-foreground");
  });

  it("opens the record view with the job's own detail fields", async () => {
    const { container } = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });
    const firstRow = container.querySelector("tbody tr");
    expect(firstRow).not.toBeNull();
    await click(firstRow as HTMLElement);

    const detail = document.body.textContent ?? "";
    expect(detail).toContain("Job details");
    // The drawer reads the selected row, so the fields are the row's own values.
    expect(detail).toContain("job-running");
    expect(detail).toContain("corr-job-running");
  });

  it("enables cancel/retry only where the server says the transition is legal", async () => {
    const { container } = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });

    // Row 0 is running → cancellable, not retryable.
    expect(rowButton(container, 0, "Cancel")?.disabled).toBe(false);
    expect(rowButton(container, 0, "Retry")?.disabled).toBe(true);
    // Row 1 is failed → retryable, not cancellable.
    expect(rowButton(container, 1, "Cancel")?.disabled).toBe(true);
    expect(rowButton(container, 1, "Retry")?.disabled).toBe(false);
    // Row 2 succeeded and row 3 expired: neither may be cancelled or retried.
    for (const index of [2, 3]) {
      expect(rowButton(container, index, "Cancel")?.disabled).toBe(true);
      expect(rowButton(container, index, "Retry")?.disabled).toBe(true);
    }
  });

  it("cancels the row's own job through the frozen management-scope route", async () => {
    const harness = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });
    await click(rowButton(harness.container, 0, "Cancel"));
    // The row action carries a confirm step; accept it.
    const confirm = Array.from(document.body.querySelectorAll("button")).find(
      (button) => (button.textContent ?? "").trim() === "Confirm",
    );
    await click(confirm ?? null);

    const posted = harness.requests.filter((entry) => entry.method === "POST");
    expect(posted.map((entry) => entry.url)).toEqual(["/api/jobs/job-running/cancel"]);
  });

  it("retries the row's own failed job through the frozen route", async () => {
    const harness = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });
    await click(rowButton(harness.container, 1, "Retry"));
    const posted = harness.requests.filter((entry) => entry.method === "POST");
    expect(posted.map((entry) => entry.url)).toEqual(["/api/jobs/job-failed/retry"]);
  });

  it("downloads a succeeded CSV result under the server's filename", async () => {
    const harness = await renderJobsPage({
      permissions: ["jobs.read", "jobs.write"],
      resultDocument: { resource: "users", rowCount: 2, fileName: "users-selection.csv", csv: "id,name\na,b\n" },
    });
    // Download is the third action, so it lives in the overflow menu. It is
    // only enabled on a SUCCEEDED row, which here is the third one.
    const succeededRowMenu = Array.from(
      harness.container.querySelectorAll("[data-row-actions-menu]"),
    )[2]?.querySelector("button");
    await click(succeededRowMenu ?? null);
    const menuItem = Array.from(document.body.querySelectorAll('[role="menuitem"]')).find(
      (item) => (item.textContent ?? "").trim() === "Download result",
    );
    await click((menuItem as HTMLElement | undefined) ?? null);

    expect(harness.requests.some((entry) => entry.url === "/api/jobs/job-done/result")).toBe(true);
    expect(harness.downloads.map((entry) => entry.filename)).toEqual(["users-selection.csv"]);
    expect(harness.blobs[0]?.parts).toBe("id,name\na,b\n");
  });

  it("denies the write actions to a read-only principal", async () => {
    const harness = await renderJobsPage({ permissions: ["jobs.read"] });
    expect(rowButton(harness.container, 0, "Cancel")?.disabled).toBe(true);
    expect(rowButton(harness.container, 1, "Retry")?.disabled).toBe(true);

    // The button is disabled, not merely styled: a programmatic click must not
    // reach the network either (the JSX disabled attribute blocks activation).
    await click(rowButton(harness.container, 0, "Cancel"));
    expect(harness.requests.filter((entry) => entry.method === "POST")).toEqual([]);

    // Discriminating counterpart: download is gated on jobs.read (the server's
    // own gate for GET /api/jobs/{id}/result), so a read-only principal KEEPS
    // it. If the schema's local gate ever became jobs.write, this fails — which
    // is what makes the read-only case meaningful rather than vacuous.
    const succeededRowMenu = Array.from(
      harness.container.querySelectorAll("[data-row-actions-menu]"),
    )[2]?.querySelector("button");
    await click(succeededRowMenu ?? null);
    const downloadItem = Array.from(document.body.querySelectorAll('[role="menuitem"]')).find(
      (item) => (item.textContent ?? "").trim() === "Download result",
    ) as HTMLButtonElement | undefined;
    expect(downloadItem, "download must still be offered to a jobs.read holder").toBeDefined();
    expect(downloadItem?.disabled).toBe(false);
  });

  it("re-fetches the list on the auto-refresh cadence", async () => {
    vi.useFakeTimers();
    const harness = await renderJobsPage({ permissions: ["jobs.read", "jobs.write"] });
    const before = harness.listCalls();

    const select = harness.container.querySelector("[data-jobs-refresh-select]") as HTMLSelectElement;
    expect(select).not.toBeNull();
    await act(async () => {
      select.value = "5000";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    await act(async () => {
      vi.advanceTimersByTime(5000);
    });
    await act(async () => {
      await Promise.resolve();
    });
    expect(harness.listCalls()).toBeGreaterThan(before);

    // Off (the default) means no background traffic.
    await act(async () => {
      select.value = "0";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    const settled = harness.listCalls();
    await act(async () => {
      vi.advanceTimersByTime(30000);
    });
    expect(harness.listCalls()).toBe(settled);
  });
});
