import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";

import {
  authFetch,
  AuthError,
  fetchMe,
  isLoginMFARequired,
  login as loginRequest,
  logout as logoutRequest,
  mfaVerify as mfaVerifyRequest,
  restoreSession,
  setAuthLostListener,
  type AuthSession,
} from "@/account/auth-client";
import {
  ACCOUNT_PROFILE_NAMESPACE,
  subscribeToConfigChanges,
} from "@/app/config-events";
import type { SessionAdapterState } from "@/host/boot";
import {
  applyReturnIntentNavigation,
  buildQueryString,
  captureReturnIntent,
  takeReturnIntent,
} from "@/host/return-intent";

export type AuthStatus = SessionAdapterState;

export interface AuthContextValue {
  status: AuthStatus;
  session: AuthSession | null;
  user: AuthSession["user"] | null;
  /** Authenticates, restores a captured return intent, and transitions to the
   * shell. When the account requires a second factor (S-10 · GOAL-017 D-002
   * §3) resolveMFA is invoked with the one-time proof and must return the
   * TOTP code (or recovery code) the user entered. */
  login: (
    username: string,
    password: string,
    captcha?: import("@/account/auth-client").LoginCaptcha,
    resolveMFA?: (proof: string) => Promise<{ code: string; recoveryCode?: string }>,
  ) => Promise<void>;
  /** Revokes the session and transitions to the login page. */
  logout: () => Promise<void>;
  /** Auth-aware fetch: attaches Bearer and refreshes once on 401. */
  authFetch: typeof fetch;
  /**
   * Re-resolves /me into the session (best-effort; keeps the current session
   * on failure). W13 T-05 follow-up: the account profile save publishes the
   * account.profile config-change event and the provider refreshes itself so
   * the shell header (avatar / display name) updates without a reload.
   */
  refreshSession: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [status, setStatus] = useState<AuthStatus>("loading");
  const [session, setSession] = useState<AuthSession | null>(null);

  useEffect(() => {
    let cancelled = false;
    setAuthLostListener(() => {
      // Mid-session session loss (refresh rotation failed / revoked): the
      // adapter reports reauth-required (ADR-0035 D4). Capture the current
      // location first so a successful re-login can restore it (ADR-0036 D6).
      if (!cancelled) {
        captureReturnIntent();
        setSession(null);
        setStatus("reauth-required");
      }
    });
    restoreSession()
      .then((restored) => {
        if (cancelled) {
          return;
        }
        if (restored.kind === "session") {
          setSession(restored.session);
          setStatus("authenticated");
        } else if (restored.kind === "reauth") {
          // A refresh token existed but the session is gone: reauth-required
          // terminal, not anonymous (ADR-0035 D4/D7).
          setSession(null);
          setStatus("reauth-required");
        } else {
          setSession(null);
          setStatus("unauthenticated");
        }
      })
      .catch(() => {
        if (!cancelled) {
          setSession(null);
          setStatus("unauthenticated");
        }
      });
    return () => {
      cancelled = true;
      setAuthLostListener(null);
    };
  }, []);

  // W13 T-05 follow-up: re-resolves the identity snapshot. Best-effort —
  // a failed /me (transient network / revoked token) keeps the current
  // session; the normal auth cycle handles real session loss.
  const refreshSession = useCallback(async () => {
    try {
      const next = await fetchMe();
      setSession(next);
    } catch {
      // keep the current session
    }
  }, []);

  // The account profile save publishes the account.profile config-change
  // event (X-Schema-UI-Config-Changed header on PATCH /api/account/profile);
  // refresh the session so the shell header avatar/name update immediately.
  useEffect(() => {
    return subscribeToConfigChanges(ACCOUNT_PROFILE_NAMESPACE, () => {
      void refreshSession();
    });
  }, [refreshSession]);

  const login = useCallback(async (
    username: string,
    password: string,
    captcha?: import("@/account/auth-client").LoginCaptcha,
    resolveMFA?: (proof: string) => Promise<{ code: string; recoveryCode?: string }>,
  ) => {
    let next: AuthSession;
    try {
      // S-10 (GOAL-017 D-002 §3): the first factor may return a second-factor
      // proof instead of tokens — resolve it through the UI callback, then
      // complete the login with /api/auth/mfa/verify.
      const first = await loginRequest(username, password, captcha);
      if (isLoginMFARequired(first)) {
        if (resolveMFA === undefined) {
          throw new AuthError("MFA_REQUIRED", "second factor required", 401);
        }
        const second = await resolveMFA(first.mfaProof);
        next = await mfaVerifyRequest(first.mfaProof, second.code, second.recoveryCode);
      } else {
        next = first;
      }
    } catch (err: unknown) {
      // GOAL-004 S4-6: 423 is the account-lock terminal (ADR-0035 D4/D7
      // locked state) — the login surface surfaces it as an error, and the
      // failure path below keeps the session adapter in sync with the server
      // terminal (locked until the window expires).
      if (err instanceof AuthError && err.status === 423) {
        setStatus("locked");
      }
      throw err;
    }
    // Consume a captured return intent (single-use, validated by the pinned
    // validator) and restore the location BEFORE the shell mounts so its
    // initial route resolution picks the restored path up.
    const intent = takeReturnIntent();
    if (intent !== null) {
      const target = intent.path + buildQueryString(intent.query);
      const current = window.location.pathname + window.location.search;
      if (target !== current) {
        applyReturnIntentNavigation(target);
      }
    }
    setSession(next);
    setStatus("authenticated");
  }, []);

  const logout = useCallback(async () => {
    await logoutRequest();
    setSession(null);
    setStatus("unauthenticated");
  }, []);

  const value = useMemo<AuthContextValue>(
    () => ({
      status,
      session,
      user: session?.user ?? null,
      login,
      logout,
      authFetch,
      refreshSession,
    }),
    [status, session, login, logout, refreshSession],
  );

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth(): AuthContextValue {
  const ctx = useContext(AuthContext);
  if (ctx === null) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return ctx;
}
