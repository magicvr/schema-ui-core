// @vitest-environment jsdom
//
// Full-path integration evidence for the R1 representative pages (GOAL-004;
// GOAL-011 S3 repoints the injected resource surface from the legacy demo to users/roles):
// uses an explicit admin-profile manifest test fixture and the real
// Go-embedded, module-owned page documents
// through the App's schema-driven default path, with the users/roles API surface
// injected. Asserts "�?Schema 即可出现页面" holds on the main path and that
// unknown / illegal inputs fail closed with observable errors.

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { App } from "@/app/App";
import { I18nProvider } from "@/i18n/runtime";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const MANIFEST_PATH = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../test-fixtures/app-manifest.admin-dogfood.json",
);
const CORE_FIXTURE_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/modules/dev/examples/schema",
);
const MODULE_FIXTURE_DIRS: Record<string, string> = {
  settings: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/modules/settings/schema",
  ),
  activity: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/modules/activity/schema",
  ),
  users: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/modules/users/schema",
  ),
  roles: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/modules/roles/schema",
  ),
};

const MIGRATED_PAGE_IDS = [
  "overview",
  "data-table",
  "search-form-table",
  "form-controls",
  "form-with-reactions",
  "form-with-upload",
  "admin-list-batch",
  "data-display",
  "users",
  "roles",
  "settings",
  "activity",
];

const USERS = {
  items: [
    {
      id: "usr-1",
      username: "alice",
      name: "Alice",
      roles: ["admin"],
      updatedAt: "2026-08-03T00:00:00.000Z",
    },
    {
      id: "usr-2",
      username: "bob",
      name: "Bob",
      roles: ["editor"],
      updatedAt: "2026-08-03T11:00:00.000Z",
    },
  ],
  total: 2,
  page: 1,
  pageSize: 10,
};

const ROLES = {
  items: [
    {
      id: "role-admin",
      key: "admin",
      name: "admin",
      system: true,
	  permissions: ["roles.assign", "roles.read", "roles.write", "users.read", "users.write"],
	  menuItems: ["menu-users", "menu-roles"],
	  assignedUsers: 1,
	  editable: false,
	  deletable: false,
      updatedAt: "2026-08-03T00:00:00.000Z",
    },
    {
      id: "role-viewer",
      key: "viewer",
      name: "viewer",
      system: true,
	  permissions: ["roles.read", "users.read"],
	  menuItems: [],
	  assignedUsers: 0,
	  editable: false,
	  deletable: false,
      updatedAt: "2026-08-03T11:00:00.000Z",
    },
  ],
  total: 2,
  page: 1,
  pageSize: 10,
};

function fixtureDocument(pageId: string): unknown {
  const directory = MODULE_FIXTURE_DIRS[pageId] ?? CORE_FIXTURE_DIR;
  return JSON.parse(readFileSync(resolve(directory, `${pageId}.json`), "utf8"));
}

function manifest(): AppManifest {
  return validateAppManifest(
    JSON.parse(readFileSync(MANIFEST_PATH, "utf8")),
  );
}

/** Serves /api/schema/* from the real fixtures and /api/users + /api/roles from demo envelopes. */
function combinedFetcher(
  fixtures: Record<string, unknown>,
  usersStatus = 200,
): typeof fetch {
  return (async (input: RequestInfo | URL) => {
    const pathname = new URL(String(input), "http://test.local").pathname;
    if (pathname.startsWith("/api/schema/")) {
      const pageId = pathname.slice("/api/schema/".length);
      const document = fixtures[pageId];
      if (document === undefined) {
        return new Response(JSON.stringify({ error: "SCHEMA_NOT_FOUND" }), { status: 404 });
      }
      return new Response(JSON.stringify(document), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/users")) {
      if (usersStatus !== 200) {
        return new Response(JSON.stringify({ error: "users down" }), { status: usersStatus });
      }
      return new Response(JSON.stringify(USERS), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname.startsWith("/api/roles")) {
      return new Response(JSON.stringify(ROLES), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      });
    }
    if (pathname === "/api/settings/default") {
      return new Response(
        JSON.stringify({
          id: "default",
          siteTitle: "Schema UI Core",
          logoUrl: "",
          defaultLocale: "auto",
          siteTimezone: "auto",
          defaultTheme: "auto",
          updatedAt: "2026-08-04T00:00:00.000Z",
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    if (pathname.startsWith("/api/settings") || pathname === "/api/branding") {
      return new Response(
        JSON.stringify(
          pathname === "/api/branding"
            ? { siteTitle: "Schema UI Core", logoUrl: "" }
            : {
                items: [
                  {
                    id: "default",
                    siteTitle: "Schema UI Core",
                    logoUrl: "",
                    updatedAt: "2026-08-04T00:00:00.000Z",
                  },
                ],
                total: 1,
                page: 1,
                pageSize: 10,
              },
        ),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    if (pathname.startsWith("/api/operations")) {
      return new Response(
        JSON.stringify({
          items: [
            {
              id: "op-1",
              event: "auth.login",
              actorId: "user-admin",
              actorName: "Admin",
              recordId: "",
              detail: "{}",
              createdAt: "2026-08-04T00:00:00.000Z",
            },
          ],
          total: 1,
          page: 1,
          pageSize: 10,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    }
    return new Response("not found", { status: 404 });
  }) as typeof fetch;
}

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderApp(
  path: string,
  context: Record<string, unknown>,
  fixtures: Record<string, unknown>,
  usersStatus = 200,
): Promise<HTMLDivElement> {
  window.history.replaceState({}, "", path);
  const fetcher = combinedFetcher(fixtures, usersStatus);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <App
          manifest={manifest()}
          navigationContext={context}
          schemaFetcher={fetcher}
          resourceFetcher={fetcher}
        />
      </I18nProvider>,
    );
  });
  return container;
}

function realFixtures(): Record<string, unknown> {
  const fixtures: Record<string, unknown> = {};
  for (const pageId of MIGRATED_PAGE_IDS) {
    fixtures[pageId] = fixtureDocument(pageId);
  }
  return fixtures;
}

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  window.history.replaceState({}, "", "/");
});

describe("representative pages through the admin manifest fixture (GOAL-004)", () => {
  it("renders a migrated list page with users via the default path", async () => {
    const container = await renderApp("/data-table", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Data table");
    expect(container.textContent).toContain("alice");
    expect(container.textContent).toContain("bob");
    expect(container.textContent).toContain("2 items · page 1 of 1");
  });

  it("renders the search + table page structure", async () => {
    const container = await renderApp("/search-form-table", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Search + table");
    expect(container.textContent).toContain("Search roles");
    expect(container.textContent).toContain("admin");
  });

  it("renders a migrated form page", async () => {
    const container = await renderApp("/form-controls", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Form controls");
    expect(container.textContent).toContain("Name (input)");
    expect(container.textContent).toContain("Kind (select)");
  });

  it("renders the CRUD lifecycle page (table + toolbar + row actions + recordView)", async () => {
    const admin = {
      user: { id: "u1", roles: ["admin"], permissions: ["users.read", "users.write"] },
    };
    const container = await renderApp("/users", admin, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Users");
    // S4 surface: create toolbar trigger, row edit/delete actions, recordView.
    expect(container.textContent).toContain("New user");
    expect(container.textContent).toContain("Edit");
    // W11 · U-05: secondary actions (Delete) live in the "⋯" overflow menu,
    // which is portaled to document.body.
    const moreTrigger = container.querySelector<HTMLButtonElement>(
      '[data-row-actions-menu] button[aria-label]',
    );
    expect(moreTrigger).not.toBeNull();
    await act(async () => moreTrigger!.click());
    expect(document.body.textContent).toContain("Delete");
    expect(container.textContent).toContain("alice");
    expect(container.textContent).toContain("bob");
    expect(container.textContent).toContain("Select a record to view details.");
  });

  it("renders the roles CRUD page from the admin manifest fixture", async () => {
    const admin = {
      user: { id: "u1", roles: ["admin"], permissions: ["roles.read", "roles.write"] },
    };
    const container = await renderApp("/roles", admin, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Roles");
    expect(container.textContent).toContain("New role");
    expect(container.textContent).toContain("Edit");
    expect(container.textContent).toContain("Delete");
    expect(container.textContent).toContain("admin");
    expect(container.textContent).toContain("viewer");
    expect(container.textContent).toContain("Select a record to view details.");
  });

  it("applies $context reactions on the migrated reactive form", async () => {
    const admin = await renderApp("/form-with-reactions", { user: { roles: ["admin"] } }, realFixtures());
    expect(admin.textContent).toContain("Approval (switch)");

    const viewer = await renderApp("/form-with-reactions", { user: { roles: ["viewer"] } }, realFixtures());
    expect(viewer.textContent).not.toContain("Approval (switch)");
  });

  it("fails closed when a representative page document has an unknown node type", async () => {
    const fixtures = realFixtures();
    fixtures["search-form-table"] = {
      meta: {
        pageId: "search-form-table",
        title: "Search + table",
        protocolVersion: "2.7",
        requiredCapabilities: ["app.manifest"],
      },
      body: { type: "slider", id: "x", props: {} },
    };
    const container = await renderApp("/search-form-table", {}, fixtures);
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "outside the registry renderer whitelist",
    );
  });

  it("fails closed when a batch toolbar trigger fires with an empty selection", async () => {
    const admin = {
      user: { id: "u1", roles: ["admin"], permissions: ["users.read", "users.write"] },
    };
    const fetchSpy = vi.fn();
    const base = combinedFetcher(realFixtures());
    const trackingFetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
      fetchSpy(String(input), init?.method ?? "GET");
      return base(input, init);
    }) as typeof fetch;

    window.history.replaceState({}, "", "/admin-list-batch");
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider>
          <App
            manifest={manifest()}
            navigationContext={admin}
            schemaFetcher={trackingFetcher}
            resourceFetcher={trackingFetcher}
          />
        </I18nProvider>,
      );
    });

    // requiresSelection keeps the button disabled with no rows selected; the
    // confirm never appears and no batch request is constructed (fail-closed).
    // S2: the toolbar label resolves through the en-US catalog ("Batch delete").
    const batchButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("Batch delete"),
    )!;
    expect(batchButton.disabled).toBe(true);
    await act(async () => {
      batchButton.click();
    });
    expect(container.textContent).not.toContain("Delete the selected users?");
    expect(
      fetchSpy.mock.calls.filter(([url, method]) => url === "/api/users/batch-delete" && method === "POST"),
    ).toHaveLength(0);
  });

  it("fails closed when the users data source is unreachable on a list page", async () => {
    const container = await renderApp("/data-table", {}, realFixtures(), 500);
    expect(container.textContent).toContain("resource fetch failed");
  });

  it("renders the data-display page (statCard + chart over /api/roles)", async () => {
    const container = await renderApp("/data-display", {}, realFixtures());
    expect(container.querySelector("h1")?.textContent).toContain("Data display");
    expect(container.textContent).toContain("Total roles");
    expect(container.textContent).toContain("admin: 1");
    expect(container.querySelector("svg[role='img']")).not.toBeNull();
  });

  it("uploads through the real UploadField + actionRef (ADR-0012)", async () => {
    const admin = {
      user: { id: "u1", roles: ["admin"], permissions: [] },
    };
    const uploadCalls: Array<{ url: string; method: string; body: unknown }> = [];
    const base = combinedFetcher(realFixtures());
    const trackingFetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      const method = init?.method ?? "GET";
      if (method === "POST" && url === "/api/upload") {
        // The real transport sends the actual file bytes in a multipart part.
        let fileName = "";
        if (init?.body instanceof FormData) {
          for (const [key, value] of init.body.entries()) {
            if (value instanceof File) {
              // jsdom File exposes name/size; bytes are asserted end-to-end by
              // the headless browser run against the real Go upload endpoint.
              fileName = `${key}=${value.name}:${value.size}`;
            }
          }
        }
        uploadCalls.push({ url, method, body: fileName });
        return new Response(
          JSON.stringify({ id: "file-abc", name: "contract.pdf", size: 9, url: "/api/files/file-abc" }),
          { status: 200, headers: { "Content-Type": "application/json" } },
        );
      }
      return base(input, init);
    }) as typeof fetch;

    window.history.replaceState({}, "", "/form-with-upload");
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider>
          <App
            manifest={manifest()}
            navigationContext={admin}
            schemaFetcher={trackingFetcher}
            resourceFetcher={trackingFetcher}
          />
        </I18nProvider>,
      );
    });

    expect(container.querySelector("h1")?.textContent).toContain("Form with upload");
    const fileInput = container.querySelector("input[type='file']") as HTMLInputElement;
    expect(fileInput).not.toBeNull();

    const file = new File(["pdf-bytes"], "contract.pdf", { type: "application/pdf" });
    Object.defineProperty(fileInput, "files", {
      configurable: true,
      value: [file],
    });
    await act(async () => {
      fileInput.dispatchEvent(new Event("change", { bubbles: true }));
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    // One multipart POST against the resolved actionRef url + field value commit.
    expect(uploadCalls).toHaveLength(1);
    expect(uploadCalls[0]!.url).toBe("/api/upload");
    expect(uploadCalls[0]!.method).toBe("POST");
    expect(uploadCalls[0]!.body).toBe("file=contract.pdf:9");
    // W9 follow-up: URL-shaped committed values render as an image preview
    // instead of the raw "Value: …" text.
    expect(container.querySelector<HTMLImageElement>("img[src='/api/files/file-abc']")).not.toBeNull();
    expect(container.textContent).not.toContain("Value: /api/files/file-abc");
  });

  it("runs the ADR-0022 batch flow end-to-end (select → confirm → request → reload clears selection)", async () => {
    const admin = {
      user: { id: "u1", roles: ["admin"], permissions: ["users.read", "users.write"] },
    };
    const batchCalls: Array<{ url: string; body: unknown }> = [];
    const base = combinedFetcher(realFixtures());
    const trackingFetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      const method = init?.method ?? "GET";
      if (method !== "GET" && url.startsWith("/api/users/batch-delete")) {
        batchCalls.push({
          url,
          body: init?.body === undefined ? null : JSON.parse(String(init.body)),
        });
        return new Response(JSON.stringify({ deleted: 2 }), {
          status: 200,
          headers: { "Content-Type": "application/json" },
        });
      }
      return base(input, init);
    }) as typeof fetch;

    window.history.replaceState({}, "", "/admin-list-batch");
    const container = document.createElement("div");
    document.body.appendChild(container);
    const root = createRoot(container);
    activeRoots.push({ root, container });
    await act(async () => {
      root.render(
        <I18nProvider>
          <App
            manifest={manifest()}
            navigationContext={admin}
            schemaFetcher={trackingFetcher}
            resourceFetcher={trackingFetcher}
          />
        </I18nProvider>,
      );
    });

    expect(container.querySelector("h1")?.textContent).toContain("List + batch");
    // S2: the toolbar label resolves through the en-US catalog.
    expect(container.textContent).toContain("Batch delete");

    // The batch button starts disabled (requiresSelection + empty selection).
    const rowCheckboxes = container.querySelectorAll("input[aria-label='Select row']");
    expect(rowCheckboxes.length).toBe(2);
    const batchButton = () =>
      [...container.querySelectorAll("button")].find((button) =>
        button.textContent?.includes("Batch delete"),
      )!;
    expect(batchButton().disabled).toBe(true);

    // Select both rows → button enabled → click → confirm dialog.
    await act(async () => {
      (rowCheckboxes[0] as HTMLInputElement).click();
      (rowCheckboxes[1] as HTMLInputElement).click();
    });
    expect(batchButton().disabled).toBe(false);
    await act(async () => {
      batchButton().click();
    });
    expect(container.textContent).toContain("Delete the selected users?");
    const confirmButton = [...container.querySelectorAll("button")].find((button) =>
      button.textContent?.includes("Confirm"),
    )!;
    await act(async () => {
      confirmButton.click();
      await new Promise((resolve) => setTimeout(resolve, 0));
    });

    // One logical POST with the normalized $selection.keys body.
    expect(batchCalls).toHaveLength(1);
    expect(batchCalls[0]!.url).toBe("/api/users/batch-delete");
    expect(batchCalls[0]!.body).toEqual({ ids: ["usr-1", "usr-2"] });

    // Success feedback + selection cleared (reload) → button disabled again.
    expect(container.textContent).toContain("Items deleted");
    expect(batchButton().disabled).toBe(true);
    expect(
      [...container.querySelectorAll("input[aria-label='Select row']")].every(
        (input) => (input as HTMLInputElement).checked === false,
      ),
    ).toBe(true);
  });
});

