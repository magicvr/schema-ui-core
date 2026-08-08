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
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";

/**
 * R2 login surface (GOAL-005) + S3 visual upgrade (workspace-006 / D-004 Sign in).
 * Uses design-system Card / Input / Label / Button primitives (not one-off inputs).
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
    <div
      data-login-surface="design-system"
      className="relative flex min-h-screen items-center justify-center bg-background px-4 text-foreground"
    >
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(ellipse_at_top,oklch(0.95_0_0),transparent_55%)] dark:bg-[radial-gradient(ellipse_at_top,oklch(0.22_0_0),transparent_55%)]" />
      <div className="relative flex min-w-0 w-full max-w-sm flex-col">
        <div className="mb-6 flex items-center justify-between gap-3">
          <div className="flex min-w-0 items-center gap-2">
            {showLogo ? (
              <img src={branding.logoUrl} alt="" className="size-8 shrink-0 object-contain" />
            ) : (
              <span
                aria-hidden="true"
                className="inline-flex size-8 items-center justify-center rounded-md bg-primary text-xs font-semibold text-primary-foreground"
              >
                {siteTitle.slice(0, 1).toUpperCase()}
              </span>
            )}
            <p className="truncate text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
              {siteTitle}
            </p>
          </div>
          <ThemeToggle />
        </div>

        <Card className="shadow-md">
          <form onSubmit={handleSubmit} aria-label="Sign in">
            <CardHeader className="space-y-1 pb-4">
              <CardTitle className="text-2xl tracking-tight">Sign in</CardTitle>
              <CardDescription>{siteTitle}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="username">Username</Label>
                <Input
                  id="username"
                  name="username"
                  autoComplete="username"
                  placeholder="Username"
                  value={username}
                  onChange={(event) => setUsername(event.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">Password</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  placeholder="Password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                />
              </div>

              {error !== null ? (
                <p role="alert" className="text-sm text-destructive">
                  {error}
                </p>
              ) : null}
            </CardContent>
            <CardFooter className="flex-col gap-3">
              <Button
                type="submit"
                disabled={submitting || username === "" || password === ""}
                className="w-full"
              >
                {submitting ? "Signing in…" : "Sign in"}
              </Button>
            </CardFooter>
          </form>
        </Card>

        {showSeedHint ? (
          <p className="mt-4 text-center text-xs text-muted-foreground">
            Local development seed: <code className="font-mono">admin / admin</code>
          </p>
        ) : null}
      </div>
    </div>
  );
}
