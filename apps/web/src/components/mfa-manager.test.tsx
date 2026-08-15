// @vitest-environment jsdom
// GOAL-018: MFA manager component flow (status → enroll → confirm →
// disable/rotate) with a mocked authFetch (global fetch).
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { MfaManager } from "@/components/mfa-manager";

// W11 · M-02: the component signs out locally after a successful disable.
// Mock useAuth so the bare renderer tests do not need a full AuthProvider.
const logoutMock = vi.fn<() => Promise<void>>().mockResolvedValue(undefined);
vi.mock("@/account/AuthContext", () => ({
  useAuth: () => ({ logout: logoutMock }),
}));

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
  logoutMock.mockClear();
  sessionStorage.clear();
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


  // A-003 F-002: the disable input accepts a recovery code (non-6-digit) and
  // the rotate flow returns a fresh set.
  it("disables with a recovery code and rotates recovery codes", async () => {
    // Flow order: status(0) → rotate(1, payload) → disable(2, 204).
    // W11 · M-02: after disable success the component reflects the local
    // disabled state and signs out (mocked logout) — no status refresh.
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(jsonResponse({ enabled: true, enrolledAt: "2026-08-15T00:00:00Z" })) // 0 status
      .mockResolvedValueOnce(jsonResponse({ recoveryCodes: ["N1", "N2"] })) // 1 rotate
      .mockResolvedValueOnce(new Response(null, { status: 204 })) // 2 disable
    vi.stubGlobal("fetch", fetchMock);

    const container = await renderManager();
    await act(async () => {});
    expect(container.querySelector('[data-mfa-active]')).not.toBeNull();

    // Rotate with a 6-digit code → fresh recovery set visible (A-004
    // recommended: unconditional, not gated on later UI state).
    const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
    const rotateInput = container.querySelector<HTMLInputElement>("#mfaRotateRecovery");
    expect(rotateInput).not.toBeNull();
    await act(async () => {
      setter?.call(rotateInput, "123456");
      rotateInput!.dispatchEvent(new Event("input", { bubbles: true }));
    });
    const rotateBtn = container.querySelector<HTMLButtonElement>('[data-mfa-active] form:nth-of-type(2) button[type="submit"]');
    expect(rotateBtn).not.toBeNull();
    await act(async () => rotateBtn!.click());
    await act(async () => {});
    expect(container.querySelector('[data-mfa-new-recovery]')).not.toBeNull();
    expect(container.textContent).toContain("N1");

    // Disable with a recovery code (non-6-digit).
    const disableInput = container.querySelector<HTMLInputElement>("#mfaDisableCode");
    expect(disableInput).not.toBeNull();
    await act(async () => {
      setter?.call(disableInput, "ABCDEFGH");
      disableInput!.dispatchEvent(new Event("input", { bubbles: true }));
    });
    const disableBtn = container.querySelector<HTMLButtonElement>('[data-mfa-active] button[type="submit"]');
    await act(async () => disableBtn!.click());
    await act(async () => {});
    // Local state reflects the disabled enrollment before the sign-out.
    expect(container.textContent).toContain("Enable MFA");
    // W11 · M-02: the notice flag is set and the local sign-out ran.
    expect(sessionStorage.getItem("mfa.disabledNotice")).toBe("1");
    expect(logoutMock).toHaveBeenCalledTimes(1);
    // The rotate request carried the 6-digit code; the disable request
    // carried the recovery code.
    const rotateBody = JSON.parse(String(fetchMock.mock.calls[1][1]?.body ?? "{}")) as { code?: string; recoveryCode?: string };
    expect(rotateBody.code).toBe("123456");
    const disableBody = JSON.parse(String(fetchMock.mock.calls[2][1]?.body ?? "{}")) as { code?: string; recoveryCode?: string };
    expect(disableBody.recoveryCode).toBe("ABCDEFGH");
  });

  // W11 · M-03: a wrong code on confirm is a 400 validation failure — inline
  // error, enroll payload retained, NO sign-out.
  it("confirm with a wrong code shows an inline error and keeps the enroll payload", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(jsonResponse({ enabled: false, enrolledAt: null })) // 0 status
      .mockResolvedValueOnce(
        jsonResponse({ secretBase32: "SECRET123", otpauthURL: "otpauth://totp/x", recoveryCodes: ["R1", "R2"] }), // 1 enroll
      )
      .mockResolvedValueOnce(jsonResponse({ error: "MFA_INVALID" }, 400)); // 2 confirm wrong code
    vi.stubGlobal("fetch", fetchMock);

    const container = await renderManager();
    await act(async () => {});
    const enroll = container.querySelector<HTMLButtonElement>('button[type="button"]');
    await act(async () => enroll!.click());
    await act(async () => {});

    const codeInput = container.querySelector<HTMLInputElement>("#mfaConfirmCode");
    expect(codeInput).not.toBeNull();
    const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
    await act(async () => {
      setter?.call(codeInput, "000000");
      codeInput!.dispatchEvent(new Event("input", { bubbles: true }));
    });
    await act(async () => {
      container.querySelector<HTMLButtonElement>('[data-mfa-enroll] button[type="submit"]')!.click();
    });
    await act(async () => {});

    // Error shown, payload retained, no sign-out.
    expect(container.textContent).toContain("Invalid second-factor code");
    expect(container.querySelector('[data-mfa-secret]')).not.toBeNull();
    expect(logoutMock).not.toHaveBeenCalled();
  });

  // W11 · M-02: a wrong code on disable is a 400 validation failure — inline
  // error, active state retained, NO sign-out.
  it("disable with a wrong code shows an inline error and does not sign out", async () => {
    const fetchMock = vi
      .fn()
      .mockResolvedValueOnce(jsonResponse({ enabled: true, enrolledAt: "2026-08-15T00:00:00Z" })) // 0 status
      .mockResolvedValueOnce(jsonResponse({ error: "MFA_INVALID" }, 400)); // 1 disable wrong code
    vi.stubGlobal("fetch", fetchMock);

    const container = await renderManager();
    await act(async () => {});
    expect(container.querySelector('[data-mfa-active]')).not.toBeNull();

    const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
    const disableInput = container.querySelector<HTMLInputElement>("#mfaDisableCode");
    await act(async () => {
      setter?.call(disableInput, "000000");
      disableInput!.dispatchEvent(new Event("input", { bubbles: true }));
    });
    await act(async () => {
      container.querySelector<HTMLButtonElement>('[data-mfa-active] button[type="submit"]')!.click();
    });
    await act(async () => {});

    expect(container.textContent).toContain("Invalid second-factor code");
    expect(container.querySelector('[data-mfa-active]')).not.toBeNull();
    expect(sessionStorage.getItem("mfa.disabledNotice")).toBeNull();
    expect(logoutMock).not.toHaveBeenCalled();
  });
  it("degrades to an unavailable placeholder when the status call fails", async () => {
    vi.stubGlobal("fetch", vi.fn().mockRejectedValue(new Error("network")));
    const container = await renderManager();
    await act(async () => {});
    expect(container.textContent).toContain("MFA status is unavailable.");
  });
});