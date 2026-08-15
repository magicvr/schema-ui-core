// MFA personal-center management component (GOAL-018): status, one-time
// enrollment (secret + recovery codes), confirm, disable and recovery
// rotation. Consumes the admin.mfa self-service API through authFetch
// (S-10 · GOAL-017 D-002 §4) and renders inside the schema-driven account
// page via the custom-node registry (renderer/custom-components.ts).
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

interface MFAStatus {
  enabled: boolean;
  enrolledAt: string | null;
}

interface EnrollPayload {
  secretBase32: string;
  otpauthURL: string;
  recoveryCodes: string[];
}

function mfaErrorKey(code: string): string {
  switch (code) {
    case "MFA_INVALID":
      return "error.mfaInvalid";
    case "MFA_PROOF_EXPIRED":
      return "error.mfaProofExpired";
    case "MFA_PROOF_EXHAUSTED":
      return "error.mfaProofExhausted";
    case "MFA_NOT_ENROLLED":
      return "error.mfaNotEnrolled";
    case "MFA_PENDING_ONLY":
      return "error.mfaPendingOnly";
    case "MFA_ALREADY_ACTIVE":
      return "error.mfaAlreadyActive";
    case "INVALID_MFA_BODY":
      return "error.invalidMfaBody";
    default:
      return "schema.account.mfa.errorGeneric";
  }
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
  if (response.status === 204) {
    return undefined as T;
  }
  return (await response.json()) as T;
}

export function MfaManager(_props: CustomComponentProps) {
  const t = useTranslate();
  const [status, setStatus] = useState<MFAStatus | null>(null);
  const [unavailable, setUnavailable] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [enrollPayload, setEnrollPayload] = useState<EnrollPayload | null>(null);
  const [confirmCode, setConfirmCode] = useState("");
  const [disableCode, setDisableCode] = useState("");
  const [rotateRecovery, setRotateRecovery] = useState("");
  const [newRecovery, setNewRecovery] = useState<string[] | null>(null);
  const [busy, setBusy] = useState(false);

  const refresh = useCallback(async () => {
    try {
      const response = await authFetch("/api/mfa/status", {
        headers: { Accept: "application/json" },
      });
      if (!response.ok) {
        setUnavailable(true);
        return;
      }
      const body = (await response.json()) as MFAStatus;
      setStatus(body);
      setUnavailable(false);
    } catch {
      setUnavailable(true);
    }
  }, []);

  useEffect(() => {
    void refresh();
  }, [refresh]);

  const run = async (action: () => Promise<void>) => {
    setBusy(true);
    setError(null);
    try {
      await action();
    } catch (err: unknown) {
      setError(err instanceof AuthError ? t(mfaErrorKey(err.code)) : t("schema.account.mfa.errorGeneric"));
    } finally {
      setBusy(false);
    }
  };

  const enroll = () =>
    run(async () => {
      const payload = await postJSON<EnrollPayload>("/api/mfa/enroll", {});
      setEnrollPayload(payload);
      setNewRecovery(null);
    });

  const confirm = (event: FormEvent) => {
    event.preventDefault();
    void run(async () => {
      await postJSON<void>("/api/mfa/confirm", { code: confirmCode.trim() });
      setEnrollPayload(null);
      setConfirmCode("");
      await refresh();
    });
  };

  const disable = (event: FormEvent) => {
    event.preventDefault();
    void run(async () => {
      await postJSON<void>("/api/mfa/disable", { code: disableCode.trim() });
      setDisableCode("");
      await refresh();
    });
  };

  const rotate = (event: FormEvent) => {
    event.preventDefault();
    void run(async () => {
      const payload = await postJSON<{ recoveryCodes: string[] }>("/api/mfa/recovery/rotate", {
        recoveryCode: rotateRecovery.trim(),
      });
      setNewRecovery(payload.recoveryCodes);
      setRotateRecovery("");
    });
  };

  if (unavailable) {
    return (
      <div className="space-y-2">
        <h3 className="text-base font-semibold">{t("schema.account.mfa.title")}</h3>
        <p className="text-sm text-muted-foreground">{t("schema.account.mfa.unavailable")}</p>
      </div>
    );
  }

  return (
    <div className="space-y-3">
      <h3 className="text-base font-semibold">{t("schema.account.mfa.title")}</h3>
      <p className="text-sm text-muted-foreground">
        {status?.enabled
          ? t("schema.account.mfa.statusEnabled")
          : t("schema.account.mfa.statusDisabled")}
      </p>

      {enrollPayload !== null ? (
        <div className="space-y-2" data-mfa-enroll>
          <Label>{t("schema.account.mfa.secret")}</Label>
          <Input readOnly value={enrollPayload.secretBase32} data-mfa-secret />
          <Label>{t("schema.account.mfa.otpauth")}</Label>
          <Input readOnly value={enrollPayload.otpauthURL} />
          <Label>{t("schema.account.mfa.recoveryCodes")}</Label>
          <textarea
            readOnly
            rows={4}
            className="w-full rounded-md border bg-background px-3 py-2 font-mono text-xs"
            value={enrollPayload.recoveryCodes.join("\n")}
            data-mfa-recovery
          />
          <form onSubmit={confirm} className="space-y-2">
            <Input
              id="mfaConfirmCode"
              placeholder={t("schema.account.mfa.confirmCode")}
              autoComplete="one-time-code"
              value={confirmCode}
              onChange={(event) => setConfirmCode(event.target.value)}
            />
            <Button type="submit" disabled={busy || confirmCode.trim() === ""} className="w-full">
              {t("schema.account.mfa.confirm")}
            </Button>
          </form>
        </div>
      ) : status?.enabled ? (
        <div className="space-y-2" data-mfa-active>
          <form onSubmit={disable} className="space-y-2">
            <Input
              id="mfaDisableCode"
              placeholder={t("schema.account.mfa.disableCode")}
              autoComplete="one-time-code"
              value={disableCode}
              onChange={(event) => setDisableCode(event.target.value)}
            />
            <Button type="submit" variant="destructive" disabled={busy || disableCode.trim() === ""} className="w-full">
              {t("schema.account.mfa.disable")}
            </Button>
          </form>
          <form onSubmit={rotate} className="space-y-2">
            <Input
              id="mfaRotateRecovery"
              placeholder={t("schema.account.mfa.rotateRecovery")}
              autoComplete="off"
              value={rotateRecovery}
              onChange={(event) => setRotateRecovery(event.target.value)}
            />
            <Button type="submit" variant="outline" disabled={busy || rotateRecovery.trim() === ""} className="w-full">
              {t("schema.account.mfa.rotate")}
            </Button>
          </form>
          {newRecovery !== null ? (
            <div className="space-y-1" data-mfa-new-recovery>
              <Label>{t("schema.account.mfa.recoveryCodes")}</Label>
              <textarea
                readOnly
                rows={4}
                className="w-full rounded-md border bg-background px-3 py-2 font-mono text-xs"
                value={newRecovery.join("\n")}
              />
            </div>
          ) : null}
        </div>
      ) : (
        <Button type="button" disabled={busy} onClick={() => void enroll()} className="w-full">
          {t("schema.account.mfa.enroll")}
        </Button>
      )}

      {error !== null ? (
        <p role="alert" className="text-sm text-destructive">
          {error}
        </p>
      ) : null}
    </div>
  );
}

// Self-registration: importing this module registers the component under the
// "mfa-manager" key used by the account page schema (GOAL-018).
registerCustomComponent("mfa-manager", MfaManager);
