// Telegram channel admin console (GOAL-006 R5, 判据 #5 补做 Admin UI tab,
// user-adjudicated 2026-09-03): the channel.telegram settings page is hosted by
// this custom component — Bot Token / Webhook Secret are edited write-only
// (GET only reports token_set/secret_set booleans; PATCH accepts new values,
// empty keeps current — F-002 / R-005), plus the captured-message counter for
// the mock outbound sink. Secrets never leave the API in plaintext or partial
// masks.
import { useCallback, useEffect, useState } from "react";

import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

interface TelegramSettingsStatus {
  configured: boolean;
  token_set: boolean;
  secret_set: boolean;
  captured_messages_count?: number;
  captured_count?: number;
}

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

export function TelegramAdminTab(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [loadState, setLoadState] = useState<"loading" | "ready" | "error">("loading");
  const [status, setStatus] = useState<TelegramSettingsStatus | null>(null);
  const [tokenInput, setTokenInput] = useState("");
  const [secretInput, setSecretInput] = useState("");
  const [saving, setSaving] = useState(false);
  const [clearing, setClearing] = useState(false);
  const [confirmClear, setConfirmClear] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const loadStatus = useCallback(async () => {
    setLoadState("loading");
    try {
      const response = await fetcher("/api/channel/telegram/settings");
      if (!response.ok) {
        setLoadState("error");
        return;
      }
      setStatus((await response.json()) as TelegramSettingsStatus);
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, [fetcher]);

  useEffect(() => {
    void loadStatus();
  }, [loadStatus]);

  async function extractError(response: Response): Promise<string> {
    try {
      const body = (await response.json()) as { detail?: string; message?: string };
      const text = typeof body.detail === "string" && body.detail !== "" ? body.detail : body.message;
      if (typeof text === "string" && text !== "") {
        return text;
      }
    } catch {
      // fall through to generic message
    }
    return t("schema.telegram.feedback.saveFailed");
  }

  async function save() {
    setSaving(true);
    setFeedback(null);
    try {
      const payload: Record<string, unknown> = {};
      if (tokenInput.trim() !== "") {
        payload.bot_token = tokenInput.trim();
      }
      if (secretInput.trim() !== "") {
        payload.webhook_secret = secretInput.trim();
      }
      const response = await fetcher("/api/channel/telegram/settings", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        setFeedback({ kind: "error", message: await extractError(response) });
        return;
      }
      setStatus((await response.json()) as TelegramSettingsStatus);
      setTokenInput("");
      setSecretInput("");
      setFeedback({ kind: "success", message: t("schema.telegram.feedback.saved") });
    } catch {
      setFeedback({ kind: "error", message: t("schema.telegram.feedback.saveFailed") });
    } finally {
      setSaving(false);
    }
  }

  // R-004 / A-002: an explicit clear action sends empty strings so the admin
  // can disable the bot; an empty input on save means "keep current" instead.
  async function clearSecrets() {
    setClearing(true);
    setFeedback(null);
    try {
      const response = await fetcher("/api/channel/telegram/settings", {
        method: "PATCH",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ bot_token: "", webhook_secret: "" }),
      });
      if (!response.ok) {
        setFeedback({ kind: "error", message: await extractError(response) });
        return;
      }
      setStatus((await response.json()) as TelegramSettingsStatus);
      setTokenInput("");
      setSecretInput("");
      setConfirmClear(false);
      setFeedback({ kind: "success", message: t("schema.telegram.feedback.cleared") });
    } catch {
      setFeedback({ kind: "error", message: t("schema.telegram.feedback.saveFailed") });
    } finally {
      setClearing(false);
    }
  }

  if (loadState === "error") {
    return (
      <section data-telegram-admin-tab className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
        <h2 className="text-sm font-semibold">{t("schema.settings.toolbar.telegram")}</h2>
        <p role="alert" className="text-sm text-destructive">{t("schema.telegram.feedback.loadFailed")}</p>
      </section>
    );
  }

  const fieldLabel = (key: string, id: string) => (
    <label className="text-sm font-medium" htmlFor={id}>
      {t(key)}
    </label>
  );

  return (
    <section data-telegram-admin-tab className="space-y-4 rounded-xl border border-border/70 bg-card/85 p-4">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold">{t("schema.settings.toolbar.telegram")}</h2>
        {status !== null ? (
          <span className="text-xs text-muted-foreground">
            {status.configured
              ? t("schema.telegram.status.configured")
              : t("schema.telegram.status.notConfigured")}
          </span>
        ) : null}
      </div>

      {loadState === "loading" ? <p className="text-sm text-muted-foreground">{t("feedback.loading")}</p> : null}

      <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
        {fieldLabel("schema.telegram.field.botToken", "telegram-bot-token")}
        <input
          id="telegram-bot-token"
          type="password"
          autoComplete="new-password"
          value={tokenInput}
          placeholder={status?.token_set ? t("schema.telegram.secret.keep") : ""}
          onChange={(event) => setTokenInput(event.target.value)}
          className={inputClass}
        />
      </div>
      <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
        {fieldLabel("schema.telegram.field.webhookSecret", "telegram-webhook-secret")}
        <input
          id="telegram-webhook-secret"
          type="password"
          autoComplete="new-password"
          value={secretInput}
          placeholder={status?.secret_set ? t("schema.telegram.secret.keep") : ""}
          onChange={(event) => setSecretInput(event.target.value)}
          className={inputClass}
        />
      </div>

      <div className="flex items-center gap-3">
        <button type="button" disabled={saving || loadState !== "ready"} onClick={() => void save()} className={buttonClass}>
          {saving ? t("feedback.submitting") : t("schema.telegram.action.save")}
        </button>
        {status?.configured ? (
          confirmClear ? (
            <span className="inline-flex items-center gap-2 text-sm">
              <span className="text-destructive">{t("schema.telegram.clear.confirm")}</span>
              <button
                type="button"
                disabled={clearing}
                onClick={() => void clearSecrets()}
                className="rounded-md border border-destructive/60 px-2.5 py-1 text-sm text-destructive hover:bg-destructive/10 disabled:cursor-not-allowed disabled:opacity-50"
              >
                {clearing ? t("feedback.submitting") : t("schema.telegram.clear.confirmAction")}
              </button>
              <button
                type="button"
                disabled={clearing}
                onClick={() => setConfirmClear(false)}
                className="rounded-md border border-input/80 px-2.5 py-1 text-sm hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
              >
                {t("schema.telegram.clear.cancel")}
              </button>
            </span>
          ) : (
            <button
              type="button"
              disabled={saving || clearing}
              onClick={() => setConfirmClear(true)}
              className="rounded-md border border-input/80 px-2.5 py-1 text-sm text-muted-foreground hover:bg-muted hover:text-foreground disabled:cursor-not-allowed disabled:opacity-50"
            >
              {t("schema.telegram.clear.action")}
            </button>
          )
        ) : null}
      </div>

      {status !== null && typeof status.captured_messages_count === "number" ? (
        <p className="text-xs text-muted-foreground">
          {t("schema.telegram.status.captured")} {status.captured_messages_count}
        </p>
      ) : null}

      {feedback !== null ? (
        <p role={feedback.kind === "error" ? "alert" : "status"} className={"text-sm " + (feedback.kind === "error" ? "text-destructive" : "text-emerald-600")}>
          {feedback.message}
        </p>
      ) : null}
    </section>
  );
}

registerCustomComponent("telegram-admin-tab", TelegramAdminTab);
