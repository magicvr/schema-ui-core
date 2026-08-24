// Account email identity card (workspace-018 R3 · GOAL-004 D-001 §4):
// the minimal binding surface inside the schema-driven account page —
// status badge (unbound / pending / verified), bind form, code entry and
// resend. Consumes the admin.account self-service email API through
// authFetch; delivery goes through the composed kernel.MailSender server-
// side. Registered as custom component "email-identity" (mfa-manager
// precedent, GOAL-018).
import { useCallback, useEffect, useState, type FormEvent } from "react";

import { AuthError, authFetch } from "@/account/auth-client";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { useTranslate } from "@/i18n/runtime";
import {
  registerCustomComponent,
  type CustomComponentProps,
} from "@/renderer/custom-components";

interface ProfileIdentity {
  email: string | null;
  emailStatus: string | null;
}

function emailErrorKey(code: string): string {
  switch (code) {
    case "EMAIL_TAKEN":
      return "schema.account.email.errTaken";
    case "EMAIL_CODE_INVALID":
    case "EMAIL_NOT_PENDING":
      return "schema.account.email.errCode";
    case "EMAIL_CODE_EXPIRED":
      return "schema.account.email.errExpired";
    case "EMAIL_RESEND_COOLDOWN":
      return "schema.account.email.errCooldown";
    default:
      return "schema.account.email.errorGeneric";
  }
}

async function getJSON<T>(url: string): Promise<T> {
  const response = await authFetch(url);
  if (!response.ok) {
    throw new AuthError("UNKNOWN", "GET failed: " + url, response.status);
  }
  return (await response.json()) as T;
}

async function postJSON<T>(url: string, body: unknown): Promise<T> {
  const response = await authFetch(url, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  if (!response.ok) {
    let code = "UNKNOWN";
    try {
      const parsed = (await response.json()) as { error?: string };
      code = parsed.error ?? "UNKNOWN";
    } catch {
      // non-JSON error body
    }
    throw new AuthError(code, "POST failed: " + url, response.status);
  }
  return (await response.json()) as T;
}

export function EmailIdentityCard(_props: CustomComponentProps) {
  const t = useTranslate();
  const [identity, setIdentity] = useState<ProfileIdentity | null>(null);
  const [unavailable, setUnavailable] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [busy, setBusy] = useState(false);
  const [address, setAddress] = useState("");
  const [code, setCode] = useState("");

  const reload = useCallback(async () => {
    try {
      const profile = await getJSON<ProfileIdentity>("/api/account/profile");
      setIdentity({
        email: profile.email ?? null,
        emailStatus: profile.emailStatus ?? null,
      });
      setUnavailable(false);
    } catch {
      setUnavailable(true);
    }
  }, []);

  useEffect(() => {
    void reload();
  }, [reload]);

  const run = useCallback(
    async (action: () => Promise<unknown>) => {
      setBusy(true);
      setError(null);
      try {
        await action();
        setCode("");
        await reload();
      } catch (err) {
        if (err instanceof AuthError) {
          setError(t(emailErrorKey(err.code)));
        } else {
          setError(t("schema.account.email.errorGeneric"));
        }
      } finally {
        setBusy(false);
      }
    },
    [reload, t],
  );

  const onBind = useCallback(
    (event: FormEvent) => {
      event.preventDefault();
      void run(() => postJSON("/api/account/email/bind", { email: address }));
    },
    [run, address],
  );

  const onVerify = useCallback(
    (event: FormEvent) => {
      event.preventDefault();
      void run(() => postJSON("/api/account/email/verify", { code }));
    },
    [run, code],
  );

  if (unavailable) {
    return (
      <div className="space-y-2">
        <p>{t("schema.account.email.unavailable")}</p>
      </div>
    );
  }

  const status = identity?.emailStatus ?? null;
  const email = identity?.email ?? null;

  return (
    <div className="space-y-3" data-testid="email-identity-card">
      <h3 className="text-sm font-semibold">{t("schema.account.email.title")}</h3>
      <p aria-live="polite">
        {status === "verified" && t("schema.account.email.statusVerified")}
        {status === "pending" && t("schema.account.email.statusPending")}
        {status === null && t("schema.account.email.statusUnbound")}
      </p>
      {error !== null && <p role="alert">{error}</p>}
      {status !== "verified" && (
        <form onSubmit={onBind} className="space-y-2" data-testid="email-bind-form">
          <Label htmlFor="email-identity-address">{t("schema.account.email.addressLabel")}</Label>
          <Input
            id="email-identity-address"
            type="email"
            value={status === "pending" && email ? email : address}
            onChange={(event) => setAddress(event.target.value)}
            readOnly={status === "pending"}
          />
          <Button type="submit" disabled={busy || address.trim() === ""}>
            {t("schema.account.email.bindSubmit")}
          </Button>
        </form>
      )}
      {status === "pending" && (
        <form onSubmit={onVerify} className="space-y-2" data-testid="email-verify-form">
          <Label htmlFor="email-identity-code">{t("schema.account.email.codeLabel")}</Label>
          <Input
            id="email-identity-code"
            inputMode="numeric"
            maxLength={6}
            value={code}
            onChange={(event) => setCode(event.target.value)}
          />
          <div className="flex gap-2">
            <Button type="submit" disabled={busy || code.length !== 6}>
              {t("schema.account.email.verifySubmit")}
            </Button>
            <Button
              type="button"
              variant="outline"
              disabled={busy}
              onClick={() => void run(() => postJSON("/api/account/email/resend", {}))}
            >
              {t("schema.account.email.resendSubmit")}
            </Button>
          </div>
        </form>
      )}
    </div>
  );
}

registerCustomComponent("email-identity", EmailIdentityCard);
