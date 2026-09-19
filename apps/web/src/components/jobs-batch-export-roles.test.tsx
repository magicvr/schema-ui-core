// @vitest-environment jsdom
//
// GOAL-045 W33 · the roles page trigger (user-reported gap, 2026-09-19).
//
// The R3 export denominator was always `users` / `roles` (GOAL-004 D-001 §1) and
// the backend has served both from the start, but only the users page declared a
// selection and the async trigger — so the roles page silently lacked the
// operation. This test pins the completed surface ON THE REAL roles document:
// a multiple-selection table, the trigger placed in the list page-actions slot,
// the submission naming `roles`, and the no-reload contract.

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
const ROLES_SCHEMA = resolve(MODULES, "roles/schema/roles.json");

const ROLES = [
  { id: "role-admin", key: "admin", name: "Admin", system: true, permissions: [], menuItems: [], assignedUsers: 1 },
  { id: "role-editor", key: "editor", name: "Editor", system: false, permissions: [], menuItems: [], assignedUsers: 2 },
  { id: "role-viewer", key: "viewer", name: "Viewer", system: false, permissions: [], menuItems: [], assignedUsers: 3 },
];

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  vi.restoreAllMocks();
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

interface RolesHarness {
  container: HTMLDivElement;
  requests: Array<{ url: string; method: string; body: string }>;
}

async function renderRolesPage(): Promise<RolesHarness> {
  const requests: RolesHarness["requests"] = [];
  const fetcher = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const url = String(input);
    const method = init?.method ?? "GET";
    requests.push({ url, method, body: typeof init?.body === "string" ? init.body : "" });
    if (url === "/api/jobs/batch-export") {
      return new Response(JSON.stringify({ id: "job-roles", status: "queued", progress: 0 }), {
        status: 202,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (url.startsWith("/api/jobs/")) {
      return new Response(JSON.stringify({ id: "job-roles", status: "succeeded", progress: 100 }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (url.startsWith("/api/roles")) {
      return new Response(
        JSON.stringify({ items: ROLES, total: ROLES.length, page: 1, pageSize: 20 }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    return new Response(JSON.stringify({ items: [], total: 0, page: 1, pageSize: 20 }), {
      status: 200,
      headers: { "Content-Type": "application/json" },
    });
  }) as typeof fetch;

  const raw = JSON.parse(readFileSync(ROLES_SCHEMA, "utf8")) as unknown;
  const page: PageEntry = {
    pageId: "roles",
    title: "Roles",
    schemaUrl: "/api/schema/roles",
    route: "/roles",
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
          context={{ user: { permissions: ["roles.read", "roles.write", "data.export", "jobs.write"] } }}
          dataFetcher={fetcher}
          tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
        />
      </I18nProvider>,
    );
  });
  await act(async () => {
    await Promise.resolve();
  });
  return { container, requests };
}

describe("W33 · roles page async export trigger", () => {
  it("places the trigger in the list page-actions slot, not above the filters", async () => {
    const harness = await renderRolesPage();
    const left = harness.container.querySelector("[data-list-page-actions-left]");
    expect(left, "the roles table must expose the slot segment").not.toBeNull();
    expect(left!.querySelector("[data-jobs-batch-export-submit]")).not.toBeNull();
  });

  it("requires a selection and submits the roles resource", async () => {
    const harness = await renderRolesPage();
    const submit = harness.container.querySelector<HTMLButtonElement>("[data-jobs-batch-export-submit]");
    expect(submit).not.toBeNull();
    expect(submit!.disabled).toBe(true);

    const checkboxes = Array.from(
      harness.container.querySelectorAll('tbody input[type="checkbox"]'),
    ) as HTMLInputElement[];
    expect(checkboxes.length, "the roles table must declare multiple selection").toBe(ROLES.length);
    for (const checkbox of checkboxes.slice(0, 2)) {
      await act(async () => {
        checkbox.click();
      });
    }

    const enabled = harness.container.querySelector<HTMLButtonElement>("[data-jobs-batch-export-submit]");
    expect(enabled!.disabled).toBe(false);
    await act(async () => {
      enabled!.dispatchEvent(new MouseEvent("click", { bubbles: true }));
    });
    await act(async () => {
      await Promise.resolve();
    });

    const submitted = harness.requests.filter((entry) => entry.url === "/api/jobs/batch-export");
    expect(submitted.length).toBe(1);
    expect(JSON.parse(submitted[0].body)).toEqual({
      resource: "roles",
      ids: ["role-admin", "role-editor"],
    });
    // No reloadList(): the list is not refetched and the selection survives.
    expect(harness.requests.filter((entry) => entry.url.startsWith("/api/roles")).length).toBe(1);
    expect(
      harness.container.querySelector("[data-jobs-batch-export-submit]")?.textContent,
    ).toContain("(2)");
  });
});
