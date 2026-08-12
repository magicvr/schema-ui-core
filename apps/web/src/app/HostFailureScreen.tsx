import { useEffect, useRef, useState } from "react";

import { Button } from "@/components/ui/button";
import { useTranslate } from "@/i18n/runtime";
import type { HostFailure } from "@/host/failure";

const KIND_MESSAGE_KEY: Record<string, string> = {
  "maintenance": "hostFailure.maintenance",
  "upgrade-required": "hostFailure.upgradeRequired",
  "authentication-required": "hostFailure.requiresAuth",
  "reauth-required": "hostFailure.reauthRequired",
  "account-locked": "hostFailure.accountLocked",
  "forbidden": "hostFailure.forbidden",
  "not-found": "hostFailure.notFound",
  "rate-limited": "hostFailure.rateLimited",
  "timeout": "hostFailure.timeout",
  "offline": "hostFailure.offline",
  "protocol-rejected": "hostFailure.protocolRejected",
  "render-failed": "hostFailure.renderFailed",
  "unavailable": "hostFailure.unavailable",
};

const ACTION_MESSAGE_KEY: Record<string, string> = {
  retry: "hostFailure.action.retry",
  reauth: "hostFailure.action.reauth",
  home: "hostFailure.action.home",
  back: "hostFailure.action.back",
  reload: "hostFailure.action.reload",
  support: "hostFailure.action.support",
};

/** auth/protocol/render terminal failures use assertive announcements. */
function isAssertive(kind: string): boolean {
  return kind === "authentication-required"
    || kind === "reauth-required"
    || kind === "account-locked"
    || kind === "forbidden"
    || kind === "protocol-rejected"
    || kind === "render-failed";
}

export interface HostFailureScreenProps {
  failure: HostFailure;
  onAction: (action: { type: string; url?: string }) => void;
  /** When true, render as a section (caller provides the landmark, e.g. the shell main). */
  bare?: boolean;
}

/**
 * Global failure surface (ADR-0036 D7 / spec 10 §3.8 behavioral conformance):
 * unique error title inside the `main` landmark, focus moved to the title on
 * first terminal entry, assertive/polite live-region announcement by kind,
 * no re-announcement for the same failureId, keyboard-reachable recovery
 * actions.
 */
export function HostFailureScreen({ failure, onAction, bare = false }: HostFailureScreenProps) {
  const t = useTranslate();
  const titleRef = useRef<HTMLHeadingElement>(null);
  const lastAnnouncedRef = useRef<string | null>(null);
  const [announcement, setAnnouncement] = useState<string>("");
  const [remainingSeconds, setRemainingSeconds] = useState<number | null>(
    failure.retry?.mode === "after" ? (failure.retry.afterSeconds ?? 0) : null,
  );

  // First terminal entry: focus the title and announce once per failureId.
  // Redraws of the same failureId (countdown ticks) never re-steal focus nor
  // re-announce; a new failureId must re-announce.
  useEffect(() => {
    if (lastAnnouncedRef.current === failure.failureId) {
      return;
    }
    lastAnnouncedRef.current = failure.failureId;
    setAnnouncement(t(KIND_MESSAGE_KEY[failure.kind] ?? "hostFailure.generic"));
    titleRef.current?.focus();
  }, [failure.failureId, failure.kind, t]);

  // Countdown for retry.mode "after" (maintenance / rate-limit style).
  useEffect(() => {
    if (failure.retry?.mode !== "after" || remainingSeconds === null) {
      return undefined;
    }
    if (remainingSeconds <= 0) {
      setRemainingSeconds(null);
      return undefined;
    }
    const timer = window.setTimeout(() => {
      setRemainingSeconds((value) => (value === null ? null : value - 1));
    }, 1000);
    return () => window.clearTimeout(timer);
  }, [failure.retry, remainingSeconds]);

  const title = t(KIND_MESSAGE_KEY[failure.kind] ?? "hostFailure.generic");
  const assertive = isAssertive(failure.kind);
  // Host dictionary miss must use the safe generic copy, never echo the raw
  // document key or server text (ADR-0036 D5). i18n interpolation is
  // string/number only; boolean params stay on the wire, not in the copy.
  const interpolateParams = failure.message.params
    ? Object.fromEntries(
        Object.entries(failure.message.params).filter(
          (entry): entry is [string, string | number] => typeof entry[1] !== "boolean",
        ),
      )
    : undefined;
  const rawMessage = t(failure.message.messageKey, interpolateParams);
  const messageText = rawMessage === failure.message.messageKey ? title : rawMessage;

  // Full-screen surfaces own the main landmark; in-shell surfaces (route
  // not-found) render as a section inside the shell's existing main.
  const Landmark = bare ? "section" : "main";
  const landmarkClassName = bare
    ? "max-w-2xl space-y-6"
    : "flex min-h-screen items-center justify-center bg-background px-6 py-12 text-foreground";

  return (
    <Landmark
      className={landmarkClassName}
      aria-labelledby="host-failure-title"
    >
      <div className="mx-auto w-full max-w-xl space-y-6">
        <div aria-live={assertive ? "assertive" : "polite"} role="status" className="sr-only">
          {announcement}
        </div>
        <div className="space-y-2">
          <p className="text-xs font-semibold uppercase tracking-[0.16em] text-muted-foreground">
            {t("hostFailure.kicker")}
          </p>
          <h1
            id="host-failure-title"
            ref={titleRef}
            tabIndex={-1}
            className="text-3xl font-semibold tracking-tight outline-none"
          >
            {title}
          </h1>
          <p className="text-sm leading-6 text-muted-foreground">
            {messageText}
          </p>
          {failure.diagnostics?.protocolCode !== undefined && (
            <p className="font-mono text-xs text-muted-foreground">
              {failure.diagnostics.protocolCode}
            </p>
          )}
          {remainingSeconds !== null && remainingSeconds > 0 && (
            <p className="text-sm text-muted-foreground">
              {t("hostFailure.countdown", { seconds: remainingSeconds })}
            </p>
          )}
        </div>
        <div className="flex flex-wrap gap-3">
          {(failure.recoveryActions ?? []).map((action, index) => (
            <Button
              key={`${action.type}-${index}`}
              type="button"
              variant={index === 0 ? "default" : "outline"}
              onClick={() => onAction(action)}
            >
              {t(ACTION_MESSAGE_KEY[action.type] ?? "hostFailure.action.retry")}
            </Button>
          ))}
        </div>
      </div>
    </Landmark>
  );
}
