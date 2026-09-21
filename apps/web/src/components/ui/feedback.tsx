import { useEffect, useRef, useState } from "react";

import { X } from "lucide-react";

import type { MessageParams } from "@/i18n/catalog";
import { useTranslate } from "@/i18n/runtime";
import { cn } from "@/lib/utils";

export type FeedbackKind = "success" | "error";
export type FeedbackSurface = "toast" | "inline";

/** Shared operation feedback contract for toast and inline recovery surfaces. */
export interface FeedbackNotice {
  kind: FeedbackKind;
  message: string;
  code?: string;
  messageKey?: string;
  params?: Record<string, unknown>;
  /** A user-triggered recovery action; callers decide whether it is safe. */
  retry?: () => void | Promise<void>;
  /** Safe support correlation; never used as the primary user message. */
  correlationId?: string;
}

function safeParams(params: Record<string, unknown> | undefined): MessageParams | undefined {
  if (params === undefined) {
    return undefined;
  }
  const safe: MessageParams = {};
  for (const [key, value] of Object.entries(params)) {
    if (typeof value === "string" || typeof value === "number") {
      safe[key] = value;
    }
  }
  return Object.keys(safe).length === 0 ? undefined : safe;
}

function messageText(feedback: FeedbackNotice, t: (key: string, params?: MessageParams) => string): string {
  if (feedback.messageKey !== undefined && feedback.messageKey !== "") {
    const localized = t(feedback.messageKey, safeParams(feedback.params));
    if (localized !== feedback.messageKey) {
      return localized;
    }
  }
  return feedback.message;
}

export interface FeedbackNoticeViewProps {
  feedback: FeedbackNotice;
  surface?: FeedbackSurface;
  /** Toasts are dismissible; inline data errors stay visible until replaced. */
  dismissible?: boolean;
  /** Form-level surfaces keep the machine code available without making it the primary copy. */
  showDiagnosticCode?: boolean;
  /** Test/consumer hook for a surface-specific retry button. */
  retryButtonDataAttribute?: string;
}

/**
 * Shared accessible feedback surface (R4 / D-001).
 *
 * Success notices auto-dismiss. Errors remain until replaced or explicitly
 * dismissed, and retry is always a user gesture guarded by a local in-flight
 * flag so a double click cannot issue duplicate recovery requests.
 */
export function FeedbackNoticeView({
  feedback,
  surface = "toast",
  dismissible = surface === "toast",
  showDiagnosticCode = false,
  retryButtonDataAttribute,
}: FeedbackNoticeViewProps) {
  const t = useTranslate();
  const [dismissed, setDismissed] = useState(false);
  const [retrying, setRetrying] = useState(false);
  const retryingRef = useRef(false);
  // Callers commonly derive the notice object during render. Use the
  // user-visible result identity rather than object identity so a retry state
  // update cannot reset the in-flight guard or restart a success timer.
  const feedbackIdentity = JSON.stringify({
    kind: feedback.kind,
    code: feedback.code ?? "",
    message: feedback.message,
    messageKey: feedback.messageKey ?? "",
    params: feedback.params ?? null,
    correlationId: feedback.correlationId ?? "",
    surface,
  });

  useEffect(() => {
    setDismissed(false);
    retryingRef.current = false;
    setRetrying(false);
    if (feedback.kind === "error") {
      return undefined;
    }
    const timer = window.setTimeout(() => setDismissed(true), FEEDBACK_TOAST_MS);
    return () => window.clearTimeout(timer);
  }, [feedbackIdentity, feedback.kind]);

  if (dismissed) {
    return null;
  }

  const text = messageText(feedback, t);
  const isToast = surface === "toast";
  const diagnosticTitle = [feedback.code, feedback.correlationId].filter(Boolean).join(" · ");

  const handleRetry = async () => {
    if (feedback.retry === undefined || retryingRef.current) {
      return;
    }
    retryingRef.current = true;
    setRetrying(true);
    try {
      await feedback.retry();
    } catch {
      // The retry owner owns the resulting error state. Never create an
      // unhandled rejection or hide the current error surface here.
    } finally {
      retryingRef.current = false;
      setRetrying(false);
    }
  };

  return (
    <div
      role={feedback.kind === "error" ? "alert" : "status"}
      aria-atomic="true"
      aria-busy={retrying || undefined}
      data-feedback-toast={isToast ? feedback.kind : undefined}
      data-feedback-surface={surface}
      className={cn(
        "flex items-start justify-between gap-3 text-sm",
        isToast
          ? "fixed right-4 top-16 z-50 max-w-sm rounded-md border px-3 py-2 shadow-lg"
          : "rounded-md border px-4 py-3",
        feedback.kind === "error"
          ? "border-destructive/40 bg-destructive/10 text-destructive"
          : "border-success/50 bg-success/10 text-success",
      )}
    >
      <span
        data-feedback-code={feedback.code ?? ""}
        title={diagnosticTitle === "" ? undefined : diagnosticTitle}
        className="min-w-0 flex-1"
      >
        {showDiagnosticCode && feedback.code !== undefined ? (
          <span className="sr-only" data-feedback-diagnostic>
            {feedback.code}: {" "}
          </span>
        ) : null}
        {text}
      </span>
      <span className="flex shrink-0 items-center gap-2">
        {feedback.retry !== undefined ? (
          <button
            type="button"
            {...(retryButtonDataAttribute === undefined
              ? {}
              : { [retryButtonDataAttribute]: "true" })}
            disabled={retrying}
            onClick={() => void handleRetry()}
            className="font-medium underline underline-offset-2 disabled:no-underline"
          >
            {retrying ? t("feedback.retrying") : t("feedback.retry")}
          </button>
        ) : null}
        {dismissible ? (
          <button
            type="button"
            aria-label={t("feedback.dismiss")}
            onClick={() => setDismissed(true)}
            className="rounded-sm opacity-70 transition-opacity hover:opacity-100 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring"
          >
            <X aria-hidden="true" className="size-4" />
          </button>
        ) : null}
      </span>
    </div>
  );
}

/** Backward-compatible toast entry point for the schema CRUD provider. */
export function FeedbackRegion({ feedback }: { feedback: FeedbackNotice }) {
  return <FeedbackNoticeView feedback={feedback} surface="toast" />;
}

/** Auto-dismiss window for successful operation feedback toasts. */
export const FEEDBACK_TOAST_MS = 4000;
