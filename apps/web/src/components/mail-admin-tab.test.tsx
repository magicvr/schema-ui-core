// @vitest-environment jsdom
//
// R7 UX refinement (user-requested 2026-08-24): the mail admin tab shows only
// the selected channel's settings, the mock record table renders ONLY under
// the mock channel, and the test composer offers subject/body inputs.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { MailAdminTab } from "@/components/mail-admin-tab";
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
  id: "mail-admin-tab",
  component: "mail-admin-tab",
  props: {},
};

const CONFIG = {
  channel: "mock",
  mockRetention: 500,
  resend: { from: "" },
  smtp: { host: "", port: 0, username: "", from: "" },
  secrets: { resendApiKeySet: true, smtpPasswordSet: false },
  updated_at: "2026-08-24T00:00:00Z",
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
        <MailAdminTab node={NODE as never} context={{}} />
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

describe("MailAdminTab (R7 UX refinement)", () => {
  it("shows only the selected channel's settings and hides resend/smtp blocks", async () => {
    const container = renderTab(async (input) => {
      const url = String(input);
      if (url.startsWith("/api/mail/config")) return jsonResponse(CONFIG);
      if (url.startsWith("/api/mail/outbox")) {
        return jsonResponse({ items: [{ id: "o1", to: "a@example.com", subject: "s", created_at: "2026-08-24T00:00:00Z" }], total: 1 });
      }
      return jsonResponse({}, 404);
    });
    await settle();

    // mock selected: retention visible, table visible, resend/smtp absent.
    expect(container.querySelector("#mail-mock-retention")).not.toBeNull();
    expect(container.querySelector("[data-mail-admin-tab] table")).not.toBeNull();
    expect(container.querySelector("#mail-resend-from")).toBeNull();
    expect(container.querySelector("#mail-smtp-host")).toBeNull();

    // Switch to resend: retention + table disappear, resend fields appear.
    const select = container.querySelector("#mail-channel") as HTMLSelectElement;
    await act(async () => {
      select.value = "resend";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    expect(container.querySelector("#mail-resend-from")).not.toBeNull();
    expect(container.querySelector("#mail-mock-retention")).toBeNull();
    expect(container.querySelector("[data-mail-admin-tab] table")).toBeNull();

    // Switch to smtp: smtp fields appear.
    await act(async () => {
      select.value = "smtp";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    expect(container.querySelector("#mail-smtp-host")).not.toBeNull();
    expect(container.querySelector("#mail-smtp-password")).not.toBeNull();
  });

  it("saves via PUT with flat channel-scoped keys and sends composed subject/body", async () => {
    const calls: Array<{ url: string; method?: string; body?: unknown }> = [];
    const container = renderTab(async (input, init) => {
      const url = String(input);
      calls.push({ url, method: init?.method, body: init?.body === undefined ? undefined : JSON.parse(String(init.body)) });
      if (url.startsWith("/api/mail/config") && init?.method === "PUT") {
        return jsonResponse({ ...CONFIG, channel: "resend", resend: { from: "no-reply@eshowy.top" }, secrets: { resendApiKeySet: true, smtpPasswordSet: false } });
      }
      if (url.startsWith("/api/mail/config")) return jsonResponse(CONFIG);
      if (url.startsWith("/api/mail/test-send")) {
        return jsonResponse({ sent: true, channel: "resend" });
      }
      return jsonResponse({ items: [], total: 0 });
    });
    await settle();

    // Switch to resend and fill the resend fields via native setters.
    const select = container.querySelector("#mail-channel") as HTMLSelectElement;
    await act(async () => {
      select.value = "resend";
      select.dispatchEvent(new Event("change", { bubbles: true }));
    });
    const setNativeValue = (el: HTMLInputElement | HTMLTextAreaElement, value: string) => {
      const proto = el instanceof HTMLTextAreaElement ? HTMLTextAreaElement.prototype : HTMLInputElement.prototype;
      Object.getOwnPropertyDescriptor(proto, "value")!.set!.call(el, value);
      el.dispatchEvent(new Event("input", { bubbles: true }));
    };
    await act(async () => {
      setNativeValue(container.querySelector("#mail-resend-key") as HTMLInputElement, "re-new");
      setNativeValue(container.querySelector("#mail-resend-from") as HTMLInputElement, "no-reply@eshowy.top");
    });

    const saveButton = [...container.querySelectorAll("button")].find((b) => b.textContent === "Save configuration");
    expect(saveButton).toBeDefined();
    await act(async () => saveButton!.click());

    const put = calls.find((c) => c.method === "PUT");
    expect(put).toBeDefined();
    expect(put!.body).toEqual({ channel: "resend", resendFrom: "no-reply@eshowy.top", resendApiKey: "re-new" });

    // Compose a test mail with custom subject/body.
    await act(async () => {
      setNativeValue(container.querySelector("#mail-test-to") as HTMLInputElement, "magicvr@hotmail.com");
      setNativeValue(container.querySelector("#mail-test-subject") as HTMLInputElement, "hello");
      setNativeValue(container.querySelector("#mail-test-body") as HTMLTextAreaElement, "world");
    });
    const sendButton = [...container.querySelectorAll("button")].find((b) => b.textContent === "Send test mail");
    expect(sendButton).toBeDefined();
    await act(async () => sendButton!.click());

    const post = calls.find((c) => c.method === "POST");
    expect(post).toBeDefined();
    expect(post!.body).toEqual({ to: "magicvr@hotmail.com", subject: "hello", body: "world" });
  });
});
