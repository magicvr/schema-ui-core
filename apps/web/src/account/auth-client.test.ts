// @vitest-environment jsdom

import { afterEach, beforeEach, describe, expect, it, vi } from "vitest";

import {
  authFetch,
  fetchMe,
  login,
  logout,
  restoreSession,
  setAuthLostListener,
} from "@/account/auth-client";
import {
  clearTokens,
  getAccessToken,
  getRefreshToken,
  setAccessToken,
  setRefreshToken,
} from "@/account/tokens";

function jsonResponse(body: unknown, status = 200): Response {
  return new Response(JSON.stringify(body), {
    status,
    headers: { "Content-Type": "application/json" },
  });
}

function emptyResponse(status = 204): Response {
  return new Response(null, { status });
}

const SESSION = {
  user: {
    id: "user-admin",
    name: "Admin",
    roles: ["admin", "editor"],
    permissions: ["users.read", "users.write", "settings.write"],
  },
  features: {},
};

describe("auth-client", () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    clearTokens();
    setAuthLostListener(null);
    fetchMock = vi.fn();
    globalThis.fetch = fetchMock as unknown as typeof fetch;
  });

  afterEach(() => {
    setAuthLostListener(null);
  });

  it("login stores the token pair and resolves features via /me", async () => {
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ accessToken: "a1", refreshToken: "r1", ...SESSION }))
      .mockResolvedValueOnce(
        jsonResponse({
          user: SESSION.user,
          features: { menu_users: true },
        }),
      );
    const session = await login("admin", "admin");
    expect(session.user.id).toBe("user-admin");
    expect(session.user.permissions).toEqual(["users.read", "users.write", "settings.write"]);
    expect(session.features).toEqual({ menu_users: true });
    expect(requireAuthorization(fetchMock.mock.calls[0][1])).toBeNull(); // login is not authed
    expect(requireBody(fetchMock.mock.calls[0][1])).toEqual({ username: "admin", password: "admin" });
    expect(String(fetchMock.mock.calls[1][0])).toContain("/api/accounts/me");
    expect(requireAuthorization(fetchMock.mock.calls[1][1])).toBe("Bearer a1");
  });

  it("fetchMe keeps permissions from /me for $context expressions (A-006 R-001)", async () => {
    setAccessToken("a1");
    setRefreshToken("r1");
    fetchMock.mockResolvedValueOnce(
      jsonResponse({
        user: {
          id: "user-admin",
          name: "Admin",
          roles: ["admin"],
          permissions: ["users.read", "roles.assign"],
        },
        features: { menu_users: true },
      }),
    );
    const session = await fetchMe();
    expect(session.user.permissions).toEqual(["users.read", "roles.assign"]);
  });

  it("login maps a 401 to INVALID_CREDENTIALS", async () => {
    fetchMock.mockResolvedValueOnce(jsonResponse({ error: "UNAUTHORIZED" }, 401));
    await expect(login("admin", "wrong")).rejects.toMatchObject({ code: "INVALID_CREDENTIALS" });
  });

  it("authFetch attaches the Bearer access token", async () => {
    setAccessToken("access-1");
    fetchMock.mockResolvedValueOnce(jsonResponse({ ok: true }));
    await authFetch("/api/users");
    const init = fetchMock.mock.calls[0][1] as RequestInit;
    expect((init.headers as Headers).get("Authorization")).toBe("Bearer access-1");
  });

  it("authFetch refreshes once on 401 and retries, without notifying auth loss", async () => {
    setAccessToken("expired");
    setRefreshToken("refresh-1");
    const lost = vi.fn();
    setAuthLostListener(lost);

    // First request: 401 (expired access). Refresh endpoint: 200 with new pair.
    // Retried request: 200.
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401))
      .mockImplementationOnce(async (url: RequestInfo) => {
        if (String(url).endsWith("/api/auth/refresh")) {
          return jsonResponse({ accessToken: "access-2", refreshToken: "refresh-2" });
        }
        return jsonResponse({});
      })
      .mockResolvedValueOnce(jsonResponse({ ok: true }));

    const res = await authFetch("/api/users");
    expect(res.ok).toBe(true);
    expect(lost).not.toHaveBeenCalled();
    // Refresh stored the rotated pair.
    expect(requireAuthorization(fetchMock.mock.calls[1][1])).toBeNull(); // refresh is not authed
    expect(requireAuthorization(fetchMock.mock.calls[2][1])).toBe("Bearer access-2");
  });

  it("authFetch notifies auth loss when the refresh fails", async () => {
    setAccessToken("expired");
    setRefreshToken("refresh-1");
    const lost = vi.fn();
    setAuthLostListener(lost);

    fetchMock
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401))
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHORIZED" }, 401));

    const res = await authFetch("/api/users");
    expect(res.status).toBe(401);
    expect(lost).toHaveBeenCalledTimes(1);
  });

  it("authFetch clears the session and notifies auth loss when the retry is still 401", async () => {
    setAccessToken("expired");
    setRefreshToken("refresh-1");
    const lost = vi.fn();
    setAuthLostListener(lost);

    // Original request 401 → refresh succeeds with a rotated pair → the retried
    // request is still 401: the fresh token was rejected, so the session is gone.
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401))
      .mockResolvedValueOnce(jsonResponse({ accessToken: "access-2", refreshToken: "refresh-2" }))
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401));

    const res = await authFetch("/api/users");
    expect(res.status).toBe(401);
    expect(lost).toHaveBeenCalledTimes(1);
    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
  });

  it("login rolls back the stored tokens and rejects when /me fails", async () => {
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ accessToken: "a1", refreshToken: "r1", ...SESSION }))
      .mockResolvedValueOnce(jsonResponse({ error: "INTERNAL" }, 500));

    await expect(login("admin", "admin")).rejects.toMatchObject({ code: "ME_FAILED" });
    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
  });

  it("login rejects when /me 401s and the refresh also fails, without keeping tokens", async () => {
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ accessToken: "a1", refreshToken: "r1", ...SESSION }))
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401))
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHORIZED" }, 401));

    await expect(login("admin", "admin")).rejects.toMatchObject({ code: "ME_FAILED" });
    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
  });

  it("restoreSession rotates a stored refresh and returns the session", async () => {
    setRefreshToken("refresh-1");
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ accessToken: "access-2", refreshToken: "refresh-2" }))
      .mockResolvedValueOnce(jsonResponse(SESSION));

    const session = await restoreSession();
    expect(session?.user.id).toBe("user-admin");
    expect(fetchMock.mock.calls[1][0]).toBe("/api/accounts/me");
  });

  it("restoreSession returns null without a refresh token", async () => {
    const session = await restoreSession();
    expect(session).toBeNull();
  });

  it("fetchMe throws ME_FAILED on a 401 when refresh also fails", async () => {
    setAccessToken("expired");
    setRefreshToken("refresh-1");
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401))
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHORIZED" }, 401));
    await expect(fetchMe()).rejects.toMatchObject({ code: "ME_FAILED" });
  });

  it("logout revokes the refresh token and clears local state", async () => {
    setAccessToken("access-1");
    setRefreshToken("refresh-1");
    fetchMock.mockResolvedValueOnce(emptyResponse(204));

    await logout();
    expect(requireBody(fetchMock.mock.calls[0][1])).toEqual({ refreshToken: "refresh-1" });
  });

  it("logout before a 401 refresh prevents any further refresh (D-001 P2)", async () => {
    setAccessToken("access-1");
    setRefreshToken("refresh-1");
    fetchMock.mockResolvedValueOnce(emptyResponse(204));

    await logout();
    // After logout there is no refresh token, so a later 401 must NOT hit the
    // refresh endpoint — the session is already closed.
    fetchMock.mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401));
    const res = await authFetch("/api/users");
    expect(res.status).toBe(401);
    const refreshCalls = fetchMock.mock.calls.filter(([url]) =>
      String(url).includes("/api/auth/refresh"),
    );
    expect(refreshCalls).toHaveLength(0);
    expect(getAccessToken()).toBeNull();
    expect(getRefreshToken()).toBeNull();
  });

  it("an in-flight refresh resolving after logout must not write tokens back (D-001 P2)", async () => {
    setAccessToken("expired");
    setRefreshToken("refresh-1");

    // First request 401s; the refresh request is kept pending.
    let resolveRefresh!: (response: Response) => void;
    const pendingRefresh = new Promise<Response>((resolve) => {
      resolveRefresh = resolve;
    });
    fetchMock
      .mockResolvedValueOnce(jsonResponse({ error: "UNAUTHENTICATED" }, 401))
      .mockImplementationOnce(() => pendingRefresh);

    const fetchPromise = authFetch("/api/users");
    await vi.waitFor(() => {
      expect(fetchMock.mock.calls.some(([url]) => String(url).includes("/api/auth/refresh"))).toBe(
        true,
      );
    });

    // Logout closes the session while the rotation is in flight.
    fetchMock.mockResolvedValueOnce(emptyResponse(204));
    await logout();

    // The stale rotation resolves with a fresh pair — it must be discarded.
    resolveRefresh(jsonResponse({ accessToken: "access-2", refreshToken: "refresh-2" }));
    await fetchPromise;
    await vi.waitFor(() => {
      expect(getAccessToken()).toBeNull();
      expect(getRefreshToken()).toBeNull();
    });
  });
});

function requireAuthorization(init: RequestInit | undefined): string | null {
  return (init?.headers as Headers | undefined)?.get("Authorization") ?? null;
}

function requireBody(init: RequestInit | undefined): unknown {
  return init?.body !== undefined ? JSON.parse(String(init.body)) : undefined;
}
