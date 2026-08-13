import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";

import {
  authFetch,
  AuthError,
  login as loginRequest,
  logout as logoutRequest,
  restoreSession,
  setAuthLostListener,
  type AuthSession,
} from "@/account/auth-client";
import type { SessionAdapterState } from "@/host/boot";
import { buildQueryString, captureReturnIntent, takeReturnIntent } from "@/host/return-intent";

export type AuthStatus = SessionAdapterState;

export interface AuthContextValue {
  status: AuthStatus;
  session: AuthSession | null;
  user: AuthSession["user"] | null;
  /** Authenticates, restores a captured return intent, and transitions to the shell. */
  login: (username: string, password: string) => Promise<void>;
  /** Revokes the session and transitions to the login page. */
  logout: () => Promise<void>;
  /** Auth-aware fetch: attaches Bearer and refreshes once on 401. */
  authFetch: typeof fetch;
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

  const login = useCallback(async (username: string, password: string) => {
    let next: AuthSession;
    try {
      next = await loginRequest(username, password);
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
        window.history.replaceState({}, "", target);
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
    }),
    [status, session, login, logout],
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
