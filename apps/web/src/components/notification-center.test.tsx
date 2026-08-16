// @vitest-environment jsdom
//
// W13 T-06 (GOAL-014): the notifications page is an interactive custom
// component — row click expands the detail inline and marks the
// notification read; the deep link ?open=<id> expands a targeted one.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { NotificationCenter } from "@/components/notification-center";
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
  id: "notifications-center",
  component: "notification-center",
  props: { targetTable: "notifications-table" },
};

const TWO_ITEMS = {
  items: [
    { id: "ntf-1", event: "account.locked", title: "Account locked", body: "Repeated failures.", read: false, createdAt: "2026-08-14T00:00:00.000Z" },
    { id: "ntf-2", event: "account.unlocked", title: "Account unlocked", body: "Welcome back.", read: true, createdAt: "2026-08-15T00:00:00.000Z" },
  ],
  total: 2,
  page: 1,
  pageSize: 10,
};

function renderCenter(routeQuery: Record<string, string> = {}) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider>
        <NotificationCenter
          node={NODE as never}
          context={{ route: { params: {}, query: routeQuery } }}
        />
      </I18nProvider>,
    );
  });
  return container;
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
}

describe("NotificationCenter (W13 T-06)", () => {
  it("renders the list and expands + marks read on row click", async () => {
    const readCalls: string[] = [];
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (url.includes("/read") && init?.method === "POST") {
        readCalls.push(url);
        return new Response(null, { status: 204 });
      }
      return jsonResponse(TWO_ITEMS);
    });
    vi.stubGlobal("fetch", fetchMock);
    const container = renderCenter();
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(container.textContent).toContain("Account locked");

    // Click the unread row → detail expands + POST /read fires.
    const row = container.querySelector("[data-notification-row=ntf-1]") as HTMLButtonElement;
    expect(row).not.toBeNull();
    await act(async () => {
      row.click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    const rowAfter = container.querySelector("[data-notification-row=ntf-1]") as HTMLButtonElement;
    expect(rowAfter.getAttribute("aria-expanded")).toBe("true");
    expect(container.querySelector("[data-notification-detail=ntf-1]")).not.toBeNull();
    expect(container.textContent).toContain("Repeated failures.");
    expect(readCalls).toContain("/api/notifications/ntf-1/read");
  });

  it("does not re-mark an already-read row on click", async () => {
    const readCalls: string[] = [];
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (url.includes("/read") && init?.method === "POST") {
        readCalls.push(url);
        return new Response(null, { status: 204 });
      }
      return jsonResponse(TWO_ITEMS);
    });
    vi.stubGlobal("fetch", fetchMock);
    const container = renderCenter();
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    const row = container.querySelector("[data-notification-row=ntf-2]") as HTMLButtonElement;
    await act(async () => {
      row.click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(container.querySelector("[data-notification-detail=ntf-2]")).not.toBeNull();
    expect(readCalls).toEqual([]);
  });

  it("deep link ?open=<id> expands and marks the targeted notification", async () => {
    const readCalls: string[] = [];
    const fetchMock = vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
      const url = String(input);
      if (url.includes("/read") && init?.method === "POST") {
        readCalls.push(url);
        return new Response(null, { status: 204 });
      }
      return jsonResponse(TWO_ITEMS);
    });
    vi.stubGlobal("fetch", fetchMock);
    const container = renderCenter({ open: "ntf-1" });
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    const row = container.querySelector("[data-notification-row=ntf-1]") as HTMLButtonElement;
    expect(row.getAttribute("aria-expanded")).toBe("true");
    expect(container.querySelector("[data-notification-detail=ntf-1]")).not.toBeNull();
    expect(readCalls).toContain("/api/notifications/ntf-1/read");
  });
});
