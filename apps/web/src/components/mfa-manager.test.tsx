// @vitest-environment jsdom
// GOAL-018: MFA manager component flow (status → enroll → confirm →
// disable/rotate) with a mocked authFetch (global fetch).
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { MfaManager } from "@/components/mfa-manager";

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

async function renderManager() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <MfaManager node={{ type: "custom", component: "mfa-manager" }} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), { status, headers: { "Content-Type": "application/json" } });
}

describe("MfaManager", () => {
  it("shows the disabled state and runs the enroll → confirm flow", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(jsonResponse({ enabled: false, enrolledAt: null })) // status
      .mockResolvedValueOnce(
        jsonResponse({ secretBase32: "SECRET123", otpauthURL: "otpauth://totp/x", recoveryCodes: ["R1", "R2"] }), // enroll
      )
      .mockResolvedValueOnce(new Response(null, { status: 204 })) // confirm
      .mockResolvedValueOnce(jsonResponse({ enabled: true, enrolledAt: "2026-08-15T00:00:00Z" })); // status after
    vi.stubGlobal("fetch", fetchMock);

    const container = await renderManager();
    await act(async () => {});
    expect(container.textContent).toContain("Enable MFA");

    const enroll = container.querySelector<HTMLButtonElement>('button[type="button"]');
    expect(enroll).not.toBeNull();
    await act(async () => enroll!.click());
    await act(async () => {});

    // One-time payload is visible.
    const secretInput = container.querySelector<HTMLInputElement>('[data-mfa-secret]');
    expect(secretInput).not.toBeNull();
    expect(secretInput!.value).toBe("SECRET123");
    expect(container.textContent).toContain("R1");

    // Confirm with a code.
    const codeInput = container.querySelector<HTMLInputElement>("#mfaConfirmCode");
    expect(codeInput).not.toBeNull();
    const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
    await act(async () => {
      setter?.call(codeInput, "123456");
      codeInput!.dispatchEvent(new Event("input", { bubbles: true }));
    });
    const confirmBtn = container.querySelector<HTMLButtonElement>('[data-mfa-enroll] button[type="submit"]');
    expect(confirmBtn).not.toBeNull();
    await act(async () => confirmBtn!.click());
    await act(async () => {});

    // Active state: disable form present.
    expect(container.querySelector('[data-mfa-active]')).not.toBeNull();
    expect(fetchMock.mock.calls.length).toBeGreaterThanOrEqual(4);
  });

  it("degrades to an unavailable placeholder when the status call fails", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("network")));
    const container = await renderManager();
    await act(async () => {});
    expect(container.textContent).toContain("MFA status is unavailable.");
  });
});