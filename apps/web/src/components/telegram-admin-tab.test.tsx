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
  bot_username: "fixture_bot",
  captured_messages_count: 3,
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
