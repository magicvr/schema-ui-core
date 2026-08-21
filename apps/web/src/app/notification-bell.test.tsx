// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { NotificationBell } from "@/app/notification-bell";
import { notifyConfigChanged } from "@/app/config-events";
import { I18nProvider } from "@/i18n/runtime";

/** F-04 bell tests (GOAL-006 S3) + W13 T-06 interactivity (GOAL-014). */

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
});

async function renderBell(fetcher: typeof fetch, onOpenItem: (id: string) => void = () => undefined) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <NotificationBell fetcher={fetcher} onViewAll={() => undefined} onOpenItem={onOpenItem} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("NotificationBell (F-04 / W13 T-06)", () => {
  it("shows the unread badge from the count endpoint", async () => {
    const fetcher = vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("unread-count")) {
        return new Response(JSON.stringify({ unread: 3 }), { status: 200, headers: { "Content-Type": "application/json" } });
      }
      return new Response(JSON.stringify({ items: [] }), { status: 200, headers: { "Content-Type": "application/json" } });
    });
    const container = await renderBell(fetcher as typeof fetch);
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(container.textContent).toContain("3");
    expect(fetcher).toHaveBeenCalledWith("/api/notifications/unread-count");
  });

  it("opens the dropdown with the latest notifications", async () => {
    const fetcher = vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("unread-count")) {
        return new Response(JSON.stringify({ unread: 1 }), { status: 200, headers: { "Content-Type": "application/json" } });
      }
      return new Response(
        JSON.stringify({
          items: [
            { id: "ntf-1", event: "account.locked", title: "Account locked", body: "Repeated failures.", read: false, createdAt: "2026-08-14T00:00:00.000Z" },
          ],
          total: 1,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    });
    const container = await renderBell(fetcher as typeof fetch);
    const button = container.querySelector("button");
    if (!(button instanceof HTMLButtonElement)) throw new Error("bell button not found");
    await act(async () => {
      button.click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(container.textContent).toContain("Account locked");
    expect(fetcher).toHaveBeenCalledWith("/api/notifications?pageSize=5");
  });

  it("fails open when the count endpoint errors", async () => {
    const fetcher = vi.fn(async () => {
      return new Response(JSON.stringify({ error: "INTERNAL", message: "down" }), { status: 500 });
    });
    const container = await renderBell(fetcher as typeof fetch);
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    // No badge; the shell surface stays intact (button still present).
    const button = container.querySelector("button");
    if (!(button instanceof HTMLButtonElement)) throw new Error("bell button not found");
    expect(container.textContent).not.toContain("99+");
  });

  it("opens a notification from the dropdown: marks read + navigates (W13 T-06)", async () => {
    const opened: string[] = [];
    const fetcher = vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("unread-count")) {
        return new Response(JSON.stringify({ unread: 1 }), { status: 200, headers: { "Content-Type": "application/json" } });
      }
      if (url.includes("/read")) {
        return new Response(null, { status: 204 });
      }
      return new Response(
        JSON.stringify({
          items: [
            { id: "ntf-1", event: "account.locked", title: "Account locked", body: "Repeated failures.", read: false, createdAt: "2026-08-14T00:00:00.000Z" },
          ],
          total: 1,
        }),
        { status: 200, headers: { "Content-Type": "application/json" } },
      );
    });
    const container = await renderBell(fetcher as typeof fetch, (id) => opened.push(id));
    const button = container.querySelector("button");
    if (!(button instanceof HTMLButtonElement)) throw new Error("bell button not found");
    await act(async () => {
      button.click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    const item = [...container.querySelectorAll("[role=menuitem]")].find((el) =>
      el.textContent?.includes("Account locked"),
    );
    expect(item).not.toBeUndefined();
    await act(async () => {
      (item as HTMLButtonElement).click();
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(opened).toEqual(["ntf-1"]);
    expect(fetcher).toHaveBeenCalledWith("/api/notifications/ntf-1/read", { method: "POST" });
  });

  it("refreshes the badge when a notification is read elsewhere (W13 T-06)", async () => {
    const fetcher = vi.fn(async (input: RequestInfo | URL) => {
      const url = String(input);
      if (url.includes("unread-count")) {
        return new Response(JSON.stringify({ unread: 3 }), { status: 200, headers: { "Content-Type": "application/json" } });
      }
      return new Response(JSON.stringify({ items: [] }), { status: 200, headers: { "Content-Type": "application/json" } });
    });
    const container = await renderBell(fetcher as typeof fetch);
    await act(async () => {
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(container.textContent).toContain("3");
    fetcher.mockClear();
    // The list page marked a notification read → config-change event.
    await act(async () => {
      notifyConfigChanged("notifications.read");
      await new Promise((resolve) => setTimeout(resolve, 30));
    });
    expect(fetcher).toHaveBeenCalledWith("/api/notifications/unread-count");
  });
});
