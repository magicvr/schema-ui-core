import { useCallback, useEffect, useRef, useState, type FormEvent } from "react";
import { Eye, EyeOff } from "lucide-react";

import { AuthError, recoveryComplete, recoveryStart, type LoginCaptcha } from "@/account/auth-client";
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
 * Maps a stable recovery error code to a catalog key (workspace-019 R2 ·
 * GOAL-003 D-001 §2; server codes are the single source of truth).
 */
function recoveryErrorKey(code: string): string {
  switch (code) {
    case "LOGIN_NETWORK":
      return "login.error.network";
    case "EMAIL_RESEND_COOLDOWN":
      return "login.recovery.error.cooldown";
    case "RECOVERY_CODE_INVALID":
      return "login.recovery.error.codeInvalid";
    case "RECOVERY_CODE_EXPIRED":
      return "login.recovery.error.codeExpired";
    case "RECOVERY_SECOND_FACTOR_REQUIRED":
      return "login.recovery.error.secondFactorRequired";
    case "MFA_INVALID":
      return "login.error.mfaInvalid";
    case "INVALID_PASSWORD":
      return "login.recovery.error.passwordPolicy";
    case "EMAIL_SEND_FAILED":
      return "login.recovery.error.sendFailed";
    default:
      return "login.error.generic";
  }
}

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
    case "MFA_INVALID":
      return "login.error.mfaInvalid";
    case "MFA_PROOF_EXPIRED":
      return "login.error.mfaProofExpired";
    case "MFA_PROOF_EXHAUSTED":
      return "login.error.mfaProofExhausted";
    case "MFA_REQUIRED":
      return "login.error.mfaRequired";
    case "LOGIN_CANCELLED":
      return "login.error.generic";
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
  onLogin: (
    username: string,
    password: string,
    captcha?: LoginCaptcha,
    resolveMFA?: (proof: string) => Promise<{ code: string; recoveryCode?: string }>,
  ) => Promise<void>;
}) {
  const t = useTranslate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [showPassword, setShowPassword] = useState(false);
  // W11 · M-02: one-time notice after a successful MFA disable (the server
  // revoked all sessions, so the app signed out locally and landed here).
  const [notice, setNotice] = useState<string | null>(() => {
    try {
      if (sessionStorage.getItem("mfa.disabledNotice") !== null) {
        sessionStorage.removeItem("mfa.disabledNotice");
        return t("login.mfaDisabledNotice");
      }
      if (sessionStorage.getItem("password.changedNotice") !== null) {
        sessionStorage.removeItem("password.changedNotice");
        return t("schema.account.passwordChangedReauth");
      }
    } catch {
      // storage unavailable — skip the notice
    }
    return null;
  });
  const [branding, setBranding] = useState<Branding>(() => defaultBranding());
  // S-11 (GOAL-011 D-002 §5): the login page preflights the captcha gate on
  // mount. When enabled it renders the arithmetic question and submits the
  // challenge with the credentials. Preflight failure is fail-open: the form
  // stays captcha-free and the server rejects login with INVALID_CAPTCHA if a
  // challenge is actually required (the error surfaces in the form).
  const [captchaChallenge, setCaptchaChallenge] = useState<{ id: string; question: string } | null>(null);
  const [captchaAnswer, setCaptchaAnswer] = useState("");
  // S-10 (GOAL-017 D-002 §3): two-step login — the password factor succeeded
  // and the server issued a one-time proof; the code stage resolves the
  // pending promise with the user's TOTP (or recovery) code.
  const [mfaPending, setMfaPending] = useState<string | null>(null);
  const [mfaCode, setMfaCode] = useState("");
  const [mfaRecovery, setMfaRecovery] = useState("");
  const mfaResolverRef = useRef<((v: { code: string; recoveryCode?: string }) => void) | null>(null);
  const mfaCancelRef = useRef<(() => void) | null>(null);
  // workspace-019 R2 (GOAL-003 D-001 §2): two-step self-recovery — step 1
  // requests the email code, step 2 carries code (+ second factor for MFA
  // accounts) and the replacement password. Success returns WITHOUT tokens:
  // the user signs in with the new password.
  const [mode, setMode] = useState<"signin" | "recover">("signin");
  const [recStep, setRecStep] = useState<1 | 2>(1);
  const [recAccount, setRecAccount] = useState("");
  const [recCode, setRecCode] = useState("");
  const [recNewPassword, setRecNewPassword] = useState("");
  const [recSecondFactor, setRecSecondFactor] = useState("");
  const [recRecoveryCode, setRecRecoveryCode] = useState("");
  const [recSecondFactorNeeded, setRecSecondFactorNeeded] = useState(false);
  const [recSubmitting, setRecSubmitting] = useState(false);
  const [recError, setRecError] = useState<string | null>(null);
  const [recDone, setRecDone] = useState(false);
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

  // S-11 (GOAL-011 D-002 §5): preflight the login captcha gate. Refreshes
  // on mount and after an INVALID_CAPTCHA failure so a consumed/expired
  // challenge is replaced without a page reload (grok A-004 F-009).
  const refreshCaptcha = useCallback(() => {
    void fetch("/api/auth/captcha", { headers: { Accept: "application/json" } })
      .then((response) => (response.ok ? response.json() : null))
      .then((body: { enabled?: boolean; challenge?: { id?: string; question?: string } } | null) => {
        if (body?.enabled === true && body.challenge?.id && body.challenge.question) {
          setCaptchaChallenge({ id: body.challenge.id, question: body.challenge.question });
          setCaptchaAnswer("");
        } else {
          setCaptchaChallenge(null);
        }
      })
      .catch(() => {
        // fail-open (D-002 §5): the server enforces the gate on login.
      });
  }, []);

  useEffect(() => {
    refreshCaptcha();
  }, [refreshCaptcha]);

  async function handleSubmit(event: FormEvent) {    event.preventDefault();
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
      // S-10 (GOAL-017 D-002 §3): when the account requires a second factor
      // the login callback pauses here until the code stage resolves.
      const resolveMFA = async (proof: string): Promise<{ code: string; recoveryCode?: string }> => {
        setMfaPending(proof);
        setMfaCode("");
        setMfaRecovery("");
        return new Promise((resolve, reject) => {
          mfaResolverRef.current = resolve;
          mfaCancelRef.current = () => {
            mfaResolverRef.current = null;
            mfaCancelRef.current = null;
            setMfaPending(null);
            setMfaCode("");
            setMfaRecovery("");
            reject(new AuthError("LOGIN_CANCELLED", "mfa cancelled"));
          };
        });
      };
      if (captcha === undefined) {
        await onLogin(username, password, undefined, resolveMFA);
      } else {
        await onLogin(username, password, captcha, resolveMFA);
      }
    } catch (err: unknown) {
      const code = err instanceof AuthError ? err.code : "LOGIN_UNKNOWN";
      // W4 P2-2: login.error.failed carries a `{status}` placeholder; the
      // AuthError keeps the real HTTP status so it is interpolated instead of
      // rendering the literal "{status}".
      const params = err instanceof AuthError && err.status !== undefined ? { status: err.status } : undefined;
      setError(t(loginErrorKey(code), params));
      // S-11 (GOAL-011 · grok A-004 F-009): a rejected captcha was consumed
      // server-side — issue a fresh challenge so the user can retry without a
      // page reload.
      if (code === "INVALID_CAPTCHA") {
        refreshCaptcha();
      }
    } finally {
      setSubmitting(false);
    }
  }

  /** workspace-019 R2: step 1 — request the reset code (enumeration-neutral:
   * the server answers dispatched even when no recovery path exists). */
  async function handleRecoveryStart(event: FormEvent) {
    event.preventDefault();
    if (recSubmitting || recAccount.trim() === "") {
      return;
    }
    setRecSubmitting(true);
    setRecError(null);
    try {
      await recoveryStart(recAccount.trim());
      setRecStep(2);
      setRecSecondFactorNeeded(false);
    } catch (err: unknown) {
      const code = err instanceof AuthError ? err.code : "LOGIN_UNKNOWN";
      setError(null);
      if (code === "EMAIL_RESEND_COOLDOWN") {
        // The code from a previous request may still be valid — jump ahead.
        setRecStep(2);
      } else {
        setRecError(t(recoveryErrorKey(code)));
      }
    } finally {
      setRecSubmitting(false);
    }
  }

  /** workspace-019 R2: step 2 — verify code (+ second factor) and rotate the
   * password. RECOVERY_SECOND_FACTOR_REQUIRED flips the extra fields in. */
  async function handleRecoveryComplete(event: FormEvent) {
    event.preventDefault();
    if (recSubmitting || recCode.trim() === "" || recNewPassword === "") {
      return;
    }
    setRecSubmitting(true);
    setRecError(null);
    try {
      await recoveryComplete({
        account: recAccount.trim(),
        code: recCode.trim(),
        newPassword: recNewPassword,
        ...(recSecondFactor.trim() === "" ? {} : { secondFactorCode: recSecondFactor.trim() }),
        ...(recRecoveryCode.trim() === "" ? {} : { recoveryCode: recRecoveryCode.trim() }),
      });
      setRecDone(true);
    } catch (err: unknown) {
      const code = err instanceof AuthError ? err.code : "LOGIN_UNKNOWN";
      if (code === "RECOVERY_SECOND_FACTOR_REQUIRED") {
        setRecSecondFactorNeeded(true);
        setRecError(t("login.recovery.error.secondFactorRequired"));
        return;
      }
      if (code === "RECOVERY_CODE_EXPIRED") {
        // The challenge is voided/expired — restart at step 1 with the hint.
        setRecStep(1);
        setRecCode("");
        setRecNewPassword("");
        setRecSecondFactor("");
        setRecRecoveryCode("");
        setRecSecondFactorNeeded(false);
        setRecError(t("login.recovery.error.codeExpired"));
        return;
      }
      setRecError(t(recoveryErrorKey(code)));
    } finally {
      setRecSubmitting(false);
    }
  }

  function backToSignIn(): void {
    setMode("signin");
    setRecStep(1);
    setRecAccount("");
    setRecCode("");
    setRecNewPassword("");
    setRecSecondFactor("");
    setRecRecoveryCode("");
    setRecSecondFactorNeeded(false);
    setRecError(null);
    setRecDone(false);
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

        {notice !== null ? (
          <div
            role="status"
            data-login-notice="mfa-disabled"
            className="mb-3 flex items-start justify-between gap-3 rounded-md border border-success/50 bg-success/10 px-3 py-2 text-sm text-success"
          >
            <span>{notice}</span>
            <button
              type="button"
              aria-label={t("feedback.cancel")}
              className="shrink-0 text-success/70 transition-colors hover:text-success"
              onClick={() => setNotice(null)}
            >
              ×
            </button>
          </div>
        ) : null}

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
                <div className="relative">
                  <Input
                    id="password"
                    name="password"
                    type={showPassword ? "text" : "password"}
                    autoComplete="current-password"
                    placeholder={t("login.password")}
                    value={password}
                    onChange={(event) => setPassword(event.target.value)}
                    className="pr-9"
                  />
                  <button
                    type="button"
                    aria-label={showPassword ? t("login.password.hide") : t("login.password.show")}
                    data-password-toggle
                    className="absolute inset-y-0 right-0 flex items-center px-2.5 text-muted-foreground transition-colors hover:text-foreground"
                    onClick={() => setShowPassword((v) => !v)}
                  >
                    {showPassword ? (
                      <EyeOff aria-hidden="true" className="size-4" />
                    ) : (
                      <Eye aria-hidden="true" className="size-4" />
                    )}
                  </button>
                </div>
              </div>

              {captchaChallenge !== null ? (
                <div className="space-y-2">
                  <div className="flex items-center justify-between gap-2">
                    <Label htmlFor="captchaAnswer">{t("login.captchaQuestion")}</Label>
                    <button
                      type="button"
                      data-captcha-refresh
                      className="text-xs text-muted-foreground underline-offset-2 transition-colors hover:text-foreground hover:underline"
                      onClick={refreshCaptcha}
                    >
                      {t("login.captchaRefresh")}
                    </button>
                  </div>
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

              {mfaPending !== null ? (
                <div className="space-y-2" data-mfa-stage>
                  <Label>{t("login.mfa.title")}</Label>
                  <p className="text-sm text-muted-foreground">{t("login.mfa.description")}</p>
                  <Input
                    id="mfaCode"
                    name="mfaCode"
                    autoComplete="one-time-code"
                    inputMode="numeric"
                    placeholder={t("login.mfa.code")}
                    value={mfaCode}
                    onChange={(event) => setMfaCode(event.target.value)}
                  />
                  <Input
                    id="mfaRecovery"
                    name="mfaRecovery"
                    autoComplete="off"
                    placeholder={t("login.mfa.recovery")}
                    value={mfaRecovery}
                    onChange={(event) => setMfaRecovery(event.target.value)}
                  />
                  <Button
                    type="button"
                    variant="outline"
                    className="w-full"
                    onClick={() => {
                      mfaCancelRef.current?.();
                    }}
                  >
                    {t("feedback.cancel")}
                  </Button>
                  <Button
                    type="button"
                    data-mfa-verify="true"
                    disabled={mfaCode.trim() === "" && mfaRecovery.trim() === ""}
                    className="w-full"
                    onClick={() => {
                      const resolve = mfaResolverRef.current;
                      if (resolve !== null) {
                        mfaResolverRef.current = null;
                        setMfaPending(null);
                        resolve({
                          code: mfaCode.trim(),
                          ...(mfaRecovery.trim() === "" ? {} : { recoveryCode: mfaRecovery.trim() }),
                        });
                      }
                    }}
                  >
                    {t("login.mfa.verify")}
                  </Button>
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
              <button
                type="button"
                data-recovery-link
                className="text-xs text-muted-foreground underline-offset-2 transition-colors hover:text-foreground hover:underline"
                onClick={() => {
                  setMode("recover");
                  setRecStep(1);
                  setRecDone(false);
                  setRecError(null);
                }}
              >
                {t("login.recovery.link")}
              </button>
            </CardFooter>
          </form>
        </Card>

        {mode === "recover" ? (
          <Card className="shadow-md" data-recovery-surface>
            {recDone ? (
              <form
                onSubmit={(event) => {
                  event.preventDefault();
                  backToSignIn();
                }}
                aria-label={t("login.recovery.title")}
              >
                <CardHeader className="space-y-1 pb-4">
                  <CardTitle className="text-2xl tracking-tight" data-recovery-done-title>
                    {t("login.recovery.title")}
                  </CardTitle>
                  <CardDescription>{t("login.recovery.success")}</CardDescription>
                </CardHeader>
                <CardContent />
                <CardFooter className="flex-col gap-3">
                  <Button type="submit" className="w-full">
                    {t("login.recovery.back")}
                  </Button>
                </CardFooter>
              </form>
            ) : recStep === 1 ? (
              <form onSubmit={handleRecoveryStart} aria-label={t("login.recovery.title")}>
                <CardHeader className="space-y-1 pb-4">
                  <CardTitle className="text-2xl tracking-tight">{t("login.recovery.title")}</CardTitle>
                  <CardDescription>{t("login.recovery.description")}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="recoveryAccount">{t("login.recovery.account")}</Label>
                    <Input
                      id="recoveryAccount"
                      name="account"
                      autoComplete="username"
                      placeholder={t("login.recovery.accountPlaceholder")}
                      value={recAccount}
                      onChange={(event) => setRecAccount(event.target.value)}
                    />
                  </div>
                  {recError !== null ? (
                    <p role="alert" data-recovery-error className="text-sm text-destructive">
                      {recError}
                    </p>
                  ) : null}
                </CardContent>
                <CardFooter className="flex-col gap-3">
                  <Button
                    type="submit"
                    data-recovery-send
                    disabled={recSubmitting || recAccount.trim() === ""}
                    className="w-full"
                  >
                    {recSubmitting ? t("login.signingIn") : t("login.recovery.sendCode")}
                  </Button>
                  <button
                    type="button"
                    className="text-xs text-muted-foreground underline-offset-2 transition-colors hover:text-foreground hover:underline"
                    onClick={backToSignIn}
                  >
                    {t("login.recovery.back")}
                  </button>
                </CardFooter>
              </form>
            ) : (
              <form onSubmit={handleRecoveryComplete} aria-label={t("login.recovery.title")}>
                <CardHeader className="space-y-1 pb-4">
                  <CardTitle className="text-2xl tracking-tight">{t("login.recovery.title")}</CardTitle>
                  <CardDescription>{t("login.recovery.codeSent")}</CardDescription>
                </CardHeader>
                <CardContent className="space-y-4">
                  <div className="space-y-2">
                    <Label htmlFor="recoveryCode">{t("login.recovery.code")}</Label>
                    <Input
                      id="recoveryCode"
                      name="code"
                      autoComplete="one-time-code"
                      inputMode="numeric"
                      maxLength={6}
                      value={recCode}
                      onChange={(event) => setRecCode(event.target.value)}
                    />
                  </div>
                  {recSecondFactorNeeded ? (
                    <>
                      <div className="space-y-2">
                        <Label htmlFor="recoverySecondFactor">{t("login.recovery.secondFactor")}</Label>
                        <Input
                          id="recoverySecondFactor"
                          name="secondFactorCode"
                          autoComplete="one-time-code"
                          inputMode="numeric"
                          value={recSecondFactor}
                          onChange={(event) => setRecSecondFactor(event.target.value)}
                        />
                      </div>
                      <div className="space-y-2">
                        <Label htmlFor="recoveryMfaCode">{t("login.recovery.recoveryCode")}</Label>
                        <Input
                          id="recoveryMfaCode"
                          name="recoveryCode"
                          autoComplete="off"
                          value={recRecoveryCode}
                          onChange={(event) => setRecRecoveryCode(event.target.value)}
                        />
                      </div>
                    </>
                  ) : null}
                  <div className="space-y-2">
                    <Label htmlFor="recoveryNewPassword">{t("login.recovery.newPassword")}</Label>
                    <Input
                      id="recoveryNewPassword"
                      name="newPassword"
                      type="password"
                      autoComplete="new-password"
                      value={recNewPassword}
                      onChange={(event) => setRecNewPassword(event.target.value)}
                    />
                  </div>
                  {recError !== null ? (
                    <p role="alert" data-recovery-error className="text-sm text-destructive">
                      {recError}
                    </p>
                  ) : null}
                </CardContent>
                <CardFooter className="flex-col gap-3">
                  <Button
                    type="submit"
                    data-recovery-reset
                    disabled={recSubmitting || recCode.trim() === "" || recNewPassword === ""}
                    className="w-full"
                  >
                    {recSubmitting ? t("login.signingIn") : t("login.recovery.complete")}
                  </Button>
                  <button
                    type="button"
                    className="text-xs text-muted-foreground underline-offset-2 transition-colors hover:text-foreground hover:underline"
                    onClick={backToSignIn}
                  >
                    {t("login.recovery.back")}
                  </button>
                </CardFooter>
              </form>
            )}
          </Card>
        ) : null}

        {showSeedHint ? (
          <p className="mt-4 text-center text-xs text-muted-foreground">
            {t("login.seedHint")} <code className="font-mono">admin / admin</code>
          </p>
        ) : null}
      </div>
    </div>
  );
}
