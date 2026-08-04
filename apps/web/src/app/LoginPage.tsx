import { useEffect, useState, type FormEvent } from "react";

import { AuthError } from "@/account/auth-client";
import {
  applyDocumentBranding,
  DEFAULT_SITE_TITLE,
  fetchBranding,
  subscribeToBrandingChanges,
  type Branding,
} from "@/app/branding";
import { ThemeToggle } from "@/components/theme-toggle";
import { Button } from "@/components/ui/button";

/**
 * R2 login surface (GOAL-005): shown when the session is unauthenticated. On a
 * successful submit the AuthProvider flips to authenticated and the shell
 * renders. Fail-closed: any non-success surfaces a stable error message.
 *
 * A-002 F-002-004 (GOAL-009 S5): the local seed credential hint only renders in
 * development builds; production must not advertise `admin / admin`.
 *
 * GOAL-013: shows site title + optional logo from public GET /api/branding.
 */
export function LoginPage({ onLogin }: { onLogin: (username: string, password: string) => Promise<void> }) {
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [branding, setBranding] = useState<Branding>({
    siteTitle: DEFAULT_SITE_TITLE,
    logoUrl: "",
  });
  const showSeedHint = import.meta.env.DEV;

  useEffect(() => {
    let cancelled = false;
    const load = () => {
      void fetchBranding().then((next) => {
        if (!cancelled) {
          setBranding(next);
          applyDocumentBranding(next);
        }
      });
    };
    load();
    const unsubscribe = subscribeToBrandingChanges(load);
    return () => {
      cancelled = true;
      unsubscribe();
    };
  }, []);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (submitting) {
      return;
    }
    setSubmitting(true);
    setError(null);
    try {
      await onLogin(username, password);
    } catch (err: unknown) {
      setError(err instanceof AuthError ? err.message : "login failed");
    } finally {
      setSubmitting(false);
    }
  }

  const showLogo = branding.logoUrl !== "";
  const siteTitle = branding.siteTitle || DEFAULT_SITE_TITLE;

  return (
    <div className="flex min-h-screen items-center justify-center bg-background px-4 text-foreground">
      <div className="flex min-w-0 w-full max-w-sm flex-col">
        <div className="mb-6 flex items-center justify-between gap-3">
          <div className="flex min-w-0 items-center gap-2">
            {showLogo ? (
              <img src={branding.logoUrl} alt="" className="size-8 shrink-0 object-contain" />
            ) : null}
            <p className="truncate text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
              {siteTitle}
            </p>
          </div>
          <ThemeToggle />
        </div>
        <form
          className="space-y-4 rounded-md border border-border bg-card p-6"
          onSubmit={handleSubmit}
          aria-label="Sign in"
        >
          <div className="space-y-1">
            <h1 className="text-2xl font-semibold tracking-tight">Sign in</h1>
            <p className="text-sm text-muted-foreground">{siteTitle}</p>
          </div>

          <div className="space-y-2">
            <label htmlFor="username" className="block text-sm font-medium">
              Username
            </label>
            <input
              id="username"
              name="username"
              autoComplete="username"
              placeholder="Username"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              className="h-9 w-full rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>

          <div className="space-y-2">
            <label htmlFor="password" className="block text-sm font-medium">
              Password
            </label>
            <input
              id="password"
              name="password"
              type="password"
              autoComplete="current-password"
              placeholder="Password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              className="h-9 w-full rounded-md border border-input bg-background px-3 text-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
            />
          </div>

          {error !== null ? (
            <p role="alert" className="text-sm text-destructive">
              {error}
            </p>
          ) : null}

          <Button type="submit" disabled={submitting || username === "" || password === ""} className="w-full">
            {submitting ? "Signing in…" : "Sign in"}
          </Button>
        </form>
        {showSeedHint ? (
          <p className="mt-4 text-center text-xs text-muted-foreground">
            Local development seed: <code className="font-mono">admin / admin</code>
          </p>
        ) : null}
      </div>
    </div>
  );
}
