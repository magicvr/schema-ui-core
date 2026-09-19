// @vitest-environment jsdom
//
// GOAL-004 R3 interaction contract, closed by GOAL-005 R4 C3 (the R3 audit left
// these as recommended findings because only the render-level denominator was
// covered).
//
// The R3 component has four behaviours that a render check cannot see, and each
// one is a real defect if it regresses:
//
//   1. an empty selection cannot submit (the button is disabled and the handler
//      bails before any request);
//   2. ONLY a 202 counts as success — a 200 with a job body, or an error, must
//      not hand the UI into the polling state;
//   3. it never calls `reloadList()`: any successful page reload clears every
//      table selection (ADR-0022 D2), which would delete the very selection the
//      operator is watching;
//   4. it polls the job to a terminal state and then downloads the CSV through
//      the shared result helper under the server's own filename.
//
// The document under test is the real users page (users.json), rendered through
// SchemaTable so the selection seam is the production one.

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
import "@/components/jobs-batch-export";

const MODULES = resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules");
const USERS_SCHEMA = resolve(MODULES, "users/schema/users.json");

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

const USERS = [
  { id: "usr-1", name: "Ada", email: "ada@example.com", roles: ["admin"], enabled: true },
  { id: "usr-2", name: "Grace", email: "grace@example.com", roles: ["editor"], enabled: true },
  { id: "usr-3", name: "Linus", email: "linus@example.com", roles: ["editor"], enabled: true },
];

class RecordingBlob extends Blob {
  readonly parts: string;
  constructor(parts: BlobPart[], options?: BlobPropertyBag) {
    super(parts, options);
    this.parts = parts.map((part) => String(part)).join("");
  }
}

interface Harness {
  container: HTMLDivElement;
  requests: Array<{ url: string; method: string; body: string }>;
  downloads: string[];
  blobs: RecordingBlob[];
  jobListCalls: () => number;
}

async function renderUsersPage(options?: {
  submitResponse?: { status: number; body: unknown };
}): Promise<Harness> {
  const requests: Harness["requests"] = [];
  const downloads: string[] = [];
  const blobs: RecordingBlob[] = [];
  let jobListCalls = 0;
  let polled = 0;

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

  const fetcher = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = String(input);
    const method = init?.method ?? "GET";
    const body = typeof init?.body === "string" ? init.body : "";
    requests.push({ url, method, body });

    if (url === "/api/jobs/batch-export") {
      const response = options?.submitResponse ?? {
        status: 202,
        body: { id: "job-1", status: "queued", progress: 0 },
      };
      return new Response(JSON.stringify(response.body), {
        status: response.status,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (url.startsWith("/api/jobs/") && url.endsWith("/result")) {
      return new Response(
        JSON.stringify({
          resource: "users",
          rowCount: 2,
          fileName: "users-selection.csv",
          csv: "id,name\nusr-1,Ada\n",
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    if (url.startsWith("/api/jobs/")) {
      jobListCalls += 1;
      polled += 1;
      // First poll still running, second poll settled — the loop must survive
      // the intermediate state and stop at the terminal one.
      const status = polled === 1 ? "running" : "succeeded";
      return new Response(
        JSON.stringify({ id: "job-1", status, progress: status === "running" ? 42 : 100 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    if (url.startsWith("/api/users")) {
      return new Response(
        JSON.stringify({ items: USERS, total: USERS.length, page: 1, pageSize: 20 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    return new Response(JSON.stringify({ items: [], total: 0, page: 1, pageSize: 20 }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;

  const raw = JSON.parse(readFileSync(USERS_SCHEMA, "utf8")) as unknown;
  const page: PageEntry = {
    pageId: "users",
    title: "Users",
    schemaUrl: "/api/schema/users",
    route: "/users",
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
          context={{ user: { permissions: ["users.read", "users.write", "data.export", "jobs.write"] } }}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </I18nProvider>,
    );
  });
  await act(async () => {
    await Promise.resolve();
  });
  return { container, requests, downloads, blobs, jobListCalls: () => jobListCalls };
}

function exportButton(container: HTMLElement): HTMLButtonElement | null {
  return container.querySelector("[data-jobs-batch-export-submit]");
}

async function selectRows(container: HTMLElement, count: number): Promise<void> {
  const checkboxes = Array.from(
    container.querySelectorAll('tbody input[type="checkbox"]'),
  ) as HTMLInputElement[];
  expect(checkboxes.length).toBeGreaterThanOrEqual(count);
  for (const checkbox of checkboxes.slice(0, count)) {
    await act(async () => {
      checkbox.click();
    });
  }
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

describe("R3 · async batch export interaction contract", () => {
  it("cannot submit an empty selection", async () => {
    const harness = await renderUsersPage();
    const button = exportButton(harness.container);
    expect(button).not.toBeNull();
    expect(button?.disabled).toBe(true);

    // Even a dispatched click must not reach the network: the handler guards on
    // the selection, so the disabled attribute is defence in depth rather than
    // the only barrier.
    await act(async () => {
      button?.dispatchEvent(new MouseEvent("click", { bubbles: true }));
    });
    expect(harness.requests.filter((entry) => entry.url === "/api/jobs/batch-export")).toEqual([]);
  });

  it("submits the selected keys as the target resource and never reloads the list", async () => {
    const harness = await renderUsersPage();
    await selectRows(harness.container, 2);

    expect(exportButton(harness.container)?.disabled).toBe(false);
    await click(exportButton(harness.container));

    const submitted = harness.requests.filter((entry) => entry.url === "/api/jobs/batch-export");
    expect(submitted.length).toBe(1);
    expect(JSON.parse(submitted[0].body)).toEqual({ resource: "users", ids: ["usr-1", "usr-2"] });

    // No reloadList(): the users list is not refetched, and the selection the
    // operator made is still published (the button still shows its count).
    expect(harness.requests.filter((entry) => entry.url.startsWith("/api/users")).length).toBe(1);
    expect(exportButton(harness.container)?.textContent).toContain("(2)");
  });

  it("accepts only 202 as a successful submission", async () => {
    const harness = await renderUsersPage({
      submitResponse: { status: 200, body: { id: "job-1", status: "queued" } },
    });
    await selectRows(harness.container, 1);
    await click(exportButton(harness.container));

    // A 200 is NOT a submit: the component surfaces an error and never polls.
    expect(harness.container.querySelector("[data-jobs-batch-export-progress]")).toBeNull();
    expect(harness.jobListCalls()).toBe(0);
    expect(harness.container.textContent).toContain("Could not start the export");
  });

  it("polls to a terminal state and downloads the CSV under the server's filename", async () => {
    vi.useFakeTimers();
    const harness = await renderUsersPage();
    await selectRows(harness.container, 2);
    await click(exportButton(harness.container));

    // First tick: still running → progress is shown from the live projection.
    await act(async () => {
      vi.advanceTimersByTime(1000);
    });
    await act(async () => {
      await Promise.resolve();
    });
    expect(harness.container.querySelector("[data-jobs-batch-export-progress]")?.textContent).toContain(
      "42%",
    );

    // Second tick: succeeded → the download affordance replaces the progress.
    await act(async () => {
      vi.advanceTimersByTime(1000);
    });
    await act(async () => {
      await Promise.resolve();
    });
    const downloadButton = harness.container.querySelector(
      "[data-jobs-batch-export-download]",
    ) as HTMLButtonElement | null;
    expect(downloadButton).not.toBeNull();
    expect(harness.container.querySelector("[data-jobs-batch-export-progress]")).toBeNull();

    await click(downloadButton);
    expect(harness.requests.some((entry) => entry.url === "/api/jobs/job-1/result")).toBe(true);
    expect(harness.downloads).toEqual(["users-selection.csv"]);
    expect(harness.blobs[0]?.parts).toBe("id,name\nusr-1,Ada\n");
  });
});
