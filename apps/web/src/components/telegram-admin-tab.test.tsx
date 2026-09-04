// @vitest-environment jsdom
//
// GOAL-006 R5 (判据 #5 补做 Admin UI tab, user-adjudicated 2026-09-03): the
// Telegram channel settings console — write-only token/secret edit (empty keeps
// current), status booleans, and the mock-captured counter. Secrets are never
// returned by GET (token_set/secret_set only), so the component must send PATCH
// with only non-empty values.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { TelegramAdminTab } from "@/components/telegram-admin-tab";
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
  id: "telegram-admin-tab",
  component: "telegram-admin-tab",
  props: {},
};

const STATUS = {
  configured: true,
  token_set: true,
  secret_set: true,
  mode: "polling",
  webhook_public_base_url: "",
  connection_state: "running",
  receiver: "polling",
  business_occupied: false,
  bot_username: "fixture_bot",
  captured_messages_count: 3,
};

const OPERATOR_STATUS = {
  ...STATUS,
  mode: "webhook",
  receiver: "webhook",
  bot_id: 42,
};

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
}

function renderTab(fetchMock: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>) {
  vi.stubGlobal("fetch", vi.fn(fetchMock));
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <I18nProvider>
        <TelegramAdminTab node={NODE as never} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

async function settle() {
  await act(async () => {
    await new Promise((resolve) => setTimeout(resolve, 30));
  });
}

async function flushPromises() {
  await act(async () => {
    await Promise.resolve();
    await Promise.resolve();
  });
}

const setNativeValue = (el: HTMLInputElement, value: string) => {
  Object.getOwnPropertyDescriptor(HTMLInputElement.prototype, "value")!.set!.call(el, value);
  el.dispatchEvent(new Event("input", { bubbles: true }));
};

describe("TelegramAdminTab (GOAL-006 R5)", () => {
  it("loads status and shows write-only token/secret inputs with keep-current placeholder", async () => {
    const container = renderTab(async (input) => {
      if (String(input).startsWith("/api/channel/telegram/settings")) return jsonResponse(STATUS);
      return jsonResponse({}, 404);
    });
    await settle();

    expect(container.querySelector("[data-telegram-admin-tab]")).not.toBeNull();
    const token = container.querySelector("#telegram-bot-token") as HTMLInputElement;
    const secret = container.querySelector("#telegram-webhook-secret") as HTMLInputElement;
    expect(token).not.toBeNull();
    expect(secret).not.toBeNull();
    expect(token.type).toBe("password");
    expect(secret.type).toBe("password");
    // Values are never pre-filled from the server (write-only).
    expect(token.value).toBe("");
    expect(secret.value).toBe("");
  });

  it("loads operator sessions and transcript while fail-closing the composer before capability proof", async () => {
    const container = renderTab(async (input) => {
      const url = String(input);
      if (url === "/api/channel/telegram/settings") return jsonResponse(OPERATOR_STATUS);
      if (url === "/api/channel/telegram/operator/sessions?page=1&pageSize=100") {
        return jsonResponse({
          items: [{ chatId: "8001", chatType: "private", title: "Alice", username: "alice", lastMessageAt: "2026-09-05T00:00:00Z" }],
          total: 1,
          page: 1,
          pageSize: 100,
        });
      }
      if (url === "/api/channel/telegram/operator/sessions/8001/messages?page=1&pageSize=100") {
        return jsonResponse({
          items: [
            { chatId: "8001", direction: "inbound", status: "received", occurredAt: "2026-09-05T00:00:00Z", updateId: "9001", text: "hello" },
            { chatId: "8001", direction: "outbound", status: "failed", occurredAt: "2026-09-05T00:01:00Z", requestId: "operator-1", text: "reply" },
          ],
          total: 2,
          page: 1,
          pageSize: 100,
        });
      }
      return jsonResponse({}, 404);
    });
    await settle();

    expect(container.querySelector("[data-telegram-session='8001']")).not.toBeNull();
    expect(container.querySelector("[data-telegram-transcript]")?.textContent).toContain("hello");
    expect(container.querySelector("[data-telegram-transcript]")?.textContent).toContain("reply");
    const composer = container.querySelector("[data-telegram-composer]") as HTMLFieldSetElement;
    expect(composer).not.toBeNull();
    expect(composer.disabled).toBe(true);
    expect(container.querySelector("[data-telegram-composer]")?.textContent).toContain("Send");
    expect(container.querySelector("[data-telegram-retry='operator-1']")).toHaveProperty("disabled", true);
  });

  it("pauses the ten-second operator refresh while hidden and refreshes immediately on visibility restore", async () => {
    vi.useFakeTimers();
    try {
      const sessionCalls: string[] = [];
      const container = renderTab(async (input) => {
        const url = String(input);
        if (url === "/api/channel/telegram/settings") return jsonResponse(OPERATOR_STATUS);
        if (url.startsWith("/api/channel/telegram/operator/sessions?")) {
          sessionCalls.push(url);
          return jsonResponse({ items: [], total: 0, page: 1, pageSize: 100 });
        }
        return jsonResponse({}, 404);
      });
      await flushPromises();
      expect(sessionCalls).toHaveLength(1);

      await act(async () => {
        await vi.advanceTimersByTimeAsync(9_999);
      });
      expect(sessionCalls).toHaveLength(1);

      await act(async () => {
        await vi.advanceTimersByTimeAsync(1);
      });
      await flushPromises();
      expect(sessionCalls).toHaveLength(2);

      Object.defineProperty(document, "hidden", { configurable: true, value: true });
      await act(async () => document.dispatchEvent(new Event("visibilitychange")));
      await act(async () => {
        await vi.advanceTimersByTimeAsync(20_000);
      });
      expect(sessionCalls).toHaveLength(2);

      Object.defineProperty(document, "hidden", { configurable: true, value: false });
      await act(async () => document.dispatchEvent(new Event("visibilitychange")));
      await flushPromises();
      expect(sessionCalls).toHaveLength(3);
      expect(container.querySelector("[data-telegram-operator]")).not.toBeNull();
    } finally {
      Object.defineProperty(document, "hidden", { configurable: true, value: false });
      vi.useRealTimers();
    }
  });

  it("coalesces a visibility-triggered refresh with an in-flight operator refresh", async () => {
    vi.useFakeTimers();
    try {
      let releaseSessions!: (response: Response) => void;
      const pendingSessions = new Promise<Response>((resolve) => {
        releaseSessions = resolve;
      });
      let sessionCalls = 0;
      const container = renderTab(async (input) => {
        const url = String(input);
        if (url === "/api/channel/telegram/settings") return jsonResponse(OPERATOR_STATUS);
        if (url.startsWith("/api/channel/telegram/operator/sessions?")) {
          sessionCalls += 1;
          return pendingSessions;
        }
        return jsonResponse({}, 404);
      });

      await flushPromises();
      expect(sessionCalls).toBe(1);

      Object.defineProperty(document, "hidden", { configurable: true, value: true });
      await act(async () => document.dispatchEvent(new Event("visibilitychange")));
      Object.defineProperty(document, "hidden", { configurable: true, value: false });
      await act(async () => document.dispatchEvent(new Event("visibilitychange")));
      await flushPromises();
      expect(sessionCalls).toBe(1);

      await act(async () => {
        releaseSessions(jsonResponse({ items: [], total: 0, page: 1, pageSize: 100 }));
        await pendingSessions;
        await Promise.resolve();
      });
      expect(sessionCalls).toBe(1);
      expect(container.querySelector("[data-telegram-operator]")).not.toBeNull();
    } finally {
      Object.defineProperty(document, "hidden", { configurable: true, value: false });
      vi.useRealTimers();
    }
  });

  it("coalesces same-chat transcript loads while the messages request is in flight", async () => {
    vi.useFakeTimers();
    try {
      let releaseTimeline!: (response: Response) => void;
      const pendingTimeline = new Promise<Response>((resolve) => {
        releaseTimeline = resolve;
      });
      let timelineCalls = 0;
      const container = renderTab(async (input) => {
        const url = String(input);
        if (url === "/api/channel/telegram/settings") return jsonResponse(OPERATOR_STATUS);
        if (url === "/api/channel/telegram/operator/sessions?page=1&pageSize=100") {
          return jsonResponse({
            items: [{ chatId: "8001", chatType: "private", title: "Alice", lastMessageAt: "2026-09-05T00:00:00Z" }],
            total: 1,
            page: 1,
            pageSize: 100,
          });
        }
        if (url === "/api/channel/telegram/operator/sessions/8001/messages?page=1&pageSize=100") {
          timelineCalls += 1;
          return pendingTimeline;
        }
        return jsonResponse({}, 404);
      });

      await flushPromises();
      await flushPromises();
      expect(timelineCalls).toBe(1);
      const session = container.querySelector("[data-telegram-session='8001']") as HTMLButtonElement;
      expect(session).not.toBeNull();

      await act(async () => {
        session.click();
        await Promise.resolve();
      });
      expect(timelineCalls).toBe(1);

      await act(async () => {
        releaseTimeline(jsonResponse({ items: [], total: 0, page: 1, pageSize: 100 }));
        await pendingTimeline;
        await Promise.resolve();
        await Promise.resolve();
      });
      expect(timelineCalls).toBe(1);
    } finally {
      Object.defineProperty(document, "hidden", { configurable: true, value: false });
      vi.useRealTimers();
    }
  });

  it("hides the operator console while business handlers occupy the Telegram receiver", async () => {
    const calls: string[] = [];
    const container = renderTab(async (input) => {
      calls.push(String(input));
      if (String(input) === "/api/channel/telegram/settings") {
        return jsonResponse({ ...OPERATOR_STATUS, mode: "polling", receiver: "polling", business_occupied: true });
      }
      return jsonResponse({}, 404);
    });
    await settle();

    expect(container.querySelector("[data-telegram-operator]")).toBeNull();
    expect(calls).not.toContain("/api/channel/telegram/lease/acquire");
  });

  it("fail-closes the operator console when the occupancy signal is missing", async () => {
    const calls: string[] = [];
    const container = renderTab(async (input) => {
      calls.push(String(input));
      if (String(input) === "/api/channel/telegram/settings") {
        return jsonResponse({ ...OPERATOR_STATUS, mode: "polling", receiver: "polling", business_occupied: undefined });
      }
      return jsonResponse({}, 404);
    });
    await settle();

    expect(container.querySelector("[data-telegram-operator]")).toBeNull();
    expect(calls).not.toContain("/api/channel/telegram/lease/acquire");
  });

  it("acquires a polling lease and renders the live connection status", async () => {
    const calls: Array<{ url: string; method?: string }> = [];
    const container = renderTab(async (input, init) => {
      const url = String(input);
      calls.push({ url, method: init?.method });
      if (url === "/api/channel/telegram/settings") {
        return jsonResponse({ ...STATUS, connection_state: "idle", receiver: "none" });
      }
      return jsonResponse({ ok: true, connection_state: "running", receiver: "polling" });
    });
    await settle();

    expect(calls).toContainEqual({
      url: "/api/channel/telegram/lease/acquire",
      method: "POST",
    });
    expect(container.querySelector("[data-telegram-connection]")).not.toBeNull();
    expect(container.textContent).toContain("Connection: Running · Polling");
    expect(container.textContent).toContain("Console lease: Active");
  });

  it("releases the polling lease when the tab unmounts", async () => {
    const calls: Array<{ url: string; method?: string }> = [];
    const container = renderTab(async (input, init) => {
      calls.push({ url: String(input), method: init?.method });
      if (String(input) === "/api/channel/telegram/settings") return jsonResponse(STATUS);
      return jsonResponse({ ok: true, connection_state: "running", receiver: "polling" });
    });
    await settle();

    const mounted = activeRoots.find((entry) => entry.container === container);
    expect(mounted).toBeDefined();
    await act(async () => mounted!.root.unmount());
    activeRoots.splice(activeRoots.indexOf(mounted!), 1);
    mounted!.container.remove();
    await settle();

    expect(calls).toContainEqual({
      url: "/api/channel/telegram/lease/release",
      method: "POST",
    });
  });

  it("PATCHes receiver mode and webhook origin as non-secret settings", async () => {
    const calls: Array<{ url: string; method?: string; body?: unknown }> = [];
    const container = renderTab(async (input, init) => {
      const url = String(input);
      calls.push({ url, method: init?.method, body: init?.body === undefined ? undefined : JSON.parse(String(init.body)) });
      if (url === "/api/channel/telegram/settings" && init?.method === "PATCH") {
        return jsonResponse({ ...STATUS, mode: "webhook", webhook_public_base_url: "https://console.example", receiver: "webhook" });
      }
      return jsonResponse(STATUS);
    });
    await settle();

    await act(async () => {
      const mode = container.querySelector("#telegram-mode") as HTMLSelectElement;
      mode.value = "webhook";
      mode.dispatchEvent(new Event("change", { bubbles: true }));
      setNativeValue(container.querySelector("#telegram-webhook-public-base-url") as HTMLInputElement, "https://console.example");
    });

    const saveButton = [...container.querySelectorAll("button")].find((b) => b.textContent === "Save settings");
    expect(saveButton).toBeDefined();
    await act(async () => saveButton!.click());

    const patch = calls.find((c) => c.method === "PATCH");
    expect(patch?.body).toEqual({ mode: "webhook", webhook_public_base_url: "https://console.example" });
  });

  it("PATCHes only non-empty values and refreshes status + clears inputs", async () => {
    const calls: Array<{ url: string; method?: string; body?: unknown }> = [];
    const container = renderTab(async (input, init) => {
      const url = String(input);
      calls.push({ url, method: init?.method, body: init?.body === undefined ? undefined : JSON.parse(String(init.body)) });
      if (url.startsWith("/api/channel/telegram/settings") && init?.method === "PATCH") {
        return jsonResponse({ ...STATUS, captured_messages_count: 4 });
      }
      return jsonResponse(STATUS);
    });
    await settle();

    await act(async () => {
      setNativeValue(container.querySelector("#telegram-bot-token") as HTMLInputElement, "new-token");
    });
    const saveButton = [...container.querySelectorAll("button")].find((b) => b.textContent === "Save settings");
    expect(saveButton).toBeDefined();
    await act(async () => saveButton!.click());

    const patch = calls.find((c) => c.method === "PATCH");
    expect(patch).toBeDefined();
    expect(patch!.body).toEqual({ bot_token: "new-token" });

    const token = container.querySelector("#telegram-bot-token") as HTMLInputElement;
    expect(token.value).toBe(""); // cleared after successful save
  });

  it("offers a two-step clear action that PATCHes empty strings (R-004)", async () => {
    const calls: Array<{ url: string; method?: string; body?: unknown }> = [];
    const container = renderTab(async (input, init) => {
      const url = String(input);
      calls.push({ url, method: init?.method, body: init?.body === undefined ? undefined : JSON.parse(String(init.body)) });
      if (url.startsWith("/api/channel/telegram/settings") && init?.method === "PATCH") {
        return jsonResponse({ configured: false, token_set: false, secret_set: false, captured_messages_count: 0 });
      }
      return jsonResponse(STATUS);
    });
    await settle();

    // Clear button only appears when configured.
    const clearButton = [...container.querySelectorAll("button")].find((b) => b.textContent === "Clear saved secrets");
    expect(clearButton).toBeDefined();

    // First click arms confirmation (no request yet).
    await act(async () => clearButton!.click());
    expect(calls.find((c) => c.method === "PATCH")).toBeUndefined();

    // Confirm sends empty strings.
    const confirmButton = [...container.querySelectorAll("button")].find((b) => b.textContent === "Clear");
    expect(confirmButton).toBeDefined();
    await act(async () => confirmButton!.click());

    const patch = calls.find((c) => c.method === "PATCH");
    expect(patch).toBeDefined();
    expect(patch!.body).toEqual({ bot_token: "", webhook_secret: "" });

    // After clearing, the configured badge flips to "Not configured".
    expect(container.textContent).toContain("Not configured");
  });
});
