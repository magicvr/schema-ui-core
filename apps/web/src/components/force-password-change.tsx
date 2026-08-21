// Forced initial-password change screen (W16-F01): rendered when the signed-in
// user still has mustChangePassword=true. The backend restricts business APIs
// to the password/profile surfaces, so this form is the only entry into the
// app until the initial/reset password is replaced.
import { useState, type FormEvent } from "react";

import { useAuth } from "@/account/AuthContext";
import { setAccessToken, setRefreshToken } from "@/account/tokens";
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

interface PasswordChangeResponse {
  accessToken?: string;
  refreshToken?: string;
}

export function ForcePasswordChange() {
  const t = useTranslate();
  const { authFetch, refreshSession, logout } = useAuth();
  const [currentPassword, setCurrentPassword] = useState("");
  const [newPassword, setNewPassword] = useState("");
  const [confirmPassword, setConfirmPassword] = useState("");
  const [busy, setBusy] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const submit = async (event: FormEvent) => {
    event.preventDefault();
    if (busy) {
      return;
    }
    if (newPassword !== confirmPassword) {
      setError(t("forcePasswordChange.mismatch"));
      return;
    }
    setBusy(true);
    setError(null);
    try {
      const response = await authFetch("/api/account/password", {
        method: "POST",
        headers: { "Content-Type": "application/json", Accept: "application/json" },
        body: JSON.stringify({ currentPassword, newPassword }),
      });
      if (!response.ok) {
        setError(t("forcePasswordChange.failed"));
        return;
      }
      if (response.status === 200) {
        const body = (await response.json()) as PasswordChangeResponse;
        if (typeof body.accessToken === "string" && typeof body.refreshToken === "string") {
          setAccessToken(body.accessToken);
          setRefreshToken(body.refreshToken);
        }
      }
      await refreshSession();
    } catch {
      setError(t("forcePasswordChange.failed"));
    } finally {
      setBusy(false);
    }
  };

  return (
    <div
      data-force-password-change
      className="flex min-h-screen items-center justify-center bg-background px-4 text-foreground"
    >
      <Card className="w-full max-w-sm shadow-md">
        <form onSubmit={submit} aria-label={t("forcePasswordChange.title")}>
          <CardHeader>
            <CardTitle>{t("forcePasswordChange.title")}</CardTitle>
            <CardDescription>{t("forcePasswordChange.description")}</CardDescription>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="currentPassword">{t("forcePasswordChange.currentPassword")}</Label>
              <Input
                id="currentPassword"
                type="password"
                autoComplete="current-password"
                value={currentPassword}
                onChange={(event) => setCurrentPassword(event.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="newPassword">{t("forcePasswordChange.newPassword")}</Label>
              <Input
                id="newPassword"
                type="password"
                autoComplete="new-password"
                value={newPassword}
                onChange={(event) => setNewPassword(event.target.value)}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="confirmPassword">{t("forcePasswordChange.confirmPassword")}</Label>
              <Input
                id="confirmPassword"
                type="password"
                autoComplete="new-password"
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
              />
            </div>
            {error !== null ? (
              <p role="alert" className="text-sm text-destructive">
                {error}
              </p>
            ) : null}
          </CardContent>
          <CardFooter className="flex-col gap-2">
            <Button
              type="submit"
              disabled={busy || currentPassword === "" || newPassword === "" || confirmPassword === ""}
              className="w-full"
            >
              {busy ? t("forcePasswordChange.saving") : t("forcePasswordChange.submit")}
            </Button>
            <Button type="button" variant="ghost" className="w-full" onClick={() => void logout()}>
              {t("forcePasswordChange.signOut")}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}
