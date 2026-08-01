/**
 * R2 token storage (GOAL-005 D-002): the short-lived JWT access token lives in
 * memory only; the opaque refresh token persists in localStorage so the session
 * survives a page reload and is rotated on every refresh. localStorage is the
 * user-accepted XSS trade-off (D-002), mitigated by short access TTL, server-side
 * revocation and HTTPS.
 */

const REFRESH_KEY = "schema-ui.refreshToken";

// Access token is intentionally NOT persisted (memory only).
let accessToken: string | null = null;

export function getAccessToken(): string | null {
  return accessToken;
}

export function setAccessToken(token: string | null): void {
  accessToken = token;
}

export function getRefreshToken(): string | null {
  try {
    return window.localStorage.getItem(REFRESH_KEY);
  } catch {
    return null; // storage unavailable (e.g. privacy mode): treat as no session
  }
}

export function setRefreshToken(token: string | null): void {
  try {
    if (token === null) {
      window.localStorage.removeItem(REFRESH_KEY);
    } else {
      window.localStorage.setItem(REFRESH_KEY, token);
    }
  } catch {
    // ignore storage failures; the in-memory access token still works this page
  }
}

export function clearTokens(): void {
  accessToken = null;
  try {
    window.localStorage.removeItem(REFRESH_KEY);
  } catch {
    // ignore
  }
}

/** Whether a session may exist (an access token in memory or a refresh token). */
export function hasSession(): boolean {
  return accessToken !== null || getRefreshToken() !== null;
}
