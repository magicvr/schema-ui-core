// @vitest-environment jsdom
// Invite issue card coverage (workspace-019 UX): role multi-select, create
// payload with role keys, and the resend flow landing the rotated link in
// the disclosure band (the schema table cannot surface response bodies).
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { InviteIssueCard } from "@/components/invite-issue-card";
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

function stubRoutes(posts: Array<Record<string, unknown>>, resendLinks: string[] = []) {
  let resendIndex = 0;
  const fetchMock = vi.fn().mockImplementation((url: RequestInfo | URL, init?: RequestInit) => {
    const path = String(url);
    if (path.includes("/api/users/invites/") && (init?.method ?? "GET") === "POST") {
      const link = resendLinks[Math.min(resendIndex, resendLinks.length - 1)] ?? "/invite/accept?token=rotated-1";
      resendIndex += 1;
      return Promise.resolve({ ok: true, status: 200, json: async () => ({ link }) });
    }
    if ((init?.method ?? "GET") === "POST") {
      posts.push(JSON.parse(String(init?.body ?? "{}")) as Record<string, unknown>);
      return Promise.resolve({ ok: true, status: 201, json: async () => ({ link: "/invite/accept?token=new-1" }) });
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

async function renderCard() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <InviteIssueCard />
      </I18nProvider>,
    );
  });
  return container;
}

describe("InviteIssueCard", () => {
  it("loads the role catalog and shows the default selection as a badge", async () => {
    stubRoutes([]);
    const container = await renderCard();
    expect(container.querySelector('[data-role-multiselect-badge="viewer"]')?.textContent).toContain("Viewer");
    expect(container.querySelector('[data-role-multiselect-popover]')).toBeNull();
  });

  it("opens the dropdown, toggles roles, and submits the selected keys", async () => {
    const posts: Array<Record<string, unknown>> = [];
    stubRoutes(posts);
    const container = await renderCard();

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-role-multiselect-trigger]")!.click();
    });
    await act(async () => {
      container.querySelector<HTMLInputElement>('[data-role-multiselect-option="admin"]')!.click();
    });
    expect(container.querySelector('[data-role-multiselect-badge="admin"]')?.textContent).toContain("Admin");

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-invite-create]")!.click();
    });
    expect(posts[0]).toEqual({ email: undefined, roles: ["viewer", "admin"], expiresInDays: 7 });
    // The one-time link from the create response is disclosed.
    expect(container.querySelector("[data-invite-link]")?.textContent).toContain("/invite/accept?token=new-1");
  });

  it("keeps create disabled until at least one role is selected", async () => {
    stubRoutes([]);
    const container = await renderCard();
    const create = container.querySelector<HTMLButtonElement>("[data-invite-create]")!;
    expect(create.disabled).toBe(false);

    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-role-multiselect-trigger]")!.click();
    });
    await act(async () => {
      container.querySelector<HTMLInputElement>('[data-role-multiselect-option="viewer"]')!.click();
    });
    expect(container.querySelector<HTMLButtonElement>("[data-invite-create]")!.disabled).toBe(true);
  });

  it("resends by invite id and discloses the rotated link", async () => {
    stubRoutes([], ["/invite/accept?token=rotated-1"]);
    const container = await renderCard();
    const idInput = container.querySelector<HTMLInputElement>("[data-invite-resend-id]")!;
    await act(async () => {
      const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
      setter?.call(idInput, "inv-42");
      idInput.dispatchEvent(new Event("input", { bubbles: true }));
    });
    await act(async () => {
      container.querySelector<HTMLButtonElement>("[data-invite-resend]")!.click();
    });
    expect(container.querySelector("[data-invite-link]")?.textContent).toContain("/invite/accept?token=rotated-1");
  });
});