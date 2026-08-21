// @vitest-environment jsdom
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { ForcePasswordChange } from "@/components/force-password-change";

const { authFetchMock, refreshSessionMock, logoutMock, setAccessTokenMock, setRefreshTokenMock } = vi.hoisted(() => ({
  authFetchMock: vi.fn(),
  refreshSessionMock: vi.fn().mockResolvedValue(undefined),
  logoutMock: vi.fn().mockResolvedValue(undefined),
  setAccessTokenMock: vi.fn(),
  setRefreshTokenMock: vi.fn(),
}));
vi.mock("@/account/AuthContext", () => ({
  useAuth: () => ({ authFetch: authFetchMock, refreshSession: refreshSessionMock, logout: logoutMock }),
}));
vi.mock("@/account/tokens", () => ({
  setAccessToken: setAccessTokenMock,
  setRefreshToken: setRefreshTokenMock,
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
  vi.clearAllMocks();
});

async function renderComponent() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <ForcePasswordChange />
      </I18nProvider>,
    );
  });
  return container;
}

function fill(container: HTMLDivElement, selector: string, value: string) {
  const el = container.querySelector<HTMLInputElement>(selector);
  const setter = Object.getOwnPropertyDescriptor(window.HTMLInputElement.prototype, "value")?.set;
  act(() => {
    setter?.call(el, value);
    el!.dispatchEvent(new Event("input", { bubbles: true }));
  });
}

describe("ForcePasswordChange", () => {
  it("submits the forced change and stores the reissued token pair", async () => {
    authFetchMock.mockResolvedValue(
      new Response(JSON.stringify({ accessToken: "at-new", refreshToken: "rt-new" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    );
    const container = await renderComponent();
    fill(container, "#currentPassword", "initial-pass");
    fill(container, "#newPassword", "new-secret-123");
    fill(container, "#confirmPassword", "new-secret-123");
    const form = container.querySelector("form");
    await act(async () => form!.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true })));
    await act(async () => {});

    expect(authFetchMock).toHaveBeenCalledWith(
      "/api/account/password",
      expect.objectContaining({ method: "POST" }),
    );
    expect(setAccessTokenMock).toHaveBeenCalledWith("at-new");
    expect(setRefreshTokenMock).toHaveBeenCalledWith("rt-new");
    expect(refreshSessionMock).toHaveBeenCalled();
  });

  it("shows a mismatch error without calling the API", async () => {
    const container = await renderComponent();
    fill(container, "#currentPassword", "initial-pass");
    fill(container, "#newPassword", "new-secret-123");
    fill(container, "#confirmPassword", "different");
    const form = container.querySelector("form");
    await act(async () => form!.dispatchEvent(new Event("submit", { bubbles: true, cancelable: true })));
    expect(authFetchMock).not.toHaveBeenCalled();
    expect(container.textContent).toContain("does not match");
  });
});
