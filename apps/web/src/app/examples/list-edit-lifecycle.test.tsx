// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it } from "vitest";

import { ListEditLifecyclePage } from "@/app/examples/list-edit-lifecycle-page";
import type { NavigationContext } from "@/protocol/app-manifest";

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", {
    configurable: true,
    value: true,
  });
  // jsdom cannot resolve relative URLs through fetch; stub the records API used
  // by the page (one row so the Edit/Delete gate buttons render).
  globalThis.fetch = (async (input: RequestInfo | URL, _init?: RequestInit) => {
    const url = String(input);
    if (url.startsWith("/api/records/")) {
      return new Response(null, { status: 204 });
    }
    return new Response(
      JSON.stringify({
        items: [
          {
            id: "rec-1",
            name: "Acme Console",
            status: "active",
            owner: "alice",
            updatedAt: "2026-07-31T00:00:00Z",
          },
        ],
        total: 1,
        page: 1,
        pageSize: 10,
      }),
      { status: 200, headers: { "Content-Type": "application/json" } },
    );
  }) as typeof fetch;
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

async function renderPage(context?: NavigationContext): Promise<HTMLDivElement> {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(<ListEditLifecyclePage context={context} />);
  });
  return container;
}

function actionButton(container: HTMLDivElement, label: string): HTMLButtonElement | null {
  const buttons = Array.from(container.querySelectorAll("button"));
  return buttons.find((button) => button.textContent?.trim() === label) ?? null;
}

describe("list-edit-lifecycle permission gate", () => {
  it("enables edit/delete for an admin session context", async () => {
    const container = await renderPage({ user: { roles: ["admin"] }, features: {} });
    expect(actionButton(container, "Edit")?.disabled).toBe(false);
    expect(actionButton(container, "Delete")?.disabled).toBe(false);
  });

  it("denies edit/delete for a non-admin context (observable rejection)", async () => {
    const container = await renderPage({ user: { roles: ["viewer"] }, features: {} });
    expect(actionButton(container, "Edit")?.disabled).toBe(true);
    expect(actionButton(container, "Delete")?.disabled).toBe(true);
  });
});
