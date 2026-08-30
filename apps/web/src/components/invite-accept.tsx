import { useCallback, useState, type FormEvent } from "react";

import { AuthError, inviteAccept } from "@/account/auth-client";
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

/**
 * Maps a stable invitation error code to a catalog key (workspace-019 R3 ·
 * GOAL-004 D-001 §3).
 */
function inviteErrorKey(code: string): string {
  switch (code) {
    case "LOGIN_NETWORK":
      return "login.error.network";
    case "INVITE_INVALID":
      return "invite.error.invalid";
    case "INVITE_ROLE_GONE":
      return "invite.error.roleGone";
    case "USERNAME_TAKEN":
      return "error.usernameTaken";
    case "INVALID_PASSWORD":
      return "login.recovery.error.passwordPolicy";
    default:
      return "login.error.generic";
  }
}

/**
 * Public invitation acceptance surface (workspace-019 R3 · GOAL-004 C4):
 * mounted at /invite/accept?token=… for unauthenticated visitors. Success
 * returns WITHOUT tokens — the new user signs in with their chosen password.
 */
export function InviteAcceptPage() {
  const t = useTranslate();
  const [token] = useState(() => {
    try {
      const params = new URLSearchParams(window.location.search);
      const value = params.get("token") ?? "";
      // W15 F-005 (GOAL-016 A-001): the invitation token is a one-time bearer
      // with a multi-day TTL — scrub it from the address bar and history the
      // moment it is read, so it never lingers in the URL, screenshots or
      // same-origin history after the visitor lands.
      if (value !== "") {
        const url = new URL(window.location.href);
        url.searchParams.delete("token");
        window.history.replaceState(window.history.state, "", url);
      }
      return value;
    } catch {
      return "";
    }
  });
  const [username, setUsername] = useState("");
  const [name, setName] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [done, setDone] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = useCallback(
    async (event: FormEvent) => {
      event.preventDefault();
      if (submitting || token.trim() === "" || username.trim() === "" || password === "") {
        return;
      }
      setSubmitting(true);
      setError(null);
      try {
        await inviteAccept({ token: token.trim(), username: username.trim(), name: name.trim(), password });
        setDone(true);
      } catch (err: unknown) {
        const code = err instanceof AuthError ? err.code : "LOGIN_UNKNOWN";
        setError(t(inviteErrorKey(code)));
      } finally {
        setSubmitting(false);
      }
    },
    [name, password, submitting, t, token, username],
  );

  return (
    <div className="relative flex min-h-screen items-center justify-center bg-background px-4 text-foreground">
      <Card className="w-full max-w-sm shadow-md" data-invite-accept-surface>
        {done ? (
          <form
            onSubmit={(event) => {
              event.preventDefault();
              window.location.href = "/";
            }}
            aria-label={t("invite.title")}
          >
            <CardHeader className="space-y-1 pb-4">
              <CardTitle className="text-2xl tracking-tight" data-invite-done-title>
                {t("invite.title")}
              </CardTitle>
              <CardDescription>{t("invite.success")}</CardDescription>
            </CardHeader>
            <CardContent />
            <CardFooter>
              <Button type="submit" className="w-full">
                {t("invite.goSignIn")}
              </Button>
            </CardFooter>
          </form>
        ) : (
          <form onSubmit={submit} aria-label={t("invite.title")}>
            <CardHeader className="space-y-1 pb-4">
              <CardTitle className="text-2xl tracking-tight">{t("invite.title")}</CardTitle>
              <CardDescription>{t("invite.description")}</CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="space-y-2">
                <Label htmlFor="inviteUsername">{t("schema.users.field.username")}</Label>
                <Input
                  id="inviteUsername"
                  autoComplete="username"
                  value={username}
                  onChange={(event) => setUsername(event.target.value)}
                />
              </div>
              <div className="space-y-2">
                <Label htmlFor="inviteName">{t("schema.users.field.name")}</Label>
                <Input id="inviteName" value={name} onChange={(event) => setName(event.target.value)} />
              </div>
              <div className="space-y-2">
                <Label htmlFor="invitePassword">{t("login.recovery.newPassword")}</Label>
                <Input
                  id="invitePassword"
                  type="password"
                  autoComplete="new-password"
                  value={password}
                  onChange={(event) => setPassword(event.target.value)}
                />
              </div>
              {error !== null ? (
                <p role="alert" data-invite-error className="text-sm text-destructive">
                  {error}
                </p>
              ) : null}
            </CardContent>
            <CardFooter className="flex-col gap-3">
              <Button
                type="submit"
                data-invite-submit
                disabled={submitting || username.trim() === "" || password === ""}
                className="w-full"
              >
                {submitting ? t("login.signingIn") : t("invite.activate")}
              </Button>
            </CardFooter>
          </form>
        )}
      </Card>
    </div>
  );
}
