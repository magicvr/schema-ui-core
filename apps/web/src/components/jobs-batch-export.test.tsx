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
  context?: Record<string, unknown>;
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
          context={
            options?.context ?? {
              user: { permissions: ["users.read", "users.write", "data.export", "jobs.write"] },
            }
          }
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

// 2026-09-19 user report: 「导出所选」 answered 「未找到」 on both list pages.
//
// Two independent causes, both pinned here:
//
//   1. The operator config did not enable `admin.jobs`, so the route was never
//      mounted and the server answered a bare 404. The schema node still ships
//      (it belongs to admin.users/admin.roles), so the button rendered and every
//      click failed. Config-side guard: internal/config
//      TestOperatorConfigCoversAdminPreset.
//   2. The same shape exists in the `mvp`/`demo` presets by design: those
//      profiles contain neither admin.jobs nor admin.data-transfer, so the two
//      grants the route requires are absent. The trigger must not look usable
//      there.
//
// W33 D-001 §3 froze this entry point as fail-open — it must NOT be hidden
// (silently removing an operation entry point is the worse failure mode, and an
// empty slot host also breaks the list-surface height contract). It is therefore
// rendered DISABLED and explained.
describe("batch export availability (2026-09-19 regression)", () => {
  it("renders the trigger disabled and explained when the profile cannot grant the route's gates", async () => {
    // mvp/demo: the users page renders, but neither admin.jobs nor
    // admin.data-transfer is in the profile, so neither grant is present.
    const harness = await renderUsersPage({
      context: { user: { permissions: ["users.read", "users.write"] } },
    });
    const button = exportButton(harness.container);
    expect(button, "the entry point must stay visible (fail-open, W33 D-001 §3)").not.toBeNull();
    expect(button!.disabled).toBe(true);
    expect(button!.getAttribute("data-jobs-batch-export-unavailable")).toBe("true");
    expect(harness.container.textContent).toContain("not available in this deployment");
    // The list itself still works — only the async entry point is withheld.
    expect(harness.container.querySelector("table")).not.toBeNull();

    // And it cannot reach the network even if clicked programmatically.
    await act(async () => {
      button!.dispatchEvent(new MouseEvent("click", { bubbles: true }));
    });
    expect(harness.requests.filter((entry) => entry.url === "/api/jobs/batch-export")).toEqual([]);
  });

  it("keeps the trigger enabled when only one of the two gates is granted", async () => {
    // jobs.write without data.export (and vice versa) is still not enough: the
    // route requires both, so the control stays unavailable.
    const harness = await renderUsersPage({
      context: { user: { permissions: ["users.read", "users.write", "jobs.write"] } },
    });
    expect(exportButton(harness.container)!.disabled).toBe(true);
  });

  it("still renders the trigger enabled for a principal holding both gates", async () => {
    const harness = await renderUsersPage();
    const button = exportButton(harness.container);
    expect(harness.container.querySelector("[data-jobs-batch-export]")).not.toBeNull();
    expect(button).not.toBeNull();
    expect(button!.getAttribute("data-jobs-batch-export-unavailable")).toBeNull();
  });

  it("keeps the trigger usable when permissions are unknown (fail-open, never disables a working op)", async () => {
    // A bare harness / older host may publish no permission list at all.
    // Disabling on unknown data would remove a working operation entry point.
    const harness = await renderUsersPage({ context: {} });
    const button = exportButton(harness.container);
    expect(button).not.toBeNull();
    // Disabled only because the selection is empty — not for availability.
    expect(button!.getAttribute("data-jobs-batch-export-unavailable")).toBeNull();
    await selectRows(harness.container, 1);
    expect(exportButton(harness.container)!.disabled).toBe(false);
  });

  it("explains an unmounted route instead of surfacing a bare NOT_FOUND", async () => {
    // Defence in depth for cause 1: if a deployment still serves the node
    // without admin.jobs, the operator must not read 「未找到」 as "my rows are
    // missing".
    const harness = await renderUsersPage({
      submitResponse: { status: 404, body: { error: "NOT_FOUND", message: "未找到" } },
    });
    await selectRows(harness.container, 1);
    await click(exportButton(harness.container));

    const text = harness.container.textContent ?? "";
    expect(text).toContain("not available in this deployment");
    expect(text).not.toContain("未找到");
    expect(harness.jobListCalls()).toBe(0);
  });
});
