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
  mode?: string;
  webhook_public_base_url?: string;
  connection_state?: string;
  receiver?: string;
  bot_id?: number;
  bot_username?: string;
  last_error?: string;
  captured_messages_count?: number;
  captured_count?: number;
}

type LeaseAction = "acquire" | "heartbeat" | "release";

interface TelegramLeaseResult {
  ok: boolean;
  connection_state?: string;
  receiver?: string;
}

const telegramPollingMode = "polling";
const telegramLeaseIntervalMs = 10_000;
const telegramLeasePaths: Record<LeaseAction, string> = {
  acquire: "/api/channel/telegram/lease/acquire",
  heartbeat: "/api/channel/telegram/lease/heartbeat",
  release: "/api/channel/telegram/lease/release",
};

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
  const [modeInput, setModeInput] = useState(telegramPollingMode);
  const [webhookPublicBaseURLInput, setWebhookPublicBaseURLInput] = useState("");
  const [saving, setSaving] = useState(false);
  const [clearing, setClearing] = useState(false);
  const [confirmClear, setConfirmClear] = useState(false);
  const [leaseState, setLeaseState] = useState<"inactive" | "acquiring" | "active" | "error">("inactive");
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const loadStatus = useCallback(async () => {
    setLoadState("loading");
    try {
      const response = await fetcher("/api/channel/telegram/settings");
      if (!response.ok) {
        setLoadState("error");
        return;
      }
      const nextStatus = (await response.json()) as TelegramSettingsStatus;
      setStatus(nextStatus);
      setModeInput(nextStatus.mode === "webhook" ? "webhook" : telegramPollingMode);
      setWebhookPublicBaseURLInput(nextStatus.webhook_public_base_url ?? "");
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, [fetcher]);

  const callLease = useCallback(
    async (action: LeaseAction): Promise<TelegramLeaseResult> => {
      try {
        const response = await fetcher(telegramLeasePaths[action], { method: "POST" });
        if (!response.ok) return { ok: false };
        const body = (await response.json()) as Omit<TelegramLeaseResult, "ok"> & { ok?: boolean };
        return {
          ok: body.ok !== false,
          connection_state: body.connection_state,
          receiver: body.receiver,
        };
      } catch {
        return { ok: false };
      }
    },
    [fetcher],
  );

  useEffect(() => {
    if (loadState !== "ready" || status?.mode !== telegramPollingMode) {
      setLeaseState("inactive");
      return;
    }

    let disposed = false;
    let leaseHeld = false;
    let timer: ReturnType<typeof setTimeout> | undefined;
    let leaseQueue: Promise<TelegramLeaseResult> = Promise.resolve({ ok: true });

    const queueLease = (action: LeaseAction) => {
      const result = leaseQueue.then(() => callLease(action));
      leaseQueue = result.then(() => ({ ok: true }), () => ({ ok: false }));
      return result;
    };

    const applyLeaseSnapshot = (result: TelegramLeaseResult) => {
      if (result.connection_state === undefined && result.receiver === undefined) return;
      setStatus((current) => current === null
        ? current
        : {
            ...current,
            connection_state: result.connection_state ?? current.connection_state,
            receiver: result.receiver ?? current.receiver,
          });
    };

    const scheduleHeartbeat = () => {
      timer = setTimeout(() => {
        void heartbeat();
      }, telegramLeaseIntervalMs);
    };

    const heartbeat = async () => {
      if (disposed) return;
      // HeartbeatLease intentionally creates an unknown lease to recover a
      // lost acquire response. Mark the lease as potentially live before the
      // request starts so cleanup always queues release after this request.
      leaseHeld = true;
      const result = await queueLease("heartbeat");
      if (disposed) return;
      if (result.ok) {
        applyLeaseSnapshot(result);
        setLeaseState("active");
      } else {
        setLeaseState("error");
      }
      scheduleHeartbeat();
    };

    const acquire = async () => {
      setLeaseState("acquiring");
      const result = await queueLease("acquire");
      if (disposed) {
        if (result.ok) void queueLease("release");
        return;
      }
      if (!result.ok) {
        setLeaseState("error");
        scheduleHeartbeat();
        return;
      }
      leaseHeld = true;
      applyLeaseSnapshot(result);
      setLeaseState("active");
      scheduleHeartbeat();
    };

    void acquire();
    return () => {
      disposed = true;
      if (timer !== undefined) clearTimeout(timer);
      if (leaseHeld) void queueLease("release");
    };
  }, [callLease, loadState, status?.mode]);

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
      if (status !== null && modeInput !== (status.mode ?? telegramPollingMode)) {
        payload.mode = modeInput;
      }
      if (status !== null && webhookPublicBaseURLInput.trim() !== (status.webhook_public_base_url ?? "")) {
        payload.webhook_public_base_url = webhookPublicBaseURLInput.trim();
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
      const nextStatus = (await response.json()) as TelegramSettingsStatus;
      setStatus(nextStatus);
      setModeInput(nextStatus.mode === "webhook" ? "webhook" : telegramPollingMode);
      setWebhookPublicBaseURLInput(nextStatus.webhook_public_base_url ?? "");
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
      const nextStatus = (await response.json()) as TelegramSettingsStatus;
      setStatus(nextStatus);
      setModeInput(nextStatus.mode === "webhook" ? "webhook" : telegramPollingMode);
      setWebhookPublicBaseURLInput(nextStatus.webhook_public_base_url ?? "");
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

  const connectionStateLabel = (() => {
    switch (status?.connection_state) {
      case "unconfigured":
        return t("schema.telegram.connection.unconfigured");
      case "starting":
        return t("schema.telegram.connection.starting");
      case "running":
        return t("schema.telegram.connection.running");
      case "stopping":
        return t("schema.telegram.connection.stopping");
      case "error":
        return t("schema.telegram.connection.error");
      case "idle":
        return t("schema.telegram.connection.idle");
      default:
        return t("schema.telegram.connection.unknown");
    }
  })();

  const receiverLabel = status?.receiver === "polling"
    ? t("schema.telegram.receiver.polling")
    : status?.receiver === "webhook"
      ? t("schema.telegram.receiver.webhook")
      : t("schema.telegram.receiver.none");

  const leaseLabel = leaseState === "active"
    ? t("schema.telegram.lease.active")
    : leaseState === "acquiring"
      ? t("schema.telegram.lease.acquiring")
      : leaseState === "error"
        ? t("schema.telegram.lease.error")
        : null;

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

      {status !== null ? (
        <div data-telegram-connection className="space-y-1 rounded-md border border-border/60 bg-muted/20 px-3 py-2 text-xs text-muted-foreground">
          <p>
            {t("schema.telegram.status.connection")} {connectionStateLabel} · {receiverLabel}
          </p>
          {status.bot_username ? <p>{t("schema.telegram.status.bot")} @{status.bot_username}</p> : null}
          {status.last_error ? <p role="alert" className="text-destructive">{status.last_error}</p> : null}
          {leaseLabel !== null && status.mode === telegramPollingMode ? (
            <p role={leaseState === "error" ? "alert" : "status"}>{t("schema.telegram.status.consoleLease")} {leaseLabel}</p>
          ) : null}
        </div>
      ) : null}

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
      <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
        {fieldLabel("schema.telegram.field.mode", "telegram-mode")}
        <select
          id="telegram-mode"
          value={modeInput}
          onChange={(event) => setModeInput(event.target.value)}
          className={inputClass}
        >
          <option value="polling">{t("schema.telegram.mode.polling")}</option>
          <option value="webhook">{t("schema.telegram.mode.webhook")}</option>
        </select>
      </div>
      <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
        {fieldLabel("schema.telegram.field.webhookPublicBaseURL", "telegram-webhook-public-base-url")}
        <input
          id="telegram-webhook-public-base-url"
          type="url"
          autoComplete="url"
          value={webhookPublicBaseURLInput}
          placeholder="https://example.com"
          onChange={(event) => setWebhookPublicBaseURLInput(event.target.value)}
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
