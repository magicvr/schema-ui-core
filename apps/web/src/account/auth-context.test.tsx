// @vitest-environment jsdom
//
// W13 T-05 follow-up: the account.profile config-change event must refresh
// the session (re-resolve /me) so the shell header avatar/name updates
// immediately after a profile save — no page reload required.

import { act } from "react";
import { createRoot, type Root } from "react-dom/client";
import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import { AuthProvider, useAuth } from "@/account/AuthContext";
import type { AuthUser } from "@/account/auth-client";
import { notifyConfigChanged } from "@/app/config-events";

const mocks = vi.hoisted(() => ({
  restoreSession: vi.fn(),
  fetchMe: vi.fn(),
}));

vi.mock("@/account/auth-client", async (importOriginal) => {
  const actual = await importOriginal<typeof import("@/account/auth-client")>();
  return {
    ...actual,
    restoreSession: mocks.restoreSession,
    fetchMe: mocks.fetchMe,
  };
});

const activeRoots: Array<{ root: Root; container: HTMLDivElement }> = [];

beforeEach(() => {
  Object.defineProperty(globalThis, "IS_REACT_ACT_ENVIRONMENT", { configurable: true, value: true });
  localStorage.clear();
  mocks.restoreSession.mockReset();
  mocks.fetchMe.mockReset();
});

afterEach(async () => {
  for (const { root, container } of activeRoots.splice(0)) {
    await act(async () => root.unmount());
    container.remove();
  }
});

function baseUser(): AuthUser {
  return { id: "user-admin", name: "Admin", roles: ["admin"] };
}

function renderProvider(userSink: (user: AuthUser | null) => void): HTMLDivElement {
  const container = document.createElement("div");
  document.body.appendChild(container);
  const root = createRoot(container);
  activeRoots.push({ root, container });
  act(() => {
    root.render(
      <AuthProvider>
        <SessionProbe onUser={userSink} />
      </AuthProvider>,
    );
  });
  return container;
}

function SessionProbe({ onUser }: { onUser: (user: AuthUser | null) => void }) {
  const { user } = useAuth();
  onUser(user);
  return <div data-session-probe />;
}

describe("AuthProvider account.profile session refresh (W13 T-05)", () => {
  it("re-resolves /me on the account.profile config-change event", async () => {
    mocks.restoreSession.mockResolvedValue({ kind: "none" });
    const seen: Array<AuthUser | null> = [];
    renderProvider((user) => seen.push(user));
    await act(async () => {
      await Promise.resolve();
    });
    expect(seen[seen.length - 1]).toBeNull();

    // The profile save publishes the event; the provider refreshes /me.
    mocks.fetchMe.mockResolvedValue({
      user: { ...baseUser(), avatarUrl: "/api/account/avatars/abc123" },
      features: { menu_users: true },
    });
    await act(async () => {
      notifyConfigChanged("account.profile");
      await Promise.resolve();
    });
    expect(mocks.fetchMe).toHaveBeenCalledTimes(1);
    expect(seen[seen.length - 1]?.avatarUrl).toBe("/api/account/avatars/abc123");
  });

  it("ignores config-change events of other namespaces", async () => {
    mocks.restoreSession.mockResolvedValue({ kind: "none" });
    renderProvider(() => undefined);
    await act(async () => {
      await Promise.resolve();
    });
    await act(async () => {
      notifyConfigChanged("settings.branding");
      await Promise.resolve();
    });
    expect(mocks.fetchMe).not.toHaveBeenCalled();
  });

  it("keeps the current session when the refresh /me fails", async () => {
    mocks.restoreSession.mockResolvedValue({ kind: "none" });
    mocks.fetchMe.mockRejectedValue(new Error("network down"));
    const seen: Array<AuthUser | null> = [];
    renderProvider((user) => seen.push(user));
    await act(async () => {
      await Promise.resolve();
    });
    await act(async () => {
      notifyConfigChanged("account.profile");
      await Promise.resolve();
    });
    // No session was established; the failed refresh must not throw.
    expect(seen.every((entry) => entry === null)).toBe(true);
  });
});
