// Outbound-mail admin console (VP-017 R7 UX refinement, user-requested
// 2026-08-24; W26 GOAL-038 D-001 §2.2: promoted to the standalone「邮件控制台」
// page — channel select drives conditional per-channel fields (mock →
// retention, resend → key/from, smtp → host/port/username/password/from) and
// the test-send composer lets the admin author their own subject/body. The
// outbound records mini-table moved to the dedicated mail-outbox page).
//
// Implemented as a custom component (data-permission-scopes precedent):
// the requirements are tightly coupled to one channel-selection state,
// which declarative sibling nodes cannot share. Secrets stay write-only —
// empty inputs mean "keep current" and no read face ever returns them
// (Root D-007).
import { useCallback, useEffect, useState } from "react";

import { useTranslate } from "@/i18n/runtime";
import { registerCustomComponent, type CustomComponentProps } from "@/renderer/custom-components";
import { useSchemaCrud } from "@/renderer/render.tsx";

interface MailConfigView {
  channel: string;
  mockRetention: number;
  resend: { from: string };
  smtp: { host: string; port: number; username: string; from: string };
  secrets: { resendApiKeySet: boolean; smtpPasswordSet: boolean };
  updated_at?: string;
}

const CHANNELS = [
  { value: "mock", labelKey: "schema.settings.option.channelMock" },
  { value: "resend", labelKey: "schema.settings.option.channelResend" },
  { value: "smtp", labelKey: "schema.settings.option.channelSMTP" },
];

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

export function MailAdminTab(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;

  const [loadState, setLoadState] = useState<"loading" | "ready" | "error">("loading");
  const [channel, setChannel] = useState("mock");
  const [secrets, setSecrets] = useState({ resendApiKeySet: false, smtpPasswordSet: false });
  const [updatedHint, setUpdatedHint] = useState("");

  // Per-channel editable fields (strings keep number inputs manageable).
  const [mockRetention, setMockRetention] = useState("500");
  const [resendFrom, setResendFrom] = useState("");
  const [resendApiKey, setResendApiKey] = useState("");
  const [smtpHost, setSmtpHost] = useState("");
  const [smtpPort, setSmtpPort] = useState("");
  const [smtpUsername, setSmtpUsername] = useState("");
  const [smtpPassword, setSmtpPassword] = useState("");
  const [smtpFrom, setSmtpFrom] = useState("");

  // Test-send composer (R7 refinement: admin authors subject/body).
  const [testTo, setTestTo] = useState("");
  const [testSubject, setTestSubject] = useState("");
  const [testBody, setTestBody] = useState("");

  const [saving, setSaving] = useState(false);
  const [sending, setSending] = useState(false);
  const [feedback, setFeedback] = useState<{ kind: "success" | "error"; message: string } | null>(null);

  const applyConfigView = useCallback((view: MailConfigView) => {
    setChannel(view.channel);
    setMockRetention(String(view.mockRetention));
    setResendFrom(view.resend?.from ?? "");
    setSmtpHost(view.smtp?.host ?? "");
    setSmtpPort(view.smtp?.port ? String(view.smtp.port) : "");
    setSmtpUsername(view.smtp?.username ?? "");
    setSmtpFrom(view.smtp?.from ?? "");
    setSecrets({
      resendApiKeySet: view.secrets?.resendApiKeySet === true,
      smtpPasswordSet: view.secrets?.smtpPasswordSet === true,
    });
    if (typeof view.updated_at === "string") {
      setUpdatedHint(view.updated_at);
    }
  }, []);

  const loadConfig = useCallback(async () => {
    setLoadState("loading");
    try {
      const response = await fetcher("/api/mail/config");
      if (!response.ok) {
        setLoadState("error");
        return;
      }
      applyConfigView((await response.json()) as MailConfigView);
      setLoadState("ready");
    } catch {
      setLoadState("error");
    }
  }, [applyConfigView, fetcher]);

  useEffect(() => {
    void loadConfig();
  }, [loadConfig]);

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
    return t("schema.mail.feedback.saveFailed");
  }

  async function save() {
    setSaving(true);
    setFeedback(null);
    try {
      const payload: Record<string, unknown> = { channel };
      if (channel === "mock") {
        payload.mockRetention = Number(mockRetention) > 0 ? Number(mockRetention) : 500;
      } else if (channel === "resend") {
        payload.resendFrom = resendFrom;
        if (resendApiKey !== "") {
          payload.resendApiKey = resendApiKey;
        }
      } else if (channel === "smtp") {
        payload.smtpHost = smtpHost;
        payload.smtpUsername = smtpUsername;
        payload.smtpFrom = smtpFrom;
        if (smtpPort.trim() !== "" && Number.isFinite(Number(smtpPort))) {
          payload.smtpPort = Number(smtpPort);
        }
        if (smtpPassword !== "") {
          payload.smtpPassword = smtpPassword;
        }
      }
      const response = await fetcher("/api/mail/config", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload),
      });
      if (!response.ok) {
        setFeedback({ kind: "error", message: await extractError(response) });
        return;
      }
      applyConfigView((await response.json()) as MailConfigView);
      setResendApiKey("");
      setSmtpPassword("");
      setFeedback({ kind: "success", message: t("schema.mail.feedback.saved") });
    } catch {
      setFeedback({ kind: "error", message: t("schema.mail.feedback.saveFailed") });
    } finally {
      setSaving(false);
    }
  }

  async function sendTest() {
    setSending(true);
    setFeedback(null);
    try {
      const response = await fetcher("/api/mail/test-send", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ to: testTo.trim(), subject: testSubject.trim(), body: testBody }),
      });
      if (!response.ok) {
        setFeedback({ kind: "error", message: t("schema.mail.feedback.sendFailed") + ": " + (await extractError(response)) });
        return;
      }
      setFeedback({ kind: "success", message: t("schema.mail.feedback.sendOk") });
    } catch {
      setFeedback({ kind: "error", message: t("schema.mail.feedback.sendFailed") });
    } finally {
      setSending(false);
    }
  }

  if (loadState === "error") {
    return (
      <section data-mail-admin-tab className="space-y-3 rounded-xl border border-border/70 bg-card/85 p-4">
        <h2 className="text-sm font-semibold">{t("schema.settings.toolbar.mail")}</h2>
        <p role="alert" className="text-sm text-destructive">{t("schema.mail.feedback.loadFailed")}</p>
      </section>
    );
  }

  const fieldLabel = (key: string, id: string) => (
    <label className="text-sm font-medium" htmlFor={id}>
      {t(key)}
    </label>
  );

  return (
    <section data-mail-admin-tab className="space-y-4 rounded-xl border border-border/70 bg-card/85 p-4">
      <div className="flex items-center justify-between gap-3">
        <h2 className="text-sm font-semibold">{t("schema.settings.toolbar.mail")}</h2>
        {updatedHint !== "" ? (
          <span className="text-xs text-muted-foreground">{new Date(updatedHint).toLocaleString()}</span>
        ) : null}
      </div>

      {loadState === "loading" ? <p className="text-sm text-muted-foreground">{t("feedback.loading")}</p> : null}

      {/* Channel selection drives every conditional block below. */}
      <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
        {fieldLabel("schema.settings.field.channel", "mail-channel")}
        <select
          id="mail-channel"
          value={channel}
          onChange={(event) => setChannel(event.target.value)}
          className={inputClass}
        >
          {CHANNELS.map((option) => (
            <option key={option.value} value={option.value}>
              {t(option.labelKey)}
            </option>
          ))}
        </select>
      </div>

      {/* Requirement 1: only the selected channel's settings render. */}
      {channel === "mock" ? (
        <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
          {fieldLabel("schema.settings.field.mockRetention", "mail-mock-retention")}
          <input
            id="mail-mock-retention"
            type="number"
            min={1}
            max={100000}
            step={1}
            value={mockRetention}
            onChange={(event) => setMockRetention(event.target.value)}
            className={inputClass}
          />
        </div>
      ) : null}

      {channel === "resend" ? (
        <>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.resendApiKey", "mail-resend-key")}
            <input
              id="mail-resend-key"
              type="password"
              autoComplete="new-password"
              value={resendApiKey}
              placeholder={secrets.resendApiKeySet ? t("schema.mail.secret.keep") : ""}
              onChange={(event) => setResendApiKey(event.target.value)}
              className={inputClass}
            />
          </div>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.resendFrom", "mail-resend-from")}
            <input id="mail-resend-from" value={resendFrom} onChange={(event) => setResendFrom(event.target.value)} className={inputClass} />
          </div>
        </>
      ) : null}

      {channel === "smtp" ? (
        <>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.smtpHost", "mail-smtp-host")}
            <input id="mail-smtp-host" value={smtpHost} onChange={(event) => setSmtpHost(event.target.value)} className={inputClass} />
          </div>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.smtpPort", "mail-smtp-port")}
            <input
              id="mail-smtp-port"
              type="number"
              min={1}
              max={65535}
              value={smtpPort}
              placeholder={t("schema.mail.smtp.portPlaceholder")}
              onChange={(event) => setSmtpPort(event.target.value)}
              className={inputClass}
            />
          </div>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.smtpUsername", "mail-smtp-user")}
            <input id="mail-smtp-user" value={smtpUsername} onChange={(event) => setSmtpUsername(event.target.value)} className={inputClass} />
          </div>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.smtpPassword", "mail-smtp-password")}
            <input
              id="mail-smtp-password"
              type="password"
              autoComplete="new-password"
              value={smtpPassword}
              placeholder={secrets.smtpPasswordSet ? t("schema.mail.secret.keep") : ""}
              onChange={(event) => setSmtpPassword(event.target.value)}
              className={inputClass}
            />
          </div>
          <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
            {fieldLabel("schema.settings.field.smtpFrom", "mail-smtp-from")}
            <input id="mail-smtp-from" value={smtpFrom} onChange={(event) => setSmtpFrom(event.target.value)} className={inputClass} />
          </div>
        </>
      ) : null}

      <button type="button" disabled={saving || loadState !== "ready"} onClick={() => void save()} className={buttonClass}>
        {saving ? t("feedback.submitting") : t("schema.mail.action.save")}
      </button>

      {/* Requirement 3: the admin authors subject/body themselves. */}
      <div className="space-y-2 rounded-lg border border-border/60 p-3">
        <h3 className="text-sm font-semibold">{t("schema.mail.test.title")}</h3>
        <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
          {fieldLabel("schema.mail.test.to", "mail-test-to")}
          <input id="mail-test-to" value={testTo} onChange={(event) => setTestTo(event.target.value)} className={inputClass} />
        </div>
        <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-center">
          {fieldLabel("schema.mail.test.subject", "mail-test-subject")}
          <input
            id="mail-test-subject"
            value={testSubject}
            placeholder={t("schema.mail.test.subjectPlaceholder")}
            onChange={(event) => setTestSubject(event.target.value)}
            className={inputClass}
          />
        </div>
        <div className="grid gap-2 sm:grid-cols-[12rem_1fr] sm:items-start">
          {fieldLabel("schema.mail.test.body", "mail-test-body")}
          <textarea
            id="mail-test-body"
            rows={4}
            value={testBody}
            placeholder={t("schema.mail.test.bodyPlaceholder")}
            onChange={(event) => setTestBody(event.target.value)}
            className={inputClass + " h-auto py-2"}
          />
        </div>
        <button type="button" disabled={sending || testTo.trim() === ""} onClick={() => void sendTest()} className={buttonClass}>
          {sending ? t("feedback.submitting") : t("schema.mail.test.submit")}
        </button>
      </div>

      {feedback !== null ? (
        <p role={feedback.kind === "error" ? "alert" : "status"} className={"text-sm " + (feedback.kind === "error" ? "text-destructive" : "text-emerald-600")}>
          {feedback.message}
        </p>
      ) : null}
    </section>
  );
}

registerCustomComponent("mail-admin-tab", MailAdminTab);
