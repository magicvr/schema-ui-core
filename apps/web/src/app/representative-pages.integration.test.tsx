// @vitest-environment jsdom
//
// Full-path integration evidence for the R1 representative pages (GOAL-004;
// GOAL-011 S3 repoints the injected resource surface from the legacy demo to users/roles):
// uses the real app manifest (`apps/web/public/.well-known/schema-ui/`) and the
// real Go-embedded page fixtures (core fixtures plus module-owned Settings and
// Activity schemas)
// through the App's schema-driven default path, with the users/roles API surface
// injected. Asserts "改 Schema 即可出现页面" holds on the main path and that
// unknown / illegal inputs fail closed with observable errors.

import { readFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { App } from "@/app/App";
import { validateAppManifest, type AppManifest } from "@/protocol/app-manifest";

const MANIFEST_PATH = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../public/.well-known/schema-ui/app-manifest.json",
);
const FIXTURE_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/internal/handler/fixtures/schema",
);
const MODULE_FIXTURE_DIRS: Record<string, string> = {
  settings: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/settings/schema",
  ),
  activity: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/activity/schema",
  ),
  users: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/users/schema",
  ),
  roles: resolve(
    dirname(fileURLToPath(import.meta.url)),
    "../../../api/internal/modules/roles/schema",
  ),
};

const MIGRATED_PAGE_IDS = [
  "overview",
  "data-table",
  "search-form-table",
  "form-controls",
  "form-with-reactions",
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
  const directory = MODULE_FIXTURE_DIRS[pageId] ?? FIXTURE_DIR;
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
      <App
        manifest={manifest()}
        navigationContext={context}
        schemaFetcher={fetcher}
        resourceFetcher={fetcher}
      />,
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

describe("representative pages through the real manifest (GOAL-004)", () => {
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
    expect(container.textContent).toContain("Delete");
    expect(container.textContent).toContain("alice");
    expect(container.textContent).toContain("bob");
    expect(container.textContent).toContain("Select a record to view details.");
  });

  it("renders the roles CRUD page from the real manifest and fixture", async () => {
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
      body: { type: "chart", id: "x", props: {} },
    };
    const container = await renderApp("/search-form-table", {}, fixtures);
    expect(container.querySelector('[role="alert"]')?.textContent).toContain(
      "outside the §5 renderer whitelist",
    );
  });

  it("fails closed when the users data source is unreachable on a list page", async () => {
    const container = await renderApp("/data-table", {}, realFixtures(), 500);
    expect(container.textContent).toContain("resource fetch failed");
  });
});
