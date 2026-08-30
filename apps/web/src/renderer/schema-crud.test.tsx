// @vitest-environment jsdom
//
// T-UI-01..10 acceptance evidence for GOAL-007 S4/S5 (I-007-003 v0.2.2 §6;
// GOAL-011 S3 repoints the CRUD driver from the legacy demo to the users resource):
// the `users` representative page drives list/search/detail/create/edit/delete
// through Schema nodes against an emulated users API, and the permission /
// error / confirm matrix holds. The renderer under test is the frozen §5 main
// path; every page behaviour below comes from the real fixture, not from source.
//
// The API emulator mirrors the frozen resource contract (GET list envelope,
// POST 201 / INVALID_CREATE_*, PATCH 200, DELETE 204, unified {error,message},
// users.read / users.write). Backend authority itself is proven by the Go
// users/roles suite; here we assert the UI renders consistently with it.

import { readFileSync, readdirSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import type { RenderPageDocument } from "@/renderer/render.types";
import { RenderPage } from "@/renderer/render.tsx";
import { SchemaTable } from "@/renderer/schema-table";
import type { ResourceItem } from "@/renderer/resource";

const CORE_FIXTURE_DIR = resolve(
  dirname(fileURLToPath(import.meta.url)),
  "../../../api/modules/dev/examples/schema",
);
// R4 C3.3: users/roles schema documents are module-owned.
const MODULE_FIXTURE_DIRS: Record<string, string> = {
  users: resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules/users/schema"),
  "users-invites": resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules/users/schema"),
  roles: resolve(dirname(fileURLToPath(import.meta.url)), "../../../api/modules/roles/schema"),
};
const WEB_SRC_DIR = resolve(dirname(fileURLToPath(import.meta.url)), "..");

function fixtureDocument(pageId: string): unknown {
  const directory = MODULE_FIXTURE_DIRS[pageId] ?? CORE_FIXTURE_DIR;
  return JSON.parse(readFileSync(resolve(directory, `${pageId}.json`), "utf8"));
}

// --- Users API emulator (resource contract) ---

interface ApiCall {
  method: string;
  url: string;
  body?: unknown;
}

function json(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function createUsersApi(
  initial: ResourceItem[],
  resourcePath = "/api/users",
  options: { denyWrites?: boolean; failPatch?: boolean } = {},
): { fetcher: typeof fetch; store: ResourceItem[]; calls: ApiCall[] } {
  const store: ResourceItem[] = initial.map((user) => ({ ...user }));
  const calls: ApiCall[] = [];
  let tick = 0;
  const now = (): string => {
    tick += 1;
    return `2026-08-02T10:00:00.${String(tick).padStart(3, "0")}Z`;
  };
  const deny = (): Response =>
    json({ error: "FORBIDDEN", message: "permission required: users.write" }, 403);

  const fetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const raw = String(input);
    const url = new URL(raw, "http://test.local");
    const method = (init?.method ?? "GET").toUpperCase();
    calls.push({
      method,
      url: raw,
      body: init?.body !== undefined ? JSON.parse(String(init.body)) : undefined,
    });
    // W11 · U-01: role options for the openRoles checkboxGroup (the schema
    // requests the full catalog via ?pageSize=100; list fetches carry the
    // standard page/pageSize query and keep routing to the resource below).
    if (url.pathname === "/api/roles" && url.searchParams.get("pageSize") === "100") {
      return json({
        items: [
          { key: "admin", name: "Administrator" },
          { key: "editor", name: "Editor" },
          { key: "viewer", name: "Viewer" },
        ],
      });
    }
    const match = new RegExp(`^${resourcePath}(?:/([^/]+))?$`).exec(url.pathname);
    if (!match) {
      return json({ error: "NOT_FOUND", message: "no such route" }, 404);
    }
    const id = match[1];

    if (method === "GET" && id === undefined) {
      let items = [...store];
      const q = url.searchParams.get("q");
      if (q !== null && q !== "") {
        items = items.filter((user) =>
          String(user.key ?? user.username ?? user.name ?? "")
            .toLowerCase()
            .includes(q.toLowerCase()),
        );
      }
      const sort = url.searchParams.get("sort");
      const order = url.searchParams.get("order");
      if (sort !== null) {
        items.sort((a, b) => {
          const left = String(a[sort as keyof ResourceItem] ?? "");
          const right = String(b[sort as keyof ResourceItem] ?? "");
          const cmp = left.localeCompare(right);
          return order === "desc" ? -cmp : cmp;
        });
      }
      const page = Number(url.searchParams.get("page") ?? "1");
      const pageSize = Number(url.searchParams.get("pageSize") ?? "10");
      return json({
        items: items.slice((page - 1) * pageSize, page * pageSize),
        total: items.length,
        page,
        pageSize,
      });
    }
    if (method === "POST" && id === undefined) {
      if (options.denyWrites) {
        return deny();
      }
      const body = init?.body !== undefined ? JSON.parse(String(init.body)) : {};
      for (const field of ["username", "name", "password"]) {
        const value = typeof body[field] === "string" ? String(body[field]).trim() : "";
        if (value === "") {
          return json({ error: "INVALID_CREATE_FIELD", message: `${field} is required` }, 400);
        }
      }
      const user: ResourceItem = {
        id: `usr-new-${store.length + 1}`,
        username: String(body.username).trim(),
        name: String(body.name).trim(),
        roles:
          typeof body.roles === "string"
            ? body.roles.split(",").map((role: string) => role.trim()).filter(Boolean)
            : Array.isArray(body.roles)
              ? body.roles
              : [],
        updatedAt: now(),
      };
      store.unshift(user);
      return json(user, 201);
    }
    if (id !== undefined && (method === "PATCH" || method === "DELETE")) {
      const index = store.findIndex((user) => user.id === id);
      if (index < 0) {
        return json({ error: "USER_NOT_FOUND", message: "no user with that id" }, 404);
      }
      if (method === "DELETE") {
        if (options.denyWrites) {
          return deny();
        }
        store.splice(index, 1);
        return new Response(null, { status: 204 });
      }
      if (options.denyWrites) {
        return deny();
      }
      if (options.failPatch) {
        return json({ error: "INVALID_PATCH_FIELD", message: "name is required" }, 400);
      }
      const body = init?.body !== undefined ? JSON.parse(String(init.body)) : {};
      const patch: Partial<ResourceItem> = {};
      if (body.name !== undefined) {
        const value = String(body.name).trim();
        if (value === "") {
          return json({ error: "INVALID_PATCH_FIELD", message: "name is required" }, 400);
        }
        patch.name = value;
      }
	  if (body.roles !== undefined) {
		patch.roles =
		  typeof body.roles === "string"
			? body.roles.split(",").map((role: string) => role.trim()).filter(Boolean)
			: body.roles;
	  }
      const updated = { ...store[index], ...patch, updatedAt: now() };
      store[index] = updated;
      return json(updated);
    }
    return json({ error: "NOT_FOUND", message: "no such route" }, 404);
  }) as typeof fetch;

  return { fetcher, store, calls };
}

function createRolesApi(
  initial: ResourceItem[],
): { fetcher: typeof fetch; store: ResourceItem[]; calls: ApiCall[] } {
  const store = initial.map((role) => ({ ...role }));
  const calls: ApiCall[] = [];
  let tick = 0;
  const now = (): string => {
    tick += 1;
    return `2026-08-03T12:00:00.${String(tick).padStart(3, "0")}Z`;
  };

  const fetcher: typeof fetch = (async (input: RequestInfo | URL, init?: RequestInit) => {
    const raw = String(input);
    const url = new URL(raw, "http://test.local");
    const method = (init?.method ?? "GET").toUpperCase();
    const body = init?.body !== undefined ? JSON.parse(String(init.body)) : undefined;
    calls.push({ method, url: raw, body });
    // W11 · U-01/U-02: RBAC catalogs served to the role form dynamic options.
    if (url.pathname === "/api/permissions") {
      return json({
        items: [
          { key: "users.read", label: "users.read" },
          { key: "users.write", label: "users.write" },
          { key: "roles.read", label: "roles.read" },
          { key: "roles.write", label: "roles.write" },
          { key: "roles.assign", label: "roles.assign" },
          { key: "settings.read", label: "settings.read" },
          { key: "settings.write", label: "settings.write" },
          { key: "operations.read", label: "operations.read" },
        ],
      });
    }
    if (url.pathname === "/api/menu-items") {
      return json({
        items: [
          { id: "menu-users", pageRef: "users", label: "Users" },
          { id: "menu-roles", pageRef: "roles", label: "Roles" },
          { id: "menu-settings", pageRef: "settings", label: "Settings" },
          { id: "menu-activity", pageRef: "activity", label: "Activity" },
        ],
      });
    }
    const match = /^\/api\/roles(?:\/([^/]+))?$/.exec(url.pathname);
    if (!match) {
      return json({ error: "NOT_FOUND", message: "no such route" }, 404);
    }
    const id = match[1];

    if (method === "GET" && id === undefined) {
      return json({ items: [...store], total: store.length, page: 1, pageSize: 10 });
    }
    if (method === "POST" && id === undefined) {
      const role: ResourceItem = {
        id: `role-${String(body?.key).trim()}`,
        key: String(body?.key).trim(),
        name: String(body?.name).trim(),
        system: false,
		permissions: Array.isArray(body?.permissions) ? body.permissions : [],
		menuItems: Array.isArray(body?.menuItems) ? body.menuItems : [],
		assignedUsers: 0,
		editable: true,
		deletable: true,
        updatedAt: now(),
      };
      store.unshift(role);
      return json(role, 201);
    }
    if (id !== undefined && (method === "PATCH" || method === "DELETE")) {
      const index = store.findIndex((role) => role.id === id);
      if (index < 0) {
        return json({ error: "ROLE_NOT_FOUND", message: "no role with that id" }, 404);
      }
      if (method === "DELETE") {
        store.splice(index, 1);
        return new Response(null, { status: 204 });
      }
      const updated = {
        ...store[index],
		name: body?.name === undefined ? store[index]?.name : String(body.name).trim(),
		permissions: Array.isArray(body?.permissions) ? body.permissions : store[index]?.permissions,
		menuItems: Array.isArray(body?.menuItems) ? body.menuItems : store[index]?.menuItems,
        updatedAt: now(),
      };
      store[index] = updated;
      return json(updated);
    }
    return json({ error: "NOT_FOUND", message: "no such route" }, 404);
  }) as typeof fetch;

  return { fetcher, store, calls };
}

// --- Render harness ---

const ADMIN = {
  user: {
    id: "u1",
    roles: ["admin"],
    permissions: ["users.read", "users.write", "roles.assign"],
  },
};
const VIEWER = {
  user: { id: "u2", roles: ["viewer"], permissions: ["users.read"] },
};
const USER_WRITER = {
  user: { id: "u4", roles: ["editor"], permissions: ["users.read", "users.write"] },
};
const ROLE_ADMIN = {
  user: {
    id: "u1",
    roles: ["admin"],
    permissions: ["roles.read", "roles.write", "roles.assign", "users.read", "users.write"],
  },
};
const ROLE_WRITER = {
  user: {
    id: "u3",
    roles: ["role-manager"],
    permissions: ["roles.read", "roles.write"],
  },
};

const USERS: ResourceItem[] = [
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
];

const ROLES: ResourceItem[] = [
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
];

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

async function renderCrud(
  fixture: unknown,
  context: Record<string, unknown>,
  fetcher: typeof fetch,
): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <RenderPage
        document={fixture as RenderPageDocument}
        context={context}
        tableRenderer={(node) => <SchemaTable node={node} fetcher={fetcher} />}
      />,
    );
  });
  return container;
}

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

function buttonByText(container: HTMLElement, text: string): HTMLButtonElement | undefined {
  return Array.from(container.querySelectorAll("button")).find(
    (button) => button.textContent?.trim() === text,
  ) as HTMLButtonElement | undefined;
}

function setFieldValue(el: Element, value: string): void {
  const proto =
    el instanceof HTMLSelectElement
      ? HTMLSelectElement.prototype
      : el instanceof HTMLTextAreaElement
        ? HTMLTextAreaElement.prototype
        : HTMLInputElement.prototype;
  Object.getOwnPropertyDescriptor(proto, "value")?.set?.call(el, value);
  el.dispatchEvent(new Event(el instanceof HTMLSelectElement ? "change" : "input", { bubbles: true }));
}

function fieldInput(container: HTMLElement, id: string): HTMLInputElement | null {
  return container.querySelector<HTMLInputElement>(`#field-${id}`);
}

function checkboxByLabel(container: HTMLElement, label: string): HTMLInputElement | null {
  const owner = Array.from(container.querySelectorAll("label")).find(
    (candidate) => candidate.textContent?.trim() === label,
  );
  return owner?.querySelector<HTMLInputElement>('input[type="checkbox"]') ?? null;
}

/**
 * W11 · U-05: finds a row-action button, opening the row's "⋯" overflow menu
 * when the action is not inline (users tables carry 8 actions; only 2 stay
 * visible). Desktop + mobile dual-end renders duplicate menus — the first
 * match wins.
 */
async function rowActionButton(container: HTMLElement, label: string): Promise<HTMLButtonElement> {
  const inline = Array.from(container.querySelectorAll<HTMLButtonElement>("button")).find(
    (button) => button.textContent?.trim() === label,
  );
  if (inline !== undefined) {
    return inline;
  }
  // The overflow menu is portaled to document.body (W11 · U-05 fix: the table's
  // overflow container would otherwise clip the absolutely-positioned menu).
  // The menu may already be open from a previous lookup in the same loop —
  // toggling the trigger again would close it.
  const alreadyOpen = Array.from(document.body.querySelectorAll<HTMLButtonElement>('[role="menuitem"]')).find(
    (button) => button.textContent?.trim() === label,
  );
  if (alreadyOpen !== undefined) {
    return alreadyOpen;
  }
  const trigger = container.querySelector<HTMLButtonElement>(
    '[data-row-actions-menu] button[aria-label]',
  );
  expect(trigger, "more-menu trigger for " + label).not.toBeNull();
  await act(async () => trigger!.click());
  const item = Array.from(document.body.querySelectorAll<HTMLButtonElement>('[role="menuitem"]')).find(
    (button) => button.textContent?.trim() === label,
  );
  expect(item, label + " in overflow menu").not.toBeUndefined();
  return item!;
}

// --- T-UI tests ---

describe("T-UI-01 · list loads and empty state", () => {
  it("loads the representative list from /api/users and shows the item count", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    expect(container.textContent).toContain("alice");
    expect(container.textContent).toContain("bob");
    expect(container.textContent).toContain("2 items · page 1 of 1");
  });

  it("shows the empty message when the data source has no rows", async () => {
    const api = createUsersApi([]);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    expect(container.textContent).toContain("No items yet.");
  });
});

describe("T-UI-02 · sort and pagination", () => {
  it("toggles a column sort header and refetches with sort/order", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    const nameHeader = Array.from(container.querySelectorAll("th button")).find((button) =>
      button.textContent?.includes("Name"),
    );
    expect(nameHeader).not.toBeUndefined();
    await act(async () => (nameHeader as HTMLButtonElement).click());
    expect((nameHeader as HTMLButtonElement).getAttribute("aria-sort")).toBe("ascending");
    const sortCalls = api.calls.filter(
      (call) => call.method === "GET" && call.url.includes("sort=name") && call.url.includes("order=asc"),
    );
    expect(sortCalls.length).toBeGreaterThan(0);
  });
});

describe("T-UI-03 · search form-to-query binding (search-form-table)", () => {
  it("binds the search form to the target table query", async () => {
    const api = createUsersApi(ROLES, "/api/roles");
    const container = await renderCrud(fixtureDocument("search-form-table"), ADMIN, api.fetcher);
    const input = fieldInput(container, "q");
    expect(input).not.toBeNull();
    await act(async () => setFieldValue(input as HTMLInputElement, "adm"));
    const searchButton = buttonByText(container, "Search");
    await act(async () => (searchButton as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.url.includes("q=adm"))).toBe(true);
    expect(container.textContent).toContain("admin");
    expect(container.textContent).not.toContain("viewer");
    expect(container.textContent).toContain("1 item · page 1 of 1");
  });
});

describe("T-UI-04 · create form (POST → 201)", () => {
  it("creates a user and shows it in the refreshed list", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New user") as HTMLButtonElement).click());
    await act(async () => setFieldValue(fieldInput(container, "username") as HTMLInputElement, "carol"));
    await act(async () => setFieldValue(fieldInput(container, "name") as HTMLInputElement, "Carol"));
    await act(async () => setFieldValue(fieldInput(container, "password") as HTMLInputElement, "pw123456"));
    await act(async () => (buttonByText(container, "Create user") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "POST" && call.url === "/api/users")).toBe(true);
    expect(container.textContent).toContain("Item created");
    expect(container.textContent).toContain("Carol");
  });

  it("surfaces a blank-field INVALID_CREATE_FIELD form error", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New user") as HTMLButtonElement).click());
    await act(async () => (buttonByText(container, "Create user") as HTMLButtonElement).click());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("INVALID_CREATE_FIELD");
  });
});

describe("T-UI-05 · edit form (PATCH → 200, prefilled from the row)", () => {
  it("prefills the edit form, patches the row, and refreshes updatedAt", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Edit") as HTMLButtonElement).click());
    const nameInput = fieldInput(container, "name") as HTMLInputElement;
    expect(nameInput.value).toBe("Alice");
    await act(async () => setFieldValue(nameInput, "Alice Rebrand"));
    await act(async () => (buttonByText(container, "Save changes") as HTMLButtonElement).click());
    const patchCall = api.calls.find((call) => call.method === "PATCH");
    expect(patchCall?.url).toBe("/api/users/usr-1");
    expect(container.textContent).toContain("Item updated");
    expect(container.textContent).toContain("Alice Rebrand");
    expect(api.store.find((user) => user.id === "usr-1")?.name).toBe("Alice Rebrand");
  });
});

describe("T-UI-06 · delete with confirmation", () => {
  it("confirms then deletes the row (DELETE 204) and removes it from the list", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (await rowActionButton(container, "Delete")).click());
    expect(container.textContent).toContain("Delete this user?");
    await act(async () => (buttonByText(container, "Confirm") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "DELETE" && call.url === "/api/users/usr-1")).toBe(true);
    expect(container.textContent).toContain("Item deleted");
    expect(container.textContent).not.toContain("alice");
  });

  it("cancelling issues no request", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (await rowActionButton(container, "Delete")).click());
    await act(async () => (buttonByText(container, "Cancel") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "DELETE")).toBe(false);
    expect(container.textContent).toContain("alice");
  });
});

describe("T-UI-07 · error envelope rendering", () => {
  it("renders a FORBIDDEN envelope when the backend denies a write", async () => {
    const api = createUsersApi(USERS, "/api/users", { denyWrites: true });
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New user") as HTMLButtonElement).click());
    await act(async () => setFieldValue(fieldInput(container, "username") as HTMLInputElement, "dave"));
    await act(async () => setFieldValue(fieldInput(container, "name") as HTMLInputElement, "Dave"));
    await act(async () => setFieldValue(fieldInput(container, "password") as HTMLInputElement, "pw123456"));
    await act(async () => (buttonByText(container, "Create user") as HTMLButtonElement).click());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("FORBIDDEN");
  });

  it("renders a USER_NOT_FOUND envelope on a stale edit", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Edit") as HTMLButtonElement).click());
    api.store.length = 0; // emulate a concurrent delete
    await act(async () => (buttonByText(container, "Save changes") as HTMLButtonElement).click());
    expect(container.querySelector('[role="alert"]')?.textContent).toContain("USER_NOT_FOUND");
  });
});

describe("T-UI-08 · permission matrix (admin vs viewer)", () => {
  it("admin sees enabled write affordances", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
    for (const label of ["New user", "Edit", "Roles", "Password", "Delete"]) {
      const button = await rowActionButton(container, label);
      expect(button, `${label} present`).not.toBeUndefined();
      expect(button.disabled, `${label} enabled for admin`).toBe(false);
    }
  });

  it("viewer/editor is read-only: write affordances are disabled", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), VIEWER, api.fetcher);
    for (const label of ["New user", "Edit", "Roles", "Password", "Delete"]) {
      const button = await rowActionButton(container, label);
      expect(button, `${label} present`).not.toBeUndefined();
      expect(button.disabled, `${label} disabled for viewer`).toBe(true);
    }
    // The read surface still works (users.read is held).
    expect(container.textContent).toContain("alice");
  });

  it("users.write without roles.assign cannot open the role assignment path", async () => {
    const api = createUsersApi(USERS);
    const container = await renderCrud(fixtureDocument("users"), USER_WRITER, api.fetcher);
    expect((buttonByText(container, "New user") as HTMLButtonElement).disabled).toBe(false);
    expect((buttonByText(container, "Edit") as HTMLButtonElement).disabled).toBe(false);
    expect((buttonByText(container, "Roles") as HTMLButtonElement).disabled).toBe(true);
    expect((await rowActionButton(container, "Password")).disabled).toBe(false);
  });
});

describe("T-UI-09 · backend authority is not replaced by frontend hiding", () => {
  it("the API denies a write even when the frontend hides the affordance", async () => {
    const viewerApi = createUsersApi(USERS, "/api/users", { denyWrites: true });
    const container = await renderCrud(fixtureDocument("users"), VIEWER, viewerApi.fetcher);
    // Frontend hides/disables the affordance for the viewer…
    const editButton = buttonByText(container, "Edit") as HTMLButtonElement;
    expect(editButton.disabled).toBe(true);
    // …and a direct write the UI never issues is still blocked by the backend.
    const response = await viewerApi.fetcher("/api/users", {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ username: "x", name: "X", password: "y12345" }),
    });
    expect(response.status).toBe(403);
    const envelope = (await response.json()) as { error: string };
    expect(envelope.error).toBe("FORBIDDEN");
  });
});

describe("T-UI-10 · dual-resource page changes are fixture-only", () => {
  it("drives roles create from the real fixture", async () => {
    const api = createRolesApi(ROLES);
    const container = await renderCrud(fixtureDocument("roles"), ROLE_ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "New role") as HTMLButtonElement).click());
    await act(async () => setFieldValue(fieldInput(container, "key") as HTMLInputElement, "auditor"));
    await act(async () => setFieldValue(fieldInput(container, "name") as HTMLInputElement, "Auditor"));
    // W11 · U-02: permission/menu options come from the RBAC catalogs.
    await act(async () => (checkboxByLabel(container, "users.read") as HTMLInputElement).click());
    await act(async () => (checkboxByLabel(container, "Users") as HTMLInputElement).click());
    await act(async () => (buttonByText(container, "Create role") as HTMLButtonElement).click());
    expect(api.calls.find((call) => call.method === "POST")).toEqual({
      method: "POST",
      url: "/api/roles",
      body: {
        key: "auditor",
        name: "Auditor",
        permissions: ["users.read"],
        menuItems: ["menu-users"],
      },
    });
    expect(container.textContent).toContain("Item created");
    expect(api.store[0]?.id).toBe("role-auditor");
  });

  it("drives roles update from the real fixture", async () => {
    const api = createRolesApi([
	  {
		id: "role-ops",
		key: "ops",
		name: "Operator",
		system: false,
		permissions: ["users.read"],
		menuItems: ["menu-users"],
		assignedUsers: 0,
		editable: true,
		deletable: true,
		updatedAt: "2026-08-03T00:00:00.000Z",
	  },
    ]);
    const container = await renderCrud(fixtureDocument("roles"), ROLE_ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Edit") as HTMLButtonElement).click());
    const nameInput = fieldInput(container, "name") as HTMLInputElement;
    expect(nameInput.value).toBe("Operator");
    await act(async () => setFieldValue(nameInput, "Operations"));
	// W11 · U-02: dynamic catalog options — toggle users.read off, roles.read on.
	await act(async () => (checkboxByLabel(container, "users.read") as HTMLInputElement).click());
	await act(async () => (checkboxByLabel(container, "roles.read") as HTMLInputElement).click());
    await act(async () => (buttonByText(container, "Save changes") as HTMLButtonElement).click());
    expect(api.calls.find((call) => call.method === "PATCH")).toEqual({
      method: "PATCH",
      url: "/api/roles/role-ops",
      body: { name: "Operations", permissions: ["roles.read"], menuItems: ["menu-users"] },
    });
    expect(container.textContent).toContain("Item updated");
  });

  it("drives roles delete from the real fixture", async () => {
    const api = createRolesApi([
	  {
		id: "role-ops",
		key: "ops",
		name: "Operator",
		system: false,
		permissions: [],
		menuItems: [],
		assignedUsers: 0,
		editable: true,
		deletable: true,
		updatedAt: "2026-08-03T00:00:00.000Z",
	  },
    ]);
    const container = await renderCrud(fixtureDocument("roles"), ROLE_ADMIN, api.fetcher);
    await act(async () => (buttonByText(container, "Delete") as HTMLButtonElement).click());
    expect(container.textContent).toContain("Delete this role?");
    await act(async () => (buttonByText(container, "Confirm") as HTMLButtonElement).click());
    expect(api.calls.some((call) => call.method === "DELETE" && call.url === "/api/roles/role-ops")).toBe(true);
    expect(api.store).toHaveLength(0);
  });

  it("disables system-role and in-use delete actions from structured row flags", async () => {
	const api = createRolesApi(ROLES);
	const container = await renderCrud(fixtureDocument("roles"), ROLE_ADMIN, api.fetcher);
	const edits = Array.from(container.querySelectorAll("button")).filter(
	  (button) => button.textContent?.trim() === "Edit",
	);
	const deletes = Array.from(container.querySelectorAll("button")).filter(
	  (button) => button.textContent?.trim() === "Delete",
	);
	// Dual-end table (desktop + mobile cards) renders each row action twice.
	expect(edits.length).toBeGreaterThanOrEqual(2);
	expect(deletes.length).toBeGreaterThanOrEqual(2);
	expect(edits.every((button) => button.disabled)).toBe(true);
	expect(deletes.every((button) => button.disabled)).toBe(true);
  });

  it("disables grant create and edit paths for non-admin role writers", async () => {
	const api = createRolesApi([
	  {
		id: "role-ops",
		key: "ops",
		name: "Operator",
		system: false,
		permissions: ["roles.read"],
		menuItems: [],
		assignedUsers: 0,
		editable: true,
		deletable: true,
		updatedAt: "2026-08-03T00:00:00.000Z",
	  },
	]);
	const container = await renderCrud(fixtureDocument("roles"), ROLE_WRITER, api.fetcher);
	expect((buttonByText(container, "New role") as HTMLButtonElement).disabled).toBe(true);
	expect((buttonByText(container, "Edit") as HTMLButtonElement).disabled).toBe(true);
  });

  it("drives user role assignment and password change from distinct real fixture actions", async () => {
	const api = createUsersApi(USERS);
	const container = await renderCrud(fixtureDocument("users"), ADMIN, api.fetcher);
	await act(async () => (buttonByText(container, "Roles") as HTMLButtonElement).click());
	// W11 · U-01: roles are a checkboxGroup with dynamic options; the row's
	// roles[] prefill checks "Administrator" (admin).
	const adminBox = checkboxByLabel(container, "Administrator");
	expect(adminBox).not.toBeNull();
	expect(adminBox!.checked).toBe(true);
	await act(async () => adminBox!.click());
	await act(async () => (checkboxByLabel(container, "Editor") as HTMLInputElement).click());
	await act(async () => (checkboxByLabel(container, "Viewer") as HTMLInputElement).click());
	await act(async () => (buttonByText(container, "Save roles") as HTMLButtonElement).click());
	expect(api.calls.find((call) => call.method === "PATCH")?.body).toEqual({
	  roles: ["editor", "viewer"],
	});

	await act(async () => (await rowActionButton(container, "Password")).click());
	const password = fieldInput(container, "password") as HTMLInputElement;
	expect(password.type).toBe("password");
	await act(async () => setFieldValue(password, "  exact-password  "));
	await act(async () => (buttonByText(container, "Change password") as HTMLButtonElement).click());
	const patchCalls = api.calls.filter((call) => call.method === "PATCH");
	expect(patchCalls.at(-1)?.body).toEqual({ password: "  exact-password  " });
  });

  it("the user and role CRUD action ids exist only in fixtures, never in renderer source", async () => {
    const fixtureTexts = ["users", "roles"].map((id) =>
      readFileSync(resolve(MODULE_FIXTURE_DIRS[id] ?? CORE_FIXTURE_DIR, `${id}.json`), "utf8"),
    );
    const actionIds = [
      "createUser",
      "updateUser",
	  "updateUserRoles",
	  "changeUserPassword",
      "deleteUser",
      "createRole",
      "updateRole",
      "deleteRole",
      "openCreate",
      "openEdit",
	  "openRoles",
	  "openPassword",
    ];
    for (const id of actionIds) {
      expect(fixtureTexts.some((fixtureText) => fixtureText.includes(id)), `${id} declared by a fixture`).toBe(true);
    }
    const srcFiles = readDirRecursive(WEB_SRC_DIR)
      .filter(
        (file) =>
          (file.endsWith(".ts") || file.endsWith(".tsx")) && !file.includes(".test."),
      );
    for (const file of srcFiles) {
      const content = readFileSync(file, "utf8");
      for (const id of actionIds) {
        expect(content, `${id} must not be hardcoded in ${file}`).not.toContain(id);
      }
    }
  });
});

function readDirRecursive(dir: string): string[] {
  const out: string[] = [];
  for (const entry of readdirSync(dir, { withFileTypes: true })) {
    const full = resolve(dir, entry.name);
    if (entry.isDirectory()) {
      out.push(...readDirRecursive(full));
    } else {
      out.push(full);
    }
  }
  return out;
}
