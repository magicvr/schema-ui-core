// @vitest-environment jsdom
// workspace-018 R3 · GOAL-004: email identity card flows (unbound → bind →
// pending → verify) with a mocked global fetch (authFetch passthrough).
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { EmailIdentityCard } from "@/components/email-identity";

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
  vi.restoreAllMocks();
});

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
}

async function renderCard(fetchMock: ReturnType<typeof vi.fn>) {
  vi.stubGlobal("fetch", fetchMock);
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <EmailIdentityCard node={{ type: "custom", component: "email-identity" }} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("EmailIdentityCard", () => {
  it("renders the unbound state and binds an address", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(jsonResponse({ name: "Admin", email: null, emailStatus: null }))
      .mockResolvedValue(jsonResponse({ status: "pending" }));
    const container = await renderCard(fetchMock);

    expect(container.querySelector('[data-testid="email-bind-form"]')).not.toBeNull();

    const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
    const input = container.querySelector<HTMLInputElement>("#email-identity-address");
    expect(input).not.toBeNull();
    await act(async () => {
      setter?.call(input, "alice@example.com");
      input?.dispatchEvent(new Event("input", { bubbles: true }));
    });
    const submit = container.querySelector<HTMLButtonElement>('[data-testid="email-bind-form"] button[type="submit"]');
    expect(submit).not.toBeNull();
    await act(async () => submit!.click());

    const bindCall = fetchMock.mock.calls.find(([url]) => String(url).includes("/api/account/email/bind"));
    expect(bindCall).toBeDefined();
    expect(String(bindCall![1]?.body)).toContain("alice@example.com");
  });

  it("shows the verify form while pending and surfaces cooldown errors", async () => {
    const fetchMock = vi.fn()
      .mockResolvedValueOnce(jsonResponse({ name: "Admin", email: "alice@example.com", emailStatus: "pending" }))
      // resend answered with a cooldown rejection
      .mockResolvedValueOnce(jsonResponse({ error: "EMAIL_RESEND_COOLDOWN" }, 429));
    const container = await renderCard(fetchMock);

    expect(container.querySelector('[data-testid="email-verify-form"]')).not.toBeNull();

    const resend = container.querySelectorAll<HTMLButtonElement>('[data-testid="email-verify-form"] button')[1];
    expect(resend).not.toBeNull();
    await act(async () => resend.click());

    const alert = container.querySelector('[role="alert"]');
    expect(alert).not.toBeNull();
    expect(alert!.textContent).toContain("Too many requests");
  });

  it("shows the verified state without forms", async () => {
    const fetchMock = vi.fn().mockResolvedValue(
      jsonResponse({ name: "Admin", email: "alice@example.com", emailStatus: "verified" }),
    );
    const container = await renderCard(fetchMock);
    expect(container.querySelector('[data-testid="email-bind-form"]')).toBeNull();
    expect(container.querySelector('[data-testid="email-verify-form"]')).toBeNull();
    expect(container.textContent).toContain("Email verified.");
  });
});
