// @vitest-environment jsdom
// Invitation-panel role multi-select coverage (workspace-019 UX polish):
// options load from /api/roles, the dropdown toggles selections, the create
// payload carries the selected role keys, and create stays disabled until at
// least one role is selected.
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { UserInvitesPanel } from "@/components/user-invites-panel";
import { I18nProvider } from "@/i18n/runtime";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.unstubAllGlobals();
});

function stubRoutes(posts: Array<Record<string, unknown>>, allRows: InviteRowFixture[] = []) {
  const fetchMock = vi.fn().mockImplementation((url: RequestInfo | URL, init?: RequestInit) => {
    const path = String(url);
    if ((init?.method ?? "GET") === "POST") {
      posts.push(JSON.parse(String(init?.body ?? "{}")) as Record<string, unknown>);
      return Promise.resolve({ ok: true, status: 201, json: async () => ({}) });
    }
    if (path.includes("/api/roles")) {
      return Promise.resolve({
        ok: true,
        json: async () => ({
          items: [
            { key: "admin", name: "Admin" },
            { key: "viewer", name: "Viewer" },
            { key: "editor", name: "Editor" },
          ],
        }),
      });
    }
    if (path.includes("/api/users/invites")) {
      const q = new URL(path, "http://x").searchParams;
      const status = q.get("status") ?? "all";
      const page = Number(q.get("page") ?? "1");
      const pageSize = Number(q.get("pageSize") ?? "10");
      const source =
        status === "consumed"
          ? [
              {
                id: "inv-1",
                roles: ["viewer"],
                status: "consumed",
                email: "u@example.com",
                expiresAt: "2026-09-01T00:00:00Z",
              },
            ]
          : allRows;
      const start = (page - 1) * pageSize;
      return Promise.resolve({
        ok: true,
        json: async () => ({
          items: source.slice(start, start + pageSize),
          total: source.length,
          page,
          pageSize,
        }),
      });
    }
    return Promise.resolve({ ok: true, json: async () => ({ items: [], total: 0 }) });
  });
  vi.stubGlobal("fetch", fetchMock);
  return fetchMock;
}

interface InviteRowFixture {
  id: string;
  roles: string[];
  status: string;
  email?: string;
  expiresAt?: string;
}

function seededRows(count: number): InviteRowFixture[] {
  return Array.from({ length: count }, (_, i) => ({
    id: `inv-${i + 1}`,
    roles: ["viewer"],
    status: "pending",
    email: `inv${i + 1}@example.com`,
    expiresAt: "2026-09-01T00:00:00Z",
  }));
}

async function renderPanel() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <UserInvitesPanel />
      </I18nProvider>,
    );
  });
  return container;
}

describe("UserInvitesPanel role multi-select", () => {
  it("loads the role catalog and shows the default selection as a badge", async () => {
    stubRoutes([]);
    const container = await renderPanel();
    // Default selection "viewer" renders as a badge on the trigger.
    expect(container.querySelector('[data-role-multiselect-badge="viewer"]')?.textContent).toContain("Viewer");
    // The popover is closed initially.
    expect(container.querySelector('[data-role-multiselect-popover]')).toBeNull();
  });

  it("opens the dropdown, toggles roles, and submits the selected keys", async () => {
    const posts: Array<Record<string, unknown>> = [];
    stubRoutes(posts);
    const container = await renderPanel();

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-role-multiselect-trigger]")!.click();
    });
    const option = container.querySelector<HTMLInputElement>('[data-role-multiselect-option="admin"]');
    expect(option).not.toBeNull();
    await act(async () => {
      option!.click();
    });
    expect(container.querySelector('[data-role-multiselect-badge="admin"]')?.textContent).toContain("Admin");

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-invite-create]")!.click();
    });
    expect(posts[0]).toEqual({
      email: undefined,
      roles: ["viewer", "admin"],
      expiresInDays: 7,
    });
  });

  it("keeps create disabled until at least one role is selected", async () => {
    stubRoutes([]);
    const container = await renderPanel();
    const create = container.querySelector<HTMLButtonElement>("[data-invite-create]")!;
    expect(create.disabled).toBe(false);

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-role-multiselect-trigger]")!.click();
    });
    await act(async () => {
      container.querySelector<HTMLInputElement>('[data-role-multiselect-option="viewer"]')!.click();
    });
    // Deselecting the only chosen role disables the create action.
    expect(container.querySelector<HTMLButtonElement>("[data-invite-create]")!.disabled).toBe(true);
  });

  it("filters the records via the server-side status parameter", async () => {
    const fetchMock = stubRoutes([]);
    const container = await renderPanel();
    // Two cards: issue + records are separate surfaces.
    expect(container.querySelector("[data-invite-issue-card]")).not.toBeNull();
    expect(container.querySelector("[data-invite-records-card]")).not.toBeNull();
    // Default: no status param, empty list.
    expect(container.querySelector("[data-invite-table]")).toBeNull();

    const select = container.querySelector<HTMLSelectElement>("[data-invite-status-filter]")!;
    await act(async () => {
      select.value = "consumed";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    await act(async () => {});
    const calledWithStatus = fetchMock.mock.calls.some(([url]) =>
      String(url).includes("/api/users/invites") && String(url).includes("status=consumed"),
    );
    expect(calledWithStatus).toBe(true);
    expect(container.querySelector("[data-invite-table]")).not.toBeNull();
    expect(container.querySelector('[data-invite-status]')?.textContent).toBe("consumed");
    expect(container.textContent).toContain("u@example.com");
  });

  it("paginates the records list server-side", async () => {
    const fetchMock = stubRoutes([], seededRows(25));
    const container = await renderPanel();
    // 25 rows at pageSize=10 → 3 pages; page 1 shows 10 rows.
    expect(container.querySelector("[data-invite-pagination]")).not.toBeNull();
    expect(container.querySelectorAll("[data-invite-table] tbody tr").length).toBe(10);
    const prev = container.querySelector<HTMLButtonElement>("[data-invite-prev-page]")!;
    const next = container.querySelector<HTMLButtonElement>("[data-invite-next-page]")!;
    expect(prev.disabled).toBe(true); // already on page 1

    await act(async () => {
      next.click();
    });
    await act(async () => {});
    const calledWithPageTwo = fetchMock.mock.calls.some(
      ([url]) => String(url).includes("/api/users/invites") && String(url).includes("page=2"),
    );
    expect(calledWithPageTwo).toBe(true);
    // Page 2 shows rows 11..20.
    expect(container.textContent).toContain("inv11@example.com");
    expect(container.textContent).not.toContain("inv1@example.com");
    expect(container.querySelector<HTMLButtonElement>("[data-invite-prev-page]")!.disabled).toBe(false);

    // pageSize switch resets to page 1 and narrows the slice.
    await act(async () => {
      const size = container.querySelector<HTMLSelectElement>("[data-invite-page-size]")!;
      size.value = "50";
      size.dispatchEvent(new Event("change", { bubbles: true }));
    });
    await act(async () => {});
    expect(container.querySelectorAll("[data-invite-table] tbody tr").length).toBe(25);
    expect(container.querySelector("[data-invite-pagination]")).toBeNull(); // single page now
  });
});