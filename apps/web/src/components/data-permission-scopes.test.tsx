// @vitest-environment jsdom
//
// W14 F-02 (GOAL-016): the data-permission page scope-assignment editor loads
// users/policies, lets an admin pick a user, and saves scope assignments.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { DataPermissionScopes } from "@/components/data-permission-scopes";
import { I18nProvider } from "@/i18n/runtime";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
  vi.unstubAllGlobals();
});

const NODE = {
  type: "custom" as const,
  id: "data-permission-scopes",
  component: "data-permission-scopes",
  props: {},
};

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
}

function renderEditor() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider>
        <DataPermissionScopes node={NODE as never} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("DataPermissionScopes (W14 F-02)", () => {
  it("loads users and policies and saves scope assignments", async () => {
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (url.startsWith("/api/users")) {
        return jsonResponse({ items: [{ id: "usr-1", username: "alice", name: "Alice" }], total: 1, page: 1, pageSize: 100 });
      }
      if (url.startsWith("/api/data-permission/policies")) {
        return jsonResponse({ items: [{ resource: "orders", ownerColumn: "owner_id", defaultScope: "self", enabled: true }], total: 1, page: 1, pageSize: 1 });
      }
      if (url.startsWith("/api/data-permission/scopes")) {
        if (init?.method === "PATCH") {
          return jsonResponse({ userId: "usr-1", updated: 1 });
        }
        return jsonResponse({ userId: "usr-1", items: [{ resource: "orders", scopeType: "all" }] });
      }
      return jsonResponse({ items: [], total: 0, page: 1, pageSize: 10 });
    });
    vi.stubGlobal("fetch", fetchMock);

    const container = renderEditor();
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(container.textContent).toContain("Scope assignments");

    const userSelect = container.querySelector("#data-permission-user") as HTMLSelectElement;
    expect(userSelect).not.toBeNull();
    await act(async () => {
      userSelect.value = "usr-1";
      userSelect.dispatchEvent(new Event("change", { bubbles: true }));
      await new Promise((resolve) => setTimeout(resolve, 30));
    });

    const scopeSelect = container.querySelector("select[aria-label='Scope for orders']") as HTMLSelectElement;
    expect(scopeSelect).not.toBeNull();
    expect(scopeSelect.value).toBe("all");

    const save = Array.from(container.querySelectorAll("button")).find((button) => button.textContent?.includes("Save scopes"));
    expect(save).not.toBeUndefined();
    await act(async () => {
      save?.click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });

    const patchCall = fetchMock.mock.calls.find(([input, init]) => String(input) === "/api/data-permission/scopes" && (init as RequestInit | undefined)?.method === "PATCH");
    expect(patchCall).toBeDefined();
    expect(JSON.parse(String((patchCall?.[1] as RequestInit | undefined)?.body))).toEqual({
      userId: "usr-1",
      scopes: { orders: "all" },
    });
  });
});
