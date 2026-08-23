// @vitest-environment jsdom

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { AuthError } from "@/account/auth-client";
import { LoginPage } from "@/app/LoginPage";
import { I18nProvider } from "@/i18n/runtime";

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

async function renderLogin(onLogin: (u: string, p: string) => Promise<void>) {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <LoginPage onLogin={onLogin} />
      </I18nProvider>,
    );
  });
  return container;
}

// React 19 tracks the input value via a property descriptor; assign through the
// native setter and dispatch `input` so onChange fires in the test DOM.
function fill(container: HTMLDivElement, selector: string, value: string) {
  const el = container.querySelector<HTMLInputElement>(selector);
  expect(el).not.toBeNull();
  const setter = Object.getOwnPropertyDescriptor(
    window.HTMLInputElement.prototype,
    "value",
  )?.set;
  act(() => {
    setter?.call(el, value);
    el!.dispatchEvent(new Event("input", { bubbles: true }));
  });
}

describe("LoginPage", () => {

  // S-10 (GOAL-017 D-002 §3 / A-007 F-001 contract): when the server demands a
  // second factor the login pauses at the code stage; the entered code
  // resolves the pending two-step promise.
  it("renders the second-factor stage and completes the two-step login", async () => {
    let entered: { code: string; recoveryCode?: string } | null = null;
    const onLogin = vi.fn(
      async (
        _u: string,
        _p: string,
        _c?: unknown,
        resolveMFA?: (proof: string) => Promise<{ code: string; recoveryCode?: string }>,
      ) => {
        if (!resolveMFA) {
          return;
        }
        entered = await resolveMFA("proof-1");
      },
    );
    const container = await renderLogin(onLogin);
    fill(container, "#username", "admin");
    fill(container, "#password", "admin");
    const submit = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    await act(async () => submit!.click());

    // The code stage is now visible while the login promise stays pending.
    expect(container.querySelector('[data-mfa-stage]')).not.toBeNull();
    fill(container, "#mfaCode", "123456");
    const verify = container.querySelector<HTMLButtonElement>("[data-mfa-verify]");
    expect(verify).not.toBeNull();
    await act(async () => verify!.click());

    expect(entered).toEqual({ code: "123456" });
    expect(container.querySelector('[data-mfa-stage]')).toBeNull();
    expect(onLogin).toHaveBeenCalledWith("admin", "admin", undefined, expect.any(Function));
  });

  it("renders the sign-in form and submits credentials", async () => {
    const onLogin = vi.fn().mockResolvedValue(undefined);
    const container = await renderLogin(onLogin);

    expect(container.textContent).toContain("Sign in");
    fill(container, "#username", "admin");
    fill(container, "#password", "admin");

    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    expect(button).not.toBeNull();
    await act(async () => button!.click());
    expect(onLogin).toHaveBeenCalledWith("admin", "admin", undefined, expect.any(Function));
  });

  it("surfaces the server error when login fails", async () => {
    const onLogin = vi.fn().mockRejectedValue(new AuthError("INVALID_CREDENTIALS", "invalid username or password"));
    const container = await renderLogin(onLogin);

    fill(container, "#username", "admin");
    fill(container, "#password", "wrong");
    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    await act(async () => button!.click());

    expect(container.textContent).toContain("invalid username or password");
    expect(container.querySelector('[role="alert"]')).not.toBeNull();
  });

  // W4 P2-2: login.error.failed carries a `{status}` placeholder; the AuthError
  // keeps the real HTTP status so it interpolates instead of rendering the
  // literal "{status}".
  it("interpolates the HTTP status into login.error.failed (no literal {status})", async () => {
    const onLogin = vi
      .fn()
      .mockRejectedValue(new AuthError("LOGIN_FAILED", "login failed: HTTP 503", 503));
    const container = await renderLogin(onLogin);

    fill(container, "#username", "admin");
    fill(container, "#password", "admin");
    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    await act(async () => button!.click());

    const text = container.textContent ?? "";
    expect(text).toContain("503");
    expect(text).not.toContain("{status}");
  });

  it("disables submit until both fields are filled", async () => {
    const container = await renderLogin(vi.fn());
    const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
    expect(button?.disabled).toBe(true);
    fill(container, "#username", "admin");
    expect(button?.disabled).toBe(true);
    fill(container, "#password", "admin");
    expect(button?.disabled).toBe(false);
  });

  it("shows the local seed hint in development builds (F-002-004)", async () => {
    const env = import.meta.env as { DEV: boolean };
    const dev = env.DEV;
    env.DEV = true;
    try {
      const container = await renderLogin(vi.fn());
      expect(container.textContent).toContain("admin / admin");
    } finally {
      env.DEV = dev;
    }
  });

  it("hides the local seed hint in production builds (F-002-004)", async () => {
    const env = import.meta.env as { DEV: boolean };
    const dev = env.DEV;
    env.DEV = false;
    try {
      const container = await renderLogin(vi.fn());
      expect(container.textContent).not.toContain("admin / admin");
      expect(container.textContent).not.toContain("Local development seed");
    } finally {
      env.DEV = dev;
    }
  });


  // S-11 (GOAL-011 D-002 §5): when the preflight reports the gate enabled, the
  // page renders the challenge and submits captchaId/captchaAnswer.
  it("renders the captcha challenge when the preflight reports it enabled", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({
        enabled: true,
        challenge: { id: "cap-1", question: "7 + 3 = ?", expiresInSeconds: 300 },
      }),
    });
    vi.stubGlobal("fetch", fetchMock);
    try {
      const onLogin = vi.fn().mockResolvedValue(undefined);
      const container = await renderLogin(onLogin);
      expect(container.querySelector('[data-captcha-question]')?.textContent).toContain("7 + 3");
      fill(container, "#username", "admin");
      fill(container, "#password", "admin");
      fill(container, "#captchaAnswer", "10");
      const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
      await act(async () => button!.click());
      expect(onLogin).toHaveBeenCalledWith("admin", "admin", { id: "cap-1", answer: "10" }, expect.any(Function));
    } finally {
      vi.unstubAllGlobals();
    }
  });

  it("refreshes the captcha challenge when the user clicks 换一题 (W16-F08)", async () => {
    let calls = 0;
    const fetchMock = vi.fn().mockImplementation(() => {
      calls += 1;
      return Promise.resolve({
        ok: true,
        json: async () => ({
          enabled: true,
          challenge: { id: "cap-" + calls, question: calls === 1 ? "1 + 1 = ?" : "2 + 2 = ?", expiresInSeconds: 300 },
        }),
      });
    });
    vi.stubGlobal("fetch", fetchMock);
    try {
      const container = await renderLogin(vi.fn().mockResolvedValue(undefined));
      const refresh = container.querySelector<HTMLButtonElement>("[data-captcha-refresh]");
      expect(refresh).not.toBeNull();
      const before = calls;
      await act(async () => refresh!.click());
      await act(async () => {});
      expect(calls).toBeGreaterThan(before);
      expect(container.querySelector('[data-captcha-question]')?.textContent).toContain("2 + 2");
    } finally {
      vi.unstubAllGlobals();
    }
  });

  it("stays captcha-free and still logs in when the preflight reports disabled", async () => {
    const fetchMock = vi.fn().mockResolvedValue({
      ok: true,
      json: async () => ({ enabled: false }),
    });
    vi.stubGlobal("fetch", fetchMock);
    try {
      const onLogin = vi.fn().mockResolvedValue(undefined);
      const container = await renderLogin(onLogin);
      expect(container.querySelector("#captchaAnswer")).toBeNull();
      fill(container, "#username", "admin");
      fill(container, "#password", "admin");
      const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
      await act(async () => button!.click());
      expect(onLogin).toHaveBeenCalledWith("admin", "admin", undefined, expect.any(Function));
    } finally {
      vi.unstubAllGlobals();
    }
  });

  it("surfaces the INVALID_CAPTCHA error and refreshes the challenge (F-009)", async () => {
    let calls = 0;
    const fetchMock = vi.fn().mockImplementation(() => {
      calls += 1;
      return Promise.resolve({
        ok: true,
        json: async () => ({
          enabled: true,
          challenge: { id: "cap-" + calls, question: "3 + 4 = ?", expiresInSeconds: 300 },
        }),
      });
    });
    vi.stubGlobal("fetch", fetchMock);
    try {
      const onLogin = vi
        .fn()
        .mockRejectedValue(new AuthError("INVALID_CAPTCHA", "captcha verification failed", 400));
      const container = await renderLogin(onLogin);
      fill(container, "#username", "admin");
      fill(container, "#password", "admin");
      const button = container.querySelector<HTMLButtonElement>('button[type="submit"]');
      await act(async () => button!.click());
      expect(container.textContent).toContain("captcha verification failed");
      // A fresh challenge was fetched after the rejection (consumed server-side).
      expect(calls).toBeGreaterThanOrEqual(2);
      expect(container.querySelector('[data-captcha-question]')?.textContent).toContain("3 + 4");
    } finally {
      vi.unstubAllGlobals();
    }
  });

  it("consumes design-system Card/Input/Label primitives (S3 / D-004 Sign in)", async () => {
    const container = await renderLogin(vi.fn());
    expect(container.querySelector('[data-login-surface="design-system"]')).not.toBeNull();
    // Input primitive uses shadow-sm + transparent bg classes from components/ui/input
    const username = container.querySelector<HTMLInputElement>("#username");
    expect(username).not.toBeNull();
    expect(username?.className).toMatch(/shadow-sm/);
    const labels = container.querySelectorAll("label");
    expect(labels.length).toBeGreaterThanOrEqual(2);
    // Card surface (rounded-lg border from ui/card)
    const card = container.querySelector(".rounded-lg.border");
    expect(card).not.toBeNull();
  });

  // W11 · M-02: after a successful MFA disable the app signs out locally and
  // the login page shows a one-time notice (flag consumed on mount).
  it("shows a one-time notice when mfa.disabledNotice is set", async () => {
    sessionStorage.setItem("mfa.disabledNotice", "1");
    const container = await renderLogin(vi.fn());
    const banner = container.querySelector('[data-login-notice="mfa-disabled"]');
    expect(banner).not.toBeNull();
    expect(banner!.textContent).toContain("MFA was disabled");
    // The flag was consumed — a second mount shows no banner.
    sessionStorage.clear();
    const container2 = await renderLogin(vi.fn());
    expect(container2.querySelector('[data-login-notice="mfa-disabled"]')).toBeNull();
  });

  // F-VUI-011: password visibility toggle
  it("password input starts as type=password", async () => {
    const container = await renderLogin(vi.fn());
    const input = container.querySelector<HTMLInputElement>("#password");
    expect(input).not.toBeNull();
    expect(input!.type).toBe("password");
  });

  it("clicking the toggle reveals the password (type changes to text) and updates aria-label", async () => {
    const container = await renderLogin(vi.fn());
    const toggle = container.querySelector<HTMLButtonElement>("[data-password-toggle]");
    expect(toggle).not.toBeNull();
    // Initial aria-label is "Show password"
    expect(toggle!.getAttribute("aria-label")).toMatch(/show/i);
    await act(async () => toggle!.click());
    const input = container.querySelector<HTMLInputElement>("#password");
    expect(input!.type).toBe("text");
    // aria-label switches to "Hide password"
    expect(toggle!.getAttribute("aria-label")).toMatch(/hide/i);
  });

  it("clicking the toggle twice restores type=password and original aria-label", async () => {
    const container = await renderLogin(vi.fn());
    const toggle = container.querySelector<HTMLButtonElement>("[data-password-toggle]");
    expect(toggle).not.toBeNull();
    await act(async () => toggle!.click());
    await act(async () => toggle!.click());
    const input = container.querySelector<HTMLInputElement>("#password");
    expect(input!.type).toBe("password");
    expect(toggle!.getAttribute("aria-label")).toMatch(/show/i);
  });
});