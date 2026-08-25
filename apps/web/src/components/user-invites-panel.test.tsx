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

function stubRoutes(posts: Array<Record<string, unknown>>) {
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
    return Promise.resolve({ ok: true, json: async () => ({ items: [], total: 0 }) });
  });
  vi.stubGlobal("fetch", fetchMock);
  return fetchMock;
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
});