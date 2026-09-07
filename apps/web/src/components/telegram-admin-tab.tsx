// Telegram channel admin console (GOAL-006 R5, 判据 #5 补做 Admin UI tab,
// user-adjudicated 2026-09-03): the channel.telegram settings page is hosted by
// these custom surfaces — Bot Token / Webhook Secret are edited write-only
// (GET only reports token_set/secret_set booleans; PATCH accepts new values,
// empty keeps current — F-002 / R-005). The operator surface also exposes the
// captured-message counter for the mock outbound sink. Secrets never leave the
// API in plaintext or partial masks.
import { useCallback, useEffect, useLayoutEffect, useRef, useState } from "react";

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
  business_occupied?: boolean;
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

interface TelegramSession {
  chatId: string;
  chatType: string;
  title: string;
  username?: string;
  lastMessageAt: string;
}

interface TelegramTimelineItem {
  chatId: string;
  direction: "inbound" | "outbound" | string;
  status: string;
  occurredAt: string;
  updateId?: string;
  messageId?: string;
  userId?: string;
  senderUsername?: string;
  requestId?: string;
  retryOf?: string | null;
  text: string;
}

interface TelegramPagedResponse<T> {
  items?: T[];
  total?: number;
  page?: number;
  pageSize?: number;
}

interface TelegramCapabilityResponse {
  chatId: string;
  canSend: boolean;
}

const telegramPollingMode = "polling";
const telegramLeaseIntervalMs = 10_000;
const telegramLeasePaths: Record<LeaseAction, string> = {
  acquire: "/api/channel/telegram/lease/acquire",
  heartbeat: "/api/channel/telegram/lease/heartbeat",
  release: "/api/channel/telegram/lease/release",
};
const telegramOperatorSessionsPath = "/api/channel/telegram/operator/sessions";
const telegramOperatorPageQuery = "?page=1&pageSize=100";
const telegramTimelineBottomThreshold = 48;

function telegramCapabilityPath(chatId: string, force = false): string {
  return `${telegramOperatorSessionsPath}/${encodeURIComponent(chatId)}/capability${force ? "?refresh=1" : ""}`;
}

function createTelegramOperatorRequestID(): string {
  return `operator-${Date.now().toString(36)}-${Math.random().toString(36).slice(2)}`;
}

function sortTelegramTimeline(items: TelegramTimelineItem[]): TelegramTimelineItem[] {
  return items
    .map((item, index) => ({ item, index, timestamp: Date.parse(item.occurredAt) }))
    .sort((left, right) => {
      const leftTimestamp = Number.isNaN(left.timestamp) ? Number.POSITIVE_INFINITY : left.timestamp;
      const rightTimestamp = Number.isNaN(right.timestamp) ? Number.POSITIVE_INFINITY : right.timestamp;
      return leftTimestamp - rightTimestamp || left.index - right.index;
    })
    .map(({ item }) => item);
}

function telegramTimelineIsNearBottom(element: HTMLElement): boolean {
  return element.scrollHeight - element.scrollTop - element.clientHeight <= telegramTimelineBottomThreshold;
}

function insertTelegramComposerLineBreak(textarea: HTMLTextAreaElement, setValue: (value: string) => void): void {
  const selectionStart = textarea.selectionStart;
  const selectionEnd = textarea.selectionEnd;
  const nextValue = `${textarea.value.slice(0, selectionStart)}\n${textarea.value.slice(selectionEnd)}`;
  const nextCaret = selectionStart + 1;
  setValue(nextValue);
  const restoreSelection = () => {
    textarea.selectionStart = nextCaret;
    textarea.selectionEnd = nextCaret;
  };
  if (typeof requestAnimationFrame === "function") {
    requestAnimationFrame(restoreSelection);
  } else {
    setTimeout(restoreSelection, 0);
  }
}

const inputClass =
  "h-9 w-full rounded-md border border-input/80 bg-background px-3 text-sm shadow-2xs outline-none transition-all duration-150 hover:border-muted-foreground/30 focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20";
const buttonClass =
  "inline-flex h-9 cursor-pointer items-center justify-center gap-1.5 rounded-md bg-primary px-3.5 text-sm font-medium text-primary-foreground shadow-sm transition-all duration-150 hover:bg-primary/90 disabled:cursor-not-allowed disabled:opacity-50";

export function TelegramAdminTab(_props: CustomComponentProps) {
  const t = useTranslate();
  const crud = useSchemaCrud();
  const fetcher = crud?.fetcher ?? globalThis.fetch;
  const isOperatorSurface = _props.node.props?.surface === "operator";

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
  const [pageVisible, setPageVisible] = useState(() => typeof document === "undefined" || !document.hidden);
  const [sessionsLoadState, setSessionsLoadState] = useState<"idle" | "loading" | "ready" | "error">("idle");
  const [sessions, setSessions] = useState<TelegramSession[]>([]);
  const [selectedChatId, setSelectedChatId] = useState<string | null>(null);
  const [timelineLoadState, setTimelineLoadState] = useState<"idle" | "loading" | "ready" | "error">("idle");
  const [timeline, setTimeline] = useState<TelegramTimelineItem[]>([]);
  const [operatorCapability, setOperatorCapability] = useState<"unknown" | "allowed" | "denied" | "error">("unknown");
  const [composerText, setComposerText] = useState("");
  const [reverseComposerShortcuts, setReverseComposerShortcuts] = useState(false);
  const [sending, setSending] = useState(false);
  const [retryingRequestID, setRetryingRequestID] = useState<string | null>(null);
  const selectedChatRef = useRef<string | null>(null);
  const operatorReadyRef = useRef(false);
  const operatorRefreshRef = useRef<Promise<void> | null>(null);
  const timelineFlightsRef = useRef(new Map<string, Promise<void>>());
  const capabilityFlightsRef = useRef(new Map<string, Promise<void>>());
  const timelineListRef = useRef<HTMLDivElement | null>(null);
  const timelineStickToBottomRef = useRef(true);

  useEffect(() => {
    selectedChatRef.current = selectedChatId;
  }, [selectedChatId]);

  useEffect(() => {
    const onVisibilityChange = () => {
      setPageVisible(typeof document === "undefined" || !document.hidden);
    };
    document.addEventListener("visibilitychange", onVisibilityChange);
    return () => document.removeEventListener("visibilitychange", onVisibilityChange);
  }, []);

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
    if (!isOperatorSurface || loadState !== "ready" || status?.mode !== telegramPollingMode || status.business_occupied !== false) {
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

    let pageIsVisible = typeof document === "undefined" || !document.hidden;

    const scheduleHeartbeat = () => {
      if (!pageIsVisible) return;
      timer = setTimeout(() => {
        void heartbeat();
      }, telegramLeaseIntervalMs);
    };

    const heartbeat = async () => {
      if (disposed || !pageIsVisible) return;
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
      if (disposed || !pageIsVisible) return;
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

    const onVisibilityChange = () => {
      pageIsVisible = typeof document === "undefined" || !document.hidden;
      if (!pageIsVisible) {
        if (timer !== undefined) clearTimeout(timer);
        return;
      }
      if (leaseHeld) {
        void heartbeat();
      } else {
        void acquire();
      }
    };

    document.addEventListener("visibilitychange", onVisibilityChange);
    void acquire();
    return () => {
      disposed = true;
      document.removeEventListener("visibilitychange", onVisibilityChange);
      if (timer !== undefined) clearTimeout(timer);
      if (leaseHeld) void queueLease("release");
    };
  }, [callLease, isOperatorSurface, loadState, status?.business_occupied, status?.mode]);

  useEffect(() => {
    void loadStatus();
  }, [loadStatus]);

  const operatorReady = isOperatorSurface
    && status !== null
    && status.configured
    && status.business_occupied === false
    && status.connection_state === "running"
    && typeof status.bot_id === "number"
    && status.bot_id > 0
    && (status.receiver !== telegramPollingMode || leaseState === "active");
  const selectedSession = sessions.find((session) => session.chatId === selectedChatId) ?? null;

  useEffect(() => {
    operatorReadyRef.current = operatorReady;
  }, [operatorReady]);

  const loadTimeline = useCallback(async (chatId: string): Promise<void> => {
    const existing = timelineFlightsRef.current.get(chatId);
    if (existing !== undefined) {
      await existing;
      return;
    }

    const request = (async () => {
      setTimelineLoadState("loading");
      try {
        const response = await fetcher(
          `${telegramOperatorSessionsPath}/${encodeURIComponent(chatId)}/messages${telegramOperatorPageQuery}`,
          { headers: { Accept: "application/json" } },
        );
        if (!response.ok) {
          if (selectedChatRef.current === chatId) setTimelineLoadState("error");
          return;
        }
        const body = (await response.json()) as TelegramPagedResponse<TelegramTimelineItem>;
        if (selectedChatRef.current === chatId) {
          setTimeline(sortTelegramTimeline(Array.isArray(body.items) ? body.items : []));
          setTimelineLoadState("ready");
        }
      } catch {
        if (selectedChatRef.current === chatId) setTimelineLoadState("error");
      }
    })();
    timelineFlightsRef.current.set(chatId, request);
    try {
      await request;
    } finally {
      if (timelineFlightsRef.current.get(chatId) === request) {
        timelineFlightsRef.current.delete(chatId);
      }
    }
  }, [fetcher]);

  const loadCapability = useCallback(async (chatId: string, force = false): Promise<void> => {
    const existing = capabilityFlightsRef.current.get(chatId);
    if (existing !== undefined) {
      await existing;
      return;
    }

    const request = (async () => {
      if (selectedChatRef.current === chatId && operatorReadyRef.current) {
        setOperatorCapability("unknown");
      }
      try {
        const response = await fetcher(telegramCapabilityPath(chatId, force), {
          headers: { Accept: "application/json" },
        });
        if (!response.ok) {
          if (selectedChatRef.current === chatId && operatorReadyRef.current) setOperatorCapability("error");
          return;
        }
        const body = (await response.json()) as Partial<TelegramCapabilityResponse>;
        if (body.chatId !== chatId || typeof body.canSend !== "boolean") {
          throw new Error("invalid Telegram capability response");
        }
        if (selectedChatRef.current === chatId && operatorReadyRef.current) {
          setOperatorCapability(body.canSend ? "allowed" : "denied");
        }
      } catch {
        if (selectedChatRef.current === chatId && operatorReadyRef.current) setOperatorCapability("error");
      }
    })();
    capabilityFlightsRef.current.set(chatId, request);
    try {
      await request;
    } finally {
      if (capabilityFlightsRef.current.get(chatId) === request) {
        capabilityFlightsRef.current.delete(chatId);
      }
    }
  }, [fetcher]);

  const loadSessions = useCallback(async (): Promise<TelegramSession[] | null> => {
    setSessionsLoadState("loading");
    try {
      const response = await fetcher(
        `${telegramOperatorSessionsPath}${telegramOperatorPageQuery}`,
        { headers: { Accept: "application/json" } },
      );
      if (!response.ok) {
        setSessionsLoadState("error");
        return null;
      }
      const body = (await response.json()) as TelegramPagedResponse<TelegramSession>;
      const nextSessions = Array.isArray(body.items) ? body.items : [];
      setSessions(nextSessions);
      const currentChatId = selectedChatRef.current;
      const nextChatId = currentChatId !== null && nextSessions.some((session) => session.chatId === currentChatId)
        ? currentChatId
        : nextSessions[0]?.chatId ?? null;
      const chatChanged = currentChatId !== nextChatId;
      selectedChatRef.current = nextChatId;
      setSelectedChatId(nextChatId);
      if (nextChatId === null) {
        timelineStickToBottomRef.current = true;
        setTimeline([]);
        setTimelineLoadState("ready");
        setOperatorCapability("unknown");
      } else if (chatChanged) {
        timelineStickToBottomRef.current = true;
        setTimeline([]);
        setTimelineLoadState("loading");
        setOperatorCapability("unknown");
      }
      setSessionsLoadState("ready");
      return nextSessions;
    } catch {
      setSessionsLoadState("error");
      return null;
    }
  }, [fetcher]);

  const refreshOperatorSurface = useCallback(async (forceCapability = false) => {
    const existing = operatorRefreshRef.current;
    if (existing !== null) {
      await existing;
      return;
    }
    const request = (async () => {
      const nextSessions = await loadSessions();
      if (nextSessions === null) return;
      const nextChatId = selectedChatRef.current;
      if (nextChatId !== null) {
        await loadTimeline(nextChatId);
        if (forceCapability) await loadCapability(nextChatId, true);
      }
    })();
    operatorRefreshRef.current = request;
    try {
      await request;
    } finally {
      if (operatorRefreshRef.current === request) operatorRefreshRef.current = null;
    }
  }, [loadCapability, loadSessions, loadTimeline]);

  useEffect(() => {
    if (!isOperatorSurface || !operatorReady || !pageVisible) return;
    let disposed = false;
    let timer: ReturnType<typeof setTimeout> | undefined;

    const refresh = async () => {
      await refreshOperatorSurface();
      if (!disposed && pageVisible && (typeof document === "undefined" || !document.hidden)) {
        timer = setTimeout(() => {
          void refresh();
        }, telegramLeaseIntervalMs);
      }
    };

    void refresh();
    return () => {
      disposed = true;
      if (timer !== undefined) clearTimeout(timer);
    };
  }, [isOperatorSurface, operatorReady, pageVisible, refreshOperatorSurface]);

  useEffect(() => {
    if (!isOperatorSurface || !operatorReady || selectedChatId === null) {
      timelineStickToBottomRef.current = true;
      setOperatorCapability("unknown");
      setTimeline([]);
      setTimelineLoadState(isOperatorSurface && operatorReady ? "ready" : "idle");
      return;
    }
    void loadTimeline(selectedChatId);
    void loadCapability(selectedChatId, true);
  }, [isOperatorSurface, loadCapability, loadTimeline, operatorReady, selectedChatId]);

  useLayoutEffect(() => {
    if (!isOperatorSurface || selectedChatId === null || timeline.length === 0 || !timelineStickToBottomRef.current) {
      return;
    }
    const element = timelineListRef.current;
    if (element !== null) {
      element.scrollTop = element.scrollHeight;
    }
  }, [isOperatorSurface, selectedChatId, timeline]);

  async function extractError(response: Response, fallbackKey = "schema.telegram.feedback.saveFailed"): Promise<string> {
    try {
      const body = (await response.json()) as { detail?: string; message?: string };
      const text = typeof body.detail === "string" && body.detail !== "" ? body.detail : body.message;
      if (typeof text === "string" && text !== "") {
        return text;
      }
    } catch {
      // fall through to generic message
    }
    return t(fallbackKey);
  }

  async function sendMessage() {
    const chatId = selectedChatRef.current;
    const text = composerText.trim();
    if (chatId === null || text === "" || !operatorReadyRef.current || operatorCapability !== "allowed" || sending) {
      return;
    }
    const requestID = createTelegramOperatorRequestID();
    setSending(true);
    setFeedback(null);
    setOperatorCapability("unknown");
    try {
      const response = await fetcher(`${telegramOperatorSessionsPath}/${encodeURIComponent(chatId)}/messages`, {
        method: "POST",
        headers: { Accept: "application/json", "Content-Type": "application/json" },
        body: JSON.stringify({ requestId: requestID, text }),
      });
      if (!response.ok) {
        if (selectedChatRef.current === chatId) {
          setFeedback({ kind: "error", message: await extractError(response, "schema.telegram.feedback.sendFailed") });
          setOperatorCapability("unknown");
        }
        await loadTimeline(chatId);
        return;
      }
      if (selectedChatRef.current === chatId && operatorReadyRef.current) {
        setComposerText("");
        setOperatorCapability("allowed");
      }
      await loadTimeline(chatId);
    } catch {
      if (selectedChatRef.current === chatId) {
        setFeedback({ kind: "error", message: t("schema.telegram.feedback.sendFailed") });
        setOperatorCapability("unknown");
      }
      await loadTimeline(chatId);
    } finally {
      setSending(false);
    }
  }

  async function retryMessage(sourceRequestID: string) {
    const chatId = selectedChatRef.current;
    if (chatId === null || !operatorReadyRef.current || operatorCapability !== "allowed" || retryingRequestID !== null) {
      return;
    }
    const requestID = createTelegramOperatorRequestID();
    setRetryingRequestID(sourceRequestID);
    setFeedback(null);
    setOperatorCapability("unknown");
    try {
      const response = await fetcher(
        `${telegramOperatorSessionsPath}/${encodeURIComponent(chatId)}/messages/${encodeURIComponent(sourceRequestID)}/retry`,
        {
          method: "POST",
          headers: { Accept: "application/json", "Content-Type": "application/json" },
          body: JSON.stringify({ requestId: requestID }),
        },
      );
      if (!response.ok) {
        if (selectedChatRef.current === chatId) {
          setFeedback({ kind: "error", message: await extractError(response, "schema.telegram.feedback.retryFailed") });
          setOperatorCapability("unknown");
        }
        await loadTimeline(chatId);
        return;
      }
      if (selectedChatRef.current === chatId && operatorReadyRef.current) setOperatorCapability("allowed");
      await loadTimeline(chatId);
    } catch {
      if (selectedChatRef.current === chatId) {
        setFeedback({ kind: "error", message: t("schema.telegram.feedback.retryFailed") });
        setOperatorCapability("unknown");
      }
      await loadTimeline(chatId);
    } finally {
      setRetryingRequestID(null);
    }
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

  const surfaceTitleKey = isOperatorSurface
    ? "schema.telegram.operator.title"
    : "schema.settings.toolbar.telegram";

  if (loadState === "error") {
    return (
      <section
        data-telegram-admin-tab
        data-telegram-operator-page={isOperatorSurface ? "true" : undefined}
        className={isOperatorSurface
          ? "flex h-full min-h-0 min-w-0 flex-1 flex-col space-y-3 overflow-hidden rounded-xl border border-border/70 bg-card/85 p-4"
          : "space-y-3 rounded-xl border border-border/70 bg-card/85 p-4"}
      >
        <h2 className="text-sm font-semibold">{t(surfaceTitleKey)}</h2>
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
    <section
      data-telegram-admin-tab
      data-telegram-operator-page={isOperatorSurface ? "true" : undefined}
      className={isOperatorSurface
        ? "flex h-full min-h-0 min-w-0 flex-1 flex-col space-y-4 overflow-hidden rounded-xl border border-border/70 bg-card/85 p-4"
        : "space-y-4 rounded-xl border border-border/70 bg-card/85 p-4"}
    >
      <div className="flex shrink-0 items-center justify-between gap-3">
        <h2 className="text-sm font-semibold">{t(surfaceTitleKey)}</h2>
        {status !== null ? (
          <span className="text-xs text-muted-foreground">
            {status.configured
              ? t("schema.telegram.status.configured")
              : t("schema.telegram.status.notConfigured")}
          </span>
        ) : null}
      </div>

      {status !== null ? (
        <div data-telegram-connection className="shrink-0 space-y-1 rounded-md border border-border/60 bg-muted/20 px-3 py-2 text-xs text-muted-foreground">
          <p>
            {t("schema.telegram.status.connection")} {connectionStateLabel} · {receiverLabel}
          </p>
          {status.mode === telegramPollingMode ? (
            <p role="alert" data-telegram-polling-warning>
              {t("schema.telegram.status.pollingSingleInstanceWarning")}
            </p>
          ) : null}
          {status.bot_username ? <p>{t("schema.telegram.status.bot")} @{status.bot_username}</p> : null}
          {status.last_error ? <p role="alert" className="text-destructive">{status.last_error}</p> : null}
          {leaseLabel !== null && status.mode === telegramPollingMode ? (
            <p role={leaseState === "error" ? "alert" : "status"}>{t("schema.telegram.status.consoleLease")} {leaseLabel}</p>
          ) : null}
        </div>
      ) : null}

      {loadState === "loading" ? <p className="shrink-0 text-sm text-muted-foreground">{t("feedback.loading")}</p> : null}

      {!isOperatorSurface ? (
        <>
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

        </>
      ) : null}

      {isOperatorSurface && status !== null && status.configured && status.business_occupied === false && typeof status.bot_id === "number" && status.bot_id > 0 ? (
        <section data-telegram-operator className="flex min-h-0 min-w-0 flex-1 flex-col gap-3 overflow-hidden rounded-md border border-border/60 bg-muted/10 p-3">
          <div className="flex items-center justify-between gap-3">
            <h3 className="text-sm font-semibold">{t("schema.telegram.operator.title")}</h3>
            <button
              type="button"
              data-telegram-operator-refresh
              disabled={!operatorReady || sessionsLoadState === "loading"}
              onClick={() => void refreshOperatorSurface(true)}
              className="rounded-md border border-input/80 px-2.5 py-1 text-xs hover:bg-muted disabled:cursor-not-allowed disabled:opacity-50"
            >
              {sessionsLoadState === "loading" ? t("feedback.loading") : t("schema.telegram.operator.refresh")}
            </button>
          </div>
          {typeof status.captured_messages_count === "number" ? (
            <p className="text-xs text-muted-foreground">
              {t("schema.telegram.status.captured")} {status.captured_messages_count}
            </p>
          ) : null}
          {sessionsLoadState === "error" && sessions.length > 0 ? (
            <p role="alert" className="shrink-0 text-xs text-destructive">{t("schema.telegram.operator.loadFailed")}</p>
          ) : null}

          {!operatorReady ? (
            <p className="text-xs text-muted-foreground">{t("schema.telegram.operator.unavailable")}</p>
          ) : sessions.length === 0 ? (
            sessionsLoadState === "error" ? (
              <p role="alert" className="shrink-0 text-xs text-destructive">{t("schema.telegram.operator.loadFailed")}</p>
            ) : sessionsLoadState === "ready" ? (
              <p className="shrink-0 text-xs text-muted-foreground">{t("schema.telegram.operator.empty")}</p>
            ) : (
              <p className="shrink-0 text-xs text-muted-foreground">{t("feedback.loading")}</p>
            )
          ) : (
            <div className="grid min-h-0 min-w-0 flex-1 grid-rows-[auto_minmax(0,1fr)] gap-3 overflow-hidden lg:grid-cols-[minmax(12rem,16rem)_minmax(0,1fr)] lg:grid-rows-[minmax(0,1fr)]">
              <nav aria-label={t("schema.telegram.operator.select")} data-telegram-sessions className="max-h-40 min-h-0 space-y-1 overflow-x-hidden overflow-y-auto overscroll-contain pr-1 lg:max-h-none">
                {sessions.map((session) => {
                  const displayName = session.title || session.username || session.chatId;
                  return (
                    <button
                      type="button"
                      key={session.chatId}
                      data-telegram-session={session.chatId}
                      aria-pressed={selectedChatId === session.chatId}
                      onClick={() => {
                        const chatChanged = selectedChatRef.current !== session.chatId;
                        selectedChatRef.current = session.chatId;
                        setSelectedChatId(session.chatId);
                        if (chatChanged) {
                          timelineStickToBottomRef.current = true;
                          setTimeline([]);
                          setTimelineLoadState("loading");
                          setOperatorCapability("unknown");
                        }
                        void loadTimeline(session.chatId);
                        void loadCapability(session.chatId, true);
                      }}
                      className="block w-full rounded-md border border-transparent px-2.5 py-2 text-left text-xs hover:bg-muted aria-pressed:border-border aria-pressed:bg-muted/60"
                    >
                      <span className="block truncate font-medium">{displayName}</span>
                      <span className="block break-all text-muted-foreground">{session.chatType} · {session.chatId}</span>
                    </button>
                  );
                })}
              </nav>

              <div data-telegram-transcript className="flex min-h-0 min-w-0 flex-col gap-3 overflow-hidden">
                <div className="flex shrink-0 items-center justify-between gap-2">
                  <h4 className="min-w-0 truncate text-xs font-semibold text-muted-foreground">{t("schema.telegram.operator.timeline")}</h4>
                  {timelineLoadState === "loading" && timeline.length > 0 ? (
                    <span role="status" data-telegram-timeline-refreshing className="shrink-0 text-xs text-muted-foreground">{t("schema.telegram.operator.timelineRefreshing")}</span>
                  ) : timelineLoadState === "error" && timeline.length > 0 ? (
                    <span role="alert" className="shrink-0 text-xs text-destructive">{t("schema.telegram.operator.timelineFailed")}</span>
                  ) : null}
                </div>
                {timeline.length > 0 ? (
                  <div
                    ref={timelineListRef}
                    data-telegram-message-list
                    onScroll={(event) => {
                      timelineStickToBottomRef.current = telegramTimelineIsNearBottom(event.currentTarget);
                    }}
                    className="min-h-0 min-w-0 flex-1 space-y-2 overflow-x-hidden overflow-y-auto overscroll-contain pr-1"
                  >
                    {timeline.map((item, index) => {
                      const isOutbound = item.direction === "outbound";
                      const directionLabel = item.direction === "inbound"
                        ? t("schema.telegram.operator.inbound")
                        : t("schema.telegram.operator.outbound");
                      const statusLabel = item.status === "pending"
                        ? t("schema.telegram.operator.pending")
                        : item.status === "sent"
                          ? t("schema.telegram.operator.sent")
                          : item.status === "failed"
                            ? t("schema.telegram.operator.failed")
                            : item.status;
                      const senderLabel = isOutbound
                        ? t("schema.telegram.operator.senderBot")
                        : item.senderUsername
                          || (selectedSession?.chatType === "private"
                            ? selectedSession.title || selectedSession.username
                            : undefined)
                          || t("schema.telegram.operator.senderUser");
                      const metadataClass = isOutbound ? "text-primary-foreground/75" : "text-muted-foreground";
                      return (
                        <div
                          key={`${item.direction}-${item.requestId ?? item.updateId ?? item.messageId ?? item.occurredAt}-${index}`}
                          className={isOutbound ? "flex justify-end" : "flex justify-start"}
                          data-telegram-message-row
                        >
                          <article
                            data-telegram-message
                            data-direction={item.direction}
                            className={isOutbound
                              ? "max-w-[85%] rounded-2xl rounded-br-md border border-primary/30 bg-primary px-3 py-2 text-xs text-primary-foreground shadow-sm"
                              : "max-w-[85%] rounded-2xl rounded-bl-md border border-border/60 bg-background px-3 py-2 text-xs shadow-sm"}
                          >
                            <div className={`flex min-w-0 items-center justify-between gap-2 ${metadataClass}`}>
                              <span data-telegram-sender className="min-w-0 truncate font-medium">{senderLabel}</span>
                              <span className="min-w-0 truncate">{directionLabel} · {statusLabel}</span>
                              <time className="max-w-[45%] shrink-0 truncate text-right" dateTime={item.occurredAt} title={item.occurredAt}>{item.occurredAt}</time>
                            </div>
                            <p className="mt-1 whitespace-pre-wrap break-words">{item.text}</p>
                            {isOutbound && item.status === "failed" ? (
                              <button
                                type="button"
                                data-telegram-retry={item.requestId}
                                disabled={operatorCapability !== "allowed" || !operatorReady || retryingRequestID !== null}
                                onClick={() => {
                                  if (item.requestId !== undefined) void retryMessage(item.requestId);
                                }}
                                className="mt-2 rounded-md border border-primary-foreground/40 px-2 py-1 text-xs text-primary-foreground hover:bg-primary-foreground/10 disabled:cursor-not-allowed disabled:opacity-50"
                              >
                                {retryingRequestID === item.requestId ? t("feedback.submitting") : t("schema.telegram.operator.retry")}
                              </button>
                            ) : null}
                          </article>
                        </div>
                      );
                    })}
                  </div>
                ) : timelineLoadState === "loading" ? (
                  <p className="shrink-0 text-xs text-muted-foreground">{t("schema.telegram.operator.timelineLoading")}</p>
                ) : timelineLoadState === "error" ? (
                  <p role="alert" className="shrink-0 text-xs text-destructive">{t("schema.telegram.operator.timelineFailed")}</p>
                ) : (
                  <p className="shrink-0 text-xs text-muted-foreground">{t("schema.telegram.operator.timelineEmpty")}</p>
                )}

                <fieldset
                  data-telegram-composer
                  disabled={operatorCapability !== "allowed" || !operatorReady}
                  className="shrink-0 space-y-2 rounded-md border border-border/50 p-3 disabled:opacity-60"
                >
                  <legend className="px-1 text-xs font-semibold">{t("schema.telegram.operator.composer")}</legend>
                  <textarea
                    aria-label={t("schema.telegram.operator.composerText")}
                    value={composerText}
                    onChange={(event) => setComposerText(event.target.value)}
                    onKeyDown={(event) => {
                      if (event.key !== "Enter") return;
                      const isPlainEnter = !event.ctrlKey && !event.altKey && !event.metaKey && !event.shiftKey;
                      const isCtrlEnter = event.ctrlKey && !event.altKey && !event.metaKey && !event.shiftKey;
                      const shouldSend = reverseComposerShortcuts ? isCtrlEnter : isPlainEnter;
                      const shouldInsertLineBreak = reverseComposerShortcuts ? isPlainEnter : isCtrlEnter;
                      if (shouldSend) {
                        event.preventDefault();
                        void sendMessage();
                      } else if (shouldInsertLineBreak) {
                        event.preventDefault();
                        insertTelegramComposerLineBreak(event.currentTarget, setComposerText);
                      }
                    }}
                    placeholder={t("schema.telegram.operator.composerPlaceholder")}
                    rows={3}
                    className="min-h-20 max-h-40 w-full resize-y rounded-md border border-input/80 bg-background px-3 py-2 text-sm outline-none focus-visible:border-ring focus-visible:ring-2 focus-visible:ring-ring/20"
                  />
                  <p id="telegram-composer-shortcut-hint" data-telegram-shortcut-hint className="text-xs text-muted-foreground">
                    {reverseComposerShortcuts
                      ? t("schema.telegram.operator.shortcutsReversed")
                      : t("schema.telegram.operator.shortcutsDefault")}
                  </p>
                  <div className="flex min-w-0 items-center justify-between gap-2">
                    <p role="status" className="min-w-0 flex-1 break-words text-xs text-muted-foreground">
                      {operatorCapability === "unknown"
                        ? t("schema.telegram.operator.capabilityPending")
                        : operatorCapability === "allowed"
                          ? t("schema.telegram.operator.capabilityAllowed")
                          : operatorCapability === "denied"
                            ? t("schema.telegram.operator.capabilityDenied")
                            : t("schema.telegram.operator.capabilityUnavailable")}
                    </p>
                    <label htmlFor="telegram-reverse-composer-shortcuts" className="inline-flex shrink-0 items-center gap-1.5 text-xs text-muted-foreground">
                      <input
                        id="telegram-reverse-composer-shortcuts"
                        data-telegram-shortcut-reverse
                        type="checkbox"
                        checked={reverseComposerShortcuts}
                        onChange={(event) => setReverseComposerShortcuts(event.target.checked)}
                        aria-describedby="telegram-composer-shortcut-hint"
                        className="size-4 cursor-pointer rounded border-input text-primary accent-primary transition-colors focus:ring-2 focus:ring-ring/20"
                      />
                      <span>{t("schema.telegram.operator.reverseShortcuts")}</span>
                    </label>
                    <button
                      type="button"
                      disabled={sending || operatorCapability !== "allowed" || !operatorReady || composerText.trim() === ""}
                      onClick={() => void sendMessage()}
                      className={`${buttonClass} shrink-0`}
                    >
                      {sending ? t("feedback.submitting") : t("schema.telegram.operator.send")}
                    </button>
                  </div>
                </fieldset>
              </div>
            </div>
          )}
        </section>
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
