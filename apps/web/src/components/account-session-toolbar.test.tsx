// @vitest-environment jsdom
import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { I18nProvider } from "@/i18n/runtime";
import { AccountSessionToolbar } from "@/components/account-session-toolbar";

const { authFetchMock, refreshSessionMock, setAccessTokenMock, setRefreshTokenMock } = vi.hoisted(() => ({
  authFetchMock: vi.fn(),
  refreshSessionMock: vi.fn().mockResolvedValue(undefined),
  setAccessTokenMock: vi.fn(),
  setRefreshTokenMock: vi.fn(),
}));
vi.mock("@/account/AuthContext", () => ({
  useAuth: () => ({ authFetch: authFetchMock, refreshSession: refreshSessionMock }),
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
  vi.unstubAllGlobals();
});

async function renderComponent() {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  await act(async () => {
    root.render(
      <I18nProvider>
        <AccountSessionToolbar node={{ type: "custom", component: "account-session-toolbar" }} context={{}} />
      </I18nProvider>,
    );
  });
  return container;
}

describe("AccountSessionToolbar", () => {
  it("revokes other sessions and stores the reissued tokens", async () => {
    vi.stubGlobal("confirm", vi.fn(() => true));
    authFetchMock.mockResolvedValue(
      new Response(JSON.stringify({ accessToken: "at-new", refreshToken: "rt-new" }), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      }),
    );
    const container = await renderComponent();
    const button = container.querySelector<HTMLButtonElement>("[data-revoke-others]");
    expect(button).not.toBeNull();
    await act(async () => button!.click());
    await act(async () => {});

    expect(authFetchMock).toHaveBeenCalledWith(
      "/api/account/sessions/revoke-others",
      expect.objectContaining({ method: "POST" }),
    );
    expect(setAccessTokenMock).toHaveBeenCalledWith("at-new");
    expect(setRefreshTokenMock).toHaveBeenCalledWith("rt-new");
    expect(refreshSessionMock).toHaveBeenCalled();
  });
});
