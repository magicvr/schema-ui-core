import { createContext, useCallback, useContext, useEffect, useMemo, useState, type ReactNode } from "react";

import {
  authFetch,
  login as loginRequest,
  logout as logoutRequest,
  restoreSession,
  setAuthLostListener,
  type AuthSession,
} from "@/account/auth-client";

export type AuthStatus = "loading" | "authenticated" | "unauthenticated";

export interface AuthContextValue {
  status: AuthStatus;
  session: AuthSession | null;
  user: AuthSession["user"] | null;
  /** Authenticates and transitions to the shell. */
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
      if (!cancelled) {
        setSession(null);
        setStatus("unauthenticated");
      }
    });
    restoreSession()
      .then((restored) => {
        if (cancelled) {
          return;
        }
        if (restored !== null) {
          setSession(restored);
          setStatus("authenticated");
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
    const next = await loginRequest(username, password);
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
