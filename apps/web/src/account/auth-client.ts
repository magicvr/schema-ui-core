/**
 * R2 auth client (GOAL-005): login / logout / session-restore / me and an
 * auth-aware fetch wrapper that attaches the Bearer access token and transparently
 * refreshes once on 401. On refresh failure the session is cleared and the
 * registered auth-lost listener fires so the UI can return to the login page.
 *
 * A-002 F-002-003 (GOAL-009 S2): a retry that is still 401 after a successful
 * refresh also clears the session and fires auth-lost (a fresh token that is
 * still rejected means the session itself is gone). A `/me` failure during
 * login rolls the stored tokens back and fails the login instead of silently
 * degrading to an empty-feature authenticated session.
 */
import {
  clearTokens,
  getAccessToken,
  getRefreshToken,
  setAccessToken,
  setRefreshToken,
} from "@/account/tokens";
import { getActiveLocale } from "@/i18n/runtime";

export interface AuthUser {
  id: string;
  name: string;
  roles: string[];
  /**
   * Permission keys resolved from persisted RBAC at identity load (`/me`).
   * Required by Schema expressions (`$context.user.permissions contains "…"`).
   * Optional on the type only because login token payloads may omit it until
   * `/me` completes; the session returned by login/restore always carries it
   * when the API provides it.
   */
  permissions?: string[];
}

export interface AuthSession {
  user: AuthUser;
  features: Record<string, boolean>;
}

/** Error with a stable code so the UI can branch (e.g. INVALID_CREDENTIALS). */
export class AuthError extends Error {
  code: string;
  /** Optional HTTP status carried for errors that localize a {status} param. */
  status?: number;
  constructor(code: string, message: string, status?: number) {
    super(message);
    this.name = "AuthError";
    this.code = code;
    if (status !== undefined) {
      this.status = status;
    }
  }
}

const LOGIN_URL = "/api/auth/login";
const REFRESH_URL = "/api/auth/refresh";
const LOGOUT_URL = "/api/auth/logout";
const ME_URL = "/api/accounts/me";

const AUTH_ENDPOINTS = new Set([LOGIN_URL, REFRESH_URL, LOGOUT_URL]);

// Notified when the session is lost (refresh failed / revoked) so the UI can
// transition to the unauthenticated state. Registered by AuthContext.
let onAuthLost: (() => void) | null = null;

export function setAuthLostListener(listener: (() => void) | null): void {
  onAuthLost = listener;
}

function jsonHeaders(init?: RequestInit): Headers {
  const headers = new Headers(init?.headers);
  headers.set("Content-Type", "application/json");
  // VP-007 S4: the server negotiates user-visible messages in the active locale.
  headers.set("Accept-Language", getActiveLocale());
  return headers;
}

function postJSON(url: string, body: unknown): Promise<Response> {
  return fetch(url, {
    method: "POST",
    headers: jsonHeaders(),
    body: JSON.stringify(body),
  });
}

/** Exchanges the stored refresh token for a new access/refresh pair (rotation). */
async function refreshAccess(): Promise<boolean> {
  const refresh = getRefreshToken();
  if (!refresh) {
    return false;
  }
  // Concurrent 401s (parallel requests after access-token expiry) must not each
  // rotate the same refresh token: the server revokes it on first use, so a
  // second concurrent rotation would fail and clear a perfectly valid session.
  // Deduplicate to a single in-flight rotation shared by all callers.
  if (inflightRefresh !== null) {
    return inflightRefresh;
  }
  const rotation = doRefresh(refresh, refreshGeneration).finally(() => {
    inflightRefresh = null;
  });
  inflightRefresh = rotation;
  return rotation;
}

let inflightRefresh: Promise<boolean> | null = null;

// Monotonic session generation (D-001 P2): logout bumps it so an in-flight
// refresh that resolves after logout cannot write a rotated token pair back
// into a session the user just closed.
let refreshGeneration = 0;

function bumpRefreshGeneration(): void {
  refreshGeneration += 1;
}

function isCurrentGeneration(generation: number): boolean {
  return generation === refreshGeneration;
}

async function doRefresh(refresh: string, generation: number): Promise<boolean> {
  if (!isCurrentGeneration(generation)) {
    return false;
  }
  let response: Response;
  try {
    response = await postJSON(REFRESH_URL, { refreshToken: refresh });
  } catch {
    if (isCurrentGeneration(generation)) {
      clearTokens();
    }
    return false;
  }
  if (!isCurrentGeneration(generation)) {
    // The session was closed while the rotation was in flight; never store the
    // rotated pair, never clear tokens the caller already cleared.
    return false;
  }
  if (!response.ok) {
    // Cross-tab rotation (A-002 F-003): another tab may have just won the
    // rotation race and written a fresh refresh token to localStorage while
    // this request was in flight. The old token is now revoked — but the
    // session is not lost. Retry once with the newer token before giving up.
    const current = getRefreshToken();
    if (current !== null && current !== refresh) {
      return doRefresh(current, generation);
    }
    clearTokens();
    return false;
  }
  const body = (await response.json()) as {
    accessToken?: string;
    refreshToken?: string;
  };
  if (!isCurrentGeneration(generation)) {
    return false;
  }
  if (typeof body.accessToken !== "string" || typeof body.refreshToken !== "string") {
    clearTokens();
    return false;
  }
  setAccessToken(body.accessToken);
  setRefreshToken(body.refreshToken);
  return true;
}

function withAuth(init?: RequestInit): RequestInit {
  const access = getAccessToken();
  const headers = new Headers(init?.headers);
  if (access !== null) {
    headers.set("Authorization", `Bearer ${access}`);
  }
  // VP-007 S4: attach the active locale so the server negotiates messages.
  headers.set("Accept-Language", getActiveLocale());
  return { ...init, headers };
}

function isAuthEndpoint(input: RequestInfo | URL): boolean {
  try {
    return AUTH_ENDPOINTS.has(new URL(String(input), window.location.origin).pathname);
  } catch {
    return false;
  }
}

/**
 * Auth-aware fetch: attaches the Bearer access token, and on a 401 (not an auth
 * endpoint) attempts one silent refresh then retries once. If the refresh fails,
 * or the retry is still 401 after a successful refresh, the session is cleared
 * and the auth-lost listener fires (UI → login page).
 */
export async function authFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  let response = await fetch(input, withAuth(init));
  if (response.status === 401 && !isAuthEndpoint(input)) {
    const refreshed = await refreshAccess();
    if (refreshed) {
      response = await fetch(input, withAuth(init));
      if (response.status === 401) {
        clearTokens();
        onAuthLost?.();
      }
    } else {
      clearTokens();
      onAuthLost?.();
    }
  }
  return response;
}

/** Captcha challenge submitted with login (S-11 · GOAL-011 D-002 §5). */
export interface LoginCaptcha {
  id: string;
  answer: string;
}

/** Authenticates with username/password and persists the token pair. */
export async function login(username: string, password: string, captcha?: LoginCaptcha): Promise<AuthSession> {
  let response: Response;
  try {
    response = await postJSON(LOGIN_URL, {
      username,
      password,
      ...(captcha === undefined
        ? {}
        : { captchaId: captcha.id, captchaAnswer: captcha.answer }),
    });
  } catch {
    throw new AuthError("LOGIN_NETWORK", "unable to reach the login service");
  }
  if (response.status === 401) {
    throw new AuthError("INVALID_CREDENTIALS", "invalid username or password");
  }
  // S-11 (GOAL-011 D-002 §2): the gate is on but the challenge was missing,
  // expired or wrong — surface the captcha error so the UI can refresh the
  // challenge and let the user retry (fail-open preflight per D-002 §5).
  if (response.status === 400) {
    let code = "LOGIN_FAILED";
    try {
      const body = (await response.json()) as { error?: string };
      if (body.error === "INVALID_CAPTCHA") {
        code = "INVALID_CAPTCHA";
      }
    } catch {
      // non-JSON error body: keep LOGIN_FAILED
    }
    throw new AuthError(code, "login failed", 400);
  }
  // GOAL-004 S4-6: 423 is the account-lock terminal (ADR-0035 D4/D7 locked
  // state), distinct from 401 credential failure.
  if (response.status === 423) {
    throw new AuthError("ACCOUNT_LOCKED", "account is temporarily locked", 423);
  }
  if (!response.ok) {
    throw new AuthError("LOGIN_FAILED", `login failed: HTTP ${response.status}`, response.status);
  }
  const body = (await response.json()) as {
    accessToken?: string;
    refreshToken?: string;
    user?: AuthUser;
  };
  if (typeof body.accessToken !== "string" || typeof body.refreshToken !== "string" || !body.user) {
    throw new AuthError("LOGIN_MALFORMED", "login response was malformed");
  }
  setAccessToken(body.accessToken);
  setRefreshToken(body.refreshToken);
  // Login token response carries user identity only; menu/feature projection
  // lives on GET /me (GOAL-006 S5). Resolve features the same way restoreSession
  // does so post-login navigation matches a restored session. A /me failure
  // rolls the tokens back and fails the login rather than silently degrading.
  try {
    return await fetchMe();
  } catch (error) {
    clearTokens();
    throw error;
  }
}

/** Revokes the refresh token (best-effort, idempotent) and clears local state. */
export async function logout(): Promise<void> {
  const refresh = getRefreshToken();
  // D-001 P2: bump the generation first so an in-flight refresh can never write
  // a rotated token pair back after logout.
  bumpRefreshGeneration();
  if (refresh !== null) {
    try {
      await postJSON(LOGOUT_URL, { refreshToken: refresh });
    } catch {
      // best-effort: the server revoke may be unreachable; local state is still
      // cleared and the revoked token will fail any later refresh server-side.
    }
  }
  clearTokens();
}

/** Restore outcomes (ADR-0035 D4 normalized adapter input, GOAL-004 S4-2). */
export type RestoreSessionResult =
  | { kind: "none" }
  | { kind: "reauth" }
  | { kind: "session"; session: AuthSession };

/**
 * Restores a session on boot: if a refresh token exists, rotate it for a fresh
 * access/refresh pair, then resolve the identity via /me. No refresh token →
 * `none` (anonymous); rotation or identity resolution failure with an existing
 * refresh token → `reauth` (the adapter reports reauth-required, ADR-0035 D4).
 */
export async function restoreSession(): Promise<RestoreSessionResult> {
  if (getRefreshToken() === null) {
    return { kind: "none" };
  }
  const refreshed = await refreshAccess();
  if (!refreshed) {
    return { kind: "reauth" };
  }
  try {
    const session = await fetchMe();
    return { kind: "session", session };
  } catch {
    return { kind: "reauth" };
  }
}

/** Normalizes a /me user snapshot so permissions are always an array when present. */
function parseAuthUser(raw: unknown): AuthUser | null {
  if (typeof raw !== "object" || raw === null || Array.isArray(raw)) {
    return null;
  }
  const record = raw as Record<string, unknown>;
  if (typeof record.id !== "string" || record.id === "") {
    return null;
  }
  const roles = Array.isArray(record.roles)
    ? record.roles.filter((entry): entry is string => typeof entry === "string")
    : [];
  const permissions = Array.isArray(record.permissions)
    ? record.permissions.filter((entry): entry is string => typeof entry === "string")
    : undefined;
  return {
    id: record.id,
    name: typeof record.name === "string" ? record.name : "",
    roles,
    ...(permissions === undefined ? {} : { permissions }),
  };
}

/** Fetches the current session from /me using the active access token. */
export async function fetchMe(): Promise<AuthSession> {
  const response = await authFetch(ME_URL);
  if (!response.ok) {
    throw new AuthError("ME_FAILED", `session fetch failed: HTTP ${response.status}`);
  }
  const body = (await response.json()) as { user?: unknown; features?: Record<string, boolean> };
  const user = parseAuthUser(body.user);
  if (user === null) {
    throw new AuthError("ME_MALFORMED", "session response was malformed");
  }
  return { user, features: body.features ?? {} };
}
