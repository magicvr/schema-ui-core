import { useEffect, useState, type FormEvent } from "react";

import { AuthError, type LoginCaptcha } from "@/account/auth-client";
import {
  applyDocumentBranding,
  defaultBranding,
  DEFAULT_SITE_TITLE,
  fetchBranding,
  subscribeToBrandingChanges,
  type Branding,
} from "@/app/branding";
import { LocaleSwitcher } from "@/components/locale-switcher";
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
import { useTranslate } from "@/i18n/runtime";
import { applySystemDefaultTheme } from "@/theme/theme";

/**
 * Maps a stable auth error code to a catalog key (frontend localization
 * floor for the login surface; the server-side catalog lands in S4).
 */
function loginErrorKey(code: string): string {
  switch (code) {
    case "LOGIN_NETWORK":
      return "login.error.network";
    case "INVALID_CREDENTIALS":
      return "login.error.invalidCredentials";
    case "LOGIN_FAILED":
      return "login.error.failed";
    case "LOGIN_MALFORMED":
      return "login.error.malformed";
    case "INVALID_CAPTCHA":
      return "login.error.invalidCaptcha";
    default:
      return "login.error.generic";
  }
}

/**
 * R2 login surface (GOAL-005) + S3 visual upgrade (workspace-006 / D-004 Sign in).
 * Uses design-system Card / Input / Label / Button primitives (not one-off inputs).
 */
export function LoginPage({
  onLogin,
}: {
  onLogin: (username: string, password: string, captcha?: LoginCaptcha) => Promise<void>;
}) {
  const t = useTranslate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [branding, setBranding] = useState<Branding>(() => defaultBranding());
  // S-11 (GOAL-011 D-002 §5): the login page preflights the captcha gate on
  // mount. When enabled it renders the arithmetic question and submits the
  // challenge with the credentials. Preflight failure is fail-open: the form
  // stays captcha-free and the server rejects login with INVALID_CAPTCHA if a
  // challenge is actually required (the error surfaces in the form).
  const [captchaChallenge, setCaptchaChallenge] = useState<{ id: string; question: string } | null>(null);
  const [captchaAnswer, setCaptchaAnswer] = useState("");
  const showSeedHint = import.meta.env.DEV;

  useEffect(() => {
    let cancelled = false;
    const load = () => {
      void fetchBranding().then((next) => {
        if (!cancelled) {
          setBranding(next);
          applyDocumentBranding(next);
          applySystemDefaultTheme(next.defaultTheme);
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

  // S-11 (GOAL-011 D-002 §5): preflight the login captcha gate once on mount.
  useEffect(() => {
    let cancelled = false;
    void fetch("/api/auth/captcha", { headers: { Accept: "application/json" } })
      .then((response) => (response.ok ? response.json() : null))
      .then((body: { enabled?: boolean; challenge?: { id?: string; question?: string } } | null) => {
        if (!cancelled && body?.enabled === true && body.challenge?.id && body.challenge.question) {
          setCaptchaChallenge({ id: body.challenge.id, question: body.challenge.question });
        }
      })
      .catch(() => {
        // fail-open (D-002 §5): the server enforces the gate on login.
      });
    return () => {
      cancelled = true;
    };
  }, []);

  async function handleSubmit(event: FormEvent) {
    event.preventDefault();
    if (submitting) {
      return;
    }
    setSubmitting(true);
    setError(null);
    const captcha: LoginCaptcha | undefined =
      captchaChallenge !== null
        ? { id: captchaChallenge.id, answer: captchaAnswer.trim() }
        : undefined;
    try {
      if (captcha === undefined) {
        await onLogin(username, password);
      } else {
        await onLogin(username, password, captcha);
      }
    } catch (err: unknown) {
      const code = err instanceof AuthError ? err.code : "LOGIN_UNKNOWN";
      // W4 P2-2: login.error.failed carries a `{status}` placeholder; the
      // AuthError keeps the real HTTP status so it is interpolated instead of
      // rendering the literal "{status}".
      const params = err instanceof AuthError && err.status !== undefined ? { status: err.status } : undefined;
      setError(t(loginErrorKey(code), params));
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
          <div className="flex items-center gap-2">
            <LocaleSwitcher />
            <ThemeToggle />
          </div>
        </div>

        <Card className="shadow-md">
          <form onSubmit={handleSubmit} aria-label={t("login.title")}>
            <CardHeader className="space-y-1 pb-4">
              <CardTitle className="text-2xl tracking-tight">{t("login.title")}</CardTitle>
              <CardDescription>{siteTitle}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="username">{t("login.username")}</Label>
                <Input
                  id="username"
                  name="username"
                  autoComplete="username"
                  placeholder={t("login.username")}
                  value={username}
                  onChange={(event) => setUsername(event.target.value)}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">{t("login.password")}</Label>
                <Input
                  id="password"
                  name="password"
                  type="password"
                  autoComplete="current-password"
                  placeholder={t("login.password")}
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                />
              </div>

              {captchaChallenge !== null ? (
                <div className="space-y-2">
                  <Label htmlFor="captchaAnswer">{t("login.captchaQuestion")}</Label>
                  <p className="text-sm font-medium" data-captcha-question>
                    {captchaChallenge.question}
                  </p>
                  <Input
                    id="captchaAnswer"
                    name="captchaAnswer"
                    autoComplete="off"
                    placeholder={t("login.captchaAnswer")}
                    value={captchaAnswer}
                    onChange={(event) => setCaptchaAnswer(event.target.value)}
                  />
                </div>
              ) : null}

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
                {submitting ? t("login.signingIn") : t("login.signIn")}
              </Button>
            </CardFooter>
          </form>
        </Card>

        {showSeedHint ? (
          <p className="mt-4 text-center text-xs text-muted-foreground">
            {t("login.seedHint")} <code className="font-mono">admin / admin</code>
          </p>
        ) : null}
      </div>
    </div>
  );
}
